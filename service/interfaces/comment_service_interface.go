package interfaces

import (
	"context"

	"Qingyu_backend/models/reader"
)

// CommentService 评论服务接口
type CommentService interface {
	// 评论管理
	PublishComment(ctx context.Context, userID, bookID, chapterID, content string, rating int) (*reader.Comment, error)
	ReplyComment(ctx context.Context, userID, parentCommentID, content string) (*reader.Comment, error)
	GetCommentList(ctx context.Context, bookID string, sortBy string, page, size int) ([]*reader.Comment, int64, error)
	GetCommentDetail(ctx context.Context, commentID string) (*reader.Comment, error)
	UpdateComment(ctx context.Context, userID, commentID, content string) error
	DeleteComment(ctx context.Context, userID, commentID string) error

	// 点赞（将与LikeService集成）
	LikeComment(ctx context.Context, userID, commentID string) error
	UnlikeComment(ctx context.Context, userID, commentID string) error

	// 审核
	AutoReviewComment(ctx context.Context, comment *reader.Comment) (string, string, error) // status, reason, error

	// 统计
	GetBookCommentStats(ctx context.Context, bookID string) (map[string]interface{}, error)
	GetUserComments(ctx context.Context, userID string, page, size int) ([]*reader.Comment, int64, error)
}
