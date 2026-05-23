package events

import (
	"context"
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

	logEventRuntime("info", "后台重试处理器已启动", map[string]interface{}{
		"check_interval": w.checkInterval,
		"batch_size":     w.batchSize,
	})

	for {
		select {
		case <-ctx.Done():
			logEventRuntime("info", "重试处理器因上下文取消而停止", nil)
			return

		case <-w.stopCh:
			logEventRuntime("info", "重试处理器收到停止信号", nil)
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
		logEventRuntime("warn", "获取重试项失败", map[string]interface{}{
			"error": err,
		})
		return
	}

	if len(items) == 0 {
		return // 没有需要重试的项
	}

	logEventRuntime("info", "获取到待重试项", map[string]interface{}{
		"count": len(items),
	})

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
	logEventRuntime("info", "处理重试项", map[string]interface{}{
		"retry_item_id": item.ID,
		"attempt":       item.Attempt + 1,
		"max_retries":   item.MaxRetries,
		"handler":       item.HandlerName,
		"event_type":    item.EventType,
	})

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
			logEventRuntime("warn", "标记重试成功失败", map[string]interface{}{
				"retry_item_id": item.ID,
				"error":         err,
			})
		} else {
			logEventRuntime("info", "重试成功", map[string]interface{}{
				"retry_item_id": item.ID,
			})
		}
		return
	}

	// 失败，判断是否继续重试
	nextAttempt := item.Attempt + 1
	if nextAttempt >= item.MaxRetries || !w.retryPolicy.ShouldRetry(err, nextAttempt) {
		// 达到最大重试次数或不应重试，移入死信队列
		logEventRuntime("warn", "重试失败，移入死信队列", map[string]interface{}{
			"retry_item_id": item.ID,
			"error":         err,
			"attempt":       nextAttempt,
		})
		if err := w.retryQueue.MarkFailed(ctx, item.ID); err != nil {
			logEventRuntime("warn", "移入死信队列失败", map[string]interface{}{
				"retry_item_id": item.ID,
				"error":         err,
			})
		}
		return
	}

	// 计算下次重试时间并更新
	nextRetryTime := time.Now().Add(w.retryPolicy.GetDelay(nextAttempt))
	if err := w.retryQueue.UpdateAttempt(ctx, item.ID, nextRetryTime, nextAttempt); err != nil {
		logEventRuntime("warn", "更新重试时间失败", map[string]interface{}{
			"retry_item_id": item.ID,
			"error":         err,
		})
	} else {
		logEventRuntime("info", "已安排下一次重试", map[string]interface{}{
			"retry_item_id": item.ID,
			"next_attempt":  nextAttempt + 1,
			"delay":         time.Until(nextRetryTime),
		})
	}
}

// reprocessEvent 重新处理事件
func (w *RetryWorker) reprocessEvent(ctx context.Context, event *base.BaseEvent, handlerName string) error {
	// 这里需要根据 handlerName 获取对应的处理器
	// 由于我们没有处理器注册表，这里暂时通过事件总线重新发布
	// 实际使用时应该维护一个处理器注册表

	logEventRuntime("info", "重新发布事件", map[string]interface{}{
		"event_type": event.GetEventType(),
		"handler":    handlerName,
	})

	// 重新发布事件（同步）
	return w.eventBus.Publish(ctx, event)
}
