# Writer Service

> 最后更新：2026-03-29

## 职责

写作工作台的核心业务层，管理项目/文档/角色/地点/时间线的完整创作生命周期，以及发布到书城的流程。不负责搜索、支付、用户认证。

## 数据流

```
API Handler → WriterServiceFactory → impl/{Port} → {子模块Service} → Repository → MongoDB
                                                        ↓
                                                 EventBus (异步事件)
                                                        ↓
                                               BookstoreClient (发布到书城)
```

统一入口是 `WriterService` 接口（factory.go），通过 `WriterServiceFactory` 组合 5 个 Port 实现：
- ProjectManagement → project/ 子模块
- DocumentManagement → document/ 子模块
- ContentManagement → Character/Location/Timeline
- Collaboration → Comment
- PublishExport → Publish/Export

## 约定 & 陷阱

- **乐观锁版本控制**：DocumentContent 有 `version` 字段，更新时必须 `UpdateWithVersion` 验证版本号匹配，否则返回 `ErrVersionConflict`
- **全局总纲分布式锁**：创建第一个卷时用 Redis 分布式锁保护全局总纲创建（TTL 5s，最多重试3次，降级为 MongoDB upsert）
- **大纲-文档双向同步**：大纲节点(OutlineNode)和文档(Document)通过 `OutlineDocumentSyncService` 双向引用，修改大纲会同步创建/更新文档，反之亦然
- **4位错误码体系**：40xx 客户端错误 / 50xx 服务端错误，所有错误实现 `WriterError` 结构，支持字段级定位和可重试判断
- **文档内容分层存储**：内容快照根据大小选择内联或外部存储（StoreSnapshot）
- **发布是异步流程**：创建发布记录(pending) → 审核(approved/rejected) → 同步到书城(published)，通过 EventBus 解耦
