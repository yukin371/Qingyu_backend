package admin

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	aiModels "Qingyu_backend/models/ai"
	"Qingyu_backend/pkg/response"
	ai "Qingyu_backend/service/ai"
)

// QuotaAdminAPI AI配额管理API处理器（管理员）
type QuotaAdminAPI struct {
	quotaService *ai.QuotaService
	adminService *ai.QuotaAdminService
}

// NewQuotaAdminAPI 创建AI配额管理API实例
func NewQuotaAdminAPI(quotaService *ai.QuotaService, adminService *ai.QuotaAdminService) *QuotaAdminAPI {
	return &QuotaAdminAPI{
		quotaService: quotaService,
		adminService: adminService,
	}
}

// ===========================
// 用户配额管理
// ===========================

// ListUserQuotas 获取用户配额列表
//
//	@Summary		获取用户配额列表
//	@Description	管理员获取用户配额列表（支持筛选和分页）
//	@Tags			管理员-AI配额管理
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			role	query		string	false	"角色筛选(reader/writer/admin)"
//	@Param			status	query		string	false	"状态筛选(active/exhausted/suspended)"
//	@Param			search	query		string	false	"搜索用户ID"
//	@Param			page	query		int		false	"页码"
//	@Param			limit	query		int		false	"每页数量"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/users [get]
func (api *QuotaAdminAPI) ListUserQuotas(c *gin.Context) {
	role := c.Query("role")
	status := c.Query("status")
	search := c.Query("search")
	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	items, total, err := api.adminService.ListUserQuotas(c.Request.Context(), role, status, search, page, limit)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取用户配额列表失败: %w", err))
		return
	}

	response.PaginatedJSON(c, "获取成功", items, total, page, limit)
}

// GetUserQuotaDetails 获取用户配额详情
//
//	@Summary		获取用户配额详情
//	@Description	管理员获取指定用户的配额详情
//	@Tags			管理员-AI配额管理
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/users/{userId} [get]
func (api *QuotaAdminAPI) GetUserQuotaDetails(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	quotas, err := api.quotaService.GetAllQuotas(c.Request.Context(), targetUserID)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取配额详情失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", quotas)
}

// GetUserQuotaReconciliation 获取单用户配额对账结果
//
//	@Summary		获取单用户配额对账结果
//	@Description	管理员查看指定用户在 backend 与 AI service 之间的配额消费差异
//	@Tags			管理员-AI配额管理
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId			path		string	true	"用户ID"
//	@Param			timeRange		query		string	false	"时间范围(day/week/month/all，默认day)"
//	@Param			workflowType	query		string	false	"工作流类型过滤"
//	@Success		200				{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/users/{userId}/reconciliation [get]
func (api *QuotaAdminAPI) GetUserQuotaReconciliation(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	result, err := api.adminService.GetUserQuotaReconciliation(
		c.Request.Context(),
		targetUserID,
		c.DefaultQuery("timeRange", "day"),
		c.Query("workflowType"),
	)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取配额对账结果失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", result)
}

// UpdateUserQuota 更新用户配额
//
//	@Summary		更新用户配额
//	@Description	管理员更新指定用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string				true	"用户ID"
//	@Param			request	body		UpdateQuotaRequest	true	"配额信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/users/{userId} [put]
func (api *QuotaAdminAPI) UpdateUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	var req UpdateQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	var quotaType aiModels.QuotaType
	switch req.QuotaType {
	case "daily":
		quotaType = aiModels.QuotaTypeDaily
	case "monthly":
		quotaType = aiModels.QuotaTypeMonthly
	case "total":
		quotaType = aiModels.QuotaTypeTotal
	default:
		response.BadRequest(c, "参数错误", "无效的配额类型")
		return
	}

	err := api.quotaService.UpdateUserQuota(c.Request.Context(), targetUserID, quotaType, req.TotalQuota)
	if err != nil {
		response.InternalError(c, fmt.Errorf("更新配额失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// RechargeUserQuota 管理员充值用户配额
//
//	@Summary		充值用户配额
//	@Description	管理员为指定用户充值AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string					true	"用户ID"
//	@Param			request	body		AdminRechargeRequest	true	"充值信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/users/{userId}/recharge [post]
func (api *QuotaAdminAPI) RechargeUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	var req AdminRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	operatorID := c.GetString("user_id")
	err := api.adminService.RechargeUserQuota(c.Request.Context(), targetUserID, req.Amount, req.QuotaType, req.Reason, operatorID)
	if err != nil {
		response.InternalError(c, fmt.Errorf("充值失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "充值成功", nil)
}

// SuspendUserQuota 暂停用户配额
//
//	@Summary		暂停用户配额
//	@Description	管理员暂停指定用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/users/{userId}/suspend [post]
func (api *QuotaAdminAPI) SuspendUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	err := api.quotaService.SuspendUserQuota(c.Request.Context(), targetUserID)
	if err != nil {
		response.InternalError(c, fmt.Errorf("暂停失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "暂停成功", nil)
}

// ActivateUserQuota 激活用户配额
//
//	@Summary		激活用户配额
//	@Description	管理员激活指定用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/users/{userId}/activate [post]
func (api *QuotaAdminAPI) ActivateUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	err := api.quotaService.ActivateUserQuota(c.Request.Context(), targetUserID)
	if err != nil {
		response.InternalError(c, fmt.Errorf("激活失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "激活成功", nil)
}

// ===========================
// 批量操作
// ===========================

// BatchRecharge 批量充值
//
//	@Summary		批量充值配额
//	@Description	管理员为多个用户批量充值AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		BatchQuotaRequest	true	"批量充值信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/batch-recharge [post]
func (api *QuotaAdminAPI) BatchRecharge(c *gin.Context) {
	var req BatchQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	operatorID := c.GetString("user_id")
	result, err := api.adminService.BatchRecharge(c.Request.Context(), req.UserIDs, req.Amount, req.QuotaType, req.Reason, operatorID)
	if err != nil {
		response.InternalError(c, fmt.Errorf("批量充值失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "批量充值完成", result)
}

// BatchUpdateQuota 批量更新配额
//
//	@Summary		批量更新配额
//	@Description	管理员为多个用户批量更新AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		BatchQuotaRequest	true	"批量更新信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/batch-update [post]
func (api *QuotaAdminAPI) BatchUpdateQuota(c *gin.Context) {
	var req BatchQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	result, err := api.adminService.BatchUpdateQuota(c.Request.Context(), req.UserIDs, req.Amount, req.QuotaType)
	if err != nil {
		response.InternalError(c, fmt.Errorf("批量更新失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "批量更新完成", result)
}

// BatchSuspend 批量暂停
//
//	@Summary		批量暂停配额
//	@Description	管理员批量暂停多个用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		BatchQuotaRequest	true	"批量暂停信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/batch-suspend [post]
func (api *QuotaAdminAPI) BatchSuspend(c *gin.Context) {
	var req BatchQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	result, err := api.adminService.BatchSuspend(c.Request.Context(), req.UserIDs)
	if err != nil {
		response.InternalError(c, fmt.Errorf("批量暂停失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "批量暂停完成", result)
}

// BatchActivate 批量激活
//
//	@Summary		批量激活配额
//	@Description	管理员批量激活多个用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		BatchQuotaRequest	true	"批量激活信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/batch-activate [post]
func (api *QuotaAdminAPI) BatchActivate(c *gin.Context) {
	var req BatchQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	result, err := api.adminService.BatchActivate(c.Request.Context(), req.UserIDs)
	if err != nil {
		response.InternalError(c, fmt.Errorf("批量激活失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "批量激活完成", result)
}
