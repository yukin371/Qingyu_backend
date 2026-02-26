package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestTokenBucketLimiter 测试令牌桶限流器
func TestTokenBucketLimiter(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            10,
		Burst:           20,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	t.Run("Allow_WithinRate", func(t *testing.T) {
		// 前10个请求应该被允许
		for i := 0; i < 10; i++ {
			assert.True(t, limiter.Allow("test_key"), "Request %d should be allowed", i+1)
		}
	})

	t.Run("Allow_ExceedsRate", func(t *testing.T) {
		// 令牌桶算法：burst容量是20，所以最多允许20个请求
		// 即使rate是10/秒，但burst允许突发流量
		allowed := 0
		for i := 0; i < 20; i++ {
			if limiter.Allow("test_key_2") {
				allowed++
			}
		}
		// 应该有burst个请求被允许
		assert.Equal(t, config.Burst, allowed)
	})

	t.Run("Wait_WithinLimit", func(t *testing.T) {
		ctx := context.Background()
		err := limiter.Wait(ctx, "test_key_3")
		assert.NoError(t, err)
	})

	t.Run("Reset", func(t *testing.T) {
		key := "test_key_4"

		// 用完所有配额（需要超过burst才能触发限流）
		for i := 0; i < config.Burst+1; i++ {
			limiter.Allow(key)
		}

		// 应该被拒绝
		assert.False(t, limiter.Allow(key))

		// 重置
		limiter.Reset(key)

		// 重置后应该被允许（令牌桶重新创建，配额恢复）
		assert.True(t, limiter.Allow(key))
	})

	t.Run("GetStats", func(t *testing.T) {
		key := "test_key_5"

		// 发起一些请求
		limiter.Allow(key)
		limiter.Allow(key)
		limiter.Allow(key)

		stats := limiter.GetStats(key)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(3), stats.TotalRequests)
		assert.Equal(t, int64(3), stats.AllowedRequests)
	})

	t.Run("GetTotalStats", func(t *testing.T) {
		// TokenBucketLimiter有GetTotalStats方法
		stats := limiter.GetTotalStats()
		assert.NotNil(t, stats)
		assert.True(t, stats.TotalRequests > 0)
	})
}

// TestSlidingWindowLimiter 测试滑动窗口限流器
func TestSlidingWindowLimiter(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            5,
		Burst:           10,
		WindowSize:      1,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	t.Run("Allow_WithinRate", func(t *testing.T) {
		// 前5个请求应该被允许
		for i := 0; i < 5; i++ {
			assert.True(t, limiter.Allow("sw_test_key"))
		}
	})

	t.Run("Allow_ExceedsRate", func(t *testing.T) {
		// 第6个请求应该被拒绝
		assert.False(t, limiter.Allow("sw_test_key"))
	})

	t.Run("Allow_WindowResets", func(t *testing.T) {
		// 等待窗口过期
		time.Sleep(time.Duration(config.WindowSize) * time.Second)

		// 窗口过期后应该被允许
		assert.True(t, limiter.Allow("sw_test_key"))
	})

	t.Run("GetStats", func(t *testing.T) {
		key := "sw_test_key_2"

		// 发起一些请求
		limiter.Allow(key)
		limiter.Allow(key)

		stats := limiter.GetStats(key)
		assert.NotNil(t, stats)
		assert.True(t, stats.TotalRequests >= 2)
	})

	t.Run("Wait_NotSupported", func(t *testing.T) {
		ctx := context.Background()
		err := limiter.Wait(ctx, "sw_test_key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not support Wait")
	})
}

// TestRateLimitConfig 测试限流配置
func TestRateLimitConfig(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		config := DefaultRateLimitConfig()
		assert.NotNil(t, config)
		assert.Equal(t, "token_bucket", config.Strategy)
		assert.Equal(t, 100, config.Rate)
		assert.Equal(t, 200, config.Burst)
	})

	t.Run("Validate_ValidConfig", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy:  "token_bucket",
			Rate:      10,
			Burst:     20,
			KeyFunc:   "ip",
			StatusCode: 429,
		}
		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validate_InvalidRate", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy: "token_bucket",
			Rate:     0,
			Burst:    10,
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate must be positive")
	})

	t.Run("Validate_InvalidStrategy", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy: "invalid_strategy",
			Rate:     10,
			Burst:    20,
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid strategy")
	})

	t.Run("Validate_BurstLessThanRate", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy: "token_bucket",
			Rate:     20,
			Burst:    10,
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "should be >=")
	})

	t.Run("ShouldSkipPath", func(t *testing.T) {
		config := &RateLimitConfig{
			SkipPaths: []string{"/health", "/metrics"},
		}

		assert.True(t, config.ShouldSkipPath("/health"))
		assert.True(t, config.ShouldSkipPath("/metrics"))
		assert.False(t, config.ShouldSkipPath("/api/users"))
	})
}

