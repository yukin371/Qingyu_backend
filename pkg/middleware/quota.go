package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/ai"
	"Qingyu_backend/pkg/quota"
	aiService "Qingyu_backend/service/ai"
)

// QuotaMiddleware 配额中间件
// 优先使用接口注入，保持向后兼容
type QuotaMiddleware struct {
	checker      quota.Checker // 接口注入（推荐）
	quotaService *aiService.QuotaService // 具体实现（向后兼容）
}

// NewQuotaMiddleware 创建配额中间件（向后兼容）
func NewQuotaMiddleware(quotaService *aiService.QuotaService) *QuotaMiddleware {
	return &QuotaMiddleware{
		quotaService: quotaService,
	}
}

// NewQuotaMiddlewareWithChecker 使用接口创建配额中间件（推荐）
// 这是Port/Adapter模式的核心：依赖接口而非具体实现
func NewQuotaMiddlewareWithChecker(checker quota.Checker) *QuotaMiddleware {
	return &QuotaMiddleware{
		checker: checker,
	}
}

// CheckQuota 检查配额中间件
// amount: 预估消耗的配额数量（Token数或次数）
func (m *QuotaMiddleware) CheckQuota(estimatedAmount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
			c.Abort()
			return
		}

		// 优先使用接口注入的checker
		var err error
		if m.checker != nil {
			err = m.checker.Check(c.Request.Context(), userID.(string), estimatedAmount)
		} else if m.quotaService != nil {
			// 向后兼容：使用具体实现
			err = m.quotaService.CheckQuota(c.Request.Context(), userID.(string), estimatedAmount)
		} else {
			shared.Error(c, http.StatusInternalServerError, "配额检查未配置", "配额检查服务未正确初始化")
			c.Abort()
			return
		}

		// 处理配额检查结果
		if err != nil {
			m.handleQuotaError(c, err)
			return
		}

		// 配额检查通过，继续处理请求
		c.Next()
	}
}

// handleQuotaError 统一处理配额错误
func (m *QuotaMiddleware) handleQuotaError(c *gin.Context, err error) {
	// 检查是否为标准配额错误
	if err == ai.ErrQuotaExhausted || err == quota.ErrQuotaExhausted {
		shared.Error(c, http.StatusTooManyRequests, "配额已用尽", "您的AI配额已用尽，请明天再试或升级会员")
		c.Abort()
		return
	}
	if err == ai.ErrQuotaSuspended || err == quota.ErrQuotaSuspended {
		shared.Error(c, http.StatusForbidden, "配额已暂停", "您的AI配额已被暂停")
		c.Abort()
		return
	}
	if err == ai.ErrInsufficientQuota || err == quota.ErrInsufficientQuota {
		shared.Error(c, http.StatusTooManyRequests, "配额不足", "您的AI配额不足以完成此操作")
		c.Abort()
		return
	}

	// 未知错误
	shared.Error(c, http.StatusInternalServerError, "配额检查失败", err.Error())
	c.Abort()
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
