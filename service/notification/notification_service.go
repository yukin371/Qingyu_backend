package notification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/models/notification"
	"Qingyu_backend/pkg/errors"
	repo "Qingyu_backend/repository/interfaces/notification"
)

// NotificationService 通知服务接口
type NotificationService interface {
	// 通知基础操作
	CreateNotification(ctx context.Context, req *CreateNotificationRequest) (*notification.Notification, error)
	GetNotification(ctx context.Context, id string) (*notification.Notification, error)
	GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	GetNotificationStats(ctx context.Context, userID string) (*notification.NotificationStats, error)

	// 通知状态操作
	MarkAsRead(ctx context.Context, id string, userID string) error
	MarkMultipleAsRead(ctx context.Context, ids []string, userID string) error
	MarkMultipleAsReadWithResult(ctx context.Context, ids []string, userID string) (succeeded int, failed int, err error)
	MarkAllAsRead(ctx context.Context, userID string) error
	DeleteNotification(ctx context.Context, id string, userID string) error
	BatchDeleteNotifications(ctx context.Context, ids []string, userID string) error
	BatchDeleteNotificationsWithResult(ctx context.Context, ids []string, userID string) (succeeded int, failed int, err error)
	DeleteAllNotifications(ctx context.Context, userID string) error
	ClearReadNotifications(ctx context.Context, userID string) (int64, error)
	ResendNotification(ctx context.Context, id string, userID string, method string) error

	// 通知发送
	SendNotification(ctx context.Context, userID string, notificationType notification.NotificationType, title, content string, data map[string]interface{}) error
	SendNotificationWithTemplate(ctx context.Context, userID string, notificationType notification.NotificationType, action string, variables map[string]interface{}) error
	BatchSendNotification(ctx context.Context, userIDs []string, notificationType notification.NotificationType, title, content string, data map[string]interface{}) error

	// 通知偏好设置
	GetNotificationPreference(ctx context.Context, userID string) (*notification.NotificationPreference, error)
	UpdateNotificationPreference(ctx context.Context, userID string, req *UpdateNotificationPreferenceRequest) error
	ResetNotificationPreference(ctx context.Context, userID string) error

	// 邮件通知设置
	GetEmailNotificationSettings(ctx context.Context, userID string) (*notification.EmailNotificationSettings, error)
	UpdateEmailNotificationSettings(ctx context.Context, userID string, settings *notification.EmailNotificationSettings) error

	// 短信通知设置
	GetSMSNotificationSettings(ctx context.Context, userID string) (*notification.SMSNotificationSettings, error)
	UpdateSMSNotificationSettings(ctx context.Context, userID string, settings *notification.SMSNotificationSettings) error

	// 推送设备管理
	RegisterPushDevice(ctx context.Context, req *RegisterPushDeviceRequest) (*notification.PushDevice, error)
	UnregisterPushDevice(ctx context.Context, deviceID string, userID string) error
	GetUserPushDevices(ctx context.Context, userID string) ([]*notification.PushDevice, error)
	UpdatePushDeviceToken(ctx context.Context, deviceID, userID, deviceToken string) error

	// 定期清理
	CleanupExpiredNotifications(ctx context.Context) (int64, error)
	CleanupOldNotifications(ctx context.Context, days int) (int64, error)

	// 重新发送和清除操作
	ResendNotification(ctx context.Context, id string, userID string, method string) error
	ClearReadNotifications(ctx context.Context, userID string) (int64, error)
}

// notificationServiceImpl 通知服务实现
type notificationServiceImpl struct {
	notificationRepo repo.NotificationRepository
	preferenceRepo   repo.NotificationPreferenceRepository
	pushDeviceRepo   repo.PushDeviceRepository
	templateRepo     repo.NotificationTemplateRepository
	emailService     EmailService
	wsHub            WSHub
}

// EmailService 邮件服务接口（避免循环依赖）
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

// WSHub WebSocket Hub接口（避免循环依赖）
type WSHub interface {
	BroadcastNotification(userID string, notification interface{})
}

