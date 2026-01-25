package search

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"Qingyu_backend/models/search"
	"Qingyu_backend/pkg/metrics"
	"Qingyu_backend/service/search/cache"
	searchengine "Qingyu_backend/service/search/engine"
	"Qingyu_backend/service/search/provider"
)

// SearchService 统一搜索服务
type SearchService struct {
	providers map[search.SearchType]provider.Provider
	cache     cache.Cache
	logger    *log.Logger
	config    *Config

	// ES 相关
	esConfig      *SearchConfig // ES 配置
	esEngine      searchengine.Engine // ES 引擎（如果启用）
	mongoEngine   searchengine.Engine // MongoDB 引擎（fallback）

	// 灰度决策
	grayscaleDecision GrayScaleDecision // 灰度决策器
	mu                sync.RWMutex     // 读写锁
}

// Config 搜索服务配置
type Config struct {
	// 是否启用缓存
	EnableCache bool
	// 默认缓存过期时间
	DefaultCacheTTL int // 秒
	// 最大并发搜索数
	MaxConcurrentSearches int
}

// NewSearchService 创建搜索服务实例
func NewSearchService(logger *log.Logger, config *Config, grayscaleDecision GrayScaleDecision) *SearchService {
	if config == nil {
		config = &Config{
			EnableCache:           true,
			DefaultCacheTTL:       300, // 5分钟
			MaxConcurrentSearches: 10,
		}
	}

	return &SearchService{
		providers:          make(map[search.SearchType]provider.Provider),
		logger:             logger,
		config:             config,
		esConfig:           nil,                  // 将通过 SetESConfig 设置
		esEngine:           nil,                  // 将通过 SetESEngine 设置
		mongoEngine:        nil,                  // 将通过 SetMongoEngine 设置
		grayscaleDecision:  grayscaleDecision,    // 灰度决策器
	}
}

// RegisterProvider 注册搜索提供者
func (s *SearchService) RegisterProvider(p provider.Provider) {
	s.providers[p.Type()] = p
	s.logger.Printf("Registered search provider: %s", p.Type())
}

// SetCache 设置缓存
func (s *SearchService) SetCache(c cache.Cache) {
	s.cache = c
}

// SetESConfig 设置 ES 配置
func (s *SearchService) SetESConfig(config *SearchConfig) {
	s.esConfig = config

	// 更新灰度百分比指标
	if config != nil && config.ES.GrayScale.Enabled {
		s.updateGrayscaleMetrics()
	}
}

// updateGrayscaleMetrics 更新灰度指标
func (s *SearchService) updateGrayscaleMetrics() {
	if s.esConfig == nil {
		return
	}

	// 为每个搜索类型更新灰度百分比
	percent := s.esConfig.ES.GrayScale.Percent
	searchTypes := []string{"books", "projects", "documents"}

	for _, searchType := range searchTypes {
		metrics.UpdateGrayscalePercent(searchType, percent)
	}

	s.logger.Printf("[SearchService] Updated grayscale metrics: percent=%d", percent)
}

// SetESEngine 设置 ES 引擎
func (s *SearchService) SetESEngine(engine searchengine.Engine) {
	s.esEngine = engine
}

// SetMongoEngine 设置 MongoDB 引擎
func (s *SearchService) SetMongoEngine(engine searchengine.Engine) {
	s.mongoEngine = engine
}

// shouldUseES 判断是否应该使用 ES（灰度逻辑）
func (s *SearchService) shouldUseES(ctx context.Context, searchType search.SearchType, userID string) bool {
	// 如果 ES 未启用，使用 MongoDB
	if s.esConfig == nil || !s.esConfig.ES.Enabled {
		return false
	}

	// 如果 ES 引擎未初始化，使用 MongoDB
	if s.esEngine == nil {
		return false
	}

	// 如果未启用灰度，直接使用 ES
	if !s.esConfig.ES.GrayScale.Enabled {
		return true
	}

	// 使用灰度决策器
	if s.grayscaleDecision != nil {
		return s.grayscaleDecision.ShouldUseES(ctx, string(searchType), userID)
	}

	// 默认使用 MongoDB
	return false
}

// getUserID 从 context 获取用户 ID
func (s *SearchService) getUserID(ctx context.Context) string {
	if uid, exists := ctx.Value("userId").(string); exists {
		return uid
	}
	if uid, exists := ctx.Value("user_id").(string); exists {
		return uid
	}
	return ""
}

