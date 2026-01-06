package reader

import (
	"Qingyu_backend/models/community"
	"context"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, comment *community.Comment) error
	GetByID(ctx context.Context, id string) (*community.Comment, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 查询操作
	GetCommentsByBookID(ctx context.Context, bookID string, page, size int) ([]*community.Comment, int64, error)
	GetCommentsByUserID(ctx context.Context, userID string, page, size int) ([]*community.Comment, int64, error)
	GetRepliesByCommentID(ctx context.Context, commentID string) ([]*community.Comment, error)
	GetCommentsByChapterID(ctx context.Context, chapterID string, page, size int) ([]*community.Comment, int64, error)

	// 排序查询
	GetCommentsByBookIDSorted(ctx context.Context, bookID string, sortBy string, page, size int) ([]*community.Comment, int64, error)

	// 审核相关
	UpdateCommentStatus(ctx context.Context, id, status, reason string) error
	GetPendingComments(ctx context.Context, page, size int) ([]*community.Comment, int64, error)

	// 统计操作
	IncrementLikeCount(ctx context.Context, id string) error
	DecrementLikeCount(ctx context.Context, id string) error
	IncrementReplyCount(ctx context.Context, id string) error
	DecrementReplyCount(ctx context.Context, id string) error
	GetBookRatingStats(ctx context.Context, bookID string) (map[string]interface{}, error)
	GetCommentCount(ctx context.Context, bookID string) (int64, error)

	// 批量操作
	GetCommentsByIDs(ctx context.Context, ids []string) ([]*community.Comment, error)
	DeleteCommentsByBookID(ctx context.Context, bookID string) error

	// 健康检查
	Health(ctx context.Context) error
}
