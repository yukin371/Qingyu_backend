package reader

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/interfaces"
	readerservice "Qingyu_backend/service/reader"
	"Qingyu_backend/pkg/response"
)

// ChapterAPI 阅读器章节API
type ChapterAPI struct {
	chapterService interfaces.ReaderChapterService
}

// NewChapterAPI 创建章节API实例
func NewChapterAPI(chapterService interfaces.ReaderChapterService) *ChapterAPI {
	return &ChapterAPI{
		chapterService: chapterService,
	}
}

// GetChapterContent 获取章节内容
//
//	@Summary		获取章节内容（阅读器专用）
//	@Description	获取章节内容，自动检查权限、保存阅读进度
//	@Tags			阅读器-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	response.APIResponse
//	@Failure		400			{object}	response.APIResponse
//	@Failure		401			{object}	response.APIResponse
//	@Failure		403			{object}	response.APIResponse
//	@Failure		404			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/api/v1/reader/books/{bookId}/chapters/{chapterId} [get]
// @Summary 获取章节内容
// @Description TODO: 补充详细描述
// @Tags reader
// @Accept json
// @Produce json
// @Security Bearer
// @Param bookId path string true "BookId"
// @Param chapterId path string true "ChapterId"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /reader/books/{bookId}/chapters/{chapterId} [get]

func (api *ChapterAPI) GetChapterContent(c *gin.Context) {
	bookID := c.Param("bookId")
	chapterID := c.Param("chapterId")

	if bookID == "" || chapterID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID和章节ID不能为空")
		return
	}

	// 获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	content, err := api.chapterService.GetChapterContent(c.Request.Context(), userID, bookID, chapterID)
	if err != nil {
		if err == readerservice.ErrChapterNotFound {
			response.NotFound(c, "章节不存在")
			return
		}
		if err == readerservice.ErrChapterNotPublished {
			response.Forbidden(c, "章节未发布")
			return
		}
		if err == readerservice.ErrAccessDenied && content != nil {
			// 返回访问拒绝响应
			response.Forbidden(c, "无权访问")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, content)
}

// GetChapterByNumber 根据章节号获取内容
//
//	@Summary		根据章节号获取内容
//	@Description	根据章节号获取章节内容
//	@Tags			阅读器-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterNum	path		int		true	"章节号"
//	@Success		200			{object}	response.APIResponse
//	@Failure		400			{object}	response.APIResponse
//	@Failure		404			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/api/v1/reader/books/{bookId}/chapters/by-number/{chapterNum} [get]
func (api *ChapterAPI) GetChapterByNumber(c *gin.Context) {
	bookID := c.Param("bookId")
	chapterNumStr := c.Param("chapterNum")

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil || chapterNum <= 0 {
		response.BadRequest(c,  "参数错误", "无效的章节号")
		return
	}

	// 获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	content, err := api.chapterService.GetChapterByNumber(c.Request.Context(), userID, bookID, chapterNum)
	if err != nil {
		if err == readerservice.ErrChapterNotFound {
			response.NotFound(c, "章节不存在")
			return
		}
		if err == readerservice.ErrAccessDenied && content != nil {
			response.Forbidden(c, "无权访问")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, content)
}

// GetNextChapter 获取下一章
//
//	@Summary		获取下一章
//	@Description	获取当前章节的下一章信息
//	@Tags			阅读器-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"当前章节ID"
//	@Success		200			{object}	response.APIResponse
//	@Failure		400			{object}	response.APIResponse
//	@Failure		404			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/api/v1/reader/books/{bookId}/chapters/{chapterId}/next [get]
func (api *ChapterAPI) GetNextChapter(c *gin.Context) {
	bookID := c.Param("bookId")
	chapterID := c.Param("chapterId")

	if bookID == "" || chapterID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID和章节ID不能为空")
		return
	}

	// 获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	nextChapter, err := api.chapterService.GetNextChapter(c.Request.Context(), userID, bookID, chapterID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	if nextChapter == nil {
		response.Success(c, nil)
		return
	}

	response.Success(c, nextChapter)
}

// GetPreviousChapter 获取上一章
//
//	@Summary		获取上一章
//	@Description	获取当前章节的上一章信息
//	@Tags			阅读器-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId		path		string	true	"书籍ID"
//	@Param			chapterId	path		string	true	"当前章节ID"
//	@Success		200			{object}	response.APIResponse
//	@Failure		400			{object}	response.APIResponse
//	@Failure		404			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/api/v1/reader/books/{bookId}/chapters/{chapterId}/previous [get]
func (api *ChapterAPI) GetPreviousChapter(c *gin.Context) {
	bookID := c.Param("bookId")
	chapterID := c.Param("chapterId")

	if bookID == "" || chapterID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID和章节ID不能为空")
		return
	}

	// 获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	prevChapter, err := api.chapterService.GetPreviousChapter(c.Request.Context(), userID, bookID, chapterID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	if prevChapter == nil {
		response.Success(c, nil)
		return
	}

	response.Success(c, prevChapter)
}

// GetChapterList 获取章节目录
//
//	@Summary		获取章节目录
//	@Description	获取书籍的章节目录列表
//	@Tags			阅读器-章节
//	@Accept			json
//	@Produce		json
//	@Param			bookId	path		string	true	"书籍ID"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			size		query		int		false	"每页数量"	default(50)
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/reader/books/{bookId}/chapters [get]
func (api *ChapterAPI) GetChapterList(c *gin.Context) {
	bookID := c.Param("bookId")

	if bookID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 50
	}

	// 获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	chapterList, err := api.chapterService.GetChapterList(c.Request.Context(), userID, bookID, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, chapterList)
}

// GetChapterInfo 获取章节信息
//
//	@Summary		获取章节信息
//	@Description	获取章节信息（不含内容）
//	@Tags			阅读器-章节
//	@Accept			json
//	@Produce		json
//	@Param			chapterId	path		string	true	"章节ID"
//	@Success		200			{object}	response.APIResponse
//	@Failure		400			{object}	response.APIResponse
//	@Failure		404			{object}	response.APIResponse
//	@Failure		500			{object}	response.APIResponse
//	@Router			/api/v1/reader/chapters/{chapterId}/info [get]
func (api *ChapterAPI) GetChapterInfo(c *gin.Context) {
	chapterID := c.Param("chapterId")

	if chapterID == "" {
		response.BadRequest(c,  "参数错误", "章节ID不能为空")
		return
	}

	// 获取用户ID
	userID := ""
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	chapterInfo, err := api.chapterService.GetChapterInfo(c.Request.Context(), userID, chapterID)
	if err != nil {
		if err == readerservice.ErrChapterNotFound {
			response.NotFound(c, "章节不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, chapterInfo)
}
