# 数据库设计文档

> 最后整理: 2026-05-22  
> 当前状态: `summary-draft`

本文档是 `design/database/` 的早期汇总稿，适合快速回看当时如何把数据库选型、架构设计和优化策略组织在同一专题下；它不是当前目录的标准入口，也不等同于当前数据库专题的 owner 文档。

## Page Role

- 这里负责：数据库历史专题总览、早期分组方式、旧版设计原则与材料导航。
- 不负责：当前数据库 owner、当前迁移策略、现行 Repository / Model 分层规范。

## Recommended Read Path

1. [README.md](./README.md)
2. [../../database/README.md](../../database/README.md)
3. [../../migration/README.md](../../migration/README.md)
4. [../../standards/layer-repository.md](../../standards/layer-repository.md)
5. [../../standards/layer-models.md](../../standards/layer-models.md)

## Boundary

- 如果你要找“design/database 目录现在怎么读”，优先看 [README.md](./README.md)。
- 如果你要找“当前数据库专题 owner 和迁移样例”，优先看 [../../database/README.md](../../database/README.md)。
- 如果你要找“文档迁移与兼容口径”，优先看 [../../migration/README.md](../../migration/README.md)。
- 如果你要找“当前 Repository / Models 分层规范”，优先看 [../../standards/layer-repository.md](../../standards/layer-repository.md) 和 [../../standards/layer-models.md](../../standards/layer-models.md)。

## 📁 文档目录

### 技术选型

- [数据库技术选型分析](./数据库技术选型分析.md) - MongoDB vs MySQL 技术选型分析和决策

### 架构设计

- [数据库设计说明书](./数据库设计说明书.md) - 数据库整体架构、表结构、索引设计
- [Repository层设计说明书](./Repository层设计说明书.md) - 数据访问层的设计规范和实现指南

### 性能优化

- [MongoDB应用层优化设计](./MongoDB应用层优化设计.md) - MongoDB的查询优化、索引策略、性能调优

## 🎯 设计原则

### 数据一致性

通过事务机制和约束保证数据的一致性和完整性。

### 性能优化

- 合理的索引设计
- 查询优化和缓存策略
- 读写分离和分片策略

### 可扩展性

- 支持水平扩展
- 灵活的Schema设计
- 版本迁移机制

## 📊 数据模型

### 核心实体

- **用户 (User)** - 用户基本信息、认证信息
- **项目 (Project)** - 写作项目、作品信息
- **文档 (Document)** - 文档内容、章节信息
- **书籍 (Book)** - 书城书籍、详情信息
- **分类 (Category)** - 书籍分类、标签系统

### 关系设计

```tree
User (1) ─────── (N) Project
             │
             └─── (N) Book

Project (1) ───── (N) Document

Book (1) ────── (N) Chapter
         │
         └────── (N) Rating
```

## 🔗 相关文档

- [核心功能设计](../core/) - 核心功能设计文档
- [Repository层重构](../重构规划/Repository层重构设计.md) - Repository层重构规划

## 📝 更新日志

- 2025-09-30: 更新数据库设计文档目录，整理现有设计文档
