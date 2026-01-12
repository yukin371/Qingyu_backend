package events

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/service/base"
)

// RetryQueue 重试队列接口
type RetryQueue interface {
	// Add 添加失败事件到重试队列
	Add(ctx context.Context, event base.Event, handlerName string, err error, attempt int) error

	// Get 获取需要重试的事件
	Get(ctx context.Context, limit int64) ([]*RetryItem, error)

	// MarkSuccess 标记重试成功
	MarkSuccess(ctx context.Context, itemID string) error

	// MarkFailed 标记重试失败（最终失败，移入死信队列）
	MarkFailed(ctx context.Context, itemID string) error

	// UpdateAttempt 更新重试次数和下次重试时间
	UpdateAttempt(ctx context.Context, itemID string, nextRetryTime time.Time, attempt int) error

	// Count 统计队列中的事件数量
	Count(ctx context.Context) (int64, error)

	// Cleanup 清理过期的重试项
	Cleanup(ctx context.Context, before time.Time) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// RetryItem 重试项
type RetryItem struct {
	ID             string      `bson:"_id" json:"id"`
	EventType      string      `bson:"event_type" json:"event_type"`
	EventData      interface{} `bson:"event_data" json:"event_data"`
	EventSource    string      `bson:"event_source" json:"event_source"`
	EventTimestamp time.Time   `bson:"event_timestamp" json:"event_timestamp"`
	HandlerName    string      `bson:"handler_name" json:"handler_name"`
	Error          string      `bson:"error" json:"error"`
	Attempt        int         `bson:"attempt" json:"attempt"`
	MaxRetries     int         `bson:"max_retries" json:"max_retries"`
	NextRetry      time.Time   `bson:"next_retry" json:"next_retry"`
	CreatedAt      time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time   `bson:"updated_at" json:"updated_at"`
}

// MongoRetryQueue MongoDB 重试队列实现
type MongoRetryQueue struct {
	collection      *mongo.Collection
	deadLetterQueue DeadLetterQueue // 死信队列引用
}

// NewMongoRetryQueue 创建 MongoDB 重试队列
func NewMongoRetryQueue(db *mongo.Database, dlq DeadLetterQueue) RetryQueue {
	queue := &MongoRetryQueue{
		collection:      db.Collection("retry_queue"),
		deadLetterQueue: dlq,
	}

	// 创建索引
	queue.ensureIndexes(context.Background())

	return queue
}

// ensureIndexes 创建必要的索引
func (q *MongoRetryQueue) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		// 下次重试时间索引（用于查询待重试项）
		{
			Keys: bson.D{{Key: "next_retry", Value: 1}},
		},
		// TTL 索引（自动删除超过7天的项）
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(7 * 24 * 3600),
		},
		// 处理器名称索引
		{
			Keys: bson.D{{Key: "handler_name", Value: 1}},
		},
		// 复合索引
		{
			Keys: bson.D{
				{Key: "next_retry", Value: 1},
				{Key: "attempt", Value: 1},
			},
		},
	}

	_, err := q.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Add 添加失败事件到重试队列
func (q *MongoRetryQueue) Add(ctx context.Context, event base.Event, handlerName string, err error, attempt int) error {
	item := &RetryItem{
		ID:             generateRetryItemID(),
		EventType:      event.GetEventType(),
		EventData:      event.GetEventData(),
		EventSource:    event.GetSource(),
		EventTimestamp: event.GetTimestamp(),
		HandlerName:    handlerName,
		Error:          err.Error(),
		Attempt:        attempt,
		MaxRetries:     3,                                                      // 默认最大重试3次
		NextRetry:      time.Now().Add(time.Second * time.Duration(attempt+1)), // 简单的延迟策略
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, mongoErr := q.collection.InsertOne(ctx, item)
	if mongoErr != nil {
		return fmt.Errorf("failed to add to retry queue: %w", mongoErr)
	}

	return nil
}

// Get 获取需要重试的事件
func (q *MongoRetryQueue) Get(ctx context.Context, limit int64) ([]*RetryItem, error) {
	filter := bson.M{
		"next_retry": bson.M{"$lte": time.Now()},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "next_retry", Value: 1}}).
		SetLimit(limit)

	cursor, err := q.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get retry items: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*RetryItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode retry items: %w", err)
	}

	return items, nil
}

// MarkSuccess 标记重试成功
func (q *MongoRetryQueue) MarkSuccess(ctx context.Context, itemID string) error {
	_, err := q.collection.DeleteOne(ctx, bson.M{"_id": itemID})
	if err != nil {
		return fmt.Errorf("failed to mark retry success: %w", err)
	}
	return nil
}

// MarkFailed 标记重试失败（移入死信队列）
func (q *MongoRetryQueue) MarkFailed(ctx context.Context, itemID string) error {
	// 先获取该项
	var item RetryItem
	err := q.collection.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
	if err != nil {
		return fmt.Errorf("failed to find retry item: %w", err)
	}

	// 删除该项
	_, err = q.collection.DeleteOne(ctx, bson.M{"_id": itemID})
	if err != nil {
		return fmt.Errorf("failed to delete retry item: %w", err)
	}

	// 添加到死信队列
	if q.deadLetterQueue != nil {
		if err := q.deadLetterQueue.Add(ctx, &item); err != nil {
			return fmt.Errorf("failed to add to dead letter queue: %w", err)
		}
	}

	return nil
}

// UpdateAttempt 更新重试次数和下次重试时间
func (q *MongoRetryQueue) UpdateAttempt(ctx context.Context, itemID string, nextRetryTime time.Time, attempt int) error {
	update := bson.M{
		"$set": bson.M{
			"attempt":    attempt,
			"next_retry": nextRetryTime,
			"updated_at": time.Now(),
		},
	}

	_, err := q.collection.UpdateOne(ctx, bson.M{"_id": itemID}, update)
	if err != nil {
		return fmt.Errorf("failed to update retry attempt: %w", err)
	}

	return nil
}

// Count 统计队列中的事件数量
func (q *MongoRetryQueue) Count(ctx context.Context) (int64, error) {
	count, err := q.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count retry items: %w", err)
	}
	return count, nil
}

// Cleanup 清理过期的重试项
func (q *MongoRetryQueue) Cleanup(ctx context.Context, before time.Time) (int64, error) {
	result, err := q.collection.DeleteMany(ctx, bson.M{
		"created_at": bson.M{"$lt": before},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup retry items: %w", err)
	}

	return result.DeletedCount, nil
}

// Health 健康检查
func (q *MongoRetryQueue) Health(ctx context.Context) error {
	return q.collection.Database().Client().Ping(ctx, nil)
}

// generateRetryItemID 生成重试项ID
func generateRetryItemID() string {
	return fmt.Sprintf("retry_%d", time.Now().UnixNano())
}
