package ai

import (
	"context"

	aiModels "Qingyu_backend/models/ai"
)

// QuotaAlertRepository 配额告警Repository接口
type QuotaAlertRepository interface {
	// Create 创建告警
	Create(ctx context.Context, alert *aiModels.QuotaAlert) error

	// GetByID 根据ID获取告警
	GetByID(ctx context.Context, id string) (*aiModels.QuotaAlert, error)

	// List 获取告警列表
	List(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error)

	// GetByUserID 获取用户的所有告警
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]*aiModels.QuotaAlert, int64, error)

	// Update 更新告警
	Update(ctx context.Context, alert *aiModels.QuotaAlert) error

	// GetRecentGlobal 获取最近的全剧告警
	GetRecentGlobal(ctx context.Context, limit int) ([]*aiModels.QuotaAlert, error)

	// CountByStatus 按状态统计告警数量
	CountByStatus(ctx context.Context) (map[aiModels.QuotaAlertStatus]int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
