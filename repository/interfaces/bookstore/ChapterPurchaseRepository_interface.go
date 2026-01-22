package bookstore

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"time"
)

// ChapterPurchaseRepository 章节购买记录仓储接口
type ChapterPurchaseRepository interface {
	// Health 健康检查
	Health(ctx context.Context) error

	// 单章购买记录
	Create(ctx context.Context, purchase *bookstore.ChapterPurchase) error
	GetByID(ctx context.Context, id string) (*bookstore.ChapterPurchase, error)
	GetByUserAndChapter(ctx context.Context, userID, chapterID string) (*bookstore.ChapterPurchase, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 查询用户的购买记录
	GetByUser(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error)
	GetByUserAndBook(ctx context.Context, userID, bookID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error)

	// 批量购买记录
	CreateBatch(ctx context.Context, batch *bookstore.ChapterPurchaseBatch) error
	GetBatchByID(ctx context.Context, id string) (*bookstore.ChapterPurchaseBatch, error)
	GetBatchesByUser(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error)
	GetBatchesByUserAndBook(ctx context.Context, userID, bookID string, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error)

	// 全书购买记录
	CreateBookPurchase(ctx context.Context, purchase *bookstore.BookPurchase) error
	GetBookPurchaseByID(ctx context.Context, id string) (*bookstore.BookPurchase, error)
	GetBookPurchaseByUserAndBook(ctx context.Context, userID, bookID string) (*bookstore.BookPurchase, error)
	GetBookPurchasesByUser(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.BookPurchase, int64, error)

	// 权限检查
	CheckUserPurchasedChapter(ctx context.Context, userID, chapterID string) (bool, error)
	CheckUserPurchasedBook(ctx context.Context, userID, bookID string) (bool, error)
	GetPurchasedChapterIDs(ctx context.Context, userID, bookID string) ([]string, error)

	// 统计
	CountByUser(ctx context.Context, userID string) (int64, error)
	CountByUserAndBook(ctx context.Context, userID, bookID string) (int64, error)
	GetTotalSpentByUser(ctx context.Context, userID string) (float64, error)
	GetTotalSpentByUserAndBook(ctx context.Context, userID, bookID string) (float64, error)

	// 时间范围查询
	GetPurchasesByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*bookstore.ChapterPurchase, error)

	// 事务支持
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
