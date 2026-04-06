package writer

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/dto"
	documentModel "Qingyu_backend/models/writer" // Import for Swagger annotations
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/writer/project"
)

// ProjectApi 项目API
type ProjectApi struct {
	projectService *project.ProjectService
}

// NewProjectApi 创建项目API
func NewProjectApi(projectService *project.ProjectService) *ProjectApi {
	return &ProjectApi{
		projectService: projectService,
	}
}

// CreateProject 创建项目
// @Summary 创建项目
// @Description 创建一个新的写作项目
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param request body dto.CreateProjectRequest true "创建项目请求"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/projects [post]
func (api *ProjectApi) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	ctx := shared.AddUserIDToContext(c)

	projectModel, err := api.projectService.CreateProject(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	// 将模型转换为DTO响应
	resp := dto.ToProjectResponse(projectModel)
	response.Created(c, resp)
}

// GetProject 获取项目详情
// @Summary 获取项目详情
// @Description 根据ID获取项目详细信息
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/projects/{id} [get]
func (api *ProjectApi) GetProject(c *gin.Context) {
	projectID := c.Param("id")

	ctx := shared.AddUserIDToContext(c)

	projectModel, err := api.projectService.GetProject(ctx, projectID)
	if err != nil {
		c.Error(err)
		return
	}

	// 将模型转换为DTO响应
	resp := dto.ToProjectResponse(projectModel)
	response.Success(c, resp)
}

// ListProjects 获取项目列表
// @Summary 获取项目列表
// @Description 获取当前用户的项目列表（支持分页和筛选）
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query string false "项目状态"
// @Param sort query string false "排序字段"
// @Param order query string false "排序方向"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects [get]
func (api *ProjectApi) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	sort := c.Query("sort")
	order := c.Query("order")

	req := &dto.ListProjectsRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Sort:     sort,
		Order:    order,
	}

	ctx := shared.AddUserIDToContext(c)

	resp, err := api.projectService.ListMyProjects(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	// 将服务层响应转换为DTO响应
	dtoResp := &dto.ProjectListResponse{
		Items:    dto.ToProjectResponseList(resp.Projects),
		Total:    resp.Total,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}
	response.Success(c, dtoResp)
}

// UpdateProject 更新项目
// @Summary 更新项目
// @Description 更新项目信息
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body object true "更新项目请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/projects/{id} [put]
func (api *ProjectApi) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")

	var req dto.UpdateProjectRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	ctx := shared.AddUserIDToContext(c)

	if err := api.projectService.UpdateProject(ctx, projectID, &req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// DeleteProject 删除项目
// @Summary 删除项目
// @Description 删除项目（软删除）
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/projects/{id} [delete]
func (api *ProjectApi) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	ctx := shared.AddUserIDToContext(c)

	if err := api.projectService.DeleteProject(ctx, projectID); err != nil {
		c.Error(err)
		return
	}

	// 返回成功响应（包含删除的项目ID）
	response.Success(c, map[string]interface{}{
		"projectId": projectID,
		"deleted":   true,
	})
}

// UpdateProjectStatistics 更新项目统计信息
// @Summary 更新项目统计信息
// @Description 更新项目的统计信息（字数、章节数等）
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects/{id}/statistics [put]
func (api *ProjectApi) UpdateProjectStatistics(c *gin.Context) {
	projectID := c.Param("id")

	ctx := shared.AddUserIDToContext(c)

	// 调用Service计算并更新统计信息
	if err := api.projectService.RecalculateProjectStatistics(ctx, projectID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

var _ = documentModel.Project{}
