package router

import (
	"Qingyu_backend/router/ai"
	"Qingyu_backend/router/document"
	"Qingyu_backend/router/system"

	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本组
	v1 := r.Group("/api/v1")

	// 注册系统路由（用户认证等）
	system.RegisterRoutes(v1)

	// 注册文档路由
	document.RegisterRoutes(v1)

	// 注册AI路由
	aiSvc := aiService.NewService()
	aiRouter := ai.NewAIRouter(aiSvc)
	aiRouter.InitAIRouter(v1)

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
