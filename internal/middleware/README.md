# 中间件系统使用文档

> **版本**: v1.0
> **最后更新**: 2026-02-26
> **状态**: 稳定

---

## 目录

- [概述](#概述)
- [快速开始](#快速开始)
- [内置中间件](#内置中间件)
- [认证授权](#认证授权)
- [限流](#限流)
- [配置管理](#配置管理)
- [自定义中间件](#自定义中间件)
- [API参考](#api参考)
- [测试指南](#测试指南)
- [迁移指南](#迁移指南)
- [最佳实践](#最佳实践)

---

## 概述

### 中间件系统介绍

青羽后端中间件系统提供统一的请求处理管道，包含认证、授权、限流、日志、监控等横切关注点。系统采用插件化架构，支持灵活的配置和扩展。

### 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP 请求                             │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                  中间件管道                               │
│  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐       │
│  │Request │  │ Logger │  │  CORS  │  │RateLimit│       │
│  │   ID   │  │        │  │        │  │         │       │
│  └───┬────┘  └───┬────┘  └───┬────┘  └───┬────┘       │
│      │           │           │           │              │
│  ┌───▼────┐  ┌───▼────┐  ┌───▼────┐  ┌───▼────┐       │
│  │Recovery│  │  Auth  │  │Permission│  │Validation│    │
│  └───┬────┘  └───┬────┘  └───┬────┘  └───┬────┘       │
│      │           │           │           │              │
└──────┼───────────┼───────────┼───────────┼──────────────┘
       │           │           │           │
       ▼           ▼           ▼           ▼
┌─────────────────────────────────────────────────────────┐
│                  业务处理器                               │
└─────────────────────────────────────────────────────────┘
```

### 核心特性

- **模块化架构**: 插件式中间件设计，易于扩展
- **优先级控制**: 灵活的中间件执行顺序管理
- **配置驱动**: YAML配置文件，支持热重载
- **类型安全**: 强类型Go接口，编译时检查
- **性能优化**: 最小化性能影响，支持对象池
- **监控集成**: 完整的指标收集和日志记录

### 目录结构

```
internal/middleware/
├── auth/              # 认证授权中间件
│   ├── jwt.go
│   ├── jwt_manager.go
│   ├── rbac_checker.go
│   ├── permission.go
│   └── blacklist.go
├── builtin/           # 内置中间件
│   ├── logger.go
│   ├── cors.go
│   ├── recovery.go
│   ├── request_id.go
│   ├── security.go
│   ├── error_handler.go
│   └── compression.go
├── core/              # 核心接口
│   ├── middleware.go
│   ├── manager.go
│   ├── registry.go
│   ├── initializer.go
│   └── context.go
├── ratelimit/         # 限流中间件
│   ├── rate_limit.go
│   ├── token_bucket.go
│   ├── sliding_window.go
│   ├── redis_limiter.go
│   └── config.go
├── monitoring/        # 监控中间件
│   └── metrics.go
├── validation/        # 验证中间件
│   └── validation.go
├── config.go          # 配置定义
├── version_routing.go # 版本路由
└── deprecation.go     # 废弃警告
```

---

## 快速开始

### 基本使用示例

#### 1. 创建管理器并注册中间件

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "Qingyu_backend/internal/middleware/core"
    "Qingyu_backend/internal/middleware/builtin"
    "Qingyu_backend/internal/middleware/auth"
    "Qingyu_backend/internal/middleware/ratelimit"
)

func main() {
    // 创建Gin引擎
    r := gin.New()

    // 创建logger
    logger, _ := zap.NewProduction()

    // 创建中间件管理器
    manager := core.NewManager(logger)

    // 注册中间件
    requestID := builtin.NewRequestIDMiddleware()
    loggerMW := builtin.NewLoggerMiddleware(logger)
    recovery := builtin.NewRecoveryMiddleware(logger)

    manager.Register(requestID)
    manager.Register(loggerMW)
    manager.Register(recovery)

    // 应用中间件到路由
    manager.ApplyToRouter(r)

    // 注册路由
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    // 启动服务
    r.Run(":8080")
}
```

#### 2. 使用认证中间件

```go
// 创建JWT管理器
jwtManager, err := auth.NewJWTManager(
    "your-secret-key",
    2*time.Hour,  // Access Token过期时间
    7*24*time.Hour, // Refresh Token过期时间
)
if err != nil {
    log.Fatal(err)
}

// 创建黑名单（可选）
blacklist := auth.NewMemoryBlacklist()

// 创建认证中间件
authMiddleware := auth.NewJWTAuthMiddleware(
    jwtManager,
    blacklist,
    logger,
)

// 注册认证中间件
manager.Register(authMiddleware)

// 创建需要认证的路由组
authenticated := r.Group("/api/v1")
authenticated.Use(authMiddleware.Handler())
{
    authenticated.GET("/profile", getProfile)
    authenticated.PUT("/profile", updateProfile)
}
```

#### 3. 使用限流中间件

```go
// 创建限流配置
config := &ratelimit.RateLimitConfig{
    Enabled: true,
    Strategy: "token_bucket",
    Rate:     100,  // 每秒100个请求
    Burst:    200,  // 突发容量200
    KeyFunc:  "client_ip", // 按客户端IP限流
    SkipPaths: []string{"/health", "/metrics"},
}

// 创建限流中间件
rateLimitMiddleware, err := ratelimit.NewRateLimitMiddleware(config, logger)
if err != nil {
    log.Fatal(err)
}

// 注册限流中间件
manager.Register(rateLimitMiddleware)
```

### 最小配置

#### 最小化中间件配置

```yaml
# config/middleware.yaml
middleware:
  request_id:
    enabled: true
    header_name: "X-Request-ID"

  logger:
    enabled: true
    skip_paths:
      - /health
      - /metrics

  recovery:
    enabled: true
    disable_print: true
```

#### 使用最小配置启动

```go
import "Qingyu_backend/internal/middleware"

func main() {
    r := gin.New()
    logger, _ := zap.NewProduction()

    // 使用默认配置初始化中间件
    middlewares, err := middleware.InitializeDefault(logger)
    if err != nil {
        log.Fatal(err)
    }

    // 应用中间件
    for _, mw := range middlewares {
        r.Use(mw.Handler())
    }

    r.Run(":8080")
}
```

### 运行第一个中间件

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "Qingyu_backend/internal/middleware/builtin"
)

func main() {
    r := gin.New()

    logger, _ := zap.NewDevelopment()

    // 只启用RequestID中间件
    requestID := builtin.NewRequestIDMiddleware()
    r.Use(requestID.Handler())

    r.GET("/test", func(c *gin.Context) {
        // 从上下文获取请求ID
        requestID := c.GetString("request_id")
        c.JSON(200, gin.H{
            "request_id": requestID,
            "message": "Hello, World!",
        })
    })

    r.Run(":8080")
}
```

测试：

```bash
$ curl http://localhost:8080/test
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Hello, World!"
}
```

---

## 内置中间件

### 中间件优先级说明

中间件按照优先级顺序执行，优先级数值越小越先执行：

| 优先级 | 层级 | 中间件 | 说明 |
|--------|------|--------|------|
| 1 | 基础设施 | RequestID | 生成唯一请求ID |
| 2 | 基础设施 | Security | 设置安全响应头 |
| 3 | 基础设施 | CORS | 处理跨域请求 |
| 5 | 基础设施 | Recovery | 恢复panic |
| 6 | 监控 | Metrics | 收集性能指标 |
| 7 | 监控 | Logger | 记录请求日志 |
| 8 | 安全 | RateLimit | 请求限流 |
| 9 | 认证 | JWT | JWT认证 |
| 10 | 授权 | Permission | 权限检查 |
| 11 | 业务 | Validation | 请求验证 |
| 12 | 业务 | Compression | 响应压缩 |

### RequestID 中间件

为每个请求生成唯一标识符。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| header_name | string | "X-Request-ID" | 请求ID的HTTP头名称 |
| force_gen | bool | false | 是否强制生成新ID |

**使用示例：**

```go
requestID := builtin.NewRequestIDMiddleware()
r.Use(requestID.Handler())

// 在处理器中使用
func handler(c *gin.Context) {
    requestID := c.GetString("request_id")
    // 使用request_id进行追踪
}
```

### Logger 中间件

记录请求和响应日志。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| skip_paths | []string | [] | 跳过日志记录的路径 |
| enable_request_id | bool | true | 是否记录请求ID |
| enable_request_body | bool | true | 是否记录请求体 |
| enable_response_body | bool | false | 是否记录响应体 |
| slow_request_threshold | int | 3000 | 慢请求阈值（毫秒） |
| max_body_size | int | 2048 | 请求体最大记录大小（字节） |
| redact_keys | []string | ["authorization", "password"] | 需要脱敏的键名 |

**使用示例：**

```go
logger := builtin.NewLoggerMiddleware(zapLogger)

// 自定义配置
config := &builtin.LoggerConfig{
    SkipPaths: []string{"/health", "/metrics"},
    SlowRequestThreshold: 5000,
    RedactKeys: []string{"password", "token"},
}
middleware := builtin.NewLoggerMiddlewareWithConfig(config, zapLogger)
```

### Recovery 中间件

恢复panic并记录错误日志。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| stack_size | int | 4096 | 堆栈大小 |
| disable_print | bool | true | 是否禁用打印 |

**使用示例：**

```go
recovery := builtin.NewRecoveryMiddleware(zapLogger)
r.Use(recovery.Handler())
```

### CORS 中间件

处理跨域资源共享请求。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| allow_origins | []string | ["*"] | 允许的来源 |
| allow_methods | []string | ["GET", "POST", "PUT", "DELETE", "OPTIONS"] | 允许的HTTP方法 |
| allow_headers | []string | ["*"] | 允许的请求头 |
| expose_headers | []string | [] | 暴露的响应头 |
| allow_credentials | bool | true | 是否允许携带凭证 |
| max_age | duration | 12h | 预检请求缓存时间 |

**使用示例：**

```go
config := &builtin.CORSConfig{
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
cors := builtin.NewCORSMiddleware(config)
r.Use(cors.Handler())
```

### Security 中间件

设置安全相关的HTTP响应头。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| enable_x_frame_options | bool | true | 是否启用X-Frame-Options |
| x_frame_options | string | "DENY" | X-Frame-Options值 |
| enable_hsts | bool | true | 是否启用HSTS |
| enable_csp | bool | true | 是否启用CSP |

**使用示例：**

```go
security := builtin.NewSecurityMiddleware()
r.Use(security.Handler())
```

### ErrorHandler 中间件

统一错误处理和响应格式化。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| enable_stack_trace | bool | true | 是否启用堆栈追踪 |
| enable_logging | bool | true | 是否记录日志 |
| log_level | string | "error" | 日志级别 |

**使用示例：**

```go
errorHandler := builtin.NewErrorHandlerMiddleware(zapLogger)
r.Use(errorHandler.Handler())
```

### Compression 中间件

压缩响应内容。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| enabled | bool | true | 是否启用压缩 |
| level | int | 5 | 压缩级别（1-9） |
| types | []string | ["application/json", "text/html"] | 压缩的内容类型 |

**使用示例：**

```go
config := &builtin.CompressionConfig{
    Enabled: true,
    Level:   5,
    Types:   []string{"application/json", "text/html"},
}
compression := builtin.NewCompressionMiddleware(config)
r.Use(compression.Handler())
```

---

## 认证授权

### JWT 认证中间件

验证JWT Token并提取用户信息。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| enabled | bool | true | 是否启用认证 |
| secret | string | "" | JWT密钥（必须设置） |
| access_expiration | duration | 2h | Access Token过期时间 |
| refresh_expiration | duration | 7d | Refresh Token过期时间 |
| issuer | string | "qingyu" | 签发者 |
| token_header | string | "Authorization" | Token所在HTTP头 |
| token_prefix | string | "Bearer" | Token前缀 |
| skip_paths | []string | ["/health", "/metrics"] | 跳过认证的路径 |

**使用示例：**

```go
// 创建JWT管理器
jwtManager, err := auth.NewJWTManager(
    "your-secret-key",
    2*time.Hour,
    7*24*time.Hour,
)

// 创建认证中间件
authMiddleware := auth.NewJWTAuthMiddleware(
    jwtManager,
    nil, // 黑名单（可选）
    zapLogger,
)

// 配置跳过路径
config := &auth.JWTConfig{
    SkipPaths: []string{
        "/api/v1/auth/login",
        "/api/v1/auth/register",
        "/health",
    },
}
authMiddleware.LoadConfig(map[string]interface{}{
    "skip_paths": config.SkipPaths,
})

// 使用中间件
r.Use(authMiddleware.Handler())
```

**在处理器中获取用户信息：**

```go
func getProfile(c *gin.Context) {
    userID := c.GetString("user_id")
    username := c.GetString("username")
    roles := c.Get("roles")

    c.JSON(200, gin.H{
        "user_id":  userID,
        "username": username,
        "roles":    roles,
    })
}
```

### RBAC 权限检查中间件

基于角色的访问控制。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| enabled | bool | true | 是否启用权限检查 |
| strategy | string | "rbac" | 权限策略 |
| config_path | string | "" | 权限规则配置文件路径 |
| skip_paths | []string | [] | 跳过权限检查的路径 |

**使用示例：**

```go
// 创建权限检查器
permissionChecker := auth.NewRBACChecker()

// 创建权限中间件
permissionMiddleware := auth.NewPermissionMiddleware(
    permissionChecker,
    zapLogger,
)

// 配置权限规则
rules := &auth.RBACRules{
    "/api/v1/admin": []string{"admin"},
    "/api/v1/users": []string{"admin", "moderator"},
    "/api/v1/content": []string{"user", "vip"},
}

// 使用中间件
r.Use(permissionMiddleware.Handler())
```

### 权限检查辅助函数

```go
// 检查用户是否有指定权限
func (h *Handler) requirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        roles := c.Get("roles")
        if !hasPermission(roles, permission) {
            c.JSON(403, gin.H{
                "code":    "403",
                "message": "权限不足",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// 使用
admin := r.Group("/api/v1/admin")
admin.Use(authMiddleware.Handler())
admin.Use(h.requirePermission("admin"))
{
    admin.GET("/users", listUsers)
}
```

---

## 限流

### 限流策略

中间件支持三种限流策略：

#### 1. 令牌桶算法 (Token Bucket)

适合突发流量场景。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| strategy | string | "token_bucket" | 策略名称 |
| rate | int | 100 | 每秒生成的令牌数 |
| burst | int | 200 | 桶的最大容量 |
| key_func | string | "client_ip" | 限流键函数 |

**使用示例：**

```go
config := &ratelimit.RateLimitConfig{
    Enabled:  true,
    Strategy: "token_bucket",
    Rate:     100,
    Burst:    200,
    KeyFunc:  "client_ip",
}
middleware, _ := ratelimit.NewRateLimitMiddleware(config, zapLogger)
```

#### 2. 滑动窗口算法 (Sliding Window)

适合精确控制请求速率的场景。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| strategy | string | "sliding_window" | 策略名称 |
| rate | int | 100 | 时间窗口内的请求数 |
| window_size | int | 60 | 时间窗口大小（秒） |
| key_func | string | "client_ip" | 限流键函数 |

**使用示例：**

```go
config := &ratelimit.RateLimitConfig{
    Enabled:    true,
    Strategy:   "sliding_window",
    Rate:       100,
    WindowSize: 60,
    KeyFunc:    "client_ip",
}
middleware, _ := ratelimit.NewRateLimitMiddleware(config, zapLogger)
```

#### 3. Redis分布式限流

适合分布式部署场景。

**配置参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| strategy | string | "redis" | 策略名称 |
| rate | int | 100 | 每秒请求数 |
| redis.addr | string | "" | Redis地址 |
| redis.password | string | "" | Redis密码 |
| redis.db | int | 0 | Redis数据库编号 |
| redis.prefix | string | "ratelimit:" | 键前缀 |

**使用示例：**

```go
config := &ratelimit.RateLimitConfig{
    Enabled:  true,
    Strategy: "redis",
    Rate:     100,
    Redis: &ratelimit.RedisConfig{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
        Prefix:   "ratelimit:",
    },
}
middleware, _ := ratelimit.NewRateLimitMiddleware(config, zapLogger)
```

### 配置示例

#### 按IP限流

```go
config := &ratelimit.RateLimitConfig{
    Enabled: true,
    Strategy: "token_bucket",
    Rate:     100,
    Burst:    200,
    KeyFunc:  "client_ip", // 按客户端IP
    SkipPaths: []string{"/health", "/metrics"},
}
```

#### 按用户限流

```go
config := &ratelimit.RateLimitConfig{
    Enabled: true,
    Strategy: "token_bucket",
    Rate:     50,
    Burst:    100,
    KeyFunc:  "user_id", // 按用户ID
    SkipPaths: []string{"/health", "/metrics"},
}
```

#### VIP用户差异化限流

```go
// 自定义键函数
func customKeyFunc(ctx *ratelimit.RateLimitContext) string {
    // 检查是否是VIP用户
    if ctx.Metadata["is_vip"] == true {
        return "vip:" + ctx.UserID
    }
    return "normal:" + ctx.ClientIP
}

config := &ratelimit.RateLimitConfig{
    Enabled:  true,
    Strategy: "token_bucket",
    Rate:     100,
    Burst:    200,
    KeyFunc:  "custom",
    SkipPaths: []string{"/health", "/metrics"},
}
```

### 使用场景

| 场景 | 策略 | 配置 |
|------|------|------|
| API网关 | token_bucket | rate: 1000, burst: 2000 |
| 登录接口 | sliding_window | rate: 10, window: 60 |
| 上传接口 | token_bucket | rate: 5, burst: 10 |
| 搜索接口 | redis | rate: 50, burst: 100 |

---

## 配置管理

### 配置文件格式

使用YAML格式的配置文件：

```yaml
# configs/middleware.yaml
middleware:
  # 请求ID配置
  request_id:
    enabled: true
    header_name: "X-Request-ID"
    force_gen: false

  # 日志配置
  logger:
    enabled: true
    skip_paths:
      - /health
      - /metrics
    enable_request_id: true
    enable_request_body: true
    enable_response_body: false
    slow_request_threshold: 3000
    max_body_size: 2048
    redact_keys:
      - authorization
      - password
      - token

  # CORS配置
  cors:
    enabled: true
    allow_origins:
      - "https://example.com"
      - "https://www.example.com"
    allow_methods:
      - GET
      - POST
      - PUT
      - DELETE
      - OPTIONS
    allow_headers:
      - Origin
      - Content-Type
      - Authorization
    allow_credentials: true
    max_age: 12h

  # 限流配置
  rate_limit:
    enabled: true
    strategy: "token_bucket"
    rate: 100
    burst: 200
    key_func: "client_ip"
    skip_paths:
      - /health
      - /metrics
    message: "请求过于频繁，请稍后再试"
    status_code: 429

  # 认证配置
  auth:
    enabled: true
    secret: "${JWT_SECRET}"
    access_expiration: 7200    # 2小时（秒）
    refresh_expiration: 604800 # 7天（秒）
    issuer: "qingyu"
    token_header: "Authorization"
    token_prefix: "Bearer"
    skip_paths:
      - /api/v1/auth/login
      - /api/v1/auth/register
      - /health
      - /metrics

  # 权限配置
  permission:
    enabled: true
    strategy: "rbac"
    config_path: "./configs/permissions.yaml"
    skip_paths: []

  # 恢复配置
  recovery:
    enabled: true
    stack_size: 4096
    disable_print: true

  # 安全配置
  security:
    enabled: true
    enable_x_frame_options: true
    x_frame_options: "DENY"
    enable_hsts: true
    enable_csp: true

  # 压缩配置
  compression:
    enabled: true
    level: 5
    types:
      - "application/json"
      - "text/html"
      - "text/plain"
```

### 热重载

支持运行时重新加载配置：

```go
import "github.com/spf13/viper"

func watchMiddlewareConfig(manager *core.Manager, configPath string) {
    v := viper.New()
    v.SetConfigFile(configPath)

    v.WatchConfig()
    v.OnConfigChange(func(e fsnotify.Event) {
        // 读取新配置
        config, err := loadMiddlewareConfig(configPath)
        if err != nil {
            log.Printf("配置加载失败: %v", err)
            return
        }

        // 重新加载中间件配置
        for _, middleware := range manager.List() {
            if hotReload, ok := middleware.(core.HotReloadMiddleware); ok {
                middlewareConfig := config.GetMiddlewareConfig(hotReload.Name())
                if middlewareConfig != nil {
                    if err := hotReload.Reload(middlewareConfig); err != nil {
                        log.Printf("中间件 %s 重载失败: %v", hotReload.Name(), err)
                    } else {
                        log.Printf("中间件 %s 重载成功", hotReload.Name())
                    }
                }
            }
        }
    })
}
```

### 配置验证

在加载配置后自动验证：

```go
// 加载配置
config, err := loadConfig("config/middleware.yaml")
if err != nil {
    log.Fatal(err)
}

// 验证配置
if err := config.Validate(); err != nil {
    log.Fatalf("配置验证失败: %v", err)
}

// 使用配置
middlewares, err := initializeFromConfig(config, logger)
if err != nil {
    log.Fatal(err)
}
```

### 环境变量支持

支持使用环境变量覆盖配置：

```bash
# 设置JWT密钥
export JWT_SECRET="your-secret-key"

# 设置限流速率
export RATE_LIMIT_RATE="100"

# 设置CORS允许的来源
export CORS_ALLOW_ORIGINS="https://example.com,https://www.example.com"
```

在配置文件中使用：

```yaml
auth:
  secret: "${JWT_SECRET}"

rate_limit:
  rate: "${RATE_LIMIT_RATE}"

cors:
  allow_origins: "${CORS_ALLOW_ORIGINS}"
```

---

## 自定义中间件

### 开发指南

#### 1. 实现核心接口

所有中间件必须实现 `Middleware` 接口：

```go
type Middleware interface {
    Name() string
    Priority() int
    Handler() gin.HandlerFunc
}
```

#### 2. 实现可配置接口（可选）

如果中间件需要配置，实现 `ConfigurableMiddleware` 接口：

```go
type ConfigurableMiddleware interface {
    Middleware
    LoadConfig(config map[string]interface{}) error
    ValidateConfig() error
}
```

#### 3. 实现热重载接口（可选）

如果中间件需要支持热重载，实现 `HotReloadMiddleware` 接口：

```go
type HotReloadMiddleware interface {
    ConfigurableMiddleware
    Reload(config map[string]interface{}) error
}
```

### 接口实现示例

```go
package custom

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "Qingyu_backend/internal/middleware/core"
)

// CustomMiddleware 自定义中间件
type CustomMiddleware struct {
    config *CustomConfig
    logger *zap.Logger
}

// CustomConfig 自定义配置
type CustomConfig struct {
    Enabled bool   `yaml:"enabled"`
    Option1 string `yaml:"option1"`
    Option2 int    `yaml:"option2"`
}

// NewCustomMiddleware 创建自定义中间件
func NewCustomMiddleware(logger *zap.Logger) *CustomMiddleware {
    return &CustomMiddleware{
        config: &CustomConfig{
            Enabled: true,
            Option1: "default",
            Option2: 100,
        },
        logger: logger,
    }
}

// Name 返回中间件名称
func (m *CustomMiddleware) Name() string {
    return "custom"
}

// Priority 返回执行优先级
func (m *CustomMiddleware) Priority() int {
    return 6 // 根据功能选择合适的优先级
}

// Handler 返回Gin处理函数
func (m *CustomMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查是否启用
        if !m.config.Enabled {
            c.Next()
            return
        }

        // 前置处理
        m.logger.Info("Custom middleware pre-processing",
            zap.String("path", c.Request.URL.Path))

        // 调用后续处理器
        c.Next()

        // 后置处理
        m.logger.Info("Custom middleware post-processing",
            zap.Int("status", c.Writer.Status()))
    }
}

// LoadConfig 从配置加载参数
func (m *CustomMiddleware) LoadConfig(config map[string]interface{}) error {
    if m.config == nil {
        m.config = &CustomConfig{}
    }

    // 加载Enabled
    if enabled, ok := config["enabled"].(bool); ok {
        m.config.Enabled = enabled
    }

    // 加载Option1
    if option1, ok := config["option1"].(string); ok {
        m.config.Option1 = option1
    }

    // 加载Option2
    if option2, ok := config["option2"].(int); ok {
        m.config.Option2 = option2
    }

    return nil
}

// ValidateConfig 验证配置有效性
func (m *CustomMiddleware) ValidateConfig() error {
    if m.config.Option1 == "" {
        return errors.New("option1不能为空")
    }

    if m.config.Option2 <= 0 {
        return errors.New("option2必须大于0")
    }

    return nil
}

// Reload 热重载配置
func (m *CustomMiddleware) Reload(config map[string]interface{}) error {
    // 保存旧配置
    oldConfig := *m.config

    // 加载新配置
    if err := m.LoadConfig(config); err != nil {
        m.config = &oldConfig
        return err
    }

    // 验证新配置
    if err := m.ValidateConfig(); err != nil {
        m.config = &oldConfig
        return err
    }

    m.logger.Info("Custom middleware reloaded")
    return nil
}

// 确保实现了接口
var _ core.Middleware = (*CustomMiddleware)(nil)
var _ core.ConfigurableMiddleware = (*CustomMiddleware)(nil)
var _ core.HotReloadMiddleware = (*CustomMiddleware)(nil)
```

### 最佳实践

#### 1. 单一职责

每个中间件只负责一个功能：

```go
// ✅ 好的设计：每个中间件负责单一功能
requestID := NewRequestIDMiddleware()
logger := NewLoggerMiddleware(logger)
recovery := NewRecoveryMiddleware(logger)

// ❌ 不好的设计：一个中间件负责多个功能
allInOne := NewAllInOneMiddleware()
```

#### 2. 无状态设计

中间件应该是无状态的，便于水平扩展：

```go
// ✅ 好的设计：无状态
func (m *Middleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 不修改中间件状态
        c.Set("key", "value")
        c.Next()
    }
}

// ❌ 不好的设计：有状态
func (m *Middleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        m.Counter++ // 修改中间件状态
        c.Next()
    }
}
```

#### 3. 错误处理

优雅处理错误，不影响其他中间件：

```go
func (m *Middleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 尝试处理
        if err := m.process(c); err != nil {
            // 记录错误但继续执行
            m.logger.Warn("Processing failed", zap.Error(err))
            c.Next()
            return
        }

        c.Next()
    }
}
```

#### 4. 性能优化

最小化中间件对请求处理的影响：

```go
// 使用对象池
var requestPool = sync.Pool{
    New: func() interface{} {
        return &RequestContext{}
    },
}

func (m *Middleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从池中获取对象
        ctx := requestPool.Get().(*RequestContext)
        defer requestPool.Put(ctx)

        // 重置对象
        *ctx = RequestContext{}

        // 使用对象
        ctx.Path = c.Request.URL.Path
        ctx.Method = c.Request.Method

        c.Next()
    }
}
```

#### 5. 配置验证

在启动时验证配置，避免运行时错误：

```go
func (m *Middleware) ValidateConfig() error {
    if m.config.Rate <= 0 {
        return errors.New("rate must be positive")
    }

    if m.config.Secret == "" {
        return errors.New("secret cannot be empty")
    }

    return nil
}

// 在创建中间件时验证
middleware := NewMiddleware(config)
if err := middleware.ValidateConfig(); err != nil {
    log.Fatalf("Invalid config: %v", err)
}
```

---

## API参考

### 核心接口

#### Middleware 接口

```go
type Middleware interface {
    // Name 返回中间件唯一标识
    Name() string

    // Priority 返回执行优先级
    // 返回值越小，中间件越先执行
    Priority() int

    // Handler 返回Gin处理函数
    Handler() gin.HandlerFunc
}
```

#### ConfigurableMiddleware 接口

```go
type ConfigurableMiddleware interface {
    Middleware

    // LoadConfig 从配置加载参数
    LoadConfig(config map[string]interface{}) error

    // ValidateConfig 验证配置有效性
    ValidateConfig() error
}
```

#### HotReloadMiddleware 接口

```go
type HotReloadMiddleware interface {
    ConfigurableMiddleware

    // Reload 热重载配置
    Reload(config map[string]interface{}) error
}
```

### 管理器接口

#### Manager 接口

```go
type Manager interface {
    // Register 注册中间件
    Register(middleware Middleware, options ...RegisterOption) error

    // Unregister 注销中间件
    Unregister(name string) error

    // Get 获取中间件
    Get(name string) (Middleware, error)

    // List 列出所有中间件
    List() []Middleware

    // Build 构建中间件链
    Build() []gin.HandlerFunc

    // ApplyToRouter 应用到Gin路由
    ApplyToRouter(router *gin.Engine, globalMiddlewares ...string) error

    // Validate 验证中间件配置
    Validate() error

    // GetExecutionOrder 获取执行顺序
    GetExecutionOrder() []string

    // GenerateOrderReport 生成顺序报告
    GenerateOrderReport() string
}
```

### 注册选项

```go
// RegisterOption 注册选项
type RegisterOption interface {
    apply(*registerConfig)
}

// WithPriorityOverride 设置优先级覆盖
func WithPriorityOverride(priority int) RegisterOption
```

### 配置结构

#### Config 总配置

```go
type Config struct {
    Middleware          MiddlewareConfigs `yaml:"middleware"`
    PriorityOverrides   map[string]int    `yaml:"priority_overrides,omitempty"`
}
```

#### MiddlewareConfigs 中间件配置集合

```go
type MiddlewareConfigs struct {
    RequestID   *RequestIDConfig   `yaml:"request_id,omitempty"`
    Recovery    *RecoveryConfig    `yaml:"recovery,omitempty"`
    Security    *SecurityConfig    `yaml:"security,omitempty"`
    Logger      *LoggerConfig      `yaml:"logger,omitempty"`
    Compression *CompressionConfig `yaml:"compression,omitempty"`
    RateLimit   *RateLimitConfig   `yaml:"rate_limit,omitempty"`
    Auth        *AuthConfig        `yaml:"auth,omitempty"`
    Permission  *PermissionConfig  `yaml:"permission,omitempty"`
}
```

### 辅助函数

#### 初始化函数

```go
// InitializeDefault 使用默认配置初始化中间件
func InitializeDefault(logger *zap.Logger) ([]Middleware, error)

// InitializeFromConfig 从配置文件初始化中间件
func InitializeFromConfig(config *Config, logger *zap.Logger) ([]Middleware, error)

// InitializeMiddleware 初始化单个中间件
func InitializeMiddleware(name string, config interface{}, logger *zap.Logger) (Middleware, error)
```

#### 辅助工具

```go
// GetMiddlewareOrder 获取中间件执行顺序
func GetMiddlewareOrder(manager Manager) []string

// ValidateMiddlewareConfig 验证中间件配置
func ValidateMiddlewareConfig(config *Config) error

// MergeConfigs 合并配置
func MergeConfigs(base, override *Config) *Config
```

---

## 测试指南

### 中间件测试

#### 基本测试结构

```go
package builtin_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "Qingyu_backend/internal/middleware/builtin"
)

func TestRequestIDMiddleware(t *testing.T) {
    // 设置Gin为测试模式
    gin.SetMode(gin.TestMode)

    // 创建测试路由
    router := gin.New()

    // 添加中间件
    requestID := builtin.NewRequestIDMiddleware()
    router.Use(requestID.Handler())

    // 添加测试处理器
    router.GET("/test", func(c *gin.Context) {
        requestID := c.GetString("request_id")
        c.JSON(200, gin.H{"request_id": requestID})
    })

    // 创建测试请求
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("X-Request-ID", "test-request-id")

    // 执行请求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 验证响应
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "test-request-id")
}
```

#### 测试配置

```go
func TestLoggerMiddleware(t *testing.T) {
    // 创建测试logger
    logger, _ := zap.NewDevelopment()

    // 创建自定义配置
    config := &builtin.LoggerConfig{
        SkipPaths:            []string{"/health"},
        SlowRequestThreshold: 100,
    }

    // 创建中间件
    middleware := builtin.NewLoggerMiddlewareWithConfig(config, logger)
    assert.NotNil(t, middleware)

    // 验证配置
    err := middleware.ValidateConfig()
    assert.NoError(t, err)
}
```

#### 测试错误处理

```go
func TestAuthMiddleware_ErrorCases(t *testing.T) {
    logger, _ := zap.NewDevelopment()
    jwtManager, _ := auth.NewJWTManager("secret", 2*time.Hour, 7*24*time.Hour)
    middleware := auth.NewJWTAuthMiddleware(jwtManager, nil, logger)

    tests := []struct {
        name           string
        token          string
        expectedStatus int
        expectedCode   string
    }{
        {
            name:           "Missing Token",
            token:          "",
            expectedStatus: 401,
            expectedCode:   "2010",
        },
        {
            name:           "Invalid Token",
            token:          "invalid-token",
            expectedStatus: 401,
            expectedCode:   "2008",
        },
        {
            name:           "Expired Token",
            token:          generateExpiredToken(),
            expectedStatus: 401,
            expectedCode:   "2007",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            router := gin.New()
            router.Use(middleware.Handler())
            router.GET("/test", func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "ok"})
            })

            req := httptest.NewRequest("GET", "/test", nil)
            if tt.token != "" {
                req.Header.Set("Authorization", "Bearer "+tt.token)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)
            assert.Contains(t, w.Body.String(), tt.expectedCode)
        })
    }
}
```

### 测试示例

#### 测试限流中间件

```go
func TestRateLimitMiddleware(t *testing.T) {
    logger, _ := zap.NewDevelopment()

    config := &ratelimit.RateLimitConfig{
        Enabled:  true,
        Strategy: "token_bucket",
        Rate:     10,
        Burst:    10,
        KeyFunc:  "client_ip",
    }

    middleware, err := ratelimit.NewRateLimitMiddleware(config, logger)
    assert.NoError(t, err)

    router := gin.New()
    router.Use(middleware.Handler())
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "ok"})
    })

    // 发送多个请求
    allowed := 0
    for i := 0; i < 15; i++ {
        req := httptest.NewRequest("GET", "/test", nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        if w.Code == 200 {
            allowed++
        }
    }

    // 验证限流效果
    assert.Equal(t, 10, allowed)
}
```

#### 测试权限中间件

```go
func TestPermissionMiddleware(t *testing.T) {
    logger, _ := zap.NewDevelopment()
    checker := auth.NewRBACChecker()

    // 设置权限规则
    checker.AddRule("/api/v1/admin", []string{"admin"})
    checker.AddRule("/api/v1/users", []string{"admin", "moderator"})

    middleware := auth.NewPermissionMiddleware(checker, logger)

    tests := []struct {
        name         string
        path         string
        roles        []string
        expectAllow  bool
    }{
        {
            name:        "Admin Access Admin Path",
            path:        "/api/v1/admin",
            roles:       []string{"admin"},
            expectAllow: true,
        },
        {
            name:        "User Cannot Access Admin Path",
            path:        "/api/v1/admin",
            roles:       []string{"user"},
            expectAllow: false,
        },
        {
            name:        "Moderator Access Users Path",
            path:        "/api/v1/users",
            roles:       []string{"moderator"},
            expectAllow: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            router := gin.New()
            router.Use(middleware.Handler())
            router.GET(tt.path, func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "ok"})
            })

            req := httptest.NewRequest("GET", tt.path, nil)
            w := httptest.NewRecorder()

            // 模拟认证中间件设置的角色
            ctx, _ := gin.CreateTestContext(w)
            ctx.Set("roles", tt.roles)
            ctx.Request = req

            router.ServeHTTP(w, req)

            if tt.expectAllow {
                assert.Equal(t, 200, w.Code)
            } else {
                assert.Equal(t, 403, w.Code)
            }
        })
    }
}
```

---

## 迁移指南

### 从 pkg/middleware 迁移

#### 1. 导入路径更新

```go
// 旧导入
import "Qingyu_backend/pkg/middleware"

