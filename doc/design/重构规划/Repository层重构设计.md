# Repository层重构设计文档

## 1. 概述

### 1.1 重构目标
基于新的分层架构设计，对现有Repository层进行重构，实现更清晰的数据访问抽象、更好的可测试性和更强的扩展性。

### 1.2 现状分析
当前Repository层存在以下问题：
- 接口定义分散，缺乏统一的基础抽象
- 事务管理复杂，与业务逻辑耦合
- 缺乏统一的错误处理机制
- 查询优化和缓存策略不统一
- 测试覆盖率不足

### 1.3 重构原则
- **接口抽象**：定义清晰的Repository接口，实现与具体数据库技术解耦
- **实现分离**：支持多种数据库实现（MongoDB、PostgreSQL等）
- **事务管理**：统一的事务处理机制
- **错误处理**：标准化的错误类型和处理流程
- **API兼容**：保持现有API的向后兼容性

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   AI Service    │  │ Document Service│  │ User Service │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Repository Layer (New)                       │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  Base Repository│  │Transaction Mgr  │  │ Query Builder│ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  User Repository│  │  AI Repository  │  │Doc Repository│ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                Database Adapters                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │MongoDB Adapter  │  │PostgreSQL Adapter│ │ Cache Adapter│ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 核心组件

#### 2.2.1 BaseRepository接口
```go
type BaseRepository[T any, ID comparable] interface {
    // 基础CRUD操作
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id ID) (*T, error)
    Update(ctx context.Context, id ID, updates map[string]interface{}) error
    Delete(ctx context.Context, id ID) error
    
    // 查询操作
    List(ctx context.Context, filter Filter) ([]*T, error)
    Count(ctx context.Context, filter Filter) (int64, error)
    Exists(ctx context.Context, id ID) (bool, error)
    
    // 批量操作
    BatchCreate(ctx context.Context, entities []*T) error
    BatchUpdate(ctx context.Context, ids []ID, updates map[string]interface{}) error
    BatchDelete(ctx context.Context, ids []ID) error
    
    // 事务支持
    WithTransaction(ctx context.Context, fn func(ctx context.Context, repo BaseRepository[T, ID]) error) error
}
```

#### 2.2.2 专用Repository接口
```go
// UserRepository 用户仓储接口
type UserRepository interface {
    BaseRepository[system.User, string]
    
    // 用户特有方法
    GetByUsername(ctx context.Context, username string) (*system.User, error)
    GetByEmail(ctx context.Context, email string) (*system.User, error)
    UpdateLastLogin(ctx context.Context, id string) error
    UpdatePassword(ctx context.Context, id string, hashedPassword string) error
    GetActiveUsers(ctx context.Context, limit int64) ([]*system.User, error)
}

// AIRepository AI相关仓储接口
type AIRepository interface {
    // 聊天会话管理
    CreateChatSession(ctx context.Context, session *ai.ChatSession) error
    GetChatSession(ctx context.Context, sessionID string) (*ai.ChatSession, error)
    UpdateChatSession(ctx context.Context, session *ai.ChatSession) error
    DeleteChatSession(ctx context.Context, sessionID string) error
    ListChatSessions(ctx context.Context, projectID string, filter Filter) ([]*ai.ChatSession, error)
    
    // 聊天消息管理
    CreateChatMessage(ctx context.Context, message *ai.ChatMessage) error
    GetChatMessages(ctx context.Context, sessionID string, filter Filter) ([]*ai.ChatMessage, error)
    
    // AI配置管理
    GetAIConfig(ctx context.Context, configID string) (*ai.AIConfig, error)
    UpdateAIConfig(ctx context.Context, config *ai.AIConfig) error
}

// DocumentRepository 文档仓储接口
type DocumentRepository interface {
    BaseRepository[document.Document, string]
    
    // 文档特有方法
    GetByProjectID(ctx context.Context, projectID string, filter Filter) ([]*document.Document, error)
    GetByUserID(ctx context.Context, userID string, filter Filter) ([]*document.Document, error)
    UpdateContent(ctx context.Context, id string, content string) error
    GetVersions(ctx context.Context, documentID string) ([]*document.Version, error)
}
```

## 3. 事务管理设计

### 3.1 TransactionManager重构
```go
type TransactionManager interface {
    // 执行单个事务
    ExecuteTransaction(ctx context.Context, fn func(ctx TransactionContext) error) error
    
    // 执行分布式事务（Saga模式）
    ExecuteSaga(ctx context.Context, saga *Saga) error
    
    // 获取事务上下文
    GetTransactionContext(ctx context.Context) (TransactionContext, error)
}

type TransactionContext interface {
    context.Context
    
    // 获取事务化的Repository
    UserRepository() UserRepository
    AIRepository() AIRepository
    DocumentRepository() DocumentRepository
    
    // 事务控制
    Commit() error
    Rollback() error
    IsInTransaction() bool
}
```

### 3.2 Saga模式支持
```go
type Saga struct {
    ID    string
    Steps []SagaStep
}

type SagaStep struct {
    Name        string
    Execute     func(ctx TransactionContext) error
    Compensate  func(ctx TransactionContext) error
    Timeout     time.Duration
}
```

## 4. 查询构建器设计

