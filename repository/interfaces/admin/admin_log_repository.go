package admin

import (
	adminModel "Qingyu_backend/models/users"
	"context"
	"time"
)

// AdminLogRepository 管理员操作日志Repository
type AdminLogRepository interface {
	// 操作日志
	CreateAdminLog(ctx context.Context, log *adminModel.AdminLog) error
	GetAdminLog(ctx context.Context, logID string) (*adminModel.AdminLog, error)
	ListAdminLogs(ctx context.Context, filter *AdminLogFilter) ([]*adminModel.AdminLog, error)
	CountAdminLogs(ctx context.Context, filter *AdminLogFilter) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// AdminLogFilter 管理员日志过滤器
type AdminLogFilter struct {
	AdminID   string
	Operation string
	StartDate time.Time
	EndDate   time.Time
	Limit     int64
	Offset    int64
}
