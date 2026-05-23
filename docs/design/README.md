# Design Legacy Hub

`Qingyu_backend/docs/design` 现在被视为“历史设计资料区”，不再作为当前实现、联调或治理文档的默认入口。

## Page Role

- legacy-root-hub
- current-owner: `Qingyu_backend/docs/design/`
- current-bounded: 后端历史设计总入口，只负责把读者分流到 current source of truth 或具体历史专题

## Recommended Read Path

1. 先看本页的 `Current Source Of Truth`。
2. 如果你需要当前计划，跳到父仓 `docs/plans/submodules/backend/README.md`。
3. 只有在需要追溯某个历史专题时，才进入本目录下的各子树。

## Boundary

- 本目录是历史设计资料区，不是当前实现与长期计划 owner。
- 当前架构事实优先看 `../architecture/`、当前落地过程优先看 `../implementation/`、当前 API 入口优先看 `../api/`。
- 子目录里保留的大量正文应默认理解为历史设计稿，除非该子树已明确标注 current mainline/topic 语义。

## Quick Section Map

- Current Status
- Topic Status Matrix
- Current Source Of Truth
- Rule

## Quick Takeaways

- `design/` 现在主要承担“历史设计资料导航”职责。
- 当前事实、当前计划和当前实施都已经有各自更合适的 owner。

## Skip Guide

- 只看当前架构：去 `../architecture/README.md`。
- 只看当前实施：去 `../implementation/README.md`。
- 只看当前计划：去父仓 `docs/plans/submodules/backend/README.md`。

## Current Status

- 根层原先平铺的历史总览、阶段进度、完成报告、模板和旧版架构概览，已迁入：
  - [../archive/legacy-2026-05/design/README.md](../archive/legacy-2026-05/design/README.md)
- 当前保留在本目录下的，主要是按专题拆开的历史设计子目录，便于按主题回溯。
- 旧的 `modules/` 九大模块镜像树也已迁入归档区，不再与专题树并行保留正文。
- 本轮已为仍保留正文的历史专题树统一补齐标准 `README.md` 入口，后续不再要求读者直接从 `README_*.md` 文件名猜目录角色。

## Topic Status Matrix

### 已归档，仅保留迁移说明

| 目录 | 当前状态 | 归档去向 |
|------|----------|----------|
| `modules/` | 迁移说明 | `../archive/legacy-2026-05/design/modules/` |
| `阅读端/` | 迁移说明 | `../archive/legacy-2026-05/design/阅读端/` |

### 仍保留正文的历史专题树

| 目录 | 当前状态 | 备注 |
|------|----------|------|
| `ai/` | `legacy-live` | AI 历史设计区；入口已按“当前主线 / 完成态专题 / 历史参考”分层 |
| `auth/` | `legacy-live` | 认证相关历史设计 |
| `communication/` | `legacy-live` | 通信与消息历史设计 |
| `core/` | `legacy-live` | 核心能力与基础架构历史设计 |
| `database/` | `legacy-live` | 数据库历史设计 |
| `middleware/` | `legacy-live` | 中间件历史设计 |
| `ops/` | `legacy-live` | 运维历史设计；当前运维 owner 已转回 `docs/ops/` |
| `platform/` | `legacy-live` | 平台层历史设计 |
| `reader/` | `legacy-live` | 阅读端历史设计；当前 owner 已转回 `architecture/api/implementation` |
| `recommendation/` | `legacy-live` | 推荐历史设计 |
| `security/` | `legacy-live` | 安全历史设计 |
| `shared/` | `legacy-live` | 共享底层服务历史设计 |
| `stats/` | `legacy-live` | 统计历史设计 |
| `testing/` | `legacy-live` | 测试设计稿 |
| `writer/` | `legacy-live` | 写作端历史设计 |
| `重构规划/` | `legacy-live` | 重构规划历史设计 |

补充说明：

- 上表中的 `legacy-live` 目录现已统一具备标准 `README.md` 入口。
- 这些 README 只负责历史索引与只读提示，不表示目录重新成为当前 owner。
- 上表中若同时存在旧 `README_*.md`，该文件只视为历史正文，不再视为目录标准入口。

## Current Source Of Truth

如果你要找的是当前文档 owner，请优先看：

- [../architecture/README.md](../architecture/README.md)
- [../standards/README.md](../standards/README.md)
- [../api/README.md](../api/README.md)
- [../ops/README.md](../ops/README.md)
- [../testing/README.md](../testing/README.md)

## Rule

1. 新的架构边界、标准和实施规则不要继续写入 `docs/design/`。
2. 只有明确属于历史设计稿、旧方案或回溯资料的内容，才继续放在本目录或其归档桶中。
3. 若未来继续做第二轮精简，优先处理其他仍保留正文的历史专题树，并继续减少 `design/` 作为正文承载区的范围。