// searchWithProvider 使用 Provider 搜索（原有逻辑）
func (s *SearchService) searchWithProvider(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	// 获取对应的 Provider
	prov, err := s.GetProvider(req.Type)
	if err != nil {
		s.logger.Printf("[SearchService] Provider not found: %s", req.Type)
		return nil, fmt.Errorf("provider not found: %s", req.Type)
	}

	// Provider 级别验证
	if err := prov.Validate(req); err != nil {
		s.logger.Printf("[SearchService] Provider validation failed: %v", err)
		return nil, fmt.Errorf("provider validation failed: %w", err)
	}

	// 执行搜索
	return prov.Search(ctx, req)
}

// searchWithES 使用 ES 搜索
func (s *SearchService) searchWithES(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	// 根据搜索类型确定索引名称
	indexName := s.getESIndexName(req.Type)
	if indexName == "" {
		s.logger.Printf("[SearchService] Unsupported search type for ES: %s", req.Type)
		return nil, fmt.Errorf("unsupported search type for ES: %s", req.Type)
	}

	// 构建搜索选项
	opts := &searchengine.SearchOptions{
		From: (req.Page - 1) * req.PageSize,
		Size: req.PageSize,
	}

	// 添加排序
	if req.Sort != nil {
		for _, sortField := range req.Sort {
			opts.Sort = append(opts.Sort, searchengine.SortField{
				Field:     sortField.Field,
				Ascending: sortField.Ascending,
			})
		}
	}

	// 添加过滤条件
	if req.Filter != nil {
		opts.Filter = req.Filter
	}

	// 执行 ES 搜索
	result, err := s.esEngine.Search(ctx, indexName, req.Query, opts)
	if err != nil {
		s.logger.Printf("[SearchService] ES search failed: %v", err)
		return nil, fmt.Errorf("ES search failed: %w", err)
	}

	// 转换为 SearchResponse
	return s.convertESToSearchResponse(result, req.Type)
}

// getESIndexName 根据搜索类型获取 ES 索引名称
func (s *SearchService) getESIndexName(searchType search.SearchType) string {
	// 使用配置的索引前缀
	prefix := ""
	if s.esConfig != nil && s.esConfig.ES.IndexPrefix != "" {
		prefix = s.esConfig.ES.IndexPrefix + "_"
	}

	switch searchType {
	case search.SearchTypeBooks:
		return prefix + "books"
	case search.SearchTypeProjects:
		return prefix + "projects"
	case search.SearchTypeDocuments:
		return prefix + "documents"
	case search.SearchTypeUsers:
		return prefix + "users"
	default:
		return ""
	}
}

// convertESToSearchResponse 将 ES 搜索结果转换为 SearchResponse
func (s *SearchService) convertESToSearchResponse(result *searchengine.SearchResult, searchType search.SearchType) (*search.SearchResponse, error) {
	results := make([]search.SearchItem, 0, len(result.Hits))

	for _, hit := range result.Hits {
		// 将 ES 的 Source 转换为对应类型的模型
		item := s.convertHitToItem(hit, searchType)
		if item != nil {
			results = append(results, *item)
		}
	}

	return &search.SearchResponse{
		Success: true,
		Data: &search.SearchData{
			Type:     searchType,
			Total:    result.Total,
			Results:  results,
			Page:     0, // 需要从外部设置
			PageSize: len(results),
			Took:     result.Took,
		},
	}, nil
}

// convertHitToItem 将 ES 命中项转换为对应的模型
func (s *SearchService) convertHitToItem(hit searchengine.Hit, searchType search.SearchType) *search.SearchItem {
	// 这里需要根据不同的搜索类型转换为对应的模型
	// 暂时返回通用的 SearchItem，实际使用时应该转换为具体的模型结构
	item := &search.SearchItem{
		ID:    hit.ID,
		Score: hit.Score,
		Data:  hit.Source,
	}

	// 添加高亮信息
	if hit.Highlight != nil && len(hit.Highlight) > 0 {
		item.Highlight = hit.Highlight
	}

	return item
}

