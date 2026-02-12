package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/pkg/quota"
)

// TestQuotaMiddlewareWithInterface 测试中间件使用接口而非具体实现
// 这是TDD重构的核心：先定义期望的接口行为
func TestQuotaMiddlewareWithInterface(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 模拟认证中间件：设置user_id
	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}

	t.Run("配额检查通过时应该继续处理请求", func(t *testing.T) {
		// 创建一个会通过检查的mock
		checker := &passingMockChecker{}

		middleware := NewQuotaMiddlewareWithChecker(checker)
		router := gin.New()

		// 注意：认证中间件必须在配额中间件之前
		router.Use(authMiddleware)
		router.Use(middleware.CheckQuota(1000))

		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("配额不足时应该返回错误", func(t *testing.T) {
		// 创建一个会失败的mock
		checker := &failingMockChecker{}

		middleware := NewQuotaMiddlewareWithChecker(checker)
		router := gin.New()

		// 注意：认证中间件必须在配额中间件之前
		router.Use(authMiddleware)
		router.Use(middleware.CheckQuota(1000))

		router.GET("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "success")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "配额")
	})
}

// passingMockChecker 总是通过的配额检查器
type passingMockChecker struct{}

func (m *passingMockChecker) Check(ctx context.Context, userID string, amount int) *quota.CheckResult {
	return &quota.CheckResult{
		Allowed:   true,
		Remaining: 1000,
		Error:     nil,
	}
}

// failingMockChecker 总是失败的配额检查器
type failingMockChecker struct{}

func (m *failingMockChecker) Check(ctx context.Context, userID string, amount int) *quota.CheckResult {
	return &quota.CheckResult{
		Allowed:   false,
		Remaining: 0,
		Error:     quota.ErrInsufficientQuota,
	}
}
