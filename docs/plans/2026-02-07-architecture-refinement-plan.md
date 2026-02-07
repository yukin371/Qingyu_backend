# Architecture Refinement Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 分阶段修复架构审查报告中的5个问题，提升代码内聚性、降低耦合度，优化系统架构健康度。

**Architecture:** 采用渐进式重构策略，按优先级(P0→P1→P2)分阶段执行，每阶段完成后进行测试验证和文档更新。

**Tech Stack:** Go 1.x, Gin, MongoDB, Redis, Milvus, MinIO

---

## 阶段划分总览

| 阶段 | 问题 | 优先级 | 预计工作量 |
|------|------|--------|------------|
| Phase 1 | 中间件跨层依赖 | P0 | 2-3天 |
| Phase 2 | 服务初始化顺序依赖 | P1 | 3-4天 |
| Phase 3 | shared模块职责过重 | P1 | 5-7天 |
| Phase 4 | WriterService耦合度高 | P2 | 4-5天 |
| Phase 5 | 事件总线持久化 | P2 | 3-4天 |

---

## Phase 1: 修复中间件跨层依赖 [P0]

### 目标
将 `pkg/middleware/quota.go` 中间件与业务服务层解耦，引入抽象层。

### 当前问题
```go
// pkg/middleware/quota.go:15
type QuotaMiddleware struct {
    quotaService *aiService.QuotaService  // 直接依赖业务服务
}
```

### 解决方案
引入配额检查接口，中间件依赖接口而非具体实现。

### Task 1.1: 创建配额检查接口

**Files:**
- Create: `pkg/quota/interface.go`
- Create: `pkg/quota/errors.go`

**Step 1: 创建接口文件**

```bash
# 创建目录
mkdir -p pkg/quota
```

**Step 2: 编写接口定义**

```go
// pkg/quota/interface.go
package quota

import "context"

// Checker 配额检查接口
// 中间件依赖此接口，而非具体的服务实现
type Checker interface {
    // Check 检查用户是否有足够配额执行操作
    Check(ctx context.Context, userID string, operation string, estimatedAmount int) error

    // Consume 消耗配额
    Consume(ctx context.Context, userID string, operation string, amount int) error

    // GetRemaining 获取剩余配额
    GetRemaining(ctx context.Context, userID string) (int, error)
}

// OperationType 操作类型常量
const (
    OperationChat     = "chat"      // AI对话
    OperationGenerate = "generate"  // 内容生成
    OperationEdit     = "edit"      // 内容编辑
)
```

**Step 3: 编写错误定义**

```go
// pkg/quota/errors.go
package quota

import "errors"

var (
    // ErrQuotaExhausted 配额已用尽
    ErrQuotaExhausted = errors.New("quota exhausted")

    // ErrQuotaSuspended 配额已暂停
    ErrQuotaSuspended = errors.New("quota suspended")

    // ErrInsufficientQuota 配额不足
    ErrInsufficientQuota = errors.New("insufficient quota")

    // ErrUserNotFound 用户不存在
    ErrUserNotFound = errors.New("user not found")
)
```

**Step 4: 提交**

```bash
git add pkg/quota/interface.go pkg/quota/errors.go
git commit -m "feat(quota): add quota checker interface and error definitions"
```

---

### Task 1.2: 实现QuotaService适配器

**Files:**
- Create: `service/ai/quota_adapter.go`
- Test: `service/ai/quota_adapter_test.go`

**Step 1: 编写适配器实现**

```go
// service/ai/quota_adapter.go
package ai

import (
    "context"
    "fmt"

    "Qingyu_backend/pkg/quota"
)

// QuotaCheckerAdapter 将QuotaService适配为quota.Checker接口
type QuotaCheckerAdapter struct {
    quotaService *QuotaService
}

// NewQuotaCheckerAdapter 创建配额检查适配器
func NewQuotaCheckerAdapter(quotaService *QuotaService) quota.Checker {
    return &QuotaCheckerAdapter{
        quotaService: quotaService,
    }
}

// Check 实现quota.Checker接口
func (a *QuotaCheckerAdapter) Check(ctx context.Context, userID string, operation string, estimatedAmount int) error {
    return a.quotaService.CheckQuota(ctx, userID, estimatedAmount)
}

// Consume 实现quota.Checker接口
func (a *QuotaCheckerAdapter) Consume(ctx context.Context, userID string, operation string, amount int) error {
    return a.quotaService.ConsumeQuota(ctx, userID, amount)
}

// GetRemaining 实现quota.Checker接口
func (a *QuotaCheckerAdapter) GetRemaining(ctx context.Context, userID string) (int, error) {
    quota, err := a.quotaService.GetUserQuota(ctx, userID)
    if err != nil {
        return 0, fmt.Errorf("get user quota: %w", err)
    }
    return quota.RemainingTokens, nil
}

// 确保实现了接口
var _ quota.Checker = (*QuotaCheckerAdapter)(nil)
```

**Step 2: 编写单元测试**

```go
// service/ai/quota_adapter_test.go
package ai

import (
    "context"
    "testing"

    "Qingyu_backend/pkg/quota"
)

func TestQuotaCheckerAdapter(t *testing.T) {
    // 创建mock QuotaService
    mockService := &MockQuotaService{
        remaining: 1000,
    }

    // 创建适配器
    checker := NewQuotaCheckerAdapter(mockService)

    // 测试Check
    t.Run("Check success", func(t *testing.T) {
        err := checker.Check(context.Background(), "user123", quota.OperationChat, 100)
        if err != nil {
            t.Errorf("expected no error, got %v", err)
        }
    })

    // 测试GetRemaining
    t.Run("GetRemaining", func(t *testing.T) {
        remaining, err := checker.GetRemaining(context.Background(), "user123")
        if err != nil {
            t.Errorf("expected no error, got %v", err)
        }
        if remaining != 1000 {
            t.Errorf("expected 1000, got %d", remaining)
        }
    })
}

// MockQuotaService 用于测试
type MockQuotaService struct {
    remaining int
}

func (m *MockQuotaService) CheckQuota(ctx context.Context, userID string, amount int) error {
    if amount > m.remaining {
        return quota.ErrInsufficientQuota
    }
    return nil
}

func (m *MockQuotaService) ConsumeQuota(ctx context.Context, userID string, amount int) error {
    m.remaining -= amount
    return nil
}

func (m *MockQuotaService) GetUserQuota(ctx context.Context, userID string) (*UserQuota, error) {
    return &UserQuota{RemainingTokens: m.remaining}, nil
}
```

**Step 3: 运行测试**

```bash
go test ./service/ai/quota_adapter_test.go -v
```

**Step 4: 提交**

```bash
git add service/ai/quota_adapter.go service/ai/quota_adapter_test.go
git commit -m "feat(ai): add QuotaCheckerAdapter to bridge service and interface"
```

---

### Task 1.3: 重构中间件使用接口

**Files:**
- Modify: `pkg/middleware/quota.go`

