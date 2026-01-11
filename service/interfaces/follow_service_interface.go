package interfaces

import (
	"context"
	socialModel "Qingyu_backend/models/social"
)

// FollowService 关注服务接口
type FollowService interface {
	// =========================
	// 用户关注
	// =========================
	FollowUser(ctx context.Context, followerID, followingID string) error
	UnfollowUser(ctx context.Context, followerID, followingID string) error
	CheckFollowStatus(ctx context.Context, followerID, followingID string) (bool, error)
	GetFollowers(ctx context.Context, userID string, page, size int) ([]*socialModel.FollowInfo, int64, error)
	GetFollowing(ctx context.Context, userID string, page, size int) ([]*socialModel.FollowingInfo, int64, error)
	GetFollowStats(ctx context.Context, userID string) (*socialModel.FollowStats, error)

	// =========================
	// 作者关注
	// =========================
	FollowAuthor(ctx context.Context, userID, authorID, authorName, authorAvatar string, notifyNewBook bool) error
	UnfollowAuthor(ctx context.Context, userID, authorID string) error
	GetFollowingAuthors(ctx context.Context, userID string, page, size int) ([]*socialModel.AuthorFollow, int64, error)
}
