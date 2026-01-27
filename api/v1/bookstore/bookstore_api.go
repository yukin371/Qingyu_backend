package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	searchModels "Qingyu_backend/models/search"
	"Qingyu_backend/models/shared/types"
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/logger"
	bookstoreService "Qingyu_backend/service/bookstore"
	"Qingyu_backend/service/search"
)

// BookstoreAPI 书城API处理器
type BookstoreAPI struct {
	service       bookstoreService.BookstoreService
	searchService *search.SearchService
	logger        *logger.Logger
}

// NewBookstoreAPI 创建书城API实例
func NewBookstoreAPI(service bookstoreService.BookstoreService, searchService *search.SearchService, logger *logger.Logger) *BookstoreAPI {
	return &BookstoreAPI{
		service:       service,
		searchService: searchService,
		logger:        logger,
	}
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse 分页响应格式
type PaginatedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	Size    int         `json:"size,omitempty"`
	Limit   int         `json:"limit,omitempty"` // 添加Limit字段
}

// GetHomepage 获取首页数据
//
//	@Summary		获取书城首页数据
//	@Description	获取书城首页的Banner、推荐书籍、精选书籍、分类等数据
//	@Tags			书城
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/homepage [get]
func (api *BookstoreAPI) GetHomepage(c *gin.Context) {
	data, err := api.service.GetHomepageData(c.Request.Context())
	if err != nil {
		shared.InternalError(c, "获取首页数据失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取首页数据成功", data)
}

// GetBooks 获取书籍列表
//
//	@Summary		获取书籍列表
//	@Description	获取所有书籍列表，支持分页
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"页码"	default(1)
//	@Param			size	query		int	false	"每页数量"	default(20)
//	@Success 200 {object} PaginatedResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/books [get]
func (api *BookstoreAPI) GetBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	books, total, err := api.service.GetAllBooks(c.Request.Context(), page, size)
	if err != nil {
		shared.InternalError(c, "获取书籍列表失败", err)
		return
	}

	// 转换为 DTO
	bookDTOs := ToBookDTOsFromPtrSlice(books)
	shared.Paginated(c, bookDTOs, total, page, size, "获取书籍列表成功")
}

// GetBookByID 根据ID获取书籍详情
//
//	@Summary		获取书籍详情
//	@Description	根据书籍ID获取书籍的详细信息
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success 200 {object} APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/books/{id} [get]
func (api *BookstoreAPI) GetBookByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	book, err := api.service.GetBookByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "book not found" || err.Error() == "book not available" {
			// 按需返回空数据而非404，便于前端容错显示
			shared.Success(c, http.StatusOK, "书籍不存在或不可用", nil)
			return
		}

		shared.InternalError(c, "获取书籍详情失败", err)
		return
	}

	// 转换为 DTO
	bookDTO := ToBookDTO(book)
	shared.Success(c, http.StatusOK, "获取书籍详情成功", bookDTO)
}

// GetBooksByCategory 根据分类获取书籍列表
//
//	@Summary		根据分类获取书籍列表
//	@Description	根据分类ID获取该分类下的书籍列表，支持分页
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			categoryId	path		string	true	"分类ID"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			size		query		int		false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/categories/{categoryId}/books [get]
func (api *BookstoreAPI) GetBooksByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		shared.BadRequest(c, "参数错误", "分类ID不能为空")
		return
	}

	// 验证ObjectID格式
	if _, err := primitive.ObjectIDFromHex(categoryID); err != nil {
		shared.BadRequest(c, "参数错误", "分类ID格式无效")
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

	books, total, err := api.service.GetBooksByCategory(c.Request.Context(), categoryID, page, size)
	if err != nil {
		shared.InternalError(c, "获取分类书籍失败", err)
		return
	}

	// 转换为 DTO
	bookDTOs := ToBookDTOsFromPtrSlice(books)
	shared.Paginated(c, bookDTOs, total, page, size, "获取分类书籍成功")
}

// GetRecommendedBooks 获取推荐书籍
//
//	@Summary		获取推荐书籍
//	@Description	获取推荐书籍列表，支持分页
//	@Tags			书籍推荐
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"页码"	default(1)
//	@Param			size	query		int	false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/books/recommended [get]
func (api *BookstoreAPI) GetRecommendedBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	books, total, err := api.service.GetRecommendedBooks(c.Request.Context(), page, size)
	if err != nil {
		shared.InternalError(c, "获取推荐书籍失败", err)
		return
	}

	// 转换为 DTO
	bookDTOs := ToBookDTOsFromPtrSlice(books)
	shared.Paginated(c, bookDTOs, total, page, size, "获取推荐书籍成功")
}

// GetFeaturedBooks 获取精选书籍
//
//	@Summary		获取精选书籍
//	@Description	获取精选书籍列表，支持分页
//	@Tags			书籍推荐
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"页码"	default(1)
//	@Param			size	query		int	false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/books/featured [get]
func (api *BookstoreAPI) GetFeaturedBooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	books, total, err := api.service.GetFeaturedBooks(c.Request.Context(), page, size)
	if err != nil {
		shared.InternalError(c, "获取精选书籍失败", err)
		return
	}

	// 转换为 DTO
	bookDTOs := ToBookDTOsFromPtrSlice(books)
	shared.Paginated(c, bookDTOs, total, page, size, "获取精选书籍成功")
}

