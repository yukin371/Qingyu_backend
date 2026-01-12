package container

import (
	"context"
	"fmt"
	"sync"
	"time"

	repoInterfaces "Qingyu_backend/repository/interfaces"
	userRepo "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userInterface "Qingyu_backend/service/interfaces/user"

	// Service implementations
	aiService "Qingyu_backend/service/ai"
	bookstoreService "Qingyu_backend/service/bookstore"
	financeService "Qingyu_backend/service/finance"
	readingService "Qingyu_backend/service/reader"
	readingStatsService "Qingyu_backend/service/reader/stats"
	socialService "Qingyu_backend/service/social"
	userService "Qingyu_backend/service/user"
	projectService "Qingyu_backend/service/writer/project"

	// Audit service
	auditSvc "Qingyu_backend/service/audit"

	// Messaging service
	messagingSvc "Qingyu_backend/service/messaging"

	// Notification service
	mongoNotification "Qingyu_backend/repository/mongodb/notification"
	notificationService "Qingyu_backend/service/notification"

	// Shared services
	"Qingyu_backend/service/admin"
	financeWalletService "Qingyu_backend/service/finance/wallet"
	"Qingyu_backend/service/recommendation"
	"Qingyu_backend/service/shared/auth"
	sharedMessaging "Qingyu_backend/service/shared/messaging"
	"Qingyu_backend/service/shared/metrics"
	"Qingyu_backend/service/shared/storage"

	adminModel "Qingyu_backend/models/users"
	adminInterface "Qingyu_backend/repository/interfaces/admin"

	// Infrastructure
	"Qingyu_backend/config"
	"Qingyu_backend/global"
	"Qingyu_backend/pkg/cache"
	"Qingyu_backend/repository/mongodb"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	authModel "Qingyu_backend/models/auth"
)

// ServiceContainer 服务容器
// 负责管理所有服务的生命周期和依赖注入
type ServiceContainer struct {
	repositoryFactory repoInterfaces.RepositoryFactory
	services          map[string]serviceInterfaces.BaseService
	initialized       bool
	mu                sync.RWMutex // 保护并发访问

	// 基础设施
	eventBus    serviceInterfaces.EventBus
	redisClient cache.RedisClient
	mongoClient *mongo.Client
	mongoDB     *mongo.Database

	// 服务指标
	serviceMetrics map[string]*metrics.ServiceMetrics

	// 业务服务
	userService           userInterface.UserService
	aiService             *aiService.Service
	bookstoreService      bookstoreService.BookstoreService
	chapterService        bookstoreService.ChapterService
	bookDetailService     bookstoreService.BookDetailService
	bookRatingService     bookstoreService.BookRatingService
	bookStatisticsService bookstoreService.BookStatisticsService
	readerService         *readingService.ReaderService
	readingStatsService   *readingStatsService.ReadingStatsService
	commentService        *socialService.CommentService
	likeService           *socialService.LikeService
	collectionService     *socialService.CollectionService
	readingHistoryService *readingService.ReadingHistoryService
	projectService        *projectService.ProjectService

	// AI 相关服务
	quotaService *aiService.QuotaService
	chatService  *aiService.ChatService
	phase3Client *aiService.Phase3Client

	// Shared services
	authService           auth.AuthService
	oauthService          *auth.OAuthService
	walletService         financeWalletService.WalletService
	recommendationService recommendation.RecommendationService
	messagingService      sharedMessaging.MessagingService
	storageService        storage.StorageService
	adminService          admin.AdminService
	announcementService   messagingSvc.AnnouncementService
	notificationService   notificationService.NotificationService
	templateService       notificationService.TemplateService

	// 财务服务
	membershipService    financeService.MembershipService
	authorRevenueService financeService.AuthorRevenueService

	// 审核服务
	auditService *auditSvc.ContentAuditService

	// 存储相关具体实现（用于API层）
	storageServiceImpl *storage.StorageServiceImpl
	multipartService   *storage.MultipartUploadService
	imageProcessor     *storage.ImageProcessor
}

// NewServiceContainer 创建服务容器
// Repository工厂将在Initialize()时自动创建
func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		services:       make(map[string]serviceInterfaces.BaseService),
		serviceMetrics: make(map[string]*metrics.ServiceMetrics),
		initialized:    false,
		eventBus:       base.NewSimpleEventBus(), // 创建事件总线
	}
}

// RegisterService 注册服务
func (c *ServiceContainer) RegisterService(name string, service serviceInterfaces.BaseService) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.services[name] != nil {
		return fmt.Errorf("服务 %s 已存在", name)
	}

	c.services[name] = service

	// 为服务创建指标
	c.serviceMetrics[name] = metrics.NewServiceMetrics(
		service.GetServiceName(),
		service.GetVersion(),
	)

	return nil
}

// GetService 获取服务
func (c *ServiceContainer) GetService(name string) (serviceInterfaces.BaseService, error) {
	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("服务 %s 不存在", name)
	}

	return service, nil
}

// ============ 业务服务获取方法 ============

// GetUserService 获取用户服务
func (c *ServiceContainer) GetUserService() (userInterface.UserService, error) {
	if c.userService == nil {
		return nil, fmt.Errorf("UserService未初始化")
	}
	return c.userService, nil
}

// GetAIService 获取AI服务
func (c *ServiceContainer) GetAIService() (*aiService.Service, error) {
	if c.aiService == nil {
		return nil, fmt.Errorf("AIService未初始化")
	}
	return c.aiService, nil
}

// GetBookstoreService 获取书城服务
func (c *ServiceContainer) GetBookstoreService() (bookstoreService.BookstoreService, error) {
	if c.bookstoreService == nil {
		return nil, fmt.Errorf("BookstoreService未初始化")
	}
	return c.bookstoreService, nil
}

// GetBookDetailService 获取书籍详情服务
func (c *ServiceContainer) GetBookDetailService() (bookstoreService.BookDetailService, error) {
	if c.bookDetailService == nil {
		return nil, fmt.Errorf("BookDetailService未初始化")
	}
	return c.bookDetailService, nil
}

