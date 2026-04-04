# PKG 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **状态**: 进行中

---

## 1. 职责边界与依赖关系

### 1.1 职责定义

PKG 层是后端系统的**通用工具包层**，负责：

1. **通用工具封装**：封装常用的工具函数和组件
2. **基础设施抽象**：提供统一的基础设施访问接口
3. **错误处理**：统一的错误定义和处理机制
4. **日志记录**：结构化日志记录工具
5. **缓存访问**：Redis 等缓存客户端封装
6. **验证工具**：参数验证器和自定义验证规则
7. **分布式组件**：熔断器、分布式锁等

### 1.2 与上下层级的交互边界

```
┌─────────────────────────────────────────────────────────┐
│                   Service/Repository/API                │
│              (业务层调用工具)                            │
└─────────────────────────────────────────────────────────┘
                         │
                    调用工具函数
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                      PKG 层                             │
│  ┌─────────────────────────────────────────────────────┤
│  │ 职责：                                              │
│  │ - 错误处理 (errors/)                               │
│  │ - 日志记录 (logger/)                               │
│  │ - 缓存访问 (cache/)                                │
│  │ - 验证工具 (validator/)                            │
│  │ - 熔断器 (circuitbreaker/)                         │
│  │ - 分布式锁 (lock/)                                 │
│  │ - 指标收集 (metrics/)                              │
│  │ - 测试工具 (testutil/)                             │
│  │ - 响应工具 (response/)                             │
│  └─────────────────────────────────────────────────────┤
│  输出：                                                  │
│  │ - 通用工具函数                                      │
│  │ - 基础设施客户端                                    │
│  │ - 统一错误类型                                      │
└─────────────────────────────────────────────────────────┘
                         │
                    调用外部依赖
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                   外部依赖                              │
│        (zap, redis, validator, mongodb等)              │
└─────────────────────────────────────────────────────────┘
```

### 1.3 依赖关系

```go
// PKG 层允许的依赖
import (
    "github.com/redis/go-redis/v9"      // Redis客户端
    "github.com/go-playground/validator/v10" // 验证器
    "go.uber.org/zap"                   // 日志
    "go.mongodb.org/mongo-driver"       // MongoDB
    "github.com/gin-gonic/gin"          // Gin框架
)

// PKG 层禁止的依赖
import (
    "Qingyu_backend/service/xxx"        // ❌ 禁止依赖 Service
    "Qingyu_backend/repository/xxx"     // ❌ 禁止依赖 Repository
    "Qingyu_backend/api/xxx"            // ❌ 禁止依赖 API
)
```

---

## 2. 命名与代码规范

### 2.1 目录命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 功能模块 | 小写单数名词 | `logger`, `cache`, `errors` |
| 复合词 | 小写下划线 | `test_util` 或单数 `testutil` |

### 2.2 文件命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 主文件 | 包名.go | `logger.go`, `cache.go` |
| 接口定义 | `interface.go` 或在主文件 | `interface.go` |
| 实现文件 | `{功能}_impl.go` | `redis_client.go` |
| 错误定义 | `errors.go` | `errors.go` |
| 测试文件 | `{文件名}_test.go` | `logger_test.go` |

### 2.3 目录组织规范

```
pkg/
├── backup/                 # 备份工具
│   └── backup.go
├── cache/                  # 缓存工具
│   ├── redis_client.go     # Redis客户端接口和实现
│   ├── strategy.go         # 缓存策略
│   └── warmer.go           # 缓存预热
├── circuitbreaker/         # 熔断器
│   ├── circuit_breaker.go  # 熔断器实现
│   └── circuit_breaker_test.go
├── config/                 # 配置工具
│   └── config.go
├── cron/                   # 定时任务
│   └── backup_task.go
├── emailcode/              # 邮件验证码
│   └── manager.go
├── errors/                 # 错误处理
│   ├── unified_error.go    # 统一错误结构
│   ├── error_factory.go    # 错误工厂
│   ├── codes.go            # 错误码定义
│   └── layer_errors.go     # 层级错误
├── grpc/                   # gRPC客户端
│   └── client.go
├── lock/                   # 分布式锁
│   ├── distributed_lock.go
│   └── document_lock.go
├── logger/                 # 日志工具
│   └── logger.go
├── metrics/                # 指标收集
│   └── db_metrics.go
├── response/               # 响应工具
│   └── response.go
├── testutil/               # 测试工具
│   └── testutil.go
├── transaction/            # 事务工具
│   └── transaction.go
├── types/                  # 类型定义
│   └── types.go
├── utils/                  # 通用工具函数
│   └── utils.go
├── validator/              # 验证器
│   └── validator.go
└── websocket/              # WebSocket工具
    └── websocket.go
```

