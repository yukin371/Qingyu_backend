package channels

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisQueueClient Redis队列客户端实现
type RedisQueueClient struct {
	client *redis.Client
}

// NewRedisQueueClient 创建Redis队列客户端
func NewRedisQueueClient(client *redis.Client) QueueClient {
	return &RedisQueueClient{
		client: client,
	}
}

// Publish 发布消息到Redis Stream
func (r *RedisQueueClient) Publish(ctx context.Context, stream string, data map[string]interface{}) (string, error) {
	// 使用 XAdd 发布消息到 Redis Stream
	result, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: data,
	}).Result()

	if err != nil {
		return "", fmt.Errorf("发布消息到Redis Stream失败: %w", err)
	}

	return result, nil
}

// Subscribe 从Redis Stream读取消息
func (r *RedisQueueClient) Subscribe(ctx context.Context, stream, group, consumer string, count int64) ([]StreamMessage, error) {
	// 使用 XReadGroup 读取消息
	streams, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    count,
		Block:    0, // 非阻塞
	}).Result()

	if err != nil {
		// 处理特殊错误
		if err == redis.Nil {
			// 没有新消息
			return []StreamMessage{}, nil
		}
		return nil, fmt.Errorf("从Redis Stream读取消息失败: %w", err)
	}

	// 转换消息格式
	var messages []StreamMessage
	for _, xstream := range streams {
		for _, msg := range xstream.Messages {
			messages = append(messages, StreamMessage{
				ID:   msg.ID,
				Data: msg.Values,
			})
		}
	}

	return messages, nil
}

// Ack 确认消息
func (r *RedisQueueClient) Ack(ctx context.Context, stream, group, messageID string) error {
	err := r.client.XAck(ctx, stream, group, messageID).Err()
	if err != nil {
		return fmt.Errorf("确认消息失败: %w", err)
	}
	return nil
}

// CreateGroup 创建消费者组
func (r *RedisQueueClient) CreateGroup(ctx context.Context, stream, group string) error {
	// 使用 MKSTREAM 选项自动创建stream
	err := r.client.XGroupCreateMkStream(ctx, stream, group, "$").Err()
	if err != nil {
		// 如果组已存在，忽略错误
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			return nil
		}
		return fmt.Errorf("创建消费者组失败: %w", err)
	}
	return nil
}

// DeleteStream 删除流
func (r *RedisQueueClient) DeleteStream(ctx context.Context, stream string) error {
	err := r.client.Del(ctx, stream).Err()
	if err != nil {
		return fmt.Errorf("删除Stream失败: %w", err)
	}
	return nil
}

// ListStreams 列出所有流
func (r *RedisQueueClient) ListStreams(ctx context.Context) ([]string, error) {
	// 使用 SCAN 命令查找所有 Stream 类型的键
	var cursor uint64
	var streams []string

	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("扫描Redis键失败: %w", err)
		}

		// 检查每个键的类型
		for _, key := range keys {
			keyType, err := r.client.Type(ctx, key).Result()
			if err != nil {
				continue
			}
			if keyType == "stream" {
				streams = append(streams, key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return streams, nil
}
