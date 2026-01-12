package notification

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType 通知类型
type NotificationType string

const (
	// NotificationTypeSystem 系统通知 - 平台公告、活动通知
	NotificationTypeSystem NotificationType = "system"
	// NotificationTypeSocial 社交通知 - 关注、点赞、评论通知
	NotificationTypeSocial NotificationType = "social"
	// NotificationTypeContent 内容通知 - 作品审核、上架、下架通知
	NotificationTypeContent NotificationType = "content"
	// NotificationTypeReward 打赏通知 - 收到打赏通知
	NotificationTypeReward NotificationType = "reward"
	// NotificationTypeMessage 私信通知 - 收到私信通知
	NotificationTypeMessage NotificationType = "message"
	// NotificationTypeUpdate 更新通知 - 关注作品更新通知
	NotificationTypeUpdate NotificationType = "update"
	// NotificationTypeMembership 会员通知 - 会员到期、续费提醒
	NotificationTypeMembership NotificationType = "membership"
)

// IsValid 验证通知类型是否有效
func (nt NotificationType) IsValid() bool {
	switch nt {
	case NotificationTypeSystem, NotificationTypeSocial, NotificationTypeContent,
		NotificationTypeReward, NotificationTypeMessage, NotificationTypeUpdate, NotificationTypeMembership:
		return true
	default:
		return false
	}
}

// NotificationPriority 通知优先级
type NotificationPriority string

const (
	// NotificationPriorityLow 低优先级
	NotificationPriorityLow NotificationPriority = "low"
	// NotificationPriorityNormal 普通优先级
	NotificationPriorityNormal NotificationPriority = "normal"
	// NotificationPriorityHigh 高优先级
	NotificationPriorityHigh NotificationPriority = "high"
	// NotificationPriorityUrgent 紧急
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// Notification 通知模型
type Notification struct {
	ID        string                 `json:"id" bson:"_id"`
	UserID    string                 `json:"userId" bson:"user_id"`
	Type      NotificationType       `json:"type" bson:"type"`
	Priority  NotificationPriority   `json:"priority" bson:"priority"`
	Title     string                 `json:"title" bson:"title"`
	Content   string                 `json:"content" bson:"content"`
	Data      map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
	Read      bool                   `json:"read" bson:"read"`
	ReadAt    *time.Time             `json:"readAt,omitempty" bson:"read_at,omitempty"`
	CreatedAt time.Time              `json:"createdAt" bson:"created_at"`
	ExpiresAt *time.Time             `json:"expiresAt,omitempty" bson:"expires_at,omitempty"`
}

// NotificationFilter 通知筛选条件
type NotificationFilter struct {
	UserID    *string               `json:"userId,omitempty"`
	Type      *NotificationType     `json:"type,omitempty"`
	Read      *bool                 `json:"read,omitempty"`
	Priority  *NotificationPriority `json:"priority,omitempty"`
	StartDate *time.Time            `json:"startDate,omitempty"`
	EndDate   *time.Time            `json:"endDate,omitempty"`
	Keyword   *string               `json:"keyword,omitempty"`
	SortBy    string                `json:"sortBy,omitempty"`    // created_at, priority, read_at
	SortOrder string                `json:"sortOrder,omitempty"` // asc, desc
	Limit     int                   `json:"limit,omitempty"`
	Offset    int                   `json:"offset,omitempty"`
}

// NotificationPreference 通知偏好设置
type NotificationPreference struct {
	ID                string                    `json:"id" bson:"_id"`
	UserID            string                    `json:"userId" bson:"user_id"`
	EnableSystem      bool                      `json:"enableSystem" bson:"enable_system"`
	EnableSocial      bool                      `json:"enableSocial" bson:"enable_social"`
	EnableContent     bool                      `json:"enableContent" bson:"enable_content"`
	EnableReward      bool                      `json:"enableReward" bson:"enable_reward"`
	EnableMessage     bool                      `json:"enableMessage" bson:"enable_message"`
	EnableUpdate      bool                      `json:"enableUpdate" bson:"enable_update"`
	EnableMembership  bool                      `json:"enableMembership" bson:"enable_membership"`
	EmailNotification EmailNotificationSettings `json:"emailNotification" bson:"email_notification"`
	SMSNotification   SMSNotificationSettings   `json:"smsNotification" bson:"sms_notification"`
	PushNotification  bool                      `json:"pushNotification" bson:"push_notification"`
	QuietHoursStart   *string                   `json:"quietHoursStart,omitempty" bson:"quiet_hours_start,omitempty"` // HH:MM格式
	QuietHoursEnd     *string                   `json:"quietHoursEnd,omitempty" bson:"quiet_hours_end,omitempty"`     // HH:MM格式
	CreatedAt         time.Time                 `json:"createdAt" bson:"created_at"`
	UpdatedAt         time.Time                 `json:"updatedAt" bson:"updated_at"`
}

// EmailNotificationSettings 邮件通知设置
type EmailNotificationSettings struct {
	Enabled   bool     `json:"enabled" bson:"enabled"`
	Types     []string `json:"types" bson:"types"`         // 启用邮件通知的通知类型列表
	Frequency string   `json:"frequency" bson:"frequency"` // immediate, hourly, daily
}

// SMSNotificationSettings 短信通知设置
type SMSNotificationSettings struct {
	Enabled bool     `json:"enabled" bson:"enabled"`
	Types   []string `json:"types" bson:"types"` // 启用短信通知的通知类型列表
}

// PushDevice 推送设备
type PushDevice struct {
	ID          string    `json:"id" bson:"_id"`
	UserID      string    `json:"userId" bson:"user_id"`
	DeviceType  string    `json:"deviceType" bson:"device_type"` // ios, android, web
	DeviceToken string    `json:"deviceToken" bson:"device_token"`
	DeviceID    string    `json:"deviceId" bson:"device_id"` // 设备唯一标识
	IsActive    bool      `json:"isActive" bson:"is_active"`
	LastUsedAt  time.Time `json:"lastUsedAt" bson:"last_used_at"`
	CreatedAt   time.Time `json:"createdAt" bson:"created_at"`
}

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	ID        string                 `json:"id" bson:"_id"`
	Type      NotificationType       `json:"type" bson:"type"`
	Action    string                 `json:"action" bson:"action"` // follow, like, comment, review, etc.
	Title     string                 `json:"title" bson:"title"`
	Content   string                 `json:"content" bson:"content"`
	Variables []string               `json:"variables" bson:"variables"` // 模板变量列表
	Data      map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"`
	Language  string                 `json:"language" bson:"language"` // zh-CN, en-US, etc.
	IsActive  bool                   `json:"isActive" bson:"is_active"`
	CreatedAt time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updated_at"`
}

