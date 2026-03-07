# 2026-03-06 设计文档审查摘要

**审查日期**: 2026-03-06
**审查范围**:
- [2026-03-06-publication-review-subsystem-design.md](./2026-03-06-publication-review-subsystem-design.md)
- [2026-03-06-publication-review-subsystem-mvp.md](./2026-03-06-publication-review-subsystem-mvp.md)
- [2026-03-06-project-book-separation-architecture.md](./2026-03-06-project-book-separation-architecture.md)
- [2026-03-06-three-tier-version-management-design.md](./2026-03-06-three-tier-version-management-design.md)
- [2026-03-06-layered-version-management-design.md](./2026-03-06-layered-version-management-design.md)
- [2026-03-06-p0-implementation-roadmap.md](./2026-03-06-p0-implementation-roadmap.md)

---

## 总结结论

| 文档 | 结论 | 说明 |
|------|------|------|
| publication-review-subsystem-design | 可作为目标架构 | 方向正确，但不能直接落地 |
| publication-review-subsystem-mvp | 可直接作为 Phase 1 | 与当前代码更匹配 |
| project-book-separation-architecture | 必做前置改造 | 是发布系统落地基础 |
| three-tier-version-management-design | 可作为长期方案 | 过重，不宜先做 |
| layered-version-management-design | 可作为优化备选 | 比三层版更激进，优先级更低 |
| p0-implementation-roadmap | 需要重排顺序 | 工期偏乐观，依赖关系需调整 |

---

## 逐项审查

### 1. publication-review-subsystem-design

**结论**: 架构方向正确，当前不可直接落地。

**优点**:

- 把发布、审核、版本拆成独立子域，边界清晰
- 引入 `PublicationSnapshot` 能避免审核期内容漂移
- 为后续版本对比、通知、重审提供了统一状态机

**问题**:

- 当前代码里 writer 发布路由仍在用 mock 服务
- 现有系统只有 `PublicationRecord`，没有 `Publication` 聚合
- `models/social/review.go` 是书评，不是审核单模型
- 当前没有 `BookVersion` / `ChapterVersion` 实体与仓储

**审查意见**:

- 保留原版作为 Phase 2/3 蓝图
- 不建议直接按原版全量开工
- 应先实施 MVP 版替换 mock

### 2. publication-review-subsystem-mvp

**结论**: 可直接进入实现排期。

**原因**:

- 复用现有 `PublishService` 和 `PublicationRecord`
- 只补最小 Book-Project 关联
- 不要求同时完成审核与版本体系
- 能在较小改动下把发布功能从 mock 变为真实链路

**前置依赖**:

- `PublicationRepository` MongoDB 实现
- `BookRepository.GetByProjectID`
- `Book.project_id` 索引与字段迁移

### 3. project-book-separation-architecture

**结论**: 高优先级，且是发布系统的前置条件。

**优点**:

- 识别了当前 `Project` 与 `Book` 没有稳定关联的真实问题
- `Book.ProjectID` 是最关键的补齐字段
- 为后续同步、审计、版本溯源提供基础

**问题**:

- 文档里的 `SyncMode` 设计对当前阶段过重
- 直接引入全量自动同步逻辑会扩大改动范围

**审查意见**:

- Phase 1 只落地 `ProjectID + SourceType + LastSyncedAt`
- `SyncMode` 延后到版本系统或审核系统稳定后再做

### 4. three-tier-version-management-design

**结论**: 长期可行，短期不建议实施。

**优点**:

- 与现有 writer/bookstore 内容分层模型有一定兼容性
- 快照引用思路比全量复制更稳健
- 可以支持文档级和书籍级版本对比

**问题**:

- 当前对外版本接口尚未完成
- 现有发布功能尚未真实落地
- 在发布审核链路未稳定前引入双域版本体系，会显著提高复杂度

**审查意见**:

- 作为完整版本子系统的长期设计保留
- 不进入最近一期实现范围

### 5. layered-version-management-design

**结论**: 技术上可行，但当前优先级低于 three-tier 版本方案。

**优点**:

- 存储效率高
- 适合大章节、小修改频繁的场景
- 对超大内容的长期成本更友好

**问题**:

- diff 重建链路复杂
- 调试成本高
- 当前业务并未证明已需要这一级别优化
- 如果先做，会把基础发布问题掩盖成底层存储工程问题

**审查意见**:

- 只作为后续性能优化方案保留
- 不建议在 P0 主路径优先实施

### 6. p0-implementation-roadmap

**结论**: 需要重排顺序，不能原样执行。

**原问题**:

- 把“发布-审核-订阅系统”和“版本管理”放在同一阶段并行推进，低估了依赖成本
- 默认项目发布主链路已经存在，但当前实际上仍是 mock

**建议顺序**:

1. 类型统一
2. 事务/仓储基础补齐
3. Project-Book 最小关联
4. 发布 MVP 真正落地
5. 审核子系统
6. BookVersion / ChapterVersion
7. 分层 diff 优化

---

## 推荐执行顺序

### Phase 1

- `project-book-separation-architecture` 的最小字段版
- `publication-review-subsystem-mvp`

### Phase 2

- `publication-review-subsystem-design` 中的审核流
- `PublicationSnapshot`

### Phase 3

- `three-tier-version-management-design`

### Phase 4

- `layered-version-management-design` 中的 diff / 压缩优化

---

## 最终判断

### 可以直接推进

- `publication-review-subsystem-mvp`
- `project-book-separation-architecture` 的最小字段子集

### 应作为后续任务保留

- `publication-review-subsystem-design`
- `three-tier-version-management-design`
- `layered-version-management-design`

### 需要更新的管理认知

- `2026-03-06-p0-implementation-roadmap.md` 不应再把完整审核系统和完整版本系统视为同一阶段可闭环任务
- 当前 P0 首要目标应从“完整平台能力”调整为“真实发布主链路替换 mock”
