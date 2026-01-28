package writer

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// ExportApi 导出API处理器
type ExportApi struct {
	exportService interfaces.ExportService
}

// NewExportApi 创建ExportApi实例
func NewExportApi(exportService interfaces.ExportService) *ExportApi {
	return &ExportApi{
		exportService: exportService,
	}
}

// ExportDocument 导出文档
// @Summary 导出文档
// @Description 将文档导出为指定格式（TXT/MD/DOCX）
// @Tags 导出管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param projectId query string true "项目ID"
// @Param request body object true "导出请求"
// @Success 202 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/documents/{id}/export [post]
func (api *ExportApi) ExportDocument(c *gin.Context) {
	documentID := c.Param("id")
	projectID := c.Query("projectId")

	if documentID == "" || projectID == "" {
		response.BadRequest(c, "参数错误", "documentId和projectId不能为空")
		return
	}

	var req interfaces.ExportDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	task, err := api.exportService.ExportDocument(c.Request.Context(), documentID, projectID, userID, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, task)
}

// ExportProject 导出项目
// @Summary 导出项目
// @Description 将整个项目导出为ZIP包
// @Tags 导出管理
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body object true "导出请求"
// @Success 202 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/projects/{id}/export [post]
func (api *ExportApi) ExportProject(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		response.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	var req interfaces.ExportProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	task, err := api.exportService.ExportProject(c.Request.Context(), projectID, userID, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, task)
}

// GetExportTask 获取导出任务状态
// @Summary 获取导出任务状态
// @Description 根据任务ID获取导出任务的详细状态
// @Tags 导出管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/exports/{id} [get]
func (api *ExportApi) GetExportTask(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		response.BadRequest(c, "参数错误", "任务ID不能为空")
		return
	}

	task, err := api.exportService.GetExportTask(c.Request.Context(), taskID)
	if err != nil {
		response.NotFound(c, "导出任务不存在")
		return
	}

	response.Success(c, task)
}

// DownloadExportFile 下载导出文件
// @Summary 下载导出文件
// @Description 下载已完成的导出文件
// @Tags 导出管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/exports/{id}/download [get]
func (api *ExportApi) DownloadExportFile(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		response.BadRequest(c, "参数错误", "任务ID不能为空")
		return
	}

	file, err := api.exportService.DownloadExportFile(c.Request.Context(), taskID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, file)
}

// ListExportTasks 获取项目的导出任务列表
// @Summary 获取导出任务列表
// @Description 获取指定项目的所有导出任务
// @Tags 导出管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Router /api/v1/writer/projects/{projectId}/exports [get]
func (api *ExportApi) ListExportTasks(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		response.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	tasks, total, err := api.exportService.ListExportTasks(c.Request.Context(), projectID, page, pageSize)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Paginated(c, tasks, total, page, pageSize, "获取成功")
}

// DeleteExportTask 删除导出任务
// @Summary 删除导出任务
// @Description 删除指定的导出任务及其文件
// @Tags 导出管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/exports/{id} [delete]
func (api *ExportApi) DeleteExportTask(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		response.BadRequest(c, "参数错误", "任务ID不能为空")
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	err := api.exportService.DeleteExportTask(c.Request.Context(), taskID, userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// CancelExportTask 取消导出任务
// @Summary 取消导出任务
// @Description 取消正在处理或等待中的导出任务
// @Tags 导出管理
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse
// @Failure 404 {object} shared.APIResponse
// @Router /api/v1/writer/exports/{id}/cancel [post]
func (api *ExportApi) CancelExportTask(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		response.BadRequest(c, "参数错误", "任务ID不能为空")
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	err := api.exportService.CancelExportTask(c.Request.Context(), taskID, userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}
