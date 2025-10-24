package router

import (
	"log"
	"time"

	adminRouter "Qingyu_backend/router/admin"
	aiRouter "Qingyu_backend/router/ai"
	bookstoreRouter "Qingyu_backend/router/bookstore"
	projectRouter "Qingyu_backend/router/project"
	readerRouter "Qingyu_backend/router/reader"
	sharedRouter "Qingyu_backend/router/shared"
	userRouter "Qingyu_backend/router/user"

	aiService "Qingyu_backend/service/ai"
	bookstoreService "Qingyu_backend/service/bookstore"
	readingService "Qingyu_backend/service/reading"
	"Qingyu_backend/service/shared/container"
	userService "Qingyu_backend/service/user"

	"Qingyu_backend/config"
	"Qingyu_backend/global"
	"Qingyu_backend/repository/mongodb"
	mongoAI "Qingyu_backend/repository/mongodb/ai"
	mongoBookstore "Qingyu_backend/repository/mongodb/bookstore"
	mongoUser "Qingyu_backend/repository/mongodb/user"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本组
	v1 := r.Group("/api/v1")

	// 注册共享服务路由（认证、钱包、存储）
	// 注意：SharedServiceContainer 的服务实现尚未完全就绪
	// 当前路由已注册，但 API 调用可能会失败，因为服务为 nil
	sharedContainer := container.NewSharedServiceContainer()

	log.Println("警告: Shared 服务容器已创建，但服务尚未完全实现")
	log.Println("Shared API 路由已注册到: /api/v1/shared/")
	log.Println("  - /api/v1/shared/auth/* (认证相关)")
	log.Println("  - /api/v1/shared/wallet/* (钱包相关)")
	log.Println("  - /api/v1/shared/storage/* (存储相关)")

	// 注册 shared API 路由组
	sharedGroup := v1.Group("/shared")
	sharedRouter.RegisterRoutes(sharedGroup, sharedContainer)

	// 注册书店路由
	// 初始化Bookstore Repositories
	dbName := config.GlobalConfig.Database.Primary.MongoDB.Database
	bookRepo := mongoBookstore.NewMongoBookRepository(global.MongoClient, dbName)
	categoryRepo := mongoBookstore.NewMongoCategoryRepository(global.MongoClient, dbName)
	bannerRepo := mongoBookstore.NewMongoBannerRepository(global.MongoClient, dbName)
	rankingRepo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, dbName)

	// 创建Bookstore Services
	bookstoreSvc := bookstoreService.NewBookstoreService(bookRepo, categoryRepo, bannerRepo, rankingRepo)
	// TODO: 初始化其他书店服务
	// bookDetailSvc := bookstoreService.NewBookDetailService(...)
	// ratingSvc := bookstoreService.NewRatingService(...)
	// statisticsSvc := bookstoreService.NewStatisticsService(...)

	// 注册书店路由
	bookstoreRouter.InitBookstoreRouter(v1, bookstoreSvc, nil, nil, nil)

	log.Println("书店路由已注册到: /api/v1/bookstore/")
	log.Println("  - /api/v1/bookstore/homepage (书城首页)")
	log.Println("  - /api/v1/bookstore/books/* (书籍列表、搜索、详情)")
	log.Println("  - /api/v1/bookstore/categories/* (分类)")
	log.Println("  - /api/v1/bookstore/rankings/* (排行榜)")

	// 注册阅读器路由
	// 创建Repository工厂
	mongoConfig := &config.MongoDBConfig{
		URI:            config.GlobalConfig.Database.Primary.MongoDB.URI,
		Database:       config.GlobalConfig.Database.Primary.MongoDB.Database,
		MaxPoolSize:    config.GlobalConfig.Database.Primary.MongoDB.MaxPoolSize,
		MinPoolSize:    config.GlobalConfig.Database.Primary.MongoDB.MinPoolSize,
		ConnectTimeout: 10 * time.Second,
		ServerTimeout:  30 * time.Second,
	}

	repoFactory, err := mongodb.NewMongoRepositoryFactory(mongoConfig)
	if err != nil {
		log.Printf("警告: 创建Repository工厂失败: %v", err)
		log.Println("阅读器路由未注册")
	} else {
		// 创建Reader相关的Repository
		chapterRepo := repoFactory.CreateChapterRepository()
		progressRepo := repoFactory.CreateReadingProgressRepository()
		annotationRepo := repoFactory.CreateAnnotationRepository()
		settingsRepo := repoFactory.CreateReadingSettingsRepository()

		// 创建ReaderService（暂不使用缓存和VIP服务）
		readerSvc := readingService.NewReaderService(
			chapterRepo,
			progressRepo,
			annotationRepo,
			settingsRepo,
			nil, // eventBus - TODO: 实现事件总线
			nil, // cacheService - TODO: 实现缓存服务
			nil, // vipService - TODO: 实现VIP服务
		)

		// 注册阅读器路由
		readerRouter.InitReaderRouter(v1, readerSvc)

		log.Println("阅读器路由已注册到: /api/v1/reader/")
		log.Println("  - /api/v1/reader/books/* (书架管理)")
		log.Println("  - /api/v1/reader/chapters/* (章节内容)")
		log.Println("  - /api/v1/reader/progress/* (阅读进度)")
		log.Println("  - /api/v1/reader/annotations/* (标注管理)")
		log.Println("  - /api/v1/reader/settings/* (阅读设置)")
	}

	// 注册用户路由（用户认证、个人信息管理）
	// 初始化UserRepository和UserService
	userRepo := mongoUser.NewMongoUserRepository(global.DB)
	userSvc := userService.NewUserService(userRepo)
	userRouter.RegisterUserRoutes(v1, userSvc)

	log.Println("用户路由已注册到: /api/v1/")
	log.Println("  - /api/v1/register (用户注册)")
	log.Println("  - /api/v1/login (用户登录)")
	log.Println("  - /api/v1/users/profile (个人信息)")
	log.Println("  - /api/v1/users/password (修改密码)")

	// 注册文档路由
	projectRouter.RegisterRoutes(v1)

	// 注册AI路由
	aiSvc := aiService.NewService()

	// 创建AI相关Repository
	quotaRepo := mongoAI.NewMongoQuotaRepository(global.DB)

	// 创建聊天Repository（使用临时实现）
	chatRepo := aiService.NewInMemoryChatRepository()

	// 创建AI服务
	quotaService := aiService.NewQuotaService(quotaRepo)
	chatService := aiService.NewChatService(aiSvc, chatRepo)

	// 注册AI路由
	aiRouter.InitAIRouter(v1, aiSvc, chatService, quotaService)

	log.Println("AI服务路由已注册到: /api/v1/ai/")
	log.Println("  - /api/v1/ai/writing/* (续写、改写)")
	log.Println("  - /api/v1/ai/chat/* (聊天)")
	log.Println("  - /api/v1/ai/quota/* (配额管理)")

	// 注册管理员路由
	// TODO: 获取审核服务实例（需要实现）
	// auditSvc := ...
	adminSvc := sharedContainer.AdminService()
	adminRouter.RegisterAdminRoutes(v1, userSvc, quotaService, nil, adminSvc)

	log.Println("管理员路由已注册到: /api/v1/admin/")
	log.Println("  - /api/v1/admin/users/* (用户管理)")
	log.Println("  - /api/v1/admin/quota/* (AI配额管理)")
	log.Println("  - /api/v1/admin/audit/* (审核管理)")
	log.Println("  - /api/v1/admin/stats (系统统计)")
	log.Println("  - /api/v1/admin/config (系统配置)")

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
