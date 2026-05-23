# Repository层设计说明书

> 最后整理: 2026-05-22  
> 当前状态: `legacy-live`

本文档是 Repository 层的历史详细设计稿，适合回看当时如何定义统一数据访问接口、MongoDB 实现方式和事务能力；它不等同于当前代码中 Repository 层实现或规范的唯一 owner。

## Page Role

- 这里负责：Repository 层架构、接口设计、MongoDB 实现、事务、性能与测试的历史方案说明。
- 不负责：当前代码目录结构事实、现行 Repository 分层规则、当前实施完成度结论。

## Recommended Read Path

1. [README.md](./README.md)
2. [../../database/README.md](../../database/README.md)
3. [../../standards/layer-repository.md](../../standards/layer-repository.md)
4. [../../standards/layer-models.md](../../standards/layer-models.md)
5. [../../implementation/infrastructure/README.md](../../implementation/infrastructure/README.md)

## Boundary

- 如果你要找“当前数据库专题 owner”，优先看 [../../database/README.md](../../database/README.md)。
- 如果你要找“现行 Repository 分层规则”，优先看 [../../standards/layer-repository.md](../../standards/layer-repository.md)。
- 如果你要找“已落地的基础设施实施记录”，优先看 [../../implementation/infrastructure/README.md](../../implementation/infrastructure/README.md)。
- 本页更适合解释设计意图和分层思路，不适合作为当前代码事实的唯一依据。

## Quick Section Map

