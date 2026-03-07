# 发布系统最小落地方案（Phase 1）

**设计日期**: 2026-03-06
**定位**: 最小可运行版本
**状态**: 新增补充方案
**对应原版**: [2026-03-06-publication-review-subsystem-design.md](./2026-03-06-publication-review-subsystem-design.md)
**后续升级依赖**:
- [2026-03-06-project-book-separation-architecture.md](./2026-03-06-project-book-separation-architecture.md)
- [2026-03-06-three-tier-version-management-design.md](./2026-03-06-three-tier-version-management-design.md)
- [2026-03-06-layered-version-management-design.md](./2026-03-06-layered-version-management-design.md)

---

## 目标

本方案不替代原版完整设计。

本方案仅解决当前 `Qingyu_backend` 的三个现实问题：

1. writer 路由中的发布功能仍使用 `MockPublishService`
2. `PublicationRepository` 尚未提供真实落地实现
3. Project 发布后缺少最小的 Book 关联与可追踪记录

本方案完成后，系统应具备：

- 真实的项目发布链路
- 真实的发布记录持久化
- Project 到 Book 的最小关联
- 为后续“审核单 / BookVersion / 订阅通知”预留扩展位

不在 Phase 1 实现：

- 审核单 `Review` 子系统
- 发布快照 `PublicationSnapshot`
- `BookVersion` / `ChapterVersion`
- 完整订阅通知流
- 全量增量版本对比

---

## 现状约束

### 已有能力

1. 已有 writer 侧 `PublishService` 业务骨架
2. 已有 `PublicationRecord` DTO 与发布状态常量
3. 已有 `Book` / `Chapter` / `ChapterContent` 读侧模型
4. 已有 MongoDB 事务能力与事件总线基础

### 当前阻塞点

1. 路由层仍注入 `MockPublishService`
2. `PublicationRepository` 仅存在接口和 mock
3. 现有发布语义是“直接推书城”，不是“提交审核后发布”
4. `Book` 模型缺少 `ProjectID` 关联字段

---

## Phase 1 目标模型

### 1. 保留现有 `PublicationRecord`

Phase 1 不新建 `models/publishing/Publication`。

继续使用现有 `dto.PublicationRecord` 作为发布记录实体，但补充字段语义：

- `Type`: 保留 `project` / `document`
- `Status`: 保留 `pending` / `published` / `failed` / `unpublished`
- `Metadata`: 增加兼容扩展字段承载 `syncMode`、`sourceProjectID`

理由：

- 现有 API、service interface、测试已围绕 `PublicationRecord` 建立
- 可以先让真实发布链路跑通，避免一次性引入完整 publishing 子域

### 2. 对 `Book` 做最小关联扩展

建议仅增加以下字段：

```go
type Book struct {
    // 现有字段省略
    ProjectID    *string    `bson:"project_id,omitempty" json:"projectId,omitempty"`
    SourceType   string     `bson:"source_type,omitempty" json:"sourceType,omitempty"` // project | manual
    LastSyncedAt *time.Time `bson:"last_synced_at,omitempty" json:"lastSyncedAt,omitempty"`
}
```

Phase 1 不引入 `SyncMode`。

理由：

- `ProjectID` 足以建立最小溯源关系
- `SourceType` 可区分后续手动创建书籍
- `LastSyncedAt` 可支持后续增量更新与运维排查

### 3. `Chapter` 暂不做复杂版本字段扩展

Phase 1 只在必要时补：

```go
ProjectID        *string `bson:"project_id,omitempty" json:"projectId,omitempty"`
ProjectChapterID *string `bson:"project_chapter_id,omitempty" json:"projectChapterId,omitempty"`
```

不在 Phase 1 引入完整 diff/version 字段。

---

## Phase 1 真实业务流程

```text
作者发布项目
  -> writer.PublishService 校验权限与项目状态
  -> PublicationRepository.Create(record)
  -> BookRepository.GetByProjectID(projectID)
  -> 若不存在则创建 Book + Chapters + ChapterContents
  -> 若存在则拒绝二次项目级发布，返回已发布状态
  -> 更新 PublicationRecord.Status=published
  -> 发布 event: project.published
```

### 设计原则

1. 只支持“首次项目发布”作为主路径
2. 暂不支持完整“重新提交审核”
3. 暂不支持增量版本发布
4. 文档级/章节级发布保持现有能力，但不升级为审核体系

---

## 实施步骤

### Step 1: 实现真实 `PublicationRepository`

新增建议：

