package writer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/document"
)

// DocumentApi 文档API
type DocumentApi struct {
	documentService *document.DocumentService
}

// NewDocumentApi 创建文档API
func NewDocumentApi(documentService *document.DocumentService) *DocumentApi {
	return &DocumentApi{
		documentService: documentService,
	}
}

// CreateDocument 创建文档
// @Summary 创建文档
// @Description 在项目中创建新文档
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param request body document.CreateDocumentRequest true "创建文档请求"
// @Success 201 {object} shared.Response{data=document.CreateDocumentResponse}
// @Failure 400 {object} shared.Response
// @Router /api/v1/projects/{projectId}/documents [post]
func (api *DocumentApi) CreateDocument(c *gin.Context) {
	projectID := c.Param("projectId")

	var req document.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	req.ProjectID = projectID

	resp, err := api.documentService.CreateDocument(c.Request.Context(), &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", resp)
}

// GetDocument 获取文档详情
// @Summary 获取文档详情
// @Description 根据ID获取文档详细信息
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} shared.Response{data=document.Document}
// @Router /api/v1/documents/{id} [get]
func (api *DocumentApi) GetDocument(c *gin.Context) {
	documentID := c.Param("id")

	doc, err := api.documentService.GetDocument(c.Request.Context(), documentID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", doc)
}

// GetDocumentTree 获取文档树
// @Summary 获取文档树
// @Description 获取项目的文档树形结构
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} shared.Response{data=document.DocumentTreeResponse}
// @Router /api/v1/projects/{projectId}/documents/tree [get]
func (api *DocumentApi) GetDocumentTree(c *gin.Context) {
	projectID := c.Param("projectId")

	resp, err := api.documentService.GetDocumentTree(c.Request.Context(), projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", resp)
}

// UpdateDocument 更新文档
// @Summary 更新文档
// @Description 更新文档信息
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body document.UpdateDocumentRequest true "更新文档请求"
// @Success 200 {object} shared.Response
// @Router /api/v1/documents/{id} [put]
func (api *DocumentApi) UpdateDocument(c *gin.Context) {
	documentID := c.Param("id")

	var req document.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if err := api.documentService.UpdateDocument(c.Request.Context(), documentID, &req); err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// DeleteDocument 删除文档
// @Summary 删除文档
// @Description 删除文档（软删除）
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} shared.Response
// @Router /api/v1/documents/{id} [delete]
func (api *DocumentApi) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")

	if err := api.documentService.DeleteDocument(c.Request.Context(), documentID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}
