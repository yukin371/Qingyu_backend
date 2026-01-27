package core

import (
	"fmt"
	"sort"
	"sync"

	"go.uber.org/zap"
)

// Registry 中间件注册器
//
// 负责中间件的注册、存储和按优先级排序。
type Registry struct {
	middlewares map[string]Middleware
	priorities  map[string]int // 覆盖的优先级配置
	mu          sync.RWMutex
	logger      *zap.Logger
}

// NewRegistry 创建注册器
func NewRegistry(logger *zap.Logger) *Registry {
	return &Registry{
		middlewares: make(map[string]Middleware),
		priorities:  make(map[string]int),
		logger:      logger,
	}
}

// Register 注册中间件
//
// 如果中间件已存在，返回错误。
// 如果存在优先级覆盖配置，使用覆盖值。
func (r *Registry) Register(mw Middleware) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := mw.Name()

	// 检查是否已注册
	if _, exists := r.middlewares[name]; exists {
		return fmt.Errorf("middleware %s already registered", name)
	}

	// 记录日志
	r.logger.Info("Registering middleware",
		zap.String("name", name),
		zap.Int("priority", mw.Priority()))

	r.middlewares[name] = mw
	return nil
}

// Unregister 注销中间件
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.middlewares[name]; !exists {
		return fmt.Errorf("middleware %s not found", name)
	}

	delete(r.middlewares, name)
	r.logger.Info("Unregistered middleware", zap.String("name", name))
	return nil
}

// Get 获取中间件
func (r *Registry) Get(name string) (Middleware, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mw, exists := r.middlewares[name]
	if !exists {
		return nil, fmt.Errorf("middleware %s not found", name)
	}

	return mw, nil
}

// List 列出所有已注册中间件的名称
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.middlewares))
	for name := range r.middlewares {
		names = append(names, name)
	}

	return names
}

// Sorted 返回按优先级排序的中间件列表
//
// 优先级值越小，越先执行。
func (r *Registry) Sorted() []Middleware {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 复制到slice避免修改原map
	middlewareList := make([]Middleware, 0, len(r.middlewares))
	for _, mw := range r.middlewares {
		middlewareList = append(middlewareList, mw)
	}

	// 按优先级排序
	sort.Slice(middlewareList, func(i, j int) bool {
		priorityI := r.getPriority(middlewareList[i])
		priorityJ := r.getPriority(middlewareList[j])
		return priorityI < priorityJ
	})

	return middlewareList
}

// SetPriorityOverride 设置优先级覆盖
//
// 允许通过配置文件覆盖中间件的默认优先级。
func (r *Registry) SetPriorityOverride(name string, priority int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.priorities[name] = priority
	r.logger.Info("Priority override set",
		zap.String("middleware", name),
		zap.Int("override_priority", priority))
}

// getPriority 获取中间件的优先级（考虑覆盖）
func (r *Registry) getPriority(mw Middleware) int {
	name := mw.Name()

	// 检查是否有优先级覆盖
	if override, exists := r.priorities[name]; exists {
		return override
	}

	// 返回默认优先级
	return mw.Priority()
}

// Validate 验证中间件配置
//
// 检查：
// - 中间件数量是否合理
// - 优先级是否合理
// - 是否存在循环依赖（未来扩展）
func (r *Registry) Validate() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 检查中间件数量
	if len(r.middlewares) == 0 {
		return fmt.Errorf("no middleware registered")
	}

	if len(r.middlewares) > 20 {
		r.logger.Warn("Too many middlewares registered",
			zap.Int("count", len(r.middlewares)))
	}

	// 检查优先级范围
	sorted := r.Sorted()
	for _, mw := range sorted {
		priority := r.getPriority(mw)
		if priority < 1 || priority > 100 {
			return fmt.Errorf("invalid priority %d for middleware %s", priority, mw.Name())
		}
	}

	// 检查优先级唯一性（警告）
	priorities := make(map[int]bool)
	for _, mw := range sorted {
		priority := r.getPriority(mw)
		if priorities[priority] {
			r.logger.Warn("Duplicate priority detected",
				zap.Int("priority", priority),
				zap.String("middleware", mw.Name()))
		}
		priorities[priority] = true
	}

	return nil
}

// Count 返回已注册中间件数量
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.middlewares)
}

// Clear 清空所有中间件
//
// 主要用于测试。
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.middlewares = make(map[string]Middleware)
	r.priorities = make(map[string]int)
	r.logger.Info("Registry cleared")
}
