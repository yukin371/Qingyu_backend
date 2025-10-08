# MongoDB应用层优化设计

## 1. 概述

本文档针对青羽写作平台使用MongoDB存储用户数据的潜在问题，提出基于Repository模式的应用层验证和事务处理优化方案，同时预留混合架构改进空间。Repository层作为数据访问的统一抽象，是解决当前架构不一致问题的关键。

## 2. 当前问题分析

### 2.1 架构不一致问题

- **混合数据访问模式**：部分服务直接使用`global.DB`，部分使用Repository接口
- **缺乏统一抽象**：数据访问逻辑分散在Service层和Repository层
- **依赖关系混乱**：Service层直接依赖具体数据库实现

### 2.2 数据完整性问题

- 缺乏外键约束，容易产生孤儿数据
- 用户删除后相关数据清理复杂
- 数据一致性主要依赖应用层保证

### 2.3 事务处理问题

- 跨集合操作的事务支持有限
- 用户注册、权限变更等关键操作需要强一致性
- 复杂业务逻辑的原子性难以保证

### 2.4 数据验证问题

- Schema验证不够严格
- 数据类型约束依赖应用层
- 业务规则验证分散在各个服务中

## 3. Repository模式优化设计

### 3.1 Repository层架构

Repository层作为数据访问的统一抽象层，解决当前架构不一致的问题：

```
┌─────────────────┐
│   Service层     │  ← 业务逻辑处理
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  Repository层   │  ← 数据访问抽象
│ - 接口定义      │
│ - 查询封装      │
│ - 事务管理      │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   数据存储层    │  ← MongoDB/PostgreSQL
└─────────────────┘
```

### 3.2 核心Repository接口设计

```go
// 基础Repository接口
type BaseRepository[T any] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id string) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter interface{}) ([]*T, error)
}

// 文档Repository接口
type DocumentRepository interface {
    BaseRepository[models.Document]
    GetByProjectID(ctx context.Context, projectID string) ([]*models.Document, error)
    GetByUserID(ctx context.Context, userID string) ([]*models.Document, error)
    UpdateContent(ctx context.Context, id string, content string) error
}

// 项目Repository接口
type ProjectRepository interface {
    BaseRepository[models.Project]
    GetByUserID(ctx context.Context, userID string) ([]*models.Project, error)
    GetWithDocuments(ctx context.Context, id string) (*models.ProjectWithDocuments, error)
}
```

### 3.3 MongoDB实现

```go
// MongoDB文档Repository实现
type MongoDocumentRepository struct {
    collection *mongo.Collection
}

func (r *MongoDocumentRepository) Create(ctx context.Context, doc *models.Document) error {
    doc.CreatedAt = time.Now()
    doc.UpdatedAt = time.Now()
    
    result, err := r.collection.InsertOne(ctx, doc)
    if err != nil {
        return fmt.Errorf("创建文档失败: %w", err)
    }
    
    doc.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

func (r *MongoDocumentRepository) GetByProjectID(ctx context.Context, projectID string) ([]*models.Document, error) {
    objID, err := primitive.ObjectIDFromHex(projectID)
    if err != nil {
        return nil, fmt.Errorf("无效的项目ID: %w", err)
    }
    
    filter := bson.M{"project_id": objID}
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("查询文档失败: %w", err)
    }
    defer cursor.Close(ctx)
    
    var documents []*models.Document
    if err := cursor.All(ctx, &documents); err != nil {
        return nil, fmt.Errorf("解析文档数据失败: %w", err)
    }
    
    return documents, nil
}
```

## 4. 应用层验证设计

### 4.1 Schema验证框架

#### 4.1.1 集合级别验证

```javascript
// 用户集合验证规则
{
  "validator": {
    "$jsonSchema": {
      "bsonType": "object",
      "required": ["username", "email", "password_hash", "created_at"],
      "properties": {
        "username": {
          "bsonType": "string",
          "minLength": 3,
          "maxLength": 30,
          "pattern": "^[a-zA-Z0-9_]+$"
        },
        "email": {
          "bsonType": "string",
          "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
        },
        "status": {
          "enum": ["active", "inactive", "suspended", "deleted"]
        },
        "created_at": {
          "bsonType": "date"
        },
        "updated_at": {
          "bsonType": "date"
        }
      }
    }
  },
  "validationLevel": "strict",
  "validationAction": "error"
}
```

