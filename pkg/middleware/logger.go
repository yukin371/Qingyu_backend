package middleware

import (
	"time"

	"Qingyu_backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(skipPaths ...string) gin.HandlerFunc {
	skip := make(map[string]bool)
	for _, path := range skipPaths {
		skip[path] = true
	}

	return func(c *gin.Context) {
		// 跳过健康检查等路径
		if skip[c.FullPath()] {
			c.Next()
			return
		}

		// 生成请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestId", requestID)

		// 记录开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()

		// 构建日志字段
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", ip),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 添加用户ID（如果存在）
		if userID, exists := c.Get("userId"); exists {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		// 根据状态码选择日志级别
		switch {
		case status >= 500:
			logger.Error("server error", fields...)
		case status >= 400:
			logger.Warn("client error", fields...)
		case status >= 300:
			logger.Info("redirect", fields...)
		default:
			logger.Info("request", fields...)
		}
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header获取request_id
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 设置到context和header
		c.Set("requestId", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息
				requestID := c.GetString("requestId")
				if requestID == "" {
					requestID = uuid.New().String()
				}

				// 记录panic日志
				logger.Error("panic recovered",
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
					zap.Any("error", err),
					zap.Stack("stack"),
				)

				// 返回500错误
				c.JSON(500, gin.H{
					"code":       5000,
					"message":    "服务器内部错误",
					"request_id": requestID,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// AccessLogMiddleware 访问日志中间件（更详细的日志）
func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestId", requestID)

		// 记录请求开始
		start := time.Now()
		logger.Debug("request started",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
		)

		// 处理请求
		c.Next()

		// 记录请求结束
		latency := time.Since(start)
		log := logger.With(
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		)

		// 添加用户ID
		if userID, exists := c.Get("userId"); exists {
			log = log.With(zap.String("user_id", userID.(string)))
		}

		// 记录响应
		msg := "request completed"
		if c.Writer.Status() >= 400 {
			msg = "request failed"
		}
		log.Info(msg)
	}
}
