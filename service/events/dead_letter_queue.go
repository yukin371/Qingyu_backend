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

// DeadLetterQueue 死信队列接口
type DeadLetterQueue interface {
	// Add 添加到死信队列
	Add(ctx context.Context, item *RetryItem) error

	// List 获取死信队列项
	List(ctx context.Context, limit, offset int64) ([]*DeadLetterItem, error)

	// Get 根据ID获取项
	Get(ctx context.Context, id string) (*DeadLetterItem, error)

	// Reprocess 重新处理
	Reprocess(ctx context.Context, id string) error

	// Delete 删除
	Delete(ctx context.Context, id string) error

	// Count 统计数量
	Count(ctx context.Context) (int64, error)

	// Cleanup 清理旧项
	Cleanup(ctx context.Context, before time.Time) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// DeadLetterItem 死信队列项
type DeadLetterItem struct {
	ID             string                 `bson:"_id" json:"id"`
	EventType      string                 `bson:"event_type" json:"event_type"`
	EventData      interface{}            `bson:"event_data" json:"event_data"`
	EventSource    string                 `bson:"event_source" json:"event_source"`
	EventTimestamp time.Time              `bson:"event_timestamp" json:"event_timestamp"`
	HandlerName    string                 `bson:"handler_name" json:"handler_name"`
	Error          string                 `bson:"error" json:"error"`
	Attempt        int                    `bson:"attempt" json:"attempt"`
	Reason         string                 `bson:"reason" json:"reason"` // 失败原因
	Processed      bool                   `bson:"processed" json:"processed"`
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
	ProcessedAt    *time.Time             `bson:"processed_at,omitempty" json:"processed_at,omitempty"`
	Metadata       map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// MongoDeadLetterQueue MongoDB 死信队列实现
type MongoDeadLetterQueue struct {
	collection *mongo.Collection
}

// NewMongoDeadLetterQueue 创建 MongoDB 死信队列
func NewMongoDeadLetterQueue(db *mongo.Database) DeadLetterQueue {
	queue := &MongoDeadLetterQueue{
		collection: db.Collection("dead_letter_queue"),
	}

	// 创建索引
	queue.ensureIndexes(context.Background())

	return queue
}

// ensureIndexes 创建必要的索引
func (q *MongoDeadLetterQueue) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		// 创建时间索引
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		// 处理状态索引
		{
			Keys: bson.D{{Key: "processed", Value: 1}},
		},
		// 事件类型索引
		{
			Keys: bson.D{{Key: "event_type", Value: 1}},
		},
		// TTL 索引（保留30天）
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(30 * 24 * 3600),
		},
	}

	_, err := q.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Add 添加到死信队列
func (q *MongoDeadLetterQueue) Add(ctx context.Context, item *RetryItem) error {
	deadLetterItem := &DeadLetterItem{
		ID:             generateDeadLetterItemID(),
		EventType:      item.EventType,
		EventData:      item.EventData,
		EventSource:    item.EventSource,
		EventTimestamp: item.EventTimestamp,
		HandlerName:    item.HandlerName,
		Error:          item.Error,
		Attempt:        item.Attempt,
		Reason:         fmt.Sprintf("重试 %d 次后仍然失败", item.Attempt),
		Processed:      false,
		CreatedAt:      time.Now(),
	}

	_, err := q.collection.InsertOne(ctx, deadLetterItem)
	if err != nil {
		return fmt.Errorf("failed to add to dead letter queue: %w", err)
	}

	return nil
}

// List 获取死信队列项
func (q *MongoDeadLetterQueue) List(ctx context.Context, limit, offset int64) ([]*DeadLetterItem, error) {
	filter := bson.M{"processed": false}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := q.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list dead letter items: %w", err)
	}
	defer cursor.Close(ctx)

	var items []*DeadLetterItem
	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode dead letter items: %w", err)
	}

	return items, nil
}

// Get 根据ID获取项
func (q *MongoDeadLetterQueue) Get(ctx context.Context, id string) (*DeadLetterItem, error) {
	var item DeadLetterItem
	err := q.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get dead letter item: %w", err)
	}

	return &item, nil
}

// Reprocess 重新处理
func (q *MongoDeadLetterQueue) Reprocess(ctx context.Context, id string) error {
	// 获取项
	item, err := q.Get(ctx, id)
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("dead letter item not found: %s", id)
	}

	// 重建事件
	event := &base.BaseEvent{
		EventType: item.EventType,
		EventData: item.EventData,
		Timestamp: item.EventTimestamp,
		Source:    item.EventSource,
	}

	// 这里需要重新发布事件到事件总线
	// 由于需要在上下文中获取事件总线的引用，这里暂时只记录日志
	fmt.Printf("[DeadLetterQueue] 重新处理事件: %s, 类型: %s\n", id, item.EventType)

	// 标记为已处理
	now := time.Now()
	_, err = q.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"processed":    true,
			"processed_at": now,
		}})

	if err != nil {
		return fmt.Errorf("failed to mark dead letter item as processed: %w", err)
	}

	return nil
}

// Delete 删除
func (q *MongoDeadLetterQueue) Delete(ctx context.Context, id string) error {
	_, err := q.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete dead letter item: %w", err)
	}
	return nil
}

// Count 统计数量
func (q *MongoDeadLetterQueue) Count(ctx context.Context) (int64, error) {
	count, err := q.collection.CountDocuments(ctx, bson.M{"processed": false})
	if err != nil {
		return 0, fmt.Errorf("failed to count dead letter items: %w", err)
	}
	return count, nil
}

// Cleanup 清理旧项
func (q *MongoDeadLetterQueue) Cleanup(ctx context.Context, before time.Time) (int64, error) {
	result, err := q.collection.DeleteMany(ctx, bson.M{
		"created_at": bson.M{"$lt": before},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup dead letter items: %w", err)
	}

	return result.DeletedCount, nil
}

// Health 健康检查
func (q *MongoDeadLetterQueue) Health(ctx context.Context) error {
	return q.collection.Database().Client().Ping(ctx, nil)
}

// generateDeadLetterItemID 生成死信项ID
func generateDeadLetterItemID() string {
	return fmt.Sprintf("dlq_%d", time.Now().UnixNano())
}
