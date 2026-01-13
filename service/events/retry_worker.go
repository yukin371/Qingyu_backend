package events

import (
	"context"
	"log"
	"sync"
	"time"

	"Qingyu_backend/service/base"
)

// RetryWorker 后台重试处理器
// 定期从重试队列中获取失败的事件并重新处理
type RetryWorker struct {
	retryQueue      RetryQueue
	retryPolicy     RetryPolicy
	eventBus        base.EventBus
	checkInterval   time.Duration
	batchSize       int64
	stopCh          chan struct{}
	wg              WaitGroupWrapper
	deadLetterQueue DeadLetterQueue
}

// WaitGroupWrapper waitGroup 包装接口
type WaitGroupWrapper interface {
	Add(delta int)
	Done()
	Wait()
}

// StandardWaitGroup 标准的 waitGroup 实现
type StandardWaitGroup struct {
	wg sync.WaitGroup
}

func (w *StandardWaitGroup) Add(delta int) {
	w.wg.Add(delta)
}

func (w *StandardWaitGroup) Done() {
	w.wg.Done()
}

func (w *StandardWaitGroup) Wait() {
	w.wg.Wait()
}

// NewRetryWorker 创建后台重试处理器
func NewRetryWorker(
	retryQueue RetryQueue,
	deadLetterQueue DeadLetterQueue,
	retryPolicy RetryPolicy,
	eventBus base.EventBus,
	checkInterval time.Duration,
	batchSize int64,
) *RetryWorker {
	if checkInterval == 0 {
		checkInterval = 10 * time.Second // 默认10秒检查一次
	}
	if batchSize == 0 {
		batchSize = 100 // 默认每次处理100个
	}

	return &RetryWorker{
		retryQueue:      retryQueue,
		deadLetterQueue: deadLetterQueue,
		retryPolicy:     retryPolicy,
		eventBus:        eventBus,
		checkInterval:   checkInterval,
		batchSize:       batchSize,
		stopCh:          make(chan struct{}),
		wg:              &StandardWaitGroup{},
	}
}

// Start 启动后台重试处理器
func (w *RetryWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.processRetryItems(ctx)
}

// Stop 停止后台重试处理器
func (w *RetryWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
}

// processRetryItems 处理重试项
func (w *RetryWorker) processRetryItems(ctx context.Context) {
	defer w.wg.Done()

	ticker := time.NewTicker(w.checkInterval)
	defer ticker.Stop()

	log.Printf("[RetryWorker] 后台重试处理器已启动，检查间隔: %v", w.checkInterval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[RetryWorker] 上下文已取消，停止重试处理器")
			return

		case <-w.stopCh:
			log.Printf("[RetryWorker] 收到停止信号，正在退出...")
			return

		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

// processBatch 批量处理重试项
func (w *RetryWorker) processBatch(ctx context.Context) {
	// 获取需要重试的项
	items, err := w.retryQueue.Get(ctx, w.batchSize)
	if err != nil {
		log.Printf("[RetryWorker] 获取重试项失败: %v", err)
		return
	}

	if len(items) == 0 {
		return // 没有需要重试的项
	}

	log.Printf("[RetryWorker] 获取到 %d 个需要重试的项", len(items))

	for _, item := range items {
		select {
		case <-ctx.Done():
			return
		default:
			w.processItem(ctx, item)
		}
	}
}

// processItem 处理单个重试项
func (w *RetryWorker) processItem(ctx context.Context, item *RetryItem) {
	log.Printf("[RetryWorker] 处理重试项: %s, 尝试次数: %d/%d",
		item.ID, item.Attempt+1, item.MaxRetries)

	// 重建事件
	event := &base.BaseEvent{
		EventType: item.EventType,
		EventData: item.EventData,
		Timestamp: item.EventTimestamp,
		Source:    item.EventSource,
	}

	// 尝试重新处理事件
	err := w.reprocessEvent(ctx, event, item.HandlerName)
	if err == nil {
		// 成功，从重试队列中删除
		if err := w.retryQueue.MarkSuccess(ctx, item.ID); err != nil {
			log.Printf("[RetryWorker] 标记重试成功失败: %v", err)
		} else {
			log.Printf("[RetryWorker] 重试成功: %s", item.ID)
		}
		return
	}

	// 失败，判断是否继续重试
	nextAttempt := item.Attempt + 1
	if nextAttempt >= item.MaxRetries || !w.retryPolicy.ShouldRetry(err, nextAttempt) {
		// 达到最大重试次数或不应重试，移入死信队列
		log.Printf("[RetryWorker] 重试失败，移入死信队列: %s, 错误: %v", item.ID, err)
		if err := w.retryQueue.MarkFailed(ctx, item.ID); err != nil {
			log.Printf("[RetryWorker] 移入死信队列失败: %v", err)
		}
		return
	}

	// 计算下次重试时间并更新
	nextRetryTime := time.Now().Add(w.retryPolicy.GetDelay(nextAttempt))
	if err := w.retryQueue.UpdateAttempt(ctx, item.ID, nextRetryTime, nextAttempt); err != nil {
		log.Printf("[RetryWorker] 更新重试时间失败: %v", err)
	} else {
		log.Printf("[RetryWorker] 将在 %v 后进行第 %d 次重试",
			time.Until(nextRetryTime), nextAttempt+1)
	}
}

// reprocessEvent 重新处理事件
func (w *RetryWorker) reprocessEvent(ctx context.Context, event *base.BaseEvent, handlerName string) error {
	// 这里需要根据 handlerName 获取对应的处理器
	// 由于我们没有处理器注册表，这里暂时通过事件总线重新发布
	// 实际使用时应该维护一个处理器注册表

	log.Printf("[RetryWorker] 重新发布事件: %s, 处理器: %s", event.GetEventType(), handlerName)

	// 重新发布事件（同步）
	return w.eventBus.Publish(ctx, event)
}
