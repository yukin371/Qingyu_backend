package reader

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/reading"
)

// ChaptersAPI 章节API
type ChaptersAPI struct {
	readerService *reading.ReaderService
}

// NewChaptersAPI 创建章节API实例
func NewChaptersAPI(readerService *reading.ReaderService) *ChaptersAPI {
	return &ChaptersAPI{
		readerService: readerService,
	}
}

// GetChapterByID 获取章节信息
//
//	@Summary	获取章节信息
//	@Tags		阅读器
//	@Param		id	path		string	true	"章节ID"
//	@Success	200	{object}	response.Response
//	@Router		/api/v1/reader/chapters/{id} [get]
func (api *ChaptersAPI) GetChapterByID(c *gin.Context) {
	chapterID := c.Param("id")

	chapter, err := api.readerService.GetChapterByID(c.Request.Context(), chapterID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "章节不存在", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", chapter)
}

// GetChapterContent 获取章节内容
//
//	@Summary	获取章节内容
//	@Tags		阅读器
//	@Param		id	path		string	true	"章节ID"
//	@Success	200	{object}	response.Response
//	@Router		/api/v1/reader/chapters/{id}/content [get]
func (api *ChaptersAPI) GetChapterContent(c *gin.Context) {
	chapterID := c.Param("id")

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	content, err := api.readerService.GetChapterContent(c.Request.Context(), userID.(string), chapterID)
	if err != nil {
		shared.Error(c, http.StatusForbidden, "获取章节内容失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{"content": content})
}

// GetBookChapters 获取书籍章节列表
//
//	@Summary	获取书籍章节列表
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Param		page	query		int		false	"页码"	default(1)
//	@Param		size	query		int		false	"每页数量"	default(20)
//	@Success	200		{object}	response.Response
//	@Router		/api/v1/reader/chapters [get]
func (api *ChaptersAPI) GetBookChapters(c *gin.Context) {
	bookID := c.Query("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	chapters, total, err := api.readerService.GetBookChaptersWithPagination(c.Request.Context(), bookID, page, size)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取章节列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"chapters": chapters,
		"total":    total,
		"page":     page,
		"size":     size,
	})
}

// GetNavigationChapters 获取章节导航（上一章/下一章）
//
//	@Summary	获取章节导航
//	@Tags		阅读器
//	@Param		bookId		query		string	true	"书籍ID"
//	@Param		chapterNum	query		int		true	"当前章节号"
//	@Success	200			{object}	response.Response
//	@Router		/api/v1/reader/chapters/navigation [get]
func (api *ChaptersAPI) GetNavigationChapters(c *gin.Context) {
	bookID := c.Query("bookId")
	chapterNumStr := c.Query("chapterNum")

	if bookID == "" || chapterNumStr == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID和章节号不能为空")
		return
	}

	chapterNum, err := strconv.Atoi(chapterNumStr)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", "章节号格式错误")
		return
	}

	prevChapter, _ := api.readerService.GetPrevChapter(c.Request.Context(), bookID, chapterNum)
	nextChapter, _ := api.readerService.GetNextChapter(c.Request.Context(), bookID, chapterNum)

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"prevChapter": prevChapter,
		"nextChapter": nextChapter,
	})
}

// GetFirstChapter 获取第一章
//
//	@Summary	获取第一章
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.Response
//	@Router		/api/v1/reader/chapters/first [get]
func (api *ChaptersAPI) GetFirstChapter(c *gin.Context) {
	bookID := c.Query("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	chapter, err := api.readerService.GetFirstChapter(c.Request.Context(), bookID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "获取第一章失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", chapter)
}

// GetLastChapter 获取最后一章
//
//	@Summary	获取最后一章
//	@Tags		阅读器
//	@Param		bookId	query		string	true	"书籍ID"
//	@Success	200		{object}	response.Response
//	@Router		/api/v1/reader/chapters/last [get]
func (api *ChaptersAPI) GetLastChapter(c *gin.Context) {
	bookID := c.Query("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	chapter, err := api.readerService.GetLastChapter(c.Request.Context(), bookID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "获取最后一章失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", chapter)
}
