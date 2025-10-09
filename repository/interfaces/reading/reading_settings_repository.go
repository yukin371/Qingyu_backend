package reading

import (
	"Qingyu_backend/models/reading/reader"
	base "Qingyu_backend/repository/interfaces/infrastructure"
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
