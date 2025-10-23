package bookstore

import (
	"Qingyu_backend/models/bookstore"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreService "Qingyu_backend/service/bookstore"
)

// BookDetailAPI 书籍详情API处理器
type BookDetailAPI struct {
	service bookstoreService.BookDetailService
}

// NewBookDetailAPI 创建书籍详情API实例
func NewBookDetailAPI(service bookstoreService.BookDetailService) *BookDetailAPI {
	return &BookDetailAPI{
		service: service,
	}
}

// GetBookDetail 获取书籍详情
//
//	@Summary		获取书籍详情
//	@Description	根据书籍ID获取详细信息
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success 200 {object} APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/books/{id} [get]
func (api *BookDetailAPI) GetBookDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	book, err := api.service.GetBookDetailByID(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    http.StatusNotFound,
				Message: "书籍不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取书籍详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    book,
	})
}

// GetBooksByTitle 根据标题搜索书籍
//
//	@Summary		根据标题搜索书籍
//	@Description	根据书籍标题进行模糊搜索
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			title	query		string	true	"书籍标题"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/search/title [get]
func (api *BookDetailAPI) GetBooksByTitle(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍标题不能为空",
		})
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

	books, total, err := api.service.GetBooksByTitle(c.Request.Context(), title, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "搜索书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    http.StatusOK,
		Message: "搜索成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// GetBooksByAuthor 根据作者搜索书籍
//
//	@Summary		根据作者搜索书籍
//	@Description	根据作者名称搜索书籍
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			author	query		string	true	"作者名称"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/search/author [get]
func (api *BookDetailAPI) GetBooksByAuthor(c *gin.Context) {
	author := c.Query("author")
	if author == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "作者名称不能为空",
		})
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

	books, total, err := api.service.GetBooksByAuthor(c.Request.Context(), author, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "搜索书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    http.StatusOK,
		Message: "搜索成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// GetBooksByCategory 根据分类获取书籍
//
//	@Summary		根据分类获取书籍
//	@Description	根据分类获取书籍列表
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			category	query		string	true	"分类名称"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			size		query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/books/category [get]
func (api *BookDetailAPI) GetBooksByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "分类名称不能为空",
		})
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

	books, total, err := api.service.GetBooksByCategory(c.Request.Context(), category, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取分类书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// GetBooksByStatus 根据状态获取书籍
//
//	@Summary		根据状态获取书籍
//	@Description	根据书籍状态获取书籍列表
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string	true	"书籍状态(serializing/completed/paused)"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/status [get]
func (api *BookDetailAPI) GetBooksByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍状态不能为空",
		})
		return
	}

	// 验证状态值
	validStatuses := []string{"serializing", "completed", "paused"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍状态",
		})
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

	books, total, err := api.service.GetBooksByStatus(c.Request.Context(), status, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// GetBooksByTags 根据标签获取书籍
//
//	@Summary		根据标签获取书籍
//	@Description	根据标签获取书籍列表
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			tags	query		string	true	"标签列表(逗号分隔)"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/tags [get]
func (api *BookDetailAPI) GetBooksByTags(c *gin.Context) {
	tagsStr := c.Query("tags")
	if tagsStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "标签不能为空",
		})
		return
	}

	tags := strings.Split(tagsStr, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	books, total, err := api.service.GetBooksByTags(c.Request.Context(), tags, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// SearchBooks 搜索书籍
//
//	@Summary		搜索书籍
//	@Description	根据关键词搜索书籍
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			keyword	query		string	true	"搜索关键词"
//	@Param			page	query		int		false	"页码"	default(1)
//	@Param			size	query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/search [get]
func (api *BookDetailAPI) SearchBooks(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "搜索关键词不能为空",
		})
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

	books, total, err := api.service.SearchBooks(c.Request.Context(), keyword, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "搜索书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    http.StatusOK,
		Message: "搜索成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// GetRecommendedBooks 获取推荐书籍
//
//	@Summary		获取推荐书籍
//	@Description	获取推荐书籍列表
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"数量限制"	default(10)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/recommended [get]
func (api *BookDetailAPI) GetRecommendedBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	books, err := api.service.GetRecommendedBooks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取推荐书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    books,
	})
}

