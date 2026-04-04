# Service 层设计说明

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

Service 层是业务逻辑层，负责：

- **业务逻辑编排**：实现复杂的业务规则和流程
- **跨 Repository 协调**：协调多个 Repository 完成业务操作
- **DTO 转换**：处理数据传输对象的转换
- **缓存策略**：实现业务级缓存逻辑
- **事件发布**：发布领域事件
- **事务管理**：协调跨 Repository 的事务

### 1.2 依赖关系图

```
┌─────────────────────────────────────────────────────────┐
│                    上层依赖方                            │
│  ┌─────────────────────────────────────────────────┐   │
│  │                    API 层                        │   │
│  └───────────────────────┬─────────────────────────┘   │
└──────────────────────────┼──────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    Service 层                           │
│  ┌─────────────────────────────────────────────────┐   │
│  │                 container/                       │   │
│  │            (依赖注入容器)                        │   │
│  └───────────────────────┬─────────────────────────┘   │
│                          │                              │
│  ┌───────────┬───────────┼───────────┬───────────┐    │
│  │  admin/   │   ai/     │  auth/    │bookstore/ │    │
│  ├───────────┼───────────┼───────────┼───────────┤    │
│  │ reader/   │ social/   │ finance/  │  writer/  │    │
│  └─────┬─────┴─────┬─────┴─────┬─────┴─────┬─────┘    │
│        │           │           │           │           │
│  ┌─────┴───────────┴───────────┴───────────┴─────┐    │
│  │              base/ (基础服务)                  │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐     │    │
│  │  │Validator │  │EventBus  │  │Container │     │    │
│  │  └──────────┘  └──────────┘  └──────────┘     │    │
│  └───────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   Repository 层                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │BookRepository│ │UserRepository│ │  ...        │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
```

### 1.3 层级边界

| 可以做 | 不应该做 |
|--------|----------|
| 业务逻辑编排 | 处理 HTTP 请求/响应 |
| 跨 Repository 协调 | 直接操作数据库 |
| DTO 转换 | 定义 Models |
| 缓存读写 | 调用其他 Service（避免循环依赖） |
| 事件发布 | 包含 HTTP 状态码 |
| 事务管理 | 访问 Request/Response |

### 1.4 目录结构

```
service/
├── admin/          # 管理后台服务
├── ai/             # AI服务（聊天、图片生成）
│   └── adapter/    # AI 适配器（Claude、GPT、Gemini等）
├── audit/          # 审计服务
├── auth/           # 认证服务（JWT、Session、OAuth）
├── base/           # 基础服务（Validator、EventBus）
├── bookstore/      # 书城服务（书籍、榜单、章节）
├── channels/       # 频道服务（消息推送）
├── container/      # 服务容器（依赖注入）
├── finance/        # 财务服务（钱包、会员）
├── interfaces/     # 接口定义
├── reader/         # 读者服务（进度、书签、笔记）
├── social/         # 社交服务（评论、点赞、关注）
├── user/           # 用户服务
└── writer/         # 作家服务（写作、项目）
```

---

## 2. 命名与代码规范

### 2.1 文件命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 服务接口 | `interfaces.go` 或 `service_interface.go` | `interfaces.go` |
| 服务实现 | 实体名 + `_service.go` | `bookstore_service.go` |
| 缓存装饰器 | `cached_` + 实体名 + `_service.go` | `cached_bookstore_service.go` |
| 测试文件 | 原文件名 + `_test.go` | `bookstore_service_test.go` |

### 2.2 接口命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 服务接口 | 实体名 + `Service` | `BookstoreService`, `UserService` |
| 缓存接口 | `CacheService` | `BookstoreCacheService` |
| 基础接口 | 功能描述 | `BaseService`, `EventBus` |

### 2.3 实现命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 基础实现 | 接口名 + `Impl` | `BookstoreServiceImpl` |
| 缓存装饰器 | `Cached` + 接口名 | `CachedBookstoreService` |
| 工厂函数 | `New` + 实现名 | `NewBookstoreService()` |

### 2.4 方法命名

| 操作类型 | 命名规范 | 示例 |
|----------|----------|------|
| 获取单个 | `Get` + 实体名 + `By` + 条件 | `GetBookByID(ctx, id)` |
| 获取列表 | `Get` + 复数名 + `By` + 条件 | `GetBooksByCategory(ctx, categoryID)` |
| 创建 | `Create` + 实体名 | `CreateBook(ctx, book)` |
| 更新 | `Update` + 实体名 | `UpdateBook(ctx, id, updates)` |
| 删除 | `Delete` + 实体名 | `DeleteBook(ctx, id)` |
| 搜索 | `Search` + 实体名 | `SearchBooks(ctx, keyword)` |
| 统计 | `Get` + 实体名 + `Stats` | `GetBookStats(ctx)` |
| 增量操作 | `Increment` + 字段名 | `IncrementViewCount(ctx, id)` |

