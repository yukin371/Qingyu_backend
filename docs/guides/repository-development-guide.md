# MongoDB Repository 开发指南

**版本**: v1.0
**创建日期**: 2026-01-30
**最后更新**: 2026-01-30

---

## 📋 概述

本文档说明如何使用 `BaseMongoRepository` 基类开发 MongoDB Repository，提供统一的 ID 转换逻辑和通用 CRUD 方法喵~

### 目录

1. [快速开始](#快速开始)
2. [BaseMongoRepository 功能](#basemongorepository-功能)
3. [最佳实践](#最佳实践)
4. [示例代码](#示例代码)
5. [常见问题](#常见问题)

---

## 🚀 快速开始

### 创建新的 Repository

```go
package reader

import (
    "Qingyu_backend/models/reader"
    "Qingyu_backend/repository/mongodb/base"
    "context"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// MongoReadingProgressRepository 阅读进度仓储MongoDB实现
type MongoReadingProgressRepository struct {
    *base.BaseMongoRepository  // 嵌入基类，继承ID转换和通用CRUD方法喵~
    db *mongo.Database          // 保留db引用以备特殊需求
}

// NewMongoReadingProgressRepository 创建阅读进度仓储实例
func NewMongoReadingProgressRepository(db *mongo.Database) *MongoReadingProgressRepository {
    return &MongoReadingProgressRepository{
        BaseMongoRepository: base.NewBaseMongoRepository(db, "reading_progress"),
        db:                 db,
    }
}
```

### 使用基类的方法

```go
// 使用 ParseID 进行 ID 转换
func (r *MongoReadingProgressRepository) GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error) {
    objectID, err := r.ParseID(id)  // 使用基类的ParseID方法喵~
    if err != nil {
        return nil, err
    }

    var progress reader.ReadingProgress
    err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&progress)
    // ...
}

// 或者直接使用基类的通用方法
func (r *MongoReadingProgressRepository) GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error) {
    var progress reader.ReadingProgress
    err := r.FindByID(ctx, id, &progress)  // 使用基类的FindByID方法
    if err != nil {
        return nil, err
    }
    return &progress, nil
}
```

---

## 📦 BaseMongoRepository 功能

### 1. ID 转换辅助方法

这些方法封装了 `models/shared/types` 中的 ID 转换逻辑，提供统一的错误处理喵~

| 方法 | 描述 | 返回值 |
|------|------|--------|
| `ParseID(id string)` | 解析 ID 字符串为 ObjectID | `(primitive.ObjectID, error)` |
| `ParseIDs(ids []string)` | 批量解析 ID 字符串 | `([]primitive.ObjectID, error)` |
| `IDToHex(id primitive.ObjectID)` | ObjectID 转换为 hex 字符串 | `string` |
| `IDsToHex(ids []primitive.ObjectID)` | 批量转换 ObjectID 为 hex 字符串 | `[]string` |
| `IsValidID(id string)` | 检查是否为有效的 ObjectID 格式 | `bool` |
| `GenerateID()` | 生成新的 ObjectID 并返回 hex 字符串 | `string` |

**示例**：
```go
// 单个 ID 转换
objectID, err := r.ParseID("507f1f77bcf86cd799439011")
if err != nil {
    return nil, err  // 已经包含友好的错误消息
}

// 批量 ID 转换
oids, err := r.ParseIDs([]string{"id1", "id2", "id3"})
if err != nil {
    return nil, err
}

// ObjectID 转换为 hex 字符串
hexID := r.IDToHex(objectID)

// 验证 ID 格式
if !r.IsValidID(someID) {
    return nil, types.ErrInvalidIDFormat
}
```

### 2. 通用 CRUD 方法

这些方法提供常用的 CRUD 操作，减少重复代码喵~

| 方法 | 描述 |
|------|------|
| `FindByID(ctx, id, result)` | 根据 ID 查找文档 |
| `FindByIDWithOpts(ctx, id, result, opts)` | 根据 ID 查找文档（支持自定义选项）|
| `UpdateByID(ctx, id, update)` | 根据 ID 更新文档 |
| `DeleteByID(ctx, id)` | 根据 ID 删除文档 |
| `Find(ctx, filter, results, opts)` | 查找多个文档 |
| `FindOne(ctx, filter, result, opts)` | 查找单个文档 |
| `Count(ctx, filter)` | 统计文档数量 |
| `Exists(ctx, id)` | 检查指定 ID 的文档是否存在 |
| `ExistsByFilter(ctx, filter)` | 根据条件检查文档是否存在 |
| `Create(ctx, document)` | 创建文档 |
| `CreateWithResult(ctx, document)` | 创建文档并返回插入的 ID |

**示例**：
```go
// FindByID - 根据ID查找
var progress reader.ReadingProgress
err := r.FindByID(ctx, id, &progress)

// Find - 查找多个文档
var progresses []*reader.ReadingProgress
filter := bson.M{"user_id": userOID}
err := r.Find(ctx, filter, &progresses)

// Count - 统计文档数量
count, err := r.Count(ctx, filter)

// Exists - 检查文档是否存在
exists, err := r.Exists(ctx, id)
```

### 3. 其他辅助方法

| 方法 | 描述 |
|------|------|
| `GetCollection()` | 获取 MongoDB 集合实例（用于子类直接访问）|

---

## ✨ 最佳实践

### 1. ID 转换

**❌ 不推荐**：直接使用 `primitive.ObjectIDFromHex`
```go
objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
    return nil, fmt.Errorf("无效的ID: %w", err)
}
```

**✅ 推荐**：使用基类的 `ParseID` 方法
```go
objectID, err := r.ParseID(id)
if err != nil {
    return nil, err  // 已经包含友好的错误消息
}
```

### 2. 错误处理

基类的 `ParseID` 方法已经返回了友好的错误消息（`types.ErrInvalidIDFormat` 或 `types.ErrEmptyID`），直接返回即可喵~

```go
objectID, err := r.ParseID(id)
if err != nil {
    return nil, err  // ✅ 直接返回，不需要额外的错误包装
}
```

### 3. 批量操作

对于批量 ID 转换，使用 `ParseIDs` 方法喵~

```go
// ✅ 推荐：使用 ParseIDs
oids, err := r.ParseIDs(ids)
if err != nil {
    return nil, err
}
filter := bson.M{"_id": bson.M{"$in": oids}}
```

### 4. 复用通用方法

对于简单的 CRUD 操作，优先使用基类提供的通用方法喵~

```go
// ✅ 推荐：使用基类的 FindByID
var progress reader.ReadingProgress
err := r.FindByID(ctx, id, &progress)

// 对于需要特殊处理的查询，可以直接使用 collection
err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}, opts).Decode(&progress)
```

### 5. 访问 Collection

使用 `GetCollection()` 方法获取集合实例喵~

```go
// ✅ 推荐
cursor, err := r.GetCollection().Find(ctx, filter, opts)

// ❌ 不推荐：不要在子类中重新定义 collection 字段
```

---

## 📝 示例代码

### 完整的 Repository 示例

```go
package reader

import (
    "Qingyu_backend/models/reader"
    "Qingyu_backend/repository/mongodb/base"
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoReadingProgressRepository 阅读进度仓储MongoDB实现
type MongoReadingProgressRepository struct {
    *base.BaseMongoRepository  // 嵌入基类
    db *mongo.Database
}

// NewMongoReadingProgressRepository 创建实例
func NewMongoReadingProgressRepository(db *mongo.Database) *MongoReadingProgressRepository {
    return &MongoReadingProgressRepository{
        BaseMongoRepository: base.NewBaseMongoRepository(db, "reading_progress"),
        db:                 db,
    }
}

// Create 创建阅读进度
func (r *MongoReadingProgressRepository) Create(ctx context.Context, progress *reader.ReadingProgress) error {
    if progress.ID.IsZero() {
        progress.ID = primitive.NewObjectID()
    }
    progress.CreatedAt = time.Now()
    progress.UpdatedAt = time.Now()

    return r.GetCollection().InsertOne(ctx, progress)
}

// GetByID 根据ID获取阅读进度（使用基类方法）
func (r *MongoReadingProgressRepository) GetByID(ctx context.Context, id string) (*reader.ReadingProgress, error) {
    var progress reader.ReadingProgress
    err := r.FindByID(ctx, id, &progress)
    if err != nil {
        return nil, err
    }
    return &progress, nil
}

// GetByUserAndBook 获取用户对特定书籍的阅读进度
func (r *MongoReadingProgressRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.ReadingProgress, error) {
    // 使用 ParseID 进行 ID 转换
    userOID, err := r.ParseID(userID)
    if err != nil {
        return nil, err
    }

    bookOID, err := r.ParseID(bookID)
    if err != nil {
        return nil, err
    }

    var progress reader.ReadingProgress
    err = r.GetCollection().FindOne(ctx, bson.M{
        "user_id": userOID,
        "book_id": bookOID,
    }).Decode(&progress)

    if err != nil {
        return nil, nil // 没有阅读记录，返回nil而不是错误
    }

    return &progress, nil
}

// Update 更新阅读进度（使用基类方法）
func (r *MongoReadingProgressRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
    updates["updated_at"] = time.Now()
    update := bson.M{"$set": updates}
    return r.UpdateByID(ctx, id, update)
}

// Delete 删除阅读进度（使用基类方法）
func (r *MongoReadingProgressRepository) Delete(ctx context.Context, id string) error {
    return r.DeleteByID(ctx, id)
}

// CountByUser 统计用户的阅读进度数量
func (r *MongoReadingProgressRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
    userOID, err := r.ParseID(userID)
    if err != nil {
        return 0, err
    }

    return r.Count(ctx, bson.M{"user_id": userOID})
}
```

---

## ❓ 常见问题

### Q1: 为什么需要 BaseMongoRepository？

**A**: 为了解决以下问题喵~

1. **ID 转换重复代码**：每个 Repository 都重复进行 `primitive.ObjectIDFromHex()` 转换
2. **错误处理不统一**：不同的错误消息格式
3. **CRUD 方法重复**：很多 Repository 实现相同的 CRUD 操作
4. **维护困难**：修改需要在多处进行

### Q2: 为什么 BaseMongoRepository 在独立的 `base` 包中？

**A**: 为了避免循环依赖喵~

```
repository/mongodb/          ← 主包（包含 factory）
repository/mongodb/base/     ← BaseMongoRepository（基类）
repository/mongodb/reader/   ← Reader Repository（使用基类）
```

如果 BaseMongoRepository 在 `repository/mongodb` 包中，而 `repository/mongodb` 包的 factory 又导入了 `reader` 包，就会形成循环依赖喵~

### Q3: 如何处理特殊的查询逻辑？

**A**: 使用 `GetCollection()` 方法直接访问集合喵~

```go
// 对于复杂的聚合查询等，可以直接使用 collection
pipeline := mongo.Pipeline{
    {{Key: "$match", Value: bson.M{"user_id": userOID}}},
    {{Key: "$group", Value: bson.M{
        "_id":   nil,
        "total": bson.M{"$sum": "$reading_time"},
    }}},
}
cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
```

### Q4: 错误消息已经是友好的了吗？

**A**: 是的喵！基类的 `ParseID` 方法使用 `models/shared/types.ParseObjectID` 进行转换，返回的错误包含：
- `types.ErrEmptyID`: "ID cannot be empty"
- `types.ErrInvalidIDFormat`: "invalid ID format: must be 24-character hex"

直接返回这些错误即可喵~

### Q5: 如何在 Service 层使用？

**A**: Service 层直接调用 Repository 的方法，传递 string ID 即可喵~

```go
// Service 层
func (s *ReadingProgressService) GetProgress(ctx context.Context, progressID string) (*reader.ReadingProgress, error) {
    // 直接传递 string ID，Repository 内部会进行转换
    return s.repo.GetByID(ctx, progressID)
}
```

### Q6: 是否需要修改现有的 Repository？

**A**: 强烈建议迁移到使用 BaseMongoRepository，但可以逐步进行喵~

**迁移步骤**：
1. 嵌入 `*base.BaseMongoRepository`
2. 修改构造函数
3. 逐步替换 ID 转换代码
4. 逐步使用基类的通用方法
5. 运行测试验证

---

## 🔗 相关文档

- [Repository 层设计规范](../archived_standards/repository层设计规范_v1.0_过时.md)（已过时，待更新）
- [Repository 层测试规范](../standards/testing/01_测试层级规范/repository_层测试规范.md)
- [MongoDB ObjectID 转换](../../models/shared/types/id.go)

---

**文档维护者**: 猫娘助手 Kore 🐱
**最后审核**: 2026-01-30
**状态**: ✅ 有效
