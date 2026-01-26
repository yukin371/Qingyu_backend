package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(redisClient *redis.Client, config RedisRateLimitConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(RedisRateLimitMiddleware(redisClient, config))

	router.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	return router
}

func TestRedisRateLimitMiddleware_AllowWithinLimit(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	config := DefaultRedisRateLimitConfig()
	config.RequestsPerMinute = 5

	router := setupTestRouter(redisClient, config)

	// 发送 5 个请求（在限制内）
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/search", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRedisRateLimitMiddleware_BlockWhenExceeded(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	config := DefaultRedisRateLimitConfig()
	config.RequestsPerMinute = 3

	router := setupTestRouter(redisClient, config)

	// 发送 3 个请求（达到限制）
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/search", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 第 4 个请求应该被拒绝
	req, _ := http.NewRequest("GET", "/search", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "搜索请求过于频繁")
}

func TestRedisRateLimitMiddleware_DifferentIPs(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	config := DefaultRedisRateLimitConfig()
	config.RequestsPerMinute = 2

	router := setupTestRouter(redisClient, config)

	// IP1 发送 2 个请求
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("GET", "/search", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}

	// IP1 的第 3 个请求应该被拒绝
	req1, _ := http.NewRequest("GET", "/search", nil)
	req1.RemoteAddr = "192.168.1.1:1234"
	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusTooManyRequests, w1.Code)

	// IP2 的请求应该正常通过（有独立的限制）
	req2, _ := http.NewRequest("GET", "/search", nil)
	req2.RemoteAddr = "192.168.1.2:1234"
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRedisRateLimitMiddleware_WithUserID(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	config := DefaultRedisRateLimitConfig()
	config.RequestsPerMinute = 2

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		// 模拟认证中间件，设置 user_id
		c.Set("user_id", "test-user-123")
		c.Next()
	})

	router.Use(RedisRateLimitMiddleware(redisClient, config))

	router.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 发送 3 个请求
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/search", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}
}

func TestRedisRateLimitMiddleware_RateLimitHeaders(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	config := DefaultRedisRateLimitConfig()
	config.RequestsPerMinute = 10

	router := setupTestRouter(redisClient, config)

	req, _ := http.NewRequest("GET", "/search", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("X-RateLimit-Limit"), "10")
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"))
}

func TestGetRateLimitStatus(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	ctx := context.Background()
	key := "search:ratelimit:test"
	limit := 10

	// 设置计数器
	redisClient.Set(ctx, key, 5, 1*time.Minute)

	current, remaining, resetTime, err := GetRateLimitStatus(ctx, redisClient, key, limit)

	require.NoError(t, err)
	assert.Equal(t, int64(5), current)
	assert.Equal(t, int64(5), remaining)
	assert.True(t, resetTime.After(time.Now()))
}

func TestDefaultRedisRateLimitConfig(t *testing.T) {
	config := DefaultRedisRateLimitConfig()

	assert.Equal(t, 60, config.RequestsPerMinute)
	assert.Equal(t, "search:ratelimit:", config.KeyPrefix)
	assert.False(t, config.SkipAuthenticated)
	assert.NotEmpty(t, config.Message)
	assert.Equal(t, http.StatusTooManyRequests, config.StatusCode)
}

func TestSearchRateLimit(t *testing.T) {
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SearchRateLimit(redisClient))

	router.GET("/search", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 使用默认配置（60 次/分钟）
	for i := 0; i < 61; i++ {
		req, _ := http.NewRequest("GET", "/search", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if i < 60 {
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}
}