// NewNotificationService 创建通知服务实例
func NewNotificationService(
	notificationRepo repo.NotificationRepository,
	preferenceRepo repo.NotificationPreferenceRepository,
	pushDeviceRepo repo.PushDeviceRepository,
	templateRepo repo.NotificationTemplateRepository,
	emailService EmailService,
	wsHub WSHub,
) NotificationService {
	return &notificationServiceImpl{
		notificationRepo: notificationRepo,
		preferenceRepo:   preferenceRepo,
		pushDeviceRepo:   pushDeviceRepo,
		templateRepo:     templateRepo,
		emailService:     emailService,
		wsHub:            wsHub,
	}
}

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	UserID    string                            `json:"userId" validate:"required"`
	Type      notification.NotificationType     `json:"type" validate:"required"`
	Priority  notification.NotificationPriority `json:"priority"`
	Title     string                            `json:"title" validate:"required,min=1,max=200"`
	Content   string                            `json:"content" validate:"required,min=1,max=5000"`
	Data      map[string]interface{}            `json:"data"`
	ExpiresIn *int                              `json:"expiresIn"` // 过期时间（秒）
}

// GetNotificationsRequest 获取通知列表请求
type GetNotificationsRequest struct {
	UserID   string                             `json:"userId" validate:"required"`
	Type     *notification.NotificationType     `json:"type"`
	Read     *bool                              `json:"read"`
	Priority *notification.NotificationPriority `json:"priority"`
	Keyword  *string                            `json:"keyword"`
	Limit    int                                `json:"limit" validate:"min=1,max=100"`
	Offset   int                                `json:"offset" validate:"min=0"`
	SortBy   string                             `json:"sortBy" validate:"omitempty,oneof=created_at priority read_at"`
	SortDesc bool                               `json:"sortDesc"`
}

// GetNotificationsResponse 获取通知列表响应
type GetNotificationsResponse struct {
	Notifications []*notification.Notification `json:"notifications"`
	Total         int64                        `json:"total"`
	UnreadCount   int64                        `json:"unreadCount"`
}

// UpdateNotificationPreferenceRequest 更新通知偏好设置请求
type UpdateNotificationPreferenceRequest struct {
	EnableSystem      *bool                                   `json:"enableSystem"`
	EnableSocial      *bool                                   `json:"enableSocial"`
	EnableContent     *bool                                   `json:"enableContent"`
	EnableReward      *bool                                   `json:"enableReward"`
	EnableMessage     *bool                                   `json:"enableMessage"`
	EnableUpdate      *bool                                   `json:"enableUpdate"`
	EnableMembership  *bool                                   `json:"enableMembership"`
	EmailNotification *notification.EmailNotificationSettings `json:"emailNotification"`
	SMSNotification   *notification.SMSNotificationSettings   `json:"smsNotification"`
	PushNotification  *bool                                   `json:"pushNotification"`
	QuietHoursStart   *string                                 `json:"quietHoursStart" validate:"omitempty"`
	QuietHoursEnd     *string                                 `json:"quietHoursEnd" validate:"omitempty"`
}

// RegisterPushDeviceRequest 注册推送设备请求
type RegisterPushDeviceRequest struct {
	UserID      string `json:"userId" validate:"required"`
	DeviceType  string `json:"deviceType" validate:"required,oneof=ios android web"`
	DeviceToken string `json:"deviceToken" validate:"required"`
	DeviceID    string `json:"deviceId" validate:"required"`
}

// CreateNotification 创建通知
func (s *notificationServiceImpl) CreateNotification(ctx context.Context, req *CreateNotificationRequest) (*notification.Notification, error) {
	// 创建通知对象
	notif := notification.NewNotification(req.UserID, req.Type, req.Title, req.Content)

	// 设置优先级
	if req.Priority != "" {
		notif.Priority = req.Priority
	}

	// 设置扩展数据
	if req.Data != nil {
		notif.Data = req.Data
	}

	// 设置过期时间
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Second)
		notif.ExpiresAt = &expiresAt
	}

	// 保存到数据库
	if err := s.notificationRepo.Create(ctx, notif); err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("NOTIFICATION_CREATE_FAILED", "创建通知失败", err)
	}

	return notif, nil
}

