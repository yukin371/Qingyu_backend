package bookstore

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/models/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// ChapterAPI 章节API处理器
type ChapterAPI struct {
	service bookstoreService.ChapterService
}

// NewChapterAPI 创建章节API实例
func NewChapterAPI(service bookstoreService.ChapterService) *ChapterAPI {
	return &ChapterAPI{
		service: service,
	}
}

// GetChapter 获取章节详情
//
//	@Summary		获取章节详情
//	@Description	根据章节ID获取章节详细信息
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"章节ID"
//	@Success 200 {object} response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters/{id} [get]
func (api *ChapterAPI) GetChapter(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	chapter, err := api.service.GetChapterByID(c.Request.Context(), id.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", chapter)
}

// GetChaptersByBookID 根据书籍ID获取章节列表
//
//	@Summary		根据书籍ID获取章节列表
//	@Description	根据书籍ID获取所有章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters [get]
func (api *ChapterAPI) GetChaptersByBookID(c *gin.Context) {
	bookIDStr := c.Param("id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	chapters, total, err := api.service.GetChaptersByBookID(c.Request.Context(), bookID.Hex(), page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Paginated(c, chapters, total, page, size, "获取成功")
}

// GetChapterByBookIDAndNumber 根据书籍ID和章节号获取章节
//
//	@Summary		根据书籍ID和章节号获取章节
//	@Description	根据书籍ID和章节号获取特定章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			book_id		path		string	true	"书籍ID"
//	@Param			chapter_num	path		int		true	"章节号"
//	@Success 200 {object} response.APIResponse
//	@Failure		400			{object}	response.APIResponse
//	@Failure		404			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/{chapter_num} [get]
func (api *ChapterAPI) GetChapterByBookIDAndNumber(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	chapterNumStr := c.Param("chapter_num")
	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的章节号格式")
		return
	}

	chapter, err := api.service.GetChapterByBookIDAndNum(c.Request.Context(), bookID.Hex(), chapterNum)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", chapter)
}

// GetFreeChapters 获取免费章节
//
//	@Summary		获取免费章节
//	@Description	根据书籍ID获取免费章节列表
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/free [get]
func (api *ChapterAPI) GetFreeChapters(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	chapters, total, err := api.service.GetFreeChaptersByBookID(c.Request.Context(), bookID.Hex(), page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Paginated(c, chapters, total, page, size, "获取成功")
}

// GetPaidChapters 获取付费章节
//
//	@Summary		获取付费章节
//	@Description	根据书籍ID获取付费章节列表
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/paid [get]
func (api *ChapterAPI) GetPaidChapters(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	chapters, total, err := api.service.GetPaidChaptersByBookID(c.Request.Context(), bookID.Hex(), page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Paginated(c, chapters, total, page, size, "获取成功")
}

// GetPublishedChapters 获取已发布章节
//
//	@Summary		获取已发布章节
//	@Description	根据书籍ID获取已发布章节列表
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/published [get]
func (api *ChapterAPI) GetPublishedChapters(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	chapters, total, err := api.service.GetPublishedChaptersByBookID(c.Request.Context(), bookID.Hex(), page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Paginated(c, chapters, total, page, size, "获取成功")
}

// GetPreviousChapter 获取上一章节
//
//	@Summary		获取上一章节
//	@Description	根据当前章节ID获取上一章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"当前章节ID"
//	@Success 200 {object} response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters/{id}/previous [get]
func (api *ChapterAPI) GetPreviousChapter(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	// 先获取当前章节信息
	currentChapter, err := api.service.GetChapterByID(c.Request.Context(), id.Hex())
	if err != nil {
		response.NotFound(c, "当前章节不存在")
		return
	}

	chapter, err := api.service.GetPreviousChapter(c.Request.Context(), currentChapter.BookID, currentChapter.ChapterNum)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "上一章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", chapter)
}

// GetNextChapter 获取下一章节
//
//	@Summary		获取下一章节
//	@Description	根据当前章节ID获取下一章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"当前章节ID"
//	@Success 200 {object} response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters/{id}/next [get]
func (api *ChapterAPI) GetNextChapter(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	// 先获取当前章节信息
	currentChapter, err := api.service.GetChapterByID(c.Request.Context(), id.Hex())
	if err != nil {
		response.NotFound(c, "当前章节不存在")
		return
	}

	chapter, err := api.service.GetNextChapter(c.Request.Context(), currentChapter.BookID, currentChapter.ChapterNum)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "下一章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", chapter)
}

// GetFirstChapter 获取第一章节
//
//	@Summary		获取第一章节
//	@Description	根据书籍ID获取第一章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/first [get]
func (api *ChapterAPI) GetFirstChapter(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	chapter, err := api.service.GetFirstChapter(c.Request.Context(), bookID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "第一章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", chapter)
}

// GetLastChapter 获取最后章节
//
//	@Summary		获取最后章节
//	@Description	根据书籍ID获取最后章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/last [get]
func (api *ChapterAPI) GetLastChapter(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	chapter, err := api.service.GetLastChapter(c.Request.Context(), bookID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "最后章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", chapter)
}

// GetChapterContent 获取章节内容
//
//	@Summary		获取章节内容
//	@Description	根据章节ID获取章节内容
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"章节ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters/{id}/content [get]
func (api *ChapterAPI) GetChapterContent(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	// 获取用户ID（从中间件设置的上下文中）
	var userID primitive.ObjectID
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(string); ok {
			userID, _ = primitive.ObjectIDFromHex(uid)
		}
	}

	content, err := api.service.GetChapterContent(c.Request.Context(), id.Hex(), userID.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "章节内容不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", map[string]interface{}{
		"chapter_id": id.Hex(),
		"content":    content,
	})
}

// SearchChapters 搜索章节
//
//	@Summary		搜索章节
//	@Description	根据关键词搜索章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query		string	true	"搜索关键词"
//	@Param			book_id	query		string	false	"书籍ID(可选)"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters/search [get]
func (api *ChapterAPI) SearchChapters(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "参数错误", "搜索关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	chapters, total, err := api.service.SearchChapters(c.Request.Context(), keyword, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Paginated(c, chapters, total, page, size, "搜索成功")
}

// GetChapterStatistics 获取章节统计信息
//
//	@Summary		获取章节统计信息
//	@Description	获取书籍的章节统计信息
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/chapters/statistics [get]
func (api *ChapterAPI) GetChapterStatistics(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	// 获取统计信息
	totalCount, err := api.service.GetChapterCountByBookID(c.Request.Context(), bookID.Hex())
	if err != nil {
		response.InternalError(c, err)
		return
	}

	freeCount, err := api.service.GetFreeChapterCountByBookID(c.Request.Context(), bookID.Hex())
	if err != nil {
		response.InternalError(c, err)
		return
	}

	paidCount, err := api.service.GetPaidChapterCountByBookID(c.Request.Context(), bookID.Hex())
	if err != nil {
		response.InternalError(c, err)
		return
	}

	totalWordCount, err := api.service.GetTotalWordCountByBookID(c.Request.Context(), bookID.Hex())
	if err != nil {
		response.InternalError(c, err)
		return
	}

	statistics := map[string]interface{}{
		"book_id":          bookID.Hex(),
		"total_chapters":   totalCount,
		"free_chapters":    freeCount,
		"paid_chapters":    paidCount,
		"total_word_count": totalWordCount,
	}

	response.SuccessWithMessage(c, "获取成功", statistics)
}

// CreateChapter 创建章节
//
//	@Summary		创建章节
//	@Description	创建新的章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			chapter	body		bookstore.Chapter	true	"章节信息"
//	@Success 201 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters [post]
func (api *ChapterAPI) CreateChapter(c *gin.Context) {
	var chapter bookstore.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		response.BadRequest(c, "参数错误", "请求参数格式错误")
		return
	}

	if err := api.service.CreateChapter(c.Request.Context(), &chapter); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, chapter)
}

// UpdateChapter 更新章节
//
//	@Summary		更新章节
//	@Description	更新章节信息
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"章节ID"
//	@Param			chapter	body		bookstore.Chapter	true	"章节信息"
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters/{id} [put]
func (api *ChapterAPI) UpdateChapter(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	var chapter bookstore.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		response.BadRequest(c, "参数错误", "请求参数格式错误")
		return
	}

	chapter.ID = id.Hex()

	if err := api.service.UpdateChapter(c.Request.Context(), &chapter); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", chapter)
}

// DeleteChapter 删除章节
//
//	@Summary		删除章节
//	@Description	删除章节
//	@Tags			章节
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"章节ID"
//	@Success		200	{object}	response.APIResponse
//	@Failure		400	{object}	response.APIResponse
//	@Failure		404	{object}	response.APIResponse
//	@Failure		500	{object}	response.APIResponse
//	@Router			/api/v1/bookstore/chapters/{id} [delete]
func (api *ChapterAPI) DeleteChapter(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "参数错误", "无效的章节ID格式")
		return
	}

	if err := api.service.DeleteChapter(c.Request.Context(), id.Hex()); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
