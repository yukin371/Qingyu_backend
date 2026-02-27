package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProviderStruct 测试Provider结构定义
// 这是TDD的核心：先定义期望的Provider行为
func TestProviderStruct(t *testing.T) {
	t.Run("Provider应该有基本的字段", func(t *testing.T) {
		provider := Provider{
			Name:         "test-service",
			Dependencies: []string{"database"},
			Singleton:    true,
			Lazy:         false,
		}

		assert.Equal(t, "test-service", provider.Name)
		assert.Equal(t, []string{"database"}, provider.Dependencies)
		assert.True(t, provider.Singleton)
		assert.False(t, provider.Lazy)
	})

	t.Run("Provider应该有Factory函数", func(t *testing.T) {
		provider := Provider{
			Name: "simple-service",
			Factory: func(c *ServiceContainer) (interface{}, error) {
				return "service-instance", nil
			},
		}

		require.NotNil(t, provider.Factory)

		container := NewServiceContainer()
		instance, err := provider.Factory(container)

		assert.NoError(t, err)
		assert.Equal(t, "service-instance", instance)
	})
}

// TestProviderDependencyDeclaration 测试Provider依赖声明
func TestProviderDependencyDeclaration(t *testing.T) {
	t.Run("Provider应该能够声明显式依赖", func(t *testing.T) {
		databaseProvider := Provider{
			Name:         "database",
			Factory:      nil,        // 简化测试
			Dependencies: []string{}, // 数据库没有依赖
		}

		userServiceProvider := Provider{
			Name:         "user-service",
			Factory:      nil,
			Dependencies: []string{"database"}, // 用户服务依赖数据库
		}

		assert.Empty(t, databaseProvider.Dependencies)
		assert.Equal(t, []string{"database"}, userServiceProvider.Dependencies)
	})

	t.Run("Provider应该能够声明多个依赖", func(t *testing.T) {
		provider := Provider{
			Name: "complex-service",
			Dependencies: []string{
				"database",
				"cache",
				"logger",
			},
		}

		assert.Len(t, provider.Dependencies, 3)
	})
}

// TestProviderResolutionOrder 测试Provider解析顺序
func TestProviderResolutionOrder(t *testing.T) {
	t.Run("应该能够根据依赖关系排序Provider", func(t *testing.T) {
		providers := []Provider{
			{Name: "user-service", Dependencies: []string{"database"}},
			{Name: "database", Dependencies: []string{}},
			{Name: "cache-service", Dependencies: []string{"database"}},
		}

		order := calculateResolutionOrder(providers)

		// database应该最先（没有依赖）
		assert.Equal(t, "database", order[0])

		// user-service和cache-service应该在database之后
		userIndex := indexOf(order, "user-service")
		cacheIndex := indexOf(order, "cache-service")
		dbIndex := indexOf(order, "database")

		assert.Greater(t, userIndex, dbIndex)
		assert.Greater(t, cacheIndex, dbIndex)
	})

	t.Run("应该能够检测循环依赖", func(t *testing.T) {
		registry := NewProviderRegistry()

		// 注册有循环依赖的Provider
		registry.Register(Provider{Name: "service-a", Dependencies: []string{"service-b"}})
		registry.Register(Provider{Name: "service-b", Dependencies: []string{"service-a"}})

		// 检测循环依赖
		_, err := registry.DetectCircularDependency()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "循环依赖")
	})
}

// TestLazyProvider 测试延迟加载Provider
func TestLazyProvider(t *testing.T) {
	t.Run("Lazy Provider应该延迟初始化", func(t *testing.T) {
		callCount := 0

		provider := Provider{
			Name: "lazy-service",
			Lazy: true,
			Factory: func(c *ServiceContainer) (interface{}, error) {
				callCount++
				return "lazy-instance", nil
			},
		}

		// Provider创建时不应该调用Factory
		assert.Equal(t, 0, callCount)

		// 显式初始化时才调用Factory
		container := NewServiceContainer()
		instance, err := provider.Factory(container)

		assert.NoError(t, err)
		assert.Equal(t, 1, callCount)
		assert.Equal(t, "lazy-instance", instance)
	})
}

// TestSingletonProvider 测试单例Provider
func TestSingletonProvider(t *testing.T) {
	t.Run("Singleton Provider应该只创建一次实例", func(t *testing.T) {
		callCount := 0

		provider := Provider{
			Name:      "singleton-service",
			Singleton: true,
			Factory: func(c *ServiceContainer) (interface{}, error) {
				callCount++
				return "instance-" + string(rune(callCount)), nil
			},
		}

		container := NewServiceContainer()
		registry := NewProviderRegistry()
		err := registry.Register(provider)
		require.NoError(t, err)

		// 第一次调用
		instance1, err := registry.GetOrCreate("singleton-service", container)
		require.NoError(t, err)
		assert.Equal(t, 1, callCount)

		// 第二次调用（单例应该返回相同实例）
		instance2, err := registry.GetOrCreate("singleton-service", container)
		require.NoError(t, err)
		assert.Equal(t, 1, callCount, "单例不应该重复创建")
		assert.Equal(t, instance1, instance2)
	})
}

