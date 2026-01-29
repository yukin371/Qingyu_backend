package search

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	searchModels "Qingyu_backend/models/search"
	"Qingyu_backend/pkg/logger"
	"Qingyu_backend/pkg/response"
	searchService "Qingyu_backend/service/search"
)

// SearchAPI 统一搜索 API
type SearchAPI struct {
	searchService *searchService.SearchService
}

// NewSearchAPI 创建搜索 API 实例
func NewSearchAPI(searchService *searchService.SearchService) *SearchAPI {
	return &SearchAPI{
		searchService: searchService,
	}
}

// SearchRequest HTTP 搜索请求
type SearchRequest struct {
	Type     string                 `json:"type" binding:"required"`     // 搜索类型：books, projects, documents, users
	Query    string                 `json:"query" binding:"required"`    // 搜索关键词
	Filter   map[string]interface{} `json:"filter"`                      // 过滤条件
	Sort     []SortFieldRequest     `json:"sort"`                        // 排序字段
	Page     int                    `json:"page"`                        // 页码，默认 1
	PageSize int                    `json:"page_size"`                   // 每页数量，默认 20
}

// SortFieldRequest 排序字段请求
type SortFieldRequest struct {
	Field     string `json:"field"`     // 排序字段名
	Ascending bool   `json:"ascending"` // 是否升序
}

// BatchSearchRequest 批量搜索请求
type BatchSearchRequest struct {
	Queries []SearchRequest `json:"queries" binding:"required,min=1"` // 搜索查询列表
}