| 如果你想看 | 直接跳到 |
|------|------|
| Repository 层目标、技术栈和整体定位 | [1. 概述](#1-概述) |
| 架构与设计模式 | [2. 架构设计](#2-架构设计) |
| 各类接口设计 | [3. 接口设计](#3-接口设计) |
| MongoDB 实现细节 | [4. MongoDB实现](#4-mongodb实现) |
| 数据模型与事务 | [5. 数据模型设计](#5-数据模型设计) / [6. 事务管理](#6-事务管理) |
| 性能、错误、测试、部署、扩展 | [7. 性能优化](#7-性能优化) / [8. 错误处理](#8-错误处理) / [9. 测试策略](#9-测试策略) / [10. 部署和维护](#10-部署和维护) / [11. 未来扩展](#11-未来扩展) |

## Quick Takeaways

- 这篇最值得先看的，是它把“接口抽象、MongoDB 实现、事务管理、测试部署”放进同一个 Repository 视角里。
- 如果你只想理解 Repository 层如何分层，优先看“2. 架构设计”“3. 接口设计”“4. MongoDB实现”。
- 后半段的测试、部署和未来扩展更偏完整方案陈述，不一定对应当前实际落地状态。

## Skip Guide

- 只想知道“Repository 层为什么这样分”：看 [2. 架构设计](#2-架构设计)。
- 只想知道“接口长什么样”：看 [3. 接口设计](#3-接口设计)。
- 只想知道“MongoDB 这一层怎么实现”：看 [4. MongoDB实现](#4-mongodb实现)。
- 如果你当前只关心现行规范或代码事实，请优先回到 `docs/standards/`、`docs/database/` 和 `docs/implementation/`。

## 1. 概述

### 1.1 文档目的
本文档详细描述青羽后端系统Repository层的设计架构、接口定义和实现方案，为开发团队提供统一的数据访问层设计规范。

### 1.2 设计目标
- **统一数据访问接口**：提供标准化的数据访问模式
- **数据库无关性**：支持多种数据库实现的切换
- **事务管理**：提供完整的事务支持
- **性能优化**：支持分页、索引管理和查询优化
- **可测试性**：便于单元测试和集成测试

### 1.3 技术栈
- **数据库**：MongoDB
- **ORM框架**：go.mongodb.org/mongo-driver
- **设计模式**：Repository Pattern + Factory Pattern
- **编程语言**：Go 1.19+

## 2. 架构设计

### 2.1 整体架构
Repository层采用分层架构设计，包含以下组件：

```
Repository层架构
├── interfaces/           # 接口定义层
│   ├── base_interface.go      # 基础接口定义
│   ├── user_repository_interface.go    # 用户仓储接口
│   ├── project_document_repository.go  # 项目文档仓储接口
│   ├── bookstore_repository.go         # 书店仓储接口
│   ├── reading_settings_repository.go  # 阅读设置仓储接口
│   ├── Error_interface.go              # 错误处理接口
│   ├── query_builder.go                # 查询构建器接口
│   └── transaction_manager_interface.go # 事务管理接口
└── mongodb/             # MongoDB实现层
    ├── factory.go            # 工厂模式实现
    ├── user_repository_mongo.go   # 用户仓储实现
    ├── project_repository_mongo.go    # 项目仓储实现
    ├── reading_settings_repository_mongo.go  # 阅读设置仓储实现
    ├── MongoQueryBuilder.go        # MongoDB查询构建器实现
    └── data_migration.go           # 数据迁移实现
```

### 2.2 设计模式

#### 2.2.1 Repository Pattern
- **目的**：封装数据访问逻辑，提供统一的数据操作接口
- **优势**：
  - 分离业务逻辑和数据访问逻辑
  - 提高代码可测试性
  - 支持多种数据源实现

#### 2.2.2 Factory Pattern
- **实现**：`MongoRepositoryFactory`
- **职责**：
  - 统一创建和管理Repository实例
  - 管理数据库连接
  - 提供健康检查功能

## 3. 接口设计

### 3.1 基础接口 (base_interface.go)

#### 3.1.1 CRUDRepository接口
```go
type CRUDRepository[T any, ID comparable] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id ID) (*T, error)
    Update(ctx context.Context, id ID, updates map[string]interface{}) error
    Delete(ctx context.Context, id ID) error
    Health(ctx context.Context) error
}
```

**设计说明**：
- 使用泛型支持不同实体类型
- 提供标准CRUD操作
- 包含健康检查接口

#### 3.1.2 查询过滤器
```go
type Filter interface {
    GetConditions() map[string]interface{}
    GetSort() map[string]int
    GetFields() []string
}

type BaseFilter struct {
    Conditions map[string]interface{}
    Sort       map[string]int
    Fields     []string
}
```

**设计说明**：
- 支持灵活的查询条件构建
- 支持排序和字段选择
- 可扩展的过滤器设计

#### 3.1.3 分页支持
```go
type Pagination struct {
    Limit  int
    Offset int
}

type PagedResult[T any] struct {
    Data       []T
    Total      int64
    Limit      int
    Offset     int
    HasMore    bool
}
```

### 3.2 用户仓储接口 (user_repository_interface.go)

#### 3.2.1 UserRepository接口
```go
type UserRepository interface {
    CRUDRepository[models.User, primitive.ObjectID]
    
    // 用户特定方法
    GetByUsername(ctx context.Context, username string) (*models.User, error)
    UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error
    Transaction(ctx context.Context, fn func(context.Context) error) error
}
```

**设计特点**：
- 继承基础CRUD接口
- 提供用户特定的查询方法
- 支持事务操作

#### 3.2.2 UserFilter
```go
type UserFilter struct {
    BaseFilter
    Username    *string
    Email       *string
    Status      *string
    CreatedTime *TimeRange
    UpdatedTime *TimeRange
}
```

### 3.3 项目仓储接口 (project_repository.go)

#### 3.3.1 ProjectRepository接口
```go
type ProjectRepository interface {
    CRUDRepository[models.Project, primitive.ObjectID]
    
    // 项目查询方法
    GetListByOwnerID(ctx context.Context, ownerID primitive.ObjectID, filter ProjectFilter, pagination Pagination) (*PagedResult[models.Project], error)
    GetByOwnerAndStatus(ctx context.Context, ownerID primitive.ObjectID, status string, limit, offset int) ([]models.Project, error)
    
    // 项目管理方法
    UpdateByOwner(ctx context.Context, ownerID, projectID primitive.ObjectID, updates map[string]interface{}) error
    IsOwner(ctx context.Context, projectID, userID primitive.ObjectID) (bool, error)
    
    // 软删除和恢复
    SoftDelete(ctx context.Context, projectID primitive.ObjectID) error
    HardDelete(ctx context.Context, projectID primitive.ObjectID) error
    RestoreByOwner(ctx context.Context, ownerID, projectID primitive.ObjectID) error
    
    // 统计方法
    CountByOwner(ctx context.Context, ownerID primitive.ObjectID) (int64, error)
    
    // 事务支持
    Transaction(ctx context.Context, fn func(context.Context) error) error
}
```

**设计特点**：
- 支持按所有者查询和管理
- 提供软删除和硬删除功能
- 支持项目恢复操作
- 包含统计功能

#### 3.3.2 相关接口
- **ProjectIndexManager**：项目索引管理
- **NodeRepository**：节点仓储接口
- **DocumentRepository**：文档仓储接口

## 4. MongoDB实现

### 4.1 工厂实现 (factory.go)

#### 4.1.1 MongoRepositoryFactory
```go
type MongoRepositoryFactory struct {
    client   *mongo.Client
    database *mongo.Database
    
    userRepo    UserRepository
    projectRepo ProjectRepository
    roleRepo    RoleRepository
}
```

**核心功能**：
- **连接管理**：统一管理MongoDB连接
- **实例创建**：创建各种Repository实例
- **健康检查**：提供数据库健康检查
- **依赖注入**：支持依赖注入模式

#### 4.1.2 主要方法
- `NewMongoRepositoryFactory()`: 创建工厂实例
- `GetUserRepository()`: 获取用户仓储
- `GetProjectRepository()`: 获取项目仓储
- `GetRoleRepository()`: 获取角色仓储
- `Health()`: 健康检查

### 4.2 用户仓储实现 (user_repository_mongo.go)

#### 4.2.1 MongoUserRepository
```go
type MongoUserRepositoryNew struct {
    collection *mongo.Collection
    client     *mongo.Client
}
```

**实现的方法**：
- **基础CRUD**：Create, GetByID, Update, Delete, Health
- **用户特定**：GetByUsername, UpdateLastLogin
- **批量操作**：List, FindWithPagination
- **事务支持**：Transaction
- **扩展功能**：Exists, HardDelete

### 4.3 项目仓储实现 (project_repository_mongo.go)

#### 4.3.1 MongoProjectRepository
```go
type MongoProjectRepository struct {
    collection *mongo.Collection
    client     *mongo.Client
}
```

**实现的方法**：
- **项目查询**：GetByOwnerID, GetByOwnerAndStatus
- **项目管理**：UpdateByOwner, DeleteByOwner
- **软删除**：SoftDelete, RestoreByOwner
- **事务支持**：CreateWithTransaction

### 4.4 新版项目仓储 (factory.go中的MongoProjectRepositoryNew)

**增强功能**：
- **批量操作**：BatchCreate, BatchUpdate, BatchDelete
- **高级查询**：List with complex filters
- **性能优化**：支持索引管理

## 5. 数据模型设计

### 5.1 用户模型 (User)
```go
type User struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Username    string            `bson:"username"`
    Email       string            `bson:"email"`
    Status      string            `bson:"status"`
    CreatedAt   time.Time         `bson:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at"`
    LastLoginAt *time.Time        `bson:"last_login_at,omitempty"`
}
```

### 5.2 项目模型 (Project)
```go
type Project struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Name        string            `bson:"name"`
    OwnerID     primitive.ObjectID `bson:"owner_id"`
    Status      string            `bson:"status"`
    IsDeleted   bool              `bson:"is_deleted"`
    CreatedAt   time.Time         `bson:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at"`
    DeletedAt   *time.Time        `bson:"deleted_at,omitempty"`
}
```

## 6. 事务管理

### 6.1 事务支持
所有Repository实现都支持MongoDB事务：

```go
func (r *MongoUserRepositoryNew) Transaction(ctx context.Context, fn func(context.Context) error) error {
    session, err := r.client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)
    
    return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
        return session.WithTransaction(sc, func(sc mongo.SessionContext) (interface{}, error) {
            return nil, fn(sc)
        })
    })
}
```

### 6.2 事务使用场景
- 跨集合的数据操作
- 复杂的业务逻辑处理
- 数据一致性保证

## 7. 性能优化

### 7.1 索引策略
- **用户集合**：username, email唯一索引
- **项目集合**：owner_id, status复合索引
- **时间字段**：created_at, updated_at索引

### 7.2 查询优化
- 支持字段选择，减少网络传输
- 分页查询，避免大结果集
- 条件过滤，提高查询效率

### 7.3 连接管理
- 连接池复用
- 健康检查机制
- 超时控制

## 8. 错误处理

### 8.1 错误类型
- **连接错误**：数据库连接失败
- **查询错误**：查询语法或条件错误
- **数据错误**：数据验证失败
- **事务错误**：事务执行失败

### 8.2 错误处理策略
- 统一错误格式
- 详细错误日志
- 优雅降级处理

## 9. 测试策略

### 9.1 单元测试
- Mock Repository接口
- 测试业务逻辑
- 覆盖边界条件

### 9.2 集成测试
- 真实数据库环境
- 完整流程测试
- 性能基准测试

## 10. 部署和维护

### 10.1 配置管理
- 数据库连接配置
- 性能参数调优
- 监控指标设置

### 10.2 监控和日志
- 查询性能监控
- 错误率统计
- 连接池状态监控

## 11. 未来扩展

### 11.1 多数据源支持
- 支持MySQL实现
- 支持Redis缓存层
- 支持分布式数据库

### 11.2 功能增强
- 读写分离支持
- 数据分片支持
- 缓存集成

## 12. 总结

Repository层设计采用了现代化的架构模式，提供了：
- **统一的数据访问接口**
- **灵活的实现切换能力**
- **完整的事务支持**
- **优秀的性能表现**
- **良好的可测试性**

该设计为青羽后端系统提供了稳定、高效、可扩展的数据访问层基础。
