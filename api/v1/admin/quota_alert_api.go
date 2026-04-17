package admin

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	ai "Qingyu_backend/service/ai"
)

// QuotaAlertAPI 配额告警API处理器
type QuotaAlertAPI struct {
	alertService *ai.QuotaAlertService
}

// NewQuotaAlertAPI 创建配额告警API实例
func NewQuotaAlertAPI(alertService *ai.QuotaAlertService) *QuotaAlertAPI {
	return &QuotaAlertAPI{
		alertService: alertService,
	}
}

// ListAlerts 获取告警列表
//
//	@Summary		获取配额告警列表
//	@Description	管理员获取配额告警列表（支持筛选和分页）
//	@Tags			管理员-配额告警
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			type	query		string	false	"告警类型(threshold/anomaly/abuse/consistency)"
//	@Param			level	query		string	false	"告警级别(info/warning/critical)"
//	@Param			status	query		string	false	"告警状态(pending/acknowledged/resolved/ignored)"
//	@Param			page	query		int		false	"页码"
//	@Param			limit	query		int		false	"每页数量"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/alerts [get]
func (api *QuotaAlertAPI) ListAlerts(c *gin.Context) {
	alertType := c.Query("type")
	level := c.Query("level")
	status := c.Query("status")
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

	alerts, total, err := api.alertService.ListAlerts(c.Request.Context(), alertType, level, status, page, limit)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取告警列表失败: %w", err))
		return
	}

	response.PaginatedJSON(c, "获取成功", alerts, total, page, limit)
}

// GetAlert 获取告警详情
//
//	@Summary		获取配额告警详情
//	@Description	管理员获取指定配额告警详情
//	@Tags			管理员-配额告警
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"告警ID"
//	@Success		200	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/alerts/{id} [get]
func (api *QuotaAlertAPI) GetAlert(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "参数错误", "告警ID不能为空")
		return
	}

	alert, err := api.alertService.GetAlert(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取告警失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", alert)
}

// AcknowledgeAlert 确认告警
//
//	@Summary		确认告警
//	@Description	管理员确认指定配额告警
//	@Tags			管理员-配额告警
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string				true	"告警ID"
//	@Param			request	body		AlertActionRequest	true	"操作信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/alerts/{id}/acknowledge [put]
func (api *QuotaAlertAPI) AcknowledgeAlert(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "参数错误", "告警ID不能为空")
		return
	}

	var req AlertActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.alertService.AcknowledgeAlert(c.Request.Context(), id, req.OperatorID); err != nil {
		response.InternalError(c, fmt.Errorf("确认告警失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "确认成功", nil)
}

// ResolveAlert 解决告警
//
//	@Summary		解决告警
//	@Description	管理员解决指定配额告警
//	@Tags			管理员-配额告警
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string				true	"告警ID"
//	@Param			request	body		AlertActionRequest	true	"操作信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/alerts/{id}/resolve [put]
func (api *QuotaAlertAPI) ResolveAlert(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "参数错误", "告警ID不能为空")
		return
	}

	var req AlertActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.alertService.ResolveAlert(c.Request.Context(), id, req.OperatorID); err != nil {
		response.InternalError(c, fmt.Errorf("解决告警失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "解决成功", nil)
}

// IgnoreAlert 忽略告警
//
//	@Summary		忽略告警
//	@Description	管理员忽略指定配额告警
//	@Tags			管理员-配额告警
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string				true	"告警ID"
//	@Param			request	body		AlertActionRequest	true	"操作信息"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/alerts/{id}/ignore [put]
func (api *QuotaAlertAPI) IgnoreAlert(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "参数错误", "告警ID不能为空")
		return
	}

	var req AlertActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := api.alertService.IgnoreAlert(c.Request.Context(), id, req.OperatorID); err != nil {
		response.InternalError(c, fmt.Errorf("忽略告警失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "已忽略", nil)
}
