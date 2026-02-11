package events

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 事件驱动系统指标
type Metrics struct {
	// 发布指标
	EventsPublished         *prometheus.CounterVec
	EventsPublishedDuration *prometheus.HistogramVec

	// 处理指标
	EventsHandled          *prometheus.CounterVec
	EventsHandlingDuration *prometheus.HistogramVec
	EventsHandlingErrors   *prometheus.CounterVec

	// 重试指标
	EventsRetryQueued *prometheus.CounterVec
	RetryQueueSize    prometheus.Gauge
	EventsRetried     *prometheus.CounterVec
	EventsRetryFailed *prometheus.CounterVec

	// 死信队列指标
	EventsDeadLettered  *prometheus.CounterVec
	DeadLetterQueueSize prometheus.Gauge
	EventsReprocessed   *prometheus.CounterVec

	// 持久化指标
	EventsStored         *prometheus.CounterVec
	EventStorageDuration *prometheus.HistogramVec
	EventStorageErrors   *prometheus.CounterVec

	// 处理器指标
	HandlerActiveCount *prometheus.GaugeVec

	// Replay 指标
	EventsReplayTotal    *prometheus.CounterVec
	EventsReplayFailed   *prometheus.CounterVec
	EventsReplayDuration *prometheus.HistogramVec
}

// NewMetrics 创建事件指标
func NewMetrics() *Metrics {
	return &Metrics{
		// 发布指标
		EventsPublished: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_published_total",
				Help: "Total number of events published.",
			},
			[]string{"event_type", "source"},
		),
		EventsPublishedDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "events_publish_duration_seconds",
				Help:    "Duration of event publishing in seconds.",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"event_type", "source"},
		),

		// 处理指标
		EventsHandled: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_handled_total",
				Help: "Total number of events handled.",
			},
			[]string{"event_type", "handler", "status"},
		),
		EventsHandlingDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "events_handling_duration_seconds",
				Help:    "Duration of event handling in seconds.",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"event_type", "handler"},
		),
		EventsHandlingErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_handling_errors_total",
				Help: "Total number of event handling errors.",
			},
			[]string{"event_type", "handler", "error_type"},
		),

		// 重试指标
		EventsRetryQueued: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_retry_queued_total",
				Help: "Total number of events queued for retry.",
			},
			[]string{"event_type", "handler"},
		),
		RetryQueueSize: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "events_retry_queue_size",
				Help: "Current size of the retry queue.",
			},
		),
		EventsRetried: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_retried_total",
				Help: "Total number of event retry attempts.",
			},
			[]string{"event_type", "handler", "attempt"},
		),
		EventsRetryFailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_retry_failed_total",
				Help: "Total number of failed event retry attempts.",
			},
			[]string{"event_type", "handler"},
		),

		// 死信队列指标
		EventsDeadLettered: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_dead_letter_total",
				Help: "Total number of events moved to dead letter queue.",
			},
			[]string{"event_type", "handler"},
		),
		DeadLetterQueueSize: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "events_dead_letter_queue_size",
				Help: "Current size of the dead letter queue.",
			},
		),
		EventsReprocessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_reprocessed_total",
				Help: "Total number of reprocessed events from dead letter queue.",
			},
			[]string{"event_type", "status"},
		),

		// 持久化指标
		EventsStored: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_stored_total",
				Help: "Total number of events stored.",
			},
			[]string{"event_type"},
		),
		EventStorageDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "events_storage_duration_seconds",
				Help:    "Duration of event storage in seconds.",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
			},
			[]string{"event_type"},
		),
		EventStorageErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_storage_errors_total",
				Help: "Total number of event storage errors.",
			},
			[]string{"event_type", "error_type"},
		),

		// 处理器指标
		HandlerActiveCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "event_handlers_active",
				Help: "Number of active event handlers.",
			},
			[]string{"handler", "event_type"},
		),

		// Replay 指标
		EventsReplayTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_replay_total",
				Help: "Total number of events replayed.",
			},
			[]string{"event_type", "status"},
		),
		EventsReplayFailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "events_replay_failed_total",
				Help: "Total number of failed event replay attempts.",
			},
			[]string{"event_type", "error_type"},
		),
		EventsReplayDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "events_replay_duration_seconds",
				Help:    "Duration of event replay in seconds.",
				Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10, 30, 60, 120, 300},
			},
			[]string{"event_type"},
		),
	}
}

