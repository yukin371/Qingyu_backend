# Middleware 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **状态**: 进行中

---

## 1. 职责边界与依赖关系

### 1.1 职责定义

Middleware 层是后端系统的**横切关注点层**，负责：

1. **请求预处理**：在业务逻辑执行前对请求进行处理
2. **响应后处理**：在业务逻辑执行后对响应进行处理
3. **认证授权**：验证用户身份和权限
4. **限流保护**：防止系统过载
5. **日志记录**：记录请求日志
6. **错误恢复**：捕获 panic 并恢复
7. **安全防护**：CORS、安全头等

### 1.2 与上下层级的交互边界

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP 请求                           │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                   Middleware 层                         │
│  ┌─────────────────────────────────────────────────────┤
│  │ 优先级 1-5: 基础设施                                │
│  │   - RequestID 生成请求ID                            │
│  │   - Recovery 恢复panic                              │
│  │   - Security 安全头                                 │
│  │   - CORS 跨域支持                                   │
│  ├─────────────────────────────────────────────────────┤
│  │ 优先级 6-8: 监控日志                                │
│  │   - Logger 请求日志                                 │
│  │   - Metrics 性能指标                                │
│  ├─────────────────────────────────────────────────────┤
│  │ 优先级 9-10: 认证授权                               │
│  │   - RateLimit 限流                                  │
│  │   - Auth JWT认证                                    │
│  │   - Permission 权限检查                             │
│  ├─────────────────────────────────────────────────────┤
│  │ 优先级 11-12: 业务层                                │
│  │   - Validation 参数验证                             │
│  │   - Compression 响应压缩                            │
│  └─────────────────────────────────────────────────────┤
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                    Router 层                            │
│              (路由分发)                                 │
└─────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│                     API 层                              │
│              (请求处理)                                 │
└─────────────────────────────────────────────────────────┘
```

### 1.3 依赖关系

```go
// Middleware 层允许的依赖
import (
    "github.com/gin-gonic/gin"
    "Qingyu_backend/pkg/logger"         // 日志
    "Qingyu_backend/pkg/cache"          // 缓存（限流、黑名单）
    "Qingyu_backend/config"             // 配置读取
)

// Middleware 层禁止的依赖
import (
    "Qingyu_backend/service/xxx"        // ❌ 禁止依赖 Service
    "Qingyu_backend/repository/xxx"     // ❌ 禁止依赖 Repository
)
```

---

## 2. 命名与代码规范

### 2.1 文件命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 中间件定义 | `{功能}.go` | `jwt.go`, `recovery.go`, `cors.go` |
| 中间件接口 | `middleware.go` | `middleware.go` |
| 管理器 | `manager.go` | `manager.go`, `manager_impl.go` |
| 测试文件 | `{文件名}_test.go` | `jwt_test.go` |

### 2.2 结构体命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 中间件结构体 | `{功能}Middleware` | `JWTAuthMiddleware`, `RateLimitMiddleware` |
| 配置结构体 | `{功能}Config` | `JWTConfig`, `RateLimitConfig` |
| 接口 | `{功能}` | `Middleware`, `Manager`, `Blacklist` |

### 2.3 目录组织规范

```
internal/middleware/
├── core/                   # 核心接口和基础实现
│   ├── middleware.go       # Middleware 接口定义
│   ├── manager.go          # Manager 接口定义
│   ├── manager_impl.go     # Manager 实现
│   ├── registry.go         # 中间件注册表
│   ├── initializer.go      # 初始化器
│   └── context.go          # 上下文工具
├── builtin/                # 内置中间件
│   ├── request_id.go       # 请求ID
│   ├── recovery.go         # 恢复
│   ├── logger.go           # 日志
│   ├── cors.go             # 跨域
│   ├── security.go         # 安全头
│   ├── compression.go      # 压缩
│   └── error_handler.go    # 错误处理
├── auth/                   # 认证授权中间件
│   ├── jwt.go              # JWT认证
│   ├── jwt_manager.go      # JWT管理器
│   ├── blacklist.go        # Token黑名单
│   ├── permission.go       # 权限检查
│   └── rbac_checker.go     # RBAC检查器
├── ratelimit/              # 限流中间件
│   ├── rate_limit.go       # 限流入口
│   ├── limiter.go          # 限流器接口
│   ├── token_bucket.go     # 令牌桶算法
│   ├── sliding_window.go   # 滑动窗口算法
│   └── redis_limiter.go    # Redis限流器
├── validation/             # 验证中间件
│   └── validation.go       # 参数验证
├── monitoring/             # 监控中间件
│   └── metrics.go          # 性能指标
└── admin/                  # 管理员中间件
    └── admin_permission.go # 管理员权限
