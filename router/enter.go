package router

import (
	"fmt"
	"os"

	adminRouter "Qingyu_backend/router/admin"
	aiRouter "Qingyu_backend/router/ai"
	bookstoreRouter "Qingyu_backend/router/bookstore"
	messagingRouter "Qingyu_backend/router/messaging"
	projectRouter "Qingyu_backend/router/project"
	readerRouter "Qingyu_backend/router/reader"
	recommendationRouter "Qingyu_backend/router/recommendation"
	sharedRouter "Qingyu_backend/router/shared"
	socialRouter "Qingyu_backend/router/social"
	systemRouter "Qingyu_backend/router/system"
	userRouter "Qingyu_backend/router/user"
	writerRouter "Qingyu_backend/router/writer"

	"Qingyu_backend/service"
	sharedService "Qingyu_backend/service/shared"

	socialApi "Qingyu_backend/api/v1/social"
	recommendationAPI "Qingyu_backend/api/v1/recommendation"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 初始化zap日志器
	logger := initRouterLogger()

	// API版本组
	v1 := r.Group("/api/v1")

	// 获取全局服务容器
	serviceContainer := service.GetServiceContainer()
	if serviceContainer == nil {
		logger.Fatal("服务容器未初始化")
	}

	logger.Info("✓ 服务容器已初始化，开始注册路由...")

	// ============ 注册共享服务路由（渐进式注册） ============
	// 尝试从服务容器获取共享服务
	authSvc, authErr := serviceContainer.GetAuthService()
	walletSvc, walletErr := serviceContainer.GetWalletService()

	// 获取存储相关服务
	storageServiceImpl, storageErr := serviceContainer.GetStorageServiceImpl()
	multipartService, multipartErr := serviceContainer.GetMultipartUploadService()
	imageProcessor, imageErr := serviceContainer.GetImageProcessor()

	// 创建共享路由组（即使部分服务不可用也创建）
	sharedGroup := v1.Group("/shared")

	// 按可用服务逐个注册（渐进式注册策略）
	registeredCount := 0

	// 1. 注册认证服务路由
	if authErr == nil && authSvc != nil {
		sharedRouter.RegisterAuthRoutes(sharedGroup, authSvc)
		logger.Info("✓ 认证服务路由已注册: /api/v1/shared/auth/*")
		registeredCount++
	} else {
		logger.Warn("⚠ AuthService未配置，跳过认证路由注册", zap.Error(authErr))
	}

	// 2. 注册钱包服务路由
	if walletErr == nil && walletSvc != nil {
		sharedRouter.RegisterWalletRoutes(sharedGroup, walletSvc)
		logger.Info("✓ 钱包服务路由已注册: /api/v1/shared/wallet/*")
		registeredCount++
	} else {
		logger.Warn("⚠ WalletService未配置，跳过钱包路由注册", zap.Error(walletErr))
	}

	// 3. 注册存储服务路由
	if storageErr == nil && storageServiceImpl != nil && multipartErr == nil && multipartService != nil && imageErr == nil && imageProcessor != nil {
		sharedRouter.RegisterStorageRoutes(sharedGroup, storageServiceImpl, multipartService, imageProcessor)
		logger.Info("✓ 存储服务路由已注册: /api/v1/shared/storage/*")
		registeredCount++
	} else {
		logger.Warn("⚠ 存储服务未完全配置，跳过存储路由注册")
		if storageErr != nil {
			logger.Warn("  - StorageService", zap.Error(storageErr))
		}
		if multipartErr != nil {
			logger.Warn("  - MultipartUploadService", zap.Error(multipartErr))
		}
		if imageErr != nil {
			logger.Warn("  - ImageProcessor", zap.Error(imageErr))
		}
	}

	// 总结注册情况
	if registeredCount > 0 {
		logger.Info(fmt.Sprintf("✓ 已注册 %d 个共享服务模块到 /api/v1/shared/", registeredCount))
	} else {
		logger.Warn("⚠ 所有共享服务均未配置，共享路由组为空")
	}

	// ============ 注册书店路由 ============
	bookstoreSvc, err := serviceContainer.GetBookstoreService()
	if err != nil {
		logger.Warn("获取书店服务失败", zap.Error(err))
		logger.Info("书店路由未注册")
	} else {
		// 初始化其他书店服务
		bookDetailSvc, _ := serviceContainer.GetBookDetailService()
		ratingSvc, _ := serviceContainer.GetBookRatingService()
		statisticsSvc, _ := serviceContainer.GetBookStatisticsService()

		bookstoreRouter.InitBookstoreRouter(v1, bookstoreSvc, bookDetailSvc, ratingSvc, statisticsSvc)

		logger.Info("✓ 书店路由已注册到: /api/v1/bookstore/")
		logger.Info("  - /api/v1/bookstore/homepage (书城首页)")
		logger.Info("  - /api/v1/bookstore/books/* (书籍列表、搜索、详情)")
		logger.Info("  - /api/v1/bookstore/categories/* (分类)")
		logger.Info("  - /api/v1/bookstore/rankings/* (排行榜)")
		if bookDetailSvc != nil {
			logger.Info("  - /api/v1/bookstore/book-details/* (书籍详情)")
		}
		if ratingSvc != nil {
			logger.Info("  - /api/v1/bookstore/ratings/* (书籍评分)")
		}
		if statisticsSvc != nil {
			logger.Info("  - /api/v1/bookstore/statistics/* (书籍统计)")
		}
	}

	// ============ 注册阅读器路由 ============
	readerSvc, err := serviceContainer.GetReaderService()
	if err != nil {
		logger.Warn("获取阅读器服务失败", zap.Error(err))
		logger.Info("阅读器路由未注册")
	} else {
		// 获取评论服务（如果可用）
		commentSvc, commentErr := serviceContainer.GetCommentService()
		if commentErr != nil {
			logger.Warn("评论服务未配置", zap.Error(commentErr))
			commentSvc = nil
		}

		// 获取点赞服务（如果可用）
		likeSvc, likeErr := serviceContainer.GetLikeService()
		if likeErr != nil {
			logger.Warn("点赞服务未配置", zap.Error(likeErr))
			likeSvc = nil
		}

		// 获取收藏服务（如果可用）
		collectionSvc, collectionErr := serviceContainer.GetCollectionService()
		if collectionErr != nil {
			logger.Warn("收藏服务未配置", zap.Error(collectionErr))
			collectionSvc = nil
		}

		// 获取阅读历史服务（如果可用）
		readingHistorySvc, historyErr := serviceContainer.GetReadingHistoryService()
		if historyErr != nil {
			logger.Warn("阅读历史服务未配置", zap.Error(historyErr))
			readingHistorySvc = nil
		}

		readerRouter.InitReaderRouter(v1, readerSvc, commentSvc, likeSvc, collectionSvc, readingHistorySvc)

		logger.Info("✓ 阅读器路由已注册到: /api/v1/reader/")
		logger.Info("  - /api/v1/reader/books/* (书架管理)")
		logger.Info("  - /api/v1/reader/chapters/* (章节内容)")
		logger.Info("  - /api/v1/reader/progress/* (阅读进度)")
		logger.Info("  - /api/v1/reader/annotations/* (标注管理)")
		logger.Info("  - /api/v1/reader/settings/* (阅读设置)")
		if commentSvc != nil {
			logger.Info("  - /api/v1/reader/comments/* (评论系统)")
		}
	}

	// ============ 注册社交路由（新统一入口） ============
	// 获取评论服务
	commentSvc, commentErr := serviceContainer.GetCommentService()
	var commentAPI *socialApi.CommentAPI
	if commentErr != nil {
		logger.Warn("评论服务未配置", zap.Error(commentErr))
	} else {
		commentAPI = socialApi.NewCommentAPI(commentSvc)
	}

	// 获取点赞服务
	likeSvc, likeErr := serviceContainer.GetLikeService()
	var likeAPI *socialApi.LikeAPI
	if likeErr != nil {
		logger.Warn("点赞服务未配置", zap.Error(likeErr))
	} else {
		likeAPI = socialApi.NewLikeAPI(likeSvc)
	}

	// 获取收藏服务
	collectionSvc, collectionErr := serviceContainer.GetCollectionService()
	var collectionAPI *socialApi.CollectionAPI
	if collectionErr != nil {
		logger.Warn("收藏服务未配置", zap.Error(collectionErr))
	} else {
		collectionAPI = socialApi.NewCollectionAPI(collectionSvc)
	}

	// 尝试获取用户关系服务（新功能）
	var relationAPI *socialApi.UserRelationAPI
	// TODO: 添加 UserRelationService 到服务容器后获取
	// relationSvc, relationErr := serviceContainer.GetUserRelationService()
	// if relationErr == nil && relationSvc != nil {
	//     relationAPI = socialApi.NewUserRelationAPI(relationSvc)
	// }

	// 注册统一社交路由
	if commentAPI != nil || likeAPI != nil || collectionAPI != nil || relationAPI != nil {
		socialRouter.RegisterSocialRoutes(v1, relationAPI, commentAPI, likeAPI, collectionAPI)

		logger.Info("✓ 社交路由已注册到: /api/v1/social/")
		if commentAPI != nil {
			logger.Info("  - /api/v1/social/comments/* (评论系统)")
		}
		if likeAPI != nil {
			logger.Info("  - /api/v1/social/books/:bookId/like (点赞系统)")
		}
		if collectionAPI != nil {
			logger.Info("  - /api/v1/social/collections/* (收藏系统)")
		}
		if relationAPI != nil {
			logger.Info("  - /api/v1/social/follow/* (关注系统)")
		}
	}

	// ============ 注册推荐系统路由 ============
	recommendationSvc, err := serviceContainer.GetRecommendationService()
	if err != nil {
		logger.Warn("获取推荐服务失败", zap.Error(err))
		logger.Info("推荐系统路由未注册")
	} else {
		// 创建推荐API
		recommendationApi := recommendationAPI.NewRecommendationAPI(recommendationSvc)
		recommendationRouter.RegisterRecommendationRoutes(v1, recommendationApi)

		logger.Info("✓ 推荐系统路由已注册到: /api/v1/recommendation/")
		logger.Info("  - /api/v1/recommendation/personalized (个性化推荐 - 需认证)")
		logger.Info("  - /api/v1/recommendation/similar (相似推荐)")
		logger.Info("  - /api/v1/recommendation/behavior (记录行为 - 需认证)")
		logger.Info("  - /api/v1/recommendation/homepage (首页推荐)")
		logger.Info("  - /api/v1/recommendation/hot (热门推荐)")
		logger.Info("  - /api/v1/recommendation/category (分类推荐)")
	}

	// ============ 注册Messaging路由 ============
	announcementSvc, announcementErr := serviceContainer.GetAnnouncementService()
	if announcementErr != nil {
		logger.Warn("获取公告服务失败", zap.Error(announcementErr))
		logger.Info("Messaging路由未注册")
	} else {
		messagingRouter.RegisterRoutes(v1, announcementSvc)
		logger.Info("✓ Messaging路由已注册到: /api/v1/announcements/")
		logger.Info("  - GET /api/v1/announcements/effective (获取有效公告)")
		logger.Info("  - GET /api/v1/announcements/:id (获取公告详情)")
		logger.Info("  - POST /api/v1/announcements/:id/view (增加查看次数)")
	}

	// ============ 注册用户路由 ============
	userSvc, err := serviceContainer.GetUserService()
	if err != nil {
		logger.Fatal("获取用户服务失败", zap.Error(err))
	}

	// 获取书店服务（可选）
	bookstoreSvc, bookstoreErr := serviceContainer.GetBookstoreService()
	if bookstoreErr != nil {
		logger.Warn("获取书店服务失败，用户作品列表功能将不可用", zap.Error(bookstoreErr))
		userRouter.RegisterUserRoutes(v1, userSvc)
	} else {
		userRouter.RegisterUserRoutesWithBookstore(v1, userSvc, bookstoreSvc)
	}

	logger.Info("✓ 用户路由已注册到: /api/v1/")
	logger.Info("  - /api/v1/register (用户注册)")
	logger.Info("  - /api/v1/login (用户登录)")
	logger.Info("  - /api/v1/users/profile (个人信息)")
	logger.Info("  - /api/v1/users/password (修改密码)")
	if bookstoreErr == nil {
		logger.Info("  - /api/v1/users/:userId/books (用户作品列表)")
	}

	// ============ 注册文档路由 ============
	projectRouter.RegisterRoutes(v1)
	logger.Info("✓ 文档路由已注册到: /api/v1/projects/")

	// ============ 注册写作端路由 ============
	writerRouter.RegisterWriterRoutes(v1)
	logger.Info("✓ 写作端路由已注册到: /api/v1/writer/")
	logger.Info("  - /api/v1/writer/projects/* (项目管理)")
	logger.Info("  - /api/v1/writer/documents/* (文档管理)")
	logger.Info("  - /api/v1/writer/versions/* (版本控制)")

	// ============ 注册AI路由 ============
	aiSvc, err := serviceContainer.GetAIService()
	if err != nil {
		logger.Warn("获取AI服务失败", zap.Error(err))
		logger.Info("AI路由未注册")
	} else {
		chatService, err := serviceContainer.GetChatService()
		if err != nil {
			logger.Warn("获取聊天服务失败", zap.Error(err))
			chatService = nil
		}

		quotaService, err := serviceContainer.GetQuotaService()
		if err != nil {
			logger.Warn("获取配额服务失败", zap.Error(err))
			quotaService = nil
		}

		// 尝试连接Phase3 gRPC服务
		phase3Client, err := serviceContainer.GetPhase3Client()
		if err != nil {
			logger.Warn("Phase3 gRPC客户端未配置", zap.Error(err))
			phase3Client = nil
		}

		aiRouter.InitAIRouter(v1, aiSvc, chatService, quotaService, phase3Client)

		logger.Info("✓ AI服务路由已注册到: /api/v1/ai/")
		logger.Info("  - /api/v1/ai/writing/* (续写、改写)")
		logger.Info("  - /api/v1/ai/chat/* (聊天)")
		logger.Info("  - /api/v1/ai/quota/* (配额管理)")
		if phase3Client != nil {
			logger.Info("  - /api/v1/ai/creative/* (Phase3创作工作流)")
		}
	}

	// ============ 注册管理员路由 ============
	// 获取配额服务（用于管理员管理）
	quotaService, _ := serviceContainer.GetQuotaService()

	// 获取 AdminService（如果可用）
	adminSvc, adminErr := serviceContainer.GetAdminService()
	if adminErr != nil {
		logger.Warn("⚠ AdminService未配置", zap.Error(adminErr))
		adminSvc = nil
	}

	// 创建配置管理服务
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "./config/config.yaml"
	}
	configSvc := sharedService.NewConfigService(configPath)

	// ✅ 获取审核服务
	auditSvc, auditErr := serviceContainer.GetAuditService()
	if auditErr != nil {
		logger.Warn("⚠ AuditService未配置", zap.Error(auditErr))
		auditSvc = nil
	}

	// 获取公告服务（用于管理员）
	announcementSvc, announcementSvcErr := serviceContainer.GetAnnouncementService()
	if announcementSvcErr != nil {
		logger.Warn("⚠ AnnouncementService未配置", zap.Error(announcementSvcErr))
		announcementSvc = nil
	}

	adminRouter.RegisterAdminRoutes(v1, userSvc, quotaService, auditSvc, adminSvc, configSvc, announcementSvc)

	logger.Info("✓ 管理员路由已注册到: /api/v1/admin/")
	logger.Info("  - /api/v1/admin/users/* (用户管理)")
	logger.Info("  - /api/v1/admin/quota/* (AI配额管理)")
	logger.Info("  - /api/v1/admin/audit/* (审核管理)")
	logger.Info("  - /api/v1/admin/stats (系统统计)")
	logger.Info("  - /api/v1/admin/config/* (配置管理)")

	// ============ 注册系统监控路由 ============
	systemRouter.InitSystemRoutes(v1)
	logger.Info("✓ 系统监控路由已注册到: /api/v1/system/")
	logger.Info("  - /api/v1/system/health (系统健康检查)")
	logger.Info("  - /api/v1/system/health/:service (服务健康检查)")
	logger.Info("  - /api/v1/system/metrics (所有服务指标)")
	logger.Info("  - /api/v1/system/metrics/:service (特定服务指标)")

	// ============ 注册Prometheus metrics端点 ============
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	logger.Info("✓ Prometheus metrics端点已注册: /metrics")

	// ============ 健康检查 ============
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	logger.Info("\n========================================")
	logger.Info("✓ 所有路由注册完成!")
	logger.Info("==========================================")
}

// initRouterLogger 初始化路由日志器
func initRouterLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色输出

	logger, err := cfg.Build()
	if err != nil {
		panic("Failed to initialize router logger: " + err.Error())
	}

	return logger
}
