package reader

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/interfaces"
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
	// 参数绑定（使用结构体+自动验证）
	var params GetChapterContentParams
	if !shared.BindParams(c, &params) {
		return
	}

	// 获取用户ID（可选）
	userID := shared.GetUserIDOptional(c)

	// 调用Service层（错误交给中间件处理）
	content, err := api.chapterService.GetChapterContent(c.Request.Context(), userID, params.BookID, params.ChapterID)
	if err != nil {
		c.Error(err)
		return
	}

	// 响应封装
	shared.Success(c, 200, "获取成功", content)
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
	// 参数绑定（使用结构体+自动验证）
	var params GetChapterByNumberParams
	if !shared.BindParams(c, &params) {
		return
	}

	// 获取用户ID（可选）
	userID := shared.GetUserIDOptional(c)

	// 调用Service层（错误交给中间件处理）
	content, err := api.chapterService.GetChapterByNumber(c.Request.Context(), userID, params.BookID, params.ChapterNum)
	if err != nil {
		c.Error(err)
		return
	}

	// 响应封装
	shared.Success(c, 200, "获取成功", content)
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
	// 参数绑定（使用结构体+自动验证）
	var params GetNextChapterParams
	if !shared.BindParams(c, &params) {
		return
	}

	// 获取用户ID（可选）
	userID := shared.GetUserIDOptional(c)

	// 调用Service层（错误交给中间件处理）
	nextChapter, err := api.chapterService.GetNextChapter(c.Request.Context(), userID, params.BookID, params.ChapterID)
	if err != nil {
		c.Error(err)
		return
	}

	// 响应封装
	shared.Success(c, 200, "获取成功", nextChapter)
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
	// 参数绑定（使用结构体+自动验证）
	var params GetPreviousChapterParams
	if !shared.BindParams(c, &params) {
		return
	}

	// 获取用户ID（可选）
	userID := shared.GetUserIDOptional(c)

	// 调用Service层（错误交给中间件处理）
	prevChapter, err := api.chapterService.GetPreviousChapter(c.Request.Context(), userID, params.BookID, params.ChapterID)
	if err != nil {
		c.Error(err)
		return
	}

	// 响应封装
	shared.Success(c, 200, "获取成功", prevChapter)
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
	// 参数绑定（使用结构体+自动验证）
	var params GetChapterListParams
	if !shared.BindParams(c, &params) {
		return
	}

	// 设置默认分页参数
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Size == 0 {
		params.Size = 50
	}

	// 获取用户ID（可选）
	userID := shared.GetUserIDOptional(c)

	// 调用Service层（错误交给中间件处理）
	chapterList, err := api.chapterService.GetChapterList(c.Request.Context(), userID, params.BookID, params.Page, params.Size)
	if err != nil {
		c.Error(err)
		return
	}

	// 响应封装
	shared.Success(c, 200, "获取成功", chapterList)
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
	// 参数绑定（使用结构体+自动验证）
	var params GetChapterInfoParams
	if !shared.BindParams(c, &params) {
		return
	}

	// 获取用户ID（可选）
	userID := shared.GetUserIDOptional(c)

	// 调用Service层（错误交给中间件处理）
	chapterInfo, err := api.chapterService.GetChapterInfo(c.Request.Context(), userID, params.ChapterID)
	if err != nil {
		c.Error(err)
		return
	}

	// 响应封装
	shared.Success(c, 200, "获取成功", chapterInfo)
}
