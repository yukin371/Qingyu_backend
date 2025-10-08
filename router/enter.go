package router

import (
	"log"

	// "Qingyu_backend/router/ai" // 临时禁用
	bookstoreRouter "Qingyu_backend/router/bookstore"
	projectRouter "Qingyu_backend/router/project"
	sharedRouter "Qingyu_backend/router/shared"
	userRouter "Qingyu_backend/router/users"

	readingAPI "Qingyu_backend/api/v1/reading"
	// aiService "Qingyu_backend/service/ai" // 临时禁用
	bookstoreService "Qingyu_backend/service/bookstore"
	"Qingyu_backend/service/shared/container"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本组
	v1 := r.Group("/api/v1")

	// 注册共享服务路由（认证、钱包、存储、管理）
	// 注意：SharedServiceContainer 的服务实现尚未完全就绪
	// 当前路由已注册，但 API 调用可能会失败，因为服务为 nil
	sharedContainer := container.NewSharedServiceContainer()

	log.Println("警告: Shared 服务容器已创建，但服务尚未实现")
	log.Println("Shared API 路由已注册到: /api/v1/shared/")
	log.Println("  - /api/v1/shared/auth/*")
	log.Println("  - /api/v1/shared/wallet/*")
	log.Println("  - /api/v1/shared/storage/*")
	log.Println("  - /api/v1/shared/admin/*")

	// 注册 shared API 路由组
	sharedGroup := v1.Group("/shared")
	sharedRouter.RegisterRoutes(sharedGroup, sharedContainer)

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

	// 注册AI路由 (临时禁用)
	// aiSvc := aiService.NewService()
	// aiRouter := ai.NewAIRouter(aiSvc)
	// aiRouter.InitAIRouter(v1)

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
