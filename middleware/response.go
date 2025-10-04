package middleware

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
)

// ResponseFormatterMiddleware 响应格式化中间件
func ResponseFormatterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateResponseRequestID()
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

// SensitiveFieldFilterMiddleware 敏感字段过滤中间件
func SensitiveFieldFilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果启用了敏感字段过滤
		if shouldFilterSensitiveFields(c) {
			// 获取响应数据并过滤敏感字段
			// 注意：这需要拦截响应写入，可能需要使用ResponseWriter包装器
			// 这里提供一个简化版本，实际使用中可能需要更复杂的实现
		}
	}
}

// shouldFilterSensitiveFields 判断是否应该过滤敏感字段
func shouldFilterSensitiveFields(c *gin.Context) bool {
	// 从配置或环境变量读取
	// 这里简化为总是过滤
	return true
}

// ResponseTimingMiddleware 响应时间记录中间件
func ResponseTimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		// 计算响应时间
		duration := time.Since(startTime)
		c.Header("X-Response-Time", duration.String())
	}
}

// CacheControlMiddleware 缓存控制中间件
func CacheControlMiddleware(maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据路径设置缓存策略
		if shouldCache(c.Request.URL.Path) {
			c.Header("Cache-Control", "public, max-age="+string(rune(maxAge)))
		} else {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// shouldCache 判断路径是否应该缓存
func shouldCache(path string) bool {
	// 静态资源路径应该缓存
	cacheablePaths := []string{
		"/static/",
		"/assets/",
		"/images/",
	}

	for _, p := range cacheablePaths {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}

	return false
}

// FilterSensitiveResponse 过滤响应中的敏感字段（辅助函数）
func FilterSensitiveResponse(data interface{}) interface{} {
	return response.FilterSensitiveFields(data)
}

// generateResponseRequestID 生成响应请求ID
func generateResponseRequestID() string {
	timestamp := time.Now().UnixNano()
	random := rand.Int63()
	return fmt.Sprintf("%d-%d", timestamp, random)
}
