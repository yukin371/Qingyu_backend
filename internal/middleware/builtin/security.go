package builtin

import (
	"fmt"
	"strconv"
	"strings"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware 安全头中间件
//
// 优先级: 3（基础设施层，在CORS之后执行）
// 用途: 添加安全相关的HTTP响应头，增强应用安全性
type SecurityMiddleware struct {
	config *SecurityConfig
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	// EnableXFrameOptions 是否启用X-Frame-Options头
	// 用于防止点击劫持攻击
	// 默认: true
	EnableXFrameOptions bool `yaml:"enable_x_frame_options"`

	// XFrameOptions X-Frame-Options头的值
	// 可选值: "DENY", "SAMEORIGIN", "ALLOW-FROM uri"
	// 默认: "DENY"
	XFrameOptions string `yaml:"x_frame_options"`

	// EnableHSTS 是否启用HTTP严格传输安全（HSTS）
	// 强制客户端使用HTTPS访问
	// 默认: true
	EnableHSTS bool `yaml:"enable_hsts"`

	// HSTSMaxAge HSTS的max-age值（秒）
	// 默认: 31536000（1年）
	HSTSMaxAge int `yaml:"hsts_max_age"`

	// HSTSIncludeSubDomains 是否在所有子域上启用HSTS
	// 默认: true
	HSTSIncludeSubDomains bool `yaml:"hsts_include_subdomains"`

	// EnableCSP 是否启用内容安全策略（CSP）
	// 用于防止XSS、数据注入等攻击
	// 默认: false（需要根据实际需求配置）
	EnableCSP bool `yaml:"enable_csp"`

	// CSPContent 内容安全策略的值
	// 示例: "default-src 'self'"
	CSPContent string `yaml:"csp_content"`

	// EnableXContentTypeOptions 是否启用X-Content-Type-Options
	// 防止MIME类型嗅探
	// 默认: true
	EnableXContentTypeOptions bool `yaml:"enable_x_content_type_options"`

	// EnableXSSProtection 是否启用X-XSS-Protection
	// 启用浏览器的XSS过滤
	// 默认: true
	EnableXSSProtection bool `yaml:"enable_x_ss_protection"`

	// EnableContentSecurityPolicyReportOnly 是否仅报告CSP违规而不阻止
	// 默认: false
	EnableContentSecurityPolicyReportOnly bool `yaml:"enable_content_security_policy_report_only"`

	// EnableReferrerPolicy 是否启用Referrer-Policy
	// 控制Referer头中包含的信息
	// 默认: true
	EnableReferrerPolicy bool `yaml:"enable_referrer_policy"`

	// ReferrerPolicy Referrer-Policy的值
	// 可选值: "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", etc.
	// 默认: "strict-origin-when-cross-origin"
	ReferrerPolicy string `yaml:"referrer_policy"`

	// EnablePermissionsPolicy 是否启用Permissions-Policy（原Feature-Policy）
	// 控制浏览器功能和API
	// 默认: true
	EnablePermissionsPolicy bool `yaml:"enable_permissions_policy"`

	// PermissionsPolicyContent Permissions-Policy的值
	// 示例: "geolocation=(), microphone=()"
	PermissionsPolicyContent string `yaml:"permissions_policy_content"`

	// CustomHeaders 自定义安全头
	// 键值对形式，可以添加任意额外的安全头
	// 示例: {"X-Custom-Header": "value"}
	CustomHeaders map[string]string `yaml:"custom_headers"`
}

// DefaultSecurityConfig 返回默认安全配置
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		EnableXFrameOptions:                true,
		XFrameOptions:                      "DENY",
		EnableHSTS:                         true,
		HSTSMaxAge:                         31536000, // 1年
		HSTSIncludeSubDomains:              true,
		EnableCSP:                          false,
		CSPContent:                         "",
		EnableXContentTypeOptions:          true,
		EnableXSSProtection:                true,
		EnableContentSecurityPolicyReportOnly: false,
		EnableReferrerPolicy:               true,
		ReferrerPolicy:                     "strict-origin-when-cross-origin",
		EnablePermissionsPolicy:            true,
		PermissionsPolicyContent:           "geolocation=(), microphone=(), camera=()",
		CustomHeaders:                      make(map[string]string),
	}
}

// NewSecurityMiddleware 创建新的安全中间件
func NewSecurityMiddleware() *SecurityMiddleware {
	return &SecurityMiddleware{
		config: DefaultSecurityConfig(),
	}
}

// Name 返回中间件名称
func (m *SecurityMiddleware) Name() string {
	return "security"
}

// Priority 返回执行优先级
//
// 返回3，确保安全头在基础设施层设置
func (m *SecurityMiddleware) Priority() int {
	return 3
}

// Handler 返回Gin处理函数
func (m *SecurityMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置X-Frame-Options
		if m.config.EnableXFrameOptions {
			c.Header("X-Frame-Options", m.config.XFrameOptions)
		}

		// 设置X-Content-Type-Options
		if m.config.EnableXContentTypeOptions {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		// 设置X-XSS-Protection
		if m.config.EnableXSSProtection {
			c.Header("X-XSS-Protection", "1; mode=block")
		}

		// 设置Strict-Transport-Security（HSTS）
		if m.config.EnableHSTS && c.Request.TLS != nil {
			hstsValue := "max-age=" + strconv.Itoa(m.config.HSTSMaxAge)
			if m.config.HSTSIncludeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}

		// 设置Content-Security-Policy
		if m.config.EnableCSP && m.config.CSPContent != "" {
			headerName := "Content-Security-Policy"
			if m.config.EnableContentSecurityPolicyReportOnly {
				headerName = "Content-Security-Policy-Report-Only"
			}
			c.Header(headerName, m.config.CSPContent)
		}

		// 设置Referrer-Policy
		if m.config.EnableReferrerPolicy {
			c.Header("Referrer-Policy", m.config.ReferrerPolicy)
		}

		// 设置Permissions-Policy
		if m.config.EnablePermissionsPolicy && m.config.PermissionsPolicyContent != "" {
			c.Header("Permissions-Policy", m.config.PermissionsPolicyContent)
		}

		// 设置自定义安全头
		for key, value := range m.config.CustomHeaders {
			c.Header(key, value)
		}

		c.Next()
	}
}

