package middleware

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseFormatterMiddleware 响应格式化中间件
// 生成请求ID并设置到上下文和响应头
func ResponseFormatterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置请求ID到上下文
		c.Set("request_id", requestID)

		// 设置请求ID到响应头
		c.Header("X-Request-ID", requestID)

		// 设置响应时间到上下文
		c.Set("request_time", time.Now())

		c.Next()
	}
}

// ResponseTimingMiddleware 响应时间记录中间件
// 记录请求处理时间并添加到响应头
func ResponseTimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		// 计算响应时间
		duration := time.Since(startTime)
		c.Header("X-Response-Time", duration.String())
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	timestamp := time.Now().UnixNano()
	random := rand.Int63()
	return fmt.Sprintf("%d-%d", timestamp, random)
}
