package container

import (
	"context"
	"fmt"

	repoInterfaces "Qingyu_backend/repository/interfaces"
	serviceInterfaces "Qingyu_backend/service/interfaces"
	"Qingyu_backend/service/user"
)

// ServiceContainer 服务容器
// 负责管理所有服务的生命周期和依赖注入
type ServiceContainer struct {
	repositoryFactory repoInterfaces.RepositoryFactory
	services          map[string]serviceInterfaces.BaseService
	initialized       bool
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

// GetUserService 获取用户服务
func (c *ServiceContainer) GetUserService() (serviceInterfaces.UserService, error) {
	service, err := c.GetService("UserService")
	if err != nil {
		return nil, err
	}

	userService, ok := service.(serviceInterfaces.UserService)
	if !ok {
		return nil, fmt.Errorf("服务类型转换失败")
	}

	return userService, nil
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
func (c *ServiceContainer) SetupDefaultServices() error {
	// 创建用户服务
	userRepo := c.repositoryFactory.CreateUserRepository()
	userService := user.NewUserService(userRepo)

	// 注册用户服务
	if err := c.RegisterService("UserService", userService); err != nil {
		return fmt.Errorf("注册用户服务失败: %w", err)
	}

	return nil
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