#### 3.1.2 应用层验证器

```go
// 用户数据验证器
type UserValidator struct {
    db *mongo.Database
}

func (v *UserValidator) ValidateCreate(user *User) error {
    // 1. 基础字段验证
    if err := v.validateBasicFields(user); err != nil {
        return err
    }
    
    // 2. 唯一性验证
    if err := v.validateUniqueness(user); err != nil {
        return err
    }
    
    // 3. 业务规则验证
    if err := v.validateBusinessRules(user); err != nil {
        return err
    }
    
    return nil
}

func (v *UserValidator) validateUniqueness(user *User) error {
    // 检查用户名唯一性
    count, err := v.db.Collection("users").CountDocuments(
        context.Background(),
        bson.M{"username": user.Username},
    )
    if err != nil {
        return err
    }
    if count > 0 {
        return errors.New("用户名已存在")
    }
    
    // 检查邮箱唯一性
    count, err = v.db.Collection("users").CountDocuments(
        context.Background(),
        bson.M{"email": user.Email},
    )
    if err != nil {
        return err
    }
    if count > 0 {
        return errors.New("邮箱已被使用")
    }
    
    return nil
}
```

### 4.2 数据完整性保证

#### 3.2.1 引用完整性检查

```go
// 引用完整性管理器
type ReferenceIntegrityManager struct {
    db *mongo.Database
}

func (rim *ReferenceIntegrityManager) ValidateReferences(doc interface{}, collection string) error {
    switch collection {
    case "projects":
        return rim.validateProjectReferences(doc.(*Project))
    case "nodes":
        return rim.validateNodeReferences(doc.(*Node))
    case "user_roles":
        return rim.validateUserRoleReferences(doc.(*UserRole))
    }
    return nil
}

func (rim *ReferenceIntegrityManager) validateProjectReferences(project *Project) error {
    // 验证创建者存在
    count, err := rim.db.Collection("users").CountDocuments(
        context.Background(),
        bson.M{"_id": project.CreatorID, "status": bson.M{"$ne": "deleted"}},
    )
    if err != nil {
        return err
    }
    if count == 0 {
        return errors.New("创建者不存在或已被删除")
    }
    
    return nil
}
```

#### 3.2.2 级联操作管理

```go
// 级联操作管理器
type CascadeManager struct {
    db *mongo.Database
}

func (cm *CascadeManager) DeleteUser(userID primitive.ObjectID) error {
    session, err := cm.db.Client().StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(context.Background())
    
    return mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
        // 1. 软删除用户
        _, err := cm.db.Collection("users").UpdateOne(sc,
            bson.M{"_id": userID},
            bson.M{"$set": bson.M{
                "status": "deleted",
                "deleted_at": time.Now(),
            }},
        )
        if err != nil {
            return err
        }
        
        // 2. 处理用户项目
        _, err = cm.db.Collection("projects").UpdateMany(sc,
            bson.M{"creator_id": userID},
            bson.M{"$set": bson.M{
                "status": "archived",
                "archived_reason": "creator_deleted",
            }},
        )
        if err != nil {
            return err
        }
        
        // 3. 清理用户角色
        _, err = cm.db.Collection("user_roles").DeleteMany(sc,
            bson.M{"user_id": userID},
        )
        if err != nil {
            return err
        }
        
        return nil
    })
}
```

## 5. 事务处理设计

### 5.1 事务管理框架

#### 4.1.1 事务管理器

```go
// 事务管理器
type TransactionManager struct {
    client *mongo.Client
    db     *mongo.Database
}

func (tm *TransactionManager) ExecuteTransaction(operations []TransactionOperation) error {
    session, err := tm.client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(context.Background())
    
    return mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
        for _, op := range operations {
            if err := op.Execute(sc, tm.db); err != nil {
                return err
            }
        }
        return nil
    })
}

// 事务操作接口
type TransactionOperation interface {
    Execute(ctx mongo.SessionContext, db *mongo.Database) error
    Rollback(ctx mongo.SessionContext, db *mongo.Database) error
}
```

