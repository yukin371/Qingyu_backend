package reading

import (
	"context"
	readerModel "Qingyu_backend/models/reader"
)

// ReaderThemeRepository 读者主题仓储接口
type ReaderThemeRepository interface {
	// 内置主题管理
	GetBuiltInThemes(ctx context.Context) ([]*readerModel.ReaderTheme, error)
	GetThemeByName(ctx context.Context, name string) (*readerModel.ReaderTheme, error)

	// 用户主题管理
	CreateTheme(ctx context.Context, theme *readerModel.ReaderTheme) error
	GetTheme(ctx context.Context, themeID string) (*readerModel.ReaderTheme, error)
	UpdateTheme(ctx context.Context, themeID string, updates map[string]interface{}) error
	DeleteTheme(ctx context.Context, themeID string) error

	// 查询用户主题
	GetUserThemes(ctx context.Context, userID string) ([]*readerModel.ReaderTheme, error)
	GetActiveTheme(ctx context.Context, userID string) (*readerModel.ReaderTheme, error)
	SetActiveTheme(ctx context.Context, userID, themeID string) error

	// 公开主题
	GetPublicThemes(ctx context.Context, page, pageSize int) ([]*readerModel.ReaderTheme, int64, error)
	IncrementUseCount(ctx context.Context, themeID string) error

	// 批量操作
	BatchGetThemes(ctx context.Context, themeIDs []string) (map[string]*readerModel.ReaderTheme, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