// 新导入
import "Qingyu_backend/internal/middleware"
```

#### 2. 中间件初始化

**旧方式：**

```go
// pkg/middleware
authMiddleware := middleware.AuthMiddleware()
rateLimitMiddleware := middleware.RateLimitMiddleware(100, 60)
```

**新方式：**

```go
// internal/middleware
jwtManager, _ := auth.NewJWTManager("secret", 2*time.Hour, 7*24*time.Hour)
authMiddleware := auth.NewJWTAuthMiddleware(jwtManager, nil, logger)

config := &ratelimit.RateLimitConfig{
    Enabled:  true,
    Strategy: "token_bucket",
    Rate:     100,
    Burst:    200,
}
rateLimitMiddleware, _ := ratelimit.NewRateLimitMiddleware(config, logger)
```

#### 3. 路由使用更新

**旧方式：**

```go
// pkg/middleware
r.Use(middleware.CORS())
r.Use(middleware.Logger())
r.Use(middleware.Auth())
```

**新方式：**

```go
// internal/middleware
manager := core.NewManager(logger)
manager.Register(builtin.NewCORSMiddleware(corsConfig))
manager.Register(builtin.NewLoggerMiddleware(logger))
manager.Register(auth.NewJWTAuthMiddleware(jwtManager, nil, logger))

