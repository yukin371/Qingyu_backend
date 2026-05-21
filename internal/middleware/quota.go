package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/quota"
	"Qingyu_backend/pkg/response"
)

// QuotaMiddleware 配额中间件。
// 仅用于后端内部 HTTP 链路，因此收口到 internal/middleware。
type QuotaMiddleware struct {
	checker quota.Checker
}

// NewQuotaMiddleware 创建配额中间件。
func NewQuotaMiddleware(checker quota.Checker) *QuotaMiddleware {
	return &QuotaMiddleware{
		checker: checker,
	}
}

// QuotaCheckMiddleware 创建标准配额检查中间件（估计消耗 1000 个 Token）。
func QuotaCheckMiddleware(checker quota.Checker) gin.HandlerFunc {
	return NewQuotaMiddleware(checker).CheckQuota(1000)
}

// LightQuotaCheckMiddleware 创建轻量级配额检查中间件（估计消耗 100 个 Token）。
func LightQuotaCheckMiddleware(checker quota.Checker) gin.HandlerFunc {
	return NewQuotaMiddleware(checker).CheckQuota(100)
}

// CheckQuota 检查配额中间件。
func (m *QuotaMiddleware) CheckQuota(estimatedAmount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			writeQuotaError(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
			c.Abort()
			return
		}

		if m.checker == nil {
			writeQuotaError(c, http.StatusServiceUnavailable, "配额检查服务未配置", "配额检查功能暂未启用")
			c.Abort()
			return
		}

		err := m.checker.Check(c.Request.Context(), userID.(string), estimatedAmount)
		if err != nil {
			if errors.Is(err, quota.ErrQuotaExhausted) || errors.Is(err, quota.ErrQuotaSuspended) {
				m.handleQuotaError(c, err)
				return
			}
			if errors.Is(err, quota.ErrInsufficientQuota) {
				m.handleQuotaError(c, err)
				return
			}
			m.handleQuotaError(c, errors.New("配额检查失败"))
			return
		}

		c.Next()
	}
}

func (m *QuotaMiddleware) handleQuotaError(c *gin.Context, err error) {
	if errors.Is(err, quota.ErrQuotaExhausted) {
		writeQuotaError(c, http.StatusTooManyRequests, "配额已用尽", "您的AI配额已用尽，请明天再试或升级会员")
	} else if errors.Is(err, quota.ErrQuotaSuspended) {
		writeQuotaError(c, http.StatusForbidden, "配额已暂停", "您的AI配额已被暂停")
	} else if errors.Is(err, quota.ErrInsufficientQuota) {
		writeQuotaError(c, http.StatusTooManyRequests, "配额不足", "您的AI配额不足以完成此操作")
	} else {
		writeQuotaError(c, http.StatusInternalServerError, "配额检查失败", err.Error())
	}
	c.Abort()
}

type quotaErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

func writeQuotaError(c *gin.Context, statusCode int, message, errorDetail string) {
	response.JSON(c, statusCode, quotaErrorResponse{
		Code:      statusCode,
		Message:   message,
		Error:     errorDetail,
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("request_id"),
	})
}