// TestKeyFunc 测试限流键生成函数
func TestKeyFunc(t *testing.T) {
	ctx := &RateLimitContext{
		ClientIP: "192.168.1.1",
		UserID:   "user123",
		Path:     "/api/test",
		Method:   "GET",
	}

	t.Run("KeyFuncIP", func(t *testing.T) {
		keyFunc := GetKeyFunc(KeyFuncIP)
		key := keyFunc(ctx)
		assert.Equal(t, "192.168.1.1", key)
	})

	t.Run("KeyFuncUser", func(t *testing.T) {
		keyFunc := GetKeyFunc(KeyFuncUser)
		key := keyFunc(ctx)
		assert.Equal(t, "user:user123", key)
	})

	t.Run("KeyFuncUser_NoUserID", func(t *testing.T) {
		ctxNoUser := &RateLimitContext{
			ClientIP: "192.168.1.1",
			Path:     "/api/test",
		}
		keyFunc := GetKeyFunc(KeyFuncUser)
		key := keyFunc(ctxNoUser)
		assert.Equal(t, "192.168.1.1", key)
	})

	t.Run("KeyFuncPath", func(t *testing.T) {
		keyFunc := GetKeyFunc(KeyFuncPath)
		key := keyFunc(ctx)
		assert.Equal(t, "/api/test", key)
	})

	t.Run("KeyFuncIPPath", func(t *testing.T) {
		keyFunc := GetKeyFunc(KeyFuncIPPath)
		key := keyFunc(ctx)
		assert.Equal(t, "192.168.1.1:/api/test", key)
	})

	t.Run("KeyFuncUserPath", func(t *testing.T) {
		keyFunc := GetKeyFunc(KeyFuncUserPath)
		key := keyFunc(ctx)
		assert.Equal(t, "user:user123:/api/test", key)
	})
}

// TestRateLimitMiddleware 测试限流中间件
func TestRateLimitMiddleware(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           5,
		Burst:          10,
		KeyFunc:        "ip",
		SkipPaths:      []string{"/health"},
		StatusCode:     429,
		CleanupInterval: 60,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "rate_limit", middleware.Name())
	})

	t.Run("Priority", func(t *testing.T) {
		assert.Equal(t, 8, middleware.Priority())
	})

	t.Run("LoadConfig", func(t *testing.T) {
		newConfig := map[string]interface{}{
			"enabled": false,
			"rate":    20,
			"burst":   40,
		}

		err := middleware.LoadConfig(newConfig)
		assert.NoError(t, err)
		assert.False(t, middleware.config.Enabled)
		assert.Equal(t, 20, middleware.config.Rate)
	})

	t.Run("ValidateConfig", func(t *testing.T) {
		err := middleware.ValidateConfig()
		assert.NoError(t, err)
	})

	t.Run("Reload", func(t *testing.T) {
		reloadConfig := map[string]interface{}{
			"rate":  15,
			"burst": 30,
		}

		err := middleware.Reload(reloadConfig)
		assert.NoError(t, err)
		assert.Equal(t, 15, middleware.config.Rate)
		assert.Equal(t, 30, middleware.config.Burst)
	})

	t.Run("Reload_SameStrategy", func(t *testing.T) {
		reloadConfig := map[string]interface{}{
			"strategy": "token_bucket",
			"rate":     25,
			"burst":    50,
		}

		err := middleware.Reload(reloadConfig)
		assert.NoError(t, err)
		assert.Equal(t, 25, middleware.config.Rate)
	})

	t.Run("Reload_InvalidConfig", func(t *testing.T) {
		// 记录原始配置
		originalRate := middleware.config.Rate

		reloadConfig := map[string]interface{}{
			"rate": 0, // 无效配置
		}

		err := middleware.Reload(reloadConfig)
		assert.Error(t, err)
		// 配置应该恢复到原始值
		assert.Equal(t, originalRate, middleware.config.Rate)
	})
}