**Step 1: 重写中间件**

```go
// pkg/middleware/quota.go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "Qingyu_backend/api/v1/shared"
    "Qingyu_backend/pkg/quota"
)

// QuotaMiddleware 配额中间件
type QuotaMiddleware struct {
    checker quota.Checker  // 依赖接口而非具体实现
}

// NewQuotaMiddleware 创建配额中间件
func NewQuotaMiddleware(checker quota.Checker) *QuotaMiddleware {
    return &QuotaMiddleware{
        checker: checker,
    }
}

// CheckQuota 检查配额中间件
func (m *QuotaMiddleware) CheckQuota(estimatedAmount int) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取用户ID
        userID, exists := c.Get("user_id")
        if !exists {
            shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
            c.Abort()
            return
        }

        // 使用接口检查配额
        err := m.checker.Check(c.Request.Context(), userID.(string), "", estimatedAmount)
        if err != nil {
            m.handleError(c, err)
            return
        }

        c.Next()
    }
}

// handleError 处理配额错误
func (m *QuotaMiddleware) handleError(c *gin.Context, err error) {
    switch err {
    case quota.ErrQuotaExhausted:
        shared.Error(c, http.StatusTooManyRequests, "配额已用尽", "您的AI配额已用尽，请明天再试或升级会员")
    case quota.ErrQuotaSuspended:
        shared.Error(c, http.StatusForbidden, "配额已暂停", "您的AI配额已被暂停")
    case quota.ErrInsufficientQuota:
        shared.Error(c, http.StatusTooManyRequests, "配额不足", "您的AI配额不足以完成此操作")
    case quota.ErrUserNotFound:
        shared.Error(c, http.StatusNotFound, "用户不存在", "未找到用户信息")
    default:
        shared.Error(c, http.StatusInternalServerError, "配额检查失败", err.Error())
    }
    c.Abort()
}

// QuotaCheckMiddleware 简化版配额检查中间件（兼容旧代码）
// Deprecated: 建议使用 NewQuotaMiddleware + CheckQuota
func QuotaCheckMiddleware(checker quota.Checker) gin.HandlerFunc {
    middleware := NewQuotaMiddleware(checker)
    return middleware.CheckQuota(1000)  // 默认1000 tokens
}

// LightQuotaCheckMiddleware 轻量级配额检查
func LightQuotaCheckMiddleware(checker quota.Checker) gin.HandlerFunc {
    middleware := NewQuotaMiddleware(checker)
    return middleware.CheckQuota(300)  // 300 tokens
}

// HeavyQuotaCheckMiddleware 重量级配额检查
func HeavyQuotaCheckMiddleware(checker quota.Checker) gin.HandlerFunc {
    middleware := NewQuotaMiddleware(checker)
    return middleware.CheckQuota(3000)  // 3000 tokens
}
```

**Step 2: 提交**

```bash
git add pkg/middleware/quota.go
git commit -m "refactor(middleware): quota middleware now depends on interface"
```

---

### Task 1.4: 更新ServiceContainer注册

**Files:**
- Modify: `service/container/service_container.go`
- Modify: `router/ai/ai_router.go` (或使用配额中间件的路由文件)

**Step 1: 在ServiceContainer中注册配额检查器**

```go
// 在 service_container.go 中添加字段
type ServiceContainer struct {
    // ... 现有字段
    quotaChecker quota.Checker  // 新增：配额检查接口
}

// 添加获取方法
func (c *ServiceContainer) GetQuotaChecker() (quota.Checker, error) {
    if c.quotaChecker == nil {
        return nil, fmt.Errorf("QuotaChecker未初始化")
    }
    return c.quotaChecker, nil
}
```

**Step 2: 在SetupDefaultServices中初始化**

```go
// 在 SetupDefaultServices 方法中，QuotaService初始化后添加
// ============ AI服务初始化 ============
quotaRepo := c.repositoryFactory.CreateQuotaRepository()
c.quotaService = aiService.NewQuotaService(quotaRepo)

// 创建配额检查适配器
c.quotaChecker = aiService.NewQuotaCheckerAdapter(c.quotaService)
```

**Step 3: 更新路由使用新接口**

```go
// 在路由初始化时使用接口
quotaChecker, _ := container.GetQuotaChecker()
aiRouter.Use(middleware.QuotaCheckMiddleware(quotaChecker))
```

**Step 4: 提交**

```bash
git add service/container/service_container.go
git commit -m "refactor(container): register QuotaChecker interface"
```

---

### Task 1.5: 验证和测试

**Step 1: 编写集成测试**

```go
// pkg/middleware/quota_integration_test.go
package middleware

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"

    "Qingyu_backend/pkg/quota"
)

func TestQuotaMiddlewareIntegration(t *testing.T) {
    // 设置Gin为测试模式
    gin.SetMode(gin.TestMode)

    // 创建mock checker
    mockChecker := &MockChecker{}

    // 创建中间件
    middleware := NewQuotaMiddleware(mockChecker)

    // 创建测试路由
    router := gin.New()
    router.Use(middleware.CheckQuota(100))
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // 测试正常情况
    t.Run("valid quota", func(t *testing.T) {
        mockChecker.shouldFail = false
        req := httptest.NewRequest("GET", "/test", nil)
        req.Header.Set("Authorization", "Bearer valid_token")
        w := httptest.NewRecorder()

        router.ServeHTTP(w, req)

        if w.Code != 200 {
            t.Errorf("expected 200, got %d", w.Code)
        }
    })

    // 测试配额不足
    t.Run("insufficient quota", func(t *testing.T) {
        mockChecker.shouldFail = true
        mockChecker.failWithError = quota.ErrInsufficientQuota
        req := httptest.NewRequest("GET", "/test", nil)
        req.Header.Set("Authorization", "Bearer valid_token")
        w := httptest.NewRecorder()

        router.ServeHTTP(w, req)

        if w.Code != http.StatusTooManyRequests {
            t.Errorf("expected 429, got %d", w.Code)
        }
    })
}

// MockChecker 用于测试
type MockChecker struct {
    shouldFail   bool
    failWithError error
}

func (m *MockChecker) Check(ctx context.Context, userID string, operation string, amount int) error {
    if m.shouldFail {
        return m.failWithError
    }
    return nil
}

func (m *MockChecker) Consume(ctx context.Context, userID string, operation string, amount int) error {
    return nil
}

func (m *MockChecker) GetRemaining(ctx context.Context, userID string) (int, error) {
    return 1000, nil
}
```

**Step 2: 运行所有测试**

```bash
go test ./pkg/middleware/... ./service/ai/... -v
```

**Step 3: 启动服务手动验证**

```bash
# 启动服务
go run cmd/main.go

# 测试API（需要有效token）
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/ai/chat
```

**Step 4: 提交测试代码**

```bash
git add pkg/middleware/quota_integration_test.go
git commit -m "test(middleware): add integration tests for quota middleware"
```

---

### Phase 1 验收标准

