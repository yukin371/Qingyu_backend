package admin

import (
	adminModel "Qingyu_backend/models/users"
	"context"
	"time"
)

// AdminLogRepository 管理员操作日志Repository（扩展版-审计追踪）
type AdminLogRepository interface {
	// 操作日志
	CreateAdminLog(ctx context.Context, log *adminModel.AdminLog) error
	GetAdminLog(ctx context.Context, logID string) (*adminModel.AdminLog, error)
	ListAdminLogs(ctx context.Context, filter *AdminLogFilter) ([]*adminModel.AdminLog, error)
	CountAdminLogs(ctx context.Context, filter *AdminLogFilter) (int64, error)

	// 审计追踪新增方法
	GetByResource(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error)
	CleanOldLogs(ctx context.Context, beforeDate time.Time) error

	// Health 健康检查
	Health(ctx context.Context) error
}

// AdminLogFilter 管理员日志过滤器（扩展版-审计追踪）
type AdminLogFilter struct {
	AdminID      string
	Operation    string
	ResourceType string // 新增：资源类型
	ResourceID   string // 新增：资源ID
	StartDate    time.Time
	EndDate      time.Time
	Keyword      string // 新增：关键词
	Limit        int64
	Offset       int64
}
