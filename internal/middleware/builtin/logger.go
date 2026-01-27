package builtin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// ResponseTimeKey 响应时间在Context中的key
	ResponseTimeKey = "response_time"
)

// LoggerMiddleware 日志中间件
//
// 优先级: 7（监控层，在基础设施之后、业务层之前）
// 用途: 记录请求和响应日志，包括请求方法、路径、状态码、耗时等信息
type LoggerMiddleware struct {
	config      *LoggerConfig
	logger      *zap.Logger
	timeFormat  string
	utc         bool
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	// SkipPaths 跳过日志记录的路径列表
	// 例如: ["/health", "/metrics"]
	// 默认: 空（记录所有请求）
	SkipPaths []string `yaml:"skip_paths"`

	// TimeFormat 时间格式
	// 例如: "2006-01-02 15:04:05"
	// 默认: ""（使用时间戳）
	TimeFormat string `yaml:"time_format"`

	// UTC 是否使用UTC时间
	// 默认: false（使用本地时间）
	UTC bool `yaml:"utc"`

	// SkipLogBody 是否跳过记录请求和响应体
	// 默认: false
	SkipLogBody bool `yaml:"skip_log_body"`

	// EnableRequestID 是否在日志中包含请求ID
	// 默认: true
	EnableRequestID bool `yaml:"enable_request_id"`

	// EnableRequestBody 是否在日志中包含请求体
	// 默认: true
	EnableRequestBody bool `yaml:"enable_request_body"`

	// EnableResponseBody 是否在日志中包含响应体
	// 默认: false（响应体可能很大）
	EnableResponseBody bool `yaml:"enable_response_body"`

	// SlowRequestThreshold 慢请求阈值（毫秒）
	// 超过此阈值的请求会被记录为WARN级别
	// 默认: 3000ms
	SlowRequestThreshold int `yaml:"slow_request_threshold"`

	// EnableColors 是否在控制台输出中使用颜色
	// 默认: true
	EnableColors bool `yaml:"enable_colors"`
}

// DefaultLoggerConfig 返回默认日志配置
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		SkipPaths:            []string{},
		TimeFormat:           "",
		UTC:                  false,
		SkipLogBody:          false,
		EnableRequestID:      true,
		EnableRequestBody:    true,
		EnableResponseBody:   false,
		SlowRequestThreshold: 3000,
		EnableColors:         true,
	}
}

// NewLoggerMiddleware 创建新的日志中间件
func NewLoggerMiddleware(logger *zap.Logger) *LoggerMiddleware {
	if logger == nil {
		// 如果没有提供logger，创建一个开发环境的logger
		logger, _ = zap.NewDevelopment()
	}

	return &LoggerMiddleware{
		config:     DefaultLoggerConfig(),
		logger:     logger,
		timeFormat: "",
		utc:        false,
	}
}

// Name 返回中间件名称
func (m *LoggerMiddleware) Name() string {
	return "logger"
}

// Priority 返回执行优先级
//
// 返回7，确保日志在监控层执行（在基础设施之后、业务层之前）
func (m *LoggerMiddleware) Priority() int {
	return 7
}

// Handler 返回Gin处理函数
func (m *LoggerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过此路径
		if m.skipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 记录开始时间
		start := time.Now()

		// 读取并记录请求体
		var requestBody string
		if m.config.EnableRequestBody && !m.config.SkipLogBody {
			requestBody = m.getRequestBody(c)
		}

		// 使用response writer包装器记录响应体
		var responseBody string
		var writer *responseWriter
		if m.config.EnableResponseBody && !m.config.SkipLogBody {
			writer = &responseWriter{
				ResponseWriter: c.Writer,
				body:          bytes.NewBufferString(""),
			}
			c.Writer = writer
		}

		// 执行后续处理
		c.Next()

		// 计算耗时
		latency := time.Since(start)
	 latencyMs := latency.Milliseconds()

		// 记录响应体
		if writer != nil {
			responseBody = writer.body.String()
		}

		// 设置响应时间到Context
		c.Set(ResponseTimeKey, latencyMs)

		// 记录日志
		m.logRequest(c, requestBody, responseBody, latency, latencyMs)
	}
}

// skipPath 检查是否跳过指定路径
func (m *LoggerMiddleware) skipPath(path string) bool {
	for _, skipPath := range m.config.SkipPaths {
		if skipPath == path {
			return true
		}
	}
	return false
}

// getRequestBody 获取请求体
func (m *LoggerMiddleware) getRequestBody(c *gin.Context) string {
	// 只对JSON和表单请求记录body
	if c.Request.Body == nil || c.Request.Body == http.NoBody {
		return ""
	}

	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "application/x-www-form-urlencoded" {
		return ""
	}

	// 读取请求体
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Sprintf("(读取失败: %v)", err)
	}

	// 恢复请求体供后续中间件使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 限制长度
	bodyStr := string(bodyBytes)
	if len(bodyStr) > 1000 {
		bodyStr = bodyStr[:1000] + "... (truncated)"
	}

	return bodyStr
}

