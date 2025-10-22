package ai

import (
	"net/http"
	"strconv"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/ai"
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

// UpdateQuotaRequest 更新配额请求
type UpdateQuotaRequest struct {
	TotalQuota int    `json:"totalQuota" binding:"required,min=0"`
	QuotaType  string `json:"quotaType" binding:"required,oneof=daily monthly total"`
}

// UpdateUserQuota 更新用户配额（管理员）
// @Summary 更新用户配额
// @Description 管理员更新指定用户的配额
// @Tags AI配额
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Param request body UpdateQuotaRequest true "配额信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/quota/:userId [put]
func (api *QuotaApi) UpdateUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	var req UpdateQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 转换配额类型
	var quotaType ai.QuotaType
	switch req.QuotaType {
	case "daily":
		quotaType = ai.QuotaTypeDaily
	case "monthly":
		quotaType = ai.QuotaTypeMonthly
	case "total":
		quotaType = ai.QuotaTypeTotal
	default:
		shared.Error(c, http.StatusBadRequest, "参数错误", "无效的配额类型")
		return
	}

	err := api.quotaService.UpdateUserQuota(c.Request.Context(), targetUserID, quotaType, req.TotalQuota)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新配额失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// SuspendUserQuota 暂停用户配额（管理员）
// @Summary 暂停用户配额
// @Description 管理员暂停指定用户的配额
// @Tags AI配额
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/quota/:userId/suspend [post]
func (api *QuotaApi) SuspendUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	err := api.quotaService.SuspendUserQuota(c.Request.Context(), targetUserID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "暂停配额失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "暂停成功", nil)
}

// ActivateUserQuota 激活用户配额（管理员）
// @Summary 激活用户配额
// @Description 管理员激活指定用户的配额
// @Tags AI配额
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/quota/:userId/activate [post]
func (api *QuotaApi) ActivateUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "用户ID不能为空")
		return
	}

	err := api.quotaService.ActivateUserQuota(c.Request.Context(), targetUserID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "激活配额失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "激活成功", nil)
}
