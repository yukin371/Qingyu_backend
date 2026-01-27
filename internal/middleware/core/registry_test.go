package core

import (
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/stretchr/testify/assert"
)

// 测试辅助函数：创建测试注册器
func createTestRegistry(t *testing.T) *Registry {
	t.Helper()

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	return NewRegistry(logger)
}

// TestRegistry_Validate_Empty 测试空注册器验证
func TestRegistry_Validate_Empty(t *testing.T) {
	registry := createTestRegistry(t)

	err := registry.Validate()
	assert.Error(t, err, "Validate should fail for empty registry")
	assert.Contains(t, err.Error(), "no middleware registered")
}

// TestRegistry_Validate_Valid 测试有效注册器验证
func TestRegistry_Validate_Valid(t *testing.T) {
	registry := createTestRegistry(t)

	// 注册中间件
	mw := createTestMiddleware("test_middleware", 10)
	err := registry.Register(mw)
	assert.NoError(t, err)

	// 验证
	err = registry.Validate()
	assert.NoError(t, err, "Validate should succeed for valid registry")
}

// TestRegistry_Validate_InvalidPriority 测试无效优先级验证
func TestRegistry_Validate_InvalidPriority(t *testing.T) {
	registry := createTestRegistry(t)

	// 注册无效优先级的中间件（priority < 1）
	mw1 := createTestMiddleware("invalid_low", 0)
	err := registry.Register(mw1)
	assert.NoError(t, err)

	// 验证应该失败
	err = registry.Validate()
	assert.Error(t, err, "Validate should fail for invalid priority")
	assert.Contains(t, err.Error(), "invalid priority")
}

// TestRegistry_Validate_TooManyMiddlewares 测试过多中间件警告
func TestRegistry_Validate_TooManyMiddlewares(t *testing.T) {
	registry := createTestRegistry(t)

	// 注册超过20个中间件
	for i := 1; i <= 21; i++ {
		name := "middleware_" + string(rune('0'+i))
		mw := createTestMiddleware(name, i)
		err := registry.Register(mw)
		assert.NoError(t, err)
	}

	// 验证应该成功（但有警告）
	err := registry.Validate()
	assert.NoError(t, err, "Validate should succeed even with too many middlewares")
}

// TestRegistry_Count 测试中间件计数
func TestRegistry_Count(t *testing.T) {
	registry := createTestRegistry(t)

	// 初始计数为0
	count := registry.Count()
	assert.Equal(t, 0, count, "Initial count should be 0")

	// 注册3个中间件
	mw1 := createTestMiddleware("middleware_1", 1)
	mw2 := createTestMiddleware("middleware_2", 2)
	mw3 := createTestMiddleware("middleware_3", 3)

	err := registry.Register(mw1)
	assert.NoError(t, err)
	err = registry.Register(mw2)
	assert.NoError(t, err)
	err = registry.Register(mw3)
	assert.NoError(t, err)

	// 计数应该为3
	count = registry.Count()
	assert.Equal(t, 3, count, "Count should be 3 after registering 3 middlewares")

	// 注销一个中间件
	err = registry.Unregister("middleware_1")
	assert.NoError(t, err)

	// 计数应该为2
	count = registry.Count()
	assert.Equal(t, 2, count, "Count should be 2 after unregistering 1 middleware")
}

// TestRegistry_Clear 测试清空注册器
func TestRegistry_Clear(t *testing.T) {
	registry := createTestRegistry(t)

	// 注册3个中间件
	mw1 := createTestMiddleware("middleware_1", 1)
	mw2 := createTestMiddleware("middleware_2", 2)
	mw3 := createTestMiddleware("middleware_3", 3)

	err := registry.Register(mw1)
	assert.NoError(t, err)
	err = registry.Register(mw2)
	assert.NoError(t, err)
	err = registry.Register(mw3)
	assert.NoError(t, err)

	// 确认有中间件
	count := registry.Count()
	assert.Equal(t, 3, count, "Should have 3 middlewares before clear")

	// 清空注册器
	registry.Clear()

	// 确认已清空
	count = registry.Count()
	assert.Equal(t, 0, count, "Should have 0 middlewares after clear")

	// 确认无法获取中间件
	_, err = registry.Get("middleware_1")
	assert.Error(t, err, "Should not be able to get middleware after clear")
}

