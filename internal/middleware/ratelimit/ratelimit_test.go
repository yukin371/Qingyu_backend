package ratelimit

import (
	"context"
	"testing"
	"time"

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
