package reader

import (
	"Qingyu_backend/models"
	"context"
)

// DeviceRepository 设备仓储接口
type DeviceRepository interface {
	// UpsertDevice 创建或更新设备记录（按 user_id + user_agent 去重）
	UpsertDevice(ctx context.Context, device *models.Device) error
	// GetByUserID 获取用户所有设备
	GetByUserID(ctx context.Context, userID string) ([]*models.Device, error)
	// DeleteByUserID 删除用户所有设备
	DeleteByUserID(ctx context.Context, userID string) error
}
