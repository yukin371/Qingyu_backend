package stats

import (
	"Qingyu_backend/models/stats"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"
)

// ReaderBehaviorRepository 读者行为Repository接口
type ReaderBehaviorRepository interface {
	// 基础CRUD
	Create(ctx context.Context, behavior *stats.ReaderBehavior) error
	GetByID(ctx context.Context, id string) (*stats.ReaderBehavior, error)
	Delete(ctx context.Context, id string) error

	// 查询方法
	GetByUserID(ctx context.Context, userID string, limit, offset int64) ([]*stats.ReaderBehavior, error)
	GetByBookID(ctx context.Context, bookID string, limit, offset int64) ([]*stats.ReaderBehavior, error)
	GetByChapterID(ctx context.Context, chapterID string, limit, offset int64) ([]*stats.ReaderBehavior, error)
	GetByBehaviorType(ctx context.Context, behaviorType string, limit, offset int64) ([]*stats.ReaderBehavior, error)

	// 时间范围查询
	GetByDateRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.ReaderBehavior, error)

	// 聚合统计
	CountUniqueReaders(ctx context.Context, bookID string) (int64, error)
	CountUniqueReadersByChapter(ctx context.Context, chapterID string) (int64, error)
	CalculateAvgReadTime(ctx context.Context, chapterID string) (float64, error)
	CalculateCompletionRate(ctx context.Context, chapterID string) (float64, error)
	CalculateDropOffRate(ctx context.Context, chapterID string) (float64, error)

	// 会话相关
	CreateSession(ctx context.Context, session *stats.ReadingSession) error
	GetSessionByID(ctx context.Context, sessionID string) (*stats.ReadingSession, error)
	GetUserSessions(ctx context.Context, userID string, limit int) ([]*stats.ReadingSession, error)

	// 留存数据
	CreateRetention(ctx context.Context, retention *stats.ReaderRetention) error
	GetRetentionByBookID(ctx context.Context, bookID string) (*stats.ReaderRetention, error)
	CalculateRetention(ctx context.Context, bookID string, days int) (float64, error)

	// 批量操作
	BatchCreate(ctx context.Context, behaviors []*stats.ReaderBehavior) error

	// 统计方法
	Count(ctx context.Context, filter infrastructure.Filter) (int64, error)
	CountByBehaviorType(ctx context.Context, behaviorType string) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
