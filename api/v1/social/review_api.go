package social

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/social"
)

// ReviewAPI 书评API处理器
type ReviewAPI struct {
	reviewService *social.ReviewService
}

// NewReviewAPI 创建书评API实例
func NewReviewAPI(reviewService *social.ReviewService) *ReviewAPI {
	return &ReviewAPI{
		reviewService: reviewService,
	}
}

// CreateReviewRequest 创建书评请求
type CreateReviewRequest struct {
	BookID     string `json:"book_id" binding:"required"`
	Title      string `json:"title" binding:"required,max=100"`
	Content    string `json:"content" binding:"required,max=5000"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	IsSpoiler  bool   `json:"is_spoiler"`
	IsPublic   bool   `json:"is_public"`
}

// CreateReview 创建书评
// @Summary 创建书评
// @Tags 社交-书评
// @Accept json
// @Produce json
// @Param request body CreateReviewRequest true "书评信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/reviews [post]
// @Security Bearer
func (api *ReviewAPI) CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 获取用户信息
	userName := ""
	userAvatar := ""
	if name, ok := c.Get("username"); ok {
		userName = name.(string)
	}

	review, err := api.reviewService.CreateReview(
		c.Request.Context(),
		req.BookID,
		userID.(string),
		userName,
		userAvatar,
		req.Title,
		req.Content,
		req.Rating,
		req.IsSpoiler,
		req.IsPublic,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建书评失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建书评成功", review)
}

// GetReviews 获取书评列表
// @Summary 获取书评列表
// @Tags 社交-书评
// @Accept json
// @Produce json
// @Param book_id query string false "书籍ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/reviews [get]
func (api *ReviewAPI) GetReviews(c *gin.Context) {
	bookID := c.Query("book_id")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	reviews, total, err := api.reviewService.GetReviews(
		c.Request.Context(),
		bookID,
		params.Page,
		params.Size,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取书评列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取书评列表成功", gin.H{
		"list":  reviews,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// GetReviewDetail 获取书评详情
// @Summary 获取书评详情
// @Tags 社交-书评
// @Accept json
// @Produce json
// @Param id path string true "书评ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/reviews/{id} [get]
func (api *ReviewAPI) GetReviewDetail(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书评ID不能为空")
		return
	}

	review, err := api.reviewService.GetReviewByID(c.Request.Context(), reviewID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取书评详情失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取书评详情成功", review)
}

// UpdateReviewRequest 更新书评请求
type UpdateReviewRequest struct {
	Title     *string `json:"title" binding:"omitempty,max=100"`
	Content   *string `json:"content" binding:"omitempty,max=5000"`
	Rating    *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	IsSpoiler *bool   `json:"is_spoiler"`
	IsPublic  *bool   `json:"is_public"`
}

// UpdateReview 更新书评
// @Summary 更新书评
// @Tags 社交-书评
// @Accept json
// @Produce json
// @Param id path string true "书评ID"
// @Param request body UpdateReviewRequest true "更新信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/reviews/{id} [put]
// @Security Bearer
func (api *ReviewAPI) UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书评ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.IsSpoiler != nil {
		updates["is_spoiler"] = *req.IsSpoiler
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if len(updates) == 0 {
		shared.Error(c, http.StatusBadRequest, "参数错误", "没有要更新的字段")
		return
	}

	err := api.reviewService.UpdateReview(c.Request.Context(), userID.(string), reviewID, updates)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新书评失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新书评成功", nil)
}

// DeleteReview 删除书评
// @Summary 删除书评
// @Tags 社交-书评
// @Accept json
// @Produce json
// @Param id path string true "书评ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/reviews/{id} [delete]
// @Security Bearer
func (api *ReviewAPI) DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书评ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.reviewService.DeleteReview(c.Request.Context(), userID.(string), reviewID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除书评失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除书评成功", nil)
}

// LikeReview 点赞书评
// @Summary 点赞书评
// @Tags 社交-书评
// @Accept json
// @Produce json
// @Param id path string true "书评ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/reviews/{id}/like [post]
// @Security Bearer
func (api *ReviewAPI) LikeReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书评ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	err := api.reviewService.LikeReview(c.Request.Context(), userID.(string), reviewID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "已经点赞过该书评" {
			shared.Error(c, http.StatusBadRequest, "操作失败", errMsg)
		} else {
			shared.Error(c, http.StatusInternalServerError, "点赞失败", errMsg)
		}
		return
	}

	shared.Success(c, http.StatusOK, "点赞成功", nil)
}
