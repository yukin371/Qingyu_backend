package utils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIP 从 gin.Context 获取客户端真实 IP 地址
// 考虑代理和负载均衡器的情况
func GetClientIP(c *gin.Context) string {
	// 1. 尝试从 X-Real-IP 头获取
	ip := c.GetHeader("X-Real-IP")
	if ip != "" && isValidIP(ip) {
		return ip
	}

	// 2. 尝试从 X-Forwarded-For 头获取（可能包含多个IP，取第一个）
	ip = c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For 格式: client, proxy1, proxy2
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0])
			if isValidIP(ip) {
				return ip
			}
		}
	}

	// 3. 尝试从 CF-Connecting-IP 获取（Cloudflare）
	ip = c.GetHeader("CF-Connecting-IP")
	if ip != "" && isValidIP(ip) {
		return ip
	}

	// 4. 尝试从 True-Client-IP 获取（Akamai、Cloudflare）
	ip = c.GetHeader("True-Client-IP")
	if ip != "" && isValidIP(ip) {
		return ip
	}

	// 5. 最后从 RemoteAddr 获取
	ip = c.ClientIP()
	if ip != "" && isValidIP(ip) {
		return ip
	}

	// 6. 兜底返回 unknown
	return "unknown"
}

// isValidIP 验证 IP 地址是否有效
func isValidIP(ip string) bool {
	// 移除可能的端口号
	if strings.Contains(ip, ":") {
		ip, _, _ = net.SplitHostPort(ip)
	}

	// 验证 IP 格式
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// 排除私有 IP（可选，根据需求决定）
	// return !isPrivateIP(parsedIP)

	return true
}

// isPrivateIP 检查是否为私有 IP 地址
func isPrivateIP(ip net.IP) bool {
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"::1/128",
		"fe80::/10",
		"fc00::/7",
	}

	for _, block := range privateBlocks {
		_, subnet, _ := net.ParseCIDR(block)
		if subnet != nil && subnet.Contains(ip) {
			return true
		}
	}

	return false
}

// GetIPInfo 获取 IP 详细信息（用于日志）
func GetIPInfo(c *gin.Context) map[string]string {
	return map[string]string{
		"client_ip":         GetClientIP(c),
		"remote_addr":       c.Request.RemoteAddr,
		"x_real_ip":         c.GetHeader("X-Real-IP"),
		"x_forwarded_for":   c.GetHeader("X-Forwarded-For"),
		"cf_connecting_ip":  c.GetHeader("CF-Connecting-IP"),
		"true_client_ip":    c.GetHeader("True-Client-IP"),
		"user_agent":        c.GetHeader("User-Agent"),
		"x_forwarded_proto": c.GetHeader("X-Forwarded-Proto"),
	}
}
