package writer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/project"
)

// VersionApi 版本API
type VersionApi struct {
	versionService *project.VersionService
}

// NewVersionApi 创建版本API
func NewVersionApi(versionService *project.VersionService) *VersionApi {
	return &VersionApi{
		versionService: versionService,
	}
}

// GetVersionHistory 获取版本历史
// @Summary 获取版本历史
// @Description 获取文档的版本历史列表
// @Tags 版本控制
// @Accept json
// @Produce json
// @Param documentId path string true "文档ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} shared.Response{data=project.VersionHistoryResponse}
// @Router /api/v1/documents/{documentId}/versions [get]
func (api *VersionApi) GetVersionHistory(c *gin.Context) {
	documentID := c.Param("documentId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	versions, err := api.versionService.GetVersionHistory(c.Request.Context(), documentID, page, pageSize)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询版本历史失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", versions)
}

// GetVersion 获取特定版本
// @Summary 获取特定版本
// @Description 获取文档的特定版本内容
// @Tags 版本控制
// @Accept json
// @Produce json
// @Param documentId path string true "文档ID"
// @Param versionId path string true "版本ID"
// @Success 200 {object} shared.Response{data=project.VersionDetail}
// @Router /api/v1/documents/{documentId}/versions/{versionId} [get]
func (api *VersionApi) GetVersion(c *gin.Context) {
	documentID := c.Param("documentId")
	versionID := c.Param("versionId")

	version, err := api.versionService.GetVersion(c.Request.Context(), documentID, versionID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询版本失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", version)
}

// CompareVersions 比较版本
// @Summary 比较版本
// @Description 比较两个版本的差异
// @Tags 版本控制
// @Accept json
// @Produce json
// @Param documentId path string true "文档ID"
// @Param fromVersion query string true "源版本ID"
// @Param toVersion query string true "目标版本ID"
// @Success 200 {object} shared.Response{data=project.VersionDiff}
// @Router /api/v1/documents/{documentId}/versions/compare [get]
func (api *VersionApi) CompareVersions(c *gin.Context) {
	documentID := c.Param("documentId")
	fromVersion := c.Query("fromVersion")
	toVersion := c.Query("toVersion")

	if fromVersion == "" || toVersion == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "fromVersion和toVersion不能为空")
		return
	}

	diff, err := api.versionService.CompareVersions(c.Request.Context(), documentID, fromVersion, toVersion)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "比较版本失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "比较成功", diff)
}

// RestoreVersion 恢复版本
// @Summary 恢复版本
// @Description 将文档恢复到特定版本
// @Tags 版本控制
// @Accept json
// @Produce json
// @Param documentId path string true "文档ID"
// @Param versionId path string true "版本ID"
// @Success 200 {object} shared.Response
// @Router /api/v1/documents/{documentId}/versions/{versionId}/restore [post]
func (api *VersionApi) RestoreVersion(c *gin.Context) {
	documentID := c.Param("documentId")
	versionID := c.Param("versionId")

	if err := api.versionService.RestoreVersion(c.Request.Context(), documentID, versionID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "恢复版本失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "恢复成功", nil)
}