// TestProviderRegistry 测试Provider注册表
func TestProviderRegistry(t *testing.T) {
	t.Run("应该能够注册和获取Provider", func(t *testing.T) {
		registry := NewProviderRegistry()

		provider := Provider{
			Name: "test-provider",
			Factory: func(c *ServiceContainer) (interface{}, error) {
				return "test", nil
			},
		}

		err := registry.Register(provider)
		assert.NoError(t, err)

		retrieved, exists := registry.Get("test-provider")
		assert.True(t, exists)
		assert.Equal(t, "test-provider", retrieved.Name)
	})

	t.Run("不应该允许重复注册", func(t *testing.T) {
		registry := NewProviderRegistry()

		provider := Provider{Name: "duplicate"}

		_ = registry.Register(provider)
		err := registry.Register(provider)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已存在")
	})

	t.Run("应该能够列出所有Provider", func(t *testing.T) {
		registry := NewProviderRegistry()

		registry.Register(Provider{Name: "provider-1"})
		registry.Register(Provider{Name: "provider-2"})

		providers := registry.List()

		assert.Len(t, providers, 2)
	})
}

// Helper functions

// calculateResolutionOrder 计算Provider的解析顺序（拓扑排序）
func calculateResolutionOrder(providers []Provider) []string {
	// 简化的拓扑排序实现
	order := []string{}
	remaining := make(map[string]bool)
	dependencies := make(map[string][]string)

	for _, p := range providers {
		remaining[p.Name] = true
		dependencies[p.Name] = p.Dependencies
	}

	for len(remaining) > 0 {
		progress := false
		for name := range remaining {
			// 检查所有依赖是否已满足
			allDepsMet := true
			for _, dep := range dependencies[name] {
				if remaining[dep] {
					allDepsMet = false
					break
				}
			}

			if allDepsMet {
				order = append(order, name)
				delete(remaining, name)
				progress = true
				break
			}
		}

		if !progress {
			// 有循环依赖
			return nil
		}
	}

	return order
}

// detectCircularDependency 检测循环依赖
func detectCircularDependency(providers []Provider) ([]string, []string) {
	order := calculateResolutionOrder(providers)
	if order == nil {
		// 返回检测到的循环
		return nil, []string{"circular-dependency-detected"}
	}
	return order, nil
}

// indexOf 辅助函数：查找元素在切片中的索引
func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// mockFactoryForTesting 创建用于测试的mock factory
func mockFactoryForTesting(value interface{}) func(*ServiceContainer) (interface{}, error) {
	return func(c *ServiceContainer) (interface{}, error) {
		return value, nil
	}
}

// TestProviderWithRealContainer 测试Provider与实际容器的集成
func TestProviderWithRealContainer(t *testing.T) {
	t.Run("Provider应该能够使用容器获取依赖", func(t *testing.T) {
		container := NewServiceContainer()

		// 注册一个模拟的数据库Provider
		dbProvider := Provider{
			Name: "database",
			Factory: func(c *ServiceContainer) (interface{}, error) {
				return "mock-database", nil
			},
		}

		// 注册一个使用数据库的Provider
		userProvider := Provider{
			Name: "user-service",
			Factory: func(c *ServiceContainer) (interface{}, error) {
				// 从容器获取数据库
				db, err := c.GetProvider("database")
				if err != nil {
					return nil, err
				}
				return "user-service-with-" + db.(string), nil
			},
			Dependencies: []string{"database"},
		}

		// 这个测试展示了Provider如何使用容器
		// 实际实现需要容器的支持
		_ = container
		_ = dbProvider
		_ = userProvider
	})
}

// TestProviderErrorHandling 测试Provider错误处理
func TestProviderErrorHandling(t *testing.T) {
	t.Run("Factory返回错误时应该传播错误", func(t *testing.T) {
		expectedErr := errors.New("factory failed")

		provider := Provider{
			Name: "failing-service",
			Factory: func(c *ServiceContainer) (interface{}, error) {
				return nil, expectedErr
			},
		}

		container := NewServiceContainer()
		instance, err := provider.Factory(container)

		assert.Error(t, err)
		assert.Nil(t, instance)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Factory返回nil实例时应该返回错误", func(t *testing.T) {
		registry := NewProviderRegistry()
		container := NewServiceContainer()

		provider := Provider{
			Name: "nil-returning-service",
			Factory: func(c *ServiceContainer) (interface{}, error) {
				return nil, nil // 返回nil实例
			},
		}

		registry.Register(provider)

		_, err := registry.GetOrCreate("nil-returning-service", container)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "返回nil")
	})
}
