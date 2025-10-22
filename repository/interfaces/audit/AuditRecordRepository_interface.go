package audit

import (
	"context"
	"time"

	"Qingyu_backend/models/audit"
	"Qingyu_backend/repository/interfaces/infrastructure"
)

// AuditRecordRepository 审核记录Repository接口
type AuditRecordRepository interface {
	// 基础CRUD
	Create(ctx context.Context, record *audit.AuditRecord) error
	GetByID(ctx context.Context, id string) (*audit.AuditRecord, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 查询方法
	List(ctx context.Context, filter infrastructure.Filter) ([]*audit.AuditRecord, error)
	Count(ctx context.Context, filter infrastructure.Filter) (int64, error)
	FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[audit.AuditRecord], error)

	// 业务查询
	GetByTargetID(ctx context.Context, targetType, targetID string) (*audit.AuditRecord, error)
	GetByAuthor(ctx context.Context, authorID string, limit int64) ([]*audit.AuditRecord, error)
	GetByStatus(ctx context.Context, status string, limit int64) ([]*audit.AuditRecord, error)
	GetByResult(ctx context.Context, result string, limit int64) ([]*audit.AuditRecord, error)
	GetPendingReview(ctx context.Context, limit int64) ([]*audit.AuditRecord, error)
	GetHighRisk(ctx context.Context, minRiskLevel int, limit int64) ([]*audit.AuditRecord, error)

	// 审核操作
	UpdateStatus(ctx context.Context, id string, status string, reviewerID string, note string) error
	UpdateAppealStatus(ctx context.Context, id string, appealStatus string) error
	BatchUpdateStatus(ctx context.Context, ids []string, status string) error

	// 统计方法
	CountByStatus(ctx context.Context, startTime, endTime time.Time) (map[string]int64, error)
	CountByAuthor(ctx context.Context, authorID string) (int64, error)
	CountHighRiskByAuthor(ctx context.Context, authorID string, minLevel int) (int64, error)
	GetRecentRecords(ctx context.Context, targetType string, limit int64) ([]*audit.AuditRecord, error)

	// 健康检查
	Health(ctx context.Context) error
}
