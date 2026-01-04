package logger

import (
	"bytes"
	"os"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/assert"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "info", config.Level)
	assert.Equal(t, "json", config.Format)
	assert.Equal(t, "stdout", config.Output)
	assert.False(t, config.Development)
}

// TestNewLogger 测试创建日志记录器
func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "默认配置",
			config: DefaultConfig(),
			wantErr: false,
		},
		{
			name: "nil配置",
			config: nil,
			wantErr: false,
		},
		{
			name: "开发模式",
			config: &Config{
				Level:       "debug",
				Format:      "console",
				Output:      "stdout",
				Development: true,
			},
			wantErr: false,
		},
		{
			name: "JSON格式",
			config: &Config{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "无效日志级别",
			config: &Config{
				Level:  "invalid",
				Format: "json",
				Output: "stdout",
			},
			wantErr: false, // zap会使用默认级别
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				assert.NotNil(t, logger.Logger)
				assert.NotNil(t, logger.sugar)
			}
		})
	}
}

// TestLoggerLevels 测试各种日志级别
func TestLoggerLevels(t *testing.T) {
	// 创建观察核心以捕获日志
	observedZapCore, logs := observer.New(zapcore.InfoLevel)
	observedLogger := &Logger{
		Logger: zap.New(observedZapCore),
		sugar:  zap.New(observedZapCore).Sugar(),
	}

	// 测试各种日志级别
	observedLogger.Debug("debug message", zap.String("key", "value"))
	observedLogger.Info("info message", zap.String("key", "value"))
	observedLogger.Warn("warn message", zap.String("key", "value"))
	observedLogger.Error("error message", zap.String("key", "value"))

	// 验证日志数量 (debug级别低于InfoLevel，不会被记录)
	assert.Equal(t, 3, logs.Len()) // info, warn, error

	// 验证日志内容
	assert.Equal(t, "info message", logs.All()[0].Message)
	assert.Equal(t, "warn message", logs.All()[1].Message)
	assert.Equal(t, "error message", logs.All()[2].Message)
}

// TestLoggerWithFields 测试带结构化字段的日志
func TestLoggerWithFields(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	// 测试 With
	logger.With(
		zap.String("service", "test-service"),
		zap.Int("port", 8080),
	).Info("started with fields")

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "started with fields", entry.Message)
	assert.Equal(t, "test-service", entry.ContextMap()["service"])
	assert.Equal(t, int64(8080), entry.ContextMap()["port"])
}

// TestLoggerWithRequest 测试带请求信息的日志
func TestLoggerWithRequest(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	requestLogger := logger.WithRequest("req-123", "GET", "/api/test", "127.0.0.1")
	requestLogger.Info("request received")

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "req-123", entry.ContextMap()["request_id"])
	assert.Equal(t, "GET", entry.ContextMap()["method"])
	assert.Equal(t, "/api/test", entry.ContextMap()["path"])
	assert.Equal(t, "127.0.0.1", entry.ContextMap()["ip"])
}

// TestLoggerWithUser 测试带用户信息的日志
func TestLoggerWithUser(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	userLogger := logger.WithUser("user-456")
	userLogger.Info("user action")

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "user-456", entry.ContextMap()["user_id"])
}

// TestLoggerWithModule 测试带模块信息的日志
func TestLoggerWithModule(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	moduleLogger := logger.WithModule("auth-service")
	moduleLogger.Info("module initialized")

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "auth-service", entry.ContextMap()["module"])
}

// TestLoggerWithError 测试带错误信息的日志
func TestLoggerWithError(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	testErr := assert.AnError
	errLogger := logger.WithError(testErr)
	errLogger.Error("operation failed")

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Contains(t, entry.ContextMap()["error"].(string), testErr.Error())
}

// TestLoggerFormatted 测试格式化日志
func TestLoggerFormatted(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	// 测试格式化日志
	logger.Debugf("debug: %s", "value")
	logger.Infof("info: %d", 42)
	logger.Warnf("warn: %v", true)
	logger.Errorf("error: %f", 3.14)

	// Debug级别不会记录（因为核心级别是Info）
	assert.Equal(t, 3, logs.Len())
	assert.Equal(t, "info: 42", logs.All()[0].Message)
	assert.Equal(t, "warn: true", logs.All()[1].Message)
	assert.Contains(t, logs.All()[2].Message, "error: 3.14")
}

