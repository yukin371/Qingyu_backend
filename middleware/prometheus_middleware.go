package middleware

import (
	"strconv"
	"time"

	"Qingyu_backend/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware Prometheus指标收集中间件
func PrometheusMiddleware() gin.HandlerFunc {
	promMetrics := metrics.GetDefaultMetrics()

	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 增加活跃请求计数
		promMetrics.HTTPActiveRequests.Inc()
		defer promMetrics.HTTPActiveRequests.Dec()

		// 记录请求大小
		if c.Request.ContentLength > 0 {
			promMetrics.HTTPRequestSize.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
			).Observe(float64(c.Request.ContentLength))
		}

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(start).Seconds()

		// 获取响应状态码
		status := strconv.Itoa(c.Writer.Status())

		// 获取路径（使用FullPath而不是Request.URL.Path以避免高基数）
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 记录HTTP请求总数
		promMetrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		// 记录HTTP请求响应时间
		promMetrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)

		// 记录响应大小
		responseSize := c.Writer.Size()
		if responseSize > 0 {
			promMetrics.HTTPResponseSize.WithLabelValues(
				c.Request.Method,
				path,
			).Observe(float64(responseSize))
		}
	}
}

// RecordServiceMetrics 记录服务指标的辅助函数
func RecordServiceMetrics(serviceName, method string, duration time.Duration, err error) {
	promMetrics := metrics.GetDefaultMetrics()

	status := "success"
	if err != nil {
		status = "error"
		// 记录错误
		errorType := "unknown"
		if err != nil {
			errorType = err.Error()
			// 限制错误类型长度，避免高基数
			if len(errorType) > 50 {
				errorType = errorType[:50]
			}
		}
		promMetrics.ServiceErrors.WithLabelValues(serviceName, errorType).Inc()
	}

	// 记录调用总数
	promMetrics.ServiceCallTotal.WithLabelValues(
		serviceName,
		method,
		status,
	).Inc()

	// 记录调用响应时间
	promMetrics.ServiceCallDuration.WithLabelValues(
		serviceName,
		method,
	).Observe(duration.Seconds())
}

// RecordAIMetrics 记录AI相关指标
func RecordAIMetrics(service, model, userID string, duration time.Duration, tokensUsed int, err error) {
	promMetrics := metrics.GetDefaultMetrics()

	status := "success"
	if err != nil {
		status = "error"
	}

	// 记录AI请求总数
	promMetrics.AIRequestsTotal.WithLabelValues(
		service,
		model,
		status,
	).Inc()

	// 记录AI请求响应时间
	promMetrics.AIRequestDuration.WithLabelValues(
		service,
		model,
	).Observe(duration.Seconds())

	// 记录Token消耗
	if tokensUsed > 0 {
		promMetrics.AITokensConsumed.WithLabelValues(
			userID,
			model,
		).Add(float64(tokensUsed))
	}
}

// UpdateAIQuotaMetrics 更新AI配额指标
func UpdateAIQuotaMetrics(userID, quotaType string, total, used, remaining int) {
	promMetrics := metrics.GetDefaultMetrics()

	promMetrics.AIQuotaTotal.WithLabelValues(userID, quotaType).Set(float64(total))
	promMetrics.AIQuotaUsage.WithLabelValues(userID, quotaType).Set(float64(used))
	promMetrics.AIQuotaRemaining.WithLabelValues(userID, quotaType).Set(float64(remaining))
}

// RecordDBMetrics 记录数据库指标
func RecordDBMetrics(operation, collection string, duration time.Duration, err error) {
	promMetrics := metrics.GetDefaultMetrics()

	status := "success"
	if err != nil {
		status = "error"
		errorType := "unknown"
		if err != nil {
			errorType = err.Error()
			if len(errorType) > 50 {
				errorType = errorType[:50]
			}
		}
		promMetrics.DBErrors.WithLabelValues(operation, errorType).Inc()
	}

	// 记录查询总数
	promMetrics.DBQueryTotal.WithLabelValues(
		operation,
		collection,
		status,
	).Inc()

	// 记录查询响应时间
	promMetrics.DBQueryDuration.WithLabelValues(
		operation,
		collection,
	).Observe(duration.Seconds())
}

// RecordRedisMetrics 记录Redis指标
func RecordRedisMetrics(command string, duration time.Duration, hit bool) {
	promMetrics := metrics.GetDefaultMetrics()

	if hit {
		promMetrics.RedisHits.Inc()
	} else {
		promMetrics.RedisMisses.Inc()
	}

	promMetrics.RedisCommandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

// RecordUserActivity 记录用户活动指标
func RecordUserActivity(activityType string) {
	promMetrics := metrics.GetDefaultMetrics()

	switch activityType {
	case "registration":
		promMetrics.UserRegistrations.Inc()
	case "login":
		promMetrics.UserLogins.Inc()
	}
}

// UpdateActiveUsers 更新活跃用户数
func UpdateActiveUsers(count int) {
	promMetrics := metrics.GetDefaultMetrics()
	promMetrics.UserActiveTotal.Set(float64(count))
}

// RecordBookMetrics 记录书城指标
func RecordBookMetrics(bookID, metricType string, value ...interface{}) {
	promMetrics := metrics.GetDefaultMetrics()

	switch metricType {
	case "view":
		promMetrics.BookViewsTotal.WithLabelValues(bookID).Inc()
	case "purchase":
		promMetrics.BookPurchasesTotal.WithLabelValues(bookID).Inc()
	case "rating":
		if len(value) > 0 {
			rating := value[0].(string)
			promMetrics.BookRatingsTotal.WithLabelValues(bookID, rating).Inc()
		}
	}
}
