# Repository / MongoDB

> 最后更新：2026-03-29

## 职责

数据访问层，通过 `MongoRepositoryFactory` 工厂模式创建所有 MongoDB Repository 实例，统一管理数据库连接和集合映射。不包含业务逻辑。

## 数据流

```
Service 层 → MongoRepositoryFactory.CreateXxxRepository() → 具体 Repository → MongoDB Driver → MongoDB
```

## 约定 & 陷阱

- **ID 类型强制规则**：所有引用 ID 在 MongoDB 中存储为 `primitive.ObjectID`，API 层转为 string，Repository 层负责转换
- **禁止重复转换**：如果字段已经是 ObjectID，不要再调 `primitive.ObjectIDFromHex()`，否则会 panic
- **工厂单例**：`MongoRepositoryFactory` 应该是应用级单例，不要重复创建
- **集合命名约定**：集合名通常为模型名的复数小写形式（如 `projects`、`documents`、`characters`）
- **事务支持**：涉及多文档写操作时使用 `session.StartTransaction()`，注意 MongoDB 副本集才能用事务
