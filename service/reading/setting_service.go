package reading

import (
	"Qingyu_backend/models/reader"
	"context"
	"errors"

	"Qingyu_backend/repository/interfaces/reading"
)

// SettingService 阅读设置服务
type SettingService struct {
	settingsRepo reading.ReadingSettingsRepository
}

// NewSettingService 创建设置服务
func NewSettingService(settingsRepo reading.ReadingSettingsRepository) *SettingService {
	return &SettingService{
		settingsRepo: settingsRepo,
	}
}

// GetSetting 获取用户阅读设置
func (s *SettingService) GetSetting(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 先尝试获取用户现有设置
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 如果用户没有设置，创建默认设置
	if settings == nil {
		settings, err = s.settingsRepo.CreateDefaultSettings(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	return settings, nil
}

// UpdateSetting 更新用户阅读设置
func (s *SettingService) UpdateSetting(ctx context.Context, userID string, settings *reader.ReadingSettings) error {
	if userID == "" {
		return errors.New("用户ID不能为空")
	}

	if settings == nil {
		return errors.New("设置信息不能为空")
	}

	// 确保设置属于指定用户
	settings.UserID = userID

	return s.settingsRepo.UpdateByUserID(ctx, userID, settings)
}

// ResetSetting 重置用户阅读设置为默认值
func (s *SettingService) ResetSetting(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 检查用户是否已有设置
	exists, err := s.settingsRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if exists {
		// 如果存在，先删除现有设置
		currentSettings, err := s.settingsRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if currentSettings != nil {
			err = s.settingsRepo.Delete(ctx, currentSettings.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	// 创建新的默认设置
	return s.settingsRepo.CreateDefaultSettings(ctx, userID)
}
