package events

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"Qingyu_backend/service/base"
)

// PersistedEventBusInterface 持久化事件总线接口
// 定义用于事件管理和回放的接口
type PersistedEventBusInterface interface {
	// Replay 事件回放（从存储中重放事件到处理器）
	Replay(ctx context.Context, handler base.EventHandler, filter EventFilter) (*ReplayResult, error)
}

// PersistedEventBus 持久化事件总线
// 在发布事件的同时持久化到存储
type PersistedEventBus struct {
	mu         sync.RWMutex
	eventBus   base.EventBus    // 内存事件总线
	store      EventStore       // 事件存储
	asyncStore bool             // 是否异步存储
	stopCh     chan struct{}    // 停止通道
	wg         sync.WaitGroup   // 等待组
	eventCh    chan *base.Event // 事件通道（用于异步存储）
}

// NewPersistedEventBus 创建持久化事件总线
func NewPersistedEventBus(eventBus base.EventBus, store EventStore, asyncStore bool) base.EventBus {
	bus := &PersistedEventBus{
		eventBus:   eventBus,
		store:      store,
		asyncStore: asyncStore,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000), // 缓冲1000个事件
	}

	// 如果启用异步存储，启动后台处理
	if asyncStore {
		bus.startAsyncProcessor()
	}

	return bus
}

// Subscribe 订阅事件
func (b *PersistedEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return b.eventBus.Subscribe(eventType, handler)
}

// Unsubscribe 取消订阅
func (b *PersistedEventBus) Unsubscribe(eventType string, handlerName string) error {
	return b.eventBus.Unsubscribe(eventType, handlerName)
}

// Publish 同步发布事件
func (b *PersistedEventBus) Publish(ctx context.Context, event base.Event) error {
	// 1. 持久化事件
	if err := b.storeEvent(ctx, event); err != nil {
		log.Printf("[PersistedEventBus] Failed to store event: %v", err)
		// 持久化失败不影响事件发布
	}

	// 2. 发布到内存总线
	return b.eventBus.Publish(ctx, event)
}

// PublishAsync 异步发布事件
func (b *PersistedEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	// 1. 持久化事件（异步）
	if b.asyncStore {
		select {
		case b.eventCh <- &event:
			// 事件已加入异步队列
		default:
			log.Printf("[PersistedEventBus] Event channel full, storing synchronously")
			// 队列满时同步存储
			if err := b.storeEvent(context.Background(), event); err != nil {
				log.Printf("[PersistedEventBus] Failed to store event: %v", err)
			}
		}
	} else {
		if err := b.storeEvent(ctx, event); err != nil {
			log.Printf("[PersistedEventBus] Failed to store event: %v", err)
		}
	}

	// 2. 发布到内存总线
	return b.eventBus.PublishAsync(ctx, event)
}

// storeEvent 存储事件
func (b *PersistedEventBus) storeEvent(ctx context.Context, event base.Event) error {
	return b.store.Store(ctx, event)
}

// startAsyncProcessor 启动异步存储处理器
func (b *PersistedEventBus) startAsyncProcessor() {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		b.processAsyncEvents()
	}()
}

