package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityConfig 安全中间件配置
type SecurityConfig struct {
	// XSS 防护
	XSSProtection      bool   `json:"xss_protection" yaml:"xss_protection"`
	ContentTypeNoSniff bool   `json:"content_type_nosniff" yaml:"content_type_nosniff"`
	FrameOptions       string `json:"frame_options" yaml:"frame_options"` // DENY, SAMEORIGIN, ALLOW-FROM

	// HSTS (HTTP Strict Transport Security)
	HSTSMaxAge            int  `json:"hsts_max_age" yaml:"hsts_max_age"`
	HSTSIncludeSubdomains bool `json:"hsts_include_subdomains" yaml:"hsts_include_subdomains"`
	HSTSPreload           bool `json:"hsts_preload" yaml:"hsts_preload"`

	// CSP (Content Security Policy)
	ContentSecurityPolicy string `json:"content_security_policy" yaml:"content_security_policy"`

	// CSRF 防护
	CSRFProtection bool     `json:"csrf_protection" yaml:"csrf_protection"`
	CSRFTokenName  string   `json:"csrf_token_name" yaml:"csrf_token_name"`
	CSRFCookieName string   `json:"csrf_cookie_name" yaml:"csrf_cookie_name"`
	CSRFSkipPaths  []string `json:"csrf_skip_paths" yaml:"csrf_skip_paths"`

	// 其他安全头
	ReferrerPolicy      string `json:"referrer_policy" yaml:"referrer_policy"`
	PermissionsPolicy   string `json:"permissions_policy" yaml:"permissions_policy"`
	CrossOriginEmbedder string `json:"cross_origin_embedder" yaml:"cross_origin_embedder"`
	CrossOriginOpener   string `json:"cross_origin_opener" yaml:"cross_origin_opener"`
	CrossOriginResource string `json:"cross_origin_resource" yaml:"cross_origin_resource"`

	// 自定义头
	CustomHeaders map[string]string `json:"custom_headers" yaml:"custom_headers"`

	// 移除的头
	RemoveHeaders []string `json:"remove_headers" yaml:"remove_headers"`
}

// DefaultSecurityConfig 默认安全配置
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		XSSProtection:         true,
		ContentTypeNoSniff:    true,
		FrameOptions:          "DENY",
		HSTSMaxAge:            31536000, // 1年
		HSTSIncludeSubdomains: true,
		HSTSPreload:           false,
		ContentSecurityPolicy: "default-src 'self'",
		CSRFProtection:        true,
		CSRFTokenName:         "X-CSRF-Token",
		CSRFCookieName:        "_csrf",
		CSRFSkipPaths:         []string{"/health", "/metrics", "/api/auth/login"},
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		PermissionsPolicy:     "geolocation=(), microphone=(), camera=()",
		CrossOriginEmbedder:   "require-corp",
		CrossOriginOpener:     "same-origin",
		CrossOriginResource:   "same-origin",
		CustomHeaders:         make(map[string]string),
		RemoveHeaders:         []string{"Server", "X-Powered-By"},
	}
}

// Security 默认安全中间件
func Security() gin.HandlerFunc {
	return SecurityWithConfig(DefaultSecurityConfig())
}

