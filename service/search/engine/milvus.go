package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"
	searchRepo "Qingyu_backend/repository/search"
)

// MilvusEngine Milvus 向量搜索引擎
type MilvusEngine struct {
	repo   *searchRepo.MilvusRepository
	logger *logger.Logger
}

// NewMilvusEngine 创建 Milvus 引擎
func NewMilvusEngine(repo *searchRepo.MilvusRepository) (*MilvusEngine, error) {
	if repo == nil {
		return nil, fmt.Errorf("milvus repository cannot be nil")
	}

	return &MilvusEngine{
		repo:   repo,
		logger: logger.Get().WithModule("milvus-engine"),
	}, nil
}

// Search 执行向量搜索
// query 参数应为 []float32 类型的向量
func (m *MilvusEngine) Search(ctx context.Context, index string, query interface{}, opts *SearchOptions) (*SearchResult, error) {
	startTime := time.Now()

	// 解析查询向量为 []float32
	var vectors []float32
	switch v := query.(type) {
	case []float32:
		vectors = v
	case []float64:
		vectors = make([]float32, len(v))
		for i, f := range v {
			vectors[i] = float32(f)
		}
	default:
		return nil, fmt.Errorf("milvus search query must be []float32 or []float64, got %T", query)
	}

	// 确定 topK
	topK := 10
	if opts != nil && opts.Size > 0 {
		topK = opts.Size
	}

	// 调用 repository 搜索
	collectionName := index
	results, err := m.repo.Search(ctx, collectionName, vectors, topK)
	if err != nil {
		m.logger.Error("Milvus search failed",
			zap.String("collection", collectionName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("milvus search failed: %w", err)
	}

	// 转换为 Engine 统一的 SearchResult
	hits := make([]Hit, 0, len(results))
	for _, result := range results {
		bookID, _ := result["book_id"].(string)
		score, _ := result["score"].(float32)

		source := map[string]interface{}{
			"book_id": bookID,
		}

		hits = append(hits, Hit{
			ID:     bookID,
			Score:  float64(score),
			Source: source,
		})
	}

	took := time.Since(startTime)

	m.logger.Info("Milvus engine search completed",
		zap.String("collection", collectionName),
		zap.Int("topK", topK),
		zap.Int("results", len(hits)),
		zap.Duration("took", took),
	)

	return &SearchResult{
		Total: int64(len(hits)),
		Hits:  hits,
		Took:  took,
	}, nil
}

// Index 批量索引向量
func (m *MilvusEngine) Index(ctx context.Context, index string, documents []Document) error {
	startTime := time.Now()

	if len(documents) == 0 {
		return nil
	}

	// 转换为 repository 需要的格式
	vectors := make([]map[string]interface{}, 0, len(documents))
	for _, doc := range documents {
		vecMap := map[string]interface{}{
			"book_id": doc.ID,
		}

		// 从 Source 提取向量和元数据
		if vec, ok := doc.Source["title_vector"].([]float32); ok {
			vecMap["title_vector"] = vec
		} else if vec64, ok := doc.Source["title_vector"].([]float64); ok {
			vec := make([]float32, len(vec64))
			for i, f := range vec64 {
				vec[i] = float32(f)
			}
			vecMap["title_vector"] = vec
		} else {
			m.logger.Warn("Document missing title_vector, skipping",
				zap.String("book_id", doc.ID),
			)
			continue
		}

		if status, ok := doc.Source["status"].(string); ok {
			vecMap["status"] = status
		} else {
			vecMap["status"] = "ongoing"
		}

		if isPrivate, ok := doc.Source["is_private"].(bool); ok {
			vecMap["is_private"] = isPrivate
		} else {
			vecMap["is_private"] = false
		}

		vectors = append(vectors, vecMap)
	}

	if len(vectors) == 0 {
		return fmt.Errorf("no valid vectors to index")
	}

	if err := m.repo.Insert(ctx, vectors); err != nil {
		m.logger.Error("Milvus bulk insert failed",
			zap.String("collection", index),
			zap.Int("count", len(vectors)),
			zap.Error(err),
		)
		return fmt.Errorf("milvus index failed: %w", err)
	}

	took := time.Since(startTime)
	m.logger.Info("Milvus bulk insert completed",
		zap.String("collection", index),
		zap.Int("count", len(vectors)),
		zap.Duration("took", took),
	)

	return nil
}

// Update 更新向量（Milvus 使用 delete + insert 策略）
func (m *MilvusEngine) Update(ctx context.Context, index string, id string, document Document) error {
	startTime := time.Now()

	// Milvus 不支持原地更新，使用 delete + insert 策略
	if err := m.Delete(ctx, index, id); err != nil {
		// 如果是文档不存在的错误，可以忽略继续插入
		m.logger.Warn("Failed to delete existing vector during update, continuing with insert",
			zap.String("book_id", id),
			zap.Error(err),
		)
	}

	// 插入新向量
	docs := []Document{document}
	if err := m.Index(ctx, index, docs); err != nil {
		m.logger.Error("Milvus update (re-insert) failed",
			zap.String("book_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("milvus update failed: %w", err)
	}

	took := time.Since(startTime)
	m.logger.Info("Milvus update completed",
		zap.String("book_id", id),
		zap.Duration("took", took),
	)

	return nil
}

// Delete 删除向量
func (m *MilvusEngine) Delete(ctx context.Context, index string, id string) error {
	startTime := time.Now()

	if err := m.repo.DeleteVectors(ctx, index, []string{id}); err != nil {
		m.logger.Error("Milvus delete failed",
			zap.String("book_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("milvus delete failed: %w", err)
	}

	took := time.Since(startTime)
	m.logger.Info("Milvus delete completed",
		zap.String("book_id", id),
		zap.Duration("took", took),
	)

	return nil
}

// CreateIndex 创建集合和 IVF_FLAT 索引
func (m *MilvusEngine) CreateIndex(ctx context.Context, index string, mapping interface{}) error {
	startTime := time.Now()

	// 1. 创建集合
	if err := m.repo.CreateCollection(ctx); err != nil {
		m.logger.Error("Failed to create Milvus collection", zap.Error(err))
		return fmt.Errorf("failed to create milvus collection: %w", err)
	}

	// 2. 创建 IVF_FLAT 索引
	idx, err := entity.NewIndexIvfFlat(entity.COSINE, 128)
	if err != nil {
		m.logger.Error("Failed to create IVF_FLAT index params", zap.Error(err))
		return fmt.Errorf("failed to create index params: %w", err)
	}

	if err := m.repo.CreateIndex(ctx, "title_vector", idx); err != nil {
		m.logger.Error("Failed to create Milvus vector index",
			zap.String("field", "title_vector"),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create vector index: %w", err)
	}

	took := time.Since(startTime)
	m.logger.Info("Milvus collection and index created", zap.Duration("took", took))

	return nil
}

// Health 健康检查
func (m *MilvusEngine) Health(ctx context.Context) error {
	return m.repo.HealthCheck(ctx)
}