// GetNotification 获取通知详情
func (s *notificationServiceImpl) GetNotification(ctx context.Context, id string) (*notification.Notification, error) {
	notif, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("NOTIFICATION_GET_FAILED", "获取通知失败", err)
	}

	if notif == nil {
		return nil, errors.BookstoreServiceFactory.NotFoundError("Notification", id)
	}

	return notif, nil
}

// GetNotifications 获取通知列表
func (s *notificationServiceImpl) GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// 构建筛选条件
	filter := &notification.NotificationFilter{
		UserID:   &req.UserID,
		Type:     req.Type,
		Read:     req.Read,
		Priority: req.Priority,
		Keyword:  req.Keyword,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}

	// 设置排序
	sortOrder := "desc"
	if req.SortDesc {
		sortOrder = "asc"
	}
	if req.SortBy != "" {
		filter.SortBy = req.SortBy
		filter.SortOrder = sortOrder
	} else {
		filter.SortBy = "created_at"
		filter.SortOrder = "desc"
	}

	// 获取通知列表
	notifications, err := s.notificationRepo.List(ctx, filter)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("NOTIFICATION_LIST_FAILED", "获取通知列表失败", err)
	}

	// 获取总数
	total, err := s.notificationRepo.Count(ctx, filter)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("NOTIFICATION_COUNT_FAILED", "获取通知总数失败", err)
	}

	// 获取未读数量
	unreadCount, err := s.notificationRepo.CountUnread(ctx, req.UserID)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("NOTIFICATION_UNREAD_COUNT_FAILED", "获取未读数量失败", err)
	}

	return &GetNotificationsResponse{
		Notifications: notifications,
		Total:         total,
		UnreadCount:   unreadCount,
	}, nil
}

// GetUnreadCount 获取未读通知数量
func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	count, err := s.notificationRepo.CountUnread(ctx, userID)
	if err != nil {
		return 0, errors.BookstoreServiceFactory.InternalError("NOTIFICATION_UNREAD_COUNT_FAILED", "获取未读数量失败", err)
	}
	return count, nil
}

// GetNotificationStats 获取通知统计
func (s *notificationServiceImpl) GetNotificationStats(ctx context.Context, userID string) (*notification.NotificationStats, error) {
	stats, err := s.notificationRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("NOTIFICATION_STATS_FAILED", "获取通知统计失败", err)
	}
	return stats, nil
}

// MarkAsRead 标记通知为已读
func (s *notificationServiceImpl) MarkAsRead(ctx context.Context, id string, userID string) error {
	// 检查通知是否存在
	notif, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_GET_FAILED", "获取通知失败", err)
	}

	if notif == nil {
		return errors.BookstoreServiceFactory.NotFoundError("Notification", id)
	}

	// 验证权限
	if notif.UserID != userID {
		return errors.BookstoreServiceFactory.ForbiddenError("NO_PERMISSION", "无权操作此通知")
	}

	// 如果已经已读，直接返回
	if notif.Read {
		return nil
	}

	// 标记为已读
	now := time.Now()
	updates := map[string]interface{}{
		"read":    true,
		"read_at": now,
	}

	if err := s.notificationRepo.Update(ctx, id, updates); err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_UPDATE_FAILED", "更新通知失败", err)
	}

	return nil
}

