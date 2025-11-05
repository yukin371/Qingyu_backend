package audit

import "time"

// AuditRecord 审核记录模型
type AuditRecord struct {
	ID           string                 `bson:"_id,omitempty" json:"id"`
	TargetType   string                 `bson:"targetType" json:"targetType" validate:"required"`     // 审核对象类型（document/chapter/comment）
	TargetID     string                 `bson:"targetId" json:"targetId" validate:"required"`         // 审核对象ID
	AuthorID     string                 `bson:"authorId" json:"authorId"`                             // 作者ID
	Content      string                 `bson:"content" json:"content"`                               // 审核内容（可能很长，考虑存储策略）
	Status       string                 `bson:"status" json:"status" validate:"required"`             // 审核状态
	Result       string                 `bson:"result" json:"result"`                                 // 审核结果
	RiskLevel    int                    `bson:"riskLevel" json:"riskLevel"`                           // 风险等级 1-5
	RiskScore    float64                `bson:"riskScore" json:"riskScore"`                           // 风险分数 0-100
	Violations   []ViolationDetail      `bson:"violations" json:"violations"`                         // 违规详情列表
	ReviewerID   string                 `bson:"reviewerId,omitempty" json:"reviewerId,omitempty"`     // 复核人ID
	ReviewNote   string                 `bson:"reviewNote,omitempty" json:"reviewNote,omitempty"`     // 复核说明
	AppealStatus string                 `bson:"appealStatus,omitempty" json:"appealStatus,omitempty"` // 申诉状态
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`         // 额外元数据
	CreatedAt    time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time              `bson:"updatedAt" json:"updatedAt"`
	ReviewedAt   *time.Time             `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"` // 复核时间
}

// AuditStatus 审核状态常量
const (
	StatusPending   = "pending"   // 待审核
	StatusReviewing = "reviewing" // 审核中
	StatusApproved  = "approved"  // 通过
	StatusRejected  = "rejected"  // 拒绝
	StatusWarning   = "warning"   // 警告（通过但有风险）
)

// AuditResult 审核结果常量
const (
	ResultPass    = "pass"    // 通过
	ResultWarning = "warning" // 警告
	ResultReject  = "reject"  // 拒绝
	ResultManual  = "manual"  // 需人工复核
)

// AppealStatus 申诉状态常量
const (
	AppealNone     = ""         // 无申诉
	AppealPending  = "pending"  // 申诉待处理
	AppealApproved = "approved" // 申诉通过
	AppealRejected = "rejected" // 申诉驳回
)

// TargetType 审核对象类型
const (
	TargetTypeDocument = "document" // 文档
	TargetTypeChapter  = "chapter"  // 章节
	TargetTypeComment  = "comment"  // 评论
)

// ViolationDetail 违规详情
type ViolationDetail struct {
	Type        string   `bson:"type" json:"type"`               // 违规类型
	Category    string   `bson:"category" json:"category"`       // 违规分类
	Level       int      `bson:"level" json:"level"`             // 严重等级
	Description string   `bson:"description" json:"description"` // 违规描述
	Position    int      `bson:"position" json:"position"`       // 违规位置（字符索引）
	Context     string   `bson:"context" json:"context"`         // 违规上下文（前后文）
	Keywords    []string `bson:"keywords" json:"keywords"`       // 命中的关键词
}

// IsApproved 是否通过审核
func (a *AuditRecord) IsApproved() bool {
	return a.Status == StatusApproved || a.Result == ResultPass
}

// IsRejected 是否被拒绝
func (a *AuditRecord) IsRejected() bool {
	return a.Status == StatusRejected || a.Result == ResultReject
}

// NeedsManualReview 是否需要人工复核
func (a *AuditRecord) NeedsManualReview() bool {
	return a.Result == ResultManual || a.RiskLevel >= 3
}

// CanAppeal 是否可以申诉
func (a *AuditRecord) CanAppeal() bool {
	return a.IsRejected() && (a.AppealStatus == "" || a.AppealStatus == AppealNone)
}

// GetRiskLevelName 获取风险等级名称
func (a *AuditRecord) GetRiskLevelName() string {
	return GetLevelName(a.RiskLevel)
}

// AddViolation 添加违规详情
func (a *AuditRecord) AddViolation(violation ViolationDetail) {
	if a.Violations == nil {
		a.Violations = []ViolationDetail{}
	}
	a.Violations = append(a.Violations, violation)

	// 更新风险等级（取最高）
	if violation.Level > a.RiskLevel {
		a.RiskLevel = violation.Level
	}
}

// CalculateRiskScore 计算风险分数
func (a *AuditRecord) CalculateRiskScore() float64 {
	if len(a.Violations) == 0 {
		return 0
	}

	// 基础分数：违规数量 * 10
	baseScore := float64(len(a.Violations)) * 10

	// 严重度加权
	levelWeight := 0.0
	for _, v := range a.Violations {
		levelWeight += float64(v.Level) * 10
	}

	// 最终分数（0-100）
	score := baseScore + levelWeight
	if score > 100 {
		score = 100
	}

	a.RiskScore = score
	return score
}
