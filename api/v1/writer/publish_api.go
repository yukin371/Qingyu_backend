package writer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/interfaces"
)

// PublishApi 发布管理API处理器
type PublishApi struct {
	publishService interfaces.PublishService
}

// NewPublishApi 创建PublishApi实例
func NewPublishApi(publishService interfaces.PublishService) *PublishApi {
	return &PublishApi{
		publishService: publishService,
	}
}

// PublishProject 发布项目到书城
// @Summary 发布项目到书城
// @Description 将项目发布到指定书城平台
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body object true "发布请求"
// @Success 202 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/projects/{id}/publish [post]
func (api *PublishApi) PublishProject(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "项目ID不能为空")
		return
	}

	var req interfaces.PublishProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	record, err := api.publishService.PublishProject(c.Request.Context(), projectID, userID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "发布项目失败", err.Error())
		return
	}

	shared.Success(c, http.StatusAccepted, "发布任务已创建", record)
}

// UnpublishProject 取消发布项目
// @Summary 取消发布项目
// @Description 将项目从书城平台下架
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/projects/{id}/unpublish [post]
func (api *PublishApi) UnpublishProject(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "项目ID不能为空")
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	err := api.publishService.UnpublishProject(c.Request.Context(), projectID, userID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "取消发布失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "取消发布成功", nil)
}

// GetProjectPublicationStatus 获取项目发布状态
// @Summary 获取项目发布状态
// @Description 获取项目的发布状态和统计信息
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/projects/{id}/publication-status [get]
func (api *PublishApi) GetProjectPublicationStatus(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "项目ID不能为空")
		return
	}

	status, err := api.publishService.GetProjectPublicationStatus(c.Request.Context(), projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取发布状态失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", status)
}

// PublishDocument 发布文档（章节）
// @Summary 发布文档（章节）
// @Description 将单个文档发布到书城平台
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param projectId query string true "项目ID"
// @Param request body object true "发布请求"
// @Success 202 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/documents/{id}/publish [post]
func (api *PublishApi) PublishDocument(c *gin.Context) {
	documentID := c.Param("id")
	projectID := c.Query("projectId")

	if documentID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "documentId和projectId不能为空")
		return
	}

	var req interfaces.PublishDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	record, err := api.publishService.PublishDocument(c.Request.Context(), documentID, projectID, userID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "发布文档失败", err.Error())
		return
	}

	shared.Success(c, http.StatusAccepted, "发布任务已创建", record)
}

// UpdateDocumentPublishStatus 更新文档发布状态
// @Summary 更新文档发布状态
// @Description 更新文档的发布状态（发布/取消发布、免费/付费等）
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param projectId query string true "项目ID"
// @Param request body object true "更新请求"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/documents/{id}/publish-status [put]
func (api *PublishApi) UpdateDocumentPublishStatus(c *gin.Context) {
	documentID := c.Param("id")
	projectID := c.Query("projectId")

	if documentID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "documentId和projectId不能为空")
		return
	}

	var req interfaces.UpdateDocumentPublishStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	err := api.publishService.UpdateDocumentPublishStatus(c.Request.Context(), documentID, projectID, userID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新发布状态失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// BatchPublishDocuments 批量发布文档
// @Summary 批量发布文档
// @Description 批量发布多个文档到书城平台
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param request body object true "批量发布请求"
// @Success 202 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/projects/{projectId}/documents/batch-publish [post]
func (api *PublishApi) BatchPublishDocuments(c *gin.Context) {
	projectID := c.Param("projectId")

	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "项目ID不能为空")
		return
	}

	var req interfaces.BatchPublishDocumentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	result, err := api.publishService.BatchPublishDocuments(c.Request.Context(), projectID, userID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "批量发布失败", err.Error())
		return
	}

	shared.Success(c, http.StatusAccepted, "批量发布任务已创建", result)
}

// GetPublicationRecords 获取发布记录列表
// @Summary 获取发布记录列表
// @Description 获取项目的所有发布记录
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Router /api/v1/writer/projects/{projectId}/publications [get]
func (api *PublishApi) GetPublicationRecords(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "项目ID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	records, total, err := api.publishService.GetPublicationRecords(c.Request.Context(), projectID, page, pageSize)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取发布记录失败", err.Error())
		return
	}

	shared.Paginated(c, records, total, page, pageSize, "获取成功")
}

// GetPublicationRecord 获取发布记录详情
// @Summary 获取发布记录详情
// @Description 根据ID获取发布记录的详细信息
// @Tags 发布管理
// @Accept json
// @Produce json
// @Param id path string true "记录ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/publications/{id} [get]
func (api *PublishApi) GetPublicationRecord(c *gin.Context) {
	recordID := c.Param("id")

	if recordID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "记录ID不能为空")
		return
	}

	record, err := api.publishService.GetPublicationRecord(c.Request.Context(), recordID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "发布记录不存在", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", record)
}
