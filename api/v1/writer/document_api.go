package writer

import (
	"context"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/writer/document"
	writerModels "Qingyu_backend/models/writer" // Import for Swagger annotations
	"Qingyu_backend/pkg/response"
	"errors"
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
// @Param request body object true "创建文档请求"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/documents [post]
func (api *DocumentApi) CreateDocument(c *gin.Context) {
	// 检查服务是否初始化
	if api.documentService == nil {
		response.InternalError(c, errors.New("服务未初始化: 文档服务未正确初始化"))
		return
	}

	projectID := c.Param("projectId")

	// 验证项目ID
	if projectID == "" {
		response.BadRequest(c,  "参数错误", "项目ID不能为空")
		return
	}

	// 获取并验证用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.BadRequest(c,  "参数错误", "无效的用户ID")
		return
	}

	// 将用户ID添加到context
	ctx := context.WithValue(c.Request.Context(), "userID", userIDStr)

	var req document.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	req.ProjectID = projectID

	resp, err := api.documentService.CreateDocument(ctx, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, resp)
}

// GetDocument 获取文档详情
// @Summary 获取文档详情
// @Description 根据ID获取文档详细信息
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/documents/{id} [get]
func (api *DocumentApi) GetDocument(c *gin.Context) {
	documentID := c.Param("id")

	doc, err := api.documentService.GetDocument(c.Request.Context(), documentID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, doc)
}

// GetDocumentTree 获取文档树
// @Summary 获取文档树
// @Description 获取项目的文档树形结构
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/documents/tree [get]
func (api *DocumentApi) GetDocumentTree(c *gin.Context) {
	projectID := c.Param("projectId")

	resp, err := api.documentService.GetDocumentTree(c.Request.Context(), projectID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, resp)
}

// UpdateDocument 更新文档
// @Summary 更新文档
// @Description 更新文档信息
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body object true "更新文档请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/documents/{id} [put]
func (api *DocumentApi) UpdateDocument(c *gin.Context) {
	documentID := c.Param("id")

	var req document.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	if err := api.documentService.UpdateDocument(c.Request.Context(), documentID, &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteDocument 删除文档
// @Summary 删除文档
// @Description 删除文档（软删除）
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/documents/{id} [delete]
func (api *DocumentApi) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")

	if err := api.documentService.DeleteDocument(c.Request.Context(), documentID); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// ListDocuments 获取文档列表
// @Summary 获取文档列表
// @Description 获取项目下的文档列表（支持分页）
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/documents [get]
func (api *DocumentApi) ListDocuments(c *gin.Context) {
	projectID := c.Param("projectId")

	page, _ := c.GetQuery("page")
	pageSize, _ := c.GetQuery("pageSize")

	req := &document.ListDocumentsRequest{
		ProjectID: projectID,
		Page:      page,
		PageSize:  pageSize,
	}

	resp, err := api.documentService.ListDocuments(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, resp)
}

// MoveDocument 移动文档
// @Summary 移动文档
// @Description 移动文档到新的父节点
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body object true "移动文档请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/documents/{id}/move [put]
func (api *DocumentApi) MoveDocument(c *gin.Context) {
	documentID := c.Param("id")

	var req document.MoveDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	req.DocumentID = documentID

	if err := api.documentService.MoveDocument(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// ReorderDocuments 重新排序文档
// @Summary 重新排序文档
// @Description 批量更新同级文档的顺序
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param projectId path string true "项目ID"
// @Param request body object true "排序请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/projects/{projectId}/documents/reorder [put]
func (api *DocumentApi) ReorderDocuments(c *gin.Context) {
	projectID := c.Param("projectId")

	var req document.ReorderDocumentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	req.ProjectID = projectID

	if err := api.documentService.ReorderDocuments(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, nil)
}

// DuplicateDocument 复制文档
// @Summary 复制文档
// @Description 复制文档到指定位置，可选择是否复制内容
// @Tags 文档管理
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body object true "复制文档请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/writer/documents/{id}/duplicate [post]
func (api *DocumentApi) DuplicateDocument(c *gin.Context) {
	documentID := c.Param("id")

	// 获取并验证用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.BadRequest(c,  "参数错误", "无效的用户ID")
		return
	}

	// 将用户ID添加到context
	ctx := context.WithValue(c.Request.Context(), "userID", userIDStr)

	var req document.DuplicateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", err.Error())
		return
	}

	resp, err := api.documentService.DuplicateDocument(ctx, documentID, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, resp)
}

var _ = writerModels.Document{}
