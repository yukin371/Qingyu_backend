# Repository层设计规范

**版本**: v2.0
**更新**: 2026-01-08
**状态**: ✅ 正式实施

---

## 一、职责定义

Repository层负责：
- ✅ 数据库操作封装
- ✅ 查询构建和优化
- ✅ 索引管理
- ✅ 缓存策略实现
- ✅ 事务处理
- ✅ 数据映射和转换

Repository层不负责：
- ❌ 业务逻辑
- ❌ 参数验证（基础验证除外）
- ❌ HTTP处理
- ❌ 业务规则判断

---

## 二、接口定义

### 2.1 基础接口

```go
// repository/interfaces/infrastructure/base_interface.go
type BaseRepository[T any, ID comparable] interface {
    // CRUD操作
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id ID) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id ID) error

    // 查询操作
    List(ctx context.Context, filter Filter) ([]*T, error)
    Count(ctx context.Context, filter Filter) (int64, error)

    // 工具方法
    Health(ctx context.Context) error
}
```

### 2.2 业务接口

```go
// repository/interfaces/bookstore/BookRepository_interface.go
type BookRepository interface {
    // 继承基础接口
    base.BaseRepository[*Book, string]

    // 书籍特定方法
    GetByTitle(ctx context.Context, title string) (*Book, error)
    GetByAuthor(ctx context.Context, author string) ([]*Book, error)
    FindWithFilter(ctx context.Context, filter *BookFilter) ([]*Book, int64, error)
    Search(ctx context.Context, keyword string, page, pageSize int) ([]*Book, int64, error)

    // 状态管理
    UpdateStatus(ctx context.Context, bookID string, status string) error
    GetBooksByStatus(ctx context.Context, status string) ([]*Book, error)

    // 事务支持
    Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

---

## 三、实现规范

### 3.1 结构定义

```go
// repository/mongodb/bookstore/bookstore_repository_mongo.go
type MongoBookRepository struct {
    db         *mongo.Database
    collection *mongo.Collection
    cache      cache.Cache
    logger     Logger
}

// 构造函数
func NewMongoBookRepository(
    db *mongo.Database,
    cache cache.Cache,
    logger Logger,
) BookRepository {
    return &MongoBookRepository{
        db:         db,
        collection: db.Collection("books"),
        cache:      cache,
        logger:     logger,
    }
}
```

### 3.2 CRUD实现

```go
// Create 创建书籍
func (r *MongoBookRepository) Create(ctx context.Context, book *Book) error {
    // 1. 设置默认值
    book.ID = primitive.NewObjectID().Hex()
    book.CreatedAt = time.Now()
    book.UpdatedAt = time.Now()

    // 2. 插入数据
    _, err := r.collection.InsertOne(ctx, book)
    if err != nil {
        if mongo.IsDuplicateKeyError(err) {
            return &RepositoryError{
                Code:    "DUPLICATE_KEY",
                Message: "书籍已存在",
                Err:     err,
            }
        }
        return &RepositoryError{
            Code:    "INTERNAL_ERROR",
            Message: "创建书籍失败",
            Err:     err,
        }
    }

    return nil
}

// GetByID 根据ID获取书籍
func (r *MongoBookRepository) GetByID(ctx context.Context, id string) (*Book, error) {
    var book Book
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, &RepositoryError{
            Code:    "INTERNAL_ERROR",
            Message: "查询书籍失败",
            Err:     err,
        }
    }

    return &book, nil
}
```

### 3.3 查询实现

```go
// FindWithFilter 根据条件查询
func (r *MongoBookRepository) FindWithFilter(
    ctx context.Context,
    filter *BookFilter,
) ([]*Book, int64, error) {
    // 1. 构建查询条件
    query := r.buildQuery(filter)

    // 2. 执行查询
    opts := r.buildQueryOptions(filter)
    cursor, err := r.collection.Find(ctx, query, opts)
    if err != nil {
        return nil, 0, &RepositoryError{
            Code:    "QUERY_ERROR",
            Message: "查询失败",
            Err:     err,
        }
    }
    defer cursor.Close(ctx)

    // 3. 解析结果
    var books []*Book
    if err = cursor.All(ctx, &books); err != nil {
        return nil, 0, &RepositoryError{
            Code:    "PARSE_ERROR",
            Message: "解析结果失败",
            Err:     err,
        }
    }

    // 4. 统计总数
    total, err := r.collection.CountDocuments(ctx, query)
    if err != nil {
        return nil, 0, &RepositoryError{
            Code:    "COUNT_ERROR",
            Message: "统计失败",
            Err:     err,
        }
    }

    return books, total, nil
}

