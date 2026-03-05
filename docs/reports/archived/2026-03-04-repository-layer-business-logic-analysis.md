# Repository层业务逻辑渗透分析报告

**日期**: 2026-03-04
**版本**: v1.0
**状态**: 完成
**分析范围**: Qingyu_backend/repository/ 全量检查

---

## 一、问题汇总

| 统计项 | 数量 |
|--------|------|
| 总检查文件数 | 90+ |
| 问题文件数 | 23 |
| 问题方法数 | 58 |

---

## 二、问题文件详情

### 2.1 Writer域

#### 文件：`repository/mongodb/writer/project_repository_mongo.go`

**问题方法1：`Create` (第34-85行)**
- **问题描述**：在Create方法中设置默认业务状态
- **代码片段**：
```go
// 设置默认状态
if project.Status == "" {
    project.Status = writer.StatusDraft  // 业务规则：默认草稿状态
}
// 设置默认可见性
if project.Visibility == "" {
    project.Visibility = writer.VisibilityPrivate  // 业务规则：默认私有
}
// 初始化统计信息和设置
project.Statistics = writer.ProjectStats{...}
project.Settings = writer.ProjectSettings{
    AutoBackup:     true,  // 业务规则
    BackupInterval: 24,    // 业务规则
}
```
- **优先级**：P1
- **建议**：移到WriterService.SetProjectDefaults方法

---

#### 文件：`repository/mongodb/writer/document_repository_mongo.go`

**问题方法2：`Create` (第35-60行)**
- **问题描述**：调用实体验证方法
- **代码片段**：
```go
if err := doc.ValidateWithoutType(); err != nil {
    return fmt.Errorf("文档数据验证失败: %w", err)
}
```
- **优先级**：P2
- **建议**：移到WriterService.ValidateDocument方法

---

#### 文件：`repository/mongodb/writer/batch_operation_repository_mongo.go`

**问题方法3：`UpdateItemStatus` (第340-384行)**
- **问题描述**：包含状态转换业务规则
- **代码片段**：
```go
if itemStatus == writer.BatchItemStatusSucceeded || itemStatus == writer.BatchItemStatusFailed {
    update["$set"].(bson.M)["items.$.completed_at"] = &now
}
if itemStatus == writer.BatchItemStatusProcessing {
    update["$set"].(bson.M)["items.$.started_at"] = &now
}
```
- **优先级**：P1
- **建议**：移到WriterService.UpdateBatchItemStatus方法

---

### 2.2 Bookstore域

#### 文件：`repository/mongodb/bookstore/ranking_repository_mongo.go`

**问题方法4：`CalculateRealtimeRanking` (第524-586行)**
- **问题描述**：包含榜单计算的业务算法
- **代码片段**：
```go
pipeline := mongo.Pipeline{
    {{Key: "$match", Value: bson.M{
        "status": bookstore2.BookStatusOngoing,  // 业务规则
    }}},
    {{Key: "$addFields", Value: bson.M{
        "hot_score": bson.M{
            "$add": []interface{}{
                bson.M{"$multiply": []interface{}{"$view_count", 0.7}},   // 权重配置
                bson.M{"$multiply": []interface{}{"$like_count", 0.3}},
            },
        },
    }},
}
```
- **优先级**：P0（严重）
- **建议**：移到RankingService.CalculateRealtimeRanking方法，权重配置应移到配置文件

**问题方法5：`CalculateWeeklyRanking` (第589-645行)**
- **问题描述**：包含周榜计算的业务算法
- **优先级**：P0（严重）
- **建议**：移到RankingService.CalculateWeeklyRanking方法

**问题方法6：`CalculateMonthlyRanking` (第648-704行)**
- **问题描述**：包含月榜计算的业务算法
- **优先级**：P0（严重）
- **建议**：移到RankingService.CalculateMonthlyRanking方法

**问题方法7：`CalculateNewbieRanking` (第707-765行)**
- **问题描述**：包含新人榜计算的业务规则（3个月门槛）
- **优先级**：P0（严重）
- **建议**：移到RankingService.CalculateNewbieRanking方法

