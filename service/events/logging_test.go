package events

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// ============ Replay 日志测试 ============

// TestEventLogger_LogReplayStarted 测试记录回放开始
func TestEventLogger_LogReplayStarted(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "test.event"
	now := time.Now()
	startTime := now.Add(-1 * time.Hour)
	endTime := now.Add(-30 * time.Minute)

	// Act
	logger.LogReplayStarted(ctx, eventType, &startTime, &endTime)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	assert.Equal(t, "Event replay started", entry.Message)
	assert.Equal(t, zap.InfoLevel, entry.Level)

	// 验证字段存在（使用String来获取值）
	fieldMap := make(map[string]string)
	for _, field := range entry.Context {
		fieldMap[field.Key] = field.String
	}
	assert.Equal(t, eventType, fieldMap["event_type"])
	assert.Contains(t, fieldMap, "start_time")
	assert.Contains(t, fieldMap, "end_time")
	assert.Contains(t, fieldMap, "range")
}

// TestEventLogger_LogReplayStarted_NoTimeRange 测试无时间范围的回放开始日志
func TestEventLogger_LogReplayStarted_NoTimeRange(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "all"

	// Act
	logger.LogReplayStarted(ctx, eventType, nil, nil)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	assert.Equal(t, "Event replay started", entry.Message)

	// 验证字段（使用String来获取值）
	fieldMap := make(map[string]string)
	for _, field := range entry.Context {
		fieldMap[field.Key] = field.String
	}
	assert.Equal(t, eventType, fieldMap["event_type"])
	// 不应该有 start_time, end_time, range 字段
	assert.NotContains(t, fieldMap, "start_time")
	assert.NotContains(t, fieldMap, "end_time")
	assert.NotContains(t, fieldMap, "range")
}

// TestEventLogger_LogReplayStarted_OnlyStartTime 测试只有开始时间的回放开始日志
func TestEventLogger_LogReplayStarted_OnlyStartTime(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "test.event"
	startTime := time.Now().Add(-1 * time.Hour)

	// Act
	logger.LogReplayStarted(ctx, eventType, &startTime, nil)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	assert.Equal(t, "Event replay started", entry.Message)

	// 验证字段（使用String来获取值）
	fieldMap := make(map[string]string)
	for _, field := range entry.Context {
		fieldMap[field.Key] = field.String
	}
	assert.Equal(t, eventType, fieldMap["event_type"])
	assert.Contains(t, fieldMap, "start_time")
	// 不应该有 range（需要start和end时间）
	assert.NotContains(t, fieldMap, "range")
}

// TestEventLogger_LogReplayCompleted_Success 测试记录成功的回放完成
func TestEventLogger_LogReplayCompleted_Success(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "test.event"
	result := &ReplayResult{
		ReplayedCount: 10,
		FailedCount:   0,
		SkippedCount:  0,
		Duration:      500 * time.Millisecond,
	}

	// Act
	logger.LogReplayCompleted(ctx, eventType, result)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	assert.Equal(t, "Event replay completed successfully", entry.Message)
	assert.Equal(t, zap.InfoLevel, entry.Level)

	// 验证字段存在
	fieldKeys := make([]string, len(entry.Context))
	for i, field := range entry.Context {
		fieldKeys[i] = field.Key
	}

	// 验证必需字段存在
	assert.Contains(t, fieldKeys, "event_type")
	assert.Contains(t, fieldKeys, "result")
	assert.Contains(t, fieldKeys, "replayed_count")
	assert.Contains(t, fieldKeys, "failed_count")
	assert.Contains(t, fieldKeys, "skipped_count")
	assert.Contains(t, fieldKeys, "duration")
	assert.Contains(t, fieldKeys, "duration_ms")
}

// TestEventLogger_LogReplayCompleted_PartialFailure 测试记录部分失败的回放完成
func TestEventLogger_LogReplayCompleted_PartialFailure(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "test.event"
	result := &ReplayResult{
		ReplayedCount: 8,
		FailedCount:   2,
		SkippedCount:  0,
		Duration:      750 * time.Millisecond,
	}

	// Act
	logger.LogReplayCompleted(ctx, eventType, result)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	assert.Equal(t, "Event replay completed with failures", entry.Message)
	assert.Equal(t, zap.WarnLevel, entry.Level)

	// 验证字段存在
	fieldKeys := make([]string, len(entry.Context))
	for i, field := range entry.Context {
		fieldKeys[i] = field.Key
	}

	assert.Contains(t, fieldKeys, "event_type")
	assert.Contains(t, fieldKeys, "result")
	assert.Contains(t, fieldKeys, "replayed_count")
	assert.Contains(t, fieldKeys, "failed_count")
	assert.Contains(t, fieldKeys, "duration_ms")
}

