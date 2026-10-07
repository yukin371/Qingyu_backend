# Writer Service

> 最后更新：2026-10-09

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

Phase 5 当前口径：writer 主线已完成 factory / compat / port adapter、context key 与主要业务参数 helper 收口；剩余是零散 helper 漏网点、service context 注入复核、运行时 stdout/debug 风险复扫与旧文档状态清理。

## 约定 & 陷阱

- **乐观锁版本控制**：DocumentContent 有 `version` 字段，更新时必须 `UpdateWithVersion` 验证版本号匹配，否则返回 `ErrVersionConflict`
- **全局总纲分布式锁**：创建第一个卷时用 Redis 分布式锁保护全局总纲创建（TTL 5s，最多重试3次，降级为 MongoDB upsert）
- **大纲-文档双向同步**：大纲节点(OutlineNode)和文档(Document)通过 `OutlineDocumentSyncService` 双向引用，修改大纲会同步创建/更新文档，反之亦然
- **4位错误码体系**：40xx 客户端错误 / 50xx 服务端错误，所有错误实现 `WriterError` 结构，支持字段级定位和可重试判断
- **文档内容分层存储**：内容快照根据大小选择内联或外部存储（StoreSnapshot）
- **发布是异步流程**：创建发布记录(pending) → 审核(approved/rejected) → 同步到书城(published)，通过 EventBus 解耦
- **历史 compat 已退场**：`service/writer/factory.go`、`service/writer/_migration/compat.go`、`service/writer/impl/` 与 `service/interfaces/writer/` 已在 2026-04-30 删除。删除前备份点为 git tag `backup/writer-compat-predelete-20260430`，如需恢复旧适配器可从该 tag 查找。
- **stdout/debug 风险**：树构建与文档链路中的已知临时 stdout 输出已在 2026-05-05 清理；后续新增排障输出必须使用 logger，不能直接写 stdout 或 dump 请求/文档结构。
- **项目 ID 完整性**：角色、地点、大纲及关系资产的服务入口在查询、权限比较或写入模型前必须严格解析 `projectID`；禁止把解析失败的全零 `ObjectID` 用于权限判断、持久化或事件发布。当前角色、地点和大纲核心读写入口已完成该边界收口，其他资产服务按 Phase 3/T3 逐模块迁移。
- **关系时序兼容边界**：关系时序创建/读取已校验项目 ID；更新/删除的事件索引与关系状态持久化仍受 `TODO(T3-072-character-relation-state-owner)` 约束，未形成唯一 repository port 前不得把兼容 stub 当作真实写入能力。
- **大纲-文档同步错误**：全局大纲查找、文档创建同步和双向引用回写必须传播项目 ID 解析及 repository 更新错误；禁止用全零 `ProjectID` 继续同步，也禁止把引用回写失败当作成功。
- **Story Harness 上下文完整性**：章节上下文入口必须先解析项目 ID，再查询投影、角色和关系；待处理变更计数失败必须向上返回，snapshot 转换只能使用已验证的项目 ObjectID，不能静默生成全零归属。
- **导入/导出身份完整性**：项目导入路由 `POST /api/v1/writer/projects/import` 传入的用户 ID 必须在 ZIP 解析和项目构造前通过 `repository.ParseID` 校验；解析失败返回校验错误，禁止把全零 `AuthorID` 写入项目。导出任务及导入服务的用户、项目和任务 ID 均不得忽略解析错误。
- **大纲查询完整性**：大纲 repository 的项目查询、根节点查询和计数入口必须显式处理项目 ID 解析错误，禁止用全零 `ObjectID` 查询；创建子节点时同级计数失败必须返回错误，不能继续以默认顺序写入。
- **发布状态读取完整性**：`GET /api/v1/writer/projects/:id/publication-status` 的发布记录、书城统计和文档统计读取失败必须向上返回错误，不能降级为“未发布”或零统计。发布执行器的异步更新与审计回写遵循 `TODO(T3-078-publish-write-error-owner)`，未定义补偿契约前不得把忽略错误当作成功。
- **层级同步错误传播**：大纲层级变为不映射文档类型时，清除文档 `outline_node_id` 的更新失败必须返回错误，不能以同步成功结束并保留脏关联。
- **发布创建前置检查**：创建项目发布记录前，既有已发布记录查询失败必须返回错误并停止创建；仓储不可用时不能继续写入 pending 记录。
- **章节发布前置检查**：创建章节发布记录前，章节自身和所属项目的既有发布记录查询失败必须返回错误并停止创建，不能继续写入无法关联书城项目的 pending 记录。
- **批注统计错误传播**：批注统计由 `repository/mongodb/writer/MongoCommentRepository.GetStats` 负责聚合；总数、已解决数、聚合和 cursor 解码错误必须向 `CommentServiceImpl`/API 传播，异常 BSON 类型不得通过强制断言触发 panic。真实 Mongo 聚合联调与错误注入测试归入 Phase 4。