// GetBookRatingService 获取书籍评分服务
func (c *ServiceContainer) GetBookRatingService() (bookstoreService.BookRatingService, error) {
	if c.bookRatingService == nil {
		return nil, fmt.Errorf("BookRatingService未初始化")
	}
	return c.bookRatingService, nil
}

// GetBookStatisticsService 获取书籍统计服务
func (c *ServiceContainer) GetBookStatisticsService() (bookstoreService.BookStatisticsService, error) {
	if c.bookStatisticsService == nil {
		return nil, fmt.Errorf("BookStatisticsService未初始化")
	}
	return c.bookStatisticsService, nil
}

// GetChapterService 获取章节服务
func (c *ServiceContainer) GetChapterService() (bookstoreService.ChapterService, error) {
	if c.chapterService == nil {
		return nil, fmt.Errorf("ChapterService未初始化")
	}
	return c.chapterService, nil
}

// getChapterService 内部方法：获取章节服务（简化版，用于依赖注入）
func (c *ServiceContainer) getChapterService() bookstoreService.ChapterService {
	return c.chapterService
}

// GetReaderService 获取阅读器服务
func (c *ServiceContainer) GetReaderService() (*readingService.ReaderService, error) {
	if c.readerService == nil {
		return nil, fmt.Errorf("ReaderService未初始化")
	}
	return c.readerService, nil
}

// GetCommentService 获取评论服务
func (c *ServiceContainer) GetCommentService() (*socialService.CommentService, error) {
	if c.commentService == nil {
		return nil, fmt.Errorf("CommentService未初始化")
	}
	return c.commentService, nil
}

// GetLikeService 获取点赞服务
func (c *ServiceContainer) GetLikeService() (*socialService.LikeService, error) {
	if c.likeService == nil {
		return nil, fmt.Errorf("LikeService未初始化")
	}
	return c.likeService, nil
}

// GetCollectionService 获取收藏服务
func (c *ServiceContainer) GetCollectionService() (*socialService.CollectionService, error) {
	if c.collectionService == nil {
		return nil, fmt.Errorf("CollectionService未初始化")
	}
	return c.collectionService, nil
}

// GetReadingHistoryService 获取阅读历史服务
func (c *ServiceContainer) GetReadingHistoryService() (*readingService.ReadingHistoryService, error) {
	if c.readingHistoryService == nil {
		return nil, fmt.Errorf("ReadingHistoryService未初始化")
	}
	return c.readingHistoryService, nil
}

// GetReadingStatsService 获取阅读统计服务
func (c *ServiceContainer) GetReadingStatsService() (*readingStatsService.ReadingStatsService, error) {
	if c.readingStatsService == nil {
		return nil, fmt.Errorf("ReadingStatsService未初始化")
	}
	return c.readingStatsService, nil
}

// GetQuotaService 获取配额服务
func (c *ServiceContainer) GetQuotaService() (*aiService.QuotaService, error) {
	if c.quotaService == nil {
		return nil, fmt.Errorf("QuotaService未初始化")
	}
	return c.quotaService, nil
}

// GetChatService 获取聊天服务
func (c *ServiceContainer) GetChatService() (*aiService.ChatService, error) {
	if c.chatService == nil {
		return nil, fmt.Errorf("ChatService未初始化")
	}
	return c.chatService, nil
}

// GetPhase3Client 获取Phase3 gRPC客户端
func (c *ServiceContainer) GetPhase3Client() (*aiService.Phase3Client, error) {
	if c.phase3Client == nil {
		return nil, fmt.Errorf("Phase3Client未初始化")
	}
	return c.phase3Client, nil
}

// ============ 共享服务获取方法 ============

// GetAuthService 获取认证服务
func (c *ServiceContainer) GetAuthService() (auth.AuthService, error) {
	if c.authService == nil {
		return nil, fmt.Errorf("AuthService未初始化")
	}
	return c.authService, nil
}

// GetOAuthService 获取OAuth服务
func (c *ServiceContainer) GetOAuthService() (*auth.OAuthService, error) {
	if c.oauthService == nil {
		return nil, fmt.Errorf("OAuthService未初始化")
	}
	return c.oauthService, nil
}

// GetWalletService 获取钱包服务
func (c *ServiceContainer) GetWalletService() (financeWalletService.WalletService, error) {
	if c.walletService == nil {
		return nil, fmt.Errorf("WalletService未初始化")
	}
	return c.walletService, nil
}

// GetRecommendationService 获取推荐服务
func (c *ServiceContainer) GetRecommendationService() (recommendation.RecommendationService, error) {
	if c.recommendationService == nil {
		return nil, fmt.Errorf("RecommendationService未初始化")
	}
	return c.recommendationService, nil
}

// GetMessagingService 获取消息服务
func (c *ServiceContainer) GetMessagingService() (sharedMessaging.MessagingService, error) {
	if c.messagingService == nil {
		return nil, fmt.Errorf("MessagingService未初始化")
	}
	return c.messagingService, nil
}

// GetStorageService 获取存储服务
func (c *ServiceContainer) GetStorageService() (storage.StorageService, error) {
	if c.storageService == nil {
		return nil, fmt.Errorf("StorageService未初始化")
	}
	return c.storageService, nil
}

// GetAdminService 获取管理服务
func (c *ServiceContainer) GetAdminService() (admin.AdminService, error) {
	if c.adminService == nil {
		return nil, fmt.Errorf("AdminService未初始化")
	}
	return c.adminService, nil
}

// GetAuditService 获取审核服务
func (c *ServiceContainer) GetAuditService() (*auditSvc.ContentAuditService, error) {
	if c.auditService == nil {
		return nil, fmt.Errorf("AuditService未初始化")
	}
	return c.auditService, nil
}