```

---

## 3. 设计模式与最佳实践

### 3.1 Middleware 接口模式

```go
// Middleware 中间件核心接口
type Middleware interface {
    // Name 返回中间件唯一标识
    Name() string

    // Priority 返回执行优先级
    // 返回值越小，中间件越先执行
    Priority() int

    // Handler 返回Gin处理函数
    Handler() gin.HandlerFunc
}

// ConfigurableMiddleware 可配置中间件接口（可选）
type ConfigurableMiddleware interface {
    Middleware

    // LoadConfig 从配置加载参数
    LoadConfig(config map[string]interface{}) error

    // ValidateConfig 验证配置有效性
    ValidateConfig() error
}

// HotReloadMiddleware 支持热更新的中间件（v2.0新增）
type HotReloadMiddleware interface {
    ConfigurableMiddleware

    // Reload 热重载配置
    Reload(config map[string]interface{}) error
}
```

### 3.2 优先级规范

| 优先级范围 | 类型 | 中间件示例 |
|------------|------|------------|
| 1-5 | 基础设施 | RequestID, Recovery, Security, CORS |
| 6-8 | 监控日志 | Timeout, Logger, Metrics |
| 9-10 | 认证授权 | RateLimit, Auth, Permission |
| 11-12 | 业务层 | Validation, Compression |

```go
const (
    // 基础设施
    RequestIDPriority = 1
    RecoveryPriority  = 2
    SecurityPriority  = 3
    CORSPriority      = 4

    // 监控日志
    TimeoutPriority = 6
    LoggerPriority  = 7
    MetricsPriority = 8

    // 认证授权
    RateLimitPriority  = 9
    AuthPriority       = 9
    PermissionPriority = 10

    // 业务层
    ValidationPriority  = 11
    CompressionPriority = 12
)
```

### 3.3 中间件实现模板

```go
// jwt.go
package auth

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "Qingyu_backend/internal/middleware/core"
)

// JWTAuthMiddleware JWT认证中间件
type JWTAuthMiddleware struct {
    config     *JWTConfig
    jwtManager JWTManager
    blacklist  Blacklist
    logger     *zap.Logger
}

// NewJWTAuthMiddleware 创建JWT认证中间件
func NewJWTAuthMiddleware(
    jwtManager JWTManager,
    blacklist Blacklist,
    logger *zap.Logger,
) *JWTAuthMiddleware {
    return &JWTAuthMiddleware{
        config:     DefaultJWTConfig(),
        jwtManager: jwtManager,
        blacklist:  blacklist,
        logger:     logger,
    }
}

// Name 返回中间件名称
func (m *JWTAuthMiddleware) Name() string {
    return "jwt_auth"
}

// Priority 返回执行优先级
func (m *JWTAuthMiddleware) Priority() int {
    return AuthPriority
}

// Handler 返回Gin处理函数
func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 检查是否跳过该路径
        if m.shouldSkipPath(c.Request.URL.Path) {
            c.Next()
            return
        }

        // 2. 提取Token
        token, err := m.extractToken(c)
        if err != nil {
            m.respondWithError(c, err)
            c.Abort()
            return
        }

        // 3. 验证Token
        claims, err := m.validateToken(token)
        if err != nil {
            m.respondWithError(c, err)
            c.Abort()
            return
        }

        // 4. 检查黑名单
        if m.blacklist != nil {
            isBlacklisted, _ := m.blacklist.IsBlacklisted(c.Request.Context(), token)
            if isBlacklisted {
                m.respondWithError(c, errors.New("token revoked"))
                c.Abort()
                return
            }
        }

        // 5. 注入用户信息到上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("roles", claims.Roles)

        c.Next()
    }
}

