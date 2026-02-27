package audit

import (
	auditInterface "Qingyu_backend/service/interfaces/audit"
	"context"

	"Qingyu_backend/models/audit"
)

// NoopContentAuditService 提供一个可用但不做实际审核的默认实现。
// 用于在审核仓储尚未完成接入时，保持路由与服务链路可用。
type NoopContentAuditService struct{}

// NewNoopContentAuditService 创建 Noop 审核服务
func NewNoopContentAuditService() *NoopContentAuditService {
	return &NoopContentAuditService{}
}

func (s *NoopContentAuditService) CheckContent(ctx context.Context, content string) (*auditInterface.AuditCheckResult, error) {
	return &auditInterface.AuditCheckResult{
		IsSafe:      true,
		RiskLevel:   0,
		RiskScore:   0,
		Violations:  []audit.ViolationDetail{},
		Suggestions: []string{},
		NeedsReview: false,
		CanPublish:  true,
	}, nil
}

func (s *NoopContentAuditService) AuditDocument(ctx context.Context, documentID string, content string, authorID string) (*audit.AuditRecord, error) {
	return &audit.AuditRecord{}, nil
}

func (s *NoopContentAuditService) GetAuditResult(ctx context.Context, targetType, targetID string) (*audit.AuditRecord, error) {
	return &audit.AuditRecord{}, nil
}

func (s *NoopContentAuditService) BatchAuditDocuments(ctx context.Context, documentIDs []string) ([]*audit.AuditRecord, error) {
	return []*audit.AuditRecord{}, nil
}

func (s *NoopContentAuditService) ReviewAudit(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	return nil
}

func (s *NoopContentAuditService) SubmitAppeal(ctx context.Context, auditID string, authorID string, reason string) error {
	return nil
}

func (s *NoopContentAuditService) ReviewAppeal(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error {
	return nil
}

func (s *NoopContentAuditService) GetUserViolations(ctx context.Context, userID string) ([]*audit.ViolationRecord, error) {
	return []*audit.ViolationRecord{}, nil
}

func (s *NoopContentAuditService) GetUserViolationSummary(ctx context.Context, userID string) (*audit.UserViolationSummary, error) {
	return &audit.UserViolationSummary{
		UserID:              userID,
		TotalViolations:     0,
		WarningCount:        0,
		RejectCount:         0,
		HighRiskCount:       0,
		ActivePenalties:     0,
		IsBanned:            false,
		IsPermanentlyBanned: false,
	}, nil
}

func (s *NoopContentAuditService) GetPendingReviews(ctx context.Context, limit int) ([]*audit.AuditRecord, error) {
	return []*audit.AuditRecord{}, nil
}

func (s *NoopContentAuditService) GetHighRiskAudits(ctx context.Context, minRiskLevel int, limit int) ([]*audit.AuditRecord, error) {
	return []*audit.AuditRecord{}, nil
}

func (s *NoopContentAuditService) GetAuditStatistics(ctx context.Context) (interface{}, error) {
	return map[string]interface{}{
		"totalAudits":    0,
		"pendingAudits":  0,
		"approvedAudits": 0,
		"rejectedAudits": 0,
		"highRiskAudits": 0,
	}, nil
}