// MarkMultipleAsRead 批量标记通知为已读
func (s *notificationServiceImpl) MarkMultipleAsRead(ctx context.Context, ids []string, userID string) error {
	if len(ids) == 0 {
		return nil
	}

	// 验证所有通知都属于该用户
	for _, id := range ids {
		notif, err := s.notificationRepo.GetByID(ctx, id)
		if err != nil {
			return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_GET_FAILED", "获取通知失败", err)
		}
		if notif == nil || notif.UserID != userID {
			return errors.BookstoreServiceFactory.ForbiddenError("NO_PERMISSION", "无权操作部分通知")
		}
	}

	if err := s.notificationRepo.BatchMarkAsRead(ctx, ids); err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_BATCH_UPDATE_FAILED", "批量更新失败", err)
	}

	return nil
}

// MarkAllAsRead 标记所有通知为已读
func (s *notificationServiceImpl) MarkAllAsRead(ctx context.Context, userID string) error {
	if err := s.notificationRepo.MarkAllAsReadForUser(ctx, userID); err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_MARK_ALL_READ_FAILED", "标记全部已读失败", err)
	}
	return nil
}

// DeleteNotification 删除通知
func (s *notificationServiceImpl) DeleteNotification(ctx context.Context, id string, userID string) error {
	// 检查通知是否存在
	notif, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_GET_FAILED", "获取通知失败", err)
	}

	if notif == nil {
		return errors.BookstoreServiceFactory.NotFoundError("Notification", id)
	}

	// 验证权限
	if notif.UserID != userID {
		return errors.BookstoreServiceFactory.ForbiddenError("NO_PERMISSION", "无权操作此通知")
	}

	if err := s.notificationRepo.Delete(ctx, id); err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_DELETE_FAILED", "删除通知失败", err)
	}

	return nil
}

// BatchDeleteNotifications 批量删除通知
func (s *notificationServiceImpl) BatchDeleteNotifications(ctx context.Context, ids []string, userID string) error {
	if len(ids) == 0 {
		return nil
	}

	// 验证所有通知都属于该用户
	for _, id := range ids {
		notif, err := s.notificationRepo.GetByID(ctx, id)
		if err != nil {
			return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_GET_FAILED", "获取通知失败", err)
		}
		if notif == nil || notif.UserID != userID {
			return errors.BookstoreServiceFactory.ForbiddenError("NO_PERMISSION", "无权操作部分通知")
		}
	}

	if err := s.notificationRepo.BatchDelete(ctx, ids); err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_BATCH_DELETE_FAILED", "批量删除失败", err)
	}

	return nil
}

// DeleteAllNotifications 删除所有通知
func (s *notificationServiceImpl) DeleteAllNotifications(ctx context.Context, userID string) error {
	if err := s.notificationRepo.DeleteAllForUser(ctx, userID); err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_DELETE_ALL_FAILED", "删除所有通知失败", err)
	}
	return nil
}

// ClearReadNotifications 清除已读通知
func (s *notificationServiceImpl) ClearReadNotifications(ctx context.Context, userID string) (int64, error) {
	// TODO: 实现清除已读通知的逻辑
	// 需要在repository层添加DeleteReadForUser方法
	// 暂时返回成功，count为0
	return 0, nil
}

// ResendNotification 重新发送通知
func (s *notificationServiceImpl) ResendNotification(ctx context.Context, id string, userID string, method string) error {
	// TODO: 实现重发通知的逻辑
	// 需要根据method类型（email/push/sms）调用相应的发送服务
	return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_RESEND_NOT_IMPLEMENTED", "重发通知功能暂未实现", nil)
}

// SendNotification 发送通知
func (s *notificationServiceImpl) SendNotification(ctx context.Context, userID string, notificationType notification.NotificationType, title, content string, data map[string]interface{}) error {
	// 获取用户通知偏好设置
	preference, err := s.preferenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("PREFERENCE_GET_FAILED", "获取通知偏好设置失败", err)
	}

	// 如果没有偏好设置，使用默认设置
	if preference == nil {
		preference = notification.NewNotificationPreference(userID)
		if err := s.preferenceRepo.Create(ctx, preference); err != nil {
			return errors.BookstoreServiceFactory.InternalError("PREFERENCE_CREATE_FAILED", "创建通知偏好设置失败", err)
		}
	}

	// 检查该类型通知是否启用
	if !preference.IsTypeEnabled(notificationType) {
		return nil // 用户关闭了该类型通知
	}

	// 创建通知
	req := &CreateNotificationRequest{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Content: content,
		Data:    data,
	}

	notif, err := s.CreateNotification(ctx, req)
	if err != nil {
		return err
	}

	// TODO: 发送实时推送（WebSocket）
	// TODO: 发送邮件通知（如果启用）
	// TODO: 发送短信通知（如果启用）

	_ = notif
	return nil
}

