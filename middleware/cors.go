package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的源（生产环境应该配置具体域名）
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// 设置CORS响应头
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID, X-Response-Time")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24小时

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORSMiddlewareWithConfig 带配置的CORS中间件
func CORSMiddlewareWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查origin是否在允许列表中
		if !isAllowedOrigin(origin, config.AllowedOrigins) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 设置CORS响应头
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", joinStrings(config.AllowedMethods))
		c.Header("Access-Control-Allow-Headers", joinStrings(config.AllowedHeaders))
		c.Header("Access-Control-Expose-Headers", joinStrings(config.ExposedHeaders))

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins   []string // 允许的源
	AllowedMethods   []string // 允许的HTTP方法
	AllowedHeaders   []string // 允许的请求头
	ExposedHeaders   []string // 暴露的响应头
	AllowCredentials bool     // 是否允许携带凭证
	MaxAge           int      // 预检请求缓存时间（秒）
}

// DefaultCORSConfig 默认CORS配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Request-ID"},
		ExposedHeaders:   []string{"Content-Length", "X-Request-ID", "X-Response-Time"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// isAllowedOrigin 检查origin是否在允许列表中
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

// joinStrings 连接字符串数组
func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += ", " + strs[i]
	}
	return result
}

// CreateCORSMiddleware 创建CORS中间件的工厂函数
func CreateCORSMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	// 可以从 config 中读取 CORS 配置
	// 暂时使用默认配置
	return CORSMiddleware(), nil
}
