package writer

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	documentModel "Qingyu_backend/models/writer" // Import for Swagger annotations
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
// @Param request body project.CreateProjectRequest true "创建项目请求"
// @Success 201 {object} shared.APIResponse{data=project.CreateProjectResponse}
// @Failure 400 {object} shared.APIResponse
// @Failure 401 {object} shared.APIResponse
// @Router /api/v1/projects [post]
func (api *ProjectApi) CreateProject(c *gin.Context) {
	var req project.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从gin.Context获取userId并添加到context.Context
	ctx := c.Request.Context()
	if userID, exists := c.Get("userId"); exists {
		ctx = context.WithValue(ctx, "userId", userID)
	}

	resp, err := api.projectService.CreateProject(ctx, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", resp)
}

// GetProject 获取项目详情
// @Summary 获取项目详情
// @Description 根据ID获取项目详细信息
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} shared.APIResponse{data=documentModel.Project}
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/projects/{id} [get]
func (api *ProjectApi) GetProject(c *gin.Context) {
	projectID := c.Param("id")

	project, err := api.projectService.GetProject(c.Request.Context(), projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", project)
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
// @Param category query string false "项目分类"
// @Success 200 {object} shared.APIResponse{data=project.ListProjectsResponse}
// @Router /api/v1/projects [get]
func (api *ProjectApi) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	category := c.Query("category")

	req := &project.ListProjectsRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Category: category,
	}

	// 从gin.Context获取userId并添加到context.Context
	ctx := c.Request.Context()
	if userID, exists := c.Get("userId"); exists {
		ctx = context.WithValue(ctx, "userId", userID)
	}

	resp, err := api.projectService.ListMyProjects(ctx, req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", resp)
}

// UpdateProject 更新项目
// @Summary 更新项目
// @Description 更新项目信息
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body project.UpdateProjectRequest true "更新项目请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 403 {object} shared.APIResponse
// @Router /api/v1/projects/{id} [put]
func (api *ProjectApi) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")

	var req project.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从gin.Context获取userId并添加到context.Context
	ctx := c.Request.Context()
	if userID, exists := c.Get("userId"); exists {
		ctx = context.WithValue(ctx, "userId", userID)
	}

	if err := api.projectService.UpdateProject(ctx, projectID, &req); err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// DeleteProject 删除项目
// @Summary 删除项目
// @Description 删除项目（软删除）
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Failure 403 {object} shared.APIResponse
// @Router /api/v1/projects/{id} [delete]
func (api *ProjectApi) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	// 从gin.Context获取userId并添加到context.Context
	ctx := c.Request.Context()
	if userID, exists := c.Get("userId"); exists {
		ctx = context.WithValue(ctx, "userId", userID)
	}

	if err := api.projectService.DeleteProject(ctx, projectID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// UpdateProjectStatistics 更新项目统计信息
// @Summary 更新项目统计信息
// @Description 更新项目的统计信息（字数、章节数等）
// @Tags 项目管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/projects/{id}/statistics [put]
func (api *ProjectApi) UpdateProjectStatistics(c *gin.Context) {
	projectID := c.Param("id")

	// 从gin.Context获取userId并添加到context.Context
	ctx := c.Request.Context()
	if userID, exists := c.Get("userId"); exists {
		ctx = context.WithValue(ctx, "userId", userID)
	}

	// 调用Service计算并更新统计信息
	if err := api.projectService.RecalculateProjectStatistics(ctx, projectID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新统计失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

var _ = documentModel.Project{}