---

## 3. 设计模式与最佳实践

### 3.1 接口抽象模式

所有服务都定义接口，便于测试和替换：

```go
// 定义接口
type BookstoreService interface {
    GetBookByID(ctx context.Context, id string) (*bookstore.Book, error)
    GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore.Book, int64, error)
    SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.Book, int64, error)
}

// 实现接口
type BookstoreServiceImpl struct {
    bookRepo     BookstoreRepo.BookRepository
    categoryRepo BookstoreRepo.CategoryRepository
}

func NewBookstoreService(
    bookRepo BookstoreRepo.BookRepository,
    categoryRepo BookstoreRepo.CategoryRepository,
) BookstoreService {
    return &BookstoreServiceImpl{
        bookRepo:     bookRepo,
        categoryRepo: categoryRepo,
    }
}
```

### 3.2 缓存装饰器模式

用装饰器添加缓存，不改变原实现：

```go
// 缓存装饰器
type CachedBookstoreService struct {
    service BookstoreService
    cache   CacheService
}

func NewCachedBookstoreService(service BookstoreService, cache CacheService) BookstoreService {
    return &CachedBookstoreService{
        service: service,
        cache:   cache,
    }
}

func (c *CachedBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("book:%s", id)
    if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
        var book bookstore.Book
        if json.Unmarshal([]byte(cached), &book) == nil {
            return &book, nil
        }
    }

    // 2. 调用原始服务
    book, err := c.service.GetBookByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    if data, err := json.Marshal(book); err == nil {
        c.cache.Set(ctx, cacheKey, data, BookCacheExpiration)
    }

    return book, nil
}
```

### 3.3 依赖注入模式

在 `container/` 统一创建和注入：

```go
// container/service_container.go
type ServiceContainer struct {
    repositoryFactory interfaces.RepositoryFactory
    redisClient       *redis.Client

    // 服务实例
    bookstoreService bookstore.BookstoreService
    userService      user.UserService
    // ...
}

func (c *ServiceContainer) SetupDefaultServices() error {
    // 1. 创建缓存服务
    var cacheService bookstore.CacheService
    if c.redisClient != nil {
        cacheService = bookstore.NewRedisCacheService(c.redisClient)
    }

    // 2. 创建基础服务
    baseBookstoreService := bookstore.NewBookstoreService(
        c.repositoryFactory.CreateBookRepository(),
        c.repositoryFactory.CreateCategoryRepository(),
    )

    // 3. 用装饰器包装（启用缓存）
    if cacheService != nil {
        c.bookstoreService = bookstore.NewCachedBookstoreService(baseBookstoreService, cacheService)
    } else {
        c.bookstoreService = baseBookstoreService
    }

    return nil
}
```

### 3.4 降级策略模式

Redis 不可用时，自动降级：

```go
func (c *CachedBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
    // 1. 尝试 Redis
    if c.cache != nil {
        if cached, err := c.cache.Get(ctx, "book:"+id); err == nil {
            // 返回缓存数据
            return cached, nil
        }
    }

    // 2. Redis 失败，降级到原始服务
    return c.service.GetBookByID(ctx, id)
}
```

### 3.5 事件驱动模式

使用 EventBus 发布领域事件：

```go
// 定义事件
type BookPublishedEvent struct {
    base.BaseEvent
    BookID    string
    AuthorID  string
    PublishedAt time.Time
}

// 发布事件
func (s *BookstoreServiceImpl) PublishBook(ctx context.Context, bookID string) error {
    // 1. 更新书籍状态
    err := s.bookRepo.Update(ctx, bookID, map[string]interface{}{
        "status":       bookstore.BookStatusOngoing,
        "published_at": time.Now(),
    })
    if err != nil {
        return err
    }

    // 2. 发布事件
    event := &BookPublishedEvent{
        BaseEvent: base.BaseEvent{
            EventType: "book.published",
            Timestamp: time.Now(),
            Source:    "bookstore-service",
        },
        BookID: bookID,
    }

    return s.eventBus.PublishAsync(ctx, event)
}
```

### 3.6 反模式警示

| 反模式 | 问题 | 正确做法 |
|--------|------|----------|
| 在 Service 中处理 HTTP | 耦合表示层 | HTTP 处理放 API 层 |
| 直接操作数据库 | 绕过 Repository | 所有数据操作通过 Repository |
| 循环依赖 | 导致编译错误 | 使用接口解耦或事件机制 |
| 硬编码缓存时间 | 难以维护 | 使用常量定义 |
| 忽略错误处理 | 导致数据不一致 | 始终处理并包装错误 |

