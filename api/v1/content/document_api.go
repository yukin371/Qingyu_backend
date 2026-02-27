package content

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/api/v1/shared"
	contentService "Qingyu_backend/service/interfaces/content"
)

// DocumentAPI 文档管理API
type DocumentAPI struct {
	documentService contentService.DocumentServicePort
}

// NewDocumentAPI 创建文档API实例
func NewDocumentAPI(documentService contentService.DocumentServicePort) *DocumentAPI {
	return &DocumentAPI{
		documentService: documentService,
	}
}

// CreateDocument 创建文档
//
//	@Summary		创建文档
//	@Description	创建新的文档节点（章节/场景/笔记）
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateDocumentRequest	true	"创建文档请求"
//	@Success		201		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/documents [post]
func (api *DocumentAPI) CreateDocument(c *gin.Context) {
	var req dto.CreateDocumentRequest
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

	result, err := api.documentService.CreateDocument(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 201, "创建成功", result)
}

// GetDocument 获取文档详情
//
//	@Summary		获取文档详情
//	@Description	根据文档ID获取文档详细信息（不含内容）
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"文档ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id} [get]
func (api *DocumentAPI) GetDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	result, err := api.documentService.GetDocument(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// UpdateDocument 更新文档
//
//	@Summary		更新文档
//	@Description	更新文档元数据信息
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"文档ID"
//	@Param			request	body		dto.UpdateDocumentRequest	true	"更新文档请求"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id} [put]
func (api *DocumentAPI) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	var req dto.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	result, err := api.documentService.UpdateDocument(c.Request.Context(), id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "更新成功", result)
}

// DeleteDocument 删除文档
//
//	@Summary		删除文档
//	@Description	软删除指定文档
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"文档ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id} [delete]
func (api *DocumentAPI) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	err := api.documentService.DeleteDocument(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "删除成功", nil)
}

// ListDocuments 获取文档列表
//
//	@Summary		获取文档列表
//	@Description	分页获取文档列表，支持按父节点和状态筛选
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			projectId	query		string	false	"项目ID"
//	@Param			parentId	query		string	false	"父节点ID"
//	@Param			status		query		string	false	"状态"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			pageSize	query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/documents [get]
func (api *DocumentAPI) ListDocuments(c *gin.Context) {
	projectID := c.Query("projectId")
	parentID := c.Query("parentId")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	req := &dto.ListDocumentsRequest{
		ProjectID: projectID,
		ParentID:  parentID,
		Page:      page,
		PageSize:  pageSize,
		Status:    status,
	}

	result, err := api.documentService.ListDocuments(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Paginated(c, result.Documents, int64(result.Total), page, pageSize, "获取成功")
}

// DuplicateDocument 复制文档
//
//	@Summary		复制文档
//	@Description	复制指定文档及其内容
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"文档ID"
//	@Param			request	body		dto.DuplicateRequest	true	"复制请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id}/duplicate [post]
func (api *DocumentAPI) DuplicateDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	var req dto.DuplicateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	result, err := api.documentService.DuplicateDocument(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "复制成功", result)
}

// MoveDocument 移动文档
//
//	@Summary		移动文档
//	@Description	移动文档到新的父节点
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"文档ID"
//	@Param			request	body		dto.MoveDocumentRequest	true	"移动请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id}/move [put]
func (api *DocumentAPI) MoveDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	var req dto.MoveDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	req.DocumentID = id

	err := api.documentService.MoveDocument(c.Request.Context(), req.DocumentID, req.NewParentID, req.Order)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "移动成功", nil)
}

// GetDocumentTree 获取文档树
//
//	@Summary		获取文档树
//	@Description	获取项目的完整文档树形结构
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			projectId	path		string	true	"项目ID"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/tree/{projectId} [get]
func (api *DocumentAPI) GetDocumentTree(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.BadRequest(c, "参数错误", "项目ID不能为空")
		return
	}

	result, err := api.documentService.GetDocumentTree(c.Request.Context(), projectID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetDocumentContent 获取文档内容
//
//	@Summary		获取文档内容
//	@Description	获取文档的实际文本内容
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"文档ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		404	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id}/content [get]
func (api *DocumentAPI) GetDocumentContent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	result, err := api.documentService.GetDocumentContent(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// UpdateDocumentContent 更新文档内容
//
//	@Summary		更新文档内容
//	@Description	更新文档的文本内容
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"文档ID"
//	@Param			request	body		dto.UpdateContentRequest	true	"更新内容请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id}/content [put]
func (api *DocumentAPI) UpdateDocumentContent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	var req dto.UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	req.DocumentID = id

	err := api.documentService.UpdateDocumentContent(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "更新成功", nil)
}

// AutoSaveDocument 自动保存文档
//
//	@Summary		自动保存文档
//	@Description	定时自动保存文档内容
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.AutoSaveRequest	true	"自动保存请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/autosave [post]
func (api *DocumentAPI) AutoSaveDocument(c *gin.Context) {
	var req dto.AutoSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.BadRequest(c, "参数错误", err.Error())
		return
	}

	result, err := api.documentService.AutoSaveDocument(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "保存成功", result)
}

// GetVersionHistory 获取版本历史
//
//	@Summary		获取版本历史
//	@Description	获取文档的版本历史记录
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"文档ID"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			pageSize	query		int		false	"每页数量"	default(20)
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		500		{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id}/versions [get]
func (api *DocumentAPI) GetVersionHistory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := api.documentService.GetVersionHistory(c.Request.Context(), id, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Paginated(c, result.Versions, int64(result.Total), page, pageSize, "获取成功")
}

// RestoreVersion 恢复版本
//
//	@Summary		恢复版本
//	@Description	将文档内容恢复到指定版本
//	@Tags			内容管理-文档
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"文档ID"
//	@Param			versionId	path		string	true	"版本ID"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/documents/{id}/versions/{versionId}/restore [post]
func (api *DocumentAPI) RestoreVersion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "文档ID不能为空")
		return
	}

	versionID := c.Param("versionId")
	if versionID == "" {
		shared.BadRequest(c, "参数错误", "版本ID不能为空")
		return
	}

	err := api.documentService.RestoreVersion(c.Request.Context(), id, versionID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "恢复成功", nil)
}
