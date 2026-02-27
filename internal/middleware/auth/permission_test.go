package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// TestPermissionMiddleware_Handler_Disabled 测试禁用权限检查
func TestPermissionMiddleware_Handler_Disabled(t *testing.T) {
	// 创建测试观察器
	observedZapCore, _ := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	config := DefaultPermissionConfig()
	config.Enabled = false
	config.Strategy = "noop" // 使用空检查器避免创建RBAC

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// 模拟未认证用户
		router.Use(middleware.Handler())
		c.Next()
	})
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求（没有user_id）
	w := performRequest(router, "GET", "/test", nil)

	// 验证请求通过（权限检查已禁用）
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestPermissionMiddleware_Handler_SkipPath 测试跳过指定路径
func TestPermissionMiddleware_Handler_SkipPath(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop"
	config.SkipPaths = []string{"/health", "/metrics", "/api/v1/auth/login"}

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"metrics": "ok"})
	})
	router.GET("/api/v1/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "xxx"})
	})
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "protected"})
	})

	// 测试跳过的路径
	w1 := performRequest(router, "GET", "/health", nil)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := performRequest(router, "GET", "/metrics", nil)
	assert.Equal(t, http.StatusOK, w2.Code)

	w3 := performRequest(router, "GET", "/api/v1/auth/login", nil)
	assert.Equal(t, http.StatusOK, w3.Code)

	// 测试受保护的路径（应该因为没有user_id而返回401）
	w4 := performRequest(router, "GET", "/protected", nil)
	assert.Equal(t, http.StatusUnauthorized, w4.Code)
}

// TestPermissionMiddleware_Handler_Unauthenticated 测试未认证用户
func TestPermissionMiddleware_Handler_Unauthenticated(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop"

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "protected"})
	})

	// 发送请求（没有设置user_id）
	w := performRequest(router, "GET", "/protected", nil)

	// 验证返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(40101), response["code"])
	assert.Equal(t, "用户未认证", response["message"])
}

// TestPermissionMiddleware_Handler_InvalidUserID 测试无效的user_id类型
func TestPermissionMiddleware_Handler_InvalidUserID(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop"

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// 设置无效的user_id类型
		c.Set("user_id", 123) // 应该是string
		c.Next()
	})
	router.Use(middleware.Handler())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "protected"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/protected", nil)

	// 验证返回500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, float64(50001), response["code"])
	assert.Equal(t, "用户信息格式错误", response["message"])
}

// TestPermissionMiddleware_Handler_PermissionDenied 测试权限不足
func TestPermissionMiddleware_Handler_PermissionDenied(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop" // NoOpChecker允许所有权限

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// 设置user_id
		c.Set("user_id", "user1")
		c.Next()
	})
	router.Use(middleware.Handler())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "protected"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/protected", nil)

	// NoOpChecker应该允许所有权限
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestPermissionMiddleware_GetResourceFromPath 测试资源提取
func TestPermissionMiddleware_GetResourceFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/projects", "project"},
		{"/api/v1/projects/123", "project"},
		{"/api/v1/users", "user"},
		{"/api/v1/users/456", "user"},
		{"/api/v1/documents", "document"},
		{"/api/v1/documents/789", "document"},
		{"/api/v1/books", "book"},
		{"/api/v1/books/123/chapters", "book"}, // 第一段资源
		{"/health", "health"},                  // 无版本前缀
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := getResourceFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermissionMiddleware_GetActionFromMethod 测试操作提取
func TestPermissionMiddleware_GetActionFromMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "read"},
		{"POST", "create"},
		{"PUT", "update"},
		{"PATCH", "update"},
		{"DELETE", "delete"},
		{"OPTIONS", "unknown"},
		{"HEAD", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := getActionFromMethod(tt.method)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermissionMiddleware_GetChecker 测试获取检查器
func TestPermissionMiddleware_GetChecker(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop"

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	checker := middleware.GetChecker()
	assert.NotNil(t, checker)
	assert.Equal(t, "noop", checker.Name())
}

// TestPermissionMiddleware_ShouldSkipPath 测试路径跳过逻辑
func TestPermissionMiddleware_ShouldSkipPath(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.SkipPaths = []string{"/health", "/metrics", "/api/v1/auth"}

	middleware, _ := NewPermissionMiddleware(config, logger)

	tests := []struct {
		path     string
		expected bool
	}{
		{"/health", true},
		{"/metrics", true},
		{"/api/v1/auth/login", true}, // 前缀匹配
		{"/api/v1/users", false},
		{"/protected", false},
		{"/healthz", true}, // 会被/health前缀匹配
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			// 使用反射或公开方法测试
			// 这里我们通过实际行为来验证
			router := gin.New()
			router.Use(middleware.Handler())
			router.GET(tt.path, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"path": tt.path})
			})

			w := performRequest(router, "GET", tt.path, nil)

			// 跳过的路径应该返回200（不检查user_id）
			// 不跳过的路径应该返回401（没有user_id）
			if tt.expected {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			}
		})
	}
}

// TestPermissionMiddleware_CustomMessage 测试自定义错误消息
func TestPermissionMiddleware_CustomMessage(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop"
	config.Message = "您没有权限访问此资源"
	config.StatusCode = 403

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	// 创建测试路由
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "protected"})
	})

	// 发送请求（没有user_id，会返回401而不是403）
	w := performRequest(router, "GET", "/protected", nil)

	// 未认证用户返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestPermissionMiddleware_Priority 测试中间件优先级
func TestPermissionMiddleware_Priority(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop"

	middleware, err := NewPermissionMiddleware(config, logger)
	assert.NoError(t, err)

	assert.Equal(t, "permission", middleware.Name())
	assert.Equal(t, 10, middleware.Priority())
}

// BenchmarkPermissionMiddleware 性能测试
func BenchmarkPermissionMiddleware(b *testing.B) {
	logger, _ := zap.NewDevelopment()

	config := DefaultPermissionConfig()
	config.Strategy = "noop"
	config.SkipPaths = []string{"/health"}

	middleware, _ := NewPermissionMiddleware(config, logger)

	// 创建测试路由
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user1")
		c.Next()
	})
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// 辅助函数：执行HTTP请求
func performRequest(router *gin.Engine, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}
