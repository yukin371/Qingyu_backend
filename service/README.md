# Service 层架构文档

## 概览

```
┌───────────────────────────────────────────────────────┐
│                      API 层                           │
│                (api/v1/... )                          │
└───────────────────────────────────────────────────────┘
                         │
                         ▼
┌───────────────────────────────────────────────────────┐
│                   Service Container                   │
│           (service/container/... ) - 依赖注入          │
└───────────────────────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
    ┌─────────┐    ┌─────────┐    ┌─────────┐
    | Bookstore|    |  User   |    |  AI    |  ...
    | Service  |    | Service |    | Service |
    └─────────┘    └─────────┘    └─────────┘
                         │
                         ▼
    ┌───────────────────────────────────────────────────────┐
    |                   Repository 层                       │
    |            (repository/mongodb/... )                  │
    └───────────────────────────────────────────────────────┘
```

## 目录结构

```
service/
├── admin/          # 管理后台服务 (用户/内容/统计分析/审核)
├── ai/             # AI服务 (聊天/文心一言/图片生成/敏感词)
├── audit/          # 审计服务
├── auth/           # 认证服务 (JWT/Session/OAuth/权限)
├── base/           # 基础服务接口 (BaseService)
├── bookstore/      # 书城服务 (书籍/榜单/章节/评分/统计)
├── channels/       # 频道服务 (消息推送)
├── container/      # 服务容器 (依赖注入/服务注册)
├── content/        # 内容服务
├── events/         # 事件服务
├── finance/        # 财务服务 (钱包/会员/收入)
├── interfaces/     # 接口定义
├── reader/         # 读者服务 (阅读进度/书签/笔记/设置)
├── recommendation/ # 推荐服务
├── search/         # 搜索服务
├── shared/         # 共享服务 (缓存/存储/指标)
├── social/         # 社交服务 (评论/点赞/收藏/关注)
├── user/           # 用户服务 (用户管理/验证/状态)
├── validation/     # 验证服务
└── writer/         # 作家服务 (写作/项目/文档)
```

## 核心设计模式

### 1. 接口抽象
所有服务都定义接口, 便于测试和替换:

```go
// 定义接口
type BookstoreService interface {
    GetHomepageData(ctx context.Context) (*HomepageData, error)
    GetRealtimeRanking(ctx context.Context, limit int) ([]*RankingItem, error)
    // ...
}

// 实现接口
type BookstoreServiceImpl struct { ... }
```

### 2. 装饰器模式 (缓存)
用装饰器添加缓存, 不改变原实现:

```
原始服务 (BookstoreServiceImpl)
       │
       ▼ 包装
缓存装饰器 (CachedBookstoreService)
       │
       ▼
  1. 查 Redis
  2. 未命中 → 调用原始服务
  3. 异步写回 Redis
```

###  3. 依赖注入 (DI)
在 `service_container.go` 统一创建和注入:

```go
// 创建缓存服务
var bookstoreCacheService bookstoreService.CacheService
if c.redisClient != nil {
    bookstoreCacheService = bookstoreService.NewRedisCacheService(redisClient, "qingyu:bookstore")
}

// 创建基础服务
baseBookstoreService := bookstoreService.NewBookstoreService(...)

// 用装饰器包装 (启用缓存)
c.bookstoreService = bookstoreService.NewCachedBookstoreService(baseBookstoreService, bookstoreCacheService)
```

### 4. 降级策略
Redis 不可用时, 自动降级到 MongoDB:

```go
func (c *CachedBookstoreService) GetRealtimeRanking(ctx context.Context, limit int) ([]*RankingItem, error) {
    // 1. 尝试 Redis
    if items, err := c.cache.GetRanking(ctx, ...); err == nil {
        return items, nil
    }
    // 2. Redis 失败, 降级到 MongoDB
    return c.service.GetRealtimeRanking(ctx, limit)
}
```

## 缓存 Key 规范

| 服务 | Key 前缀 | 示例 |
|-----|---------|------|
| 书城 | `qingyu:bookstore:` | `qingyu:bookstore:homepage` |
| 读者 | `qingyu:` | `qingyu:reader:progress:{uid}` |
| 会话 | `session:` | `session:{sessionID}` |
| Token黑名单 | `token:blacklist:` | `token:blacklist:{tokenHash}` |

## 服务生命周期

```
┌─────────────────────────────────────────┐
│            应用启动                │
└─────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│   1. 创建 RepositoryFactory           │
│   2. 创建 ServiceContainer            │
│   3. 调用 SetupDefaultServices()        │
│      - 初始化数据库连接              │
│      - 创建缓存服务                │
│      - 创建业务服务                │
│      - 用装饰器包装 (启用缓存)        │
│   4. 调用 SetupBusinessServices()       │
│      - 创建 AI /搜索/推荐等服务        │
└─────────────────────────────────────────┘
```

## 新服务开发指南

### 1. 定义接口
```go
// interfaces/xxx/service_interface.go
type XxxService interface {
    DoSomething(ctx context.Context, id string) (*Model, error)
}
```

### 2. 实现服务
```go
// xxx/xxx_service.go
type XxxServiceImpl struct {
    repo Repository
    cache CacheService  // 可选
}

func NewXxxService(repo Repository, cache CacheService) XxxService {
    return &XxxServiceImpl{repo: repo, cache: cache}
}
```

### 3. 注册到容器
```go
// container/service_container.go

// 1. 添加字段
type ServiceContainer struct {
    // ...
    xxxService xxxService.XxxService
}

// 2. 在 SetupDefaultServices 或新方法中创建
func (c *ServiceContainer) SetupXxxService() error {
    repo := c.repositoryFactory.CreateXxxRepository()

    // 创建缓存服务 (可选)
    var cacheService xxxService.CacheService
    if c.redisClient != nil {
        // ...
    }

    // 创建并注册
    c.xxxService = xxxService.NewXxxService(repo, cacheService)
    return nil
}
```

### 4. 注入到 API
```go
// api/v1/xxx/xxx_api.go
type XxxAPI struct {
    service xxxService.XxxService
}

func NewXxxAPI(service xxxService.XxxService) *XxxAPI {
    return &XxxAPI{service: service}
}
```

## 常见问题

### Q: 如何判断是否需要缓存?
A: 高频读取 + 计算成本高 + 数据变化不频繁 = 需要缓存

### Q: 缓存过期时间如何设置?
A: 在 `cached_bookstore_service.go` 顶部定义:
```go
const (
    HomepageCacheExpiration = 5 * time.Minute
    RankingCacheExpiration = 10 * time.Minute
    BookCacheExpiration = 1 * time.Hour
)
```

### Q: 如何处理缓存失败?
A: 缓存读取失败时, 降级到直接查数据库, 不影响业务:
```go
if items, err := c.cache.GetXxx(...); err == nil {
    return items, nil
}
// 降级: 忽略错误, 继续查 DB
return c.service.GetXxx(...)
```
