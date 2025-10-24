package container

import (
	"context"
	"fmt"

	"Qingyu_backend/service/shared/admin"
	"Qingyu_backend/service/shared/auth"
	"Qingyu_backend/service/shared/messaging"
	"Qingyu_backend/service/shared/recommendation"
	"Qingyu_backend/service/shared/storage"
	"Qingyu_backend/service/shared/wallet"
)

// DEPRECATED: SharedServiceContainer 已废弃
// 请使用 service/container/ServiceContainer 统一管理所有服务（包括共享服务）
// 此容器将在下一版本中删除
//
// 迁移指南:
// 1. 使用 service.GetServiceContainer() 获取全局服务容器
// 2. 通过 container.GetAuthService(), container.GetWalletService() 等方法获取共享服务
// 3. 不再需要单独创建 SharedServiceContainer 实例
//
// SharedServiceContainer 共享服务容器
// 统一管理所有共享底层服务
type SharedServiceContainer struct {
	// 核心服务
	authService           auth.AuthService
	walletService         wallet.WalletService
	recommendationService recommendation.RecommendationService
	messagingService      messaging.MessagingService
	storageService        storage.StorageService
	adminService          admin.AdminService

	// 初始化状态
	initialized bool
}

// SharedServiceConfig 共享服务配置
type SharedServiceConfig struct {
	// Auth配置
	JWTSecret    string
	JWTExpiresIn int64

	// Storage配置
	StorageBasePath string
	StorageBaseURL  string

	// Messaging配置
	RedisAddr string
	RedisPass string

	// 其他配置...
}

// NewSharedServiceContainer 创建共享服务容器
// DEPRECATED: 此方法已废弃，请使用 service.GetServiceContainer() 获取全局服务容器
func NewSharedServiceContainer() *SharedServiceContainer {
	return &SharedServiceContainer{
		initialized: false,
	}
}

// SetAuthService 设置认证服务
func (c *SharedServiceContainer) SetAuthService(service auth.AuthService) {
	c.authService = service
}

// SetWalletService 设置钱包服务
func (c *SharedServiceContainer) SetWalletService(service wallet.WalletService) {
	c.walletService = service
}

// SetRecommendationService 设置推荐服务
func (c *SharedServiceContainer) SetRecommendationService(service recommendation.RecommendationService) {
	c.recommendationService = service
}

// SetMessagingService 设置消息服务
func (c *SharedServiceContainer) SetMessagingService(service messaging.MessagingService) {
	c.messagingService = service
}

// SetStorageService 设置存储服务
func (c *SharedServiceContainer) SetStorageService(service storage.StorageService) {
	c.storageService = service
}

// SetAdminService 设置管理服务
func (c *SharedServiceContainer) SetAdminService(service admin.AdminService) {
	c.adminService = service
}

// ============ 获取服务 ============

// AuthService 获取认证服务
func (c *SharedServiceContainer) AuthService() auth.AuthService {
	return c.authService
}

// WalletService 获取钱包服务
func (c *SharedServiceContainer) WalletService() wallet.WalletService {
	return c.walletService
}

// RecommendationService 获取推荐服务
func (c *SharedServiceContainer) RecommendationService() recommendation.RecommendationService {
	return c.recommendationService
}

// MessagingService 获取消息服务
func (c *SharedServiceContainer) MessagingService() messaging.MessagingService {
	return c.messagingService
}

// StorageService 获取存储服务
func (c *SharedServiceContainer) StorageService() storage.StorageService {
	return c.storageService
}

// AdminService 获取管理服务
func (c *SharedServiceContainer) AdminService() admin.AdminService {
	return c.adminService
}

// ============ 生命周期管理 ============

// Initialize 初始化所有服务
func (c *SharedServiceContainer) Initialize(ctx context.Context) error {
	if c.initialized {
		return nil
	}

	// 检查所有服务是否已注册
	if c.authService == nil {
		return fmt.Errorf("Auth服务未注册")
	}
	if c.walletService == nil {
		return fmt.Errorf("Wallet服务未注册")
	}
	if c.recommendationService == nil {
		return fmt.Errorf("Recommendation服务未注册")
	}
	if c.messagingService == nil {
		return fmt.Errorf("Messaging服务未注册")
	}
	if c.storageService == nil {
		return fmt.Errorf("Storage服务未注册")
	}
	if c.adminService == nil {
		return fmt.Errorf("Admin服务未注册")
	}

	// 执行健康检查
	if err := c.Health(ctx); err != nil {
		return fmt.Errorf("服务健康检查失败: %w", err)
	}

	c.initialized = true
	return nil
}

// Health 健康检查
func (c *SharedServiceContainer) Health(ctx context.Context) error {
	// 检查所有服务健康状态
	services := map[string]interface {
		Health(context.Context) error
	}{
		"Auth":           c.authService,
		"Wallet":         c.walletService,
		"Recommendation": c.recommendationService,
		"Messaging":      c.messagingService,
		"Storage":        c.storageService,
		"Admin":          c.adminService,
	}

	for name, service := range services {
		if service == nil {
			return fmt.Errorf("服务 %s 未初始化", name)
		}
		if err := service.Health(ctx); err != nil {
			return fmt.Errorf("服务 %s 健康检查失败: %w", name, err)
		}
	}

	return nil
}

// IsInitialized 检查是否已初始化
func (c *SharedServiceContainer) IsInitialized() bool {
	return c.initialized
}

// GetServiceStatus 获取服务状态
func (c *SharedServiceContainer) GetServiceStatus(ctx context.Context) map[string]string {
	status := make(map[string]string)

	services := map[string]interface {
		Health(context.Context) error
	}{
		"Auth":           c.authService,
		"Wallet":         c.walletService,
		"Recommendation": c.recommendationService,
		"Messaging":      c.messagingService,
		"Storage":        c.storageService,
		"Admin":          c.adminService,
	}

	for name, service := range services {
		if service == nil {
			status[name] = "not_registered"
			continue
		}
		if err := service.Health(ctx); err != nil {
			status[name] = fmt.Sprintf("unhealthy: %v", err)
		} else {
			status[name] = "healthy"
		}
	}

	return status
}