// GetAnnouncementService 获取公告服务
func (c *ServiceContainer) GetAnnouncementService() (messagingSvc.AnnouncementService, error) {
	if c.announcementService == nil {
		return nil, fmt.Errorf("AnnouncementService未初始化")
	}
	return c.announcementService, nil
}

// GetNotificationService 获取通知服务
func (c *ServiceContainer) GetNotificationService() (notificationService.NotificationService, error) {
	if c.notificationService == nil {
		return nil, fmt.Errorf("NotificationService未初始化")
	}
	return c.notificationService, nil
}

// GetTemplateService 获取通知模板服务
func (c *ServiceContainer) GetTemplateService() (notificationService.TemplateService, error) {
	if c.templateService == nil {
		return nil, fmt.Errorf("TemplateService未初始化")
	}
	return c.templateService, nil
}

// GetMembershipService 获取会员服务
func (c *ServiceContainer) GetMembershipService() (financeService.MembershipService, error) {
	if c.membershipService == nil {
		return nil, fmt.Errorf("MembershipService未初始化")
	}
	return c.membershipService, nil
}

// GetAuthorRevenueService 获取作者收入服务
func (c *ServiceContainer) GetAuthorRevenueService() (financeService.AuthorRevenueService, error) {
	if c.authorRevenueService == nil {
		return nil, fmt.Errorf("AuthorRevenueService未初始化")
	}
	return c.authorRevenueService, nil
}

// GetEventBus 获取事件总线
func (c *ServiceContainer) GetEventBus() serviceInterfaces.EventBus {
	return c.eventBus
}

// Initialize 初始化所有服务
func (c *ServiceContainer) Initialize(ctx context.Context) error {
	if c.initialized {
		return nil
	}

	// 1. 初始化MongoDB（优先级最高）
	if err := c.initMongoDB(); err != nil {
		return fmt.Errorf("MongoDB初始化失败: %w", err)
	}

	// 2. 创建Repository工厂（使用容器的MongoDB连接）
	c.repositoryFactory = mongodb.NewMongoRepositoryFactoryWithClient(
		c.mongoClient,
		c.mongoDB,
	)

	// 3. 初始化Redis客户端（失败不阻塞）
	if err := c.initRedis(); err != nil {
		// Redis初始化失败不阻塞启动，但记录错误
		fmt.Printf("警告: Redis客户端初始化失败: %v\n", err)
	}

	// 4. 初始化Repository工厂健康检查
	if err := c.repositoryFactory.Health(ctx); err != nil {
		return fmt.Errorf("Repository工厂健康检查失败: %w", err)
	}

	// 5. 初始化所有服务
	for name, service := range c.services {
		if err := service.Initialize(ctx); err != nil {
			return fmt.Errorf("初始化服务 %s 失败: %w", name, err)
		}
	}

	c.initialized = true
	return nil
}

// initMongoDB 初始化MongoDB客户端
func (c *ServiceContainer) initMongoDB() error {
	cfg := config.GlobalConfig.Database
	if cfg == nil || cfg.Primary.MongoDB == nil {
		return fmt.Errorf("MongoDB配置未找到")
	}

	mongoCfg := cfg.Primary.MongoDB

	clientOptions := options.Client().
		ApplyURI(mongoCfg.URI).
		SetMaxPoolSize(mongoCfg.MaxPoolSize).
		SetMinPoolSize(mongoCfg.MinPoolSize).
		SetConnectTimeout(mongoCfg.ConnectTimeout).
		SetServerSelectionTimeout(mongoCfg.ServerTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), mongoCfg.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("连接MongoDB失败: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(ctx)
		return fmt.Errorf("MongoDB连接测试失败: %w", err)
	}

	c.mongoClient = client
	c.mongoDB = client.Database(mongoCfg.Database)

	// 设置全局DB变量（兼容测试和旧代码）
	// 注意：需要导入global包
	return c.setGlobalDB()
}

// initRedis 初始化Redis客户端
func (c *ServiceContainer) initRedis() error {
	cfg := config.GetRedisConfig()
	if cfg == nil {
		return fmt.Errorf("Redis配置未找到")
	}

	client, err := cache.NewRedisClient(cfg)
	if err != nil {
		return fmt.Errorf("创建Redis客户端失败: %w", err)
	}

	c.redisClient = client
	return nil
}

// setGlobalDB 设置全局DB变量（用于测试和向后兼容）
func (c *ServiceContainer) setGlobalDB() error {
	if c.mongoDB == nil {
		return fmt.Errorf("MongoDB未初始化")
	}
	global.DB = c.mongoDB
	return nil
}

// Health 检查所有服务健康状态
func (c *ServiceContainer) Health(ctx context.Context) error {
	// 检查Repository工厂健康状态
	if err := c.repositoryFactory.Health(ctx); err != nil {
		return fmt.Errorf("Repository工厂健康检查失败: %w", err)
	}

	// 检查所有服务健康状态
	for name, service := range c.services {
		if err := service.Health(ctx); err != nil {
			return fmt.Errorf("服务 %s 健康检查失败: %w", name, err)
		}
	}

	return nil
}

// Close 关闭所有服务
func (c *ServiceContainer) Close(ctx context.Context) error {
	var lastErr error

	// 1. 关闭Redis客户端
	if c.redisClient != nil {
		if err := c.redisClient.Close(); err != nil {
			fmt.Printf("警告: 关闭Redis客户端失败: %v\n", err)
			lastErr = fmt.Errorf("关闭Redis客户端失败: %w", err)
		}
	}

	// 2. 关闭所有服务
	for name, service := range c.services {
		if err := service.Close(ctx); err != nil {
			lastErr = fmt.Errorf("关闭服务 %s 失败: %w", name, err)
		}
	}

	// 3. 关闭MongoDB（在Repository工厂之前关闭）
	if c.mongoClient != nil {
		if err := c.mongoClient.Disconnect(ctx); err != nil {
			fmt.Printf("警告: 关闭MongoDB客户端失败: %v\n", err)
			lastErr = fmt.Errorf("关闭MongoDB客户端失败: %w", err)
		}
	}

	// 4. 关闭Repository工厂
	if c.repositoryFactory != nil {
		if err := c.repositoryFactory.Close(); err != nil {
			lastErr = fmt.Errorf("关闭Repository工厂失败: %w", err)
		}
	}

	c.initialized = false
	return lastErr
}

