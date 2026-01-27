package writer

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/writer/document"
)

// TemplateAPI 模板API
type TemplateAPI struct {
	service *document.TemplateService
	logger  *zap.Logger
}

// NewTemplateAPI 创建模板API实例
func NewTemplateAPI(service *document.TemplateService, logger *zap.Logger) *TemplateAPI {
	return &TemplateAPI{
		service: service,
		logger:  logger,
	}
}

// CreateTemplate 创建模板
// @Summary 创建模板
// @Description 创建新的文档模板
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param request body CreateTemplateRequest true "创建模板请求"
// @Success 201 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 401 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/templates [post]
func (api *TemplateAPI) CreateTemplate(c *gin.Context) {
	// 检查服务是否初始化
	if api.service == nil {
		shared.Error(c, http.StatusInternalServerError, "服务未初始化", "模板服务未正确初始化")
		return
	}

	// 获取并验证用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	// 将用户ID添加到context
	ctx := context.WithValue(c.Request.Context(), "userID", userIDStr)

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 转换为service请求
	serviceReq := &document.CreateTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Type:        writer.TemplateType(req.Type),
		Category:    req.Category,
		Content:     req.Content,
		Variables:   convertVariables(req.Variables),
		CreatedBy:   userIDStr,
	}

	// 处理ProjectID
	if req.ProjectID != "" {
		projectID, err := primitive.ObjectIDFromHex(req.ProjectID)
		if err != nil {
			shared.Error(c, http.StatusBadRequest, "参数错误", "无效的项目ID")
			return
		}
		serviceReq.ProjectID = &projectID
	}

	// 调用service创建模板
	template, err := api.service.CreateTemplate(ctx, serviceReq)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", template)
}

// ListTemplates 列出模板
// @Summary 列出模板
// @Description 获取模板列表（支持分页和过滤）
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param projectId query string false "项目ID（为空则查询全局模板）"
// @Param type query string false "模板类型 (chapter/outline/setting)"
// @Param status query string false "模板状态 (active/deprecated)"
// @Param category query string false "模板分类"
// @Param keyword query string false "关键词搜索"
// @Param isSystem query bool false "是否为系统模板"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Param sortBy query string false "排序字段" default(created_at)
// @Param sortOrder query int false "排序方向 (1=升序, -1=降序)" default(-1)
// @Success 200 {object} shared.PaginatedResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 401 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/templates [get]
func (api *TemplateAPI) ListTemplates(c *gin.Context) {
	// 解析查询参数
	req := &document.ListTemplatesRequest{}

	// 项目ID（可选）
	if projectID := c.Query("projectId"); projectID != "" {
		id, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			shared.Error(c, http.StatusBadRequest, "参数错误", "无效的项目ID")
			return
		}
		req.ProjectID = &id
	}

	// 模板类型（可选）
	if templateType := c.Query("type"); templateType != "" {
		ttype := writer.TemplateType(templateType)
		req.Type = &ttype
	}

	// 模板状态（可选）
	if status := c.Query("status"); status != "" {
		s := writer.TemplateStatus(status)
		req.Status = &s
	}

	// 其他参数
	req.Category = c.Query("category")
	req.Keyword = c.Query("keyword")

	// 是否系统模板
	if isSystemStr := c.Query("isSystem"); isSystemStr != "" {
		isSystem, err := strconv.ParseBool(isSystemStr)
		if err == nil {
			req.IsSystem = &isSystem
		}
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	req.Page = page
	req.PageSize = pageSize

	// 排序参数
	req.SortBy = c.DefaultQuery("sortBy", "created_at")
	sortOrder, _ := strconv.Atoi(c.DefaultQuery("sortOrder", "-1"))
	req.SortOrder = sortOrder

	// 调用service查询
	templates, total, err := api.service.ListTemplates(c.Request.Context(), req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	shared.Paginated(c, templates, total, page, pageSize, "查询成功")
}

// GetTemplate 获取模板详情
// @Summary 获取模板详情
// @Description 根据ID获取模板详细信息
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/templates/{id} [get]
func (api *TemplateAPI) GetTemplate(c *gin.Context) {
	id := c.Param("id")

	templateID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", "无效的模板ID")
		return
	}

	template, err := api.service.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", template)
}

