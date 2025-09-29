package bookstore

import (
	"Qingyu_backend/models/reading/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookStatisticsRepository 书籍统计仓储接口
type BookStatisticsRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.BookStatistics, primitive.ObjectID]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// 统计特定查询方法
	GetByBookID(ctx context.Context, bookID primitive.ObjectID) (*bookstore.BookStatistics, error)
	GetTopViewed(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error)
	GetTopFavorited(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error)
	GetTopRated(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error)
	GetHottest(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error)
	GetMostCommented(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error)
	GetMostShared(ctx context.Context, limit, offset int) ([]*bookstore.BookStatistics, error)
	GetTrendingBooks(ctx context.Context, days int, limit, offset int) ([]*bookstore.BookStatistics, error)

	// 统计操作
	IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementFavoriteCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementShareCount(ctx context.Context, bookID primitive.ObjectID) error
	UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating int) error
	RemoveRating(ctx context.Context, bookID primitive.ObjectID, rating int) error
	UpdateHotScore(ctx context.Context, bookID primitive.ObjectID, hotScore float64) error

	// 批量统计操作
	BatchIncrementViewCount(ctx context.Context, increments map[primitive.ObjectID]int64) error
	BatchUpdateViewCount(ctx context.Context, bookIDs []primitive.ObjectID, increment int64) error
	BatchUpdateHotScore(ctx context.Context, bookIDs []primitive.ObjectID) error
	BatchRecalculateStatistics(ctx context.Context, bookIDs []primitive.ObjectID) error

	// 聚合统计
	GetTotalViews(ctx context.Context) (int64, error)
	GetTotalFavorites(ctx context.Context) (int64, error)
	GetTotalComments(ctx context.Context) (int64, error)
	GetTotalShares(ctx context.Context) (int64, error)
	GetAverageRating(ctx context.Context) (float64, error)
	GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error)

	// 时间范围统计
	GetViewsInRange(ctx context.Context, bookID primitive.ObjectID, days int) (int64, error)
	GetFavoritesInRange(ctx context.Context, bookID primitive.ObjectID, days int) (int64, error)
	GetCommentsInRange(ctx context.Context, bookID primitive.ObjectID, days int) (int64, error)
	GetStatisticsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*bookstore.BookStatistics, error)

	// 搜索功能
	Search(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.BookStatistics, int64, error)
	SearchByFilter(ctx context.Context, filter *BookStatisticsFilter, page, pageSize int) ([]*bookstore.BookStatistics, int64, error)

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BookStatisticsFilter 书籍统计过滤器
type BookStatisticsFilter struct {
	BookID           *primitive.ObjectID `json:"book_id,omitempty"`
	MinViewCount     *int64              `json:"min_view_count,omitempty"`
	MaxViewCount     *int64              `json:"max_view_count,omitempty"`
	MinFavoriteCount *int64              `json:"min_favorite_count,omitempty"`
	MaxFavoriteCount *int64              `json:"max_favorite_count,omitempty"`
	MinCommentCount  *int64              `json:"min_comment_count,omitempty"`
	MaxCommentCount  *int64              `json:"max_comment_count,omitempty"`
	MinShareCount    *int64              `json:"min_share_count,omitempty"`
	MaxShareCount    *int64              `json:"max_share_count,omitempty"`
	MinAverageRating *float64            `json:"min_average_rating,omitempty"`
	MaxAverageRating *float64            `json:"max_average_rating,omitempty"`
	MinRatingCount   *int64              `json:"min_rating_count,omitempty"`
	MaxRatingCount   *int64              `json:"max_rating_count,omitempty"`
	MinHotScore      *float64            `json:"min_hot_score,omitempty"`
	MaxHotScore      *float64            `json:"max_hot_score,omitempty"`
	MinRating        *float64            `json:"min_rating,omitempty"`
	MaxRating        *float64            `json:"max_rating,omitempty"`
	Limit            int                 `json:"limit,omitempty"`
	Offset           int                 `json:"offset,omitempty"`
	SortBy           string              `json:"sort_by,omitempty"`
	SortOrder        string              `json:"sort_order,omitempty"`
}