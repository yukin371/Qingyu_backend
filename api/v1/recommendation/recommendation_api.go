package recommendation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/shared/recommendation"
)

// RecommendationAPI 推荐API
type RecommendationAPI struct {
	recoService recommendation.RecommendationService
}

// NewRecommendationAPI 创建推荐API实例
func NewRecommendationAPI(recoService recommendation.RecommendationService) *RecommendationAPI {
	return &RecommendationAPI{
		recoService: recoService,
	}
}

// GetPersonalizedRecommendations 获取个性化推荐
//
//	@Summary	获取个性化推荐
//	@Tags		推荐系统
//	@Param		limit	query		int	false	"推荐数量"	default(10)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/recommendation/personalized [get]
func (api *RecommendationAPI) GetPersonalizedRecommendations(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
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
	recommendations, err := api.recoService.GetPersonalizedRecommendations(c.Request.Context(), userID.(string), limit)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取推荐失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
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
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/recommendation/similar [get]
func (api *RecommendationAPI) GetSimilarItems(c *gin.Context) {
	itemID := c.Query("itemId")
	if itemID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "itemId参数不能为空")
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
		shared.Error(c, http.StatusInternalServerError, "获取推荐失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"similar_items": similarItems,
		"count":         len(similarItems),
	})
}

// RecordBehavior 记录用户行为
//
//	@Summary	记录用户行为
//	@Tags		推荐系统
//	@Param		body	body		object	true	"行为数据"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/recommendation/behavior [post]
func (api *RecommendationAPI) RecordBehavior(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 解析请求体
	var req struct {
		ItemID       string                 `json:"itemId" binding:"required"`
		ItemType     string                 `json:"itemType"`                        // book/article等，默认book
		BehaviorType string                 `json:"behaviorType" binding:"required"` // view/click/favorite/read等
		Duration     int64                  `json:"duration"`                        // 阅读时长（秒）
		Metadata     map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 默认类型为book
	if req.ItemType == "" {
		req.ItemType = "book"
	}

	// 构建行为记录请求
	behaviorReq := &recommendation.RecordBehaviorRequest{
		UserID:     userID.(string),
		ItemID:     req.ItemID,
		ItemType:   req.ItemType,
		ActionType: req.BehaviorType,
		Duration:   req.Duration,
		Metadata:   req.Metadata,
	}

	// 记录行为
	err := api.recoService.RecordUserBehavior(c.Request.Context(), behaviorReq)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "记录行为失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "记录成功", nil)
}

// GetHomepageRecommendations 获取首页推荐
//
//	@Summary	获取首页推荐（混合推荐策略：个性化+热门）
//	@Tags		推荐系统
//	@Param		limit	query		int	false	"推荐数量"	default(20)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/recommendation/homepage [get]
func (api *RecommendationAPI) GetHomepageRecommendations(c *gin.Context) {
	// 获取用户ID（可选）
	userID, _ := c.Get("userId")
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}

	// 获取limit参数
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var recommendations []*recommendation.RecommendedItem
	var err error

	// 如果用户已登录，优先返回个性化推荐
	if userIDStr != "" {
		recommendations, err = api.recoService.GetPersonalizedRecommendations(c.Request.Context(), userIDStr, limit)
	}

	// 如果未登录或个性化推荐失败，返回热门推荐
	if userIDStr == "" || err != nil {
		recommendations, err = api.recoService.GetHotItems(c.Request.Context(), "book", limit)
		if err != nil {
			shared.Error(c, http.StatusInternalServerError, "获取首页推荐失败", err.Error())
			return
		}
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetHotRecommendations 获取热门推荐
//
//	@Summary	获取热门推荐
//	@Tags		推荐系统
//	@Param		limit	query		int		false	"推荐数量"	default(20)
//	@Param		type	query		string	false	"物品类型"	default(book)
//	@Success	200		{object}	shared.APIResponse
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
		shared.Error(c, http.StatusInternalServerError, "获取热门推荐失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
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
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/recommendation/category [get]
func (api *RecommendationAPI) GetCategoryRecommendations(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "category参数不能为空")
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
		shared.Error(c, http.StatusInternalServerError, "获取分类推荐失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
		"category":        category,
	})
}
