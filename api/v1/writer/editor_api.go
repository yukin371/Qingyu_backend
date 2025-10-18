package writer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	documentModel "Qingyu_backend/models/document"
	"Qingyu_backend/service/document"
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
// @Param request body document.AutoSaveRequest true "自动保存请求"
// @Success 200 {object} shared.Response{data=document.AutoSaveResponse}
// @Failure 409 {object} shared.Response "版本冲突"
// @Router /api/v1/documents/{id}/autosave [post]
func (api *EditorApi) AutoSaveDocument(c *gin.Context) {
	documentID := c.Param("id")

	var req document.AutoSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	req.DocumentID = documentID

	resp, err := api.documentService.AutoSaveDocument(c.Request.Context(), &req)
	if err != nil {
		// 检查是否是版本冲突
		if err.Error() == "版本冲突" {
			shared.Error(c, http.StatusConflict, "版本冲突", "文档已被其他用户修改，请刷新后重试")
			return
		}
		shared.Error(c, http.StatusInternalServerError, "自动保存失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "保存成功", resp)
}

// GetSaveStatus 获取保存状态
// @Summary 获取保存状态
// @Description 获取文档的保存状态和最后保存时间
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} shared.Response{data=document.SaveStatusResponse}
// @Router /api/v1/documents/{id}/save-status [get]
func (api *EditorApi) GetSaveStatus(c *gin.Context) {
	documentID := c.Param("id")

	status, err := api.documentService.GetSaveStatus(c.Request.Context(), documentID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "查询保存状态失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", status)
}

// GetDocumentContent 获取文档内容
// @Summary 获取文档内容
// @Description 获取文档的完整内容（用于编辑器加载）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Success 200 {object} shared.Response{data=document.DocumentContentResponse}
// @Router /api/v1/documents/{id}/content [get]
func (api *EditorApi) GetDocumentContent(c *gin.Context) {
	documentID := c.Param("id")

	content, err := api.documentService.GetDocumentContent(c.Request.Context(), documentID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取文档内容失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", content)
}

// UpdateDocumentContent 更新文档内容
// @Summary 更新文档内容
// @Description 手动更新文档内容（非自动保存）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body document.UpdateContentRequest true "更新内容请求"
// @Success 200 {object} shared.Response
// @Router /api/v1/documents/{id}/content [put]
func (api *EditorApi) UpdateDocumentContent(c *gin.Context) {
	documentID := c.Param("id")

	var req document.UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	req.DocumentID = documentID

	if err := api.documentService.UpdateDocumentContent(c.Request.Context(), &req); err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新内容失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// CalculateWordCount 计算字数
// @Summary 计算字数
// @Description 计算文档内容的字数统计（支持Markdown过滤）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param request body WordCountRequest true "字数统计请求"
// @Success 200 {object} shared.Response{data=document.WordCountResult}
// @Router /api/v1/documents/{id}/word-count [post]
func (api *EditorApi) CalculateWordCount(c *gin.Context) {
	var req WordCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	var result *document.WordCountResult
	if req.FilterMarkdown {
		result = api.wordCountService.CalculateWordCountWithMarkdown(req.Content)
	} else {
		result = api.wordCountService.CalculateWordCount(req.Content)
	}

	shared.Success(c, http.StatusOK, "计算成功", result)
}

// GetUserShortcuts 获取用户快捷键配置
// @Summary 获取用户快捷键配置
// @Description 获取当前用户的快捷键配置（包括自定义和默认）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Success 200 {object} shared.Response{data=document.ShortcutConfig}
// @Router /api/v1/user/shortcuts [get]
func (api *EditorApi) GetUserShortcuts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	config, err := api.shortcutService.GetUserShortcuts(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取快捷键配置失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", config)
}

// UpdateUserShortcuts 更新用户快捷键配置
// @Summary 更新用户快捷键配置
// @Description 更新用户的自定义快捷键配置
// @Tags 编辑器
// @Accept json
// @Produce json
// @Param request body UpdateShortcutsRequest true "快捷键配置"
// @Success 200 {object} shared.Response
// @Router /api/v1/user/shortcuts [put]
func (api *EditorApi) UpdateUserShortcuts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	var req UpdateShortcutsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if err := api.shortcutService.UpdateUserShortcuts(c.Request.Context(), userID.(string), req.Shortcuts); err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新快捷键配置失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// ResetUserShortcuts 重置用户快捷键配置
// @Summary 重置用户快捷键配置
// @Description 重置用户快捷键为默认配置
// @Tags 编辑器
// @Accept json
// @Produce json
// @Success 200 {object} shared.Response
// @Router /api/v1/user/shortcuts/reset [post]
func (api *EditorApi) ResetUserShortcuts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	if err := api.shortcutService.ResetUserShortcuts(c.Request.Context(), userID.(string)); err != nil {
		shared.Error(c, http.StatusInternalServerError, "重置快捷键配置失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "重置成功", nil)
}

// GetShortcutHelp 获取快捷键帮助
// @Summary 获取快捷键帮助
// @Description 获取快捷键帮助文档（按分类）
// @Tags 编辑器
// @Accept json
// @Produce json
// @Success 200 {object} shared.Response{data=[]document.ShortcutCategory}
// @Router /api/v1/user/shortcuts/help [get]
func (api *EditorApi) GetShortcutHelp(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	help, err := api.shortcutService.GetShortcutHelp(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取快捷键帮助失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", help)
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
