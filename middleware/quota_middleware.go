package middleware

import (
	"net/http"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// QuotaMiddleware 配额中间件
type QuotaMiddleware struct {
	quotaService *aiService.QuotaService
}

// NewQuotaMiddleware 创建配额中间件
func NewQuotaMiddleware(quotaService *aiService.QuotaService) *QuotaMiddleware {
	return &QuotaMiddleware{
		quotaService: quotaService,
	}
}

// CheckQuota 检查配额中间件
// amount: 预估消耗的配额数量（Token数或次数）
func (m *QuotaMiddleware) CheckQuota(estimatedAmount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("userId")
		if !exists {
			shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
			c.Abort()
			return
		}

		// 检查配额
		err := m.quotaService.CheckQuota(c.Request.Context(), userID.(string), estimatedAmount)
		if err != nil {
			if err == ai.ErrQuotaExhausted {
				shared.Error(c, http.StatusTooManyRequests, "配额已用尽", "您的AI配额已用尽，请明天再试或升级会员")
				c.Abort()
				return
			}
			if err == ai.ErrQuotaSuspended {
				shared.Error(c, http.StatusForbidden, "配额已暂停", "您的AI配额已被暂停")
				c.Abort()
				return
			}
			if err == ai.ErrInsufficientQuota {
				shared.Error(c, http.StatusTooManyRequests, "配额不足", "您的AI配额不足以完成此操作")
				c.Abort()
				return
			}

			shared.Error(c, http.StatusInternalServerError, "配额检查失败", err.Error())
			c.Abort()
			return
		}

		// 配额检查通过，继续处理请求
		c.Next()
	}
}

// ConsumeQuotaAfter 消费配额（在请求处理后）
// 这个中间件应该在实际的AI调用之后执行
func (m *QuotaMiddleware) ConsumeQuotaAfter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先处理请求
		c.Next()

		// 请求处理完成后消费配额
		// 从context中获取实际使用的Token数
		if tokensUsed, exists := c.Get("tokensUsed"); exists {
			userID, _ := c.Get("userId")
			service, _ := c.Get("aiService")
			model, _ := c.Get("aiModel")
			requestID, _ := c.Get("requestID")

			// 异步消费配额（不阻塞响应）
			go func() {
				ctx := c.Request.Context()
				_ = m.quotaService.ConsumeQuota(
					ctx,
					userID.(string),
					tokensUsed.(int),
					service.(string),
					model.(string),
					requestID.(string),
				)
			}()
		}
	}
}

// QuotaCheckMiddleware 简化版配额检查中间件
// 适用于大多数AI接口，使用默认的预估值
func QuotaCheckMiddleware(quotaService *aiService.QuotaService) gin.HandlerFunc {
	middleware := NewQuotaMiddleware(quotaService)

	// 默认预估1000 tokens（约500字）
	return middleware.CheckQuota(1000)
}

// LightQuotaCheckMiddleware 轻量级配额检查（聊天接口）
func LightQuotaCheckMiddleware(quotaService *aiService.QuotaService) gin.HandlerFunc {
	middleware := NewQuotaMiddleware(quotaService)

	// 聊天接口预估300 tokens
	return middleware.CheckQuota(300)
}

// HeavyQuotaCheckMiddleware 重量级配额检查（长文本生成）
func HeavyQuotaCheckMiddleware(quotaService *aiService.QuotaService) gin.HandlerFunc {
	middleware := NewQuotaMiddleware(quotaService)

	// 长文本生成预估3000 tokens
	return middleware.CheckQuota(3000)
}