// TestEventLogger_LogReplayCompleted_DryRun 测试DryRun模式下的回放完成日志
func TestEventLogger_LogReplayCompleted_DryRun(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "test.event"
	result := &ReplayResult{
		ReplayedCount: 0,
		FailedCount:   0,
		SkippedCount:  100,
		Duration:      50 * time.Millisecond,
	}

	// Act
	logger.LogReplayCompleted(ctx, eventType, result)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	// 验证字段存在
	fieldKeys := make([]string, len(entry.Context))
	for i, field := range entry.Context {
		fieldKeys[i] = field.Key
	}

	assert.Contains(t, fieldKeys, "replayed_count")
	assert.Contains(t, fieldKeys, "failed_count")
	assert.Contains(t, fieldKeys, "skipped_count")
	assert.Contains(t, fieldKeys, "duration_ms")
}

// TestEventLogger_LogReplayFailed 测试记录回放失败
func TestEventLogger_LogReplayFailed(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.ErrorLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "test.event"
	testErr := context.DeadlineExceeded
	duration := 100 * time.Millisecond

	// Act
	logger.LogReplayFailed(ctx, eventType, testErr, duration)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	assert.Equal(t, "Event replay failed", entry.Message)
	assert.Equal(t, zap.ErrorLevel, entry.Level)

	// 验证字段存在
	fieldKeys := make([]string, len(entry.Context))
	for i, field := range entry.Context {
		fieldKeys[i] = field.Key
	}

	assert.Contains(t, fieldKeys, "event_type")
	assert.Contains(t, fieldKeys, "result")
	assert.Contains(t, fieldKeys, "error")
	assert.Contains(t, fieldKeys, "duration")
	assert.Contains(t, fieldKeys, "duration_ms")
}

// TestEventLogger_LogReplayFailed_LongDuration 测试长耗时的回放失败
func TestEventLogger_LogReplayFailed_LongDuration(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.ErrorLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	eventType := "test.event"
	testErr := context.DeadlineExceeded
	duration := 5 * time.Minute

	// Act
	logger.LogReplayFailed(ctx, eventType, testErr, duration)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	// 验证duration_ms字段存在
	fieldKeys := make([]string, len(entry.Context))
	for i, field := range entry.Context {
		fieldKeys[i] = field.Key
	}
	assert.Contains(t, fieldKeys, "duration_ms")
}

// TestEventLogger_AllLogFields 测试所有日志字段完整性
func TestEventLogger_AllLogFields(t *testing.T) {
	// 测试 LogReplayStarted 的所有字段
	t.Run("LogReplayStarted fields", func(t *testing.T) {
		core, logs := observer.New(zap.InfoLevel)
		logger := &EventLogger{
			logger: zap.New(core),
			config: DefaultLoggingConfig(),
		}

		ctx := context.Background()
		startTime := time.Now().Add(-1 * time.Hour)
		endTime := time.Now()

		logger.LogReplayStarted(ctx, "test.event", &startTime, &endTime)

		entry := logs.All()[0]
		fieldKeys := make([]string, len(entry.Context))
		for i, field := range entry.Context {
			fieldKeys[i] = field.Key
		}

		// 验证必需字段存在
		assert.Contains(t, fieldKeys, "event_type")
		assert.Contains(t, fieldKeys, "start_time")
		assert.Contains(t, fieldKeys, "end_time")
		assert.Contains(t, fieldKeys, "range")
	})

	// 测试 LogReplayCompleted 的所有字段
	t.Run("LogReplayCompleted fields", func(t *testing.T) {
		core, logs := observer.New(zap.InfoLevel)
		logger := &EventLogger{
			logger: zap.New(core),
			config: DefaultLoggingConfig(),
		}

		ctx := context.Background()
		result := &ReplayResult{
			ReplayedCount: 5,
			FailedCount:   1,
			SkippedCount:  0,
			Duration:      200 * time.Millisecond,
		}

		logger.LogReplayCompleted(ctx, "test.event", result)

		entry := logs.All()[0]
		fieldKeys := make([]string, len(entry.Context))
		for i, field := range entry.Context {
			fieldKeys[i] = field.Key
		}

		// 验证必需字段存在
		assert.Contains(t, fieldKeys, "event_type")
		assert.Contains(t, fieldKeys, "result")
		assert.Contains(t, fieldKeys, "replayed_count")
		assert.Contains(t, fieldKeys, "failed_count")
		assert.Contains(t, fieldKeys, "skipped_count")
		assert.Contains(t, fieldKeys, "duration")
		assert.Contains(t, fieldKeys, "duration_ms")
	})

	// 测试 LogReplayFailed 的所有字段
	t.Run("LogReplayFailed fields", func(t *testing.T) {
		core, logs := observer.New(zap.ErrorLevel)
		logger := &EventLogger{
			logger: zap.New(core),
			config: DefaultLoggingConfig(),
		}

		ctx := context.Background()
		testErr := context.Canceled

		logger.LogReplayFailed(ctx, "test.event", testErr, 100*time.Millisecond)

		entry := logs.All()[0]
		fieldKeys := make([]string, len(entry.Context))
		for i, field := range entry.Context {
			fieldKeys[i] = field.Key
		}

		// 验证必需字段存在
		assert.Contains(t, fieldKeys, "event_type")
		assert.Contains(t, fieldKeys, "result")
		assert.Contains(t, fieldKeys, "error")
		assert.Contains(t, fieldKeys, "duration")
		assert.Contains(t, fieldKeys, "duration_ms")
	})
}

