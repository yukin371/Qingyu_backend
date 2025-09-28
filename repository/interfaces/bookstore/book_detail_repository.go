package bookstore

import (
	"Qingyu_backend/models/reading/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookDetailRepository 书籍详情仓储接口
type BookDetailRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.BookDetail, primitive.ObjectID]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// 书籍详情特定查询方法
	GetByTitle(ctx context.Context, title string) (*bookstore.BookDetail, error)
	GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.BookDetail, error)
	GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore.BookDetail, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*bookstore.BookDetail, error)
	GetByStatus(ctx context.Context, status bookstore.BookStatus, limit, offset int) ([]*bookstore.BookDetail, error)
	GetByISBN(ctx context.Context, isbn string) (*bookstore.BookDetail, error)
	GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstore.BookDetail, error)

	// 搜索方法
	Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.BookDetail, error)
	SearchByFilter(ctx context.Context, filter *BookDetailFilter) ([]*bookstore.BookDetail, error)

	// 统计方法
	CountByCategory(ctx context.Context, category string) (int64, error)
	CountByAuthor(ctx context.Context, author string) (int64, error)
	CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error)
	CountByTags(ctx context.Context, tags []string) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore.BookStatus) error
	BatchUpdateCategories(ctx context.Context, bookIDs []primitive.ObjectID, categories []string) error
	BatchUpdateTags(ctx context.Context, bookIDs []primitive.ObjectID, tags []string) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// BookDetailFilter 书籍详情过滤器
type BookDetailFilter struct {
	Title       string                   `json:"title,omitempty"`
	Author      string                   `json:"author,omitempty"`
	AuthorID    *primitive.ObjectID      `json:"author_id,omitempty"`
	Categories  []string                 `json:"categories,omitempty"`
	Tags        []string                 `json:"tags,omitempty"`
	Status      *bookstore.BookStatus    `json:"status,omitempty"`
	IsFree      *bool                    `json:"is_free,omitempty"`
	MinPrice    *float64                 `json:"min_price,omitempty"`
	MaxPrice    *float64                 `json:"max_price,omitempty"`
	MinWordCount *int64                  `json:"min_word_count,omitempty"`
	MaxWordCount *int64                  `json:"max_word_count,omitempty"`
	Publisher   string                   `json:"publisher,omitempty"`
	Limit       int                      `json:"limit,omitempty"`
	Offset      int                      `json:"offset,omitempty"`
	SortBy      string                   `json:"sort_by,omitempty"`
	SortOrder   string                   `json:"sort_order,omitempty"`
}