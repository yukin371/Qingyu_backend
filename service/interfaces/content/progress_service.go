package content

import (
	"context"

	"Qingyu_backend/models/dto"
)

// ============================================================================
// 阅读进度服务接口详细定义
// ============================================================================

// ReadingProgressService 阅读进度服务接口
// 提供阅读进度的保存、查询、统计和历史记录管理功能
type ReadingProgressService interface {
	// === 进度管理 ===

	// GetProgress 获取阅读进度
	// 返回用户对指定书籍的阅读进度信息
	GetProgress(ctx context.Context, userID, bookID string) (*dto.ReadingProgressResponse, error)

	// SaveProgress 保存阅读进度
	// 保存用户的阅读位置和进度百分比
	SaveProgress(ctx context.Context, req *dto.SaveProgressRequest) error

	// UpdateReadingTime 更新阅读时长
	// 累加用户的阅读时长
	UpdateReadingTime(ctx context.Context, userID, bookID string, duration int) error

	// DeleteProgress 删除阅读进度
	// 删除用户对指定书籍的阅读进度记录
	DeleteProgress(ctx context.Context, userID, bookID string) error

	// === 书籍状态管理 ===

	// UpdateBookStatus 更新书籍状态
	// 更新书籍的阅读状态（在读/想读/读完）
	UpdateBookStatus(ctx context.Context, userID, bookID, status string) error

	// BatchUpdateBookStatus 批量更新书籍状态
	// 批量更新多本书籍的阅读状态
	BatchUpdateBookStatus(ctx context.Context, userID string, bookIDs []string, status string) error

	// GetBooksByStatus 根据状态获取书籍
	// 按阅读状态筛选书籍列表
	GetBooksByStatus(ctx context.Context, userID, status string, page, pageSize int) (*dto.RecentBooksResponse, error)

	// === 阅读统计 ===

	// GetReadingStats 获取阅读统计
	// 返回用户的阅读统计数据，包括总时长、总字数、已读书籍数等
	GetReadingStats(ctx context.Context, userID string) (*dto.ReadingStatsResponse, error)

	// GetTotalReadingTime 获取总阅读时长
	// 返回用户的累计阅读时长（秒）
	GetTotalReadingTime(ctx context.Context, userID string) (int64, error)

	// GetReadingTimeByPeriod 获取时间段阅读时长
	// 返回指定时间范围内的阅读时长
	GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime int64) (int64, error)

	// === 阅读历史 ===

	// GetRecentBooks 获取最近阅读的书籍
	// 返回用户最近阅读的书籍列表，按最后阅读时间排序
	GetRecentBooks(ctx context.Context, userID string, limit int) (*dto.RecentBooksResponse, error)

	// GetReadingHistory 获取阅读历史
	// 返回用户的阅读历史记录，支持分页
	GetReadingHistory(ctx context.Context, userID string, page, pageSize int) (*dto.ReadingHistoryResponse, error)

	// GetUnfinishedBooks 获取未读完的书籍
	// 返回所有未读完的书籍列表
	GetUnfinishedBooks(ctx context.Context, userID string) ([]*dto.BookProgressInfo, error)

	// GetFinishedBooks 获取已读完的书籍
	// 返回所有已读完的书籍列表
	GetFinishedBooks(ctx context.Context, userID string) ([]*dto.BookProgressInfo, error)

	// === 阅读分析 ===

	// GetReadingTrends 获取阅读趋势
	// 返回用户的阅读趋势数据（按天/周/月统计）
	GetReadingTrends(ctx context.Context, userID string, period string, days int) (*dto.ReadingTrendsResponse, error)

	// GetReadingStreak 获取连续阅读天数
	// 返回用户的连续阅读天数记录
	GetReadingStreak(ctx context.Context, userID string) (*dto.ReadingStreakResponse, error)

	// GetLongestBooks 获取阅读字数最多的书籍
	// 返回用户阅读字数最多的书籍列表
	GetLongestBooks(ctx context.Context, userID string, limit int) ([]*dto.BookProgressInfo, error)

	// === 同步相关 ===

	// GetProgressSyncData 获取进度同步数据
	// 返回需要同步的进度数据，用于多端同步
	GetProgressSyncData(ctx context.Context, userID string, lastSyncTime int64) (*dto.ProgressSyncData, error)

	// MergeProgress 合并进度数据
	// 合并来自不同端的进度数据
	MergeProgress(ctx context.Context, userID string, progressData []*dto.ReadingProgressResponse) error
}
