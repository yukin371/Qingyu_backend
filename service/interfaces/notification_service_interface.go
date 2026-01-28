package interfaces

import (
	"context"

	notificationModel "Qingyu_backend/models/notification"
)

// NotificationService 通知服务接口
type NotificationService interface {
	// =========================
	// 通知基础操作
	// =========================
	CreateNotification(ctx context.Context, req *CreateNotificationRequest) (*notificationModel.Notification, error)
	GetNotification(ctx context.Context, id string) (*notificationModel.Notification, error)
	GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	GetNotificationStats(ctx context.Context, userID string) (*notificationModel.NotificationStats, error)

	// =========================
	// 通知状态操作
	// =========================
	MarkAsRead(ctx context.Context, id string, userID string) error
	MarkMultipleAsRead(ctx context.Context, ids []string, userID string) error
	MarkMultipleAsReadWithResult(ctx context.Context, ids []string, userID string) (succeeded int, failed int, err error)
	MarkAllAsRead(ctx context.Context, userID string) error
	DeleteNotification(ctx context.Context, id string, userID string) error
	BatchDeleteNotifications(ctx context.Context, ids []string, userID string) error
	BatchDeleteNotificationsWithResult(ctx context.Context, ids []string, userID string) (succeeded int, failed int, err error)
	DeleteAllNotifications(ctx context.Context, userID string) error

	// =========================
	// 通知发送
	// =========================
	SendNotification(ctx context.Context, userID string, notificationType notificationModel.NotificationType, title, content string, data map[string]interface{}) error
	SendNotificationWithTemplate(ctx context.Context, userID string, notificationType notificationModel.NotificationType, action string, variables map[string]interface{}) error
	BatchSendNotification(ctx context.Context, userIDs []string, notificationType notificationModel.NotificationType, title, content string, data map[string]interface{}) error

	// =========================
	// 通知偏好设置
	// =========================
	GetNotificationPreference(ctx context.Context, userID string) (*notificationModel.NotificationPreference, error)
	UpdateNotificationPreference(ctx context.Context, userID string, req *UpdateNotificationPreferenceRequest) error
	ResetNotificationPreference(ctx context.Context, userID string) error

	// =========================
	// 邮件通知设置
	// =========================
	GetEmailNotificationSettings(ctx context.Context, userID string) (*notificationModel.EmailNotificationSettings, error)
	UpdateEmailNotificationSettings(ctx context.Context, userID string, settings *notificationModel.EmailNotificationSettings) error

	// =========================
	// 短信通知设置
	// =========================
	GetSMSNotificationSettings(ctx context.Context, userID string) (*notificationModel.SMSNotificationSettings, error)
	UpdateSMSNotificationSettings(ctx context.Context, userID string, settings *notificationModel.SMSNotificationSettings) error

	// =========================
	// 推送设备管理
	// =========================
	RegisterPushDevice(ctx context.Context, req *RegisterPushDeviceRequest) (*notificationModel.PushDevice, error)
	UnregisterPushDevice(ctx context.Context, deviceID string, userID string) error
	GetUserPushDevices(ctx context.Context, userID string) ([]*notificationModel.PushDevice, error)
	UpdatePushDeviceToken(ctx context.Context, deviceID, userID, deviceToken string) error

	// =========================
	// 定期清理
	// =========================
	CleanupExpiredNotifications(ctx context.Context) (int64, error)
	CleanupOldNotifications(ctx context.Context, days int) (int64, error)
	ClearReadNotifications(ctx context.Context, userID string) (int64, error)
	ResendNotification(ctx context.Context, id string, userID string, method string) error
}

// =========================
// 请求/响应类型定义
// =========================

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	UserID    string                                 `json:"userId" validate:"required"`
	Type      notificationModel.NotificationType     `json:"type" validate:"required"`
	Priority  notificationModel.NotificationPriority `json:"priority"`
	Title     string                                 `json:"title" validate:"required,min=1,max=200"`
	Content   string                                 `json:"content" validate:"required,min=1,max=5000"`
	Data      map[string]interface{}                 `json:"data"`
	ExpiresIn *int                                   `json:"expiresIn"` // 过期时间（秒）
}

// GetNotificationsRequest 获取通知列表请求
type GetNotificationsRequest struct {
	UserID   string                                  `json:"userId" validate:"required"`
	Type     *notificationModel.NotificationType     `json:"type"`
	Read     *bool                                   `json:"read"`
	Priority *notificationModel.NotificationPriority `json:"priority"`
	Keyword  *string                                 `json:"keyword"`
	Limit    int                                     `json:"limit" validate:"min=1,max=100"`
	Offset   int                                     `json:"offset" validate:"min=0"`
	SortBy   string                                  `json:"sortBy" validate:"omitempty,oneof=created_at priority read_at"`
	SortDesc bool                                    `json:"sortDesc"`
}

// GetNotificationsResponse 获取通知列表响应
type GetNotificationsResponse struct {
	Notifications []*notificationModel.Notification `json:"notifications"`
	Total         int64                             `json:"total"`
	UnreadCount   int64                             `json:"unreadCount"`
}

// UpdateNotificationPreferenceRequest 更新通知偏好设置请求
type UpdateNotificationPreferenceRequest struct {
	EnableSystem      *bool                                        `json:"enableSystem"`
	EnableSocial      *bool                                        `json:"enableSocial"`
	EnableContent     *bool                                        `json:"enableContent"`
	EnableReward      *bool                                        `json:"enableReward"`
	EnableMessage     *bool                                        `json:"enableMessage"`
	EnableUpdate      *bool                                        `json:"enableUpdate"`
	EnableMembership  *bool                                        `json:"enableMembership"`
	EmailNotification *notificationModel.EmailNotificationSettings `json:"emailNotification"`
	SMSNotification   *notificationModel.SMSNotificationSettings   `json:"smsNotification"`
	PushNotification  *bool                                        `json:"pushNotification"`
	QuietHoursStart   *string                                      `json:"quietHoursStart" validate:"omitempty"`
	QuietHoursEnd     *string                                      `json:"quietHoursEnd" validate:"omitempty"`
}

// RegisterPushDeviceRequest 注册推送设备请求
type RegisterPushDeviceRequest struct {
	UserID      string `json:"userId" validate:"required"`
	DeviceType  string `json:"deviceType" validate:"required,oneof=ios android web"`
	DeviceToken string `json:"deviceToken" validate:"required"`
	DeviceID    string `json:"deviceId" validate:"required"`
}