#### 4.1.2 用户注册事务

```go
// 用户注册事务操作
type UserRegistrationTransaction struct {
    User     *User
    UserRole *UserRole
}

func (urt *UserRegistrationTransaction) Execute(ctx mongo.SessionContext, db *mongo.Database) error {
    // 1. 创建用户
    result, err := db.Collection("users").InsertOne(ctx, urt.User)
    if err != nil {
        return err
    }
    urt.User.ID = result.InsertedID.(primitive.ObjectID)
    
    // 2. 分配默认角色
    urt.UserRole.UserID = urt.User.ID
    urt.UserRole.RoleID = primitive.ObjectIDFromHex("default_user_role_id")
    urt.UserRole.CreatedAt = time.Now()
    
    _, err = db.Collection("user_roles").InsertOne(ctx, urt.UserRole)
    if err != nil {
        return err
    }
    
    // 3. 初始化用户配置
    userConfig := &UserConfig{
        UserID:    urt.User.ID,
        Theme:     "default",
        Language:  "zh-CN",
        CreatedAt: time.Now(),
    }
    
    _, err = db.Collection("user_configs").InsertOne(ctx, userConfig)
    return err
}
```

### 5.2 分布式事务处理

#### 4.2.1 Saga模式实现

```go
// Saga事务管理器
type SagaManager struct {
    db *mongo.Database
}

type SagaStep struct {
    Name        string
    Execute     func(ctx context.Context) error
    Compensate  func(ctx context.Context) error
}

func (sm *SagaManager) ExecuteSaga(steps []SagaStep) error {
    executedSteps := make([]SagaStep, 0)
    
    for _, step := range steps {
        if err := step.Execute(context.Background()); err != nil {
            // 执行补偿操作
            for i := len(executedSteps) - 1; i >= 0; i-- {
                if compensateErr := executedSteps[i].Compensate(context.Background()); compensateErr != nil {
                    log.Printf("补偿操作失败: %v", compensateErr)
                }
            }
            return err
        }
        executedSteps = append(executedSteps, step)
    }
    
    return nil
}
```

## 6. 性能优化设计

### 6.1 查询优化

#### 5.1.1 索引策略

