package interfaces

import (
	"Qingyu_backend/models"
	readerModel "Qingyu_backend/models/reader"
	"context"
	"time"
)

// ReaderDevicePort 是阅读器跨设备进度场景所需的最小设备接口。
// 具体持久化实现由组合根注入，API 层不依赖完整仓储接口。
type ReaderDevicePort interface {
	UpsertDevice(ctx context.Context, device *models.Device) error
	GetByUserID(ctx context.Context, userID string) ([]*models.Device, error)
}

// ReaderService 阅读器服务接口
type ReaderService interface {
	// =========================
	// 章节相关方法
	// =========================
	GetChapterContent(ctx context.Context, userID, chapterID string) (string, error)
	GetChapterByID(ctx context.Context, chapterID string) (interface{}, error)
	GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error)

	// =========================
	// 阅读进度相关方法
	// =========================
	GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModel.ReadingProgress, error)
	SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error
	UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error
	GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModel.ReadingProgress, error)
	GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*readerModel.ReadingProgress, int64, error)
	GetTotalReadingTime(ctx context.Context, userID string) (int64, error)
	GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error)
	GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error)
	GetFinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error)
	DeleteReadingProgress(ctx context.Context, userID, bookID string) error
	UpdateBookStatus(ctx context.Context, userID, bookID, status string) error
	BatchUpdateBookStatus(ctx context.Context, userID string, bookIDs []string, status string) error

	// =========================
	// 标注相关方法
	// =========================
	CreateAnnotation(ctx context.Context, annotation *readerModel.Annotation) error
	UpdateAnnotation(ctx context.Context, userID, annotationID string, updates map[string]interface{}) error
	DeleteAnnotation(ctx context.Context, userID, annotationID string) error
	GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModel.Annotation, error)
	GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)
	GetNotes(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)
	SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModel.Annotation, error)
	GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)
	GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModel.Annotation, error)
	GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error)
	GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModel.Annotation, error)
	GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModel.Annotation, error)
	GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error)

	// =========================
	// 阅读设置相关方法
	// =========================
	GetReadingSettings(ctx context.Context, userID string) (*readerModel.ReadingSettings, error)
	SaveReadingSettings(ctx context.Context, settings *readerModel.ReadingSettings) error
	UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error

	// =========================
	// 批量操作方法
	// =========================
	BatchCreateAnnotations(ctx context.Context, annotations []*readerModel.Annotation) error
	BatchDeleteAnnotations(ctx context.Context, userID string, annotationIDs []string) error

	// =========================
	// 同步方法
	// =========================
	SyncAnnotations(ctx context.Context, userID string, req interface{}) (map[string]interface{}, error)
}
