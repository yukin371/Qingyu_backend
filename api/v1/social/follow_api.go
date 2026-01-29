package social

import (

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// FollowAPI 关注API处理器
type FollowAPI struct {
	followService interfaces.FollowService
}

// NewFollowAPI 创建关注API实例
func NewFollowAPI(followService interfaces.FollowService) *FollowAPI {
	return &FollowAPI{
		followService: followService,
	}
}

// =========================
// 用户关注
// =========================

// FollowUser 关注用户
// @Summary 关注用户
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param userId path string true "被关注用户ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/users/{userId}/follow [post]
// @Security Bearer
func (api *FollowAPI) FollowUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	err := api.followService.FollowUser(c.Request.Context(), currentUserID.(string), userID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "不能关注自己" {
			response.BadRequest(c, "操作失败", errMsg)
		} else if errMsg == "已经关注过该用户" {
			response.BadRequest(c, "操作失败", errMsg)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// UnfollowUser 取消关注用户
// @Summary 取消关注用户
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param userId path string true "被关注用户ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/users/{userId}/unfollow [delete]
// @Security Bearer
func (api *FollowAPI) UnfollowUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	err := api.followService.UnfollowUser(c.Request.Context(), currentUserID.(string), userID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "未关注该用户" {
			response.BadRequest(c, "操作失败", errMsg)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// GetFollowers 获取粉丝列表
// @Summary 获取粉丝列表
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/users/{userId}/followers [get]
// @Security Bearer
func (api *FollowAPI) GetFollowers(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	followers, total, err := api.followService.GetFollowers(
		c.Request.Context(),
		userID,
		params.Page,
		params.Size,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  followers,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// GetFollowing 获取关注列表
// @Summary 获取关注列表
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/users/{userId}/following [get]
// @Security Bearer
func (api *FollowAPI) GetFollowing(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	following, total, err := api.followService.GetFollowing(
		c.Request.Context(),
		userID,
		params.Page,
		params.Size,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  following,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}

// CheckFollowStatus 检查关注状态
// @Summary 检查关注状态
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/users/{userId}/follow-status [get]
// @Security Bearer
func (api *FollowAPI) CheckFollowStatus(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	isFollowing, err := api.followService.CheckFollowStatus(
		c.Request.Context(),
		currentUserID.(string),
		userID,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"is_following": isFollowing,
	})
}

// =========================
// 作者关注
// =========================

// FollowAuthorRequest 关注作者请求
type FollowAuthorRequest struct {
	AuthorName    string `json:"author_name" binding:"required"`
	AuthorAvatar  string `json:"author_avatar"`
	NotifyNewBook bool   `json:"notify_new_book"`
}

// FollowAuthor 关注作者
// @Summary 关注作者
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param authorId path string true "作者ID"
// @Param request body FollowAuthorRequest true "关注信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/authors/{authorId}/follow [post]
// @Security Bearer
func (api *FollowAPI) FollowAuthor(c *gin.Context) {
	authorID := c.Param("authorId")
	if authorID == "" {
		response.BadRequest(c, "参数错误", "作者ID不能为空")
		return
	}

	var req FollowAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	err := api.followService.FollowAuthor(
		c.Request.Context(),
		currentUserID.(string),
		authorID,
		req.AuthorName,
		req.AuthorAvatar,
		req.NotifyNewBook,
	)

	if err != nil {
		errMsg := err.Error()
		if errMsg == "已经关注过该作者" {
			response.BadRequest(c, "操作失败", errMsg)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// UnfollowAuthor 取消关注作者
// @Summary 取消关注作者
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param authorId path string true "作者ID"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/authors/{authorId}/unfollow [delete]
// @Security Bearer
func (api *FollowAPI) UnfollowAuthor(c *gin.Context) {
	authorID := c.Param("authorId")
	if authorID == "" {
		response.BadRequest(c, "参数错误", "作者ID不能为空")
		return
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	err := api.followService.UnfollowAuthor(
		c.Request.Context(),
		currentUserID.(string),
		authorID,
	)

	if err != nil {
		errMsg := err.Error()
		if errMsg == "未关注该作者" {
			response.BadRequest(c, "操作失败", errMsg)
		} else {
			response.InternalError(c, err)
		}
		return
	}

	response.Success(c, nil)
}

// GetFollowingAuthors 获取关注的作者列表
// @Summary 获取关注的作者列表
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/social/following/authors [get]
// @Security Bearer
func (api *FollowAPI) GetFollowingAuthors(c *gin.Context) {
	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	var params struct {
		Page int `form:"page" binding:"min=1"`
		Size int `form:"size" binding:"min=1,max=100"`
	}
	params.Page = 1
	params.Size = 20

	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	authors, total, err := api.followService.GetFollowingAuthors(
		c.Request.Context(),
		currentUserID.(string),
		params.Page,
		params.Size,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  authors,
		"total": total,
		"page":  params.Page,
		"size":  params.Size,
	})
}