// SecurityWithConfig 带配置的安全中间件
func SecurityWithConfig(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		setSecurityHeaders(c, config)

		// 移除指定头
		removeHeaders(c, config.RemoveHeaders)

		// CSRF 防护
		if config.CSRFProtection && !shouldSkipCSRF(c.Request.URL.Path, config.CSRFSkipPaths) {
			if !validateCSRFToken(c, config) {
				c.JSON(http.StatusForbidden, gin.H{
					"code":      40301,
					"message":   "CSRF token validation failed",
					"timestamp": time.Now().Unix(),
					"data":      nil,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// setSecurityHeaders 设置安全头
func setSecurityHeaders(c *gin.Context, config SecurityConfig) {
	// XSS 防护
	if config.XSSProtection {
		c.Header("X-XSS-Protection", "1; mode=block")
	}

	// 内容类型嗅探防护
	if config.ContentTypeNoSniff {
		c.Header("X-Content-Type-Options", "nosniff")
	}

	// 点击劫持防护
	if config.FrameOptions != "" {
		c.Header("X-Frame-Options", config.FrameOptions)
	}

	// HSTS
	if config.HSTSMaxAge > 0 {
		hstsValue := fmt.Sprintf("max-age=%d", config.HSTSMaxAge)
		if config.HSTSIncludeSubdomains {
			hstsValue += "; includeSubDomains"
		}
		if config.HSTSPreload {
			hstsValue += "; preload"
		}
		c.Header("Strict-Transport-Security", hstsValue)
	}

	// CSP
	if config.ContentSecurityPolicy != "" {
		c.Header("Content-Security-Policy", config.ContentSecurityPolicy)
	}

	// Referrer Policy
	if config.ReferrerPolicy != "" {
		c.Header("Referrer-Policy", config.ReferrerPolicy)
	}

	// Permissions Policy
	if config.PermissionsPolicy != "" {
		c.Header("Permissions-Policy", config.PermissionsPolicy)
	}

	// Cross-Origin 头
	if config.CrossOriginEmbedder != "" {
		c.Header("Cross-Origin-Embedder-Policy", config.CrossOriginEmbedder)
	}
	if config.CrossOriginOpener != "" {
		c.Header("Cross-Origin-Opener-Policy", config.CrossOriginOpener)
	}
	if config.CrossOriginResource != "" {
		c.Header("Cross-Origin-Resource-Policy", config.CrossOriginResource)
	}

	// 自定义头
	for key, value := range config.CustomHeaders {
		c.Header(key, value)
	}
}

// removeHeaders 移除指定头
func removeHeaders(c *gin.Context, headers []string) {
	for _, header := range headers {
		c.Header(header, "")
	}
}

// shouldSkipCSRF 检查是否应该跳过CSRF检查
func shouldSkipCSRF(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// validateCSRFToken 验证CSRF令牌
func validateCSRFToken(c *gin.Context, config SecurityConfig) bool {
	// GET、HEAD、OPTIONS 请求通常不需要CSRF保护
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
		return true
	}

	// 从请求头获取CSRF令牌
	token := c.GetHeader(config.CSRFTokenName)
	if token == "" {
		// 从表单数据获取
		token = c.PostForm(strings.ToLower(config.CSRFTokenName))
	}

	if token == "" {
		return false
	}

	// 从Cookie获取期望的令牌
	expectedToken, err := c.Cookie(config.CSRFCookieName)
	if err != nil {
		return false
	}

	return token == expectedToken
}

// GenerateCSRFToken 生成CSRF令牌
func GenerateCSRFToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SetCSRFToken 设置CSRF令牌到Cookie
func SetCSRFToken(c *gin.Context, config SecurityConfig) {
	token := GenerateCSRFToken()
	c.SetCookie(
		config.CSRFCookieName,
		token,
		3600, // 1小时
		"/",
		"",
		false, // 不要求HTTPS（开发环境）
		true,  // HttpOnly
	)
	c.Header(config.CSRFTokenName, token)
}

// CSRFToken 获取CSRF令牌的中间件
func CSRFToken() gin.HandlerFunc {
	config := DefaultSecurityConfig()
	return func(c *gin.Context) {
		SetCSRFToken(c, config)
		c.Next()
	}
}

// CSRFTokenWithConfig 带配置的CSRF令牌中间件
func CSRFTokenWithConfig(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetCSRFToken(c, config)
		c.Next()
	}
}

// CreateSecurityMiddleware 创建安全中间件（用于中间件工厂）
func CreateSecurityMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	securityConfig := DefaultSecurityConfig()

	// 解析配置
	if xssProtection, ok := config["xss_protection"].(bool); ok {
		securityConfig.XSSProtection = xssProtection
	}
	if contentTypeNoSniff, ok := config["content_type_nosniff"].(bool); ok {
		securityConfig.ContentTypeNoSniff = contentTypeNoSniff
	}
	if frameOptions, ok := config["frame_options"].(string); ok {
		securityConfig.FrameOptions = frameOptions
	}
	if hstsMaxAge, ok := config["hsts_max_age"].(int); ok {
		securityConfig.HSTSMaxAge = hstsMaxAge
	}
	if hstsIncludeSubdomains, ok := config["hsts_include_subdomains"].(bool); ok {
		securityConfig.HSTSIncludeSubdomains = hstsIncludeSubdomains
	}
	if hstsPreload, ok := config["hsts_preload"].(bool); ok {
		securityConfig.HSTSPreload = hstsPreload
	}
	if csp, ok := config["content_security_policy"].(string); ok {
		securityConfig.ContentSecurityPolicy = csp
	}
	if csrfProtection, ok := config["csrf_protection"].(bool); ok {
		securityConfig.CSRFProtection = csrfProtection
	}
	if csrfTokenName, ok := config["csrf_token_name"].(string); ok {
		securityConfig.CSRFTokenName = csrfTokenName
	}
	if csrfCookieName, ok := config["csrf_cookie_name"].(string); ok {
		securityConfig.CSRFCookieName = csrfCookieName
	}
	if csrfSkipPaths, ok := config["csrf_skip_paths"].([]string); ok {
		securityConfig.CSRFSkipPaths = csrfSkipPaths
	}
	if referrerPolicy, ok := config["referrer_policy"].(string); ok {
		securityConfig.ReferrerPolicy = referrerPolicy
	}
	if permissionsPolicy, ok := config["permissions_policy"].(string); ok {
		securityConfig.PermissionsPolicy = permissionsPolicy
	}
	if crossOriginEmbedder, ok := config["cross_origin_embedder"].(string); ok {
		securityConfig.CrossOriginEmbedder = crossOriginEmbedder
	}
	if crossOriginOpener, ok := config["cross_origin_opener"].(string); ok {
		securityConfig.CrossOriginOpener = crossOriginOpener
	}
	if crossOriginResource, ok := config["cross_origin_resource"].(string); ok {
		securityConfig.CrossOriginResource = crossOriginResource
	}
	if customHeaders, ok := config["custom_headers"].(map[string]string); ok {
		securityConfig.CustomHeaders = customHeaders
	}
	if removeHeaders, ok := config["remove_headers"].([]string); ok {
		securityConfig.RemoveHeaders = removeHeaders
	}

	return SecurityWithConfig(securityConfig), nil
}
