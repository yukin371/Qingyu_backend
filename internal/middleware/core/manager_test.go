package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// mockMiddleware 模拟中间件实现
type mockMiddleware struct {
	name     string
	priority int
	handler  gin.HandlerFunc
}

func (m *mockMiddleware) Name() string {
	return m.name
}

func (m *mockMiddleware) Priority() int {
	return m.priority
}

func (m *mockMiddleware) Handler() gin.HandlerFunc {
	return m.handler
}

// 测试辅助函数：创建测试管理器
func createTestManager(t *testing.T) Manager {
	t.Helper()

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	return NewManager(logger)
}

// 测试辅助函数：创建测试中间件
func createTestMiddleware(name string, priority int) Middleware {
	return &mockMiddleware{
		name:     name,
		priority: priority,
		handler: func(c *gin.Context) {
			c.Next()
		},
	}
}

// 测试辅助函数：创建带标记的中间件（用于验证执行顺序）
func createMarkingMiddleware(name string, priority int, order *[]string) Middleware {
	return &mockMiddleware{
		name:     name,
		priority: priority,
		handler: func(c *gin.Context) {
			*order = append(*order, name)
			c.Next()
		},
	}
}

// TestManager_Register 测试中间件注册
func TestManager_Register(t *testing.T) {
	manager := createTestManager(t)

	// 测试正常注册
	middleware := createTestMiddleware("test_middleware", 10)
	err := manager.Register(middleware)
	assert.NoError(t, err, "Register should succeed")

	// 测试重复注册
	err = manager.Register(middleware)
	assert.Error(t, err, "Register should fail for duplicate middleware")
	assert.Contains(t, err.Error(), "already registered")
}

// TestManager_Register_WithPriorityOverride 测试带优先级覆盖的注册
func TestManager_Register_WithPriorityOverride(t *testing.T) {
	manager := createTestManager(t)

	// 创建优先级为10的中间件
	middleware := createTestMiddleware("test_middleware", 10)

	// 注册时覆盖优先级为5
	err := manager.Register(middleware, WithPriority(5))
	assert.NoError(t, err, "Register with priority override should succeed")

	// 验证中间件已注册
	retrieved, err := manager.Get("test_middleware")
	assert.NoError(t, err, "Get should succeed")
	assert.Equal(t, "test_middleware", retrieved.Name())
}

// TestManager_Unregister 测试中间件注销
func TestManager_Unregister(t *testing.T) {
	manager := createTestManager(t)

	// 注册中间件
	middleware := createTestMiddleware("test_middleware", 10)
	err := manager.Register(middleware)
	assert.NoError(t, err)

	// 注销中间件
	err = manager.Unregister("test_middleware")
	assert.NoError(t, err, "Unregister should succeed")

	// 验证中间件已不存在
	_, err = manager.Get("test_middleware")
	assert.Error(t, err, "Get should fail for unregistered middleware")
	assert.Contains(t, err.Error(), "not found")
}

// TestManager_Unregister_NotFound 测试注销不存在的中间件
func TestManager_Unregister_NotFound(t *testing.T) {
	manager := createTestManager(t)

	// 尝试注销不存在的中间件
	err := manager.Unregister("nonexistent")
	assert.Error(t, err, "Unregister should fail for nonexistent middleware")
	assert.Contains(t, err.Error(), "not found")
}

// TestManager_Get 测试获取中间件
func TestManager_Get(t *testing.T) {
	manager := createTestManager(t)

	// 测试获取不存在的中间件
	_, err := manager.Get("nonexistent")
	assert.Error(t, err, "Get should fail for nonexistent middleware")

	// 注册中间件
	middleware := createTestMiddleware("test_middleware", 10)
	err = manager.Register(middleware)
	assert.NoError(t, err)

	// 测试获取存在的中间件
	retrieved, err := manager.Get("test_middleware")
	assert.NoError(t, err, "Get should succeed")
	assert.Equal(t, "test_middleware", retrieved.Name())
	assert.Equal(t, 10, retrieved.Priority())
}

