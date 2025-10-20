package reading

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// BookstoreAPI 书城API处理器
type BookstoreAPI struct {
	service bookstoreService.BookstoreService
}

// NewBookstoreAPI 创建书城API实例
func NewBookstoreAPI(service bookstoreService.BookstoreService) *BookstoreAPI {
	return &BookstoreAPI{
		service: service,
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取首页数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取首页数据成功",
		Data:    data,
	})
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	book, err := api.service.GetBookByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "book not found" || err.Error() == "book not available" {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    404,
				Message: "书籍不存在或不可用",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取书籍详情失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取书籍详情成功",
		Data:    book,
	})
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
	categoryID := c.Param("categoryId")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "分类ID不能为空",
		})
		return
	}

	// 验证ObjectID格式
	if _, err := primitive.ObjectIDFromHex(categoryID); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "分类ID格式无效",
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

	books, total, err := api.service.GetBooksByCategory(c.Request.Context(), categoryID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取分类书籍失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取分类书籍成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// GetRecommendedBooks 获取推荐书籍
//
//	@Summary		获取推荐书籍
//	@Description	获取推荐书籍列表，支持分页
//	@Tags			书籍
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

	books, err := api.service.GetRecommendedBooks(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取推荐书籍失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取推荐书籍成功",
		Data:    books,
		Page:    page,
		Size:    size,
	})
}

// GetFeaturedBooks 获取精选书籍
//
//	@Summary		获取精选书籍
//	@Description	获取精选书籍列表，支持分页
//	@Tags			书籍
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

	books, err := api.service.GetFeaturedBooks(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取精选书籍失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "获取精选书籍成功",
		Data:    books,
		Page:    page,
		Size:    size,
	})
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
	keyword := c.Query("keyword")

	// 构建过滤器
	filter := &bookstore.BookFilter{}

	if categoryID := c.Query("categoryId"); categoryID != "" {
		// 这里需要转换为ObjectID，简化处理
		filter.CategoryID = nil // 实际实现中需要转换
	}

	if author := c.Query("author"); author != "" {
		filter.Author = &author
	}

	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	filter.SortBy = c.DefaultQuery("sortBy", "created_at")
	filter.SortOrder = c.DefaultQuery("sortOrder", "desc")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	filter.Limit = size
	filter.Offset = (page - 1) * size

	if keyword == "" && filter.CategoryID == nil && filter.Author == nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请提供搜索关键词或过滤条件",
		})
		return
	}

	// 设置关键词
	if keyword != "" {
		filter.Keyword = &keyword
	}

	books, total, err := api.service.SearchBooksWithFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "搜索书籍失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    200,
		Message: "搜索书籍成功",
		Data:    books,
		Total:   total,
		Page:    page,
		Size:    size,
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取分类树失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取分类树成功",
		Data:    tree,
	})
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "分类ID不能为空",
		})
		return
	}

	category, err := api.service.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "category not found" || err.Error() == "category not available" {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    404,
				Message: "分类不存在或不可用",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取分类详情失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取分类详情成功",
		Data:    category,
	})
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取Banner列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取Banner列表成功",
		Data:    banners,
	})
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "书籍ID不能为空",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的书籍ID格式",
		})
		return
	}

	err = api.service.IncrementBookView(c.Request.Context(), id.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "增加浏览量失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "浏览量增加成功",
	})
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
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Banner ID不能为空",
		})
		return
	}

	err := api.service.IncrementBannerClick(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "banner not found" || err.Error() == "banner not available" {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:    404,
				Message: "Banner不存在或不可用",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "增加点击次数失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "点击次数增加成功",
	})
}

// GetRealtimeRanking 获取实时榜
//
//	@Summary		获取实时榜单
//	@Description	获取当日实时榜单数据
//	@Tags			榜单
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取实时榜失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取实时榜成功",
		Data:    rankings,
	})
}

// GetWeeklyRanking 获取周榜
//
//	@Summary		获取周榜单
//	@Description	获取指定周期的周榜单数据
//	@Tags			榜单
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取周榜失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取周榜成功",
		Data:    rankings,
	})
}

// GetMonthlyRanking 获取月榜
//
//	@Summary		获取月榜单
//	@Description	获取指定月份的月榜单数据
//	@Tags			榜单
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取月榜失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取月榜成功",
		Data:    rankings,
	})
}

// GetNewbieRanking 获取新人榜
//
//	@Summary		获取新人榜单
//	@Description	获取指定月份的新人榜单数据
//	@Tags			榜单
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
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取新人榜失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取新人榜成功",
		Data:    rankings,
	})
}

// GetRankingByType 根据类型获取榜单
//
//	@Summary		根据类型获取榜单
//	@Description	根据榜单类型获取指定周期的榜单数据
//	@Tags			榜单
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
	if !bookstore.IsValidRankingType(rankingType) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的榜单类型",
		})
		return
	}

	period := c.Query("period")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rankings, err := api.service.GetRankingByType(c.Request.Context(), bookstore.RankingType(rankingType), period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取榜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取榜单成功",
		Data:    rankings,
	})
}
