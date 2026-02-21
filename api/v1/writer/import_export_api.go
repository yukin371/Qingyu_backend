package writer

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// ImportExportApi 导入导出API处理器
type ImportExportApi struct {
	exportService interfaces.ExportService
}

// NewImportExportApi 创建ImportExportApi实例
func NewImportExportApi(exportService interfaces.ExportService) *ImportExportApi {
	return &ImportExportApi{
		exportService: exportService,
	}
}

// ExportProject 导出项目为ZIP
// @Summary 导出项目为ZIP
// @Description 将整个项目导出为ZIP压缩包（直接下载）
// @Tags 导入导出
// @Accept json
// @Produce application/zip
// @Param id path string true "项目ID"
// @Success 200 {file} binary "ZIP文件"
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/projects/{id}/export [get]
func (api *ImportExportApi) ExportProject(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		response.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	// 调用服务导出项目为ZIP
	zipData, err := api.exportService.ExportProjectAsZip(c.Request.Context(), projectID, userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 使用项目ID作为默认文件名
	filename := fmt.Sprintf("project_%s.zip", projectID[:8])

	// 设置响应头
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(zipData)))

	// 直接返回ZIP数据
	c.Data(http.StatusOK, "application/zip", zipData)
}

// ImportProject 从ZIP导入项目
// @Summary 从ZIP导入项目
// @Description 从ZIP压缩包导入项目
// @Tags 导入导出
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "ZIP文件"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/projects/import [post]
func (api *ImportExportApi) ImportProject(c *gin.Context) {
	// 从上下文获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "参数错误", "请上传ZIP文件")
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".zip" {
		response.BadRequest(c, "参数错误", "只支持ZIP格式文件")
		return
	}

	// 读取文件内容
	fileSize := header.Size
	if fileSize > 100*1024*1024 { // 100MB 限制
		response.BadRequest(c, "参数错误", "文件大小不能超过100MB")
		return
	}

	// 读取文件数据
	zipData := make([]byte, fileSize)
	_, err = file.Read(zipData)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 调用服务导入项目
	result, err := api.exportService.ImportProject(c.Request.Context(), userID, zipData)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, result)
}

// sanitizeFilename 清理文件名
func sanitizeFilename(name string) string {
	// 替换不安全字符
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	result := replacer.Replace(name)
	// 移除首尾空格
	result = strings.TrimSpace(result)
	// 限制长度
	if len(result) > 100 {
		result = result[:100]
	}
	if result == "" {
		result = "project"
	}
	return result
}