- [ ] 中间件不再直接依赖 `aiService.QuotaService`
- [ ] 所有测试通过（单元测试 + 集成测试）
- [ ] 服务启动正常，AI接口可用
- [ ] 架构文档已更新

---

## Phase 2: 服务初始化顺序依赖 [P1]

### 目标
实现服务依赖的显式声明，容器自动解析初始化顺序，消除隐式依赖。

### 当前问题
```go
// service/container/service_container.go:804
// AIService需要ProjectService来构建上下文（隐式依赖）
c.aiService = aiService.NewServiceWithDependencies(c.projectService)
```

### 解决方案
引入服务注册表，显式声明依赖关系，使用拓扑排序解析初始化顺序。

### Task 2.1: 定义服务注册表结构

**Files:**
- Create: `service/container/registry.go`

**Step 1: 创建注册表**

```go
// service/container/registry.go
package container

import (
    "fmt"
    "sort"
)

// ServiceDef 服务定义
type ServiceDef struct {
    // Name 服务名称（唯一标识）
    Name string

    // Factory 服务工厂函数
    Factory func(*ServiceContainer) (interface{}, error)

    // Dependencies 依赖的服务名称列表
    Dependencies []string

    // Optional 是否为可选服务（初始化失败不阻塞）
    Optional bool

    // Priority 优先级（数值越小越先初始化，同级别按依赖排序）
    Priority int
}

// ServiceRegistry 服务注册表
type ServiceRegistry struct {
    services map[string]*ServiceDef
}

// NewServiceRegistry 创建服务注册表
func NewServiceRegistry() *ServiceRegistry {
    return &ServiceRegistry{
        services: make(map[string]*ServiceDef),
    }
}

// Register 注册服务定义
func (r *ServiceRegistry) Register(def *ServiceDef) error {
    if def.Name == "" {
        return fmt.Errorf("服务名称不能为空")
    }

    if _, exists := r.services[def.Name]; exists {
        return fmt.Errorf("服务 %s 已存在", def.Name)
    }

    r.services[def.Name] = def
    return nil
}

// Get 获取服务定义
func (r *ServiceRegistry) Get(name string) (*ServiceDef, bool) {
    def, exists := r.services[name]
    return def, exists
}

// GetInitializationOrder 获取服务初始化顺序（拓扑排序）
func (r *ServiceRegistry) GetInitializationOrder() ([]string, error) {
    // 使用Kahn算法进行拓扑排序
    inDegree := make(map[string]int)
    adjacency := make(map[string][]string)

    // 初始化
    for name := range r.services {
        inDegree[name] = 0
        adjacency[name] = []string{}
    }

    // 计算入度和邻接表
    for name, def := range r.services {
        for _, dep := range def.Dependencies {
            if _, exists := r.services[dep]; !exists {
                return nil, fmt.Errorf("服务 %s 依赖的 %s 不存在", name, dep)
            }
            adjacency[dep] = append(adjacency[dep], name)
            inDegree[name]++
        }
    }

    // 按优先级分组
    priorityGroups := make(map[int][]string)
    for name := range r.services {
        priority := r.services[name].Priority
        priorityGroups[priority] = append(priorityGroups[priority], name)
    }

    // 获取所有优先级并排序
    priorities := make([]int, 0, len(priorityGroups))
    for p := range priorityGroups {
        priorities = append(priorities, p)
    }
    sort.Ints(priorities)

    // 按优先级进行拓扑排序
    var order []string
    for _, priority := range priorities {
        names := priorityGroups[priority]
        groupOrder, err := r.topologicalSortGroup(names, inDegree, adjacency)
        if err != nil {
            return nil, err
        }
        order = append(order, groupOrder...)
    }

    return order, nil
}

// topologicalSortGroup 对一组服务进行拓扑排序
func (r *ServiceRegistry) topologicalSortGroup(
    names []string,
    inDegree map[string]int,
    adjacency map[string][]string,
) ([]string, error) {
    // 创建入度副本
    inDegreeCopy := make(map[string]int)
    for _, name := range names {
        inDegreeCopy[name] = inDegree[name]
    }

    var queue []string
    for _, name := range names {
        if inDegreeCopy[name] == 0 {
            queue = append(queue, name)
        }
    }

    var order []string
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        order = append(order, current)

        for _, neighbor := range adjacency[current] {
            // 只处理同优先级的
            if inDegreeCopy[neighbor] > 0 {
                inDegreeCopy[neighbor]--
                if inDegreeCopy[neighbor] == 0 {
                    queue = append(queue, neighbor)
                }
            }
        }
    }

    if len(order) != len(names) {
        return nil, fmt.Errorf("检测到循环依赖")
    }

    return order, nil
}

// Validate 验证注册表
func (r *ServiceRegistry) Validate() error {
    // 检查循环依赖
    visited := make(map[string]bool)
    recStack := make(map[string]bool)

    for name := range r.services {
        if !visited[name] {
            if err := r.detectCycle(name, visited, recStack); err != nil {
                return err
            }
        }
    }

    // 验证依赖存在
    for name, def := range r.services {
        for _, dep := range def.Dependencies {
            if _, exists := r.services[dep]; !exists {
                return fmt.Errorf("服务 %s 依赖的 %s 不存在", name, dep)
            }
        }
    }

    return nil
}

// detectCycle 检测循环依赖
func (r *ServiceRegistry) detectCycle(name string, visited, recStack map[string]bool) error {
    visited[name] = true
    recStack[name] = true

    def := r.services[name]
    for _, dep := range def.Dependencies {
        if depDef, exists := r.services[dep]; exists {
            if !visited[dep] {
                if err := r.detectCycle(dep, visited, recStack); err != nil {
                    return err
                }
            } else if recStack[dep] {
                return fmt.Errorf("检测到循环依赖: %s -> %s", name, dep)
            }
            _ = depDef // 使用变量避免unused警告
        }
    }

    recStack[name] = false
    return nil
}
```

**Step 2: 提交**

```bash
git add service/container/registry.go
git commit -m "feat(container): add ServiceRegistry with dependency resolution"
```

---

### Task 2.2: 重构ServiceContainer使用注册表

**Files:**
- Modify: `service/container/service_container.go`

**Step 1: 添加注册表字段**

```go
// 在 ServiceContainer 结构体中添加
type ServiceContainer struct {
    // ... 现有字段
    registry *ServiceRegistry
}
```

**Step 2: 修改NewServiceContainer**

```go
func NewServiceContainer() *ServiceContainer {
    return &ServiceContainer{
        services:       make(map[string]serviceInterfaces.BaseService),
        serviceMetrics: make(map[string]*metrics.ServiceMetrics),
        initialized:    false,
        eventBus:       base.NewSimpleEventBus(),
        registry:       NewServiceRegistry(),  // 新增
    }
}
```

**Step 3: 创建服务注册方法**

