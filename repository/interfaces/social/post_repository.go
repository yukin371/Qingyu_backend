package social

import (
	"context"

	"Qingyu_backend/models/social"
)

// PostRepository 动态仓储接口
type PostRepository interface {
	// ========== 动态管理 ==========

	// Create 创建动态
	Create(ctx context.Context, post *social.Post) error

	// GetByID 根据ID获取动态
	GetByID(ctx context.Context, id string) (*social.Post, error)

	// List 获取动态列表
	// page: 页码(从1开始)
	// size: 每页数量
	// topic: 话题筛选(可选，空字符串表示不过滤)
	// sort: 排序方式("latest"按时间倒序，"hottest"按点赞数倒序)
	List(ctx context.Context, page, size int, topic, sort string) ([]*social.Post, int64, error)

	// Update 更新动态
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete 删除动态
	Delete(ctx context.Context, id string) error

	// GetByUser 获取用户发布的动态
	GetByUser(ctx context.Context, userID string, page, size int) ([]*social.Post, int64, error)

	// ========== 动态点赞 ==========

	// Like 点赞动态
	Like(ctx context.Context, postID, userID string) error

	// Unlike 取消点赞动态
	Unlike(ctx context.Context, postID, userID string) error

	// IsLiked 检查是否已点赞
	IsLiked(ctx context.Context, postID, userID string) (bool, error)

	// GetLikedPostIDs 批量获取动态的点赞状态
	// 返回 map[postID]bool，表示每个动态是否被该用户点赞
	GetLikedPostIDs(ctx context.Context, userID string, postIDs []string) (map[string]bool, error)

	// IncrementLikeCount 增加点赞数
	IncrementLikeCount(ctx context.Context, postID string) error

	// DecrementLikeCount 减少点赞数
	DecrementLikeCount(ctx context.Context, postID string) error

	// IncrementCommentCount 增加评论数
	IncrementCommentCount(ctx context.Context, postID string) error

	// DecrementCommentCount 减少评论数
	DecrementCommentCount(ctx context.Context, postID string) error

	// RunInTransaction 在事务中执行操作
	RunInTransaction(ctx context.Context, fn func(context.Context) error) error

	// Health 健康检查
	Health(ctx context.Context) error
}