// GetRedisClient 获取Redis客户端
func (c *ServiceContainer) GetRedisClient() cache.RedisClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.redisClient
}

// GetMongoDB 获取MongoDB数据库实例
func (c *ServiceContainer) GetMongoDB() *mongo.Database {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mongoDB
}

// GetMongoClient 获取MongoDB客户端
func (c *ServiceContainer) GetMongoClient() *mongo.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mongoClient
}

// GetRepositoryFactory 获取Repository工厂
func (c *ServiceContainer) GetRepositoryFactory() repoInterfaces.RepositoryFactory {
	return c.repositoryFactory
}

// SetupDefaultServices 设置默认服务
// 创建并注册所有核心业务服务
func (c *ServiceContainer) SetupDefaultServices() error {
	// ============ 1. 创建用户服务 ============
	userRepo := c.repositoryFactory.CreateUserRepository()
	authRepo := c.repositoryFactory.CreateAuthRepository()
	c.userService = userService.NewUserService(userRepo, authRepo)
	// 用户服务实现了BaseService接口，可以注册
	if err := c.RegisterService("UserService", c.userService); err != nil {
		return fmt.Errorf("注册用户服务失败: %w", err)
	}

	// ============ 2. 创建书城服务 ============
	// RepositoryFactory已返回具体类型，无需类型断言
	bookRepo := c.repositoryFactory.CreateBookRepository()
	categoryRepo := c.repositoryFactory.CreateCategoryRepository()
	bannerRepo := c.repositoryFactory.CreateBannerRepository()
	rankingRepo := c.repositoryFactory.CreateRankingRepository()
	chapterRepo := c.repositoryFactory.CreateBookstoreChapterRepository()      // ← 创建章节仓储
	chapterContentRepo := c.repositoryFactory.CreateChapterContentRepository() // ← 创建章节内容仓储

	c.bookstoreService = bookstoreService.NewBookstoreService(
		bookRepo,
		categoryRepo,
		bannerRepo,
		rankingRepo,
	)
	// 注意：BookstoreService 不完全实现 BaseService，不注册到 services map

	// 创建章节服务（注入 ChapterContentRepository）
	c.chapterService = bookstoreService.NewChapterService(chapterRepo, chapterContentRepo, nil) // ← 创建 ChapterService (CacheService暂时为nil)

	// 创建书店详细服务
	bookDetailRepo := c.repositoryFactory.CreateBookDetailRepository()
	bookRatingRepo := c.repositoryFactory.CreateBookRatingRepository()
	bookStatisticsRepo := c.repositoryFactory.CreateBookStatisticsRepository()

	// 这些服务也需要 CacheService，暂时传 nil
	c.bookDetailService = bookstoreService.NewBookDetailService(bookDetailRepo, nil)
	c.bookRatingService = bookstoreService.NewBookRatingService(bookRatingRepo, nil)
	c.bookStatisticsService = bookstoreService.NewBookStatisticsService(bookStatisticsRepo, nil)

	// ============ 3. 创建阅读器服务 ============
	progressRepo := c.repositoryFactory.CreateReadingProgressRepository()
	annotationRepo := c.repositoryFactory.CreateAnnotationRepository()
	settingsRepo := c.repositoryFactory.CreateReadingSettingsRepository()
	chapterService := c.getChapterService() // ← 获取 ChapterService

	// 创建缓存服务和VIP服务
	var cacheService readingService.ReaderCacheService
	var vipService readingService.VIPPermissionService

	if c.redisClient != nil {
		// 获取原始 Redis 客户端
		rawClient := c.redisClient.GetClient()
		if redisClient, ok := rawClient.(*redis.Client); ok {
			cacheService = readingService.NewRedisReaderCacheService(redisClient, "qingyu")
			vipService = readingService.NewVIPPermissionService(redisClient, "qingyu")
		}
	}

	c.readerService = readingService.NewReaderService(
		progressRepo,
		annotationRepo,
		settingsRepo,
		chapterService, // ← 注入 ChapterService
		c.eventBus,     // 注入事件总线
		cacheService,   // 注入缓存服务
		vipService,     // 注入VIP服务
	)
	// 注意：ReaderService 不完全实现 BaseService，不注册到 services map

	// ============ 4. 创建评论服务 ============
	commentRepo := c.repositoryFactory.CreateCommentRepository()
	sensitiveWordRepo := c.repositoryFactory.CreateSensitiveWordRepository()

	c.commentService = socialService.NewCommentService(
		commentRepo,
		sensitiveWordRepo, // 可以为nil，表示不启用敏感词检测
		c.eventBus,        // 注入事件总线
	)
	c.services["CommentService"] = c.commentService

	// ============ 4.5 创建点赞服务 ============
	likeRepo := c.repositoryFactory.CreateLikeRepository()

	c.likeService = socialService.NewLikeService(
		likeRepo,
		commentRepo, // 用于更新评论点赞数
		c.eventBus,  // 注入事件总线
	)
	c.services["LikeService"] = c.likeService

	// ============ 4.6 创建收藏服务 ============
	collectionRepo := c.repositoryFactory.CreateCollectionRepository()

	c.collectionService = socialService.NewCollectionService(
		collectionRepo,
		c.eventBus, // 注入事件总线
	)
	c.services["CollectionService"] = c.collectionService

	// ============ 4.7 创建阅读历史服务 ============
	readingHistoryRepo := c.repositoryFactory.CreateReadingHistoryRepository()

	c.readingHistoryService = readingService.NewReadingHistoryService(
		readingHistoryRepo,
		c.eventBus, // 注入事件总线
	)
	c.services["ReadingHistoryService"] = c.readingHistoryService

	// ============ 4.8 创建阅读统计服务 ============
	chapterStatsRepo := c.repositoryFactory.CreateChapterStatsRepository()
	readerBehaviorRepo := c.repositoryFactory.CreateReaderBehaviorRepository()
	bookStatsRepo := c.repositoryFactory.CreateBookStatsRepository()
	c.readingStatsService = readingStatsService.NewReadingStatsService(
		chapterStatsRepo,
		readerBehaviorRepo,
		bookStatsRepo,
	)
	c.services["ReadingStatsService"] = c.readingStatsService

	// ============ 4.9 创建项目服务 ============
	projectRepo := c.repositoryFactory.CreateProjectRepository()
	c.projectService = projectService.NewProjectService(
		projectRepo,
		c.eventBus,
	)
	// 注册ProjectService
	if err := c.RegisterService("ProjectService", c.projectService); err != nil {
		return fmt.Errorf("注册项目服务失败: %w", err)
	}

	// ============ 5. 创建AI服务 ============
	// AIService需要ProjectService来构建上下文
	c.aiService = aiService.NewServiceWithDependencies(c.projectService)

	// 创建AI配额服务
	quotaRepo := c.repositoryFactory.CreateQuotaRepository()
	c.quotaService = aiService.NewQuotaService(quotaRepo)
	// 注意：QuotaService 不完全实现 BaseService，不注册到 services map

	// 创建聊天服务（使用临时内存存储）
	chatRepo := aiService.NewInMemoryChatRepository()
	c.chatService = aiService.NewChatService(c.aiService, chatRepo)
	// 注意：ChatService 不完全实现 BaseService，不注册到 services map

	// ============ 5. 共享服务初始化 ============

	// 5.1 创建 WalletService（简单版，只需要 WalletRepository）
	walletRepo := c.repositoryFactory.CreateWalletRepository()
	walletSvc := financeWalletService.NewUnifiedWalletService(walletRepo)
	c.walletService = walletSvc // 保存为接口类型

	// 类型断言为 BaseService，以便注册到服务映射
	if baseWalletSvc, ok := walletSvc.(serviceInterfaces.BaseService); ok {
		if err := c.RegisterService("WalletService", baseWalletSvc); err != nil {
			return fmt.Errorf("注册钱包服务失败: %w", err)
		}
	} else {
		return fmt.Errorf("WalletService 未实现 BaseService 接口")
	}

	// 5.2 创建 AuthService（完整实现，支持Redis降级）
	authRepo = c.repositoryFactory.CreateAuthRepository()
	oauthRepo := c.repositoryFactory.CreateOAuthRepository()

	// 创建Redis适配器或内存降级方案
	// 注意：使用具体类型而不是接口，以便同时实现RedisClient和CacheClient接口
	var redisAdapter interface{}
	if c.redisClient != nil {
		// 使用Redis Token黑名单
		redisImpl := auth.NewRedisAdapter(c.redisClient)
		redisAdapter = redisImpl
		fmt.Println("✓ AuthService使用Redis Token黑名单")
	} else {
		// 降级到内存Token黑名单
		redisImpl := auth.NewInMemoryTokenBlacklist()
		redisAdapter = redisImpl
		fmt.Println("⚠ Redis不可用，AuthService使用内存Token黑名单（降级模式）")
		fmt.Println("  注意：内存模式不支持分布式部署，服务器重启后黑名单会丢失")
	}

	// 创建子服务
	jwtService := auth.NewJWTService(config.GetJWTConfigEnhanced(), redisAdapter.(auth.RedisClient))
	roleService := auth.NewRoleService(authRepo)

	// 类型断言为CacheClient
	cacheClient, ok := redisAdapter.(auth.CacheClient)
	if !ok {
		return fmt.Errorf("redisAdapter does not implement CacheClient")
	}
	permissionService := auth.NewPermissionService(authRepo, cacheClient)
	sessionService := auth.NewSessionService(cacheClient)

	// 创建AuthService
	c.authService = auth.NewAuthService(
		jwtService,
		roleService,
		permissionService,
		authRepo,
		oauthRepo,
		c.userService,
		sessionService,
	)

	// 类型断言为BaseService，以便注册到服务映射
	if baseAuthSvc, ok := c.authService.(serviceInterfaces.BaseService); ok {
		if err := c.RegisterService("AuthService", baseAuthSvc); err != nil {
			return fmt.Errorf("注册认证服务失败: %w", err)
		}
	}

	// 5.2.1 创建 OAuthService（可选，需要配置）
	// 获取OAuth配置
	oauthConfigs := make(map[string]*authModel.OAuthConfig)
	// TODO: 从配置文件加载OAuth配置
	// oauthConfigs["google"] = &authModel.OAuthConfig{
	//     Enabled: true,
	//     ClientID: config.GetOAuthConfig("google").ClientID,
	//     ClientSecret: config.GetOAuthConfig("google").ClientSecret,
	//     RedirectURI: config.GetOAuthConfig("google").RedirectURI,
	//     Scopes: "openid profile email",
	// }

	if len(oauthConfigs) > 0 {
		// 如果有OAuth配置，创建OAuthService
		logger, err := zap.NewProduction()
		if err == nil {
			oauthSvc, err := auth.NewOAuthService(logger, oauthRepo, oauthConfigs)
			if err != nil {
				fmt.Printf("警告: 创建OAuthService失败: %v\n", err)
			} else {
				c.oauthService = oauthSvc
				fmt.Println("  ✓ OAuthService初始化完成")
			}
		}
	} else {
		fmt.Println("  ℹ OAuth配置为空，跳过OAuthService创建（OAuth登录功能将不可用）")
	}

	// 5.3 创建 RecommendationService
	if c.redisClient != nil {
		recRepo := c.repositoryFactory.CreateRecommendationRepository()
		recAdapter := recommendation.NewRedisAdapter(c.redisClient)
		recSvc := recommendation.NewRecommendationService(recRepo, recAdapter)
		c.recommendationService = recSvc

		// 类型断言为 BaseService，以便注册到服务映射
		if baseRecSvc, ok := recSvc.(serviceInterfaces.BaseService); ok {
			if err := c.RegisterService("RecommendationService", baseRecSvc); err != nil {
				return fmt.Errorf("注册推荐服务失败: %w", err)
			}
		}
	} else {
		fmt.Println("警告: Redis客户端未初始化，跳过RecommendationService创建")
	}

	// 5.4 StorageService（Phase2快速通道）
	fmt.Println("初始化 StorageService...")
	storageRepo := c.repositoryFactory.CreateStorageRepository()

	// 使用本地文件系统Backend（快速通道方案）
	localBackend := storage.NewLocalBackend("./uploads", "http://localhost:8080/api/v1/files")

	// 适配StorageRepository到FileRepository接口
	fileRepo := storage.NewRepositoryAdapter(storageRepo)
	storageSvc := storage.NewStorageService(localBackend, fileRepo)
	c.storageServiceImpl = storageSvc.(*storage.StorageServiceImpl)
	c.storageService = storageSvc

	// 注册为BaseService
	if baseStorageSvc, ok := storageSvc.(serviceInterfaces.BaseService); ok {
		if err := c.RegisterService("StorageService", baseStorageSvc); err != nil {
			return fmt.Errorf("注册存储服务失败: %w", err)
		}
		fmt.Println("  ✓ StorageService 已注册")
	}

	// 初始化MultipartUploadService
	multipartSvc := storage.NewMultipartUploadService(localBackend, storageRepo)
	c.multipartService = multipartSvc

	// 初始化ImageProcessor
	imageProcessor := storage.NewImageProcessor(localBackend)
	c.imageProcessor = imageProcessor

	fmt.Println("  ✓ StorageService完整初始化完成（LocalBackend）")

	// 5.5 AdminService
	// 使用 RepositoryFactory 获取 Admin Repository
	adminAuditRepoImpl := c.repositoryFactory.CreateAuditRepository()
	adminLogRepoImpl := c.repositoryFactory.CreateAdminLogRepository()

	// 创建适配器
	adminAuditRepo := &adminAuditRepositoryAdapter{repo: adminAuditRepoImpl}
	adminLogRepo := &adminLogRepositoryAdapter{repo: adminLogRepoImpl}

	// 创建 AdminService 需要的简化 UserRepository 适配器
	adminUserRepo := &adminUserRepositoryAdapter{userRepo: c.repositoryFactory.CreateUserRepository()}

	// 创建 AdminService
	adminSvc := admin.NewAdminService(adminAuditRepo, adminLogRepo, adminUserRepo)
	c.adminService = adminSvc

	if baseAdminSvc, ok := adminSvc.(serviceInterfaces.BaseService); ok {
		if err := c.RegisterService("AdminService", baseAdminSvc); err != nil {
			return fmt.Errorf("注册管理服务失败: %w", err)
		}
	}
	fmt.Println("  ✓ AdminService初始化完成")

	// 5.6 MessagingService
	if c.redisClient != nil {
		rawClient := c.redisClient.GetClient()
		if redisClient, ok := rawClient.(*redis.Client); ok {
			queueClient := sharedMessaging.NewRedisQueueClient(redisClient)
			messagingSvc := sharedMessaging.NewMessagingService(queueClient)
			c.messagingService = messagingSvc

			if baseMessagingSvc, ok := messagingSvc.(serviceInterfaces.BaseService); ok {
				if err := c.RegisterService("MessagingService", baseMessagingSvc); err != nil {
					return fmt.Errorf("注册消息服务失败: %w", err)
				}
			}
			fmt.Println("  ✓ MessagingService初始化完成（Redis Stream）")
		} else {
			fmt.Println("警告: Redis客户端类型转换失败，跳过MessagingService创建")
		}
	} else {
		fmt.Println("警告: Redis客户端未初始化，跳过MessagingService创建")
	}

	// 5.7 AnnouncementService
	announcementRepo := c.repositoryFactory.CreateAnnouncementRepository()
	announcementSvc := messagingSvc.NewAnnouncementService(announcementRepo)
	c.announcementService = announcementSvc

	if baseAnnouncementSvc, ok := announcementSvc.(serviceInterfaces.BaseService); ok {
		if err := c.RegisterService("AnnouncementService", baseAnnouncementSvc); err != nil {
			return fmt.Errorf("注册公告服务失败: %w", err)
		}
	}
	fmt.Println("  ✓ AnnouncementService初始化完成")

	// 5.8 NotificationService
	notificationRepo := mongoNotification.NewNotificationRepository(c.mongoDB)
	preferenceRepo := mongoNotification.NewNotificationPreferenceRepository(c.mongoDB)
	pushDeviceRepo := mongoNotification.NewPushDeviceRepository(c.mongoDB)
	templateRepo := mongoNotification.NewNotificationTemplateRepository(c.mongoDB)

	notificationSvc := notificationService.NewNotificationService(
		notificationRepo,
		preferenceRepo,
		pushDeviceRepo,
		templateRepo,
	)
	c.notificationService = notificationSvc

	templateSvc := notificationService.NewTemplateService(templateRepo)
	c.templateService = templateSvc

	// 初始化默认模板
	if err := templateSvc.InitializeDefaultTemplates(context.Background()); err != nil {
		fmt.Printf("警告: 初始化默认通知模板失败: %v\n", err)
	}

	if baseNotificationSvc, ok := notificationSvc.(serviceInterfaces.BaseService); ok {
		if err := c.RegisterService("NotificationService", baseNotificationSvc); err != nil {
			return fmt.Errorf("注册通知服务失败: %w", err)
		}
	}
	fmt.Println("  ✓ NotificationService初始化完成")

	// 5.10 Finance Services
	membershipRepo := c.repositoryFactory.CreateMembershipRepository()
	c.membershipService = financeService.NewMembershipService(membershipRepo)

	authorRevenueRepo := c.repositoryFactory.CreateAuthorRevenueRepository()
	c.authorRevenueService = financeService.NewAuthorRevenueService(authorRevenueRepo)

	fmt.Println("  ✓ Finance服务初始化完成")

	// 5.11 AuditService - 暂时为可选，在service/audit实现完成后再完整初始化
	// TODO: 完整的AuditService初始化逻辑需要在service/audit完全实现后添加
	fmt.Println("  ℹ AuditService初始化跳过（标记为可选）")

	// ============ 6. 初始化所有已注册的服务 ============
	// 注意：SetupDefaultServices 在 Initialize 之后调用，所以这里需要手动初始化新注册的服务
	ctx := context.Background()
	for name, service := range c.services {
		if err := service.Initialize(ctx); err != nil {
			return fmt.Errorf("初始化服务 %s 失败: %w", name, err)
		}
	}

	return nil
}