```go
// RegisterServiceDefinition 注册服务定义
func (c *ServiceContainer) RegisterServiceDefinition(def *ServiceDef) error {
    return c.registry.Register(def)
}

// InitializeServices 按依赖顺序初始化服务
func (c *ServiceContainer) InitializeServices(ctx context.Context) error {
    // 验证注册表
    if err := c.registry.Validate(); err != nil {
        return fmt.Errorf("服务注册表验证失败: %w", err)
    }

    // 获取初始化顺序
    order, err := c.registry.GetInitializationOrder()
    if err != nil {
        return fmt.Errorf("获取初始化顺序失败: %w", err)
    }

    // 按顺序初始化
    for _, name := range order {
        def, _ := c.registry.Get(name)

        instance, err := def.Factory(c)
        if err != nil {
            if def.Optional {
                fmt.Printf("警告: 可选服务 %s 初始化失败: %v\n", name, err)
                continue
            }
            return fmt.Errorf("初始化服务 %s 失败: %w", name, err)
        }

        // 存储服务实例
        if err := c.setServiceInstance(name, instance); err != nil {
            return fmt.Errorf("存储服务 %s 失败: %w", name, err)
        }

        fmt.Printf("✓ 服务 %s 初始化完成\n", name)
    }

    return nil
}

// setServiceInstance 存储服务实例到对应字段
func (c *ServiceContainer) setServiceInstance(name string, instance interface{}) error {
    switch name {
    case "UserService":
        c.userService = instance.(userInterface.UserService)
    case "AuthService":
        c.authService = instance.(auth.AuthService)
    case "AIService":
        c.aiService = instance.(*aiService.Service)
    case "QuotaService":
        c.quotaService = instance.(*aiService.QuotaService)
    case "QuotaChecker":
        c.quotaChecker = instance.(quota.Checker)
    // ... 其他服务映射
    default:
        // 对于未显式映射的服务，尝试作为BaseService存储
        if bs, ok := instance.(serviceInterfaces.BaseService); ok {
            c.services[name] = bs
        }
    }
    return nil
}
```

**Step 4: 提交**

```bash
git add service/container/service_container.go
git commit -m "refactor(container): integrate ServiceRegistry for dependency resolution"
```

---

### Task 2.3: 将现有服务迁移到注册表

**Files:**
- Create: `service/container/service_registration.go`

**Step 1: 创建服务注册文件**

```go
// service/container/service_registration.go
package container

import (
    "Qingyu_backend/service/ai"
    "Qingyu_backend/service/user"
    "Qingyu_backend/service/shared/auth"
    // ... 其他服务导入
)

// RegisterDefaultServices 注册所有默认服务到注册表
func (c *ServiceContainer) RegisterDefaultServices() error {
    // ============ 基础设施服务 (优先级: 0) ============
    c.RegisterServiceDefinition(&ServiceDef{
        Name:     "EventBus",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            return c.eventBus, nil
        },
        Dependencies: []string{},
        Priority:     0,
    })

    // ============ 核心服务 (优先级: 10) ============
    c.RegisterServiceDefinition(&ServiceDef{
        Name: "UserService",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            userRepo := c.repositoryFactory.CreateUserRepository()
            authRepo := c.repositoryFactory.CreateAuthRepository()
            return userService.NewUserService(userRepo, authRepo), nil
        },
        Dependencies: []string{},
        Priority:     10,
    })

    c.RegisterServiceDefinition(&ServiceDef{
        Name: "AuthService",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            authRepo := c.repositoryFactory.CreateAuthRepository()
            oauthRepo := c.repositoryFactory.CreateOAuthRepository()

            var redisAdapter interface{}
            if c.redisClient != nil {
                redisAdapter = auth.NewRedisAdapter(c.redisClient)
            } else {
                redisAdapter = auth.NewInMemoryTokenBlacklist()
            }

            jwtService := auth.NewJWTService(
                config.GetJWTConfigEnhanced(),
                redisAdapter.(auth.RedisClient),
            )
            roleService := auth.NewRoleService(authRepo)

            cacheClient, ok := redisAdapter.(auth.CacheClient)
            if !ok {
                return nil, fmt.Errorf("redisAdapter does not implement CacheClient")
            }

            permissionService := auth.NewPermissionService(authRepo, cacheClient, nil)
            sessionService := auth.NewSessionService(cacheClient)

            return auth.NewAuthService(
                jwtService,
                roleService,
                permissionService,
                authRepo,
                oauthRepo,
                c.userService,
                sessionService,
            ), nil
        },
        Dependencies: []string{"UserService"},
        Priority:     10,
    })

    // ============ 业务服务 (优先级: 20) ============
    c.RegisterServiceDefinition(&ServiceDef{
        Name: "ProjectService",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            projectRepo := c.repositoryFactory.CreateProjectRepository()
            return projectService.NewProjectService(
                projectRepo,
                c.eventBus,
            ), nil
        },
        Dependencies: []string{"EventBus"},
        Priority:     20,
    })

    c.RegisterServiceDefinition(&ServiceDef{
        Name: "AIService",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            return aiService.NewServiceWithDependencies(c.projectService), nil
        },
        Dependencies: []string{"ProjectService"},
        Priority:     20,
    })

    c.RegisterServiceDefinition(&ServiceDef{
        Name: "QuotaService",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            quotaRepo := c.repositoryFactory.CreateQuotaRepository()
            return aiService.NewQuotaService(quotaRepo), nil
        },
        Dependencies: []string{},
        Priority:     20,
    })

    c.RegisterServiceDefinition(&ServiceDef{
        Name: "QuotaChecker",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            return aiService.NewQuotaCheckerAdapter(c.quotaService), nil
        },
        Dependencies: []string{"QuotaService"},
        Priority:     25,
    })

    // ============ 可选服务 (优先级: 30) ============
    c.RegisterServiceDefinition(&ServiceDef{
        Name: "OAuthService",
        Factory: func(c *ServiceContainer) (interface{}, error) {
            oauthConfigMgr := config.NewOAuthConfigManager()
            oauthConfigMgr.LoadFromEnv()
            if config.GlobalConfig != nil {
                oauthConfigMgr.LoadFromConfig(config.GlobalConfig)
            }

            oauthConfigs := oauthConfigMgr.GetConfigs()
            if len(oauthConfigs) == 0 {
                return nil, nil // 无配置，不创建
            }

            oauthRepo := c.repositoryFactory.CreateOAuthRepository()
            logger, err := zap.NewProduction()
            if err != nil {
                return nil, err
            }

            return auth.NewOAuthService(logger, oauthRepo, oauthConfigs)
        },
        Dependencies: []string{"AuthService"},
        Priority:     30,
        Optional:     true,
    })

    // ... 其他服务注册

    return nil
}
```

**Step 2: 提交**

```bash
git add service/container/service_registration.go
git commit -m "feat(container): add service registration with explicit dependencies"
```

---

### Task 2.4: 更新初始化流程

**Files:**
- Modify: `service/container/service_container.go`

**Step 1: 修改Initialize方法**

