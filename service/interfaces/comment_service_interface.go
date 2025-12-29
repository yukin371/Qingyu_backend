package interfaces

import (
	"Qingyu_backend/models/community"
	"context"
)

// CommentService 评论服务接口
type CommentService interface {
	// 评论管理
	PublishComment(ctx context.Context, userID, bookID, chapterID, content string, rating int) (*community.Comment, error)
	ReplyComment(ctx context.Context, userID, parentCommentID, content string) (*community.Comment, error)
	GetCommentList(ctx context.Context, bookID string, sortBy string, page, size int) ([]*community.Comment, int64, error)
	GetCommentDetail(ctx context.Context, commentID string) (*community.Comment, error)
	UpdateComment(ctx context.Context, userID, commentID, content string) error
	DeleteComment(ctx context.Context, userID, commentID string) error

	// 点赞（将与LikeService集成）
	LikeComment(ctx context.Context, userID, commentID string) error
	UnlikeComment(ctx context.Context, userID, commentID string) error

	// 审核
	AutoReviewComment(ctx context.Context, comment *community.Comment) (string, string, error) // status, reason, error

	// 统计
	GetBookCommentStats(ctx context.Context, bookID string) (map[string]interface{}, error)
	GetUserComments(ctx context.Context, userID string, page, size int) ([]*community.Comment, int64, error)

	// 高级查询
	GetCommentThread(ctx context.Context, commentID string) (*community.CommentThread, error)                     // 获取评论树状结构
	GetTopComments(ctx context.Context, bookID string, limit int) ([]*community.Comment, error)                   // 获取热门评论
	GetCommentReplies(ctx context.Context, commentID string, page, size int) ([]*community.Comment, int64, error) // 获取评论回复列表
}
