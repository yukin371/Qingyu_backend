package router

import (
	"log"

	adminRouter "Qingyu_backend/router/admin"
	aiRouter "Qingyu_backend/router/ai"
	bookstoreRouter "Qingyu_backend/router/bookstore"
	projectRouter "Qingyu_backend/router/project"
	readerRouter "Qingyu_backend/router/reader"
	sharedRouter "Qingyu_backend/router/shared"
	systemRouter "Qingyu_backend/router/system"
	userRouter "Qingyu_backend/router/user"

	"Qingyu_backend/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API版本组
	v1 := r.Group("/api/v1")

	// 获取全局服务容器
	serviceContainer := service.GetServiceContainer()
	if serviceContainer == nil {
		log.Fatal("服务容器未初始化")
	}

	log.Println("✓ 服务容器已初始化，开始注册路由...")

	// ============ 注册共享服务路由（如果已配置） ============
	// 尝试从服务容器获取共享服务
	authSvc, authErr := serviceContainer.GetAuthService()
	walletSvc, walletErr := serviceContainer.GetWalletService()
	storageSvc, storageErr := serviceContainer.GetStorageService()

	// 只有当所有共享服务都可用时，才注册共享服务路由
	if authErr == nil && walletErr == nil && storageErr == nil {
		sharedGroup := v1.Group("/shared")
		sharedRouter.RegisterRoutes(sharedGroup, authSvc, walletSvc, storageSvc)
		log.Println("✓ 共享服务路由已注册到: /api/v1/shared/")
		log.Println("  - /api/v1/shared/auth/* (认证服务)")
		log.Println("  - /api/v1/shared/wallet/* (钱包服务)")
		log.Println("  - /api/v1/shared/storage/* (存储服务)")
	} else {
		log.Println("⚠ 共享服务路由未注册（服务未配置）")
		if authErr != nil {
			log.Printf("  - AuthService: %v", authErr)
		}
		if walletErr != nil {
			log.Printf("  - WalletService: %v", walletErr)
		}
		if storageErr != nil {
			log.Printf("  - StorageService: %v", storageErr)
		}
	}

	// ============ 注册书店路由 ============
	bookstoreSvc, err := serviceContainer.GetBookstoreService()
	if err != nil {
		log.Printf("获取书店服务失败: %v", err)
		log.Println("书店路由未注册")
	} else {
		// TODO: 初始化其他书店服务
		// bookDetailSvc := serviceContainer.GetBookDetailService()
		// ratingSvc := serviceContainer.GetRatingService()
		// statisticsSvc := serviceContainer.GetStatisticsService()

		bookstoreRouter.InitBookstoreRouter(v1, bookstoreSvc, nil, nil, nil)

		log.Println("✓ 书店路由已注册到: /api/v1/bookstore/")
		log.Println("  - /api/v1/bookstore/homepage (书城首页)")
		log.Println("  - /api/v1/bookstore/books/* (书籍列表、搜索、详情)")
		log.Println("  - /api/v1/bookstore/categories/* (分类)")
		log.Println("  - /api/v1/bookstore/rankings/* (排行榜)")
	}

	// ============ 注册阅读器路由 ============
	readerSvc, err := serviceContainer.GetReaderService()
	if err != nil {
		log.Printf("获取阅读器服务失败: %v", err)
		log.Println("阅读器路由未注册")
	} else {
		readerRouter.InitReaderRouter(v1, readerSvc)

		log.Println("✓ 阅读器路由已注册到: /api/v1/reader/")
		log.Println("  - /api/v1/reader/books/* (书架管理)")
		log.Println("  - /api/v1/reader/chapters/* (章节内容)")
		log.Println("  - /api/v1/reader/progress/* (阅读进度)")
		log.Println("  - /api/v1/reader/annotations/* (标注管理)")
		log.Println("  - /api/v1/reader/settings/* (阅读设置)")
	}

	// ============ 注册用户路由 ============
	userSvc, err := serviceContainer.GetUserService()
	if err != nil {
		log.Fatalf("获取用户服务失败: %v", err)
	}

	userRouter.RegisterUserRoutes(v1, userSvc)

	log.Println("✓ 用户路由已注册到: /api/v1/")
	log.Println("  - /api/v1/register (用户注册)")
	log.Println("  - /api/v1/login (用户登录)")
	log.Println("  - /api/v1/users/profile (个人信息)")
	log.Println("  - /api/v1/users/password (修改密码)")

	// ============ 注册文档路由 ============
	projectRouter.RegisterRoutes(v1)
	log.Println("✓ 文档路由已注册到: /api/v1/projects/")

	// ============ 注册AI路由 ============
	aiSvc, err := serviceContainer.GetAIService()
	if err != nil {
		log.Printf("获取AI服务失败: %v", err)
		log.Println("AI路由未注册")
	} else {
		chatService, err := serviceContainer.GetChatService()
		if err != nil {
			log.Printf("获取聊天服务失败: %v", err)
			chatService = nil
		}

		quotaService, err := serviceContainer.GetQuotaService()
		if err != nil {
			log.Printf("获取配额服务失败: %v", err)
			quotaService = nil
		}

		aiRouter.InitAIRouter(v1, aiSvc, chatService, quotaService)

		log.Println("✓ AI服务路由已注册到: /api/v1/ai/")
		log.Println("  - /api/v1/ai/writing/* (续写、改写)")
		log.Println("  - /api/v1/ai/chat/* (聊天)")
		log.Println("  - /api/v1/ai/quota/* (配额管理)")
	}

	// ============ 注册管理员路由 ============
	// 获取配额服务（用于管理员管理）
	quotaService, _ := serviceContainer.GetQuotaService()

	// 获取 AdminService（如果可用）
	adminSvc, adminErr := serviceContainer.GetAdminService()
	if adminErr != nil {
		log.Printf("⚠ AdminService未配置: %v", adminErr)
		adminSvc = nil
	}

	// TODO: 获取审核服务实例（需要实现）
	// auditSvc := serviceContainer.GetAuditService()
	adminRouter.RegisterAdminRoutes(v1, userSvc, quotaService, nil, adminSvc)

	log.Println("✓ 管理员路由已注册到: /api/v1/admin/")
	log.Println("  - /api/v1/admin/users/* (用户管理)")
	log.Println("  - /api/v1/admin/quota/* (AI配额管理)")
	log.Println("  - /api/v1/admin/audit/* (审核管理)")
	log.Println("  - /api/v1/admin/stats (系统统计)")
	log.Println("  - /api/v1/admin/config (系统配置)")

	// ============ 注册系统监控路由 ============
	systemRouter.InitSystemRoutes(v1)
	log.Println("✓ 系统监控路由已注册到: /api/v1/system/")
	log.Println("  - /api/v1/system/health (系统健康检查)")
	log.Println("  - /api/v1/system/health/:service (服务健康检查)")
	log.Println("  - /api/v1/system/metrics (所有服务指标)")
	log.Println("  - /api/v1/system/metrics/:service (特定服务指标)")

	// ============ 健康检查 ============
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	log.Println("\n========================================")
	log.Println("✓ 所有路由注册完成!")
	log.Println("==========================================")
}