```go
// Initialize 初始化所有服务（使用注册表）
func (c *ServiceContainer) Initialize(ctx context.Context) error {
    if c.initialized {
        return nil
    }

    // 1. 初始化MongoDB
    if err := c.initMongoDB(); err != nil {
        return fmt.Errorf("MongoDB初始化失败: %w", err)
    }

    // 2. 创建Repository工厂
    c.repositoryFactory = mongodb.NewMongoRepositoryFactoryWithClient(
        c.mongoClient,
        c.mongoDB,
    )

    // 3. 初始化Redis
    if err := c.initRedis(); err != nil {
        fmt.Printf("警告: Redis客户端初始化失败: %v\n", err)
    }

    // 4. 健康检查
    if err := c.repositoryFactory.Health(ctx); err != nil {
        return fmt.Errorf("Repository工厂健康检查失败: %w", err)
    }

    // 5. 注册所有服务定义
    if err := c.RegisterDefaultServices(); err != nil {
        return fmt.Errorf("注册服务定义失败: %w", err)
    }

    // 6. 按依赖顺序初始化服务
    if err := c.InitializeServices(ctx); err != nil {
        return fmt.Errorf("初始化服务失败: %w", err)
    }

    // 7. 预热缓存
    if err := c.warmUpCache(ctx); err != nil {
        fmt.Printf("警告: 缓存预热失败: %v\n", err)
    }

    c.initialized = true
    return nil
}
```

**Step 2: 添加依赖可视化方法**

```go
// PrintDependencyGraph 打印服务依赖图（用于调试）
func (c *ServiceContainer) PrintDependencyGraph() string {
    order, err := c.registry.GetInitializationOrder()
    if err != nil {
        return fmt.Sprintf("错误: %v", err)
    }

    var result string
    result += "服务初始化顺序:\n"
    result += "==================\n"

    for i, name := range order {
        def, _ := c.registry.Get(name)
        result += fmt.Sprintf("%d. %s (优先级: %d)\n", i+1, name, def.Priority)
        if len(def.Dependencies) > 0 {
            result += fmt.Sprintf("   依赖: %v\n", def.Dependencies)
        }
        if def.Optional {
            result += "   [可选]\n"
        }
    }

    return result
}
```

**Step 3: 提交**

```bash
git add service/container/service_container.go
git commit -m "refactor(container): update initialization to use registry-based approach"
```

---

### Task 2.5: 测试和验证

**Files:**
- Create: `service/container/registry_test.go`

**Step 1: 编写测试**

```go
// service/container/registry_test.go
package container

import (
    "testing"
)

func TestServiceRegistry(t *testing.T) {
    registry := NewServiceRegistry()

    // 测试注册
    t.Run("Register services", func(t *testing.T) {
        err := registry.Register(&ServiceDef{
            Name:         "ServiceA",
            Dependencies: []string{},
            Priority:     0,
        })
        if err != nil {
            t.Errorf("expected no error, got %v", err)
        }

        // 测试重复注册
        err = registry.Register(&ServiceDef{
            Name:         "ServiceA",
            Dependencies: []string{},
            Priority:     0,
        })
        if err == nil {
            t.Error("expected error for duplicate registration")
        }
    })

    // 测试依赖解析
    t.Run("Dependency resolution", func(t *testing.T) {
        registry := NewServiceRegistry()

        // A -> B -> C
        registry.Register(&ServiceDef{Name: "C", Dependencies: []string{}, Priority: 0})
        registry.Register(&ServiceDef{Name: "B", Dependencies: []string{"C"}, Priority: 0})
        registry.Register(&ServiceDef{Name: "A", Dependencies: []string{"B"}, Priority: 0})

        order, err := registry.GetInitializationOrder()
        if err != nil {
            t.Errorf("expected no error, got %v", err)
        }

        // 验证顺序: C -> B -> A
        if order[0] != "C" || order[1] != "B" || order[2] != "A" {
            t.Errorf("wrong order: %v", order)
        }
    })

    // 测试循环依赖检测
    t.Run("Circular dependency detection", func(t *testing.T) {
        registry := NewServiceRegistry()

        // A -> B -> A (循环)
        registry.Register(&ServiceDef{Name: "A", Dependencies: []string{"B"}, Priority: 0})
        registry.Register(&ServiceDef{Name: "B", Dependencies: []string{"A"}, Priority: 0})

        err := registry.Validate()
        if err == nil {
            t.Error("expected error for circular dependency")
        }
    })

    // 测试优先级排序
    t.Run("Priority ordering", func(t *testing.T) {
        registry := NewServiceRegistry()

        registry.Register(&ServiceDef{Name: "A", Dependencies: []string{}, Priority: 10})
        registry.Register(&ServiceDef{Name: "B", Dependencies: []string{}, Priority: 0})
        registry.Register(&ServiceDef{Name: "C", Dependencies: []string{}, Priority: 5})

        order, err := registry.GetInitializationOrder()
        if err != nil {
            t.Errorf("expected no error, got %v", err)
        }

        // 验证顺序: B(0) -> C(5) -> A(10)
        if order[0] != "B" || order[1] != "C" || order[2] != "A" {
            t.Errorf("wrong order: %v", order)
        }
    })
}
```

**Step 2: 运行测试**

```bash
go test ./service/container/... -v
```

**Step 3: 启动服务验证**

```bash
# 查看依赖图
go run cmd/main.go --show-dependencies

# 正常启动
go run cmd/main.go
```

**Step 4: 提交**

```bash
git add service/container/registry_test.go
git commit -m "test(container): add ServiceRegistry tests"
```

---

### Phase 2 验收标准

- [ ] 服务依赖显式声明在注册表中
- [ ] 容器自动解析初始化顺序
- [ ] 循环依赖能被检测并报错
- [ ] 可选服务初始化失败不阻塞系统启动
- [ ] 所有测试通过

---

## Phase 3: shared模块职责过重 [P1]

### 目标
将 `service/shared/` 模块按职责拆分为独立模块。

### 当前结构
```
service/shared/
├── auth/           # 认证服务
├── cache/          # 缓存服务（只有redis_cache_service.go）
├── messaging/      # 消息服务
├── metrics/        # 服务指标
├── storage/        # 存储服务
├── stats/          # 统计服务
└── ...
```

### 目标结构
```
service/
├── auth/           # 独立认证模块
├── cache/          # 独立缓存模块
├── storage/        # 独立存储模块
├── messaging/      # 独立消息模块
├── metrics/        # 独立指标模块
└── ...
```

### Task 3.1: 创建新模块结构

**Files:**
- Create: `service/auth/auth_service.go`
- Create: `service/cache/cache_service.go`
- Create: `service/storage/storage_service.go`

**Step 1: 移动auth模块**

```bash
# 创建目标目录
mkdir -p service/auth

# 移动文件（使用git mv保留历史）
git mv service/shared/auth service/auth
```

**Step 2: 更新导入路径**

