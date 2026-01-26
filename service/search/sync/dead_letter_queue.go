package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"Qingyu_backend/models/search"
)

// DeadLetterQueue 死信队列，处理失败的同步事件
type DeadLetterQueue struct {
	redisClient *redis.Client
	zapLogger   *zap.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	stats       struct {
		totalReceived int64
		totalSent     int64
		totalRetried  int64
	}
	config *DLQConfig
}

// DLQConfig 死信队列配置
type DLQConfig struct {
	MaxRetries      int           // 最大重试次数
	RetryInterval   time.Duration // 重试间隔
	Expiration      time.Duration // 消息过期时间
}

// DefaultDLQConfig 默认配置
func DefaultDLQConfig() *DLQConfig {
	return &DLQConfig{
		MaxRetries:    3,
		RetryInterval: 5 * time.Minute,
		Expiration:    24 * time.Hour,
	}
}

// NewDeadLetterQueue 创建死信队列
func NewDeadLetterQueue(redisClient *redis.Client, zapLogger *zap.Logger, config *DLQConfig) *DeadLetterQueue {
	if config == nil {
		config = DefaultDLQConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &DeadLetterQueue{
		redisClient: redisClient,
		zapLogger:   zapLogger,
		ctx:         ctx,
		cancel:      cancel,
		config:      config,
	}
}

// Send 发送失败事件到死信队列
func (dlq *DeadLetterQueue) Send(event *search.SyncEvent, err error) {
	atomic.AddInt64(&dlq.stats.totalReceived, 1)

	dlqEvent := &DeadLetterEvent{
		Event:        event,
		ErrorMessage: err.Error(),
		RetryCount:   0,
	}

	// 序列化事件
	eventJSON, err := json.Marshal(dlqEvent)
	if err != nil {
		if dlq.zapLogger != nil {
			dlq.zapLogger.Error("Failed to marshal DLQ event",
				zap.String("event_id", event.ID),
				zap.Error(err),
			)
		}
		return
	}

	// 发送到 Redis 死信队列
	ctx, cancel := context.WithTimeout(dlq.ctx, 5*time.Second)
	defer cancel()

	// 使用 RPUSH 添加到队列
	if err := dlq.redisClient.RPush(ctx, "search:sync:dlq", eventJSON).Err(); err != nil {
		if dlq.zapLogger != nil {
			dlq.zapLogger.Error("Failed to send event to DLQ",
				zap.String("event_id", event.ID),
				zap.Error(err),
			)
		}
		return
	}

	atomic.AddInt64(&dlq.stats.totalSent, 1)

	if dlq.zapLogger != nil {
		dlq.zapLogger.Warn("Event sent to DLQ",
			zap.String("event_id", event.ID),
			zap.String("type", string(event.Type)),
			zap.String("index", event.Index),
			zap.String("error", err.Error()),
		)
	}
}

// StartRetryWorker 启动重试工作器
func (dlq *DeadLetterQueue) StartRetryWorker(processor func(*search.SyncEvent) error) {
	go dlq.retryWorker(processor)
}

// retryWorker 重试工作器
func (dlq *DeadLetterQueue) retryWorker(processor func(*search.SyncEvent) error) {
	ticker := time.NewTicker(dlq.config.RetryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-dlq.ctx.Done():
			if dlq.zapLogger != nil {
				dlq.zapLogger.Info("DLQ retry worker stopped")
			}
			return
		case <-ticker.C:
			dlq.retryEvents(processor)
		}
	}
}