// TestManager_List 测试列出所有中间件
func TestManager_List(t *testing.T) {
	manager := createTestManager(t)

	// 测试空列表
	list := manager.List()
	assert.Empty(t, list, "List should return empty slice initially")

	// 注册多个中间件
	mw1 := createTestMiddleware("middleware_1", 30)
	mw2 := createTestMiddleware("middleware_2", 10)
	mw3 := createTestMiddleware("middleware_3", 20)

	err := manager.Register(mw1)
	assert.NoError(t, err)
	err = manager.Register(mw2)
	assert.NoError(t, err)
	err = manager.Register(mw3)
	assert.NoError(t, err)

	// 列出所有中间件
	list = manager.List()
	assert.Len(t, list, 3, "List should return 3 middlewares")

	// 验证排序（按优先级升序）
	assert.Equal(t, "middleware_2", list[0].Name(), "First middleware should have priority 10")
	assert.Equal(t, "middleware_3", list[1].Name(), "Second middleware should have priority 20")
	assert.Equal(t, "middleware_1", list[2].Name(), "Third middleware should have priority 30")
}

// TestManager_Build 测试构建中间件链
func TestManager_Build(t *testing.T) {
	manager := createTestManager(t)

	// 测试空中间件链
	chain := manager.Build()
	assert.Empty(t, chain, "Build should return empty slice initially")

	// 注册多个中间件
	mw1 := createTestMiddleware("middleware_1", 30)
	mw2 := createTestMiddleware("middleware_2", 10)
	mw3 := createTestMiddleware("middleware_3", 20)

	err := manager.Register(mw1)
	assert.NoError(t, err)
	err = manager.Register(mw2)
	assert.NoError(t, err)
	err = manager.Register(mw3)
	assert.NoError(t, err)

	// 构建中间件链
	chain = manager.Build()
	assert.Len(t, chain, 3, "Build should return 3 handlers")
}

// TestManager_ApplyToRouter_AllMiddlewares 测试应用所有中间件到路由
func TestManager_ApplyToRouter_AllMiddlewares(t *testing.T) {
	gin.SetMode(gin.TestMode)
	manager := createTestManager(t)

	// 注册中间件
	mw1 := createTestMiddleware("middleware_1", 10)
	mw2 := createTestMiddleware("middleware_2", 20)

	err := manager.Register(mw1)
	assert.NoError(t, err)
	err = manager.Register(mw2)
	assert.NoError(t, err)

	// 创建路由并应用所有中间件
	router := gin.New()
	err = manager.ApplyToRouter(router)
	assert.NoError(t, err, "ApplyToRouter should succeed")

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestManager_ApplyToRouter_SelectiveMiddlewares 测试应用选定的中间件到路由
func TestManager_ApplyToRouter_SelectiveMiddlewares(t *testing.T) {
	gin.SetMode(gin.TestMode)
	manager := createTestManager(t)

	// 注册3个中间件
	mw1 := createTestMiddleware("middleware_1", 10)
	mw2 := createTestMiddleware("middleware_2", 20)
	mw3 := createTestMiddleware("middleware_3", 30)

	err := manager.Register(mw1)
	assert.NoError(t, err)
	err = manager.Register(mw2)
	assert.NoError(t, err)
	err = manager.Register(mw3)
	assert.NoError(t, err)

	// 创建路由并应用选定的中间件
	router := gin.New()
	err = manager.ApplyToRouter(router, "middleware_1", "middleware_3")
	assert.NoError(t, err, "ApplyToRouter should succeed")

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestManager_ApplyToRouter_NonexistentMiddleware 测试应用不存在的中间件
func TestManager_ApplyToRouter_NonexistentMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	manager := createTestManager(t)

	// 注册中间件
	mw1 := createTestMiddleware("middleware_1", 10)
	err := manager.Register(mw1)
	assert.NoError(t, err)

	// 创建路由并尝试应用不存在的中间件
	router := gin.New()
	err = manager.ApplyToRouter(router, "middleware_1", "nonexistent")
	assert.Error(t, err, "ApplyToRouter should fail for nonexistent middleware")
	assert.Contains(t, err.Error(), "not found")
}

// TestManager_ExecutionOrder 测试中间件执行顺序
func TestManager_ExecutionOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	manager := createTestManager(t)

	// 记录执行顺序
	var order []string

	// 注册3个中间件（优先级：30, 10, 20）
	mw1 := createMarkingMiddleware("mw1", 30, &order)
	mw2 := createMarkingMiddleware("mw2", 10, &order)
	mw3 := createMarkingMiddleware("mw3", 20, &order)

	err := manager.Register(mw1)
	assert.NoError(t, err)
	err = manager.Register(mw2)
	assert.NoError(t, err)
	err = manager.Register(mw3)
	assert.NoError(t, err)

	// 构建中间件链
	chain := manager.Build()
	assert.Len(t, chain, 3)

	// 创建Gin引擎并应用中间件
	router := gin.New()
	router.Use(chain...)

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 执行请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// 验证执行顺序（按优先级升序：10, 20, 30）
	assert.Len(t, order, 3)
	assert.Equal(t, "mw2", order[0], "First middleware should have priority 10")
	assert.Equal(t, "mw3", order[1], "Second middleware should have priority 20")
	assert.Equal(t, "mw1", order[2], "Third middleware should have priority 30")
}

