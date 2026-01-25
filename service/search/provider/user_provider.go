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
	// usersCollection MongoDB 用户集合名称
	usersCollection = "users"
	// userDefaultPageSize 默认每页数量
	userDefaultPageSize = 20
	// userMaxPageSize 最大每页数量
	userMaxPageSize = 100
)

// UserProvider 用户搜索提供者
type UserProvider struct {
	engine searchengine.Engine
	config *UserProviderConfig
	logger *logger.Logger
}

// UserProviderConfig 用户提供者配置
type UserProviderConfig struct {
	// UserProvider 不需要强制过滤规则
	// 所有过滤条件由用户搜索请求中的 Filter 参数控制
}

// NewUserProvider 创建用户搜索提供者
func NewUserProvider(eng searchengine.Engine, config *UserProviderConfig) (*UserProvider, error) {
	if eng == nil {
		return nil, fmt.Errorf("engine cannot be nil")
	}
	if config == nil {
		// 使用默认配置（无强制过滤）
		config = &UserProviderConfig{}
	}

	return &UserProvider{
		engine: eng,
		config: config,
		logger: logger.Get().WithModule("user-provider"),
	}, nil
}

// Search 搜索用户
func (p *UserProvider) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
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

	// 构建查询参数 - 使用带权重的多字段搜索
	// Elasticsearch 使用 ^ 符号指定字段权重
	query := map[string]interface{}{
		"must": []map[string]interface{}{
			{
				"multi_match": map[string]interface{}{
					"query":  req.Query,
					"fields": []string{"username^10", "nickname^8", "bio^1"},
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

	// 应用用户自定义过滤条件
	// UserProvider 不需要认证，不需要添加强制权限过滤
	if req.Filter != nil {
		userFilters := p.buildUserFilters(req.Filter)
		if len(userFilters) > 0 {
			query["filter"] = userFilters
		}
	}

	// 构建排序
	sortFields := p.buildSortOptions(req.Sort)
	opts.Sort = sortFields

	// 执行搜索
	result, err := p.engine.Search(ctx, usersCollection, query, opts)
	if err != nil {
		p.logger.Error("User search failed",
			zap.Error(err),
			zap.String("query", req.Query),
			zap.Any("filters", req.Filter),
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
	p.logger.Info("User search completed",
		zap.String("provider_type", string(p.Type())),
		zap.String("query", req.Query),
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
func (p *UserProvider) Type() search.SearchType {
	return search.SearchTypeUsers
}

// Validate 验证搜索参数
func (p *UserProvider) Validate(req *search.SearchRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// 验证 Query
	if req.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}

	// 验证分页参数（不使用 getPage/getPageSize，直接验证原始值）
	if req.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	if req.PageSize < 1 {
		return fmt.Errorf("page_size must be greater than 0")
	}

	if req.PageSize > userMaxPageSize {
		return fmt.Errorf("page_size cannot exceed %d", userMaxPageSize)
	}

	return nil
}

// GetByID 根据 ID 获取单个用户
func (p *UserProvider) GetByID(ctx context.Context, id string) (*search.SearchItem, error) {
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

	result, err := p.engine.Search(ctx, usersCollection, query, opts)
	if err != nil {
		p.logger.Error("Get user by ID failed",
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

// GetBatch 批量获取用户
func (p *UserProvider) GetBatch(ctx context.Context, ids []string) ([]search.SearchItem, error) {
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

	result, err := p.engine.Search(ctx, usersCollection, query, opts)
	if err != nil {
		p.logger.Error("Batch get users failed",
			zap.Error(err),
			zap.Int("count", len(ids)),
		)
		return nil, fmt.Errorf("batch search failed: %w", err)
	}

	// 转换结果
	items := p.convertHitsToSearchItems(result.Hits)

	p.logger.Info("Batch get users completed",
		zap.Int("requested", len(ids)),
		zap.Int("returned", len(items)),
	)

	return items, nil
}

// buildUserFilters 构建用户自定义过滤条件
func (p *UserProvider) buildUserFilters(filter map[string]interface{}) []map[string]interface{} {
	filters := make([]map[string]interface{}, 0)

	// 角色过滤
	if role, ok := filter["role"].(string); ok && role != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"roles": role,
			},
		})
	}

	// 认证状态过滤
	if isVerified, ok := filter["is_verified"].(bool); ok {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"email_verified": isVerified,
			},
		})
	}

	// 用户状态过滤
	if status, ok := filter["status"].(string); ok && status != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"status": status,
			},
		})
	}

	return filters
}

// buildSortOptions 构建排序选项
func (p *UserProvider) buildSortOptions(sortFields []search.SortField) []searchengine.SortField {
	if len(sortFields) == 0 {
		// 默认排序：按评分降序（相关性优先）
		return []searchengine.SortField{
			{Field: "_score", Ascending: false},
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
			{Field: "_score", Ascending: false},
		}
	}

	return engineSortFields
}

// validateSortField 验证排序字段
func (p *UserProvider) validateSortField(field string) string {
	validFields := map[string]string{
		"username":    "username",
		"nickname":    "nickname",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
		"vip_level":   "vip_level",
		"last_login":  "last_login_at",
		"_score":      "_score",
	}

	if normalized, ok := validFields[field]; ok {
		return normalized
	}
	return ""
}

// convertHitsToSearchItems 转换搜索命中项为 SearchItem
func (p *UserProvider) convertHitsToSearchItems(hits []searchengine.Hit) []search.SearchItem {
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
func (p *UserProvider) calculateOffset(page, pageSize int) int {
	parsedPage := p.getPage(page)
	parsedPageSize := p.getPageSize(pageSize)
	return (parsedPage - 1) * parsedPageSize
}

// getPage 获取页码（默认 1）
func (p *UserProvider) getPage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

// getPageSize 获取每页数量（默认 20，最大 100）
func (p *UserProvider) getPageSize(pageSize int) int {
	if pageSize < 1 {
		return userDefaultPageSize
	}
	if pageSize > userMaxPageSize {
		return userMaxPageSize
	}
	return pageSize
}
