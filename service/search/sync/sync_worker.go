package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"Qingyu_backend/models/search"
	searchengine "Qingyu_backend/service/search/engine"
)

// WorkerImpl 同步 Worker 实现
type WorkerImpl struct {
	mongoClient   *mongo.Client
	mongoDB       *mongo.Database
	redisClient   *redis.Client
	esEngine      searchengine.Engine
	zapLogger     *zap.Logger
	logger        *log.Logger       // 保持兼容性
	dlq           *DeadLetterQueue
	ctx           context.Context
	cancel        context.CancelFunc
	status        *SyncStatus
	eventCounts   struct {
		total   int64
		success int64
		failed  int64
	}
	config        *WorkerConfig
}

// WorkerConfig Worker 配置
type WorkerConfig struct {
	WorkerCount       int           // 工作器数量
	PollInterval      time.Duration // 轮询间隔
	RetryMaxAttempts  int           // 最大重试次数
	RetryDelay        time.Duration // 重试延迟
}

// DefaultWorkerConfig 默认配置
func DefaultWorkerConfig() *WorkerConfig {
	return &WorkerConfig{
		WorkerCount:       5,
		PollInterval:      100 * time.Millisecond,
		RetryMaxAttempts:  3,
		RetryDelay:        1 * time.Second,
	}
}

// NewSyncWorker 创建同步 Worker
func NewSyncWorker(
	mongoClient *mongo.Client,
	mongoDB *mongo.Database,
	redisClient *redis.Client,
	esEngine searchengine.Engine,
	zapLogger *zap.Logger,
	logger *log.Logger,
	config *WorkerConfig,
) (*WorkerImpl, error) {
	if config == nil {
		config = DefaultWorkerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerImpl{
		mongoClient: mongoClient,
		mongoDB:     mongoDB,
		redisClient: redisClient,
		esEngine:    esEngine,
		zapLogger:   zapLogger,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
		status: &SyncStatus{
			Running: false,
		},
		config: config,
	}, nil
}

// SetDeadLetterQueue 设置死信队列
func (w *WorkerImpl) SetDeadLetterQueue(dlq *DeadLetterQueue) {
	w.dlq = dlq
}

// Start 启动 Worker
func (w *WorkerImpl) Start(parentCtx context.Context) error {
	w.status.Running = true

	// 启动多个工作 goroutine
	for i := 0; i < w.config.WorkerCount; i++ {
		go w.processEvents(parentCtx, i)
	}

	if w.logger != nil {
		w.logger.Printf("Sync worker started with %d workers", w.config.WorkerCount)
	}
	if w.zapLogger != nil {
		w.zapLogger.Info("Sync worker started",
			zap.Int("worker_count", w.config.WorkerCount),
		)
	}

	return nil
}

// processEvents 处理同步事件
func (w *WorkerImpl) processEvents(parentCtx context.Context, workerID int) {
	ticker := time.NewTicker(w.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-parentCtx.Done():
			if w.zapLogger != nil {
				w.zapLogger.Debug("Worker context done",
					zap.Int("worker_id", workerID),
				)
			}
			return
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			// 从 Redis 队列中获取事件
			ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
			eventJSON, err := w.redisClient.LPop(ctx, "search:sync:events").Result()
			cancel()

			if err != nil {
				if err != redis.Nil {
					if w.zapLogger != nil {
						w.zapLogger.Error("Failed to pop event from Redis",
							zap.Int("worker_id", workerID),
							zap.Error(err),
						)
					}
				}
				continue
			}

			// 解析事件
			var event search.SyncEvent
			if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
				if w.zapLogger != nil {
					w.zapLogger.Error("Failed to unmarshal event",
						zap.Int("worker_id", workerID),
						zap.String("event_json", eventJSON),
						zap.Error(err),
					)
				}
				continue
			}

			// 处理事件（带重试）
			w.processEventWithRetry(parentCtx, &event, workerID)
		}
	}
}

