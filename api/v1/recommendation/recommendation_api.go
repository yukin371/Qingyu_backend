package recommendation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	reco "Qingyu_backend/models/recommendation/reco"
	"Qingyu_backend/service/recommendation"
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
//	@Success	200		{object}	response.Response
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
//	@Success	200		{object}	response.Response
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
//	@Success	200		{object}	response.Response
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
		ChapterID    string                 `json:"chapterId"`
		BehaviorType string                 `json:"behaviorType" binding:"required"` // view/click/collect/read/finish/like/share
		Value        float64                `json:"value"`
		Metadata     map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 构建行为对象
	behavior := &reco.Behavior{
		UserID:       userID.(string),
		ItemID:       req.ItemID,
		ChapterID:    req.ChapterID,
		BehaviorType: req.BehaviorType,
		Value:        req.Value,
		Metadata:     req.Metadata,
	}

	// 记录行为
	err := api.recoService.RecordBehavior(c.Request.Context(), behavior)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "记录行为失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "记录成功", nil)
}