```bash
# 批量替换导入路径
find . -name "*.go" -type f -exec sed -i 's|Qingyu_backend/service/shared/auth|Qingyu_backend/service/auth|g' {} \;
```

**Step 3: 移动cache服务**

```go
// service/cache/cache_service.go
package cache

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

// CacheService 缓存服务接口
type CacheService interface {
    // Get 获取缓存
    Get(ctx context.Context, key string) (string, error)

    // Set 设置缓存
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

    // Delete 删除缓存
    Delete(ctx context.Context, keys ...string) error

    // Exists 检查缓存是否存在
    Exists(ctx context.Context, key string) (bool, error)
}

// RedisCacheService Redis缓存服务实现
type RedisCacheService struct {
    client *redis.Client
}

// NewRedisCacheService 创建Redis缓存服务
func NewRedisCacheService(client *redis.Client) CacheService {
    return &RedisCacheService{
        client: client,
    }
}

// Get 实现
func (s *RedisCacheService) Get(ctx context.Context, key string) (string, error) {
    return s.client.Get(ctx, key).Result()
}

// Set 实现
func (s *RedisCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    return s.client.Set(ctx, key, value, expiration).Err()
}

// Delete 实现
func (s *RedisCacheService) Delete(ctx context.Context, keys ...string) error {
    return s.client.Del(ctx, keys...).Err()
}

// Exists 实现
func (s *RedisCacheService) Exists(ctx context.Context, key string) (bool, error) {
    n, err := s.client.Exists(ctx, key).Result()
    return n > 0, err
}
```

**Step 4: 移动storage服务**

```bash
# storage已经是独立实现，只需要移动位置
git mv service/shared/storage service/storage
```

**Step 5: 更新导入**

```bash
find . -name "*.go" -type f -exec sed -i 's|Qingyu_backend/service/shared/storage|Qingyu_backend/service/storage|g' {} \;
find . -name "*.go" -type f -exec sed -i 's|Qingyu_backend/service/shared/messaging|Qingyu_backend/service/messaging|g' {} \;
```

**Step 6: 提交**

```bash
git add .
git commit -m "refactor(service): split shared module into independent services"
```

---

### Task 3.2: 更新ServiceContainer

**Files:**
- Modify: `service/container/service_container.go`

**Step 1: 更新导入**

```go
import (
    // ...
    // 旧: "Qingyu_backend/service/shared/auth"
    "Qingyu_backend/service/auth"
    "Qingyu_backend/service/cache"
    "Qingyu_backend/service/storage"
    // ...
)
```

**Step 2: 更新服务字段类型**

```go
type ServiceContainer struct {
    // ...
    // 旧: authService auth.AuthService
    authService service.AuthService  // 使用service包别名避免冲突

    // 新增缓存服务
    cacheService cache.CacheService

    // storage已经是独立类型
    storageService storage.StorageService
    // ...
}
```

**Step 3: 更新初始化代码**

```go
// 在service_registration.go中更新
import (
    svcAuth "Qingyu_backend/service/auth"
    svcCache "Qingyu_backend/service/cache"
    svcStorage "Qingyu_backend/service/storage"
)

c.RegisterServiceDefinition(&ServiceDef{
    Name: "AuthService",
    Factory: func(c *ServiceContainer) (interface{}, error) {
        // 使用svcAuth包
        return svcAuth.NewAuthService(...), nil
    },
    // ...
})
```

**Step 4: 提交**

```bash
git add service/container/service_container.go
git commit -m "refactor(container): update imports after module split"
```

---

### Task 3.3: 清理shared目录

**Files:**
- Delete: `service/shared/` (保留必要的)

**Step 1: 保留必要文件**

```
service/shared/
├── config_service.go    # 保留：配置服务
├── permission_service.go  # 移动到auth
└── stats/               # 移动到独立metrics
```

**Step 2: 移动剩余文件**

```bash
# 移动permission
git mv service/shared/permission_service.go service/auth/

# 移动stats
git mv service/shared/stats service/metrics
```

**Step 3: 删除空目录**

```bash
# 确认没有文件后删除
rm -rf service/shared/
```

**Step 4: 更新文档**

```bash
# 更新架构文档
# docs/architecture/system_architecture.md
```

**Step 5: 提交**

```bash
git add -A
git commit -m "refactor(service): remove shared module, split into independent services"
```

---

### Phase 3 验收标准

- [ ] shared目录已删除或只保留config
- [ ] 所有服务按职责独立到自己的目录
- [ ] 导入路径全部更新
- [ ] 所有测试通过
- [ ] 服务启动正常

---

## Phase 4: WriterService耦合度高 [P2]

### 目标
使用领域事件解耦WriterService，降低其依赖数量。

### 当前依赖
```go
// WriterService依赖:
// - BookstoreService (书籍管理)
// - AIService (AI辅助)
// - EventService (事件发布)
// - NotificationService (通知)
// - FinanceService (财务结算)
```

### 解决方案
采用领域事件模式，WriterService发布事件，其他服务订阅处理。

### Task 4.1: 定义Writer领域事件

**Files:**
- Create: `service/writer/events/events.go`

**Step 1: 创建事件定义**

```go
// service/writer/events/events.go
package events

import (
    "time"
)

// ChapterCreatedEvent 章节创建事件
type ChapterCreatedEvent struct {
    ChapterID     string
    ProjectID     string
    AuthorID      string
    ChapterTitle  string
    WordCount     int
    CreatedAt     time.Time
}

// ChapterPublishedEvent 章节发布事件
type ChapterPublishedEvent struct {
    ChapterID     string
    ProjectID     string
    BookID        string
    AuthorID      string
    ChapterTitle  string
    ChapterNumber int
    PublishedAt   time.Time
}

// BookCompletedEvent 书籍完结事件
type BookCompletedEvent struct {
    BookID      string
    ProjectID   string
    AuthorID    string
    BookTitle   string
    TotalWords  int
    CompletedAt time.Time
}

// AIContentRequestEvent AI内容请求事件
type AIContentRequestEvent struct {
    RequestID   string
    AuthorID    string
    ProjectID   string
    ContentType string // "outline", "chapter", "dialogue"
    Prompt      string
    RequestedAt time.Time
}
```

**Step 2: 提交**

```bash
git add service/writer/events/events.go
git commit -m "feat(writer): add domain events for writer service"
```

---

### Task 4.2: 创建事件处理器

**Files:**
- Create: `service/writer/handlers/notification_handler.go`
- Create: `service/writer/handlers/finance_handler.go`
- Create: `service/writer/handlers/analytics_handler.go`

**Step 1: 通知处理器**

