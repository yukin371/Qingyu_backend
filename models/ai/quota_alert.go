package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuotaAlertType 配额告警类型
type QuotaAlertType string

const (
	QuotaAlertTypeThreshold    QuotaAlertType = "threshold"    // 阈值告警
	QuotaAlertTypeAnomaly      QuotaAlertType = "anomaly"     // 异常告警
	QuotaAlertTypeAbuse        QuotaAlertType = "abuse"       // 滥用告警
	QuotaAlertTypeConsistency  QuotaAlertType = "consistency" // 一致性告警
)

// QuotaAlertLevel 配额告警级别
type QuotaAlertLevel string

const (
	QuotaAlertLevelInfo     QuotaAlertLevel = "info"     // 信息
	QuotaAlertLevelWarning  QuotaAlertLevel = "warning"  // 警告
	QuotaAlertLevelCritical QuotaAlertLevel = "critical" // 严重
)

// QuotaAlertStatus 配额告警状态
type QuotaAlertStatus string

const (
	QuotaAlertStatusPending       QuotaAlertStatus = "pending"       // 待处理
	QuotaAlertStatusAcknowledged  QuotaAlertStatus = "acknowledged"  // 已确认
	QuotaAlertStatusResolved     QuotaAlertStatus = "resolved"      // 已解决
	QuotaAlertStatusIgnored      QuotaAlertStatus = "ignored"       // 已忽略
)

// QuotaAlert 配额告警模型
type QuotaAlert struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Type        QuotaAlertType         `json:"type" bson:"type"`                                   // 告警类型
	UserID      string                 `json:"userId,omitempty" bson:"user_id,omitempty"`          // 用户ID（空表示全局告警）
	Level       QuotaAlertLevel        `json:"level" bson:"level"`                                 // 告警级别
	Title       string                 `json:"title" bson:"title"`                                 // 告警标题
	Message     string                 `json:"message" bson:"message"`                            // 告警消息
	Data        map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`              // 扩展数据
	Status      QuotaAlertStatus       `json:"status" bson:"status"`                               // 告警状态
	ResolvedBy  string                 `json:"resolvedBy,omitempty" bson:"resolved_by,omitempty"` // 处理人
	ResolvedAt  *time.Time             `json:"resolvedAt,omitempty" bson:"resolved_at,omitempty"` // 处理时间
	CreatedAt   time.Time              `json:"createdAt" bson:"created_at"`
}

// CollectionName 指定集合名
func (QuotaAlert) CollectionName() string {
	return "ai_quota_alerts"
}

// BeforeCreate MongoDB钩子 - 创建前
func (a *QuotaAlert) BeforeCreate() {
	if a.ID.IsZero() {
		a.ID = primitive.NewObjectID()
	}
	if a.Status == "" {
		a.Status = QuotaAlertStatusPending
	}
	a.CreatedAt = time.Now()
}

// IsGlobal 检查是否为全局告警
func (a *QuotaAlert) IsGlobal() bool {
	return a.UserID == ""
}

// IsPending 检查告警是否待处理
func (a *QuotaAlert) IsPending() bool {
	return a.Status == QuotaAlertStatusPending
}

// Resolve 标记告警为已解决
func (a *QuotaAlert) Resolve(resolvedBy string) {
	a.Status = QuotaAlertStatusResolved
	a.ResolvedBy = resolvedBy
	now := time.Now()
	a.ResolvedAt = &now
}