// SetAuthService 设置认证服务
func (c *ServiceContainer) SetAuthService(service auth.AuthService) {
	c.authService = service
}

// SetOAuthService 设置OAuth服务
func (c *ServiceContainer) SetOAuthService(service *auth.OAuthService) {
	c.oauthService = service
}

// SetWalletService 设置钱包服务
func (c *ServiceContainer) SetWalletService(service financeWalletService.WalletService) {
	c.walletService = service
}

// SetRecommendationService 设置推荐服务
func (c *ServiceContainer) SetRecommendationService(service recommendation.RecommendationService) {
	c.recommendationService = service
}

// SetMessagingService 设置消息服务
func (c *ServiceContainer) SetMessagingService(service sharedMessaging.MessagingService) {
	c.messagingService = service
}

// SetStorageService 设置存储服务
func (c *ServiceContainer) SetStorageService(service storage.StorageService) {
	c.storageService = service
}

// SetStorageServiceImpl 设置存储服务实现
func (c *ServiceContainer) SetStorageServiceImpl(service *storage.StorageServiceImpl) {
	c.storageServiceImpl = service
}

// SetMultipartUploadService 设置分片上传服务
func (c *ServiceContainer) SetMultipartUploadService(service *storage.MultipartUploadService) {
	c.multipartService = service
}