// processAsyncEvents 处理异步事件
func (b *PersistedEventBus) processAsyncEvents() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	batch := make([]base.Event, 0, 100)

	for {
		select {
		case <-b.stopCh:
			// 处理剩余事件
			if len(batch) > 0 {
				b.storeBatch(context.Background(), batch)
			}
			return

		case event := <-b.eventCh:
			batch = append(batch, *event)

			// 批量存储
			if len(batch) >= 100 {
				b.storeBatch(context.Background(), batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			// 定期刷新
			if len(batch) > 0 {
				b.storeBatch(context.Background(), batch)
				batch = batch[:0]
			}
		}
	}
}

// storeBatch 批量存储事件
func (b *PersistedEventBus) storeBatch(ctx context.Context, events []base.Event) {
	if len(events) == 0 {
		return
	}

	if err := b.store.StoreBatch(ctx, events); err != nil {
		log.Printf("[PersistedEventBus] Failed to store event batch: %v", err)
		// 批量失败时，尝试逐个存储
		for _, event := range events {
			if err := b.store.Store(ctx, event); err != nil {
				log.Printf("[PersistedEventBus] Failed to store event: %v", err)
			}
		}
	}
}

// Replay 从存储中重放事件到处理器
// 该方法将请求转发到底层的事件存储实现
//
// 参数:
//   ctx - 上下文，用于超时控制和取消
//   handler - 事件处理器，用于处理每个重放的事件
//   filter - 过滤条件，可选的类型、来源、时间范围等
//
// 返回:
//   *ReplayResult - 重放结果统计
//   error - 错误信息，如果有的话
//
// 注意:
//   - 该方法线程安全，转发给底层store实现
//   - 具体的并发安全性取决于底层store实现
//   - 指标和日志由底层store实现处理
//   - 请参考底层store的Replay方法文档了解详细的参数验证和错误处理
//
// 示例:
//   result, err := bus.Replay(ctx, myHandler, EventFilter{
//       EventType: "ChapterPublishedEvent",
//       Limit: 1000,
//   })
func (b *PersistedEventBus) Replay(ctx context.Context, handler base.EventHandler, filter EventFilter) (*ReplayResult, error) {
	// 指标和日志由底层store实现，这里只做转发
	return b.store.Replay(ctx, handler, filter)
}

// Close 关闭事件总线
func (b *PersistedEventBus) Close() error {
	close(b.stopCh)
	b.wg.Wait()
	return nil
}

// RetryableEventBus 可重试的事件总线
// 在处理器失败时自动重试
type RetryableEventBus struct {
	mu              sync.RWMutex
	eventBus        base.EventBus
	retryQueue      RetryQueue
	deadLetterQueue DeadLetterQueue
	retryPolicy     RetryPolicy
	maxRetries      int
}

// NewRetryableEventBus 创建可重试的事件总线
func NewRetryableEventBus(eventBus base.EventBus, retryQueue RetryQueue,
	deadLetterQueue DeadLetterQueue, retryPolicy RetryPolicy, maxRetries int) base.EventBus {
	return &RetryableEventBus{
		eventBus:        eventBus,
		retryQueue:      retryQueue,
		deadLetterQueue: deadLetterQueue,
		retryPolicy:     retryPolicy,
		maxRetries:      maxRetries,
	}
}

// Subscribe 订阅事件
func (b *RetryableEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return b.eventBus.Subscribe(eventType, handler)
}

// Unsubscribe 取消订阅
func (b *RetryableEventBus) Unsubscribe(eventType string, handlerName string) error {
	return b.eventBus.Unsubscribe(eventType, handlerName)
}

// Publish 同步发布事件
func (b *RetryableEventBus) Publish(ctx context.Context, event base.Event) error {
	return b.eventBus.Publish(ctx, event)
}

// PublishAsync 异步发布事件
func (b *RetryableEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	return b.eventBus.PublishAsync(ctx, event)
}

// RetryableEventHandler 可重试的事件处理器包装器
type RetryableEventHandler struct {
	handler     base.EventHandler
	retryQueue  RetryQueue
	retryPolicy RetryPolicy
	maxRetries  int
}

// NewRetryableEventHandler 创建可重试的事件处理器
func NewRetryableEventHandler(handler base.EventHandler, retryQueue RetryQueue,
	retryPolicy RetryPolicy, maxRetries int) base.EventHandler {
	return &RetryableEventHandler{
		handler:     handler,
		retryQueue:  retryQueue,
		retryPolicy: retryPolicy,
		maxRetries:  maxRetries,
	}
}

// Handle 处理事件（带重试）
func (h *RetryableEventHandler) Handle(ctx context.Context, event base.Event) error {
	var lastErr error

	for attempt := 0; attempt <= h.maxRetries; attempt++ {
		err := h.handler.Handle(ctx, event)
		if err == nil {
			return nil // 成功
		}

		lastErr = err

		// 检查是否应该重试
		if !h.retryPolicy.ShouldRetry(err, attempt) {
			break
		}

		// 添加到重试队列
		if h.retryQueue != nil {
			_ = h.retryQueue.Add(ctx, event, h.handler.GetHandlerName(), err, attempt+1)
		}

		// 等待后重试
		delay := h.retryPolicy.GetDelay(attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// 重试
		}
	}

	return fmt.Errorf("handler failed after %d retries: %w", h.maxRetries, lastErr)
}

// GetHandlerName 获取处理器名称
func (h *RetryableEventHandler) GetHandlerName() string {
	return h.handler.GetHandlerName()
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *RetryableEventHandler) GetSupportedEventTypes() []string {
	return h.handler.GetSupportedEventTypes()
}
