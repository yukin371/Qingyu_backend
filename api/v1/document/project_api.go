package document

import (
	svc "Qingyu_backend/service/project"
	"net/http"
	"strconv"

	model "Qingyu_backend/models/document"

	"github.com/gin-gonic/gin"
)

// ProjectApi 项目相关API
type ProjectApi struct {
	service *svc.ProjectService
}

// NewProjectApi 创建项目相关API实例
func NewProjectApi() *ProjectApi {
	return &ProjectApi{service: &svc.ProjectService{}}
}

// CreateProject 创建项目
// POST /api/v1/document/project
func (a *ProjectApi) CreateProject(c *gin.Context) {
	var req svc.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	created, err := a.service.CreateProject(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetProjectList 获取项目列表
// GET /api/v1/document/project
func (a *ProjectApi) GetProjectList(c *gin.Context) {
	// 简化：从查询参数取 userId、limit、offset
	userID := c.Query("userId")
	// 省略参数校验与转换细节，默认分页
	limit, offset := int64(50), int64(0)
	if l := c.Query("limit"); l != "" {
		if l, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = l
		}
	}
	if o := c.Query("offset"); o != "" {
		if o, err := strconv.ParseInt(o, 10, 64); err == nil {
			offset = o
		}
	}
	projects, err := a.service.GetProjectList(c.Request.Context(), userID, "active", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// GetProjectByID 获取项目详情
// GET /api/v1/document/project/:id
func (a *ProjectApi) GetProjectByID(c *gin.Context) {
	projectID := c.Param("id")
	p, err := a.service.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// UpdateProjectByID 更新项目
// PUT /api/v1/document/project/:id
func (a *ProjectApi) UpdateProjectByID(c *gin.Context) {
	projectID := c.Param("id")
	var req model.Project
	if err := a.service.UpdateProjectByID(c.Request.Context(), projectID, c.GetString("userID"), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project updated"})
}

// DeleteProjectByID 软删除项目
// DELETE /api/v1/document/project/:id
func (a *ProjectApi) DeleteProjectByID(c *gin.Context) {
	projectID := c.Param("id")
	if err := a.service.DeleteProjectByID(c.Request.Context(), projectID, c.GetString("userID")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project deleted"})
}

// RestoreProjectByID 恢复项目
// POST /api/v1/document/project/:id/restore
func (a *ProjectApi) RestoreProjectByID(c *gin.Context) {
	projectID := c.Param("id")
	if err := a.service.RestoreProjectByID(c.Request.Context(), projectID, c.GetString("userID")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project restored"})
}

// DeleteHard 硬删除
// DELETE /api/v1/document/project/:id/hard
func (a *ProjectApi) DeleteHard(c *gin.Context) {
	projectID := c.Param("id")
	if err := a.service.DeleteHard(c.Request.Context(), projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "project deleted hard"})
}
