package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookDetailRepository 书籍详情仓储接口 - 专注于书籍详情页面的完整信息管理
// 用于书籍详情页面、章节管理、统计数据等详细场景
type BookDetailRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore2.BookDetail, primitive.ObjectID]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// 详情查询方法 - 用于详情页面
	GetByTitle(ctx context.Context, title string) (*bookstore2.BookDetail, error)
	GetByISBN(ctx context.Context, isbn string) (*bookstore2.BookDetail, error)
	GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore2.BookDetail, error)
	GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore2.BookDetail, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*bookstore2.BookDetail, error)
	GetByStatus(ctx context.Context, status bookstore2.BookStatus, limit, offset int) ([]*bookstore2.BookDetail, error)
	GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstore2.BookDetail, error)
	GetByPublisher(ctx context.Context, publisher string, limit, offset int) ([]*bookstore2.BookDetail, error)
	GetByBookID(ctx context.Context, bookID primitive.ObjectID) (*bookstore2.BookDetail, error)
	GetByBookIDs(ctx context.Context, bookIDs []primitive.ObjectID) ([]*bookstore2.BookDetail, error)
	UpdateAuthor(ctx context.Context, bookID primitive.ObjectID, authorID primitive.ObjectID, authorName string) error
	GetSimilarBooks(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstore2.BookDetail, error)

	// 搜索方法 - 用于详细搜索
	Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore2.BookDetail, error)
	SearchByFilter(ctx context.Context, filter *BookDetailFilter) ([]*bookstore2.BookDetail, error)

	// 统计方法 - 用于详情统计
	CountByCategory(ctx context.Context, category string) (int64, error)
	CountByAuthor(ctx context.Context, author string) (int64, error)
	CountByStatus(ctx context.Context, status bookstore2.BookStatus) (int64, error)
	CountByTags(ctx context.Context, tags []string) (int64, error)
	CountByPublisher(ctx context.Context, publisher string) (int64, error)

	// 统计数据更新 - 用于详情页面交互
	IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	DecrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error
	IncrementShareCount(ctx context.Context, bookID primitive.ObjectID) error
	UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64, ratingCount int64) error
	UpdateLastChapter(ctx context.Context, bookID primitive.ObjectID, chapterTitle string) error

	// 批量操作 - 用于管理
	BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore2.BookStatus) error
	BatchUpdatePublisher(ctx context.Context, bookIDs []primitive.ObjectID, publisher string) error
	BatchUpdateTags(ctx context.Context, bookIDs []primitive.ObjectID, tags []string) error
	BatchUpdateCategories(ctx context.Context, bookIDs []primitive.ObjectID, categoryIDs []string) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BookDetailFilter 书籍详情筛选条件 - 适用于网络小说平台
type BookDetailFilter struct {
	Title          string                `json:"title,omitempty"`
	Author         string                `json:"author,omitempty"`
	AuthorID       *primitive.ObjectID   `json:"author_id,omitempty"`
	CategoryIDs    []primitive.ObjectID  `json:"category_ids,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Status         *bookstore2.BookStatus `json:"status,omitempty"`
	IsFree         *bool                  `json:"is_free,omitempty"`
	MinPrice       *float64              `json:"min_price,omitempty"`
	MaxPrice       *float64              `json:"max_price,omitempty"`
	MinRating      *float64              `json:"min_rating,omitempty"`
	MaxRating      *float64              `json:"max_rating,omitempty"`
	MinWordCount   *int64                `json:"min_word_count,omitempty"`
	MaxWordCount   *int64                `json:"max_word_count,omitempty"`
	SerializedFrom *time.Time            `json:"serialized_from,omitempty"` // 开始连载时间范围
	SerializedTo   *time.Time            `json:"serialized_to,omitempty"`   // 开始连载时间范围
	CompletedFrom  *time.Time            `json:"completed_from,omitempty"`  // 完结时间范围
	CompletedTo    *time.Time            `json:"completed_to,omitempty"`    // 完结时间范围
	CreatedAtFrom  *time.Time            `json:"created_at_from,omitempty"`
	CreatedAtTo    *time.Time            `json:"created_at_to,omitempty"`
	UpdatedAtFrom  *time.Time            `json:"updated_at_from,omitempty"`
	UpdatedAtTo    *time.Time            `json:"updated_at_to,omitempty"`
	SortBy         string                `json:"sort_by,omitempty"`    // created_at, updated_at, serialized_at, rating, word_count, view_count, like_count, collect_count
	SortOrder      string                `json:"sort_order,omitempty"` // asc, desc
}
