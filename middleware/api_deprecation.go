package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 废弃响应头常量
const (
	HeaderDeprecated      = "X-API-Deprecated"
	HeaderSunsetDate      = "X-API-Sunset-Date"
	HeaderReplacement     = "X-API-Replacement"
	HeaderWarning         = "Warning"
)

// DeprecationConfig 废弃配置
type DeprecationConfig struct {
	Enabled     bool       `json:"enabled" yaml:"enabled"`                   // 是否启用废弃标记
	SunsetDate  *time.Time `json:"sunset_date" yaml:"sunset_date"`           // 废除日期（可选）
	Replacement string     `json:"replacement" yaml:"replacement"`           // 替代端点
	WarningMsg  string     `json:"warning_msg" yaml:"warning_msg"`           // 警告消息
}

// DeprecationMiddleware 废弃中间件
// 为已废弃的API端点添加相应的响应头
func DeprecationMiddleware(config *DeprecationConfig) gin.HandlerFunc {
	if config == nil || !config.Enabled {
		// 如果未启用废弃标记，直接放行
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// 设置废弃响应头
		SetDeprecationHeaders(c, config)

		// 记录废弃API调用日志
		logDeprecatedAPIAccess(c, config)

		c.Next()
	}
}

// SetDeprecationHeaders 设置废弃响应头
func SetDeprecationHeaders(c *gin.Context, config *DeprecationConfig) {
	// 1. 设置废弃标记
	c.Header(HeaderDeprecated, "true")

	// 2. 设置废除日期（如果有）
	if config.SunsetDate != nil {
		sunsetDate := config.SunsetDate.Format("2006-01-02")
		c.Header(HeaderSunsetDate, sunsetDate)
	}

	// 3. 设置替代端点（如果有）
	if config.Replacement != "" {
		c.Header(HeaderReplacement, config.Replacement)
	}

	// 4. 设置警告消息
	warning := buildWarningMessage(config)
	c.Header(HeaderWarning, warning)
}

// buildWarningMessage 构建警告消息
// Warning格式: 299 - "message"
// 参考: RFC 7234 Section 5.5
func buildWarningMessage(config *DeprecationConfig) string {
	message := config.WarningMsg
	if message == "" {
		message = "This API is deprecated"
	}

	if config.SunsetDate != nil {
		sunsetDate := config.SunsetDate.Format("2006-01-02")
		message = fmt.Sprintf("%s and will be removed on %s", message, sunsetDate)
		if config.Replacement != "" {
			message = fmt.Sprintf("%s. Use %s instead.", message, config.Replacement)
		}
	} else if config.Replacement != "" {
		message = fmt.Sprintf("%s. Use %s instead.", message, config.Replacement)
	}

	return fmt.Sprintf(`299 - "%s"`, message)
}

// logDeprecatedAPIAccess 记录废弃API调用日志
func logDeprecatedAPIAccess(c *gin.Context, config *DeprecationConfig) {
	logger := c.MustGet("logger").(*zap.Logger)

	logger.Warn("Deprecated API accessed",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	)

	if config.SunsetDate != nil {
		logger.Warn("API will be removed",
			zap.Time("sunset_date", *config.SunsetDate),
			zap.Int("days_remaining", int(config.SunsetDate.Sub(time.Now()).Hours()/24)),
		)
	}

	if config.Replacement != "" {
		logger.Info("Recommended replacement endpoint",
			zap.String("replacement", config.Replacement),
		)
	}
}

// DeprecationRegistry 废弃端点注册表
type DeprecationRegistry struct {
	endpoints map[string]*DeprecationConfig // key: "METHOD:path"
}

// NewDeprecationRegistry 创建废弃端点注册表
func NewDeprecationRegistry() *DeprecationRegistry {
	return &DeprecationRegistry{
		endpoints: make(map[string]*DeprecationConfig),
	}
}

// RegisterEndpoint 注册废弃端点
func (r *DeprecationRegistry) RegisterEndpoint(method, path string, config *DeprecationConfig) {
	key := buildEndpointKey(method, path)
	r.endpoints[key] = config
}

