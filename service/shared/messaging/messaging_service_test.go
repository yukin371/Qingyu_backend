package messaging

import (
	"context"
	"errors"
	"testing"
	"time"
)

// ============ Mock QueueClient ============

type MockQueueClient struct {
	publishFunc      func(ctx context.Context, stream string, data map[string]interface{}) (string, error)
	subscribeFunc    func(ctx context.Context, stream, group, consumer string, count int64) ([]StreamMessage, error)
	ackFunc          func(ctx context.Context, stream, group, messageID string) error
	createGroupFunc  func(ctx context.Context, stream, group string) error
	deleteStreamFunc func(ctx context.Context, stream string) error
	listStreamsFunc  func(ctx context.Context) ([]string, error)
}

func (m *MockQueueClient) Publish(ctx context.Context, stream string, data map[string]interface{}) (string, error) {
	if m.publishFunc != nil {
		return m.publishFunc(ctx, stream, data)
	}
	return "msg-1", nil
}

func (m *MockQueueClient) Subscribe(ctx context.Context, stream, group, consumer string, count int64) ([]StreamMessage, error) {
	if m.subscribeFunc != nil {
		return m.subscribeFunc(ctx, stream, group, consumer, count)
	}
	return []StreamMessage{}, nil
}

func (m *MockQueueClient) Ack(ctx context.Context, stream, group, messageID string) error {
	if m.ackFunc != nil {
		return m.ackFunc(ctx, stream, group, messageID)
	}
	return nil
}

func (m *MockQueueClient) CreateGroup(ctx context.Context, stream, group string) error {
	if m.createGroupFunc != nil {
		return m.createGroupFunc(ctx, stream, group)
	}
	return nil
}

func (m *MockQueueClient) DeleteStream(ctx context.Context, stream string) error {
	if m.deleteStreamFunc != nil {
		return m.deleteStreamFunc(ctx, stream)
	}
	return nil
}

func (m *MockQueueClient) ListStreams(ctx context.Context) ([]string, error) {
	if m.listStreamsFunc != nil {
		return m.listStreamsFunc(ctx)
	}
	return []string{}, nil
}

// ============ 测试用例 ============

func TestPublish(t *testing.T) {
	mockClient := &MockQueueClient{
		publishFunc: func(ctx context.Context, stream string, data map[string]interface{}) (string, error) {
			if stream == "" {
				return "", errors.New("stream不能为空")
			}
			return "msg-123", nil
		},
	}

	service := NewMessagingService(mockClient)

	// 测试发布成功
	err := service.Publish(context.Background(), "test_topic", map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		t.Errorf("发布消息失败: %v", err)
	}
}

func TestPublishFailed(t *testing.T) {
	mockClient := &MockQueueClient{
		publishFunc: func(ctx context.Context, stream string, data map[string]interface{}) (string, error) {
			return "", errors.New("Redis连接失败")
		},
	}

	service := NewMessagingService(mockClient)

	err := service.Publish(context.Background(), "test_topic", map[string]interface{}{
		"key": "value",
	})
	if err == nil {
		t.Error("期望发布失败，但成功了")
	}
}

func TestPublishDelayed(t *testing.T) {
	mockClient := &MockQueueClient{}
	service := NewMessagingService(mockClient)

	start := time.Now()
	delay := 100 * time.Millisecond

	err := service.PublishDelayed(context.Background(), "test_topic", map[string]interface{}{
		"key": "value",
	}, delay)

	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("延迟发布消息失败: %v", err)
	}

	if elapsed < delay {
		t.Errorf("延迟时间不足，期望 %v，实际 %v", delay, elapsed)
	}
}