```go
// service/writer/handlers/notification_handler.go
package handlers

import (
    "context"
    "fmt"

    "Qingyu_backend/service/writer/events"
    writerEvents "Qingyu_backend/service/writer/events"
)

// NotificationEventHandler 通知事件处理器
type NotificationEventHandler struct {
    // 不再直接依赖NotificationService
    // 通过事件总线发送通知
}

// NewNotificationEventHandler 创建处理器
func NewNotificationEventHandler() *NotificationEventHandler {
    return &NotificationEventHandler{}
}

// HandleChapterPublished 处理章节发布事件
func (h *NotificationEventHandler) HandleChapterPublished(
    ctx context.Context,
    event *writerEvents.ChapterPublishedEvent,
) error {
    // 发布通知事件到全局事件总线
    // 不再直接调用NotificationService

    // 构造通知事件
    notificationEvent := map[string]interface{}{
        "type":      "chapter_published",
        "chapter_id": event.ChapterID,
        "book_id":   event.BookID,
        "author_id": event.AuthorID,
        "title":     event.ChapterTitle,
        "number":    event.ChapterNumber,
    }

    // 通过事件总线发送
    // eventBus.Publish("notification:chapter_published", notificationEvent)

    fmt.Printf("[通知] 章节 %s 已发布\n", event.ChapterTitle)
    return nil
}

// HandleBookCompleted 处理书籍完结事件
func (h *NotificationEventHandler) HandleBookCompleted(
    ctx context.Context,
    event *writerEvents.ChapterPublishedEvent,
) error {
    fmt.Printf("[通知] 书籍 %s 已完结\n", event.BookTitle)
    return nil
}
```

**Step 2: 财务处理器**

```go
// service/writer/handlers/finance_handler.go
package handlers

import (
    "context"
    "fmt"

    writerEvents "Qingyu_backend/service/writer/events"
)

// FinanceEventHandler 财务事件处理器
type FinanceEventHandler struct{}

// NewFinanceEventHandler 创建处理器
func NewFinanceEventHandler() *FinanceEventHandler {
    return &FinanceEventHandler{}
}

// HandleChapterPublished 处理章节发布（计算收益）
func (h *FinanceEventHandler) HandleChapterPublished(
    ctx context.Context,
    event *writerEvents.ChapterPublishedEvent,
) error {
    // 发布财务计算事件
    financeEvent := map[string]interface{}{
        "type":       "chapter_revenue",
        "chapter_id": event.ChapterID,
        "author_id":  event.AuthorID,
        "word_count": 0, // 需要获取
    }

    // eventBus.Publish("finance:calculate", financeEvent)
    fmt.Printf("[财务] 计算章节 %s 收益\n", event.ChapterTitle)
    return nil
}
```

**Step 3: 提交**

```bash
git add service/writer/handlers/
git commit -m "feat(writer): add event handlers for notification and finance"
```

---

### Task 4.3: 重构WriterService

**Files:**
- Modify: `service/writer/writer_service.go`

**Step 1: 简化WriterService依赖**

```go
// service/writer/writer_service.go
package writer

import (
    "context"
    "fmt"

    writerEvents "Qingyu_backend/service/writer/events"
)

// WriterService 作者服务（简化版）
type WriterService struct {
    // 必要依赖
    projectRepo    ProjectRepository
    chapterRepo    ChapterRepository
    eventBus       EventBus

    // 移除的直接依赖:
    // - bookstoreService
    // - aiService
    // - notificationService
    // - financeService
}

// NewWriterService 创建作者服务
func NewWriterService(
    projectRepo ProjectRepository,
    chapterRepo ChapterRepository,
    eventBus EventBus,
) *WriterService {
    return &WriterService{
        projectRepo: projectRepo,
        chapterRepo: chapterRepo,
        eventBus:    eventBus,
    }
}

// PublishChapter 发布章节（使用事件解耦）
func (s *WriterService) PublishChapter(
    ctx context.Context,
    chapterID string,
) error {
    // 1. 获取章节
    chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
    if err != nil {
        return fmt.Errorf("获取章节失败: %w", err)
    }

    // 2. 更新状态为已发布
    chapter.Status = "published"
    if err := s.chapterRepo.Update(ctx, chapter); err != nil {
        return fmt.Errorf("更新章节失败: %w", err)
    }

    // 3. 发布领域事件
    event := &writerEvents.ChapterPublishedEvent{
        ChapterID:     chapter.ID,
        ProjectID:     chapter.ProjectID,
        BookID:        chapter.BookID,
        AuthorID:      chapter.AuthorID,
        ChapterTitle:  chapter.Title,
        ChapterNumber: chapter.Number,
        PublishedAt:   chapter.PublishedAt,
    }

    // 通过事件总线发布
    if err := s.eventBus.Publish("chapter:published", event); err != nil {
        return fmt.Errorf("发布事件失败: %w", err)
    }

    // 其他服务通过订阅此事件来处理:
    // - NotificationService: 发送通知
    // - FinanceService: 计算收益
    // - BookstoreService: 更新书籍

    return nil
}
```

**Step 2: 提交**

```bash
git add service/writer/writer_service.go
git commit -m "refactor(writer): simplify dependencies using domain events"
```

---

### Task 4.4: 注册事件订阅

**Files:**
- Create: `service/writer/subscriptions.go`

**Step 1: 创建订阅注册**

```go
// service/writer/subscriptions.go
package writer

import (
    "Qingyu_backend/service/writer/events"
    "Qingyu_backend/service/writer/handlers"
)

// RegisterEventSubscriptions 注册事件订阅
func RegisterEventSubscriptions(
    eventBus EventBus,
    notificationHandler *handlers.NotificationEventHandler,
    financeHandler *handlers.FinanceEventHandler,
) {
    // 订阅章节发布事件
    eventBus.Subscribe("chapter:published", func(event interface{}) error {
        if e, ok := event.(*events.ChapterPublishedEvent); ok {
            // 通知处理
            notificationHandler.HandleChapterPublished(context.Background(), e)
            // 财务处理
            financeHandler.HandleChapterPublished(context.Background(), e)
        }
        return nil
    })

    // 订阅书籍完结事件
    eventBus.Subscribe("book:completed", func(event interface{}) error {
        if e, ok := event.(*events.BookCompletedEvent); ok {
            notificationHandler.HandleBookCompleted(context.Background(), e)
        }
        return nil
    })
}
```

**Step 2: 在ServiceContainer中注册**

```go
// 在 service_container.go 的 InitializeServices 中
func (c *ServiceContainer) InitializeServices(ctx context.Context) error {
    // ... 服务初始化

    // 注册事件订阅
    notificationHandler := handlers.NewNotificationEventHandler()
    financeHandler := handlers.NewFinanceEventHandler()

    writer.RegisterEventSubscriptions(c.eventBus, notificationHandler, financeHandler)

    return nil
}
```

**Step 3: 提交**

```bash
git add service/writer/subscriptions.go
git commit -m "feat(writer): register event subscriptions for cross-service communication"
```

---

### Phase 4 验收标准

- [ ] WriterService依赖减少到核心依赖（Project、Chapter、EventBus）
- [ ] 跨模块通信通过事件完成
- [ ] 单元测试覆盖新增事件处理
- [ ] 集成测试验证事件流转

---

## Phase 5: 事件总线持久化 [P2]

