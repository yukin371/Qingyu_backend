package reader

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/reading"
)

// CommentAPI 评论API处理器
type CommentAPI struct {
	commentService *reading.CommentService
}

// NewCommentAPI 创建评论API实例
func NewCommentAPI(commentService *reading.CommentService) *CommentAPI {
	return &CommentAPI{
		commentService: commentService,
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	BookID    string `json:"book_id" binding:"required"`
	ChapterID string `json:"chapter_id"`
	Content   string `json:"content" binding:"required,min=10,max=500"`
	Rating    int    `json:"rating" binding:"min=0,max=5"`
}

// ReplyCommentRequest 回复评论请求
type ReplyCommentRequest struct {
	Content string `json:"content" binding:"required,min=10,max=500"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=10,max=500"`
}

// CreateComment 发表评论
//
//	@Summary	发表评论
//	@Tags		评论
//	@Param		request	body		CreateCommentRequest	true	"评论信息"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments [post]
func (api *CommentAPI) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 发表评论
	comment, err := api.commentService.PublishComment(
		c.Request.Context(),
		userID.(string),
		req.BookID,
		req.ChapterID,
		req.Content,
		req.Rating,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "发表评论失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "发表评论成功", comment)
}

// GetCommentList 获取评论列表
//
//	@Summary	获取评论列表
//	@Tags		评论
//	@Param		book_id	query		string	true	"书籍ID"
//	@Param		sortBy	query		string	false	"排序方式(latest/hot)"	default(latest)
//	@Param		page	query		int		false	"页码"					default(1)
//	@Param		size	query		int		false	"每页数量"				default(20)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments [get]
func (api *CommentAPI) GetCommentList(c *gin.Context) {
	bookID := c.Query("book_id")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	sortBy := c.DefaultQuery("sortBy", "latest")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	comments, total, err := api.commentService.GetCommentList(
		c.Request.Context(),
		bookID,
		sortBy,
		page,
		size,
	)

	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取评论列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"size":     size,
	})
}

// GetCommentDetail 获取评论详情
//
//	@Summary	获取评论详情
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id} [get]
func (api *CommentAPI) GetCommentDetail(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	comment, err := api.commentService.GetCommentDetail(c.Request.Context(), commentID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "获取评论详情失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", comment)
}

// UpdateComment 更新评论
//
//	@Summary	更新评论
//	@Tags		评论
//	@Param		id		path		string					true	"评论ID"
//	@Param		request	body		UpdateCommentRequest	true	"评论内容"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id} [put]
func (api *CommentAPI) UpdateComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 更新评论
	err := api.commentService.UpdateComment(
		c.Request.Context(),
		userID.(string),
		commentID,
		req.Content,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "更新评论失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// DeleteComment 删除评论
//
//	@Summary	删除评论
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id} [delete]
func (api *CommentAPI) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 删除评论
	err := api.commentService.DeleteComment(
		c.Request.Context(),
		userID.(string),
		commentID,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "删除评论失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// ReplyComment 回复评论
//
//	@Summary	回复评论
//	@Tags		评论
//	@Param		id		path		string					true	"父评论ID"
//	@Param		request	body		ReplyCommentRequest	true	"回复内容"
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id}/reply [post]
func (api *CommentAPI) ReplyComment(c *gin.Context) {
	parentCommentID := c.Param("id")
	if parentCommentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "父评论ID不能为空")
		return
	}

	var req ReplyCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 回复评论
	comment, err := api.commentService.ReplyComment(
		c.Request.Context(),
		userID.(string),
		parentCommentID,
		req.Content,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "回复评论失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "回复成功", comment)
}

// LikeComment 点赞评论
//
//	@Summary	点赞评论
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id}/like [post]
func (api *CommentAPI) LikeComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 点赞评论（实际实现将在LikeService中）
	err := api.commentService.LikeComment(
		c.Request.Context(),
		userID.(string),
		commentID,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "点赞失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "点赞成功", nil)
}

// UnlikeComment 取消点赞评论
//
//	@Summary	取消点赞评论
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id}/like [delete]
func (api *CommentAPI) UnlikeComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	// 取消点赞评论（实际实现将在LikeService中）
	err := api.commentService.UnlikeComment(
		c.Request.Context(),
		userID.(string),
		commentID,
	)

	if err != nil {
		shared.Error(c, http.StatusBadRequest, "取消点赞失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "取消点赞成功", nil)
}

// GetCommentThread 获取评论完整线程（包含所有回复）
//
//	@Summary	获取评论完整线程
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id}/thread [get]
func (api *CommentAPI) GetCommentThread(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	// 获取评论及其所有回复（树状结构）
	// TODO: 需要在 CommentService 中实现 GetCommentThread 方法
	// thread, err := api.commentService.GetCommentThread(c.Request.Context(), commentID)
	// if err != nil {
	// 	shared.Error(c, http.StatusInternalServerError, "获取评论线程失败", err.Error())
	// 	return
	// }

	// 临时实现：返回评论详情
	comment, err := api.commentService.GetCommentDetail(c.Request.Context(), commentID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "获取评论失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", comment)
}

// GetTopComments 获取热门评论
//
//	@Summary	获取热门评论
//	@Tags		评论
//	@Param		book_id	query		string	true	"书籍ID"
//	@Param		limit	query		int		false	"返回数量"	default(10)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/top [get]
func (api *CommentAPI) GetTopComments(c *gin.Context) {
	bookID := c.Query("book_id")
	if bookID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "书籍ID不能为空")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50 // 限制最大返回数量
	}

	// 获取热门评论（按点赞数、回复数等排序）
	// TODO: 需要在 CommentService 中实现 GetTopComments 方法
	// topComments, err := api.commentService.GetTopComments(c.Request.Context(), bookID, limit)
	// if err != nil {
	// 	shared.Error(c, http.StatusInternalServerError, "获取热门评论失败", err.Error())
	// 	return
	// }

	// 临时实现：返回常规评论列表，按热度排序
	comments, total, err := api.commentService.GetCommentList(c.Request.Context(), bookID, "hot", 1, limit)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取热门评论失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"comments": comments,
		"total":    total,
	})
}

// GetCommentReplies 获取评论的回复列表（分页）
//
//	@Summary	获取评论回复列表
//	@Tags		评论
//	@Param		id		path		string	true	"评论ID"
//	@Param		page	query		int		false	"页码"		default(1)
//	@Param		size	query		int		false	"每页数量"	default(20)
//	@Success	200		{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{id}/replies [get]
func (api *CommentAPI) GetCommentReplies(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "评论ID不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 获取评论的直接回复（分页）
	// TODO: 需要在 CommentService 中实现 GetCommentReplies 方法
	// replies, total, err := api.commentService.GetCommentReplies(c.Request.Context(), commentID, page, size)
	// if err != nil {
	// 	shared.Error(c, http.StatusInternalServerError, "获取回复列表失败", err.Error())
	// 	return
	// }

	// 临时实现：返回空列表（前端可以正常渲染）
	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"replies": []interface{}{},
		"total":   0,
		"page":    page,
		"size":    size,
	})
}