// retryEvents 重试死信队列中的事件
func (dlq *DeadLetterQueue) retryEvents(processor func(*search.SyncEvent) error) {
	ctx, cancel := context.WithTimeout(dlq.ctx, 30*time.Second)
	defer cancel()

	// 获取 DLQ 长度
	length, err := dlq.redisClient.LLen(ctx, "search:sync:dlq").Result()
	if err != nil {
		if dlq.zapLogger != nil {
			dlq.zapLogger.Error("Failed to get DLQ length", zap.Error(err))
		}
		return
	}

	if length == 0 {
		return
	}

	if dlq.zapLogger != nil {
		dlq.zapLogger.Info("Retrying DLQ events", zap.Int64("count", length))
	}

	// 批量处理事件（最多 100 个）
	batchSize := int64(100)
	if length < batchSize {
		batchSize = length
	}

	for i := int64(0); i < batchSize; i++ {
		// 从队列左侧取出事件
		eventJSON, err := dlq.redisClient.LPop(ctx, "search:sync:dlq").Result()
		if err != nil {
			if err != redis.Nil && dlq.zapLogger != nil {
				dlq.zapLogger.Error("Failed to pop event from DLQ", zap.Error(err))
			}
			continue
		}

		// 解析事件
		var dlqEvent DeadLetterEvent
		if err := json.Unmarshal([]byte(eventJSON), &dlqEvent); err != nil {
			if dlq.zapLogger != nil {
				dlq.zapLogger.Error("Failed to unmarshal DLQ event",
					zap.String("event_json", eventJSON),
					zap.Error(err),
				)
			}
			continue
		}

		// 检查重试次数
		if dlqEvent.RetryCount >= dlq.config.MaxRetries {
			if dlq.zapLogger != nil {
				dlq.zapLogger.Error("Event max retries exceeded",
					zap.String("event_id", dlqEvent.Event.ID),
					zap.Int("retry_count", dlqEvent.RetryCount),
				)
			}
			// 可以选择持久化到数据库或发送告警
			continue
		}

		// 增加重试计数
		dlqEvent.RetryCount++
		dlqEvent.LastRetryAt = time.Now().Unix()

		// 尝试重新处理
		if err := processor(dlqEvent.Event); err != nil {
			// 处理失败，重新放回队列
			dlqEventJSON, _ := json.Marshal(dlqEvent)
			if err := dlq.redisClient.RPush(ctx, "search:sync:dlq", dlqEventJSON).Err(); err != nil {
				if dlq.zapLogger != nil {
					dlq.zapLogger.Error("Failed to re-queue event",
						zap.String("event_id", dlqEvent.Event.ID),
						zap.Error(err),
					)
				}
			}

			if dlq.zapLogger != nil {
				dlq.zapLogger.Warn("Event retry failed",
					zap.String("event_id", dlqEvent.Event.ID),
					zap.Int("retry_count", dlqEvent.RetryCount),
					zap.Error(err),
				)
			}
		} else {
			atomic.AddInt64(&dlq.stats.totalRetried, 1)
			if dlq.zapLogger != nil {
				dlq.zapLogger.Info("Event retry succeeded",
					zap.String("event_id", dlqEvent.Event.ID),
					zap.Int("retry_count", dlqEvent.RetryCount),
				)
			}
		}
	}
}

// GetStats 获取统计信息
func (dlq *DeadLetterQueue) GetStats(ctx context.Context) (map[string]interface{}, error) {
	queueLen, err := dlq.redisClient.LLen(ctx, "search:sync:dlq").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get DLQ length: %w", err)
	}

	return map[string]interface{}{
		"queue_length":    queueLen,
		"total_received":  atomic.LoadInt64(&dlq.stats.totalReceived),
		"total_sent":      atomic.LoadInt64(&dlq.stats.totalSent),
		"total_retried":   atomic.LoadInt64(&dlq.stats.totalRetried),
		"max_retries":     dlq.config.MaxRetries,
		"retry_interval":  dlq.config.RetryInterval.String(),
	}, nil
}

// Clear 清空死信队列
func (dlq *DeadLetterQueue) Clear(ctx context.Context) error {
	if err := dlq.redisClient.Del(ctx, "search:sync:dlq").Err(); err != nil {
		return fmt.Errorf("failed to clear DLQ: %w", err)
	}

	if dlq.zapLogger != nil {
		dlq.zapLogger.Info("DLQ cleared")
	}

	return nil
}

// Stop 停止死信队列
func (dlq *DeadLetterQueue) Stop() {
	dlq.cancel()
	if dlq.zapLogger != nil {
		dlq.zapLogger.Info("DLQ stopped")
	}
}

// SendToDLQ 发送事件到死信队列（兼容旧接口）
func (dlq *DeadLetterQueue) SendToDLQ(event *search.SyncEvent, errMsg string) {
	dlq.Send(event, fmt.Errorf("%s", errMsg))
}