func TestSubscribe(t *testing.T) {
	// 创建一个带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	mockClient := &MockQueueClient{
		subscribeFunc: func(ctx context.Context, stream, group, consumer string, count int64) ([]StreamMessage, error) {
			return []StreamMessage{
				{
					ID: "msg-1",
					Data: map[string]interface{}{
						"payload":    `{"test":"data"}`,
						"created_at": time.Now().Unix(),
					},
				},
			}, nil
		},
	}

	service := NewMessagingService(mockClient)

	handlerCalled := false
	handler := func(ctx context.Context, message *Message) error {
		handlerCalled = true
		return nil
	}

	// 启动订阅（会被context超时终止）
	err := service.Subscribe(ctx, "test_topic", handler)

	// 应该因为context超时而返回错误
	if err == nil {
		t.Error("期望订阅因超时而终止")
	}

	if !handlerCalled {
		t.Error("处理函数未被调用")
	}
}

func TestCreateTopic(t *testing.T) {
	mockClient := &MockQueueClient{
		createGroupFunc: func(ctx context.Context, stream, group string) error {
			if stream == "" {
				return errors.New("stream不能为空")
			}
			return nil
		},
	}

	service := NewMessagingService(mockClient)

	err := service.CreateTopic(context.Background(), "new_topic")
	if err != nil {
		t.Errorf("创建主题失败: %v", err)
	}
}

func TestDeleteTopic(t *testing.T) {
	mockClient := &MockQueueClient{
		deleteStreamFunc: func(ctx context.Context, stream string) error {
			if stream == "protected_topic" {
				return errors.New("不能删除受保护的主题")
			}
			return nil
		},
	}

	service := NewMessagingService(mockClient)

	// 测试删除成功
	err := service.DeleteTopic(context.Background(), "test_topic")
	if err != nil {
		t.Errorf("删除主题失败: %v", err)
	}

	// 测试删除失败
	err = service.DeleteTopic(context.Background(), "protected_topic")
	if err == nil {
		t.Error("期望删除失败，但成功了")
	}
}

func TestListTopics(t *testing.T) {
	mockClient := &MockQueueClient{
		listStreamsFunc: func(ctx context.Context) ([]string, error) {
			return []string{"topic1", "topic2", "topic3"}, nil
		},
	}

	service := NewMessagingService(mockClient)

	topics, err := service.ListTopics(context.Background())
	if err != nil {
		t.Errorf("列出主题失败: %v", err)
	}

	if len(topics) != 3 {
		t.Errorf("期望3个主题，实际 %d 个", len(topics))
	}
}

func TestHealth(t *testing.T) {
	mockClient := &MockQueueClient{
		publishFunc: func(ctx context.Context, stream string, data map[string]interface{}) (string, error) {
			return "health-check-msg", nil
		},
	}

	service := NewMessagingService(mockClient)

	// 初始化服务（使用类型断言访问实现类的Initialize方法）
	if serviceImpl, ok := service.(*MessagingServiceImpl); ok {
		err := serviceImpl.Initialize(context.Background())
		if err != nil {
			t.Fatalf("初始化服务失败: %v", err)
		}
	}

	err := service.Health(context.Background())
	if err != nil {
		t.Errorf("健康检查失败: %v", err)
	}
}

func TestHealthFailed(t *testing.T) {
	mockClient := &MockQueueClient{
		publishFunc: func(ctx context.Context, stream string, data map[string]interface{}) (string, error) {
			return "", errors.New("Redis不可用")
		},
	}

	service := NewMessagingService(mockClient)

	// 初始化服务（使用类型断言访问实现类的Initialize方法）
	if serviceImpl, ok := service.(*MessagingServiceImpl); ok {
		err := serviceImpl.Initialize(context.Background())
		if err != nil {
			t.Fatalf("初始化服务失败: %v", err)
		}
	}

	err := service.Health(context.Background())
	if err == nil {
		t.Error("期望健康检查失败，但成功了")
	}
}

func TestUnsubscribe(t *testing.T) {
	mockClient := &MockQueueClient{}
	service := NewMessagingService(mockClient)

	// 测试取消订阅
	err := service.Unsubscribe(context.Background(), "test_topic")
	if err != nil {
		t.Errorf("取消订阅失败: %v", err)
	}
}
