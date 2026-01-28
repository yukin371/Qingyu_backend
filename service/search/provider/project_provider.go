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
	// projectsCollection MongoDB 项目集合名称
	projectsCollection = "projects"
	// projectDefaultPageSize 默认每页数量
	projectDefaultPageSize = 20
	// projectMaxPageSize 最大每页数量
	projectMaxPageSize = 100
)

// ProjectProvider 创作项目搜索提供者
type ProjectProvider struct {
	engine searchengine.Engine
	config *ProjectProviderConfig
	logger *logger.Logger
}

// ProjectProviderConfig 项目提供者配置
type ProjectProviderConfig struct {
	// ProjectProvider 需要认证，权限过滤在运行时强制添加
	// 不需要配置默认过滤规则，所有过滤都基于用户认证
}

// NewProjectProvider 创建项目搜索提供者
func NewProjectProvider(eng searchengine.Engine, config *ProjectProviderConfig) (*ProjectProvider, error) {
	if eng == nil {
		return nil, fmt.Errorf("engine cannot be nil")
	}
	if config == nil {
		// 使用默认配置
		config = &ProjectProviderConfig{}
	}

	return &ProjectProvider{
		engine: eng,
		config: config,
		logger: logger.Get().WithModule("project-provider"),
	}, nil
}

// Search 搜索项目
func (p *ProjectProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
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
		p.logger.Warn("Project search missing user_id",
			zap.String("query", req.Query),
		)
		return &search.SearchResponse{
			Success: false,
			Error: &search.ErrorInfo{
				Code:    search.ErrCodeUnauthorized,
				Message: "user_id is required for project search",
				Details: "ProjectProvider requires authentication",
			},
		}, nil
	}

	// 构建查询参数
	// 使用多字段搜索：标题和描述
	query := map[string]interface{}{
		"must": []map[string]interface{}{
			{
				"multi_match": map[string]interface{}{
					"query":  req.Query,
					"fields": []string{"title^3", "summary^1"},
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

	// 1. 强制添加权限过滤：只能搜索自己的项目
	filters = append(filters, map[string]interface{}{
		"term": map[string]interface{}{
			"author_id": userID,
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
	result, err := p.engine.Search(ctx, projectsCollection, query, opts)
	if err != nil {
		p.logger.Error("Project search failed",
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
	p.logger.Info("Project search completed",
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
func (p *ProjectProvider) Type() search.SearchType {
	return search.SearchTypeProjects
}

// Validate 验证搜索参数
func (p *ProjectProvider) Validate(req *search.SearchRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// 验证 Query
	if req.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	// 验证分页参数（直接验证原始值，不使用 getPage/getPageSize）
	if req.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	if req.PageSize < 1 {
		return fmt.Errorf("page_size must be greater than 0")
	}

	if req.PageSize > projectMaxPageSize {
		return fmt.Errorf("page_size cannot exceed %d", projectMaxPageSize)
	}

	return nil
}

// GetByID 根据 ID 获取单个项目
func (p *ProjectProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
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

	result, err := p.engine.Search(ctx, projectsCollection, query, opts)
	if err != nil {
		p.logger.Error("Get project by ID failed",
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

// GetBatch 批量获取项目
func (p *ProjectProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
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

	result, err := p.engine.Search(ctx, projectsCollection, query, opts)
	if err != nil {
		p.logger.Error("Batch get projects failed",
			zap.Error(err),
			zap.Int("count", len(ids)),
		)
		return nil, fmt.Errorf("batch search failed: %w", err)
	}

	// 转换结果
	items := p.convertHitsToSearchItems(result.Hits)

	p.logger.Info("Batch get projects completed",
		zap.Int("requested", len(ids)),
		zap.Int("returned", len(items)),
	)

	return items, nil
}

// extractUserID 从请求中提取用户 ID
func (p *ProjectProvider) extractUserID(req *search.SearchRequest) (string, bool) {
	if req.Filter == nil {
		return "", false
	}

	userID, ok := req.Filter["user_id"].(string)
	return userID, ok
}

// buildUserFilters 构建用户自定义过滤条件
func (p *ProjectProvider) buildUserFilters(filter map[string]interface{}) []map[string]interface{} {
	filters := make([]map[string]interface{}, 0)

	// 状态过滤
	if status, ok := filter["status"].(string); ok && status != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"status": status,
			},
		})
	}

	// 可见性过滤（is_public 映射到 visibility 字段）
	if isPublic, ok := filter["is_public"].(bool); ok {
		visibility := "private"
		if isPublic {
			visibility = "public"
		}
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"visibility": visibility,
			},
		})
	}

	// 分类过滤
	if category, ok := filter["category"].(string); ok && category != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"category": category,
			},
		})
	}

	// 写作类型过滤
	if writingType, ok := filter["writing_type"].(string); ok && writingType != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"writing_type": writingType,
			},
		})
	}

	return filters
}

// buildSortOptions 构建排序选项
func (p *ProjectProvider) buildSortOptions(sortFields []search.SortField) []searchengine.SortField {
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
func (p *ProjectProvider) validateSortField(field string) string {
	validFields := map[string]string{
		"created_at":     "created_at",
		"updated_at":     "updated_at",
		"published_at":   "published_at",
		"title":          "title",
		"total_words":    "statistics.total_words",
		"chapter_count":  "statistics.chapter_count",
		"document_count": "statistics.document_count",
		"last_update_at": "statistics.last_update_at",
	}

	if normalized, ok := validFields[field]; ok {
		return normalized
	}
	return ""
}

// convertHitsToSearchItems 转换搜索命中项为 SearchItem
func (p *ProjectProvider) convertHitsToSearchItems(hits []searchengine.Hit) []search.SearchItem {
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
func (p *ProjectProvider) calculateOffset(page, pageSize int) int {
	parsedPage := p.getPage(page)
	parsedPageSize := p.getPageSize(pageSize)
	return (parsedPage - 1) * parsedPageSize
}

// getPage 获取页码（默认 1）
func (p *ProjectProvider) getPage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

// getPageSize 获取每页数量（默认 20，最大 100）
func (p *ProjectProvider) getPageSize(pageSize int) int {
	if pageSize < 1 {
		return projectDefaultPageSize
	}
	if pageSize > projectMaxPageSize {
		return projectMaxPageSize
	}
	return pageSize
}
