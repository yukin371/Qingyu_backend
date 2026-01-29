package reader

import (
	readerModels "Qingyu_backend/models/reader"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	"errors"
)

// BatchCreateAnnotationsRequest 批量创建注记请求
type BatchCreateAnnotationsRequest struct {
	Annotations []CreateAnnotationRequest `json:"annotations" binding:"required,min=1,max=50"`
}

// AnnotationUpdate 单个注记更新
type AnnotationUpdate struct {
	ID      string                  `json:"id" binding:"required"`
	Updates UpdateAnnotationRequest `json:"updates"`
}

// BatchUpdateAnnotationsRequest 批量更新注记请求
type BatchUpdateAnnotationsRequest struct {
	Updates []AnnotationUpdate `json:"updates" binding:"required,min=1,max=50"`
}

// BatchDeleteAnnotationsRequest 批量删除注记请求
type BatchDeleteAnnotationsRequest struct {
	IDs []string `json:"ids" binding:"required,min=1,max=100"`
}

// GetAnnotationStatsResponse 注记统计响应
type GetAnnotationStatsResponse struct {
	TotalCount     int `json:"totalCount"`
	BookmarkCount  int `json:"bookmarkCount"`
	HighlightCount int `json:"highlightCount"`
	NoteCount      int `json:"noteCount"`
}

// BatchCreateAnnotations 批量创建注记
//
//	@Summary	批量创建注记
//	@Tags		阅读器
//	@Param		request	body		BatchCreateAnnotationsRequest	true	"批量创建注记请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/batch [post]
func (api *AnnotationsAPI) BatchCreateAnnotations(c *gin.Context) {
	var req BatchCreateAnnotationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("用户ID类型错误: "))
		return
	}

	// 转换为注记模型
	annotations := make([]*readerModels.Annotation, len(req.Annotations))
	for i, annReq := range req.Annotations {
		userOID, _ := primitive.ObjectIDFromHex(userIDStr)
		bookOID, _ := primitive.ObjectIDFromHex(annReq.BookID)
		chapterOID, _ := primitive.ObjectIDFromHex(annReq.ChapterID)
		
		annotations[i] = &readerModels.Annotation{
			UserID:    userOID,
			BookID:    bookOID,
			ChapterID: chapterOID,
			Type:      annReq.Type,
			Text:      annReq.Text,
			Note:      annReq.Note,
			Range:     annReq.Range,
		}
	}

	// 批量创建
	err := api.readerService.BatchCreateAnnotations(c.Request.Context(), annotations)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusCreated, "批量创建成功", gin.H{
		"count":       len(annotations),
		"annotations": annotations,
	})
}

// BatchUpdateAnnotations 批量更新注记
//
//	@Summary	批量更新注记
//	@Tags		阅读器
//	@Param		request	body		BatchUpdateAnnotationsRequest	true	"批量更新注记请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/batch [put]
func (api *AnnotationsAPI) BatchUpdateAnnotations(c *gin.Context) {
	var req BatchUpdateAnnotationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 批量更新
	successCount := 0
	for _, update := range req.Updates {
		updates := make(map[string]interface{})
		if update.Updates.Text != nil {
			updates["text"] = *update.Updates.Text
		}
		if update.Updates.Note != nil {
			updates["note"] = *update.Updates.Note
		}
		if update.Updates.Range != nil {
			updates["range"] = *update.Updates.Range
		}

		err := api.readerService.UpdateAnnotation(c.Request.Context(), update.ID, updates)
		if err == nil {
			successCount++
		}
	}

	shared.Success(c, http.StatusOK, "批量更新完成", gin.H{
		"total":   len(req.Updates),
		"success": successCount,
	})
}

