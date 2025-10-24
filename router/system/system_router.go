package system

import (
	systemAPI "Qingyu_backend/api/v1/system"

	"github.com/gin-gonic/gin"
)

// InitSystemRoutes 初始化系统路由
func InitSystemRoutes(router *gin.RouterGroup) {
	healthAPI := systemAPI.NewHealthAPI()

	systemGroup := router.Group("/system")
	{
		// 健康检查路由（公开）
		systemGroup.GET("/health", healthAPI.SystemHealth)
		systemGroup.GET("/health/:service", healthAPI.ServiceHealth)

		// 指标路由（公开，后续可以添加权限控制）
		systemGroup.GET("/metrics", healthAPI.AllMetrics)
		systemGroup.GET("/metrics/:service", healthAPI.ServiceMetrics)
	}
}
