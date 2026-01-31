package search

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"
	"Qingyu_backend/pkg/response"
	searchService "Qingyu_backend/service/search"
)

// GrayscaleAPI 灰度配置 API
type GrayscaleAPI struct {
	searchService *searchService.SearchService
}

// NewGrayscaleAPI 创建灰度配置 API 实例
func NewGrayscaleAPI(searchService *searchService.SearchService) *GrayscaleAPI {
	return &GrayscaleAPI{
		searchService: searchService,
	}
}

// UpdateGrayscaleConfigRequest 更新灰度配置请求
type UpdateGrayscaleConfigRequest struct {
	Enabled bool   `json:"enabled" binding:"required"` // 是否启用灰度
	Percent int    `json:"percent" binding:"required,min=0,max=100"` // 灰度百分比(0-100)
	Reason  string `json:"reason"`                     // 更新原因
}

// GrayscaleStatusResponse 灰度状态响应
type GrayscaleStatusResponse struct {
	Enabled         bool    `json:"enabled"`           // 是否启用灰度
	Percent         int     `json:"percent"`           // 灰度百分比
	ESCount         int64   `json:"es_count"`          // ES 使用次数
	MongoDBCount    int64   `json:"mongodb_count"`     // MongoDB 使用次数
	ESTraffic       float64 `json:"es_traffic"`        // ES 流量百分比
	MongoDBTraffic  float64 `json:"mongodb_traffic"`   // MongoDB 流量百分比
	ESAvgTook       int64   `json:"es_avg_took_ms"`    // ES 平均耗时（毫秒）
	MongoDBAvgTook  int64   `json:"mongodb_avg_took_ms"` // MongoDB 平均耗时（毫秒）
}

// GrayscaleMetricsResponse 灰度指标响应
type GrayscaleMetricsResponse struct {
	// 配置
	Config GrayscaleConfig `json:"config"`

	// 使用统计
	ESCount        int64   `json:"es_count"`
	MongoDBCount   int64   `json:"mongodb_count"`

	// 流量分配
	ESTraffic      float64 `json:"es_traffic"`
	MongoDBTraffic float64 `json:"mongodb_traffic"`

	// 性能指标
	ESAvgTook      float64 `json:"es_avg_took_ms"`
	MongoDBAvgTook float64 `json:"mongodb_avg_took_ms"`

	// 决策统计
	TotalDecisions  int64 `json:"total_decisions"`
	ESDecisions     int64 `json:"es_decisions"`
	MongoDBDecisions int64 `json:"mongodb_decisions"`
}

// GrayscaleConfig 灰度配置
type GrayscaleConfig struct {
	Enabled bool `json:"enabled"`
	Percent int  `json:"percent"`
}

// TrafficDistributionResponse 流量分配响应
type TrafficDistributionResponse struct {
	ESTraffic      float64 `json:"es_traffic"`
	MongoDBTraffic float64 `json:"mongodb_traffic"`
	TotalRequests  int64   `json:"total_requests"`
}

