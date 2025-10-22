package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QuotaType 配额类型
type QuotaType string

const (
	QuotaTypeDaily   QuotaType = "daily"   // 日配额
	QuotaTypeMonthly QuotaType = "monthly" // 月配额
	QuotaTypeTotal   QuotaType = "total"   // 总配额
)

// QuotaStatus 配额状态
type QuotaStatus string

const (
	QuotaStatusActive    QuotaStatus = "active"    // 激活
	QuotaStatusExhausted QuotaStatus = "exhausted" // 用尽
	QuotaStatusSuspended QuotaStatus = "suspended" // 暂停
)

// UserQuota 用户配额模型
type UserQuota struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         string             `json:"userId" bson:"user_id"`                           // 用户ID
	QuotaType      QuotaType          `json:"quotaType" bson:"quota_type"`                     // 配额类型
	TotalQuota     int                `json:"totalQuota" bson:"total_quota"`                   // 总配额（Token数或次数）
	UsedQuota      int                `json:"usedQuota" bson:"used_quota"`                     // 已用配额
	RemainingQuota int                `json:"remainingQuota" bson:"remaining_quota"`           // 剩余配额
	Status         QuotaStatus        `json:"status" bson:"status"`                            // 配额状态
	ResetAt        time.Time          `json:"resetAt" bson:"reset_at"`                         // 重置时间
	ExpiresAt      *time.Time         `json:"expiresAt,omitempty" bson:"expires_at,omitempty"` // 过期时间
	Metadata       *QuotaMetadata     `json:"metadata,omitempty" bson:"metadata,omitempty"`    // 元数据
	CreatedAt      time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updated_at"`
}

// QuotaMetadata 配额元数据
type QuotaMetadata struct {
	UserRole          string                 `json:"userRole" bson:"user_role"`                                  // 用户角色（reader/writer）
	MembershipLevel   string                 `json:"membershipLevel" bson:"membership_level"`                    // 会员等级
	LastConsumedAt    *time.Time             `json:"lastConsumedAt,omitempty" bson:"last_consumed_at,omitempty"` // 最后消费时间
	TotalConsumptions int                    `json:"totalConsumptions" bson:"total_consumptions"`                // 总消费次数
	AveragePerDay     float64                `json:"averagePerDay" bson:"average_per_day"`                       // 日均消费
	CustomFields      map[string]interface{} `json:"customFields,omitempty" bson:"custom_fields,omitempty"`      // 自定义字段
}

// QuotaTransaction 配额消费记录
type QuotaTransaction struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID        string             `json:"userId" bson:"user_id"`
	QuotaType     QuotaType          `json:"quotaType" bson:"quota_type"`
	Amount        int                `json:"amount" bson:"amount"`                               // 消费数量（可为负数表示恢复）
	Type          string             `json:"type" bson:"type"`                                   // consume/restore/reset
	Service       string             `json:"service" bson:"service"`                             // 服务类型（text_generation/chat/image）
	Model         string             `json:"model,omitempty" bson:"model,omitempty"`             // 使用的模型
	RequestID     string             `json:"requestId,omitempty" bson:"request_id,omitempty"`    // 请求ID
	Description   string             `json:"description,omitempty" bson:"description,omitempty"` // 描述
	BeforeBalance int                `json:"beforeBalance" bson:"before_balance"`                // 消费前余额
	AfterBalance  int                `json:"afterBalance" bson:"after_balance"`                  // 消费后余额
	Timestamp     time.Time          `json:"timestamp" bson:"timestamp"`
}

// CollectionName 指定集合名
func (UserQuota) CollectionName() string {
	return "ai_user_quotas"
}

// CollectionName 指定集合名
func (QuotaTransaction) CollectionName() string {
	return "ai_quota_transactions"
}

// BeforeCreate MongoDB钩子 - 创建前
func (q *UserQuota) BeforeCreate() {
	if q.ID.IsZero() {
		q.ID = primitive.NewObjectID()
	}
	if q.Status == "" {
		q.Status = QuotaStatusActive
	}
	q.RemainingQuota = q.TotalQuota - q.UsedQuota
	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()
}

// BeforeUpdate MongoDB钩子 - 更新前
func (q *UserQuota) BeforeUpdate() {
	q.RemainingQuota = q.TotalQuota - q.UsedQuota
	q.UpdatedAt = time.Now()

	// 更新状态
	if q.RemainingQuota <= 0 {
		q.Status = QuotaStatusExhausted
	} else {
		q.Status = QuotaStatusActive
	}
}

