package core

import (
	"fmt"

	"Qingyu_backend/service"
)

// InitDB 初始化数据库连接
// Deprecated: MongoDB初始化已迁移到ServiceContainer
// 保留此函数仅用于向后兼容，实际初始化在ServiceContainer.Initialize()中完成
func InitDB() error {
	// MongoDB初始化已迁移到ServiceContainer
	// 此函数保留为空，避免破坏现有调用
	// 实际的MongoDB连接在service.InitializeServices()中由ServiceContainer自动创建

	// 可选：为了兼容性，设置全局变量指向ServiceContainer的连接
	// 但这需要在InitServices之后调用，所以这里暂时返回nil

	fmt.Println("InitDB: MongoDB初始化已迁移到ServiceContainer")
	return nil
}

// InitServices 初始化所有服务
// ServiceContainer会自动创建MongoDB连接和Repository工厂
func InitServices() error {
	// 直接初始化服务容器
	// ServiceContainer.Initialize()会自动：
	// 1. 创建MongoDB连接
	// 2. 创建Repository工厂
	// 3. 初始化所有服务
	if err := service.InitializeServices(); err != nil {
		return fmt.Errorf("初始化服务失败: %w", err)
	}

	fmt.Println("Successfully initialized all services")
	return nil
}
