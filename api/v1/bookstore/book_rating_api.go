package bookstore

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/models/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// BookRatingAPI 图书评分API控制器
type BookRatingAPI struct {
	BookRatingService bookstoreService.BookRatingService
}

// NewBookRatingAPI 创建新的图书评分API实例
func NewBookRatingAPI(bookRatingService bookstoreService.BookRatingService) *BookRatingAPI {
	return &BookRatingAPI{
		BookRatingService: bookRatingService,
	}
}

// GetBookRating 获取评分详情
// @Summary 获取评分详情
// @Description 根据评分ID获取评分详情信息
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param id path string true "评分ID"
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 404 {object} APIResponse "评分不存在"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/rating [get]
func (api *BookRatingAPI) GetBookRating(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "评分ID不能为空", nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "无效的评分ID格式", nil)
		return
	}

	rating, err := api.BookRatingService.GetRatingByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "评分不存在")
		return
	}

	response.SuccessWithMessage(c, "获取成功", rating)
}

// GetRatingsByBookID 获取图书的所有评分
// @Summary 获取图书的所有评分
// @Description 根据图书ID获取该图书的所有评分列表
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param book_id path string true "图书ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/ratings [get]
func (api *BookRatingAPI) GetRatingsByBookID(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "图书ID不能为空", nil)
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "无效的图书ID格式", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	ratings, total, err := api.BookRatingService.GetRatingsByBookID(c.Request.Context(), bookID, page, limit)
	if err != nil {
		err = errors.New("获取评分列表失败")
		c.Error(err)
		return
	}

	response.Paginated(c, ratings, total, page, limit, "获取成功")
}

// GetRatingsByUserID 获取用户的所有评分
// @Summary 获取用户的所有评分
// @Description 根据用户ID获取该用户的所有评分列表
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/ratings/user/{id} [get]
func (api *BookRatingAPI) GetRatingsByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		response.BadRequest(c, "用户ID不能为空", nil)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		response.BadRequest(c, "无效的用户ID格式", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	ratings, total, err := api.BookRatingService.GetRatingsByUserID(c.Request.Context(), userID, page, limit)
	if err != nil {
		err = errors.New("获取评分列表失败")
		c.Error(err)
		return
	}

	response.Paginated(c, ratings, total, page, limit, "获取成功")
}

// GetAverageRating 获取图书平均评分
// @Summary 获取图书平均评分
// @Description 根据图书ID获取该图书的平均评分
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param book_id path string true "图书ID"
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/average-rating [get]
func (api *BookRatingAPI) GetAverageRating(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "图书ID不能为空", nil)
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "无效的图书ID格式", nil)
		return
	}

	avgRating, err := api.BookRatingService.GetAverageRating(c.Request.Context(), bookID)
	if err != nil {
		err = errors.New("获取平均评分失败")
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", avgRating)
}

// GetRatingDistribution 获取图书评分分布
// @Summary 获取图书评分分布
// @Description 根据图书ID获取该图书的评分分布统计
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param book_id path string true "图书ID"
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/rating-distribution [get]
func (api *BookRatingAPI) GetRatingDistribution(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		response.BadRequest(c, "图书ID不能为空", nil)
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
	if err != nil {
		response.BadRequest(c, "无效的图书ID格式", nil)
		return
	}

	distribution, err := api.BookRatingService.GetRatingDistribution(c.Request.Context(), bookID)
	if err != nil {
		err = errors.New("获取评分分布失败")
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "获取成功", distribution)
}

// CreateRating 创建评分
// @Summary 创建评分
// @Description 为图书创建新的评分
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param request body object true "评分信息"
// @Success 201 {object} APIResponse "创建成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/rating [post]
func (api *BookRatingAPI) CreateRating(c *gin.Context) {
	var rating bookstore.BookRating
	if err := c.ShouldBindJSON(&rating); err != nil {
		response.BadRequest(c, "请求参数格式错误", nil)
		return
	}

	if err := api.BookRatingService.CreateRating(c.Request.Context(), &rating); err != nil {
		c.Error(err)
		return
	}

	response.Created(c, rating)
}

// UpdateRating 更新评分
// @Summary 更新评分
// @Description 更新指定ID的评分信息
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param id path string true "评分ID"
// @Param request body object true "评分信息"
// @Success 200 {object} APIResponse "更新成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 404 {object} APIResponse "评分不存在"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/rating [put]
func (api *BookRatingAPI) UpdateRating(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "评分ID不能为空", nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "无效的评分ID格式", nil)
		return
	}

	var rating bookstore.BookRating
	if err := c.ShouldBindJSON(&rating); err != nil {
		response.BadRequest(c, "请求参数格式错误", nil)
		return
	}

	rating.ID = id
	err = api.BookRatingService.UpdateRating(c.Request.Context(), &rating)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "更新成功", rating)
}

// DeleteRating 删除评分
// @Summary 删除评分
// @Description 删除指定ID的评分
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param id path string true "评分ID"
// @Success 200 {object} APIResponse "删除成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 404 {object} APIResponse "评分不存在"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/bookstore/books/{id}/rating [delete]
func (api *BookRatingAPI) DeleteRating(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "评分ID不能为空", nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "无效的评分ID格式", nil)
		return
	}

	err = api.BookRatingService.DeleteRating(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// LikeRating 点赞评分
// @Summary 点赞评分
// @Description 为指定评分点赞
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param id path string true "评分ID"
// @Success 200 {object} APIResponse "点赞成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/ratings/{id}/like [post]
func (api *BookRatingAPI) LikeRating(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "评分ID不能为空", nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "无效的评分ID格式", nil)
		return
	}

	// 从上下文获取用户ID（假设已在中间件中设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	userObjID, ok := userID.(primitive.ObjectID)
	if !ok {
		// 尝试从string转换
		userIDStr, ok := userID.(string)
		if !ok {
			err = errors.New("用户ID格式错误")
			c.Error(err)
			return
		}
		userObjID, err = primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			err = errors.New("用户ID格式错误")
			c.Error(err)
			return
		}
	}

	err = api.BookRatingService.LikeRating(c.Request.Context(), id, userObjID)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "点赞成功", nil)
}

// UnlikeRating 取消点赞评分
// @Summary 取消点赞评分
// @Description 取消对指定评分的点赞
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param id path string true "评分ID"
// @Success 200 {object} APIResponse "取消点赞成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/ratings/{id}/unlike [post]
func (api *BookRatingAPI) UnlikeRating(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "评分ID不能为空", nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		response.BadRequest(c, "无效的评分ID格式", nil)
		return
	}

	// 从上下文获取用户ID（假设已在中间件中设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	userObjID, ok := userID.(primitive.ObjectID)
	if !ok {
		// 尝试从string转换
		userIDStr, ok := userID.(string)
		if !ok {
			err = errors.New("用户ID格式错误")
			c.Error(err)
			return
		}
		userObjID, err = primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			err = errors.New("用户ID格式错误")
			c.Error(err)
			return
		}
	}

	err = api.BookRatingService.UnlikeRating(c.Request.Context(), id, userObjID)
	if err != nil {
		c.Error(err)
		return
	}

	response.SuccessWithMessage(c, "取消点赞成功", nil)
}

// SearchRatings 搜索评分
// @Summary 搜索评分
// @Description 根据关键词搜索评分
// @Tags 图书评分
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/ratings/search [get]
func (api *BookRatingAPI) SearchRatings(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "搜索关键词不能为空", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// TODO: SearchByKeyword方法尚未在Service层实现
	// 暂时返回空结果
	response.SuccessWithMessage(c, "搜索功能开发中", []bookstore.BookRating{})
}
