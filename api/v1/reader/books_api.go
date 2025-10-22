package reader

import (
	"fmt"
	"net/http"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/reading"

	"github.com/gin-gonic/gin"
)

// BooksAPI 书架API
type BooksAPI struct {
	readerService *reading.ReaderService
}

// NewBooksAPI 创建书架API实例
func NewBooksAPI(readerService *reading.ReaderService) *BooksAPI {
	return &BooksAPI{
		readerService: readerService,
	}
}

// GetBookshelf 获取书架（基于阅读进度）
//
//	@Summary	获取书架
//	@Tags		阅读器-书架
//	@Param		page	query	int	false	"页码"	default(1)
//	@Param		size	query	int	false	"每页数量"	default(20)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/books [get]
func (api *BooksAPI) GetBookshelf(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	// 获取分页参数
	page := 1
	size := 20
	if p, ok := c.GetQuery("page"); ok {
		if n, err := parseInt(p); err == nil && n > 0 {
			page = n
		}
	}
	if s, ok := c.GetQuery("size"); ok {
		if n, err := parseInt(s); err == nil && n > 0 && n <= 100 {
			size = n
		}
	}

	// 获取阅读历史（作为书架）
	progresses, total, err := api.readerService.GetReadingHistory(c.Request.Context(), userID.(string), page, size)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取书架失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"books": progresses,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// AddToBookshelf 添加到书架
//
//	@Summary	添加到书架
//	@Tags		阅读器-书架
//	@Param		bookId	path	string	true	"书籍ID"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/books/{bookId} [post]
func (api *BooksAPI) AddToBookshelf(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	bookID := c.Param("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	// 通过保存初始进度来"添加到书架"
	err := api.readerService.SaveReadingProgress(c.Request.Context(), userID.(string), bookID, "", 0)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "添加到书架失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "添加成功", nil)
}

// RemoveFromBookshelf 从书架移除
//
//	@Summary	从书架移除
//	@Tags		阅读器-书架
//	@Param		bookId	path	string	true	"书籍ID"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/books/{bookId} [delete]
func (api *BooksAPI) RemoveFromBookshelf(c *gin.Context) {
	// 获取用户ID
	_, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	bookID := c.Param("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	// TODO: 实现删除阅读进度的方法
	// 暂时返回成功，实际需要在Repository中添加Delete方法
	shared.Success(c, http.StatusOK, "移除成功", nil)
}

// GetRecentReading 获取最近阅读
//
//	@Summary	获取最近阅读
//	@Tags		阅读器-书架
//	@Param		limit	query	int	false	"数量限制"	default(10)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/books/recent [get]
func (api *BooksAPI) GetRecentReading(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	// 获取限制参数
	limit := 10
	if l, ok := c.GetQuery("limit"); ok {
		if n, err := parseInt(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	progresses, err := api.readerService.GetRecentReading(c.Request.Context(), userID.(string), limit)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取最近阅读失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", progresses)
}

// GetUnfinishedBooks 获取未读完的书
//
//	@Summary	获取未读完的书
//	@Tags		阅读器-书架
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/books/unfinished [get]
func (api *BooksAPI) GetUnfinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	progresses, err := api.readerService.GetUnfinishedBooks(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取未读完书籍失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", progresses)
}

// GetFinishedBooks 获取已读完的书
//
//	@Summary	获取已读完的书
//	@Tags		阅读器-书架
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/books/finished [get]
func (api *BooksAPI) GetFinishedBooks(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	progresses, err := api.readerService.GetFinishedBooks(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取已读完书籍失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", progresses)
}

// parseInt 解析整数
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
