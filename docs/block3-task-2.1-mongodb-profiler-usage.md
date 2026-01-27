# Task 2.1: MongoDB Profiler 配置说明

## 概述

本任务为 Qingyu Backend 添加了 MongoDB Profiler 慢查询监控配置，用于记录和分析超过阈值的数据库查询操作。

## 实现内容

### 1. 配置扩展 (config/database.go)

在 `MongoDBConfig` 结构体中添加了三个新字段：

```go
type MongoDBConfig struct {
    // ... 原有字段

    // Profiling配置
    ProfilingLevel int   `yaml:"profiling_level" json:"profiling_level" mapstructure:"profiling_level"`
    SlowMS         int64 `yaml:"slow_ms" json:"slow_ms" mapstructure:"slow_ms"`
    ProfilerSizeMB int64 `yaml:"profiler_size_mb" json:"profiler_size_mb" mapstructure:"profiler_size_mb"`
}
```

**字段说明：**
- `ProfilingLevel`: Profiling级别 (0=关闭, 1=仅慢查询, 2=全部查询)
- `SlowMS`: 慢查询阈值（毫秒）
- `ProfilerSizeMB`: Profiler集合大小限制（MB）

### 2. 默认配置

默认值在 `getDefaultDatabaseConfig()` 中设置：

```go
MongoDB: &MongoDBConfig{
    // ... 原有配置
    ProfilingLevel: 1,
    SlowMS:         100,
    ProfilerSizeMB: 100,
},
```

**默认配置说明：**
- Profiling级别：1（仅记录慢查询）
- 慢查询阈值：100毫秒
- 存储限制：100MB（循环覆盖）

### 3. 配置验证

在 `MongoDBConfig.Validate()` 方法中添加了验证逻辑：

```go
if c.ProfilingLevel < 0 || c.ProfilingLevel > 2 {
    return fmt.Errorf("ProfilingLevel必须在0-2之间")
}
if c.SlowMS < 0 {
    return fmt.Errorf("SlowMS必须非负")
}
if c.ProfilerSizeMB < 1 {
    return fmt.Errorf("ProfilerSizeMB必须至少为1MB")
}
```

### 4. JavaScript配置脚本

创建了 `scripts/db/enable_profiling.js` 脚本，用于在MongoDB中启用Profiler：

```javascript
// 设置profiling级别为1 (仅记录慢查询)，阈值为100ms
db.setProfilingLevel(1, { slowms: 100 });

// 限制profiler集合大小为100MB
db.system.profile.convertToCapped({
    size: 104857600  // 100MB
});
```

### 5. 测试覆盖

创建了 `config/database_test.go`，包含以下测试：

- `TestMongoDBProfilingDefaults`: 验证默认配置值
- `TestMongoDBProfilingValidation`: 验证配置验证逻辑
  - 有效配置（级别0、1、2）
  - 无效级别（负数、超过2）
  - 无效SlowMS（负数）
  - 无效ProfilerSizeMB（零、负数）
- `TestMongoDBConfigProfilingDefaults`: 验证配置默认值填充
- `TestMongoDBConfigToRepositoryConfig`: 验证配置转换

## 使用方法

### 方法1: 使用YAML配置文件

在 `config/config.yaml` 中添加：

```yaml
database:
  type: mongodb
  primary:
    type: mongodb
    mongodb:
      uri: mongodb://localhost:27017
      database: qingyu
      max_pool_size: 100
      min_pool_size: 5
      connect_timeout: 10s
      server_timeout: 30s

      # Profiling配置
      profiling_level: 1  # 0=关闭, 1=慢查询, 2=全部
      slow_ms: 100        # 慢查询阈值(ms)
      profiler_size_mb: 100  # 存储大小限制(MB)
```

### 方法2: 使用环境变量

```bash
export MONGODB_PROFILING_LEVEL=1
export MONGODB_SLOW_MS=100
export MONGODB_PROFILER_SIZE_MB=100
```

### 方法3: 使用JavaScript脚本手动配置

```bash
# 连接到MongoDB
mongosh mongodb://localhost:27017/qingyu

# 执行配置脚本
load scripts/db/enable_profiling.js
```

## 查询慢查询

### 查看最近的慢查询

```javascript
// 在MongoDB shell中执行
db.system.profile.find().sort({ts: -1}).limit(10)
```

### 查找特定集合的慢查询

```javascript
db.system.profile.find({ns: "qingyu.books"}).sort({ts: -1})
```

### 查找超过特定时间的查询

```javascript
db.system.profile.find({millis: {$gt: 200}}).sort({millis: -1})
```

### 统计慢查询数量

```javascript
db.system.profile.count()
```

## Profiling级别说明

| 级别 | 说明 | 性能影响 |
|------|------|----------|
| 0 | 关闭Profiling | 无 |
| 1 | 仅记录超过slowms阈值的查询 | 较小 |
| 2 | 记录所有查询 | 较大 |

**生产环境建议：** 使用级别1，设置合理的slowms阈值（如100ms）

## 性能建议

1. **存储限制**：Profiler使用固定大小集合，旧数据会自动覆盖
   - 建议大小：100MB - 1GB
   - 太小会丢失历史数据，太大会占用磁盘空间

2. **慢查询阈值**：
   - 开发环境：50-100ms
   - 生产环境：100-200ms
   - 根据实际业务调整

3. **性能影响**：
   - Level 1 的性能影响很小（<1%）
   - Level 2 会记录所有查询，性能影响较大
   - 建议生产环境使用 Level 1

## 监控和告警

### 持续监控

建议定期检查慢查询数量和模式：

```javascript
// 最近1小时的慢查询
var oneHourAgo = new Date(Date.now() - 60*60*1000);
db.system.profile.find({ts: {$gte: oneHourAgo}}).count()
```

### 告警规则示例

可以基于以下指标设置告警：
- 慢查询数量突增
- 单个查询耗时超过500ms
- 特定集合的慢查询比例过高

## 故障排除

### 检查Profiling状态

```javascript
db.getProfilingStatus()
```

### 禁用Profiling

```javascript
db.setProfilingLevel(0)
```

### 清空Profiler数据

```javascript
db.system.profile.drop()
```

然后重新配置：

```javascript
db.setProfilingLevel(1, { slowms: 100 })
db.system.profile.convertToCapped({ size: 104857600 })
```

## 测试验证

运行测试验证配置：

```bash
cd Qingyu_backend-block3-optimization
go test ./config/... -v -run TestMongoDBProfiling
```

预期输出：

```
=== RUN   TestMongoDBProfilingDefaults
--- PASS: TestMongoDBProfilingDefaults (0.00s)
=== RUN   TestMongoDBProfilingValidation
--- PASS: TestMongoDBProfilingValidation (0.00s)
=== RUN   TestMongoDBConfigProfilingDefaults
--- PASS: TestMongoDBConfigProfilingDefaults (0.00s)
PASS
```

## 后续步骤

本任务为 Task 2.2（创建慢查询分析工具）提供了配置基础。接下来需要：

1. 创建Go工具查询和分析system.profile集合
2. 实现慢查询报告生成
3. 集成到监控系统

## 参考资料

- [MongoDB Database Profiler](https://www.mongodb.com/docs/manual/administration/analyze-mongodb-performance/#database-profiler)
- [system.profile集合](https://www.mongodb.com/docs/manual/reference/system-collections/#system.profile)

---

**实现者**: 猫娘助手Kore
**任务ID**: Task 2.1
**提交哈希**: 09d0169
**日期**: 2026-01-27
