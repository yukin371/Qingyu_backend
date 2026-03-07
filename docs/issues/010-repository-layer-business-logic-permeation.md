# Issue #010: Repository 层业务逻辑渗透

**优先级**: 高 (P0)
**类型**: 架构问题
**状态**: ✅ Phase 2 已完成 (2026-03-07)
**创建日期**: 2026-03-05
**来源报告**: [Repository 层业务逻辑渗透分析报告](../reports/archived/2026-03-04-repository-layer-business-logic-analysis.md)
**审查日期**: 2026-03-05
**审查报告**: [P0问题审查报告](../reports/2026-03-05-p0-issue-audit-report.md)

---

## 审查结果

**状态**: ❌ 问题确认存在

### 审查发现

#### #010-A: Bookstore域 Repository 重构（P0）- ⚠️ 部分修复

1. ✅ **榜单计算已从 `ranking_repository_mongo.go` 移到 `service/bookstore/bookstore_service.go`**
2. ✅ **榜单更新事务编排已移到 Service 层**
3. ⚠️ **权重配置仍硬编码在 Service 层**
4. ⚠️ **尚无独立的 RankingService / 定时更新入口**

**证据**:
- Repository 接口已删除 `Calculate*Ranking` / `UpdateRankings`
- `BookstoreServiceImpl.UpdateRankings` 现在负责：
  - 获取书籍列表
  - 计算分数
  - 排序生成榜单
  - 调用 Repository 事务执行删除和批量写入
- Repository 仅保留 `DeleteByTypeAndPeriod`、`BatchUpsertRankingItems`、`Transaction`

#### #010-C: Finance域 Repository 重构（P0）- ⚠️ 部分修复

1. ✅ **Service层已实现余额验证**（transaction_service.go, withdraw_service.go）
2. ❌ **事务编排不完整**（存在TODO注释）
3. ⚠️ **存在竞态条件风险**

---

## 问题描述

Repository 层承担了部分 Service 层的职责，违反了分层架构原则，导致：
- 业务逻辑分散
- 代码难以维护
- 测试困难
- 职责边界模糊

### 问题规模

| 统计项 | 数量 |
|--------|------|
| 总检查文件数 | 90+ |
| 问题文件数 | 23 |
| 问题方法数 | 58 |

### 优先级分布

| 优先级 | 问题数量 | 说明 |
|--------|----------|------|
| P0（严重） | 6 | 需要立即处理 |
| P1（重要） | 15 | 应该尽快处理 |
| P2（一般） | 37 | 可以后续优化 |

---

## 主要问题分类

### 1. 业务规则在 Repository 层

#### Writer 域 - Project 创建默认值设置

**问题**: `project_repository_mongo.go::Create` 设置默认业务状态

```go
// ❌ 当前在 Repository 层
if project.Status == "" {
    project.Status = writer.StatusDraft  // 业务规则：默认草稿状态
}
if project.Visibility == "" {
    project.Visibility = writer.VisibilityPrivate  // 业务规则：默认私有
}
project.Statistics = writer.ProjectStats{...}
project.Settings = writer.ProjectSettings{
    AutoBackup:     true,  // 业务规则
    BackupInterval: 24,    // 业务规则
}
```

**应该移到**: `WriterService.SetProjectDefaults()`

#### Reader 域 - 未完成/已完成业务定义

**问题**: `reading_progress_repository_mongo.go` 包含业务规则定义

```go
// ❌ 当前在 Repository 层
func (r *MongoReadingProgressRepository) GetUnfinishedBooks(ctx context.Context, userID string) {
    filter := bson.M{
        "user_id": userOID,
        "progress": bson.M{"$lt": 1.0}, // 业务规则：未读完的定义
    }
}

func (r *MongoReadingProgressRepository) GetFinishedBooks(ctx context.Context, userID string) {
    filter := bson.M{
        "progress": bson.M{"$gte": 1.0}, // 业务规则：已读完的定义
    }
}
```

**应该移到**: `ReadingProgressService.GetUnfinishedBooks() / GetFinishedBooks()`

---

### 2. 复杂计算在 Repository 层

