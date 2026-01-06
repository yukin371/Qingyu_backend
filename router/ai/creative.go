package ai

import (
	aiAPI "Qingyu_backend/api/v1/ai"
	"Qingyu_backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitCreativeRoutes 初始化Phase3创作相关路由
func InitCreativeRoutes(router *gin.RouterGroup, creativeAPI *aiAPI.CreativeAPI) {
	// Phase3创作路由组
	creative := router.Group("/creative")
	{
		// 健康检查（公开）
		creative.GET("/health", creativeAPI.HealthCheck)

		// 需要认证的路由
		authGroup := creative.Group("").Use(middleware.JWTAuth())
		{
			// 大纲生成
			authGroup.POST("/outline", creativeAPI.GenerateOutline)

			// 角色生成
			authGroup.POST("/characters", creativeAPI.GenerateCharacters)

			// 情节生成
			authGroup.POST("/plot", creativeAPI.GeneratePlot)

			// 完整工作流
			authGroup.POST("/workflow", creativeAPI.ExecuteCreativeWorkflow)
		}
	}
}