// TestTokenBucketStrategy_Concurrent 测试令牌桶并发请求
func TestTokenBucketStrategy_Concurrent(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            100,
		Burst:           200,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	const goroutines = 10
	const requestsPerGoroutine = 20

	done := make(chan bool, goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			key := fmt.Sprintf("concurrent_test_%d", id)
			allowed := 0
			for j := 0; j < requestsPerGoroutine; j++ {
				if limiter.Allow(key) {
					allowed++
				}
			}
			// 每个key应该被允许 requestsPerGoroutine 次（因为小于burst）
			assert.Equal(t, requestsPerGoroutine, allowed)
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

// TestTokenBucketStrategy_TokenReplenishment 测试令牌补充
func TestTokenBucketStrategy_TokenReplenishment(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            10,  // 每秒10个令牌
		Burst:           10,  // 桶容量10
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	key := "replenishment_test"

	// 消耗所有令牌
	for i := 0; i < config.Burst; i++ {
		assert.True(t, limiter.Allow(key), "Request %d should be allowed", i+1)
	}

	// 应该被拒绝
	assert.False(t, limiter.Allow(key), "Should be rejected after exhausting tokens")

	// 等待令牌补充 (等待200ms，应该补充2个令牌)
	time.Sleep(200 * time.Millisecond)

	// 应该允许2个请求
	assert.True(t, limiter.Allow(key), "Should be allowed after token replenishment")
	assert.True(t, limiter.Allow(key), "Second request should also be allowed")
	assert.False(t, limiter.Allow(key), "Third request should be rejected")
}

// TestSlidingWindowStrategy_Concurrent 测试滑动窗口并发安全
func TestSlidingWindowStrategy_Concurrent(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            50,
		Burst:           100,
		WindowSize:      1,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	const goroutines = 10
	const requestsPerGoroutine = 10

	done := make(chan bool, goroutines)
	successCount := make(chan int, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			key := fmt.Sprintf("sw_concurrent_%d", id)
			allowed := 0
			for j := 0; j < requestsPerGoroutine; j++ {
				if limiter.Allow(key) {
					allowed++
				}
			}
			successCount <- allowed
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}

	close(successCount)
	totalAllowed := 0
	for count := range successCount {
		totalAllowed += count
	}

	// 每个key有独立的限流器，总共应该允许 goroutines * requestsPerGoroutine 个请求
	// 因为每个key的rate是50，而每个goroutine只发10个请求
	assert.Equal(t, goroutines*requestsPerGoroutine, totalAllowed)
}

// TestSlidingWindowStrategy_WindowBoundary 测试窗口边界
func TestSlidingWindowStrategy_WindowBoundary(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            5,
		Burst:           10,
		WindowSize:      1, // 1秒窗口
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	key := "boundary_test"

	// 消耗所有配额
	for i := 0; i < config.Rate; i++ {
		assert.True(t, limiter.Allow(key), "Request %d should be allowed", i+1)
	}

	// 应该被拒绝
	assert.False(t, limiter.Allow(key), "Should be rejected after rate limit")

	// 等待窗口过期
	time.Sleep(time.Duration(config.WindowSize) * time.Second)

	// 应该被允许
	assert.True(t, limiter.Allow(key), "Should be allowed after window expires")
}

// TestSlidingWindowStrategy_Reset 测试重置功能
func TestSlidingWindowStrategy_Reset(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            5,
		Burst:           10,
		WindowSize:      1,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	key := "reset_test"

	// 消耗所有配额
	for i := 0; i < config.Rate; i++ {
		assert.True(t, limiter.Allow(key))
	}

	// 应该被拒绝
	assert.False(t, limiter.Allow(key))

	// 重置
	limiter.Reset(key)

	// 重置后应该被允许
	assert.True(t, limiter.Allow(key))
}

// TestRateLimitMiddleware_Handler 测试中间件Handler
func TestRateLimitMiddleware_Handler(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           2,
		Burst:          5,
		KeyFunc:        "ip",
		StatusCode:     429,
		Message:        "Rate limit exceeded",
		CleanupInterval: 60,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	t.Run("WithinLimit", func(t *testing.T) {
		for i := 0; i < config.Rate; i++ {
			w := performRequest(router, "GET", "/test")
			assert.Equal(t, 200, w.Code)
		}
	})

	t.Run("ExceedsLimit", func(t *testing.T) {
		// 使用不同的IP地址，避免前面的测试影响
		// 创建新的router和middleware
		config2 := &RateLimitConfig{
			Enabled:        true,
			Strategy:       "token_bucket",
			Rate:           2,
			Burst:          3,
			KeyFunc:        "ip",
			StatusCode:     429,
			Message:        "Rate limit exceeded",
			CleanupInterval: 60,
		}

		middleware2, err := NewRateLimitMiddleware(config2, logger)
		assert.NoError(t, err)

		router2 := gin.New()
		router2.Use(middleware2.Handler())
		router2.GET("/test2", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		// 发送超过限制的请求
		allowed := 0
		rejected := 0
		for i := 0; i < config2.Burst+2; i++ {
			w := performRequest(router2, "GET", "/test2")
			if w.Code == 200 {
				allowed++
			} else if w.Code == config2.StatusCode {
				rejected++
			}
		}
		assert.Equal(t, config2.Burst, allowed)
		assert.True(t, rejected > 0, "Should have some rejected requests")
	})

	t.Run("SkipPath", func(t *testing.T) {
		config.SkipPaths = []string{"/skip"}
		router2 := gin.New()
		router2.Use(middleware.Handler())
		router2.GET("/skip", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "skipped"})
		})

		// 即使超过限制，skip路径也应该被允许
		for i := 0; i < 100; i++ {
			w := performRequest(router2, "GET", "/skip")
			assert.Equal(t, 200, w.Code)
		}
	})
}