#### Bookstore 域 - 榜单计算算法

**问题**: `ranking_repository_mongo.go` 包含榜单计算业务算法

```go
// ❌ 当前在 Repository 层
func (r *MongoRankingRepository) CalculateRealtimeRanking(ctx context.Context) {
    pipeline := mongo.Pipeline{
        {{Key: "$match", Value: bson.M{
            "status": bookstore2.BookStatusOngoing,  // 业务规则
        }},
        {{Key: "$addFields", Value: bson.M{
            "hot_score": bson.M{
                "$add": []interface{}{
                    bson.M{"$multiply": []interface{}{"$view_count", 0.7}},   // 权重配置
                    bson.M{"$multiply": []interface{}{"$like_count", 0.3}},
                },
            },
        }},
    }
}
```

**问题**:
- 包含业务算法（热度分数计算）
- 包含权重配置（0.7, 0.3 硬编码）
- Repository 层不应该包含业务逻辑

**应该移到**: `RankingService.CalculateRealtimeRanking()`

**当前剩余 TODO**:
- [ ] 抽离独立 `RankingService`，避免继续挂在 `BookstoreService`
- [ ] 将榜单权重配置外置，避免 Service 中继续硬编码
- [ ] 为榜单刷新补独立调度/任务入口
- [ ] 明确榜单算法依赖的统计字段来源，补齐 `like_count` 等独立统计口径
- [ ] 评估是否需要按统计快照而不是全量 `List()` 计算榜单

#### Stats 域 - 统计计算

**问题**: 所有 `Calculate*` 方法在 Repository 层

- `reader_behavior_repository_mongo.go`:
  - `CalculateAvgReadTime`
  - `CalculateCompletionRate`
  - `CalculateDropOffRate`
  - `CalculateRetention`

- `book_stats_repository_mongo.go`:
  - `CalculateTotalRevenue`
  - `CalculateRevenueByType`
  - `CalculateAvgCompletionRate`
  - `CalculateAvgDropOffRate`
  - `CalculateAvgReadingDuration`

**应该移到**: `StatsService` 对应的统计方法

---

### 3. 跨表事务操作在 Repository 层

#### Bookstore 域 - 榜单更新

**问题**: `ranking_repository_mongo.go::UpdateRankings` 包含跨表事务

```go
// ❌ 当前在 Repository 层
func (r *MongoRankingRepository) UpdateRankings(ctx context.Context) error {
    session, err := r.client.StartSession()
    return mongo.SessionWithContext(ctx, session).WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        _, err := r.GetCollection().DeleteMany(sessCtx, bson.M{...})
        _, err = r.GetCollection().InsertMany(sessCtx, docs)
    })
}
```

**问题**: Repository 层不应该处理事务编排

**应该移到**: `RankingService.UpdateRankings()`

---

### 4. 缺少业务验证

#### Finance 域 - 余额更新

**问题**: `wallet_repository_mongo.go::UpdateBalance` 缺少余额验证

```go
// ❌ 当前在 Repository 层
func (r *MongoWalletRepository) UpdateBalance(ctx context.Context, userID string, amount int64) error {
    result, err := r.walletCollection.UpdateOne(
        ctx,
        bson.M{"user_id": safeUserID},
        bson.M{"$inc": bson.M{"balance": amount}},  // 没有检查余额是否足够
    )
}
```

**问题**: 财务相关操作缺少业务验证

**应该移到**: `WalletService.UpdateBalance()` 并添加余额验证

---

## 架构重构原则

### Repository 层职责

