package reader

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readerModels "Qingyu_backend/models/reader"
	"Qingyu_backend/service/interfaces"
	"Qingyu_backend/pkg/response"
)

// AnnotationsAPI 标注API
type AnnotationsAPI struct {
	readerService interfaces.ReaderService
}

// NewAnnotationsAPI 创建标注API实例
func NewAnnotationsAPI(readerService interfaces.ReaderService) *AnnotationsAPI {
	return &AnnotationsAPI{
		readerService: readerService,
	}
}

// CreateAnnotationRequest 创建标注请求
type CreateAnnotationRequest struct {
	BookID    string `json:"bookId" binding:"required"`
	ChapterID string `json:"chapterId" binding:"required"`
	Type      string `json:"type" binding:"required"` // bookmark(书签) | highlight(高亮) | note(笔记)
	Text      string `json:"text"`                    // 标注文本
	Note      string `json:"note"`                    // 注释内容
	Range     string `json:"range"`                   // 标注范围：start-end
}

// UpdateAnnotationRequest 更新标注请求
type UpdateAnnotationRequest struct {
	Text  *string `json:"text"`  // 标注文本
	Note  *string `json:"note"`  // 注释内容
	Range *string `json:"range"` // 标注范围
}

// getUserID 从gin.Context中获取并验证用户ID
func (api *AnnotationsAPI) getUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return "", false
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.Error(errors.New("用户ID类型错误"))
		return "", false
	}

	return userIDStr, true
}

// requireQueryParam 验证查询参数是否存在
func (api *AnnotationsAPI) requireQueryParam(c *gin.Context, key string) (string, bool) {
	value := c.Query(key)
	if value == "" {
		response.BadRequest(c, "参数错误", key+"不能为空")
		return "", false
	}
	return value, true
}

// CreateAnnotation 创建标注
//
//	@Summary	创建标注
//	@Tags		阅读器
//	@Param		request	body		CreateAnnotationRequest	true	"创建标注请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations [post]
func (api *AnnotationsAPI) CreateAnnotation(c *gin.Context) {
	var req CreateAnnotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.Error(errors.New("用户ID类型错误"))
		return
	}

	userOID, _ := primitive.ObjectIDFromHex(userIDStr)
	bookOID, _ := primitive.ObjectIDFromHex(req.BookID)
	chapterOID, _ := primitive.ObjectIDFromHex(req.ChapterID)
	
	annotation := &readerModels.Annotation{
		UserID:    userOID,
		BookID:    bookOID,
		ChapterID: chapterOID,
		Type:      req.Type,
		Text:      req.Text,
		Note:      req.Note,
		Range:     req.Range,
	}

	err := api.readerService.CreateAnnotation(c.Request.Context(), annotation)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, annotation)
}

