package interfaces

import (
	"context"

	"Qingyu_backend/models/audit"
)

// ContentAuditService 内容审核服务接口
type ContentAuditService interface {
	// 实时检测
	CheckContent(ctx context.Context, content string) (*AuditCheckResult, error)

	// 全文审核
	AuditDocument(ctx context.Context, documentID string, content string, authorID string) (*audit.AuditRecord, error)

	// 获取审核结果
	GetAuditResult(ctx context.Context, targetType, targetID string) (*audit.AuditRecord, error)

	// 批量审核
	BatchAuditDocuments(ctx context.Context, documentIDs []string) ([]*audit.AuditRecord, error)

	// 复核
	ReviewAudit(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error

	// 申诉
	SubmitAppeal(ctx context.Context, auditID string, authorID string, reason string) error
	ReviewAppeal(ctx context.Context, auditID string, reviewerID string, approved bool, note string) error

	// 用户违规查询
	GetUserViolations(ctx context.Context, userID string) ([]*audit.ViolationRecord, error)
	GetUserViolationSummary(ctx context.Context, userID string) (*audit.UserViolationSummary, error)

	// 管理方法
	GetPendingReviews(ctx context.Context, limit int) ([]*audit.AuditRecord, error)
	GetHighRiskAudits(ctx context.Context, minRiskLevel int, limit int) ([]*audit.AuditRecord, error)
}

// AuditCheckResult 审核检测结果
type AuditCheckResult struct {
	IsSafe      bool                    `json:"isSafe"`      // 是否安全
	RiskLevel   int                     `json:"riskLevel"`   // 风险等级 0-5
	RiskScore   float64                 `json:"riskScore"`   // 风险分数 0-100
	Violations  []audit.ViolationDetail `json:"violations"`  // 违规详情
	Suggestions []string                `json:"suggestions"` // 修改建议
	NeedsReview bool                    `json:"needsReview"` // 是否需要人工复核
	CanPublish  bool                    `json:"canPublish"`  // 是否可以发布
}
