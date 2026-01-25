package provider

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/search"
	"Qingyu_backend/pkg/logger"
	searchengine "Qingyu_backend/service/search/engine"
)

const (
	// booksCollection MongoDB 书籍集合名称
	booksCollection = "books"
	// defaultPageSize 默认每页数量
	defaultPageSize = 20
	// maxPageSize 最大每页数量
	maxPageSize = 100
)

// BookProvider 书籍搜索提供者
type BookProvider struct {
	engine searchengine.Engine
	config *BookProviderConfig
	logger *logger.Logger
}

// BookProviderConfig 书籍提供者配置
type BookProviderConfig struct {
	AllowedStatuses []string // 允许搜索的书籍状态
	AllowedPrivacy  []bool   // 允许搜索的隐私设置
}

// NewBookProvider 创建书籍搜索提供者
func NewBookProvider(eng searchengine.Engine, config *BookProviderConfig) (*BookProvider, error) {
	if eng == nil {
		return nil, fmt.Errorf("engine cannot be nil")
	}
	if config == nil {
		// 使用默认配置
		config = &BookProviderConfig{
			AllowedStatuses: []string{
				string(bookstore.BookStatusOngoing),
				string(bookstore.BookStatusCompleted),
			},
			AllowedPrivacy: []bool{false}, // 只允许公开书籍
		}
	}

	return &BookProvider{
		engine: eng,
		config: config,
		logger: logger.Get().WithModule("book-provider"),
	}, nil
}

