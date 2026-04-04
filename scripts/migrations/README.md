# MongoDB 迁移脚本

> **版本**: v1.0
> **创建日期**: 2026-04-04

---

## 概述

本目录包含 MongoDB 数据库迁移脚本，用于创建和维护数据库索引。

## 迁移脚本清单

| 脚本 | 描述 | 执行时机 |
|------|------|----------|
| `001_create_indexes.js` | 创建所有集合的索引 | 首次部署 |

## 执行方法

### 在 Docker 环境中执行

```bash
# 1. 确保 MongoDB 容器正在运行
docker ps | grep mongodb

# 2. 执行迁移脚本
docker exec -i qingyu-mongodb mongosh -u admin -p "$MONGO_PASSWORD" \
  --authenticationDatabase admin Qingyu_writer < scripts/migrations/001_create_indexes.js

# 3. 验证索引创建
docker exec qingyu-mongodb mongosh -u admin -p "$MONGO_PASSWORD" \
  --authenticationDatabase admin Qingyu_writer --eval "db.documents.getIndexes()"
```

### 直接连接 MongoDB 执行

```bash
mongosh -u admin -p "$MONGO_PASSWORD" --authenticationDatabase admin Qingyu_writer < 001_create_indexes.js
```

## 索引列表

### Writer 模块
- `documents`: 2 个索引（project_id+parent_id+order_key, project_id+type+updated_at）
- `document_contents`: 1 个索引（document_id+version）
- `versions`: 1 个索引（document_id+created_at）
- `characters`: 1 个索引（project_id+name）
- `character_relations`: 1 个索引（project_id+source_id）
- `locations`: 1 个索引（project_id+name）
- `outlines`: 1 个索引（project_id+type+order_key）
- `timeline_events`: 1 个索引（project_id+timestamp）

### Bookstore 模块
- `books`: 5 个索引（status+published_at, author_id+status, category+status, stats.hot_score, stats.rating_avg）
- `chapters`: 2 个索引（book_id+order_key, book_id+status）
- `book_ratings`: 1 个唯一索引（book_id+user_id）
- `comments`: 1 个索引（target_id+target_type+created_at）

### Reader 模块
- `reading_progress`: 2 个索引（1 个唯一索引 user_id+book_id, user_id+updated_at）
- `bookmarks`: 1 个索引（user_id+book_id+chapter_id）

### User 模块
- `users`: 3 个索引（username 唯一, email 唯一, phone 稀疏）
- `sessions`: 2 个索引（user_id+expires_at, expires_at TTL）

### Notification 模块
- `notifications`: 2 个索引（user_id+is_read+created_at, user_id+type）

### Audit 模块
- `audit_logs`: 2 个索引（user_id+created_at, action+created_at）

## 注意事项

1. **首次部署**: 必须执行 `001_create_indexes.js` 创建所有索引
2. **生产环境**: 建议在低峰期执行，避免影响线上服务
3. **索引监控**: 执行后监控查询性能，确保索引被正确使用

## 回滚

如需回滚索引创建，可以删除特定索引：

```javascript
// 示例：删除 documents 集合的某个索引
db.documents.dropIndex("project_id_1_parent_id_1_order_key_1")
```

---

**最后更新**: 2026-04-04
