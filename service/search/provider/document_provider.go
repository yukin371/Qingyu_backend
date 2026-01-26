package provider

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"Qingyu_backend/models/search"
	"Qingyu_backend/pkg/logger"
	searchengine "Qingyu_backend/service/search/engine"
)

const (
	// documentsCollection MongoDB 文档集合名称
	documentsCollection = "documents"
	// documentDefaultPageSize 默认每页数量
	documentDefaultPageSize = 20
	// documentMaxPageSize 最大每页数量
	documentMaxPageSize = 100
)

// DocumentProvider 文档搜索提供者
type DocumentProvider struct {
	engine searchengine.Engine
	config *DocumentProviderConfig
	logger *logger.Logger
}

// DocumentProviderConfig 文档提供者配置
type DocumentProviderConfig struct {
	// DocumentProvider 需要认证，权限过滤在运行时强制添加
	// 不需要配置默认过滤规则，所有过滤都基于用户认证
}

// NewDocumentProvider 创建文档搜索提供者
func NewDocumentProvider(eng searchengine.Engine, config *DocumentProviderConfig) (*DocumentProvider, error) {
	if eng == nil {
		return nil, fmt.Errorf("engine cannot be nil")
	}
	if config == nil {
		// 使用默认配置
		config = &DocumentProviderConfig{}
	}

	return &DocumentProvider{
		engine: eng,
		config: config,
		logger: logger.Get().WithModule("document-provider"),
	}, nil
}

// Search 搜索文档
func (p *DocumentProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	startTime := time.Now()

	// 验证请求
	if err := p.Validate(req); err != nil {
		p.logger.Error("Invalid search request",
			zap.Error(err),
			zap.String("query", req.Query),
		)
		return &search.SearchResponse{
			Success: false,
			Error: &search.ErrorInfo{
				Code:    search.ErrCodeInvalidRequest,
				Message: err.Error(),
			},
		}, nil
	}

	// 验证用户认证：从 Filter 中获取 user_id
	userID, ok := p.extractUserID(req)
	if !ok || userID == "" {
		p.logger.Warn("Document search missing user_id",
			zap.String("query", req.Query),
		)
		return &search.SearchResponse{
			Success: false,
			Error: &search.ErrorInfo{
				Code:    search.ErrCodeUnauthorized,
				Message: "user_id is required for document search",
				Details: "DocumentProvider requires authentication",
			},
		}, nil
	}

	// 构建查询参数
	// 使用多字段搜索：标题和内容
	query := map[string]interface{}{
		"must": []map[string]interface{}{
			{
				"multi_match": map[string]interface{}{
					"query":  req.Query,
					"fields": []string{"title^3", "content^1"},
					"type":   "best_fields",
				},
			},
		},
	}

	// 构建搜索选项
	opts := &searchengine.SearchOptions{
		From: p.calculateOffset(req.Page, req.PageSize),
		Size: p.getPageSize(req.PageSize),
	}

	// 构建过滤条件
	filters := make([]map[string]interface{}, 0)

	// 1. 强制添加权限过滤：只能搜索自己项目的文档
	filters = append(filters, map[string]interface{}{
		"term": map[string]interface{}{
			"user_id": userID,
		},
	})

	// 2. 应用用户自定义过滤条件
	if req.Filter != nil {
		userFilters := p.buildUserFilters(req.Filter)
		filters = append(filters, userFilters...)
	}

	// 合并所有过滤条件
	if len(filters) > 0 {
		if len(filters) == 1 {
			query["filter"] = filters[0]
		} else {
			query["filter"] = map[string]interface{}{
				"bool": map[string]interface{}{
					"must": filters,
				},
			}
		}
	}

	// 构建排序
	sortFields := p.buildSortOptions(req.Sort)
	opts.Sort = sortFields

	// 执行搜索
	result, err := p.engine.Search(ctx, documentsCollection, query, opts)
	if err != nil {
		p.logger.Error("Document search failed",
			zap.Error(err),
			zap.String("query", req.Query),
			zap.String("user_id", userID),
			zap.Any("filters", opts.Filter),
		)
		return &search.SearchResponse{
			Success: false,
			Error: &search.ErrorInfo{
				Code:    search.ErrCodeEngineFailure,
				Message: "Search engine failed",
				Details: err.Error(),
			},
		}, nil
	}

	// 转换搜索结果
	searchItems := p.convertHitsToSearchItems(result.Hits)

	took := time.Since(startTime)

	// 记录搜索日志
	p.logger.Info("Document search completed",
		zap.String("provider_type", string(p.Type())),
		zap.String("query", req.Query),
		zap.String("user_id", userID),
		zap.Int64("total", result.Total),
		zap.Int("returned", len(searchItems)),
		zap.Duration("took_ms", took),
	)

	// 构建响应
	return &search.SearchResponse{
		Success: true,
		Data: &search.SearchData{
			Type:     p.Type(),
			Total:    result.Total,
			Page:     p.getPage(req.Page),
			PageSize: opts.Size,
			Results:  searchItems,
			Took:     took,
		},
		Meta: &search.MetaInfo{
			TookMs: took.Milliseconds(),
		},
	}, nil
}

// Type 获取搜索类型
func (p *DocumentProvider) Type() search.SearchType {
	return search.SearchTypeDocuments
}

// Validate 验证搜索参数
func (p *DocumentProvider) Validate(req *search.SearchRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// 验证 Query
	if req.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	// 验证分页参数 - 检查原始值而不是默认值
	if req.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	if req.PageSize < 1 {
		return fmt.Errorf("page_size must be greater than 0")
	}

	if req.PageSize > documentMaxPageSize {
		return fmt.Errorf("page_size cannot exceed %d", documentMaxPageSize)
	}

	return nil
}

