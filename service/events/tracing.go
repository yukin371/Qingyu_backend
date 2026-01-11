package events

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"Qingyu_backend/service/base"
)

// Tracing 事件追踪器
type Tracing struct {
	tracer             trace.Tracer
	enabled            bool
	serviceName        string
}

// NewTracing 创建事件追踪器
func NewTracing(serviceName string, enabled bool) *Tracing {
	return &Tracing{
		tracer:      otel.Tracer(serviceName),
		enabled:     enabled,
		serviceName: serviceName,
	}
}

// StartSpan 开始一个追踪span
func (t *Tracing) StartSpan(ctx context.Context, spanName string, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	if !t.enabled {
		return ctx, trace.SpanFromContext(ctx)
	}

	// 添加通用属性
	allAttrs := append([]attribute.KeyValue{
		attribute.String("service", t.serviceName),
	}, attributes...)

	ctx, span := t.tracer.Start(ctx, spanName, trace.WithAttributes(allAttrs...))
	return ctx, span
}

// RecordError 记录错误到span
func (t *Tracing) RecordError(span trace.Span, err error) {
	if !t.enabled || span == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// RecordEventPublish 记录事件发布追踪
func (t *Tracing) RecordEventPublish(ctx context.Context, eventType, source string) (context.Context, trace.Span) {
	return t.StartSpan(ctx, "EventBus.Publish",
		attribute.String("event.type", eventType),
		attribute.String("event.source", source),
	)
}

// RecordEventHandle 记录事件处理追踪
func (t *Tracing) RecordEventHandle(ctx context.Context, eventType, handler string) (context.Context, trace.Span) {
	return t.StartSpan(ctx, "EventHandler.Handle",
		attribute.String("event.type", eventType),
		attribute.String("handler.name", handler),
	)
}

// RecordEventStorage 记录事件存储追踪
func (t *Tracing) RecordEventStorage(ctx context.Context, eventType string) (context.Context, trace.Span) {
	return t.StartSpan(ctx, "EventStore.Store",
		attribute.String("event.type", eventType),
	)
}

// RecordRetryAttempt 记录重试尝试追踪
func (t *Tracing) RecordRetryAttempt(ctx context.Context, eventType, handler string, attempt int) (context.Context, trace.Span) {
	return t.StartSpan(ctx, "EventRetry",
		attribute.String("event.type", eventType),
		attribute.String("handler.name", handler),
		attribute.Int("retry.attempt", attempt),
	)
}

// TraceableEventBus 可追踪的事件总线包装器
type TraceableEventBus struct {
	eventBus base.EventBus
	tracing  *Tracing
	metrics  *Metrics
}

// NewTraceableEventBus 创建可追踪的事件总线
func NewTraceableEventBus(eventBus base.EventBus, tracing *Tracing, metrics *Metrics) base.EventBus {
	return &TraceableEventBus{
		eventBus: eventBus,
		tracing:  tracing,
		metrics:  metrics,
	}
}

// Subscribe 订阅事件
func (b *TraceableEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return b.eventBus.Subscribe(eventType, handler)
}

// Unsubscribe 取消订阅
func (b *TraceableEventBus) Unsubscribe(eventType string, handlerName string) error {
	return b.eventBus.Unsubscribe(eventType, handlerName)
}

// Publish 同步发布事件（带追踪）
func (b *TraceableEventBus) Publish(ctx context.Context, event base.Event) error {
	// 开始发布追踪
	ctx, span := b.tracing.RecordEventPublish(ctx, event.GetEventType(), event.GetSource())
	defer span.End()

	startTime := time.Now()

	// 调用底层发布
	err := b.eventBus.Publish(ctx, event)

	// 记录指标
	duration := time.Since(startTime)
	if b.metrics != nil {
		b.metrics.RecordEventPublished(event.GetEventType(), event.GetSource(), duration)
	}

	// 记录错误
	if err != nil && b.tracing.enabled {
		b.tracing.RecordError(span, err)
	} else {
		span.SetStatus(codes.Ok, "published")
	}

	return err
}

// PublishAsync 异步发布事件（带追踪）
func (b *TraceableEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	// 异步发布也需要追踪，但不等待完成
	ctx, span := b.tracing.RecordEventPublish(ctx, event.GetEventType(), event.GetSource())

	// 在后台goroutine中完成span
	go func() {
		defer span.End()
	}()

	startTime := time.Now()

	// 调用底层发布
	err := b.eventBus.PublishAsync(ctx, event)

	// 记录指标
	duration := time.Since(startTime)
	if b.metrics != nil {
		b.metrics.RecordEventPublished(event.GetEventType(), event.GetSource(), duration)
	}

	// 记录错误
	if err != nil && b.tracing.enabled {
		b.tracing.RecordError(span, err)
	} else {
		span.SetStatus(codes.Ok, "published")
	}

	return err
}

// TraceableEventHandler 可追踪的事件处理器包装器
type TraceableEventHandler struct {
	handler base.EventHandler
	tracing  *Tracing
	metrics  *Metrics
}

// NewTraceableEventHandler 创建可追踪的事件处理器
func NewTraceableEventHandler(handler base.EventHandler, tracing *Tracing, metrics *Metrics) base.EventHandler {
	return &TraceableEventHandler{
		handler: handler,
		tracing:  tracing,
		metrics:  metrics,
	}
}

// Handle 处理事件（带追踪）
func (h *TraceableEventHandler) Handle(ctx context.Context, event base.Event) error {
	// 开始处理追踪
	ctx, span := h.tracing.RecordEventHandle(ctx, event.GetEventType(), h.GetHandlerName())
	defer span.End()

	startTime := time.Now()

	// 增加活跃处理器计数
	if h.metrics != nil {
		h.metrics.IncrementHandlerActive(h.GetHandlerName(), event.GetEventType())
		defer h.metrics.DecrementHandlerActive(h.GetHandlerName(), event.GetEventType())
	}

	// 调用底层处理器
	err := h.handler.Handle(ctx, event)

	// 记录指标
	duration := time.Since(startTime)
	if h.metrics != nil {
		h.metrics.RecordEventHandled(event.GetEventType(), h.GetHandlerName(), duration, err)
	}

	// 记录错误
	if err != nil && h.tracing.enabled {
		h.tracing.RecordError(span, err)
	} else {
		span.SetStatus(codes.Ok, "handled")
	}

	return err
}

// GetHandlerName 获取处理器名称
func (h *TraceableEventHandler) GetHandlerName() string {
	return h.handler.GetHandlerName()
}

// GetSupportedEventTypes 获取支持的事件类型
func (h *TraceableEventHandler) GetSupportedEventTypes() []string {
	return h.handler.GetSupportedEventTypes()
}
