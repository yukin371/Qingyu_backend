package internalapi

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/config"
	"Qingyu_backend/pkg/logger"
)

const (
	AIServiceKeyHeader = "X-AI-Service-Key"
)

// AIAuthMiddleware AI服务认证中间件
func AIAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GlobalConfig
		if cfg == nil || cfg.AI == nil || cfg.AI.AIService == nil {
			logger.Warn("AI service config not found")
			c.JSON(500, gin.H{"error": "server configuration error"})
			c.Abort()
			return
		}

		apiKey := c.GetHeader(AIServiceKeyHeader)
		expectedKey := strings.TrimSpace(cfg.AI.AIService.InternalAPIKey)
		if apiKey == "" || expectedKey == "" || apiKey != expectedKey {
			logger.Warn("AI service authentication failed",
				zap.String("client_ip", c.ClientIP()),
				zap.Bool("has_key", apiKey != ""))
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		clientIP := c.ClientIP()
		if !isAllowedIP(clientIP, cfg.AI.AIService.AllowedIPs) {
			logger.Warn("AI service IP not in whitelist",
				zap.String("client_ip", clientIP))
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isAllowedIP 检查IP是否在白名单中
func isAllowedIP(clientIP string, allowedIPs []string) bool {
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}
	for _, allowed := range allowedIPs {
		candidate := strings.TrimSpace(allowed)
		if candidate == "" {
			continue
		}
		if strings.Contains(candidate, "/") {
			_, ipNet, err := net.ParseCIDR(candidate)
			if err != nil {
				continue
			}
			if ipNet.Contains(ip) {
				return true
			}
			continue
		}
		allowedIP := net.ParseIP(candidate)
		if allowedIP != nil && allowedIP.Equal(ip) {
			return true
		} else {
			// 精确匹配
			allowedIP := net.ParseIP(candidate)
			if allowedIP != nil && allowedIP.Equal(ip) {
				return true
			}
		}
	}
	return false
}
