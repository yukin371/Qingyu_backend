package system

import (

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service"

	"github.com/gin-gonic/gin"
	"Qingyu_backend/pkg/response"
	"errors"
)

// HealthAPI 健康检查API
type HealthAPI struct{}

// NewHealthAPI 创建健康检查API
func NewHealthAPI() *HealthAPI {
	return &HealthAPI{}
}

// SystemHealth 系统整体健康检查
// @Summary 系统健康检查
// @Description 检查系统整体健康状态
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse "健康"
// @Failure 500 {object} shared.ErrorResponse "不健康"
// @Router /system/health [get]
func (api *HealthAPI) SystemHealth(c *gin.Context) {
	container := service.GetServiceContainer()
	if container == nil {
		response.InternalError(c, errors.New("服务容器未初始化: "))
		return
	}

	// 获取所有服务健康状态
	healthStatus := container.GetAllServicesHealth(c.Request.Context())

	// 检查是否所有服务都健康
	allHealthy := true
	for _, healthy := range healthStatus {
		if !healthy {
			allHealthy = false
			break
		}
	}

	if allHealthy {
		response.SuccessWithMessage(c, "系统健康", gin.H{
			"status":   "healthy",
			"services": healthStatus,
		})
	} else {
		shared.Error(c, http.StatusServiceUnavailable, "部分服务不健康", "部分服务健康检查失败")
	}
}

// ServiceHealth 特定服务健康检查
// @Summary 特定服务健康检查
// @Description 检查指定服务的健康状态
// @Tags 系统监控
// @Accept json
// @Produce json
// @Param service path string true "服务名称"
// @Success 200 {object} shared.APIResponse "健康"
// @Failure 404 {object} shared.ErrorResponse "服务不存在"
// @Failure 500 {object} shared.ErrorResponse "不健康"
// @Router /system/health/{service} [get]
func (api *HealthAPI) ServiceHealth(c *gin.Context) {
	serviceName := c.Param("service")

	container := service.GetServiceContainer()
	if container == nil {
		response.InternalError(c, errors.New("服务容器未初始化: "))
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
// @Success 200 {object} shared.APIResponse "成功"
// @Failure 500 {object} shared.ErrorResponse "失败"
// @Router /system/metrics [get]
func (api *HealthAPI) AllMetrics(c *gin.Context) {
	container := service.GetServiceContainer()
	if container == nil {
		response.InternalError(c, errors.New("服务容器未初始化: "))
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
// @Success 200 {object} shared.APIResponse "成功"
// @Failure 404 {object} shared.ErrorResponse "服务不存在"
// @Failure 500 {object} shared.ErrorResponse "失败"
// @Router /system/metrics/{service} [get]
func (api *HealthAPI) ServiceMetrics(c *gin.Context) {
	serviceName := c.Param("service")

	container := service.GetServiceContainer()
	if container == nil {
		response.InternalError(c, errors.New("服务容器未初始化: "))
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
