package audit

import "time"

// ViolationRecord 违规记录模型（用于统计和封号判断）
type ViolationRecord struct {
	ID              string     `bson:"_id,omitempty" json:"id"`
	UserID          string     `bson:"user_id" json:"userId" validate:"required"`         // 用户ID
	AuditRecordID   string     `bson:"audit_record_id" json:"auditRecordId"`               // 关联的审核记录ID
	TargetType      string     `bson:"target_type" json:"targetType"`                     // 违规对象类型
	TargetID        string     `bson:"target_id" json:"targetId"`                         // 违规对象ID
	ViolationType   string     `bson:"violation_type" json:"violationType"`               // 违规类型
	ViolationLevel  int        `bson:"violation_level" json:"violationLevel"`             // 违规等级
	ViolationCount  int        `bson:"violation_count" json:"violationCount"`             // 违规次数（累计）
	PenaltyType     string     `bson:"penalty_type,omitempty" json:"penaltyType"`         // 处罚类型
	PenaltyDuration int        `bson:"penalty_duration,omitempty" json:"penaltyDuration"` // 处罚时长（天）
	IsPenalized     bool       `bson:"is_penalized" json:"isPenalized"`                   // 是否已处罚
	PenalizedAt     *time.Time `bson:"penalized_at,omitempty" json:"penalizedAt"`         // 处罚时间
	ExpiresAt       *time.Time `bson:"expires_at,omitempty" json:"expiresAt"`             // 处罚到期时间
	Description     string     `bson:"description" json:"description"`                   // 违规描述
	CreatedAt       time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `bson:"updated_at" json:"updatedAt"`
}

// PenaltyType 处罚类型常量
const (
	PenaltyWarning       = "warning"        // 警告
	PenaltyContentHidden = "content_hidden" // 内容隐藏
	PenaltyAccountMuted  = "account_muted"  // 禁言
	PenaltyAccountBanned = "account_banned" // 封号
	PenaltyPermanentBan  = "permanent_ban"  // 永久封号
)

// GetPenaltyTypeName 获取处罚类型名称
func GetPenaltyTypeName(penaltyType string) string {
	names := map[string]string{
		PenaltyWarning:       "警告",
		PenaltyContentHidden: "内容隐藏",
		PenaltyAccountMuted:  "禁言",
		PenaltyAccountBanned: "封号",
		PenaltyPermanentBan:  "永久封号",
	}
	if name, ok := names[penaltyType]; ok {
		return name
	}
	return "未知处罚"
}

// IsActive 处罚是否生效中
func (v *ViolationRecord) IsActive() bool {
	if !v.IsPenalized || v.ExpiresAt == nil {
		return false
	}
	return time.Now().Before(*v.ExpiresAt)
}

// IsPermanentBan 是否永久封号
func (v *ViolationRecord) IsPermanentBan() bool {
	return v.PenaltyType == PenaltyPermanentBan
}

// ShouldEscalatePenalty 是否应该升级处罚
func (v *ViolationRecord) ShouldEscalatePenalty() bool {
	// 违规3次以上，或者违规等级4级以上
	return v.ViolationCount >= 3 || v.ViolationLevel >= LevelCritical
}

// UserViolationSummary 用户违规统计
type UserViolationSummary struct {
	UserID              string    `bson:"_id" json:"userId"`
	TotalViolations     int       `bson:"total_violations" json:"totalViolations"`
	WarningCount        int       `bson:"warning_count" json:"warningCount"`
	RejectCount         int       `bson:"reject_count" json:"rejectCount"`
	HighRiskCount       int       `bson:"high_risk_count" json:"highRiskCount"`
	LastViolationAt     time.Time `bson:"last_violation_at" json:"lastViolationAt"`
	ActivePenalties     int       `bson:"active_penalties" json:"activePenalties"`
	IsBanned            bool      `bson:"is_banned" json:"isBanned"`
	IsPermanentlyBanned bool      `bson:"is_permanently_banned" json:"isPermanentlyBanned"`
}

// IsHighRiskUser 是否高风险用户
func (s *UserViolationSummary) IsHighRiskUser() bool {
	return s.HighRiskCount >= 2 || s.TotalViolations >= 5
}

// ShouldBan 是否应该封号
func (s *UserViolationSummary) ShouldBan() bool {
	return s.HighRiskCount >= 3 || s.TotalViolations >= 10
}