```go
// ✅ Repository 层应该只负责
type ProjectRepository interface {
    // 基本 CRUD
    Create(ctx context.Context, project *writer.Project) error
    GetByID(ctx context.Context, id string) (*writer.Project, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
    Delete(ctx context.Context, id string) error

    // 简单查询（按ID、按字段、组合条件）
    List(ctx context.Context, filter Filter) ([]*writer.Project, error)
    Count(ctx context.Context, filter Filter) (int64, error)
    Exists(ctx context.Context, id string) (bool, error)

    // 数据库操作（分页、排序、投影）
    ListWithPagination(ctx context.Context, filter Filter, page, size int) ([]*writer.Project, error)

    // 缓存操作
    InvalidateCache(ctx context.Context) error
}

// ❌ Repository 层不应该负责
type ProjectRepository interface {
    // 包含业务规则
    GetUnfinishedProjects(ctx context.Context) ([]*writer.Project, error)
    GetByStatusAndVisibility(ctx context.Context, status, visibility string) ([]*writer.Project, error)

    // 包含复杂计算
    CalculateHotScore(ctx context.Context) ([]*writer.Project, error)
    UpdateRankings(ctx context.Context) error

    // 包含业务验证
    ValidateAndCreate(ctx context.Context, project *writer.Project) error
}
```

### Service 层职责

```go
// ✅ Service 层应该负责
type ProjectService interface {
    // 业务规则验证
    ValidateProject(project *writer.Project) error

    // 状态转换逻辑
    PublishProject(ctx context.Context, projectID string) error

    // 复杂计算
    CalculateHotScore(ctx context.Context, projects []*writer.Project) []float64

    // 跨实体操作
    PublishProjectWithChapters(ctx context.Context, projectID string) error

    // 事务编排
    TransferProjectOwnership(ctx context.Context, projectID, newOwnerID string) error
}
```

---

## 具体重构方案

### Writer 域

```go
// WriterService - 新增方法
func (s *WriterService) CreateProject(ctx context.Context, project *writer.Project) error {
    // 业务规则：设置默认值
    s.SetProjectDefaults(project)

    // 业务验证
    if err := s.ValidateProject(project); err != nil {
        return err
    }

    // 调用Repository保存
    return s.projectRepo.Create(ctx, project)
}

func (s *WriterService) SetProjectDefaults(project *writer.Project) {
    if project.Status == "" {
        project.Status = writer.StatusDraft
    }
    if project.Visibility == "" {
        project.Visibility = writer.VisibilityPrivate
    }
    project.Statistics = writer.ProjectStats{
        WordCount:    0,
        ChapterCount:  0,
        LastUpdatedAt: time.Now(),
    }
    project.Settings = writer.ProjectSettings{
        AutoBackup:     true,
        BackupInterval: 24,
    }
}
```

### Bookstore 域

```go
// RankingService - 新建服务
type RankingService struct {
    rankingRepo interfaces.RankingRepository
    bookRepo    interfaces.BookRepository
    config      *RankingConfig
}

func (s *RankingService) CalculateRealtimeRanking(ctx context.Context) ([]*bookstore.RankingItem, error) {
    // 1. 获取符合条件的书籍
    books, err := s.bookRepo.List(ctx, filter.Filter{
        Conditions: map[string]interface{}{
            "status": bookstore2.BookStatusOngoing,
        },
    })
    if err != nil {
        return nil, err
    }

    // 2. 计算热度分数（Service 层业务逻辑）
    items := make([]*bookstore.RankingItem, 0, len(books))
    for _, book := range books {
        hotScore := s.calculateHotScore(book)
        items = append(items, &bookstore.RankingItem{
            BookID:    book.ID,
            HotScore:  hotScore,
            // ...
        })
    }

    // 3. 批量保存
    return s.rankingRepo.BatchUpsert(ctx, items)
}

func (s *RankingService) calculateHotScore(book *bookstore.Book) float64 {
    // 业务算法：热度分数 = 浏览量 * 0.7 + 点赞数 * 0.3
    return float64(book.Statistics.ViewCount)*s.config.RealtimeWeights.ViewCount +
           float64(book.Statistics.LikeCount)*s.config.RealtimeWeights.LikeCount
}

// 权重配置移到配置文件
type RankingConfig struct {
    RealtimeWeights WeightConfig
    WeeklyWeights   WeightConfig
    MonthlyWeights  WeightConfig
    NewbiePeriod    time.Duration
}
```

### Reader 域