// SetImageProcessor 设置图片处理器
func (c *ServiceContainer) SetImageProcessor(processor *storage.ImageProcessor) {
	c.imageProcessor = processor
}

// GetStorageServiceImpl 获取存储服务实现
func (c *ServiceContainer) GetStorageServiceImpl() (*storage.StorageServiceImpl, error) {
	if c.storageServiceImpl == nil {
		return nil, fmt.Errorf("StorageServiceImpl未初始化")
	}
	return c.storageServiceImpl, nil
}

// GetMultipartUploadService 获取分片上传服务
func (c *ServiceContainer) GetMultipartUploadService() (*storage.MultipartUploadService, error) {
	if c.multipartService == nil {
		return nil, fmt.Errorf("MultipartUploadService未初始化")
	}
	return c.multipartService, nil
}

// GetImageProcessor 获取图片处理器
func (c *ServiceContainer) GetImageProcessor() (*storage.ImageProcessor, error) {
	if c.imageProcessor == nil {
		return nil, fmt.Errorf("ImageProcessor未初始化")
	}
	return c.imageProcessor, nil
}

// SetAdminService 设置管理服务
func (c *ServiceContainer) SetAdminService(service admin.AdminService) {
	c.adminService = service
}

// GetServiceNames 获取所有服务名称
func (c *ServiceContainer) GetServiceNames() []string {
	names := make([]string, 0, len(c.services))
	for name := range c.services {
		names = append(names, name)
	}
	return names
}