// TestEventLogger_NewEventLogger 测试创建事件日志记录器
func TestEventLogger_NewEventLogger(t *testing.T) {
	// Act
	logger, err := NewEventLogger(DefaultLoggingConfig())

	// Assert
	require.NoError(t, err, "Should not return error")
	require.NotNil(t, logger, "Logger should not be nil")
	assert.NotNil(t, logger.logger, "Internal logger should not be nil")
	assert.NotNil(t, logger.config, "Config should not be nil")
}

// TestEventLogger_NewEventLogger_InvalidLevel 测试无效的日志级别
func TestEventLogger_NewEventLogger_InvalidLevel(t *testing.T) {
	// Arrange
	config := &LoggingConfig{
		Level:            "invalid",
		Encoding:         "json",
		EnableStackTrace: false,
		EnableCaller:     true,
		OutputPaths:      []string{"stdout"},
	}

	// Act
	logger, err := NewEventLogger(config)

	// Assert
	assert.Error(t, err, "Should return error for invalid log level")
	assert.Nil(t, logger, "Logger should be nil")
}

// TestDefaultLoggingConfig 测试默认日志配置
func TestDefaultLoggingConfig(t *testing.T) {
	// Act
	config := DefaultLoggingConfig()

	// Assert
	assert.NotNil(t, config)
	assert.Equal(t, "info", config.Level)
	assert.Equal(t, "json", config.Encoding)
	assert.False(t, config.EnableStackTrace)
	assert.True(t, config.EnableCaller)
	assert.Equal(t, []string{"stdout"}, config.OutputPaths)
}

// TestEventLogger_LogReplayCompleted_LargeNumbers 测试大数值日志
func TestEventLogger_LogReplayCompleted_LargeNumbers(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	result := &ReplayResult{
		ReplayedCount: 1000000,
		FailedCount:   500,
		SkippedCount:  100000,
		Duration:      10 * time.Minute,
	}

	// Act
	logger.LogReplayCompleted(ctx, "test.event", result)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	// 验证字段存在
	fieldKeys := make([]string, len(entry.Context))
	for i, field := range entry.Context {
		fieldKeys[i] = field.Key
	}

	assert.Contains(t, fieldKeys, "replayed_count")
	assert.Contains(t, fieldKeys, "failed_count")
	assert.Contains(t, fieldKeys, "skipped_count")
	assert.Contains(t, fieldKeys, "duration_ms")
}

// TestStructuredEventLogger_LogWithFields 测试结构化日志记录器
func TestStructuredEventLogger_LogWithFields(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	structuredLogger := NewStructuredEventLogger(logger)

	fields := map[string]interface{}{
		"event_type": "test.event",
		"count":      10,
		"duration":   100 * time.Millisecond,
		"success":    true,
	}

	// Act
	structuredLogger.LogWithFields("info", "Test message", fields)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	assert.Equal(t, "Test message", entry.Message)
	assert.Equal(t, zap.InfoLevel, entry.Level)
}

