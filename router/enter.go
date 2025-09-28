package router

import (
	"Qingyu_backend/router/ai"
	bookstoreRouter "Qingyu_backend/router/bookstore"
	projectRouter "Qingyu_backend/router/project"
	userRouter "Qingyu_backend/router/users"

	readingAPI "Qingyu_backend/api/v1/reading"
	aiService "Qingyu_backend/service/ai"
	bookstoreService "Qingyu_backend/service/bookstore"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本组
	v1 := r.Group("/api/v1")

	// 注册书城路由
	bookstoreSvc := bookstoreService.NewBookstoreService(nil,
		nil,
		nil,
		nil) // TODO: 注入实际的依赖
	bookstoreAPI := readingAPI.NewBookstoreAPI(bookstoreSvc)
	bookstoreRouterInstance := bookstoreRouter.NewBookstoreRouter(bookstoreAPI)
	bookstoreRouterInstance.RegisterRoutes(v1)
	bookstoreRouterInstance.RegisterPublicRoutes(v1)
	bookstoreRouterInstance.RegisterPrivateRoutes(v1)

	// 注册系统路由（用户认证等）
	userRouter.RegisterRoutes(v1)

	// 注册文档路由
	projectRouter.RegisterRoutes(v1)

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
