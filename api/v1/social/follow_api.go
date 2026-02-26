package social

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/users/{userId}/follow [post]
// @Security Bearer
func (api *FollowAPI) FollowUser(c *gin.Context) {
	userID, ok := shared.GetRequiredParam(c, "userId", "用户ID")
	if !ok {
		return
	}

	currentUserID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.followService.FollowUser(c.Request.Context(), currentUserID, userID)
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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/users/{userId}/unfollow [delete]
// @Security Bearer
func (api *FollowAPI) UnfollowUser(c *gin.Context) {
	userID, ok := shared.GetRequiredParam(c, "userId", "用户ID")
	if !ok {
		return
	}

	currentUserID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.followService.UnfollowUser(c.Request.Context(), currentUserID, userID)
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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/users/{userId}/followers [get]
// @Security Bearer
func (api *FollowAPI) GetFollowers(c *gin.Context) {
	userID, ok := shared.GetRequiredParam(c, "userId", "用户ID")
	if !ok {
		return
	}

	params := shared.GetPaginationParamsStandard(c)

	followers, total, err := api.followService.GetFollowers(
		c.Request.Context(),
		userID,
		params.Page,
		params.PageSize,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  followers,
		"total": total,
		"page":  params.Page,
		"size":  params.PageSize,
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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/users/{userId}/following [get]
// @Security Bearer
func (api *FollowAPI) GetFollowing(c *gin.Context) {
	userID, ok := shared.GetRequiredParam(c, "userId", "用户ID")
	if !ok {
		return
	}

	params := shared.GetPaginationParamsStandard(c)

	following, total, err := api.followService.GetFollowing(
		c.Request.Context(),
		userID,
		params.Page,
		params.PageSize,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  following,
		"total": total,
		"page":  params.Page,
		"size":  params.PageSize,
	})
}

// CheckFollowStatus 检查关注状态
// @Summary 检查关注状态
// @Tags 社交-关注
// @Accept json
// @Produce json
// @Param userId path string true "用户ID"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/users/{userId}/follow-status [get]
// @Security Bearer
func (api *FollowAPI) CheckFollowStatus(c *gin.Context) {
	userID, ok := shared.GetRequiredParam(c, "userId", "用户ID")
	if !ok {
		return
	}

	currentUserID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	isFollowing, err := api.followService.CheckFollowStatus(
		c.Request.Context(),
		currentUserID,
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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/authors/{authorId}/follow [post]
// @Security Bearer
func (api *FollowAPI) FollowAuthor(c *gin.Context) {
	authorID, ok := shared.GetRequiredParam(c, "authorId", "作者ID")
	if !ok {
		return
	}

	var req FollowAuthorRequest
	if !shared.BindJSON(c, &req) {
		return
	}

	currentUserID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.followService.FollowAuthor(
		c.Request.Context(),
		currentUserID,
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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/authors/{authorId}/unfollow [delete]
// @Security Bearer
func (api *FollowAPI) UnfollowAuthor(c *gin.Context) {
	authorID, ok := shared.GetRequiredParam(c, "authorId", "作者ID")
	if !ok {
		return
	}

	currentUserID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	err := api.followService.UnfollowAuthor(
		c.Request.Context(),
		currentUserID,
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
// @Success 200 {object} response.APIResponse
// @Router /api/v1/social/following/authors [get]
// @Security Bearer
func (api *FollowAPI) GetFollowingAuthors(c *gin.Context) {
	currentUserID, ok := shared.GetUserID(c)
	if !ok {
		return
	}

	params := shared.GetPaginationParamsStandard(c)

	authors, total, err := api.followService.GetFollowingAuthors(
		c.Request.Context(),
		currentUserID,
		params.Page,
		params.PageSize,
	)

	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  authors,
		"total": total,
		"page":  params.Page,
		"size":  params.PageSize,
	})
}
