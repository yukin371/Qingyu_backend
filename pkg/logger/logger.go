package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *Logger
	once         sync.Once
)

// Logger 结构化日志记录器
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// Config 日志配置
type Config struct {
	Level       string `json:"level"`        // debug/info/warn/error/dpanic/panic/fatal
	Format      string `json:"format"`       // json/console
	Output      string `json:"output"`       // stdout/stderr/file
	Filename    string `json:"filename"`     // 日志文件路径
	MaxSize     int    `json:"maxSize"`      // 单个日志文件最大大小(MB)
	MaxBackups  int    `json:"maxBackups"`   // 保留的旧日志文件最大数量
	MaxAge      int    `json:"maxAge"`       // 保留旧日志文件的最大天数
	Compress    bool   `json:"compress"`     // 是否压缩旧日志文件
	Development bool   `json:"development"`  // 是否开发模式
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		Development: false,
	}
}

// Init 初始化全局日志记录器
func Init(config *Config) error {
	var initErr error
	once.Do(func() {
		logger, err := NewLogger(config)
		if err != nil {
			initErr = err
			return
		}
		globalLogger = logger
	})
	return initErr
}

// NewLogger 创建新的日志记录器
func NewLogger(config *Config) (*Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 解析日志级别
	level := zapcore.InfoLevel
	if config.Development {
		level = zapcore.DebugLevel
	}
	if err := level.UnmarshalText([]byte(config.Level)); err == nil {
		// 使用配置的级别
	}

	// 编码器配置
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

	if config.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// 选择编码器
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 输出
	var writeSyncer zapcore.WriteSyncer
	switch config.Output {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	case "file":
		if config.Filename == "" {
			config.Filename = "logs/app.log"
		}
		// TODO: 支持日志轮转
		file, err := os.OpenFile(config.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建Core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建Logger
	zapLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	if config.Development {
		zapLogger = zapLogger.WithOptions(zap.Development())
	}

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}, nil
}

// Get 获取全局日志记录器
func Get() *Logger {
	if globalLogger == nil {
		// 使用默认配置初始化
		_ = Init(DefaultConfig())
	}
	return globalLogger
}

// With 添加结构化字段
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		sugar:  l.sugar,
	}
}

// WithRequest 添加请求相关字段
func (l *Logger) WithRequest(requestID, method, path, ip string) *Logger {
	return l.With(
		zap.String("request_id", requestID),
		zap.String("method", method),
		zap.String("path", path),
		zap.String("ip", ip),
	)
}

// WithUser 添加用户相关字段
func (l *Logger) WithUser(userID string) *Logger {
	return l.With(zap.String("user_id", userID))
}

// WithModule 添加模块字段
func (l *Logger) WithModule(module string) *Logger {
	return l.With(zap.String("module", module))
}

// WithError 添加错误字段
func (l *Logger) WithError(err error) *Logger {
	return l.With(zap.Error(err))
}

// Debug 调试日志
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Info 信息日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Warn 警告日志
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error 错误日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Panic Panic日志
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.Logger.Panic(msg, fields...)
}

// DPanic 开发模式Panic日志
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.Logger.DPanic(msg, fields...)
}

// Debugf 格式化调试日志
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Infof 格式化信息日志
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warnf 格式化警告日志
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Fatalf 格式化致命错误日志
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

// Panicf 格式化Panic日志
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

// Sync 刷新缓冲区
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// 便捷函数 - 使用全局Logger

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// With 添加字段
func With(fields ...zap.Field) *Logger {
	return Get().With(fields...)
}

// WithRequest 添加请求字段
func WithRequest(requestID, method, path, ip string) *Logger {
	return Get().WithRequest(requestID, method, path, ip)
}

// WithUser 添加用户字段
func WithUser(userID string) *Logger {
	return Get().WithUser(userID)
}

// WithModule 添加模块字段
func WithModule(module string) *Logger {
	return Get().WithModule(module)
}
