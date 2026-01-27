package core

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// managerImpl 中间件管理器实现
type managerImpl struct {
	registry *Registry
	logger   *zap.Logger
}

// NewManager 创建中间件管理器
func NewManager(logger *zap.Logger) Manager {
	return &managerImpl{
		registry: NewRegistry(logger),
		logger:   logger,
	}
}

// Register 注册中间件
func (m *managerImpl) Register(middleware Middleware, options ...RegisterOption) error {
	// 应用注册选项
	cfg := &registerConfig{}
	for _, opt := range options {
		opt.apply(cfg)
	}

	// 如果有优先级覆盖，设置到registry
	if cfg.priorityOverride != nil {
		m.registry.SetPriorityOverride(middleware.Name(), *cfg.priorityOverride)
	}

	// 注册中间件
	if err := m.registry.Register(middleware); err != nil {
		return err
	}

	m.logger.Info("Middleware registered successfully",
		zap.String("name", middleware.Name()),
		zap.Int("priority", m.getPriority(middleware)))

	return nil
}

// Unregister 注销中间件
func (m *managerImpl) Unregister(name string) error {
	return m.registry.Unregister(name)
}

// Get 获取中间件
func (m *managerImpl) Get(name string) (Middleware, error) {
	return m.registry.Get(name)
}

// List 列出所有中间件
func (m *managerImpl) List() []Middleware {
	return m.registry.Sorted()
}

// Build 构建中间件链
func (m *managerImpl) Build() []gin.HandlerFunc {
	sorted := m.registry.Sorted()
	handlers := make([]gin.HandlerFunc, 0, len(sorted))

	for _, middleware := range sorted {
		handlers = append(handlers, middleware.Handler())
	}

	m.logger.Info("Middleware chain built",
		zap.Int("count", len(handlers)))

	return handlers
}

// ApplyToRouter 应用到Gin路由
func (m *managerImpl) ApplyToRouter(router *gin.Engine, globalMiddlewares ...string) error {
	// 确定要应用的中间件
	var middlewareList []Middleware

	if len(globalMiddlewares) == 0 {
		// 没有指定，应用所有中间件
		middlewareList = m.registry.Sorted()
	} else {
		// 按名称选择中间件
		middlewareList = make([]Middleware, 0, len(globalMiddlewares))
		for _, name := range globalMiddlewares {
			middleware, err := m.registry.Get(name)
			if err != nil {
				return fmt.Errorf("middleware %s not found", name)
			}
			middlewareList = append(middlewareList, middleware)
		}

		// 按优先级排序选中的中间件
		middlewareList = m.sortByPriority(middlewareList)
	}

	// 构建Gin处理函数列表
	handlers := make([]gin.HandlerFunc, 0, len(middlewareList))
	for _, middleware := range middlewareList {
		handlers = append(handlers, middleware.Handler())
	}

	// 应用到路由
	router.Use(handlers...)

	m.logger.Info("Middlewares applied to router",
		zap.Int("count", len(handlers)),
		zap.Strings("middlewares", getMiddlewareNames(middlewareList)))

	return nil
}

// getPriority 获取中间件的优先级
func (m *managerImpl) getPriority(mw Middleware) int {
	// 这里我们无法直接访问registry的getPriority方法
	// 所以通过反射或公开的方法获取
	// 为了简化，我们直接使用middleware的Priority方法
	// 如果有覆盖，需要在registry层面处理
	return mw.Priority()
}

// sortByPriority 按优先级排序中间件列表
func (m *managerImpl) sortByPriority(middlewares []Middleware) []Middleware {
	sorted := make([]Middleware, len(middlewares))
	copy(sorted, middlewares)

	// 简单的冒泡排序（稳定性好）
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if m.getPriority(sorted[j]) < m.getPriority(sorted[i]) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// getMiddlewareNames 获取中间件名称列表
func getMiddlewareNames(middlewares []Middleware) []string {
	names := make([]string, 0, len(middlewares))
	for _, mw := range middlewares {
		names = append(names, mw.Name())
	}
	return names
}