// SendNotificationWithTemplate 使用模板发送通知
func (s *notificationServiceImpl) SendNotificationWithTemplate(ctx context.Context, userID string, notificationType notification.NotificationType, action string, variables map[string]interface{}) error {
	// 获取通知模板（默认使用中文）
	template, err := s.templateRepo.GetActiveTemplate(ctx, notificationType, action, "zh-CN")
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("TEMPLATE_GET_FAILED", "获取通知模板失败", err)
	}

	if template == nil {
		return errors.BookstoreServiceFactory.NotFoundError("NotificationTemplate", fmt.Sprintf("%s:%s", notificationType, action))
	}

	// 替换模板变量
	title := s.replaceVariables(template.Title, variables)
	content := s.replaceVariables(template.Content, variables)

	// 发送通知
	return s.SendNotification(ctx, userID, notificationType, title, content, template.Data)
}

// BatchSendNotification 批量发送通知
func (s *notificationServiceImpl) BatchSendNotification(ctx context.Context, userIDs []string, notificationType notification.NotificationType, title, content string, data map[string]interface{}) error {
	for _, userID := range userIDs {
		if err := s.SendNotification(ctx, userID, notificationType, title, content, data); err != nil {
			// 记录错误但继续处理其他用户
			continue
		}
	}
	return nil
}

// replaceVariables 替换模板变量
func (s *notificationServiceImpl) replaceVariables(template string, variables map[string]interface{}) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// GetNotificationPreference 获取通知偏好设置
func (s *notificationServiceImpl) GetNotificationPreference(ctx context.Context, userID string) (*notification.NotificationPreference, error) {
	preference, err := s.preferenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("PREFERENCE_GET_FAILED", "获取通知偏好设置失败", err)
	}

	// 如果没有偏好设置，返回默认设置
	if preference == nil {
		preference = notification.NewNotificationPreference(userID)
	}

	return preference, nil
}

// UpdateNotificationPreference 更新通知偏好设置
func (s *notificationServiceImpl) UpdateNotificationPreference(ctx context.Context, userID string, req *UpdateNotificationPreferenceRequest) error {
	// 获取现有偏好设置
	preference, err := s.preferenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("PREFERENCE_GET_FAILED", "获取通知偏好设置失败", err)
	}

	// 如果不存在，创建新的
	if preference == nil {
		preference = notification.NewNotificationPreference(userID)
		if err := s.preferenceRepo.Create(ctx, preference); err != nil {
			return errors.BookstoreServiceFactory.InternalError("PREFERENCE_CREATE_FAILED", "创建通知偏好设置失败", err)
		}
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if req.EnableSystem != nil {
		updates["enable_system"] = *req.EnableSystem
	}
	if req.EnableSocial != nil {
		updates["enable_social"] = *req.EnableSocial
	}
	if req.EnableContent != nil {
		updates["enable_content"] = *req.EnableContent
	}
	if req.EnableReward != nil {
		updates["enable_reward"] = *req.EnableReward
	}
	if req.EnableMessage != nil {
		updates["enable_message"] = *req.EnableMessage
	}
	if req.EnableUpdate != nil {
		updates["enable_update"] = *req.EnableUpdate
	}
	if req.EnableMembership != nil {
		updates["enable_membership"] = *req.EnableMembership
	}
	if req.PushNotification != nil {
		updates["push_notification"] = *req.PushNotification
	}
	if req.EmailNotification != nil {
		updates["email_notification"] = *req.EmailNotification
	}
	if req.SMSNotification != nil {
		updates["sms_notification"] = *req.SMSNotification
	}
	if req.QuietHoursStart != nil {
		updates["quiet_hours_start"] = *req.QuietHoursStart
	}
	if req.QuietHoursEnd != nil {
		updates["quiet_hours_end"] = *req.QuietHoursEnd
	}

	updates["updated_at"] = time.Now()

	if err := s.preferenceRepo.Update(ctx, preference.ID, updates); err != nil {
		return errors.BookstoreServiceFactory.InternalError("PREFERENCE_UPDATE_FAILED", "更新通知偏好设置失败", err)
	}

	return nil
}