manager.ApplyToRouter(r)
```

### 兼容性说明

#### 已废弃的中间件

以下中间件已废弃，请使用新的替代方案：

| 旧中间件 | 新中间件 | 说明 |
|---------|---------|------|
| pkg/middleware.Auth | internal/middleware/auth.JWTAuthMiddleware | JWT认证 |
| pkg/middleware.RateLimit | internal/middleware/ratelimit.RateLimitMiddleware | 限流 |
| pkg/middleware.CORS | internal/middleware/builtin.CORSMiddleware | CORS |
| pkg/middleware.Logger | internal/middleware/builtin.LoggerMiddleware | 日志 |

#### API 变更

```go
// 旧API
middleware.AuthMiddleware() gin.HandlerFunc

// 新API
auth.NewJWTAuthMiddleware(jwtManager, blacklist, logger) *JWTAuthMiddleware
// 需要调用 Handler() 方法获取 gin.HandlerFunc
```

#### 配置变更

```go
// 旧配置（代码配置）
middleware.SetRateLimit(100, 60)

// 新配置（YAML配置）
rate_limit:
  enabled: true
  strategy: "token_bucket"
  rate: 100
  burst: 200
```

### 迁移步骤

1. **备份现有代码**
   ```bash
   git checkout -b middleware-migration
   ```

2. **更新导入路径**
   ```bash
   find . -name "*.go" -type f -exec sed -i 's|Qingyu_backend/pkg/middleware|Qingyu_backend/internal/middleware|g' {} +
   ```

3. **更新初始化代码**
   - 创建中间件管理器
   - 更新中间件创建方式
   - 添加必要的配置

4. **测试验证**
   - 运行单元测试
   - 运行集成测试
   - 手动测试关键功能

5. **清理旧代码**
   ```bash
   rm -rf pkg/middleware
   ```

---

## 最佳实践

### 中间件使用建议

#### 1. 优先级设置

根据功能正确设置中间件优先级：

```go
// 基础设施：优先级1-5
RequestID:    1
Security:     2
CORS:         3
Recovery:     5

