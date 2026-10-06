package recommendation

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/recommendation"
)

// RecommendationAPI 推荐API
type RecommendationAPI struct {
	recoService  recommendation.RecommendationService
	tableService recommendation.RecommendationTableService
	bookCatalog  recommendation.BookCatalogPort
}

// NewRecommendationAPI 创建推荐API实例
func NewRecommendationAPI(recoService recommendation.RecommendationService) *RecommendationAPI {
	return &RecommendationAPI{
		recoService: recoService,
	}
}

// WithTableService 设置推荐表格服务
func (api *RecommendationAPI) WithTableService(tableService recommendation.RecommendationTableService) *RecommendationAPI {
	api.tableService = tableService
	return api
}

// WithBookCatalog 设置推荐模块读取书籍详情和列表所需的最小端口。
func (api *RecommendationAPI) WithBookCatalog(bookCatalog recommendation.BookCatalogPort) *RecommendationAPI {
	api.bookCatalog = bookCatalog
	return api
}

// WithBookRepository 保留旧装配入口，迁移期转发到最小书籍目录端口。
// Deprecated: use WithBookCatalog.
func (api *RecommendationAPI) WithBookRepository(bookCatalog recommendation.BookCatalogPort) *RecommendationAPI {
	return api.WithBookCatalog(bookCatalog)
}

func authenticatedUserID(c *gin.Context) (string, bool) {
	return shared.GetUserIDWithMessage(c, "请先登录")
}