// logRequest 记录请求日志
func (m *LoggerMiddleware) logRequest(c *gin.Context, requestBody, responseBody string, latency time.Duration, latencyMs int64) {
	// 确定日志级别
	level := zapcore.InfoLevel
	statusCode := c.Writer.Status()

	// 根据状态码和耗时确定日志级别
	switch {
	case statusCode >= 500:
		level = zapcore.ErrorLevel
	case statusCode >= 400:
		level = zapcore.WarnLevel
	case latencyMs >= int64(m.config.SlowRequestThreshold):
		level = zapcore.WarnLevel
	}

	// 构建日志字段
	fields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", statusCode),
		zap.Int64("latency_ms", latencyMs),
		zap.String("client_ip", c.ClientIP()),
		zap.String("protocol", c.Request.Proto),
	}

	// 添加请求ID（如果启用）
	if m.config.EnableRequestID {
		// 直接从Context中获取请求ID
		if requestID, exists := c.Get("request_id"); exists {
			if requestIDStr, ok := requestID.(string); ok {
				fields = append(fields, zap.String("request_id", requestIDStr))
			}
		}
	}

	// 添加查询参数
	if c.Request.URL.RawQuery != "" {
		fields = append(fields, zap.String("query", c.Request.URL.RawQuery))
	}

	// 添加请求体（如果启用且有内容）
	if m.config.EnableRequestBody && !m.config.SkipLogBody && requestBody != "" {
		fields = append(fields, zap.String("request_body", requestBody))
	}

	// 添加响应体（如果启用且有内容）
	if m.config.EnableResponseBody && !m.config.SkipLogBody && responseBody != "" {
		responseBodyStr := responseBody
		if len(responseBodyStr) > 1000 {
			responseBodyStr = responseBodyStr[:1000] + "... (truncated)"
		}
		fields = append(fields, zap.String("response_body", responseBodyStr))
	}

	// 添加User-Agent
	if userAgent := c.Request.UserAgent(); userAgent != "" {
		fields = append(fields, zap.String("user_agent", userAgent))
	}

	// 添加错误信息（如果有）
	if len(c.Errors) > 0 {
		fields = append(fields, zap.Any("errors", c.Errors.String()))
	}

	// 记录日志
	msg := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)

	switch level {
	case zapcore.DebugLevel:
		m.logger.Debug(msg, fields...)
	case zapcore.InfoLevel:
		m.logger.Info(msg, fields...)
	case zapcore.WarnLevel:
		m.logger.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		m.logger.Error(msg, fields...)
	case zapcore.FatalLevel:
		m.logger.Fatal(msg, fields...)
	case zapcore.PanicLevel:
		m.logger.Panic(msg, fields...)
	}
}

// LoadConfig 从配置加载参数
//
// 实现ConfigurableMiddleware接口
func (m *LoggerMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &LoggerConfig{}
	}

	// 加载SkipPaths
	if skipPaths, ok := config["skip_paths"].([]interface{}); ok {
		m.config.SkipPaths = make([]string, len(skipPaths))
		for i, v := range skipPaths {
			if str, ok := v.(string); ok {
				m.config.SkipPaths[i] = str
			}
		}
	}

	// 加载TimeFormat
	if timeFormat, ok := config["time_format"].(string); ok {
		m.config.TimeFormat = timeFormat
		m.timeFormat = timeFormat
	}

	// 加载UTC
	if utc, ok := config["utc"].(bool); ok {
		m.config.UTC = utc
		m.utc = utc
	}

	// 加载SkipLogBody
	if skipLogBody, ok := config["skip_log_body"].(bool); ok {
		m.config.SkipLogBody = skipLogBody
	}

	// 加载EnableRequestID
	if enableRequestID, ok := config["enable_request_id"].(bool); ok {
		m.config.EnableRequestID = enableRequestID
	}

	// 加载EnableRequestBody
	if enableRequestBody, ok := config["enable_request_body"].(bool); ok {
		m.config.EnableRequestBody = enableRequestBody
	}

	// 加载EnableResponseBody
	if enableResponseBody, ok := config["enable_response_body"].(bool); ok {
		m.config.EnableResponseBody = enableResponseBody
	}

	// 加载SlowRequestThreshold
	if threshold, ok := config["slow_request_threshold"].(int); ok {
		m.config.SlowRequestThreshold = threshold
	}

	// 加载EnableColors
	if enableColors, ok := config["enable_colors"].(bool); ok {
		m.config.EnableColors = enableColors
	}

	return nil
}

// ValidateConfig 验证配置有效性
//
// 实现ConfigurableMiddleware接口
func (m *LoggerMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultLoggerConfig()
	}

	// 验证SlowRequestThreshold
	if m.config.SlowRequestThreshold < 0 {
		return fmt.Errorf("slow_request_threshold不能为负数")
	}

	return nil
}

// SetSlowRequestThreshold 设置慢请求阈值
//
// 方便在测试中修改配置
func (m *LoggerMiddleware) SetSlowRequestThreshold(threshold int) {
	if m.config == nil {
		m.config = DefaultLoggerConfig()
	}
	m.config.SlowRequestThreshold = threshold
}

// GetConfig 获取当前配置（只读）
//
// 返回配置的副本，防止外部修改
func (m *LoggerMiddleware) GetConfig() LoggerConfig {
	if m.config == nil {
		return *DefaultLoggerConfig()
	}
	return *m.config
}

// responseWriter 响应写入器包装器
// 用于捕获响应体内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应数据
func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// WriteString 写入字符串响应数据
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 确保LoggerMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*LoggerMiddleware)(nil)