// GetGrayscaleStatus 获取灰度状态
//
//	@Summary		获取灰度状态
//	@Description	获取当前灰度配置状态和流量分配统计
//	@Tags			搜索-灰度
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/search/grayscale/status [get]
func (api *GrayscaleAPI) GetGrayscaleStatus(c *gin.Context) {
	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 构建日志记录器
	apiLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	// 获取灰度决策器
	grayscaleDecision := api.searchService.GetGrayscaleDecision()
	if grayscaleDecision == nil {
		apiLogger.WithModule("search").Warn("灰度决策器未初始化")
		response.InternalError(c, nil)
		return
	}

	// 获取灰度配置
	config := grayscaleDecision.GetConfig()

	// 获取灰度指标
	metrics := grayscaleDecision.GetMetrics()

	// 计算流量分配
	esTraffic, mongoTraffic := metrics.GetTrafficDistribution()

	// 转换平均耗时为毫秒
	esAvgTook := metrics.ESAvgTook.Milliseconds()
	mongoAvgTook := metrics.MongoDBAvgTook.Milliseconds()

	// 构造响应
	statusResponse := GrayscaleStatusResponse{
		Enabled:        config.Enabled,
		Percent:        config.Percent,
		ESCount:        metrics.ESCount,
		MongoDBCount:   metrics.MongoDBCount,
		ESTraffic:      esTraffic,
		MongoDBTraffic: mongoTraffic,
		ESAvgTook:      esAvgTook,
		MongoDBAvgTook: mongoAvgTook,
	}

	apiLogger.WithModule("search").Info("获取灰度状态成功",
		zap.Bool("enabled", config.Enabled),
		zap.Int("percent", config.Percent),
		zap.Int64("es_count", metrics.ESCount),
		zap.Int64("mongodb_count", metrics.MongoDBCount),
	)

	response.SuccessWithMessage(c, "获取灰度状态成功", statusResponse)
}

// UpdateGrayscaleConfig 更新灰度配置
//
//	@Summary		更新灰度配置
//	@Description	更新灰度配置，支持运行时热更新，无需重启服务
//	@Tags			搜索-灰度
//	@Accept			json
//	@Produce		json
//	@Param			request	body		UpdateGrayscaleConfigRequest	true	"灰度配置更新请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/search/grayscale/config [post]
func (api *GrayscaleAPI) UpdateGrayscaleConfig(c *gin.Context) {
	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 获取用户ID
	var userID string
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(string)
	}

	// 构建日志记录器
	apiLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	if userID != "" {
		apiLogger = apiLogger.WithUser(userID)
	}

	// 绑定请求参数
	var req UpdateGrayscaleConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiLogger.WithModule("search").Warn("灰度配置更新参数绑定失败",
			zap.Error(err),
		)
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 参数验证（虽然 gin binding 已经处理了，但为了安全再检查一次）
	if req.Percent < 0 || req.Percent > 100 {
		apiLogger.WithModule("search").Warn("灰度百分比超出范围",
			zap.Int("percent", req.Percent),
		)
		response.BadRequest(c, "参数错误", "灰度百分比必须在 0-100 之间")
		return
	}

	// 获取灰度决策器
	grayscaleDecision := api.searchService.GetGrayscaleDecision()
	if grayscaleDecision == nil {
		apiLogger.WithModule("search").Warn("灰度决策器未初始化")
		response.InternalError(c, nil)
		return
	}

	// 记录更新前的配置
	oldConfig := grayscaleDecision.GetConfig()

	// 更新灰度配置
	if err := grayscaleDecision.UpdateConfig(req.Enabled, req.Percent); err != nil {
		apiLogger.WithModule("search").Error("灰度配置更新失败",
			zap.Error(err),
		)
		response.InternalError(c, err)
		return
	}

	// 记录配置更新日志
	apiLogger.WithModule("search").Info("灰度配置已更新",
		zap.Bool("old_enabled", oldConfig.Enabled),
		zap.Bool("new_enabled", req.Enabled),
		zap.Int("old_percent", oldConfig.Percent),
		zap.Int("new_percent", req.Percent),
		zap.String("reason", req.Reason),
		zap.String("operator", userID),
	)

	// 获取更新后的指标
	metrics := grayscaleDecision.GetMetrics()
	esTraffic, mongoTraffic := metrics.GetTrafficDistribution()

	// 构造响应
	updateResponse := gin.H{
		"old_config": gin.H{
			"enabled": oldConfig.Enabled,
			"percent": oldConfig.Percent,
		},
		"new_config": gin.H{
			"enabled": req.Enabled,
			"percent": req.Percent,
		},
		"current_stats": gin.H{
			"es_count":        metrics.ESCount,
			"mongodb_count":   metrics.MongoDBCount,
			"es_traffic":      esTraffic,
			"mongodb_traffic": mongoTraffic,
		},
		"reason": req.Reason,
	}

	response.SuccessWithMessage(c, "灰度配置更新成功", updateResponse)
}

