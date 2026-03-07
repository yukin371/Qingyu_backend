package messaging

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/shared"
)

// Message 消息模型
type Message struct {
	ID       primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Topic    string                 `json:"topic" bson:"topic"`
	Payload  map[string]interface{} `json:"payload" bson:"payload"`
	Metadata map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Status   string                 `json:"status" bson:"status"`
	Retry    int                    `json:"retry" bson:"retry"`
	MaxRetry int                    `json:"max_retry" bson:"max_retry"`
	Error    string                 `json:"error,omitempty" bson:"error,omitempty"`
	shared.BaseEntity
	ProcessedAt *time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
}

// MessageTemplate 消息模板
type MessageTemplate struct {
	ID        primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Name      string                 `json:"name" bson:"name"`
	Type      string                 `json:"type" bson:"type"`
	Subject   string                 `json:"subject,omitempty" bson:"subject,omitempty"`
	Content   string                 `json:"content" bson:"content"`
	Variables []string               `json:"variables" bson:"variables"`
	IsActive  bool                   `json:"is_active" bson:"is_active"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	shared.BaseEntity
}

// NotificationDelivery 外部通知发送记录（email/sms/push）
type NotificationDelivery struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID  string             `json:"userId" bson:"user_id"`
	Type    string             `json:"type" bson:"type"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
	Status  string             `json:"status" bson:"status"`
	shared.ReadStatus
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt time.Time              `json:"createdAt" bson:"created_at"`
	SentAt    *time.Time             `json:"sentAt,omitempty" bson:"sent_at,omitempty"`
}

const (
	MessageStatusPending    = "pending"
	MessageStatusProcessing = "processing"
	MessageStatusCompleted  = "completed"
	MessageStatusFailed     = "failed"
)

const (
	NotificationTypeEmail  = "email"
	NotificationTypeSMS    = "sms"
	NotificationTypePush   = "push"
	NotificationTypeSystem = "system"
)

const (
	NotificationStatusSent    = "sent"
	NotificationStatusPending = "pending"
	NotificationStatusFailed  = "failed"
	NotificationStatusDeleted = "deleted"
)

const (
	MessageTemplateTypeEmail = "email"
	MessageTemplateTypeSMS   = "sms"
	MessageTemplateTypePush  = "push"
)

type NotificationDeliveryType string

type NotificationType = NotificationDeliveryType
