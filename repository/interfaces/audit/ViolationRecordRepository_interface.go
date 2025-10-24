package audit

import (
	"context"
	"time"

	"Qingyu_backend/models/audit"
	"Qingyu_backend/repository/interfaces/infrastructure"
)

// ViolationRecordRepository 违规记录Repository接口
type ViolationRecordRepository interface {
	// 基础CRUD
	Create(ctx context.Context, record *audit.ViolationRecord) error
	GetByID(ctx context.Context, id string) (*audit.ViolationRecord, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error

	// 查询方法
	List(ctx context.Context, filter infrastructure.Filter) ([]*audit.ViolationRecord, error)
	Count(ctx context.Context, filter infrastructure.Filter) (int64, error)
	FindWithPagination(ctx context.Context, filter infrastructure.Filter, pagination infrastructure.Pagination) (*infrastructure.PagedResult[audit.ViolationRecord], error)

	// 业务查询
	GetByUserID(ctx context.Context, userID string) ([]*audit.ViolationRecord, error)
	GetActiveViolations(ctx context.Context, userID string) ([]*audit.ViolationRecord, error)
	GetRecentViolations(ctx context.Context, userID string, since time.Time) ([]*audit.ViolationRecord, error)

	// 统计方法
	GetUserSummary(ctx context.Context, userID string) (*audit.UserViolationSummary, error)
	CountByUser(ctx context.Context, userID string) (int64, error)
	CountHighRiskByUser(ctx context.Context, userID string, minLevel int) (int64, error)
	GetTopViolators(ctx context.Context, limit int) ([]*audit.UserViolationSummary, error)

	// 处罚操作
	ApplyPenalty(ctx context.Context, userID string, penaltyType string, duration int, description string) error
	RemovePenalty(ctx context.Context, id string) error
	GetActivePenalties(ctx context.Context) ([]*audit.ViolationRecord, error)
	CleanExpiredPenalties(ctx context.Context) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
