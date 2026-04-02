package writer

import (
	documentModel "Qingyu_backend/models/writer"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/writer/document"
)

// EditorApi 编辑器API
type EditorApi struct {
	documentService  *document.DocumentService
	wordCountService *document.WordCountService
	shortcutService  *document.ShortcutService
}

// NewEditorApi 创建编辑器API
func NewEditorApi(documentService *document.DocumentService) *EditorApi {
	return &EditorApi{
		documentService:  documentService,
		wordCountService: document.NewWordCountService(),
		shortcutService:  document.NewShortcutService(),
	}
}

// AutoSaveDocument 自动保存文档
// @Summary 自动保存文档
// @Description 自动保存文档内容（支持版本冲突检测）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body object true "自动保存请求"
// @Success 200 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse "版本冲突"
// @Router /api/v1/writer/documents/{id}/autosave [post]
func (api *EditorApi) AutoSaveDocument(c *gin.Context) {
	documentID := c.Param("id")

	var req document.AutoSaveRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	req.DocumentID = documentID

	ctx := shared.AddUserIDToContext(c)

	resp, err := api.documentService.AutoSaveDocument(ctx, &req)
	if err != nil {
		// 检查是否是版本冲突
		if err.Error() == "版本冲突" {
			response.Conflict(c, "版本冲突", "文档已被其他用户修改，请刷新后重试")
			return
		}
		c.Error(err)
		return
	}

	response.Success(c, resp)
}

// GetSaveStatus 获取保存状态
// @Summary 获取保存状态
// @Description 获取文档的保存状态和最后保存时间
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/documents/{id}/save-status [get]
func (api *EditorApi) GetSaveStatus(c *gin.Context) {
	documentID := c.Param("id")

	ctx := shared.AddUserIDToContext(c)

	status, err := api.documentService.GetSaveStatus(ctx, documentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, status)
}

// GetDocumentContent 获取文档内容
// Deprecated: 使用 GET /api/v1/writer/documents/{id}/contents
// @Summary 获取文档内容
// @Description 获取文档的完整内容（用于编辑器加载）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/documents/{id}/content [get]
func (api *EditorApi) GetDocumentContent(c *gin.Context) {
	documentID := c.Param("id")
	c.Header("X-API-Deprecated", "true")
	c.Header("X-API-Replacement", "/api/v1/writer/documents/{id}/contents")

	ctx := shared.AddUserIDToContext(c)

	content, err := api.documentService.GetDocumentContent(ctx, documentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, content)
}

// UpdateDocumentContent 更新文档内容
// Deprecated: 使用 PUT /api/v1/writer/documents/{id}/contents
// @Summary 更新文档内容
// @Description 手动更新文档内容（非自动保存）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body object true "更新内容请求"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/documents/{id}/content [put]
func (api *EditorApi) UpdateDocumentContent(c *gin.Context) {
	documentID := c.Param("id")
	c.Header("X-API-Deprecated", "true")
	c.Header("X-API-Replacement", "/api/v1/writer/documents/{id}/contents")

	var req document.UpdateContentRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	req.DocumentID = documentID

	ctx := shared.AddUserIDToContext(c)

	if err := api.documentService.UpdateDocumentContent(ctx, &req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetDocumentContents 获取文档分段内容
func (api *EditorApi) GetDocumentContents(c *gin.Context) {
	documentID := c.Param("id")

	ctx := shared.AddUserIDToContext(c)

	contents, err := api.documentService.GetDocumentContents(ctx, documentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, contents)
}

// ReplaceDocumentContents 批量替换文档分段内容
func (api *EditorApi) ReplaceDocumentContents(c *gin.Context) {
	documentID := c.Param("id")

	var req document.ReplaceDocumentContentsRequest
	if !shared.BindJSON(c, &req) {
		return
	}
	req.DocumentID = documentID

	ctx := shared.AddUserIDToContext(c)

	resp, err := api.documentService.ReplaceDocumentContents(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, resp)
}

// ReindexDocumentContents 重建段落顺序
func (api *EditorApi) ReindexDocumentContents(c *gin.Context) {
	documentID := c.Param("id")

	ctx := shared.AddUserIDToContext(c)

	resp, err := api.documentService.ReindexDocumentContents(ctx, documentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, resp)
}

// CalculateWordCount 计算字数
// @Summary 计算字数
// @Description 计算文档内容的字数统计（支持Markdown过滤）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body WordCountRequest false "字数统计请求（可选，不传则自动获取文档内容）"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/documents/{id}/word-count [post]
func (api *EditorApi) CalculateWordCount(c *gin.Context) {
	var req WordCountRequest
	// 允许空请求体，不传则自动从文档获取内容
	_ = c.ShouldBindJSON(&req)

	content := req.Content
	if content == "" {
		// 从文档内容自动获取
		docID := c.Param("id")
		docContent, err := api.documentService.GetDocumentContent(c.Request.Context(), docID)
		if err != nil {
			response.BadRequest(c, "获取文档内容失败", err.Error())
			return
		}
		content = docContent.Content
	}

	var result *document.WordCountResult
	if req.FilterMarkdown {
		result = api.wordCountService.CalculateWordCountWithMarkdown(content)
	} else {
		result = api.wordCountService.CalculateWordCount(content)
	}

	response.Success(c, result)
}

// GetUserShortcuts 获取用户快捷键配置
// @Summary 获取用户快捷键配置
// @Description 获取当前用户的快捷键配置（包括自定义和默认）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/user/shortcuts [get]
func (api *EditorApi) GetUserShortcuts(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未授权")
		return
	}

	config, err := api.shortcutService.GetUserShortcuts(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, config)
}

// UpdateUserShortcuts 更新用户快捷键配置
// @Summary 更新用户快捷键配置
// @Description 更新用户的自定义快捷键配置
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param request body UpdateShortcutsRequest true "快捷键配置"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/user/shortcuts [put]
func (api *EditorApi) UpdateUserShortcuts(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未授权")
		return
	}

	var req UpdateShortcutsRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	if err := api.shortcutService.UpdateUserShortcuts(c.Request.Context(), userID, req.Shortcuts); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// ResetUserShortcuts 重置用户快捷键配置
// @Summary 重置用户快捷键配置
// @Description 重置用户快捷键为默认配置
// @Tags 编辑器
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/user/shortcuts/reset [post]
func (api *EditorApi) ResetUserShortcuts(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未授权")
		return
	}

	if err := api.shortcutService.ResetUserShortcuts(c.Request.Context(), userID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetShortcutHelp 获取快捷键帮助
// @Summary 获取快捷键帮助
// @Description 获取快捷键帮助文档（按分类）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /api/v1/writer/user/shortcuts/help [get]
func (api *EditorApi) GetShortcutHelp(c *gin.Context) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "未授权")
		return
	}

	help, err := api.shortcutService.GetShortcutHelp(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, help)
}

// WordCountRequest 字数统计请求
type WordCountRequest struct {
	Content        string `json:"content" validate:"required"`
	FilterMarkdown bool   `json:"filterMarkdown"` // 是否过滤Markdown语法
}

// UpdateShortcutsRequest 更新快捷键请求
type UpdateShortcutsRequest struct {
	Shortcuts map[string]documentModel.Shortcut `json:"shortcuts" validate:"required"`
}

