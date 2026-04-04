package provider

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"Qingyu_backend/models/search"
	"Qingyu_backend/pkg/logger"
	searchengine "Qingyu_backend/service/search/engine"
)

const (
	// defaultVectorTopK 默认向量搜索返回数量
	defaultVectorTopK = 10
	// maxVectorTopK 最大向量搜索返回数量
	maxVectorTopTopK = 100
	// vectorDimension 默认向量维度
	vectorDimension = 1536
)

// VectorProvider 向量搜索提供者
type VectorProvider struct {
	engine  searchengine.Engine
	typeVal search.SearchType
	logger  *logger.Logger
}

// NewVectorProvider 创建向量搜索提供者
func NewVectorProvider(eng searchengine.Engine) (*VectorProvider, error) {
	if eng == nil {
		return nil, fmt.Errorf("engine cannot be nil")
	}

	return &VectorProvider{
		engine:  eng,
		typeVal: search.SearchTypeVector,
		logger:  logger.Get().WithModule("vector-provider"),
	}, nil
}

// Search 执行向量搜索
// SearchRequest.Query 在此 provider 中不直接使用作为关键词
// 而是通过 SearchRequest.Options["vector"] 传入 []float32 向量
func (p *VectorProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	startTime := time.Now()

	// 验证请求
	if err := p.Validate(req); err != nil {
		p.logger.Error("Invalid vector search request",
			zap.Error(err),
		)
		return &search.SearchResponse{
			Success: false,
			Error: &search.ErrorInfo{
				Code:    search.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		}, nil
	}

	// 从 Options 中提取向量数据
	vector, err := p.extractVector(req)
	if err != nil {
		return &search.SearchResponse{
			Success: false,
			Error: &search.ErrorInfo{
				Code:    search.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		}, nil
	}

	// 构建 Milvus 集合名
	collection := "qingyu_books"
	if colName, ok := req.Options["collection"].(string); ok && colName != "" {
		collection = colName
	}

	// 构建搜索选项
	topK := p.getTopK(req.PageSize)
	opts := &searchengine.SearchOptions{
		From: 0,
		Size: topK,
	}

	// 应用过滤条件
	if req.Filter != nil {
		opts.Filter = req.Filter
	}

	// 执行 Milvus 向量搜索
	result, err := p.engine.Search(ctx, collection, vector, opts)
	if err != nil {
		p.logger.Error("Vector search failed",
			zap.Error(err),
			zap.String("collection", collection),
		)
		return &search.SearchResponse{
			Success: false,
			Error: &search.ErrorInfo{
				Code:    search.ErrCodeEngineFailure,
				Message: "Vector search engine failed",
				Details: err.Error(),
			},
		}, nil
	}

	// 转换搜索结果
	searchItems := p.convertHitsToSearchItems(result.Hits)

	took := time.Since(startTime)

	p.logger.Info("Vector search completed",
		zap.String("collection", collection),
		zap.Int("topK", topK),
		zap.Int64("total", result.Total),
		zap.Int("returned", len(searchItems)),
		zap.Duration("took", took),
	)

	return &search.SearchResponse{
		Success: true,
		Data: &search.SearchData{
			Type:     p.Type(),
			Total:    result.Total,
			Page:     1,
			PageSize: topK,
			Results:  searchItems,
			Took:     took,
		},
		Meta: &search.MetaInfo{
			TookMs: took.Milliseconds(),
		},
	}, nil
}

// Type 获取搜索类型
func (p *VectorProvider) Type() search.SearchType {
	return p.typeVal
}

// Validate 验证搜索参数
func (p *VectorProvider) Validate(req *search.SearchRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// 验证向量数据存在
	if req.Options == nil {
		return fmt.Errorf("options with vector data is required for vector search")
	}

	_, hasVector := req.Options["vector"]
	if !hasVector {
		// Query 也可以直接作为占位符使用，但 Options["vector"] 是必须的
		return fmt.Errorf("options[\"vector\"] is required for vector search, provide []float32")
	}

	return nil
}

// GetByID 根据 ID 获取单个文档
// 向量搜索不直接支持按 ID 获取，这里返回错误提示
func (p *VectorProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
	return nil, fmt.Errorf("vector provider does not support GetByID, use book provider instead")
}

// GetBatch 批量获取文档
// 向量搜索不直接支持按 ID 批量获取
func (p *VectorProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
	return nil, fmt.Errorf("vector provider does not support GetBatch, use book provider instead")
}

// extractVector 从请求中提取向量数据
func (p *VectorProvider) extractVector(req *search.SearchRequest) ([]float32, error) {
	vecData := req.Options["vector"]

	switch v := vecData.(type) {
	case []float32:
		if len(v) != vectorDimension {
			return nil, fmt.Errorf("vector dimension must be %d, got %d", vectorDimension, len(v))
		}
		return v, nil
	case []float64:
		if len(v) != vectorDimension {
			return nil, fmt.Errorf("vector dimension must be %d, got %d", vectorDimension, len(v))
		}
		result := make([]float32, len(v))
		for i, f := range v {
			result[i] = float32(f)
		}
		return result, nil
	case []interface{}:
		result := make([]float32, len(v))
		for i, val := range v {
			switch f := val.(type) {
			case float64:
				result[i] = float32(f)
			case float32:
				result[i] = f
			case int:
				result[i] = float32(f)
			default:
				return nil, fmt.Errorf("invalid vector element type at index %d: %T", i, val)
			}
		}
		if len(result) != vectorDimension {
			return nil, fmt.Errorf("vector dimension must be %d, got %d", vectorDimension, len(result))
		}
		return result, nil
	default:
		return nil, fmt.Errorf("vector must be []float32 or []float64, got %T", vecData)
	}
}

// getTopK 获取 topK 值
func (p *VectorProvider) getTopK(pageSize int) int {
	if pageSize <= 0 {
		return defaultVectorTopK
	}
	if pageSize > maxVectorTopTopK {
		return maxVectorTopTopK
	}
	return pageSize
}

// convertHitsToSearchItems 转换搜索命中项为 SearchItem
func (p *VectorProvider) convertHitsToSearchItems(hits []searchengine.Hit) []search.SearchItem {
	items := make([]search.SearchItem, 0, len(hits))

	for _, hit := range hits {
		item := search.SearchItem{
			ID:    hit.ID,
			Score: hit.Score,
			Data:  hit.Source,
		}
		items = append(items, item)
	}

	return items
}