- `repository/interfaces/writer/publication_repository.go`
- `repository/mongodb/writer/publication_repository_mongo.go`

集合建议：

- `publication_records`

最小接口：

```go
type PublicationRepository interface {
    Create(ctx context.Context, record *dto.PublicationRecord) error
    FindByID(ctx context.Context, id string) (*dto.PublicationRecord, error)
    FindByProjectID(ctx context.Context, projectID string, page, pageSize int) ([]*dto.PublicationRecord, int64, error)
    FindPublishedByProjectID(ctx context.Context, projectID string) (*dto.PublicationRecord, error)
    Update(ctx context.Context, record *dto.PublicationRecord) error
}
```

### Step 2: 为 `BookRepository` 增加 `GetByProjectID`

新增建议：

```go
GetByProjectID(ctx context.Context, projectID string) (*bookstore.Book, error)
```

并为 `books.project_id` 建唯一索引（允许 `null`）。

### Step 3: 发布服务改为真实创建书籍

保留现有 `service/writer/publish_service.go`，但替换 `BookstoreClient` 语义：

- 当前不是推外部书城服务
- 当前应直接创建本地 `Book` / `Chapter` / `ChapterContent`

更准确的做法是新增内部 adapter，例如：

- `service/writer/bookstore_local_adapter.go`

由 adapter 封装：

- 创建 Book
- 同步项目章节为 Chapter
- 写入 ChapterContent

这样不破坏 `PublishService` 对 `BookstoreClient` 的既有抽象。

### Step 4: writer 路由切换真实服务

替换：

- `router/writer/writer_router_init.go`

从 `NewMockPublishService()` 切到真实 `NewPublishService(...)`。

前提：

- `PublicationRepository` 可用
- 本地 `BookstoreClient` adapter 可用
- Book/Chapter 仓储可用

### Step 5: 事件只做最小发布

Phase 1 只保留：

- `project.published`
- `chapter.published`

不做：

- 审核事件
- 版本发布事件
- 订阅用户通知扇出

---

## 数据迁移建议

### Book 集合

为现有 `books` 集合补充可选字段：

- `project_id`
- `source_type`
- `last_synced_at`

### Chapter 集合

按需补充：

- `project_id`
- `project_chapter_id`

### 索引

建议新增：

```text
books: { project_id: 1 } unique sparse
chapters: { project_id: 1, project_chapter_id: 1 }
publication_records: { resourceId: 1, status: 1 }
publication_records: { createdBy: 1, createdAt: -1 }
```

---

## 与原版设计的关系

本方案是原版的 Phase 1，不替代原版。

升级路径：

1. 先完成本方案，替换 mock，跑通真实发布链路
2. 再引入 `Publication` 聚合替代 `PublicationRecord`
3. 再引入审核单 `Review`
4. 再引入 `BookVersion` / `ChapterVersion`
5. 最后接入订阅通知与版本对比

原版中以下能力全部延期到后续：

- `PublicationSnapshot`
- 审核流转历史
- 拒绝原因与问题列表
- 审核通过生成 `BookVersion`
- 已订阅用户更新通知
- 版本差异查看

---

## 验收标准

### 最低标准

- writer 发布路由不再使用 mock
- 发布项目后能在 `books` 集合看到关联记录
- `books.project_id` 能反查到唯一来源 Project
- `publication_records` 能记录发布成功/失败状态

### 一般标准

- 支持项目首次发布
- 支持查看项目发布状态
- 支持取消发布
- 有基础集成测试覆盖项目发布主链路

### 不作为本阶段验收项

- 审核功能
- 版本对比
- 订阅通知
- 多轮重新发布审批

---

## 风险

| 风险 | 影响 | 应对 |
|------|------|------|
| 现有 `BookstoreClient` 抽象语义偏外部系统 | 中 | 先加本地 adapter，不直接推翻接口 |
| 书籍与项目历史数据无法自动关联 | 中 | 只要求新发布数据建立关联，旧数据单独迁移 |
| 后续审核体系需要替换 `PublicationRecord` | 低 | 明确本方案为过渡层，保留兼容映射 |
| 章节内容来源复杂 | 中 | Phase 1 只覆盖项目主文档/章节发布主路径 |

---

## 推荐实施顺序

1. `PublicationRepository` 真正落地
2. `Book` 增加 `ProjectID`
3. 本地 `BookstoreClient` adapter
4. 路由切换真实 `PublishService`
5. 项目发布集成测试
6. 再立项原版完整审核与版本系统