### 4.1 QueryBuilder接口
```go
type QueryBuilder interface {
    // 条件构建
    Where(field string, operator string, value interface{}) QueryBuilder
    WhereIn(field string, values []interface{}) QueryBuilder
    WhereNotIn(field string, values []interface{}) QueryBuilder
    WhereBetween(field string, start, end interface{}) QueryBuilder
    WhereNull(field string) QueryBuilder
    WhereNotNull(field string) QueryBuilder
    
    // 逻辑操作
    And() QueryBuilder
    Or() QueryBuilder
    Not() QueryBuilder
    
    // 排序和分页
    OrderBy(field string, direction string) QueryBuilder
    Limit(limit int64) QueryBuilder
    Offset(offset int64) QueryBuilder
    
    // 聚合操作
    GroupBy(fields ...string) QueryBuilder
    Having(field string, operator string, value interface{}) QueryBuilder
    
    // 构建查询
    Build() (interface{}, error)
    BuildCount() (interface{}, error)
}
```

### 4.2 Filter统一接口
```go
type Filter interface {
    ToQuery() (interface{}, error)
    GetLimit() int64
    GetOffset() int64
    GetSort() map[string]int
}

type BaseFilter struct {
    Limit  int64             `json:"limit,omitempty"`
    Offset int64             `json:"offset,omitempty"`
    Sort   map[string]int    `json:"sort,omitempty"`
    Fields map[string]interface{} `json:"fields,omitempty"`
}
```

## 5. 错误处理设计

### 5.1 错误类型定义
```go
type RepositoryError struct {
    Type      ErrorType `json:"type"`
    Message   string    `json:"message"`
    Code      string    `json:"code"`
    Details   map[string]interface{} `json:"details,omitempty"`
    Cause     error     `json:"-"`
    Timestamp time.Time `json:"timestamp"`
}

type ErrorType string

const (
    ErrorTypeNotFound     ErrorType = "NOT_FOUND"
    ErrorTypeDuplicate    ErrorType = "DUPLICATE"
    ErrorTypeValidation   ErrorType = "VALIDATION"
    ErrorTypeConnection   ErrorType = "CONNECTION"
    ErrorTypeTransaction  ErrorType = "TRANSACTION"
    ErrorTypeTimeout      ErrorType = "TIMEOUT"
    ErrorTypePermission   ErrorType = "PERMISSION"
    ErrorTypeInternal     ErrorType = "INTERNAL"
)
```

### 5.2 错误处理工具
```go
type ErrorHandler interface {
    Handle(err error) *RepositoryError
    IsRetryable(err error) bool
    ShouldLog(err error) bool
}

func NewRepositoryError(errorType ErrorType, message string, cause error) *RepositoryError
func IsNotFoundError(err error) bool
func IsDuplicateError(err error) bool
func IsValidationError(err error) bool
```

## 6. 缓存策略设计

### 6.1 缓存接口
```go
type CacheRepository interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    DeletePattern(ctx context.Context, pattern string) error
    Exists(ctx context.Context, key string) (bool, error)
    
    // 批量操作
    MGet(ctx context.Context, keys []string) (map[string]interface{}, error)
    MSet(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
    MDelete(ctx context.Context, keys []string) error
}
```

### 6.2 缓存策略
```go
type CacheStrategy interface {
    GenerateKey(entity interface{}, operation string) string
    GetTTL(entity interface{}) time.Duration
    ShouldCache(entity interface{}, operation string) bool
    ShouldInvalidate(entity interface{}, operation string) bool
}
```

## 7. 实现计划

### 7.1 第一阶段：基础架构
1. 实现BaseRepository接口和基础实现
2. 重构TransactionManager
3. 实现QueryBuilder和Filter系统
4. 统一错误处理机制

### 7.2 第二阶段：专用Repository
1. 重构UserRepository实现
2. 实现AIRepository
3. 重构DocumentRepository
4. 添加缓存支持

### 7.3 第三阶段：高级特性
1. 实现Saga事务模式
2. 添加查询优化
3. 实现读写分离
4. 添加监控和指标

## 8. API兼容性保证

### 8.1 兼容性策略
- 保持现有公共API接口不变
- 通过适配器模式包装新实现
- 渐进式迁移，支持新旧实现并存
- 提供迁移工具和文档

### 8.2 迁移路径
```go
// 旧接口保持兼容
type LegacyUserRepository interface {
    Create(ctx context.Context, user *system.User) error
    GetByID(ctx context.Context, id string) (*system.User, error)
    // ... 其他现有方法
}

// 适配器实现
type UserRepositoryAdapter struct {
    newRepo UserRepository
}

func (a *UserRepositoryAdapter) Create(ctx context.Context, user *system.User) error {
    return a.newRepo.Create(ctx, user)
}
```

## 9. 测试策略

### 9.1 单元测试
- Repository接口的Mock实现
- 各种数据库适配器的单元测试
- 事务管理器的测试
- 查询构建器的测试

### 9.2 集成测试
- 数据库连接和操作测试
- 事务完整性测试
- 缓存一致性测试
- 性能基准测试

### 9.3 兼容性测试
- 现有API的回归测试
- 数据迁移测试
- 并发访问测试

## 10. 监控和指标

### 10.1 性能指标
- 查询响应时间
- 事务成功率
- 缓存命中率
- 连接池使用率

### 10.2 业务指标
- 操作成功率
- 错误类型分布
- 数据一致性检查
- 资源使用情况

这个重构设计确保了系统的可扩展性、可维护性和性能，同时保持了API的向后兼容性。