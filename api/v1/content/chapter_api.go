package content

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	contentService "Qingyu_backend/service/interfaces/content"
)

// ChapterAPI 章节内容API
type ChapterAPI struct {
	chapterService contentService.ChapterServicePort
}

// NewChapterAPI 创建章节API实例
func NewChapterAPI(chapterService contentService.ChapterServicePort) *ChapterAPI {
	return &ChapterAPI{
		chapterService: chapterService,
	}
}

// GetChapter 获取章节内容
//
//	@Summary		获取章节内容
//	@Description	获取指定章节的完整内容
//	@Tags			内容管理-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/books/{bookId}/chapters/{chapterId} [get]
func (api *ChapterAPI) GetChapter(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	result, err := api.chapterService.GetChapter(c.Request.Context(), bookID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// ListChapters 获取章节列表
//
//	@Summary		获取章节列表
//	@Description	获取书籍的所有章节列表
//	@Tags			内容管理-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId	path		string	true	"书籍ID"
//	@Success		200	{object}	shared.APIResponse
//	@Failure		400	{object}	shared.APIResponse
//	@Failure		500	{object}	shared.APIResponse
//	@Router			/api/v1/content/books/{bookId}/chapters [get]
func (api *ChapterAPI) ListChapters(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	result, err := api.chapterService.ListChapters(c.Request.Context(), bookID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetNextChapter 获取下一章
//
//	@Summary		获取下一章
//	@Description	获取当前章节的下一章内容
//	@Tags			内容管理-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/books/{bookId}/chapters/{chapterId}/next [get]
func (api *ChapterAPI) GetNextChapter(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	result, err := api.chapterService.GetNextChapter(c.Request.Context(), bookID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetPreviousChapter 获取上一章
//
//	@Summary		获取上一章
//	@Description	获取当前章节的上一章内容
//	@Tags			内容管理-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/books/{bookId}/chapters/{chapterId}/previous [get]
func (api *ChapterAPI) GetPreviousChapter(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	result, err := api.chapterService.GetPreviousChapter(c.Request.Context(), bookID, chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetChapterByNumber 根据章节号获取章节
//
//	@Summary		根据章节号获取章节
//	@Description	通过章节序号获取章节内容
//	@Tags			内容管理-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterNum	path		int		true	"章节号"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/books/{bookId}/chapters/by-number/{chapterNum} [get]
func (api *ChapterAPI) GetChapterByNumber(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	chapterNumStr := c.Param("chapterNum")
	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的章节号格式")
		return
	}

	result, err := api.chapterService.GetChapterByNumber(c.Request.Context(), bookID, chapterNum)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}

// GetChapterInfo 获取章节信息
//
//	@Summary		获取章节信息
//	@Description	获取章节基本信息，不含正文内容
//	@Tags			内容管理-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/content/books/{bookId}/chapters/{chapterId}/info [get]
func (api *ChapterAPI) GetChapterInfo(c *gin.Context) {
	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.BadRequest(c, "参数错误", "章节ID不能为空")
		return
	}

	result, err := api.chapterService.GetChapterInfo(c.Request.Context(), chapterID)
	if err != nil {
		c.Error(err)
		return
	}

	shared.Success(c, 200, "获取成功", result)
}
