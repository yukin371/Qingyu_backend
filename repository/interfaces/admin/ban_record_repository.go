package admin

import (
	"context"

	"Qingyu_backend/models/admin"
)

// BanRecordRepository 封禁记录仓储接口
type BanRecordRepository interface {
	// Create 创建封禁记录
	Create(ctx context.Context, record *admin.BanRecord) error

	// GetByUserID 获取用户的封禁历史
	GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*admin.BanRecord, int64, error)

	// GetActiveBan 获取用户当前生效的封禁记录（最后一条ban操作）
	GetActiveBan(ctx context.Context, userID string) (*admin.BanRecord, error)
}
