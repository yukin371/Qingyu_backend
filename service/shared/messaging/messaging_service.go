package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// TODO: 实现 BaseService 接口方法 (Initialize, Health, Close, GetServiceName, GetVersion)
// TODO: 完善消息队列功能（消息持久化、重试机制、死信队列）
// TODO: 添加消息优先级支持
// TODO: 实现消息批量发送和接收

// MessagingServiceImpl 消息服务实现
type MessagingServiceImpl struct {
	queueClient QueueClient // Redis队列客户端
}

// QueueClient 队列客户端接口（Redis）
type QueueClient interface {
	// 发布消息到Redis Stream
	Publish(ctx context.Context, stream string, data map[string]interface{}) (string, error)
	// 从Redis Stream读取消息
	Subscribe(ctx context.Context, stream, group, consumer string, count int64) ([]StreamMessage, error)
	// 确认消息
	Ack(ctx context.Context, stream, group, messageID string) error
	// 创建消费者组
	CreateGroup(ctx context.Context, stream, group string) error
	// 删除流
	DeleteStream(ctx context.Context, stream string) error
	// 列出所有流
	ListStreams(ctx context.Context) ([]string, error)
}

// StreamMessage Redis Stream消息
type StreamMessage struct {
	ID   string
	Data map[string]interface{}
}

// NewMessagingService 创建消息服务
func NewMessagingService(queueClient QueueClient) MessagingService {
	return &MessagingServiceImpl{
		queueClient: queueClient,
	}
}

// ============ 发布消息 ============

// Publish 发布消息
func (s *MessagingServiceImpl) Publish(ctx context.Context, topic string, message interface{}) error {
	// 1. 序列化消息
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	// 2. 构建消息数据
	data := map[string]interface{}{
		"topic":      topic,
		"payload":    string(payload),
		"created_at": time.Now().Unix(),
	}

	// 3. 发布到Redis Stream
	messageID, err := s.queueClient.Publish(ctx, topic, data)
	if err != nil {
		return fmt.Errorf("发布消息失败: %w", err)
	}

	_ = messageID // 可以记录消息ID用于追踪
	return nil
}

// PublishDelayed 发布延迟消息
func (s *MessagingServiceImpl) PublishDelayed(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	// 简化实现：直接发布，实际应使用延迟队列
	// TODO: 实现真正的延迟队列（可使用Redis Sorted Set）
	time.Sleep(delay) // 简化处理
	return s.Publish(ctx, topic, message)
}

// ============ 订阅消息 ============

// Subscribe 订阅消息
func (s *MessagingServiceImpl) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	// 1. 创建消费者组
	group := fmt.Sprintf("%s_group", topic)
	consumer := fmt.Sprintf("consumer_%d", time.Now().Unix())

	err := s.queueClient.CreateGroup(ctx, topic, group)
	if err != nil {
		// 组可能已存在，忽略错误
	}

	// 2. 持续读取消息
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 读取消息
			messages, err := s.queueClient.Subscribe(ctx, topic, group, consumer, 10)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			// 处理每条消息
			for _, msg := range messages {
				// 构建消息对象
				message := &Message{
					ID:    msg.ID,
					Topic: topic,
				}

				if payload, ok := msg.Data["payload"].(string); ok {
					var p interface{}
					json.Unmarshal([]byte(payload), &p)
					message.Payload = p
				}

				if ts, ok := msg.Data["created_at"].(int64); ok {
					message.CreatedAt = time.Unix(ts, 0)
				}

				// 调用处理函数
				if err := handler(ctx, message); err != nil {
					// 处理失败，可以考虑重试
					continue
				}

				// 确认消息
				s.queueClient.Ack(ctx, topic, group, msg.ID)
			}

			// 避免CPU占用过高
			if len(messages) == 0 {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// Unsubscribe 取消订阅
func (s *MessagingServiceImpl) Unsubscribe(ctx context.Context, topic string) error {
	// 简化实现：实际应该管理订阅状态
	return nil
}

// ============ 管理 ============

// CreateTopic 创建主题
func (s *MessagingServiceImpl) CreateTopic(ctx context.Context, topic string) error {
	// Redis Stream会在第一次发布消息时自动创建
	// 这里可以预先创建消费者组
	return s.queueClient.CreateGroup(ctx, topic, fmt.Sprintf("%s_group", topic))
}

// DeleteTopic 删除主题
func (s *MessagingServiceImpl) DeleteTopic(ctx context.Context, topic string) error {
	return s.queueClient.DeleteStream(ctx, topic)
}

// ListTopics 列出所有主题
func (s *MessagingServiceImpl) ListTopics(ctx context.Context) ([]string, error) {
	return s.queueClient.ListStreams(ctx)
}

// ============ 健康检查 ============

// Health 健康检查
func (s *MessagingServiceImpl) Health(ctx context.Context) error {
	// 尝试发布一条测试消息
	testMsg := map[string]interface{}{
		"type": "health_check",
		"time": time.Now().Unix(),
	}

	_, err := s.queueClient.Publish(ctx, "health_check", testMsg)
	return err
}
