# Repository 层设计说明

> **版本**: v1.0
> **创建日期**: 2026-03-19
> **最后更新**: 2026-03-19

---

## 目录

1. [职责边界与依赖关系](#1-职责边界与依赖关系)
2. [命名与代码规范](#2-命名与代码规范)
3. [设计模式与最佳实践](#3-设计模式与最佳实践)
4. [接口与契约规范](#4-接口与契约规范)
5. [测试策略](#5-测试策略)
6. [完整代码示例](#6-完整代码示例)

---

## 1. 职责边界与依赖关系

### 1.1 核心职责

Repository 层是数据访问层，负责：

- **数据持久化**：封装所有数据库操作（MongoDB、Redis）
- **CRUD 操作**：提供标准的增删改查接口
- **查询构建**：构建复杂查询条件和聚合管道
- **ID 转换**：处理 string ID 与 ObjectID 之间的转换
- **缓存策略**：实现数据缓存和失效逻辑
- **事务支持**：提供数据库事务能力

### 1.2 依赖关系图

```
┌─────────────────────────────────────────────────────────┐
│                    上层依赖方                            │
│  ┌─────────────────────────────────────────────────┐   │
│  │                   Service 层                     │   │
│  └───────────────────────┬─────────────────────────┘   │
└──────────────────────────┼──────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   Repository 层                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │                 interfaces/                      │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐           │   │
│  │  │bookstore│ │ reader  │ │  auth   │ ...       │   │
│  │  └─────────┘ └─────────┘ └─────────┘           │   │
│  └───────────────────────┬─────────────────────────┘   │
│                          │                              │
│  ┌───────────────────────┴─────────────────────────┐   │
│  │                  mongodb/                        │   │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐           │   │
│  │  │  base/  │ │bookstore│ │ reader  │ ...       │   │
│  │  └─────────┘ └─────────┘ └─────────┘           │   │
│  └───────────────────────┬─────────────────────────┘   │
│                          │                              │
│  ┌───────────┐  ┌────────┴────────┐  ┌───────────┐    │
│  │  cache/   │  │   querybuilder/ │  │  redis/   │    │
│  └───────────┘  └─────────────────┘  └───────────┘    │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    外部依赖                              │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │  MongoDB Driver │  │   Redis Client  │              │
│  └─────────────────┘  └─────────────────┘              │
└─────────────────────────────────────────────────────────┘
```

### 1.3 层级边界

| 可以做 | 不应该做 |
|--------|----------|
| 数据库 CRUD 操作 | 包含业务逻辑 |
| 查询条件构建 | 调用 Service 层 |
| ID 类型转换 | 处理 HTTP 请求/响应 |
| 缓存读写 | 包含验证逻辑 |
| 事务管理 | 直接返回 HTTP 错误码 |
| 聚合查询 | 修改 Models 定义 |

### 1.4 目录结构

```
repository/
├── cache/                  # 缓存层
│   ├── cached_repository.go
│   └── metrics.go
├── interfaces/             # 接口定义
│   ├── admin/              # 管理员接口
│   ├── ai/                 # AI服务接口
│   ├── auth/               # 认证接口
│   ├── bookstore/          # 书城接口
│   ├── infrastructure/     # 基础设施接口
│   ├── reader/             # 阅读接口
│   └── ...
├── mongodb/                # MongoDB 实现
│   ├── base/               # 基础 Repository
│   ├── bookstore/          # 书城实现
│   ├── reader/             # 阅读实现
│   └── factory.go          # 工厂函数
├── querybuilder/           # 查询构建器
│   └── mongo_query_builder.go
├── redis/                  # Redis 实现
├── search/                 # 搜索实现
├── errors.go               # 错误定义
└── id_converter.go         # ID 转换工具
```

---

## 2. 命名与代码规范

### 2.1 文件命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 接口定义 | 实体名 + `Repository_interface.go` | `BookStoreRepository_interface.go` |
| MongoDB 实现 | 实体名 + `_repository_mongo.go` | `bookstore_repository_mongo.go` |
| Redis 实现 | 实体名 + `_repository_redis.go` | `session_repository_redis.go` |
| 缓存装饰器 | `cached_` + 实体名 + `_repository.go` | `cached_book_repository.go` |
| 测试文件 | 原文件名 + `_test.go` | `bookstore_repository_test.go` |

### 2.2 接口命名

| 类型 | 规范 | 示例 |
|------|------|------|
| Repository 接口 | 实体名 + `Repository` | `BookRepository`, `ChapterRepository` |
| 基础接口 | 功能描述 | `CRUDRepository`, `HealthRepository` |
| 工厂接口 | 领域名 + `Factory` | `BookstoreFactory` |

### 2.3 实现命名

| 类型 | 规范 | 示例 |
|------|------|------|
| MongoDB 实现 | `Mongo` + 接口名 | `MongoBookRepository` |
| Redis 实现 | `Redis` + 接口名 | `RedisSessionRepository` |
| 缓存装饰器 | `Cached` + 接口名 | `CachedBookRepository` |
| 基类 | `Base` + `MongoRepository` | `BaseMongoRepository` |

### 2.4 方法命名

| 操作类型 | 命名规范 | 示例 |
|----------|----------|------|
| 创建 | `Create` | `Create(ctx, entity)` |
| 批量创建 | `BatchCreate` / `CreateMany` | `BatchCreate(ctx, entities)` |
| 按ID查询 | `GetByID` / `FindByID` | `GetByID(ctx, id)` |
| 按条件查询 | `GetBy` + 条件 | `GetByCategory(ctx, categoryID)` |
| 列表查询 | `List` / `GetAll` | `List(ctx, filter)` |
| 更新 | `Update` | `Update(ctx, id, updates)` |
| 批量更新 | `BatchUpdate` + 字段 | `BatchUpdateStatus(ctx, ids, status)` |
| 删除 | `Delete` | `Delete(ctx, id)` |
| 统计 | `Count` + 条件 | `CountByStatus(ctx, status)` |
| 存在检查 | `Exists` + 条件 | `Exists(ctx, id)` |
| 健康检查 | `Health` | `Health(ctx)` |

---

## 3. 设计模式与最佳实践

### 3.1 Repository 接口模式

接口定义在 `interfaces/` 目录，实现放在 `mongodb/` 目录：

```go
// interfaces/bookstore/BookStoreRepository_interface.go
package bookstore

import (
    "context"
    "Qingyu_backend/models/bookstore"
    base "Qingyu_backend/repository/interfaces/infrastructure"
)

// BookRepository 书籍仓储接口
type BookRepository interface {
    // 继承基础CRUD接口
    base.CRUDRepository[*bookstore.Book, string]
    base.HealthRepository

    // 领域特定方法
    GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*bookstore.Book, error)
    Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Book, error)
    CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error)
}
```

### 3.2 基类嵌入模式

MongoDB 实现嵌入 `BaseMongoRepository` 获取通用方法：

```go
// mongodb/bookstore/bookstore_repository_mongo.go
package mongodb

import (
    "Qingyu_backend/repository/mongodb/base"
    BookstoreInterface "Qingyu_backend/repository/interfaces/bookstore"
)

// MongoBookRepository MongoDB书籍仓储实现
type MongoBookRepository struct {
    *base.BaseMongoRepository  // 嵌入基类
    client *mongo.Client
}

// NewMongoBookRepository 创建实例
func NewMongoBookRepository(client *mongo.Client, database string) BookstoreInterface.BookRepository {
    db := client.Database(database)
    return &MongoBookRepository{
        BaseMongoRepository: base.NewBaseMongoRepository(db, "books"),
        client:              client,
    }
}
```

### 3.3 ID 转换模式

使用统一的 ID 转换函数：

```go
// repository/id_converter.go

// ParseID 解析必需的ID，空字符串返回错误
func ParseID(id string) (primitive.ObjectID, error) {
    if id == "" {
        return primitive.NilObjectID, ErrEmptyID
    }
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return primitive.NilObjectID, fmt.Errorf("%w: %s", ErrInvalidIDFormat, id)
    }
    return oid, nil
}

// ParseOptionalID 解析可选ID，空字符串返回nil
func ParseOptionalID(id string) (*primitive.ObjectID, error) {
    if id == "" {
        return nil, nil
    }
    oid, err := ParseID(id)
    if err != nil {
        return nil, err
    }
    return &oid, nil
}

// ParseIDs 批量解析ID列表
func ParseIDs(ids []string) ([]primitive.ObjectID, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    result := make([]primitive.ObjectID, 0, len(ids))
    for i, id := range ids {
        oid, err := ParseID(id)
        if err != nil {
            return nil, fmt.Errorf("ids[%d]: %w", i, err)
        }
        result = append(result, oid)
    }
    return result, nil
}
```

### 3.4 缓存装饰器模式

为 Repository 添加缓存层：

```go
// repository/cache/cached_repository.go
type CachedBookRepository struct {
    repo   BookRepository
    cache  redis.Client
    ttl    time.Duration
}

func NewCachedBookRepository(repo BookRepository, cache redis.Client) *CachedBookRepository {
    return &CachedBookRepository{
        repo:  repo,
        cache: cache,
        ttl:   5 * time.Minute,
    }
}

func (c *CachedBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("book:%s", id)
    cached, err := c.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var book bookstore.Book
        if json.Unmarshal([]byte(cached), &book) == nil {
            return &book, nil
        }
    }

    // 2. 从数据库获取
    book, err := c.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    if data, err := json.Marshal(book); err == nil {
        c.cache.Set(ctx, cacheKey, data, c.ttl)
    }

    return book, nil
}
```

### 3.5 查询构建器模式

使用 Query Builder 构建复杂查询：

```go
// repository/querybuilder/mongo_query_builder.go
type MongoQueryBuilder struct {
    filter bson.M
    sort   bson.D
    limit  int64
    skip   int64
}

func NewMongoQueryBuilder() *MongoQueryBuilder {
    return &MongoQueryBuilder{
        filter: bson.M{},
    }
}

func (b *MongoQueryBuilder) WithCategory(categoryID string) *MongoQueryBuilder {
    if categoryID != "" {
        oid, _ := primitive.ObjectIDFromHex(categoryID)
        b.filter["category_ids"] = oid
    }
    return b
}

func (b *MongoQueryBuilder) WithStatus(status bookstore.BookStatus) *MongoQueryBuilder {
    if status != "" {
        b.filter["status"] = status
    }
    return b
}

func (b *MongoQueryBuilder) WithPriceRange(min, max float64) *MongoQueryBuilder {
    if min > 0 || max > 0 {
        priceFilter := bson.M{}
        if min > 0 {
            priceFilter["$gte"] = min
        }
        if max > 0 {
            priceFilter["$lte"] = max
        }
        b.filter["price"] = priceFilter
    }
    return b
}

func (b *MongoQueryBuilder) SortBy(field string, ascending bool) *MongoQueryBuilder {
    order := 1
    if !ascending {
        order = -1
    }
    b.sort = append(b.sort, bson.E{Key: field, Value: order})
    return b
}

func (b *MongoQueryBuilder) Build() (bson.M, *options.FindOptions) {
    opts := options.Find()
    if len(b.sort) > 0 {
        opts.SetSort(b.sort)
    }
    if b.limit > 0 {
        opts.SetLimit(b.limit)
    }
    if b.skip > 0 {
        opts.SetSkip(b.skip)
    }
    return b.filter, opts
}
```

### 3.6 反模式警示

| 反模式 | 问题 | 正确做法 |
|--------|------|----------|
| 在 Repository 中写业务逻辑 | 违反单一职责 | 业务逻辑放 Service 层 |
| 直接返回 HTTP 错误 | 耦合表示层 | 返回领域错误，由 API 层转换 |
| 硬编码 SQL/查询 | 难以维护 | 使用 Query Builder |
| 忽略错误处理 | 导致数据不一致 | 始终处理并包装错误 |
| 在循环中查询 | N+1 问题 | 使用批量查询 |

---

## 4. 接口与契约规范

### 4.1 基础接口定义

```go
// repository/interfaces/infrastructure/base.go
package infrastructure

import "context"

// CRUDRepository 通用CRUD接口
type CRUDRepository[T any, ID any] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id ID) (*T, error)
    Update(ctx context.Context, id ID, updates map[string]interface{}) error
    Delete(ctx context.Context, id ID) error
    List(ctx context.Context, filter Filter) ([]*T, error)
    Count(ctx context.Context, filter Filter) (int64, error)
    Exists(ctx context.Context, id ID) (bool, error)
}

// HealthRepository 健康检查接口
type HealthRepository interface {
    Health(ctx context.Context) error
}

// Filter 通用过滤条件
type Filter map[string]interface{}
```

### 4.2 返回值规范

| 场景 | 返回值 | 说明 |
|------|--------|------|
| 单条查询 | `(*Entity, error)` | 未找到返回 `nil, nil` |
| 列表查询 | `([]*Entity, error)` | 空列表返回 `[]*Entity{}, nil` |
| 创建操作 | `error` | 通过指针修改实体ID |
| 更新操作 | `error` | 包含 "not found" 错误 |
| 删除操作 | `error` | 包含 "not found" 错误 |
| 统计操作 | `(int64, error)` | - |
| 存在检查 | `(bool, error)` | - |

### 4.3 错误定义

```go
// repository/errors.go
package repository

import "errors"

var (
    // ErrEmptyID ID为空
    ErrEmptyID = errors.New("ID cannot be empty")

    // ErrInvalidIDFormat ID格式无效
    ErrInvalidIDFormat = errors.New("invalid ID format")

    // ErrNotFound 记录未找到
    ErrNotFound = errors.New("record not found")

    // ErrDuplicateKey 重复键
    ErrDuplicateKey = errors.New("duplicate key error")
)

// IsIDError 判断是否为ID相关错误
func IsIDError(err error) bool {
    return errors.Is(err, ErrEmptyID) || errors.Is(err, ErrInvalidIDFormat)
}
```

### 4.4 分页规范

```go
// 分页参数
type Pagination struct {
    Limit  int    `json:"limit"`           // 每页数量，默认 20
    Offset int    `json:"offset"`          // 偏移量
    SortBy string `json:"sortBy"`          // 排序字段
    Order  string `json:"order"`           // asc / desc
}

// 分页结果
type PaginatedResult[T any] struct {
    Items      []*T   `json:"items"`
    Total      int64  `json:"total"`
    Page       int    `json:"page"`
    PageSize   int    `json:"pageSize"`
    TotalPages int    `json:"totalPages"`
}
```

---

## 5. 测试策略

### 5.1 测试文件组织

```
repository/
├── mongodb/
│   └── bookstore/
│       ├── bookstore_repository_mongo.go
│       ├── bookstore_repository_test.go      # 集成测试
│       └── book_benchmark_test.go            # 性能测试
└── id_converter_test.go                       # 单元测试
```

### 5.2 测试类型

| 测试类型 | 覆盖内容 | 工具 |
|----------|----------|------|
| 单元测试 | ID 转换、查询构建 | 标准测试 |
| 集成测试 | MongoDB 操作 | testcontainers |
| 性能测试 | 查询性能 | benchmark |

### 5.3 测试示例

```go
// bookstore_repository_test.go
package mongodb

import (
    "context"
    "testing"

    "github.com/testcontainers/testcontainers-go"
    "go.mongodb.org/mongo-driver/mongo"
)

func TestBookRepository_Create(t *testing.T) {
    // 使用 testcontainers 启动 MongoDB
    ctx := context.Background()
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image: "mongo:7",
            ExposedPorts: []string{"27017/tcp"},
        },
    })
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }
    defer container.Terminate(ctx)

    // 创建 Repository
    client := createMongoClient(t, container)
    repo := NewMongoBookRepository(client, "testdb")

    // 测试创建
    book := &bookstore.Book{
        Title:  "测试书籍",
        Author: "测试作者",
        Status: bookstore.BookStatusDraft,
    }

    err = repo.Create(ctx, book)
    if err != nil {
        t.Fatalf("Create failed: %v", err)
    }

    if book.ID.IsZero() {
        t.Error("Expected ID to be set after creation")
    }
}

func TestBookRepository_GetByID_NotFound(t *testing.T) {
    repo := setupTestRepository(t)
    ctx := context.Background()

    // 查询不存在的ID
    book, err := repo.GetByID(ctx, "507f1f77bcf86cd799439011")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if book != nil {
        t.Error("Expected nil for non-existent book")
    }
}
```

### 5.4 测试覆盖率要求

| 类型 | 最低覆盖率 |
|------|-----------|
| CRUD 操作 | 90% |
| 查询构建 | 85% |
| ID 转换 | 100% |
| 错误处理 | 80% |

---

## 6. 完整代码示例

### 6.1 完整的 Repository 接口

```go
// interfaces/bookstore/ChapterRepository_interface.go
package bookstore

import (
    "context"

    "Qingyu_backend/models/bookstore"
    base "Qingyu_backend/repository/interfaces/infrastructure"
)

// ChapterRepository 章节仓储接口
type ChapterRepository interface {
    // 基础CRUD
    base.CRUDRepository[*bookstore.Chapter, string]
    base.HealthRepository

    // 按书籍查询
    GetByBookID(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error)
    GetByBookIDAndNumber(ctx context.Context, bookID string, chapterNum int) (*bookstore.Chapter, error)
    CountByBookID(ctx context.Context, bookID string) (int64, error)

    // 状态查询
    GetPublishedByBook(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error)
    GetVipChaptersByBook(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error)

    // 批量操作
    BatchUpdateStatus(ctx context.Context, chapterIDs []string, status bookstore.ChapterStatus) error
    BatchUpdateVipStatus(ctx context.Context, chapterIDs []string, isVip bool) error

    // 排序操作
    ReorderChapters(ctx context.Context, bookID string, chapterIDs []string) error

    // 统计
    GetTotalWordCount(ctx context.Context, bookID string) (int64, error)

    // 事务
    Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### 6.2 完整的 MongoDB 实现

```go
// mongodb/bookstore/chapter_repository_mongo.go
package mongodb

import (
    "context"
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "Qingyu_backend/models/bookstore"
    "Qingyu_backend/repository/mongodb/base"
    ChapterInterface "Qingyu_backend/repository/interfaces/bookstore"
)

// MongoChapterRepository MongoDB章节仓储实现
type MongoChapterRepository struct {
    *base.BaseMongoRepository
    client *mongo.Client
}

// NewMongoChapterRepository 创建章节仓储实例
func NewMongoChapterRepository(client *mongo.Client, database string) ChapterInterface.ChapterRepository {
    db := client.Database(database)
    return &MongoChapterRepository{
        BaseMongoRepository: base.NewBaseMongoRepository(db, "chapters"),
        client:              client,
    }
}

// Create 创建章节
func (r *MongoChapterRepository) Create(ctx context.Context, chapter *bookstore.Chapter) error {
    if chapter == nil {
        return errors.New("chapter cannot be nil")
    }

    now := time.Now()
    chapter.CreatedAt = now
    chapter.UpdatedAt = now

    if chapter.ID.IsZero() {
        chapter.ID = primitive.NewObjectID()
    }

    _, err := r.GetCollection().InsertOne(ctx, chapter)
    return err
}

// GetByID 根据ID获取章节
func (r *MongoChapterRepository) GetByID(ctx context.Context, id string) (*bookstore.Chapter, error) {
    oid, err := r.ParseID(id)
    if err != nil {
        return nil, err
    }

    var chapter bookstore.Chapter
    err = r.GetCollection().FindOne(ctx, bson.M{"_id": oid}).Decode(&chapter)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }

    return &chapter, nil
}

// Update 更新章节
func (r *MongoChapterRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
    oid, err := r.ParseID(id)
    if err != nil {
        return err
    }

    updates["updated_at"] = time.Now()

    result, err := r.GetCollection().UpdateOne(
        ctx,
        bson.M{"_id": oid},
        bson.M{"$set": updates},
    )
    if err != nil {
        return err
    }

    if result.MatchedCount == 0 {
        return errors.New("chapter not found")
    }

    return nil
}

// Delete 删除章节
func (r *MongoChapterRepository) Delete(ctx context.Context, id string) error {
    oid, err := r.ParseID(id)
    if err != nil {
        return err
    }

    result, err := r.GetCollection().DeleteOne(ctx, bson.M{"_id": oid})
    if err != nil {
        return err
    }

    if result.DeletedCount == 0 {
        return errors.New("chapter not found")
    }

    return nil
}

// GetByBookID 根据书籍ID获取章节列表
func (r *MongoChapterRepository) GetByBookID(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
    bookOID, err := r.ParseID(bookID)
    if err != nil {
        return nil, err
    }

    opts := options.Find().
        SetSort(bson.D{{Key: "chapter_num", Value: 1}}).
        SetSkip(int64(offset)).
        SetLimit(int64(limit))

    cursor, err := r.GetCollection().Find(ctx, bson.M{"book_id": bookOID}, opts)
    if err != nil {
        return nil, err
    }

    var chapters []*bookstore.Chapter
    if err := cursor.All(ctx, &chapters); err != nil {
        return nil, err
    }

    return chapters, nil
}

// CountByBookID 统计书籍章节数
func (r *MongoChapterRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
    bookOID, err := r.ParseID(bookID)
    if err != nil {
        return 0, err
    }

    return r.GetCollection().CountDocuments(ctx, bson.M{"book_id": bookOID})
}

// Health 健康检查
func (r *MongoChapterRepository) Health(ctx context.Context) error {
    return r.GetDB().Client().Ping(ctx, nil)
}

// Transaction 事务支持
func (r *MongoChapterRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
    session, err := r.client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        return nil, fn(sessCtx)
    })

    return err
}
```

### 6.3 BaseMongoRepository 关键方法

```go
// mongodb/base/base_repository.go (关键部分)
package base

type BaseMongoRepository struct {
    db         *mongo.Database
    collection *mongo.Collection
}

func NewBaseMongoRepository(db *mongo.Database, collectionName string) *BaseMongoRepository {
    return &BaseMongoRepository{
        db:         db,
        collection: db.Collection(collectionName),
    }
}

func (b *BaseMongoRepository) GetDB() *mongo.Database {
    return b.db
}

func (b *BaseMongoRepository) GetCollection() *mongo.Collection {
    return b.collection
}

func (b *BaseMongoRepository) ParseID(id string) (primitive.ObjectID, error) {
    return types.ParseObjectID(id)
}

func (b *BaseMongoRepository) FindByID(ctx context.Context, id string, result interface{}) error {
    oid, err := b.ParseID(id)
    if err != nil {
        return err
    }
    return b.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(result)
}

func (b *BaseMongoRepository) UpdateByID(ctx context.Context, id string, update bson.M) error {
    oid, err := b.ParseID(id)
    if err != nil {
        return err
    }
    _, err = b.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
    return err
}

func (b *BaseMongoRepository) DeleteByID(ctx context.Context, id string) error {
    oid, err := b.ParseID(id)
    if err != nil {
        return err
    }
    _, err = b.collection.DeleteOne(ctx, bson.M{"_id": oid})
    return err
}

func (b *BaseMongoRepository) Find(ctx context.Context, filter bson.M, results interface{}, opts ...*options.FindOptions) error {
    cursor, err := b.collection.Find(ctx, filter, opts...)
    if err != nil {
        return err
    }
    return cursor.All(ctx, results)
}

func (b *BaseMongoRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
    return b.collection.CountDocuments(ctx, filter)
}

func (b *BaseMongoRepository) Exists(ctx context.Context, id string) (bool, error) {
    oid, err := b.ParseID(id)
    if err != nil {
        return false, err
    }
    count, err := b.collection.CountDocuments(ctx, bson.M{"_id": oid})
    return count > 0, err
}

func (b *BaseMongoRepository) Create(ctx context.Context, document interface{}) error {
    _, err := b.collection.InsertOne(ctx, document)
    return err
}
```

---

## 附录

### A. 相关文档

- [Models 层设计说明](./layer-models.md)
- [Service 层设计说明](./layer-service.md)
- [API 层设计说明](./layer-api.md)

### B. 参考资源

- [MongoDB Go Driver 文档](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [go-redis 文档](https://redis.uptrace.dev/)
- [testcontainers-go](https://github.com/testcontainers/testcontainers-go)

---

*最后更新：2026-03-19*
