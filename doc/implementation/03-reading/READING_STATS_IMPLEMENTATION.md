# 阅读统计模块实施文档

## 概述

本文档记录了青羽平台阅读统计 (Reading Stats) 模块的完整实施过程。

**实施日期**: 2026-01-06 ~ 2026-01-07
**状态**: ✅ 已完成

## 功能范围

阅读统计模块提供以下核心功能：

1. **章节统计** - 章节阅读量、完读率、跳出率
2. **作品统计** - 作品总阅读量、收入趋势
3. **读者行为** - 阅读时长、阅读历史记录
4. **热力图生成** - 阅读热力图数据
5. **趋势分析** - 阅读量和收入趋势
6. **用户排行** - 阅读时长排行、活跃度排行

## 架构设计

### 分层架构

```
┌─────────────────────────────────────────┐
│         API Layer (reading-stats)       │
│  - ReadingStatsAPI                      │
│  - 路由: /api/v1/reading-stats/         │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│       Service Layer (reader/stats)      │
│  - ReadingStatsService                  │
│  - 实现 BaseService 接口                │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│    Repository Layer (interfaces/stats)  │
│  - ChapterStatsRepository               │
│  - ReaderBehaviorRepository             │
│  - BookStatsRepository                  │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│    MongoDB Repository (mongodb/stats)   │
│  - MongoChapterStatsRepository          │
│  - MongoReaderBehaviorRepository        │
│  - MongoBookStatsRepository             │
└─────────────────────────────────────────┘
```

## 实施步骤

### 步骤 1: 添加仓储接口 (feat: repository)

**提交**: `4129a26 feat(repository): 在 RepositoryFactory 中添加 stats 模块仓储`

#### 修改文件
- `repository/interfaces/RepoFactory_interface.go`
- `repository/mongodb/factory.go`

#### 实现内容

1. 在 RepoFactory 接口中添加三个 stats 仓储创建方法：

```go
// Stats相关Repository
CreateChapterStatsRepository() StatsInterfaces.ChapterStatsRepository
CreateReaderBehaviorRepository() StatsInterfaces.ReaderBehaviorRepository
CreateBookStatsRepository() StatsInterfaces.BookStatsRepository
```

2. 在 MongoRepositoryFactory 中实现这些方法：

```go
// CreateChapterStatsRepository 创建章节统计Repository
func (f *MongoRepositoryFactory) CreateChapterStatsRepository() statsRepo.ChapterStatsRepository {
    return mongoStats.NewMongoChapterStatsRepository(f.database)
}

// CreateReaderBehaviorRepository 创建读者行为Repository
func (f *MongoRepositoryFactory) CreateReaderBehaviorRepository() statsRepo.ReaderBehaviorRepository {
    return mongoStats.NewMongoReaderBehaviorRepository(f.database)
}

// CreateBookStatsRepository 创建作品统计Repository
func (f *MongoRepositoryFactory) CreateBookStatsRepository() statsRepo.BookStatsRepository {
    return mongoStats.NewMongoBookStatsRepository(f.database)
}
```

### 步骤 2: 实现 BaseService 接口 (feat: service)

**提交**: `a5a8762 feat(service): ReadingStatsService 实现 BaseService 接口`

#### 修改文件
- `service/reader/stats/reading_stats_service.go`

#### 实现内容

添加缺失的 BaseService 接口方法：

```go
// Initialize 初始化服务
func (s *ReadingStatsService) Initialize(ctx context.Context) error {
    // ReadingStatsService 无需特殊初始化
    return nil
}

// Close 关闭服务
func (s *ReadingStatsService) Close(ctx context.Context) error {
    // ReadingStatsService 无需清理资源
    return nil
}
```

### 步骤 3: 启用服务注册 (feat: service container)

**提交**: `3c738ae feat(service): 启用 ReadingStatsService 并注册 reading-stats 路由`

#### 修改文件
- `service/container/service_container.go`
- `router/enter.go`

#### 实现内容

1. 在 ServiceContainer 中启用 ReadingStatsService 初始化：

```go
// ============ 4.8 创建阅读统计服务 ============
chapterStatsRepo := c.repositoryFactory.CreateChapterStatsRepository()
readerBehaviorRepo := c.repositoryFactory.CreateReaderBehaviorRepository()
bookStatsRepo := c.repositoryFactory.CreateBookStatsRepository()
c.readingStatsService = readingStatsService.NewReadingStatsService(
    chapterStatsRepo,
    readerBehaviorRepo,
    bookStatsRepo,
)
c.services["ReadingStatsService"] = c.readingStatsService
```

2. 在路由中启用 reading-stats 注册：

```go
// ============ 注册阅读统计路由 ============
readingStatsSvc, readingStatsErr := serviceContainer.GetReadingStatsService()
if readingStatsErr != nil {
    logger.Warn("获取阅读统计服务失败", zap.Error(readingStatsErr))
    logger.Info("阅读统计路由未注册")
} else {
    readingstatsRouter.RegisterReadingStatsRoutes(v1, readingStatsSvc)
    logger.Info("✓ 阅读统计路由已注册到: /api/v1/reading-stats/")
}
```

### 步骤 4: 修复路由冲突 (fix: router)

**提交**: `dd5f79a fix(router): 修复重复路由注册冲突`

#### 修改文件
- `router/enter.go`
- `router/usermanagement/usermanagement_router.go`

#### 问题修复

1. **问题 1**: 重复的 /metrics 端点
   - 移除 `router/enter.go` 中的重复注册
   - 保留 `core/server.go` 中的注册