// Search 统一搜索入口
func (s *SearchService) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	startTime := time.Now()

	// 1. 验证请求参数
	if err := s.validateRequest(req); err != nil {
		s.logger.Printf("[SearchService] Invalid request: %v", err)
		return s.errorResponse(search.ErrInvalidRequest), nil
	}

	// 2. 获取用户 ID
	userID := s.getUserID(ctx)

	// 3. 判断是否使用 ES（灰度逻辑）
	useES := s.shouldUseES(ctx, req.Type, userID)

	// 记录灰度决策
	metrics.RecordGrayscaleDecision(string(req.Type), useES)

	// 4. 生成缓存键（基于是否使用 ES）
	cacheKey := s.generateCacheKeyWithEngine(req, useES)

	// 5. 检查缓存（如果启用）
	var cacheHit bool
	if s.config.EnableCache && s.cache != nil {
		cachedData, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			cacheHit = true
			engineName := "MongoDB"
			if useES {
				engineName = "ES"
			}
			s.logger.Printf("[SearchService] Cache hit for key: %s (engine: %s)", cacheKey, engineName)
			// 从缓存反序列化响应
			resp := &search.SearchResponse{}
			if err := s.unmarshalResponse(cachedData, resp); err == nil {
				resp.Meta = &search.MetaInfo{
					RequestID: s.generateRequestID(),
					TookMs:    0, // 缓存命中不计算耗时
				}
				resp.Data.Took = 0 // 缓存命中耗时为 0
				return resp, nil
			}
		}
	}

	// 6. 执行搜索
	start := time.Now()
	var resp *search.SearchResponse
	var err error

	if useES && s.esEngine != nil {
		// 使用 ES 搜索
		resp, err = s.searchWithES(ctx, req)
		searchTook := time.Since(startTime)

		// 记录 ES 使用情况（原有灰度决策器）
		if s.grayscaleDecision != nil {
			s.grayscaleDecision.RecordUsage("elasticsearch", searchTook)
		}

		// 记录 Prometheus 指标
		metrics.RecordSearch("elasticsearch", searchTook)
	} else {
		// 使用 Provider（原有逻辑）
		resp, err = s.searchWithProvider(ctx, req)
		searchTook := time.Since(startTime)

		// 记录 MongoDB 使用情况（原有灰度决策器）
		if s.grayscaleDecision != nil {
			s.grayscaleDecision.RecordUsage("mongodb", searchTook)
		}

		// 记录 Prometheus 指标
		metrics.RecordSearch("mongodb", searchTook)
	}

	if err != nil {
		engineName := "MongoDB"
		if useES {
			engineName = "ES"
		}
		s.logger.Printf("[SearchService] Search failed: type=%s, engine=%s, error=%v", req.Type, engineName, err)
		return nil, search.WrapError(err, search.ErrCodeEngineFailure, "Search execution failed")
	}

	// 7. 记录耗时
	took := time.Since(start)
	if resp.Data != nil {
		resp.Data.Took = took
	}

	// 8. 更新缓存
	if s.config.EnableCache && s.cache != nil && !cacheHit && resp.Success && resp.Data != nil {
		if data, err := s.marshalResponse(resp); err == nil {
			ttl := s.calculateTTL(req)
			if err := s.cache.Set(ctx, cacheKey, data, ttl); err != nil {
				s.logger.Printf("[SearchService] Failed to set cache: %v", err)
			} else {
				s.logger.Printf("[SearchService] Cached response with TTL: %v", ttl)
			}
		}
	}

	// 9. 添加元信息
	if resp.Meta == nil {
		resp.Meta = &search.MetaInfo{}
	}
	resp.Meta.RequestID = s.generateRequestID()
	resp.Meta.TookMs = took.Milliseconds()

	// 10. 记录搜索日志
	engineName := "MongoDB"
	if useES {
		engineName = "ES"
	}
	s.logger.Printf("[SearchService] Search completed: type=%s, engine=%s, query=%s, took=%v, total=%d, cache_hit=%v",
		req.Type, engineName, req.Query, took, resp.Data.Total, cacheHit)

	return resp, nil
}

// SearchBatch 批量搜索（并发执行）
func (s *SearchService) SearchBatch(ctx context.Context, reqs []*search.SearchRequest) ([]*search.SearchResponse, error) {
	if len(reqs) == 0 {
		return []*search.SearchResponse{}, nil
	}

	// 限制并发数
	semaphore := make(chan struct{}, s.config.MaxConcurrentSearches)
	var wg sync.WaitGroup

	responses := make([]*search.SearchResponse, len(reqs))
	errors := make([]error, len(reqs))

	for i, req := range reqs {
		wg.Add(1)
		go func(idx int, r *search.SearchRequest) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			resp, err := s.Search(ctx, r)
			if err != nil {
				errors[idx] = err
				responses[idx] = s.errorResponse(search.WrapError(err, search.ErrCodeEngineFailure, "Batch search failed"))
			} else {
				responses[idx] = resp
			}
		}(i, req)
	}

	wg.Wait()

	// 检查是否有错误
	for _, err := range errors {
		if err != nil {
			return responses, err
		}
	}

	return responses, nil
}

// GetProvider 获取指定类型的 Provider
func (s *SearchService) GetProvider(searchType search.SearchType) (provider.Provider, error) {
	p, ok := s.providers[searchType]
	if !ok {
		return nil, fmt.Errorf("provider not found for type: %s", searchType)
	}
	return p, nil
}

// ListProviders 列出所有已注册的 Provider
func (s *SearchService) ListProviders() []search.SearchType {
	types := make([]search.SearchType, 0, len(s.providers))
	for t := range s.providers {
		types = append(types, t)
	}
	return types
}

