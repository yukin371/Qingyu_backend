package events

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"

	"Qingyu_backend/service/base"
)

// ObservabilityConfig 可观测性配置
type ObservabilityConfig struct {
	// Prometheus配置
	PrometheusEnabled bool   // 是否启用Prometheus
	PrometheusPort    int    // Prometheus端口
	PrometheusPath    string // Prometheus路径

	// 追踪配置
	TracingEnabled    bool    // 是否启用追踪
	TracingExporter   string  // 追踪导出器: jaeger/stdout
	TracingSampleRate float64 // 追踪采样率
	JaegerEndpoint    string  // Jaeger endpoint
	JaegerAgentHost   string  // Jaeger agent host
	JaegerAgentPort   int     // Jaeger agent port

	// 日志配置
	LogLevel         string   // 日志级别
	LogEncoding      string   // 日志编码: json/console
	LogOutputPaths   []string // 日志输出路径
	EnableStackTrace bool     // 是否启用堆栈跟踪

	// 服务信息
	ServiceName    string
	ServiceVersion string
	Environment    string
}

// DefaultObservabilityConfig 默认可观测性配置
func DefaultObservabilityConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		PrometheusEnabled: true,
		PrometheusPort:    9090,
		PrometheusPath:    "/metrics",
		TracingEnabled:    true,
		TracingExporter:   "stdout", // 默认使用stdout，生产环境使用jaeger
		TracingSampleRate: 1.0,      // 100%采样
		LogLevel:          "info",
		LogEncoding:       "json",
		LogOutputPaths:    []string{"stdout"},
		EnableStackTrace:  false,
		ServiceName:       "qingyu-backend",
		ServiceVersion:    "1.0.0",
		Environment:       "development",
	}
}

// ObservabilityManager 可观测性管理器
type ObservabilityManager struct {
	config        *ObservabilityConfig
	metrics       *Metrics
	tracing       *Tracing
	logger        *zap.Logger
	eventLogger   *EventLogger
	shutdownFuncs []func() error
}

// NewObservabilityManager 创建可观测性管理器
func NewObservabilityManager(config *ObservabilityConfig) (*ObservabilityManager, error) {
	if config == nil {
		config = DefaultObservabilityConfig()
	}

	manager := &ObservabilityManager{
		config:        config,
		shutdownFuncs: make([]func() error, 0),
	}

	// 初始化日志
	if err := manager.initLogging(); err != nil {
		return nil, fmt.Errorf("failed to initialize logging: %w", err)
	}

	// 初始化指标
	manager.metrics = NewMetrics()

	// 初始化追踪
	if err := manager.initTracing(); err != nil {
		manager.logger.Warn("Failed to initialize tracing", zap.Error(err))
	}

	manager.logger.Info("Observability manager initialized",
		zap.String("service", config.ServiceName),
		zap.Bool("prometheus_enabled", config.PrometheusEnabled),
		zap.Bool("tracing_enabled", config.TracingEnabled),
	)

	return manager, nil
}

// initLogging 初始化日志
func (m *ObservabilityManager) initLogging() error {
	logConfig := &LoggingConfig{
		Level:            m.config.LogLevel,
		Encoding:         m.config.LogEncoding,
		EnableStackTrace: m.config.EnableStackTrace,
		EnableCaller:     true,
		OutputPaths:      m.config.LogOutputPaths,
	}

	eventLogger, err := NewEventLogger(logConfig)
	if err != nil {
		return err
	}

	m.logger = eventLogger.logger
	m.eventLogger = eventLogger

	return nil
}

// initTracing 初始化追踪
func (m *ObservabilityManager) initTracing() error {
	if !m.config.TracingEnabled {
		m.logger.Info("Tracing disabled")
		return nil
	}

	m.tracing = NewTracing(m.config.ServiceName, true)

	// 创建资源
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.config.ServiceName),
			semconv.ServiceVersionKey.String(m.config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(m.config.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建导出器（目前只支持stdout导出器，用于开发环境）
	// 如需使用Jaeger或OTLP，需要添加相应的依赖包
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return fmt.Errorf("failed to create stdout exporter: %w", err)
	}

	// 创建批量处理器
	batchProcessor := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithBatchTimeout(5*time.Second),
		sdktrace.WithMaxExportBatchSize(100),
	)

	// 创建追踪提供者
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(m.config.TracingSampleRate)),
		sdktrace.WithSpanProcessor(batchProcessor),
	)

	// 注册全局追踪提供者
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// 添加关闭函数
	m.shutdownFuncs = append(m.shutdownFuncs, func() error {
		return tp.Shutdown(context.Background())
	})

	m.logger.Info("Tracing initialized",
		zap.String("exporter", m.config.TracingExporter),
		zap.Float64("sample_rate", m.config.TracingSampleRate),
	)

	return nil
}

// GetMetrics 获取指标
func (m *ObservabilityManager) GetMetrics() *Metrics {
	return m.metrics
}

// GetTracing 获取追踪
func (m *ObservabilityManager) GetTracing() *Tracing {
	return m.tracing
}

// GetLogger 获取日志记录器
func (m *ObservabilityManager) GetLogger() *zap.Logger {
	return m.logger
}

