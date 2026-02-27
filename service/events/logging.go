package events

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"Qingyu_backend/service/base"
)

// EventLogger 事件日志记录器
type EventLogger struct {
	logger *zap.Logger
	config *LoggingConfig
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level            string   // 日志级别: debug/info/warn/error
	Encoding         string   // 编码格式: json/console
	EnableStackTrace bool     // 是否启用堆栈跟踪
	EnableCaller     bool     // 是否显示调用位置
	OutputPaths      []string // 输出路径
}

// DefaultLoggingConfig 默认日志配置
func DefaultLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		Level:            "info",
		Encoding:         "json",
		EnableStackTrace: false,
		EnableCaller:     true,
		OutputPaths:      []string{"stdout"},
	}
}

// NewEventLogger 创建事件日志记录器
func NewEventLogger(config *LoggingConfig) (*EventLogger, error) {
	if config == nil {
		config = DefaultLoggingConfig()
	}

	// 解析日志级别
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置输出
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.Lock(zapcore.AddSync(GetOutputWriter(config.OutputPaths...)))),
		level,
	)

	// 创建logger
	logger := zap.New(core,
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &EventLogger{
		logger: logger,
		config: config,
	}, nil
}

// LogEventPublished 记录事件发布
func (l *EventLogger) LogEventPublished(ctx context.Context, event base.Event, duration time.Duration) {
	l.logger.Info("Event published",
		zap.String("event_type", event.GetEventType()),
		zap.String("source", event.GetSource()),
		zap.Duration("duration", duration),
		zap.Time("timestamp", event.GetTimestamp()),
	)
}

// LogEventHandled 记录事件处理
func (l *EventLogger) LogEventHandled(ctx context.Context, event base.Event, handler string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("event_type", event.GetEventType()),
		zap.String("handler", handler),
		zap.Duration("duration", duration),
		zap.Time("timestamp", event.GetTimestamp()),
	}

	if err != nil {
		fields = append(fields,
			zap.Error(err),
			zap.String("status", "error"),
		)
		l.logger.Error("Event handle failed", fields...)
	} else {
		fields = append(fields,
			zap.String("status", "success"),
		)
		l.logger.Info("Event handled successfully", fields...)
	}
}

// LogEventStored 记录事件存储
func (l *EventLogger) LogEventStored(ctx context.Context, event base.Event, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("event_type", event.GetEventType()),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		l.logger.Error("Event storage failed", fields...)
	} else {
		l.logger.Info("Event stored successfully", fields...)
	}
}

// LogEventRetryQueued 记录事件加入重试队列
func (l *EventLogger) LogEventRetryQueued(ctx context.Context, event base.Event, handler string, attempt int, err error) {
	l.logger.Warn("Event queued for retry",
		zap.String("event_type", event.GetEventType()),
		zap.String("handler", handler),
		zap.Int("attempt", attempt),
		zap.Error(err),
	)
}

// LogEventRetryAttempt 记录重试尝试
func (l *EventLogger) LogEventRetryAttempt(ctx context.Context, event base.Event, handler string, attempt int) {
	l.logger.Info("Event retry attempt",
		zap.String("event_type", event.GetEventType()),
		zap.String("handler", handler),
		zap.Int("attempt", attempt),
	)
}

// LogEventDeadLettered 记录事件移入死信队列
func (l *EventLogger) LogEventDeadLettered(ctx context.Context, event base.Event, handler string, attempt int, err error) {
	l.logger.Error("Event moved to dead letter queue",
		zap.String("event_type", event.GetEventType()),
		zap.String("handler", handler),
		zap.Int("final_attempt", attempt),
		zap.Error(err),
	)
}

// LogBatchProcess 记录批量处理
func (l *EventLogger) LogBatchProcess(ctx context.Context, batchSize int, queueType string) {
	l.logger.Info("Batch processing started",
		zap.Int("batch_size", batchSize),
		zap.String("queue_type", queueType),
	)
}

// LogBatchComplete 记录批量处理完成
func (l *EventLogger) LogBatchComplete(ctx context.Context, batchSize int, successCount, failCount int, queueType string) {
	l.logger.Info("Batch processing completed",
		zap.Int("batch_size", batchSize),
		zap.Int("success_count", successCount),
		zap.Int("fail_count", failCount),
		zap.String("queue_type", queueType),
	)
}

// LogReplayStarted 记录事件回放开始
func (l *EventLogger) LogReplayStarted(ctx context.Context, eventType string, startTime, endTime *time.Time) {
	fields := []zap.Field{
		zap.String("event_type", eventType),
	}

	// 添加trace_id（如果可用）
	if traceID := getTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// 添加时间范围
	if startTime != nil {
		fields = append(fields, zap.Time("start_time", *startTime))
	}
	if endTime != nil {
		fields = append(fields, zap.Time("end_time", *endTime))
	}

	// 添加时间范围描述
	if startTime != nil && endTime != nil {
		duration := endTime.Sub(*startTime)
		fields = append(fields, zap.Duration("range", duration))
	}

	l.logger.Info("Event replay started", fields...)
}

