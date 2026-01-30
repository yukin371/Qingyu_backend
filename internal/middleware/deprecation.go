package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// DeprecationConfig API废弃配置
type DeprecationConfig struct {
	Enabled     bool        `json:"enabled"`      // 是否启用废弃标记
	SunsetDate  *time.Time  `json:"sunset_date"`  // 废弃日期
	Replacement string      `json:"replacement"`  // 替代API路径
	Message     string      `json:"message"`      // 废弃警告消息
}

// DeprecationRegistry API废弃注册表
type DeprecationRegistry struct {
	endpoints map[string]map[string]*DeprecationConfig // [method][path]
}

// NewDeprecationConfig 创建废弃配置
func NewDeprecationConfig(sunsetDate *time.Time, replacement string, message string) *DeprecationConfig {
	return &DeprecationConfig{
		Enabled:     true,
		SunsetDate:  sunsetDate,
		Replacement: replacement,
		Message:     message,
	}
}

// NewDeprecationRegistry 创建废弃注册表
func NewDeprecationRegistry() *DeprecationRegistry {
	return &DeprecationRegistry{
		endpoints: make(map[string]map[string]*DeprecationConfig),
	}
}

// RegisterEndpoint 注册废弃端点
func (r *DeprecationRegistry) RegisterEndpoint(method string, path string, config *DeprecationConfig) {
	if r.endpoints[method] == nil {
		r.endpoints[method] = make(map[string]*DeprecationConfig)
	}
	r.endpoints[method][path] = config
}

// IsDeprecated 检查端点是否被废弃
func (r *DeprecationRegistry) IsDeprecated(method string, path string) (*DeprecationConfig, bool) {
	if methodEndpoints, ok := r.endpoints[method]; ok {
		if config, exists := methodEndpoints[path]; exists {
			return config, true
		}
	}
	return nil, false
}

// DeprecationMiddleware 废弃中间件
// 为API端点添加废弃响应头
func DeprecationMiddleware(config *DeprecationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config != nil && config.Enabled {
			setDeprecationHeaders(c, config)
		}
		c.Next()
	}
}

// setDeprecationHeaders 设置废弃响应头
func setDeprecationHeaders(c *gin.Context, config *DeprecationConfig) {
	// X-API-Deprecated: 标记API已废弃
	c.Header("X-API-Deprecated", "true")

	// X-API-Sunset-Date: 废弃日期
	if config.SunsetDate != nil {
		c.Header("X-API-Sunset-Date", config.SunsetDate.Format(time.RFC3339))
	}

	// X-API-Replacement: 替代API
	if config.Replacement != "" {
		c.Header("X-API-Replacement", config.Replacement)
	}

	// Warning: HTTP警告头
	// 格式: 299 - "message"
	warningMsg := config.Message
	if warningMsg == "" {
		warningMsg = "This API is deprecated"
	}
	if config.Replacement != "" {
		warningMsg += fmt.Sprintf(". Use %s instead", config.Replacement)
	}
	c.Header("Warning", fmt.Sprintf(`299 - "%s"`, warningMsg))
}

// DeprecationOption 废弃配置选项函数类型
type DeprecationOption func(*DeprecationConfig)

// WithSunsetDate 设置废弃日期选项
func WithSunsetDate(date time.Time) DeprecationOption {
	return func(c *DeprecationConfig) {
		c.SunsetDate = &date
	}
}

// WithReplacement 设置替代API选项
func WithReplacement(path string) DeprecationOption {
	return func(c *DeprecationConfig) {
		c.Replacement = path
	}
}

// WithWarningMessage 设置警告消息选项
func WithWarningMessage(msg string) DeprecationOption {
	return func(c *DeprecationConfig) {
		c.Message = msg
	}
}

// SetDeprecationHeadersWithOptions 使用选项模式设置废弃响应头
// 允许在处理函数中动态设置废弃标记
func SetDeprecationHeadersWithOptions(c *gin.Context, opts ...DeprecationOption) {
	config := &DeprecationConfig{
		Enabled: true,
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(config)
	}

	setDeprecationHeaders(c, config)
}
