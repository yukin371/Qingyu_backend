package admin

import (
	adminModel "Qingyu_backend/models/users"
	"context"
	"time"
)

// AuditRepository 审核记录Repository
type AuditRepository interface {
	// 审核记录
	CreateAuditRecord(ctx context.Context, record *adminModel.AuditRecord) error
	GetAuditRecord(ctx context.Context, contentID, contentType string) (*adminModel.AuditRecord, error)
	UpdateAuditRecord(ctx context.Context, recordID string, updates map[string]interface{}) error
	ListAuditRecords(ctx context.Context, filter *AuditFilter) ([]*adminModel.AuditRecord, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// AuditFilter 审核记录过滤器
type AuditFilter struct {
	ContentType string
	Status      string
	ReviewerID  string
	StartDate   time.Time
	EndDate     time.Time
	Limit       int64
	Offset      int64
}