// TestInitGlobalLogger 测试全局日志初始化
func TestInitGlobalLogger(t *testing.T) {
	// 保存旧的全局logger
	oldLogger := globalLogger
	defer func() {
		globalLogger = oldLogger
	}()

	// 重置once以便可以重新初始化
	once = *new(sync.Once)

	config := &Config{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}

	err := Init(config)
	assert.NoError(t, err)
	assert.NotNil(t, globalLogger)
}

// TestGetLogger 测试获取全局日志记录器
func TestGetLogger(t *testing.T) {
	// 保存旧的全局logger
	oldLogger := globalLogger
	defer func() {
		globalLogger = oldLogger
		once = *new(sync.Once)
	}()

	// 重置
	globalLogger = nil
	once = *new(sync.Once)

	logger := Get()
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.Logger)
	assert.NotNil(t, logger.sugar)

	// 多次调用应该返回同一个实例
	logger2 := Get()
	assert.Same(t, logger.Logger, logger2.Logger)
}

// TestGlobalConvenienceFunctions 测试全局便捷函数
func TestGlobalConvenienceFunctions(t *testing.T) {
	// 保存旧的全局logger
	oldLogger := globalLogger
	defer func() {
		globalLogger = oldLogger
		once = *new(sync.Once)
	}()

	// 重置并初始化
	globalLogger = nil
	once = *new(sync.Once)
	_ = Init(&Config{Level: "debug", Format: "json", Output: "stdout"})

	// 测试全局函数（不应该panic）
	assert.NotPanics(t, func() {
		Debug("debug message")
		Info("info message")
		Warn("warn message")
		Error("error message")
	})

	// 测试全局logger的格式化函数
	logger := Get()
	assert.NotPanics(t, func() {
		logger.Debugf("debug: %s", "test")
		logger.Infof("info: %d", 123)
		logger.Warnf("warn: %v", true)
		logger.Errorf("error: %f", 3.14)
	})
}

// TestGlobalWithFunctions 测试全局With函数
func TestGlobalWithFunctions(t *testing.T) {
	// 保存旧的全局logger
	oldLogger := globalLogger
	defer func() {
		globalLogger = oldLogger
		once = *new(sync.Once)
	}()

	// 重置并初始化
	globalLogger = nil
	once = *new(sync.Once)
	_ = Init(&Config{Level: "info", Format: "json", Output: "stdout"})

	// 测试全局With函数
	assert.NotPanics(t, func() {
		With(zap.String("key", "value"))
		WithRequest("req-123", "GET", "/api/test", "127.0.0.1")
		WithUser("user-456")
		WithModule("test-module")
	})
}

// TestLoggerOutputToFile 测试日志输出到文件
func TestLoggerOutputToFile(t *testing.T) {
	// 创建临时文件
	tempFile := "test_temp.log"
	defer os.Remove(tempFile)

	config := &Config{
		Level:    "info",
		Format:   "json",
		Output:   "file",
		Filename: tempFile,
	}

	logger, err := NewLogger(config)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	logger.Info("test message to file")
	logger.Sync()

	// 读取文件并验证内容
	content, err := os.ReadFile(tempFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test message to file")
	assert.Contains(t, string(content), `"level":"info"`)
}

// TestLoggerConsoleFormat 测试控制台格式
func TestLoggerConsoleFormat(t *testing.T) {
	config := &Config{
		Level:       "info",
		Format:      "console",
		Output:      "stdout",
		Development: true,
	}

	logger, err := NewLogger(config)
	assert.NoError(t, err)

	// 只验证logger能正常工作，不捕获输出
	assert.NotPanics(t, func() {
		logger.Info("console message")
		logger.Sync()
	})
}

// TestLoggerDevelopmentMode 测试开发模式
func TestLoggerDevelopmentMode(t *testing.T) {
	config := &Config{
		Level:       "debug",
		Format:      "console",
		Output:      "stdout",
		Development: true,
	}

	logger, err := NewLogger(config)
	assert.NoError(t, err)

	// 开发模式应该使用sugar的格式化功能
	assert.NotNil(t, logger.sugar)
	assert.NotNil(t, logger.Logger)
}

// TestLoggerMultipleOutputs 测试多个日志输出
func TestLoggerMultipleOutputs(t *testing.T) {
	t.Skip("需要实现多输出支持后才能测试")
	// TODO: 当支持多个输出路径时测试
}

// BenchmarkLoggerLogging 性能测试 - 结构化日志
func BenchmarkLoggerLogging(b *testing.B) {
	observedCore, _ := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message",
			zap.String("user_id", "12345"),
			zap.String("action", "login"),
			zap.Int("attempt", i),
			zap.Duration("latency", time.Millisecond),
		)
	}
}

