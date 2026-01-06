package messaging

import (
	msgModel "Qingyu_backend/models/messaging"
	"context"
	"fmt"
	"strings"
	"time"
)

// NotificationServiceImpl 通知服务实现
type NotificationServiceImpl struct {
	messagingService MessagingService
	templateStore    TemplateStore
}

// TemplateStore 模板存储接口
type TemplateStore interface {
	GetTemplate(ctx context.Context, name string) (*msgModel.MessageTemplate, error)
	SaveTemplate(ctx context.Context, template *msgModel.MessageTemplate) error
}

// NewNotificationService 创建通知服务
func NewNotificationService(messagingService MessagingService, templateStore TemplateStore) *NotificationServiceImpl {
	return &NotificationServiceImpl{
		messagingService: messagingService,
		templateStore:    templateStore,
	}
}

// ============ 发送通知 ============

// SendEmail 发送邮件通知
func (s *NotificationServiceImpl) SendEmail(ctx context.Context, to, subject, content string) error {
	notification := map[string]interface{}{
		"type":    "email",
		"to":      to,
		"subject": subject,
		"content": content,
		"time":    time.Now().Unix(),
	}

	return s.messagingService.Publish(ctx, TopicNotificationEmail, notification)
}

// SendEmailWithTemplate 使用模板发送邮件
func (s *NotificationServiceImpl) SendEmailWithTemplate(ctx context.Context, to, templateName string, variables map[string]string) error {
	// 1. 获取模板
	template, err := s.templateStore.GetTemplate(ctx, templateName)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}

	// 2. 渲染模板
	content := s.renderTemplate(template.Content, variables)
	subject := s.renderTemplate(template.Subject, variables)

	// 3. 发送邮件
	return s.SendEmail(ctx, to, subject, content)
}

// SendSMS 发送短信通知
func (s *NotificationServiceImpl) SendSMS(ctx context.Context, phone, content string) error {
	notification := map[string]interface{}{
		"type":    "sms",
		"phone":   phone,
		"content": content,
		"time":    time.Now().Unix(),
	}

	return s.messagingService.Publish(ctx, TopicNotificationSMS, notification)
}

// SendPush 发送推送通知
func (s *NotificationServiceImpl) SendPush(ctx context.Context, userID, title, content string) error {
	notification := map[string]interface{}{
		"type":    "push",
		"user_id": userID,
		"title":   title,
		"content": content,
		"time":    time.Now().Unix(),
	}

	return s.messagingService.Publish(ctx, TopicNotificationPush, notification)
}

// SendSystemNotification 发送系统通知
func (s *NotificationServiceImpl) SendSystemNotification(ctx context.Context, userID, title, content string) error {
	notification := &msgModel.Notification{
		UserID:  userID,
		Type:    msgModel.NotificationTypeSystem,
		Title:   title,
		Content: content,
		Status:  msgModel.NotificationStatusSent,
	}
	notification.IsRead = false
	notification.CreatedAt = time.Now()

	// 直接保存到数据库或通过消息队列
	return s.messagingService.Publish(ctx, "notification.system", notification)
}

// ============ 批量发送 ============

// SendBatchEmail 批量发送邮件
func (s *NotificationServiceImpl) SendBatchEmail(ctx context.Context, recipients []string, subject, content string) error {
	for _, to := range recipients {
		if err := s.SendEmail(ctx, to, subject, content); err != nil {
			// 记录错误但继续发送
			continue
		}
	}
	return nil
}

// SendBatchSMS 批量发送短信
func (s *NotificationServiceImpl) SendBatchSMS(ctx context.Context, phones []string, content string) error {
	for _, phone := range phones {
		if err := s.SendSMS(ctx, phone, content); err != nil {
			continue
		}
	}
	return nil
}

// ============ 模板管理 ============

// CreateTemplate 创建消息模板
func (s *NotificationServiceImpl) CreateTemplate(ctx context.Context, template *msgModel.MessageTemplate) error {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	return s.templateStore.SaveTemplate(ctx, template)
}

// ============ 辅助方法 ============

// renderTemplate 渲染模板
func (s *NotificationServiceImpl) renderTemplate(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// ============ 事件监听 ============

// ListenToEvents 监听业务事件并发送通知
func (s *NotificationServiceImpl) ListenToEvents(ctx context.Context) error {
	// 监听用户注册事件
	go s.messagingService.Subscribe(ctx, TopicUserRegistered, func(ctx context.Context, msg *Message) error {
		// 发送欢迎邮件
		if userEmail, ok := msg.Payload.(map[string]interface{})["email"].(string); ok {
			s.SendEmailWithTemplate(ctx, userEmail, "welcome_email", map[string]string{
				"username": msg.Payload.(map[string]interface{})["username"].(string),
			})
		}
		return nil
	})

	// 监听钱包充值事件
	go s.messagingService.Subscribe(ctx, TopicWalletRecharged, func(ctx context.Context, msg *Message) error {
		// 发送充值成功通知
		if userID, ok := msg.Payload.(map[string]interface{})["user_id"].(string); ok {
			amount := msg.Payload.(map[string]interface{})["amount"]
			s.SendSystemNotification(ctx, userID, "充值成功", fmt.Sprintf("您已成功充值 %.2f 元", amount))
		}
		return nil
	})

	return nil
}