// IsDeprecated 检查端点是否已废弃
func (r *DeprecationRegistry) IsDeprecated(method, path string) (*DeprecationConfig, bool) {
	key := buildEndpointKey(method, path)
	config, exists := r.endpoints[key]
	return config, exists
}

// GetDeprecationMiddleware 获取废弃中间件
// 如果端点已废弃，返回相应的中间件；否则返回空中间件
func (r *DeprecationRegistry) GetDeprecationMiddleware(method, path string) gin.HandlerFunc {
	config, deprecated := r.IsDeprecated(method, path)
	if !deprecated {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return DeprecationMiddleware(config)
}

// buildEndpointKey 构建端点key
func buildEndpointKey(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

// NewDeprecationConfig 创建废弃配置（辅助函数）
func NewDeprecationConfig(sunsetDate *time.Time, replacement, warningMsg string) *DeprecationConfig {
	return &DeprecationConfig{
		Enabled:     true,
		SunsetDate:  sunsetDate,
		Replacement: replacement,
		WarningMsg:  warningMsg,
	}
}

// NewDeprecationConfigWithDays 创建废弃配置，指定剩余天数
func NewDeprecationConfigWithDays(days int, replacement, warningMsg string) *DeprecationConfig {
	sunsetDate := time.Now().AddDate(0, 0, days)
	return NewDeprecationConfig(&sunsetDate, replacement, warningMsg)
}

// CheckDeprecationStatus 检查废弃状态（辅助函数，用于健康检查）
func CheckDeprecationStatus(config *DeprecationConfig) map[string]interface{} {
	status := make(map[string]interface{})

	status["deprecated"] = config.Enabled

	if config.SunsetDate != nil {
		daysRemaining := int(config.SunsetDate.Sub(time.Now()).Hours()/24)
		status["sunset_date"] = config.SunsetDate.Format("2006-01-02")
		status["days_remaining"] = daysRemaining
		status["urgent"] = daysRemaining <= 30 // 距离废除少于30天视为紧急
	}

	if config.Replacement != "" {
		status["replacement"] = config.Replacement
	}

	return status
}

// SetDeprecationHeadersWithOptions 使用选项模式设置废弃响应头
// 提供更灵活的配置方式
func SetDeprecationHeadersWithOptions(c *gin.Context, options ...DeprecationOption) {
	config := &DeprecationConfig{Enabled: true}

	for _, option := range options {
		option(config)
	}

	if config.Enabled {
		SetDeprecationHeaders(c, config)
	}
}

// DeprecationOption 废弃配置选项
type DeprecationOption func(*DeprecationConfig)

// WithSunsetDate 设置废除日期
func WithSunsetDate(date time.Time) DeprecationOption {
	return func(config *DeprecationConfig) {
		config.SunsetDate = &date
	}
}

// WithReplacement 设置替代端点
func WithReplacement(replacement string) DeprecationOption {
	return func(config *DeprecationConfig) {
		config.Replacement = replacement
	}
}

// WithWarningMessage 设置警告消息
func WithWarningMessage(message string) DeprecationOption {
	return func(config *DeprecationConfig) {
		config.WarningMsg = message
	}
}

// DeprecateUntil 废除到指定日期
func DeprecateUntil(date time.Time, replacement string) DeprecationOption {
	return func(config *DeprecationConfig) {
		config.SunsetDate = &date
		config.Replacement = replacement
	}
}

// Example: 在Handler中使用
//
//	// 方式1: 使用中间件
//	deprecationConfig := NewDeprecationConfigWithDays(90, "/api/v2/users", "This API is deprecated")
//	router.GET("/users", middleware.DeprecationMiddleware(deprecationConfig), userHandler)
//
//	// 方式2: 在Handler中直接设置响应头
//	func GetUsers(c *gin.Context) {
//	    middleware.SetDeprecationHeadersWithOptions(c,
//	        middleware.WithSunsetDate(time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)),
//	        middleware.WithReplacement("/api/v2/users"),
//	        middleware.WithWarningMessage("Upgrade to v2 API"),
//	    )
//	    // ... handler logic
//	}