---

## 4. 接口与契约规范

### 4.1 基础接口定义

```go
// service/interfaces/base/base.go
package base

import "context"

// BaseService 基础服务接口
type BaseService interface {
    Initialize(ctx context.Context) error
    Health(ctx context.Context) error
    Close(ctx context.Context) error
}

// EventBus 事件总线接口
type EventBus interface {
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handlerName string) error
    Publish(ctx context.Context, event Event) error
    PublishAsync(ctx context.Context, event Event) error
}

// Event 事件接口
type Event interface {
    GetEventType() string
    GetEventData() interface{}
    GetTimestamp() time.Time
    GetSource() string
}

// EventHandler 事件处理器接口
type EventHandler interface {
    GetHandlerName() string
    Handle(ctx context.Context, event Event) error
}
```

### 4.2 返回值规范

| 场景 | 返回值 | 说明 |
|------|--------|------|
| 单条查询 | `(*Entity, error)` | 未找到返回 `nil, ErrNotFound` |
| 列表查询 | `([]*Entity, int64, error)` | 返回列表、总数、错误 |
| 创建操作 | `(*Entity, error)` 或 `error` | 返回创建的实体或仅错误 |
| 更新操作 | `error` | 包含 "not found" 错误 |
| 删除操作 | `error` | 包含 "not found" 错误 |
| 搜索操作 | `([]*Entity, int64, error)` | 返回结果、总数、错误 |

### 4.3 分页规范

```go
// 分页参数
type PaginationParams struct {
    Page     int `json:"page"`      // 页码，从1开始
    PageSize int `json:"pageSize"`  // 每页数量，默认20
}

// 分页结果
type PaginatedResult[T any] struct {
    Items      []*T  `json:"items"`
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PageSize   int   `json:"pageSize"`
    TotalPages int   `json:"totalPages"`
}

// 计算偏移量
func (p *PaginationParams) Offset() int {
    if p.Page < 1 {
        p.Page = 1
    }
    return (p.Page - 1) * p.PageSize
}
```

### 4.4 缓存 Key 规范

| 服务 | Key 前缀 | 示例 |
|------|---------|------|
| 书城 | `qingyu:bookstore:` | `qingyu:bookstore:homepage` |
| 读者 | `qingyu:reader:` | `qingyu:reader:progress:{uid}` |
| 会话 | `session:` | `session:{sessionID}` |
| Token 黑名单 | `token:blacklist:` | `token:blacklist:{tokenHash}` |

---

## 5. 测试策略

### 5.1 测试文件组织

```
service/
├── bookstore/
│   ├── bookstore_service.go
│   ├── bookstore_service_test.go      # 单元测试
│   └── cached_bookstore_service_test.go # 缓存测试
```

### 5.2 测试类型

| 测试类型 | 覆盖内容 | 工具 |
|----------|----------|------|
| 单元测试 | 业务逻辑 | Mock Repository |
| 集成测试 | 多服务协作 | testify |
| 缓存测试 | 缓存逻辑 | miniredis |

### 5.3 测试示例

```go
// bookstore_service_test.go
package bookstore

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockBookRepository 模拟书籍仓储
type MockBookRepository struct {
    mock.Mock
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*bookstore.Book), args.Error(1)
}

func TestBookstoreService_GetBookByID(t *testing.T) {
    // 1. 准备
    mockRepo := new(MockBookRepository)
    service := NewBookstoreService(mockRepo, nil, nil, nil, nil)

    expectedBook := &bookstore.Book{
        Title:  "测试书籍",
        Author: "测试作者",
    }
    mockRepo.On("GetByID", mock.Anything, "book123").Return(expectedBook, nil)

    // 2. 执行
    book, err := service.GetBookByID(context.Background(), "book123")

    // 3. 验证
    assert.NoError(t, err)
    assert.Equal(t, "测试书籍", book.Title)
    mockRepo.AssertExpectations(t)
}

func TestBookstoreService_GetBookByID_NotFound(t *testing.T) {
    mockRepo := new(MockBookRepository)
    service := NewBookstoreService(mockRepo, nil, nil, nil, nil)

    mockRepo.On("GetByID", mock.Anything, "notexist").Return(nil, nil)

    book, err := service.GetBookByID(context.Background(), "notexist")

    assert.Error(t, err)
    assert.Nil(t, book)
    assert.Contains(t, err.Error(), "not found")
}
```

