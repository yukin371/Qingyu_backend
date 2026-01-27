package builtin

import (
	"fmt"
	"net/http"
	"strings"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS跨域中间件
//
// 优先级: 4（最外层基础设施，在所有业务中间件之前执行）
// 用途: 处理跨域请求，包括预检请求和实际请求
type CORSMiddleware struct {
	config *CORSConfig
}

// CORSConfig CORS配置
type CORSConfig struct {
	// AllowedOrigins 允许的源列表
	// 支持 "*" 通配符表示允许所有源
	// 示例: ["*"] 或 ["https://example.com", "https://api.example.com"]
	AllowedOrigins []string `yaml:"allowed_origins"`

	// AllowedMethods 允许的HTTP方法列表
	// 示例: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
	AllowedMethods []string `yaml:"allowed_methods"`

	// AllowedHeaders 允许的请求头列表
	// 示例: ["Origin", "Content-Type", "Authorization"]
	AllowedHeaders []string `yaml:"allowed_headers"`

	// ExposedHeaders 暴露的响应头列表
	// 示例: ["Content-Length", "X-Request-ID"]
	ExposedHeaders []string `yaml:"exposed_headers"`

	// AllowCredentials 是否允许携带凭证（Cookie、Authorization等）
	// 注意: 当设置为true时，AllowedOrigins不能使用"*"
	AllowCredentials bool `yaml:"allow_credentials"`

	// MaxAge 预检请求缓存时间（秒）
	// 示例: 86400（24小时）
	MaxAge int `yaml:"max_age"`
}

// DefaultCORSConfig 返回默认CORS配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Request-ID"},
		ExposedHeaders:   []string{"Content-Length", "X-Request-ID", "X-Response-Time"},
		AllowCredentials: true,
		MaxAge:           86400, // 24小时
	}
}

// NewCORSMiddleware 创建新的CORS中间件
func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{
		config: DefaultCORSConfig(),
	}
}

// Name 返回中间件名称
func (m *CORSMiddleware) Name() string {
	return "cors"
}

// Priority 返回执行优先级
//
// 返回4，确保CORS在所有业务中间件之前执行（最外层基础设施）
func (m *CORSMiddleware) Priority() int {
	return 4
}

// Handler 返回Gin处理函数
func (m *CORSMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查origin是否在允许列表中
		if !m.isAllowedOrigin(origin) {
			// 如果origin不在允许列表中，拒绝请求
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 设置CORS响应头
		c.Header("Access-Control-Allow-Origin", m.getAllowOrigin(origin))
		c.Header("Access-Control-Allow-Methods", strings.Join(m.config.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(m.config.AllowedHeaders, ", "))

		if len(m.config.ExposedHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(m.config.ExposedHeaders, ", "))
		}

		if m.config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if m.config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", m.config.MaxAge))
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin 检查origin是否在允许列表中
func (m *CORSMiddleware) isAllowedOrigin(origin string) bool {
	for _, allowed := range m.config.AllowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
	}
	return false
}

// getAllowOrigin 获取允许的origin值
//
// 如果允许列表中包含"*"且不要求凭证，返回"*"
// 否则返回请求的origin
func (m *CORSMiddleware) getAllowOrigin(origin string) string {
	// 如果允许所有来源且不要求凭证，返回"*"
	for _, allowed := range m.config.AllowedOrigins {
		if allowed == "*" && !m.config.AllowCredentials {
			return "*"
		}
	}
	// 否则返回具体的origin
	return origin
}

// LoadConfig 从配置加载参数
//
// 实现ConfigurableMiddleware接口
func (m *CORSMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &CORSConfig{}
	}

	// 加载AllowedOrigins
	if allowedOrigins, ok := config["allowed_origins"].([]interface{}); ok {
		m.config.AllowedOrigins = make([]string, len(allowedOrigins))
		for i, v := range allowedOrigins {
			if str, ok := v.(string); ok {
				m.config.AllowedOrigins[i] = str
			}
		}
	}

	// 加载AllowedMethods
	if allowedMethods, ok := config["allowed_methods"].([]interface{}); ok {
		m.config.AllowedMethods = make([]string, len(allowedMethods))
		for i, v := range allowedMethods {
			if str, ok := v.(string); ok {
				m.config.AllowedMethods[i] = str
			}
		}
	}

	// 加载AllowedHeaders
	if allowedHeaders, ok := config["allowed_headers"].([]interface{}); ok {
		m.config.AllowedHeaders = make([]string, len(allowedHeaders))
		for i, v := range allowedHeaders {
			if str, ok := v.(string); ok {
				m.config.AllowedHeaders[i] = str
			}
		}
	}

	// 加载ExposedHeaders
	if exposedHeaders, ok := config["exposed_headers"].([]interface{}); ok {
		m.config.ExposedHeaders = make([]string, len(exposedHeaders))
		for i, v := range exposedHeaders {
			if str, ok := v.(string); ok {
				m.config.ExposedHeaders[i] = str
			}
		}
	}

	// 加载AllowCredentials
	if allowCredentials, ok := config["allow_credentials"].(bool); ok {
		m.config.AllowCredentials = allowCredentials
	}

	// 加载MaxAge
	if maxAge, ok := config["max_age"].(int); ok {
		m.config.MaxAge = maxAge
	}

	return nil
}

// ValidateConfig 验证配置有效性
//
// 实现ConfigurableMiddleware接口
func (m *CORSMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultCORSConfig()
	}

	// 验证AllowedOrigins
	if len(m.config.AllowedOrigins) == 0 {
		return fmt.Errorf("allowed_origins不能为空")
	}

	// 验证AllowCredentials与通配符的兼容性
	if m.config.AllowCredentials {
		for _, origin := range m.config.AllowedOrigins {
			if origin == "*" {
				return fmt.Errorf("allow_credentials为true时，allowed_origins不能使用通配符'*'")
			}
		}
	}

	// 验证AllowedMethods
	if len(m.config.AllowedMethods) == 0 {
		return fmt.Errorf("allowed_methods不能为空")
	}

	// 验证MaxAge
	if m.config.MaxAge < 0 {
		return fmt.Errorf("max_age不能为负数")
	}

	return nil
}

// 确保CORSMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*CORSMiddleware)(nil)
