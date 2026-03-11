package bookstore

import (
	"context"
	"errors"
	"strconv"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/pkg/response"
	bookstoreService "Qingyu_backend/service/bookstore"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookRatingAPI 图书评分API控制器
type BookRatingAPI struct {
	BookRatingService bookstoreService.BookRatingService
}

type upsertRatingRequest struct {
	BookID string   `json:"bookId"`
	Score  int      `json:"score"`
	Review string   `json:"review"`
	Tags   []string `json:"tags"`
}

type bookRatingSummaryResponse struct {
	BookID          string           `json:"bookId"`
	AverageRating   float64          `json:"averageRating"`
	AverageScore    float64          `json:"averageScore"`
	TotalRatings    int64            `json:"totalRatings"`
	TotalCount      int64            `json:"totalCount"`
	Distribution    map[string]int   `json:"distribution"`
	RawDistribution map[string]int64 `json:"rawDistribution,omitempty"`
}

type userBookRatingResponse struct {
	ID        string   `json:"id"`
	BookID    string   `json:"bookId"`
	UserID    string   `json:"userId"`
	Score     int      `json:"score"`
	Review    string   `json:"review,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
}

// NewBookRatingAPI 创建新的图书评分API实例
func NewBookRatingAPI(bookRatingService bookstoreService.BookRatingService) *BookRatingAPI {
	return &BookRatingAPI{
		BookRatingService: bookRatingService,
	}
}

// GetBookRating 获取书籍评分摘要
// @Summary 获取书籍评分摘要
// @Description 根据图书ID获取评分统计信息
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param id path string true "图书ID"
// @Success 200 {object} response.APIResponse "成功"
// @Failure 400 {object} response.APIResponse "请求参数错误"
// @Failure 500 {object} response.APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/rating [get]
func (api *BookRatingAPI) GetBookRating(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "id", "图书ID")
	if !ok {
		return
	}

	summary, err := api.buildBookRatingSummary(c.Request.Context(), bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", summary)
}

// GetBookRatingCompat 兼容旧版前端 GET /bookstore/ratings/book/:bookId
func (api *BookRatingAPI) GetBookRatingCompat(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "bookId", "图书ID")
	if !ok {
		return
	}

	summary, err := api.buildBookRatingSummary(c.Request.Context(), bookID)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", summary)
}

// GetRatingsByBookID 获取图书的所有评分
func (api *BookRatingAPI) GetRatingsByBookID(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "id", "图书ID")
	if !ok {
		return
	}

	page, limit := parsePageAndLimit(c)
	ratings, total, err := api.BookRatingService.GetRatingsByBookID(c.Request.Context(), bookID, page, limit)
	if err != nil {
		c.Error(errors.New("获取评分列表失败"))
		return
	}

	response.Paginated(c, ratings, total, page, limit, "获取成功")
}

// GetRatingsByUserID 获取用户的所有评分
func (api *BookRatingAPI) GetRatingsByUserID(c *gin.Context) {
	userID, ok := api.parseObjectIDParam(c, "id", "用户ID")
	if !ok {
		return
	}

	page, limit := parsePageAndLimit(c)
	ratings, total, err := api.BookRatingService.GetRatingsByUserID(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.Error(errors.New("获取评分列表失败"))
		return
	}

	response.Paginated(c, ratings, total, page, limit, "获取成功")
}

// GetAverageRating 获取图书平均评分
func (api *BookRatingAPI) GetAverageRating(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "id", "图书ID")
	if !ok {
		return
	}

	avgRating, err := api.BookRatingService.GetAverageRating(c.Request.Context(), bookID)
	if err != nil {
		c.Error(errors.New("获取平均评分失败"))
		return
	}

	response.SuccessWithMessage(c, "获取成功", avgRating)
}

// GetRatingDistribution 获取图书评分分布
func (api *BookRatingAPI) GetRatingDistribution(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "id", "图书ID")
	if !ok {
		return
	}

	distribution, err := api.BookRatingService.GetRatingDistribution(c.Request.Context(), bookID)
	if err != nil {
		c.Error(errors.New("获取评分分布失败"))
		return
	}

	response.SuccessWithMessage(c, "获取成功", distribution)
}

// CreateRating 创建当前用户对书籍的评分
func (api *BookRatingAPI) CreateRating(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "id", "图书ID")
	if !ok {
		return
	}

	api.upsertCurrentUserRatingByBookID(c, bookID)
}

// CreateOrUpdateRatingCompat 兼容旧版前端 POST /bookstore/ratings
func (api *BookRatingAPI) CreateOrUpdateRatingCompat(c *gin.Context) {
	var req upsertRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数格式错误", nil)
		return
	}

	bookID, err := primitive.ObjectIDFromHex(req.BookID)
	if err != nil {
		response.BadRequest(c, "无效的图书ID格式", nil)
		return
	}

	api.upsertCurrentUserRating(c, bookID, req)
}

// UpdateRating 更新当前用户对书籍的评分
func (api *BookRatingAPI) UpdateRating(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "id", "图书ID")
	if !ok {
		return
	}

	api.upsertCurrentUserRatingByBookID(c, bookID)
}

// UpdateUserRatingCompat 兼容旧版前端 PUT /bookstore/ratings/:id
// 前端历史实现会把 bookId 传进 :id，因此这里优先按 bookId 处理，再回退 ratingId。
func (api *BookRatingAPI) UpdateUserRatingCompat(c *gin.Context) {
	currentUserID, ok := api.getCurrentUserObjectID(c)
	if !ok {
		return
	}

	var req upsertRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数格式错误", nil)
		return
	}

	rawID, ok := shared.GetRequiredParam(c, "id", "评分ID")
	if !ok {
		return
	}
	objectID, err := primitive.ObjectIDFromHex(rawID)
	if err != nil {
		response.BadRequest(c, "无效的ID格式", nil)
		return
	}

	if existing, err := api.BookRatingService.GetUserRatingForBook(c.Request.Context(), objectID, currentUserID); err == nil && existing != nil {
		req.BookID = objectID.Hex()
		api.upsertCurrentUserRating(c, objectID, req)
		return
	}

	rating, err := api.BookRatingService.GetRatingByID(c.Request.Context(), objectID)
	if err != nil || rating == nil {
		response.NotFound(c, "评分不存在")
		return
	}
	if rating.UserID != currentUserID {
		response.Forbidden(c, "无权修改他人的评分")
		return
	}

	req.BookID = rating.BookID.Hex()
	api.upsertCurrentUserRating(c, rating.BookID, req)
}

// DeleteRating 删除当前用户对书籍的评分
func (api *BookRatingAPI) DeleteRating(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "id", "图书ID")
	if !ok {
		return
	}

	api.deleteCurrentUserRatingByBookID(c, bookID)
}

// DeleteUserRatingCompat 兼容旧版前端 DELETE /bookstore/ratings/:id
// 历史上 :id 可能是 bookId，也可能是 ratingId。
func (api *BookRatingAPI) DeleteUserRatingCompat(c *gin.Context) {
	currentUserID, ok := api.getCurrentUserObjectID(c)
	if !ok {
		return
	}

	rawID, ok := shared.GetRequiredParam(c, "id", "评分ID")
	if !ok {
		return
	}
	objectID, err := primitive.ObjectIDFromHex(rawID)
	if err != nil {
		response.BadRequest(c, "无效的ID格式", nil)
		return
	}

	if existing, err := api.BookRatingService.GetUserRatingForBook(c.Request.Context(), objectID, currentUserID); err == nil && existing != nil {
		api.deleteCurrentUserRatingByBookID(c, objectID)
		return
	}

	rating, err := api.BookRatingService.GetRatingByID(c.Request.Context(), objectID)
	if err != nil || rating == nil {
		response.NotFound(c, "评分不存在")
		return
	}
	if rating.UserID != currentUserID {
		response.Forbidden(c, "无权删除他人的评分")
		return
	}

	if err := api.BookRatingService.DeleteRating(c.Request.Context(), rating.ID); err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetCurrentUserBookRatingCompat 兼容旧版前端 GET /bookstore/ratings/user/me/book/:bookId
func (api *BookRatingAPI) GetCurrentUserBookRatingCompat(c *gin.Context) {
	bookID, ok := api.parseObjectIDParam(c, "bookId", "图书ID")
	if !ok {
		return
	}
	currentUserID, ok := api.getCurrentUserObjectID(c)
	if !ok {
		return
	}

	rating, err := api.BookRatingService.GetUserRatingForBook(c.Request.Context(), bookID, currentUserID)
	if err != nil {
		c.Error(err)
		return
	}
	if rating == nil {
		response.SuccessWithMessage(c, "获取成功", nil)
		return
	}

	response.SuccessWithMessage(c, "获取成功", api.toUserBookRatingResponse(rating))
}

// LikeRating 点赞评分
func (api *BookRatingAPI) LikeRating(c *gin.Context) {
	ratingID, ok := api.parseObjectIDParam(c, "id", "评分ID")
	if !ok {
		return
	}
	userID, ok := api.getCurrentUserObjectID(c)
	if !ok {
		return
	}

	if err := api.BookRatingService.LikeRating(c.Request.Context(), ratingID, userID); err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "点赞成功", nil)
}

// UnlikeRating 取消点赞评分
func (api *BookRatingAPI) UnlikeRating(c *gin.Context) {
	ratingID, ok := api.parseObjectIDParam(c, "id", "评分ID")
	if !ok {
		return
	}
	userID, ok := api.getCurrentUserObjectID(c)
	if !ok {
		return
	}

	if err := api.BookRatingService.UnlikeRating(c.Request.Context(), ratingID, userID); err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "取消点赞成功", nil)
}

// SearchRatings 搜索评分
func (api *BookRatingAPI) SearchRatings(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "搜索关键词不能为空", nil)
		return
	}

	page, limit := parsePageAndLimit(c)
	ratings, total, err := api.BookRatingService.SearchRatings(c.Request.Context(), keyword, page, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.Paginated(c, ratings, total, page, limit, "获取成功")
}

func (api *BookRatingAPI) upsertCurrentUserRatingByBookID(c *gin.Context, bookID primitive.ObjectID) {
	var req upsertRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数格式错误", nil)
		return
	}

	api.upsertCurrentUserRating(c, bookID, req)
}

func (api *BookRatingAPI) upsertCurrentUserRating(c *gin.Context, bookID primitive.ObjectID, req upsertRatingRequest) {
	currentUserID, ok := api.getCurrentUserObjectID(c)
	if !ok {
		return
	}

	if req.Score < 1 || req.Score > 5 {
		response.BadRequest(c, "评分必须在1到5之间", nil)
		return
	}

	if err := api.BookRatingService.UpdateUserRating(c.Request.Context(), bookID, currentUserID, float64(req.Score), req.Review, req.Tags); err != nil {
		c.Error(err)
		return
	}

	rating, err := api.BookRatingService.GetUserRatingForBook(c.Request.Context(), bookID, currentUserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "操作成功", api.toUserBookRatingResponse(rating))
}

func (api *BookRatingAPI) deleteCurrentUserRatingByBookID(c *gin.Context, bookID primitive.ObjectID) {
	currentUserID, ok := api.getCurrentUserObjectID(c)
	if !ok {
		return
	}

	if err := api.BookRatingService.DeleteUserRating(c.Request.Context(), bookID, currentUserID); err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

func (api *BookRatingAPI) buildBookRatingSummary(ctx context.Context, bookID primitive.ObjectID) (*bookRatingSummaryResponse, error) {
	stats, err := api.BookRatingService.GetRatingStats(ctx, bookID)
	if err != nil {
		return nil, err
	}

	average := api.toFloat64(stats["average_rating"])
	total := api.toInt64(stats["total_ratings"])
	distribution := api.toStringIntMap(stats["rating_distribution"])

	return &bookRatingSummaryResponse{
		BookID:          bookID.Hex(),
		AverageRating:   average,
		AverageScore:    average,
		TotalRatings:    total,
		TotalCount:      total,
		Distribution:    distribution,
		RawDistribution: api.toStringInt64Map(stats["rating_distribution"]),
	}, nil
}

func (api *BookRatingAPI) toUserBookRatingResponse(rating *bookstore.BookRating) *userBookRatingResponse {
	if rating == nil {
		return nil
	}

	return &userBookRatingResponse{
		ID:        rating.ID.Hex(),
		BookID:    rating.BookID.Hex(),
		UserID:    rating.UserID.Hex(),
		Score:     rating.Rating,
		Review:    rating.Comment,
		Tags:      rating.Tags,
		CreatedAt: rating.CreatedAt.Format(timeLayout),
		UpdatedAt: rating.UpdatedAt.Format(timeLayout),
	}
}

func (api *BookRatingAPI) parseObjectIDParam(c *gin.Context, key, displayName string) (primitive.ObjectID, bool) {
	value, ok := shared.GetRequiredParam(c, key, displayName)
	if !ok {
		return primitive.NilObjectID, false
	}

	objectID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		response.BadRequest(c, "无效的"+displayName+"格式", nil)
		return primitive.NilObjectID, false
	}

	return objectID, true
}

func (api *BookRatingAPI) getCurrentUserObjectID(c *gin.Context) (primitive.ObjectID, bool) {
	userID, ok := shared.GetUserID(c)
	if !ok {
		return primitive.NilObjectID, false
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		response.Unauthorized(c, "用户ID格式错误")
		return primitive.NilObjectID, false
	}

	return objectID, true
}

func (api *BookRatingAPI) toFloat64(value interface{}) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case float32:
		return float64(typed)
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	default:
		return 0
	}
}

func (api *BookRatingAPI) toInt64(value interface{}) int64 {
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int32:
		return int64(typed)
	case int64:
		return typed
	case float64:
		return int64(typed)
	default:
		return 0
	}
}

func (api *BookRatingAPI) toStringIntMap(value interface{}) map[string]int {
	result := map[string]int{
		"1": 0,
		"2": 0,
		"3": 0,
		"4": 0,
		"5": 0,
	}

	switch typed := value.(type) {
	case map[string]int64:
		for key, count := range typed {
			result[key] = int(count)
		}
	case map[string]int:
		for key, count := range typed {
			result[key] = count
		}
	case map[string]interface{}:
		for key, count := range typed {
			result[key] = int(api.toInt64(count))
		}
	}

	return result
}

func (api *BookRatingAPI) toStringInt64Map(value interface{}) map[string]int64 {
	result := map[string]int64{
		"1": 0,
		"2": 0,
		"3": 0,
		"4": 0,
		"5": 0,
	}

	switch typed := value.(type) {
	case map[string]int64:
		for key, count := range typed {
			result[key] = count
		}
	case map[string]int:
		for key, count := range typed {
			result[key] = int64(count)
		}
	case map[string]interface{}:
		for key, count := range typed {
			result[key] = api.toInt64(count)
		}
	}

	return result
}

func parsePageAndLimit(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return page, limit
}

const timeLayout = "2006-01-02T15:04:05Z07:00"
