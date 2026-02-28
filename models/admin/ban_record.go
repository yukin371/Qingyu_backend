package admin

import (
	"Qingyu_backend/models/shared"
)

// BanRecord 封禁记录模型
// 用于记录用户的封禁/解封历史，支持审计追踪
type BanRecord struct {
	shared.IdentifiedEntity `bson:",inline"`
	shared.BaseEntity       `bson:",inline"`

	UserID     string `bson:"user_id" json:"userId"`
	Action     string `bson:"action" json:"action"`     // "ban" 或 "unban"
	Reason     string `bson:"reason" json:"reason"`
	OperatorID string `bson:"operator_id" json:"operatorId"`
	BanDuration *int  `bson:"ban_duration,omitempty" json:"banDuration,omitempty"` // 封禁时长（天）
}

// IsBanAction 判断是否为封禁操作
func (r *BanRecord) IsBanAction() bool {
	return r.Action == "ban"
}

// IsUnbanAction 判断是否为解封操作
func (r *BanRecord) IsUnbanAction() bool {
	return r.Action == "unban"
}
