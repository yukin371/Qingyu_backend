# 书单系统模块实施文档

## 概述

本文档记录了青羽平台书单系统 (BookList) 模块的完整实施过程。

**实施日期**: 2026-01-05
**状态**: ✅ 已完成

## 功能范围

书单系统模块提供以下核心功能：

1. **书单管理** - 创建、编辑、删除书单
2. **书单浏览** - 书单广场、收藏的书单
3. **书单详情** - 查看书单详情和书籍列表
4. **我的书单** - 查看和编辑自己的书单
5. **书单收藏** - 收藏/取消收藏书单

## 架构设计

### 分层架构

```
┌─────────────────────────────────────────┐
│         API Layer (booklist)            │
│  - BookListAPI                          │
│  - 路由: /api/v1/booklists/             │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│       Service Layer (social)            │
│  - BookListService                      │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│    Repository Layer (interfaces/social) │
│  - BookListRepository                   │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│    MongoDB Repository (mongodb/social)  │
│  - MongoBookListRepository              │
└─────────────────────────────────────────┘
```

## 实施步骤

### 步骤 1: 实现 MongoDB 仓储 (feat: social)

**提交**: `28018e7 feat(social): 实现 BookListRepository MongoDB 仓储`

#### 创建文件
- `repository/mongodb/social/booklist_repository_mongo.go`

#### 实现内容

实现了完整的 MongoDB 仓储层，包括：

1. **基础 CRUD 操作**
   - `Create` - 创建书单
   - `FindByID` - 根据 ID 查找
   - `Update` - 更新书单
   - `Delete` - 删除书单
   - `List` - 列出书单

2. **查询操作**
   - `FindByCreatorID` - 查找用户创建的书单
   - `FindPublic` - 查找公开书单
   - `FindByCategory` - 按分类查找
   - `SearchByName` - 搜索书单名称
   - `FindTrending` - 查找热门书单

3. **统计操作**
   - `CountByCreator` - 统计用户创建的书单数
   - `IncrementViewCount` - 增加浏览次数
   - `IncrementFavoriteCount` - 增加收藏次数
   - `DecrementFavoriteCount` - 减少收藏次数

4. **收藏管理**
   - `AddToFavorites` - 添加到收藏
   - `RemoveFromFavorites` - 从收藏移除
   - `FindFavoritesByUserID` - 查找用户收藏

#### 核心代码示例

```go
// Create 创建书单
func (r *MongoBookListRepository) Create(ctx context.Context, booklist *interfaces.BookList) error {
    booklist.ID = primitive.NewObjectID().Hex()
    booklist.CreatedAt = time.Now()
    booklist.UpdatedAt = time.Now()

    collection := r.database.Collection("booklists")
    _, err := collection.InsertOne(ctx, booklist)
    return err
}

// FindTrending 查找热门书单
func (r *MongoBookListRepository) FindTrending(ctx context.Context, limit int) ([]*interfaces.BookList, error) {
    collection := r.database.Collection("booklists")

    pipeline := mongo.Pipeline{
        {{"$match", bson.D{
            {"is_public", true},
            {"status", "published"},
        }}},
        {{"$sort", bson.D{
            {"favorite_count", -1},
            {"view_count", -1},
        }}},
        {{"$limit", limit}},
    }

    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }

    var booklists []*interfaces.BookList
    if err = cursor.All(ctx, &booklists); err != nil {
        return nil, err
    }

    return booklists, nil
}
```

### 步骤 2: 添加仓储工厂方法 (refactor: repository)

**提交**: `4c72d79 refactor(repository): 在 RepositoryFactory 中添加 BookListRepository`

#### 修改文件
- `repository/interfaces/RepoFactory_interface.go`
- `repository/mongodb/factory.go`

#### 实现内容

1. 在 RepoFactory 接口中添加方法：

```go
CreateBookListRepository() socialInterfaces.BookListRepository
```

2. 在 MongoRepositoryFactory 中实现：

```go
func (f *MongoRepositoryFactory) CreateBookListRepository() socialRepo.BookListRepository {
    return social.NewMongoBookListRepository(f.database)
}
```

### 步骤 3: 启用服务注册 (feat: service)

**提交**: `bb97a60 feat(service): 添加 ReadingStatsService 并启用 booklist 路由`

#### 修改文件
- `service/container/service_container.go`
- `router/enter.go`

#### 实现内容

1. 在 ServiceContainer 中初始化 BookListService：

```go
// ============ 4.6 创建书单服务 ============
bookListRepo := c.repositoryFactory.CreateBookListRepository()
c.bookListService = socialService.NewBookListService(bookListRepo)
c.services["BookListService"] = c.bookListService
```

2. 在路由中注册 booklist 路由：

```go
// ============ 注册书单路由 ============
booklistSvc, booklistErr := serviceContainer.GetBookListService()
if booklistErr != nil {
    logger.Warn("获取书单服务失败", zap.Error(booklistErr))
    logger.Info("书单路由未注册")
} else {
    booklistRouter.RegisterBookListRoutes(v1, booklistSvc)
    logger.Info("✓ 书单路由已注册到: /api/v1/booklists/")
}
```