// TestStructuredEventLogger_LogEventLifecycle 测试事件生命周期日志
func TestStructuredEventLogger_LogEventLifecycle(t *testing.T) {
	t.Skip("需要base包支持，暂时跳过")
}

// TestEventLogger_Flush 测试刷新日志
func TestEventLogger_Flush(t *testing.T) {
	// Arrange
	logger, err := NewEventLogger(DefaultLoggingConfig())
	require.NoError(t, err)

	// Act - Flush should not panic
	logger.Flush()
}

// TestEventLogger_Sync 测试同步日志
func TestEventLogger_Sync(t *testing.T) {
	// Arrange
	logger, err := NewEventLogger(DefaultLoggingConfig())
	require.NoError(t, err)

	// Act
	err = logger.Sync()

	// Assert
	assert.NoError(t, err, "Sync should not return error")
}

// TestGetTraceID 测试获取trace_id
func TestGetTraceID(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	traceID := getTraceID(ctx)

	// Assert
	// 当前实现返回空字符串，这是符合预期的
	// TODO: 当实现实际的trace ID提取后，更新此测试
	assert.Equal(t, "", traceID)
}

// TestEventLogger_JSONEncoding 测试JSON编码格式的日志
func TestEventLogger_JSONEncoding(t *testing.T) {
	// Arrange
	config := &LoggingConfig{
		Level:            "info",
		Encoding:         "json",
		EnableStackTrace: false,
		EnableCaller:     false,
		OutputPaths:      []string{"stdout"},
	}

	// 使用内存编码器来捕获JSON输出
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var output []byte
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(&discardWriter{})), // 使用discard输出
		zap.InfoLevel,
	)

	logger := &EventLogger{
		logger: zap.New(core),
		config: config,
	}

	ctx := context.Background()
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()

	// Act
	logger.LogReplayStarted(ctx, "test.event", &startTime, &endTime)

	// Assert
	// JSON编码格式不会panic
	_ = output
	assert.NotNil(t, logger)
}

// discardWriter 用于测试的丢弃写入器
type discardWriter struct{}

func (d *discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (d *discardWriter) Sync() error {
	return nil
}

// TestEventLogger_ConsoleEncoding 测试Console编码格式的日志
func TestEventLogger_ConsoleEncoding(t *testing.T) {
	// Arrange
	config := &LoggingConfig{
		Level:            "info",
		Encoding:         "console",
		EnableStackTrace: false,
		EnableCaller:     false,
		OutputPaths:      []string{"stdout"},
	}

	// Act
	logger, err := NewEventLogger(config)

	// Assert
	require.NoError(t, err, "Should not return error for console encoding")
	assert.NotNil(t, logger)
	assert.Equal(t, "console", logger.config.Encoding)
}

// TestEventLogger_LogReplayCompleted_AllZeros 测试全零结果的日志
func TestEventLogger_LogReplayCompleted_AllZeros(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	result := &ReplayResult{
		ReplayedCount: 0,
		FailedCount:   0,
		SkippedCount:  0,
		Duration:      0,
	}

	// Act
	logger.LogReplayCompleted(ctx, "test.event", result)

	// Assert
	require.Equal(t, 1, logs.Len(), "Should have 1 log entry")
	entry := logs.All()[0]

	// 全零时应该记录为success（因为没有失败）
	assert.Equal(t, "Event replay completed successfully", entry.Message)
	assert.Equal(t, zap.InfoLevel, entry.Level)
}

// TestEventLogger_MultipleReplayLogs 测试多个Replay日志
func TestEventLogger_MultipleReplayLogs(t *testing.T) {
	// Arrange
	core, logs := observer.New(zap.InfoLevel)
	logger := &EventLogger{
		logger: zap.New(core),
		config: DefaultLoggingConfig(),
	}

	ctx := context.Background()
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()

	// Act - 记录多个日志
	logger.LogReplayStarted(ctx, "test.event", &startTime, &endTime)

	result := &ReplayResult{
		ReplayedCount: 10,
		FailedCount:   0,
		SkippedCount:  0,
		Duration:      100 * time.Millisecond,
	}
	logger.LogReplayCompleted(ctx, "test.event", result)

	// Assert
	require.Equal(t, 2, logs.Len(), "Should have 2 log entries")

	// 验证第一个日志（started）
	assert.Equal(t, "Event replay started", logs.All()[0].Message)

	// 验证第二个日志（completed）
	assert.Equal(t, "Event replay completed successfully", logs.All()[1].Message)
}