// ============ AdminService 适配器 ============

// adminUserRepositoryAdapter 将 UserRepository 适配为 AdminService 需要的接口
type adminUserRepositoryAdapter struct {
	userRepo userRepo.UserRepository
}

// 确保实现了 admin.UserRepository 接口
var _ admin.UserRepository = (*adminUserRepositoryAdapter)(nil)

// GetStatistics 获取用户统计信息
func (a *adminUserRepositoryAdapter) GetStatistics(ctx context.Context, userID string) (*admin.UserStatistics, error) {
	// 简化实现，返回基本统计信息
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &admin.UserStatistics{
		UserID:           user.ID,
		TotalBooks:       0,   // TODO: 从 BookRepository 获取
		TotalChapters:    0,   // TODO: 从 ChapterRepository 获取
		TotalWords:       0,   // TODO: 从统计数据获取
		TotalReads:       0,   // TODO: 从 ReadingProgress 获取
		TotalIncome:      0.0, // TODO: 从 Wallet 获取
		RegistrationDate: user.CreatedAt,
		LastLoginDate:    user.LastLoginAt,
	}, nil
}

// BanUser 封禁用户
func (a *adminUserRepositoryAdapter) BanUser(ctx context.Context, userID, reason string, until time.Time) error {
	// 使用 UserRepository 的 Update 方法
	updates := map[string]interface{}{
		"status":     "banned",
		"ban_reason": reason,
		"ban_until":  until,
		"updated_at": time.Now(),
	}
	return a.userRepo.Update(ctx, userID, updates)
}

// UnbanUser 解封用户
func (a *adminUserRepositoryAdapter) UnbanUser(ctx context.Context, userID string) error {
	// 使用 UserRepository 的 Update 方法
	updates := map[string]interface{}{
		"status":     "active",
		"ban_reason": "",
		"ban_until":  nil,
		"updated_at": time.Now(),
	}
	return a.userRepo.Update(ctx, userID, updates)
}

// ============ Admin Repository Adapters ============
// 将新的 admin repository 接口适配为 shared/admin 服务需要的接口

// adminAuditRepositoryAdapter 适配器
type adminAuditRepositoryAdapter struct {
	repo adminInterface.AuditRepository
}

var _ admin.AuditRepository = (*adminAuditRepositoryAdapter)(nil)