// GetEventLogger 获取事件日志记录器
func (m *ObservabilityManager) GetEventLogger() *EventLogger {
	return m.eventLogger
}

// WrapEventBus 包装事件总线（添加可观测性）
func (m *ObservabilityManager) WrapEventBus(eventBus base.EventBus) base.EventBus {
	// 添加追踪
	bus := NewTraceableEventBus(eventBus, m.tracing, m.metrics)

	m.logger.Info("Event bus wrapped with observability")

	return bus
}

// WrapEventHandler 包装事件处理器（添加可观测性）
func (m *ObservabilityManager) WrapEventHandler(handler base.EventHandler) base.EventHandler {
	// 添加追踪
	return NewTraceableEventHandler(handler, m.tracing, m.metrics)
}

// GetPrometheusHandler 获取Prometheus HTTP处理器
func (m *ObservabilityManager) GetPrometheusHandler() http.Handler {
	if !m.config.PrometheusEnabled {
		return nil
	}
	return promhttp.Handler()
}

// Shutdown 关闭可观测性管理器
func (m *ObservabilityManager) Shutdown(ctx context.Context) error {
	m.logger.Info("Shutting down observability manager")

	// 刷新日志
	if m.eventLogger != nil {
		m.eventLogger.Flush()
	}

	// 关闭追踪
	var errs []error
	for _, fn := range m.shutdownFuncs {
		if err := fn(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	m.logger.Info("Observability manager shutdown complete")
	return nil
}

// GetMetricsSummary 获取指标摘要（用于监控面板）
func (m *ObservabilityManager) GetMetricsSummary(ctx context.Context) (*MetricsSummary, error) {
	// 这里可以从各种指标收集中获取当前值
	// 简化版本，实际使用时应该从Prometheus查询

	summary := &MetricsSummary{
		EventsPublishedTotal: m.getCounterTotal(m.metrics.EventsPublished),
		EventsHandledTotal:   m.getCounterTotal(m.metrics.EventsHandled),
		RetryQueueSize:       m.getGaugeValue(m.metrics.RetryQueueSize),
		DeadLetterQueueSize:  m.getGaugeValue(m.metrics.DeadLetterQueueSize),
		ActiveHandlers:       m.getActiveHandlerCount(),
		CollectedAt:          time.Now(),
	}

	return summary, nil
}

// MetricsSummary 指标摘要
type MetricsSummary struct {
	EventsPublishedTotal map[string]int64
	EventsHandledTotal   map[string]int64
	RetryQueueSize       float64
	DeadLetterQueueSize  float64
	ActiveHandlers       int
	CollectedAt          time.Time
}

// 辅助方法
func (m *ObservabilityManager) getCounterTotal(counterVec *prometheus.CounterVec) map[string]int64 {
	result := make(map[string]int64)
	// 简化实现，实际应该使用prometheus的metrics API
	return result
}

func (m *ObservabilityManager) getGaugeValue(gauge prometheus.Gauge) float64 {
	// 简化实现，实际应该使用prometheus的metrics API
	return 0
}

func (m *ObservabilityManager) getActiveHandlerCount() int {
	// 简化实现
	return 0
}

// MonitorService 监控服务
// 提供监控数据查询和仪表板功能
type MonitorService struct {
	manager         *ObservabilityManager
	retryQueue      RetryQueue
	deadLetterQueue DeadLetterQueue
	eventStore      EventStore
}

// NewMonitorService 创建监控服务
func NewMonitorService(
	manager *ObservabilityManager,
	retryQueue RetryQueue,
	deadLetterQueue DeadLetterQueue,
	eventStore EventStore,
) *MonitorService {
	return &MonitorService{
		manager:         manager,
		retryQueue:      retryQueue,
		deadLetterQueue: deadLetterQueue,
		eventStore:      eventStore,
	}
}

// GetDashboardData 获取仪表板数据
func (s *MonitorService) GetDashboardData(ctx context.Context) (*DashboardData, error) {
	data := &DashboardData{
		CollectedAt: time.Now(),
	}

	// 获取队列大小
	if s.retryQueue != nil {
		count, _ := s.retryQueue.Count(ctx)
		data.RetryQueueSize = count
	}

	if s.deadLetterQueue != nil {
		count, _ := s.deadLetterQueue.Count(ctx)
		data.DeadLetterQueueSize = count
	}

	// 获取事件统计
	if s.eventStore != nil {
		// 获取过去1小时的事件数量
		start := time.Now().Add(-1 * time.Hour)
		count, _ := s.eventStore.Count(ctx, EventFilter{
			StartTime: &start,
		})
		data.EventsLastHour = count
	}

	return data, nil
}

// DashboardData 仪表板数据
type DashboardData struct {
	CollectedAt         time.Time `json:"collected_at"`
	RetryQueueSize      int64     `json:"retry_queue_size"`
	DeadLetterQueueSize int64     `json:"dead_letter_queue_size"`
	EventsLastHour      int64     `json:"events_last_hour"`
	EventsToday         int64     `json:"events_today"`
	ActiveHandlers      int       `json:"active_handlers"`
	// 可以添加更多指标...
}
