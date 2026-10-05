package ai

import (
	"strconv"

	"Qingyu_backend/api/v1/shared"
	aiService "Qingyu_backend/service/ai"

	"Qingyu_backend/pkg/response"

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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/quota [get]
func (api *QuotaApi) GetQuotaInfo(c *gin.Context) {
	userID, ok := shared.GetUserIDWithMessage(c, "无法获取用户信息")
	if !ok {
		return
	}

	quota, err := api.quotaService.GetQuotaInfo(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, quota)
}

// GetAllQuotas 获取所有配额
// @Summary 获取所有配额
// @Description 获取当前用户的所有配额类型
// @Tags AI配额
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/quota/all [get]
func (api *QuotaApi) GetAllQuotas(c *gin.Context) {
	userID, ok := shared.GetUserIDWithMessage(c, "无法获取用户信息")
	if !ok {
		return
	}

	quotas, err := api.quotaService.GetAllQuotas(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, quotas)
}

// GetQuotaStatistics 获取配额统计
// @Summary 获取配额统计
// @Description 获取当前用户的配额使用统计
// @Tags AI配额
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/quota/statistics [get]
func (api *QuotaApi) GetQuotaStatistics(c *gin.Context) {
	userID, ok := shared.GetUserIDWithMessage(c, "无法获取用户信息")
	if !ok {
		return
	}

	stats, err := api.quotaService.GetQuotaStatistics(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, stats)
}

// GetTransactionHistory 获取配额事务历史
// @Summary 获取配额事务历史
// @Description 获取当前用户的配额消费记录
// @Tags AI配额
// @Accept json
// @Produce json
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/quota/transactions [get]
func (api *QuotaApi) GetTransactionHistory(c *gin.Context) {
	userID, ok := shared.GetUserIDWithMessage(c, "无法获取用户信息")
	if !ok {
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

	transactions, err := api.quotaService.GetTransactionHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, transactions)
}

// RechargeRequest 充值请求
type RechargeRequest struct {
	Amount int    `json:"amount" binding:"required,gt=0"` // 充值数量（必须 > 0）
	Reason string `json:"reason"`                         // 充值原因
}

// RechargeQuota 用户自助充值配额（需要购买或兑换）
// @Summary 配额充值
// @Description 用户使用积分或购买充值配额
// @Tags AI配额
// @Accept json
// @Produce json
// @Param request body RechargeRequest true "充值请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/ai/quota/recharge [post]
func (api *QuotaApi) RechargeQuota(c *gin.Context) {
	userID, ok := shared.GetUserIDWithMessage(c, "无法获取用户信息")
	if !ok {
		return
	}

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 调用充值服务
	err := api.quotaService.RechargeQuota(c.Request.Context(), userID, req.Amount, req.Reason, userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// 注意：管理员配额管理功能已迁移到 admin 模块
// 参见: api/v1/admin/quota_admin_api.go
