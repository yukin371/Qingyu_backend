package community

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"

	"Qingyu_backend/models/social"
)

// PostAPI 动态API处理器
type PostAPI struct {
	postService interfaces.PostService
}

// NewPostAPI 创建动态API实例
func NewPostAPI(postService interfaces.PostService) *PostAPI {
	return &PostAPI{
		postService: postService,
	}
}

// CreatePostRequest 创建动态请求
type CreatePostRequest struct {
	Type      social.PostType `json:"type" binding:"required"`
	Content   string          `json:"content" binding:"required,max=5000"`
	Images    []string        `json:"images"`
	BookID    string          `json:"bookId"`
	BookTitle string          `json:"bookTitle"`
	BookCover string          `json:"bookCover"`
	BookAuthor string         `json:"bookAuthor"`
	ChapterID string          `json:"chapterId"`
	ChapterTitle string       `json:"chapterTitle"`
	Progress  int             `json:"progress"`
	Topics    []string        `json:"topics"`
}

// CreatePost 创建动态
// @Summary 创建动态
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "动态信息"
// @Success 201 {object} response.APIResponse
// @Router /api/v1/community/posts [post]
// @Security Bearer
func (api *PostAPI) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	// 获取用户信息
	userName := ""
	userAvatar := ""
	userLevel := 0

	if name, ok := c.Get("username"); ok {
		userName = name.(string)
	}
	if avatar, ok := c.Get("avatar"); ok {
		userAvatar = avatar.(string)
	}
	if level, ok := c.Get("level"); ok {
		userLevel, _ = level.(int)
	}

	post, err := api.postService.CreatePost(
		c.Request.Context(),
		userID,
		userName,
		userAvatar,
		userLevel,
		req.Type,
		req.Content,
		req.Images,
		req.BookID,
		req.BookTitle,
		req.BookCover,
		req.BookAuthor,
		req.ChapterID,
		req.ChapterTitle,
		req.Progress,
		req.Topics,
	)

	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, post)
}

// GetPosts 获取动态列表
// @Summary 获取动态列表
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param topic query string false "话题筛选"
// @Param sort query string false "排序方式(latest/hottest)" default(latest)
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} response.APIResponse
// @Router /api/v1/community/posts [get]
func (api *PostAPI) GetPosts(c *gin.Context) {
	topic := c.Query("topic")
	sort := c.DefaultQuery("sort", "latest")

	pagination := shared.GetPaginationParamsStandard(c)

	// 获取当前用户ID（可选，用于获取点赞状态）
	userID := shared.GetUserIDOptional(c)

	posts, total, err := api.postService.GetPosts(
		c.Request.Context(),
		userID,
		pagination.Page,
		pagination.PageSize,
		topic,
		sort,
	)

	if err != nil {
		c.Error(err)
		return
	}

	response.Paginated(c, posts, total, pagination.Page, pagination.PageSize, "获取动态列表成功")
}

// GetPostDetail 获取动态详情
// @Summary 获取动态详情
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param id path string true "动态ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/community/posts/{id} [get]
func (api *PostAPI) GetPostDetail(c *gin.Context) {
	postID, ok := shared.GetRequiredParam(c, "id", "动态ID")
	if !ok {
		return
	}

	// 获取当前用户ID（可选，用于获取点赞状态）
	userID := shared.GetUserIDOptional(c)

	post, err := api.postService.GetPostByID(c.Request.Context(), userID, postID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, post)
}

// UpdatePostRequest 更新动态请求
type UpdatePostRequest struct {
	Content string   `json:"content" binding:"required,max=5000"`
	Topics  []string `json:"topics"`
}

// UpdatePost 更新动态
// @Summary 更新动态
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param id path string true "动态ID"
// @Param request body UpdatePostRequest true "更新信息"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/community/posts/{id} [put]
// @Security Bearer
func (api *PostAPI) UpdatePost(c *gin.Context) {
	postID, ok := shared.GetRequiredParam(c, "id", "动态ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	var req UpdatePostRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	err := api.postService.UpdatePost(c.Request.Context(), userID, postID, req.Content, req.Topics)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// DeletePost 删除动态
// @Summary 删除动态
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param id path string true "动态ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/community/posts/{id} [delete]
// @Security Bearer
func (api *PostAPI) DeletePost(c *gin.Context) {
	postID, ok := shared.GetRequiredParam(c, "id", "动态ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.postService.DeletePost(c.Request.Context(), userID, postID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, nil)
}

// LikePost 点赞动态
// @Summary 点赞动态
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param id path string true "动态ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/community/posts/{id}/like [post]
// @Security Bearer
func (api *PostAPI) LikePost(c *gin.Context) {
	postID, ok := shared.GetRequiredParam(c, "id", "动态ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.postService.Like(c.Request.Context(), userID, postID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "已经点赞过该动态" {
			response.BadRequest(c, "操作失败", errMsg)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, gin.H{"isLiked": true})
}

// UnlikePost 取消点赞动态
// @Summary 取消点赞动态
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param id path string true "动态ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/community/posts/{id}/like [delete]
// @Security Bearer
func (api *PostAPI) UnlikePost(c *gin.Context) {
	postID, ok := shared.GetRequiredParam(c, "id", "动态ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.postService.Unlike(c.Request.Context(), userID, postID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "未点赞该动态" {
			response.BadRequest(c, "操作失败", errMsg)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, gin.H{"isLiked": false})
}

// ToggleLikePost 切换点赞状态
// @Summary 切换点赞状态
// @Tags 社区-动态
// @Accept json
// @Produce json
// @Param id path string true "动态ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/community/posts/{id}/like [patch]
// @Security Bearer
func (api *PostAPI) ToggleLikePost(c *gin.Context) {
	postID, ok := shared.GetRequiredParam(c, "id", "动态ID")
	if !ok {
		return
	}

	// 获取用户ID
	userID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	isLiked, err := api.postService.ToggleLike(c.Request.Context(), userID, postID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"isLiked": isLiked})
}
