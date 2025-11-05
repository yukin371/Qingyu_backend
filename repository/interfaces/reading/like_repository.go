package reading

import (
	"context"

	"Qingyu_backend/models/reader"
)

// LikeRepository 点赞仓储接口
type LikeRepository interface {
	// 基础操作
	AddLike(ctx context.Context, like *reader.Like) error
	RemoveLike(ctx context.Context, userID, targetType, targetID string) error
	IsLiked(ctx context.Context, userID, targetType, targetID string) (bool, error)
	GetByID(ctx context.Context, id string) (*reader.Like, error)

	// 查询操作
	GetUserLikes(ctx context.Context, userID, targetType string, page, size int) ([]*reader.Like, int64, error)
	GetLikeCount(ctx context.Context, targetType, targetID string) (int64, error)

	// 批量操作
	GetLikesCountBatch(ctx context.Context, targetType string, targetIDs []string) (map[string]int64, error)
	GetUserLikeStatusBatch(ctx context.Context, userID, targetType string, targetIDs []string) (map[string]bool, error)

	// 统计操作
	CountUserLikes(ctx context.Context, userID string) (int64, error)
	CountTargetLikes(ctx context.Context, targetType, targetID string) (int64, error)

	// 健康检查
	Health(ctx context.Context) error
}
