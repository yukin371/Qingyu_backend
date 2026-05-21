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

func TestQuotaMiddlewareWithInterface(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	}

	t.Run("配额检查通过时应该继续处理请求", func(t *testing.T) {
		checker := &passingMockChecker{}

		middleware := NewQuotaMiddleware(checker)
		router := gin.New()
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
		checker := &failingMockChecker{}

		middleware := NewQuotaMiddleware(checker)
		router := gin.New()
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

type passingMockChecker struct{}

func (m *passingMockChecker) Check(ctx context.Context, userID string, amount int) error {
	return nil
}

type failingMockChecker struct{}

func (m *failingMockChecker) Check(ctx context.Context, userID string, amount int) error {
	return quota.ErrInsufficientQuota
}
