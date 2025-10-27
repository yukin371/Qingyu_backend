package messaging

import (
	"context"
	"fmt"
	"time"

	messagingModel "Qingyu_backend/models/shared/messaging"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// NotificationService 完整的通知服务接口
type NotificationService interface {
	// 创建通知
	CreateNotification(ctx context.Context, notification *messagingModel.Notification) error
	// 获取通知
	GetNotification(ctx context.Context, notificationID string) (*messagingModel.Notification, error)
	// 列出用户通知
	ListNotifications(ctx context.Context, userID string, page, pageSize int) ([]*messagingModel.Notification, int64, error)
	// 标记为已读
	MarkAsRead(ctx context.Context, notificationID string) error
	// 标记所有为已读
	MarkAllAsRead(ctx context.Context, userID string) error
	// 删除通知
	DeleteNotification(ctx context.Context, notificationID string) error
	// 获取未读数量
	GetUnreadCount(ctx context.Context, userID string) (int64, error)

	// 发送邮件通知
	SendEmailNotification(ctx context.Context, to, subject, content string) error
	// 使用模板发送邮件
	SendTemplateEmail(ctx context.Context, to, templateName string, variables map[string]string) error

	// 批量创建通知
	CreateBatchNotifications(ctx context.Context, notifications []*messagingModel.Notification) error
}

// NotificationServiceComplete 完整的通知服务实现
type NotificationServiceComplete struct {
	messageRepo  sharedRepo.MessageRepository
	emailService EmailService
	msgService   MessagingService
}

// NewNotificationServiceComplete 创建完整的通知服务
func NewNotificationServiceComplete(
	messageRepo sharedRepo.MessageRepository,
	emailService EmailService,
	msgService MessagingService,
) NotificationService {
	return &NotificationServiceComplete{
		messageRepo:  messageRepo,
		emailService: emailService,
		msgService:   msgService,
	}
}

// CreateNotification 创建通知
func (s *NotificationServiceComplete) CreateNotification(ctx context.Context, notification *messagingModel.Notification) error {
	// 1. 参数验证
	if notification.UserID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if notification.Title == "" {
		return fmt.Errorf("通知标题不能为空")
	}

	// 2. 设置默认值
	notification.Status = messagingModel.NotificationStatusSent
	notification.IsRead = false
	notification.CreatedAt = time.Now()

	// 3. 保存到数据库
	return s.messageRepo.CreateNotification(ctx, notification)
}

// GetNotification 获取通知
func (s *NotificationServiceComplete) GetNotification(ctx context.Context, notificationID string) (*messagingModel.Notification, error) {
	if notificationID == "" {
		return nil, fmt.Errorf("通知ID不能为空")
	}
	return s.messageRepo.GetNotification(ctx, notificationID)
}

// ListNotifications 列出用户通知
func (s *NotificationServiceComplete) ListNotifications(ctx context.Context, userID string, page, pageSize int) ([]*messagingModel.Notification, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	// 构建过滤器
	filter := &sharedRepo.NotificationFilter{
		UserID: userID,
		Limit:  int64(pageSize),
		Offset: int64((page - 1) * pageSize),
	}

	return s.messageRepo.ListNotifications(ctx, filter)
}

// MarkAsRead 标记为已读
func (s *NotificationServiceComplete) MarkAsRead(ctx context.Context, notificationID string) error {
	if notificationID == "" {
		return fmt.Errorf("通知ID不能为空")
	}

	// 使用UpdateNotification更新已读状态
	updates := map[string]interface{}{
		"is_read":    true,
		"read_at":    time.Now(),
		"updated_at": time.Now(),
	}

	return s.messageRepo.UpdateNotification(ctx, notificationID, updates)
}

// MarkAllAsRead 标记所有为已读
func (s *NotificationServiceComplete) MarkAllAsRead(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	// TODO(Phase3): 实现批量更新接口
	// 当前使用逐个更新的方式（简化实现）
	// 生产环境应使用UpdateMany批量更新

	// 获取所有未读通知
	isRead := false
	filter := &sharedRepo.NotificationFilter{
		UserID: userID,
		IsRead: &isRead,
	}

	notifications, _, err := s.messageRepo.ListNotifications(ctx, filter)
	if err != nil {
		return err
	}

	// 逐个标记为已读
	for _, notification := range notifications {
		updates := map[string]interface{}{
			"is_read":    true,
			"read_at":    time.Now(),
			"updated_at": time.Now(),
		}
		if err := s.messageRepo.UpdateNotification(ctx, notification.ID, updates); err != nil {
			// 记录错误但继续
			continue
		}
	}

	return nil
}

// DeleteNotification 删除通知
func (s *NotificationServiceComplete) DeleteNotification(ctx context.Context, notificationID string) error {
	if notificationID == "" {
		return fmt.Errorf("通知ID不能为空")
	}

	// 软删除：更新状态
	updates := map[string]interface{}{
		"status":     messagingModel.NotificationStatusDeleted,
		"updated_at": time.Now(),
	}

	return s.messageRepo.UpdateNotification(ctx, notificationID, updates)
}

// GetUnreadCount 获取未读数量
func (s *NotificationServiceComplete) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}

	isRead := false
	filter := &sharedRepo.NotificationFilter{
		UserID: userID,
		IsRead: &isRead,
	}

	return s.messageRepo.CountNotifications(ctx, filter)
}