// processEventWithRetry 带重试的事件处理
func (w *WorkerImpl) processEventWithRetry(ctx context.Context, event *search.SyncEvent, workerID int) {
	var lastErr error

	for attempt := 0; attempt <= w.config.RetryMaxAttempts; attempt++ {
		if attempt > 0 {
			// 重试前等待
			time.Sleep(w.config.RetryDelay)
			if w.zapLogger != nil {
				w.zapLogger.Info("Retrying event",
					zap.Int("worker_id", workerID),
					zap.String("event_id", event.ID),
					zap.Int("attempt", attempt),
				)
			}
		}

		err := w.ProcessEvent(ctx, event)
		if err == nil {
			atomic.AddInt64(&w.eventCounts.success, 1)
			if w.zapLogger != nil {
				w.zapLogger.Debug("Event processed successfully",
					zap.Int("worker_id", workerID),
					zap.String("event_id", event.ID),
					zap.String("type", string(event.Type)),
					zap.String("collection", event.Collection),
				)
			}
			return
		}

		lastErr = err
		if w.zapLogger != nil {
			w.zapLogger.Warn("Event processing failed",
				zap.Int("worker_id", workerID),
				zap.String("event_id", event.ID),
				zap.Int("attempt", attempt),
				zap.Error(err),
			)
		}
	}

	// 所有重试都失败
	atomic.AddInt64(&w.eventCounts.failed, 1)

	// 发送到死信队列
	if w.dlq != nil {
		w.dlq.Send(&event, lastErr)
	}
}

// ProcessEvent 处理单个同步事件
func (w *WorkerImpl) ProcessEvent(ctx context.Context, event *search.SyncEvent) error {
	atomic.AddInt64(&w.eventCounts.total, 1)

	if w.zapLogger != nil {
		w.zapLogger.Info("Processing event",
			zap.String("event_id", event.ID),
			zap.String("type", string(event.Type)),
			zap.String("collection", event.Collection),
			zap.String("document_id", event.DocumentID),
		)
	}

	// 根据事件类型处理
	switch event.Type {
	case search.SyncEventInsert:
		return w.handleInsert(ctx, event)
	case search.SyncEventUpdate:
		return w.handleUpdate(ctx, event)
	case search.SyncEventDelete:
		return w.handleDelete(ctx, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// handleInsert 处理插入事件
func (w *WorkerImpl) handleInsert(ctx context.Context, event *search.SyncEvent) error {
	collection := w.mongoDB.Collection(event.Collection)
	objectID, err := primitive.ObjectIDFromHex(event.DocumentID)
	if err != nil {
		return fmt.Errorf("invalid document ID: %w", err)
	}

	var doc bson.Raw
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if w.zapLogger != nil {
				w.zapLogger.Warn("Document not found in MongoDB",
					zap.String("document_id", event.DocumentID),
				)
			}
			return nil // 文档已被删除，忽略
		}
		return fmt.Errorf("failed to find document: %w", err)
	}

	// 转换为 Document
	docMap := make(map[string]interface{})
	if err := doc.Unmarshal(&docMap); err != nil {
		return fmt.Errorf("failed to unmarshal document: %w", err)
	}

	// 删除 _id 字段
	delete(docMap, "_id")

	// 确定索引名称
	indexName := w.getElasticsearchIndex(event.Collection)

	// 索引到 Elasticsearch
	documents := []searchengine.Document{
		{
			ID:     event.DocumentID,
			Fields: docMap,
		},
	}

	if err := w.esEngine.Index(ctx, indexName, documents); err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}

	if w.zapLogger != nil {
		w.zapLogger.Info("Document indexed",
			zap.String("document_id", event.DocumentID),
			zap.String("index", indexName),
		)
	}

	return nil
}