// GetGrayscaleMetrics 获取灰度指标
//
//	@Summary		获取灰度指标
//	@Description	获取完整的灰度监控指标数据，包括配置、使用统计、流量分配、性能指标等
//	@Tags			搜索-灰度
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/search/grayscale/metrics [get]
func (api *GrayscaleAPI) GetGrayscaleMetrics(c *gin.Context) {
	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 构建日志记录器
	apiLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	// 获取灰度决策器
	grayscaleDecision := api.searchService.GetGrayscaleDecision()
	if grayscaleDecision == nil {
		apiLogger.WithModule("search").Warn("灰度决策器未初始化")
		response.InternalError(c, nil)
		return
	}

	// 获取灰度配置
	config := grayscaleDecision.GetConfig()

	// 获取灰度指标
	metrics := grayscaleDecision.GetMetrics()

	// 计算流量分配
	esTraffic, mongoTraffic := metrics.GetTrafficDistribution()

	// 转换平均耗时为毫秒
	esAvgTook := float64(metrics.ESAvgTook.Milliseconds())
	mongoAvgTook := float64(metrics.MongoDBAvgTook.Milliseconds())

	// 计算总决策次数
	totalDecisions := metrics.ESCount + metrics.MongoDBCount

	// 构造响应
	metricsResponse := GrayscaleMetricsResponse{
		Config: GrayscaleConfig{
			Enabled: config.Enabled,
			Percent: config.Percent,
		},
		ESCount:         metrics.ESCount,
		MongoDBCount:    metrics.MongoDBCount,
		ESTraffic:       esTraffic,
		MongoDBTraffic:  mongoTraffic,
		ESAvgTook:       esAvgTook,
		MongoDBAvgTook:  mongoAvgTook,
		TotalDecisions:  totalDecisions,
		ESDecisions:     metrics.ESCount,
		MongoDBDecisions: metrics.MongoDBCount,
	}

	apiLogger.WithModule("search").Info("获取灰度指标成功",
		zap.Bool("enabled", config.Enabled),
		zap.Int("percent", config.Percent),
		zap.Int64("total_decisions", totalDecisions),
	)

	response.SuccessWithMessage(c, "获取灰度指标成功", metricsResponse)
}

// GetTrafficDistribution 获取流量分配
//
//	@Summary		获取流量分配
//	@Description	获取当前流量分配比例数据
//	@Tags			搜索-灰度
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/search/grayscale/traffic [get]
func (api *GrayscaleAPI) GetTrafficDistribution(c *gin.Context) {
	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 构建日志记录器
	apiLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	// 获取灰度决策器
	grayscaleDecision := api.searchService.GetGrayscaleDecision()
	if grayscaleDecision == nil {
		apiLogger.WithModule("search").Warn("灰度决策器未初始化")
		response.InternalError(c, nil)
		return
	}

	// 获取灰度指标
	metrics := grayscaleDecision.GetMetrics()

	// 计算流量分配
	esTraffic, mongoTraffic := metrics.GetTrafficDistribution()

	// 计算总请求数
	totalRequests := metrics.ESCount + metrics.MongoDBCount

	// 构造响应
	trafficResponse := TrafficDistributionResponse{
		ESTraffic:      esTraffic,
		MongoDBTraffic: mongoTraffic,
		TotalRequests:  totalRequests,
	}

	apiLogger.WithModule("search").Info("获取流量分配成功",
		zap.Float64("es_traffic", esTraffic),
		zap.Float64("mongodb_traffic", mongoTraffic),
		zap.Int64("total_requests", totalRequests),
	)

	response.SuccessWithMessage(c, "获取流量分配成功", trafficResponse)
}