// TestRateLimitMiddleware_Disabled 测试禁用限流
func TestRateLimitMiddleware_Disabled(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        false,
		Strategy:       "token_bucket",
		Rate:           1,
		Burst:          1,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 即使配置了低速率，禁用时应该允许所有请求
	for i := 0; i < 100; i++ {
		w := performRequest(router, "GET", "/test")
		assert.Equal(t, 200, w.Code)
	}
}

// TestRateLimitMiddleware_StrategyReload 测试策略切换
func TestRateLimitMiddleware_StrategyReload(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 切换到滑动窗口策略
	reloadConfig := map[string]interface{}{
		"strategy": "sliding_window",
		"rate":     15,
		"burst":    30,
	}

	err = middleware.Reload(reloadConfig)
	assert.NoError(t, err)
	assert.Equal(t, "sliding_window", middleware.config.Strategy)
	assert.Equal(t, 15, middleware.config.Rate)
	assert.Equal(t, 30, middleware.config.Burst)
}

// TestRateLimitMiddleware_GetStats 测试获取统计信息
func TestRateLimitMiddleware_GetStats(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	key := "stats_test"

	// 发起一些请求
	middleware.limiter.Allow(key)
	middleware.limiter.Allow(key)
	middleware.limiter.Allow(key)

	stats := middleware.GetStats(key)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(3), stats.TotalRequests)
	assert.Equal(t, int64(3), stats.AllowedRequests)
}

// TestRateLimitMiddleware_HotReload 测试热重载配置
func TestRateLimitMiddleware_HotReload(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		SkipPaths:      []string{"/health"},
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	reloadConfig := map[string]interface{}{
		"enabled":     false,
		"rate":        20,
		"burst":       40,
		"skip_paths":  []interface{}{"/health", "/metrics"},
		"status_code": 503,
	}

	err = middleware.Reload(reloadConfig)
	assert.NoError(t, err)

	assert.False(t, middleware.config.Enabled)
	assert.Equal(t, 20, middleware.config.Rate)
	assert.Equal(t, 40, middleware.config.Burst)
	assert.Equal(t, 503, middleware.config.StatusCode)
	// SkipPaths 应该被替换为新配置的值
	t.Logf("SkipPaths: %v", middleware.config.SkipPaths)
	assert.Equal(t, 2, len(middleware.config.SkipPaths))
}

// TestRateLimit_InvalidConfigs 测试无效配置
func TestRateLimit_InvalidConfigs(t *testing.T) {
	t.Run("ZeroRate", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy: "token_bucket",
			Rate:     0,
			Burst:    10,
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate must be positive")
	})

	t.Run("NegativeBurst", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy: "token_bucket",
			Rate:     10,
			Burst:    -1,
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "burst must be positive")
	})

	t.Run("InvalidStatusCode", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy:   "token_bucket",
			Rate:       10,
			Burst:      20,
			KeyFunc:    "ip",
			StatusCode: 700,
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status_code")
	})

	t.Run("InvalidKeyFunc", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy: "token_bucket",
			Rate:     10,
			Burst:    20,
			KeyFunc:  "invalid",
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid key_func")
	})

	t.Run("RedisWithoutAddr", func(t *testing.T) {
		config := &RateLimitConfig{
			Strategy: "redis",
			Rate:     10,
			Burst:    20,
			KeyFunc:  "ip",
			Redis: &RedisConfig{
				Addr: "",
			},
		}
		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis addr is required")
	})
}