// GetPersonalizedRecommendations 获取个性化推荐
//
//	@Summary	获取个性化推荐
//	@Tags		推荐系统
//	@Param		limit	query		int	false	"推荐数量"	default(10)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/recommendation/personalized [get]
func (api *RecommendationAPI) GetPersonalizedRecommendations(c *gin.Context) {
	// 获取用户ID
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	// 获取limit参数
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 获取推荐
	recommendations, err := api.recoService.GetPersonalizedRecommendations(c.Request.Context(), userID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetSimilarItems 获取相似物品推荐
//
//	@Summary	获取相似物品推荐
//	@Tags		推荐系统
//	@Param		itemId	query		string	true	"物品ID（书籍ID）"
//	@Param		limit	query		int		false	"推荐数量"	default(10)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/recommendation/similar [get]
func (api *RecommendationAPI) GetSimilarItems(c *gin.Context) {
	itemID := c.Query("itemId")
	if itemID == "" {
		response.BadRequest(c, "参数错误", "itemId参数不能为空")
		return
	}

	// 获取limit参数
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 获取相似物品
	similarItems, err := api.recoService.GetSimilarItems(c.Request.Context(), itemID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", gin.H{
		"similar_items": similarItems,
		"count":         len(similarItems),
	})
}

// RecordBehavior 记录用户行为
//
//	@Summary	记录用户行为
//	@Tags		推荐系统
//	@Param		body	body		object	true	"行为数据"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/recommendation/behavior [post]
func (api *RecommendationAPI) RecordBehavior(c *gin.Context) {
	// 获取用户ID
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	// 解析请求体
	var req struct {
		ItemID       string                 `json:"itemId" binding:"required"`
		ItemType     string                 `json:"itemType"`                        // book/article等，默认book
		BehaviorType string                 `json:"behaviorType" binding:"required"` // view/click/collect/read/finish等
		Duration     int64                  `json:"duration"`                        // 阅读时长（秒）
		Metadata     map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 默认类型为book
	if req.ItemType == "" {
		req.ItemType = "book"
	}

	// 构建行为记录请求
	behaviorReq := &recommendation.RecordBehaviorRequest{
		UserID:     userID,
		ItemID:     req.ItemID,
		ItemType:   req.ItemType,
		ActionType: req.BehaviorType,
		Duration:   req.Duration,
		Metadata:   req.Metadata,
	}

	// 记录行为
	err := api.recoService.RecordUserBehavior(c.Request.Context(), behaviorReq)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "记录成功", nil)
}

// BookRecommendationItem 书籍推荐项（包含书籍详情）
type BookRecommendationItem struct {
	ItemID   string  `json:"item_id"`
	ItemType string  `json:"item_type"`
	Score    float64 `json:"score"`
	Reason   string  `json:"reason"`
	Rank     int     `json:"rank"`
	// 书籍详情（从bookstore服务补充）
	BookID string `json:"book_id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Cover  string `json:"cover,omitempty"`
	Rating string `json:"rating,omitempty"`
}

// HomepageResponse 首页推荐响应
type HomepageResponse struct {
	Recommendations []BookRecommendationItem `json:"recommendations"`
	NewBooks        []*bookstoreModel.Book   `json:"new_books"`
	EditorPicks     []*bookstoreModel.Book   `json:"editor_picks"`
}

// GetHomepageRecommendations 获取首页推荐
//
//	@Summary	获取首页推荐（混合推荐策略：个性化+热门）
//	@Description	返回推荐列表、新书上架、编辑推荐
//	@Tags		推荐系统
//	@Param		limit	query		int	false	"推荐数量"	default(20)
//	@Success	200		{object}	response.APIResponse{data=HomepageResponse}
//	@Router		/api/v1/recommendation/homepage [get]
func (api *RecommendationAPI) GetHomepageRecommendations(c *gin.Context) {
	ctx := c.Request.Context()

	// 获取用户ID（可选）
	userIDStr := shared.GetUserIDOptional(c)

	// 获取limit参数
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 1. 获取推荐列表
	var recoItems []*recommendation.RecommendedItem
	var err error

	if userIDStr != "" {
		recoItems, err = api.recoService.GetPersonalizedRecommendations(ctx, userIDStr, limit)
	}

	if userIDStr == "" || err != nil {
		recoItems, err = api.recoService.GetHotItems(ctx, "book", limit)
		if err != nil {
			response.InternalError(c, err)
			return
		}
	}

	// 2. 转换推荐项并尝试填充书籍详情
	recommendations := make([]BookRecommendationItem, 0, len(recoItems))
	for _, item := range recoItems {
		bookRec := BookRecommendationItem{
			ItemID:   item.ItemID,
			ItemType: item.ItemType,
			Score:    item.Score,
			Reason:   item.Reason,
			Rank:     item.Rank,
		}

		// 如果有书籍仓储且ItemID看起来像真实ID，尝试获取书籍详情
		if api.bookCatalog != nil && item.ItemID != "" {
			book, err := api.bookCatalog.GetByID(ctx, item.ItemID)
			if err == nil && book != nil {
				bookRec.BookID = book.ID.Hex()
				bookRec.Title = book.Title
				bookRec.Author = book.Author
				bookRec.Cover = book.Cover
				bookRec.Rating = book.Rating.String()
			}
		}

		recommendations = append(recommendations, bookRec)
	}

	// 3. 获取新书上架
	var newBooks []*bookstoreModel.Book
	if api.bookCatalog != nil {
		books, err := api.bookCatalog.GetNewReleases(ctx, limit, 0)
		if err == nil {
			newBooks = books
		}
	}

	// 4. 获取编辑推荐
	var editorPicks []*bookstoreModel.Book
	if api.bookCatalog != nil {
		books, err := api.bookCatalog.GetFeatured(ctx, limit, 0)
		if err == nil {
			editorPicks = books
		}
	}

	response.Success(c, HomepageResponse{
		Recommendations: recommendations,
		NewBooks:        newBooks,
		EditorPicks:     editorPicks,
	})
}

// GetHotRecommendations 获取热门推荐
//
//	@Summary	获取热门推荐
//	@Tags		推荐系统
//	@Param		limit	query		int		false	"推荐数量"	default(20)
//	@Param		type	query		string	false	"物品类型"	default(book)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/recommendation/hot [get]
func (api *RecommendationAPI) GetHotRecommendations(c *gin.Context) {
	// 获取limit参数
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 获取type参数，默认为book
	itemType := c.DefaultQuery("type", "book")

	// 获取热门推荐
	recommendations, err := api.recoService.GetHotItems(c.Request.Context(), itemType, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetCategoryRecommendations 获取分类推荐
//
//	@Summary	获取分类推荐（当前使用热门推荐）
//	@Tags		推荐系统
//	@Param		category	query		string	true	"分类名称"
//	@Param		limit		query		int		false	"推荐数量"	default(20)
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/recommendation/category [get]
func (api *RecommendationAPI) GetCategoryRecommendations(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		response.BadRequest(c, "参数错误", "category参数不能为空")
		return
	}

	// 获取limit参数
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 注意：当前使用热门推荐作为分类推荐
	// TODO: 后续可以基于category参数实现真正的分类推荐
	recommendations, err := api.recoService.GetHotItems(c.Request.Context(), "book", limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
		"category":        category,
	})
}