// UpdateAnnotation 更新标注
//
//	@Summary	更新标注
//	@Tags		阅读器
//	@Param		id		path		string					true	"标注ID"
//	@Param		request	body		UpdateAnnotationRequest	true	"更新标注请求"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/{id} [put]
func (api *AnnotationsAPI) UpdateAnnotation(c *gin.Context) {
	annotationID := c.Param("id")

	var req UpdateAnnotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Text != nil {
		updates["text"] = *req.Text
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.Range != nil {
		updates["range"] = *req.Range
	}

	err := api.readerService.UpdateAnnotation(c.Request.Context(), annotationID, updates)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// DeleteAnnotation 删除标注
//
//	@Summary	删除标注
//	@Tags		阅读器
//	@Param		id	path		string	true	"标注ID"
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/{id} [delete]
func (api *AnnotationsAPI) DeleteAnnotation(c *gin.Context) {
	annotationID := c.Param("id")

	err := api.readerService.DeleteAnnotation(c.Request.Context(), annotationID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// GetAnnotationsByChapter 获取章节标注
//
//	@Summary	获取章节标注
//	@Tags		阅读器
//	@Param		bookId		query		string	true	"书籍ID"
//	@Param		chapterId	query		string	true	"章节ID"
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/chapter [get]
func (api *AnnotationsAPI) GetAnnotationsByChapter(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	bookID, ok := api.requireQueryParam(c, "bookId")
	if !ok {
		return
	}

	chapterID, ok := api.requireQueryParam(c, "chapterId")
	if !ok {
		return
	}

	annotations, err := api.readerService.GetAnnotationsByChapter(c.Request.Context(), userIDStr, bookID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, annotations)
}

// GetAnnotationsByBook 获取书籍标注
//
//	@Summary	获取书籍标注
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/book [get]
func (api *AnnotationsAPI) GetAnnotationsByBook(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	bookID, ok := api.requireQueryParam(c, "bookId")
	if !ok {
		return
	}

	annotations, err := api.readerService.GetAnnotationsByBook(c.Request.Context(), userIDStr, bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, annotations)
}

// GetNotes 获取笔记
//
//	@Summary	获取笔记
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/notes [get]
func (api *AnnotationsAPI) GetNotes(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	bookID, ok := api.requireQueryParam(c, "bookId")
	if !ok {
		return
	}

	notes, err := api.readerService.GetNotes(c.Request.Context(), userIDStr, bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, notes)
}

// SearchNotes 搜索笔记
//
//	@Summary	搜索笔记
//	@Tags		阅读器
//	@Param		keyword	query		string	true	"搜索关键词"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/notes/search [get]
func (api *AnnotationsAPI) SearchNotes(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	keyword, ok := api.requireQueryParam(c, "keyword")
	if !ok {
		return
	}

	notes, err := api.readerService.SearchNotes(c.Request.Context(), userIDStr, keyword)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, notes)
}

// GetBookmarks 获取书签
//
//	@Summary	获取书签
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/bookmarks [get]
func (api *AnnotationsAPI) GetBookmarks(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	bookID, ok := api.requireQueryParam(c, "bookId")
	if !ok {
		return
	}

	bookmarks, err := api.readerService.GetBookmarks(c.Request.Context(), userIDStr, bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, bookmarks)
}

// GetLatestBookmark 获取最新书签
//
//	@Summary	获取最新书签
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/bookmarks/latest [get]
func (api *AnnotationsAPI) GetLatestBookmark(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	bookID, ok := api.requireQueryParam(c, "bookId")
	if !ok {
		return
	}

	bookmark, err := api.readerService.GetLatestBookmark(c.Request.Context(), userIDStr, bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, bookmark)
}

// GetHighlights 获取高亮
//
//	@Summary	获取高亮
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/highlights [get]
func (api *AnnotationsAPI) GetHighlights(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	bookID, ok := api.requireQueryParam(c, "bookId")
	if !ok {
		return
	}

	highlights, err := api.readerService.GetHighlights(c.Request.Context(), userIDStr, bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, highlights)
}

// GetRecentAnnotations 获取最近标注
//
//	@Summary	获取最近标注
//	@Tags		阅读器
//	@Param		limit	query		int	false	"数量限制"	default(20)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/recent [get]
func (api *AnnotationsAPI) GetRecentAnnotations(c *gin.Context) {
	userIDStr, ok := api.getUserID(c)
	if !ok {
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	annotations, err := api.readerService.GetRecentAnnotations(c.Request.Context(), userIDStr, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, annotations)
}

// GetPublicAnnotations 获取公开标注
//
//	@Summary	获取公开标注
//	@Tags		阅读器
//	@Param		bookId		query		string	true	"书籍ID"
//	@Param		chapterId	query		string	true	"章节ID"
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/annotations/public [get]
func (api *AnnotationsAPI) GetPublicAnnotations(c *gin.Context) {
	bookID, ok := api.requireQueryParam(c, "bookId")
	if !ok {
		return
	}

	chapterID, ok := api.requireQueryParam(c, "chapterId")
	if !ok {
		return
	}

	annotations, err := api.readerService.GetPublicAnnotations(c.Request.Context(), bookID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, annotations)
}