// GetByID 根据 ID 获取单个文档
func (p *DocumentProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
	// 验证 ID
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	// 转换 ID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// 构建查询
	query := bson.M{"_id": objectID}

	// 执行搜索
	opts := &searchengine.SearchOptions{
		From: 0,
		Size: 1,
	}

	result, err := p.engine.Search(ctx, documentsCollection, query, opts)
	if err != nil {
		p.logger.Error("Get document by ID failed",
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(result.Hits) == 0 {
		return nil, search.ErrDocumentNotFound
	}

	// 转换结果
	items := p.convertHitsToSearchItems(result.Hits)
	if len(items) == 0 {
		return nil, search.ErrDocumentNotFound
	}

	return &items[0], nil
}

// GetBatch 批量获取文档
func (p *DocumentProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
	if len(ids) == 0 {
		return []search.SearchItem{}, nil
	}

	// 转换 IDs
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			p.logger.Warn("Invalid ID format, skipping",
				zap.String("id", id),
				zap.Error(err),
			)
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return []search.SearchItem{}, nil
	}

	// 构建查询
	query := bson.M{"_id": bson.M{"$in": objectIDs}}

	// 执行搜索
	opts := &searchengine.SearchOptions{
		From: 0,
		Size: len(ids),
	}

	result, err := p.engine.Search(ctx, documentsCollection, query, opts)
	if err != nil {
		p.logger.Error("Batch get documents failed",
			zap.Error(err),
			zap.Int("count", len(ids)),
		)
		return nil, fmt.Errorf("batch search failed: %w", err)
	}

	// 转换结果
	items := p.convertHitsToSearchItems(result.Hits)

	p.logger.Info("Batch get documents completed",
		zap.Int("requested", len(ids)),
		zap.Int("returned", len(items)),
	)

	return items, nil
}

// extractUserID 从请求中提取用户 ID
func (p *DocumentProvider) extractUserID(req *search.SearchRequest) (string, bool) {
	if req.Filter == nil {
		return "", false
	}

	userID, ok := req.Filter["user_id"].(string)
	return userID, ok
}

// buildUserFilters 构建用户自定义过滤条件
func (p *DocumentProvider) buildUserFilters(filter map[string]interface{}) []map[string]interface{} {
	filters := make([]map[string]interface{}, 0)

	// 项目 ID 过滤
	if projectID, ok := filter["project_id"].(string); ok && projectID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"project_id": projectID,
			},
		})
	}

	// 状态过滤
	if status, ok := filter["status"].(string); ok && status != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"status": status,
			},
		})
	}

	// 类型过滤
	if docType, ok := filter["type"].(string); ok && docType != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"type": docType,
			},
		})
	}

	// 层级过滤
	if level, ok := filter["level"].(int); ok {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"level": level,
			},
		})
	}

	// 字数范围过滤
	if minWordCount, ok := filter["word_count_min"].(int); ok && minWordCount > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"word_count": map[string]interface{}{
					"gte": minWordCount,
				},
			},
		})
	}

	if maxWordCount, ok := filter["word_count_max"].(int); ok && maxWordCount > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"word_count": map[string]interface{}{
					"lte": maxWordCount,
				},
			},
		})
	}

	return filters
}

// buildSortOptions 构建排序选项
func (p *DocumentProvider) buildSortOptions(sortFields []search.SortField) []searchengine.SortField {
	if len(sortFields) == 0 {
		// 默认排序：按更新时间降序
		return []searchengine.SortField{
			{Field: "updated_at", Ascending: false},
		}
	}

	engineSortFields := make([]searchengine.SortField, 0, len(sortFields))
	for _, sf := range sortFields {
		// 验证字段名
		validField := p.validateSortField(sf.Field)
		if validField != "" {
			engineSortFields = append(engineSortFields, searchengine.SortField{
				Field:     validField,
				Ascending: sf.Ascending,
			})
		}
	}

	// 如果没有有效排序字段，使用默认排序
	if len(engineSortFields) == 0 {
		return []searchengine.SortField{
			{Field: "updated_at", Ascending: false},
		}
	}

	return engineSortFields
}

// validateSortField 验证排序字段
func (p *DocumentProvider) validateSortField(field string) string {
	validFields := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"title":      "title",
		"word_count": "word_count",
		"level":      "level",
		"order":      "order",
	}

	if normalized, ok := validFields[field]; ok {
		return normalized
	}
	return ""
}

// convertHitsToSearchItems 转换搜索命中项为 SearchItem
func (p *DocumentProvider) convertHitsToSearchItems(hits []searchengine.Hit) []search.SearchItem {
	items := make([]search.SearchItem, 0, len(hits))

	for _, hit := range hits {
		item := search.SearchItem{
			ID:        hit.ID,
			Score:     hit.Score,
			Data:      hit.Source,
			Highlight: hit.Highlight,
		}
		items = append(items, item)
	}

	return items
}

// calculateOffset 计算分页偏移量
func (p *DocumentProvider) calculateOffset(page, pageSize int) int {
	parsedPage := p.getPage(page)
	parsedPageSize := p.getPageSize(pageSize)
	return (parsedPage - 1) * parsedPageSize
}

// getPage 获取页码（默认 1）
func (p *DocumentProvider) getPage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

// getPageSize 获取每页数量（默认 20，最大 100）
func (p *DocumentProvider) getPageSize(pageSize int) int {
	if pageSize < 1 {
		return documentDefaultPageSize
	}
	if pageSize > documentMaxPageSize {
		return documentMaxPageSize
	}
	return pageSize
}
