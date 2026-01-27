package builtin

import (
	"fmt"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// DefaultRequestIDHeader 默认的请求ID头名称
	DefaultRequestIDHeader = "X-Request-ID"
	// RequestIDKey 请求ID在Context中的key
	RequestIDKey = "request_id"
)

// RequestIDMiddleware 请求ID中间件
//
// 优先级: 1（最先执行，确保所有后续中间件都能获取到请求ID）
// 用途: 为每个请求生成唯一标识，便于追踪和日志关联
type RequestIDMiddleware struct {
	config *RequestIDConfig
}

// RequestIDConfig 请求ID配置
type RequestIDConfig struct {
	// HeaderName 请求ID头名称
	// 默认: "X-Request-ID"
	// 示例: "X-Request-ID", "X-Trace-ID"
	HeaderName string `yaml:"header_name"`

	// ForceGen 是否强制生成新的请求ID
	// 如果为false，当请求中已存在请求ID时，使用请求中的ID
	// 如果为true，忽略请求中的ID，始终生成新的ID
	// 默认: false
	ForceGen bool `yaml:"force_gen"`
}

// DefaultRequestIDConfig 返回默认请求ID配置
func DefaultRequestIDConfig() *RequestIDConfig {
	return &RequestIDConfig{
		HeaderName: DefaultRequestIDHeader,
		ForceGen:   false,
	}
}

// NewRequestIDMiddleware 创建新的请求ID中间件
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{
		config: DefaultRequestIDConfig(),
	}
}

// Name 返回中间件名称
func (m *RequestIDMiddleware) Name() string {
	return "request_id"
}

// Priority 返回执行优先级
//
// 返回1，确保请求ID中间件最先执行
// 这样所有后续中间件和处理器都能获取到请求ID
func (m *RequestIDMiddleware) Priority() int {
	return 1
}

// Handler 返回Gin处理函数
func (m *RequestIDMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取或生成请求ID
		requestID := m.getOrGenerateRequestID(c)

		// 设置到Context中，供后续中间件使用
		c.Set(RequestIDKey, requestID)

		// 设置到响应头中，便于客户端追踪
		c.Header(m.config.HeaderName, requestID)

		c.Next()
	}
}

// getOrGenerateRequestID 获取或生成请求ID
//
// 如果ForceGen为false，首先尝试从请求头中获取请求ID
// 如果请求头中没有或ForceGen为true，则生成新的UUID
func (m *RequestIDMiddleware) getOrGenerateRequestID(c *gin.Context) string {
	// 如果不强制生成，尝试从请求头中获取
	if !m.config.ForceGen {
		if requestID := c.GetHeader(m.config.HeaderName); requestID != "" {
			return requestID
		}
	}

	// 生成新的UUID
	return uuid.New().String()
}

// GetRequestID 从Context中获取请求ID
//
// 这是一个辅助函数，供其他中间件和处理器使用
// 如果Context中没有请求ID，返回空字符串
//
// 示例:
//
//	requestID := builtin.GetRequestID(c)
//	logger.Info("Processing request", zap.String("request_id", requestID))
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if str, ok := requestID.(string); ok {
			return str
		}
	}
	return ""
}

// LoadConfig 从配置加载参数
//
// 实现ConfigurableMiddleware接口
func (m *RequestIDMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &RequestIDConfig{}
	}

	// 加载HeaderName
	if headerName, ok := config["header_name"].(string); ok {
		m.config.HeaderName = headerName
	}

	// 加载ForceGen
	if forceGen, ok := config["force_gen"].(bool); ok {
		m.config.ForceGen = forceGen
	}

	return nil
}

// ValidateConfig 验证配置有效性
//
// 实现ConfigurableMiddleware接口
func (m *RequestIDMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultRequestIDConfig()
	}

	// 验证HeaderName
	if m.config.HeaderName == "" {
		return fmt.Errorf("header_name不能为空")
	}

	return nil
}

// 确保RequestIDMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*RequestIDMiddleware)(nil)
