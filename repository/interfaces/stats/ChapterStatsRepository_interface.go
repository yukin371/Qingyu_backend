package stats

import (
	"Qingyu_backend/models/stats"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"
)

// ChapterStatsRepository 章节统计Repository接口
type ChapterStatsRepository interface {
	// 基础CRUD
	Create(ctx context.Context, chapterStats *stats.ChapterStats) error
	GetByID(ctx context.Context, id string) (*stats.ChapterStats, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	
	// 查询方法
	GetByChapterID(ctx context.Context, chapterID string) (*stats.ChapterStats, error)
	GetByBookID(ctx context.Context, bookID string, limit, offset int64) ([]*stats.ChapterStats, error)
	GetByDateRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.ChapterStats, error)
	
	// 聚合查询
	GetChapterStatsAggregate(ctx context.Context, bookID string) ([]*stats.ChapterStatsAggregate, error)
	GetTopViewedChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error)
	GetTopRevenueChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error)
	GetLowestCompletionChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error)
	GetHighestDropOffChapters(ctx context.Context, bookID string, limit int) ([]*stats.ChapterStatsAggregate, error)
	
	// 热力图数据
	GenerateHeatmap(ctx context.Context, bookID string) ([]*stats.HeatmapPoint, error)
	
	// 时间范围统计
	GetTimeRangeStats(ctx context.Context, bookID string, startDate, endDate time.Time) (*stats.TimeRangeStats, error)
	
	// 批量操作
	BatchCreate(ctx context.Context, chapterStats []*stats.ChapterStats) error
	BatchUpdate(ctx context.Context, updates []map[string]interface{}) error
	
	// 统计方法
	Count(ctx context.Context, filter infrastructure.Filter) (int64, error)
	CountByBook(ctx context.Context, bookID string) (int64, error)
	
	// 健康检查
	Health(ctx context.Context) error
}

