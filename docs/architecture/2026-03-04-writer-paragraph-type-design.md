# Writer 段落类型设计（2026-03-04）

## 1. 结论

- 引入 `Paragraph` 领域类型（Value Object），用于承载编辑器分段语义。
- 不新增 `paragraph` 独立集合/表，持久化仍使用 `document_content`。
- 继续采用批量提交接口（按章节提交多个段落）。

## 2. 为什么不新增实体

- 当前项目处于内部快速迭代阶段，优先控制复杂度与交付速度。
- 现有 `document_content` 已可表达段落：`document_id + paragraph_order + content_type=tiptap`。
- 新增实体会引入额外迁移、事务边界、索引和兼容成本，当前收益不足。

## 3. 模型边界

`Paragraph`（领域层）：
- `id`
- `documentId`
- `order`
- `content`
- `contentType`
- `version`
- `updatedAt`

`DocumentContent`（持久化层）：
- 继续作为 Mongo 文档模型，不改变底层集合。

## 4. 服务层职责

- Service 负责：
- 段落顺序归一化（`order<=0` 自动重排）。
- 内容类型默认值（空值默认 `tiptap`）。
- DTO 与领域模型映射。
- Repository 负责：
- 基础 CRUD 与查询，不承载段落业务规则。

## 5. 后续升级触发条件

只有出现以下需求时，才升级为独立 `paragraph` 实体：
- 段落级独立权限控制。
- 段落级审计/版本回溯策略显著复杂化。
- 段落跨文档复用或引用关系。
- 段落级高并发写入需要独立分片策略。
