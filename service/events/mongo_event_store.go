package events

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"Qingyu_backend/service/base"
)

// MongoEventStore MongoDB 事件存储实现
type MongoEventStore struct {
	db         *mongo.Database
	collection *mongo.Collection
	config     *EventStoreConfig
	metrics    *Metrics
	logger     *EventLogger
}

// NewMongoEventStore 创建 MongoDB 事件存储
func NewMongoEventStore(db *mongo.Database, config *EventStoreConfig) EventStore {
	if config == nil {
		config = DefaultEventStoreConfig()
	}

	logger, _ := NewEventLogger(DefaultLoggingConfig())

	store := &MongoEventStore{
		db:         db,
		collection: db.Collection("events_log"),
		config:     config,
		metrics:    GetGlobalMetrics(),
		logger:     logger,
	}

	// 创建索引
	store.ensureIndexes(context.Background())

	return store
}

// ensureIndexes 创建必要的索引
func (s *MongoEventStore) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		// 事件类型索引
		{
			Keys: bson.D{
				{Key: "event_type", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
		// 来源索引
		{
			Keys: bson.D{
				{Key: "source", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
		// TTL 索引（自动删除过期事件）
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
		// 复合索引用于查询
		{
			Keys: bson.D{
				{Key: "event_type", Value: 1},
				{Key: "processed", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
	}

	_, err := s.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Store 存储事件
func (s *MongoEventStore) Store(ctx context.Context, event base.Event) error {
	if !s.config.Enabled {
		return nil
	}

	storedEvent := &StoredEvent{
		ID:        generateEventID(),
		EventType: event.GetEventType(),
		EventData: event.GetEventData(),
		Timestamp: event.GetTimestamp(),
		Source:    event.GetSource(),
		Processed: false,
		CreatedAt: time.Now(),
	}

	// 设置过期时间
	if s.config.RetentionDuration > 0 {
		storedEvent.ExpiresAt = time.Now().Add(s.config.RetentionDuration)
	}

	_, err := s.collection.InsertOne(ctx, storedEvent)
	if err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	return nil
}

// StoreBatch 批量存储事件
func (s *MongoEventStore) StoreBatch(ctx context.Context, events []base.Event) error {
	if !s.config.Enabled || len(events) == 0 {
		return nil
	}

	documents := make([]interface{}, 0, len(events))
	now := time.Now()

	for _, event := range events {
		storedEvent := &StoredEvent{
			ID:        generateEventID(),
			EventType: event.GetEventType(),
			EventData: event.GetEventData(),
			Timestamp: event.GetTimestamp(),
			Source:    event.GetSource(),
			Processed: false,
			CreatedAt: now,
		}

		// 设置过期时间
		if s.config.RetentionDuration > 0 {
			storedEvent.ExpiresAt = now.Add(s.config.RetentionDuration)
		}

		documents = append(documents, storedEvent)
	}

	// 分批插入
	batchSize := s.config.BatchSize
	if batchSize <= 0 {
		batchSize = 100
	}

	for i := 0; i < len(documents); i += batchSize {
		end := i + batchSize
		if end > len(documents) {
			end = len(documents)
		}

		batch := documents[i:end]
		_, err := s.collection.InsertMany(ctx, batch)
		if err != nil {
			return fmt.Errorf("failed to store event batch: %w", err)
		}
	}

	return nil
}

// GetByID 根据ID获取事件
func (s *MongoEventStore) GetByID(ctx context.Context, eventID string) (*StoredEvent, error) {
	var event StoredEvent
	err := s.collection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event by ID: %w", err)
	}

	return &event, nil
}

// GetByType 根据类型获取事件
func (s *MongoEventStore) GetByType(ctx context.Context, eventType string, limit, offset int64) ([]*StoredEvent, error) {
	filter := bson.M{"event_type": eventType}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}
	defer cursor.Close(ctx)

	var events []*StoredEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}

// GetBySource 根据来源获取事件
func (s *MongoEventStore) GetBySource(ctx context.Context, source string, limit, offset int64) ([]*StoredEvent, error) {
	filter := bson.M{"source": source}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by source: %w", err)
	}
	defer cursor.Close(ctx)

	var events []*StoredEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}

// GetByTimeRange 根据时间范围获取事件
func (s *MongoEventStore) GetByTimeRange(ctx context.Context, start, end time.Time, limit, offset int64) ([]*StoredEvent, error) {
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by time range: %w", err)
	}
	defer cursor.Close(ctx)

	var events []*StoredEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}

// GetByTypeAndTimeRange 根据类型和时间范围获取事件
func (s *MongoEventStore) GetByTypeAndTimeRange(ctx context.Context, eventType string, start, end time.Time, limit, offset int64) ([]*StoredEvent, error) {
	filter := bson.M{
		"event_type": eventType,
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by type and time range: %w", err)
	}
	defer cursor.Close(ctx)

	var events []*StoredEvent
	if err = cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}

// Replay 事件回放
func (s *MongoEventStore) Replay(ctx context.Context, handler base.EventHandler, filter EventFilter) (*ReplayResult, error) {
	startTime := time.Now()

	// 记录回放开始
	eventType := filter.EventType
	if eventType == "" {
		eventType = "all"
	}
	s.logger.LogReplayStarted(ctx, eventType, filter.StartTime, filter.EndTime)

	// 构建查询条件
	mongoFilter := bson.M{}
	if filter.EventType != "" {
		mongoFilter["event_type"] = filter.EventType
	}
	if filter.Source != "" {
		mongoFilter["source"] = filter.Source
	}
	if filter.StartTime != nil || filter.EndTime != nil {
		timeFilter := bson.M{}
		if filter.StartTime != nil {
			timeFilter["$gte"] = filter.StartTime
		}
		if filter.EndTime != nil {
			timeFilter["$lte"] = filter.EndTime
		}
		mongoFilter["timestamp"] = timeFilter
	}
	if filter.Processed != nil {
		mongoFilter["processed"] = *filter.Processed
	}

	// 设置查询选项
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	if filter.Limit > 0 {
		opts.SetLimit(filter.Limit)
	}
	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	cursor, err := s.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		duration := time.Since(startTime)
		s.logger.LogReplayFailed(ctx, eventType, err, duration)
		s.metrics.RecordEventReplayFailed(eventType, "query_error")
		return nil, fmt.Errorf("failed to replay events: %w", err)
	}
	defer cursor.Close(ctx)

	result := &ReplayResult{
		ReplayedCount: 0,
		FailedCount:   0,
		SkippedCount:  0,
	}

	// 逐个处理事件
	for cursor.Next(ctx) {
		var storedEvent StoredEvent
		if err := cursor.Decode(&storedEvent); err != nil {
			result.FailedCount++
			continue
		}

		// DryRun模式：只统计不执行
		if filter.DryRun {
			result.SkippedCount++
			continue
		}

		// 重建事件对象
		event := &base.BaseEvent{
			EventType: storedEvent.EventType,
			EventData: storedEvent.EventData,
			Timestamp: storedEvent.Timestamp,
			Source:    storedEvent.Source,
		}

		// 调用处理器处理事件
		if err := handler.Handle(ctx, event); err != nil {
			result.FailedCount++
			// 继续处理下一个事件，不中断整个回放过程
			continue
		}

		result.ReplayedCount++

		// 标记事件为已处理
		_, err := s.collection.UpdateOne(ctx,
			bson.M{"_id": storedEvent.ID},
			bson.M{"$set": bson.M{"processed": true}})
		if err != nil {
			// 使用EventLogger记录标记失败错误
			s.logger.logger.Error("Failed to mark event as processed",
				zap.String("event_id", storedEvent.ID),
				zap.Error(err))
		}
	}

	result.Duration = time.Since(startTime)

	// 记录回放完成
	s.logger.LogReplayCompleted(ctx, eventType, result)

	// 记录指标
	s.metrics.RecordEventReplayed(eventType, result.ReplayedCount, result.FailedCount, result.Duration)

	return result, nil
}

// Cleanup 清理过期事件
func (s *MongoEventStore) Cleanup(ctx context.Context, before time.Time) (int64, error) {
	result, err := s.collection.DeleteMany(ctx, bson.M{
		"created_at": bson.M{"$lt": before},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup events: %w", err)
	}

	return result.DeletedCount, nil
}

// Count 统计事件数量
func (s *MongoEventStore) Count(ctx context.Context, filter EventFilter) (int64, error) {
	mongoFilter := bson.M{}
	if filter.EventType != "" {
		mongoFilter["event_type"] = filter.EventType
	}
	if filter.Source != "" {
		mongoFilter["source"] = filter.Source
	}
	if filter.StartTime != nil || filter.EndTime != nil {
		timeFilter := bson.M{}
		if filter.StartTime != nil {
			timeFilter["$gte"] = filter.StartTime
		}
		if filter.EndTime != nil {
			timeFilter["$lte"] = filter.EndTime
		}
		mongoFilter["timestamp"] = timeFilter
	}
	if filter.Processed != nil {
		mongoFilter["processed"] = *filter.Processed
	}

	count, err := s.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}

	return count, nil
}

// Health 健康检查
func (s *MongoEventStore) Health(ctx context.Context) error {
	return s.db.Client().Ping(ctx, nil)
}

// generateEventID 生成事件ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// MongoEventReplayer MongoDB 事件回放器
type MongoEventReplayer struct {
	store EventStore
}

// NewMongoEventReplayer 创建 MongoDB 事件回放器
func NewMongoEventReplayer(store EventStore) EventReplayer {
	return &MongoEventReplayer{store: store}
}

// ReplayFromTimestamp 从指定时间戳回放事件
func (r *MongoEventReplayer) ReplayFromTimestamp(ctx context.Context, timestamp time.Time, handler base.EventHandler) error {
	filter := EventFilter{
		StartTime: &timestamp,
	}
	_, err := r.store.Replay(ctx, handler, filter)
	return err
}

// ReplayFromEventID 从指定事件ID回放事件
func (r *MongoEventReplayer) ReplayFromEventID(ctx context.Context, eventID string, handler base.EventHandler) error {
	// 首先获取事件
	event, err := r.store.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event not found: %s", eventID)
	}

	// 然后回放该事件之后的所有事件
	filter := EventFilter{
		StartTime: &event.Timestamp,
	}
	_, err = r.store.Replay(ctx, handler, filter)
	return err
}

// ReplayWithType 回放指定类型的事件
func (r *MongoEventReplayer) ReplayWithType(ctx context.Context, eventType string, handler base.EventHandler) error {
	filter := EventFilter{
		EventType: eventType,
	}
	_, err := r.store.Replay(ctx, handler, filter)
	return err
}
