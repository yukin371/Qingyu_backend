# 章节内容数据迁移指南

## 概述

此迁移工具用于将现有 `chapters` 集合中的 `content` 字段迁移到新的 `chapter_contents` 集合，实现内容与元数据的分离。

## 迁移前准备

### 1. 备份数据库

**重要：在执行迁移前务必备份数据库！**

```bash
# 使用 mongodump 备份
mongodump --uri="mongodb://localhost:27017" --db=qingyu --out=./backup_$(date +%Y%m%d_%H%M%S)

# 或者在 MongoDB Compass 中导出数据
```

### 2. 检查当前数据

```bash
# 连接到 MongoDB
mongosh mongodb://localhost:27017/qingyu

# 查看章节总数
db.chapters.countDocuments()

# 查看包含 content 字段的章节数
db.chapters.countDocuments({content: {$exists: true, $ne: ""}})

# 随机查看一个章节示例
db.chapters.findOne({content: {$exists: true, $ne: ""}})
```

## 使用迁移工具

### 方式 1: 直接运行（开发环境）

```bash
cd E:\Github\Qingyu\Qingyu_backend

# 修改配置（如果需要）
# 编辑 scripts/migrate_chapter_content.go 中的 Config 结构体

# 运行迁移
go run scripts/migrate_chapter_content.go
```

### 方式 2: 编译后运行（推荐生产环境）

```bash
# 编译迁移工具
cd E:\Github\Qingyu\Qingyu_backend
go build -o ./bin/migrate_chapter_content scripts/migrate_chapter_content.go

# 运行迁移
./bin/migrate_chapter_content
```

### 方式 3: 试运行模式（Dry Run）

在执行实际迁移前，建议先运行试运行模式：

```go
// 在 scripts/migrate_chapter_content.go 中修改配置
config := Config{
    MongoURI:   "mongodb://localhost:27017",
    Database:   "qingyu",
    DryRun:     true,  // 设置为 true
    BatchSize:  100,
    SkipErrors: true,
}
```

试运行模式会：
- 统计需要迁移的章节数
- 打印迁移计划
- **不会实际修改数据库**

## 迁移过程

### 执行步骤

1. **试运行验证**
   ```bash
   # 确保 DryRun: true
   go run scripts/migrate_chapter_content.go
   ```

2. **确认无误后执行真实迁移**
   ```go
   // 修改配置 DryRun: false
   config := Config{
       DryRun: false,  // 改为 false
       ...
   }
   ```
   ```bash
   go run scripts/migrate_chapter_content.go
   ```

3. **输入确认**
   迁移工具会提示确认，输入 `yes` 继续：
   ```
   警告：此操作将修改数据库！确认继续？(yes/no): yes
   ```

### 迁移过程说明

迁移工具会执行以下操作：

1. **统计阶段**
   - 统计包含 content 字段的章节数
   - 显示需要迁移的总数

2. **迁移阶段**（批量处理）
   - 每批处理 100 个章节（可配置）
   - 为每个章节创建 ChapterContent 文档
   - 更新 Chapter 文档，添加内容引用信息
   - 跳过已迁移的章节

3. **验证阶段**
   - 统计 chapter_contents 集合文档数
   - 检查 chapters 集合中是否仍有 content 字段
   - 计算迁移成功率

4. **报告阶段**
   - 打印详细统计信息
   - 生成 JSON 格式迁移报告

## 配置说明

### Config 参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `MongoURI` | string | mongodb://localhost:27017 | MongoDB 连接 URI |
| `Database` | string | qingyu | 数据库名称 |
| `DryRun` | bool | false | 是否为试运行模式 |
| `BatchSize` | int | 100 | 批量处理大小 |
| `SkipErrors` | bool | true | 遇到错误是否继续 |

### 生产环境配置示例

```go
config := Config{
    MongoURI:   "mongodb://username:password@production-server:27017/qingyu?authSource=admin",
    Database:   "qingyu",
    DryRun:     false,
    BatchSize:  50,  // 生产环境建议减小批次
    SkipErrors: false, // 生产环境建议遇到错误停止
}
```

## 验证迁移结果

### 1. 使用迁移工具内置验证

迁移完成后，工具会自动显示：
- 迁移总数
- 成功数量
- 失败数量
- 成功率

### 2. 手动验证

