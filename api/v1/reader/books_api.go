package reader

import (
	"fmt"
	"net/http"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/interfaces"

	"github.com/gin-gonic/gin"
	"Qingyu_backend/pkg/response"
)

// BooksAPI 书架API
type BooksAPI struct {
	readerService interfaces.ReaderService
}

// NewBooksAPI 创建书架API实例
func NewBooksAPI(readerService interfaces.ReaderService) *BooksAPI {
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
		response.InternalError(c, err)
		return
	}

	// 转换为 DTO
	progressDTOs := ToReadingProgressDTOsFromPtrSlice(progresses)
	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"books": progressDTOs,
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
		response.BadRequest(c,  "参数错误", "书籍ID不能为空")
		return
	}

	// 通过保存初始进度来"添加到书架"
	err := api.readerService.SaveReadingProgress(c.Request.Context(), userID.(string), bookID, "", 0)
	if err != nil {
		response.InternalError(c, err)
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
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	bookID := c.Param("bookId")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID不能为空")
		return
	}

	// 调用Service删除阅读进度
	err := api.readerService.DeleteReadingProgress(c.Request.Context(), userID.(string), bookID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

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
		response.InternalError(c, err)
		return
	}

	// 转换为 DTO
	progressDTOs := ToReadingProgressDTOsFromPtrSlice(progresses)
	shared.Success(c, http.StatusOK, "获取成功", progressDTOs)
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
		response.InternalError(c, err)
		return
	}

	// 转换为 DTO
	progressDTOs := ToReadingProgressDTOsFromPtrSlice(progresses)
	shared.Success(c, http.StatusOK, "获取成功", progressDTOs)
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
		response.InternalError(c, err)
		return
	}

	// 转换为 DTO
	progressDTOs := ToReadingProgressDTOsFromPtrSlice(progresses)
	shared.Success(c, http.StatusOK, "获取成功", progressDTOs)
}

// UpdateBookStatusRequest 更新书籍状态请求
type UpdateBookStatusRequest struct {
	Status string `json:"status" binding:"required"` // reading(在读), want_read(想读), finished(读完)
}

// UpdateBookStatus 更新书籍状态
//
//	@Summary	更新书籍状态
//	@Tags		阅读器-书架
//	@Param		bookId	path	string	true	"书籍ID"
//	@Param		request	body	UpdateBookStatusRequest	true	"状态信息"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/books/{bookId}/status [put]
func (api *BooksAPI) UpdateBookStatus(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	bookID := c.Param("bookId")
	if bookID == "" {
		response.BadRequest(c,  "参数错误", "书籍ID不能为空")
		return
	}

	// 解析请求参数
	var req UpdateBookStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", "请求参数格式错误")
		return
	}

	// 调用服务层更新状态
	err := api.readerService.UpdateBookStatus(c.Request.Context(), userID.(string), bookID, req.Status)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// BatchUpdateBookStatusRequest 批量更新书籍状态请求
type BatchUpdateBookStatusRequest struct {
	BookIDs []string `json:"bookIds" binding:"required"` // 书籍ID列表
	Status  string   `json:"status" binding:"required"`  // reading(在读), want_read(想读), finished(读完)
}

// BatchUpdateBookStatus 批量更新书籍状态
//
//	@Summary	批量更新书籍状态
//	@Tags		阅读器-书架
//	@Param		request	body	BatchUpdateBookStatusRequest	true	"批量状态信息"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/books/batch/status [put]
func (api *BooksAPI) BatchUpdateBookStatus(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
		return
	}

	// 解析请求参数
	var req BatchUpdateBookStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c,  "参数错误", "请求参数格式错误")
		return
	}

	// 验证书籍ID列表
	if len(req.BookIDs) == 0 {
		response.BadRequest(c,  "参数错误", "书籍ID列表不能为空")
		return
	}

	if len(req.BookIDs) > 50 {
		response.BadRequest(c,  "参数错误", "批量更新数量不能超过50个")
		return
	}

	// 调用服务层批量更新状态
	err := api.readerService.BatchUpdateBookStatus(c.Request.Context(), userID.(string), req.BookIDs, req.Status)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	shared.Success(c, http.StatusOK, "批量更新成功", gin.H{
		"count": len(req.BookIDs),
	})
}

// parseInt 解析整数
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
