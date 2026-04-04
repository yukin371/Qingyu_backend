# PKG 层快速参考

## 职责

通用工具包层，负责通用工具封装、基础设施抽象、错误处理、日志记录、缓存访问、验证工具、分布式组件。

## 目录结构

```
pkg/
├── errors/                 # 错误处理
│   ├── unified_error.go    # 统一错误结构
│   ├── error_factory.go    # 错误工厂
│   └── codes.go            # 错误码定义
├── logger/                 # 日志工具
│   └── logger.go
├── cache/                  # 缓存工具
│   ├── redis_client.go     # Redis客户端
│   └── strategy.go         # 缓存策略
├── validator/              # 验证器
│   └── validator.go
├── circuitbreaker/         # 熔断器
│   └── circuit_breaker.go
├── lock/                   # 分布式锁
├── metrics/                # 指标收集
├── testutil/               # 测试工具
├── response/               # 响应工具
├── utils/                  # 通用工具函数
│   ├── ip.go               # IP处理
│   └── sanitizer.go        # 脱敏工具
└── ...
```

## 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 目录 | 小写单数名词 | `logger`, `cache`, `errors` |
| 主文件 | 包名.go | `logger.go` |
| 测试文件 | `{文件名}_test.go` | `logger_test.go` |

## 核心设计模式

### 1. 单例模式

```go
// logger
logger := logger.Get()
logger.Info("message")

// validator
v := validator.GetValidator()
err := v.Struct(data)
```

### 2. Builder 模式

```go
err := errors.NewErrorBuilder().
    WithCode("1001").
    WithCategory(errors.CategoryValidation).
    WithMessage("参数验证失败").
    WithCause(originalErr).
    Build()
```

### 3. 接口抽象模式

```go
// 定义接口
type RedisClient interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    // ...
}

// 通过工厂函数创建
client, err := cache.NewRedisClient(cfg)
```

### 4. 链式调用模式

```go
logger.Get().
    WithRequest(requestID, "GET", "/api/users", ip).
    WithUser(userID).
    WithModule("user_service").
    Info("用户登录成功")
```

## 常用模块

### 错误处理 (errors/)

```go
// 统一错误结构
type UnifiedError struct {
    Code       string
    Category   ErrorCategory
    Message    string
    HTTPStatus int
    // ...
}

// 使用 Builder 构建
err := errors.NewErrorBuilder().
    WithCode("1001").
    WithCategory(errors.CategoryValidation).
    WithMessage("用户名不能为空").
    Build()
```

### 日志 (logger/)

```go
// 全局日志
logger.Info("message", zap.String("key", "value"))

// 链式调用
logger.Get().
    WithModule("service").
    WithUser(userID).
    Info("操作完成")

// 配置
config := &logger.Config{
    Level:  "info",
    Format: "json",
    Output: "stdout",
}
logger.Init(config)
```

### 缓存 (cache/)

```go
// 创建客户端
client, err := cache.NewRedisClient(config.RedisConfig())

// 基础操作
client.Set(ctx, "key", "value", time.Minute)
val, err := client.Get(ctx, "key")

// 错误处理
if errors.Is(err, cache.ErrRedisNil) {
    // key 不存在
}
```

### 验证器 (validator/)

```go
// 验证结构体
err := validator.ValidateStruct(user)

// 获取友好错误
errs := validator.ValidateStructWithErrors(user)
for _, e := range errs {
    fmt.Printf("%s: %s\n", e.Field, e.Message)
}
```

### 熔断器 (circuitbreaker/)

```go
// 创建熔断器
cb := circuitbreaker.NewCircuitBreaker(5, time.Minute, 3)

// 使用
if cb.AllowRequest() {
    err := callExternalService()
    if err != nil {
        cb.RecordFailure()
    } else {
        cb.RecordSuccess()
    }
}
```

### 工具函数 (utils/)

```go
// IP 处理
ip := utils.GetClientIP(c)

// 脱敏
masked := utils.MaskEmail("user@example.com")  // us**@example.com
masked := utils.MaskPhone("13812345678")        // 138****5678
```

## 禁止事项

- ❌ 在 pkg 中调用 Service/Repository/API
- ❌ pkg 中包含业务逻辑
- ❌ 全局可变状态（非单例）
- ❌ 硬编码配置值

## 详见

完整设计文档: [docs/standards/layer-pkg.md](../docs/standards/layer-pkg.md)