2. **问题 2**: 重复的 /api/v1/admin/users 路由
   - 移除 `usermanagement_router.go` 中的管理员路由
   - 这些路由已在 `admin` 路由器中统一注册

## API 端点

### 我的统计 (需要认证)

#### 获取我的统计
```
GET /api/v1/reading-stats/my/stats
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "start_date": "2025-12-08",
    "end_date": "2026-01-07",
    "total_views": 0,
    "total_reading_time": 0,
    "books_read": 0,
    "chapters_read": 0
  }
}
```

#### 获取每日统计
```
GET /api/v1/reading-stats/my/daily
Authorization: Bearer {token}
```

#### 获取排名
```
GET /api/v1/reading-stats/my/ranking
Authorization: Bearer {token}
```

#### 获取阅读时长
```
GET /api/v1/reading-stats/my/reading-time
Authorization: Bearer {token}
```

#### 获取阅读历史
```
GET /api/v1/reading-stats/my/history
Authorization: Bearer {token}
```

### 推荐系统

#### 获取推荐
```
GET /api/v1/reading-stats/recommendations
Authorization: Bearer {token}
```

## 测试验证

### 测试环境
- 后端服务运行在 `http://localhost:8080`
- 测试用户: `testuser`

### 测试结果

#### 1. 健康检查
```bash
curl http://localhost:8080/health
```

**结果**: ✅ ReadingStatsService 显示为 `true`

#### 2. 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/user-management/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

**结果**: ✅ 注册成功，获得 JWT token

#### 3. 获取统计信息
```bash
curl http://localhost:8080/api/v1/reading-stats/my/stats \
  -H "Authorization: Bearer {token}"
```

**结果**: ✅ 返回 200 状态码和统计数据

#### 4. 获取每日统计
```bash
curl http://localhost:8080/api/v1/reading-stats/my/daily \
  -H "Authorization: Bearer {token}"
```

**结果**: ✅ 返回 30 天每日统计数据

#### 5. 获取排名
```bash
curl http://localhost:8080/api/v1/reading-stats/my/ranking \
  -H "Authorization: Bearer {token}"
```

**结果**: ✅ 返回排名信息

## 数据模型

### ChapterStats (章节统计)
```go
type ChapterStats struct {
    ID              string
    ChapterID       string
    BookID          string
    UniqueViewers   int64
    AvgReadTime     float64
    CompletionRate  float64
    DropOffRate     float64
    UpdatedAt       time.Time
}
```

### BookStats (作品统计)
```go
type BookStats struct {
    ID                 string
    BookID             string
    UniqueReaders      int64
    AvgCompletionRate  float64
    AvgDropOffRate     float64
    AvgReadingDuration float64
    TotalRevenue       float64
    ViewTrend          string
    RevenueTrend       string
    UpdatedAt          time.Time
}
```

### ReaderBehavior (读者行为)
```go
type ReaderBehavior struct {
    ID          string
    UserID      string
    BookID      string
    ChapterID   string
    Action      string
    Duration    int64
    Progress    float64
    CreatedAt   time.Time
}
```

## 核心服务方法

### ReadingStatsService

#### CalculateChapterStats
计算章节统计数据
- 统计独立读者数
- 计算平均阅读时长
- 计算完读率
- 计算跳出率

#### CalculateBookStats
计算作品统计数据
- 统计独立读者数
- 计算平均完读率
- 计算平均跳出率
- 计算总收入
- 分析趋势

#### RecordReaderBehavior
记录读者行为
- 保存行为记录
- 异步更新统计数据

#### GenerateHeatmap
生成阅读热力图
- 获取各章节阅读量
- 计算热度分数 (0-100)

## 技术要点

### 1. 异步统计更新
使用 goroutine 异步更新统计数据，避免阻塞主请求：

```go
go func() {
    bgCtx := context.Background()
    _, _ = s.CalculateChapterStats(bgCtx, behavior.ChapterID)
    _, _ = s.CalculateBookStats(bgCtx, behavior.BookID)
}()
```

### 2. 热度分数计算
综合考虑多个因素计算热度：

```go
viewScore := float64(point.ViewCount) / float64(maxViews) * 50
completionScore := point.CompletionRate * 30
dropOffScore := (1 - point.DropOffRate) * 20
point.HeatScore = viewScore + completionScore + dropOffScore
```

### 3. 趋势分析
自动分析数据趋势：
- `TrendRising`: 上升
- `TrendStable`: 稳定
- `TrendDeclining`: 下降

## 性能考虑

1. **异步处理**: 统计计算异步执行，不阻塞用户请求
2. **批量更新**: 支持批量更新统计数据
3. **缓存策略**: 可添加 Redis 缓存热门作品统计
4. **索引优化**: MongoDB 索引优化查询性能

## 后续优化

1. **缓存层**: 添加 Redis 缓存减少数据库查询
2. **实时推送**: 使用 WebSocket 实时推送统计更新
3. **批量导入**: 支持历史数据批量导入
4. **数据清理**: 定期清理过期的行为记录
5. **性能监控**: 添加 Prometheus 指标监控

## 相关文档

- [书城 API 实现](../feature/BOOKSTORE_API_IMPLEMENTATION_SUMMARY.md)
- [P0 中间件集成](MIDDLEWARE_INTEGRATION.md)
- [项目结构总结](../docs/项目结构总结.md)

## 提交历史

```
dd5f79a - fix(router): 修复重复路由注册冲突
3c738ae - feat(service): 启用 ReadingStatsService 并注册 reading-stats 路由
a5a8762 - feat(service): ReadingStatsService 实现 BaseService 接口
4129a26 - feat(repository): 在 RepositoryFactory 中添加 stats 模块仓储
```