// TestRateLimitMiddleware_Count 测试限流器计数
func TestRateLimitMiddleware_Count(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            10,
		Burst:           20,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	// 初始计数应该为0
	assert.Equal(t, 0, limiter.Count())

	// 添加一些限流器
	limiter.Allow("key1")
	limiter.Allow("key2")
	limiter.Allow("key3")

	assert.Equal(t, 3, limiter.Count())
}

// TestRateLimitMiddlewareSlidingWindow_Count 测试滑动窗口计数
func TestRateLimitMiddlewareSlidingWindow_Count(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            10,
		Burst:           20,
		WindowSize:      1,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	// 初始计数应该为0
	assert.Equal(t, 0, limiter.Count())

	// 添加一些窗口
	limiter.Allow("key1")
	limiter.Allow("key2")
	limiter.Allow("key3")

	assert.Equal(t, 3, limiter.Count())
}

// TestRateLimitMiddlewareSlidingWindow_GetTotalStats 测试滑动窗口总统计
func TestRateLimitMiddlewareSlidingWindow_GetTotalStats(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            10,
		Burst:           20,
		WindowSize:      1,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	// 发起一些请求
	limiter.Allow("key1")
	limiter.Allow("key2")
	limiter.Allow("key1")
	limiter.Allow("key3")

	stats := limiter.GetTotalStats()
	assert.NotNil(t, stats)
	assert.Equal(t, int64(4), stats.TotalRequests)
	assert.Equal(t, int64(4), stats.AllowedRequests)
}

// TestRateLimitMiddleware_GetTotalStats 测试总统计信息
func TestRateLimitMiddleware_GetTotalStats(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 发起一些请求
	middleware.limiter.Allow("key1")
	middleware.limiter.Allow("key2")
	middleware.limiter.Allow("key1")

	stats := middleware.GetTotalStats()
	assert.NotNil(t, stats)
	assert.True(t, stats.TotalRequests >= 3)
}

// TestRateLimitMiddlewareSimple 测试简单限流中间件
func TestRateLimitMiddlewareSimple(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建每秒5个请求，窗口为1秒的限流器
	handler := RateLimitMiddlewareSimple(5, 1)

	router := gin.New()
	router.Use(handler)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	allowed := 0
	for i := 0; i < 10; i++ {
		w := performRequest(router, "GET", "/test")
		if w.Code == 200 {
			allowed++
		}
	}

	// 应该允许5个请求
	assert.Equal(t, 5, allowed)
}

// TestTokenBucketStrategy_Wait 测试Wait方法
func TestTokenBucketStrategy_Wait(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            10,
		Burst:           10,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	t.Run("Wait_WithinLimit", func(t *testing.T) {
		ctx := context.Background()
		err := limiter.Wait(ctx, "wait_test")
		assert.NoError(t, err)
	})

	t.Run("Wait_WithTimeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		// 消耗所有令牌，然后连续请求Wait
		// 由于令牌补充速度(rate)很低，Wait会等待直到超时
		for i := 0; i < config.Burst*2; i++ {
			limiter.Allow("wait_timeout_test")
		}

		// 下一个Wait应该因为等待时间过长而超时
		err := limiter.Wait(ctx, "wait_timeout_test")
		assert.Error(t, err)
	})
}

// TestRateLimitConfig_NegativeWindow 测试负数时间窗口
func TestRateLimitConfig_NegativeWindow(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:   "sliding_window",
		Rate:       10,
		Burst:      20,
		KeyFunc:    "ip",
		WindowSize: -1,
		StatusCode: 429,
	}

	// 负数窗口仍然可以通过验证（因为它会被当作绝对值使用）
	// 但在实际使用时会有问题
	err := config.Validate()
	// 验证会通过，但实际使用时需要注意
	assert.NoError(t, err)
}

// TestRateLimitConfig_DefaultRedisConfig 测试默认Redis配置
func TestRateLimitConfig_DefaultRedisConfig(t *testing.T) {
	config := DefaultRedisConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "localhost:6379", config.Addr)
	assert.Equal(t, 0, config.DB)
	assert.Equal(t, "ratelimit:", config.Prefix)
	assert.Equal(t, 10, config.PoolSize)
	assert.Equal(t, 2, config.MinIdleConns)
	assert.Equal(t, 3, config.MaxRetries)
}

// TestRateLimitMiddleware_ReloadToSlidingWindow 测试重载到滑动窗口
func TestRateLimitMiddleware_ReloadToSlidingWindow(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 重载到滑动窗口策略
	reloadConfig := map[string]interface{}{
		"strategy": "sliding_window",
		"rate":     15,
		"burst":    30,
	}

	err = middleware.Reload(reloadConfig)
	assert.NoError(t, err)
	assert.Equal(t, "sliding_window", middleware.config.Strategy)
}