// SearchBooks 搜索书籍
//
//	@Summary		搜索书籍
//	@Description	根据关键词和过滤条件搜索书籍
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			keyword		query		string		false	"搜索关键词"
//	@Param			categoryId	query		string		false	"分类ID"
//	@Param			author		query		string		false	"作者"
//	@Param			minRating	query		number		false	"最低评分"
//	@Param			tags		query		[]string	false	"标签"
//	@Param			sortBy		query		string		false	"排序字段"	Enums(created_at, updated_at, view_count, like_count, rating)
//	@Param			sortOrder	query		string		false	"排序方向"	Enums(asc, desc)
//	@Param			page		query		int			false	"页码"	default(1)
//	@Param			size		query		int			false	"每页数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/bookstore/books/search [get]
func (api *BookstoreAPI) SearchBooks(c *gin.Context) {
	startTime := time.Now()

	// 获取请求ID
	requestID := c.GetString("requestId")
	if requestID == "" {
		requestID = c.GetHeader("X-Request-ID")
	}

	// 获取用户ID
	var userID string
	if uid, exists := c.Get("userId"); exists {
		userID = uid.(string)
	}

	// 构建日志记录器
	searchLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	if userID != "" {
		searchLogger = searchLogger.WithUser(userID)
	}

	// 获取搜索参数
	keyword := c.Query("keyword")
	author := c.Query("author")
	categoryID := c.Query("categoryId")
	status := c.Query("status")
	tags := c.QueryArray("tags")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 记录搜索请求
	searchLogger.WithModule("search").Info("搜索请求",
		zap.String("keyword", keyword),
		zap.String("author", author),
		zap.String("category_id", categoryID),
		zap.String("status", status),
		zap.Strings("tags", tags),
		zap.String("sort_by", sortBy),
		zap.String("sort_order", sortOrder),
		zap.Int("page", page),
		zap.Int("page_size", size),
	)

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 优先尝试新路径（SearchService）
	if api.searchService != nil {
		// 检查是否有搜索关键词
		if keyword == "" && categoryID == "" && author == "" {
			searchLogger.WithModule("search").Warn("搜索参数不完整",
				zap.String("keyword", keyword),
				zap.Bool("has_category", categoryID != ""),
				zap.Bool("has_author", author != ""),
			)
			shared.BadRequest(c, "参数错误", "请提供搜索关键词或过滤条件")
			return
		}

		// 构建新的搜索请求
		newReq := &searchModels.SearchRequest{
			Type:     searchModels.SearchTypeBooks,
			Query:    keyword,
			Filter:   api.buildSearchFilter(categoryID, author, status, tags),
			Sort:     api.buildSearchSort(sortBy, sortOrder),
			Page:     page,
			PageSize: size,
		}

		// 如果 query 为空但有 filter，使用通配符查询
		if newReq.Query == "" && (newReq.Filter != nil || categoryID != "" || author != "") {
			newReq.Query = "*"
		}

		// 尝试新路径
		newResp, newErr := api.searchService.Search(c.Request.Context(), newReq)
		duration := time.Since(startTime)

		if newErr == nil && newResp != nil && newResp.Success && newResp.Data != nil {
			// 新路径成功
			searchLogger.WithModule("search").Info("搜索成功",
				zap.String("path", "new_search"),
				zap.Int64("total", newResp.Data.Total),
				zap.Int("returned", len(newResp.Data.Results)),
				zap.Duration("duration", duration),
				zap.Duration("took", newResp.Data.Took),
			)

			// 转换响应
			books := api.convertSearchResponseToBooks(newResp.Data.Results)
			bookDTOs := ToBookDTOsFromPtrSlice(books)

			responseData := map[string]interface{}{
				"books": bookDTOs,
				"total": newResp.Data.Total,
			}

			c.JSON(http.StatusOK, shared.APIResponse{
				Code:      http.StatusOK,
				Message:   "搜索书籍成功",
				Data:      responseData,
				Timestamp: 0,
			})
			return
		}

		// 新路径失败，记录日志并 fallback 到旧路径
		fallbackReason := "unknown"
		if newErr != nil {
			fallbackReason = newErr.Error()
		} else if newResp != nil && newResp.Error != nil {
			fallbackReason = newResp.Error.Message
		} else if newResp != nil && !newResp.Success {
			fallbackReason = "search failed"
		}

		searchLogger.WithModule("search").Warn("新路径失败，fallback 到旧路径",
			zap.String("path", "new_search"),
			zap.String("status", "fallback"),
			zap.String("fallback_reason", fallbackReason),
			zap.Duration("duration", duration),
		)
	}

	// 旧路径（原有实现）
	searchLogger.WithModule("search").Info("使用旧路径搜索",
		zap.String("path", "old_search"),
	)

	// 构建过滤器
	filter := &bookstore2.BookFilter{}

	if categoryID != "" {
		// 转换为ObjectID
		if _, err := primitive.ObjectIDFromHex(categoryID); err == nil {
			filter.CategoryID = &categoryID
		}
	}

	if author != "" {
		filter.Author = &author
	}

	// 处理status参数 - 前端使用"serializing"，后端使用"ongoing"
	if status != "" {
		// 映射前端状态值到后端状态值
		var backendStatus string
		switch status {
		case "serializing":
			backendStatus = "ongoing"
		case "completed", "paused":
			backendStatus = status
		default:
			// 其他状态值也尝试使用
			backendStatus = status
		}
		bookStatus := bookstore2.BookStatus(backendStatus)
		filter.Status = &bookStatus
	}

	if len(tags) > 0 {
		filter.Tags = tags
	}

	filter.SortBy = sortBy
	filter.SortOrder = sortOrder
	filter.Limit = size
	filter.Offset = (page - 1) * size

	if keyword == "" && filter.CategoryID == nil && filter.Author == nil {
		searchLogger.WithModule("search").Warn("搜索参数不完整",
			zap.String("keyword", keyword),
			zap.Bool("has_category", filter.CategoryID != nil),
			zap.Bool("has_author", filter.Author != nil),
		)
		shared.BadRequest(c, "参数错误", "请提供搜索关键词或过滤条件")
		return
	}

	// 设置关键词
	if keyword != "" {
		filter.Keyword = &keyword
	}

	// 执行搜索
	books, total, err := api.service.SearchBooksWithFilter(c.Request.Context(), filter)

	// 计算耗时
	duration := time.Since(startTime)

	if err != nil {
		searchLogger.WithModule("search").Error("搜索失败",
			zap.String("path", "old_search"),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		shared.InternalError(c, "搜索书籍失败", err)
		return
	}

	// 确保返回空数组而不是nil
	if books == nil {
		books = make([]*bookstore2.Book, 0)
	}

	// 记录搜索结果
	searchLogger.WithModule("search").Info("搜索成功",
		zap.String("path", "old_search"),
		zap.Int64("total", total),
		zap.Int("returned", len(books)),
		zap.Duration("duration", duration),
	)

	// 如果结果为空，记录警告
	if total == 0 {
		searchLogger.WithModule("search").Debug("搜索无结果",
			zap.String("keyword", keyword),
			zap.String("author", author),
		)
	}

	// 转换为 DTO
	bookDTOs := ToBookDTOsFromPtrSlice(books)

	// 构造响应数据
	responseData := map[string]interface{}{
		"books": bookDTOs,
		"total": total,
	}

	// 使用shared包的分页响应，带上额外数据
	c.JSON(http.StatusOK, shared.APIResponse{
		Code:      http.StatusOK,
		Message:   "搜索书籍成功",
		Data:      responseData,
		Timestamp: 0, // shared.Success会设置这个
	})
}

// SearchByTitle 按标题搜索书籍
//
//	@Summary     按标题搜索书籍
//	@Description 根据书籍标题进行模糊搜索，支持分页。优先使用SearchService (Milvus向量搜索)，失败或空结果时fallback到MongoDB
//	@Tags        书籍搜索
//	@Accept      json
//	@Produce     json
//	@Param       title query string true "标题关键词"
//	@Param       page query int false "页码" default(1)
//	@Param       size query int false "每页数量" default(20)
//	@Success     200 {object} shared.PaginatedResponse
//	@Failure     400 {object} shared.APIResponse
//	@Failure     500 {object} shared.APIResponse
//	@Router      /api/v1/bookstore/books/search/title [get]
func (api *BookstoreAPI) SearchByTitle(c *gin.Context) {
	// 构建日志记录器
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = primitive.NewObjectID().Hex()
	}

	searchLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	title := c.Query("title")
	if title == "" {
		shared.BadRequest(c, "参数错误", "标题不能为空")
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

	startTime := time.Now()

	// 记录搜索请求
	searchLogger.WithModule("search").Info("按标题搜索请求",
		zap.String("title", title),
		zap.Int("page", page),
		zap.Int("size", size),
	)

	// 优先使用新路径 (SearchService - Milvus向量搜索)
	if api.searchService != nil {
		req := &searchModels.SearchRequest{
			Type:  searchModels.SearchTypeBooks,
			Query: title,
			Filter: map[string]interface{}{
				"search_fields": []string{"title"},
			},
			Sort: []searchModels.SortField{
				{Field: "view_count", Ascending: false},
			},
			Page:     page,
			PageSize: size,
		}

		resp, err := api.searchService.Search(c.Request.Context(), req)
		duration := time.Since(startTime)

		// v1.2补充：完整的fallback触发条件
		shouldFallback := err != nil ||
			resp == nil ||
			!resp.Success ||
			resp.Data == nil ||
			resp.Data.Total == 0 // ⚠️ 空结果也触发fallback

		if !shouldFallback {
			// 新路径成功
			searchLogger.WithModule("search").Info("按标题搜索成功",
				zap.String("path", "search_by_title"),
				zap.Int64("total", resp.Data.Total),
				zap.Int("returned", len(resp.Data.Results)),
				zap.Duration("duration", duration),
				zap.Duration("took", resp.Data.Took),
			)

			books := api.convertSearchResponseToBooks(resp.Data.Results)
			bookDTOs := ToBookDTOsFromPtrSlice(books)

			responseData := map[string]interface{}{
				"books": bookDTOs,
				"total": resp.Data.Total,
			}

			c.JSON(http.StatusOK, shared.APIResponse{
				Code:      http.StatusOK,
				Message:   "搜索成功",
				Data:      responseData,
				Timestamp: 0,
			})
			return
		}

		// 记录fallback日志
		fallbackReason := "unknown"
		if err != nil {
			fallbackReason = err.Error()
		} else if resp != nil && resp.Error != nil {
			fallbackReason = resp.Error.Message
		} else if resp != nil && !resp.Success {
			fallbackReason = "search failed"
		} else if resp != nil && resp.Data != nil && resp.Data.Total == 0 {
			fallbackReason = "empty results"
		}

		searchLogger.WithModule("search").Warn("SearchService失败，fallback到MongoDB",
			zap.String("search_type", "title"),
			zap.String("fallback_reason", fallbackReason),
			zap.Duration("duration", duration),
		)
	}

	// Fallback 到旧路径 (MongoDB查询)
	searchLogger.WithModule("search").Info("使用MongoDB按标题搜索",
		zap.String("path", "mongodb_fallback"),
	)

	filter := &bookstore2.BookFilter{
		Keyword:   &title,
		SortBy:    "view_count",
		SortOrder: "desc",
		Limit:     size,
		Offset:    (page - 1) * size,
	}

	books, total, err := api.service.SearchBooksWithFilter(c.Request.Context(), filter)
	duration := time.Since(startTime)

	if err != nil {
		searchLogger.WithModule("search").Error("MongoDB搜索失败",
			zap.String("path", "mongodb_fallback"),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		shared.InternalError(c, "搜索失败", err)
		return
	}

	// 确保返回空数组而不是nil
	if books == nil {
		books = make([]*bookstore2.Book, 0)
	}

	searchLogger.WithModule("search").Info("MongoDB搜索成功",
		zap.String("path", "mongodb_fallback"),
		zap.Int64("total", total),
		zap.Int("returned", len(books)),
		zap.Duration("duration", duration),
	)

	bookDTOs := ToBookDTOsFromPtrSlice(books)
	responseData := map[string]interface{}{
		"books": bookDTOs,
		"total": total,
	}

	c.JSON(http.StatusOK, shared.APIResponse{
		Code:      http.StatusOK,
		Message:   "搜索成功",
		Data:      responseData,
		Timestamp: 0,
	})
}

// SearchByAuthor 按作者搜索书籍
//
//	@Summary     按作者搜索书籍
//	@Description 根据作者姓名进行模糊搜索，支持分页。优先使用SearchService (Milvus向量搜索)，失败或空结果时fallback到MongoDB
//	@Tags        书籍搜索
//	@Accept      json
//	@Produce     json
//	@Param       author query string true "作者姓名"
//	@Param       page query int false "页码" default(1)
//	@Param       size query int false "每页数量" default(20)
//	@Success     200 {object} shared.PaginatedResponse
//	@Failure     400 {object} shared.APIResponse
//	@Failure     500 {object} shared.APIResponse
//	@Router      /api/v1/bookstore/books/search/author [get]
func (api *BookstoreAPI) SearchByAuthor(c *gin.Context) {
	// 构建日志记录器
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = primitive.NewObjectID().Hex()
	}

	searchLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	author := c.Query("author")
	if author == "" {
		shared.BadRequest(c, "参数错误", "作者姓名不能为空")
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

	startTime := time.Now()

	// 记录搜索请求
	searchLogger.WithModule("search").Info("按作者搜索请求",
		zap.String("author", author),
		zap.Int("page", page),
		zap.Int("size", size),
	)

	// 优先使用新路径 (SearchService - Milvus向量搜索)
	if api.searchService != nil {
		req := &searchModels.SearchRequest{
			Type:  searchModels.SearchTypeBooks,
			Query: author,
			Filter: map[string]interface{}{
				"search_fields": []string{"author"},
			},
			Sort: []searchModels.SortField{
				{Field: "view_count", Ascending: false},
			},
			Page:     page,
			PageSize: size,
		}

		resp, err := api.searchService.Search(c.Request.Context(), req)
		duration := time.Since(startTime)

		// v1.2补充：完整的fallback触发条件
		shouldFallback := err != nil ||
			resp == nil ||
			!resp.Success ||
			resp.Data == nil ||
			resp.Data.Total == 0 // ⚠️ 空结果也触发fallback

		if !shouldFallback {
			// 新路径成功
			searchLogger.WithModule("search").Info("按作者搜索成功",
				zap.String("path", "search_by_author"),
				zap.Int64("total", resp.Data.Total),
				zap.Int("returned", len(resp.Data.Results)),
				zap.Duration("duration", duration),
				zap.Duration("took", resp.Data.Took),
			)

			books := api.convertSearchResponseToBooks(resp.Data.Results)
			bookDTOs := ToBookDTOsFromPtrSlice(books)

			responseData := map[string]interface{}{
				"books": bookDTOs,
				"total": resp.Data.Total,
			}

			c.JSON(http.StatusOK, shared.APIResponse{
				Code:      http.StatusOK,
				Message:   "搜索成功",
				Data:      responseData,
				Timestamp: 0,
			})
			return
		}

		// 记录fallback日志
		fallbackReason := "unknown"
		if err != nil {
			fallbackReason = err.Error()
		} else if resp != nil && resp.Error != nil {
			fallbackReason = resp.Error.Message
		} else if resp != nil && !resp.Success {
			fallbackReason = "search failed"
		} else if resp != nil && resp.Data != nil && resp.Data.Total == 0 {
			fallbackReason = "empty results"
		}

		searchLogger.WithModule("search").Warn("SearchService失败，fallback到MongoDB",
			zap.String("search_type", "author"),
			zap.String("fallback_reason", fallbackReason),
			zap.Duration("duration", duration),
		)
	}

	// Fallback 到旧路径 (MongoDB查询)
	searchLogger.WithModule("search").Info("使用MongoDB按作者搜索",
		zap.String("path", "mongodb_fallback"),
	)

	filter := &bookstore2.BookFilter{
		Author:    &author,
		SortBy:    "view_count",
		SortOrder: "desc",
		Limit:     size,
		Offset:    (page - 1) * size,
	}

	books, total, err := api.service.SearchBooksWithFilter(c.Request.Context(), filter)
	duration := time.Since(startTime)

	if err != nil {
		searchLogger.WithModule("search").Error("MongoDB搜索失败",
			zap.String("path", "mongodb_fallback"),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		shared.InternalError(c, "搜索失败", err)
		return
	}

	// 确保返回空数组而不是nil
	if books == nil {
		books = make([]*bookstore2.Book, 0)
	}

	searchLogger.WithModule("search").Info("MongoDB搜索成功",
		zap.String("path", "mongodb_fallback"),
		zap.Int64("total", total),
		zap.Int("returned", len(books)),
		zap.Duration("duration", duration),
	)

	bookDTOs := ToBookDTOsFromPtrSlice(books)
	responseData := map[string]interface{}{
		"books": bookDTOs,
		"total": total,
	}

	c.JSON(http.StatusOK, shared.APIResponse{
		Code:      http.StatusOK,
		Message:   "搜索成功",
		Data:      responseData,
		Timestamp: 0,
	})
}

// GetCategoryTree 获取分类树
//
//	@Summary		获取分类树
//	@Description	获取完整的分类树结构
//	@Tags			分类
//	@Accept			json
//	@Produce		json
//	@Success 200 {object} APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/categories/tree [get]
func (api *BookstoreAPI) GetCategoryTree(c *gin.Context) {
	tree, err := api.service.GetCategoryTree(c.Request.Context())
	if err != nil {
		shared.InternalError(c, "获取分类树失败", err)
		return
	}

	// 确保返回空数组而不是 nil
	if tree == nil {
		tree = []*bookstore2.CategoryTree{}
	}

	shared.Success(c, http.StatusOK, "获取分类树成功", tree)
}

// GetCategoryByID 根据ID获取分类详情
//
//	@Summary		获取分类详情
//	@Description	根据分类ID获取分类的详细信息
//	@Tags			分类
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"分类ID"
//	@Success 200 {object} APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/categories/{id} [get]
func (api *BookstoreAPI) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "分类ID不能为空")
		return
	}

	category, err := api.service.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "category not found" || err.Error() == "category not available" {
			shared.NotFound(c, "分类不存在或不可用")
			return
		}

		shared.InternalError(c, "获取分类详情失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取分类详情成功", category)
}

