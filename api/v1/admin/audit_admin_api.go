package admin

import (
	"Qingyu_backend/service/interfaces/audit"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
)

// AuditAdminAPI 审核管理API处理器（管理员）
type AuditAdminAPI struct {
	auditService audit.ContentAuditService
}

// NewAuditAdminAPI 创建审核管理API实例
func NewAuditAdminAPI(auditService audit.ContentAuditService) *AuditAdminAPI {
	return &AuditAdminAPI{
		auditService: auditService,
	}
}

// GetPendingAudits 获取待审核内容列表（管理员）
//
//	@Summary		获取待审核内容
//	@Description	获取需要人工审核的内容列表
//	@Tags			管理员-审核管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			targetType	query		string	false	"目标类型(document/chapter/comment)"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			limit		query		int		false	"每页数量"	default(50)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.ErrorResponse
//	@Failure		403			{object}	shared.ErrorResponse
//	@Failure		500			{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/audit/pending [get]
func (api *AuditAdminAPI) GetPendingAudits(c *gin.Context) {
	targetType := c.Query("targetType")
	page := 1
	limit := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// 调用Service层获取待审核列表
	records, err := api.auditService.GetPendingReviews(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	// 如果指定了targetType，进行过滤
	// TODO: 这里应该在Service层实现过滤逻辑
	_ = targetType
	_ = page

	response.Success(c, records)
}

// ReviewAudit 审核通过/拒绝（管理员）
//
//	@Summary		审核内容
//	@Description	管理员审核内容，通过或拒绝
//	@Tags			管理员-审核管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string				true	"审核记录ID"
//	@Param			request	body		ReviewAuditRequest	true	"审核信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/audit/{id}/review [post]
func (api *AuditAdminAPI) ReviewAudit(c *gin.Context) {
	auditID := c.Param("id")
	if auditID == "" {
		response.BadRequest(c, "参数错误", "审核记录ID不能为空")
		return
	}

	var req ReviewAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从context获取管理员ID
	reviewerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 判断审核动作
	approved := req.Action == "approve"

	// 调用Service层
	err := api.auditService.ReviewAudit(c.Request.Context(), auditID, reviewerID.(string), approved, req.ReviewNote)
	if err != nil {
		c.Error(err)
		return
	}

	message := "审核已拒绝"
	if approved {
		message = "审核已通过"
	}

	response.SuccessWithMessage(c, message, nil)
}

// ReviewAppeal 审核申诉（管理员）
//
//	@Summary		审核申诉
//	@Description	管理员审核用户的申诉请求
//	@Tags			管理员-审核管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string				true	"审核记录ID"
//	@Param			request	body		ReviewAppealRequest	true	"审核申诉信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.ErrorResponse
//	@Failure		401		{object}	shared.ErrorResponse
//	@Failure		403		{object}	shared.ErrorResponse
//	@Failure		500		{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/audit/{id}/appeal/review [post]
func (api *AuditAdminAPI) ReviewAppeal(c *gin.Context) {
	auditID := c.Param("id")
	if auditID == "" {
		response.BadRequest(c, "参数错误", "审核记录ID不能为空")
		return
	}

	var req ReviewAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从context获取管理员ID
	reviewerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 判断审核动作
	approved := req.Action == "approve"

	// 调用Service层
	err := api.auditService.ReviewAppeal(c.Request.Context(), auditID, reviewerID.(string), approved, req.ReviewNote)
	if err != nil {
		c.Error(err)
		return
	}

	message := "申诉已驳回"
	if approved {
		message = "申诉已通过"
	}

	response.SuccessWithMessage(c, message, nil)
}

// GetHighRiskAudits 获取高风险审核记录（管理员）
//
//	@Summary		获取高风险审核记录
//	@Description	获取高风险审核记录列表
//	@Tags			管理员-审核管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			minRiskLevel	query		int	false	"最低风险等级"	default(3)
//	@Param			limit			query		int	false	"数量限制"		default(50)
//	@Success		200				{object}	shared.APIResponse
//	@Failure		401				{object}	shared.ErrorResponse
//	@Failure		403				{object}	shared.ErrorResponse
//	@Failure		500				{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/audit/high-risk [get]
func (api *AuditAdminAPI) GetHighRiskAudits(c *gin.Context) {
	minRiskLevel := 3
	limit := 50

	if level := c.Query("minRiskLevel"); level != "" {
		if l, err := strconv.Atoi(level); err == nil {
			minRiskLevel = l
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// 调用Service层
	records, err := api.auditService.GetHighRiskAudits(c.Request.Context(), minRiskLevel, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, records)
}

// GetAuditStatistics 获取审核统计（管理员）
//
//	@Summary		获取审核统计
//	@Description	获取审核相关的统计数据
//	@Tags			管理员-审核管理
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	shared.APIResponse
//	@Failure		401	{object}	shared.ErrorResponse
//	@Failure		403	{object}	shared.ErrorResponse
//	@Failure		500	{object}	shared.ErrorResponse
//	@Router			/api/v1/admin/audit/statistics [get]
func (api *AuditAdminAPI) GetAuditStatistics(c *gin.Context) {
	// 调用AuditService获取审核统计
	stats, err := api.auditService.GetAuditStatistics(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, stats)
}