// SendEmailNotification 发送邮件通知
func (s *NotificationServiceComplete) SendEmailNotification(ctx context.Context, to, subject, content string) error {
	if s.emailService == nil {
		return fmt.Errorf("邮件服务未初始化")
	}

	return s.emailService.SendEmail(ctx, &EmailRequest{
		To:      []string{to},
		Subject: subject,
		Body:    content,
		IsHTML:  true,
	})
}

// SendTemplateEmail 使用模板发送邮件
func (s *NotificationServiceComplete) SendTemplateEmail(ctx context.Context, to, templateName string, variables map[string]string) error {
	if s.emailService == nil {
		return fmt.Errorf("邮件服务未初始化")
	}

	// 1. 获取模板
	template, err := s.messageRepo.GetTemplateByName(ctx, templateName)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}

	// 2. 使用模板发送
	return s.emailService.SendWithTemplate(ctx, []string{to}, template, variables)
}

// CreateBatchNotifications 批量创建通知
func (s *NotificationServiceComplete) CreateBatchNotifications(ctx context.Context, notifications []*messagingModel.Notification) error {
	// 批量创建通知
	for _, notification := range notifications {
		notification.Status = messagingModel.NotificationStatusSent
		notification.IsRead = false
		notification.CreatedAt = time.Now()

		if err := s.messageRepo.CreateNotification(ctx, notification); err != nil {
			// 记录错误但继续
			// TODO: 使用logger记录
			continue
		}
	}

	return nil
}

// ============ BaseService接口实现（可选）============

// Initialize 初始化服务
func (s *NotificationServiceComplete) Initialize(ctx context.Context) error {
	// 验证依赖
	if s.messageRepo == nil {
		return fmt.Errorf("messageRepo未初始化")
	}

	// 健康检查
	return s.messageRepo.Health(ctx)
}

// Health 健康检查
func (s *NotificationServiceComplete) Health(ctx context.Context) error {
	if s.messageRepo == nil {
		return fmt.Errorf("messageRepo未初始化")
	}
	return s.messageRepo.Health(ctx)
}

// Close 关闭服务
func (s *NotificationServiceComplete) Close(ctx context.Context) error {
	// 无需清理
	return nil
}

// GetServiceName 获取服务名称
func (s *NotificationServiceComplete) GetServiceName() string {
	return "NotificationService"
}

// GetVersion 获取服务版本
func (s *NotificationServiceComplete) GetVersion() string {
	return "v1.0.0"
}
