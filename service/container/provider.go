package container

import (
	"fmt"
	"sync"
)

// Provider 服务提供者定义
// 用于声明式地注册和管理服务及其依赖关系
type Provider struct {
	// Name 服务名称，必须唯一
	Name string

	// Factory 服务工厂函数，用于创建服务实例
	// 返回值：服务实例和可能的错误
	Factory func(*ServiceContainer) (interface{}, error)

	// Dependencies 显式声明的依赖项
	// 列出此Provider依赖的其他Provider名称
	Dependencies []string

	// Singleton 是否为单例服务
	// true表示整个容器生命周期内只创建一次
	// false表示每次获取都创建新实例
	Singleton bool

	// Lazy 是否延迟初始化
	// true表示在首次使用时才初始化
	// false表示容器初始化时就创建
	Lazy bool
}

// ProviderRegistry Provider注册表
// 用于管理所有已注册的Provider
type ProviderRegistry struct {
	providers map[string]*Provider
	instances map[string]interface{} // 单例缓存
	mu        sync.RWMutex
}

// NewProviderRegistry 创建新的Provider注册表
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]*Provider),
		instances: make(map[string]interface{}),
	}
}

// Register 注册一个新的Provider
func (r *ProviderRegistry) Register(provider Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[provider.Name]; exists {
		return fmt.Errorf("Provider %s 已存在", provider.Name)
	}

	r.providers[provider.Name] = &provider
	return nil
}

// Get 获取指定名称的Provider
func (r *ProviderRegistry) Get(name string) (*Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	return provider, exists
}

// List 列出所有已注册的Provider
func (r *ProviderRegistry) List() []*Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Provider, 0, len(r.providers))
	for _, p := range r.providers {
		result = append(result, p)
	}
	return result
}

// GetOrCreate 获取或创建服务实例
func (r *ProviderRegistry) GetOrCreate(name string, container *ServiceContainer) (interface{}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("Provider %s 不存在", name)
	}

	// 如果是单例且已创建，直接返回
	if provider.Singleton {
		if instance, exists := r.instances[name]; exists {
			return instance, nil
		}
	}

	// 调用Factory创建实例
	instance, err := provider.Factory(container)
	if err != nil {
		return nil, fmt.Errorf("创建 %s 失败: %w", name, err)
	}

	if instance == nil {
		return nil, fmt.Errorf("Provider %s 返回nil实例", name)
	}

	// 缓存单例
	if provider.Singleton {
		r.instances[name] = instance
	}

	return instance, nil
}

// ResolutionOrder 计算Provider的解析顺序（拓扑排序）
func (r *ProviderRegistry) ResolutionOrder() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 构建依赖图
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// 初始化所有节点
	for name := range r.providers {
		graph[name] = []string{}
		inDegree[name] = 0
	}

	// 构建边
	for _, provider := range r.providers {
		for _, dep := range provider.Dependencies {
			if _, depExists := r.providers[dep]; depExists {
				graph[dep] = append(graph[dep], provider.Name)
				inDegree[provider.Name]++
			}
		}
	}

	// 拓扑排序
	order := []string{}
	queue := []string{}

	// 找到所有入度为0的节点
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	for len(queue) > 0 {
		// 取出一个节点
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		// 减少依赖此节点的其他节点的入度
		for _, dependent := range graph[current] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	// 检查是否有循环依赖
	if len(order) != len(r.providers) {
		return nil, fmt.Errorf("检测到循环依赖，已解析 %d/%d 个Provider", len(order), len(r.providers))
	}

	return order, nil
}

// DetectCircularDependency 检测循环依赖
func (r *ProviderRegistry) DetectCircularDependency() ([]string, error) {
	order, err := r.ResolutionOrder()
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetProvider 获取指定Provider的实例
// 这是一个辅助方法，用于从ProviderRegistry获取服务实例
func (c *ServiceContainer) GetProvider(name string) (interface{}, error) {
	if c.providerRegistry == nil {
		return nil, fmt.Errorf("ProviderRegistry未初始化")
	}
	return c.providerRegistry.GetOrCreate(name, c)
}

// RegisterProvider 注册Provider到容器
func (c *ServiceContainer) RegisterProvider(provider Provider) error {
	if c.providerRegistry == nil {
		c.providerRegistry = NewProviderRegistry()
	}
	return c.providerRegistry.Register(provider)
}

// ValidateProviders 验证所有Provider的配置
func (c *ServiceContainer) ValidateProviders() error {
	if c.providerRegistry == nil {
		return fmt.Errorf("ProviderRegistry未初始化")
	}

	// 验证依赖是否存在
	for _, provider := range c.providerRegistry.List() {
		for _, dep := range provider.Dependencies {
			if _, exists := c.providerRegistry.providers[dep]; !exists {
				return fmt.Errorf("Provider %s 依赖的 %s 不存在", provider.Name, dep)
			}
		}
	}

	// 检测循环依赖
	_, err := c.providerRegistry.ResolutionOrder()
	if err != nil {
		return err
	}

	return nil
}
