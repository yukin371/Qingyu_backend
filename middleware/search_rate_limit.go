package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimitConfig Redis 限流配置
type RedisRateLimitConfig struct {
	// 每分钟请求限制
	RequestsPerMinute int
	// 限流键前缀
	KeyPrefix string
	// 是否跳过已认证用户（已认证用户可能有更高的限制）
	SkipAuthenticated bool
	// 自定义消息
	Message string
	// HTTP 状态码
	StatusCode int
}

// DefaultRedisRateLimitConfig 默认 Redis 限流配置
func DefaultRedisRateLimitConfig() RedisRateLimitConfig {
	return RedisRateLimitConfig{
		RequestsPerMinute: 60,
		KeyPrefix:         "search:ratelimit:",
		SkipAuthenticated: false,
		Message:           "搜索请求过于频繁，请稍后再试",
		StatusCode:        http.StatusTooManyRequests,
	}
}

// RedisRateLimitMiddleware Redis 搜索限流中间件
// 基于 Redis 实现分布式限流，支持多实例部署
func RedisRateLimitMiddleware(redisClient *redis.Client, config RedisRateLimitConfig) gin.HandlerFunc {
	if config.KeyPrefix == "" {
		config.KeyPrefix = "search:ratelimit:"
	}
	if config.Message == "" {
		config.Message = "搜索请求过于频繁，请稍后再试"
	}
	if config.StatusCode == 0 {
		config.StatusCode = http.StatusTooManyRequests
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 获取用户标识：优先使用 user_id，否则使用 IP
		key := config.KeyPrefix
		if userID, exists := c.Get("user_id"); exists && !config.SkipAuthenticated {
			key += fmt.Sprintf("user:%v", userID)
		} else {
			key += "ip:" + c.ClientIP()
		}

		// 尝试递增计数器
		current, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// Redis 出错，允许通过（降级策略）
			c.Next()
			return
		}

		// 首次请求，设置过期时间
		if current == 1 {
			// 1 分钟过期
			redisClient.Expire(ctx, key, 1*time.Minute)
		}

		// 检查是否超过限制
		if current > int64(config.RequestsPerMinute) {
			ttl, _ := redisClient.TTL(ctx, key).Result()
			retryAfter := int64(0)
			if ttl > 0 {
				retryAfter = int64(ttl.Seconds())
			}

			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerMinute))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))

			c.JSON(config.StatusCode, gin.H{
				"code":      42901,
				"message":   config.Message,
				"retry_after": retryAfter,
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}

		// 添加限流信息到响应头
		ttl, _ := redisClient.TTL(ctx, key).Result()
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", config.RequestsPerMinute-int(current)))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

		c.Next()
	}
}

// SearchRateLimit 搜索限流中间件（使用默认配置）
func SearchRateLimit(redisClient *redis.Client) gin.HandlerFunc {
	return RedisRateLimitMiddleware(redisClient, DefaultRedisRateLimitConfig())
}

// SearchRateLimitWithConfig 带配置的搜索限流中间件
func SearchRateLimitWithConfig(redisClient *redis.Client, requestsPerMinute int, message string) gin.HandlerFunc {
	config := DefaultRedisRateLimitConfig()
	config.RequestsPerMinute = requestsPerMinute
	if message != "" {
		config.Message = message
	}
	return RedisRateLimitMiddleware(redisClient, config)
}

// GetRateLimitStatus 获取限流状态（用于监控）
func GetRateLimitStatus(ctx context.Context, redisClient *redis.Client, key string, limit int) (current int64, remaining int64, resetTime time.Time, err error) {
	current, err = redisClient.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		return 0, 0, time.Time{}, err
	}

	if current == 0 {
		current = 0
	}

	remaining = int64(limit) - current
	if remaining < 0 {
		remaining = 0
	}

	ttl, _ := redisClient.TTL(ctx, key).Result()
	resetTime = time.Now().Add(ttl)

	return current, remaining, resetTime, nil
}
