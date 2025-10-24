package container

import (
	"context"
	"fmt"
	"sync"

	repoInterfaces "Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userInterface "Qingyu_backend/service/interfaces/user"

	// Service implementations
	aiService "Qingyu_backend/service/ai"
	bookstoreService "Qingyu_backend/service/bookstore"
	readingService "Qingyu_backend/service/reading"
	userService "Qingyu_backend/service/user"

	// Shared services
	"Qingyu_backend/service/shared/admin"
	"Qingyu_backend/service/shared/auth"
	"Qingyu_backend/service/shared/messaging"
	"Qingyu_backend/service/shared/metrics"
	"Qingyu_backend/service/shared/recommendation"
	"Qingyu_backend/service/shared/storage"
	"Qingyu_backend/service/shared/wallet"

	// Infrastructure
	"Qingyu_backend/config"
	"Qingyu_backend/pkg/cache"
	"Qingyu_backend/repository/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	userService      userInterface.UserService
	aiService        *aiService.Service
	bookstoreService bookstoreService.BookstoreService
	readerService    *readingService.ReaderService

	// AI 相关服务
	quotaService *aiService.QuotaService
	chatService  *aiService.ChatService

	// 共享服务
	authService           auth.AuthService
	walletService         wallet.WalletService
	recommendationService recommendation.RecommendationService
	messagingService      messaging.MessagingService
	storageService        storage.StorageService
	adminService          admin.AdminService
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

// GetReaderService 获取阅读器服务
func (c *ServiceContainer) GetReaderService() (*readingService.ReaderService, error) {
	if c.readerService == nil {
		return nil, fmt.Errorf("ReaderService未初始化")
	}
	return c.readerService, nil
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

// ============ 共享服务获取方法 ============

// GetAuthService 获取认证服务
func (c *ServiceContainer) GetAuthService() (auth.AuthService, error) {
	if c.authService == nil {
		return nil, fmt.Errorf("AuthService未初始化")
	}
	return c.authService, nil
}

// GetWalletService 获取钱包服务
func (c *ServiceContainer) GetWalletService() (wallet.WalletService, error) {
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
func (c *ServiceContainer) GetMessagingService() (messaging.MessagingService, error) {
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
	return nil
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
	c.userService = userService.NewUserService(userRepo)
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

	c.bookstoreService = bookstoreService.NewBookstoreService(
		bookRepo,
		categoryRepo,
		bannerRepo,
		rankingRepo,
	)
	// 注意：BookstoreService 不完全实现 BaseService，不注册到 services map

	// ============ 3. 创建阅读器服务 ============
	chapterRepo := c.repositoryFactory.CreateChapterRepository()
	progressRepo := c.repositoryFactory.CreateReadingProgressRepository()
	annotationRepo := c.repositoryFactory.CreateAnnotationRepository()
	settingsRepo := c.repositoryFactory.CreateReadingSettingsRepository()

	c.readerService = readingService.NewReaderService(
		chapterRepo,
		progressRepo,
		annotationRepo,
		settingsRepo,
		c.eventBus, // ✅ 注入事件总线
		nil,        // cacheService - TODO: 实现缓存服务
		nil,        // vipService - TODO: 实现VIP服务
	)
	// 注意：ReaderService 不完全实现 BaseService，不注册到 services map

	// ============ 4. 创建AI服务 ============
	c.aiService = aiService.NewService()

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
	walletSvc := wallet.NewUnifiedWalletService(walletRepo)
	c.walletService = walletSvc // 保存为接口类型

	// 类型断言为 BaseService，以便注册到服务映射
	if baseWalletSvc, ok := walletSvc.(serviceInterfaces.BaseService); ok {
		if err := c.RegisterService("WalletService", baseWalletSvc); err != nil {
			return fmt.Errorf("注册钱包服务失败: %w", err)
		}
	} else {
		return fmt.Errorf("WalletService 未实现 BaseService 接口")
	}

	// 5.2 创建 AuthService（复杂，需要多个子服务）
	// TODO: 完整实现 AuthService 需要配置以下依赖：
	//   - JWTService: 需要 JWT配置 和 RedisClient
	//   - RoleService: 需要 AuthRepository
	//   - PermissionService: 需要 AuthRepository
	//   - SessionService: 需要 RedisClient
	//   - UserService: 已在上面创建
	//
	// 示例实现：
	// authRepo := c.repositoryFactory.CreateAuthRepository()
	// jwtService := auth.NewJWTService(config.GetJWTConfigEnhanced(), nil) // nil 表示暂不使用 Redis
	// roleService := auth.NewRoleService(authRepo)
	// permissionService := auth.NewPermissionService(authRepo)
	// sessionService := auth.NewSessionService(nil) // nil 表示暂不使用 Redis
	// c.authService = auth.NewAuthService(
	//     jwtService,
	//     roleService,
	//     permissionService,
	//     authRepo,
	//     c.userService,
	//     sessionService,
	// )
	// if err := c.RegisterService("AuthService", c.authService); err != nil {
	//     return fmt.Errorf("注册认证服务失败: %w", err)
	// }
	//
	// 注意：完整实现需要 Redis 客户端，暂时跳过，留待后续配置

	// 5.3 其他共享服务（暂未实现）
	// TODO: RecommendationService - 需要 RecommendationRepository 和 Redis
	// TODO: MessagingService - 需要 Redis/RabbitMQ 等消息队列
	// TODO: StorageService - 需要 StorageBackend 和 FileRepository
	// TODO: AdminService - 需要 AuditRepository, LogRepository, UserRepository

	return nil
}

// SetAuthService 设置认证服务
func (c *ServiceContainer) SetAuthService(service auth.AuthService) {
	c.authService = service
}

// SetWalletService 设置钱包服务
func (c *ServiceContainer) SetWalletService(service wallet.WalletService) {
	c.walletService = service
}

// SetRecommendationService 设置推荐服务
func (c *ServiceContainer) SetRecommendationService(service recommendation.RecommendationService) {
	c.recommendationService = service
}

// SetMessagingService 设置消息服务
func (c *ServiceContainer) SetMessagingService(service messaging.MessagingService) {
	c.messagingService = service
}

// SetStorageService 设置存储服务
func (c *ServiceContainer) SetStorageService(service storage.StorageService) {
	c.storageService = service
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