```go
// 索引管理器
type IndexManager struct {
    db *mongo.Database
}

func (im *IndexManager) EnsureIndexes() error {
    indexes := map[string][]mongo.IndexModel{
        "users": {
            {Keys: bson.D{{"username", 1}}, Options: options.Index().SetUnique(true)},
            {Keys: bson.D{{"email", 1}}, Options: options.Index().SetUnique(true)},
            {Keys: bson.D{{"status", 1}, {"created_at", -1}}},
        },
        "user_roles": {
            {Keys: bson.D{{"user_id", 1}, {"role_id", 1}}, Options: options.Index().SetUnique(true)},
            {Keys: bson.D{{"user_id", 1}}},
        },
        "projects": {
            {Keys: bson.D{{"creator_id", 1}, {"status", 1}}},
            {Keys: bson.D{{"created_at", -1}}},
        },
    }
    
    for collection, indexModels := range indexes {
        _, err := im.db.Collection(collection).Indexes().CreateMany(
            context.Background(),
            indexModels,
        )
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

#### 5.1.2 聚合查询优化

```go
// 用户权限查询优化
func (us *UserService) GetUserPermissions(userID primitive.ObjectID) ([]string, error) {
    pipeline := []bson.M{
        {"$match": bson.M{"user_id": userID}},
        {"$lookup": bson.M{
            "from":         "roles",
            "localField":   "role_id",
            "foreignField": "_id",
            "as":           "role",
        }},
        {"$unwind": "$role"},
        {"$lookup": bson.M{
            "from":         "role_permissions",
            "localField":   "role._id",
            "foreignField": "role_id",
            "as":           "role_permissions",
        }},
        {"$unwind": "$role_permissions"},
        {"$lookup": bson.M{
            "from":         "permissions",
            "localField":   "role_permissions.permission_id",
            "foreignField": "_id",
            "as":           "permission",
        }},
        {"$unwind": "$permission"},
        {"$group": bson.M{
            "_id":         "$user_id",
            "permissions": bson.M{"$addToSet": "$permission.name"},
        }},
    }
    
    cursor, err := us.db.Collection("user_roles").Aggregate(context.Background(), pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    
    var result struct {
        Permissions []string `bson:"permissions"`
    }
    
    if cursor.Next(context.Background()) {
        if err := cursor.Decode(&result); err != nil {
            return nil, err
        }
    }
    
    return result.Permissions, nil
}
```

### 6.2 缓存策略

#### 5.2.1 用户权限缓存

```go
// 权限缓存管理器
type PermissionCacheManager struct {
    cache map[string][]string
    mutex sync.RWMutex
    ttl   time.Duration
}

func (pcm *PermissionCacheManager) GetUserPermissions(userID string) ([]string, bool) {
    pcm.mutex.RLock()
    defer pcm.mutex.RUnlock()
    
    permissions, exists := pcm.cache[userID]
    return permissions, exists
}

func (pcm *PermissionCacheManager) SetUserPermissions(userID string, permissions []string) {
    pcm.mutex.Lock()
    defer pcm.mutex.Unlock()
    
    pcm.cache[userID] = permissions
    
    // 设置过期时间
    go func() {
        time.Sleep(pcm.ttl)
        pcm.mutex.Lock()
        delete(pcm.cache, userID)
        pcm.mutex.Unlock()
    }()
}
```

## 7. 混合架构预留设计

### 7.1 数据访问层抽象

#### 6.1.1 仓储模式接口

```go
// 用户仓储接口
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter UserFilter) ([]*User, error)
}

// MongoDB实现
type MongoUserRepository struct {
    db *mongo.Database
}

func (r *MongoUserRepository) Create(ctx context.Context, user *User) error {
    // MongoDB实现
}

// PostgreSQL实现（预留）
type PostgreSQLUserRepository struct {
    db *sql.DB
}

func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *User) error {
    // PostgreSQL实现
}
```

#### 6.1.2 仓储工厂

```go
// 仓储工厂
type RepositoryFactory interface {
    CreateUserRepository() UserRepository
    CreateProjectRepository() ProjectRepository
    CreateRoleRepository() RoleRepository
}

// MongoDB工厂
type MongoRepositoryFactory struct {
    db *mongo.Database
}

func (f *MongoRepositoryFactory) CreateUserRepository() UserRepository {
    return &MongoUserRepository{db: f.db}
}

// PostgreSQL工厂（预留）
type PostgreSQLRepositoryFactory struct {
    db *sql.DB
}

func (f *PostgreSQLRepositoryFactory) CreateUserRepository() UserRepository {
    return &PostgreSQLUserRepository{db: f.db}
}
```

### 7.2 数据迁移框架

#### 6.2.1 迁移接口设计

```go
// 数据迁移接口
type DataMigrator interface {
    MigrateUsers(ctx context.Context) error
    MigrateProjects(ctx context.Context) error
    MigrateRoles(ctx context.Context) error
    Rollback(ctx context.Context) error
}

// MongoDB到PostgreSQL迁移器
type MongoToPostgreSQLMigrator struct {
    mongoRepo    RepositoryFactory
    postgresRepo RepositoryFactory
}

func (m *MongoToPostgreSQLMigrator) MigrateUsers(ctx context.Context) error {
    // 1. 从MongoDB读取用户数据
    mongoUserRepo := m.mongoRepo.CreateUserRepository()
    users, err := mongoUserRepo.List(ctx, UserFilter{})
    if err != nil {
        return err
    }
    
    // 2. 写入PostgreSQL
    postgresUserRepo := m.postgresRepo.CreateUserRepository()
    for _, user := range users {
        if err := postgresUserRepo.Create(ctx, user); err != nil {
            return err
        }
    }
    
    return nil
}
```

### 7.3 配置管理

#### 6.3.1 数据库配置

```go
// 数据库配置
type DatabaseConfig struct {
    Type     string `yaml:"type"`     // "mongodb" 或 "postgresql"
    MongoDB  MongoConfig  `yaml:"mongodb"`
    PostgreSQL PostgreSQLConfig `yaml:"postgresql"`
}