---

## 3. 设计模式与最佳实践

### 3.1 单例模式

用于全局唯一实例，如 Logger、Validator。

```go
// logger/logger.go
var (
    globalLogger   *Logger
    globalLoggerMu sync.RWMutex
    once           sync.Once
)

// Get 获取全局日志记录器
func Get() *Logger {
    globalLoggerMu.RLock()
    logger := globalLogger
    globalLoggerMu.RUnlock()
    if logger != nil {
        return logger
    }

    globalLoggerMu.Lock()
    defer globalLoggerMu.Unlock()

    if globalLogger != nil {
        return globalLogger
    }

    // 使用默认配置初始化
    if created, err := NewLogger(DefaultConfig()); err == nil {
        globalLogger = created
        return globalLogger
    }

    // 最后兜底
    nop := zap.NewNop()
    globalLogger = &Logger{
        Logger: nop,
        sugar:  nop.Sugar(),
    }
    return globalLogger
}

// validator/validator.go
var (
    once     sync.Once
    validate *validator.Validate
)

// GetValidator 获取全局验证器实例
func GetValidator() *validator.Validate {
    once.Do(func() {
        validate = validator.New()
        RegisterCustomValidators(validate)
    })
    return validate
}
```

### 3.2 接口抽象模式

所有外部依赖通过接口抽象，便于测试和替换。

```go
// cache/redis_client.go

// RedisClient Redis客户端接口
type RedisClient interface {
    // 基础操作
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, keys ...string) (int64, error)

    // 批量操作
    MGet(ctx context.Context, keys ...string) ([]interface{}, error)
    MSet(ctx context.Context, pairs ...interface{}) error

    // Hash操作
    HGet(ctx context.Context, key, field string) (string, error)
    HSet(ctx context.Context, key string, values ...interface{}) error
    HGetAll(ctx context.Context, key string) (map[string]string, error)
    HDel(ctx context.Context, key string, fields ...string) error

    // Set操作
    SAdd(ctx context.Context, key string, members ...interface{}) error
    SMembers(ctx context.Context, key string) ([]string, error)
    SRem(ctx context.Context, key string, members ...interface{}) error

    // 原子操作
    Incr(ctx context.Context, key string) (int64, error)
    Decr(ctx context.Context, key string) (int64, error)

    // 生命周期
    Ping(ctx context.Context) error
    Close() error
    GetClient() interface{}
}
```

### 3.3 Builder 模式

用于构建复杂对象，如 UnifiedError。

```go
// errors/unified_error.go

// ErrorBuilder 错误构建器
type ErrorBuilder struct {
    error *UnifiedError
}

// NewErrorBuilder 创建错误构建器
func NewErrorBuilder() *ErrorBuilder {
    return &ErrorBuilder{
        error: &UnifiedError{
            Timestamp: time.Now(),
            Metadata:  make(map[string]interface{}),
        },
    }
}

// WithCode 设置错误代码
func (b *ErrorBuilder) WithCode(code string) *ErrorBuilder {
    b.error.Code = code
    return b
}

// WithCategory 设置错误分类
func (b *ErrorBuilder) WithCategory(category ErrorCategory) *ErrorBuilder {
    b.error.Category = category
    return b
}

// WithMessage 设置错误消息
func (b *ErrorBuilder) WithMessage(message string) *ErrorBuilder {
    b.error.Message = message
    return b
}

// WithCause 设置原因错误
func (b *ErrorBuilder) WithCause(cause error) *ErrorBuilder {
    b.error.Cause = cause
    return b
}

// WithStack 设置堆栈信息
func (b *ErrorBuilder) WithStack() *ErrorBuilder {
    buf := make([]byte, 4096)
    n := runtime.Stack(buf, false)
    b.error.Stack = string(buf[:n])
    return b
}

// Build 构建错误
func (b *ErrorBuilder) Build() *UnifiedError {
    return b.error
}

// 使用示例
err := errors.NewErrorBuilder().
    WithCode("1001").
    WithCategory(errors.CategoryValidation).
    WithMessage("参数验证失败").
    WithDetails("用户名不能为空").
    WithCause(originalErr).
    Build()
```

