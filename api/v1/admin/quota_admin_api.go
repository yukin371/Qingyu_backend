package admin

import (

	"github.com/gin-gonic/gin"

	aiModels "Qingyu_backend/models/ai"
	ai "Qingyu_backend/service/ai"
	"Qingyu_backend/pkg/response"
)

// QuotaAdminAPI AI配额管理API处理器（管理员）
type QuotaAdminAPI struct {
	quotaService *ai.QuotaService
}

// NewQuotaAdminAPI 创建AI配额管理API实例
func NewQuotaAdminAPI(quotaService *ai.QuotaService) *QuotaAdminAPI {
	return &QuotaAdminAPI{
		quotaService: quotaService,
	}
}

// UpdateUserQuota 更新用户配额（管理员）
//
//	@Summary		更新用户配额
//	@Description	管理员更新指定用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string				true	"用户ID"
//	@Param			request	body		UpdateQuotaRequest	true	"配额信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/quota/{userId} [put]
func (api *QuotaAdminAPI) UpdateUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c,  "参数错误", "用户ID不能为空")
		return
	}

	var req UpdateQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	// 转换配额类型
	var quotaType aiModels.QuotaType
	switch req.QuotaType {
	case "daily":
		quotaType = aiModels.QuotaTypeDaily
	case "monthly":
		quotaType = aiModels.QuotaTypeMonthly
	case "total":
		quotaType = aiModels.QuotaTypeTotal
	default:
		response.BadRequest(c,  "参数错误", "无效的配额类型")
		return
	}

	err := api.quotaService.UpdateUserQuota(c.Request.Context(), targetUserID, quotaType, req.TotalQuota)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// SuspendUserQuota 暂停用户配额（管理员）
//
//	@Summary		暂停用户配额
//	@Description	管理员暂停指定用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/quota/{userId}/suspend [post]
func (api *QuotaAdminAPI) SuspendUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c,  "参数错误", "用户ID不能为空")
		return
	}

	err := api.quotaService.SuspendUserQuota(c.Request.Context(), targetUserID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "暂停成功", nil)
}

// ActivateUserQuota 激活用户配额（管理员）
//
//	@Summary		激活用户配额
//	@Description	管理员激活指定用户的AI配额
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/quota/{userId}/activate [post]
func (api *QuotaAdminAPI) ActivateUserQuota(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c,  "参数错误", "用户ID不能为空")
		return
	}

	err := api.quotaService.ActivateUserQuota(c.Request.Context(), targetUserID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "激活成功", nil)
}

// GetUserQuotaDetails 获取用户配额详情（管理员）
//
//	@Summary		获取用户配额详情
//	@Description	管理员获取指定用户的配额详情
//	@Tags			管理员-AI配额管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			userId	path		string	true	"用户ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/quota/{userId} [get]
func (api *QuotaAdminAPI) GetUserQuotaDetails(c *gin.Context) {
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		response.BadRequest(c,  "参数错误", "用户ID不能为空")
		return
	}

	// 获取所有类型的配额
	quotas, err := api.quotaService.GetAllQuotas(c.Request.Context(), targetUserID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", quotas)
}