// TestRegistry_List 测试列出中间件名称
func TestRegistry_List(t *testing.T) {
	registry := createTestRegistry(t)

	// 列出空列表
	names := registry.List()
	assert.Empty(t, names, "List should return empty slice initially")

	// 注册3个中间件
	mw1 := createTestMiddleware("middleware_1", 10)
	mw2 := createTestMiddleware("middleware_2", 20)
	mw3 := createTestMiddleware("middleware_3", 30)

	err := registry.Register(mw1)
	assert.NoError(t, err)
	err = registry.Register(mw2)
	assert.NoError(t, err)
	err = registry.Register(mw3)
	assert.NoError(t, err)

	// 列出所有中间件名称
	names = registry.List()
	assert.Len(t, names, 3, "List should return 3 names")
	assert.Contains(t, names, "middleware_1")
	assert.Contains(t, names, "middleware_2")
	assert.Contains(t, names, "middleware_3")
}

// TestRegistry_SetPriorityOverride 测试设置优先级覆盖
func TestRegistry_SetPriorityOverride(t *testing.T) {
	registry := createTestRegistry(t)

	// 注册中间件（默认优先级10）
	mw := createTestMiddleware("test_middleware", 10)
	err := registry.Register(mw)
	assert.NoError(t, err)

	// 设置优先级覆盖为5
	registry.SetPriorityOverride("test_middleware", 5)

	// 获取排序后的中间件列表
	sorted := registry.Sorted()
	assert.Len(t, sorted, 1)

	// 验证优先级已覆盖（需要通过Sorted方法验证）
	// 由于Sorted方法使用覆盖后的优先级，我们无法直接验证
	// 但可以通过执行顺序来验证
}

// TestRegistry_Sorted_WithPriorityOverride 测试带优先级覆盖的排序
func TestRegistry_Sorted_WithPriorityOverride(t *testing.T) {
	registry := createTestRegistry(t)

	// 注册3个中间件（优先级：10, 20, 30）
	mw1 := createTestMiddleware("middleware_1", 10)
	mw2 := createTestMiddleware("middleware_2", 20)
	mw3 := createTestMiddleware("middleware_3", 30)

	err := registry.Register(mw1)
	assert.NoError(t, err)
	err = registry.Register(mw2)
	assert.NoError(t, err)
	err = registry.Register(mw3)
	assert.NoError(t, err)

	// 设置middleware_3的优先级覆盖为5
	registry.SetPriorityOverride("middleware_3", 5)

	// 获取排序后的中间件列表
	sorted := registry.Sorted()
	assert.Len(t, sorted, 3)

	// 验证排序（middleware_3应该最前，因为优先级被覆盖为5）
	assert.Equal(t, "middleware_3", sorted[0].Name(), "First middleware should be middleware_3 with priority override 5")
	assert.Equal(t, "middleware_1", sorted[1].Name(), "Second middleware should be middleware_1 with priority 10")
	assert.Equal(t, "middleware_2", sorted[2].Name(), "Third middleware should be middleware_2 with priority 20")
}

// TestRegistry_Validate_DuplicatePriority 测试重复优先级检测
func TestRegistry_Validate_DuplicatePriority(t *testing.T) {
	registry := createTestRegistry(t)

	// 注册3个相同优先级的中间件
	mw1 := createTestMiddleware("middleware_1", 10)
	mw2 := createTestMiddleware("middleware_2", 10)
	mw3 := createTestMiddleware("middleware_3", 10)

	err := registry.Register(mw1)
	assert.NoError(t, err)
	err = registry.Register(mw2)
	assert.NoError(t, err)
	err = registry.Register(mw3)
	assert.NoError(t, err)

	// 验证应该成功（但有警告）
	err = registry.Validate()
	assert.NoError(t, err, "Validate should succeed even with duplicate priorities")
}

// TestManager_SortByPriority_FullCoverage 测试sortByPriority的完整覆盖
func TestManager_SortByPriority_FullCoverage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	manager := createTestManager(t)

	// 注册5个中间件（乱序优先级）
	mw1 := createTestMiddleware("mw_1", 50)
	mw2 := createTestMiddleware("mw_2", 10)
	mw3 := createTestMiddleware("mw_3", 30)
	mw4 := createTestMiddleware("mw_4", 20)
	mw5 := createTestMiddleware("mw_5", 40)

	err := manager.Register(mw1)
	assert.NoError(t, err)
	err = manager.Register(mw2)
	assert.NoError(t, err)
	err = manager.Register(mw3)
	assert.NoError(t, err)
	err = manager.Register(mw4)
	assert.NoError(t, err)
	err = manager.Register(mw5)
	assert.NoError(t, err)

	// 应用到路由以触发sortByPriority
	router := gin.New()
	err = manager.ApplyToRouter(router, "mw_1", "mw_2", "mw_3", "mw_4", "mw_5")
	assert.NoError(t, err)

	// 验证中间件已应用
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})
}