// Search 搜索书籍
func (p *BookProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
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

	// 构建查询参数
	query := req.Query

	// 构建搜索选项
	opts := &searchengine.SearchOptions{
		From: p.calculateOffset(req.Page, req.PageSize),
		Size: p.getPageSize(req.PageSize),
	}

	// 应用强制过滤规则（来自配置）
	filters := p.buildConfigFilters()

	// 应用用户自定义过滤条件
	if req.Filter != nil {
		userFilters := p.buildUserFilters(req.Filter)
		filters = append(filters, userFilters...)
	}

	// 合并所有过滤条件
	if len(filters) > 0 {
		opts.Filter = p.mergeFilters(filters)
	}

	// 构建排序
	sortFields := p.buildSortOptions(req.Sort)
	opts.Sort = sortFields

	// 执行搜索
	engineReq := p.buildEngineRequest(query, opts)
	result, err := p.engine.Search(ctx, booksCollection, engineReq, opts)
	if err != nil {
		p.logger.Error("Book search failed",
			zap.Error(err),
			zap.String("query", query),
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
	p.logger.Info("Book search completed",
		zap.String("provider_type", string(p.Type())),
		zap.String("query", query),
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
func (p *BookProvider) Type() search.SearchType {
	return search.SearchTypeBooks
}

// Validate 验证搜索参数
func (p *BookProvider) Validate(req *search.SearchRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// 验证 Query
	if req.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	// 验证分页参数
	page := p.getPage(req.Page)
	if page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	pageSize := p.getPageSize(req.PageSize)
	if pageSize < 1 {
		return fmt.Errorf("page_size must be greater than 0")
	}

	if pageSize > maxPageSize {
		return fmt.Errorf("page_size cannot exceed %d", maxPageSize)
	}

	return nil
}

// GetByID 根据 ID 获取单个文档
func (p *BookProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
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

	// 应用配置过滤
	filters := p.buildConfigFilters()
	for _, filter := range filters {
		for k, v := range filter {
			query[k] = v
		}
	}

	// 执行搜索
	opts := &searchengine.SearchOptions{
		From: 0,
		Size: 1,
	}

	result, err := p.engine.Search(ctx, booksCollection, query, opts)
	if err != nil {
		p.logger.Error("Get book by ID failed",
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
func (p *BookProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
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

	// 应用配置过滤
	filters := p.buildConfigFilters()
	for _, filter := range filters {
		for k, v := range filter {
			query[k] = v
		}
	}

	// 执行搜索
	opts := &searchengine.SearchOptions{
		From: 0,
		Size: len(ids),
	}

	result, err := p.engine.Search(ctx, booksCollection, query, opts)
	if err != nil {
		p.logger.Error("Batch get books failed",
			zap.Error(err),
			zap.Int("count", len(ids)),
		)
		return nil, fmt.Errorf("batch search failed: %w", err)
	}

	// 转换结果
	items := p.convertHitsToSearchItems(result.Hits)

	p.logger.Info("Batch get books completed",
		zap.Int("requested", len(ids)),
		zap.Int("returned", len(items)),
	)

	return items, nil
}

// buildConfigFilters 构建配置过滤规则（强制过滤）
func (p *BookProvider) buildConfigFilters() []bson.M {
	filters := make([]bson.M, 0)

	// 状态过滤：只搜索允许的状态
	if len(p.config.AllowedStatuses) > 0 {
		filters = append(filters, bson.M{
			"status": bson.M{"$in": p.config.AllowedStatuses},
		})
	}

	// 隐私过滤：只搜索公开书籍
	if len(p.config.AllowedPrivacy) > 0 {
		// 如果配置中只允许 false（公开），则添加过滤
		onlyPublic := true
		for _, privacy := range p.config.AllowedPrivacy {
			if privacy { // 如果允许私密
				onlyPublic = false
				break
			}
		}
		if onlyPublic {
			// Book 模型中没有 is_private 字段，通过状态控制
			// 已发布的状态就是公开的
		}
	}

	return filters
}

// buildUserFilters 构建用户自定义过滤条件
func (p *BookProvider) buildUserFilters(filter map[string]interface{}) []bson.M {
	filters := make([]bson.M, 0)

	// 分类过滤
	if categoryID, ok := filter["category_id"].(string); ok && categoryID != "" {
		objectID, err := primitive.ObjectIDFromHex(categoryID)
		if err == nil {
			filters = append(filters, bson.M{
				"category_ids": bson.M{"$in": []primitive.ObjectID{objectID}},
			})
		}
	}

	// 作者过滤
	if author, ok := filter["author"].(string); ok && author != "" {
		filters = append(filters, bson.M{"author": author})
	}

	// 标签过滤
	if tags, ok := filter["tags"].([]string); ok && len(tags) > 0 {
		filters = append(filters, bson.M{
			"tags": bson.M{"$in": tags},
		})
	}

	// 状态过滤
	if status, ok := filter["status"].(string); ok && status != "" {
		filters = append(filters, bson.M{"status": status})
	}

	// 字数范围过滤
	if minWordCount, ok := filter["word_count_min"].(int); ok && minWordCount > 0 {
		filters = append(filters, bson.M{
			"word_count": bson.M{"$gte": minWordCount},
		})
	}

	if maxWordCount, ok := filter["word_count_max"].(int); ok && maxWordCount > 0 {
		filters = append(filters, bson.M{
			"word_count": bson.M{"$lte": maxWordCount},
		})
	}

	// 评分过滤
	if minRating, ok := filter["rating_min"].(float64); ok && minRating > 0 {
		filters = append(filters, bson.M{
			"rating": bson.M{"$gte": minRating},
		})
	}

	return filters
}

// mergeFilters 合并过滤条件
func (p *BookProvider) mergeFilters(filters []bson.M) map[string]interface{} {
	if len(filters) == 0 {
		return nil
	}

	// 使用 $and 合并所有过滤条件
	result := make(map[string]interface{})
	if len(filters) == 1 {
		result = filters[0]
	} else {
		result["$and"] = filters
	}

	return result
}

// buildSortOptions 构建排序选项
func (p *BookProvider) buildSortOptions(sortFields []search.SortField) []searchengine.SortField {
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
func (p *BookProvider) validateSortField(field string) string {
	validFields := map[string]string{
		"created_at":     "created_at",
		"updated_at":     "updated_at",
		"published_at":   "published_at",
		"word_count":     "word_count",
		"chapter_count":  "chapter_count",
		"view_count":     "view_count",
		"rating":         "rating",
		"rating_count":   "rating_count",
		"title":          "title",
		"last_update_at": "last_update_at",
	}

	if normalized, ok := validFields[field]; ok {
		return normalized
	}
	return ""
}

// convertHitsToSearchItems 转换搜索命中项为 SearchItem
func (p *BookProvider) convertHitsToSearchItems(hits []searchengine.Hit) []search.SearchItem {
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

// buildEngineRequest 构建引擎请求
func (p *BookProvider) buildEngineRequest(query string, opts *searchengine.SearchOptions) interface{} {
	// MongoDB 引擎使用字符串作为查询
	return query
}

// calculateOffset 计算分页偏移量
func (p *BookProvider) calculateOffset(page, pageSize int) int {
	parsedPage := p.getPage(page)
	parsedPageSize := p.getPageSize(pageSize)
	return (parsedPage - 1) * parsedPageSize
}

// getPage 获取页码（默认 1）
func (p *BookProvider) getPage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

// getPageSize 获取每页数量（默认 20，最大 100）
func (p *BookProvider) getPageSize(pageSize int) int {
	if pageSize < 1 {
		return defaultPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}
	return pageSize
}