// Create 创建审核记录
func (a *adminAuditRepositoryAdapter) Create(ctx context.Context, record *admin.AuditRecord) error {
	adminRecord := &adminModel.AuditRecord{
		ID:          record.ID,
		ContentID:   record.ContentID,
		ContentType: record.ContentType,
		Status:      record.Status,
		ReviewerID:  record.ReviewerID,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}
	return a.repo.CreateAuditRecord(ctx, adminRecord)
}

// Get 获取审核记录
func (a *adminAuditRepositoryAdapter) Get(ctx context.Context, recordID string) (*admin.AuditRecord, error) {
	adminRecord, err := a.repo.GetAuditRecord(ctx, recordID, "")
	if err != nil {
		return nil, err
	}
	return &admin.AuditRecord{
		ID:          adminRecord.ID,
		ContentID:   adminRecord.ContentID,
		ContentType: adminRecord.ContentType,
		Status:      adminRecord.Status,
		ReviewerID:  adminRecord.ReviewerID,
		Reason:      "",          // 从 adminModel 中没有这个字段
		ReviewedAt:  time.Time{}, // 从 adminModel 中没有这个字段
		CreatedAt:   adminRecord.CreatedAt,
		UpdatedAt:   adminRecord.UpdatedAt,
	}, nil
}

// Update 更新审核记录
func (a *adminAuditRepositoryAdapter) Update(ctx context.Context, recordID string, updates map[string]interface{}) error {
	return a.repo.UpdateAuditRecord(ctx, recordID, updates)
}

// ListByStatus 按状态列出审核记录
func (a *adminAuditRepositoryAdapter) ListByStatus(ctx context.Context, contentType, status string) ([]*admin.AuditRecord, error) {
	filter := &adminInterface.AuditFilter{
		ContentType: contentType,
		Status:      status,
	}
	records, err := a.repo.ListAuditRecords(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]*admin.AuditRecord, len(records))
	for i, r := range records {
		result[i] = &admin.AuditRecord{
			ID:          r.ID,
			ContentID:   r.ContentID,
			ContentType: r.ContentType,
			Status:      r.Status,
			ReviewerID:  r.ReviewerID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}
	return result, nil
}

// ListByContent 按内容列出审核记录
func (a *adminAuditRepositoryAdapter) ListByContent(ctx context.Context, contentID string) ([]*admin.AuditRecord, error) {
	record, err := a.repo.GetAuditRecord(ctx, contentID, "")
	if err != nil {
		return nil, err
	}

	return []*admin.AuditRecord{{
		ID:          record.ID,
		ContentID:   record.ContentID,
		ContentType: record.ContentType,
		Status:      record.Status,
		ReviewerID:  record.ReviewerID,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}}, nil
}

// adminLogRepositoryAdapter 适配器
type adminLogRepositoryAdapter struct {
	repo adminInterface.AdminLogRepository
}

var _ admin.LogRepository = (*adminLogRepositoryAdapter)(nil)

// Create 创建日志
func (a *adminLogRepositoryAdapter) Create(ctx context.Context, log *admin.AdminLog) error {
	adminLog := &adminModel.AdminLog{
		ID:        log.ID,
		AdminID:   log.AdminID,
		Operation: log.Operation,
		Target:    log.Target, // shared/admin 只有一个 Target 字段
		Details:   log.Details,
		CreatedAt: log.CreatedAt,
	}
	return a.repo.CreateAdminLog(ctx, adminLog)
}

// List 列出日志
func (a *adminLogRepositoryAdapter) List(ctx context.Context, filter *admin.LogFilter) ([]*admin.AdminLog, error) {
	adminFilter := &adminInterface.AdminLogFilter{
		AdminID:   filter.AdminID,
		Operation: filter.Operation,
		StartDate: filter.StartDate,
		EndDate:   filter.EndDate,
		Limit:     int64(filter.PageSize),                     // 将 PageSize 转换为 Limit
		Offset:    int64((filter.Page - 1) * filter.PageSize), // 将 Page 转换为 Offset
	}

	logs, err := a.repo.ListAdminLogs(ctx, adminFilter)
	if err != nil {
		return nil, err
	}

	result := make([]*admin.AdminLog, len(logs))
	for i, l := range logs {
		result[i] = &admin.AdminLog{
			ID:        l.ID,
			AdminID:   l.AdminID,
			Operation: l.Operation,
			Target:    l.Target, // adminModel 只有一个 Target 字段
			Details:   l.Details,
			IP:        l.IP, // adminModel 使用 IP 字段
			CreatedAt: l.CreatedAt,
		}
	}
	return result, nil
}

// IsInitialized 检查是否已初始化
func (c *ServiceContainer) IsInitialized() bool {
	return c.initialized
}

// ============ 指标管理方法 ============

// GetServiceMetrics 获取指定服务的指标
func (c *ServiceContainer) GetServiceMetrics(serviceName string) (*metrics.ServiceMetrics, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metric, exists := c.serviceMetrics[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务 %s 的指标不存在", serviceName)
	}

	return metric, nil
}

// GetAllServicesMetrics 获取所有服务的指标
func (c *ServiceContainer) GetAllServicesMetrics() map[string]*metrics.ServiceMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*metrics.ServiceMetrics)
	for name, metric := range c.serviceMetrics {
		result[name] = metric
	}

	return result
}

// GetAllServicesHealth 获取所有服务的健康状态
func (c *ServiceContainer) GetAllServicesHealth(ctx context.Context) map[string]bool {
	c.mu.RLock()
	services := make(map[string]serviceInterfaces.BaseService)
	for name, service := range c.services {
		services[name] = service
	}
	c.mu.RUnlock()

	result := make(map[string]bool)
	for name, service := range services {
		err := service.Health(ctx)
		healthy := err == nil
		result[name] = healthy

		// 更新指标
		if metric, exists := c.serviceMetrics[name]; exists {
			metric.RecordHealthCheck(healthy)
		}
	}

	return result
}