// BenchmarkLoggerFormatted 性能测试 - 格式化日志
func BenchmarkLoggerFormatted(b *testing.B) {
	observedCore, _ := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("User %s performed action %s (attempt %d)", "12345", "login", i)
	}
}

// BenchmarkLoggerWithFields 性能测试 - With字段
func BenchmarkLoggerWithFields(b *testing.B) {
	observedCore, _ := observer.New(zapcore.InfoLevel)
	baseLogger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}
	logger := baseLogger.With(
		zap.String("service", "benchmark"),
		zap.String("version", "1.0.0"),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("processing request", zap.Int("request_id", i))
	}
}

// TestLoggerChaining 测试Logger链式调用
func TestLoggerChaining(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	// 测试链式调用
	logger.
		WithModule("test").
		WithUser("user-123").
		WithRequest("req-456", "POST", "/api/data", "10.0.0.1").
		Info("chained log entry")

	assert.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "test", entry.ContextMap()["module"])
	assert.Equal(t, "user-123", entry.ContextMap()["user_id"])
	assert.Equal(t, "req-456", entry.ContextMap()["request_id"])
}

// TestLoggerLevelsFiltering 测试日志级别过滤
func TestLoggerLevelsFiltering(t *testing.T) {
	tests := []struct {
		name         string
		configLevel  string
		expectedLogs int // 应该记录的日志数量
	}{
		{"Debug级别", "debug", 5},
		{"Info级别", "info", 4},
		{"Warn级别", "warn", 3},
		{"Error级别", "error", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var level zapcore.Level
			level.UnmarshalText([]byte(tt.configLevel))
			observedCore, logs := observer.New(level)
			logger := &Logger{
				Logger: zap.New(observedCore),
				sugar:  zap.New(observedCore).Sugar(),
			}

			logger.Debug("debug")
			logger.Info("info")
			logger.Warn("warn")
			logger.Error("error")
			logger.DPanic("dpanic")

			// 验证日志数量
			assert.LessOrEqual(t, logs.Len(), tt.expectedLogs)
		})
	}
}

// TestLoggerJSONOutput 测试JSON格式输出
func TestLoggerJSONOutput(t *testing.T) {
	var buf bytes.Buffer

	// 创建写入core
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)

	logger := &Logger{
		Logger: zap.New(core),
		sugar:  zap.New(core).Sugar(),
	}

	logger.Info("json message", zap.String("key", "value"))
	logger.Sync()

	output := buf.String()

	// 验证JSON格式
	assert.Contains(t, output, `"level":"info"`)
	assert.Contains(t, output, `"message":"json message"`)
	assert.Contains(t, output, `"key":"value"`)
	assert.Contains(t, output, `"timestamp"`)
}

// TestLoggerConcurrency 测试并发安全性
func TestLoggerConcurrency(t *testing.T) {
	observedCore, logs := observer.New(zapcore.InfoLevel)
	logger := &Logger{
		Logger: zap.New(observedCore),
		sugar:  zap.New(observedCore).Sugar(),
	}

	// 并发写入日志
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(n int) {
			logger.Info("concurrent message", zap.Int("goroutine", n))
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 100; i++ {
		<-done
	}

	// 验证所有日志都被记录
	assert.Equal(t, 100, logs.Len())
}
