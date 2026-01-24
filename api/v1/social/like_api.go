package social

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	socialModels "Qingyu_backend/models/social" // Import for Swagger annotations
	"Qingyu_backend/service/interfaces"
)

// LikeAPI 点赞API处理器
type LikeAPI struct {
	likeService interfaces.LikeService
}

// NewLikeAPI 创建点赞API实例
func NewLikeAPI(likeService interfaces.LikeService) *LikeAPI {
	return &LikeAPI{
		likeService: likeService,
	}
}

// =========================
// 书籍点赞
// =========================

// LikeBook 点赞书籍
// @Summary 点赞书籍
// @Tags 阅读端-点赞
// @Accept json
// @Produce json
// @Param bookId path string true "书籍ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/reader/books/{bookId}/like [post]
// @Security Bearer
func (api *LikeAPI) LikeBook(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 点赞
	if err := api.likeService.LikeBook(c.Request.Context(), userID.(string), bookID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "点赞失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "点赞成功", gin.H{
		"book_id": bookID,
	})
}

// UnlikeBook 取消点赞书籍
// @Summary 取消点赞书籍
// @Tags 阅读端-点赞
// @Accept json
// @Produce json
// @Param bookId path string true "书籍ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/reader/books/{bookId}/like [delete]
// @Security Bearer
func (api *LikeAPI) UnlikeBook(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 取消点赞
	if err := api.likeService.UnlikeBook(c.Request.Context(), userID.(string), bookID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "取消点赞失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "取消点赞成功", gin.H{
		"book_id": bookID,
	})
}

// GetBookLikeInfo 获取书籍点赞信息
// @Summary 获取书籍点赞信息
// @Tags 阅读端-点赞
// @Accept json
// @Produce json
// @Param bookId path string true "书籍ID"
// @Success 200 {object} shared.APIResponse{data=object{like_count=int64,is_liked=bool}}
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/reader/books/{bookId}/like/info [get]
func (api *LikeAPI) GetBookLikeInfo(c *gin.Context) {
	bookID := c.Param("bookId")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	// 获取点赞数
	likeCount, err := api.likeService.GetBookLikeCount(c.Request.Context(), bookID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取点赞数失败", err.Error())
		return
	}

	result := gin.H{
		"like_count": likeCount,
		"is_liked":   false,
	}

	// 如果已登录，检查点赞状态
	if userID, exists := c.Get("user_id"); exists {
		isLiked, err := api.likeService.IsBookLiked(c.Request.Context(), userID.(string), bookID)
		if err == nil {
			result["is_liked"] = isLiked
		}
	}

	shared.Success(c, http.StatusOK, "获取点赞信息成功", result)
}

// =========================
// 评论点赞
// =========================

// LikeComment 点赞评论
// @Summary 点赞评论
// @Tags 阅读端-点赞
// @Accept json
// @Produce json
// @Param id path string true "评论ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/reader/comments/{id}/like [post]
// @Security Bearer
func (api *LikeAPI) LikeComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 点赞
	if err := api.likeService.LikeComment(c.Request.Context(), userID.(string), commentID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "点赞失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "点赞成功", gin.H{
		"comment_id": commentID,
	})
}

// UnlikeComment 取消点赞评论
// @Summary 取消点赞评论
// @Tags 阅读端-点赞
// @Accept json
// @Produce json
// @Param id path string true "评论ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/reader/comments/{id}/like [delete]
// @Security Bearer
func (api *LikeAPI) UnlikeComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 取消点赞
	if err := api.likeService.UnlikeComment(c.Request.Context(), userID.(string), commentID); err != nil {
		shared.Error(c, http.StatusInternalServerError, "取消点赞失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "取消点赞成功", gin.H{
		"comment_id": commentID,
	})
}

// =========================
// 用户点赞列表
// =========================

// GetUserLikedBooks 获取用户点赞的书籍列表
// @Summary 获取用户点赞的书籍列表
// @Tags 阅读端-点赞
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse{data=object{list=[]socialModels.Like,total=int64,page=int,size=int}}
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/reader/likes/books [get]
// @Security Bearer
func (api *LikeAPI) GetUserLikedBooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 获取分页参数
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

	// 查询
	likes, total, err := api.likeService.GetUserLikedBooks(c.Request.Context(), userID.(string), params.Page, params.Size)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取点赞列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取点赞列表成功", gin.H{
		"list":  likes,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// GetUserLikeStats 获取用户点赞统计
// @Summary 获取用户点赞统计
// @Tags 阅读端-点赞
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse{data=object{total_likes=int64}}
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /api/v1/reader/likes/stats [get]
// @Security Bearer
func (api *LikeAPI) GetUserLikeStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	stats, err := api.likeService.GetUserLikeStats(c.Request.Context(), userID.(string))
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取点赞统计失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取点赞统计成功", stats)
}