// BatchDeleteAnnotations 批量删除注记
//
//	@Summary	批量删除注记
//	@Tags		阅读器
//	@Param		request	body		BatchDeleteAnnotationsRequest	true	"批量删除注记请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/batch [delete]
func (api *AnnotationsAPI) BatchDeleteAnnotations(c *gin.Context) {
	var req BatchDeleteAnnotationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 批量删除
	err := api.readerService.BatchDeleteAnnotations(c.Request.Context(), req.IDs)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "批量删除成功", gin.H{
		"count": len(req.IDs),
	})
}

// GetAnnotationStats 获取注记统计
//
//	@Summary	获取注记统计
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/stats [get]
func (api *AnnotationsAPI) GetAnnotationStats(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("用户ID类型错误: "))
		return
	}

	bookID := c.Query("bookId")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID不能为空")
		return
	}

	// 获取统计数据
	stats, err := api.readerService.GetAnnotationStats(c.Request.Context(), userIDStr, bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", stats)
}

// ExportAnnotations 导出注记
//
//	@Summary	导出注记
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Param		format	query		string	false	"导出格式 (json/markdown/txt)"	default(json)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/export [get]
func (api *AnnotationsAPI) ExportAnnotations(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("用户ID类型错误: "))
		return
	}

	bookID := c.Query("bookId")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID不能为空")
		return
	}

	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "markdown" && format != "txt" {
		response.BadRequest(c,  "参数错误", "不支持的导出格式")
		return
	}

	// 获取注记数据
	annotations, err := api.readerService.GetAnnotationsByBook(c.Request.Context(), userIDStr, bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// 根据格式导出
	var content string
	var contentType string
	var filename string

	switch format {
	case "json":
		content = api.exportAsJSON(annotations)
		contentType = "application/json"
		filename = "annotations.json"
	case "markdown":
		content = api.exportAsMarkdown(annotations)
		contentType = "text/markdown"
		filename = "annotations.md"
	case "txt":
		content = api.exportAsText(annotations)
		contentType = "text/plain"
		filename = "annotations.txt"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.String(http.StatusOK, content)
}

// 导出为JSON格式
func (api *AnnotationsAPI) exportAsJSON(_ []*readerModels.Annotation) string {
	// 简单实现，实际可以使用json.Marshal
	return `{"annotations": []}`
}

// 导出为Markdown格式
func (api *AnnotationsAPI) exportAsMarkdown(annotations []*readerModels.Annotation) string {
	content := "# 我的笔记\n\n"
	for _, ann := range annotations {
		content += "## " + ann.Type + "\n\n"
		if ann.Text != "" {
			content += "> " + ann.Text + "\n\n"
		}
		if ann.Note != "" {
			content += ann.Note + "\n\n"
		}
		content += "---\n\n"
	}
	return content
}

// 导出为文本格式
func (api *AnnotationsAPI) exportAsText(annotations []*readerModels.Annotation) string {
	content := "我的笔记\n\n"
	for _, ann := range annotations {
		content += ann.Type + "\n"
		if ann.Text != "" {
			content += "引用: " + ann.Text + "\n"
		}
		if ann.Note != "" {
			content += "笔记: " + ann.Note + "\n"
		}
		content += "\n"
	}
	return content
}

// SyncAnnotations 同步注记（多端同步）
//
//	@Summary	同步注记
//	@Tags		阅读器
//	@Param		request	body		object	true	"同步注记请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/sync [post]
func (api *AnnotationsAPI) SyncAnnotations(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		response.InternalError(c, errors.New("用户ID类型错误: "))
		return
	}

	var req SyncAnnotationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 执行同步逻辑
	result, err := api.readerService.SyncAnnotations(c.Request.Context(), userIDStr, &req)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "同步成功", result)
}

// SyncAnnotationsRequest 同步注记请求
type SyncAnnotationsRequest struct {
	BookID           string                     `json:"bookId" binding:"required"`
	LastSyncTime     int64                      `json:"lastSyncTime"` // Unix时间戳
	LocalAnnotations []*readerModels.Annotation `json:"localAnnotations"`
}