// IsAvailable 检查配额是否可用
func (q *UserQuota) IsAvailable() bool {
	now := time.Now()

	// 检查状态
	if q.Status != QuotaStatusActive {
		return false
	}

	// 检查是否已过期
	if q.ExpiresAt != nil && now.After(*q.ExpiresAt) {
		return false
	}

	// 检查是否需要重置
	if now.After(q.ResetAt) {
		return false // 需要先重置
	}

	// 检查剩余配额
	return q.RemainingQuota > 0
}

// CanConsume 检查是否可以消费指定数量
func (q *UserQuota) CanConsume(amount int) bool {
	return q.IsAvailable() && q.RemainingQuota >= amount
}

// Consume 消费配额
func (q *UserQuota) Consume(amount int) error {
	if !q.CanConsume(amount) {
		return ErrQuotaExhausted
	}

	q.UsedQuota += amount
	q.RemainingQuota -= amount

	if q.Metadata != nil {
		q.Metadata.TotalConsumptions++
		now := time.Now()
		q.Metadata.LastConsumedAt = &now
	}

	q.BeforeUpdate()
	return nil
}

// Restore 恢复配额
func (q *UserQuota) Restore(amount int) {
	q.UsedQuota -= amount
	if q.UsedQuota < 0 {
		q.UsedQuota = 0
	}

	q.RemainingQuota = q.TotalQuota - q.UsedQuota
	q.BeforeUpdate()
}

// Reset 重置配额
func (q *UserQuota) Reset() {
	q.UsedQuota = 0
	q.RemainingQuota = q.TotalQuota

	// 计算下次重置时间
	now := time.Now()
	switch q.QuotaType {
	case QuotaTypeDaily:
		q.ResetAt = now.AddDate(0, 0, 1)
	case QuotaTypeMonthly:
		q.ResetAt = now.AddDate(0, 1, 0)
	default:
		// 总配额不重置
	}

	q.BeforeUpdate()
}

// ShouldReset 检查是否应该重置
func (q *UserQuota) ShouldReset() bool {
	return time.Now().After(q.ResetAt)
}

// GetUsagePercentage 获取使用百分比
func (q *UserQuota) GetUsagePercentage() float64 {
	if q.TotalQuota == 0 {
		return 0
	}
	return float64(q.UsedQuota) / float64(q.TotalQuota) * 100
}

// QuotaError 配额错误
type QuotaError struct {
	Code    string
	Message string
}

func (e *QuotaError) Error() string {
	return e.Message
}

// 配额错误常量
var (
	ErrQuotaNotFound     = &QuotaError{Code: "QUOTA_NOT_FOUND", Message: "配额记录不存在"}
	ErrQuotaExhausted    = &QuotaError{Code: "QUOTA_EXHAUSTED", Message: "配额已用尽"}
	ErrQuotaExpired      = &QuotaError{Code: "QUOTA_EXPIRED", Message: "配额已过期"}
	ErrQuotaSuspended    = &QuotaError{Code: "QUOTA_SUSPENDED", Message: "配额已暂停"}
	ErrInsufficientQuota = &QuotaError{Code: "INSUFFICIENT_QUOTA", Message: "配额不足"}
)

// QuotaConfig 配额配置
type QuotaConfig struct {
	// 读者配额
	ReaderDailyQuota    int // 普通读者日配额
	VIPReaderDailyQuota int // VIP读者日配额

	// 作者配额
	NoviceWriterDailyQuota int // 新手作者日配额
	SignedWriterDailyQuota int // 签约作者日配额
	MasterWriterDailyQuota int // 大神作者日配额（-1表示无限）
}

// DefaultQuotaConfig 默认配额配置
var DefaultQuotaConfig = &QuotaConfig{
	ReaderDailyQuota:       5,   // 普通读者：5次/日
	VIPReaderDailyQuota:    50,  // VIP读者：50次/日
	NoviceWriterDailyQuota: 10,  // 新手作者：10次/日
	SignedWriterDailyQuota: 100, // 签约作者：100次/日
	MasterWriterDailyQuota: -1,  // 大神作者：无限
}

// GetDefaultQuota 根据用户角色和等级获取默认配额
func GetDefaultQuota(userRole, membershipLevel string) int {
	config := DefaultQuotaConfig

	switch userRole {
	case "reader":
		if membershipLevel == "vip" {
			return config.VIPReaderDailyQuota
		}
		return config.ReaderDailyQuota

	case "writer":
		switch membershipLevel {
		case "novice":
			return config.NoviceWriterDailyQuota
		case "signed":
			return config.SignedWriterDailyQuota
		case "master":
			return config.MasterWriterDailyQuota
		default:
			return config.NoviceWriterDailyQuota
		}

	default:
		return config.ReaderDailyQuota
	}
}