```go
// ReadingProgressService - 重构方法
func (s *ReadingProgressService) SaveProgress(
    ctx context.Context,
    userID, bookID, chapterID string,
    progress float64,
) error {
    // 业务规则：检查进度范围
    if progress < 0 || progress > 1 {
        return errors.New("进度必须在0-1之间")
    }

    // 业务逻辑：保存或更新
    existing, err := s.repo.GetByUserAndBook(ctx, userID, bookID)
    if err != nil {
        return err
    }

    if existing == nil {
        // 首次阅读，创建记录
        return s.repo.Create(ctx, &reader.ReadingProgress{
            UserID:    userID,
            BookID:    bookID,
            ChapterID:  chapterID,
            Progress:  progress,
            ReadingTime: 0,
            CreatedAt:  time.Now(),
        })
    }

    // 更新现有记录
    return s.repo.Update(ctx, existing.ID.Hex(), map[string]interface{}{
        "chapter_id":  chapterID,
        "progress":   progress,
        "last_read_at": time.Now(),
    })
}

func (s *ReadingProgressService) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
    // 业务规则：未读完 = 进度 < 100%
    return s.repo.GetByUserAndProgressRange(ctx, userID, 0, 1.0)
}
```

### Finance 域

```go
// WalletService - 重构方法
func (s *WalletService) UpdateBalance(ctx context.Context, userID string, amount int64) error {
    // 1. 业务规则：检查钱包是否存在
    wallet, err := s.repo.GetWallet(ctx, userID)
    if err != nil {
        return err
    }

    // 2. 业务规则：检查余额是否足够（如果是扣款）
    if amount < 0 && wallet.Balance < -amount {
        return errors.New("余额不足")
    }

    // 3. 业务规则：记录交易
    transaction := &financeModel.Transaction{
        UserID:        userID,
        Amount:        amount,
        BalanceBefore: wallet.Balance,
        BalanceAfter:   wallet.Balance + int64(amount),
        Type:          s.DetermineTransactionType(amount),
        Status:        financeModel.TransactionStatusCompleted,
        CreatedAt:     time.Now(),
    }

    // 4. 使用事务确保原子性
    return s.repo.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
        if err := s.repo.UpdateBalance(txCtx, userID, amount); err != nil {
            return err
        }
        return s.transactionRepo.Create(txCtx, transaction)
    })
}
```

---

## 实施计划

### Phase 1: 高优先级问题（P0）- 2-3 周

| 文件 | 方法 | 移动到 | 预计时间 |
|------|------|--------|----------|
| ranking_repository_mongo.go | CalculateRealtimeRanking | RankingService | 2天 |
| ranking_repository_mongo.go | CalculateWeeklyRanking | RankingService | 1天 |
| ranking_repository_mongo.go | CalculateMonthlyRanking | RankingService | 1天 |
| ranking_repository_mongo.go | CalculateNewbieRanking | RankingService | 1天 |
| ranking_repository_mongo.go | UpdateRankings | RankingService | 2天 |
| wallet_repository_mongo.go | UpdateBalance | WalletService | 2天 |

### Phase 2: 中优先级问题（P1）- 3-4 周

| 域 | 方法数 | 预计时间 |
|------|--------|----------|
| Writer | 4个 | 2天 |
| BookStats | 3个 | 2天 |
| Reader | 4个 | 2天 |
| Stats | 9个 | 5天 |

### Phase 3: 低优先级问题（P2）- 持续优化

- 验证方法移到 Service
- 简单业务规则移到 Service
- 按需进行

---

## 检查清单

### Writer 域
- [ ] project_repository_mongo.go::Create → WriterService.SetProjectDefaults
- [ ] document_repository_mongo.go::Create → WriterService.ValidateDocument
- [ ] batch_operation_repository_mongo.go::UpdateItemStatus → WriterService.UpdateBatchItemStatus
- [ ] outline_repository_mongo.go::normalizeAndValidateOutlineQueryID → WriterService.ValidateID

