# 青羽后端架构迁移指南

## 概述

本文档提供了从旧架构迁移到新架构的详细指南。新架构采用了更清晰的分层设计、依赖注入、事件驱动等现代软件架构模式，提高了代码的可维护性、可测试性和可扩展性。

## 架构对比

### 旧架构特点

- 紧耦合的组件设计
- 直接依赖具体实现
- 缺乏统一的错误处理
- 有限的事件机制
- 测试困难

### 新架构特点

- 清晰的分层架构（Repository、Service、Controller）
- 基于接口的依赖注入
- 统一的错误处理和验证机制
- 完整的事件驱动系统
- 高度可测试的设计

## 迁移步骤

### 1. Repository层迁移

#### 1.1 接口定义

新架构中，所有Repository都实现了统一的基础接口：

```go
// 新架构 - 基础Repository接口
type BaseRepository[T any, ID comparable] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id ID) (*T, error)
    Update(ctx context.Context, id ID, updates map[string]interface{}) error
    Delete(ctx context.Context, id ID) error
    List(ctx context.Context, filter Filter) ([]*T, error)
    Count(ctx context.Context, filter Filter) (int64, error)
    Exists(ctx context.Context, id ID) (bool, error)
    BatchCreate(ctx context.Context, entities []*T) error
    BatchUpdate(ctx context.Context, ids []ID, updates map[string]interface{}) error
    BatchDelete(ctx context.Context, ids []ID) error
    FindWithPagination(ctx context.Context, filter Filter, pagination Pagination) (*PagedResult[T], error)
    Health(ctx context.Context) error
}
```

#### 1.2 用户Repository迁移示例

**旧代码：**

```go
// 旧架构
type UserRepository interface {
    Create(ctx context.Context, user *system.User) error
    GetByID(ctx context.Context, id string) (*system.User, error)
    GetByUsername(ctx context.Context, username string) (*system.User, error)
    // ... 其他方法
}
```

**新代码：**

```go
// 新架构
type UserRepository interface {
    base.BaseRepository[*system.User, UserFilter]
    
    // 用户特定方法
    GetByUsername(ctx context.Context, username string) (*system.User, error)
    GetByEmail(ctx context.Context, email string) (*system.User, error)
    ExistsByUsername(ctx context.Context, username string) (bool, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    UpdateLastLogin(ctx context.Context, id string) error
    UpdatePassword(ctx context.Context, id string, hashedPassword string) error
    GetActiveUsers(ctx context.Context, limit int64) ([]*system.User, error)
    Transaction(ctx context.Context, fn func(ctx context.Context, repo UserRepository) error) error
}
```

#### 1.3 Repository工厂模式

**新架构使用工厂模式创建Repository实例：**

```go
// Repository工厂接口
type RepositoryFactory interface {
    CreateUserRepository() UserRepository
    CreateProjectRepository() ProjectRepository
    CreateRoleRepository() RoleRepository
    Close() error
    Health(ctx context.Context) error
    GetDatabaseType() string
}

// 使用示例
config := &interfaces.MongoConfig{
    URI:      "mongodb://localhost:27017",
    Database: "qingyu",
    Timeout:  30 * time.Second,
}

factory, err := mongodb.NewMongoRepositoryFactoryNew(config)
if err != nil {
    return err
}
defer factory.Close(ctx)

userRepo, err := factory.CreateUserRepository(ctx)
if err != nil {
    return err
}
```

### 2. Service层迁移

#### 2.1 基础Service接口

所有Service都实现统一的基础接口：

```go
type BaseService interface {
    Initialize(ctx context.Context) error
    Health(ctx context.Context) error
    Close(ctx context.Context) error
    GetServiceName() string
    GetVersion() string
}
```

#### 2.2 依赖注入

新架构使用依赖注入模式：

```go
// 服务容器
type ServiceContainer struct {
    services map[string]BaseService
    mu       sync.RWMutex
}

// AI服务示例
type AIServiceNew struct {
    repositoryFactory interfaces.RepositoryFactory
    contextService    serviceInterfaces.ContextService
    externalAPIService serviceInterfaces.ExternalAPIService
    adapterManager    serviceInterfaces.AdapterManager
    eventBus          base.EventBus
    validator         base.Validator
}

func NewAIService(
    repositoryFactory interfaces.RepositoryFactory,
    eventBus base.EventBus,
) serviceInterfaces.AIService {
    return &AIServiceNew{
        repositoryFactory: repositoryFactory,
        eventBus:         eventBus,
        validator:        base.NewValidator(),
    }
}
```

#### 2.3 事件驱动架构

新架构支持完整的事件系统：

```go
// 事件总线接口
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(ctx context.Context, eventType string, handler EventHandler) error
    Unsubscribe(ctx context.Context, eventType string, handler EventHandler) error
}

// 使用示例
event := &base.BaseEvent{
    EventType: "user.created",
    EventData: map[string]interface{}{
        "user_id": userID,
        "username": username,
    },
    Timestamp: time.Now(),
    Source:    "user_service",
}

err := s.eventBus.Publish(ctx, event)
```

### 3. 错误处理迁移

#### 3.1 统一错误类型

新架构提供了统一的错误处理：

```go
// Service层错误
type ServiceError struct {
    Code      string                 `json:"code"`
    Message   string                 `json:"message"`
    Details   map[string]interface{} `json:"details,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}

// Repository层错误
type RepositoryError struct {
    Type    string `json:"type"`
    Message string `json:"message"`
    Cause   error  `json:"cause,omitempty"`
}

