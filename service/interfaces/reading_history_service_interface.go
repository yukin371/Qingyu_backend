package interfaces

import (
	"context"
	"time"

	readerModels "Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/interfaces/reader"
	readerservice "Qingyu_backend/service/reader"
)

// ReadingHistoryService 阅读历史服务接口
type ReadingHistoryService interface {
	// RecordReading 记录阅读历史
	RecordReading(ctx context.Context, userID, bookID, chapterID string, startTime, endTime time.Time, progress float64, deviceType, deviceID string) error

	// GetUserHistories 获取用户阅读历史列表
	GetUserHistories(ctx context.Context, userID string, page, pageSize int) ([]*readerModels.ReadingHistory, *readerservice.PaginationInfo, error)

	// GetUserHistoriesByBook 获取用户某本书的阅读历史
	GetUserHistoriesByBook(ctx context.Context, userID, bookID string, page, pageSize int) ([]*readerModels.ReadingHistory, *readerservice.PaginationInfo, error)

	// GetUserReadingStats 获取用户阅读统计
	GetUserReadingStats(ctx context.Context, userID string) (*readerRepo.ReadingStats, error)

	// GetUserDailyReadingStats 获取用户每日阅读统计
	GetUserDailyReadingStats(ctx context.Context, userID string, days int) ([]readerRepo.DailyReadingStats, error)

	// DeleteHistory 删除单条历史记录
	DeleteHistory(ctx context.Context, userID, historyID string) error

	// ClearUserHistories 清空用户所有历史记录
	ClearUserHistories(ctx context.Context, userID string) error
}