### 目标
使用Redis Stream实现持久化事件总线，支持事件重放和溯源。

### Task 5.1: 设计持久化事件总线

**Files:**
- Create: `service/events/persistent_event_bus.go`

**Step 1: 定义接口**

```go
// service/events/persistent_event_bus.go
package events

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/redis/go-redis/v9"
)

// PersistentEventBus 持久化事件总线
type PersistentEventBus struct {
    client     *redis.Client
    streamName string
    consumers  map[string][]EventHandler
}

// NewPersistentEventBus 创建持久化事件总线
func NewPersistentEventBus(client *redis.Client, streamName string) *PersistentEventBus {
    return &PersistentEventBus{
        client:     client,
        streamName: streamName,
        consumers:  make(map[string][]EventHandler),
    }
}

// Publish 发布事件（持久化到Redis Stream）
func (b *PersistentEventBus) Publish(eventType string, event interface{}) error {
    ctx := context.Background()

    // 序列化事件
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("序列化事件失败: %w", err)
    }

    // 写入Redis Stream
    _, err = b.client.XAdd(ctx, &redis.XAddArgs{
        Stream: b.streamName,
        Values: map[string]interface{}{
            "type": eventType,
            "data": data,
            "timestamp": time.Now().Unix(),
        },
    }).Result()

    if err != nil {
        return fmt.Errorf("写入Stream失败: %w", err)
    }

    // 通知内存订阅者
    b.notifySubscribers(eventType, event)

    return nil
}

// Subscribe 订阅事件
func (b *PersistentEventBus) Subscribe(eventType string, handler EventHandler) func() {
    b.consumers[eventType] = append(b.consumers[eventType], handler)

    // 返回取消订阅函数
    return func() {
        b.Unsubscribe(eventType, handler)
    }
}

// Unsubscribe 取消订阅
func (b *PersistentEventBus) Unsubscribe(eventType string, handler EventHandler) {
    handlers := b.consumers[eventType]
    for i, h := range handlers {
        if h == handler {
            b.consumers[eventType] = append(handlers[:i], handlers[i+1:]...)
            break
        }
    }
}

// Replay 重放事件
func (b *PersistentEventBus) Replay(ctx context.Context, count int64) error {
    // 从Stream读取历史事件
    streams, err := b.client.XRead(ctx, &redis.XReadArgs{
        Streams: []string{b.streamName, "0"},
        Count:   count,
    }).Result()

    if err != nil {
        return fmt.Errorf("读取Stream失败: %w", err)
    }

    for _, stream := range streams {
        for _, message := range stream.Messages {
            eventType := message.Values["type"].(string)
            data := message.Values["data"].(string)

            // 反序列化并处理
            var event interface{}
            if err := json.Unmarshal([]byte(data), &event); err != nil {
                fmt.Printf("反序列化事件失败: %v\n", err)
                continue
            }

            b.notifySubscribers(eventType, event)
        }
    }

    return nil
}

// notifySubscribers 通知订阅者
func (b *PersistentEventBus) notifySubscribers(eventType string, event interface{}) {
    handlers, exists := b.consumers[eventType]
    if !exists {
        return
    }

    for _, handler := range handlers {
        go func(h EventHandler) {
            if err := h(event); err != nil {
                fmt.Printf("事件处理失败: %v\n", err)
            }
        }(handler)
    }
}
```

**Step 2: 提交**

```bash
git add service/events/persistent_event_bus.go
git commit -m "feat(events): add persistent event bus with Redis Stream"
```

---

### Task 5.2: 集成到ServiceContainer

**Files:**
- Modify: `service/container/service_container.go`

**Step 1: 修改事件总线初始化**

```go
// 在 Initialize 方法中
func (c *ServiceContainer) Initialize(ctx context.Context) error {
    // ...

    // 创建持久化事件总线
    if c.redisClient != nil {
        rawClient := c.redisClient.GetClient()
        if redisClient, ok := rawClient.(*redis.Client); ok {
            c.eventBus = events.NewPersistentEventBus(redisClient, "qingyu:events")
            fmt.Println("✓ 使用持久化事件总线 (Redis Stream)")
        }
    } else {
        c.eventBus = base.NewSimpleEventBus()
        fmt.Println("⚠ 使用内存事件总线（服务重启后事件丢失）")
    }

    // ...
}
```

**Step 2: 提交**

```bash
git add service/container/service_container.go
git commit -m "refactor(container): integrate persistent event bus"
```

---

### Task 5.3: 添加管理接口

**Files:**
- Create: `api/v1/admin/events_api.go`

**Step 1: 创建管理API**

```go
// api/v1/admin/events_api.go
package admin

import (
    "github.com/gin-gonic/gin"
)

// EventsAPI 事件管理API
type EventsAPI struct {
    container *container.ServiceContainer
}

// NewEventsAPI 创建API
func NewEventsAPI(container *container.ServiceContainer) *EventsAPI {
    return &EventsAPI{container: container}
}

// ReplayEvents 重放事件
// POST /admin/events/replay
func (api *EventsAPI) ReplayEvents(c *gin.Context) {
    var req struct {
        Count int64 `form:"count" binding:"required,min=1,max=10000"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    eventBus := api.container.GetEventBus()

    // 检查是否为持久化事件总线
    persistentBus, ok := eventBus.(*events.PersistentEventBus)
    if !ok {
        c.JSON(400, gin.H{"error": "当前使用内存事件总线，不支持重放"})
        return
    }

    if err := persistentBus.Replay(c.Request.Context(), req.Count); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "事件重放完成"})
}
```

**Step 2: 提交**

```bash
git add api/v1/admin/events_api.go
git commit -m "feat(admin): add event replay API"
```

---

### Phase 5 验收标准

- [ ] 事件持久化到Redis Stream
- [ ] 支持事件重放
- [ ] 管理API可用
- [ ] 降级到内存模式（Redis不可用时）

---

## 总体验收标准

### 功能验收
- [ ] Phase 1: 中间件解耦完成，所有API正常
- [ ] Phase 2: 服务启动顺序自动解析
- [ ] Phase 3: shared模块拆分完成
- [ ] Phase 4: WriterService依赖减少
- [ ] Phase 5: 事件总线持久化

### 质量验收
- [ ] 所有单元测试通过
- [ ] E2E测试通过
- [ ] 代码覆盖率 > 80%
- [ ] 性能测试无退化

### 文档验收
- [ ] 架构文档更新
- [ ] API文档更新
- [ ] 迁移指南编写

---

## 附录: 相关技能参考

在执行此计划时，以下技能可能有用：

- @superpowers:brainstorming - 创建新功能前进行头脑风暴
- @superpowers:test-driven-development - 使用TDD开发新功能
- @superpowers:systematic-debugging - 遇到问题时系统化调试
- @superpowers:verification-before-completion - 完成前验证
- @codex-review - 代码审查
- @codex-test - 运行测试

---

**计划版本**: v1.0
**创建日期**: 2026-02-07
**维护者**: yukin371