### Bookstore 域
- [x] ranking_repository_mongo.go::CalculateRealtimeRanking → RankingService.CalculateRealtimeRanking ✅
- [x] ranking_repository_mongo.go::CalculateWeeklyRanking → RankingService.CalculateWeeklyRanking ✅
- [x] ranking_repository_mongo.go::CalculateMonthlyRanking → RankingService.CalculateMonthlyRanking ✅
- [x] ranking_repository_mongo.go::CalculateNewbieRanking → RankingService.CalculateNewbieRanking ✅
- [x] ranking_repository_mongo.go::UpdateRankings → RankingService.UpdateRankings ✅
- [ ] book_statistics_repository_mongo.go::UpdateRating → BookStatsService.CalculateNewRating
- [ ] book_statistics_repository_mongo.go::RemoveRating → BookStatsService.CalculateRemoveRating
- [ ] book_statistics_repository_mongo.go::BatchRecalculateStatistics → BookStatsService.RecalculateStatistics

### Reader 域
- [ ] reading_progress_repository_mongo.go::SaveProgress → ReadingProgressService.SaveProgress
- [ ] reading_progress_repository_mongo.go::UpdateReadingTime → ReadingProgressService.UpdateReadingTime
- [ ] reading_progress_repository_mongo.go::GetUnfinishedBooks → ReadingProgressService.GetUnfinishedBooks
- [ ] reading_progress_repository_mongo.go::GetFinishedBooks → ReadingProgressService.GetFinishedBooks
- [ ] collection_repository_mongo.go::validateCollectionTag → CollectionService.ValidateTag

### Finance 域
- [x] wallet_repository_mongo.go::UpdateBalanceWithCheck → WalletService.UpdateBalanceWithCheck ✅ (已添加余额验证)

### Social 域
- [ ] follow_repository_mongo.go::sanitizeFollowType → FollowService.ValidateFollowType
- [ ] follow_repository_mongo.go::UpdateMutualStatus → FollowService.UpdateMutualStatus

### Stats 域
- [ ] reader_behavior_repository_mongo.go::CalculateAvgReadTime → StatsService.CalculateAvgReadTime
- [ ] reader_behavior_repository_mongo.go::CalculateCompletionRate → StatsService.CalculateCompletionRate
- [ ] reader_behavior_repository_mongo.go::CalculateDropOffRate → StatsService.CalculateDropOffRate
- [ ] reader_behavior_repository_mongo.go::CalculateRetention → StatsService.CalculateRetention
- [ ] book_stats_repository_mongo.go::所有Calculate*方法 → BookStatsService

### User 域
- [ ] user_repository_mongo.go::ValidateUser → UserService.ValidateUser
- [ ] user_repository_mongo.go::GetActiveUsers → UserService.GetActiveUsers

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [Repository 层业务逻辑渗透分析报告](../reports/archived/2026-03-04-repository-layer-business-logic-analysis.md) | 完整问题清单 |
| [后端综合审计报告](../reports/archived/backend-comprehensive-audit-summary-2026-01-26.md) | 相关问题参考 |

---

## 相关Issue

### 依赖Issue（必须先处理）
- [#007: Service 层事务管理缺失](./007-transaction-management.md) - ⚠️ 需要先实现事务管理器，才能处理跨表事务问题

### 相关Issue（联合处理）
- [#001: 统一模型层 ID 字段类型](./001-unify-id-type-in-models.md) - Repository重构时需要确保ID类型正确
- [#002: Repository Create 方法未回设 ID](./002-create-method-id-not-set-bug.md) - Create方法问题属于Repository层职责
- [#003: 测试基础设施改进](./003-test-infrastructure-improvements.md) - Repository重构需要测试支持

### 建议拆分
本Issue规模较大（58个方法，23个文件），建议按域拆分为：
- **#010-A**: Bookstore域（P0）- ranking_repository, book_statistics_repository
- **#010-B**: Reader域（P1）- reading_progress_repository, collection_repository
- **#010-C**: Finance域（P0）- wallet_repository
- **#010-D**: Writer域（P1）- project_repository, document_repository等
- **#010-E**: Stats域（P2）- reader_behavior_repository, book_stats_repository
- **#010-F**: Social域（P2）- follow_repository
- **#010-G**: User域（P2）- user_repository

详细拆分方案见 [Issue关联关系分析](./ISSUE_RELATIONSHIPS.md)