### 3.4 状态机模式

用于有状态转换的组件，如熔断器。

```go
// circuitbreaker/circuit_breaker.go

// CircuitState 熔断器状态
type CircuitState int

const (
    StateClosed   CircuitState = 0 // 正常（关闭）
    StateOpen     CircuitState = 1 // 熔断（打开）
    StateHalfOpen CircuitState = 2 // 半开（尝试恢复）
)

// CircuitBreaker 熔断器
type CircuitBreaker struct {
    mu               sync.Mutex
    state            CircuitState
    failureCount     int
    failureThreshold int
    successCount     int
    successThreshold int
    lastFailureTime  time.Time
    timeout          time.Duration
}

// AllowRequest 判断是否允许请求
func (cb *CircuitBreaker) AllowRequest() bool {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    // 如果熔断超时，尝试进入半开状态
    if cb.state == StateOpen &&
        time.Since(cb.lastFailureTime) > cb.timeout {
        cb.setState(StateHalfOpen)
        return true
    }

    return cb.state != StateOpen
}

// RecordSuccess 记录成功
func (cb *CircuitBreaker) RecordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.failureCount = 0

    if cb.state == StateHalfOpen {
        cb.successCount++
        if cb.successCount >= cb.successThreshold {
            cb.setState(StateClosed)
        }
    }
}

// RecordFailure 记录失败
func (cb *CircuitBreaker) RecordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    cb.failureCount++
    cb.lastFailureTime = time.Now()

    if cb.state == StateHalfOpen {
        cb.setState(StateOpen)
    } else if cb.failureCount >= cb.failureThreshold {
        cb.setState(StateOpen)
    }
}
```

### 3.5 工厂函数模式

用于创建复杂对象，提供默认配置。

```go
// logger/logger.go

// Config 日志配置
type Config struct {
    Level       string `json:"level"`
    Format      string `json:"format"`
    Output      string `json:"output"`
    Filename    string `json:"filename"`
    Development bool   `json:"development"`
    StrictMode  bool   `json:"strict_mode"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
    return &Config{
        Level:       "info",
        Format:      "json",
        Output:      "stdout",
        Development: false,
        StrictMode:  false,
    }
}

// NewLogger 创建新的日志记录器
func NewLogger(config *Config) (*Logger, error) {
    if config == nil {
        config = DefaultConfig()
    }
    // ... 创建逻辑
}

// cache/redis_client.go

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg *config.RedisConfig) (RedisClient, error) {
    if cfg == nil {
        cfg = config.DefaultRedisConfig()
    }
    // ... 创建逻辑
}
```

### 3.6 错误包装模式

统一包装底层错误，提供一致的错误处理。

```go
// cache/redis_client.go

// Redis错误定义
var (
    ErrRedisNil         = errors.New("redis: nil returned")
    ErrKeyNotFound      = errors.New("redis: key not found")
    ErrConnectionFailed = errors.New("redis: connection failed")
    ErrTimeout          = errors.New("redis: operation timeout")
)

// wrapRedisError 包装Redis错误
func wrapRedisError(err error) error {
    if err == nil {
        return nil
    }
    if err == redis.Nil {
        return ErrRedisNil
    }
    if strings.Contains(err.Error(), "timeout") {
        return ErrTimeout
    }
    return err
}

// 使用
func (r *redisClientImpl) Get(ctx context.Context, key string) (string, error) {
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return "", wrapRedisError(err)
    }
    return val, nil
}
```

### 3.7 链式调用模式

提供流畅的 API 设计。

```go
// logger/logger.go

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

// 使用示例
logger.Get().
    WithRequest(requestID, "GET", "/api/users", clientIP).
    WithUser(userID).
    WithModule("user_service").
    Info("用户登录成功")
```

### 3.8 反模式警示

❌ **禁止事项**：

```go
// ❌ 禁止：在 pkg 中调用 Service
func SomeUtil() {
    service.GetUserByID(...)  // 禁止依赖业务层
}

// ❌ 禁止：pkg 中包含业务逻辑
func ValidateUser(user *User) error {
    if user.Age < 18 {  // 业务规则，应在 Service 层
        return errors.New("年龄不足")
    }
}

