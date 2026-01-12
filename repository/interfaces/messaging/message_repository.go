package messaging

import (
	messagingModel "Qingyu_backend/models/messaging"
	"context"
	"time"
)

// MessageRepository 消息相关Repository
type MessageRepository interface {
	// 消息队列管理
	CreateMessage(ctx context.Context, message *messagingModel.Message) error
	GetMessage(ctx context.Context, messageID string) (*messagingModel.Message, error)
	UpdateMessage(ctx context.Context, messageID string, updates map[string]interface{}) error
	DeleteMessage(ctx context.Context, messageID string) error
	ListMessages(ctx context.Context, filter *MessageFilter) ([]*messagingModel.Message, error)
	CountMessages(ctx context.Context, filter *MessageFilter) (int64, error)

	// 通知记录管理
	CreateNotification(ctx context.Context, notification *messagingModel.Notification) error
	GetNotification(ctx context.Context, notificationID string) (*messagingModel.Notification, error)
	UpdateNotification(ctx context.Context, notificationID string, updates map[string]interface{}) error
	DeleteNotification(ctx context.Context, notificationID string) error
	ListNotifications(ctx context.Context, filter *NotificationFilter) ([]*messagingModel.Notification, int64, error)
	CountNotifications(ctx context.Context, filter *NotificationFilter) (int64, error)

	// 通知已读状态管理
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkMultipleAsRead(ctx context.Context, userID string, notificationIDs []string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)

	// 批量操作
	BatchDeleteNotifications(ctx context.Context, userID string, notificationIDs []string) error
	ClearAllNotifications(ctx context.Context, userID string) error

	// 消息模板管理
	CreateTemplate(ctx context.Context, template *messagingModel.MessageTemplate) error
	GetTemplate(ctx context.Context, templateID string) (*messagingModel.MessageTemplate, error)
	GetTemplateByName(ctx context.Context, name string) (*messagingModel.MessageTemplate, error)
	UpdateTemplate(ctx context.Context, templateID string, updates map[string]interface{}) error
	DeleteTemplate(ctx context.Context, templateID string) error
	ListTemplates(ctx context.Context, templateType string, isActive *bool) ([]*messagingModel.MessageTemplate, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// MessageFilter 消息过滤器
type MessageFilter struct {
	Topic     string
	Status    string
	StartDate time.Time
	EndDate   time.Time
	Limit     int64
	Offset    int64
}

// NotificationFilter 通知过滤器
type NotificationFilter struct {
	UserID    string
	Type      string
	Status    string
	IsRead    *bool
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	PageSize  int
	Limit     int64
	Offset    int64
}
