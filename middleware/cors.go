package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS中间件配置
type CORSConfig struct {
	AllowOrigins     []string      `json:"allow_origins" yaml:"allow_origins"`
	AllowMethods     []string      `json:"allow_methods" yaml:"allow_methods"`
	AllowHeaders     []string      `json:"allow_headers" yaml:"allow_headers"`
	ExposeHeaders    []string      `json:"expose_headers" yaml:"expose_headers"`
	AllowCredentials bool          `json:"allow_credentials" yaml:"allow_credentials"`
	MaxAge           time.Duration `json:"max_age" yaml:"max_age"`
	AllowWildcard    bool          `json:"allow_wildcard" yaml:"allow_wildcard"`
	AllowBrowserExt  bool          `json:"allow_browser_ext" yaml:"allow_browser_ext"`
	AllowWebSockets  bool          `json:"allow_websockets" yaml:"allow_websockets"`
	AllowFiles       bool          `json:"allow_files" yaml:"allow_files"`
}

// DefaultCORSConfig 默认CORS配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
		AllowBrowserExt:  false,
		AllowWebSockets:  false,
		AllowFiles:       false,
	}
}

// CORS 默认CORS中间件
func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig())
}

// CORSWithConfig 带配置的CORS中间件
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否允许该源
		if !isOriginAllowed(origin, config) {
			c.Next()
			return
		}

		// 设置CORS头部
		setCORSHeaders(c, origin, config)

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			handlePreflightRequest(c, config)
			return
		}

		c.Next()
	}
}

// isOriginAllowed 检查源是否被允许
func isOriginAllowed(origin string, config CORSConfig) bool {
	if origin == "" {
		return true
	}

	// 检查通配符
	if config.AllowWildcard {
		for _, allowedOrigin := range config.AllowOrigins {
			if allowedOrigin == "*" {
				return true
			}
		}
	}

	// 检查精确匹配
	for _, allowedOrigin := range config.AllowOrigins {
		if allowedOrigin == origin {
			return true
		}
		
		// 支持子域名通配符，如 *.example.com
		if strings.HasPrefix(allowedOrigin, "*.") {
			domain := allowedOrigin[2:]
			if strings.HasSuffix(origin, "."+domain) || origin == domain {
				return true
			}
		}
	}

	// 检查浏览器扩展
	if config.AllowBrowserExt {
		if strings.HasPrefix(origin, "chrome-extension://") ||
			strings.HasPrefix(origin, "moz-extension://") ||
			strings.HasPrefix(origin, "safari-extension://") {
			return true
		}
	}

	// 检查文件协议
	if config.AllowFiles && strings.HasPrefix(origin, "file://") {
		return true
	}

	return false
}

// setCORSHeaders 设置CORS头部
func setCORSHeaders(c *gin.Context, origin string, config CORSConfig) {
	// 设置允许的源
	if origin != "" && isOriginAllowed(origin, config) {
		c.Header("Access-Control-Allow-Origin", origin)
	} else if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
		c.Header("Access-Control-Allow-Origin", "*")
	}

	// 设置允许的方法
	if len(config.AllowMethods) > 0 {
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
	}

	// 设置允许的头部
	if len(config.AllowHeaders) > 0 {
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
	}

	// 设置暴露的头部
	if len(config.ExposeHeaders) > 0 {
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
	}

	// 设置是否允许凭据
	if config.AllowCredentials {
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	// 设置预检缓存时间
	if config.MaxAge > 0 {
		c.Header("Access-Control-Max-Age", strconv.Itoa(int(config.MaxAge.Seconds())))
	}

	// WebSocket支持
	if config.AllowWebSockets {
		if c.Request.Header.Get("Upgrade") == "websocket" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
	}
}

// handlePreflightRequest 处理预检请求
func handlePreflightRequest(c *gin.Context, config CORSConfig) {
	// 检查请求的方法是否被允许
	requestMethod := c.Request.Header.Get("Access-Control-Request-Method")
	if requestMethod != "" && !isMethodAllowed(requestMethod, config.AllowMethods) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	// 检查请求的头部是否被允许
	requestHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
	if requestHeaders != "" && !areHeadersAllowed(requestHeaders, config.AllowHeaders) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// isMethodAllowed 检查方法是否被允许
func isMethodAllowed(method string, allowedMethods []string) bool {
	for _, allowedMethod := range allowedMethods {
		if strings.EqualFold(method, allowedMethod) {
			return true
		}
	}
	return false
}

// areHeadersAllowed 检查头部是否被允许
func areHeadersAllowed(requestHeaders string, allowedHeaders []string) bool {
	headers := strings.Split(requestHeaders, ",")
	for _, header := range headers {
		header = strings.TrimSpace(header)
		if !isHeaderAllowed(header, allowedHeaders) {
			return false
		}
	}
	return true
}

// isHeaderAllowed 检查单个头部是否被允许
func isHeaderAllowed(header string, allowedHeaders []string) bool {
	// 检查通配符
	for _, allowedHeader := range allowedHeaders {
		if allowedHeader == "*" {
			return true
		}
		if strings.EqualFold(header, allowedHeader) {
			return true
		}
	}
	return false
}

// CreateCORSMiddleware 创建CORS中间件（用于中间件工厂）
func CreateCORSMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	corsConfig := DefaultCORSConfig()
	
	// 解析配置
	if allowOrigins, ok := config["allow_origins"].([]string); ok {
		corsConfig.AllowOrigins = allowOrigins
	}
	if allowMethods, ok := config["allow_methods"].([]string); ok {
		corsConfig.AllowMethods = allowMethods
	}
	if allowHeaders, ok := config["allow_headers"].([]string); ok {
		corsConfig.AllowHeaders = allowHeaders
	}
	if exposeHeaders, ok := config["expose_headers"].([]string); ok {
		corsConfig.ExposeHeaders = exposeHeaders
	}
	if allowCredentials, ok := config["allow_credentials"].(bool); ok {
		corsConfig.AllowCredentials = allowCredentials
	}
	if maxAge, ok := config["max_age"].(time.Duration); ok {
		corsConfig.MaxAge = maxAge
	}
	if allowWildcard, ok := config["allow_wildcard"].(bool); ok {
		corsConfig.AllowWildcard = allowWildcard
	}
	if allowBrowserExt, ok := config["allow_browser_ext"].(bool); ok {
		corsConfig.AllowBrowserExt = allowBrowserExt
	}
	if allowWebSockets, ok := config["allow_websockets"].(bool); ok {
		corsConfig.AllowWebSockets = allowWebSockets
	}
	if allowFiles, ok := config["allow_files"].(bool); ok {
		corsConfig.AllowFiles = allowFiles
	}
	
	return CORSWithConfig(corsConfig), nil
}