```bash
# 连接到 MongoDB
mongosh mongodb://localhost:27017/qingyu

# 1. 检查 chapter_contents 集合数量
db.chapter_contents.countDocuments()

# 2. 检查章节是否已移除 content 字段
db.chapters.countDocuments({content: {$exists: true, $ne: ""}})

# 3. 验证章节元数据包含内容引用
db.chapters.findOne({}, {
  content_url: 1,
  content_size: 1,
  content_hash: 1,
  content_version: 1
})

# 4. 验证内容完整性
db.chapter_contents.findOne()

# 5. 检查章节和内容的关联
db.chapters.aggregate([
  {
    $lookup: {
      from: "chapter_contents",
      localField: "_id",
      foreignField: "chapter_id",
      as: "content"
    }
  },
  {
    $match: {
      content: { $eq: [] }
    }
  },
  {
    $count: "orphans"
  }
])
# 结果应该为 0（没有孤立章节）
```

### 3. 应用层验证

启动后端服务，测试：

```bash
# 测试获取章节列表（应该更快）
curl http://localhost:8080/api/v1/bookstore/books/{book_id}/chapters

# 测试获取章节内容
curl http://localhost:8080/api/v1/bookstore/chapters/{chapter_id}/content
```

## 回滚方案

如果迁移出现问题，可以按以下步骤回滚：

### 方案 1: 从备份恢复（推荐）

```bash
# 停止应用服务
# 停止迁移工具（如果还在运行）

# 恢复数据库
mongorestore --uri="mongodb://localhost:27017" --db=qingyu ./backup_YYYYMMDD_HHMMSS

# 重启应用服务
```

### 方案 2: 手动回滚

如果备份不方便，可以手动迁移回去：

```javascript
// 在 mongosh 中执行
use qingyu;

// 1. 从 chapter_contents 恢复 content 字段到 chapters
db.chapter_contents.find().forEach(function(content) {
  db.chapters.updateOne(
    { _id: content.chapter_id },
    { $set: { content: content.content } }
  );
});

// 2. 删除内容引用字段
db.chapters.updateMany(
  {},
  {
    $unset: {
      content_url: "",
      content_size: "",
      content_hash: "",
      content_version: ""
    }
  }
);

// 3. 删除 chapter_contents 集合
db.chapter_contents.drop();
```

## 故障排查

### 问题 1: 连接失败

```
错误: 连接 MongoDB 失败: context deadline exceeded
```

**解决方案**：
- 检查 MongoDB 服务是否运行
- 验证 MongoURI 是否正确
- 检查网络连接和防火墙

### 问题 2: 迁移后仍有 content 字段

**可能原因**：
- 部分章节迁移失败
- 批处理大小不合适

**解决方案**：
```bash
# 查看哪些章节还有 content 字段
db.chapters.find({content: {$exists: true, $ne: ""}}, {_id: 1, title: 1})

# 重新运行迁移工具（会跳过已迁移的）
go run scripts/migrate_chapter_content.go
```

### 问题 3: 内容关联丢失

**验证**：
```javascript
// 查找孤立内容（chapter_id 对应的章节不存在）
db.chapter_contents.find({
  chapter_id: {$nin: db.chapters.distinct("_id")}
})
```

**解决方案**：删除孤立内容或重新创建对应章节

## 性能建议

### 大数据集优化

如果有大量章节（>10,000）：

1. **增加批次大小**
   ```go
   BatchSize: 500  // 根据内存调整
   ```

2. **使用索引**
   ```javascript
   // 确保有这些索引
   db.chapters.createIndex({content: 1})
   db.chapter_contents.createIndex({chapter_id: 1})
   ```

3. **分时段迁移**
   - 在业务低峰期执行
   - 考虑分多次迁移（按 book_id 过滤）

## 监控和日志

迁移过程中的关键指标：

| 指标 | 说明 | 正常值 |
|------|------|--------|
| 批次处理时间 | 每批次的处理耗时 | < 5秒/批 |
| 成功率 | 成功迁移的章节比例 | > 99% |
| 剩余 content 字段 | 迁移后仍存在的数量 | 0 |
| 孤立内容 | 无章节对应的内容数 | 0 |

## 联系支持

如果遇到问题：
1. 检查本文档的故障排查部分
2. 查看迁移工具的详细日志
3. 保留完整的错误信息和堆栈跟踪
4. 联系技术支持团队

---

**注意**: 此迁移工具是 P2 重构的一部分。详细的重构文档请查看 `models/P2_CONTENT_SEPARATION_ANALYSIS.md`。