// UpdateTemplate 更新模板
// @Summary 更新模板
// @Description 更新模板信息
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Param request body UpdateTemplateRequest true "更新模板请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 403 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/templates/{id} [put]
func (api *TemplateAPI) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")

	templateID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", "无效的模板ID")
		return
	}

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 转换为service请求
	serviceReq := &document.UpdateTemplateRequest{}

	if req.Name != nil {
		serviceReq.Name = req.Name
	}
	if req.Description != nil {
		serviceReq.Description = req.Description
	}
	if req.Category != nil {
		serviceReq.Category = req.Category
	}
	if req.Content != nil {
		serviceReq.Content = req.Content
	}
	if req.Variables != nil {
		vars := convertVariables(*req.Variables)
		serviceReq.Variables = &vars
	}

	// 调用service更新
	template, err := api.service.UpdateTemplate(c.Request.Context(), templateID, serviceReq)
	if err != nil {
		// 判断是否为权限问题
		if err.Error() == "系统模板不允许修改" {
			shared.Error(c, http.StatusForbidden, "禁止操作", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", template)
}

// DeleteTemplate 删除模板
// @Summary 删除模板
// @Description 删除模板（软删除）
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 401 {object} shared.APIResponse
// @Failure 403 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/templates/{id} [delete]
func (api *TemplateAPI) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")

	templateID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", "无效的模板ID")
		return
	}

	// 获取并验证用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "无效的用户ID")
		return
	}

	// 调用service删除
	err = api.service.DeleteTemplate(c.Request.Context(), templateID, userIDStr)
	if err != nil {
		// 判断错误类型
		if err.Error() == "系统模板不允许删除" {
			shared.Error(c, http.StatusForbidden, "禁止操作", err.Error())
			return
		}
		if err.Error() == "无权限删除此模板" {
			shared.Error(c, http.StatusForbidden, "权限不足", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "删除失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// ApplyTemplate 应用模板
// @Summary 应用模板
// @Description 渲染模板并应用到文档
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param id path string true "模板ID"
// @Param request body ApplyTemplateRequest true "应用模板请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Failure 500 {object} shared.APIResponse
// @Router /api/v1/writer/templates/{id}/apply [post]
func (api *TemplateAPI) ApplyTemplate(c *gin.Context) {
	id := c.Param("id")

	templateID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", "无效的模板ID")
		return
	}

	var req ApplyTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取模板
	template, err := api.service.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "模板不存在", err.Error())
		return
	}

	// 创建渲染器
	renderer := document.NewSimpleTemplateRenderer()

	// 渲染模板
	renderedContent, err := renderer.Render(template.Content, req.Variables)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "渲染失败", err.Error())
		return
	}

	// 返回渲染结果
	response := &ApplyTemplateResponse{
		TemplateID:      templateID,
		RenderedContent: renderedContent,
		Variables:       req.Variables,
	}

	shared.Success(c, http.StatusOK, "应用成功", response)
}

// ==================== DTO 定义 ====================

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	ProjectID   string                  `json:"projectId,omitempty" validate:"omitempty,hexadecimal"`
	Name        string                  `json:"name" validate:"required,max=200"`
	Description string                  `json:"description" validate:"max=1000"`
	Type        string                  `json:"type" validate:"required,oneof=chapter outline setting"`
	Category    string                  `json:"category" validate:"max=100"`
	Content     string                  `json:"content" validate:"required"`
	Variables   []TemplateVariableDTO   `json:"variables"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name        *string                   `json:"name,omitempty" validate:"omitempty,max=200"`
	Description *string                   `json:"description,omitempty" validate:"omitempty,max=1000"`
	Category    *string                   `json:"category,omitempty" validate:"omitempty,max=100"`
	Content     *string                   `json:"content,omitempty"`
	Variables   *[]TemplateVariableDTO    `json:"variables,omitempty"`
}

// ApplyTemplateRequest 应用模板请求
type ApplyTemplateRequest struct {
	DocumentID string            `json:"documentId" validate:"required,hexadecimal"`
	Variables  map[string]string `json:"variables"`
}

// ApplyTemplateResponse 应用模板响应
type ApplyTemplateResponse struct {
	TemplateID      primitive.ObjectID `json:"templateId"`
	RenderedContent string            `json:"renderedContent"`
	Variables       map[string]string `json:"variables"`
}

// TemplateVariableDTO 模板变量DTO
type TemplateVariableDTO struct {
	Name         string  `json:"name" validate:"required,max=50"`
	Label        string  `json:"label" validate:"required,max=100"`
	Type         string  `json:"type" validate:"required,oneof=text textarea select number"`
	Placeholder  string  `json:"placeholder" validate:"max=200"`
	DefaultValue string  `json:"defaultValue,omitempty"`
	Required     bool    `json:"required"`
	Order        int     `json:"order"`
	Options      []OptionDTO `json:"options,omitempty"`
}

// OptionDTO 选项DTO
type OptionDTO struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// convertVariables 转换变量DTO到模型
func convertVariables(dtos []TemplateVariableDTO) []writer.TemplateVariable {
	variables := make([]writer.TemplateVariable, len(dtos))
	for i, dto := range dtos {
		variables[i] = writer.TemplateVariable{
			Name:         dto.Name,
			Label:        dto.Label,
			Type:         dto.Type,
			Placeholder:  dto.Placeholder,
			DefaultValue: dto.DefaultValue,
			Required:     dto.Required,
			Order:        dto.Order,
		}

		if len(dto.Options) > 0 {
			variables[i].Options = make([]writer.Option, len(dto.Options))
			for j, opt := range dto.Options {
				variables[i].Options[j] = writer.Option{
					Label: opt.Label,
					Value: opt.Value,
				}
			}
		}
	}
	return variables
}
