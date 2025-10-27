package messaging

import "time"

// Message 消息模型
type Message struct {
	ID          string                 `json:"id" bson:"_id,omitempty"`
	Topic       string                 `json:"topic" bson:"topic"`                           // 主题
	Payload     map[string]interface{} `json:"payload" bson:"payload"`                       // 消息内容
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"` // 元数据
	Status      string                 `json:"status" bson:"status"`                         // pending, processing, completed, failed
	Retry       int                    `json:"retry" bson:"retry"`                           // 重试次数
	MaxRetry    int                    `json:"max_retry" bson:"max_retry"`                   // 最大重试次数
	Error       string                 `json:"error,omitempty" bson:"error,omitempty"`       // 错误信息
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
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
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at"`
}

// Notification 通知记录
type Notification struct {
	ID        string                 `json:"id" bson:"_id,omitempty"`
	UserID    string                 `json:"user_id" bson:"user_id"`
	Type      string                 `json:"type" bson:"type"` // email, sms, push, system
	Title     string                 `json:"title" bson:"title"`
	Content   string                 `json:"content" bson:"content"`
	Status    string                 `json:"status" bson:"status"` // sent, pending, failed
	IsRead    bool                   `json:"is_read" bson:"is_read"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	ReadAt    *time.Time             `json:"read_at,omitempty" bson:"read_at,omitempty"`
	SentAt    *time.Time             `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
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

// NotificationType 通知类型（类型别名）
type NotificationType string