// 监控层：优先级6-8
Metrics:      6
Logger:       7
RateLimit:    8

// 业务层：优先级9-12
JWT:          9
Permission:   10
Validation:   11
Compression:  12
```

#### 2. 路径跳过

合理配置跳过路径，避免不必要的处理：

```go
config := &auth.JWTConfig{
    SkipPaths: []string{
        "/health",       // 健康检查
        "/metrics",      // 监控指标
        "/api/v1/auth",  // 认证相关
    },
}
```

#### 3. 错误处理

统一错误处理格式：

```go
func (m *Middleware) respondWithError(c *gin.Context, code, message string) {
    c.JSON(401, gin.H{
        "code":    code,
        "message": message,
    })
    c.Abort()
}
```

#### 4. 日志记录

合理使用日志级别：

```go
// Debug: 详细信息
logger.Debug("Request details", zap.String("path", path))

// Info: 正常流程
logger.Info("Request processed", zap.Int("status", status))

// Warn: 警告信息
logger.Warn("Slow request", zap.Int64("latency_ms", latencyMs))

// Error: 错误信息
logger.Error("Processing failed", zap.Error(err))
```

### 性能优化

#### 1. 对象池使用

```go
var contextPool = sync.Pool{
    New: func() interface{} {
        return &RequestContext{}
    },
}

