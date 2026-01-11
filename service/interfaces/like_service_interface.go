package interfaces

import (
	"Qingyu_backend/models/social"
	"context"
)

// LikeService 点赞服务接口
type LikeService interface {
	// 书籍点赞
	LikeBook(ctx context.Context, userID, bookID string) error
	UnlikeBook(ctx context.Context, userID, bookID string) error
	GetBookLikeCount(ctx context.Context, bookID string) (int64, error)
	IsBookLiked(ctx context.Context, userID, bookID string) (bool, error)

	// 评论点赞
	LikeComment(ctx context.Context, userID, commentID string) error
	UnlikeComment(ctx context.Context, userID, commentID string) error

	// 用户点赞列表
	GetUserLikedBooks(ctx context.Context, userID string, page, size int) ([]*social.Like, int64, error)
	GetUserLikedComments(ctx context.Context, userID string, page, size int) ([]*social.Like, int64, error)

	// 批量查询
	GetBooksLikeCount(ctx context.Context, bookIDs []string) (map[string]int64, error)
	GetUserLikeStatus(ctx context.Context, userID string, bookIDs []string) (map[string]bool, error)

	// 统计
	GetUserLikeStats(ctx context.Context, userID string) (map[string]interface{}, error)
}
