package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/ai"
	"Qingyu_backend/pkg/quota"
	aiService "Qingyu_backend/service/ai"
)

// QuotaMiddleware 配额中间件
// 优先使用接口注入，保持向后兼容
// 修复 P0-1: 移除对 container 的直接依赖
type QuotaMiddleware struct {
	checker quota.Checker // 配额检查器接口注入
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
// 修复：使用 checker 接口而非直接从 container 获取服务
func (m *QuotaMiddleware) CheckQuota(estimatedAmount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
			c.Abort()
			return
		}

		// 检查 checker 是否可用
		if m.checker == nil {
			// checker 不可用，返回配额未配置错误
			shared.Error(c, http.StatusServiceUnavailable, "配额检查服务未配置", "配额检查功能暂未启用")
			c.Abort()
			return
		}

		// 使用 checker 检查配额
		allowed, remaining, err := m.checker.Check(c.Request.Context(), userID.(string), estimatedAmount)

		if err != nil {
			// 检查错误类型
			if errors.Is(err, quota.ErrQuotaExhausted) || errors.Is(err, quota.ErrQuotaSuspended) {
				// 已用尽或暂停，禁止继续
				m.handleQuotaError(c, err)
				return
			}
			if errors.Is(err, quota.ErrInsufficientQuota) {
				// 配额不足
				m.handleQuotaError(c, err)
				return
			}
			// 其他错误视为内部错误
			m.handleQuotaError(c, errors.New("配额检查失败"))
			return
		}

		// 配额检查通过，继续处理请求
		c.Set("quota_allowed", allowed)
		c.Set("quota_remaining", remaining)
		c.Next()
	}
}

// handleQuotaError 统一处理配额错误
func (m *QuotaMiddleware) handleQuotaError(c *gin.Context, err error) {
	// 检查错误类型并设置合适的 HTTP 状态码
	if errors.Is(err, quota.ErrQuotaExhausted) {
		shared.Error(c, http.StatusTooManyRequests, "配额已用尽", "您的AI配额已用尽，请明天再试或升级会员")
	} else if errors.Is(err, quota.ErrQuotaSuspended) {
		shared.Error(c, http.StatusForbidden, "配额已暂停", "您的AI配额已被暂停")
	} else if errors.Is(err, quota.ErrInsufficientQuota) {
		shared.Error(c, http.StatusTooManyRequests, "配额不足", "您的AI配额不足以完成此操作")
	} else {
		shared.Error(c, http.StatusInternalServerError, "配额检查失败", err.Error())
	}
	c.Abort()
}
