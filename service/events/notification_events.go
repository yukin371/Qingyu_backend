package events

import (
	"context"
	"log"
	"time"

	"Qingyu_backend/service/base"
)

// ============ 通知相关事件 ============

// 通知事件类型常量
const (
	EventNotificationCreated = "notification.created"
	EventNotificationSent    = "notification.sent"
	EventNotificationRead    = "notification.read"
	EventNotificationDeleted = "notification.deleted"
)

// NotificationEventData 通知事件数据
type NotificationEventData struct {
	NotificationID   string                 `json:"notification_id"`
	UserID           string                 `json:"user_id"`
	Type             string                 `json:"type"`             // system/comment/follow/reward/purchase
	Title            string                 `json:"title"`
	Content          string                 `json:"content"`
	Action           string                 `json:"action"`
	TargetType       string                 `json:"target_type,omitempty"`
	TargetID         string                 `json:"target_id,omitempty"`
	Priority         string                 `json:"priority,omitempty"` // low/normal/high/urgent
	Channels         []string               `json:"channels,omitempty"` // in_app/email/sms/push
	Time             time.Time              `json:"time"`
	ReadTime         time.Time              `json:"read_time,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// ============ 通知事件工厂函数 ============

// NewNotificationCreatedEvent 创建通知事件
func NewNotificationCreatedEvent(notificationID, userID, notificationType, title, content string) base.Event {
	return &base.BaseEvent{
		EventType: EventNotificationCreated,
		EventData: NotificationEventData{
			NotificationID: notificationID,
			UserID:         userID,
			Type:           notificationType,
			Title:          title,
			Content:        content,
			Action:         "created",
			Priority:       "normal",
			Channels:       []string{"in_app"},
			Time:           time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "NotificationService",
	}
}

// NewNotificationSentEvent 创建通知发送事件
func NewNotificationSentEvent(notificationID, userID string, channels []string) base.Event {
	return &base.BaseEvent{
		EventType: EventNotificationSent,
		EventData: NotificationEventData{
			NotificationID: notificationID,
			UserID:         userID,
			Action:         "sent",
			Channels:       channels,
			Time:           time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "NotificationService",
	}
}

// NewNotificationReadEvent 创建通知已读事件
func NewNotificationReadEvent(notificationID, userID string) base.Event {
	return &base.BaseEvent{
		EventType: EventNotificationRead,
		EventData: NotificationEventData{
			NotificationID: notificationID,
			UserID:         userID,
			Action:         "read",
			Time:           time.Now(),
			ReadTime:       time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "NotificationService",
	}
}

// NewNotificationDeletedEvent 创建通知删除事件
func NewNotificationDeletedEvent(notificationID, userID string) base.Event {
	return &base.BaseEvent{
		EventType: EventNotificationDeleted,
		EventData: NotificationEventData{
			NotificationID: notificationID,
			UserID:         userID,
			Action:         "deleted",
			Time:           time.Now(),
		},
		Timestamp: time.Now(),
		Source:    "NotificationService",
	}
}

// ============ 事件处理器 ============

// EmailNotificationHandler 邮件通知处理器
// 发送邮件通知
type EmailNotificationHandler struct {
	name string
}

// NewEmailNotificationHandler 创建邮件通知处理器
func NewEmailNotificationHandler() *EmailNotificationHandler {
	return &EmailNotificationHandler{
		name: "EmailNotificationHandler",
	}
}

// Handle 处理事件
func (h *EmailNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	data, ok := event.GetEventData().(NotificationEventData)
	if !ok {
		return nil
	}

	// 检查是否需要发送邮件通知
	needEmail := false
	for _, channel := range data.Channels {
		if channel == "email" {
			needEmail = true
			break
		}
	}

	if !needEmail {
		return nil
	}

	log.Printf("[EmailNotification] 发送邮件给用户 %s: %s - %s", data.UserID, data.Title, data.Content)
	// 实际项目中这里应该调用邮件服务
	// emailService.Send(data.UserID, data.Title, data.Content)

	return nil
}

// GetHandlerName 获取处理器名称
func (h *EmailNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *EmailNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventNotificationCreated,
	}
}

// SMSNotificationHandler 短信通知处理器
// 发送短信通知
type SMSNotificationHandler struct {
	name string
}

// NewSMSNotificationHandler 创建短信通知处理器
func NewSMSNotificationHandler() *SMSNotificationHandler {
	return &SMSNotificationHandler{
		name: "SMSNotificationHandler",
	}
}

// Handle 处理事件
func (h *SMSNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	data, ok := event.GetEventData().(NotificationEventData)
	if !ok {
		return nil
	}

	// 检查是否需要发送短信通知
	needSMS := false
	for _, channel := range data.Channels {
		if channel == "sms" {
			needSMS = true
			break
		}
	}

	if !needSMS {
		return nil
	}

	// 只发送重要通知的短信
	if data.Priority == "high" || data.Priority == "urgent" {
		log.Printf("[SMSNotification] 发送短信给用户 %s: %s", data.UserID, data.Title)
		// 实际项目中这里应该调用短信服务
		// smsService.Send(data.UserID, data.Content)
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *SMSNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *SMSNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventNotificationCreated,
	}
}

// PushNotificationHandler 推送通知处理器
// 发送APP推送通知
type PushNotificationHandler struct {
	name string
}

// NewPushNotificationHandler 创建推送通知处理器
func NewPushNotificationHandler() *PushNotificationHandler {
	return &PushNotificationHandler{
		name: "PushNotificationHandler",
	}
}

// Handle 处理事件
func (h *PushNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	data, ok := event.GetEventData().(NotificationEventData)
	if !ok {
		return nil
	}

	// 检查是否需要发送推送通知
	needPush := false
	for _, channel := range data.Channels {
		if channel == "push" {
			needPush = true
			break
		}
	}

	if !needPush {
		return nil
	}

	log.Printf("[PushNotification] 发送推送通知给用户 %s: %s - %s", data.UserID, data.Title, data.Content)
	// 实际项目中这里应该调用推送服务
	// pushService.Send(data.UserID, data.Title, data.Content, data.TargetType, data.TargetID)

	return nil
}

// GetHandlerName 获取处理器名称
func (h *PushNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *PushNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventNotificationCreated,
	}
}

// InAppNotificationHandler 站内通知处理器
// 存储站内通知
type InAppNotificationHandler struct {
	name string
}

// NewInAppNotificationHandler 创建站内通知处理器
func NewInAppNotificationHandler() *InAppNotificationHandler {
	return &InAppNotificationHandler{
		name: "InAppNotificationHandler",
	}
}

// Handle 处理事件
func (h *InAppNotificationHandler) Handle(ctx context.Context, event base.Event) error {
	data, ok := event.GetEventData().(NotificationEventData)
	if !ok {
		return nil
	}

	// 检查是否需要存储站内通知
	needInApp := false
	for _, channel := range data.Channels {
		if channel == "in_app" {
			needInApp = true
			break
		}
	}

	if !needInApp {
		return nil
	}

	log.Printf("[InAppNotification] 创建站内通知给用户 %s: %s", data.UserID, data.Title)
	// 实际项目中这里应该存储到数据库
	// notificationRepo.Create(&Notification{
	//     UserID: data.UserID,
	//     Type:   data.Type,
	//     Title:  data.Title,
	//     Content: data.Content,
	//     TargetType: data.TargetType,
	//     TargetID: data.TargetID,
	//     Priority: data.Priority,
	//     Read: false,
	// })

	return nil
}

// GetHandlerName 获取处理器名称
func (h *InAppNotificationHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *InAppNotificationHandler) GetSupportedEventTypes() []string {
	return []string{
		EventNotificationCreated,
	}
}

// NotificationStatisticsHandler 通知统计处理器
// 统计通知数据
type NotificationStatisticsHandler struct {
	name string
}

// NewNotificationStatisticsHandler 创建通知统计处理器
func NewNotificationStatisticsHandler() *NotificationStatisticsHandler {
	return &NotificationStatisticsHandler{
		name: "NotificationStatisticsHandler",
	}
}

// Handle 处理事件
func (h *NotificationStatisticsHandler) Handle(ctx context.Context, event base.Event) error {
	switch event.GetEventType() {
	case EventNotificationCreated:
		data, _ := event.GetEventData().(NotificationEventData)
		log.Printf("[NotificationStatistics] 用户 %s 创建通知，类型: %s", data.UserID, data.Type)
		// 统计通知创建数量

	case EventNotificationSent:
		data, _ := event.GetEventData().(NotificationEventData)
		log.Printf("[NotificationStatistics] 用户 %s 通知已发送，渠道: %v", data.UserID, data.Channels)
		// 统计通知发送数量

	case EventNotificationRead:
		data, _ := event.GetEventData().(NotificationEventData)
		log.Printf("[NotificationStatistics] 用户 %s 通知已读", data.UserID)
		// 统计通知已读数量
	}

	return nil
}

// GetHandlerName 获取处理器名称
func (h *NotificationStatisticsHandler) GetHandlerName() string {
	return h.name
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *NotificationStatisticsHandler) GetSupportedEventTypes() []string {
	return []string{
		EventNotificationCreated,
		EventNotificationSent,
		EventNotificationRead,
	}
}
