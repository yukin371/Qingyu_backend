package api

import (
	"Qingyu_backend/handlers"
	"Qingyu_backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(router *gin.Engine) {
	// 基础路由
	router.GET("/ping", PingHandler)

	// 用户认证路由 - 不需要认证
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/login", handlers.LoginHandler)
	}

	// 文档路由 - 需要认证
	docs := router.Group("/api/documents")
	docs.Use(middleware.AuthMiddleware())
	{
		docs.POST("", handlers.CreateDocumentHandler)       // 创建文档
		docs.GET("", handlers.GetDocumentsHandler)          // 获取所有文档
		docs.GET("/:id", handlers.GetDocumentHandler)       // 获取单个文档
		docs.PUT("/:id", handlers.UpdateDocumentHandler)    // 更新文档
		docs.DELETE("/:id", handlers.DeleteDocumentHandler) // 删除文档
	}
}
