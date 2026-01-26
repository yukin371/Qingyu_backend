package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"Qingyu_backend/models/search"
)

const (
	// Redis 队列名称
	redisSyncQueue = "search:sync:events"
	// 批量缓冲区大小
	bufferSize = 100
	// 批量刷新超时
	flushTimeout = time.Second
)

// ChangeStreamListenerImpl MongoDB Change Stream 监听器实现
type ChangeStreamListenerImpl struct {
	mongoClient    *mongo.Client
	mongoDB        *mongo.Database
	redisClient    *redis.Client
	zapLogger      *zap.Logger
	logger         *log.Logger     // 保持兼容性
	collections    []string        // 监听的集合列表
	eventBuffer    []search.SyncEvent // 事件缓冲区
	bufferMutex    sync.Mutex      // 缓冲区互斥锁
	resumeTokenMap map[string][]byte // 每个集合的 resume token
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewChangeStreamListener 创建 Change Stream 监听器
func NewChangeStreamListener(
	mongoClient *mongo.Client,
	mongoDB *mongo.Database,
	redisClient *redis.Client,
	logger *zap.Logger,
) (*ChangeStreamListenerImpl, error) {
	if mongoClient == nil {
		return nil, fmt.Errorf("mongoClient is required")
	}
	if mongoDB == nil {
		return nil, fmt.Errorf("mongoDB is required")
	}
	if redisClient == nil {
		return nil, fmt.Errorf("redisClient is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ChangeStreamListenerImpl{
		mongoClient:    mongoClient,
		mongoDB:        mongoDB,
		redisClient:    redisClient,
		zapLogger:      logger,
		logger:         log.New(log.Writer(), "", log.LstdFlags),
		collections:    []string{"books", "projects", "documents", "users"},
		eventBuffer:    make([]search.SyncEvent, 0, bufferSize),
		resumeTokenMap: make(map[string][]byte),
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

// Watch 监听指定集合的变更（实现接口）
func (l *ChangeStreamListenerImpl) Watch(ctx context.Context, db, collection, index string) error {
	l.zapLogger.Info("Watching collection",
		zap.String("database", db),
		zap.String("collection", collection),
		zap.String("index", index))

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		if err := l.watchCollection(ctx, collection); err != nil {
			l.zapLogger.Error("Failed to watch collection",
				zap.String("collection", collection),
				zap.Error(err))
		}
	}()

	return nil
}

// WatchChanges 启动监听所有集合（阻塞方法，应该在 goroutine 中运行）
func (l *ChangeStreamListenerImpl) WatchChanges(ctx context.Context) error {
	l.zapLogger.Info("Starting MongoDB Change Stream listener",
		zap.Strings("collections", l.collections))

	// 为每个集合启动独立的监听器
	for _, collection := range l.collections {
		l.wg.Add(1)
		go func(colName string) {
			defer l.wg.Done()
			if err := l.watchCollection(ctx, colName); err != nil {
				l.zapLogger.Error("Failed to watch collection",
					zap.String("collection", colName),
					zap.Error(err))
			}
		}(collection)
	}

	// 启动定时刷新 goroutine
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.flushTicker(ctx)
	}()

	// 等待所有监听器完成
	l.wg.Wait()

	return nil
}

// watchCollection 监听单个集合的变更流
func (l *ChangeStreamListenerImpl) watchCollection(ctx context.Context, collection string) error {
	l.zapLogger.Info("Watching collection",
		zap.String("collection", collection),
		zap.String("database", l.mongoDB.Name()))

	coll := l.mongoDB.Collection(collection)

	// 构建变更流选项
	opts := options.ChangeStream()
	opts.SetFullDocument(options.UpdateLookup)

	// 如果有 resume token，使用它来断点续传
	if token, exists := l.resumeTokenMap[collection]; exists {
		l.zapLogger.Info("Resuming change stream",
			zap.String("collection", collection),
			zap.Int("token_size", len(token)))
		opts.SetResumeAfter(token)
	}

	// 创建变更流
	stream, err := coll.Watch(ctx, mongo.Pipeline{}, opts)
	if err != nil {
		return fmt.Errorf("failed to create change stream for %s: %w", collection, err)
	}
	defer stream.Close(ctx)

	l.zapLogger.Info("Change stream created successfully",
		zap.String("collection", collection))

	// 处理变更事件
	for {
		select {
		case <-ctx.Done():
			l.zapLogger.Info("Stopping change stream watcher",
				zap.String("collection", collection))
			return ctx.Err()

		default:
			// 检查是否有下一个事件
			if !stream.Next(ctx) {
				if err := stream.Err(); err != nil {
					l.zapLogger.Error("Change stream error",
						zap.String("collection", collection),
						zap.Error(err))
					return err
				}
				// 流已关闭
				l.zapLogger.Info("Change stream closed",
					zap.String("collection", collection))
				return nil
			}

			// 获取原始事件
			rawEvent := stream.Current

			// 转换为 SyncEvent
			syncEvent, err := l.ConvertEventToSyncEvent(collection, &rawEvent)
			if err != nil {
				l.zapLogger.Error("Failed to convert event",
					zap.String("collection", collection),
					zap.Error(err))
				continue
			}

			// 添加到缓冲区
			l.bufferMutex.Lock()
			l.eventBuffer = append(l.eventBuffer, *syncEvent)
			shouldFlush := len(l.eventBuffer) >= bufferSize
			l.bufferMutex.Unlock()

			l.zapLogger.Debug("Event added to buffer",
				zap.String("collection", collection),
				zap.String("event_type", string(syncEvent.Type)),
				zap.String("document_id", syncEvent.ID),
				zap.Int("buffer_size", len(l.eventBuffer)))

			// 如果缓冲区满了，立即刷新
			if shouldFlush {
				if err := l.FlushEvents(); err != nil {
					l.zapLogger.Error("Failed to flush events",
						zap.String("collection", collection),
						zap.Error(err))
				}
			}

			// 保存 resume token
			token := stream.ResumeToken()
			if len(token) > 0 {
				l.SaveResumeToken(collection, token)
			}
		}
	}
}

// ProcessChange 处理变更事件（实现接口）
func (l *ChangeStreamListenerImpl) ProcessChange(ctx context.Context, change *ChangeData) error {
	l.zapLogger.Info("Processing change",
		zap.String("operation", change.Operation),
		zap.String("collection", change.Collection))

	// 转换为 SyncEvent
	syncEvent := &search.SyncEvent{
		ID:        change.ID,
		Type:      search.SyncEventType(change.Operation),
		Index:     change.Collection,
		OpType:    search.SyncEventType(change.Operation),
		Timestamp: change.UpdateTime,
	}

	// 添加到缓冲区
	l.bufferMutex.Lock()
	l.eventBuffer = append(l.eventBuffer, *syncEvent)
	l.bufferMutex.Unlock()

	return nil
}

// ConvertEventToSyncEvent 将 MongoDB 变更事件转换为 SyncEvent
func (l *ChangeStreamListenerImpl) ConvertEventToSyncEvent(collection string, event *bson.Raw) (*search.SyncEvent, error) {
	var changeEvent struct {
		ID             primitive.ObjectID `bson:"_id"`
		OperationType  string             `bson:"operationType"`
		FullDocument   bson.M             `bson:"fullDocument,omitempty"`
		DocumentKey    bson.M             `bson:"documentKey"`
		NS             struct {
			DB         string `bson:"db"`
			Collection string `bson:"collection"`
		} `bson:"ns"`
		UpdateDescription *struct {
			UpdatedFields bson.M `bson:"updatedFields"`
			RemovedFields []string `bson:"removedFields"`
		} `bson:"updateDescription,omitempty"`
		ClusterTime primitive.Timestamp `bson:"clusterTime"`
	}

	if err := bson.Unmarshal(*event, &changeEvent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal change event: %w", err)
	}

	// 确定事件类型
	var eventType search.SyncEventType
	switch changeEvent.OperationType {
	case "insert":
		eventType = search.SyncEventInsert
	case "update", "replace":
		eventType = search.SyncEventUpdate
	case "delete":
		eventType = search.SyncEventDelete
	default:
		return nil, fmt.Errorf("unknown operation type: %s", changeEvent.OperationType)
	}

	// 提取文档 ID
	docID, ok := changeEvent.DocumentKey["_id"].(primitive.ObjectID)
	if !ok {
		// 尝试字符串类型
		if docIDStr, ok := changeEvent.DocumentKey["_id"].(string); ok {
			return &search.SyncEvent{
				ID:        docIDStr,
				Type:      eventType,
				Index:     collection,
				OpType:    eventType,
				Timestamp: time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("invalid document key type")
	}

	syncEvent := &search.SyncEvent{
		ID:        docID.Hex(),
		Type:      eventType,
		Index:     collection,
		OpType:    eventType,
		Timestamp: time.Now(),
	}

	// 根据事件类型优化数据
	switch eventType {
	case search.SyncEventInsert:
		// insert: 只携带关键字段
		syncEvent.FullDocument = l.extractEssentialFields(collection, changeEvent.FullDocument)
	case search.SyncEventUpdate:
		// update: 携带变更的字段
		if changeEvent.UpdateDescription != nil && changeEvent.UpdateDescription.UpdatedFields != nil {
			syncEvent.FullDocument = changeEvent.UpdateDescription.UpdatedFields
			// 提取变更字段列表
			fields := make([]string, 0, len(changeEvent.UpdateDescription.UpdatedFields))
			for key := range changeEvent.UpdateDescription.UpdatedFields {
				fields = append(fields, key)
			}
			syncEvent.ChangedFields = fields
		}
	case search.SyncEventDelete:
		// delete: 只携带 ID
		syncEvent.FullDocument = bson.M{"_id": docID.Hex()}
	}

	return syncEvent, nil
}

// extractEssentialFields 从文档中提取关键字段
func (l *ChangeStreamListenerImpl) extractEssentialFields(collection string, doc bson.M) bson.M {
	essential := bson.M{"_id": doc["_id"]}

	switch collection {
	case "books":
		// 书籍关键字段
		if title, ok := doc["title"]; ok {
			essential["title"] = title
		}
		if authorID, ok := doc["author_id"]; ok {
			essential["author_id"] = authorID
		}
		if createdAt, ok := doc["created_at"]; ok {
			essential["created_at"] = createdAt
		}
		if status, ok := doc["status"]; ok {
			essential["status"] = status
		}

	case "projects":
		// 项目关键字段
		if title, ok := doc["title"]; ok {
			essential["title"] = title
		}
		if userID, ok := doc["user_id"]; ok {
			essential["user_id"] = userID
		}
		if createdAt, ok := doc["created_at"]; ok {
			essential["created_at"] = createdAt
		}
		if status, ok := doc["status"]; ok {
			essential["status"] = status
		}

	case "documents":
		// 文档关键字段
		if title, ok := doc["title"]; ok {
			essential["title"] = title
		}
		if projectID, ok := doc["project_id"]; ok {
			essential["project_id"] = projectID
		}
		if userID, ok := doc["user_id"]; ok {
			essential["user_id"] = userID
		}
		if createdAt, ok := doc["created_at"]; ok {
			essential["created_at"] = createdAt
		}

	case "users":
		// 用户关键字段
		if username, ok := doc["username"]; ok {
			essential["username"] = username
		}
		if nickname, ok := doc["nickname"]; ok {
			essential["nickname"] = nickname
		}
		if createdAt, ok := doc["created_at"]; ok {
			essential["created_at"] = createdAt
		}
	}

	return essential
}

// FlushEvents 刷新事件缓冲区到 Redis
func (l *ChangeStreamListenerImpl) FlushEvents() error {
	l.bufferMutex.Lock()
	defer l.bufferMutex.Unlock()

	if len(l.eventBuffer) == 0 {
		return nil
	}

	l.zapLogger.Debug("Flushing events to Redis",
		zap.Int("event_count", len(l.eventBuffer)))

	ctx := context.Background()
	pipe := l.redisClient.Pipeline()

	// 批量添加到 Redis 队列
	for _, event := range l.eventBuffer {
		eventJSON, err := json.Marshal(event)
		if err != nil {
			l.zapLogger.Error("Failed to marshal event",
				zap.String("event_id", event.ID),
				zap.Error(err))
			continue
		}

		pipe.RPush(ctx, redisSyncQueue, eventJSON)
	}

	// 执行批量操作
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute redis pipeline: %w", err)
	}

	l.zapLogger.Info("Events flushed to Redis successfully",
		zap.Int("event_count", len(l.eventBuffer)))

	// 清空缓冲区
	l.eventBuffer = make([]search.SyncEvent, 0, bufferSize)

	return nil
}

// flushTicker 定时刷新缓冲区
func (l *ChangeStreamListenerImpl) flushTicker(ctx context.Context) {
	ticker := time.NewTicker(flushTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.zapLogger.Info("Flush ticker stopped")
			// 最后刷新一次
			if err := l.FlushEvents(); err != nil {
				l.zapLogger.Error("Failed to flush events on shutdown", zap.Error(err))
			}
			return
		case <-ticker.C:
			if err := l.FlushEvents(); err != nil {
				l.zapLogger.Error("Failed to flush events", zap.Error(err))
			}
		}
	}
}

// SaveResumeToken 保存 resume token
func (l *ChangeStreamListenerImpl) SaveResumeToken(collection string, token []byte) {
	l.resumeTokenMap[collection] = token
	l.zapLogger.Debug("Resume token saved",
		zap.String("collection", collection),
		zap.Int("token_size", len(token)))
}

// Stop 停止监听器
func (l *ChangeStreamListenerImpl) Stop() error {
	l.zapLogger.Info("Stopping Change Stream listener")

	// 取消上下文
	l.cancel()

	// 等待所有 goroutine 完成
	l.wg.Wait()

	// 最后刷新一次缓冲区
	if err := l.FlushEvents(); err != nil {
		l.zapLogger.Error("Failed to flush events on stop", zap.Error(err))
	}

	l.zapLogger.Info("Change Stream listener stopped")
	return nil
}

// GetBufferLength 获取当前缓冲区长度
func (l *ChangeStreamListenerImpl) GetBufferLength() int {
	l.bufferMutex.Lock()
	defer l.bufferMutex.Unlock()
	return len(l.eventBuffer)
}

// GetResumeTokens 获取所有 resume tokens
func (l *ChangeStreamListenerImpl) GetResumeTokens() map[string][]byte {
	return l.resumeTokenMap
}

// SetResumeTokens 设置 resume tokens（用于恢复监听）
func (l *ChangeStreamListenerImpl) SetResumeTokens(tokens map[string][]byte) {
	l.resumeTokenMap = tokens
}