// LoadConfig 从配置加载参数
//
// 实现ConfigurableMiddleware接口
func (m *SecurityMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &SecurityConfig{}
	}

	// 加载EnableXFrameOptions
	if enable, ok := config["enable_x_frame_options"].(bool); ok {
		m.config.EnableXFrameOptions = enable
	}

	// 加载XFrameOptions
	if xFrameOptions, ok := config["x_frame_options"].(string); ok {
		m.config.XFrameOptions = xFrameOptions
	}

	// 加载EnableHSTS
	if enable, ok := config["enable_hsts"].(bool); ok {
		m.config.EnableHSTS = enable
	}

	// 加载HSTSMaxAge
	if maxAge, ok := config["hsts_max_age"].(int); ok {
		m.config.HSTSMaxAge = maxAge
	}

	// 加载HSTSIncludeSubDomains
	if include, ok := config["hsts_include_subdomains"].(bool); ok {
		m.config.HSTSIncludeSubDomains = include
	}

	// 加载EnableCSP
	if enable, ok := config["enable_csp"].(bool); ok {
		m.config.EnableCSP = enable
	}

	// 加载CSPContent
	if cspContent, ok := config["csp_content"].(string); ok {
		m.config.CSPContent = cspContent
	}

	// 加载EnableXContentTypeOptions
	if enable, ok := config["enable_x_content_type_options"].(bool); ok {
		m.config.EnableXContentTypeOptions = enable
	}

	// 加载EnableXSSProtection
	if enable, ok := config["enable_x_ss_protection"].(bool); ok {
		m.config.EnableXSSProtection = enable
	}

	// 加载EnableContentSecurityPolicyReportOnly
	if enable, ok := config["enable_content_security_policy_report_only"].(bool); ok {
		m.config.EnableContentSecurityPolicyReportOnly = enable
	}

	// 加载EnableReferrerPolicy
	if enable, ok := config["enable_referrer_policy"].(bool); ok {
		m.config.EnableReferrerPolicy = enable
	}

	// 加载ReferrerPolicy
	if policy, ok := config["referrer_policy"].(string); ok {
		m.config.ReferrerPolicy = policy
	}

	// 加载EnablePermissionsPolicy
	if enable, ok := config["enable_permissions_policy"].(bool); ok {
		m.config.EnablePermissionsPolicy = enable
	}

	// 加载PermissionsPolicyContent
	if content, ok := config["permissions_policy_content"].(string); ok {
		m.config.PermissionsPolicyContent = content
	}

	// 加载CustomHeaders
	if customHeaders, ok := config["custom_headers"].(map[string]interface{}); ok {
		m.config.CustomHeaders = make(map[string]string)
		for key, value := range customHeaders {
			if str, ok := value.(string); ok {
				m.config.CustomHeaders[key] = str
			}
		}
	}

	return nil
}

// ValidateConfig 验证配置有效性
//
// 实现ConfigurableMiddleware接口
func (m *SecurityMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultSecurityConfig()
	}

	// 验证XFrameOptions的值
	if m.config.EnableXFrameOptions {
		validXFrameOptions := []string{"DENY", "SAMEORIGIN"}
		validXFrameOptions = append(validXFrameOptions, "ALLOW-FROM")

		isValid := false
		for _, valid := range validXFrameOptions {
			if m.config.XFrameOptions == valid || strings.HasPrefix(m.config.XFrameOptions, "ALLOW-FROM") {
				isValid = true
				break
			}
		}

		if !isValid {
			return fmt.Errorf("无效的x_frame_options值: %s，可选值: DENY, SAMEORIGIN, ALLOW-FROM uri", m.config.XFrameOptions)
		}
	}

	// 验证HSTSMaxAge
	if m.config.EnableHSTS && m.config.HSTSMaxAge < 0 {
		return fmt.Errorf("hsts_max_age不能为负数")
	}

	// 验证CSPContent
	if m.config.EnableCSP && m.config.CSPContent == "" {
		return fmt.Errorf("启用CSP时，csp_content不能为空")
	}

	// 验证ReferrerPolicy的值
	if m.config.EnableReferrerPolicy {
		validReferrerPolicies := []string{
			"no-referrer",
			"no-referrer-when-downgrade",
			"origin",
			"origin-when-cross-origin",
			"same-origin",
			"strict-origin",
			"strict-origin-when-cross-origin",
			"unsafe-url",
		}

		isValid := false
		for _, valid := range validReferrerPolicies {
			if m.config.ReferrerPolicy == valid {
				isValid = true
				break
			}
		}

		if !isValid {
			return fmt.Errorf("无效的referrer_policy值: %s", m.config.ReferrerPolicy)
		}
	}

	// 验证PermissionsPolicyContent
	if m.config.EnablePermissionsPolicy && m.config.PermissionsPolicyContent == "" {
		return fmt.Errorf("启用Permissions-Policy时，permissions_policy_content不能为空")
	}

	return nil
}

// 确保SecurityMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*SecurityMiddleware)(nil)
