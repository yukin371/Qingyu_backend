package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"Qingyu_backend/internal/middleware/auth"
	middlewareAuth "Qingyu_backend/internal/middleware/auth"

	"github.com/gin-gonic/gin"
)

// TestPermissionAPIEndToEnd 测试权限API端到端功能
func TestPermissionAPIEndToEnd(t *testing.T) {
	if os.Getenv("TEST_MODE") != "true" {
		t.Skip("跳过端到端测试（TEST_MODE未设置）")
	}

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	r := gin.New()

	// 1. 添加权限中间件
	permConfig := &auth.PermissionConfig{
		Enabled:    true,
		Strategy:   "rbac",
		ConfigPath: "../../configs/permissions.yaml",
		SkipPaths: []string{
			"/api/v1/auth/login",
			"/api/v1/public",
		},
		Message:    "权限不足，无法访问该资源",
		StatusCode: 403,
	}

	logger := zap.NewNop()
	permMiddleware, err := auth.NewPermissionMiddleware(permConfig, logger)
	require.NoError(t, err, "创建权限中间件失败")

	// 2. 创建RBACChecker并设置权限
	checker, err := middlewareAuth.NewRBACChecker(nil)
	require.NoError(t, err)
	rbacChecker := checker.(*middlewareAuth.RBACChecker)

	// 设置测试权限：admin用户有所有权限，reader用户只有读权限
	rbacChecker.GrantPermission("admin", "*:*")
	rbacChecker.GrantPermission("reader", "book:read")
	rbacChecker.GrantPermission("reader", "chapter:read")
	rbacChecker.AssignRole("admin_user", "admin")
	rbacChecker.AssignRole("reader_user", "reader")

	// 3. 设置中间件到权限中间件（用于动态检查）
	permMiddleware.SetChecker(rbacChecker)

	// 4. 创建模拟的认证中间件（设置用户ID到context）
	authMiddleware := createMockAuthMiddleware()

	// 5. 添加路由
	api := r.Group("/api/v1")
	api.Use(authMiddleware)
	api.Use(permMiddleware.Handler())
	{
		// 公开路由
		api.GET("/public", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "公开访问成功"})
		})

		// 需要认证的路由
		api.GET("/books", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "获取书籍列表成功"})
		})

		api.POST("/books", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "创建书籍成功"})
		})

		api.DELETE("/books/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "删除书籍成功"})
		})

		api.GET("/chapters", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "获取章节列表成功"})
		})

		api.POST("/chapters", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "创建章节成功"})
		})
	}

	t.Run("公开访问不需要权限", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/public", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "公开访问成功")
	})

	t.Run("Admin用户可以访问所有API", func(t *testing.T) {
		tests := []struct {
			method   string
			url      string
			expected int
		}{
			{"GET", "/api/v1/books", 200},
			{"POST", "/api/v1/books", 200},
			{"DELETE", "/api/v1/books/123", 200},
			{"GET", "/api/v1/chapters", 200},
			{"POST", "/api/v1/chapters", 200},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s_%s", tt.method, tt.url), func(t *testing.T) {
				req, _ := http.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("X-User-ID", "admin_user")
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				assert.Equal(t, tt.expected, w.Code, fmt.Sprintf("%s %s should return %d", tt.method, tt.url, tt.expected))
			})
		}
	})

	t.Run("Reader用户只能访问读API", func(t *testing.T) {
		tests := []struct {
			name     string
			method   string
			url      string
			expected int
		}{
			{"读取书籍", "GET", "/api/v1/books", 200},
			{"创建书籍", "POST", "/api/v1/books", 403},
			{"删除书籍", "DELETE", "/api/v1/books/123", 403},
			{"读取章节", "GET", "/api/v1/chapters", 200},
			{"创建章节", "POST", "/api/v1/chapters", 403},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest(tt.method, tt.url, nil)
				req.Header.Set("X-User-ID", "reader_user")
				w := httptest.NewRecorder()

				r.ServeHTTP(w, req)

				assert.Equal(t, tt.expected, w.Code, fmt.Sprintf("%s should return %d", tt.name, tt.expected))
			})
		}
	})

	t.Run("无权限用户被拒绝访问", func(t *testing.T) {
		// 为无权限用户创建请求
		rbacChecker.AssignRole("no_perm_user", "guest") // guest角色没有任何权限

		req, _ := http.NewRequest("GET", "/api/v1/books", nil)
		req.Header.Set("X-User-ID", "no_perm_user")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "权限不足")
	})

	t.Run("无用户ID时被拒绝", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/books", nil)
		// 不设置X-User-ID
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code) // 认证失败
	})

	t.Run("通配符权限测试", func(t *testing.T) {
		// 测试admin的*:*通配符权限
		rbacChecker.GrantPermission("wildcard_user", "book:*")
		rbacChecker.AssignRole("wildcard_test", "wildcard_user")

		// 应该可以访问book相关的所有操作
		req1, _ := http.NewRequest("GET", "/api/v1/books", nil)
		req1.Header.Set("X-User-ID", "wildcard_test")
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)
		assert.Equal(t, 200, w1.Code, "book:*应该允许读取书籍")

		req2, _ := http.NewRequest("POST", "/api/v1/books", nil)
		req2.Header.Set("X-User-ID", "wildcard_test")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)
		assert.Equal(t, 200, w2.Code, "book:*应该允许创建书籍")

		// 但不能访问chapter
		req3, _ := http.NewRequest("GET", "/api/v1/chapters", nil)
		req3.Header.Set("X-User-ID", "wildcard_test")
		w3 := httptest.NewRecorder()
		r.ServeHTTP(w3, req3)
		assert.Equal(t, 403, w3.Code, "book:*不应该允许访问章节")
	})
}

