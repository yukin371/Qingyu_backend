package middleware

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"Qingyu_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AccessLogConfig 访问日志配置
type AccessLogConfig struct {
	SkipPaths      []string
	BodyAllowPaths []string
	EnableBody     bool
	MaxBodySize    int
	RedactKeys     []string
	Mode           string // normal|strict
}

// DefaultAccessLogConfig 默认访问日志配置
func DefaultAccessLogConfig() AccessLogConfig {
	return AccessLogConfig{
		SkipPaths:      []string{"/health", "/metrics", "/swagger"},
		BodyAllowPaths: []string{},
		EnableBody:     false,
		MaxBodySize:    2048,
		RedactKeys:     []string{"authorization", "password", "token", "cookie"},
		Mode:           "normal",
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	maxSize int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.body != nil && w.body.Len() < w.maxSize {
		remain := w.maxSize - w.body.Len()
		if remain > 0 {
			if len(b) > remain {
				_, _ = w.body.Write(b[:remain])
			} else {
				_, _ = w.body.Write(b)
			}
		}
	}
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware(cfg AccessLogConfig) gin.HandlerFunc {
	if cfg.MaxBodySize <= 0 {
		cfg.MaxBodySize = DefaultAccessLogConfig().MaxBodySize
	}

	skip := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skip[path] = true
	}

	return func(c *gin.Context) {
		// 跳过健康检查等路径
		if skip[c.FullPath()] || skip[c.Request.URL.Path] {
			c.Next()
			return
		}

		// 生成请求ID
		requestID := c.GetString("requestId")
		if requestID == "" {
			requestID = c.GetHeader("X-Request-ID")
		}
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestId", requestID)
		c.Header("X-Request-ID", requestID)

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		allowBody := cfg.EnableBody && isPathAllowed(path, cfg.BodyAllowPaths)

		var reqBody string
		if allowBody {
			reqBody = readRequestBody(c, int64(cfg.MaxBodySize))
			reqBody = redactJSONKeys(reqBody, cfg.RedactKeys)
		}

		var respWriter *responseWriter
		if allowBody {
			respWriter = &responseWriter{
				ResponseWriter: c.Writer,
				body:           bytes.NewBuffer([]byte{}),
				maxSize:        cfg.MaxBodySize,
			}
			c.Writer = respWriter
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", ip),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("log_mode", cfg.Mode),
		}

		if allowBody && reqBody != "" {
			fields = append(fields, zap.String("req_body", reqBody))
		}

		if allowBody && respWriter != nil && respWriter.body.Len() > 0 {
			respBody := redactJSONKeys(respWriter.body.String(), cfg.RedactKeys)
			fields = append(fields, zap.String("resp_body", respBody))
		}

		if userID, exists := c.Get("userId"); exists {
			fields = append(fields, zap.String("user_id", fmt.Sprint(userID)))
		}

		if strings.EqualFold(cfg.Mode, "strict") {
			violations := validateAccessLogFields(requestID, method, path)
			if len(violations) > 0 {
				logger.Error("log_policy_violation", zap.Strings("violations", violations))
			}
		}

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
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

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
				requestID := c.GetString("requestId")
				if requestID == "" {
					requestID = uuid.New().String()
				}

				logger.Error("panic recovered",
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
					zap.Any("error", err),
					zap.Stack("stack"),
				)

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

// readRequestBody 读取请求体
func readRequestBody(c *gin.Context, maxSize int64) string {
	if c.Request.Body == nil {
		return ""
	}

	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxSize))
	if err != nil {
		return ""
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

func isPathAllowed(path string, allowList []string) bool {
	if len(allowList) == 0 {
		return true
	}
	for _, allow := range allowList {
		if strings.HasPrefix(path, allow) {
			return true
		}
	}
	return false
}

func redactJSONKeys(input string, keys []string) string {
	if input == "" || len(keys) == 0 {
		return input
	}
	redacted := input
	for _, key := range keys {
		if key == "" {
			continue
		}
		re := regexp.MustCompile(`(?i)("` + regexp.QuoteMeta(key) + `"\\s*:\\s*")([^"]*)(")`)
		redacted = re.ReplaceAllString(redacted, `${1}***${3}`)
	}
	return redacted
}

func validateAccessLogFields(requestID, method, path string) []string {
	var violations []string
	if requestID == "" {
		violations = append(violations, "request_id_missing")
	}
	if method == "" {
		violations = append(violations, "method_missing")
	}
	if path == "" {
		violations = append(violations, "path_missing")
	}
	return violations
}