func getContext() *RequestContext {
    return contextPool.Get().(*RequestContext)
}

func putContext(ctx *RequestContext) {
    *ctx = RequestContext{}
    contextPool.Put(ctx)
}
```

#### 2. 避免重复解析

```go
// ✅ 好：解析一次，存储到上下文
func (m *Middleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        claims := validateToken(token)
        c.Set("claims", claims)
        c.Next()
    }
}

// ❌ 差：每次都重新解析
func (m *Middleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 每次都解析token
        claims := validateToken(extractToken(c))
        c.Next()
    }
}
```

#### 3. 并发安全

```go
// 使用读写锁保护共享数据
type Middleware struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (m *Middleware) Get(key string) (interface{}, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    val, ok := m.data[key]
    return val, ok
}
```

### 安全建议

#### 1. 敏感信息保护

```go
// 脱敏敏感信息
func redactSensitive(data map[string]interface{}) {
    sensitiveKeys := []string{"password", "token", "secret"}
    for _, key := range sensitiveKeys {
        if _, exists := data[key]; exists {
            data[key] = "***"
        }
    }
}
```

#### 2. 输入验证

```go
func validateInput(input string) error {
    // 验证长度
    if len(input) > 1000 {
        return errors.New("input too long")
    }

    // 验证字符
    matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", input)
    if !matched {
        return errors.New("invalid characters")
    }

    return nil
}
```

#### 3. 速率限制

```go
// 对敏感操作实施更严格的限流
sensitiveOperationConfig := &ratelimit.RateLimitConfig{
    Enabled:  true,
    Strategy: "sliding_window",
    Rate:     10,
    Burst:    20,
    KeyFunc:  "user_id", // 按用户限流
}
```

---

## 相关文档

- [设计文档](../../docs/design/middleware/README_中间件设计文档.md)
- [API架构文档](../../docs/architecture/api_architecture.md)
- [API v1文档](../../api/v1/README.md)
- [实施计划](../../docs/plans/api_refactor_plan.md)

---

## 版本历史

| 版本 | 日期 | 变更内容 |
|------|------|----------|
| v1.0 | 2026-02-26 | 初始版本 |

---

**维护者**: Backend Team
**最后更新**: 2026-02-26