// 确保实现接口
var _ core.ConfigurableMiddleware = (*JWTAuthMiddleware)(nil)
```

### 3.4 Manager 模式

```go
// Manager 中间件管理器接口
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
}

// 使用示例
manager := core.NewManager()
manager.Register(requestID.NewMiddleware())
manager.Register(recovery.NewMiddleware())
manager.Register(logger.NewMiddleware())
manager.Register(auth.NewJWTAuthMiddleware(...))

// 应用到路由
manager.ApplyToRouter(engine)
```

### 3.5 反模式警示

❌ **禁止事项**：

```go
// ❌ 禁止：在中间件中调用 Service
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        user, _ := m.userService.GetUser(...)  // 禁止
    }
}

// ❌ 禁止：中间件中包含业务逻辑
func (m *ValidationMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        if user.Age < 18 {  // 业务规则，应该在 Service 层
            c.AbortWithStatusJSON(400, ...)
        }
    }
}

// ❌ 禁止：直接操作数据库
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        m.db.Collection("users").FindOne(...)  // 禁止
    }
}
```

---

## 4. 接口与契约规范

### 4.1 上下文数据传递

```go
// 标准上下文键
const (
    ContextKeyUserID    = "user_id"
    ContextKeyUsername  = "username"
    ContextKeyRoles     = "roles"
    ContextKeyTokenType = "token_type"
    ContextKeyRequestID = "request_id"
)

// 设置上下文数据
c.Set("user_id", claims.UserID)
c.Set("username", claims.Username)
c.Set("roles", claims.Roles)

// 获取上下文数据
userID, exists := c.Get("user_id")
if !exists {
    // 未认证
}
```

### 4.2 错误响应格式

```go
// 认证错误响应
{
    "code": "2007",           // 业务错误码
    "message": "Token已过期"
}

// 权限错误响应
{
    "code": "1003",
    "message": "权限不足"
}

// 限流错误响应
{
    "code": "4290",
    "message": "请求过于频繁，请稍后重试"
}
```

### 4.3 配置加载规范

```go
// 从配置文件加载
func (m *JWTAuthMiddleware) LoadConfig(config map[string]interface{}) error {
    if enabled, ok := config["enabled"].(bool); ok {
        m.config.Enabled = enabled
    }
    if secret, ok := config["secret"].(string); ok {
        m.config.Secret = secret
    }
    if skipPaths, ok := config["skip_paths"].([]interface{}); ok {
        m.config.SkipPaths = make([]string, len(skipPaths))
        for i, v := range skipPaths {
            if str, ok := v.(string); ok {
                m.config.SkipPaths[i] = str
            }
        }
    }
    return nil
}
```

---

## 5. 测试策略

### 5.1 单元测试编写指南

```go
// jwt_test.go
package auth

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockJWTManager 模拟JWT管理器
type MockJWTManager struct {
    mock.Mock
}

func (m *MockJWTManager) ValidateToken(token string) (*Claims, error) {
    args := m.Called(token)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Claims), args.Error(1)
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
    // 1. 准备
    gin.SetMode(gin.TestMode)

    mockManager := new(MockJWTManager)
    mockManager.On("ValidateToken", "valid_token").
        Return(&Claims{UserID: "123", Username: "test"}, nil)

    middleware := NewJWTAuthMiddleware(mockManager, nil, nil)

    // 2. 执行
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/test", nil)
    c.Request.Header.Set("Authorization", "Bearer valid_token")

    middleware.Handler()(c)

    // 3. 验证
    assert.False(t, c.IsAborted())
    assert.Equal(t, "123", c.GetString("user_id"))
    mockManager.AssertExpectations(t)
}