// GetSimilarBooks 获取相似书籍
//
//	@Summary		获取相似书籍
//	@Description	根据书籍ID获取相似书籍
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"书籍ID"
//	@Param			limit	query		int		false	"数量限制"	default(10)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/{id}/similar [get]
func (api *BookDetailAPI) GetSimilarBooks(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	books, err := api.service.GetSimilarBooks(c.Request.Context(), id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取相似书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    books,
	})
}

// GetPopularBooks 获取热门书籍
//
//	@Summary		获取热门书籍
//	@Description	获取热门书籍列表
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"数量限制"	default(10)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/popular [get]
func (api *BookDetailAPI) GetPopularBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	books, err := api.service.GetPopularBooks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取热门书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    books,
	})
}

// GetLatestBooks 获取最新书籍
//
//	@Summary		获取最新书籍
//	@Description	获取最新发布的书籍列表
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"数量限制"	default(10)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/latest [get]
func (api *BookDetailAPI) GetLatestBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	books, err := api.service.GetLatestBooks(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取最新书籍失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    books,
	})
}

// GetBookStatistics 获取书籍统计信息
//
//	@Summary		获取书籍统计信息
//	@Description	获取书籍的浏览量、收藏量等统计信息
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/books/{id}/statistics [get]
func (api *BookDetailAPI) GetBookStatistics(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	// 获取统计信息
	totalCount, err := api.service.CountBooksByCategory(c.Request.Context(), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取统计信息失败",
		})
		return
	}

	// 这里可以扩展更多统计信息
	statistics := map[string]interface{}{
		"total_books": totalCount,
		"book_id":     id.Hex(),
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    statistics,
	})
}

// CreateBookDetail 创建书籍详情
//
//	@Summary		创建书籍详情
//	@Description	创建新的书籍详情信息
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			book	body		bookstore.BookDetail	true	"书籍详情信息"
//	@Success 201 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books [post]
func (api *BookDetailAPI) CreateBookDetail(c *gin.Context) {
	var book bookstore.BookDetail
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数格式错误",
		})
		return
	}

	if err := api.service.CreateBookDetail(c.Request.Context(), &book); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "创建书籍详情失败",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    http.StatusCreated,
		Message: "创建成功",
		Data:    book,
	})
}

// UpdateBookDetail 更新书籍详情
//
//	@Summary		更新书籍详情
//	@Description	更新书籍详情信息
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"书籍ID"
//	@Param			book	body		bookstore.BookDetail	true	"书籍详情信息"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		404		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/books/{id} [put]
func (api *BookDetailAPI) UpdateBookDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	var book bookstore.BookDetail
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数格式错误",
		})
		return
	}

	book.ID = id

	if err := api.service.UpdateBookDetail(c.Request.Context(), &book); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    http.StatusNotFound,
				Message: "书籍不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "更新书籍详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    book,
	})
}

// DeleteBookDetail 删除书籍详情
//
//	@Summary		删除书籍详情
//	@Description	删除书籍详情信息
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/books/{id} [delete]
func (api *BookDetailAPI) DeleteBookDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	if err := api.service.DeleteBookDetail(c.Request.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    http.StatusNotFound,
				Message: "书籍不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "删除书籍详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "书籍详情删除成功",
	})
}

// IncrementViewCount 增加书籍浏览量
//
//	@Summary		增加书籍浏览量
//	@Description	记录用户浏览书籍详情，增加浏览量统计
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/books/{id}/view [post]
func (api *BookDetailAPI) IncrementViewCount(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	err = api.service.IncrementViewCount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "增加浏览量失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "浏览量增加成功",
	})
}

// IncrementLikeCount 增加书籍点赞数
//
//	@Summary		增加书籍点赞数
//	@Description	用户点赞书籍，增加点赞数统计
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/books/{id}/like [post]
func (api *BookDetailAPI) IncrementLikeCount(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	err = api.service.IncrementLikeCount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "增加点赞数失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "点赞成功",
	})
}

// DecrementLikeCount 减少书籍点赞数
//
//	@Summary		减少书籍点赞数
//	@Description	用户取消点赞书籍，减少点赞数统计
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/books/{id}/unlike [post]
func (api *BookDetailAPI) DecrementLikeCount(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    http.StatusBadRequest,
			Message: "无效的书籍ID格式",
		})
		return
	}

	err = api.service.DecrementLikeCount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    http.StatusInternalServerError,
			Message: "取消点赞失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Message: "取消点赞成功",
	})
}
