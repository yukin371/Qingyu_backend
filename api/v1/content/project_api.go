package content

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/api/v1/shared"
	contentService "Qingyu_backend/service/interfaces/content"
)

// ProjectAPI 项目管理API
type ProjectAPI struct {
	projectService contentService.ProjectServicePort
}

// NewProjectAPI 创建项目API实例
func NewProjectAPI(projectService contentService.ProjectServicePort) *ProjectAPI {
	return &ProjectAPI{
		projectService: projectService,
	}
}

// CreateProject 创建项目
//
//	@Summary		创建项目
//	@Description	创建新的写作项目
//	@Tags			内容管理-项目
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateProjectRequest	true	"创建项目请求"
//	@Success		201		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/projects [post]
func (api *ProjectAPI) CreateProject(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Unauthorized(c, "请先登录")
		return
	}
	_ = userID // TODO: 使用userID

	result, err := api.projectService.CreateProject(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 201, "创建成功", result)
}

// GetProject 获取项目详情
//
//	@Summary		获取项目详情
//	@Description	根据项目ID获取项目详细信息
//	@Tags			内容管理-项目
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"项目ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/projects/{id} [get]
func (api *ProjectAPI) GetProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	result, err := api.projectService.GetProject(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// UpdateProject 更新项目
//
//	@Summary		更新项目
//	@Description	更新项目信息
//	@Tags			内容管理-项目
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"项目ID"
//	@Param			request	body		dto.UpdateProjectRequest	true	"更新项目请求"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/projects/{id} [put]
func (api *ProjectAPI) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	result, err := api.projectService.UpdateProject(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "更新成功", result)
}

// DeleteProject 删除项目
//
//	@Summary		删除项目
//	@Description	软删除指定项目
//	@Tags			内容管理-项目
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"项目ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/projects/{id} [delete]
func (api *ProjectAPI) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	err := api.projectService.DeleteProject(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "删除成功", nil)
}

// ListProjects 获取项目列表
//
//	@Summary		获取项目列表
//	@Description	分页获取项目列表，支持按状态和分类筛选
//	@Tags			内容管理-项目
//	@Accept			json
//	@Produce		json
//	@Param			status		query		string	false	"状态"
//	@Param			category	query		string	false	"分类"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			pageSize	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/projects [get]
func (api *ProjectAPI) ListProjects(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	req := &dto.ListProjectsRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Category: category,
	}

	result, err := api.projectService.ListProjects(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Paginated(c, result.Projects, result.Total, page, pageSize, "获取成功")
}

// GetProjectStatistics 获取项目统计
//
//	@Summary		获取项目统计
//	@Description	获取项目的详细统计数据
//	@Tags			内容管理-项目
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"项目ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/projects/{id}/statistics [get]
func (api *ProjectAPI) GetProjectStatistics(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	result, err := api.projectService.GetProjectStatistics(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// UpdateProjectStatistics 更新项目统计
//
//	@Summary		更新项目统计
//	@Description	更新项目的统计数据
//	@Tags			内容管理-项目
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"项目ID"
//	@Param			request	body		dto.ProjectStatistics	true	"统计数据"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/projects/{id}/statistics [put]
func (api *ProjectAPI) UpdateProjectStatistics(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	var stats dto.ProjectStatistics
	if err := c.ShouldBindJSON(&stats); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	err := api.projectService.UpdateProjectStatistics(c.Request.Context(), id, &stats)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "更新成功", nil)
}

// 以下是扩展功能的占位实现，使用适配器委托给现有服务

// DuplicateProject 复制项目
func (api *ProjectAPI) DuplicateProject(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// ArchiveProject 归档项目
func (api *ProjectAPI) ArchiveProject(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// UnarchiveProject 取消归档项目
func (api *ProjectAPI) UnarchiveProject(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// ListCollaborators 获取协作者列表
func (api *ProjectAPI) ListCollaborators(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// AddCollaborator 添加协作者
func (api *ProjectAPI) AddCollaborator(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// RemoveCollaborator 移除协作者
func (api *ProjectAPI) RemoveCollaborator(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// UpdateCollaboratorRole 更新协作者角色
func (api *ProjectAPI) UpdateCollaboratorRole(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// GetExportTasks 获取导出任务列表
func (api *ProjectAPI) GetExportTasks(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}

// ExportProject 导出项目
func (api *ProjectAPI) ExportProject(c *gin.Context) {
	c.JSON(501, gin.H{"code": 501, "message": "功能暂未实现"})
}