// handleUpdate 处理更新事件
func (w *WorkerImpl) handleUpdate(ctx context.Context, event *search.SyncEvent) error {
	collection := w.mongoDB.Collection(event.Collection)
	objectID, err := primitive.ObjectIDFromHex(event.DocumentID)
	if err != nil {
		return fmt.Errorf("invalid document ID: %w", err)
	}

	var doc bson.Raw
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 文档不存在，尝试从 ES 删除
			return w.handleDelete(ctx, event)
		}
		return fmt.Errorf("failed to find document: %w", err)
	}

	docMap := make(map[string]interface{})
	if err := doc.Unmarshal(&docMap); err != nil {
		return fmt.Errorf("failed to unmarshal document: %w", err)
	}

	delete(docMap, "_id")

	indexName := w.getElasticsearchIndex(event.Collection)

	document := searchengine.Document{
		ID:     event.DocumentID,
		Fields: docMap,
	}

	if err := w.esEngine.Update(ctx, indexName, event.DocumentID, document); err != nil {
		// 如果更新失败，尝试创建
		return w.handleInsert(ctx, event)
	}

	if w.zapLogger != nil {
		w.zapLogger.Info("Document updated",
			zap.String("document_id", event.DocumentID),
			zap.String("index", indexName),
		)
	}

	return nil
}

// handleDelete 处理删除事件
func (w *WorkerImpl) handleDelete(ctx context.Context, event *search.SyncEvent) error {
	indexName := w.getElasticsearchIndex(event.Collection)

	if err := w.esEngine.Delete(ctx, indexName, event.DocumentID); err != nil {
		// 忽略文档不存在的错误
		if err.Error() != "document not found" {
			return fmt.Errorf("failed to delete document: %w", err)
		}
	}

	if w.zapLogger != nil {
		w.zapLogger.Info("Document deleted",
			zap.String("document_id", event.DocumentID),
			zap.String("index", indexName),
		)
	}

	return nil
}

// getElasticsearchIndex 获取 ES 索引名称
func (w *WorkerImpl) getElasticsearchIndex(collection string) string {
	indexMap := map[string]string{
		"books":     "books_search",
		"projects":  "projects_search",
		"documents": "documents_search",
		"users":     "users_search",
	}

	if indexName, ok := indexMap[collection]; ok {
		return indexName
	}

	return collection + "_search"
}

// Stop 停止 Worker
func (w *WorkerImpl) Stop() error {
	w.status.Running = false
	w.cancel()

	if w.logger != nil {
		w.logger.Println("Sync worker stopped")
	}
	if w.zapLogger != nil {
		w.zapLogger.Info("Sync worker stopped")
	}

	return nil
}

// Status 获取同步状态
func (w *WorkerImpl) Status() (*SyncStatus, error) {
	w.status.TotalEvents = atomic.LoadInt64(&w.eventCounts.total)
	w.status.SuccessEvents = atomic.LoadInt64(&w.eventCounts.success)
	w.status.FailedEvents = atomic.LoadInt64(&w.eventCounts.failed)
	return w.status, nil
}

// GetStats 获取详细统计信息
func (w *WorkerImpl) GetStats(ctx context.Context) (map[string]interface{}, error) {
	// 获取队列长度
	queueLen, err := w.redisClient.LLen(ctx, "search:sync:events").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get queue length: %w", err)
	}

	// 获取死信队列长度
	dlqLen := int64(0)
	if w.dlq != nil {
		dlqLen, err = w.redisClient.LLen(ctx, "search:sync:dlq").Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get DLQ length: %w", err)
		}
	}

	return map[string]interface{}{
		"queue_length": queueLen,
		"dlq_length":   dlqLen,
		"total_events": atomic.LoadInt64(&w.eventCounts.total),
		"success_events": atomic.LoadInt64(&w.eventCounts.success),
		"failed_events": atomic.LoadInt64(&w.eventCounts.failed),
		"running": w.status.Running,
	}, nil
}