// ❌ 禁止：全局可变状态（非单例）
var GlobalCache = make(map[string]interface{})  // 应该通过接口访问

// ❌ 禁止：硬编码配置
func NewClient() *Client {
    return &Client{
        Timeout: 30 * time.Second,  // 应该从配置读取
    }
}
```

---

## 4. 模块规范

### 4.1 错误处理模块 (errors/)

```go
// 统一错误结构
type UnifiedError struct {
    ID         string                 `json:"id"`
    Code       string                 `json:"code"`
    Category   ErrorCategory          `json:"category"`
    Level      ErrorLevel             `json:"level"`
    Message    string                 `json:"message"`
    Details    string                 `json:"details,omitempty"`
    Cause      error                  `json:"-"`
    Stack      string                 `json:"stack,omitempty"`
    Metadata   map[string]interface{} `json:"metadata,omitempty"`
    HTTPStatus int                    `json:"http_status"`
    Retryable  bool                   `json:"retryable"`
}

// 错误分类
const (
    CategoryValidation ErrorCategory = "validation"
    CategoryBusiness   ErrorCategory = "business"
    CategorySystem     ErrorCategory = "system"
    CategoryExternal   ErrorCategory = "external"
    CategoryNetwork    ErrorCategory = "network"
    CategoryAuth       ErrorCategory = "auth"
    CategoryDatabase   ErrorCategory = "database"
    CategoryCache      ErrorCategory = "cache"
)

// 错误级别
const (
    LevelInfo     ErrorLevel = "info"
    LevelWarning  ErrorLevel = "warning"
    LevelError    ErrorLevel = "error"
    LevelCritical ErrorLevel = "critical"
)
```

### 4.2 日志模块 (logger/)

```go
// Logger 结构化日志记录器
type Logger struct {
    *zap.Logger
    sugar *zap.SugaredLogger
}

// 配置
type Config struct {
    Level       string // debug/info/warn/error
    Format      string // json/console
    Output      string // stdout/stderr/file/dual
    Filename    string // 日志文件路径
    Development bool   // 开发模式
    StrictMode  bool   // 严格模式
}

// 全局函数
func Get() *Logger
func Debug(msg string, fields ...zap.Field)
func Info(msg string, fields ...zap.Field)
func Warn(msg string, fields ...zap.Field)
func Error(msg string, fields ...zap.Field)
func Fatal(msg string, fields ...zap.Field)

// 链式调用
func (l *Logger) With(fields ...zap.Field) *Logger
func (l *Logger) WithRequest(requestID, method, path, ip string) *Logger
func (l *Logger) WithUser(userID string) *Logger
func (l *Logger) WithModule(module string) *Logger
func (l *Logger) WithError(err error) *Logger
```

### 4.3 缓存模块 (cache/)

```go
// RedisClient 接口
type RedisClient interface {
    // 基础操作
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, keys ...string) error
    Exists(ctx context.Context, keys ...string) (int64, error)

    // Hash/Set/原子操作...
    // 生命周期
    Ping(ctx context.Context) error
    Close() error
}

// 工厂函数
func NewRedisClient(cfg *config.RedisConfig) (RedisClient, error)

// 错误定义
var (
    ErrRedisNil         = errors.New("redis: nil returned")
    ErrKeyNotFound      = errors.New("redis: key not found")
    ErrConnectionFailed = errors.New("redis: connection failed")
    ErrTimeout          = errors.New("redis: operation timeout")
)
```

### 4.4 验证器模块 (validator/)

```go
// 全局验证器
func GetValidator() *validator.Validate
func ValidateStruct(s interface{}) error
func ValidateStructWithErrors(s interface{}) ValidationErrors

// 自定义验证器注册
func RegisterCustomValidators(v *validator.Validate) RegistrationStatus

// 注册状态
type RegistrationStatus struct {
    Total      int
    Success    int
    Failed     int
    FailedTags []string
    Errors     map[string]error
}
```

### 4.5 熔断器模块 (circuitbreaker/)

```go
// 状态
type CircuitState int
const (
    StateClosed   CircuitState = 0
    StateOpen     CircuitState = 1
    StateHalfOpen CircuitState = 2
)

// 熔断器
type CircuitBreaker struct {
    // ...
}

// 工厂函数
func NewCircuitBreaker(failureThreshold int, timeout time.Duration, successThreshold int) *CircuitBreaker