// TestPermissionMiddlewarePriority 测试权限中间件优先级
func TestPermissionMiddlewarePriority(t *testing.T) {
	permConfig := &auth.PermissionConfig{
		Enabled:  true,
		Strategy: "rbac",
	}

	logger := zap.NewNop()
	permMiddleware, err := auth.NewPermissionMiddleware(permConfig, logger)
	require.NoError(t, err)

	// 验证优先级
	assert.Equal(t, 10, permMiddleware.Priority(), "权限中间件优先级应该是10（在认证之后）")
}

// TestPermissionHotReload 测试权限热更新（不实际更新，只测试接口）
func TestPermissionHotReload(t *testing.T) {
	permConfig := &auth.PermissionConfig{
		Enabled:    true,
		Strategy:   "rbac",
		ConfigPath: "../../configs/permissions.yaml",
	}

	logger := zap.NewNop()
	permMiddleware, err := auth.NewPermissionMiddleware(permConfig, logger)
	require.NoError(t, err)

	// 测试Reload方法
	err = permMiddleware.Reload()
	assert.NoError(t, err, "Reload方法应该不会失败")
}

// createMockAuthMiddleware 创建模拟认证中间件
func createMockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 设置用户ID到context（权限中间件会读取）
		c.Set("user_id", userID)
		c.Next()
	}
}

// BenchmarkPermissionCheck 性能测试
func BenchmarkPermissionCheck(b *testing.B) {
	gin.SetMode(gin.TestMode)

	permConfig := &auth.PermissionConfig{
		Enabled:  true,
		Strategy: "rbac",
	}

	logger := zap.NewNop()
	permMiddleware, _ := auth.NewPermissionMiddleware(permConfig, logger)
	checker, _ := middlewareAuth.NewRBACChecker(nil)
	rbacChecker := checker.(*middlewareAuth.RBACChecker)

	// 设置大量权限
	for i := 0; i < 100; i++ {
		rbacChecker.GrantPermission(fmt.Sprintf("role_%d", i), fmt.Sprintf("resource_%d:action_%d", i, i))
	}
	rbacChecker.AssignRole("test_user", "role_0")

	permMiddleware.SetChecker(rbacChecker)

	r := gin.New()
	authMiddleware := createMockAuthMiddleware()
	r.Use(authMiddleware)
	r.Use(permMiddleware.Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-User-ID", "test_user")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
