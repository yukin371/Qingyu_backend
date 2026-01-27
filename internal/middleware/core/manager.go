package core

import (
	"github.com/gin-gonic/gin"
)

// Manager 中间件管理器接口
//
// 负责中间件的注册、注销、排序和应用到Gin路由。
type Manager interface {
	// Register 注册中间件
	//
	// 支持通过RegisterOption配置中间件（如优先级覆盖）。
	// 如果中间件已存在，返回错误。
	Register(middleware Middleware, options ...RegisterOption) error

	// Unregister 注销中间件
	//
	// 从管理器中移除中间件。
	// 如果中间件不存在，返回错误。
	Unregister(name string) error

	// Get 获取中间件
	//
	// 根据名称获取已注册的中间件。
	// 如果中间件不存在，返回错误。
	Get(name string) (Middleware, error)

	// List 列出所有中间件
	//
	// 返回所有已注册的中间件，按优先级排序。
	List() []Middleware

	// Build 构建中间件链
	//
	// 返回Gin中间件函数列表，按优先级排序。
	// 可以直接用于Gin路由的Use()方法。
	Build() []gin.HandlerFunc

	// ApplyToRouter 应用到Gin路由
	//
	// 将中间件应用到Gin引擎。
	// globalMiddlewares指定要应用的全局中间件名称（为空则应用所有）。
	ApplyToRouter(router *gin.Engine, globalMiddlewares ...string) error
}

// RegisterOption 中间件注册选项
//
// 用于配置中间件注册行为。
type RegisterOption interface {
	apply(*registerConfig)
}

// registerConfig 注册配置
type registerConfig struct {
	priorityOverride *int
}

// funcRegisterOption 函数式选项
type funcRegisterOption func(*registerConfig)

func (f funcRegisterOption) apply(cfg *registerConfig) {
	f(cfg)
}

// WithPriority 设置优先级覆盖
//
// 允许在注册时覆盖中间件的默认优先级。
// 示例：
//   manager.Register(middleware, WithPriority(5))
func WithPriority(priority int) RegisterOption {
	return funcRegisterOption(func(cfg *registerConfig) {
		cfg.priorityOverride = &priority
	})
}
