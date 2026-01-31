package social

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/models/social"
)

// UserRelationAPI 用户关系API
type UserRelationAPI struct {
	relationService UserRelationServiceInterface
}

// UserRelationServiceInterface 用户关系服务接口
type UserRelationServiceInterface interface {
	FollowUser(ctx interface{}, followerID, followeeID string) error
	UnfollowUser(ctx interface{}, followerID, followeeID string) error
	IsFollowing(ctx interface{}, followerID, followeeID string) (bool, error)
	GetFollowers(ctx interface{}, userID string, page, pageSize int) ([]*social.UserRelation, int64, error)
	GetFollowing(ctx interface{}, userID string, page, pageSize int) ([]*social.UserRelation, int64, error)
	GetFollowerCount(ctx interface{}, userID string) (int64, error)
	GetFollowingCount(ctx interface{}, userID string) (int64, error)
}

// NewUserRelationAPI 创建用户关系API
func NewUserRelationAPI(relationService UserRelationServiceInterface) *UserRelationAPI {
	return &UserRelationAPI{
		relationService: relationService,
	}
}

// FollowUser 关注用户
//
//	@Summary		关注用户
//	@Description	关注指定用户
//	@Tags			社交-关注
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	string	true	"被关注用户ID"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/social/follow/{userId} [post]
func (api *UserRelationAPI) FollowUser(c *gin.Context) {
	// 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	followerID := userIDInterface.(string)

	// 获取被关注用户ID
	followeeID := c.Param("userId")
	if followeeID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 执行关注操作
	err := api.relationService.FollowUser(c.Request.Context(), followerID, followeeID)
	if err != nil {
		if err.Error() == "已经关注了该用户" {
			response.Conflict(c, "关注失败", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"followed": true,
	})
}

// UnfollowUser 取消关注用户
//
//	@Summary		取消关注用户
//	@Description	取消关注指定用户
//	@Tags			社交-关注
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	string	true	"被取消关注用户ID"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/social/follow/{userId} [delete]
func (api *UserRelationAPI) UnfollowUser(c *gin.Context) {
	// 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	followerID := userIDInterface.(string)

	// 获取被取消关注用户ID
	followeeID := c.Param("userId")
	if followeeID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 执行取消关注操作
	err := api.relationService.UnfollowUser(c.Request.Context(), followerID, followeeID)
	if err != nil {
		if err.Error() == "关注关系不存在" || err.Error() == "已经取消关注了" {
			response.Conflict(c, "取消关注失败", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"followed": false,
	})
}

// CheckIsFollowing 检查是否关注
//
//	@Summary		检查关注状态
//	@Description	检查当前用户是否关注指定用户
//	@Tags			社交-关注
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	string	true	"目标用户ID"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		401		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/social/follow/{userId}/status [get]
func (api *UserRelationAPI) CheckIsFollowing(c *gin.Context) {
	// 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	followerID := userIDInterface.(string)

	// 获取目标用户ID
	followeeID := c.Param("userId")
	if followeeID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 检查关注状态
	isFollowing, err := api.relationService.IsFollowing(c.Request.Context(), followerID, followeeID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"is_following": isFollowing,
	})
}

// GetFollowers 获取粉丝列表
//
//	@Summary		获取粉丝列表
//	@Description	获取指定用户的粉丝列表
//	@Tags			社交-关注
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	string	true	"用户ID"
//	@Param			page	query	int		false	"页码"	default(1)
//	@Param			size	query	int		false	"每页数量"	default(20)
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/social/users/{userId}/followers [get]
func (api *UserRelationAPI) GetFollowers(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 获取粉丝列表
	relations, total, err := api.relationService.GetFollowers(c.Request.Context(), userID, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"followers": relations,
		"total":     total,
		"page":      page,
		"size":      size,
	})
}

// GetFollowing 获取关注列表
//
//	@Summary		获取关注列表
//	@Description	获取指定用户的关注列表
//	@Tags			社交-关注
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	string	true	"用户ID"
//	@Param			page	query	int		false	"页码"	default(1)
//	@Param			size	query	int		false	"每页数量"	default(20)
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/social/users/{userId}/following [get]
func (api *UserRelationAPI) GetFollowing(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 获取关注列表
	relations, total, err := api.relationService.GetFollowing(c.Request.Context(), userID, page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"following": relations,
		"total":     total,
		"page":      page,
		"size":      size,
	})
}

// GetFollowStats 获取关注统计
//
//	@Summary		获取关注统计
//	@Description	获取指定用户的粉丝数和关注数
//	@Tags			社交-关注
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	string	true	"用户ID"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/social/users/{userId}/follow-stats [get]
func (api *UserRelationAPI) GetFollowStats(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("userId")
	if userID == "" {
		response.BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	// 获取统计信息
	followerCount, err := api.relationService.GetFollowerCount(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	followingCount, err := api.relationService.GetFollowingCount(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"follower_count":  followerCount,
		"following_count": followingCount,
	})
}