// GetActiveBanners 获取激活的Banner列表
//
//	@Summary		获取激活的Banner列表
//	@Description	获取当前激活的Banner列表
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"数量限制"	default(5)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/banners [get]
func (api *BookstoreAPI) GetActiveBanners(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit < 1 || limit > 20 {
		limit = 5
	}

	banners, err := api.service.GetActiveBanners(c.Request.Context(), limit)
	if err != nil {
		shared.InternalError(c, "获取Banner列表失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取Banner列表成功", banners)
}

// IncrementBookView 增加书籍浏览量
//
//	@Summary		增加书籍浏览量
//	@Description	记录用户浏览书籍，增加浏览量统计
//	@Tags			书籍
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/books/{id}/view [post]
func (api *BookstoreAPI) IncrementBookView(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		shared.BadRequest(c, "参数错误", "无效的书籍ID格式")
		return
	}

	err = api.service.IncrementBookView(c.Request.Context(), id.Hex())
	if err != nil {
		shared.InternalError(c, "增加浏览量失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "浏览量增加成功", nil)
}

// IncrementBannerClick 增加Banner点击次数
//
//	@Summary		增加Banner点击次数
//	@Description	记录用户点击Banner，增加点击次数统计
//	@Tags			Banner
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Banner ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/bookstore/banners/{id}/click [post]
func (api *BookstoreAPI) IncrementBannerClick(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "Banner ID不能为空")
		return
	}

	err := api.service.IncrementBannerClick(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "banner not found" || err.Error() == "banner not available" {
			shared.NotFound(c, "Banner不存在或不可用")
			return
		}

		shared.InternalError(c, "增加点击次数失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "点击次数增加成功", nil)
}

// GetRealtimeRanking 获取实时榜
//
//	@Summary		获取实时榜单
//	@Description	获取当日实时榜单数据
//	@Tags			书籍排行榜
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"限制数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/rankings/realtime [get]
func (api *BookstoreAPI) GetRealtimeRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rankings, err := api.service.GetRealtimeRanking(c.Request.Context(), limit)
	if err != nil {
		shared.InternalError(c, "获取实时榜失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取实时榜成功", rankings)
}

// GetWeeklyRanking 获取周榜
//
//	@Summary		获取周榜单
//	@Description	获取指定周期的周榜单数据
//	@Tags			书籍排行榜
//	@Accept			json
//	@Produce		json
//	@Param			period	query		string	false	"周期 (格式: 2024-W01)"
//	@Param			limit	query		int		false	"限制数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/rankings/weekly [get]
func (api *BookstoreAPI) GetWeeklyRanking(c *gin.Context) {
	period := c.Query("period")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rankings, err := api.service.GetWeeklyRanking(c.Request.Context(), period, limit)
	if err != nil {
		shared.InternalError(c, "获取周榜失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取周榜成功", rankings)
}

// GetMonthlyRanking 获取月榜
//
//	@Summary		获取月榜单
//	@Description	获取指定月份的月榜单数据
//	@Tags			书籍排行榜
//	@Accept			json
//	@Produce		json
//	@Param			period	query		string	false	"月份 (格式: 2024-01)"
//	@Param			limit	query		int		false	"限制数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/rankings/monthly [get]
func (api *BookstoreAPI) GetMonthlyRanking(c *gin.Context) {
	period := c.Query("period")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rankings, err := api.service.GetMonthlyRanking(c.Request.Context(), period, limit)
	if err != nil {
		shared.InternalError(c, "获取月榜失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取月榜成功", rankings)
}

// GetNewbieRanking 获取新人榜
//
//	@Summary		获取新人榜单
//	@Description	获取指定月份的新人榜单数据
//	@Tags			书籍排行榜
//	@Accept			json
//	@Produce		json
//	@Param			period	query		string	false	"月份 (格式: 2024-01)"
//	@Param			limit	query		int		false	"限制数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/rankings/newbie [get]
func (api *BookstoreAPI) GetNewbieRanking(c *gin.Context) {
	period := c.Query("period")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rankings, err := api.service.GetNewbieRanking(c.Request.Context(), period, limit)
	if err != nil {
		shared.InternalError(c, "获取新人榜失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取新人榜成功", rankings)
}

// GetRankingByType 根据类型获取榜单
//
//	@Summary		根据类型获取榜单
//	@Description	根据榜单类型获取指定周期的榜单数据
//	@Tags			书籍排行榜
//	@Accept			json
//	@Produce		json
//	@Param			type	path		string	true	"榜单类型"	Enums(realtime,weekly,monthly,newbie)
//	@Param			period	query		string	false	"周期"
//	@Param			limit	query		int		false	"限制数量"	default(20)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/bookstore/rankings/{type} [get]
func (api *BookstoreAPI) GetRankingByType(c *gin.Context) {
	rankingType := c.Param("type")
	if !bookstore2.IsValidRankingType(rankingType) {
		shared.BadRequest(c, "参数错误", "无效的榜单类型")
		return
	}

	period := c.Query("period")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rankings, err := api.service.GetRankingByType(c.Request.Context(), bookstore2.RankingType(rankingType), period, limit)
	if err != nil {
		shared.InternalError(c, "获取榜单失败", err)
		return
	}

	shared.Success(c, http.StatusOK, "获取榜单成功", rankings)
}

// buildSearchFilter 构建搜索过滤条件
func (api *BookstoreAPI) buildSearchFilter(categoryID, author, status string, tags []string) map[string]interface{} {
	filter := make(map[string]interface{})

	if categoryID != "" {
		filter["category_id"] = categoryID
	}

	if author != "" {
		filter["author"] = author
	}

	if status != "" {
		// 映射前端状态值到后端状态值
		var backendStatus string
		switch status {
		case "serializing":
			backendStatus = "ongoing"
		case "completed", "paused":
			backendStatus = status
		default:
			backendStatus = status
		}
		filter["status"] = backendStatus
	}

	if len(tags) > 0 {
		filter["tags"] = tags
	}

	return filter
}

// buildSearchSort 构建搜索排序条件
func (api *BookstoreAPI) buildSearchSort(sortBy, sortOrder string) []searchModels.SortField {
	var ascending bool
	if sortOrder == "asc" {
		ascending = true
	}

	return []searchModels.SortField{
		{
			Field:     sortBy,
			Ascending: ascending,
		},
	}
}

// convertSearchResponseToBooks 将搜索响应转换为 Book 切片
func (api *BookstoreAPI) convertSearchResponseToBooks(items []searchModels.SearchItem) []*bookstore2.Book {
	books := make([]*bookstore2.Book, 0, len(items))

	for _, item := range items {
		book := &bookstore2.Book{}

		// 从 Data 中提取字段
		if id, ok := item.Data["id"].(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(id); err == nil {
				book.ID = objectID
			}
		}
		if title, ok := item.Data["title"].(string); ok {
			book.Title = title
		}
		if author, ok := item.Data["author"].(string); ok {
			book.Author = author
		}
		if intro, ok := item.Data["introduction"].(string); ok {
			book.Introduction = intro
		}
		if coverURL, ok := item.Data["cover_url"].(string); ok {
			book.Cover = coverURL
		}
		if viewCount, ok := item.Data["view_count"].(int64); ok {
			book.ViewCount = viewCount
		}
		if rating, ok := item.Data["rating"].(float64); ok {
			book.Rating = types.Rating(rating)
		}
		if wordCount, ok := item.Data["word_count"].(int64); ok {
			book.WordCount = wordCount
		}
		if status, ok := item.Data["status"].(string); ok {
			book.Status = bookstore2.BookStatus(status)
		}

		books = append(books, book)
	}

	return books
}

// GetBooksByTags 按标签筛选书籍
//
//	@Summary     按标签筛选书籍
//	@Description 根据一个或多个标签筛选书籍，ANY语义（命中任一即可）
//	@Tags        书籍筛选
//	@Accept      json
//	@Produce     json
//	@Param       tags query string true "标签列表（逗号分隔）"
//	@Param       page query int false "页码" default(1)
//	@Param       size query int false "每页数量" default(20)
//	@Success     200 {object} shared.PaginatedResponse
//	@Failure     400 {object} shared.APIResponse
//	@Failure     500 {object} shared.APIResponse
//	@Router      /api/v1/bookstore/books/tags [get]
func (api *BookstoreAPI) GetBooksByTags(c *gin.Context) {
	// 构建日志记录器
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = primitive.NewObjectID().Hex()
	}

	searchLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	tagsStr := c.Query("tags")
	if tagsStr == "" {
		shared.BadRequest(c, "参数错误", "标签不能为空")
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

	startTime := time.Now()

	// 解析标签（逗号分隔）
	tags := strings.Split(tagsStr, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	// 记录筛选请求
	searchLogger.WithModule("bookstore").Info("按标签筛选请求",
		zap.Strings("tags", tags),
		zap.Int("page", page),
		zap.Int("size", size),
	)

	// 使用 Service 层进行查询
	filter := &bookstore2.BookFilter{
		Tags:      tags, // ANY语义：Service层实现
		SortBy:    "view_count",
		SortOrder: "desc",
		Limit:     size,
		Offset:    (page - 1) * size,
	}

	books, total, err := api.service.SearchBooksWithFilter(c.Request.Context(), filter)
	if err != nil {
		searchLogger.WithModule("bookstore").Error("按标签筛选失败",
			zap.Error(err),
			zap.Strings("tags", tags),
		)
		shared.InternalError(c, "获取书籍列表失败", err)
		return
	}

	// 转换为 DTO
	bookDTOs := ToBookDTOsFromPtrSlice(books)

	duration := time.Since(startTime)
	searchLogger.WithModule("bookstore").Info("按标签筛选成功",
		zap.Int64("total", total),
		zap.Int("count", len(bookDTOs)),
		zap.Duration("duration", duration),
	)

	shared.Paginated(c, bookDTOs, total, page, size, "获取书籍列表成功")
}

// GetBooksByStatus 按状态筛选书籍
//
//	@Summary     按状态筛选书籍
//	@Description 根据书籍连载状态筛选书籍，支持分页
//	@Tags        书籍筛选
//	@Accept      json
//	@Produce     json
//	@Param       status query string true "书籍状态" Enums(ongoing,completed,paused)
//	@Param       page query int false "页码" default(1)
//	@Param       size query int false "每页数量" default(20)
//	@Success     200 {object} shared.PaginatedResponse
//	@Failure     400 {object} shared.APIResponse
//	@Failure     500 {object} shared.APIResponse
//	@Router      /api/v1/bookstore/books/status [get]
func (api *BookstoreAPI) GetBooksByStatus(c *gin.Context) {
	// 构建日志记录器
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = primitive.NewObjectID().Hex()
	}

	searchLogger := logger.WithRequest(
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		c.ClientIP(),
	)

	status := c.Query("status")
	if status == "" {
		shared.BadRequest(c, "参数错误", "状态不能为空")
		return
	}

	// 验证状态值
	if status != "ongoing" && status != "completed" && status != "paused" {
		shared.BadRequest(c, "参数错误", "无效的状态值")
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

	startTime := time.Now()

	// 记录筛选请求
	searchLogger.WithModule("bookstore").Info("按状态筛选请求",
		zap.String("status", status),
		zap.Int("page", page),
		zap.Int("size", size),
	)

	bookStatus := bookstore2.BookStatus(status)
	filter := &bookstore2.BookFilter{
		Status:    &bookStatus,
		SortBy:    "updated_at", // 状态筛选按更新时间排序
		SortOrder: "desc",
		Limit:     size,
		Offset:    (page - 1) * size,
	}

	books, total, err := api.service.SearchBooksWithFilter(c.Request.Context(), filter)
	if err != nil {
		searchLogger.WithModule("bookstore").Error("按状态筛选失败",
			zap.Error(err),
			zap.String("status", status),
		)
		shared.InternalError(c, "获取书籍列表失败", err)
		return
	}

	// 转换为 DTO
	bookDTOs := ToBookDTOsFromPtrSlice(books)

	duration := time.Since(startTime)
	searchLogger.WithModule("bookstore").Info("按状态筛选成功",
		zap.Int64("total", total),
		zap.Int("count", len(bookDTOs)),
		zap.Duration("duration", duration),
	)

	shared.Paginated(c, bookDTOs, total, page, size, "获取书籍列表成功")
}

// GetSimilarBooks 获取相似书籍推荐
//
//	@Summary     获取相似书籍推荐
//	@Description 基于书籍分类、标签等推荐相似书籍，有四层降级策略
//	@Tags        书籍交互
//	@Accept      json
//	@Produce     json
//	@Param       id path string true "书籍ID"
//	@Param       limit query int false "返回数量" default(10)
//	@Success     200 {object} APIResponse
//	@Failure     400 {object} APIResponse
//	@Failure     404 {object} APIResponse
//	@Failure     500 {object} APIResponse
//	@Router      /api/v1/bookstore/books/{id}/similar [get]
func (api *BookstoreAPI) GetSimilarBooks(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		shared.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 获取原书籍信息
	book, err := api.service.GetBookByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "book not found" {
			shared.NotFound(c, "书籍不存在")
			return
		}
		shared.InternalError(c, "获取书籍信息失败", err)
		return
	}

	// v1.2实现：四层降级策略
	// 策略1: 同分类 + 标签
	result := api.searchSimilarBooks(c.Request.Context(), book, limit, true, true)

	// 策略2: 如果结果不足，尝试同分类
	if len(result) < limit {
		additional := api.searchSimilarBooks(c.Request.Context(), book, limit-len(result), true, false)
		result = append(result, additional...)
	}

	// 策略3: 如果还不足，尝试标签匹配
	if len(result) < limit {
		additional := api.searchSimilarBooks(c.Request.Context(), book, limit-len(result), false, true)
		result = append(result, additional...)
	}

	// 策略4: 兜底 - 返回热门书籍（禁止返回空列表）
	if len(result) == 0 {
		popularBooks, _, err := api.service.GetHotBooks(c.Request.Context(), 1, limit)
		if err == nil && len(popularBooks) > 0 {
			result = popularBooks
		} else {
			// 如果连热门书籍都获取失败，尝试推荐书籍
			recommendedBooks, _, err := api.service.GetRecommendedBooks(c.Request.Context(), 1, limit)
			if err == nil && len(recommendedBooks) > 0 {
				result = recommendedBooks
			}
		}
	}

	// 去重并限制数量
	uniqueResult := api.deduplicateAndLimitBooks(result, book.ID, limit)

	// 确保不返回空列表
	if len(uniqueResult) == 0 {
		shared.Success(c, http.StatusOK, "暂无相似书籍", []*bookstore2.Book{})
		return
	}

	bookDTOs := ToBookDTOsFromPtrSlice(uniqueResult)
	shared.Success(c, http.StatusOK, "获取相似书籍成功", bookDTOs)
}

// searchSimilarBooks 辅助函数：搜索相似书籍
func (api *BookstoreAPI) searchSimilarBooks(
	ctx context.Context,
	book *bookstore2.Book,
	limit int,
	useCategory bool,
	useTags bool,
) []*bookstore2.Book {

	filter := &bookstore2.BookFilter{
		Limit:      limit + 1, // 多查一条用于排除当前书籍
		Offset:     0,
		SortBy:     "view_count",
		SortOrder:  "desc",
	}

	// 注意：Book 模型使用 CategoryIDs 数组而非单个 CategoryID
	// 如果书籍有分类ID，使用第一个分类ID进行查询
	if useCategory && len(book.CategoryIDs) > 0 {
		firstCategoryID := book.CategoryIDs[0].Hex()
		filter.CategoryID = &firstCategoryID
	}

	if useTags && len(book.Tags) > 0 {
		filter.Tags = book.Tags
	}

	books, _, err := api.service.SearchBooksWithFilter(ctx, filter)
	if err != nil {
		api.logger.Warn("搜索相似书籍失败",
			zap.String("book_id", book.ID.Hex()),
			zap.Error(err),
		)
		return []*bookstore2.Book{}
	}

	return books
}

// deduplicateAndLimitBooks 去重并限制数量
func (api *BookstoreAPI) deduplicateAndLimitBooks(
	books []*bookstore2.Book,
	excludeID primitive.ObjectID,
	limit int,
) []*bookstore2.Book {

	seen := make(map[string]bool)
	result := make([]*bookstore2.Book, 0, len(books))

	for _, b := range books {
		bookID := b.ID.Hex()
		// 排除当前书籍和已添加的书籍
		if b.ID != excludeID && !seen[bookID] {
			result = append(result, b)
			seen[bookID] = true
		}
		if len(result) >= limit {
			break
		}
	}
	return result
}
