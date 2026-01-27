# 中间件架构培训文档

## 培训信息

- **培训时长**: 2小时
- **培训对象**: 后端开发团队
- **培训日期**: 2026-01-27
- **培训讲师**: 架构师团队
- **目标**: 让团队成员理解新的中间件架构设计，掌握重构后的开发模式

---

## 目录

1. [重构背景](#重构背景)
2. [新架构设计](#新架构设计)
3. [配置分类](#配置分类)
4. [优先级管理](#优先级管理)
5. [错误处理策略](#错误处理策略)
6. [Phase实施计划](#phase实施计划)
7. [快速上手示例](#快速上手示例)
8. [Q&A](#qa)
9. [参考资料](#参考资料)

---

## 重构背景

### 当前问题

#### 1. 代码重复严重
每个中间件都有独立的配置加载、错误处理、日志记录逻辑，导致代码高度重复：

```go
// 每个中间件都重复这段代码
if config.Enable {
    // 中间件逻辑
}
```

**影响**：
- 维护成本高（修改一个功能需要改多个文件）
- 代码膨胀（`internal/middleware/` 目录超过 8000 行）
- 容易出现不一致的行为

#### 2. 配置分散且混乱
- 全局配置（`middleware.Config`）和独立配置混合
- 配置来源不统一（Viper、环境变量、硬编码）
- 配置验证逻辑缺失

**影响**：
- 配置冲突风险高
- 难以追溯配置来源
- 运行时配置错误难以排查

#### 3. 优先级管理混乱
- `priority.go` 定义优先级，但缺乏强制执行
- 依赖关系不明确（如：Recovery 应该最先执行）
- 手动排序容易出错

**影响**：
- 中间件执行顺序不确定
- 关键中间件可能被意外绕过
- 调试困难

#### 4. 错误处理不统一
```go
// 有些中间件直接 panic
if err != nil {
    panic(err)
}

// 有些返回 JSON
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
}

// 有些只记录日志
if err != nil {
    logger.Error(err)
}
```

**影响**：
- 客户端收到不一致的错误响应
- 错误日志难以分析
- 可能泄露敏感信息

#### 5. 测试覆盖率低
- 单元测试覆盖率不足 30%
- 缺乏集成测试
- Mock 依赖不完整

### 重构目标

#### 目标1: 代码复用率提升至 80%
通过引入 `MiddlewareCore` 统一管理配置加载、错误处理、日志记录等通用逻辑。

**指标**：
- 减少代码行数 50%（从 8000+ 行降至 4000 行）
- 新增中间件开发时间缩短 60%

#### 目标2: 配置管理标准化
实现统一的配置加载、验证、热更新机制。

**指标**：
- 配置冲突率降至 0%
- 配置错误启动时 100% 检测

#### 目标3: 优先级自动管理
基于依赖关系自动计算中间件执行顺序。

**指标**：
- 优先级错误率降至 0%
- 支持动态调整优先级

#### 目标4: 错误处理统一化
使用统一的错误处理中间件和错误码体系。

**指标**：
- 错误响应格式一致性 100%
- 错误日志可追溯性 100%

#### 目标5: 测试覆盖率提升至 90%
完善单元测试、集成测试、端到端测试。

**指标**：
- 单元测试覆盖率 ≥ 90%
- 集成测试覆盖率 ≥ 80%
- 关键路径 E2E 测试覆盖率 100%

---

## 新架构设计

### 目录结构

```
internal/middleware/
├── core/                    # 核心组件
│   ├── core.go             # MiddlewareCore 实现
│   ├── config.go           # 配置管理器
│   ├── errors.go           # 错误处理器
│   ├── logger.go           # 日志记录器
│   └── metrics.go          # 性能监控
├── base/                    # 基础中间件
│   ├── recovery.go         # Recovery（优先级：1000）
│   ├── request_id.go       # RequestID（优先级：900）
│   ├── timeout.go          # Timeout（优先级：850）
│   └── cors.go             # CORS（优先级：800）
├── auth/                    # 认证相关
│   ├── jwt.go              # JWT认证（优先级：700）
│   ├── rbac.go             # RBAC权限（优先级：600）
│   └── oauth.go            # OAuth（优先级：650）
├── security/                # 安全相关
│   ├── rate_limit.go       # 限流（优先级：500）
│   ├── csrf.go             # CSRF防护（优先级：450）
│   └── security.go         # 安全头（优先级：400）
├── monitoring/              # 监控相关
│   ├── access_log.go       # 访问日志（优先级：300）
│   └── prometheus.go       # Prometheus指标（优先级：250）
├── performance/             # 性能相关
│   ├── cache.go            # 缓存（优先级：200）
│   └── compress.go         # 压缩（优先级：150）
├── error_handler/           # 错误处理
│   └── error_handler.go    # 统一错误处理（优先级：100）
├── registry.go              # 中间件注册表
├── manager.go               # 中间件管理器
└── examples/                # 示例代码
    ├── static_example.go    # 静态配置示例
    └── dynamic_example.go   # 动态配置示例
```

### 核心组件

#### 1. MiddlewareCore

**职责**：
- 统一配置加载与验证
- 统一错误处理
- 统一日志记录
- 统一性能监控

**接口定义**：
```go
type MiddlewareCore interface {
    // 配置管理
    LoadConfig(configPath string) error
    ValidateConfig() error
    GetConfig(key string) interface{}
    UpdateConfig(key string, value interface{}) error

    // 错误处理
    HandleError(c *gin.Context, err error)
    HandlePanic(c *gin.Context, recovered interface{})

    // 日志记录
    LogRequest(c *gin.Context)
    LogResponse(c *gin.Context, duration time.Duration)
    LogError(c *gin.Context, err error)

    // 性能监控
    RecordMetric(name string, value float64, tags ...string)
    GetMetrics() map[string]interface{}
}
```

**使用示例**：
```go
// 在中间件中使用 core
core := middleware.GetCore()

func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 加载配置
        config := core.GetConfig("my_middleware")
        enabled := config.(map[string]interface{})["enabled"].(bool)

        if !enabled {
            c.Next()
            return
        }

        // 记录请求
        core.LogRequest(c)

        // 执行业务逻辑
        start := time.Now()
        c.Next()
        duration := time.Since(start)

        // 记录响应
        core.LogResponse(c, duration)

        // 记录指标
        core.RecordMetric("my_middleware.duration", duration.Seconds())
    }
}
```

#### 2. MiddlewareRegistry

**职责**：
- 管理所有已注册的中间件
- 维护中间件元数据（名称、优先级、依赖）
- 自动计算执行顺序

**接口定义**：
```go
type MiddlewareRegistry interface {
    // 注册中间件
    Register(name string, middleware MiddlewareConfig) error

    // 获取中间件
    Get(name string) (MiddlewareConfig, error)

    // 获取所有中间件（按优先级排序）
    GetAll() []MiddlewareConfig

    // 计算执行顺序
    CalculateOrder() ([]string, error)

    // 验证依赖关系
    ValidateDependencies() error
}

type MiddlewareConfig struct {
    Name        string
    Priority    int
    Handler     gin.HandlerFunc
    Enabled     bool
    Static      bool          // 是否为静态配置
    Config      interface{}   // 配置对象
    Dependencies []string     // 依赖的中间件
    Metadata    map[string]interface{}
}
```

**使用示例**：
```go
// 注册中间件
registry := middleware.NewRegistry()

registry.Register("recovery", middleware.MiddlewareConfig{
    Name:     "recovery",
    Priority: 1000,
    Handler:  Recovery(),
    Enabled:  true,
    Static:   true,
})

registry.Register("request_id", middleware.MiddlewareConfig{
    Name:     "request_id",
    Priority: 900,
    Handler:  RequestID(),
    Enabled:  true,
    Static:   true,
    Dependencies: []string{"recovery"}, // 依赖 recovery
})

// 计算执行顺序
order, _ := registry.CalculateOrder()
// 输出: ["recovery", "request_id", ...]
```

#### 3. MiddlewareManager

**职责**：
- 加载和管理中间件配置
- 动态启用/禁用中间件
- 提供中间件状态查询

**接口定义**：
```go
type MiddlewareManager interface {
    // 加载配置
    LoadConfig(configPath string) error

    // 启用/禁用中间件
    Enable(name string) error
    Disable(name string) error

    // 查询状态
    IsEnabled(name string) bool
    GetStatus(name string) MiddlewareStatus

    // 获取所有中间件状态
    GetAllStatus() map[string]MiddlewareStatus
}

type MiddlewareStatus struct {
    Name      string
    Enabled   bool
    Priority  int
    Config    interface{}
    Metrics   map[string]interface{}
}
```

**使用示例**：
```go
// 创建管理器
manager := middleware.NewManager(registry, core)

// 加载配置
manager.LoadConfig("config/middleware.yaml")

// 动态启用/禁用
manager.Enable("rate_limit")
manager.Disable("csrf")

// 查询状态
status := manager.GetStatus("rate_limit")
fmt.Printf("Enabled: %v, Priority: %d\n", status.Enabled, status.Priority)
```

### 职责分离

#### 静态配置中间件（Base）

**特点**：
- 启动时加载配置
- 运行时不可修改
- 优先级固定
- 无外部依赖

**示例**：
- `Recovery`: 恢复 panic
- `RequestID`: 生成请求 ID
- `CORS`: 跨域资源共享

#### 动态配置中间件（Auth/Security/Monitoring）

**特点**：
- 运行时可修改配置
- 支持热更新
- 优先级可调整
- 可能有外部依赖（如 Redis）

**示例**：
- `RateLimit`: 限流配置
- `RBAC`: 权限规则
- `AccessLog`: 日志级别

---

## 配置分类

### 静态配置

**定义**：启动时加载，运行时不可修改的配置。

**配置示例** (`config/middleware.yaml`)：
```yaml
# 静态配置中间件
static:
  # Recovery 中间件
  recovery:
    enabled: true
    stack: true
    stack_size: 4 << 10  # 4KB

  # RequestID 中间件
  request_id:
    enabled: true
    generator: "uuid"     # uuid, snowflake
    header: "X-Request-ID"

  # CORS 中间件
  cors:
    enabled: true
    allow_origins:
      - "http://localhost:3000"
      - "https://qingyu.example.com"
    allow_methods:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
    allow_headers:
      - "Content-Type"
      - "Authorization"
    expose_headers:
      - "X-Request-ID"
    max_age: 86400
    allow_credentials: true
```

**代码实现**：
```go
// 定义配置结构体
type RecoveryConfig struct {
    Enabled  bool `json:"enabled" yaml:"enabled"`
    Stack    bool `json:"stack" yaml:"stack"`
    StackSize int  `json:"stack_size" yaml:"stack_size"`
}

// 加载配置
func (r *RecoveryMiddleware) LoadConfig(config map[string]interface{}) error {
    r.config.Enabled = config["enabled"].(bool)
    r.config.Stack = config["stack"].(bool)
    r.config.StackSize = config["stack_size"].(int)
    return nil
}
```

### 动态配置

**定义**：运行时可以修改的配置，支持热更新。

**配置示例** (`config/middleware.yaml`)：
```yaml
# 动态配置中间件
dynamic:
  # RateLimit 中间件
  rate_limit:
    enabled: true
    driver: "redis"          # memory, redis
    # 全局限速
    global:
      rate: 1000             # 每秒1000个请求
      burst: 2000
    # 按IP限速
    per_ip:
      rate: 100              # 每IP每秒100个请求
      burst: 200
    # 按用户限速
    per_user:
      rate: 50               # 每用户每秒50个请求
      burst: 100

  # RBAC 中间件
  rbac:
    enabled: true
    driver: "redis"          # memory, redis, database
    cache_ttl: 300           # 缓存5分钟
    # 权限规则
    rules:
      - role: "admin"
        resources:
          - pattern: "/api/admin/*"
            actions: ["*"]
      - role: "user"
        resources:
          - pattern: "/api/user/profile"
            actions: ["read", "update"]
```

**热更新示例**：
```go
// 监听配置文件变化
watcher := fsnotify.NewWatcher()
watcher.Add("config/middleware.yaml")

for {
    select {
    case event := <-watcher.Events:
        if event.Op&fsnotify.Write == fsnotify.Write {
            // 重新加载配置
            manager.LoadConfig("config/middleware.yaml")
            log.Info("配置已热更新")
        }
    }
}
```

---

## 优先级管理

### 优先级定义

**规则**：
- 优先级范围：0-1000
- 数值越大，优先级越高（越先执行）
- 相同优先级按注册顺序执行

**标准优先级**：

| 优先级 | 中间件类型 | 示例 |
|-------|----------|-----|
| 1000 | 系统级 | Recovery |
| 900-800 | 基础设施 | RequestID, CORS, Timeout |
| 700-600 | 认证授权 | JWT, RBAC, OAuth |
| 500-400 | 安全防护 | RateLimit, CSRF, Security |
| 300-200 | 监控性能 | AccessLog, Prometheus, Cache |
| 100 | 错误处理 | ErrorHandler |

### 依赖管理

**依赖声明**：
```go
registry.Register("rbac", middleware.MiddlewareConfig{
    Name:         "rbac",
    Priority:     600,
    Handler:      RBAC(),
    Enabled:      true,
    Dependencies: []string{"jwt"}, // 依赖 JWT 认证
})
```

**依赖解析算法**：
```go
// 计算执行顺序
func (r *Registry) CalculateOrder() ([]string, error) {
    // 使用拓扑排序
    visited := make(map[string]bool)
    order := []string{}

    var visit func(name string) error
    visit = func(name string) error {
        if visited[name] {
            return nil
        }

        mw, _ := r.Get(name)

        // 先访问依赖
        for _, dep := range mw.Dependencies {
            if err := visit(dep); err != nil {
                return err
            }
        }

        visited[name] = true
        order = append(order, name)
        return nil
    }

    for name := range r.middlewares {
        if err := visit(name); err != nil {
            return nil, err
        }
    }

    return order, nil
}
```

### 优先级验证

**验证规则**：
1. Recovery 必须是第一个（优先级 1000）
2. ErrorHandler 必须是最后一个（优先级 100）
3. 认证中间件必须在授权之前
4. 限流中间件必须在业务逻辑之前

**验证代码**：
```go
func (r *Registry) Validate() error {
    middlewares := r.GetAll()

    // 检查 Recovery 是否第一个
    if middlewares[0].Name != "recovery" {
        return errors.New("recovery 必须是第一个中间件")
    }

    // 检查 ErrorHandler 是否最后一个
    if middlewares[len(middlewares)-1].Name != "error_handler" {
        return errors.New("error_handler 必须是最后一个中间件")
    }

    // 检查认证在授权之前
    jwtIndex := -1
    rbacIndex := -1
    for i, mw := range middlewares {
        if mw.Name == "jwt" {
            jwtIndex = i
        }
        if mw.Name == "rbac" {
            rbacIndex = i
        }
    }
    if jwtIndex == -1 || rbacIndex == -1 || jwtIndex > rbacIndex {
        return errors.New("JWT 必须在 RBAC 之前执行")
    }

    return nil
}
```

---

## 错误处理策略

### 统一错误处理

**错误分类**：
```go
// 错误类型
type ErrorType int

const (
    ErrorTypeValidation ErrorType = iota // 参数验证错误
    ErrorTypeAuth                       // 认证错误
    ErrorTypePermission                 // 权限错误
    ErrorTypeNotFound                   // 资源不存在
    ErrorTypeInternal                   // 内部错误
    ErrorTypePanic                      // Panic 错误
)

// 错误响应
type ErrorResponse struct {
    Code    string      `json:"code"`    // 错误码
    Message string      `json:"message"` // 错误消息
    Details interface{} `json:"details,omitempty"` // 详细信息
    RequestID string    `json:"request_id,omitempty"` // 请求ID
}
```

**错误处理器**：
```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 收集所有错误
        errs := c.Errors
        if len(errs) == 0 {
            return
        }

        // 处理第一个错误
        err := errs[0].Err

        // 根据错误类型返回响应
        switch e := err.(type) {
        case *ValidationError:
            c.JSON(400, ErrorResponse{
                Code:    "VALIDATION_ERROR",
                Message: e.Message,
                Details: e.FieldErrors,
                RequestID: c.GetString("request_id"),
            })

        case *AuthError:
            c.JSON(401, ErrorResponse{
                Code:    "AUTH_ERROR",
                Message: e.Message,
                RequestID: c.GetString("request_id"),
            })

        case *PermissionError:
            c.JSON(403, ErrorResponse{
                Code:    "PERMISSION_DENIED",
                Message: e.Message,
                RequestID: c.GetString("request_id"),
            })

        case *NotFoundError:
            c.JSON(404, ErrorResponse{
                Code:    "NOT_FOUND",
                Message: e.Message,
                RequestID: c.GetString("request_id"),
            })

        default:
            // 未知错误，记录日志
            core.LogError(c, err)

            c.JSON(500, ErrorResponse{
                Code:    "INTERNAL_ERROR",
                Message: "服务器内部错误",
                RequestID: c.GetString("request_id"),
            })
        }
    }
}
```

### Panic 恢复

**Recovery 中间件**：
```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录 panic
                core.HandlePanic(c, err)

                // 返回统一错误响应
                c.JSON(500, ErrorResponse{
                    Code:    "PANIC",
                    Message: "服务器内部错误",
                    RequestID: c.GetString("request_id"),
                })

                c.Abort()
            }
        }()

        c.Next()
    }
}
```

---

## Phase实施计划

### Phase 1: 基础设施搭建（Week 1-2）

**目标**：搭建核心框架，完成基础中间件迁移。

**任务**：
1. 创建 `internal/middleware/core/` 目录
2. 实现 `MiddlewareCore`
3. 实现 `MiddlewareRegistry`
4. 实现 `MiddlewareManager`
5. 迁移基础中间件（Recovery, RequestID, CORS）

**验收标准**：
- [ ] Core 组件单元测试覆盖率 ≥ 90%
- [ ] 基础中间件迁移完成
- [ ] 集成测试通过
- [ ] 文档更新

### Phase 2: 认证授权迁移（Week 3-4）

**目标**：迁移认证授权相关中间件。

**任务**：
1. 迁移 JWT 中间件
2. 迁移 RBAC 中间件
3. 迁移 OAuth 中间件
4. 完善单元测试

**验收标准**：
- [ ] 认证中间件迁移完成
- [ ] 测试覆盖率 ≥ 90%
- [ ] 性能测试通过
- [ ] 文档更新

### Phase 3: 安全监控迁移（Week 5-6）

**目标**：迁移安全监控相关中间件。

**任务**：
1. 迁移限流中间件
2. 迁移 CSRF 中间件
3. 迁移访问日志中间件
4. 迁移 Prometheus 中间件

**验收标准**：
- [ ] 安全中间件迁移完成
- [ ] 监控指标正常
- [ ] 性能测试通过
- [ ] 文档更新

### Phase 4: 性能优化迁移（Week 7-8）

**目标**：迁移性能优化相关中间件。

**任务**：
1. 迁移缓存中间件
2. 迁移压缩中间件
3. 迁移超时控制中间件
4. 性能测试与优化

**验收标准**：
- [ ] 性能中间件迁移完成
- [ ] 性能测试通过
- [ ] 基准测试对比完成
- [ ] 文档更新

### Phase 5: 错误处理完善（Week 9）

**目标**：完善错误处理和测试。

**任务**：
1. 实现统一错误处理
2. 完善错误码体系
3. 补充 E2E 测试
4. 性能优化验证

**验收标准**：
- [ ] 错误处理统一
- [ ] E2E 测试覆盖率 100%
- [ ] 性能基线达标
- [ ] 文档完整

---

## 快速上手示例

### 示例1: 创建静态配置中间件

**需求**：创建一个 `X-Response-Time` 响应头中间件。

**步骤**：

#### 1. 定义配置结构体

```go
// internal/middleware/base/response_time.go
package base

type ResponseTimeConfig struct {
    Enabled bool   `json:"enabled" yaml:"enabled"`
    Header  string `json:"header" yaml:"header"`
}

func DefaultResponseTimeConfig() *ResponseTimeConfig {
    return &ResponseTimeConfig{
        Enabled: true,
        Header:  "X-Response-Time",
    }
}
```

#### 2. 实现中间件

```go
type ResponseTimeMiddleware struct {
    core   middleware.MiddlewareCore
    config *ResponseTimeConfig
}

func NewResponseTimeMiddleware(core middleware.MiddlewareCore) *ResponseTimeMiddleware {
    return &ResponseTimeMiddleware{
        core:   core,
        config: DefaultResponseTimeConfig(),
    }
}

func (rt *ResponseTimeMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !rt.config.Enabled {
            c.Next()
            return
        }

        start := time.Now()
        c.Next()
        duration := time.Since(start)

        // 设置响应头
        c.Header(rt.config.Header, duration.String())

        // 记录指标
        rt.core.RecordMetric("response_time", duration.Seconds())
    }
}

func (rt *ResponseTimeMiddleware) LoadConfig(config map[string]interface{}) error {
    rt.config.Enabled = config["enabled"].(bool)
    rt.config.Header = config["header"].(string)
    return nil
}

func (rt *ResponseTimeMiddleware) ValidateConfig() error {
    if rt.config.Header == "" {
        return errors.New("header 不能为空")
    }
    return nil
}
```

#### 3. 注册中间件

```go
// internal/middleware/registry.go
registry.Register("response_time", middleware.MiddlewareConfig{
    Name:     "response_time",
    Priority: 350,
    Handler:  responseTime.Handler(),
    Enabled:  true,
    Static:   true,
})
```

#### 4. 编写单元测试

```go
// internal/middleware/base/response_time_test.go
func TestResponseTimeMiddleware(t *testing.T) {
    core := mock.NewMockCore()
    mw := NewResponseTimeMiddleware(core)

    // 测试禁用状态
    mw.config.Enabled = false
    w := performRequest(mw.Handler())
    assert.Equal(t, "", w.Header().Get("X-Response-Time"))

    // 测试启用状态
    mw.config.Enabled = true
    w = performRequest(mw.Handler())
    assert.NotEqual(t, "", w.Header().Get("X-Response-Time"))
}
```

### 示例2: 创建动态配置中间件

**需求**：创建一个 IP 黑名单中间件，支持动态更新。

**步骤**：

#### 1. 定义配置结构体

```go
// internal/middleware/security/ip_blacklist.go
package security

type IPBlacklistConfig struct {
    Enabled   bool     `json:"enabled" yaml:"enabled"`
    Blacklist []string `json:"blacklist" yaml:"blacklist"`
    Driver    string   `json:"driver" yaml:"driver"`    // memory, redis
    TTL       int      `json:"ttl" yaml:"ttl"`          // Redis TTL (秒)
}

type IPBlacklistMiddleware struct {
    core   middleware.MiddlewareCore
    config *IPBlacklistConfig
    store  BlacklistStore
}

type BlacklistStore interface {
    IsBlocked(ip string) bool
    Add(ip string) error
    Remove(ip string) error
}
```

#### 2. 实现中间件

```go
func NewIPBlacklistMiddleware(core middleware.MiddlewareCore, store BlacklistStore) *IPBlacklistMiddleware {
    return &IPBlacklistMiddleware{
        core:   core,
        config: DefaultIPBlacklistConfig(),
        store:  store,
    }
}

func (ib *IPBlacklistMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !ib.config.Enabled {
            c.Next()
            return
        }

        // 获取客户端IP
        ip := c.ClientIP()

        // 检查是否在黑名单
        if ib.store.IsBlocked(ip) {
            ib.core.HandleError(c, errors.New("IP blocked"))
            c.JSON(403, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func (ib *IPBlacklistMiddleware) LoadConfig(config map[string]interface{}) error {
    ib.config.Enabled = config["enabled"].(bool)
    ib.config.Driver = config["driver"].(string)

    // 加载黑名单
    blacklist := config["blacklist"].([]interface{})
    for _, ip := range blacklist {
        ib.store.Add(ip.(string))
    }

    return nil
}
```

#### 3. 实现热更新

```go
// 监听 Redis Pub/Sub
func (ib *IPBlacklistMiddleware) WatchChanges() {
    if ib.config.Driver != "redis" {
        return
    }

    pubsub := redisClient.Subscribe(context.Background(), "ip_blacklist:updates")

    for msg := range pubsub.Channel() {
        var update struct {
            Action string `json:"action"` // add, remove
            IP     string `json:"ip"`
        }

        json.Unmarshal([]byte(msg.Payload), &update)

        switch update.Action {
        case "add":
            ib.store.Add(update.IP)
        case "remove":
            ib.store.Remove(update.IP)
        }
    }
}
```

#### 4. 编写单元测试

```go
// internal/middleware/security/ip_blacklist_test.go
func TestIPBlacklistMiddleware(t *testing.T) {
    core := mock.NewMockCore()
    store := mock.NewMockStore()
    mw := NewIPBlacklistMiddleware(core, store)

    // 添加黑名单
    store.Add("192.168.1.100")

    // 测试被阻止的 IP
    w := performRequestWithIP(mw.Handler(), "192.168.1.100")
    assert.Equal(t, 403, w.Code)

    // 测试正常 IP
    w = performRequestWithIP(mw.Handler(), "192.168.1.101")
    assert.Equal(t, 200, w.Code)
}
```

---

## Q&A

### Q1: 如何调试中间件执行顺序？

**A**：使用 `MiddlewareManager.GetAllStatus()` 查看当前顺序和优先级。

```go
status := manager.GetAllStatus()
for name, st := range status {
    fmt.Printf("%s: Priority=%d, Enabled=%v\n", name, st.Priority, st.Enabled)
}
```

### Q2: 如何处理中间件之间的依赖冲突？

**A**：
1. 检查依赖声明是否正确
2. 使用 `ValidateDependencies()` 验证
3. 调整优先级或依赖关系

### Q3: 动态配置是否会影响性能？

**A**：
- 配置加载：启动时一次性加载
- 热更新：异步处理，不阻塞请求
- 读写分离：配置读多写少，使用读写锁优化

### Q4: 如何回滚到旧版中间件？

**A**：
1. 保留旧版代码在 `internal/middleware/legacy/`
2. 使用 feature flag 控制
3. 逐步迁移，保证平滑过渡

### Q5: 如何监控中间件性能？

**A**：
1. 使用 Prometheus 指标
2. 查看 `MiddlewareCore.GetMetrics()`
3. 集成 APM 工具（如 New Relic）

---

## 参考资料

### 设计文档
- [中间件总体设计](../design/middleware/中间件总体设计.md)
- [核心接口设计](../design/middleware/core-interface-design.md)

### 实施指南
- [Phase 1 实施计划](../plans/middleware-refactor-phase1.md)
- [测试指南](../testing/middleware-testing-guide.md)

### 代码示例
- [静态中间件示例](../examples/static_example.go)
- [动态中间件示例](../examples/dynamic_example.go)

### 外部资源
- [Gin 中间件文档](https://gin-gonic.com/docs/examples/custom-http-method/)
- [Go 最佳实践](https://go.dev/doc/effective_go)
- [微服务模式](https://microservices.io/patterns/microservices.html)

---

**文档版本**: v1.0
**最后更新**: 2026-01-27
**维护者**: 架构师团队
