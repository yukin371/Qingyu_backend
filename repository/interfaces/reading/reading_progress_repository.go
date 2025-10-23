package reading

import (
	"Qingyu_backend/models/reader"
	"context"
	"time"
)

// ReadingProgressRepository 阅读进度仓储接口
type ReadingProgressRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, progress *reader.ReadingProgress) error
	GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 进度查询
	GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error)
	GetByUser(ctx context.Context, userID string) ([]*reader.ReadingProgress, error)
	GetRecentReadingByUser(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error)

	// 进度保存和更新
	SaveProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error
	UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error
	UpdateLastReadAt(ctx context.Context, userID, bookID string) error

	// 批量更新
	BatchUpdateProgress(ctx context.Context, progresses []*reader.ReadingProgress) error

	// 统计查询
	GetTotalReadingTime(ctx context.Context, userID string) (int64, error)
	GetReadingTimeByBook(ctx context.Context, userID, bookID string) (int64, error)
	GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error)
	CountReadingBooks(ctx context.Context, userID string) (int64, error)

	// 阅读记录
	GetReadingHistory(ctx context.Context, userID string, limit, offset int) ([]*reader.ReadingProgress, error)
	GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error)
	GetFinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error)

	// 数据同步
	SyncProgress(ctx context.Context, userID string, progresses []*reader.ReadingProgress) error
	GetProgressesByUser(ctx context.Context, userID string, updatedAfter time.Time) ([]*reader.ReadingProgress, error)

	// 清理操作
	DeleteOldProgress(ctx context.Context, beforeTime time.Time) error
	DeleteByBook(ctx context.Context, bookID string) error

	// 健康检查
	Health(ctx context.Context) error
}
