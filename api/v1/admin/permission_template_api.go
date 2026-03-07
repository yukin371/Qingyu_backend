package admin

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	authService "Qingyu_backend/service/auth"
)

// PermissionTemplateAPI 权限模板API
type PermissionTemplateAPI struct {
	templateService authService.PermissionTemplateService
}

// NewPermissionTemplateAPI 创建权限模板API
func NewPermissionTemplateAPI(templateService authService.PermissionTemplateService) *PermissionTemplateAPI {
	return &PermissionTemplateAPI{
		templateService: templateService,
	}
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Code        string   `json:"code" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" binding:"required"`
	Category    string   `json:"category"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	Category    string   `json:"category"`
}

// ApplyTemplateRequest 应用模板请求
type ApplyTemplateRequest struct {
	RoleID string `json:"roleId" binding:"required"`
}

// CreateTemplate 创建权限模板
// @Summary 创建权限模板
// @Description 创建新的权限模板
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param request body CreateTemplateRequest true "创建请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates [post]
func (api *PermissionTemplateAPI) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效", err.Error())
		return
	}

	// 转换为服务请求
	serviceReq := &authService.CreateTemplateRequest{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Permissions: req.Permissions,
		Category:    req.Category,
	}

	// 调用服务
	resp, err := api.templateService.CreateTemplate(c.Request.Context(), serviceReq)
	if err != nil {
		response.BadRequest(c, "创建权限模板失败", err.Error())
		return
	}

	response.Success(c, resp)
}

// GetTemplate 获取权限模板
// @Summary 获取权限模板
// @Description 根据ID获取权限模板详情
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates/{id} [get]
func (api *PermissionTemplateAPI) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		response.BadRequest(c, "模板ID不能为空", nil)
		return
	}

	resp, err := api.templateService.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetTemplateByCode 根据代码获取权限模板
// @Summary 根据代码获取权限模板
// @Description 根据模板代码获取权限模板详情
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param code path string true "模板代码"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates/code/{code} [get]
func (api *PermissionTemplateAPI) GetTemplateByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "模板代码不能为空", nil)
		return
	}

	resp, err := api.templateService.GetTemplateByCode(c.Request.Context(), code)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// UpdateTemplate 更新权限模板
// @Summary 更新权限模板
// @Description 更新权限模板信息
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Param request body UpdateTemplateRequest true "更新请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates/{id} [put]
func (api *PermissionTemplateAPI) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		response.BadRequest(c, "模板ID不能为空", nil)
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效", err.Error())
		return
	}

	// 转换为服务请求
	serviceReq := &authService.UpdateTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		Category:    req.Category,
	}

	err := api.templateService.UpdateTemplate(c.Request.Context(), templateID, serviceReq)
	if err != nil {
		response.BadRequest(c, "更新权限模板失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteTemplate 删除权限模板
// @Summary 删除权限模板
// @Description 删除指定的权限模板
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates/{id} [delete]
func (api *PermissionTemplateAPI) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		response.BadRequest(c, "模板ID不能为空", nil)
		return
	}

	err := api.templateService.DeleteTemplate(c.Request.Context(), templateID)
	if err != nil {
		response.BadRequest(c, "删除权限模板失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// ListTemplates 列出权限模板
// @Summary 列出权限模板
// @Description 获取权限模板列表，可按分类筛选
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param category query string false "分类(reader/author/admin)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates [get]
func (api *PermissionTemplateAPI) ListTemplates(c *gin.Context) {
	category := c.Query("category")

	templates, err := api.templateService.ListTemplates(c.Request.Context(), category)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"items": templates,
		"total": len(templates),
	})
}

// ApplyTemplate 应用模板到角色
// @Summary 应用模板到角色
// @Description 将权限模板应用到指定角色
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Param roleId body ApplyTemplateRequest true "角色ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates/{id}/apply [post]
func (api *PermissionTemplateAPI) ApplyTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		response.BadRequest(c, "模板ID不能为空", nil)
		return
	}

	var req struct {
		RoleID string `json:"roleId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效", err.Error())
		return
	}

	err := api.templateService.ApplyTemplate(c.Request.Context(), templateID, req.RoleID)
	if err != nil {
		response.BadRequest(c, "应用权限模板失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// InitializeSystemTemplates 初始化系统模板
// @Summary 初始化系统模板
// @Description 初始化系统预设的权限模板
// @Tags 权限模板
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /admin/permission-templates/initialize [post]
func (api *PermissionTemplateAPI) InitializeSystemTemplates(c *gin.Context) {
	err := api.templateService.InitializeSystemTemplates(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// RegisterRoutes 注册路由
func (api *PermissionTemplateAPI) RegisterRoutes(router *gin.RouterGroup) {
	templates := router.Group("/permission-templates")
	{
		templates.POST("", api.CreateTemplate)
		templates.GET("", api.ListTemplates)
		// 特定路由必须在参数路由之前注册
		templates.POST("/initialize", api.InitializeSystemTemplates)
		templates.GET("/code/:code", api.GetTemplateByCode)
		templates.GET("/:id", api.GetTemplate)
		templates.PUT("/:id", api.UpdateTemplate)
		templates.DELETE("/:id", api.DeleteTemplate)
		templates.POST("/:id/apply", api.ApplyTemplate)
	}
}

// ============ 辅助类型 ============