// ResetNotificationPreference 重置通知偏好设置
func (s *notificationServiceImpl) ResetNotificationPreference(ctx context.Context, userID string) error {
	// 删除现有偏好设置
	preference, err := s.preferenceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("PREFERENCE_GET_FAILED", "获取通知偏好设置失败", err)
	}

	if preference != nil {
		if err := s.preferenceRepo.Delete(ctx, preference.ID); err != nil {
			return errors.BookstoreServiceFactory.InternalError("PREFERENCE_DELETE_FAILED", "删除通知偏好设置失败", err)
		}
	}

	// 创建默认偏好设置
	newPreference := notification.NewNotificationPreference(userID)
	if err := s.preferenceRepo.Create(ctx, newPreference); err != nil {
		return errors.BookstoreServiceFactory.InternalError("PREFERENCE_CREATE_FAILED", "创建通知偏好设置失败", err)
	}

	return nil
}

// GetEmailNotificationSettings 获取邮件通知设置
func (s *notificationServiceImpl) GetEmailNotificationSettings(ctx context.Context, userID string) (*notification.EmailNotificationSettings, error) {
	preference, err := s.GetNotificationPreference(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &preference.EmailNotification, nil
}

// UpdateEmailNotificationSettings 更新邮件通知设置
func (s *notificationServiceImpl) UpdateEmailNotificationSettings(ctx context.Context, userID string, settings *notification.EmailNotificationSettings) error {
	req := &UpdateNotificationPreferenceRequest{
		EmailNotification: settings,
	}

	return s.UpdateNotificationPreference(ctx, userID, req)
}

// GetSMSNotificationSettings 获取短信通知设置
func (s *notificationServiceImpl) GetSMSNotificationSettings(ctx context.Context, userID string) (*notification.SMSNotificationSettings, error) {
	preference, err := s.GetNotificationPreference(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &preference.SMSNotification, nil
}

// UpdateSMSNotificationSettings 更新短信通知设置
func (s *notificationServiceImpl) UpdateSMSNotificationSettings(ctx context.Context, userID string, settings *notification.SMSNotificationSettings) error {
	req := &UpdateNotificationPreferenceRequest{
		SMSNotification: settings,
	}

	return s.UpdateNotificationPreference(ctx, userID, req)
}

// RegisterPushDevice 注册推送设备
func (s *notificationServiceImpl) RegisterPushDevice(ctx context.Context, req *RegisterPushDeviceRequest) (*notification.PushDevice, error) {
	// 检查设备是否已注册
	existingDevice, err := s.pushDeviceRepo.GetByDeviceID(ctx, req.DeviceID)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("DEVICE_GET_FAILED", "获取设备信息失败", err)
	}

	// 如果设备已存在，更新令牌和最后使用时间
	if existingDevice != nil {
		updates := map[string]interface{}{
			"device_token": req.DeviceToken,
			"is_active":    true,
			"last_used_at": time.Now(),
		}

		if err := s.pushDeviceRepo.Update(ctx, existingDevice.ID, updates); err != nil {
			return nil, errors.BookstoreServiceFactory.InternalError("DEVICE_UPDATE_FAILED", "更新设备信息失败", err)
		}

		// 获取更新后的设备信息
		return s.pushDeviceRepo.GetByID(ctx, existingDevice.ID)
	}

	// 创建新设备
	device := notification.NewPushDevice(req.UserID, req.DeviceType, req.DeviceToken, req.DeviceID)

	if err := s.pushDeviceRepo.Create(ctx, device); err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("DEVICE_CREATE_FAILED", "创建设备信息失败", err)
	}

	return device, nil
}