### 5.4 测试覆盖率要求

| 类型 | 最低覆盖率 |
|------|-----------|
| 业务逻辑 | 80% |
| 缓存逻辑 | 90% |
| 错误处理 | 70% |

---

## 6. 完整代码示例

### 6.1 完整的服务接口

```go
// service/bookstore/interfaces.go
package bookstore

import (
    "context"

    "Qingyu_backend/models/bookstore"
)

// BookstoreService 书城服务接口
type BookstoreService interface {
    // 书籍列表
    GetAllBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, int64, error)
    GetBookByID(ctx context.Context, id string) (*bookstore.Book, error)
    GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore.Book, int64, error)
    GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, int64, error)
    SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.Book, int64, error)

    // 分类
    GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error)
    GetCategoryByID(ctx context.Context, id string) (*bookstore.Category, error)

    // 榜单
    GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore.RankingItem, error)
    UpdateRankings(ctx context.Context, rankingType bookstore.RankingType, period string) error

    // 首页聚合
    GetHomepageData(ctx context.Context) (*HomepageData, error)

    // 统计
    GetBookStats(ctx context.Context) (*bookstore.BookStats, error)
    IncrementBookView(ctx context.Context, bookID string) error
}

// CacheService 缓存服务接口
type CacheService interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
    GetRanking(ctx context.Context, key string) ([]*bookstore.RankingItem, error)
    SetRanking(ctx context.Context, key string, items []*bookstore.RankingItem, expiration time.Duration) error
}

// HomepageData 首页数据
type HomepageData struct {
    Banners          []*bookstore.Banner                 `json:"banners"`
    RecommendedBooks []*bookstore.Book                   `json:"recommendedBooks"`
    FeaturedBooks    []*bookstore.Book                   `json:"featuredBooks"`
    Categories       []*bookstore.Category               `json:"categories"`
    Stats            *bookstore.BookStats                `json:"stats"`
    Rankings         map[string][]*bookstore.RankingItem `json:"rankings"`
}
```

### 6.2 完整的服务实现

```go
// service/bookstore/bookstore_service.go
package bookstore

import (
    "context"
    "errors"
    "fmt"
    "time"

    "Qingyu_backend/models/bookstore"
    BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// BookstoreServiceImpl 书城服务实现
type BookstoreServiceImpl struct {
    bookRepo       BookstoreRepo.BookRepository
    categoryRepo   BookstoreRepo.CategoryRepository
    bannerRepo     BookstoreRepo.BannerRepository
    rankingRepo    BookstoreRepo.RankingRepository
}

// NewBookstoreService 创建书城服务实例
func NewBookstoreService(
    bookRepo BookstoreRepo.BookRepository,
    categoryRepo BookstoreRepo.CategoryRepository,
    bannerRepo BookstoreRepo.BannerRepository,
    rankingRepo BookstoreRepo.RankingRepository,
) BookstoreService {
    return &BookstoreServiceImpl{
        bookRepo:     bookRepo,
        categoryRepo: categoryRepo,
        bannerRepo:   bannerRepo,
        rankingRepo:  rankingRepo,
    }
}

// GetBookByID 根据ID获取书籍详情
func (s *BookstoreServiceImpl) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
    if id == "" {
        return nil, errors.New("book ID cannot be empty")
    }

    book, err := s.bookRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get book: %w", err)
    }

    if book == nil {
        return nil, errors.New("book not found")
    }

    // 只有公开状态的书籍可以访问
    if !book.Status.IsPublic() {
        return nil, errors.New("book is not available")
    }

    return book, nil
}

// GetBooksByCategory 根据分类获取书籍列表
func (s *BookstoreServiceImpl) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore.Book, int64, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    offset := (page - 1) * pageSize

    // 获取书籍列表
    books, err := s.bookRepo.GetByCategory(ctx, categoryID, pageSize, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get books by category: %w", err)
    }

    // 获取总数
    total, err := s.bookRepo.CountByCategory(ctx, categoryID)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count books: %w", err)
    }

    return books, total, nil
}

// GetHomepageData 获取首页聚合数据
func (s *BookstoreServiceImpl) GetHomepageData(ctx context.Context) (*HomepageData, error) {
    data := &HomepageData{
        Rankings: make(map[string][]*bookstore.RankingItem),
    }

    // 并行获取各类数据
    var wg sync.WaitGroup
    var mu sync.Mutex
    errs := make([]error, 0)

    // 获取 Banner
    wg.Add(1)
    go func() {
        defer wg.Done()
        banners, err := s.bannerRepo.GetActive(ctx, 5)
        if err != nil {
            mu.Lock()
            errs = append(errs, fmt.Errorf("failed to get banners: %w", err))
            mu.Unlock()
            return
        }
        mu.Lock()
        data.Banners = banners
        mu.Unlock()
    }()

    // 获取推荐书籍
    wg.Add(1)
    go func() {
        defer wg.Done()
        books, _, err := s.bookRepo.GetRecommended(ctx, 10, 0)
        if err != nil {
            mu.Lock()
            errs = append(errs, fmt.Errorf("failed to get recommended books: %w", err))
            mu.Unlock()
            return
        }
        mu.Lock()
        data.RecommendedBooks = books
        mu.Unlock()
    }()

    // 获取实时榜单
    wg.Add(1)
    go func() {
        defer wg.Done()
        ranking, err := s.rankingRepo.GetByType(ctx, bookstore.RankingTypeRealtime, "daily", 10)
        if err != nil {
            mu.Lock()
            errs = append(errs, fmt.Errorf("failed to get ranking: %w", err))
            mu.Unlock()
            return
        }
        mu.Lock()
        data.Rankings["realtime"] = ranking
        mu.Unlock()
    }()

    wg.Wait()

    if len(errs) > 0 {
        return nil, errs[0]
    }

    return data, nil
}

// IncrementBookView 增加书籍浏览量
func (s *BookstoreServiceImpl) IncrementBookView(ctx context.Context, bookID string) error {
    return s.bookRepo.IncrementViewCount(ctx, bookID)
}
```

