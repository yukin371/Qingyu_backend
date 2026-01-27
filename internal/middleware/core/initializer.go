package core

import "Qingyu_backend/internal/middleware"

// Initializer 中间件初始化器接口
//
// 职责：
//   1. 从配置文件加载中间件配置
//   2. 创建和初始化中间件实例
//   3. 提供中间件实例查询接口
//
// 不负责：
//   - 中间件的存储和管理（由Manager负责）
type Initializer interface {
	// LoadFromConfig 从配置文件加载中间件配置
	//
	// configPath: YAML配置文件路径，如 "configs/middleware.yaml"
	// 返回: 加载的配置对象
	// 错误: 文件不存在、格式错误、验证失败等
	LoadFromConfig(configPath string) (*middleware.Config, error)

	// Initialize 初始化所有中间件
	//
	// 根据加载的配置创建中间件实例。
	// 返回: 中间件实例列表
	// 错误: 初始化失败、配置验证失败等
	Initialize() ([]Middleware, error)

	// GetMiddleware 获取指定中间件实例
	//
	// name: 中间件名称，如 "rate_limit", "request_id"
	// 返回: 中间件实例
	// 错误: 中间件不存在或未初始化
	GetMiddleware(name string) (Middleware, error)

	// ListMiddlewares 列出所有已初始化的中间件
	//
	// 返回: 中间件名称列表
	ListMiddlewares() []string
}