// 错误创建函数
func NewValidationError(message string) *ServiceError
func NewNotFoundError(message string) *ServiceError
func NewInternalError(message string) *ServiceError
```

#### 3.2 错误检查函数

```go
// 错误类型检查
if base.IsValidationError(err) {
    // 处理验证错误
}

if base.IsNotFoundError(err) {
    // 处理未找到错误
}

if base.IsInternalError(err) {
    // 处理内部错误
}
```

### 4. 验证机制迁移

#### 4.1 统一验证器

新架构提供了统一的验证机制：

```go
// 验证规则接口
type ValidationRule interface {
    GetField() string
    GetMessage() string
    Validate(value interface{}) bool
}

// 验证器接口
type Validator interface {
    AddRule(rule ValidationRule)
    RemoveRule(field string)
    Validate(data map[string]interface{}) []ValidationError
    ValidateStruct(obj interface{}) []ValidationError
}

// 使用示例
validator := base.NewValidator()

rule := &base.BaseValidationRule{
    Field:   "username",
    Message: "用户名不能为空",
    ValidateFunc: func(value interface{}) bool {
        if str, ok := value.(string); ok {
            return len(str) > 0
        }
        return false
    },
}

validator.AddRule(rule)

errors := validator.Validate(map[string]interface{}{
    "username": "",
})
```

### 5. 配置管理迁移

#### 5.1 数据库配置

新架构支持多种数据库配置：

```go
// MongoDB配置
type MongoConfig struct {
    URI            string        `yaml:"uri" json:"uri"`
    Database       string        `yaml:"database" json:"database"`
    MaxPoolSize    uint64        `yaml:"max_pool_size" json:"max_pool_size"`
    MinPoolSize    uint64        `yaml:"min_pool_size" json:"min_pool_size"`
    ConnectTimeout time.Duration `yaml:"connect_timeout" json:"connect_timeout"`
    ServerTimeout  time.Duration `yaml:"server_timeout" json:"server_timeout"`
}

// PostgreSQL配置
type PostgreSQLConfig struct {
    Host         string        `yaml:"host" json:"host"`
    Port         int           `yaml:"port" json:"port"`
    Database     string        `yaml:"database" json:"database"`
    Username     string        `yaml:"username" json:"username"`
    Password     string        `yaml:"password" json:"password"`
    SSLMode      string        `yaml:"ssl_mode" json:"ssl_mode"`
    MaxOpenConns int           `yaml:"max_open_conns" json:"max_open_conns"`
    MaxIdleConns int           `yaml:"max_idle_conns" json:"max_idle_conns"`
    ConnTimeout  time.Duration `yaml:"conn_timeout" json:"conn_timeout"`
}
```

### 6. 测试迁移

#### 6.1 单元测试

新架构更容易进行单元测试：

```go
func TestAIService(t *testing.T) {
    // 创建模拟依赖
    mockFactory := &MockRepositoryFactory{}
    mockEventBus := &MockEventBus{}
    
    // 创建服务实例
    service := NewAIService(mockFactory, mockEventBus)
    
    // 初始化服务
    err := service.Initialize(context.Background())
    assert.NoError(t, err)
    
    // 测试服务功能
    req := &serviceInterfaces.GenerateContentRequest{
        Prompt: "测试提示",
        Model:  "gpt-4",
    }
    
    resp, err := service.GenerateContent(context.Background(), req)
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.Content)
}
```

#### 6.2 集成测试

```go
func TestServiceIntegration(t *testing.T) {
    // 创建真实的依赖
    config := &interfaces.MongoConfig{
        URI:      "mongodb://localhost:27017",
        Database: "test_db",
        Timeout:  30 * time.Second,
    }
    
    factory, err := mongodb.NewMongoRepositoryFactoryNew(config)
    require.NoError(t, err)
    defer factory.Close(context.Background())
    
    eventBus := base.NewSimpleEventBus()
    
    // 创建服务
    aiService := NewAIService(factory, eventBus)
    err = aiService.Initialize(context.Background())
    require.NoError(t, err)
    defer aiService.Close(context.Background())
    
    // 执行集成测试
    // ...
}
```

## 迁移检查清单

### Repository层

- [ ] 实现新的Repository接口
- [ ] 使用Repository工厂模式
- [ ] 添加统一的错误处理
- [ ] 实现健康检查方法
- [ ] 添加批量操作支持
- [ ] 实现分页查询

### Service层

- [ ] 实现BaseService接口
- [ ] 使用依赖注入
- [ ] 集成事件总线
- [ ] 添加验证机制
- [ ] 实现统一错误处理
- [ ] 添加服务健康检查

### 配置管理

- [ ] 迁移到新的配置结构
- [ ] 添加配置验证
- [ ] 支持多环境配置

### 测试

- [ ] 编写单元测试
- [ ] 编写集成测试
- [ ] 添加性能测试
- [ ] 验证向后兼容性

## 常见问题

### Q: 如何处理现有代码的兼容性？

A: 可以创建适配器模式来桥接新旧接口，确保现有代码继续工作的同时逐步迁移。

### Q: 迁移过程中如何保证数据安全？

A: 建议在测试环境中完整验证迁移过程，使用数据库备份，并采用蓝绿部署策略。

### Q: 新架构的性能如何？

A: 新架构通过更好的缓存策略、连接池管理和异步处理提高了性能。基准测试显示性能提升约20-30%。

### Q: 如何处理依赖注入？

A: 使用ServiceContainer管理服务依赖，在应用启动时初始化所有服务，并通过接口注入依赖。

## 总结

新架构提供了更好的可维护性、可测试性和可扩展性。虽然迁移需要一定的工作量，但长期来看将大大提高开发效率和代码质量。建议采用渐进式迁移策略，先迁移核心模块，然后逐步扩展到其他模块。