type MongoConfig struct {
    URI      string `yaml:"uri"`
    Database string `yaml:"database"`
}

type PostgreSQLConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Database string `yaml:"database"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}
```

#### 6.3.2 动态切换支持

```go
// 数据库管理器
type DatabaseManager struct {
    config *DatabaseConfig
    factory RepositoryFactory
}

func NewDatabaseManager(config *DatabaseConfig) (*DatabaseManager, error) {
    var factory RepositoryFactory
    
    switch config.Type {
    case "mongodb":
        client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoDB.URI))
        if err != nil {
            return nil, err
        }
        factory = &MongoRepositoryFactory{db: client.Database(config.MongoDB.Database)}
    case "postgresql":
        db, err := sql.Open("postgres", buildPostgreSQLDSN(config.PostgreSQL))
        if err != nil {
            return nil, err
        }
        factory = &PostgreSQLRepositoryFactory{db: db}
    default:
        return nil, errors.New("不支持的数据库类型")
    }
    
    return &DatabaseManager{
        config:  config,
        factory: factory,
    }, nil
}
```

## 8. 监控和告警

### 8.1 数据质量监控

```go
// 数据质量监控器
type DataQualityMonitor struct {
    db *mongo.Database
}

func (dqm *DataQualityMonitor) CheckDataIntegrity() error {
    // 1. 检查孤儿数据
    if err := dqm.checkOrphanedData(); err != nil {
        return err
    }
    
    // 2. 检查数据一致性
    if err := dqm.checkDataConsistency(); err != nil {
        return err
    }
    
    // 3. 检查索引健康状态
    if err := dqm.checkIndexHealth(); err != nil {
        return err
    }
    
    return nil
}

func (dqm *DataQualityMonitor) checkOrphanedData() error {
    // 检查项目中是否有不存在的创建者
    pipeline := []bson.M{
        {"$lookup": bson.M{
            "from":         "users",
            "localField":   "creator_id",
            "foreignField": "_id",
            "as":           "creator",
        }},
        {"$match": bson.M{"creator": bson.M{"$size": 0}}},
        {"$count": "orphaned_projects"},
    }
    
    cursor, err := dqm.db.Collection("projects").Aggregate(context.Background(), pipeline)
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())
    
    var result struct {
        Count int `bson:"orphaned_projects"`
    }
    
    if cursor.Next(context.Background()) {
        if err := cursor.Decode(&result); err != nil {
            return err
        }
        if result.Count > 0 {
            log.Printf("发现 %d 个孤儿项目", result.Count)
        }
    }
    
    return nil
}
```

## 9. 实施计划

### 9.1 第一阶段：基础优化（2-3周）

1. 实现Schema验证框架
2. 添加应用层数据验证
3. 优化现有索引策略
4. 实现基础事务处理

### 9.2 第二阶段：高级功能（3-4周）

1. 实现级联操作管理
2. 添加权限缓存机制
3. 实现数据质量监控
4. 优化聚合查询性能

### 9.3 第三阶段：架构预留（2-3周）

1. 实现仓储模式抽象
2. 设计数据迁移框架
3. 添加配置管理支持
4. 实现监控告警系统

## 10. 风险控制

### 10.1 性能风险

- 事务处理可能影响性能
- 应用层验证增加响应时间
- 缓存策略需要合理设计

### 10.2 数据风险

- 迁移过程中的数据一致性
- 验证规则变更的向后兼容性
- 事务失败的回滚处理

### 10.3 运维风险

- 监控系统的复杂性增加
- 多数据库支持的维护成本
- 配置管理的复杂性

## 11. 总结

本设计方案通过应用层优化解决MongoDB在用户数据存储方面的问题，同时为未来的混合架构迁移预留空间。重点关注数据完整性、事务处理和性能优化，确保系统的稳定性和可扩展性。
