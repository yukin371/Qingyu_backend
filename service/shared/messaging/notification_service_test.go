package messaging

import (
	"context"
	"errors"
	"testing"
	"time"

	msgModel "Qingyu_backend/models/shared/messaging"
)

// ============ Mock Services ============

type MockMessagingService struct {
	publishFunc   func(ctx context.Context, topic string, message interface{}) error
	subscribeFunc func(ctx context.Context, topic string, handler MessageHandler) error
}

func (m *MockMessagingService) Publish(ctx context.Context, topic string, message interface{}) error {
	if m.publishFunc != nil {
		return m.publishFunc(ctx, topic, message)
	}
	return nil
}

func (m *MockMessagingService) PublishDelayed(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	return nil
}

func (m *MockMessagingService) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	if m.subscribeFunc != nil {
		return m.subscribeFunc(ctx, topic, handler)
	}
	return nil
}

func (m *MockMessagingService) Unsubscribe(ctx context.Context, topic string) error {
	return nil
}

func (m *MockMessagingService) CreateTopic(ctx context.Context, topic string) error {
	return nil
}

func (m *MockMessagingService) DeleteTopic(ctx context.Context, topic string) error {
	return nil
}

func (m *MockMessagingService) ListTopics(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

func (m *MockMessagingService) Health(ctx context.Context) error {
	return nil
}

type MockTemplateStore struct {
	getTemplateFunc  func(ctx context.Context, name string) (*msgModel.MessageTemplate, error)
	saveTemplateFunc func(ctx context.Context, template *msgModel.MessageTemplate) error
}

func (m *MockTemplateStore) GetTemplate(ctx context.Context, name string) (*msgModel.MessageTemplate, error) {
	if m.getTemplateFunc != nil {
		return m.getTemplateFunc(ctx, name)
	}
	return nil, errors.New("模板不存在")
}

func (m *MockTemplateStore) SaveTemplate(ctx context.Context, template *msgModel.MessageTemplate) error {
	if m.saveTemplateFunc != nil {
		return m.saveTemplateFunc(ctx, template)
	}
	return nil
}

// ============ 测试用例 ============

func TestSendEmail(t *testing.T) {
	mockMessaging := &MockMessagingService{
		publishFunc: func(ctx context.Context, topic string, message interface{}) error {
			if topic != TopicNotificationEmail {
				t.Errorf("期望主题 %s，实际 %s", TopicNotificationEmail, topic)
			}
			return nil
		},
	}

	service := NewNotificationService(mockMessaging, &MockTemplateStore{})

	err := service.SendEmail(context.Background(), "test@example.com", "测试邮件", "这是测试内容")
	if err != nil {
		t.Errorf("发送邮件失败: %v", err)
	}
}

func TestSendEmailFailed(t *testing.T) {
	mockMessaging := &MockMessagingService{
		publishFunc: func(ctx context.Context, topic string, message interface{}) error {
			return errors.New("发送失败")
		},
	}

	service := NewNotificationService(mockMessaging, &MockTemplateStore{})

	err := service.SendEmail(context.Background(), "test@example.com", "测试邮件", "这是测试内容")
	if err == nil {
		t.Error("期望发送失败，但成功了")
	}
}

func TestSendEmailWithTemplate(t *testing.T) {
	mockMessaging := &MockMessagingService{}
	mockTemplateStore := &MockTemplateStore{
		getTemplateFunc: func(ctx context.Context, name string) (*msgModel.MessageTemplate, error) {
			return &msgModel.MessageTemplate{
				Name:    "welcome_email",
				Type:    "email",
				Subject: "欢迎 {{username}}",
				Content: "您好 {{username}}，欢迎加入我们！",
			}, nil
		},
	}

	service := NewNotificationService(mockMessaging, mockTemplateStore)

	err := service.SendEmailWithTemplate(context.Background(), "test@example.com", "welcome_email", map[string]string{
		"username": "张三",
	})

	if err != nil {
		t.Errorf("使用模板发送邮件失败: %v", err)
	}
}

func TestSendEmailWithTemplateNotFound(t *testing.T) {
	mockMessaging := &MockMessagingService{}
	mockTemplateStore := &MockTemplateStore{
		getTemplateFunc: func(ctx context.Context, name string) (*msgModel.MessageTemplate, error) {
			return nil, errors.New("模板不存在")
		},
	}

	service := NewNotificationService(mockMessaging, mockTemplateStore)

	err := service.SendEmailWithTemplate(context.Background(), "test@example.com", "non_existent", map[string]string{})

	if err == nil {
		t.Error("期望模板不存在错误，但成功了")
	}
}

func TestSendSMS(t *testing.T) {
	mockMessaging := &MockMessagingService{
		publishFunc: func(ctx context.Context, topic string, message interface{}) error {
			if topic != TopicNotificationSMS {
				t.Errorf("期望主题 %s，实际 %s", TopicNotificationSMS, topic)
			}
			return nil
		},
	}

	service := NewNotificationService(mockMessaging, &MockTemplateStore{})

	err := service.SendSMS(context.Background(), "13800138000", "这是测试短信")
	if err != nil {
		t.Errorf("发送短信失败: %v", err)
	}
}

func TestSendPush(t *testing.T) {
	mockMessaging := &MockMessagingService{
		publishFunc: func(ctx context.Context, topic string, message interface{}) error {
			if topic != TopicNotificationPush {
				t.Errorf("期望主题 %s，实际 %s", TopicNotificationPush, topic)
			}
			return nil
		},
	}

	service := NewNotificationService(mockMessaging, &MockTemplateStore{})

	err := service.SendPush(context.Background(), "user123", "测试通知", "这是推送内容")
	if err != nil {
		t.Errorf("发送推送失败: %v", err)
	}
}

func TestSendSystemNotification(t *testing.T) {
	mockMessaging := &MockMessagingService{
		publishFunc: func(ctx context.Context, topic string, message interface{}) error {
			return nil
		},
	}

	service := NewNotificationService(mockMessaging, &MockTemplateStore{})

	err := service.SendSystemNotification(context.Background(), "user123", "系统通知", "这是系统消息")
	if err != nil {
		t.Errorf("发送系统通知失败: %v", err)
	}
}

func TestSendBatchEmail(t *testing.T) {
	sendCount := 0
	mockMessaging := &MockMessagingService{
		publishFunc: func(ctx context.Context, topic string, message interface{}) error {
			sendCount++
			return nil
		},
	}

	service := NewNotificationService(mockMessaging, &MockTemplateStore{})

	recipients := []string{"test1@example.com", "test2@example.com", "test3@example.com"}
	err := service.SendBatchEmail(context.Background(), recipients, "批量邮件", "这是批量内容")

	if err != nil {
		t.Errorf("批量发送邮件失败: %v", err)
	}

	if sendCount != len(recipients) {
		t.Errorf("期望发送 %d 封邮件，实际发送 %d 封", len(recipients), sendCount)
	}
}

func TestSendBatchSMS(t *testing.T) {
	sendCount := 0
	mockMessaging := &MockMessagingService{
		publishFunc: func(ctx context.Context, topic string, message interface{}) error {
			sendCount++
			return nil
		},
	}

	service := NewNotificationService(mockMessaging, &MockTemplateStore{})

	phones := []string{"13800138001", "13800138002", "13800138003"}
	err := service.SendBatchSMS(context.Background(), phones, "这是批量短信")

	if err != nil {
		t.Errorf("批量发送短信失败: %v", err)
	}

	if sendCount != len(phones) {
		t.Errorf("期望发送 %d 条短信，实际发送 %d 条", len(phones), sendCount)
	}
}

func TestCreateTemplate(t *testing.T) {
	mockTemplateStore := &MockTemplateStore{
		saveTemplateFunc: func(ctx context.Context, template *msgModel.MessageTemplate) error {
			if template.Name == "" {
				return errors.New("模板名称不能为空")
			}
			return nil
		},
	}

	service := NewNotificationService(&MockMessagingService{}, mockTemplateStore)

	template := &msgModel.MessageTemplate{
		Name:      "test_template",
		Type:      "email",
		Subject:   "测试模板",
		Content:   "这是测试内容",
		Variables: []string{"username"},
		IsActive:  true,
	}

	err := service.CreateTemplate(context.Background(), template)
	if err != nil {
		t.Errorf("创建模板失败: %v", err)
	}

	if template.CreatedAt.IsZero() {
		t.Error("CreatedAt未设置")
	}
}

func TestRenderTemplate(t *testing.T) {
	service := NewNotificationService(&MockMessagingService{}, &MockTemplateStore{})

	template := "你好 {{username}}，你的订单 {{order_id}} 已创建"
	variables := map[string]string{
		"username": "张三",
		"order_id": "12345",
	}

	result := service.renderTemplate(template, variables)
	expected := "你好 张三，你的订单 12345 已创建"

	if result != expected {
		t.Errorf("模板渲染错误，期望: %s，实际: %s", expected, result)
	}
}
