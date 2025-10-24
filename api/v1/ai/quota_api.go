package ai

import (
	"net/http"
	"strconv"

	"Qingyu_backend/api/v1/shared"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
)

// QuotaApi 配额API
type QuotaApi struct {
	quotaService *aiService.QuotaService
}

// NewQuotaApi 创建配额API实例
func NewQuotaApi(quotaService *aiService.QuotaService) *QuotaApi {
	return &QuotaApi{
		quotaService: quotaService,
	}
}

// GetQuotaInfo 获取配额信息
// @Summary 获取配额信息
// @Description 获取当前用户的配额信息
// @Tags AI配额
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=ai.UserQuota}
// @Router /api/v1/ai/quota [get]
func (api *QuotaApi) GetQuotaInfo(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	quota, err := api.quotaService.GetQuotaInfo(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取配额信息失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", quota)
}

// GetAllQuotas 获取所有配额
// @Summary 获取所有配额
// @Description 获取当前用户的所有配额类型
// @Tags AI配额
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]ai.UserQuota}
// @Router /api/v1/ai/quota/all [get]
func (api *QuotaApi) GetAllQuotas(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	quotas, err := api.quotaService.GetAllQuotas(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取配额列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", quotas)
}

// GetQuotaStatistics 获取配额统计
// @Summary 获取配额统计
// @Description 获取当前用户的配额使用统计
// @Tags AI配额
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=interfaces.QuotaStatistics}
// @Router /api/v1/ai/quota/statistics [get]
func (api *QuotaApi) GetQuotaStatistics(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	stats, err := api.quotaService.GetQuotaStatistics(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取统计信息失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", stats)
}

// GetTransactionHistory 获取配额事务历史
// @Summary 获取配额事务历史
// @Description 获取当前用户的配额消费记录
// @Tags AI配额
// @Accept json
// @Produce json
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.Response{data=[]ai.QuotaTransaction}
// @Router /api/v1/ai/quota/transactions [get]
func (api *QuotaApi) GetTransactionHistory(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	// 获取分页参数
	limit := 20
	offset := 0
	if l, ok := c.GetQuery("limit"); ok {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	if o, ok := c.GetQuery("offset"); ok {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	transactions, err := api.quotaService.GetTransactionHistory(c.Request.Context(), userID.(string), limit, offset)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取事务历史失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", transactions)
}

// 注意：管理员配额管理功能已迁移到 admin 模块
// 参见: api/v1/admin/quota_admin_api.go
