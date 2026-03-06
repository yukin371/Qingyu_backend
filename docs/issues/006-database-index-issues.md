# Issue #006: 数据库索引问题

**优先级**: 高 (P0)
**类型**: 性能问题
**状态**: 待处理
**创建日期**: 2026-03-05
**来源报告**: [后端综合审计报告](../reports/archived/backend-comprehensive-audit-summary-2026-01-26.md)、[后端数据库分析](../reports/archived/backend-database-analysis-2026-01-26.md)

---

## 问题描述

规范中定义的数据库索引很多未创建，导致查询性能问题。

### 具体问题

#### 1. 索引实现差距 🔴 P0

**问题**: 根据 `docs/database/indexes.yaml` 定义的索引规范，实际数据库中缺少大量索引。

**影响**:
- 慢查询频发
- 数据库 CPU 使用率高
- API 响应时间过长
- 大量数据时性能急剧下降

#### 2. 慢查询监控缺失 🟡 P1

**问题**: 没有慢查询监控和告警机制。

**影响**:
- 无法及时发现性能问题
- 问题定位困难
- 性能优化缺乏数据支撑

---

## 解决方案

### 1. 创建缺失的索引

```go
// cmd/migrate/create_indexes.go
package main

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func createIndexes(ctx context.Context, db *mongo.Database) error {
    // Books 索引
    booksCol := db.Collection("books")
    booksModels := []mongo.IndexModel{
        {Keys: bson.D{{Key: "author_id", Value: 1}}},
        {Keys: bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}}},
        {Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}}},
        {Keys: bson.D{{Key: "categories", Value: 1}}},
        {Keys: bson.D{{Key: "statistics.view_count", Value: -1}}},
    }
    _, err := booksCol.Indexes().CreateMany(ctx, booksModels)
    if err != nil {
        return err
    }

    // Chapters 索引
    chaptersCol := db.Collection("chapters")
    chaptersModels := []mongo.IndexModel{
        {Keys: bson.D{{Key: "book_id", Value: 1}}},
        {Keys: bson.D{{Key: "book_id", Value: 1}, {Key: "chapter_number", Value: 1}}},
        {Keys: bson.D{{Key: "published_at", Value: -1}}},
    }
    _, err = chaptersCol.Indexes().CreateMany(ctx, chaptersModels)
    if err != nil {
        return err
    }

    // ... 其他集合的索引

    return nil
}
```

### 2. 索引验证工具

```go
// cmd/tools/verify_indexes.go
package main

type IndexVerification struct {
    Collection  string
    Expected    []string
    Actual      []string
    Missing     []string
    Extra       []string
}

func verifyIndexes(ctx context.Context, db *mongo.Database) []IndexVerification {
    // 读取规范定义
    specIndexes := loadIndexSpec("docs/database/indexes.yaml")

    // 验证每个集合
    var results []IndexVerification
    for _, colName := range specIndexes.Collections {
        expected := specIndexes.Indexes[colName]
        actual := getActualIndexes(ctx, db, colName)

        results = append(results, IndexVerification{
            Collection: colName,
            Expected:   expected,
            Actual:     actual,
            Missing:    difference(expected, actual),
            Extra:      difference(actual, expected),
        })
    }

    return results
}
```

### 3. 慢查询监控

```go
// pkg/monitor/query_logger.go
package monitor

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/event"
)

type QueryLogger struct {
    slowThreshold time.Duration
}

func (ql *QueryLogger) Monitor() *event.CommandMonitor {
    return &event.CommandMonitor{
        Started: func(ctx context.Context, event *event.CommandStartedEvent) {
            ctx = context.WithValue(ctx, "startTime", time.Now())
        },
        Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
            startTime := ctx.Value("startTime").(time.Time)
            duration := time.Since(startTime)

            if duration > ql.slowThreshold {
                logSlowQuery(event, duration)
            }
        },
    }
}

func logSlowQuery(event *event.CommandSucceededEvent, duration time.Duration) {
    log.Warn(
        "Slow query detected",
        "collection", event.CollectionName,
        "command", event.CommandName,
        "duration", duration,
        "query", event.Request,
    )
}
```

---

## 实施计划

### Phase 1: 索引现状评估（1 天）

1. 读取索引规范定义
2. 导出现有索引
3. 对比分析，生成缺失索引清单

### Phase 2: 创建缺失索引（1-2 天）

1. 在测试环境创建索引
2. 验证索引效果
3. 在生产环境创建（选择低峰期）

**注意事项**:
- 大集合创建索引需要较长时间
- 使用 `background: true` 选项避免阻塞
- 监控创建过程中的性能影响

### Phase 3: 部署监控工具（2-3 天）

1. 集成查询日志监控
2. 设置慢查询阈值（建议 100ms）
3. 配置告警规则

### Phase 4: 持续优化（持续）

1. 定期审查慢查询日志
2. 优化慢查询
3. 调整索引策略

---

## 索引创建优先级

### 高优先级（立即创建）

| 集合 | 索引 | 理由 |
|------|------|------|
| books | author_id | 按作者查询图书列表 |
| books | status + created_at | 按状态和时间过滤 |
| chapters | book_id + chapter_number | 章节顺序查询 |
| users | username | 用户登录 |
| comments | target_id | 评论列表查询 |

### 中优先级（近期创建）

| 集合 | 索引 | 理由 |
|------|------|------|
| books | categories | 分类筛选 |
| reading_progress | user_id + book_id | 阅读进度查询 |
| notifications | user_id + created_at | 通知列表 |

### 低优先级（按需创建）

- 文本搜索索引
- 统计数据索引
- 冷数据索引

---

## 检查清单

### 评估阶段
- [ ] 读取索引规范文档
- [ ] 导出现有索引
- [ ] 生成对比报告

### 实施阶段
- [ ] 在测试环境验证
- [ ] 制定生产环境创建计划
- [ ] 执行索引创建
- [ ] 验证索引效果

### 监控阶段
- [ ] 部署慢查询监控
- [ ] 设置告警规则
- [ ] 定期审查慢查询

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [后端数据库分析](../reports/archived/backend-database-analysis-2026-01-26.md) | 数据库问题详细分析 |
| [ID 类型诊断报告](../reports/archived/2026-02-01-id-type-diagnosis-report.md) | ID 类型与索引关系 |

---

## 相关Issue

### 相关Issue（联合处理）
- [#001: 统一模型层 ID 字段类型](./001-unify-id-type-in-models.md) - ID类型统一后，索引策略需要相应调整
- [#005: API 标准化问题](./005-api-standardization-issues.md) - API性能优化需要索引支持
- [#011: 前后端数据类型不一致](./011-frontend-backend-data-type-inconsistency.md) - 数据类型统一后，查询模式可能改变，需要索引优化

### 关联问题
- 规范中定义的数据库索引很多未创建
- 慢查询监控缺失
- 查询性能问题
