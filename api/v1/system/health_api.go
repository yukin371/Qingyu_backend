package system

import (
	"net/http"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service"
	"Qingyu_backend/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
)

// HealthAPI 健康检查API
type HealthAPI struct{}

// NewHealthAPI 创建健康检查API
func NewHealthAPI() *HealthAPI {
	return &HealthAPI{}
}

// SystemHealth 系统整体健康检查
// @Summary 系统健康检查
// @Description 检查系统整体健康状态，包括 MongoDB、Redis、Milvus、AI gRPC 服务
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse "健康"
// @Failure 503 {object} response.APIResponse "不健康"
// @Router /api/v1/system/health [get]
func (api *HealthAPI) SystemHealth(c *gin.Context) {
	svcContainer := service.GetServiceContainer()
	if svcContainer == nil {
		c.Error(errors.New("服务容器未初始化"))
		return
	}

	// 获取基础设施健康状态（含延迟检测）
	infraHealth := svcContainer.GetInfrastructureHealth(c.Request.Context())

	// 根据整体状态返回不同的 HTTP 状态码
	switch infraHealth.Status {
	case "unhealthy":
		// 关键服务（MongoDB）不健康 → 503
		response.JSON(c, http.StatusServiceUnavailable, gin.H{
			"code":    http.StatusServiceUnavailable,
			"message": "服务不健康：关键依赖不可用",
			"data": gin.H{
				"status":         infraHealth.Status,
				"services":       infraHealth.Services,
				"version":        infraHealth.Version,
				"uptime_seconds": infraHealth.UptimeSeconds,
			},
		})
	case "degraded":
		// 非关键服务不健康但关键服务正常 → 200 + degraded
		response.SuccessWithMessage(c, "服务降级：部分非关键依赖不可用", gin.H{
			"status":         infraHealth.Status,
			"services":       infraHealth.Services,
			"version":        infraHealth.Version,
			"uptime_seconds": infraHealth.UptimeSeconds,
		})
	default:
		// 所有服务健康 → 200 + healthy
		response.SuccessWithMessage(c, "系统健康", gin.H{
			"status":         infraHealth.Status,
			"services":       infraHealth.Services,
			"version":        infraHealth.Version,
			"uptime_seconds": infraHealth.UptimeSeconds,
		})
	}
}

// ServiceHealth 特定服务健康检查
// @Summary 特定服务健康检查
// @Description 检查指定服务的健康状态
// @Tags 系统监控
// @Accept json
// @Produce json
// @Param service path string true "服务名称"
// @Success 200 {object} response.APIResponse "健康"
// @Failure 404 {object} response.APIResponse "服务不存在"
// @Failure 500 {object} response.APIResponse "不健康"
// @Router /api/v1/system/health/{service} [get]
func (api *HealthAPI) ServiceHealth(c *gin.Context) {
	serviceName := c.Param("service")

	container := service.GetServiceContainer()
	if container == nil {
		c.Error(errors.New("服务容器未初始化: "))
		return
	}

	// 获取服务
	svc, err := container.GetService(serviceName)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "服务不存在: "+serviceName, "")
		return
	}

	// 检查健康状态
	err = svc.Health(c.Request.Context())
	if err != nil {
		shared.Error(c, http.StatusServiceUnavailable, "服务不健康", err.Error())
		return
	}

	response.SuccessWithMessage(c, "服务健康", gin.H{
		"service": serviceName,
		"status":  "healthy",
	})
}

// AllMetrics 获取所有服务指标
// @Summary 获取所有服务指标
// @Description 获取所有服务的运行指标
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse "成功"
// @Failure 500 {object} response.APIResponse "失败"
// @Router /api/v1/system/metrics [get]
func (api *HealthAPI) AllMetrics(c *gin.Context) {
	container := service.GetServiceContainer()
	if container == nil {
		c.Error(errors.New("服务容器未初始化: "))
		return
	}

	// 获取所有服务指标
	allMetrics := container.GetAllServicesMetrics()

	// 转换为响应格式
	metricsResponse := make(map[string]interface{})
	for name, metric := range allMetrics {
		metricsResponse[name] = metric.ToResponse()
	}

	response.SuccessWithMessage(c, "获取成功", gin.H{
		"total_services": len(allMetrics),
		"metrics":        metricsResponse,
	})
}

// ServiceMetrics 获取特定服务指标
// @Summary 获取特定服务指标
// @Description 获取指定服务的运行指标
// @Tags 系统监控
// @Accept json
// @Produce json
// @Param service path string true "服务名称"
// @Success 200 {object} response.APIResponse "成功"
// @Failure 404 {object} response.APIResponse "服务不存在"
// @Failure 500 {object} response.APIResponse "失败"
// @Router /api/v1/system/metrics/{service} [get]
func (api *HealthAPI) ServiceMetrics(c *gin.Context) {
	serviceName := c.Param("service")

	container := service.GetServiceContainer()
	if container == nil {
		c.Error(errors.New("服务容器未初始化: "))
		return
	}

	// 获取服务指标
	metric, err := container.GetServiceMetrics(serviceName)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "服务不存在: "+serviceName, "")
		return
	}

	response.SuccessWithMessage(c, "获取成功", metric.ToResponse())
}