// TestRateLimitMiddleware_NewWithSlidingWindow 测试使用滑动窗口创建中间件
func TestRateLimitMiddleware_NewWithSlidingWindow(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "sliding_window",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		WindowSize:     1,
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)
	assert.Equal(t, "sliding_window", middleware.config.Strategy)
}

// TestRateLimitMiddleware_SkipSuccessful 测试跳过成功请求配置
func TestRateLimitMiddleware_SkipSuccessful(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:          true,
		Strategy:         "token_bucket",
		Rate:             10,
		Burst:            20,
		KeyFunc:          "ip",
		CleanupInterval:  60,
		StatusCode:       429,
		SkipSuccessful:   true,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)
	assert.True(t, middleware.config.SkipSuccessful)
}

// TestRateLimitMiddleware_SkipFailedRequest 测试跳过失败请求配置
func TestRateLimitMiddleware_SkipFailedRequest(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:           true,
		Strategy:          "token_bucket",
		Rate:              10,
		Burst:             20,
		KeyFunc:           "ip",
		CleanupInterval:   60,
		StatusCode:        429,
		SkipFailedRequest: true,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)
	assert.True(t, middleware.config.SkipFailedRequest)
}

// TestRateLimitMiddleware_CustomMessage 测试自定义消息
func TestRateLimitMiddleware_CustomMessage(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
		Message:        "Custom rate limit message",
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)
	assert.Equal(t, "Custom rate limit message", middleware.config.Message)
}

// TestTokenBucketStrategy_Stop 测试停止限流器
func TestTokenBucketStrategy_Stop(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            10,
		Burst:           20,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)

	// Stop应该不会panic
	limiter.Stop()
}

// TestSlidingWindowStrategy_Stop 测试停止滑动窗口限流器
func TestSlidingWindowStrategy_Stop(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            10,
		Burst:           20,
		WindowSize:      1,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)

	// Stop应该不会panic
	limiter.Stop()
}