// buildQuery 构建查询条件
func (r *MongoBookRepository) buildQuery(filter *BookFilter) bson.M {
    query := bson.M{}

    if filter.Status != "" {
        query["status"] = filter.Status
    }

    if filter.Author != "" {
        query["author"] = filter.Author
    }

    if len(filter.Tags) > 0 {
        query["tags"] = bson.M{"$in": filter.Tags}
    }

    return query
}
```

---

## 四、索引管理

### 4.1 索引定义

```go
// EnsureIndexes 创建索引
func (r *MongoBookRepository) EnsureIndexes(ctx context.Context) error {
    models := []mongo.IndexModel{
        // 单字段索引
        {Keys: bson.D{{Key: "title", Value: 1}}},
        {Keys: bson.D{{Key: "author", Value: 1}}},
        {Keys: bson.D{{Key: "status", Value: 1}}},

        // 复合索引
        {Keys: bson.D{
            {Key: "status", Value: 1},
            {Key: "created_at", Value: -1},
        }},

        // 文本索引
        {Keys: bson.D{
            {Key: "title", Value: "text"},
            {Key: "description", Value: "text"},
        }},

        // 唯一索引
        {
            Keys:    bson.D{{Key: "slug", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    }

    _, err := r.collection.Indexes().CreateMany(ctx, models)
    return err
}
```

### 4.2 索引优化

**原则**：
- ✅ 为常用查询字段创建索引
- ✅ 合理使用复合索引
- ✅ 避免过多索引影响写入性能
- ✅ 定期Review索引使用情况

---

## 五、缓存策略

### 5.1 缓存实现

```go
// GetByID 带缓存查询
func (r *MongoBookRepository) GetByID(ctx context.Context, id string) (*Book, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("book:%s", id)
    if r.cache != nil {
        cached, err := r.cache.Get(ctx, cacheKey)
        if err == nil && cached != nil {
            var book Book
            if err := json.Unmarshal([]byte(cached), &book); err == nil {
                return &book, nil
            }
        }
    }

    // 2. 从数据库查询
    book, err := r.getFromDB(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    if r.cache != nil && book != nil {
        data, _ := json.Marshal(book)
        r.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
    }

    return book, nil
}
```

### 5.2 缓存失效

```go
// Update 更新并失效缓存
func (r *MongoBookRepository) Update(ctx context.Context, book *Book) error {
    // 1. 更新数据库
    book.UpdatedAt = time.Now()
    _, err := r.collection.UpdateByID(ctx, book.ID, bson.M{"$set": book})
    if err != nil {
        return err
    }

    // 2. 删除缓存
    if r.cache != nil {
        cacheKey := fmt.Sprintf("book:%s", book.ID)
        r.cache.Delete(ctx, cacheKey)
    }

    return nil
}
```

---

## 六、事务处理

### 6.1 事务支持

```go
// Transaction 执行事务
func (r *MongoBookRepository) Transaction(
    ctx context.Context,
    fn func(ctx context.Context) error,
) error {
    session, err := r.db.Client().StartSession()
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

### 6.2 事务使用

```go
// Service层使用事务
func (s *BookService) PurchaseBook(ctx context.Context, userID, bookID string) error {
    return s.bookRepo.Transaction(ctx, func(txCtx context.Context) error {
        // 1. 扣除余额
        if err := s.userRepo.DeductBalance(txCtx, userID, price); err != nil {
            return err
        }

        // 2. 添加购买记录
        if err := s.purchaseRepo.Create(txCtx, purchase); err != nil {
            return err
        }

        return nil
    })
}
```

---

## 七、错误处理

### 7.1 错误定义

```go
type RepositoryError struct {
    Code    string
    Message string
    Err     error
}

func (e *RepositoryError) Error() string {
    return e.Message
}

func (e *RepositoryError) Unwrap() error {
    return e.Err
}
```

### 7.2 错误类型

| Code | 说明 | HTTP状态码 |
|------|------|-----------|
| NOT_FOUND | 资源不存在 | 404 |
| DUPLICATE_KEY | 唯一键冲突 | 409 |
| VALIDATION_ERROR | 参数验证失败 | 400 |
| INTERNAL_ERROR | 内部错误 | 500 |
| QUERY_ERROR | 查询错误 | 500 |

---

## 八、性能优化

### 8.1 查询优化

```go
// 只返回需要的字段
opts := options.Find().SetProjection(bson.M{
    "title": 1,
    "author": 1,
    "status": 1,
})

// 使用索引
query := bson.M{"status": "published"}  // status有索引

// 限制返回数量
opts.SetLimit(100)

// 使用聚合管道优化复杂查询
pipeline := mongo.Pipeline{
    {{"$match", filter}},
    {{"$sort", bson.D{{"created_at", -1}}}},
    {{"$skip", offset}},
    {{"$limit", limit}},
}
```

### 8.2 批量操作

```go
// BatchCreate 批量创建
func (r *MongoBookRepository) BatchCreate(ctx context.Context, books []*Book) error {
    docs := make([]interface{}, len(books))
    for i, book := range books {
        book.ID = primitive.NewObjectID().Hex()
        book.CreatedAt = time.Now()
        docs[i] = book
    }

    _, err := r.collection.InsertMany(ctx, docs)
    return err
}
```

---

## 九、测试规范

### 9.1 单元测试

```go
func TestMongoBookRepository_Create(t *testing.T) {
    // Setup
    ctx := context.Background()
    client := setupTestDB()
    defer client.Disconnect(ctx)

    repo := NewMongoBookRepository(client.Database("test"), nil, nil)

    // Test
    book := &Book{
        Title:   "测试书籍",
        Author:  "测试作者",
        Status:  "draft",
    }

    err := repo.Create(ctx, book)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, book.ID)
}
```

更详细请参考[repository层测试规范](../testing/01_测试层级规范/repository_层测试规范.md)

---

## 十、最佳实践

### 10.1 设计原则

✅ **推荐**：
- 接口与实现分离
- 使用依赖注入
- 完善的错误处理
- 合理的缓存策略
- 索引优化查询

❌ **避免**：
- 在Repository中实现业务逻辑
- 直接返回数据库错误
- 忽略索引管理
- 过度使用缓存

### 10.2 代码质量

- [ ] 接口定义清晰
- [ ] 实现与数据库解耦
- [ ] 错误处理完善
- [ ] 支持事务操作
- [ ] 有索引管理
- [ ] 有单元测试

---

**相关文档**：
- [架构设计规范](../architecture/架构设计规范.md)
- [测试规范](../testing/测试架构设计规范.md)