// TestManager_MultipleRegistrationsWithPriorityOverride 测试多次注册和优先级覆盖
func TestManager_MultipleRegistrationsWithPriorityOverride(t *testing.T) {
	manager := createTestManager(t)

	// 注册中间件（默认优先级10）
	mw1 := createTestMiddleware("middleware_1", 10)
	err := manager.Register(mw1)
	assert.NoError(t, err)

	// 注册另一个中间件（默认优先级20）
	mw2 := createTestMiddleware("middleware_2", 20)
	err = manager.Register(mw2)
	assert.NoError(t, err)

	// 注册第三个中间件（默认优先级30，但覆盖为5）
	mw3 := createTestMiddleware("middleware_3", 30)
	err = manager.Register(mw3, WithPriority(5))
	assert.NoError(t, err)

	// 列出所有中间件
	list := manager.List()
	assert.Len(t, list, 3)

	// 验证排序（mw3应该最前，因为优先级被覆盖为5）
	assert.Equal(t, "middleware_3", list[0].Name(), "First middleware should be middleware_3 with priority 5")
	assert.Equal(t, "middleware_1", list[1].Name(), "Second middleware should be middleware_1 with priority 10")
	assert.Equal(t, "middleware_2", list[2].Name(), "Third middleware should be middleware_2 with priority 20")
}

// TestManager_List_AfterUnregister 测试注销后的列表
func TestManager_List_AfterUnregister(t *testing.T) {
	manager := createTestManager(t)

	// 注册3个中间件
	mw1 := createTestMiddleware("middleware_1", 10)
	mw2 := createTestMiddleware("middleware_2", 20)
	mw3 := createTestMiddleware("middleware_3", 30)

	err := manager.Register(mw1)
	assert.NoError(t, err)
	err = manager.Register(mw2)
	assert.NoError(t, err)
	err = manager.Register(mw3)
	assert.NoError(t, err)

	// 注销中间件2
	err = manager.Unregister("middleware_2")
	assert.NoError(t, err)

	// 列出所有中间件
	list := manager.List()
	assert.Len(t, list, 2, "List should return 2 middlewares after unregister")

	// 验证剩余中间件
	names := make([]string, 0, 2)
	for _, mw := range list {
		names = append(names, mw.Name())
	}
	assert.Contains(t, names, "middleware_1")
	assert.Contains(t, names, "middleware_3")
	assert.NotContains(t, names, "middleware_2")
}
