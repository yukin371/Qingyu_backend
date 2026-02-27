package writer

import (
	"github.com/gin-gonic/gin"

	writerModels "Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// OutlineApi 大纲API处理器
type OutlineApi struct {
	outlineService interfaces.OutlineService
}

// NewOutlineApi 创建OutlineApi实例
func NewOutlineApi(outlineService interfaces.OutlineService) *OutlineApi {
	return &OutlineApi{
		outlineService: outlineService,
	}
}

// CreateOutline 创建大纲节点
// @Summary 创建大纲节点
// @Description 在项目中创建一个新的大纲节点
// @Tags 大纲管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param request body interfaces.CreateOutlineRequest true "创建大纲请求"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/writer/projects/{projectId}/outlines [post]
func (api *OutlineApi) CreateOutline(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateOutlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	outline, err := api.outlineService.Create(c.Request.Context(), projectID, userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, outline)
}

// GetOutline 获取大纲详情
// @Summary 获取大纲详情
// @Description 根据ID获取大纲节点详细信息
// @Tags 大纲管理
// @Accept json
// @Produce json
// @Param outlineId path string true "大纲ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/writer/outlines/{outlineId} [get]
func (api *OutlineApi) GetOutline(c *gin.Context) {
	outlineID := c.Param("outlineId")
	projectID := c.Query("projectId")

	if outlineID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "outlineId和projectId不能为空")
		return
	}

	outline, err := api.outlineService.GetByID(c.Request.Context(), outlineID, projectID)
	if err != nil {
		response.NotFound(c, "大纲不存在")
		return
	}

	response.Success(c, outline)
}

// ListOutlines 获取项目大纲列表
// @Summary 获取项目大纲列表
// @Description 获取指定项目的所有大纲节点（扁平列表）
// @Tags 大纲管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} response.Response
// @Router /api/v1/writer/projects/{projectId}/outlines [get]
func (api *OutlineApi) ListOutlines(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	outlines, err := api.outlineService.List(c.Request.Context(), projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, outlines)
}

// GetOutlineTree 获取大纲树
// @Summary 获取大纲树
// @Description 获取项目的完整大纲树形结构
// @Tags 大纲管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} response.Response
// @Router /api/v1/writer/projects/{projectId}/outlines/tree [get]
func (api *OutlineApi) GetOutlineTree(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	tree, err := api.outlineService.GetTree(c.Request.Context(), projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, tree)
}

// UpdateOutline 更新大纲
// @Summary 更新大纲
// @Description 更新大纲节点信息
// @Tags 大纲管理
// @Accept json
// @Produce json
// @Param outlineId path string true "大纲ID"
// @Param projectId query string true "项目ID"
// @Param request body interfaces.UpdateOutlineRequest true "更新大纲请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/writer/outlines/{outlineId} [put]
func (api *OutlineApi) UpdateOutline(c *gin.Context) {
	outlineID := c.Param("outlineId")
	projectID := c.Query("projectId")

	if outlineID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "outlineId和projectId不能为空")
		return
	}

	var req interfaces.UpdateOutlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	outline, err := api.outlineService.Update(c.Request.Context(), outlineID, projectID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, outline)
}

// DeleteOutline 删除大纲
// @Summary 删除大纲
// @Description 删除指定大纲节点（如果有子节点则拒绝删除）
// @Tags 大纲管理
// @Accept json
// @Produce json
// @Param outlineId path string true "大纲ID"
// @Param projectId query string true "项目ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/writer/outlines/{outlineId} [delete]
func (api *OutlineApi) DeleteOutline(c *gin.Context) {
	outlineID := c.Param("outlineId")
	projectID := c.Query("projectId")

	if outlineID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "outlineId和projectId不能为空")
		return
	}

	err := api.outlineService.Delete(c.Request.Context(), outlineID, projectID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetOutlineChildren 获取子节点列表
// @Summary 获取子节点列表
// @Description 获取指定父节点的所有子节点
// @Tags 大纲管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param parentId query string false "父节点ID（不传则返回根节点）"
// @Success 200 {object} response.Response
// @Router /api/v1/writer/projects/{projectId}/outlines/children [get]
func (api *OutlineApi) GetOutlineChildren(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	parentID := c.Query("parentId")

	children, err := api.outlineService.GetChildren(c.Request.Context(), projectID, parentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, children)
}

var _ = writerModels.OutlineNode{}
