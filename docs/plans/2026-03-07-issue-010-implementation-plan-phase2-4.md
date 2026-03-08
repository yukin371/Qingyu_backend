# Issue #010 实施计划 Phase 2-4

**创建日期**: 2026-03-07
**状态**: 计划中
**优先级**: P1-P2
**预计总工期**: 4-6 周

---

## 概述

本文档详细规划了 Issue #010 (Repository 层业务逻辑渗透) Phase 2-4 的实施步骤，Phase 1 已于 2026-03-07 完成。

---

## Phase 1 完成总结 ✅

### 已完成的工作
1. **RankingService 创建** (`service/bookstore/ranking_service.go`)
   - RankingConfig 配置结构体（权重配置化）
   - CalculateRealtimeRanking/Weekly/Monthly/Newbie 方法
   - UpdateRankings 方法（含事务编排）
   - 榜单查询方法

2. **WalletRepository 余额验证** (`repository/mongodb/finance/wallet_repository_mongo.go`)
   - UpdateBalanceWithCheck 方法
   - MongoDB 条件更新防止负数余额

---

## Phase 2: 中优先级问题（P1）- 预计 1-2 周

### 目标
将 Writer 域和 Reader 域、业务逻辑从 Repository 移到 Service 层

### 2.1 Writer 域重构（3 天）

| 优先级 | 文件 | 方法 | 移动到 | 复杂度 |
|--------|------|------|--------|--------|
| 🔴 P1 | project_repository_mongo.go | Create (设置默认值) | WriterService.SetProjectDefaults | 中 |
| 🔴 P1 | document_repository_mongo.go | Create (设置默认值) | WriterService.ValidateDocument | 低 |
| 🟡 P2 | batch_operation_repository_mongo.go | UpdateItemStatus | WriterService.UpdateBatchItemStatus | 低 |
| 🟡 P2 | outline_repository_mongo.go | normalizeAndValidateOutlineQueryID | WriterService.ValidateID | 低 |

**实施步骤**:
1. 在 `service/writer/project_service.go` 添加 `SetProjectDefaults` 方法
2. 在 `service/writer/document_service.go` 添加 `ValidateDocument` 方法
3. 修改 Repository 的 Create 方法，移除默认值设置逻辑
4. 添加单元测试验证

### 2.2 Reader 域重构（3 天）

| 优先级 | 文件 | 方法 | 移动到 | 复杂度 |
|--------|------|------|--------|--------|
| 🔴 P1 | reading_progress_repository_mongo.go | SaveProgress | ReadingProgressService.SaveProgress | 中 |
| 🔴 P1 | reading_progress_repository_mongo.go | UpdateReadingTime | ReadingProgressService.UpdateReadingTime | 低 |
| 🟡 P2 | reading_progress_repository_mongo.go | GetUnfinishedBooks | ReadingProgressService.GetUnfinishedBooks | 低 |
| 🟡 P2 | reading_progress_repository_mongo.go | GetFinishedBooks | ReadingProgressService.GetFinishedBooks | 低 |
| 🟡 P2 | collection_repository_mongo.go | validateCollectionTag | CollectionService.ValidateTag | 低 |

**实施步骤**:
1. 在 `service/reader/reading_progress_service.go` 添加相应方法
2. 修改 Repository 方法，移除业务规则判断
3. 添加单元测试验证

### 2.3 Bookstore 域剩余重构（2 天）

| 优先级 | 文件 | 方法 | 移动到 | 复杂度 |
|--------|------|------|--------|--------|
| 🟡 P2 | book_statistics_repository_mongo.go | UpdateRating | BookStatsService.CalculateNewRating | 中 |
| 🟡 P2 | book_statistics_repository_mongo.go | RemoveRating | BookStatsService.CalculateRemoveRating | 中 |
| 🟡 P2 | book_statistics_repository_mongo.go | BatchRecalculateStatistics | BookStatsService.RecalculateStatistics | 高 |

---

## Phase 3: 低优先级问题（P2）- 预计 1-2 周

### 目标
Social 域和 User 域的业务逻辑重构

### 3.1 Social 域重构（1 天）

| 优先级 | 文件 | 方法 | 移动到 | 复杂度 |
|--------|------|------|--------|--------|
| 🟡 P2 | follow_repository_mongo.go | sanitizeFollowType | FollowService.ValidateFollowType | 低 |
| 🟡 P2 | follow_repository_mongo.go | UpdateMutualStatus | FollowService.UpdateMutualStatus | 低 |

### 3.2 User 域重构（1 天）

| 优先级 | 文件 | 方法 | 移动到 | 复杂度 |
|--------|------|------|--------|--------|
| 🟡 P2 | user_repository_mongo.go | ValidateUser | UserService.ValidateUser | 低 |
| 🟡 P2 | user_repository_mongo.go | GetActiveUsers | UserService.GetActiveUsers | 低 |

---

## Phase 4: Stats 域重构（P2）- 预计 1 周

### 目标
将所有 Calculate* 方法移到 Service 层

### 4.1 Stats 域重构（1 周）

| 优先级 | 文件 | 方法数 | 说明 |
|--------|------|--------|------|
| 🟡 P2 | reader_behavior_repository_mongo.go | 4个 | CalculateAvgReadTime/CompletionRate/DropOffRate/Retention |
| 🟡 P2 | book_stats_repository_mongo.go | 5个 | CalculateTotalRevenue/RevenueByType/AvgCompletionRate/AvgDropOffRate/AvgReadingDuration |

**实施步骤**:
1. 创建 `service/stats/stats_service.go`
2. 将所有 Calculate* 方法移到 StatsService
3. Repository 只保留原始数据查询
4. 添加单元测试

---

## 宧行优先级排序原则

### 为什么这样排序？

1. **影响范围**: Writer/Reader 域影响用户核心功能，2. **风险程度**: 业务规则在 Repository 层会导致难以维护
3. **依赖关系**: Stats 域依赖其他域的数据，4. **工作量**: 从易到难，5. **ROI**: Writer/Reader 重构带来的价值最高

### 优先级判断标准

| 标准 | 权重 | 说明 |
|------|------|------|
| 影响用户核心功能 | 40% | Writer/Reader > Stats |
| 代码修改复杂度 | 30% | 简单 > 复杂 |
| 测试覆盖难度 | 20% | 已有测试 > 需要新测试 |
| 依赖关系 | 10% | 独立 > 有依赖 |

---

## 检查点

### Phase 2 检查点
- [ ] Writer 域所有测试通过
- [ ] Reader 域所有测试通过
- [ ] Bookstore 域剩余测试通过
- [ ] 代码审查通过

### Phase 3 检查点
- [ ] Social 域所有测试通过
- [ ] User 域所有测试通过
- [ ] 代码审查通过

### Phase 4 检查点
- [ ] Stats 域所有测试通过
- [ ] 性能测试通过
- [ ] 代码审查通过

---

## 验收标准

### 每个 Phase 完成标准
1. 所有测试通过（单元测试 + 集成测试）
2. 代码审查通过（无业务逻辑在 Repository 层）
3. 文档更新（Issue 状态更新）
4. PR 合并

### 最终验收标准
1. 所有 37 个方法重构完成
2. 整体测试覆盖率 > 70%
3. 无业务逻辑在 Repository 层
4. 性能无明显下降

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [Issue #010 主文档](./010-repository-layer-business-logic-permeation.md) | 问题描述和检查清单 |
| [Repository 层业务逻辑渗透分析报告](../reports/archived/2026-03-04-repository-layer-business-logic-analysis.md) | 完整问题清单 |
