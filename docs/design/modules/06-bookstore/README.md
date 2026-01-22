# 06 - 书城模块

> **模块编号**: 06
> **模块名称**: Bookstore
> **负责功能**: 书籍浏览、搜索、分类、榜单、章节购买
> **完成度**: 🟢 80%

## 📋 目录结构

```
书城模块/
├── api/v1/
│   └── bookstore/                # 书城API
│       ├── book_api.go          # 书籍管理
│       ├── category_api.go      # 分类管理
│       ├── search_api.go        # 搜索功能
│       ├── ranking_api.go       # 榜单系统
│       └── purchase_api.go      # 购买管理
├── service/bookstore/            # 书城服务层
│   ├── book_service.go         # 书籍服务
│   ├── category_service.go     # 分类服务
│   ├── search_service.go       # 搜索服务
│   ├── ranking_service.go      # 榜单服务
│   └── purchase_service.go     # 购买服务
├── repository/interfaces/bookstore/ # 仓储接口
├── repository/mongodb/bookstore/    # MongoDB仓储实现
└── models/bookstore/                # 数据模型
    ├── book.go                    # 书籍
    ├── chapter.go                 # 章节
    ├── category.go                # 分类
    └── purchase.go                # 购买记录
```

## 🎯 核心功能

### 1. 书籍展示

- **书籍列表**: 分页获取书籍列表
- **书籍详情**: 书籍详细信息
- **作者作品**: 作者的其他作品
- **相关推荐**: 相似书籍推荐

### 2. 分类浏览

- **分类导航**: 按类型浏览
- **标签筛选**: 按标签筛选
- **多条件筛选**: 综合筛选条件
- **分类排序**: 按热度、时间、评分排序

### 3. 搜索功能

- **全文搜索**: 书名、作者、简介搜索
- **搜索建议**: 搜索关键词提示
- **热门搜索**: 热门搜索词
- **搜索历史**: 用户搜索历史

### 4. 榜单系统

- **热销榜**: 按销量排行
- **推荐榜**: 编辑推荐
- **新书榜**: 最新发布
- **完结榜**: 已完结作品
- **收藏榜**: 收藏数量排行

### 5. 章节购买

- **付费章节**: 购买付费章节
- **会员特权**: 会员免费阅读
- **购买记录**: 购买历史查询
- **自动订阅**: 自动购买后续章节

## 📊 数据模型

### Book (书籍)

```go
type Book struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    AuthorID        primitive.ObjectID   `bson:"author_id" json:"authorId"`
    Title           string               `bson:"title" json:"title"`
    Description     string               `bson:"description" json:"description"`
    Cover           string               `bson:"cover" json:"cover"`
    Category        BookCategory         `bson:"category" json:"category"`
    Tags            []string             `bson:"tags" json:"tags"`

    // 统计信息
    ViewCount       int64                `bson:"view_count" json:"viewCount"`
    ReadCount       int64                `bson:"read_count" json:"readCount"`
    CollectCount    int64                `bson:"collect_count" json:"collectCount"`
    LikeCount       int64                `bson:"like_count" json:"likeCount"`
    CommentCount    int64                `bson:"comment_count" json:"commentCount"`

    // 状态
    Status          BookStatus           `bson:"status" json:"status"`
    IsCompleted     bool                 `bson:"is_completed" json:"isCompleted"`
    IsVip           bool                 `bson:"is_vip" json:"isVip"`

    // 章节信息
    TotalChapters   int                  `bson:"total_chapters" json:"totalChapters"`
    FreeChapters    int                  `bson:"free_chapters" json:"freeChapters"`
    WordCount       int64                `bson:"word_count" json:"wordCount"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
    PublishedAt     *time.Time           `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
    CompletedAt     *time.Time           `bson:"completed_at,omitempty" json:"completedAt,omitempty"`
}

type BookCategory string
const (
    CategoryFantasy   BookCategory = "fantasy"    // 玄幻
    CategoryUrban     BookCategory = "urban"      // 都市
    CategoryRomance   BookCategory = "romance"    // 言情
    CategoryHistory   BookCategory = "history"    // 历史
    CategorySciFi     BookCategory = "scifi"      // 科幻
    CategoryMilitary  BookCategory = "military"   // 军事
    CategoryGame      BookCategory = "game"       // 游戏
    CategorySports    BookCategory = "sports"     // 体育
)

type BookStatus string
const (
    BookStatusDraft     BookStatus = "draft"
    BookStatusOngoing   BookStatus = "ongoing"
    BookStatusCompleted BookStatus = "completed"
    BookStatusPaused    BookStatus = "paused"
)
```

### Chapter (章节)

```go
type Chapter struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    BookID          primitive.ObjectID   `bson:"book_id" json:"bookId"`
    Title           string               `bson:"title" json:"title"`
    Content         string               `bson:"content" json:"content"`
    ChapterNumber   int                  `bson:"chapter_number" json:"chapterNumber"`
    VolumeNumber    int                  `bson:"volume_number" json:"volumeNumber"`
    WordCount       int                  `bson:"word_count" json:"wordCount"`

    // 付费信息
    IsFree          bool                 `bson:"is_free" json:"isFree"`
    Price           int                  `bson:"price" json:"price"`           // 价格（分）

    // 统计信息
    ViewCount       int64                `bson:"view_count" json:"viewCount"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
    PublishedAt     *time.Time           `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
}
```

### Purchase (购买记录)

```go
type Purchase struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    BookID          primitive.ObjectID   `bson:"book_id" json:"bookId"`
    ChapterID       primitive.ObjectID   `bson:"chapter_id" json:"chapterId"`

    // 购买信息
    OrderID         string               `bson:"order_id" json:"orderId"`
    Amount          int                  `bson:"amount" json:"amount"`           // 金额（分）
    PaymentMethod   string               `bson:"payment_method" json:"paymentMethod"`

    // 时间戳
    PurchasedAt     time.Time            `bson:"purchased_at" json:"purchasedAt"`
}
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/bookstore/books | 获取书籍列表 | 否 |
| GET | /api/v1/bookstore/books/:id | 获取书籍详情 | 否 |
| GET | /api/v1/bookstore/books/:id/chapters | 获取章节列表 | 否 |
| GET | /api/v1/bookstore/chapters/:id | 获取章节内容 | 是 |
| GET | /api/v1/bookstore/categories | 获取分类列表 | 否 |
| GET | /api/v1/bookstore/categories/:id/books | 按分类获取书籍 | 否 |
| GET | /api/v1/bookstore/search | 搜索书籍 | 否 |
| GET | /api/v1/bookstore/search/suggestions | 搜索建议 | 否 |
| GET | /api/v1/bookstore/rankings | 获取榜单 | 否 |
| POST | /api/v1/bookstore/chapters/:id/purchase | 购买章节 | 是 |
| GET | /api/v1/bookstore/purchases | 获取购买记录 | 是 |

## 🔧 依赖关系

### 依赖的模块
- **01 - 认证授权**: 用户身份验证
- **02 - 写作创作**: 获取作品内容
- **08 - 财务**: 处理支付

### 被依赖的模块
- **03 - 阅读器**: 获取书籍内容阅读
- **04 - 社交互动**: 分享书籍

## 📈 扩展点

1. **推荐系统**
   - 个性化推荐
   - 协同过滤
   - 基于内容的推荐

2. **专题活动**
   - 专题页面
   - 活动推荐
   - 限时免费

3. **阅读引导**
   - 新书推荐
   - 阅读排行
   - 编辑推荐

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/bookstore/`
**相关设计**: [书城系统设计](../../reading/书城系统设计.md)
