package reading

import (
	"Qingyu_backend/models/reading/reader"
	base "Qingyu_backend/repository/interfaces"
	"context"
)

// ReadingSettingsRepository 阅读设置仓储接口
type ReadingSettingsRepository interface {
	// 继承基础Repository接口
	base.CRUDRepository[*reader.ReadingSettings, string]

	// 阅读设置特定方法
	GetByUserID(ctx context.Context, userID string) (*reader.ReadingSettings, error)
	UpdateByUserID(ctx context.Context, userID string, settings *reader.ReadingSettings) error
	CreateDefaultSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
}

// ReadingProgressRepository 阅读进度仓储接口
type ReadingProgressRepository interface {
	base.CRUDRepository[*reader.ReadingProgress, string]

	// 阅读进度特定方法
	GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error)
	UpdateProgress(ctx context.Context, userID, bookID string, progress float64) error
	GetUserReadingHistory(ctx context.Context, userID string, limit int) ([]*reader.ReadingProgress, error)
	GetBookReaders(ctx context.Context, bookID string, limit int) ([]*reader.ReadingProgress, error)
}

// AnnotationRepository 标注仓储接口
type AnnotationRepository interface {
	base.CRUDRepository[*reader.Annotation, string]

	// 标注特定方法
	GetByUserAndBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error)
	GetByUserAndChapter(ctx context.Context, userID, chapterID string) ([]*reader.Annotation, error)
	GetByType(ctx context.Context, userID string, annotationType reader.AnnotationType) ([]*reader.Annotation, error)
	DeleteByUserAndBook(ctx context.Context, userID, bookID string) error
}