### 6.3 缓存装饰器实现

```go
// service/bookstore/cached_bookstore_service.go
package bookstore

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "Qingyu_backend/models/bookstore"
)

const (
    HomepageCacheExpiration = 5 * time.Minute
    BookCacheExpiration     = 1 * time.Hour
    RankingCacheExpiration  = 10 * time.Minute
)

// CachedBookstoreService 缓存装饰器
type CachedBookstoreService struct {
    service BookstoreService
    cache   CacheService
}

// NewCachedBookstoreService 创建缓存装饰器
func NewCachedBookstoreService(service BookstoreService, cache CacheService) BookstoreService {
    return &CachedBookstoreService{
        service: service,
        cache:   cache,
    }
}

// GetBookByID 获取书籍（带缓存）
func (c *CachedBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
    cacheKey := fmt.Sprintf("book:%s", id)

    // 1. 尝试从缓存获取
    if c.cache != nil {
        if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
            var book bookstore.Book
            if json.Unmarshal([]byte(cached), &book) == nil {
                return &book, nil
            }
        }
    }

    // 2. 调用原始服务
    book, err := c.service.GetBookByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    if c.cache != nil {
        if data, err := json.Marshal(book); err == nil {
            c.cache.Set(ctx, cacheKey, string(data), BookCacheExpiration)
        }
    }

    return book, nil
}

// GetHomepageData 获取首页数据（带缓存）
func (c *CachedBookstoreService) GetHomepageData(ctx context.Context) (*HomepageData, error) {
    cacheKey := "homepage:data"

    // 1. 尝试从缓存获取
    if c.cache != nil {
        if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
            var data HomepageData
            if json.Unmarshal([]byte(cached), &data) == nil {
                return &data, nil
            }
        }
    }

    // 2. 调用原始服务
    data, err := c.service.GetHomepageData(ctx)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    if c.cache != nil {
        if bytes, err := json.Marshal(data); err == nil {
            c.cache.Set(ctx, cacheKey, string(bytes), HomepageCacheExpiration)
        }
    }

    return data, nil
}

// GetRealtimeRanking 获取实时榜单（带缓存）
func (c *CachedBookstoreService) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore.RankingItem, error) {
    cacheKey := fmt.Sprintf("ranking:realtime:%d", limit)

    // 1. 尝试从缓存获取
    if c.cache != nil {
        if items, err := c.cache.GetRanking(ctx, cacheKey); err == nil && len(items) > 0 {
            return items, nil
        }
    }

    // 2. 调用原始服务
    items, err := c.service.GetRealtimeRanking(ctx, limit)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    if c.cache != nil {
        c.cache.SetRanking(ctx, cacheKey, items, RankingCacheExpiration)
    }

    return items, nil
}
```

---

## 附录

### A. 相关文档

- [Models 层设计说明](./layer-models.md)
- [Repository 层设计说明](./layer-repository.md)
- [API 层设计说明](./layer-api.md)

### B. 参考资源

- [Go 依赖注入](https://github.com/google/wire)
- [testify 测试框架](https://github.com/stretchr/testify)
- [miniredis Redis 模拟](https://github.com/alicebob/miniredis)

---

*最后更新：2026-03-19*
