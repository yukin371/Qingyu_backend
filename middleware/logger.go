package middleware

import (
	"fmt"
	"time"

	"Qingyu_backend/config"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 请求路径
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取服务器配置
		cfg := config.GlobalConfig.Server
		if cfg == nil || cfg.Mode == "release" {
			// 生产环境使用简洁的日志格式
			fmt.Printf("[%s] %s %s %d %v\n",
				end.Format("2006/01/02 - 15:04:05"),
				c.Request.Method,
				path,
				c.Writer.Status(),
				latency,
			)
		} else {
			// 开发环境使用详细的日志格式
			fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				c.Writer.Status(),
				latency,
				c.ClientIP(),
				c.Request.Method,
				path,
				c.Errors.String(),
			)
		}
	}
}