// Search 统一搜索入口
//
//	@Summary		统一搜索
//	@Description	根据类型和关键词进行统一搜索，支持书籍、项目、文档、用户等多种类型
//	@Tags			搜索
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SearchRequest	true	"搜索请求"
//	@Success		200		{object}	object
//	@Failure		400		{object}	object
//	@Failure		500		{object}	object
//	@Router			/api/v1/search/search [post]
func (api *SearchAPI) Search(c *gin.Context) {
	startTime := time.Now()

	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 获取用户ID
	var userID string
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	// 构建日志记录器
	searchLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	if userID != "" {
		searchLogger = searchLogger.WithUser(userID)
	}

	// 绑定请求参数
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		searchLogger.WithModule("search").Warn("搜索参数绑定失败",
			zap.Error(err),
		)
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 参数验证
	if req.Type == "" {
		searchLogger.WithModule("search").Warn("搜索类型为空")
		response.BadRequest(c, "参数错误", "搜索类型不能为空")
		return
	}

	if req.Query == "" {
		searchLogger.WithModule("search").Warn("搜索关键词为空")
		response.BadRequest(c, "参数错误", "搜索关键词不能为空")
		return
	}

	// 设置默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 记录搜索请求
	searchLogger.WithModule("search").Info("统一搜索请求",
		zap.String("type", req.Type),
		zap.String("query", req.Query),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
	)

	// 转换为 Service 层请求
	serviceReq := &searchModels.SearchRequest{
		Type:     searchModels.SearchType(req.Type),
		Query:    req.Query,
		Filter:   req.Filter,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 转换排序字段
	if len(req.Sort) > 0 {
		serviceReq.Sort = make([]searchModels.SortField, len(req.Sort))
		for i, s := range req.Sort {
			serviceReq.Sort[i] = searchModels.SortField{
				Field:     s.Field,
				Ascending: s.Ascending,
			}
		}
	}

	// 执行搜索
	resp, err := api.searchService.Search(c.Request.Context(), serviceReq)
	if err != nil {
		duration := time.Since(startTime)
		searchLogger.WithModule("search").Error("搜索执行失败",
			zap.Error(err),
			zap.Duration("took", duration),
		)
		response.InternalError(c, err)
		return
	}

	// 计算耗时
	duration := time.Since(startTime)

	// 检查搜索是否成功
	if !resp.Success {
		searchLogger.WithModule("search").Warn("搜索返回失败",
			zap.String("error_code", resp.Error.Code),
			zap.String("error_message", resp.Error.Message),
			zap.Duration("took", duration),
		)
		response.BadRequest(c, resp.Error.Message, resp.Error.Details)
		return
	}

	// 记录搜索结果
	searchLogger.WithModule("search").Info("搜索成功",
		zap.String("type", string(resp.Data.Type)),
		zap.Int64("total", resp.Data.Total),
		zap.Int("returned", len(resp.Data.Results)),
		zap.Duration("took", duration),
		zap.Int64("took_ms", resp.Meta.TookMs),
	)

	// 构造响应数据
	responseData := gin.H{
		"type":      string(resp.Data.Type),
		"total":     resp.Data.Total,
		"page":      resp.Data.Page,
		"page_size": resp.Data.PageSize,
		"results":   resp.Data.Results,
		"took":      resp.Data.Took.String(),
	}

	response.SuccessWithMessage(c, "搜索成功", responseData)
}

// SearchBatch 批量搜索
//
//	@Summary		批量搜索
//	@Description	一次请求执行多个搜索查询，并发执行
//	@Tags			搜索
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchSearchRequest	true	"批量搜索请求"
//	@Success		200		{object}	object
//	@Failure		400		{object}	object
//	@Failure		500		{object}	object
//	@Router			/api/v1/search/batch [post]
func (api *SearchAPI) SearchBatch(c *gin.Context) {
	startTime := time.Now()

	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 构建日志记录器
	searchLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	// 绑定请求参数
	var req BatchSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		searchLogger.WithModule("search").Warn("批量搜索参数绑定失败",
			zap.Error(err),
		)
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 参数验证
	if len(req.Queries) == 0 {
		searchLogger.WithModule("search").Warn("批量搜索查询列表为空")
		response.BadRequest(c, "参数错误", "查询列表不能为空")
		return
	}

	// 记录批量搜索请求
	searchLogger.WithModule("search").Info("批量搜索请求",
		zap.Int("query_count", len(req.Queries)),
	)

	// 转换为 Service 层请求
	serviceReqs := make([]*searchModels.SearchRequest, len(req.Queries))
	for i, q := range req.Queries {
		// 参数验证
		if q.Type == "" || q.Query == "" {
			searchLogger.WithModule("search").Warn("批量搜索查询参数不完整",
				zap.Int("index", i),
				zap.String("type", q.Type),
				zap.String("query", q.Query),
			)
			response.BadRequest(c, "参数错误", "查询类型和关键词不能为空")
			return
		}

		// 设置默认分页参数
		if q.Page < 1 {
			q.Page = 1
		}
		if q.PageSize < 1 {
			q.PageSize = 20
		}
		if q.PageSize > 100 {
			q.PageSize = 100
		}

		serviceReqs[i] = &searchModels.SearchRequest{
			Type:     searchModels.SearchType(q.Type),
			Query:    q.Query,
			Filter:   q.Filter,
			Page:     q.Page,
			PageSize: q.PageSize,
		}

		// 转换排序字段
		if len(q.Sort) > 0 {
			serviceReqs[i].Sort = make([]searchModels.SortField, len(q.Sort))
			for j, s := range q.Sort {
				serviceReqs[i].Sort[j] = searchModels.SortField{
					Field:     s.Field,
					Ascending: s.Ascending,
				}
			}
		}
	}

	// 执行批量搜索
	responses, err := api.searchService.SearchBatch(c.Request.Context(), serviceReqs)
	if err != nil {
		duration := time.Since(startTime)
		searchLogger.WithModule("search").Error("批量搜索执行失败",
			zap.Error(err),
			zap.Duration("took", duration),
		)
		response.InternalError(c, err)
		return
	}

	// 计算耗时
	duration := time.Since(startTime)

	// 记录批量搜索结果
	searchLogger.WithModule("search").Info("批量搜索成功",
		zap.Int("total_queries", len(responses)),
		zap.Duration("took", duration),
	)

	// 转换响应格式
	results := make([]gin.H, len(responses))
	for i, resp := range responses {
		result := gin.H{
			"success": resp.Success,
		}

		if resp.Success && resp.Data != nil {
			result["data"] = gin.H{
				"type":      string(resp.Data.Type),
				"total":     resp.Data.Total,
				"page":      resp.Data.Page,
				"page_size": resp.Data.PageSize,
				"results":   resp.Data.Results,
				"took":      resp.Data.Took.String(),
			}
		}

		if resp.Error != nil {
			result["error"] = gin.H{
				"code":    resp.Error.Code,
				"message": resp.Error.Message,
				"details": resp.Error.Details,
			}
		}

		if resp.Meta != nil {
			result["request_id"] = resp.Meta.RequestID
			result["took_ms"] = resp.Meta.TookMs
		}

		results[i] = result
	}

	response.SuccessWithMessage(c, "批量搜索成功", results)
}

// Health 健康检查
//
//	@Summary		搜索服务健康检查
//	@Description	检查搜索服务及其各组件的健康状态
//	@Tags			搜索
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	object
//	@Router			/api/v1/search/health [get]
func (api *SearchAPI) Health(c *gin.Context) {
	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 执行健康检查
	status := api.searchService.Health(c.Request.Context())

	// 构造响应
	healthStatus := gin.H{
		"status": "healthy",
	}

	// 检查是否有错误
	hasError := false
	for component, err := range status {
		if err != nil {
			hasError = true
			healthStatus[component] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			healthStatus[component] = gin.H{
				"status": "healthy",
			}
		}
	}

	if hasError {
		healthStatus["status"] = "degraded"
	}

	response.SuccessWithMessage(c, "健康检查完成", healthStatus)
}
