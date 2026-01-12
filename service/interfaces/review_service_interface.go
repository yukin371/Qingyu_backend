package interfaces

import (
	socialModel "Qingyu_backend/models/social"
	"context"
)

// ReviewService 书评服务接口
type ReviewService interface {
	// =========================
	// 书评管理
	// =========================
	CreateReview(ctx context.Context, bookID, userID, userName, userAvatar, title, content string, rating int, isSpoiler, isPublic bool) (*socialModel.Review, error)
	GetReviews(ctx context.Context, bookID string, page, size int) ([]*socialModel.Review, int64, error)
	GetReviewByID(ctx context.Context, reviewID string) (*socialModel.Review, error)
	UpdateReview(ctx context.Context, userID, reviewID string, updates map[string]interface{}) error
	DeleteReview(ctx context.Context, userID, reviewID string) error

	// =========================
	// 书评点赞
	// =========================
	LikeReview(ctx context.Context, userID, reviewID string) error
	UnlikeReview(ctx context.Context, userID, reviewID string) error
}