// LogReplayCompleted 记录事件回放完成
func (l *EventLogger) LogReplayCompleted(ctx context.Context, eventType string, result *ReplayResult) {
	fields := []zap.Field{
		zap.String("event_type", eventType),
		zap.String("result", "success"),
		zap.Int64("replayed_count", result.ReplayedCount),
		zap.Int64("failed_count", result.FailedCount),
		zap.Int64("skipped_count", result.SkippedCount),
		zap.Duration("duration", result.Duration),
		zap.Int64("duration_ms", result.Duration.Milliseconds()),
	}

	// 添加trace_id（如果可用）
	if traceID := getTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// 根据结果确定日志级别
	if result.FailedCount > 0 {
		l.logger.Warn("Event replay completed with failures", fields...)
	} else {
		l.logger.Info("Event replay completed successfully", fields...)
	}
}

// LogReplayFailed 记录事件回放失败
func (l *EventLogger) LogReplayFailed(ctx context.Context, eventType string, err error, duration time.Duration) {
	fields := []zap.Field{
		zap.String("event_type", eventType),
		zap.String("result", "failed"),
		zap.Error(err),
		zap.Duration("duration", duration),
		zap.Int64("duration_ms", duration.Milliseconds()),
	}

	// 添加trace_id（如果可用）
	if traceID := getTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	l.logger.Error("Event replay failed", fields...)
}

// getTraceID 从context中获取trace_id
func getTraceID(ctx context.Context) string {
	// TODO: 从context中提取trace_id
	// 这里可以根据实际使用的tracing库（如OpenTelemetry）来实现
	// 例如：
	// if span := trace.SpanFromContext(ctx); span != nil {
	//     return span.SpanContext().TraceID().String()
	// }
	return ""
}

// GetOutputWriter 获取输出写入器
func GetOutputWriter(paths ...string) zapcore.WriteSyncer {
	var writers []zapcore.WriteSyncer
	for _, path := range paths {
		if path == "stdout" {
			writers = append(writers, zapcore.Lock(os.Stdout))
		} else if path == "stderr" {
			writers = append(writers, zapcore.Lock(os.Stderr))
		} else {
			// 文件输出
			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				// 如果打开失败，使用stdout
				writers = append(writers, zapcore.Lock(os.Stdout))
			} else {
				writers = append(writers, zapcore.Lock(file))
			}
		}
	}

	if len(writers) == 0 {
		return zapcore.AddSync(os.Stdout)
	}

	return zapcore.NewMultiWriteSyncer(writers...)
}

// StructuredEventLogger 结构化事件日志记录器
// 提供更丰富的日志功能
type StructuredEventLogger struct {
	logger *zap.Logger
}

// NewStructuredEventLogger 创建结构化事件日志记录器
func NewStructuredEventLogger(logger *zap.Logger) *StructuredEventLogger {
	return &StructuredEventLogger{
		logger: logger,
	}
}

// LogWithFields 使用字段记录日志
func (l *StructuredEventLogger) LogWithFields(level string, msg string, fields map[string]interface{}) {
	zapFields := make([]zap.Field, 0, len(fields))

	for k, v := range fields {
		switch val := v.(type) {
		case string:
			zapFields = append(zapFields, zap.String(k, val))
		case int:
			zapFields = append(zapFields, zap.Int(k, val))
		case int64:
			zapFields = append(zapFields, zap.Int64(k, val))
		case float64:
			zapFields = append(zapFields, zap.Float64(k, val))
		case bool:
			zapFields = append(zapFields, zap.Bool(k, val))
		case time.Time:
			zapFields = append(zapFields, zap.Time(k, val))
		case time.Duration:
			zapFields = append(zapFields, zap.Duration(k, val))
		case error:
			if val != nil {
				zapFields = append(zapFields, zap.Error(val))
			} else {
				zapFields = append(zapFields, zap.String(k, ""))
			}
		default:
			// 其他类型转为JSON
			jsonBytes, _ := json.Marshal(v)
			zapFields = append(zapFields, zap.String(k, string(jsonBytes)))
		}
	}

	switch level {
	case "debug":
		l.logger.Debug(msg, zapFields...)
	case "info":
		l.logger.Info(msg, zapFields...)
	case "warn":
		l.logger.Warn(msg, zapFields...)
	case "error":
		l.logger.Error(msg, zapFields...)
	default:
		l.logger.Info(msg, zapFields...)
	}
}

// LogEventLifecycle 记录事件生命周期
func (l *StructuredEventLogger) LogEventLifecycle(lifecycleStage string, event base.Event, fields map[string]interface{}) {
	allFields := make(map[string]interface{})
	for k, v := range fields {
		allFields[k] = v
	}

	// 添加事件基本信息
	allFields["lifecycle_stage"] = lifecycleStage
	allFields["event_type"] = event.GetEventType()
	allFields["event_source"] = event.GetSource()
	allFields["event_timestamp"] = event.GetTimestamp()

	level := "info"
	if lifecycleStage == "failed" {
		level = "error"
	}

	l.LogWithFields(level, fmt.Sprintf("Event %s", lifecycleStage), allFields)
}

// Flush 刷新日志缓冲
func (l *EventLogger) Flush() {
	_ = l.logger.Sync()
}

// Sync 同步日志
func (l *EventLogger) Sync() error {
	err := l.logger.Sync()
	if err == nil {
		return nil
	}

	// stdout/stderr 在部分环境不支持 sync；这类错误可安全忽略。
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "invalid argument") ||
		strings.Contains(msg, "inappropriate ioctl for device") ||
		strings.Contains(msg, "bad file descriptor") {
		return nil
	}

	return err
}
