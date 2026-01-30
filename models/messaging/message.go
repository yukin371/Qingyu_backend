package messaging

import (
	"time"

	"Qingyu_backend/models/shared"
)

// Message 消息模型
type Message struct {
	ID                string                 `json:"id" bson:"_id,omitempty"`
	Topic             string                 `json:"topic" bson:"topic"`                           // 主题
	Payload           map[string]interface{} `json:"payload" bson:"payload"`                       // 消息内容
	Metadata          map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"` // 元数据
	Status            string                 `json:"status" bson:"status"`                         // pending, processing, completed, failed
	Retry             int                    `json:"retry" bson:"retry"`                           // 重试次数
	MaxRetry          int                    `json:"max_retry" bson:"max_retry"`                   // 最大重试次数
	Error             string                 `json:"error,omitempty" bson:"error,omitempty"`       // 错误信息
	shared.BaseEntity                        // 嵌入时间戳字段
	ProcessedAt       *time.Time             `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
}

// MessageTemplate 消息模板
type MessageTemplate struct {
	ID        string                 `json:"id" bson:"_id,omitempty"`
	Name      string                 `json:"name" bson:"name"`                           // 模板名称
	Type      string                 `json:"type" bson:"type"`                           // email, sms, push
	Subject   string                 `json:"subject,omitempty" bson:"subject,omitempty"` // 邮件主题
	Content   string                 `json:"content" bson:"content"`                     // 模板内容
	Variables []string               `json:"variables" bson:"variables"`                 // 可用变量
	IsActive  bool                   `json:"is_active" bson:"is_active"`                 // 是否启用
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	shared.BaseEntity
}

// NotificationDelivery 外部通知发送记录（email/sms/push）
type NotificationDelivery struct {
	ID                string                 `json:"id" bson:"_id,omitempty"`
	UserID            string                 `json:"userId" bson:"user_id"`
	Type              string                 `json:"type" bson:"type"` // email, sms, push, system
	Title             string                 `json:"title" bson:"title"`
	Content           string                 `json:"content" bson:"content"`
	Status            string                 `json:"status" bson:"status"` // sent, pending, failed
	shared.ReadStatus                        // 嵌入已读状态
	Metadata          map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"createdAt" bson:"created_at"`
	SentAt            *time.Time             `json:"sentAt,omitempty" bson:"sent_at,omitempty"`
}

// 消息状态
const (
	MessageStatusPending    = "pending"
	MessageStatusProcessing = "processing"
	MessageStatusCompleted  = "completed"
	MessageStatusFailed     = "failed"
)

// 通知类型
const (
	NotificationTypeEmail  = "email"
	NotificationTypeSMS    = "sms"
	NotificationTypePush   = "push"
	NotificationTypeSystem = "system"
)

// 通知状态
const (
	NotificationStatusSent    = "sent"
	NotificationStatusPending = "pending"
	NotificationStatusFailed  = "failed"
	NotificationStatusDeleted = "deleted"
)

// 消息模板类型
const (
	MessageTemplateTypeEmail = "email"
	MessageTemplateTypeSMS   = "sms"
	MessageTemplateTypePush  = "push"
)

// NotificationDeliveryType 通知发送类型（类型别名）
type NotificationDeliveryType string

// 向后兼容的类型别名
type NotificationType = NotificationDeliveryType
