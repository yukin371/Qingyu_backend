package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
// @Success 200 {object} authService.TemplateResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates [post]
func (api *PermissionTemplateAPI) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "请求参数无效",
			Details: err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "CREATE_TEMPLATE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetTemplate 获取权限模板
// @Summary 获取权限模板
// @Description 根据ID获取权限模板详情
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 200 {object} authService.TemplateResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates/{id} [get]
func (api *PermissionTemplateAPI) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_ID",
			Message: "模板ID不能为空",
		})
		return
	}

	resp, err := api.templateService.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "TEMPLATE_NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetTemplateByCode 根据代码获取权限模板
// @Summary 根据代码获取权限模板
// @Description 根据模板代码获取权限模板详情
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param code path string true "模板代码"
// @Success 200 {object} authService.TemplateResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates/code/{code} [get]
func (api *PermissionTemplateAPI) GetTemplateByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_CODE",
			Message: "模板代码不能为空",
		})
		return
	}

	resp, err := api.templateService.GetTemplateByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    "TEMPLATE_NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// UpdateTemplate 更新权限模板
// @Summary 更新权限模板
// @Description 更新权限模板信息
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Param request body UpdateTemplateRequest true "更新请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates/{id} [put]
func (api *PermissionTemplateAPI) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_ID",
			Message: "模板ID不能为空",
		})
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "请求参数无效",
			Details: err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "UPDATE_TEMPLATE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// DeleteTemplate 删除权限模板
// @Summary 删除权限模板
// @Description 删除指定的权限模板
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates/{id} [delete]
func (api *PermissionTemplateAPI) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_ID",
			Message: "模板ID不能为空",
		})
		return
	}

	err := api.templateService.DeleteTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "DELETE_TEMPLATE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// ListTemplates 列出权限模板
// @Summary 列出权限模板
// @Description 获取权限模板列表，可按分类筛选
// @Tags 权限模板
// @Accept json
// @Produce json
// @Param category query string false "分类(reader/author/admin)"
// @Success 200 {object} ListResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates [get]
func (api *PermissionTemplateAPI) ListTemplates(c *gin.Context) {
	category := c.Query("category")

	templates, err := api.templateService.ListTemplates(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "LIST_TEMPLATES_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Total: len(templates),
		Items: templates,
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
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates/{id}/apply [post]
func (api *PermissionTemplateAPI) ApplyTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_ID",
			Message: "模板ID不能为空",
		})
		return
	}

	var req struct {
		RoleID string `json:"roleId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "请求参数无效",
			Details: err.Error(),
		})
		return
	}

	err := api.templateService.ApplyTemplate(c.Request.Context(), templateID, req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "APPLY_TEMPLATE_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// InitializeSystemTemplates 初始化系统模板
// @Summary 初始化系统模板
// @Description 初始化系统预设的权限模板
// @Tags 权限模板
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /admin/permission-templates/initialize [post]
func (api *PermissionTemplateAPI) InitializeSystemTemplates(c *gin.Context) {
	err := api.templateService.InitializeSystemTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "INITIALIZE_TEMPLATES_FAILED",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
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

// ListResponse 列表响应
type ListResponse struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