// 方法
func (cb *CircuitBreaker) AllowRequest() bool
func (cb *CircuitBreaker) RecordSuccess()
func (cb *CircuitBreaker) RecordFailure()
func (cb *CircuitBreaker) State() CircuitState
func (cb *CircuitBreaker) Stats() CircuitStats
```

---

## 5. 测试策略

### 5.1 单元测试编写指南

```go
// cache/redis_client_test.go
package cache

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRedisClient_SetAndGet(t *testing.T) {
    // 使用 miniredis 模拟 Redis
    // 或者使用 test container
    client, err := NewRedisClient(testConfig)
    require.NoError(t, err)
    defer client.Close()

    ctx := context.Background()

    // 测试 Set
    err = client.Set(ctx, "test_key", "test_value", time.Minute)
    assert.NoError(t, err)

    // 测试 Get
    val, err := client.Get(ctx, "test_key")
    assert.NoError(t, err)
    assert.Equal(t, "test_value", val)

    // 测试不存在
    _, err = client.Get(ctx, "nonexistent")
    assert.ErrorIs(t, err, ErrRedisNil)
}

// circuitbreaker/circuit_breaker_test.go
func TestCircuitBreaker_OpenOnFailures(t *testing.T) {
    cb := NewCircuitBreaker(3, time.Second, 2)

    // 初始状态应该是关闭的
    assert.Equal(t, StateClosed, cb.State())

    // 记录3次失败
    for i := 0; i < 3; i++ {
        cb.RecordFailure()
    }

    // 现在应该是打开的
    assert.Equal(t, StateOpen, cb.State())
    assert.False(t, cb.AllowRequest())
}
```

### 5.2 测试覆盖率要求

| 模块类型 | 覆盖率要求 |
|----------|------------|
| 错误处理 | ≥ 90% |
| 日志工具 | ≥ 80% |
| 缓存工具 | ≥ 85% |
| 熔断器 | ≥ 90% |
| 验证器 | ≥ 85% |

---

## 6. 完整代码示例

### 6.1 完整 Logger 模块示例

```go
// logger/logger.go
package logger

import (
    "os"
    "path"
    "sync"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    globalLogger   *Logger
    globalLoggerMu sync.RWMutex
    once           sync.Once
)

// Logger 结构化日志记录器
type Logger struct {
    *zap.Logger
    sugar *zap.SugaredLogger
}

// Config 日志配置
type Config struct {
    Level       string `json:"level"`
    Format      string `json:"format"`
    Output      string `json:"output"`
    Filename    string `json:"filename"`
    Development bool   `json:"development"`
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
        globalLoggerMu.Lock()
        globalLogger = logger
        globalLoggerMu.Unlock()
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
        logDir := path.Dir(config.Filename)
        if logDir != "" && logDir != "." {
            if err := os.MkdirAll(logDir, 0755); err != nil {
                return nil, err
            }
        }
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

    return &Logger{
        Logger: zapLogger,
        sugar:  zapLogger.Sugar(),
    }, nil
}

// Get 获取全局日志记录器
func Get() *Logger {
    globalLoggerMu.RLock()
    logger := globalLogger
    globalLoggerMu.RUnlock()
    if logger != nil {
        return logger
    }

    globalLoggerMu.Lock()
    defer globalLoggerMu.Unlock()

    if globalLogger != nil {
        return globalLogger
    }

    // 使用默认配置初始化
    if created, err := NewLogger(DefaultConfig()); err == nil {
        globalLogger = created
        return globalLogger
    }

    // 最后兜底
    nop := zap.NewNop()
    globalLogger = &Logger{
        Logger: nop,
        sugar:  nop.Sugar(),
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

// Info 信息日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
    l.Logger.Info(msg, fields...)
}

// Error 错误日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
    l.Logger.Error(msg, fields...)
}

// 全局便捷函数
func Info(msg string, fields ...zap.Field) {
    Get().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    Get().Error(msg, fields...)
}
```

---

## 7. 参考资料

- [PKG 层快速参考](../pkg/README.md)
- [Logger 文档](https://pkg.go.dev/go.uber.org/zap)
- [Validator 文档](https://github.com/go-playground/validator)
- [Redis 客户端文档](https://github.com/redis/go-redis)
- [Service 层设计说明](./layer-service.md)

---

*最后更新：2026-03-19*
