package service

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/user"
	"context"
	"fmt"

	repoInterfaces "Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/container"
)

// ServiceManager 服务管理器
// 全局服务管理器，负责管理所有服务的生命周期
var ServiceManager *container.ServiceContainer

// InitializeServices 初始化所有服务
func InitializeServices(repositoryFactory repoInterfaces.RepositoryFactory) error {
	// 创建服务容器
	ServiceManager = container.NewServiceContainer(repositoryFactory)

	// 设置默认服务
	if err := ServiceManager.SetupDefaultServices(); err != nil {
		return err
	}

	// 初始化所有服务
	if err := ServiceManager.Initialize(context.Background()); err != nil {
		return err
	}

	return nil
}

// GetUserService 获取用户服务
func GetUserService() (serviceInterfaces.UserService, error) {
	if ServiceManager == nil {
		return nil, fmt.Errorf("服务管理器未初始化")
	}

	return ServiceManager.GetUserService()
}

// GetServiceContainer 获取服务容器
func GetServiceContainer() *container.ServiceContainer {
	return ServiceManager
}

// CloseServices 关闭所有服务
func CloseServices(ctx context.Context) error {
	if ServiceManager == nil {
		return nil
	}

	return ServiceManager.Close(ctx)
}

// HealthCheck 健康检查
func HealthCheck(ctx context.Context) error {
	if ServiceManager == nil {
		return fmt.Errorf("服务管理器未初始化")
	}

	return ServiceManager.Health(ctx)
}
