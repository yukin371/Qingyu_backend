# Writer Service

> 最后更新：2026-04-30

## 职责

写作工作台的核心业务层，管理项目/文档/角色/地点/时间线的完整创作生命周期，以及发布到书城的流程。不负责搜索、支付、用户认证。

## 数据流

```
Router/API Handler → {子模块Service} → Repository → MongoDB
                                                        ↓
                                                 EventBus (异步事件)
                                                        ↓
                                               BookstoreClient (发布到书城)
```

当前 writer 路由直接注入各子模块 service，不再通过历史 `WriterServiceFactory` / `_migration.WriterServiceAdapter` 做二次组合：
- project/：项目生命周期
- document/：文档、版本、批量操作与编辑链路
- Character/Location/Timeline：内容资产
- Comment/Template/Dashboard：协作与工作台辅助能力
- Publish/Export：发布与导出

## 约定 & 陷阱

- **乐观锁版本控制**：DocumentContent 有 `version` 字段，更新时必须 `UpdateWithVersion` 验证版本号匹配，否则返回 `ErrVersionConflict`
- **全局总纲分布式锁**：创建第一个卷时用 Redis 分布式锁保护全局总纲创建（TTL 5s，最多重试3次，降级为 MongoDB upsert）
- **大纲-文档双向同步**：大纲节点(OutlineNode)和文档(Document)通过 `OutlineDocumentSyncService` 双向引用，修改大纲会同步创建/更新文档，反之亦然
- **4位错误码体系**：40xx 客户端错误 / 50xx 服务端错误，所有错误实现 `WriterError` 结构，支持字段级定位和可重试判断
- **文档内容分层存储**：内容快照根据大小选择内联或外部存储（StoreSnapshot）
- **发布是异步流程**：创建发布记录(pending) → 审核(approved/rejected) → 同步到书城(published)，通过 EventBus 解耦
- **历史 compat 已退场**：`service/writer/factory.go`、`service/writer/_migration/compat.go`、`service/writer/impl/` 与 `service/interfaces/writer/` 已在 2026-04-30 删除。删除前备份点为 git tag `backup/writer-compat-predelete-20260430`，如需恢复旧适配器可从该 tag 查找。
