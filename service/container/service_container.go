package container

import (
	"context"
	"fmt"

	repoInterfaces "Qingyu_backend/repository/interfaces"
	aiRepoInterfaces "Qingyu_backend/repository/interfaces/ai"
	bookstoreRepoInterfaces "Qingyu_backend/repository/interfaces/bookstore"
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
	"Qingyu_backend/service/shared/recommendation"
	"Qingyu_backend/service/shared/storage"
	"Qingyu_backend/service/shared/wallet"
)

// ServiceContainer 服务容器
// 负责管理所有服务的生命周期和依赖注入
type ServiceContainer struct {
	repositoryFactory repoInterfaces.RepositoryFactory
	services          map[string]serviceInterfaces.BaseService
	initialized       bool

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
func NewServiceContainer(repositoryFactory repoInterfaces.RepositoryFactory) *ServiceContainer {
	return &ServiceContainer{
		repositoryFactory: repositoryFactory,
		services:          make(map[string]serviceInterfaces.BaseService),
		initialized:       false,
	}
}

// RegisterService 注册服务
func (c *ServiceContainer) RegisterService(name string, service serviceInterfaces.BaseService) error {
	if c.services[name] != nil {
		return fmt.Errorf("服务 %s 已存在", name)
	}

	c.services[name] = service
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

// Initialize 初始化所有服务
func (c *ServiceContainer) Initialize(ctx context.Context) error {
	if c.initialized {
		return nil
	}

	// 初始化Repository工厂
	if err := c.repositoryFactory.Health(ctx); err != nil {
		return fmt.Errorf("Repository工厂健康检查失败: %w", err)
	}

	// 初始化所有服务
	for name, service := range c.services {
		if err := service.Initialize(ctx); err != nil {
			return fmt.Errorf("初始化服务 %s 失败: %w", name, err)
		}
	}

	c.initialized = true
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

	// 关闭所有服务
	for name, service := range c.services {
		if err := service.Close(ctx); err != nil {
			lastErr = fmt.Errorf("关闭服务 %s 失败: %w", name, err)
		}
	}

	// 关闭Repository工厂
	if err := c.repositoryFactory.Close(); err != nil {
		lastErr = fmt.Errorf("关闭Repository工厂失败: %w", err)
	}

	c.initialized = false
	return lastErr
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
	// 需要进行类型断言
	bookRepoRaw := c.repositoryFactory.CreateBookRepository()
	categoryRepoRaw := c.repositoryFactory.CreateCategoryRepository()
	bannerRepoRaw := c.repositoryFactory.CreateBannerRepository()
	rankingRepoRaw := c.repositoryFactory.CreateRankingRepository()

	// 类型断言
	bookRepo, ok := bookRepoRaw.(bookstoreRepoInterfaces.BookRepository)
	if !ok {
		return fmt.Errorf("BookRepository类型转换失败")
	}
	categoryRepo, ok := categoryRepoRaw.(bookstoreRepoInterfaces.CategoryRepository)
	if !ok {
		return fmt.Errorf("CategoryRepository类型转换失败")
	}
	bannerRepo, ok := bannerRepoRaw.(bookstoreRepoInterfaces.BannerRepository)
	if !ok {
		return fmt.Errorf("BannerRepository类型转换失败")
	}
	rankingRepo, ok := rankingRepoRaw.(bookstoreRepoInterfaces.RankingRepository)
	if !ok {
		return fmt.Errorf("RankingRepository类型转换失败")
	}

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
		nil, // eventBus - TODO: 实现事件总线
		nil, // cacheService - TODO: 实现缓存服务
		nil, // vipService - TODO: 实现VIP服务
	)
	// 注意：ReaderService 不完全实现 BaseService，不注册到 services map

	// ============ 4. 创建AI服务 ============
	c.aiService = aiService.NewService()

	// 创建AI配额服务
	quotaRepoRaw := c.repositoryFactory.CreateQuotaRepository()
	quotaRepo, ok := quotaRepoRaw.(aiRepoInterfaces.QuotaRepository)
	if !ok {
		return fmt.Errorf("QuotaRepository类型转换失败")
	}
	c.quotaService = aiService.NewQuotaService(quotaRepo)
	// 注意：QuotaService 不完全实现 BaseService，不注册到 services map

	// 创建聊天服务（使用临时内存存储）
	chatRepo := aiService.NewInMemoryChatRepository()
	c.chatService = aiService.NewChatService(c.aiService, chatRepo)
	// 注意：ChatService 不完全实现 BaseService，不注册到 services map

	// ============ 5. 共享服务（可选，未完全实现） ============
	// 注意：共享服务尚未完全实现，暂时设置为nil
	// TODO: 实现并注册共享服务

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
