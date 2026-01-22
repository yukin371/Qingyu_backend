package bookstore

import (
	"Qingyu_backend/models/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookRatingRepository 书籍评分仓储接口
type BookRatingRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.BookRating, primitive.ObjectID]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// 评分特定查询方法
	GetByBookID(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.BookRating, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*bookstore.BookRating, error)
	GetByBookIDAndUserID(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstore.BookRating, error)
	GetByRating(ctx context.Context, rating int, limit, offset int) ([]*bookstore.BookRating, error)
	GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstore.BookRating, error)
	GetTopRated(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.BookRating, error)
	GetMostLiked(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.BookRating, error)

	// 搜索方法
	Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.BookRating, error)
	SearchByFilter(ctx context.Context, filter *BookRatingFilter) ([]*bookstore.BookRating, error)

	// 统计方法
	CountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	CountByUserID(ctx context.Context, userID primitive.ObjectID) (int64, error)
	CountByRating(ctx context.Context, bookID primitive.ObjectID, rating int) (int64, error)
	GetAverageRating(ctx context.Context, bookID primitive.ObjectID) (float64, error)
	GetRatingDistribution(ctx context.Context, bookID primitive.ObjectID) (map[string]int64, error)
	GetTotalLikes(ctx context.Context, bookID primitive.ObjectID) (int64, error)

	// 评分操作
	IncrementLikes(ctx context.Context, ratingID primitive.ObjectID) error
	DecrementLikes(ctx context.Context, ratingID primitive.ObjectID) error
	UpdateRating(ctx context.Context, ratingID primitive.ObjectID, rating int) error
	UpdateComment(ctx context.Context, ratingID primitive.ObjectID, comment string) error
	UpdateTags(ctx context.Context, ratingID primitive.ObjectID, tags []string) error

	// 批量操作
	BatchDelete(ctx context.Context, ratingIDs []primitive.ObjectID) error
	BatchUpdateTags(ctx context.Context, ratingIDs []primitive.ObjectID, tags []string) error
	BatchDeleteByBookID(ctx context.Context, bookID primitive.ObjectID) error
	BatchDeleteByUserID(ctx context.Context, userID primitive.ObjectID) error
	DeleteByBookID(ctx context.Context, bookID primitive.ObjectID) error
	DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BookRatingFilter 书籍评分过滤器
type BookRatingFilter struct {
	BookID     *primitive.ObjectID `json:"book_id,omitempty"`
	UserID     *primitive.ObjectID `json:"user_id,omitempty"`
	MinRating  *int                `json:"min_rating,omitempty"`
	MaxRating  *int                `json:"max_rating,omitempty"`
	Tags       []string            `json:"tags,omitempty"`
	HasComment *bool               `json:"has_comment,omitempty"`
	MinLikes   *int                `json:"min_likes,omitempty"`
	MaxLikes   *int                `json:"max_likes,omitempty"`
	Limit      int                 `json:"limit,omitempty"`
	Offset     int                 `json:"offset,omitempty"`
	SortBy     string              `json:"sort_by,omitempty"`
	SortOrder  string              `json:"sort_order,omitempty"`
}
