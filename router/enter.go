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
	serviceInterfaces "Qingyu_backend/service/interfaces"
	"Qingyu_backend/service/shared/auth"
	"Qingyu_backend/service/shared/container"
	userService "Qingyu_backend/service/user"

	"Qingyu_backend/config"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	"Qingyu_backend/repository/mongodb"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本组
	v1 := r.Group("/api/v1")

	// 注册共享服务路由（认证、钱包、存储、管理）
	sharedContainer := container.NewSharedServiceContainer()

	log.Println("Shared API 路由已注册到: /api/v1/shared/")
	log.Println("  - /api/v1/shared/auth/*")
	log.Println("  - /api/v1/shared/wallet/*")
	log.Println("  - /api/v1/shared/storage/*")
	log.Println("  - /api/v1/shared/admin/*")

	// 初始化Repository工厂
	mongoConfig := config.GlobalConfig.Database.Primary.MongoDB
	repoFactory, err := mongodb.NewMongoRepositoryFactory(mongoConfig)
	if err != nil {
		log.Printf("警告: 创建Repository工厂失败: %v", err)
		log.Println("书城服务将使用 nil repositories，部分功能可能不可用")
	}

	// 创建Repository实例（使用正确的接口类型）
	var bookRepo bookstoreRepo.BookRepository
	var categoryRepo bookstoreRepo.CategoryRepository
	var bannerRepo bookstoreRepo.BannerRepository
	var rankingRepo bookstoreRepo.RankingRepository

	if repoFactory != nil {
		bookRepo = repoFactory.CreateBookRepository()
		categoryRepo = repoFactory.CreateCategoryRepository()
		bannerRepo = repoFactory.CreateBannerRepository()
		// rankingRepo 保持为 nil，因为暂未实现
		log.Println("书城 Repository 已初始化")
	}

	// 注册书城路由
	bookstoreSvc := bookstoreService.NewBookstoreService(
		bookRepo,
		categoryRepo,
		bannerRepo,
		rankingRepo) // RankingRepository 暂未实现，传入nil
	bookstoreAPI := readingAPI.NewBookstoreAPI(bookstoreSvc)
	bookstoreRouterInstance := bookstoreRouter.NewBookstoreRouter(bookstoreAPI)
	// 只注册主路由（已包含所有功能）
	bookstoreRouterInstance.RegisterRoutes(v1)

	// 初始化用户服务
	var userSvc serviceInterfaces.UserService
	if repoFactory != nil {
		userRepo := repoFactory.CreateUserRepository()
		userSvc = userService.NewUserService(userRepo)
		log.Println("用户服务已初始化")
	} else {
		log.Println("警告: 用户服务未初始化，用户认证功能将不可用")
	}

	// 初始化Shared认证服务
	if repoFactory != nil && userSvc != nil {
		// 创建AuthRepository
		authRepo := repoFactory.CreateAuthRepository()

		// 创建JWT配置
		jwtConfig := config.GetJWTConfigEnhanced()

		// 创建JWTService（暂时不使用Redis）
		jwtService := auth.NewJWTService(jwtConfig, nil)

		// 创建RoleService
		roleService := auth.NewRoleService(authRepo)

		// 创建PermissionService（暂时不使用Cache）
		permissionService := auth.NewPermissionService(authRepo, nil)

		// 创建AuthService
		authService := auth.NewAuthService(jwtService, roleService, permissionService, authRepo, userSvc)

		// 设置到SharedServiceContainer
		sharedContainer.SetAuthService(authService)

		log.Println("Shared认证服务已初始化")
	} else {
		log.Println("警告: Shared认证服务未初始化")
	}

	// 注册 shared API 路由组
	sharedGroup := v1.Group("/shared")
	sharedRouter.RegisterRoutes(sharedGroup, sharedContainer)

	// 注册系统路由（用户认证等）
	userRouter.RegisterUserRoutes(v1, userSvc)

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
