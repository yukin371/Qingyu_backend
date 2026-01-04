package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware HTTP请求指标收集中间件
func Middleware() gin.HandlerFunc {
	metrics := GetMetrics()
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 获取请求大小
		requestSize := c.Request.ContentLength

		// 处理请求
		c.Next()

		// 计算持续时间
		duration := time.Since(start)

		// 获取响应大小
		responseSize := c.Writer.Size()
		if responseSize > 0 {
			responseSize = c.Writer.Size()
		}

		// 路径标准化（避免路径参数导致的高基数问题）
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 记录指标
		metrics.RecordHttpRequest(
			c.Request.Method,
			path,
			c.Writer.Status(),
			duration,
			requestSize,
			int64(responseSize),
		)
	}
}

// MiddlewareWithCustomLabels 自定义标签的HTTP请求指标中间件
func MiddlewareWithCustomLabels(labelFunc func(c *gin.Context) []string) gin.HandlerFunc {
	metrics := GetMetrics()
	return func(c *gin.Context) {
		start := time.Now()
		requestSize := c.Request.ContentLength

		c.Next()

		duration := time.Since(start)
		responseSize := c.Writer.Size()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		labels := []string{c.Request.Method, path, strconv.Itoa(c.Writer.Status())}
		if labelFunc != nil {
			customLabels := labelFunc(c)
			labels = append(labels, customLabels...)
		}

		metrics.httpRequestTotal.WithLabelValues(labels...).Inc()
		metrics.httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration.Seconds())
		if requestSize > 0 {
			metrics.httpRequestSize.WithLabelValues(c.Request.Method, path).Observe(float64(requestSize))
		}
		if responseSize > 0 {
			metrics.httpResponseSize.WithLabelValues(c.Request.Method, path).Observe(float64(responseSize))
		}
	}
}

// RequestSizeMiddleware 请求大小限制中间件
// 同时记录请求大小指标
func RequestSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(413, gin.H{
				"error": "request entity too large",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