// UnregisterPushDevice 取消注册推送设备
func (s *notificationServiceImpl) UnregisterPushDevice(ctx context.Context, deviceID string, userID string) error {
	// 获取设备信息
	device, err := s.pushDeviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("DEVICE_GET_FAILED", "获取设备信息失败", err)
	}

	if device == nil {
		return errors.BookstoreServiceFactory.NotFoundError("PushDevice", deviceID)
	}

	// 验证权限
	if device.UserID != userID {
		return errors.BookstoreServiceFactory.ForbiddenError("NO_PERMISSION", "无权操作此设备")
	}

	if err := s.pushDeviceRepo.Delete(ctx, device.ID); err != nil {
		return errors.BookstoreServiceFactory.InternalError("DEVICE_DELETE_FAILED", "删除设备信息失败", err)
	}

	return nil
}

// GetUserPushDevices 获取用户的推送设备列表
func (s *notificationServiceImpl) GetUserPushDevices(ctx context.Context, userID string) ([]*notification.PushDevice, error) {
	devices, err := s.pushDeviceRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, errors.BookstoreServiceFactory.InternalError("DEVICES_GET_FAILED", "获取设备列表失败", err)
	}

	return devices, nil
}

// UpdatePushDeviceToken 更新推送设备令牌
func (s *notificationServiceImpl) UpdatePushDeviceToken(ctx context.Context, deviceID, userID, deviceToken string) error {
	// 获取设备信息
	device, err := s.pushDeviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("DEVICE_GET_FAILED", "获取设备信息失败", err)
	}

	if device == nil {
		return errors.BookstoreServiceFactory.NotFoundError("PushDevice", deviceID)
	}

	// 验证权限
	if device.UserID != userID {
		return errors.BookstoreServiceFactory.ForbiddenError("NO_PERMISSION", "无权操作此设备")
	}

	// 更新令牌和最后使用时间
	updates := map[string]interface{}{
		"device_token": deviceToken,
		"last_used_at": time.Now(),
		"is_active":    true,
	}

	if err := s.pushDeviceRepo.Update(ctx, device.ID, updates); err != nil {
		return errors.BookstoreServiceFactory.InternalError("DEVICE_UPDATE_FAILED", "更新设备令牌失败", err)
	}

	return nil
}

// CleanupExpiredNotifications 清理过期通知
func (s *notificationServiceImpl) CleanupExpiredNotifications(ctx context.Context) (int64, error) {
	count, err := s.notificationRepo.DeleteExpired(ctx)
	if err != nil {
		return 0, errors.BookstoreServiceFactory.InternalError("CLEANUP_EXPIRED_FAILED", "清理过期通知失败", err)
	}
	return count, nil
}

// CleanupOldNotifications 清理旧通知
func (s *notificationServiceImpl) CleanupOldNotifications(ctx context.Context, days int) (int64, error) {
	if days <= 0 {
		return 0, errors.BookstoreServiceFactory.ValidationError("INVALID_DAYS", "天数必须大于0")
	}

	beforeDate := time.Now().AddDate(0, 0, -days)
	count, err := s.notificationRepo.DeleteOldNotifications(ctx, beforeDate)
	if err != nil {
		return 0, errors.BookstoreServiceFactory.InternalError("CLEANUP_OLD_FAILED", "清理旧通知失败", err)
	}
	return count, nil
}