func TestJWTAuthMiddleware_MissingToken(t *testing.T) {
    gin.SetMode(gin.TestMode)

    middleware := NewJWTAuthMiddleware(nil, nil, nil)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/test", nil)

    middleware.Handler()(c)

    assert.True(t, c.IsAborted())
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
```

### 5.2 测试覆盖率要求

| 测试类型 | 覆盖率要求 |
|----------|------------|
| 单元测试 | ≥ 80% |
| 边界条件 | 100% |
| 错误处理 | ≥ 90% |

---

## 6. 完整代码示例

### 6.1 完整中间件示例

```go
// ratelimit/rate_limit.go
package ratelimit

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "Qingyu_backend/internal/middleware/core"
)

const (
    RateLimitPriority = 9
)

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
    config    *RateLimitConfig
    limiter   Limiter
    logger    *zap.Logger
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
    Enabled        bool          `yaml:"enabled"`
    RequestsPerSec float64       `yaml:"requests_per_sec"`
    Burst          int           `yaml:"burst"`
    SkipPaths      []string      `yaml:"skip_paths"`
    WindowSize     time.Duration `yaml:"window_size"`
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() *RateLimitConfig {
    return &RateLimitConfig{
        Enabled:        true,
        RequestsPerSec: 100,
        Burst:          200,
        SkipPaths:      []string{"/health", "/metrics"},
        WindowSize:     time.Minute,
    }
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(limiter Limiter, logger *zap.Logger) *RateLimitMiddleware {
    return &RateLimitMiddleware{
        config:  DefaultRateLimitConfig(),
        limiter: limiter,
        logger:  logger,
    }
}

// Name 返回中间件名称
func (m *RateLimitMiddleware) Name() string {
    return "rate_limit"
}

// Priority 返回执行优先级
func (m *RateLimitMiddleware) Priority() int {
    return RateLimitPriority
}

// Handler 返回Gin处理函数
func (m *RateLimitMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !m.config.Enabled {
            c.Next()
            return
        }

        // 检查是否跳过
        if m.shouldSkipPath(c.Request.URL.Path) {
            c.Next()
            return
        }

        // 获取客户端标识（IP或用户ID）
        key := m.getClientKey(c)

        // 检查限流
        allowed, remaining, resetTime := m.limiter.Allow(key)
        if !allowed {
            // 设置响应头
            c.Header("X-RateLimit-Remaining", "0")
            c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
            c.Header("Retry-After", resetTime.Sub(time.Now()).String())

            c.JSON(http.StatusTooManyRequests, gin.H{
                "code":    4290,
                "message": "请求过于频繁，请稍后重试",
            })
            c.Abort()
            return
        }

        // 设置响应头
        c.Header("X-RateLimit-Remaining", string(remaining))
        c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

        c.Next()
    }
}

// getClientKey 获取客户端标识
func (m *RateLimitMiddleware) getClientKey(c *gin.Context) string {
    // 优先使用用户ID
    if userID, exists := c.Get("user_id"); exists {
        return "user:" + userID.(string)
    }
    // 否则使用IP
    return "ip:" + c.ClientIP()
}

// shouldSkipPath 检查是否跳过限流
func (m *RateLimitMiddleware) shouldSkipPath(path string) bool {
    for _, skipPath := range m.config.SkipPaths {
        if path == skipPath {
            return true
        }
    }
    return false
}

// LoadConfig 加载配置
func (m *RateLimitMiddleware) LoadConfig(config map[string]interface{}) error {
    if enabled, ok := config["enabled"].(bool); ok {
        m.config.Enabled = enabled
    }
    if rps, ok := config["requests_per_sec"].(float64); ok {
        m.config.RequestsPerSec = rps
    }
    if burst, ok := config["burst"].(int); ok {
        m.config.Burst = burst
    }
    return nil
}

// ValidateConfig 验证配置
func (m *RateLimitMiddleware) ValidateConfig() error {
    if m.config.RequestsPerSec <= 0 {
        return errors.New("requests_per_sec must be positive")
    }
    if m.config.Burst < 0 {
        return errors.New("burst cannot be negative")
    }
    return nil
}

// 确保实现接口
var _ core.ConfigurableMiddleware = (*RateLimitMiddleware)(nil)
```

---

## 7. 参考资料

- [Middleware 层快速参考](../internal/middleware/README.md)
- [Router 层设计说明](./layer-router.md)
- [Config 层设计说明](./layer-config.md)

---

*最后更新：2026-03-19*