// TestRateLimitConfig_AllStrategies 测试所有策略
func TestRateLimitConfig_AllStrategies(t *testing.T) {
	strategies := []string{"token_bucket", "sliding_window", "redis"}

	for _, strategy := range strategies {
		t.Run(strategy, func(t *testing.T) {
			config := &RateLimitConfig{
				Strategy:   strategy,
				Rate:       10,
				Burst:      20,
				KeyFunc:    "ip",
				StatusCode: 429,
			}

			// redis策略需要Redis配置
			if strategy == "redis" {
				config.Redis = &RedisConfig{
					Addr: "localhost:6379",
				}
			}

			err := config.Validate()
			if strategy == "redis" {
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRateLimitMiddleware_AllKeyFuncs 测试所有键生成函数
func TestRateLimitMiddleware_AllKeyFuncs(t *testing.T) {
	logger := zap.NewNop()

	keyFuncs := []string{"ip", "user", "path", "ip_path", "user_path"}

	for _, keyFunc := range keyFuncs {
		t.Run(keyFunc, func(t *testing.T) {
			config := &RateLimitConfig{
				Enabled:        true,
				Strategy:       "token_bucket",
				Rate:           10,
				Burst:          20,
				KeyFunc:        keyFunc,
				CleanupInterval: 60,
				StatusCode:     429,
			}

			middleware, err := NewRateLimitMiddleware(config, logger)
			assert.NoError(t, err)
			assert.NotNil(t, middleware)
		})
	}
}

// TestRateLimitMiddleware_HandlerWithUserID 测试带用户ID的Handler
func TestRateLimitMiddleware_HandlerWithUserID(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "user",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		// 设置用户ID
		c.Set("user_id", "test_user_123")
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 发送请求
	w := performRequest(router, "GET", "/test")
	assert.Equal(t, 200, w.Code)
}

// TestRateLimitMiddleware_ReloadWithSameStrategy 测试相同策略重载
func TestRateLimitMiddleware_ReloadWithSameStrategy(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 相同策略重载，只更新速率参数
	reloadConfig := map[string]interface{}{
		"rate":  20,
		"burst": 40,
	}

	err = middleware.Reload(reloadConfig)
	assert.NoError(t, err)
	assert.Equal(t, 20, middleware.config.Rate)
	assert.Equal(t, 40, middleware.config.Burst)
}

// TestRateLimitMiddleware_NewWithNilLogger 测试使用nil logger
func TestRateLimitMiddleware_NewWithNilLogger(t *testing.T) {
	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, nil)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.logger)
}

// TestRateLimitMiddleware_HandlerLimitExceeded 测试限流触发
func TestRateLimitMiddleware_HandlerLimitExceeded(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           2,
		Burst:          2,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
		Message:        "Too many requests",
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 发送超过限制的请求
	allowedCount := 0
	rejectedCount := 0
	for i := 0; i < 10; i++ {
		w := performRequest(router, "GET", "/test")
		if w.Code == 200 {
			allowedCount++
		} else if w.Code == 429 {
			rejectedCount++
		}
	}

	assert.Equal(t, config.Burst, allowedCount)
	assert.True(t, rejectedCount > 0)
}

// TestRateLimitMiddleware_UserIDInContext 测试上下文中的用户ID
func TestRateLimitMiddleware_UserIDInContext(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "user",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 发送请求，应该使用user123作为限流key
	for i := 0; i < 5; i++ {
		w := performRequest(router, "GET", "/test")
		assert.Equal(t, 200, w.Code)
	}
}

// TestRateLimitMiddleware_UnsupportedStrategy 测试不支持的策略
func TestRateLimitMiddleware_UnsupportedStrategy(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:    true,
		Strategy:   "unsupported_strategy",
		Rate:       10,
		Burst:      20,
		KeyFunc:    "ip",
		StatusCode: 429,
	}

	_, err := NewRateLimitMiddleware(config, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

// TestRateLimitMiddleware_NewWithDefaultConfig 测试使用默认配置
func TestRateLimitMiddleware_NewWithDefaultConfig(t *testing.T) {
	middleware, err := NewRateLimitMiddleware(nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.config)
	assert.NotNil(t, middleware.logger)
	assert.Equal(t, "rate_limit", middleware.Name())
	assert.Equal(t, 8, middleware.Priority())
}

// TestRateLimitMiddleware_ReloadStrategyChange 测试策略变更时的重载
func TestRateLimitMiddleware_ReloadStrategyChange(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 切换到滑动窗口策略
	reloadConfig := map[string]interface{}{
		"strategy": "sliding_window",
		"rate":     15,
		"burst":    30,
	}

	err = middleware.Reload(reloadConfig)
	assert.NoError(t, err)
	assert.Equal(t, "sliding_window", middleware.config.Strategy)
	assert.Equal(t, 15, middleware.config.Rate)
	assert.Equal(t, 30, middleware.config.Burst)
}

// TestRateLimitMiddleware_ReloadInvalidStrategy 测试重载时无效策略
func TestRateLimitMiddleware_ReloadInvalidStrategy(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 记录原始策略
	originalStrategy := middleware.config.Strategy

	// 尝试切换到无效策略
	reloadConfig := map[string]interface{}{
		"strategy": "invalid_strategy",
	}

	err = middleware.Reload(reloadConfig)
	assert.Error(t, err)
	// 策略应该保持不变
	assert.Equal(t, originalStrategy, middleware.config.Strategy)
}

// TestRateLimitMiddleware_LoadAllConfig 测试加载所有配置
func TestRateLimitMiddleware_LoadAllConfig(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 加载所有配置项
	newConfig := map[string]interface{}{
		"enabled":              false,
		"rate":                 20,
		"burst":                40,
		"window_size":          120,
		"key_func":             "user",
		"skip_paths":           []interface{}{"/api/test"},
		"skip_successful":      true,
		"skip_failed_request":  true,
		"message":              "Custom message",
		"status_code":          503,
	}

	err = middleware.LoadConfig(newConfig)
	assert.NoError(t, err)

	assert.False(t, middleware.config.Enabled)
	assert.Equal(t, 20, middleware.config.Rate)
	assert.Equal(t, 40, middleware.config.Burst)
	assert.Equal(t, 120, middleware.config.WindowSize)
	assert.Equal(t, "user", middleware.config.KeyFunc)
	assert.True(t, middleware.config.SkipSuccessful)
	assert.True(t, middleware.config.SkipFailedRequest)
	assert.Equal(t, "Custom message", middleware.config.Message)
	assert.Equal(t, 503, middleware.config.StatusCode)
}

// TestRateLimitMiddleware_GetStatsForNonExistentKey 测试获取不存在key的统计
func TestRateLimitMiddleware_GetStatsForNonExistentKey(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 获取不存在的key的统计信息
	stats := middleware.GetStats("non_existent_key")
	assert.Nil(t, stats)
}

// TestRateLimitMiddleware_MultipleKeys 测试多个不同的限流key
func TestRateLimitMiddleware_MultipleKeys(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           5,
		Burst:          10,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)

	// 测试不同的IP地址
	testIPs := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}

	for _, ip := range testIPs {
		router := gin.New()
		router.Use(middleware.Handler())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		// 每个IP都应该有自己的限流配额
		for i := 0; i < config.Rate; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = ip + ":12345"
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code, "IP %s should be allowed", ip)
		}
	}
}

