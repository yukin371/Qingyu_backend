package sync

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ChangeStreamListenerImpl MongoDB Change Stream 监听器实现
type ChangeStreamListenerImpl struct {
	// TODO: 注入 MongoDB client, Redis client
	logger   *log.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	watchers map[string]context.CancelFunc
}

// NewChangeStreamListener 创建 Change Stream 监听器
func NewChangeStreamListener(logger *log.Logger) (*ChangeStreamListenerImpl, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChangeStreamListenerImpl{
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
		watchers: make(map[string]context.CancelFunc),
	}, nil
}

// Watch 监听指定集合的变更
func (l *ChangeStreamListenerImpl) Watch(ctx context.Context, db, collection, index string) error {
	// TODO: 实现 MongoDB Change Stream 监听
	l.logger.Printf("Watching collection: %s.%s -> %s", db, collection, index)

	watcherCtx, watcherCancel := context.WithCancel(ctx)
	l.watchers[collection] = watcherCancel

	go l.watchCollection(watcherCtx, db, collection, index)

	return nil
}

// ProcessChange 处理变更事件
func (l *ChangeStreamListenerImpl) ProcessChange(ctx context.Context, change *ChangeData) error {
	// TODO: 实现变更处理逻辑
	// 1. 根据变更类型创建 SyncEvent
	// 2. 发送到 Redis 队列
	// 3. 处理失败情况

	l.logger.Printf("Processing change: %s on %s.%s", change.Operation, change.Database, change.Collection)
	return fmt.Errorf("process change not implemented")
}

// Stop 停止监听器
func (l *ChangeStreamListenerImpl) Stop() error {
	// 取消所有 watcher
	for collection, cancel := range l.watchers {
		l.logger.Printf("Stopping watcher for collection: %s", collection)
		cancel()
	}
	l.cancel()
	return nil
}

// watchCollection 监听指定集合（内部方法）
func (l *ChangeStreamListenerImpl) watchCollection(ctx context.Context, db, collection, index string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.logger.Printf("Stopped watching collection: %s.%s", db, collection)
			return
		case <-ticker.C:
			// TODO: 实际的监听逻辑
			// 1. 从 MongoDB 获取 Change Stream
			// 2. 解析变更事件
			// 3. 调用 ProcessChange 处理
		}
	}
}
