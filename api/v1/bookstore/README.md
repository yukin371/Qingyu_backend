# Bookstore API 模块 - 书店

## 📁 模块职责

**Bookstore（书店）**模块负责所有与书籍浏览、发现和购买相关的功能，类似于一个线上书城。

## 🎯 核心功能

### 1. 书城首页
- Banner轮播图
- 推荐书籍
- 精选书籍
- 热门分类

### 2. 书籍浏览
- 书籍列表展示
- 分类浏览
- 标签筛选
- 搜索功能
- 排行榜

### 3. 书籍详情
- 基本信息
- 内容简介
- 作者介绍
- 章节目录
- 相关推荐

### 4. 书籍评分
- 查看评分
- 用户评分
- 评分统计

### 5. 书籍统计
- 阅读量
- 收藏数
- 评论数
- 分享数

### 6. 章节预览
- 前几章免费预览
- 章节基本信息

## 📦 文件结构

```
api/v1/bookstore/
├── bookstore_api.go          # 书城主要功能（首页、列表、搜索）
├── book_detail_api.go        # 书籍详情
├── book_rating_api.go        # 书籍评分
├── book_statistics_api.go    # 书籍统计
├── chapter_api.go            # 章节预览
├── types.go                  # 共享类型定义
└── README.md                 # 本文档
```

## 🌐 API路由

### 公开接口（无需认证）

```
GET  /api/v1/bookstore/homepage              # 书城首页
GET  /api/v1/bookstore/books                 # 书籍列表
GET  /api/v1/bookstore/books/search          # 搜索书籍
GET  /api/v1/bookstore/books/:id             # 书籍详情
GET  /api/v1/bookstore/books/:id/chapters    # 书籍章节目录
GET  /api/v1/bookstore/books/:id/related     # 相关推荐
GET  /api/v1/bookstore/chapters/:id          # 章节预览
GET  /api/v1/bookstore/categories            # 分类列表
GET  /api/v1/bookstore/categories/:id/books  # 分类下的书籍
GET  /api/v1/bookstore/tags                  # 标签列表
GET  /api/v1/bookstore/tags/:id/books        # 标签下的书籍
GET  /api/v1/bookstore/rankings              # 排行榜
GET  /api/v1/bookstore/rankings/:type        # 指定类型排行榜
GET  /api/v1/bookstore/books/:id/statistics  # 书籍统计
```

### 需要认证的接口

```
GET    /api/v1/bookstore/books/:id/rating    # 获取书籍评分
POST   /api/v1/bookstore/books/:id/rating    # 评分
PUT    /api/v1/bookstore/books/:id/rating    # 更新评分
DELETE /api/v1/bookstore/books/:id/rating    # 删除评分
GET    /api/v1/bookstore/my/ratings          # 我的评分记录
GET    /api/v1/bookstore/books/:id/favorite  # 收藏状态
```

## 🔄 与Reader模块的区别

| 功能 | Bookstore（书店） | Reader（阅读器） |
|------|------------------|-----------------|
| **定位** | 发现和浏览 | 阅读和学习 |
| **用户场景** | 找书、选书 | 读书、记笔记 |
| **核心功能** | 搜索、推荐、详情 | 阅读、进度、标注 |
| **章节** | 预览（前几章） | 完整内容 |
| **认证要求** | 多为公开 | 必须认证 |
| **数据存储** | 书籍元数据 | 用户阅读数据 |

## 🎨 使用场景

### 场景1：新用户找书
```
1. 访问书城首页 → GET /bookstore/homepage
2. 浏览推荐书籍
3. 点击感兴趣的书 → GET /bookstore/books/:id
4. 查看章节目录和前几章预览
5. 决定加入书架（跳转到Reader模块）
```

### 场景2：搜索特定书籍
```
1. 搜索关键词 → GET /bookstore/books/search?keyword=xxx
2. 筛选分类和标签
3. 查看书籍详情
4. 查看其他读者的评分
```

### 场景3：浏览排行榜
```
1. 访问排行榜 → GET /bookstore/rankings
2. 选择排行榜类型（热度、收藏、评分）
3. 浏览榜单书籍
4. 点击查看详情
```

## 🔧 技术特点

### 1. 缓存优化
- 首页数据缓存
- 热门书籍缓存
- 分类和标签缓存

### 2. 性能优化
- 分页加载
- 图片懒加载
- CDN加速

### 3. SEO友好
- 书籍详情页面静态化
- 元数据优化
- 结构化数据

### 4. 数据分析
- 访问统计
- 转化率追踪
- 用户行为分析

## 📊 数据模型

### Book（书籍）
```go
type Book struct {
    ID          string
    Title       string
    Author      string
    Cover       string
    Description string
    Category    string
    Tags        []string
    Status      string
    Statistics  BookStatistics
}
```

### BookStatistics（书籍统计）
```go
type BookStatistics struct {
    ViewCount     int64
    FavoriteCount int64
    CommentCount  int64
    ShareCount    int64
    AverageRating float64
}
```

### BookRating（书籍评分）
```go
type BookRating struct {
    BookID    string
    UserID    string
    Rating    float64
    Comment   string
    CreatedAt time.Time
}
```

## 🚀 后续规划

### Phase 1（已完成）
- ✅ 书城首页
- ✅ 书籍列表和搜索
- ✅ 书籍详情
- ✅ 分类和标签

### Phase 2（进行中）
- 🔄 评分和评论
- 🔄 排行榜
- 🔄 个性化推荐

### Phase 3（计划中）
- 📋 书单功能
- 📋 作者主页
- 📋 社区互动
- 📋 付费购买

## 📚 相关文档

- [Reader API 模块](../reader/README.md)
- [Bookstore Service 设计](../../../doc/design/bookstore/README.md)
- [数据库设计](../../../doc/database/bookstore_schema.md)

---

**版本**: v2.0  
**更新日期**: 2025-10-22  
**维护者**: Bookstore模块开发组

