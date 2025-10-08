package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitConfig 限流中间件配置
type RateLimitConfig struct {
	RequestsPerSecond int      `json:"requests_per_second" yaml:"requests_per_second"`
	RequestsPerMinute int      `json:"requests_per_minute" yaml:"requests_per_minute"`
	RequestsPerHour   int      `json:"requests_per_hour" yaml:"requests_per_hour"`
	BurstSize         int      `json:"burst_size" yaml:"burst_size"`
	KeyFunc           string   `json:"key_func" yaml:"key_func"`
	SkipSuccessful    bool     `json:"skip_successful" yaml:"skip_successful"`
	SkipFailedRequest bool     `json:"skip_failed_request" yaml:"skip_failed_request"`
	SkipPaths         []string `json:"skip_paths" yaml:"skip_paths"`
	Message           string   `json:"message" yaml:"message"`
	StatusCode        int      `json:"status_code" yaml:"status_code"`
}

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
	Wait(key string) error
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	limiters map[string]*rate.Limiter
	mutex    sync.RWMutex
	rps      rate.Limit
	burst    int
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(rps int, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
}

// RateLimitMiddleware 简单的速率限制中间件
// limit: 允许的请求数, window: 时间窗口（秒）
func RateLimitMiddleware(limit int, window int) gin.HandlerFunc {
	// 计算每秒速率
	rps := float64(limit) / float64(window)
	limiter := NewTokenBucketLimiter(int(rps*100), limit) // 使用更精细的速率控制

	return func(c *gin.Context) {
		// 使用IP作为限流key
		key := c.ClientIP()

		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": fmt.Sprintf("请求过于频繁，每%d秒最多%d次请求", window, limit),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow 检查是否允许请求
func (tbl *TokenBucketLimiter) Allow(key string) bool {
	limiter := tbl.getLimiter(key)
	return limiter.Allow()
}

// Wait 等待直到可以处理请求
func (tbl *TokenBucketLimiter) Wait(key string) error {
	limiter := tbl.getLimiter(key)
	return limiter.Wait(context.Background())
}

// getLimiter 获取或创建限流器
func (tbl *TokenBucketLimiter) getLimiter(key string) *rate.Limiter {
	tbl.mutex.RLock()
	limiter, exists := tbl.limiters[key]
	tbl.mutex.RUnlock()

	if !exists {
		tbl.mutex.Lock()
		// 双重检查
		if limiter, exists = tbl.limiters[key]; !exists {
			limiter = rate.NewLimiter(tbl.rps, tbl.burst)
			tbl.limiters[key] = limiter
		}
		tbl.mutex.Unlock()
	}

	return limiter
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond: 100,
		RequestsPerMinute: 1000,
		RequestsPerHour:   10000,
		BurstSize:         200,
		KeyFunc:           "ip",
		SkipSuccessful:    false,
		SkipFailedRequest: false,
		SkipPaths:         []string{"/health", "/metrics"},
		Message:           "请求过于频繁，请稍后再试",
		StatusCode:        http.StatusTooManyRequests,
	}
}

// RateLimit 默认限流中间件
func RateLimit() gin.HandlerFunc {
	return RateLimitWithConfig(DefaultRateLimitConfig())
}

// RateLimitWithConfig 带配置的限流中间件
func RateLimitWithConfig(config RateLimitConfig) gin.HandlerFunc {
	limiter := NewTokenBucketLimiter(config.RequestsPerSecond, config.BurstSize)

	return func(c *gin.Context) {
		// 检查是否跳过限流
		if shouldSkipRateLimit(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// 获取限流键
		key := getRateLimitKey(c, config.KeyFunc)

		// 检查是否允许请求
		if !limiter.Allow(key) {
			c.JSON(config.StatusCode, gin.H{
				"code":      42901,
				"message":   config.Message,
				"timestamp": time.Now().Unix(),
				"data":      nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// shouldSkipRateLimit 检查是否应该跳过限流
func shouldSkipRateLimit(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// getRateLimitKey 获取限流键
func getRateLimitKey(c *gin.Context, keyFunc string) string {
	switch keyFunc {
	case "ip":
		return c.ClientIP()
	case "user":
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("user:%v", userID)
		}
		return c.ClientIP()
	case "path":
		return c.Request.URL.Path
	case "ip_path":
		return fmt.Sprintf("%s:%s", c.ClientIP(), c.Request.URL.Path)
	case "user_path":
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("user:%v:%s", userID, c.Request.URL.Path)
		}
		return fmt.Sprintf("%s:%s", c.ClientIP(), c.Request.URL.Path)
	default:
		return c.ClientIP()
	}
}

// CreateRateLimitMiddleware 创建限流中间件（用于中间件工厂）
func CreateRateLimitMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	rateLimitConfig := DefaultRateLimitConfig()

	// 解析配置
	if rps, ok := config["requests_per_second"].(int); ok {
		rateLimitConfig.RequestsPerSecond = rps
	}
	if rpm, ok := config["requests_per_minute"].(int); ok {
		rateLimitConfig.RequestsPerMinute = rpm
	}
	if rph, ok := config["requests_per_hour"].(int); ok {
		rateLimitConfig.RequestsPerHour = rph
	}
	if burst, ok := config["burst_size"].(int); ok {
		rateLimitConfig.BurstSize = burst
	}
	if keyFunc, ok := config["key_func"].(string); ok {
		rateLimitConfig.KeyFunc = keyFunc
	}
	if skipSuccessful, ok := config["skip_successful"].(bool); ok {
		rateLimitConfig.SkipSuccessful = skipSuccessful
	}
	if skipFailed, ok := config["skip_failed_request"].(bool); ok {
		rateLimitConfig.SkipFailedRequest = skipFailed
	}
	if skipPaths, ok := config["skip_paths"].([]string); ok {
		rateLimitConfig.SkipPaths = skipPaths
	}
	if message, ok := config["message"].(string); ok {
		rateLimitConfig.Message = message
	}
	if statusCode, ok := config["status_code"].(int); ok {
		rateLimitConfig.StatusCode = statusCode
	}

	return RateLimitWithConfig(rateLimitConfig), nil
}