// ResendNotification 重新发送通知
func (s *notificationServiceImpl) ResendNotification(ctx context.Context, id string, userID string, method string) error {
	// 检查通知是否存在
	notif, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("NOTIFICATION_GET_FAILED", "获取通知失败", err)
	}

	if notif == nil {
		return errors.BookstoreServiceFactory.NotFoundError("Notification", id)
	}

	// 验证权限
	if notif.UserID != userID {
		return errors.BookstoreServiceFactory.ForbiddenError("NO_PERMISSION", "无权操作此通知")
	}

	// 根据类型重新发送
	switch method {
	case "email":
		return s.resendEmail(ctx, notif)
	case "push":
		return s.resendPush(ctx, notif)
	default:
		return errors.BookstoreServiceFactory.ValidationError("UNSUPPORTED_METHOD", "不支持的发送方式")
	}
}

// resendEmail 通过邮件重新发送通知
func (s *notificationServiceImpl) resendEmail(ctx context.Context, notif *notification.Notification) error {
	// 检查是否有邮件服务
	if s.emailService == nil {
		return errors.BookstoreServiceFactory.InternalError("EMAIL_SERVICE_NOT_AVAILABLE", "邮件服务不可用", nil)
	}

	// 调用邮件服务发送通知
	err := s.emailService.SendEmail(ctx, notif.UserID, notif.Title, notif.Content)
	if err != nil {
		return errors.BookstoreServiceFactory.InternalError("EMAIL_SEND_FAILED", "发送邮件失败", err)
	}

	return nil
}

// resendPush 通过推送重新发送通知
func (s *notificationServiceImpl) resendPush(ctx context.Context, notif *notification.Notification) error {
	// 检查是否有WebSocket Hub
	if s.wsHub == nil {
		return errors.BookstoreServiceFactory.InternalError("WS_HUB_NOT_AVAILABLE", "WebSocket服务不可用", nil)
	}

	// 通过WebSocket Hub推送通知
	s.wsHub.BroadcastNotification(notif.UserID, notif)

	return nil
}


// ClearReadNotifications 清除已读通知
func (s *notificationServiceImpl) ClearReadNotifications(ctx context.Context, userID string) (int64, error) {
	count, err := s.notificationRepo.DeleteReadForUser(ctx, userID)
	if err != nil {
		return 0, errors.BookstoreServiceFactory.InternalError("CLEAR_READ_FAILED", "清除已读通知失败", err)
	}
	return count, nil
}

// MarkMultipleAsReadWithResult 批量标记通知为已读（返回结果）
func (s *notificationServiceImpl) MarkMultipleAsReadWithResult(ctx context.Context, ids []string, userID string) (succeeded int, failed int, err error) {
	if len(ids) == 0 {
		return 0, 0, nil
	}

	now := time.Now()
	succeeded = 0
	failed = 0

	for _, id := range ids {
		// 验证通知存在且属于该用户
		notif, err := s.notificationRepo.GetByID(ctx, id)
		if err != nil {
			failed++
			continue
		}
		if notif == nil || notif.UserID != userID {
			failed++
			continue
		}

		// 更新已读状态
		updates := map[string]interface{}{
			"read":     true,
			"read_at":  &now,
			"updated_at": &now,
		}

		if err := s.notificationRepo.Update(ctx, id, updates); err != nil {
			failed++
			continue
		}

		succeeded++
	}

	return succeeded, failed, nil
}

// BatchDeleteNotificationsWithResult 批量删除通知（返回结果）
func (s *notificationServiceImpl) BatchDeleteNotificationsWithResult(ctx context.Context, ids []string, userID string) (succeeded int, failed int, err error) {
	if len(ids) == 0 {
		return 0, 0, nil
	}

	succeeded = 0
	failed = 0

	for _, id := range ids {
		// 验证通知存在且属于该用户
		notif, err := s.notificationRepo.GetByID(ctx, id)
		if err != nil {
			failed++
			continue
		}
		if notif == nil || notif.UserID != userID {
			failed++
			continue
		}

		// 删除通知
		if err := s.notificationRepo.Delete(ctx, id); err != nil {
			failed++
			continue
		}

		succeeded++
	}

	return succeeded, failed, nil
}

