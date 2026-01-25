package notification

import (
	"context"
	"time"

	"Qingyu_backend/models/notification"
)

// NotificationRepository 通知仓储接口
type NotificationRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, notification *notification.Notification) error
	GetByID(ctx context.Context, id string) (*notification.Notification, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)

	// 查询操作
	List(ctx context.Context, filter *notification.NotificationFilter) ([]*notification.Notification, error)
	Count(ctx context.Context, filter *notification.NotificationFilter) (int64, error)

	// 批量操作
	BatchMarkAsRead(ctx context.Context, ids []string) error
	MarkAllAsReadForUser(ctx context.Context, userID string) error
	BatchDelete(ctx context.Context, ids []string) error
	DeleteAllForUser(ctx context.Context, userID string) error
	DeleteReadForUser(ctx context.Context, userID string) (int64, error)

	// 统计操作
	CountUnread(ctx context.Context, userID string) (int64, error)
	GetStats(ctx context.Context, userID string) (*notification.NotificationStats, error)
	GetUnreadByType(ctx context.Context, userID string, notificationType notification.NotificationType) ([]*notification.Notification, error)

	// 清理操作
	DeleteExpired(ctx context.Context) (int64, error)
	DeleteOldNotifications(ctx context.Context, beforeDate time.Time) (int64, error)
}

// NotificationPreferenceRepository 通知偏好设置仓储接口
type NotificationPreferenceRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, preference *notification.NotificationPreference) error
	GetByID(ctx context.Context, id string) (*notification.NotificationPreference, error)
	GetByUserID(ctx context.Context, userID string) (*notification.NotificationPreference, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, userID string) (bool, error)

	// 批量操作
	BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error
}

// PushDeviceRepository 推送设备仓储接口
type PushDeviceRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, device *notification.PushDevice) error
	GetByID(ctx context.Context, id string) (*notification.PushDevice, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)

	// 查询操作
	GetByUserID(ctx context.Context, userID string) ([]*notification.PushDevice, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*notification.PushDevice, error)
	GetActiveByUserID(ctx context.Context, userID string) ([]*notification.PushDevice, error)

	// 批量操作
	BatchDelete(ctx context.Context, ids []string) error
	DeactivateAllForUser(ctx context.Context, userID string) error
	UpdateLastUsed(ctx context.Context, id string) error

	// 清理操作
	DeleteInactiveDevices(ctx context.Context, beforeDate time.Time) (int64, error)
}

// NotificationTemplateRepository 通知模板仓储接口
type NotificationTemplateRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, template *notification.NotificationTemplate) error
	GetByID(ctx context.Context, id string) (*notification.NotificationTemplate, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)

	// 查询操作
	List(ctx context.Context, filter *TemplateFilter) ([]*notification.NotificationTemplate, error)
	GetByTypeAndAction(ctx context.Context, templateType notification.NotificationType, action string) ([]*notification.NotificationTemplate, error)
	GetActiveTemplate(ctx context.Context, templateType notification.NotificationType, action string, language string) (*notification.NotificationTemplate, error)
}

// TemplateFilter 模板筛选条件
type TemplateFilter struct {
	Type     *notification.NotificationType `json:"type,omitempty"`
	Action   *string                        `json:"action,omitempty"`
	Language *string                        `json:"language,omitempty"`
	IsActive *bool                          `json:"isActive,omitempty"`
	Limit    int                            `json:"limit,omitempty"`
	Offset   int                            `json:"offset,omitempty"`
}