// NotificationStats 通知统计
type NotificationStats struct {
	TotalCount     int64                          `json:"totalCount"`
	UnreadCount    int64                          `json:"unreadCount"`
	TypeCounts     map[NotificationType]int64     `json:"typeCounts"`
	PriorityCounts map[NotificationPriority]int64 `json:"priorityCounts"`
}

// NewNotification 创建新通知
func NewNotification(userID string, notificationType NotificationType, title, content string) *Notification {
	now := time.Now()
	return &Notification{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		Type:      notificationType,
		Priority:  NotificationPriorityNormal,
		Title:     title,
		Content:   content,
		Data:      make(map[string]interface{}),
		Read:      false,
		CreatedAt: now,
	}
}

// MarkAsRead 标记为已读
func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.Read = true
	n.ReadAt = &now
}

// IsExpired 检查通知是否已过期
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*n.ExpiresAt)
}

// NewNotificationPreference 创建默认通知偏好设置
func NewNotificationPreference(userID string) *NotificationPreference {
	now := time.Now()
	return &NotificationPreference{
		ID:               primitive.NewObjectID().Hex(),
		UserID:           userID,
		EnableSystem:     true,
		EnableSocial:     true,
		EnableContent:    true,
		EnableReward:     true,
		EnableMessage:    true,
		EnableUpdate:     true,
		EnableMembership: true,
		EmailNotification: EmailNotificationSettings{
			Enabled:   false,
			Types:     []string{},
			Frequency: "immediate",
		},
		SMSNotification: SMSNotificationSettings{
			Enabled: false,
			Types:   []string{},
		},
		PushNotification: true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// IsTypeEnabled 检查指定类型的通知是否启用
func (np *NotificationPreference) IsTypeEnabled(notificationType NotificationType) bool {
	switch notificationType {
	case NotificationTypeSystem:
		return np.EnableSystem
	case NotificationTypeSocial:
		return np.EnableSocial
	case NotificationTypeContent:
		return np.EnableContent
	case NotificationTypeReward:
		return np.EnableReward
	case NotificationTypeMessage:
		return np.EnableMessage
	case NotificationTypeUpdate:
		return np.EnableUpdate
	case NotificationTypeMembership:
		return np.EnableMembership
	default:
		return true
	}
}

// IsEmailEnabledForType 检查指定类型的邮件通知是否启用
func (np *NotificationPreference) IsEmailEnabledForType(notificationType NotificationType) bool {
	if !np.EmailNotification.Enabled {
		return false
	}
	for _, t := range np.EmailNotification.Types {
		if t == string(notificationType) {
			return true
		}
	}
	return false
}

// IsSMSEnabledForType 检查指定类型的短信通知是否启用
func (np *NotificationPreference) IsSMSEnabledForType(notificationType NotificationType) bool {
	if !np.SMSNotification.Enabled {
		return false
	}
	for _, t := range np.SMSNotification.Types {
		if t == string(notificationType) {
			return true
		}
	}
	return false
}

// NewPushDevice 创建新的推送设备
func NewPushDevice(userID, deviceType, deviceToken, deviceID string) *PushDevice {
	now := time.Now()
	return &PushDevice{
		ID:          primitive.NewObjectID().Hex(),
		UserID:      userID,
		DeviceType:  deviceType,
		DeviceToken: deviceToken,
		DeviceID:    deviceID,
		IsActive:    true,
		LastUsedAt:  now,
		CreatedAt:   now,
	}
}
