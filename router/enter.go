package router

import (
	"context"
	"fmt"
	"os"
	"time"

	adminRouter "Qingyu_backend/router/admin"
	aiRouter "Qingyu_backend/router/ai"
	announcementsRouter "Qingyu_backend/router/announcements"
	bookstoreRouter "Qingyu_backend/router/bookstore"
	financeRouter "Qingyu_backend/router/finance"
	notificationsRouter "Qingyu_backend/router/notifications"
	readerRouter "Qingyu_backend/router/reader"
	readingstatsRouter "Qingyu_backend/router/reading-stats"
	recommendationRouter "Qingyu_backend/router/recommendation"
	searchRouter "Qingyu_backend/api/v1/search"
	sharedRouter "Qingyu_backend/router/shared"
	socialRouter "Qingyu_backend/router/social"
	systemRouter "Qingyu_backend/router/system"
	userRouter "Qingyu_backend/router/user"
	writerRouter "Qingyu_backend/router/writer"

	adminrep "Qingyu_backend/repository/mongodb/admin"
	authRep "Qingyu_backend/repository/mongodb/auth"
	userRepo "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/service"
	"Qingyu_backend/service/container"
	adminservice "Qingyu_backend/service/admin"
	bookstore "Qingyu_backend/service/bookstore"
	sharedService "Qingyu_backend/service/shared"
	searchService "Qingyu_backend/service/search"
	searchengine "Qingyu_backend/service/search/engine"
	searchprovider "Qingyu_backend/service/search/provider"
	statsService "Qingyu_backend/service/shared/stats"

	financeApi "Qingyu_backend/api/v1/finance"
	messagesApi "Qingyu_backend/api/v1/messages"
	notificationsAPI "Qingyu_backend/api/v1/notifications"
	recommendationAPI "Qingyu_backend/api/v1/recommendation"
	socialApi "Qingyu_backend/api/v1/social"
	syncService "Qingyu_backend/pkg/sync"
	readerservice "Qingyu_backend/service/reader"
	messagingService "Qingyu_backend/service/messaging"
	modelsMessaging "Qingyu_backend/models/messaging"
	websocketHub "Qingyu_backend/realtime/websocket"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
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
	oauthSvc, oauthErr := serviceContainer.GetOAuthService()

	// 获取存储相关服务
	storageServiceImpl, storageErr := serviceContainer.GetStorageServiceImpl()
	multipartService, multipartErr := serviceContainer.GetMultipartUploadService()
	imageProcessor, imageErr := serviceContainer.GetImageProcessor()

	// 搜索服务（SearchService）将在writer路由初始化时创建
	// TODO: 将SearchService添加到服务容器，以便在shared路由中注册搜索建议功能

	// 创建共享路由组（即使部分服务不可用也创建）
	sharedGroup := v1.Group("/shared")

	// 按可用服务逐个注册（渐进式注册策略）
	registeredCount := 0

	// 1. 注册认证服务路由
	if authErr == nil && authSvc != nil {
		// OAuthService是可选的，如果没有配置，传入nil
		sharedRouter.RegisterAuthRoutes(sharedGroup, authSvc, oauthSvc, logger)
		logger.Info("✓ 认证服务路由已注册: /api/v1/shared/auth/*")
		if oauthErr != nil {
			logger.Warn("⚠ OAuthService未配置，OAuth登录功能将不可用", zap.Error(oauthErr))
		}
		registeredCount++
	} else {
		logger.Warn("⚠ AuthService未配置，跳过认证路由注册", zap.Error(authErr))
	}

	// 2. 注册存储服务路由
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

	// ============ 注册财务路由 ============
	walletSvc, walletErr := serviceContainer.GetWalletService()
	if walletErr != nil {
		logger.Warn("获取钱包服务失败", zap.Error(walletErr))
		logger.Info("财务路由未注册")
	} else {
		walletAPI := financeApi.NewWalletAPI(walletSvc)

		// 获取会员服务
		membershipSvc, membershipErr := serviceContainer.GetMembershipService()
		var membershipAPI *financeApi.MembershipAPI
		if membershipErr != nil {
			logger.Warn("获取会员服务失败", zap.Error(membershipErr))
		} else {
			membershipAPI = financeApi.NewMembershipAPI(membershipSvc)
		}

		// 获取作者收入服务
		authorRevenueSvc, revenueErr := serviceContainer.GetAuthorRevenueService()
		var authorRevenueAPI *financeApi.AuthorRevenueAPI
		if revenueErr != nil {
			logger.Warn("获取作者收入服务失败", zap.Error(revenueErr))
		} else {
			authorRevenueAPI = financeApi.NewAuthorRevenueAPI(authorRevenueSvc)
		}

		// 注册财务路由
		financeRouter.RegisterFinanceRoutes(v1, walletAPI, membershipAPI, authorRevenueAPI)
		logger.Info("✓ 财务路由已注册到: /api/v1/finance/")
		logger.Info("  - /api/v1/finance/wallet/* (钱包管理)")
		if membershipAPI != nil {
			logger.Info("  - /api/v1/finance/membership/* (会员系统)")
		}
		if authorRevenueAPI != nil {
			logger.Info("  - /api/v1/finance/author/* (作者收入)")
		}
		logger.Info("  - ⚠️  旧路由 /api/v1/shared/wallet/* 继续保留以向后兼容")
	}

	// ============ 初始化搜索服务（需要在书店路由之前）============
	// 创建 MongoEngine、BookProvider，并注册到 SearchService
	searchSvc := initSearchService(serviceContainer, logger)

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

		// 获取章节服务（如果可用）
		var chapterSvc bookstore.ChapterService
		if svc, err := serviceContainer.GetChapterService(); err == nil {
			chapterSvc = svc
		}

		// 获取章节购买服务（如果可用）
		var chapterPurchaseSvc bookstore.ChapterPurchaseService
		// TODO: 添加 GetChapterPurchaseService 到服务容器
		// if svc, err := serviceContainer.GetChapterPurchaseService(); err == nil {
		//     chapterPurchaseSvc = svc
		// }

		// 注册书店路由，传入搜索服务
		bookstoreRouter.InitBookstoreRouter(v1, bookstoreSvc, bookDetailSvc, ratingSvc, statisticsSvc, chapterSvc, chapterPurchaseSvc, searchSvc, logger)

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
		// 获取章节服务（用于章节阅读）
		var chapterSvc bookstore.ChapterService
		if svc, err := serviceContainer.GetChapterService(); err == nil {
			chapterSvc = svc
		} else {
			logger.Warn("章节服务未配置，章节阅读功能将不可用", zap.Error(err))
		}

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

		// 获取书签服务（如果可用）
		var bookmarkSvc readerservice.BookmarkService
		bookmarkSvc, bookmarkErr := serviceContainer.GetBookmarkService()
		if bookmarkErr != nil {
			logger.Warn("书签服务未配置", zap.Error(bookmarkErr))
			bookmarkSvc = nil
		} else {
			logger.Info("✓ 书签服务已配置")
		}

		// 进度同步服务（TODO: 需要添加到服务容器）
		var progressSyncSvc *syncService.ProgressSyncService = nil

		readerRouter.InitReaderRouter(v1, readerSvc, chapterSvc, commentSvc, likeSvc, collectionSvc, readingHistorySvc, progressSyncSvc, bookmarkSvc)

		logger.Info("✓ 阅读器路由已注册到: /api/v1/reader/")
		logger.Info("  - /api/v1/reader/books/* (书架管理)")
		if chapterSvc != nil {
			logger.Info("  - /api/v1/reader/books/:bookId/chapters/* (章节阅读)")
			logger.Info("    - GET /:chapterId (获取章节内容)")
			logger.Info("    - GET /:chapterId/next (下一章)")
			logger.Info("    - GET /:chapterId/previous (上一章)")
			logger.Info("    - GET /by-number/:chapterNum (按章节号获取)")
		}
		logger.Info("  - /api/v1/reader/progress/* (阅读进度)")
		//nolint:nilness // 待实现服务
		if progressSyncSvc != nil {
			logger.Info("  - /api/v1/reader/progress/ws (WebSocket同步)")
			logger.Info("  - /api/v1/reader/progress/sync (进度同步)")
			logger.Info("  - /api/v1/reader/progress/merge (合并离线进度)")
			logger.Info("  - /api/v1/reader/progress/sync-status (同步状态)")
		}
		logger.Info("  - /api/v1/reader/annotations/* (标注管理)")
		logger.Info("  - /api/v1/reader/settings/* (阅读设置)")
		if commentSvc != nil {
			logger.Info("  - /api/v1/reader/comments/* (评论系统)")
		}

		// ============ 注册阅读统计路由 ============
		readingStatsSvc, readingStatsErr := serviceContainer.GetReadingStatsService()
		if readingStatsErr != nil {
			logger.Warn("获取阅读统计服务失败", zap.Error(readingStatsErr))
			logger.Info("阅读统计路由未注册")
		} else {
			readingstatsRouter.RegisterReadingStatsRoutes(v1, readingStatsSvc)
			logger.Info("✓ 阅读统计路由已注册到: /api/v1/reading-stats/")
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

	// 关注服务
	followSvc, followErr := serviceContainer.GetFollowService()
	var followAPI *socialApi.FollowAPI
	if followErr == nil && followSvc != nil {
		followAPI = socialApi.NewFollowAPI(followSvc)
	}

	// 新增社交 API（待实现）
	var messageAPI *messagesApi.MessageAPI //nolint:ineffassign // 待实现
	var messageAPIV2 *socialApi.MessageAPIV2
	var reviewAPI *socialApi.ReviewAPI     //nolint:ineffassign // 待实现
	var booklistAPI *socialApi.BookListAPI //nolint:ineffassign // 待实现

	// 初始化MessageAPIV2（消息服务V2）
	messagingWSHub, wsHubErr := serviceContainer.GetMessagingWSHub()

	if wsHubErr == nil && messagingWSHub != nil {
		// 获取MongoDB数据库
		mongoDB := serviceContainer.GetMongoDB()
		repositoryFactory := serviceContainer.GetRepositoryFactory()
		if mongoDB != nil && repositoryFactory != nil {
			// 创建消息服务和会话服务（使用models/messaging中的Repository）
			messageRepo := repositoryFactory.CreateMessageRepository()
			conversationRepo := modelsMessaging.NewMongoConversationRepository(mongoDB)

			// 先创建ConversationService
			conversationSvc := messagingService.NewConversationService(conversationRepo, messageRepo)
			// 再创建MessageService（需要ConversationService）
			messageSvc := messagingService.NewMessageService(messageRepo, conversationSvc)

			messageAPIV2 = socialApi.NewMessageAPIV2(messageSvc, conversationSvc, messagingWSHub)
			logger.Info("✓ MessageAPIV2初始化完成")
		}
	} else {
		logger.Warn("MessagingWSHub未配置", zap.Error(wsHubErr))
	}

	// 注册WebSocket路由
	if messagingWSHub != nil {
		r.GET("/ws/messages", messagingWSHub.HandleMessagingWebSocket)
		logger.Info("✓ WebSocket路由已注册: /ws/messages")
	}

	// 注册统一社交路由
	if commentAPI != nil || likeAPI != nil || collectionAPI != nil { //nolint:nilness // 待实现服务已排除
		socialRouter.RegisterSocialRoutes(v1, relationAPI, commentAPI, likeAPI, collectionAPI, followAPI, messageAPI, messageAPIV2, reviewAPI, booklistAPI)

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
		//nolint:nilness // 待实现服务
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

	// ============ 注册Announcements路由 ============
	announcementSvc, announcementErr := serviceContainer.GetAnnouncementService()
	if announcementErr != nil {
		logger.Warn("获取公告服务失败", zap.Error(announcementErr))
		logger.Info("Announcements路由未注册")
	} else {
		announcementsRouter.RegisterAnnouncementRoutes(v1, announcementSvc)
		logger.Info("✓ Announcements路由已注册到: /api/v1/announcements/")
		logger.Info("  - GET /api/v1/announcements/effective (获取有效公告)")
		logger.Info("  - GET /api/v1/announcements/:id (获取公告详情)")
		logger.Info("  - POST /api/v1/announcements/:id/view (增加查看次数)")
	}

	// ============ 注册通知路由 ============
	notificationSvc, notificationErr := serviceContainer.GetNotificationService()
	if notificationErr != nil {
		logger.Warn("获取通知服务失败", zap.Error(notificationErr))
		logger.Info("通知路由未注册")
	} else {
		notificationAPI := notificationsAPI.NewNotificationAPI(notificationSvc)
		notificationsRouter.RegisterRoutes(v1, notificationAPI)
		notificationsRouter.RegisterUserManagementRoutes(v1, notificationAPI)
		logger.Info("✓ 通知路由已注册到: /api/v1/notifications/")
		logger.Info("  - GET /api/v1/notifications (获取通知列表)")
		logger.Info("  - GET /api/v1/notifications/:id (获取通知详情)")
		logger.Info("  - PUT /api/v1/notifications/:id/read (标记为已读)")
		logger.Info("  - PUT /api/v1/notifications/read-all (全部标记为已读)")
		logger.Info("  - DELETE /api/v1/notifications/:id (删除通知)")
		logger.Info("  - DELETE /api/v1/notifications/batch-delete (批量删除)")
		logger.Info("  - /api/v1/notifications/preferences (通知偏好设置)")
		logger.Info("  - /api/v1/notifications/push/* (推送设备管理)")
		logger.Info("  - /api/v1/user-management/email-notifications (邮件通知设置)")
		logger.Info("  - /api/v1/user-management/sms-notifications (短信通知设置)")
	}

	// ============ 注册写作端路由 ============
	writerRouter.RegisterWriterRoutes(v1)
	logger.Info("✓ 写作端路由已注册到: /api/v1/writer/")
	logger.Info("  - /api/v1/writer/projects/* (项目管理)")
	logger.Info("  - /api/v1/writer/documents/* (文档管理)")
	logger.Info("  - /api/v1/writer/versions/* (版本控制)")
	logger.Info("  - /api/v1/writer/search/documents (文档搜索)")
	logger.Info("  ⚠️  注意: 测试需要使用 /api/v1/writer/ 前缀而非 /api/v1/")

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

	// ============ 注册统一搜索路由 ============
	if searchSvc != nil {
		// 注册搜索路由
		searchRouter.RegisterSearchRoutes(v1, searchSvc)
		logger.Info("✓ 统一搜索路由已注册到: /api/v1/search/")
		logger.Info("  - /api/v1/search/search (统一搜索)")
		logger.Info("  - /api/v1/search/batch (批量搜索)")
		logger.Info("  - /api/v1/search/health (健康检查)")
	} else {
		logger.Warn("⚠ 搜索服务初始化失败，搜索路由未注册")
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

	// 获取用户服务（用于管理员和用户管理路由）
	userSvc, err := serviceContainer.GetUserService()
	if err != nil {
		logger.Fatal("获取用户服务失败", zap.Error(err))
	}

	// 创建用户管理服务（UserAdminService - 管理员专用）
	// 和权限管理服务（PermissionService）
	// 获取 MongoDB 数据库
	mongoDB := serviceContainer.GetMongoDB()
	if mongoDB != nil {
		// 创建用户管理仓储
		userAdminRepo := adminrep.NewMongoUserAdminRepository(mongoDB)

		// 创建用户管理服务
		userAdminSvc := adminservice.NewUserAdminService(userAdminRepo)

		// 创建权限管理仓储和服务
		permissionRepo := authRep.NewMongoPermissionRepository(mongoDB)
		permissionSvc := sharedService.NewPermissionService(permissionRepo)

		// 注册管理员路由（包含用户管理和权限管理）
		adminRouter.RegisterAdminRoutes(v1, userSvc, quotaService, auditSvc, adminSvc, configSvc, announcementSvc, userAdminSvc, permissionSvc)
	} else {
		// 如果 MongoDB 不可用，不注册用户管理和权限管理路由
		adminRouter.RegisterAdminRoutes(v1, userSvc, quotaService, auditSvc, adminSvc, configSvc, announcementSvc, nil, nil)
	}

	logger.Info("✓ 管理员路由已注册到: /api/v1/admin/")
	logger.Info("  - /api/v1/admin/users/* (用户管理)")
	logger.Info("  - /api/v1/admin/quota/* (AI配额管理)")
	logger.Info("  - /api/v1/admin/audit/* (审核管理)")
	logger.Info("  - /api/v1/admin/stats (系统统计)")
	logger.Info("  - /api/v1/admin/config/* (配置管理)")
	logger.Info("  - /api/v1/admin/permissions/* (权限管理)")
	logger.Info("  - /api/v1/admin/roles/* (角色管理)")

	// ============ 注册新的用户管理路由（按功能领域组织） ============
	// ⭐ 新架构：按功能领域组织，而非按角色组织
	// 获取书店服务（用于用户作品列表功能）
	bookstoreSvcForUM, bookstoreErrForUM := serviceContainer.GetBookstoreService()
	var bookstoreSvcInterface interface{}
	if bookstoreErrForUM == nil && bookstoreSvcForUM != nil {
		bookstoreSvcInterface = bookstoreSvcForUM
	}

	// 获取统计服务（用于用户统计功能）
	var statsSvc statsService.PlatformStatsService
	var userRepoInstance userRepo.UserRepository // 用于用户路由
	repositoryFactory := serviceContainer.GetRepositoryFactory()
	if repositoryFactory != nil {
		// 创建统计服务所需的 Repository
		userRepoInstance = repositoryFactory.CreateUserRepository()
		bookRepo := repositoryFactory.CreateBookRepository()
		projectRepo := repositoryFactory.CreateProjectRepository()
		chapterRepo := repositoryFactory.CreateBookstoreChapterRepository()

		if userRepoInstance != nil && bookRepo != nil && projectRepo != nil && chapterRepo != nil {
			statsSvc = statsService.NewPlatformStatsService(userRepoInstance, bookRepo, projectRepo, chapterRepo)
		}
	}

	// 注册新的 user 路由
	userRouter.RegisterUserRoutes(v1, userSvc, userRepoInstance, bookstoreSvcInterface, statsSvc)

	logger.Info("✓ 用户路由已注册到: /api/v1/user/")
	logger.Info("  - /api/v1/user/auth/register (用户注册)")
	logger.Info("  - /api/v1/user/auth/login (用户登录)")
	logger.Info("  - /api/v1/user/users/:id (获取用户信息-公开)")
	logger.Info("  - /api/v1/user/users/:id/profile (用户详细资料-公开)")
	logger.Info("  - /api/v1/user/users/:id/books (用户作品列表-公开)")
	logger.Info("  - /api/v1/user/profile (当前用户信息-需认证)")
	logger.Info("  - /api/v1/user/password (修改密码-需认证)")
	logger.Info("  - /api/v1/user/stats/my/* (用户统计-需认证)")

	// ============ 注册系统监控路由 ============
	systemRouter.InitSystemRoutes(v1)
	logger.Info("✓ 系统监控路由已注册到: /api/v1/system/")
	logger.Info("  - /api/v1/system/health (系统健康检查)")
	logger.Info("  - /api/v1/system/health/:service (服务健康检查)")
	logger.Info("  - /api/v1/system/metrics (所有服务指标)")
	logger.Info("  - /api/v1/system/metrics/:service (特定服务指标)")

	// ============ 健康检查 ============
	// 注意: /metrics 端点已在 core/server.go 中注册

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

// initSearchService 初始化搜索服务
// 创建 MongoEngine、BookProvider，并注册到 SearchService
// 如果启用 Elasticsearch，则同时创建 ES 引擎并配置灰度
func initSearchService(container *container.ServiceContainer, logger *zap.Logger) *searchService.SearchService {
	// 获取 MongoDB 客户端和数据库
	mongoClient := container.GetMongoClient()
	mongoDB := container.GetMongoDB()

	if mongoClient == nil || mongoDB == nil {
		logger.Warn("MongoDB client 或 database 未初始化，无法创建搜索服务")
		return nil
	}

	// 创建 MongoEngine
	mongoEngine, err := searchengine.NewMongoEngine(mongoClient, mongoDB)
	if err != nil {
		logger.Error("创建 MongoEngine 失败", zap.Error(err))
		return nil
	}
	logger.Info("✓ MongoEngine 创建成功")

	// 创建 BookProvider 配置
	bookProviderConfig := &searchprovider.BookProviderConfig{
		AllowedStatuses: []string{"ongoing", "completed"},
		AllowedPrivacy:  []bool{false}, // 只允许公开书籍
	}

	// 创建 BookProvider
	bookProvider, err := searchprovider.NewBookProvider(mongoEngine, bookProviderConfig)
	if err != nil {
		logger.Error("创建 BookProvider 失败", zap.Error(err))
		return nil
	}
	logger.Info("✓ BookProvider 创建成功",
		zap.Strings("allowed_statuses", bookProviderConfig.AllowedStatuses),
	)

	// 创建 UserProvider
	userProviderConfig := &searchprovider.UserProviderConfig{}
	userProvider, err := searchprovider.NewUserProvider(mongoEngine, userProviderConfig)
	if err != nil {
		logger.Error("创建 UserProvider 失败", zap.Error(err))
		return nil
	}
	logger.Info("✓ UserProvider 创建成功")

	// 创建 ProjectProvider
	projectProviderConfig := &searchprovider.ProjectProviderConfig{}
	projectProvider, err := searchprovider.NewProjectProvider(mongoEngine, projectProviderConfig)
	if err != nil {
		logger.Error("创建 ProjectProvider 失败", zap.Error(err))
		return nil
	}
	logger.Info("✓ ProjectProvider 创建成功")

	// 创建 DocumentProvider
	documentProviderConfig := &searchprovider.DocumentProviderConfig{}
	documentProvider, err := searchprovider.NewDocumentProvider(mongoEngine, documentProviderConfig)
	if err != nil {
		logger.Error("创建 DocumentProvider 失败", zap.Error(err))
		return nil
	}
	logger.Info("✓ DocumentProvider 创建成功")

	// 创建 SearchService 配置
	searchConfig := &searchService.Config{
		EnableCache:           true,
		DefaultCacheTTL:       300, // 5分钟
		MaxConcurrentSearches: 10,
	}

	// 尝试初始化 Elasticsearch，获取灰度配置
	var grayscaleConfig *searchService.GrayScaleConfig
	esConfig, esEngine := initElasticsearch(logger)
	if esConfig != nil {
		grayscaleConfig = &esConfig.ES.GrayScale
	}

	// 创建灰度决策器
	var grayscaleDecision searchService.GrayScaleDecision
	if grayscaleConfig != nil && grayscaleConfig.Enabled {
		grayscaleDecision = searchService.NewGrayScaleDecision(grayscaleConfig, logger)
		logger.Info("✓ 灰度决策器已创建",
			zap.Int("percent", grayscaleConfig.Percent),
		)
	} else {
		logger.Info("⚠ 灰度未启用，创建默认灰度决策器")
		grayscaleDecision = searchService.NewGrayScaleDecision(nil, logger)
	}

	// 创建 SearchService（传入灰度决策器）
	searchSvc := searchService.NewSearchService(log.Default(), searchConfig, grayscaleDecision)
	logger.Info("✓ SearchService 创建成功（已集成灰度决策器）")

	// 设置 MongoDB 引擎（作为 fallback）
	searchSvc.SetMongoEngine(mongoEngine)
	logger.Info("✓ MongoDB 引擎已设置到 SearchService")

	// 注册 BookProvider 到 SearchService
	searchSvc.RegisterProvider(bookProvider)
	logger.Info("✓ BookProvider 已注册到 SearchService")

	// 注册 UserProvider 到 SearchService
	searchSvc.RegisterProvider(userProvider)
	logger.Info("✓ UserProvider 已注册到 SearchService")

	// 注册 ProjectProvider 到 SearchService
	searchSvc.RegisterProvider(projectProvider)
	logger.Info("✓ ProjectProvider 已注册到 SearchService")

	// 注册 DocumentProvider 到 SearchService
	searchSvc.RegisterProvider(documentProvider)
	logger.Info("✓ DocumentProvider 已注册到 SearchService")

	// 设置 ES 配置和引擎
	if esEngine != nil {
		searchSvc.SetESConfig(esConfig)
		searchSvc.SetESEngine(esEngine)
		logger.Info("✓ Elasticsearch 已集成到 SearchService")

		// 记录灰度配置
		if esConfig.ES.GrayScale.Enabled {
			logger.Info("✓ ES 灰度模式已启用",
				zap.Int("grayscale_percent", esConfig.ES.GrayScale.Percent),
			)
		} else {
			logger.Info("✓ ES 全量模式已启用")
		}
	} else {
		logger.Info("⚠ Elasticsearch 未配置或初始化失败，使用 MongoDB 搜索")
	}

	return searchSvc
}

// initElasticsearch 初始化 Elasticsearch 客户端和引擎
func initElasticsearch(logger *zap.Logger) (*searchService.SearchConfig, searchengine.Engine) {
	// 从环境变量读取 ES 配置
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200" // 默认地址
	}

	esEnabled := os.Getenv("ELASTICSEARCH_ENABLED")
	if esEnabled == "" || esEnabled == "false" {
		logger.Info("Elasticsearch 未启用")
		return nil, nil
	}

	esIndexPrefix := os.Getenv("ELASTICSEARCH_INDEX_PREFIX")
	if esIndexPrefix == "" {
		esIndexPrefix = "qingyu" // 默认前缀
	}

	// 灰度配置
	grayscaleEnabled := os.Getenv("ELASTICSEARCH_GRAYSCALE_ENABLED")
	grayscalePercent := 0
	if grayscaleEnabled == "true" {
		percentStr := os.Getenv("ELASTICSEARCH_GRAYSCALE_PERCENT")
		if percentStr != "" {
			fmt.Sscanf(percentStr, "%d", &grayscalePercent)
		}
		if grayscalePercent == 0 {
			grayscalePercent = 10 // 默认 10%
		}
		logger.Info("Elasticsearch 灰度模式配置",
			zap.Int("percent", grayscalePercent),
		)
	}

	// 创建 ES 客户端
	esClient, err := initElasticsearchClient(esURL)
	if err != nil {
		logger.Error("创建 Elasticsearch 客户端失败", zap.Error(err))
		return nil, nil
	}
	logger.Info("✓ Elasticsearch 客户端创建成功", zap.String("url", esURL))

	// 创建 ElasticsearchEngine
	esEngine, err := searchengine.NewElasticsearchEngine(esClient)
	if err != nil {
		logger.Error("创建 ElasticsearchEngine 失败", zap.Error(err))
		return nil, nil
	}

	// 检查 ES 健康状态
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := esEngine.Health(ctx); err != nil {
		logger.Warn("Elasticsearch 健康检查失败，将使用 MongoDB fallback", zap.Error(err))
		return nil, nil
	}
	logger.Info("✓ Elasticsearch 健康检查通过")

	// 构建 SearchConfig
	searchConfig := &searchService.SearchConfig{
		ES: searchService.ESConfig{
			Enabled:    true,
			URL:        esURL,
			IndexPrefix: esIndexPrefix,
			GrayScale: searchService.GrayScaleConfig{
				Enabled: grayscaleEnabled == "true",
				Percent: grayscalePercent,
			},
		},
	}

	return searchConfig, esEngine
}

// initElasticsearchClient 初始化 Elasticsearch 客户端
func initElasticsearchClient(esURL string) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return client, nil
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
