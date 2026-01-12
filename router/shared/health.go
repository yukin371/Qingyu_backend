package shared

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string            `json:"status"`    // overall status: healthy, degraded, unhealthy
	Timestamp string            `json:"timestamp"` // ISO 8601 timestamp
	Services  map[string]string `json:"services"`  // individual service status
}

// ServiceStatus 服务状态详情
type ServiceStatus struct {
	Status  string `json:"status"`            // up, down, degraded
	Message string `json:"message,omitempty"` // error message if down
}

// RegisterHealthCheckRoutes 注册健康检查路由
func RegisterHealthCheckRoutes(r *gin.RouterGroup, redisClient interface{}, mongoClient interface{}) {
	r.GET("/health", func(c *gin.Context) {
		response := HealthCheckResponse{
			Status:    "unknown",
			Timestamp: time.Now().Format(time.RFC3339),
			Services:  make(map[string]string),
		}

		// 检查Redis
		if redisClient != nil {
			// 尝试调用Ping方法（Redis客户端或内存黑名单都有Ping方法）
			if pinger, ok := redisClient.(interface{ Ping(interface{}) error }); ok {
				if err := pinger.Ping(c.Request.Context()); err != nil {
					response.Services["redis"] = "down"
					zap.L().Warn("Redis health check failed", zap.Error(err))
				} else {
					response.Services["redis"] = "up"
				}
			} else {
				response.Services["redis"] = "unknown"
			}
		} else {
			response.Services["redis"] = "not configured"
		}

		// 检查MongoDB
		if mongoClient != nil {
			// 尝试调用Ping方法
			if pinger, ok := mongoClient.(interface {
				Ping(interface{}, ...interface{}) error
			}); ok {
				if err := pinger.Ping(c.Request.Context(), nil); err != nil {
					response.Services["mongodb"] = "down"
					zap.L().Warn("MongoDB health check failed", zap.Error(err))
				} else {
					response.Services["mongodb"] = "up"
				}
			} else {
				response.Services["mongodb"] = "unknown"
			}
		} else {
			response.Services["mongodb"] = "not configured"
		}

		// 计算整体状态
		upCount := 0
		downCount := 0
		for _, status := range response.Services {
			switch status {
			case "up":
				upCount++
			case "down":
				downCount++
			}
		}

		// 根据服务状态确定整体状态
		if downCount > 0 {
			response.Status = "unhealthy"
			c.JSON(http.StatusServiceUnavailable, response)
		} else if upCount > 0 {
			response.Status = "healthy"
			c.JSON(http.StatusOK, response)
		} else {
			response.Status = "degraded" // 没有任何服务在运行
			c.JSON(http.StatusServiceUnavailable, response)
		}
	})

	// 简单的存活检查（不依赖任何服务）
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 就绪检查（检查关键服务是否可用）
	r.GET("/ready", func(c *gin.Context) {
		// MongoDB必须可用才能认为就绪
		if mongoClient != nil {
			if pinger, ok := mongoClient.(interface {
				Ping(interface{}, ...interface{}) error
			}); ok {
				if err := pinger.Ping(c.Request.Context(), nil); err != nil {
					c.JSON(http.StatusServiceUnavailable, gin.H{
						"status": "not ready",
						"reason": "MongoDB is not available",
					})
					return
				}
			}
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"reason": "MongoDB is not configured",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
}
