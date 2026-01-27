# MongoDB Profiler 配置指南

## 概述

MongoDB Profiler是监控数据库性能的重要工具，可以记录慢查询信息，帮助定位性能瓶颈。

## 配置说明

### Profiling级别

- **Level 0**: 关闭profiling（不记录任何查询）
- **Level 1**: 仅记录慢查询（默认，推荐用于生产环境）
- **Level 2**: 记录所有查询（仅用于调试，不推荐用于生产环境）

### 配置参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `profiling_level` | 1 | Profiling级别（0-2） |
| `slow_ms` | 100 | 慢查询阈值（毫秒） |
| `profiler_size_mb` | 100 | Profiler存储上限（MB） |

## 使用方法

### 1. 使用JavaScript脚本配置（推荐）

```bash
# 方法1: 使用mongosh直接执行脚本
mongosh qingyu --file scripts/enable_profiling.js

# 方法2: 先连接mongosh，再执行脚本
mongosh qingyu
> load("scripts/enable_profiling.js")
```

### 2. 手动配置

```javascript
// 连接到数据库
mongosh qingyu

// 设置profiling级别和慢查询阈值
db.setProfilingLevel(1, { slowms: 100 })

// 设置profiler collection大小上限
db.system.profile.convertToCapped({
    size: 104857600  // 100MB
})

// 验证配置
db.getProfilingStatus()
```

### 3. 通过Go代码配置

在应用启动时自动配置：

```go
import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func configureProfiler(ctx context.Context, client *mongo.Client, dbName string) error {
    db := client.Database(dbName)

    // 设置profiling级别
    command := bson.D{
        {Key: "profile", Value: 1},
        {Key: "slowms", Value: 100},
    }
    result := db.RunCommand(ctx, command)
    if result.Err() != nil {
        return result.Err()
    }

    // 设置collection大小
    profileCollection := db.Collection("system.profile")
    // ... convertToCapped逻辑

    return nil
}
```

## 配置文件示例

在 `config/config.yaml` 中配置：

```yaml
database:
  primary:
    type: mongodb
    mongodb:
      uri: mongodb://localhost:27017
      database: qingyu
      profiling_level: 1   # 0=off, 1=slow only, 2=all
      slow_ms: 100         # 慢查询阈值（毫秒）
      profiler_size_mb: 100 # Profiler存储上限（MB）
```

## 环境变量配置

可以通过环境变量覆盖配置：

```bash
export MONGODB_PROFILING_LEVEL=1
export MONGODB_SLOW_MS=100
export MONGODB_PROFILER_SIZE_MB=100
```

## 查询Profiler数据

### 查看最慢的查询

```javascript
// 最慢的5个查询
db.system.profile.find().sort({millis: -1}).limit(5).pretty()

// 慢于200ms的查询
db.system.profile.find({millis: {$gt: 200}}).sort({millis: -1}).pretty()
```

### 按Collection过滤

```javascript
// 特定collection的慢查询
db.system.profile.find({ns: "qingyu.users"}).sort({millis: -1}).limit(5)

// 查看涉及特定操作的查询
db.system.profile.find({op: "query"}).sort({millis: -1}).limit(5)
```

### 按时间范围过滤

```javascript
// 最近1小时的慢查询
var oneHourAgo = new Date(Date.now() - 60*60*1000)
db.system.profile.find({ts: {$gt: oneHourAgo}}).sort({millis: -1})
```

### 常用聚合查询

```javascript
// 统计每个collection的慢查询数量
db.system.profile.aggregate([
    {$group: {_id: "$ns", count: {$sum: 1}, avgMs: {$avg: "$millis"}}},
    {$sort: {count: -1}}
])

// 查找最慢的操作类型
db.system.profile.aggregate([
    {$group: {_id: "$op", count: {$sum: 1}, maxMs: {$max: "$millis"}}},
    {$sort: {maxMs: -1}}
])
```

## Profiler数据结构

```javascript
{
    "op": "query",              // 操作类型: query, insert, update, remove, command
    "ns": "qingyu.users",       // 命名空间 (数据库.集合)
    "query": {...},             // 查询条件
    "millis": 145,              // 执行时间（毫秒）
    "ts": ISODate("..."),       // 时间戳
    "client": "127.0.0.1",      // 客户端地址
    "allUsers": [...],          // 执行用户信息
    "nreturned": 10,            // 返回文档数量
    "docsExamined": 1000,       // 检查的文档数量
    "keysExamined": 500,        // 检查的索引键数量
    "execStats": {...}          // 执行计划统计
}
```

## 性能优化建议

1. **生产环境建议使用Level 1**：仅记录慢查询，减少性能影响
2. **合理设置slow_ms阈值**：太低会产生大量数据，太高可能漏掉问题查询
3. **定期清理profiler数据**：虽然是capped collection，但仍需监控
4. **结合索引优化**：分析慢查询后，添加相应索引

## 故障排查

### Profiler未记录数据

```javascript
// 检查profiling状态
db.getProfilingStatus()

// 检查system.profile集合
db.system.profile.count()

// 检查是否有写入权限
db.system.profile.insertOne({test: 1})
```

### Profiler性能影响过大

1. 降低profiling级别（从2改为1）
2. 提高slow_ms阈值
3. 减小profiler collection大小

## 相关文档

- [MongoDB Database Profiler官方文档](https://www.mongodb.com/docs/manual/administration/analyze-mongodb-performance/)
- [慢查询优化指南](../../docs/guides/performance/slow-query-optimization.md)
