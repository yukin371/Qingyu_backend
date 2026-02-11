package impl

import (
	"context"

	"Qingyu_backend/models/reader"
	serviceReader "Qingyu_backend/service/interfaces/reader"
	readerService "Qingyu_backend/service/reader"
)

// ReaderSettingsImpl 阅读设置管理端口实现
type ReaderSettingsImpl struct {
	settingService *readerService.SettingService
	readerService  *readerService.ReaderService
	serviceName    string
	version        string
}

// NewReaderSettingsImpl 创建阅读设置管理端口实现
func NewReaderSettingsImpl(
	settingService *readerService.SettingService,
	readerService *readerService.ReaderService,
) serviceReader.ReaderSettingsPort {
	return &ReaderSettingsImpl{
		settingService: settingService,
		readerService:  readerService,
		serviceName:    "ReaderSettingsPort",
		version:        "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (r *ReaderSettingsImpl) Initialize(ctx context.Context) error {
	return r.readerService.Initialize(ctx)
}

func (r *ReaderSettingsImpl) Health(ctx context.Context) error {
	return r.readerService.Health(ctx)
}

func (r *ReaderSettingsImpl) Close(ctx context.Context) error {
	return r.readerService.Close(ctx)
}

func (r *ReaderSettingsImpl) GetServiceName() string {
	return r.serviceName
}

func (r *ReaderSettingsImpl) GetVersion() string {
	return r.version
}

// ============================================================================
// ReaderSettingsPort 方法实现
// ============================================================================

// GetReadingSettings 获取阅读设置
func (r *ReaderSettingsImpl) GetReadingSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	return r.readerService.GetReadingSettings(ctx, userID)
}

// SaveReadingSettings 保存阅读设置
func (r *ReaderSettingsImpl) SaveReadingSettings(ctx context.Context, settings *reader.ReadingSettings) error {
	return r.readerService.SaveReadingSettings(ctx, settings)
}

// UpdateReadingSettings 更新阅读设置
func (r *ReaderSettingsImpl) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	return r.readerService.UpdateReadingSettings(ctx, userID, updates)
}
