package social

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/social"
)

// RatingAPI 评分API处理器
type RatingAPI struct {
	ratingService social.RatingService
}

// NewRatingAPI 创建评分API实例
func NewRatingAPI(ratingService social.RatingService) *RatingAPI {
	return &RatingAPI{
		ratingService: ratingService,
	}
}

// GetRatingStats 获取评分统计
//
//	@Summary		获取评分统计
//	@Description	获取指定目标（评论/书评/书籍）的评分统计信息
//	@Tags			评分
//	@Accept			json
//	@Produce		json
//	@Param			targetType	path		string	true	"目标类型"	Enums(comment, review, book)
//	@Param			targetId	path		string	true	"目标ID"
//	@Success		200			{object}	response.APIResponse	"成功返回评分统计"
//	@Failure		400			{object}	response.APIResponse	"参数错误"
//	@Failure		500			{object}	response.APIResponse	"服务器内部错误"
//	@Router			/api/v1/ratings/{targetType}/{targetId}/stats [get]
func (api *RatingAPI) GetRatingStats(c *gin.Context) {
	targetType := c.Param("targetType")
	targetID := c.Param("targetId")

	// 参数验证
	if targetType == "" || targetID == "" {
		response.BadRequest(c, "参数错误", "targetType和targetId不能为空")
		return
	}

	// 验证目标类型
	validTypes := map[string]bool{
		"comment": true,
		"review":  true,
		"book":    true,
	}
	if !validTypes[targetType] {
		response.BadRequest(c, "参数错误", "targetType必须是comment、review或book")
		return
	}

	// 获取评分统计
	stats, err := api.ratingService.GetRatingStats(c.Request.Context(), targetType, targetID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, stats)
}

// GetUserRating 获取用户评分
//
//	@Summary		获取用户评分
//	@Description	获取当前用户对指定目标的评分
//	@Tags			评分
//	@Accept			json
//	@Produce		json
//	@Param			targetType	path		string	true	"目标类型"	Enums(book, review)
//	@Param			targetId	path		string	true	"目标ID"
//	@Success		200			{object}	response.APIResponse	"成功返回用户评分"
//	@Failure		400			{object}	response.APIResponse	"参数错误"
//	@Failure		401			{object}	response.APIResponse	"未授权"
//	@Failure		500			{object}	response.APIResponse	"服务器内部错误"
//	@Router			/api/v1/ratings/{targetType}/{targetId}/user-rating [get]
func (api *RatingAPI) GetUserRating(c *gin.Context) {
	// 验证用户身份
	userID := c.GetString("userId")
	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	targetType := c.Param("targetType")
	targetID := c.Param("targetId")

	// 参数验证
	if targetType == "" || targetID == "" {
		response.BadRequest(c, "参数错误", "targetType和targetId不能为空")
		return
	}

	// 验证目标类型
	validTypes := map[string]bool{
		"book":   true,
		"review": true,
	}
	if !validTypes[targetType] {
		response.BadRequest(c, "参数错误", "targetType必须是book或review")
		return
	}

	// 获取用户评分
	rating, err := api.ratingService.GetUserRating(c.Request.Context(), userID, targetType, targetID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"rating": rating,
	})
}
