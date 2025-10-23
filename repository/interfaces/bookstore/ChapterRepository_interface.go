package bookstore

import (
	"Qingyu_backend/models/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterRepository 章节仓储接口
type ChapterRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*bookstore.Chapter, primitive.ObjectID]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// 章节特定查询方法
	GetByBookID(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error)
	GetByBookIDAndChapterNum(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error)
	GetByTitle(ctx context.Context, title string, limit, offset int) ([]*bookstore.Chapter, error)
	GetFreeChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error)
	GetPaidChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error)
	GetPublishedChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error)
	GetChapterRange(ctx context.Context, bookID primitive.ObjectID, startChapter, endChapter int) ([]*bookstore.Chapter, error)

	// 搜索方法
	Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Chapter, error)
	SearchByFilter(ctx context.Context, filter *ChapterFilter) ([]*bookstore.Chapter, error)

	// 统计方法
	CountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	CountFreeChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	CountPaidChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	CountPublishedChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error)
	GetTotalWordCount(ctx context.Context, bookID primitive.ObjectID) (int64, error)

	// 章节排序和导航
	GetPreviousChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error)
	GetNextChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error)
	GetFirstChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error)
	GetLastChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error)

	// 批量操作
	BatchUpdatePrice(ctx context.Context, chapterIDs []primitive.ObjectID, price float64) error
	BatchDelete(ctx context.Context, chapterIDs []primitive.ObjectID) error
	BatchUpdateFreeStatus(ctx context.Context, chapterIDs []primitive.ObjectID, isFree bool) error
	BatchUpdatePublishTime(ctx context.Context, chapterIDs []primitive.ObjectID, publishTime time.Time) error
	DeleteByBookID(ctx context.Context, bookID primitive.ObjectID) error

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// ChapterFilter 章节过滤器
type ChapterFilter struct {
	BookID          *primitive.ObjectID `json:"book_id,omitempty"`
	Title           string              `json:"title,omitempty"`
	IsFree          *bool               `json:"is_free,omitempty"`
	IsPublished     *bool               `json:"is_published,omitempty"`
	MinPrice        *float64            `json:"min_price,omitempty"`
	MaxPrice        *float64            `json:"max_price,omitempty"`
	MinWordCount    *int                `json:"min_word_count,omitempty"`
	MaxWordCount    *int                `json:"max_word_count,omitempty"`
	StartChapter    *int                `json:"start_chapter,omitempty"`
	EndChapter      *int                `json:"end_chapter,omitempty"`
	PublishedAfter  *time.Time          `json:"published_after,omitempty"`
	PublishedBefore *time.Time          `json:"published_before,omitempty"`
	MinChapterNum   *int                `json:"min_chapter_num,omitempty"`
	MaxChapterNum   *int                `json:"max_chapter_num,omitempty"`
	Limit           int                 `json:"limit,omitempty"`
	Offset          int                 `json:"offset,omitempty"`
	SortBy          string              `json:"sort_by,omitempty"`
	SortOrder       string              `json:"sort_order,omitempty"`
}
