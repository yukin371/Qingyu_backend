package audit

import "time"

// ViolationRecord 违规记录模型（用于统计和封号判断）
type ViolationRecord struct {
	ID              string     `bson:"_id,omitempty" json:"id"`
	UserID          string     `bson:"userId" json:"userId" validate:"required"`         // 用户ID
	AuditRecordID   string     `bson:"auditRecordId" json:"auditRecordId"`               // 关联的审核记录ID
	TargetType      string     `bson:"targetType" json:"targetType"`                     // 违规对象类型
	TargetID        string     `bson:"targetId" json:"targetId"`                         // 违规对象ID
	ViolationType   string     `bson:"violationType" json:"violationType"`               // 违规类型
	ViolationLevel  int        `bson:"violationLevel" json:"violationLevel"`             // 违规等级
	ViolationCount  int        `bson:"violationCount" json:"violationCount"`             // 违规次数（累计）
	PenaltyType     string     `bson:"penaltyType,omitempty" json:"penaltyType"`         // 处罚类型
	PenaltyDuration int        `bson:"penaltyDuration,omitempty" json:"penaltyDuration"` // 处罚时长（天）
	IsPenalized     bool       `bson:"isPenalized" json:"isPenalized"`                   // 是否已处罚
	PenalizedAt     *time.Time `bson:"penalizedAt,omitempty" json:"penalizedAt"`         // 处罚时间
	ExpiresAt       *time.Time `bson:"expiresAt,omitempty" json:"expiresAt"`             // 处罚到期时间
	Description     string     `bson:"description" json:"description"`                   // 违规描述
	CreatedAt       time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time  `bson:"updatedAt" json:"updatedAt"`
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
	TotalViolations     int       `bson:"totalViolations" json:"totalViolations"`
	WarningCount        int       `bson:"warningCount" json:"warningCount"`
	RejectCount         int       `bson:"rejectCount" json:"rejectCount"`
	HighRiskCount       int       `bson:"highRiskCount" json:"highRiskCount"`
	LastViolationAt     time.Time `bson:"lastViolationAt" json:"lastViolationAt"`
	ActivePenalties     int       `bson:"activePenalties" json:"activePenalties"`
	IsBanned            bool      `bson:"isBanned" json:"isBanned"`
	IsPermanentlyBanned bool      `bson:"isPermanentlyBanned" json:"isPermanentlyBanned"`
}

// IsHighRiskUser 是否高风险用户
func (s *UserViolationSummary) IsHighRiskUser() bool {
	return s.HighRiskCount >= 2 || s.TotalViolations >= 5
}

// ShouldBan 是否应该封号
func (s *UserViolationSummary) ShouldBan() bool {
	return s.HighRiskCount >= 3 || s.TotalViolations >= 10
}
