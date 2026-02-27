package social

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// CommentAPI 评论API处理器
type CommentAPI struct {
	commentService interfaces.CommentService
}

// NewCommentAPI 创建评论API实例
func NewCommentAPI(commentService interfaces.CommentService) *CommentAPI {
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
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/comments [post]
func (api *CommentAPI) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if !shared.BindAndValidate(c, &req) {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 发表评论
	comment, err := api.commentService.PublishComment(
		c.Request.Context(),
		userID,
		req.BookID,
		req.ChapterID,
		req.Content,
		req.Rating,
	)

	if err != nil {
		response.BadRequest(c, "发表评论失败", err.Error())
		return
	}

	response.Created(c, comment)
}

// GetCommentList 获取评论列表
//
//	@Summary	获取评论列表
//	@Tags		评论
//	@Param		book_id	query		string	true	"书籍ID"
//	@Param		sortBy	query		string	false	"排序方式(latest/hot)"	default(latest)
//	@Param		page	query		int		false	"页码"					default(1)
//	@Param		size	query		int		false	"每页数量"				default(20)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/comments [get]
func (api *CommentAPI) GetCommentList(c *gin.Context) {
	bookID := c.Query("book_id")
	if bookID == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
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
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
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
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id} [get]
func (api *CommentAPI) GetCommentDetail(c *gin.Context) {
	commentID, ok := shared.GetRequiredParam(c, "id", "评论ID")
	if !ok {
		return
	}

	comment, err := api.commentService.GetCommentDetail(c.Request.Context(), commentID)
	if err != nil {
		response.NotFound(c, "获取评论详情失败")
		return
	}

	response.Success(c, comment)
}

// UpdateComment 更新评论
//
//	@Summary	更新评论
//	@Tags		评论
//	@Param		id		path		string					true	"评论ID"
//	@Param		request	body		UpdateCommentRequest	true	"评论内容"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id} [put]
func (api *CommentAPI) UpdateComment(c *gin.Context) {
	commentID, ok := shared.GetRequiredParam(c, "id", "评论ID")
	if !ok {
		return
	}

	var req UpdateCommentRequest
	if !shared.BindAndValidate(c, &req) {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 更新评论
	err := api.commentService.UpdateComment(
		c.Request.Context(),
		userID,
		commentID,
		req.Content,
	)

	if err != nil {
		response.BadRequest(c, "更新评论失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteComment 删除评论
//
//	@Summary	删除评论
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id} [delete]
func (api *CommentAPI) DeleteComment(c *gin.Context) {
	commentID, ok := shared.GetRequiredParam(c, "id", "评论ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 删除评论
	err := api.commentService.DeleteComment(
		c.Request.Context(),
		userID,
		commentID,
	)

	if err != nil {
		response.BadRequest(c, "删除评论失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// ReplyComment 回复评论
//
//	@Summary	回复评论
//	@Tags		评论
//	@Param		id		path		string					true	"父评论ID"
//	@Param		request	body		ReplyCommentRequest	true	"回复内容"
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id}/reply [post]
func (api *CommentAPI) ReplyComment(c *gin.Context) {
	parentCommentID, ok := shared.GetRequiredParam(c, "id", "父评论ID")
	if !ok {
		return
	}

	var req ReplyCommentRequest
	if !shared.BindAndValidate(c, &req) {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 回复评论
	comment, err := api.commentService.ReplyComment(
		c.Request.Context(),
		userID,
		parentCommentID,
		req.Content,
	)

	if err != nil {
		response.BadRequest(c, "回复评论失败", err.Error())
		return
	}

	response.Created(c, comment)
}

// LikeComment 点赞评论
//
//	@Summary	点赞评论
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id}/like [post]
func (api *CommentAPI) LikeComment(c *gin.Context) {
	commentID, ok := shared.GetRequiredParam(c, "id", "评论ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 点赞评论（实际实现将在LikeService中）
	err := api.commentService.LikeComment(
		c.Request.Context(),
		userID,
		commentID,
	)

	if err != nil {
		response.BadRequest(c, "点赞失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// UnlikeComment 取消点赞评论
//
//	@Summary	取消点赞评论
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id}/like [delete]
func (api *CommentAPI) UnlikeComment(c *gin.Context) {
	commentID, ok := shared.GetRequiredParam(c, "id", "评论ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 取消点赞评论（实际实现将在LikeService中）
	err := api.commentService.UnlikeComment(
		c.Request.Context(),
		userID,
		commentID,
	)

	if err != nil {
		response.BadRequest(c, "取消点赞失败", err.Error())
		return
	}

	response.Success(c, nil)
}

// GetCommentThread 获取评论完整线程（包含所有回复）
//
//	@Summary	获取评论完整线程
//	@Tags		评论
//	@Param		id	path		string	true	"评论ID"
//	@Success	200	{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id}/thread [get]
func (api *CommentAPI) GetCommentThread(c *gin.Context) {
	commentID, ok := shared.GetRequiredParam(c, "id", "评论ID")
	if !ok {
		return
	}

	// 获取评论及其所有回复（树状结构）
	thread, err := api.commentService.GetCommentThread(c.Request.Context(), commentID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, thread)
}

// GetTopComments 获取热门评论
//
//	@Summary	获取热门评论
//	@Tags		评论
//	@Param		book_id	query		string	true	"书籍ID"
//	@Param		limit	query		int		false	"返回数量"	default(10)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/top [get]
func (api *CommentAPI) GetTopComments(c *gin.Context) {
	bookID := c.Query("book_id")
	if bookID == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50 // 限制最大返回数量
	}

	// 获取热门评论（按点赞数、回复数等排序）
	topComments, err := api.commentService.GetTopComments(c.Request.Context(), bookID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"comments": topComments,
		"total":    len(topComments),
	})
}

// GetCommentReplies 获取评论的回复列表（分页）
//
//	@Summary	获取评论回复列表
//	@Tags		评论
//	@Param		id		path		string	true	"评论ID"
//	@Param		page	query		int		false	"页码"		default(1)
//	@Param		size	query		int		false	"每页数量"	default(20)
//	@Success	200		{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{id}/replies [get]
func (api *CommentAPI) GetCommentReplies(c *gin.Context) {
	commentID, ok := shared.GetRequiredParam(c, "id", "评论ID")
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 获取评论的直接回复（分页）
	replies, total, err := api.commentService.GetCommentReplies(c.Request.Context(), commentID, page, size)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"replies": replies,
		"total":   total,
		"page":    page,
		"size":    size,
	})
}
