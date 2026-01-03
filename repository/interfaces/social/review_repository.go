package social

import (
	"context"

	"Qingyu_backend/models/social"
)

// ReviewRepository 书评仓储接口
type ReviewRepository interface {
	// ========== 书评管理 ==========

	// CreateReview 创建书评
	CreateReview(ctx context.Context, review *social.Review) error

	// GetReviewByID 根据ID获取书评
	GetReviewByID(ctx context.Context, reviewID string) (*social.Review, error)

	// GetReviewsByBook 获取书籍的书评列表
	GetReviewsByBook(ctx context.Context, bookID string, page, size int) ([]*social.Review, int64, error)

	// GetReviewsByUser 获取用户的书评列表
	GetReviewsByUser(ctx context.Context, userID string, page, size int) ([]*social.Review, int64, error)

	// GetPublicReviews 获取公开书评列表
	GetPublicReviews(ctx context.Context, page, size int) ([]*social.Review, int64, error)

	// GetReviewsByRating 根据评分获取书评
	GetReviewsByRating(ctx context.Context, bookID string, rating int, page, size int) ([]*social.Review, int64, error)

	// UpdateReview 更新书评
	UpdateReview(ctx context.Context, reviewID string, updates map[string]interface{}) error

	// DeleteReview 删除书评
	DeleteReview(ctx context.Context, reviewID string) error

	// ========== 书评点赞 ==========

	// CreateReviewLike 创建书评点赞
	CreateReviewLike(ctx context.Context, reviewLike *social.ReviewLike) error

	// DeleteReviewLike 删除书评点赞
	DeleteReviewLike(ctx context.Context, reviewID, userID string) error

	// GetReviewLike 获取书评点赞记录
	GetReviewLike(ctx context.Context, reviewID, userID string) (*social.ReviewLike, error)

	// IsReviewLiked 检查是否已点赞
	IsReviewLiked(ctx context.Context, reviewID, userID string) (bool, error)

	// GetReviewLikes 获取书评点赞列表
	GetReviewLikes(ctx context.Context, reviewID string, page, size int) ([]*social.ReviewLike, int64, error)

	// IncrementReviewLikeCount 增加书评点赞数
	IncrementReviewLikeCount(ctx context.Context, reviewID string) error

	// DecrementReviewLikeCount 减少书评点赞数
	DecrementReviewLikeCount(ctx context.Context, reviewID string) error

	// ========== 统计 ==========

	// GetAverageRating 获取书籍平均评分
	GetAverageRating(ctx context.Context, bookID string) (float64, error)

	// GetRatingDistribution 获取评分分布
	GetRatingDistribution(ctx context.Context, bookID string) (map[int]int64, error)

	// CountReviews 统计书评数
	CountReviews(ctx context.Context, bookID string) (int64, error)

	// CountUserReviews 统计用户书评数
	CountUserReviews(ctx context.Context, userID string) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
