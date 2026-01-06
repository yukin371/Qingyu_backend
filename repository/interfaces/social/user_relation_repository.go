package social

import (
	"context"

	socialModel "Qingyu_backend/models/social"
)

// UserRelationRepository 用户关系仓储接口
type UserRelationRepository interface {
	// 基础CRUD
	Create(ctx context.Context, relation *socialModel.UserRelation) error
	GetByID(ctx context.Context, id string) (*socialModel.UserRelation, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 关系查询
	GetRelation(ctx context.Context, followerID, followeeID string) (*socialModel.UserRelation, error)
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)

	// 列表查询
	GetFollowers(ctx context.Context, followeeID string, limit, offset int) ([]*socialModel.UserRelation, int64, error)
	GetFollowing(ctx context.Context, followerID string, limit, offset int) ([]*socialModel.UserRelation, int64, error)

	// 统计
	CountFollowers(ctx context.Context, userID string) (int64, error)
	CountFollowing(ctx context.Context, userID string) (int64, error)

	// 批量操作
	BatchCreate(ctx context.Context, relations []*socialModel.UserRelation) error
}
