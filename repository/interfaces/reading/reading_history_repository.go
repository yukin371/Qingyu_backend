package reading

import (
	"context"
	"time"

	"Qingyu_backend/models/reader"
)

// ReadingHistoryRepository 阅读历史Repository接口
type ReadingHistoryRepository interface {
	// 基础CRUD
	Create(ctx context.Context, history *reader.ReadingHistory) error
	GetByID(ctx context.Context, id string) (*reader.ReadingHistory, error)
	Delete(ctx context.Context, id string) error

	// 查询方法
	GetUserHistories(ctx context.Context, userID string, limit, offset int) ([]*reader.ReadingHistory, error)
	GetUserHistoriesByBook(ctx context.Context, userID, bookID string, limit, offset int) ([]*reader.ReadingHistory, error)
	GetUserHistoriesByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time, limit, offset int) ([]*reader.ReadingHistory, error)

	// 统计方法
	CountUserHistories(ctx context.Context, userID string) (int64, error)
	GetUserReadingStats(ctx context.Context, userID string) (*ReadingStats, error)
	GetUserDailyReadingStats(ctx context.Context, userID string, days int) ([]DailyReadingStats, error)

	// 批量操作
	DeleteUserHistories(ctx context.Context, userID string) error
	DeleteOldHistories(ctx context.Context, beforeDate time.Time) (int64, error)

	// 健康检查
	Health(ctx context.Context) error
}

// ReadingStats 阅读统计
type ReadingStats struct {
	TotalDuration    int       `json:"total_duration"`     // 总阅读时长（秒）
	TotalBooks       int       `json:"total_books"`        // 阅读书籍数
	TotalChapters    int       `json:"total_chapters"`     // 阅读章节数
	LastReadTime     time.Time `json:"last_read_time"`     // 最后阅读时间
	AvgDailyDuration int       `json:"avg_daily_duration"` // 日均阅读时长（秒）
}

// DailyReadingStats 每日阅读统计
type DailyReadingStats struct {
	Date     string `json:"date"`     // 日期 (YYYY-MM-DD)
	Duration int    `json:"duration"` // 当日阅读时长（秒）
	Books    int    `json:"books"`    // 当日阅读书籍数
}
