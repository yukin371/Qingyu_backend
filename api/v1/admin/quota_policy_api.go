package admin

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	aiModels "Qingyu_backend/models/ai"
	"Qingyu_backend/pkg/response"
	ai "Qingyu_backend/service/ai"
)

// QuotaPolicyAPI 配额策略API处理器
type QuotaPolicyAPI struct {
	policyService *ai.QuotaPolicyService
}

// NewQuotaPolicyAPI 创建配额策略API实例
func NewQuotaPolicyAPI(policyService *ai.QuotaPolicyService) *QuotaPolicyAPI {
	return &QuotaPolicyAPI{
		policyService: policyService,
	}
}

// CreatePolicy 创建配额策略
//
//	@Summary		创建配额策略
//	@Description	管理员创建新的配额策略
//	@Tags			管理员-配额策略
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		PolicyRequest	true	"策略信息"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/policies [post]
func (api *QuotaPolicyAPI) CreatePolicy(c *gin.Context) {
	var req PolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	policy := &aiModels.QuotaPolicy{
		Name:            req.Name,
		UserRole:        aiModels.UserRole(req.UserRole),
		MembershipLevel: aiModels.MembershipLevel(req.MembershipLevel),
		DailyQuota:      req.DailyQuota,
		MonthlyQuota:    req.MonthlyQuota,
		TotalQuota:      req.TotalQuota,
		Description:     req.Description,
	}

	if err := api.policyService.CreatePolicy(c.Request.Context(), policy); err != nil {
		response.InternalError(c, fmt.Errorf("创建策略失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "创建成功", policy)
}

// GetPolicy 获取策略详情
//
//	@Summary		获取配额策略详情
//	@Description	管理员获取指定配额策略详情
//	@Tags			管理员-配额策略
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"策略ID"
//	@Success		200	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/policies/{id} [get]
func (api *QuotaPolicyAPI) GetPolicy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "参数错误", "策略ID不能为空")
		return
	}

	policy, err := api.policyService.GetPolicy(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取策略失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "获取成功", policy)
}

// ListPolicies 获取策略列表
//
//	@Summary		获取配额策略列表
//	@Description	管理员获取配额策略列表（支持筛选和分页）
//	@Tags			管理员-配额策略
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			role	query		string	false	"角色筛选"
//	@Param			status	query		string	false	"状态筛选"
//	@Param			page	query		int		false	"页码"
//	@Param			limit	query		int		false	"每页数量"
//	@Success		200		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/policies [get]
func (api *QuotaPolicyAPI) ListPolicies(c *gin.Context) {
	role := c.Query("role")
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

	policies, total, err := api.policyService.ListPolicies(c.Request.Context(), role, status, page, limit)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取策略列表失败: %w", err))
		return
	}

	response.PaginatedJSON(c, "获取成功", policies, total, page, limit)
}

// UpdatePolicy 更新配额策略
//
//	@Summary		更新配额策略
//	@Description	管理员更新指定配额策略
//	@Tags			管理员-配额策略
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string			true	"策略ID"
//	@Param			request	body		PolicyRequest	true	"策略信息"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/policies/{id} [put]
func (api *QuotaPolicyAPI) UpdatePolicy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "参数错误", "策略ID不能为空")
		return
	}

	var req PolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	policy, err := api.policyService.GetPolicy(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取策略失败: %w", err))
		return
	}

	policy.Name = req.Name
	policy.UserRole = aiModels.UserRole(req.UserRole)
	policy.MembershipLevel = aiModels.MembershipLevel(req.MembershipLevel)
	policy.DailyQuota = req.DailyQuota
	policy.MonthlyQuota = req.MonthlyQuota
	policy.TotalQuota = req.TotalQuota
	policy.Description = req.Description

	if err := api.policyService.UpdatePolicy(c.Request.Context(), policy); err != nil {
		response.InternalError(c, fmt.Errorf("更新策略失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "更新成功", policy)
}

// DeletePolicy 删除配额策略（软删除）
//
//	@Summary		删除配额策略
//	@Description	管理员删除指定配额策略（软删除）
//	@Tags			管理员-配额策略
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"策略ID"
//	@Success		200	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/policies/{id} [delete]
func (api *QuotaPolicyAPI) DeletePolicy(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "参数错误", "策略ID不能为空")
		return
	}

	if err := api.policyService.DeletePolicy(c.Request.Context(), id); err != nil {
		response.InternalError(c, fmt.Errorf("删除策略失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// InitializeDefaultPolicies 初始化默认策略
//
//	@Summary		初始化默认策略
//	@Description	管理员初始化系统默认配额策略
//	@Tags			管理员-配额策略
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.APIResponse
//	@Router			/api/v1/admin/quota/policies/initialize [post]
func (api *QuotaPolicyAPI) InitializeDefaultPolicies(c *gin.Context) {
	if err := api.policyService.InitializeDefaultPolicies(c.Request.Context()); err != nil {
		response.InternalError(c, fmt.Errorf("初始化默认策略失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "初始化成功", nil)
}
