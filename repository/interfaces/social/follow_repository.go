package social

import (
	"context"

	"Qingyu_backend/models/social"
)

// FollowRepository 关注仓储接口
type FollowRepository interface {
	// ========== 用户关注 ==========

	// CreateFollow 创建关注关系
	CreateFollow(ctx context.Context, follow *social.Follow) error

	// DeleteFollow 删除关注关系
	DeleteFollow(ctx context.Context, followerID, followingID, followType string) error

	// GetFollow 获取关注关系
	GetFollow(ctx context.Context, followerID, followingID, followType string) (*social.Follow, error)

	// IsFollowing 检查是否已关注
	IsFollowing(ctx context.Context, followerID, followingID, followType string) (bool, error)

	// GetFollowers 获取粉丝列表
	GetFollowers(ctx context.Context, userID string, followType string, page, size int) ([]*social.FollowInfo, int64, error)

	// GetFollowing 获取关注列表
	GetFollowing(ctx context.Context, userID string, followType string, page, size int) ([]*social.FollowingInfo, int64, error)

	// UpdateMutualStatus 更新互相关注状态
	UpdateMutualStatus(ctx context.Context, followerID, followingID, followType string, isMutual bool) error

	// ========== 作者关注 ==========

	// CreateAuthorFollow 创建作者关注
	CreateAuthorFollow(ctx context.Context, authorFollow *social.AuthorFollow) error

	// DeleteAuthorFollow 删除作者关注
	DeleteAuthorFollow(ctx context.Context, userID, authorID string) error

	// GetAuthorFollow 获取作者关注
	GetAuthorFollow(ctx context.Context, userID, authorID string) (*social.AuthorFollow, error)

	// GetAuthorFollowers 获取作者的粉丝列表
	GetAuthorFollowers(ctx context.Context, authorID string, page, size int) ([]*social.FollowInfo, int64, error)

	// GetUserFollowingAuthors 获取用户关注的作者列表
	GetUserFollowingAuthors(ctx context.Context, userID string, page, size int) ([]*social.AuthorFollow, int64, error)

	// ========== 统计 ==========

	// GetFollowStats 获取关注统计
	GetFollowStats(ctx context.Context, userID string) (*social.FollowStats, error)

	// UpdateFollowStats 更新关注统计
	UpdateFollowStats(ctx context.Context, userID string, followerDelta, followingDelta int) error

	// CountFollowers 统计粉丝数
	CountFollowers(ctx context.Context, userID, followType string) (int64, error)

	// CountFollowing 统计关注数
	CountFollowing(ctx context.Context, userID, followType string) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