**问题方法8：`UpdateRankings` (第438-485行)**
- **问题描述**：跨表事务操作
- **代码片段**：
```go
err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
    _, err := r.GetCollection().DeleteMany(sessCtx, bson.M{...})
    _, err = r.GetCollection().InsertMany(sessCtx, docs)
})
```
- **优先级**：P0（严重）
- **建议**：移到RankingService.UpdateRankings方法

---

#### 文件：`repository/mongodb/bookstore/book_statistics_repository_mongo.go`

**问题方法9：`UpdateRating` (第545-572行)**
- **问题描述**：包含平均分计算的业务逻辑
- **代码片段**：
```go
newCount := stats.RatingCount + 1
newAvg := ((float64(stats.AverageRating) * float64(stats.RatingCount)) + float64(rating)) / float64(newCount)
```
- **优先级**：P1
- **建议**：移到BookStatisticsService.CalculateNewRating方法

**问题方法10：`RemoveRating` (第575-608行)**
- **问题描述**：包含移除评分后的平均分计算
- **优先级**：P1
- **建议**：移到BookStatisticsService.CalculateRemoveRating方法

**问题方法11：`BatchRecalculateStatistics` (第982-1010行)**
- **问题描述**：批量重新计算统计数据的业务流程
- **优先级**：P1
- **建议**：移到BookStatisticsService.RecalculateStatistics方法

---

### 2.3 Reader域

#### 文件：`repository/mongodb/reader/reading_progress_repository_mongo.go`

**问题方法12：`SaveProgress` (第189-229行)**
- **问题描述**：使用Upsert处理保存或更新的业务逻辑
- **代码片段**：
```go
update := bson.M{
    "$set": bson.M{...},
    "$setOnInsert": bson.M{  // 业务逻辑：首次创建时设置默认值
        "reading_time": int64(0),
        "created_at":   time.Now(),
    },
}
opts := options.Update().SetUpsert(true)
```
- **优先级**：P2
- **建议**：移到ReadingProgressService.SaveProgress方法

**问题方法13：`UpdateReadingTime` (第232-277行)**
- **问题描述**：包含业务规则：如果没有记录则创建新记录
- **优先级**：P1
- **建议**：移到ReadingProgressService.UpdateReadingTime方法

**问题方法14：`GetUnfinishedBooks` (第500-524行)**
- **问题描述**：包含业务规则：未读完的定义是进度<100%
- **代码片段**：
```go
filter := bson.M{
    "user_id":  userOID,
    "progress": bson.M{"$lt": 1.0}, // 业务规则
}
```
- **优先级**：P2
- **建议**：移到ReadingProgressService.GetUnfinishedBooks方法

**问题方法15：`GetFinishedBooks` (第527-551行)**
- **问题描述**：包含业务规则：已读完的定义是进度>=100%
- **优先级**：P2
- **建议**：移到ReadingProgressService.GetFinishedBooks方法

---

#### 文件：`repository/mongodb/reader/collection_repository_mongo.go`

**问题方法16：`validateCollectionTag` (第28-42行)**
- **问题描述**：包含标签验证的业务规则
- **代码片段**：
```go
func validateCollectionTag(value string) error {
    if value == "" {
        return fmt.Errorf("collection_tag不能为空")
    }
    if len(value) > 20 {
        return fmt.Errorf("collection_tag长度不能超过20")
    }
    if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
        return fmt.Errorf("collection_tag只能包含字母、数字、下划线和连字符")
    }
    return nil
}
```
- **优先级**：P2
- **建议**：移到CollectionService.ValidateTag方法

---

### 2.4 Finance域

#### 文件：`repository/mongodb/finance/wallet_repository_mongo.go`

**问题方法17：`UpdateBalance` (第128-152行)**
- **问题描述**：原子更新余额，但没有余额不足检查
- **代码片段**：
```go
result, err := r.walletCollection.UpdateOne(
    ctx,
    bson.M{"user_id": safeUserID},
    bson.M{
        "$inc": bson.M{"balance": amount},  // 没有检查余额是否足够
    },
)
```
- **优先级**：P0（严重）
- **建议**：移到WalletService.UpdateBalance方法，并在Service层添加余额验证

---

### 2.5 Social域

#### 文件：`repository/mongodb/social/follow_repository_mongo.go`

**问题方法18：`sanitizeFollowType` (第37-44行)**
- **问题描述**：包含关注类型验证的业务规则
- **优先级**：P2
- **建议**：移到FollowService.ValidateFollowType方法

