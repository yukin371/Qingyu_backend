package writer

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/writer"
)

// CommentRepository 批注仓储接口
type CommentRepository interface {
	// Create 创建批注
	Create(ctx context.Context, comment *writer.DocumentComment) error

	// GetByID 根据ID获取批注
	GetByID(ctx context.Context, id primitive.ObjectID) (*writer.DocumentComment, error)

	// Update 更新批注
	Update(ctx context.Context, id primitive.ObjectID, comment *writer.DocumentComment) error

	// Delete 删除批注（软删除）
	Delete(ctx context.Context, id primitive.ObjectID) error

	// HardDelete 硬删除批注
	HardDelete(ctx context.Context, id primitive.ObjectID) error

	// List 查询批注列表（支持分页和筛选）
	List(ctx context.Context, filter *writer.CommentFilter, page, pageSize int) ([]*writer.DocumentComment, int64, error)

	// GetByDocument 获取文档的所有批注
	GetByDocument(ctx context.Context, documentID primitive.ObjectID, includeResolved bool) ([]*writer.DocumentComment, error)

	// GetByChapter 获取章节的所有批注
	GetByChapter(ctx context.Context, chapterID primitive.ObjectID, includeResolved bool) ([]*writer.DocumentComment, error)

	// GetThread 获取批注线程
	GetThread(ctx context.Context, threadID primitive.ObjectID) (*writer.CommentThread, error)

	// GetReplies 获取批注的回复
	GetReplies(ctx context.Context, parentID primitive.ObjectID) ([]*writer.DocumentComment, error)

	// MarkAsResolved 标记为已解决
	MarkAsResolved(ctx context.Context, id primitive.ObjectID, resolvedBy primitive.ObjectID) error

	// MarkAsUnresolved 标记为未解决
	MarkAsUnresolved(ctx context.Context, id primitive.ObjectID) error

	// GetStats 获取批注统计
	GetStats(ctx context.Context, documentID primitive.ObjectID) (*writer.CommentStats, error)

	// GetUserComments 获取用户的批注
	GetUserComments(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*writer.DocumentComment, int64, error)

	// BatchDelete 批量删除批注
	BatchDelete(ctx context.Context, ids []primitive.ObjectID) error

	// Search 搜索批注
	Search(ctx context.Context, keyword string, documentID primitive.ObjectID, page, pageSize int) ([]*writer.DocumentComment, int64, error)
}
