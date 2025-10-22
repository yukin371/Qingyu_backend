package middleware

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig 日志中间件配置
type LoggerConfig struct {
	EnableColor       bool          `json:"enable_color" yaml:"enable_color"`
	EnableReqBody     bool          `json:"enable_req_body" yaml:"enable_req_body"`
	EnableRespBody    bool          `json:"enable_resp_body" yaml:"enable_resp_body"`
	MaxReqBodySize    int           `json:"max_req_body_size" yaml:"max_req_body_size"`
	MaxRespBodySize   int           `json:"max_resp_body_size" yaml:"max_resp_body_size"`
	SkipPaths         []string      `json:"skip_paths" yaml:"skip_paths"`
	SlowThreshold     time.Duration `json:"slow_threshold" yaml:"slow_threshold"`
	EnablePerformance bool          `json:"enable_performance" yaml:"enable_performance"`
}

// RequestInfo 请求信息
type RequestInfo struct {
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	Query     string            `json:"query,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Body      string            `json:"body,omitempty"`
	ClientIP  string            `json:"client_ip"`
	UserAgent string            `json:"user_agent"`
	UserID    string            `json:"user_id,omitempty"`
	RequestID string            `json:"request_id"`
	Timestamp time.Time         `json:"timestamp"`
	BodySize  int64             `json:"body_size"`
}

// ResponseInfo 响应信息
type ResponseInfo struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
	BodySize   int64             `json:"body_size"`
	Duration   time.Duration     `json:"duration"`
	Error      string            `json:"error,omitempty"`
}

// responseWriter 响应写入器包装
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger 默认日志中间件
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		EnableColor:     true,
		MaxReqBodySize:  1024,
		MaxRespBodySize: 1024,
		SlowThreshold:   200 * time.Millisecond,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
	})
}

// LoggerWithConfig 带配置的日志中间件
func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	// 初始化zap日志器
	logger := initZapLogger()

	return func(c *gin.Context) {
		// 检查是否跳过日志记录
		if shouldSkipLogging(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// 记录开始时间
		start := time.Now()

		// 生成请求ID
		requestID := generateRequestID()
		c.Set("request_id", requestID)

		// 读取请求体
		var reqBody string
		if config.EnableReqBody && c.Request.Body != nil {
			reqBody = readRequestBody(c, int64(config.MaxReqBodySize))
		}

		// 包装响应写入器
		var respWriter *responseWriter
		if config.EnableRespBody {
			respWriter = &responseWriter{
				ResponseWriter: c.Writer,
				body:           bytes.NewBuffer([]byte{}),
			}
			c.Writer = respWriter
		}

		// 构建请求信息
		reqInfo := buildRequestInfo(c, reqBody, requestID, start)

		// 记录请求开始
		logger.Info("Request started",
			zap.String("request_id", requestID),
			zap.String("method", reqInfo.Method),
			zap.String("path", reqInfo.Path),
			zap.String("client_ip", reqInfo.ClientIP),
			zap.String("user_agent", reqInfo.UserAgent),
			zap.String("user_id", reqInfo.UserID),
		)

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(start)

		// 构建响应信息
		respInfo := buildResponseInfo(c, respWriter, duration, config)

		// 记录请求完成
		logLevel := determineLogLevel(respInfo.StatusCode, duration, config.SlowThreshold)

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", reqInfo.Method),
			zap.String("path", reqInfo.Path),
			zap.String("client_ip", reqInfo.ClientIP),
			zap.Int("status_code", respInfo.StatusCode),
			zap.Duration("duration", duration),
			zap.Int64("req_body_size", reqInfo.BodySize),
			zap.Int64("resp_body_size", respInfo.BodySize),
		}

		// 添加用户信息
		if reqInfo.UserID != "" {
			fields = append(fields, zap.String("user_id", reqInfo.UserID))
		}

		// 添加错误信息
		if respInfo.Error != "" {
			fields = append(fields, zap.String("error", respInfo.Error))
		}

		// 添加请求体
		if config.EnableReqBody && reqBody != "" {
			fields = append(fields, zap.String("req_body", reqBody))
		}

		// 添加响应体
		if config.EnableRespBody && respInfo.Body != "" {
			fields = append(fields, zap.String("resp_body", respInfo.Body))
		}

		// 记录日志
		switch logLevel {
		case zapcore.ErrorLevel:
			logger.Error("Request completed with error", fields...)
		case zapcore.WarnLevel:
			logger.Warn("Request completed slowly", fields...)
		default:
			logger.Info("Request completed", fields...)
		}
	}
}

// 辅助函数

// initZapLogger 初始化zap日志器
func initZapLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger: %v", err))
	}

	return logger
}

// shouldSkipLogging 检查是否应该跳过日志记录
func shouldSkipLogging(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
}

// readRequestBody 读取请求体
func readRequestBody(c *gin.Context, maxSize int64) string {
	if c.Request.Body == nil {
		return ""
	}

	// 限制读取大小
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxSize))
	if err != nil {
		return ""
	}

	// 重新设置请求体，以便后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}

// buildRequestInfo 构建请求信息
func buildRequestInfo(c *gin.Context, reqBody, requestID string, start time.Time) RequestInfo {
	userID := ""
	if user, exists := c.Get("user"); exists {
		if userCtx, ok := user.(*UserContext); ok {
			userID = userCtx.UserID
		}
	}

	return RequestInfo{
		RequestID: requestID,
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		Query:     c.Request.URL.RawQuery,
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		UserID:    userID,
		Body:      reqBody,
		BodySize:  int64(len(reqBody)),
		Timestamp: start,
	}
}

// buildResponseInfo 构建响应信息
func buildResponseInfo(c *gin.Context, respWriter *responseWriter, duration time.Duration, config LoggerConfig) ResponseInfo {
	respBody := ""
	respBodySize := int64(0)

	if config.EnableRespBody && respWriter != nil {
		respBody = respWriter.body.String()
		respBodySize = int64(len(respBody))
	}

	// 获取错误信息
	errorMsg := ""
	if len(c.Errors) > 0 {
		errorMsg = c.Errors.String()
	}

	return ResponseInfo{
		StatusCode: c.Writer.Status(),
		Body:       respBody,
		BodySize:   respBodySize,
		Duration:   duration,
		Error:      errorMsg,
	}
}

// determineLogLevel 确定日志级别
func determineLogLevel(statusCode int, duration time.Duration, slowThreshold time.Duration) zapcore.Level {
	if statusCode >= 500 {
		return zapcore.ErrorLevel
	}
	if statusCode >= 400 {
		return zapcore.WarnLevel
	}
	if duration > slowThreshold {
		return zapcore.WarnLevel
	}
	return zapcore.InfoLevel
}
