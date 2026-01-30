package ratelimit

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddlewareSimple 简单的限流中间件
// 提供与旧middleware兼容的API
// limit: 允许的请求数, window: 时间窗口（秒）
func RateLimitMiddlewareSimple(limit int, window int) gin.HandlerFunc {
	// 计算每秒速率
	rps := float64(limit) / float64(window)
	limiter := rate.NewLimiter(rate.Limit(rps), limit)

	return func(c *gin.Context) {
		// 使用IP作为限流key
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    42901,
				"message": fmt.Sprintf("请求过于频繁，每%d秒最多%d次请求", window, limit),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