// RecordEventPublished 记录事件发布
func (m *Metrics) RecordEventPublished(eventType, source string, duration time.Duration) {
	m.EventsPublished.WithLabelValues(eventType, source).Inc()
	m.EventsPublishedDuration.WithLabelValues(eventType, source).Observe(duration.Seconds())
}

// RecordEventHandled 记录事件处理
func (m *Metrics) RecordEventHandled(eventType, handler string, duration time.Duration, err error) {
	if err != nil {
		status := "error"
		m.EventsHandled.WithLabelValues(eventType, handler, status).Inc()
		m.EventsHandlingErrors.WithLabelValues(eventType, handler, "unknown").Inc()
	} else {
		status := "success"
		m.EventsHandled.WithLabelValues(eventType, handler, status).Inc()
	}
	m.EventsHandlingDuration.WithLabelValues(eventType, handler).Observe(duration.Seconds())
}

// RecordEventRetryQueued 记录事件加入重试队列
func (m *Metrics) RecordEventRetryQueued(eventType, handler string) {
	m.EventsRetryQueued.WithLabelValues(eventType, handler).Inc()
}

// UpdateRetryQueueSize 更新重试队列大小
func (m *Metrics) UpdateRetryQueueSize(size float64) {
	m.RetryQueueSize.Set(size)
}

// RecordEventRetried 记录事件重试
func (m *Metrics) RecordEventRetried(eventType, handler string, attempt int) {
	m.EventsRetried.WithLabelValues(eventType, handler, string(rune(attempt+'0'))).Inc()
}

// RecordEventRetryFailed 记录重试失败
func (m *Metrics) RecordEventRetryFailed(eventType, handler string) {
	m.EventsRetryFailed.WithLabelValues(eventType, handler).Inc()
}

// RecordEventDeadLettered 记录事件移入死信队列
func (m *Metrics) RecordEventDeadLettered(eventType, handler string) {
	m.EventsDeadLettered.WithLabelValues(eventType, handler).Inc()
}

// UpdateDeadLetterQueueSize 更新死信队列大小
func (m *Metrics) UpdateDeadLetterQueueSize(size float64) {
	m.DeadLetterQueueSize.Set(size)
}

// RecordEventReprocessed 记录事件重新处理
func (m *Metrics) RecordEventReprocessed(eventType string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	m.EventsReprocessed.WithLabelValues(eventType, status).Inc()
}

// RecordEventStored 记录事件存储
func (m *Metrics) RecordEventStored(eventType string, duration time.Duration, err error) {
	if err != nil {
		m.EventStorageErrors.WithLabelValues(eventType, "unknown").Inc()
	} else {
		m.EventsStored.WithLabelValues(eventType).Inc()
	}
	m.EventStorageDuration.WithLabelValues(eventType).Observe(duration.Seconds())
}

// IncrementHandlerActive 增加活跃处理器计数
func (m *Metrics) IncrementHandlerActive(handler, eventType string) {
	m.HandlerActiveCount.WithLabelValues(handler, eventType).Inc()
}

// DecrementHandlerActive 减少活跃处理器计数
func (m *Metrics) DecrementHandlerActive(handler, eventType string) {
	m.HandlerActiveCount.WithLabelValues(handler, eventType).Dec()
}

// RecordEventReplayed 记录事件回放
func (m *Metrics) RecordEventReplayed(eventType string, replayedCount, failedCount int64, duration time.Duration) {
	status := "success"
	if failedCount > 0 {
		status = "partial_failure"
	}
	m.EventsReplayTotal.WithLabelValues(eventType, status).Add(float64(replayedCount))
	m.EventsReplayDuration.WithLabelValues(eventType).Observe(duration.Seconds())
}

// RecordEventReplayFailed 记录事件回放失败
func (m *Metrics) RecordEventReplayFailed(eventType, errorType string) {
	m.EventsReplayFailed.WithLabelValues(eventType, errorType).Inc()
}

// 全局指标实例
var globalMetrics *Metrics

func init() {
	globalMetrics = NewMetrics()
}

// GetGlobalMetrics 获取全局指标实例
func GetGlobalMetrics() *Metrics {
	return globalMetrics
}