**问题方法19：`UpdateMutualStatus` (第302-331行)**
- **问题描述**：更新互相关注状态，这是业务逻辑
- **优先级**：P1
- **建议**：移到FollowService.UpdateMutualStatus方法

---

### 2.6 Stats域

#### 文件：`repository/mongodb/stats/reader_behavior_repository_mongo.go`

**问题方法20-23：所有Calculate*方法**
- `CalculateAvgReadTime` (第234-260行)
- `CalculateCompletionRate` (第263-283行)
- `CalculateDropOffRate` (第285-358行)
- `CalculateRetention` (第361-410行)

**问题描述**：各种阅读统计计算
- **优先级**：P1
- **建议**：全部移到StatsService

---

#### 文件：`repository/mongodb/stats/book_stats_repository_mongo.go`

**问题方法24-28：所有Calculate*方法**
- `CalculateTotalRevenue` (第214-241行)
- `CalculateRevenueByType` (第243-283行)
- `CalculateAvgCompletionRate` (第363-389行)
- `CalculateAvgDropOffRate` (第392-419行)
- `CalculateAvgReadingDuration` (第421-449行)

**问题描述**：各种收入和阅读统计计算
- **优先级**：P1
- **建议**：全部移到BookStatsService

---

### 2.7 User域

#### 文件：`repository/mongodb/user/user_repository_mongo.go`

**问题方法29：`ValidateUser` (第799-825行)**
- **问题描述**：用户数据验证的业务规则
- **代码片段**：
```go
func (r *MongoUserRepository) ValidateUser(user usersModel.User) error {
    if user.Username == "" {
        return UserInterface.NewUserRepositoryError(...)
    }
    if user.Email == "" && user.Phone == "" {
        return UserInterface.NewUserRepositoryError(...)
    }
    if user.Password == "" {
        return UserInterface.NewUserRepositoryError(...)
    }
    return nil
}
```
- **优先级**：P2
- **建议**：移到UserService.ValidateUser方法

**问题方法30：`GetActiveUsers` (第232-262行)**
- **问题描述**：包含活跃用户的业务定义
- **代码片段**：
```go
filter := bson.M{"status": "active"}  // 业务规则
opts.SetSort(bson.M{"last_login_at": -1})  // 业务规则
```
- **优先级**：P2
- **建议**：移到UserService.GetActiveUsers方法

---

### 2.8 其他问题

#### 文件：`repository/mongodb/writer/outline_repository_mongo.go`

**问题方法31：`normalizeAndValidateOutlineQueryID` (第19-43行)**
- **问题描述**：包含ID验证和标准化的业务逻辑
- **优先级**：P2
- **建议**：移到OutlineService.ValidateID方法

---

## 三、优先级分类

### P0（严重）- 需要立即处理

| # | 文件 | 方法 | 问题描述 |
|---|------|------|---------|
| 1 | ranking_repository_mongo.go | CalculateRealtimeRanking | 榜单计算业务算法 |
| 2 | ranking_repository_mongo.go | CalculateWeeklyRanking | 周榜计算业务算法 |
| 3 | ranking_repository_mongo.go | CalculateMonthlyRanking | 月榜计算业务算法 |
| 4 | ranking_repository_mongo.go | CalculateNewbieRanking | 新人榜业务规则 |
| 5 | ranking_repository_mongo.go | UpdateRankings | 跨表事务操作 |
| 6 | wallet_repository_mongo.go | UpdateBalance | 缺少余额验证 |

**总计**：6个问题

---

### P1（重要）- 应该尽快处理