## API 端点

### 公开端点

#### 获取书单广场
```
GET /api/v1/booklists
Query: page, pageSize, sort, category
```

#### 获取书单详情
```
GET /api/v1/booklists/:id
```

#### 搜索书单
```
GET /api/v1/booklists/search
Query: q, page, pageSize
```

#### 获取热门书单
```
GET /api/v1/booklists/trending
Query: limit
```

### 需要认证的端点

#### 创建书单
```
POST /api/v1/booklists
Authorization: Bearer {token}
Body: {
  "name": "我的书单",
  "description": "书单描述",
  "cover_image": "图片URL",
  "category": "玄幻",
  "is_public": true
}
```

#### 更新书单
```
PUT /api/v1/booklists/:id
Authorization: Bearer {token}
```

#### 删除书单
```
DELETE /api/v1/booklists/:id
Authorization: Bearer {token}
```

#### 获取我的书单
```
GET /api/v1/booklists/my
Authorization: Bearer {token}
Query: page, pageSize, status
```

#### 添加书籍到书单
```
POST /api/v1/booklists/:id/books
Authorization: Bearer {token}
Body: {
  "book_id": "书籍ID",
  "order": 1,
  "note": "推荐语"
}
```

#### 移除书籍
```
DELETE /api/v1/booklists/:id/books/:bookId
Authorization: Bearer {token}
```

#### 收藏书单
```
POST /api/v1/booklists/:id/favorite
Authorization: Bearer {token}
```

#### 取消收藏
```
DELETE /api/v1/booklists/:id/favorite
Authorization: Bearer {token}
```

#### 获取收藏的书单
```
GET /api/v1/booklists/favorites
Authorization: Bearer {token}
Query: page, pageSize
```

## 数据模型

### BookList (书单)
```go
type BookList struct {
    ID               string
    CreatorID        string
    Name             string
    Description      string
    CoverImage       string
    Category         string
    Tags            []string
    Books           []BookListItem
    IsPublic         bool
    Status           string  // draft, published, archived
    ViewCount        int64
    FavoriteCount    int64
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

### BookListItem (书单书籍项)
```go
type BookListItem struct {
    BookID       string
    Order        int
    Note         string
    AddedAt      time.Time
}
```

## 服务方法

### BookListService

#### CreateBookList
创建新书单
- 验证用户权限
- 生成书单 ID
- 设置时间戳

#### UpdateBookList
更新书单信息
- 验证所有权
- 更新字段
- 更新时间戳

#### DeleteBookList
删除书单
- 验证所有权
- 删除关联收藏
- 删除书单

#### AddBook
添加书籍到书单
- 验证书单所有权
- 检查书籍是否存在
- 添加到书单

#### RemoveBook
从书单移除书籍
- 验证书单所有权
- 移除书籍

#### ToggleFavorite
切换收藏状态
- 检查是否已收藏
- 添加或移除收藏
- 更新收藏计数

## 技术要点

### 1. 聚合管道查询
使用 MongoDB 聚合管道实现复杂查询：

```go
pipeline := mongo.Pipeline{
    {{"$match", bson.D{...}}},
    {{"$sort", bson.D{...}}},
    {{"$limit", limit}},
}
cursor, err := collection.Aggregate(ctx, pipeline)
```

### 2. 原子操作
使用原子操作更新计数：

```go
update := bson.D{
    {"$inc", bson.D{
        {"favorite_count", 1},
    }},
    {"$set", bson.D{
        {"updated_at", time.Now()},
    }},
}
```

### 3. 索引优化
建议创建以下索引：

```javascript
db.booklists.createIndex({ creator_id: 1, created_at: -1 })
db.booklists.createIndex({ is_public: 1, status: 1, favorite_count: -1 })
db.booklists.createIndex({ category: 1, view_count: -1 })
db.booklists.createIndex({ name: "text", description: "text" })
```

## 性能考虑

1. **分页查询**: 所有列表接口支持分页
2. **计数优化**: 使用原子操作更新计数
3. **缓存策略**: 热门书单可以缓存
4. **全文搜索**: 使用 MongoDB 全文搜索

## 后续优化

1. **推荐算法**: 基于用户偏好推荐书单
2. **协作编辑**: 支持多人协作编辑书单
3. **书单模板**: 预设书单模板
4. **导入导出**: 支持书单导入导出
5. **统计分析**: 书单数据分析

## 相关文档

- [项目结构总结](../docs/项目结构总结.md)
- [阅读统计模块](READING_STATS_IMPLEMENTATION.md)
- [P0 中间件集成](MIDDLEWARE_INTEGRATION.md)

## 提交历史

```
bb97a60 - feat(service): 添加 ReadingStatsService 并启用 booklist 路由
4c72d79 - refactor(repository): 在 RepositoryFactory 中添加 BookListRepository
28018e7 - feat(social): 实现 BookListRepository MongoDB 仓储
```