// validateRequest 验证请求参数
func (s *SearchService) validateRequest(req *search.SearchRequest) error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}
	if req.Type == "" {
		return fmt.Errorf("search type is required")
	}
	if req.Query == "" {
		return fmt.Errorf("query is required")
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	return nil
}

// generateCacheKey 生成缓存键
func (s *SearchService) generateCacheKey(req *search.SearchRequest) string {
	return s.generateCacheKeyWithEngine(req, false)
}

// generateCacheKeyWithEngine 生成包含引擎信息的缓存键
func (s *SearchService) generateCacheKeyWithEngine(req *search.SearchRequest, useES bool) string {
	// 序列化过滤条件
	filterData, _ := json.Marshal(req.Filter)
	filterHash := md5.Sum(filterData)

	// 序列化排序条件
	sortData, _ := json.Marshal(req.Sort)
	sortHash := md5.Sum(sortData)

	// 引擎标识
	engine := "mongodb"
	if useES {
		engine = "es"
	}

	// 生成缓存键: search:{engine}:{type}:{query}:{filter_hash}:{sort_hash}:{page}:{page_size}
	return fmt.Sprintf("search:%s:%s:%s:%x:%x:%d:%d",
		engine, req.Type, req.Query, filterHash, sortHash, req.Page, req.PageSize)
}

// calculateTTL 计算缓存过期时间
func (s *SearchService) calculateTTL(req *search.SearchRequest) time.Duration {
	// 第一页结果缓存时间更长
	if req.Page == 1 {
		return time.Duration(s.config.DefaultCacheTTL) * time.Second * 2
	}
	return time.Duration(s.config.DefaultCacheTTL) * time.Second
}

// generateRequestID 生成请求 ID
func (s *SearchService) generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// marshalResponse 序列化响应
func (s *SearchService) marshalResponse(resp *search.SearchResponse) ([]byte, error) {
	return json.Marshal(resp)
}

// unmarshalResponse 反序列化响应
func (s *SearchService) unmarshalResponse(data []byte, resp *search.SearchResponse) error {
	return json.Unmarshal(data, resp)
}

// errorResponse 创建错误响应
func (s *SearchService) errorResponse(err *search.SearchError) *search.SearchResponse {
	return &search.SearchResponse{
		Success: false,
		Error: &search.ErrorInfo{
			Code:    err.Code,
			Message: err.Message,
			Details: err.Error(),
		},
		Meta: &search.MetaInfo{
			RequestID: s.generateRequestID(),
			TookMs:    0,
		},
	}
}

// Health 健康检查
func (s *SearchService) Health(ctx context.Context) map[string]error {
	status := make(map[string]error)

	// 检查缓存健康状态
	if s.cache != nil {
		if err := s.cache.Ping(ctx); err != nil {
			status["cache"] = err
		} else {
			status["cache"] = nil
		}
	}

	// 检查各 Provider 健康状态
	for searchType := range s.providers {
		// 这里假设 Provider 实现了 Health 方法
		// 如果没有，可以跳过或返回 nil
		status[string(searchType)] = nil
	}

	return status
}

// Stats 获取搜索服务统计信息
func (s *SearchService) Stats(ctx context.Context) map[string]interface{} {
	stats := make(map[string]interface{})

	// Provider 统计
	stats["providers"] = s.ListProviders()

	// 缓存统计
	if s.cache != nil {
		if cacheStats, ok := s.cache.(cache.CacheStatsProvider); ok {
			if statsData, err := cacheStats.Stats(ctx); err == nil {
				stats["cache"] = statsData
			}
		}
	}

	// 配置信息
	stats["config"] = map[string]interface{}{
		"enable_cache":            s.config.EnableCache,
		"default_cache_ttl":       s.config.DefaultCacheTTL,
		"max_concurrent_searches": s.config.MaxConcurrentSearches,
	}

	return stats
}

// InvalidateCache 使缓存失效
func (s *SearchService) InvalidateCache(ctx context.Context, searchType search.SearchType) error {
	if s.cache == nil {
		return fmt.Errorf("cache not initialized")
	}

	pattern := fmt.Sprintf("search:%s:*", searchType)
	return s.cache.DeletePattern(ctx, pattern)
}

// ClearAllCache 清空所有搜索缓存
func (s *SearchService) ClearAllCache(ctx context.Context) error {
	if s.cache == nil {
		return fmt.Errorf("cache not initialized")
	}

	return s.cache.Clear(ctx)
}

// Close 关闭搜索服务
func (s *SearchService) Close() error {
	var errs []error

	// 关闭缓存连接
	if s.cache != nil {
		if err := s.cache.Close(); err != nil {
			errs = append(errs, fmt.Errorf("cache close error: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiple errors during close: %v", errs)
	}

	return nil
}