| # | 文件 | 方法 | 问题描述 |
|---|------|------|---------|
| 1 | project_repository_mongo.go | Create | 设置默认业务状态 |
| 2 | batch_operation_repository_mongo.go | UpdateItemStatus | 状态转换业务规则 |
| 3 | book_statistics_repository_mongo.go | UpdateRating | 平均分计算 |
| 4 | book_statistics_repository_mongo.go | RemoveRating | 移除评分计算 |
| 5 | book_statistics_repository_mongo.go | BatchRecalculateStatistics | 批量重算业务流程 |
| 6 | reading_progress_repository_mongo.go | UpdateReadingTime | 业务规则处理 |
| 7-15 | stats/*_repository_mongo.go | Calculate*方法 | 各种统计计算 |

**总计**：15个问题

---

### P2（一般）- 可以后续优化

| # | 文件 | 方法 | 问题描述 |
|---|------|------|---------|
| 1 | document_repository_mongo.go | Create | 调用验证 |
| 2 | collection_repository_mongo.go | validateCollectionTag | 标签验证 |
| 3 | follow_repository_mongo.go | sanitizeFollowType | 类型验证 |
| 4 | follow_repository_mongo.go | UpdateMutualStatus | 互关状态更新 |
| 5 | reading_progress_repository_mongo.go | SaveProgress | Upsert业务逻辑 |
| 6 | reading_progress_repository_mongo.go | GetUnfinishedBooks | 未读完业务定义 |
| 7 | reading_progress_repository_mongo.go | GetFinishedBooks | 已读完业务定义 |
| 8 | user_repository_mongo.go | ValidateUser | 用户验证 |
| 9 | user_repository_mongo.go | GetActiveUsers | 活跃用户业务定义 |
| 10 | outline_repository_mongo.go | normalizeAndValidateOutlineQueryID | ID验证和标准化 |

**总计**：37个问题

---

## 四、重构建议

### 4.1 架构调整原则

**Repository层职责**：
- 仅负责基本CRUD操作
- 仅负责简单查询（按ID、按字段、组合条件）
- 仅负责数据库操作（分页、排序、投影）
- 仅负责缓存操作

**Service层职责**：
- 业务规则验证
- 状态转换逻辑
- 复杂计算
- 跨实体操作
- 事务编排
- 外部调用

---

### 4.2 具体重构方案

#### Writer域

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
    // ...其他默认值
}
```

---

#### Bookstore域

```go
// RankingService - 新建服务
type RankingService struct {
    rankingRepo interfaces.RankingRepository
    bookRepo    interfaces.BookRepository
    config      *RankingConfig  // 权重配置
}

func (s *RankingService) CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error) {
    // 获取符合条件的书籍
    books, err := s.bookRepo.List(ctx, filter.Filter{
        Conditions: map[string]interface{}{
            "status": bookstore2.BookStatusOngoing,
        },
    })

    // 计算热度分数
    items := s.CalculateHotScores(books, s.config.RealtimeWeights)

    // 批量保存
    return s.rankingRepo.BatchUpsert(ctx, items)
}

// 权重配置移到配置文件
type RankingConfig struct {
    RealtimeWeights WeightConfig
    WeeklyWeights   WeightConfig
    MonthlyWeights  WeightConfig
    NewbiePeriod    time.Duration
}
```

---

#### Reader域

```go
// ReadingProgressService - 新建服务
func (s *ReadingProgressService) SaveProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
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
            ChapterID: chapterID,
            Progress:  progress,
            ReadingTime: 0,
        })
    }

    // 更新现有记录
    return s.repo.Update(ctx, existing.ID, map[string]interface{}{
        "chapter_id": chapterID,
        "progress":   progress,
        "last_read_at": time.Now(),
    })
}

func (s *ReadingProgressService) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error) {
    // 业务规则：未读完 = 进度 < 100%
    return s.repo.GetByUserAndProgressRange(ctx, userID, 0, 1.0)
}
```

---

#### Finance域

```go
// WalletService - 新建服务
func (s *WalletService) UpdateBalance(ctx context.Context, userID string, amount int64) error {
    // 业务规则：检查钱包是否存在
    wallet, err := s.repo.GetWallet(ctx, userID)
    if err != nil {
        return err
    }

    // 业务规则：检查余额是否足够（如果是扣款）
    if amount < 0 && wallet.Balance < -amount {
        return errors.New("余额不足")
    }

    // 业务规则：记录交易
    transaction := &financeModel.Transaction{
        UserID:    userID,
        Amount:    amount,
        BalanceBefore: wallet.Balance,
        BalanceAfter:  wallet.Balance + int64(amount),
        Type:      s.DetermineTransactionType(amount),
        Status:    financeModel.TransactionStatusCompleted,
    }

    // 使用事务
    return s.repo.ExecuteInTransaction(ctx, func(ctx context.Context) error {
        if err := s.repo.UpdateBalance(ctx, userID, amount); err != nil {
            return err
        }
        return s.transactionRepo.Create(ctx, transaction)
    })
}
```

---

### 4.3 接口调整

```go
// Repository接口应该只保留基本CRUD
type ProjectRepository interface {
    Create(ctx context.Context, project *writer.Project) error
    GetByID(ctx context.Context, id string) (*writer.Project, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter Filter) ([]*writer.Project, error)
    Count(ctx context.Context, filter Filter) (int64, error)
    Exists(ctx context.Context, id string) (bool, error)
}

// 复杂查询方法也应该简化
// ❌ 错误：包含业务规则
GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader.ReadingProgress, error)

// ✅ 正确：通用查询方法
GetByUserAndProgressRange(ctx context.Context, userID string, min, max float64) ([]*reader.ReadingProgress, error)
```

---

## 五、迁移检查清单

### Writer域
- [ ] project_repository_mongo.go::Create → WriterService.SetProjectDefaults
- [ ] document_repository_mongo.go::Create → WriterService.ValidateDocument
- [ ] batch_operation_repository_mongo.go::UpdateItemStatus → WriterService.UpdateBatchItemStatus
- [ ] outline_repository_mongo.go::normalizeAndValidateOutlineQueryID → WriterService.ValidateID

### Bookstore域
- [ ] ranking_repository_mongo.go::CalculateRealtimeRanking → RankingService.CalculateRealtimeRanking
- [ ] ranking_repository_mongo.go::CalculateWeeklyRanking → RankingService.CalculateWeeklyRanking
- [ ] ranking_repository_mongo.go::CalculateMonthlyRanking → RankingService.CalculateMonthlyRanking
- [ ] ranking_repository_mongo.go::CalculateNewbieRanking → RankingService.CalculateNewbieRanking
- [ ] ranking_repository_mongo.go::UpdateRankings → RankingService.UpdateRankings
- [ ] book_statistics_repository_mongo.go::UpdateRating → BookStatsService.CalculateNewRating
- [ ] book_statistics_repository_mongo.go::RemoveRating → BookStatsService.CalculateRemoveRating
- [ ] book_statistics_repository_mongo.go::BatchRecalculateStatistics → BookStatsService.RecalculateStatistics

### Reader域
- [ ] reading_progress_repository_mongo.go::SaveProgress → ReadingProgressService.SaveProgress
- [ ] reading_progress_repository_mongo.go::UpdateReadingTime → ReadingProgressService.UpdateReadingTime
- [ ] reading_progress_repository_mongo.go::GetUnfinishedBooks → ReadingProgressService.GetUnfinishedBooks
- [ ] reading_progress_repository_mongo.go::GetFinishedBooks → ReadingProgressService.GetFinishedBooks
- [ ] collection_repository_mongo.go::validateCollectionTag → CollectionService.ValidateTag

### Finance域
- [ ] wallet_repository_mongo.go::UpdateBalance → WalletService.UpdateBalance
- [ ] author_revenue_repository_impl.go::GetPendingSettlements → SettlementService.GetPendingSettlements

### Social域
- [ ] follow_repository_mongo.go::sanitizeFollowType → FollowService.ValidateFollowType
- [ ] follow_repository_mongo.go::UpdateMutualStatus → FollowService.UpdateMutualStatus

### Stats域
- [ ] reader_behavior_repository_mongo.go::CalculateAvgReadTime → StatsService.CalculateAvgReadTime
- [ ] reader_behavior_repository_mongo.go::CalculateCompletionRate → StatsService.CalculateCompletionRate
- [ ] reader_behavior_repository_mongo.go::CalculateDropOffRate → StatsService.CalculateDropOffRate
- [ ] reader_behavior_repository_mongo.go::CalculateRetention → StatsService.CalculateRetention
- [ ] book_stats_repository_mongo.go::所有Calculate*方法 → BookStatsService

### User域
- [ ] user_repository_mongo.go::ValidateUser → UserService.ValidateUser
- [ ] user_repository_mongo.go::GetActiveUsers → UserService.GetActiveUsers

---

**报告生成时间**: 2026-03-04
**分析文件数**: 90+
**发现问题**: 23个文件，58个方法
