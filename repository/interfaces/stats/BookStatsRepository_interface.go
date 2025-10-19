package stats

import (
	"Qingyu_backend/models/stats"
	"Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"time"
)

// BookStatsRepository 作品统计Repository接口
type BookStatsRepository interface {
	// 基础CRUD
	Create(ctx context.Context, bookStats *stats.BookStats) error
	GetByID(ctx context.Context, id string) (*stats.BookStats, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	
	// 查询方法
	GetByBookID(ctx context.Context, bookID string) (*stats.BookStats, error)
	GetByAuthorID(ctx context.Context, authorID string, limit, offset int64) ([]*stats.BookStats, error)
	GetByDateRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.BookStats, error)
	
	// 每日统计
	CreateDailyStats(ctx context.Context, dailyStats *stats.BookStatsDaily) error
	GetDailyStats(ctx context.Context, bookID string, date time.Time) (*stats.BookStatsDaily, error)
	GetDailyStatsRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*stats.BookStatsDaily, error)
	
	// 收入统计
	GetRevenueBreakdown(ctx context.Context, bookID string, startDate, endDate time.Time) (*stats.RevenueBreakdown, error)
	CalculateTotalRevenue(ctx context.Context, bookID string) (float64, error)
	CalculateRevenueByType(ctx context.Context, bookID string, revenueType string) (float64, error)
	
	// 热门章节
	GetTopChapters(ctx context.Context, bookID string) (*stats.TopChapters, error)
	
	// 趋势分析
	AnalyzeViewTrend(ctx context.Context, bookID string, days int) (string, error)
	AnalyzeRevenueTrend(ctx context.Context, bookID string, days int) (string, error)
	
	// 聚合统计
	CalculateAvgCompletionRate(ctx context.Context, bookID string) (float64, error)
	CalculateAvgDropOffRate(ctx context.Context, bookID string) (float64, error)
	CalculateAvgReadingDuration(ctx context.Context, bookID string) (float64, error)
	
	// 排名查询
	GetTopBooksByViews(ctx context.Context, limit int) ([]*stats.BookStats, error)
	GetTopBooksByRevenue(ctx context.Context, limit int) ([]*stats.BookStats, error)
	GetTopBooksByCompletion(ctx context.Context, limit int) ([]*stats.BookStats, error)
	
	// 批量操作
	BatchCreate(ctx context.Context, bookStats []*stats.BookStats) error
	BatchUpdate(ctx context.Context, updates []map[string]interface{}) error
	
	// 统计方法
	Count(ctx context.Context, filter infrastructure.Filter) (int64, error)
	CountByAuthor(ctx context.Context, authorID string) (int64, error)
	
	// 健康检查
	Health(ctx context.Context) error
}