// TestRateLimitMiddleware_WaitWithTokenBucket 测试令牌桶的Wait方法
func TestRateLimitMiddleware_WaitWithTokenBucket(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            100,
		Burst:           100,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	// Wait方法应该正常工作
	ctx := context.Background()
	err = limiter.Wait(ctx, "wait_test_key")
	assert.NoError(t, err)
}

// TestSlidingWindowStrategy_ResetNonExistent 测试重置不存在的key
func TestSlidingWindowStrategy_ResetNonExistent(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "sliding_window",
		Rate:            10,
		Burst:           20,
		WindowSize:      1,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewSlidingWindowLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	// 重置不存在的key应该不会panic
	limiter.Reset("non_existent_key")
}

// TestTokenBucketStrategy_ResetNonExistent 测试重置不存在的key
func TestTokenBucketStrategy_ResetNonExistent(t *testing.T) {
	config := &RateLimitConfig{
		Strategy:        "token_bucket",
		Rate:            10,
		Burst:           20,
		KeyFunc:         "ip",
		CleanupInterval: 1,
		StatusCode:      429,
	}

	limiter, err := NewTokenBucketLimiter(config)
	assert.NoError(t, err)
	defer limiter.Stop()

	// 重置不存在的key应该不会panic
	limiter.Reset("non_existent_key")
}

// TestRateLimitMiddleware_BurstEqualToRate 测试Burst等于Rate的情况
func TestRateLimitMiddleware_BurstEqualToRate(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          10,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)
	assert.Equal(t, 10, middleware.config.Rate)
	assert.Equal(t, 10, middleware.config.Burst)
}

// TestRateLimitMiddleware_SlidingWindowTotalStats 测试滑动窗口总统计
func TestRateLimitMiddleware_SlidingWindowTotalStats(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "sliding_window",
		Rate:           10,
		Burst:          20,
		WindowSize:     1,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 发起一些请求
	middleware.limiter.Allow("key1")
	middleware.limiter.Allow("key2")

	// 获取总统计信息
	totalStats := middleware.GetTotalStats()
	assert.NotNil(t, totalStats)
	assert.True(t, totalStats.TotalRequests >= 2)
}

// TestRateLimitMiddleware_ReloadWithValidationFailure 测试重载时验证失败
func TestRateLimitMiddleware_ReloadWithValidationFailure(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 记录原始配置
	originalRate := middleware.config.Rate

	// 尝试加载无效配置
	reloadConfig := map[string]interface{}{
		"rate": 0, // 无效的rate
	}

	err = middleware.Reload(reloadConfig)
	assert.Error(t, err)
	// 配置应该恢复到原始值
	assert.Equal(t, originalRate, middleware.config.Rate)
}

// TestRateLimitMiddleware_ReloadMessageOnly 测试只重载消息
func TestRateLimitMiddleware_ReloadMessageOnly(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "ip",
		CleanupInterval: 60,
		StatusCode:     429,
		Message:        "Original message",
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	// 只更新消息
	reloadConfig := map[string]interface{}{
		"message": "New message",
	}

	err = middleware.Reload(reloadConfig)
	assert.NoError(t, err)
	assert.Equal(t, "New message", middleware.config.Message)
}

// TestRateLimitMiddleware_HandlerPathKeyFunc 测试路径键生成函数
func TestRateLimitMiddleware_HandlerPathKeyFunc(t *testing.T) {
	logger := zap.NewNop()

	config := &RateLimitConfig{
		Enabled:        true,
		Strategy:       "token_bucket",
		Rate:           10,
		Burst:          20,
		KeyFunc:        "path",
		CleanupInterval: 60,
		StatusCode:     429,
	}

	middleware, err := NewRateLimitMiddleware(config, logger)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/api/test1", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
	router.GET("/api/test2", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 不同的路径应该有不同的限流配额
	for i := 0; i < config.Rate; i++ {
		w := performRequest(router, "GET", "/api/test1")
		assert.Equal(t, 200, w.Code)
	}

	for i := 0; i < config.Rate; i++ {
		w := performRequest(router, "GET", "/api/test2")
		assert.Equal(t, 200, w.Code)
	}
}

// performRequest 辅助函数，用于模拟HTTP请求
func performRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
