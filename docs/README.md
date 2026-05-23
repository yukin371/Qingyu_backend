# Qingyu Backend Docs Hub

本目录是 `Qingyu_backend` 的文档总入口，目标是先帮人判断“哪份是当前入口、哪份是历史资料、哪份已经迁移”，避免继续在重复目录里新增内容。

## First Read Path

1. [../README.md](../README.md)
2. [../MODULE.md](../MODULE.md)
3. [architecture/README.md](./architecture/README.md)
4. [standards/README.md](./standards/README.md)
5. [guides/README.md](./guides/README.md)
6. [testing/README.md](./testing/README.md)

## Root File Roles

- [../README.md](../README.md): 人类开发者的项目入口、启动方式、主要文档导航。
- [../MODULE.md](../MODULE.md): 后端模块职责、owner、边界和文档同步规则。
- [../qingyubackend.md](../qingyubackend.md): AI/规则文件，不作为人工优先阅读入口。

## Current Source Of Truth

- [architecture/](./architecture/): 当前架构边界、运行时链路、模块地图。
- [standards/](./standards/): 现行标准与分层规则。
- [guides/](./guides/): 人类/AI 快速上手与操作指南。
- [api/](./api/): API 参考、Swagger 导出说明、接口文档。
- [ops/](./ops/): 当前运维与部署文档入口。
- [testing/](./testing/): 当前测试指南、测试报告、测试专题入口。
- [review/](./review/): 审查结论和阶段性风险说明。
- [analysis/](./analysis/): 分析、评估、治理专项记录。
- [issues/](./issues/): 问题台账与归档问题。
- [database/](./database/): 数据库与种子数据操作文档。
- [implementation/](./implementation/): 已落地实现说明与归档实现记录。
- [migration/](./migration/): 文档迁移或依赖治理类说明，不等同于代码里的 `migration/` 或 `migrations/` 包。
- [archive/](./archive/): 已归档的历史目录、镜像目录和阶段资料。

补充规范：

- [standards/documentation-taxonomy.md](./standards/documentation-taxonomy.md): 后端 docs 根目录分类标准、状态模型和落盘规则。

## Directory Status Matrix

### Current Owners

| 目录 | 状态 | 说明 |
|------|------|------|
| `architecture/` | `current-owner` | 当前架构与模块边界入口 |
| `standards/` | `current-owner` | 现行标准入口 |
| `guides/` | `current-owner` | 上手与操作指南入口 |
| `api/` | `current-owner` | API 入口 |
| `ops/` | `current-owner` | 运维入口 |
| `testing/` | `current-owner` | 测试入口 |

### Current But Bounded

| 目录 | 状态 | 说明 |
|------|------|------|
| `analysis/` | `current-bounded` | 分析、审计、治理评估 |
| `review/` | `current-bounded` | 审查结论与风险结论 |
| `issues/` | `current-bounded` | 问题台账 |
| `database/` | `current-bounded` | 数据库专题 |
| `deployment/` | `current-bounded` | 部署专题；总入口仍归 `ops/` |
| `implementation/` | `current-bounded` | 已落地实现说明 |
| `migration/` | `current-bounded` | 迁移说明与兼容清理 |
| `tools/` | `current-bounded` | 工具专题说明 |
| `middleware/` | `current-bounded` | 中间件集成专题目录，后续需继续收口 |

### Legacy Or Archive

| 目录 | 状态 | 说明 |
|------|------|------|
| `design/` | `legacy-readonly` | 历史设计资料区 |
| `engineering/` | `legacy-readonly` | 历史需求/SRS 路径 |
| `usage/` | `legacy-readonly` | 历史使用说明路径 |
| `test/` | `legacy-readonly` | 旧测试目录路径 |
| `plans/` | `legacy-readonly` | 本地长期计划已迁往父仓 |
| `archive/` | `archive-only` | 历史镜像和归档批次 |

## Duplicate And Legacy Paths

以下目录当前存在历史残留或仅保留迁移说明，不应继续作为新文档主入口：

- `guides/api/`: 仅保留迁移说明；历史镜像已归档到 `archive/legacy-2026-05/guides/api/`。
- `guides/ops/`: 仅保留迁移说明；历史镜像已归档到 `archive/legacy-2026-05/guides/ops/`。
- `guides/usage/`: 仅保留迁移说明；历史镜像已归档到 `archive/legacy-2026-05/guides/usage/`。
- `usage/`: 仅保留迁移说明；历史使用说明已归档到 `archive/legacy-2026-05/usage/`。
- `test/`: 仅保留迁移说明；旧测试目录已归档到 `archive/legacy-2026-05/test/`。
- `engineering/`: 仅保留迁移说明；历史需求/SRS 文档已归档到 `archive/legacy-2026-05/engineering/`。
- `review/`: 当前只保留精简后的现状入口；2025-10 批次阶段性 review 已归档。
- `design/`: 当前已收口为历史设计资料区；根层散落历史文件已归档到 `archive/legacy-2026-05/design/root-files/`，保留中的历史专题树已统一补标准 `README.md` 入口。
- `plans/`: 已迁移到父仓库 [../../docs/plans/submodules/backend/README.md](../../docs/plans/submodules/backend/README.md)，本地不再新增 backend 计划文档。

## Practical Rules

1. 新增后端文档前，先判断它属于 `architecture`、`standards`、`guides`、`api`、`ops`、`testing` 中的哪一类，避免在重复目录落盘。
2. 如果文档只是历史报告或阶段结论，不要把它放进主入口 README；应放在 `review/`、`analysis/`、`issues/archived/` 或 `archive/` 对应归档区。
3. 发现 `guides/api`、`guides/ops`、`guides/usage`、`usage`、`test`、`engineering` 中仍有被引用的旧资料，优先补迁移说明，不要恢复双写结构。
4. 新增跨仓库 plan、长期治理方案或 ADR 时，落到父仓库 `docs/`，不要回写本子模块 `docs/plans/`。
5. 新增目录前先检查 [standards/documentation-taxonomy.md](./standards/documentation-taxonomy.md)，优先归入已存在 owner，避免再长出新的平行入口。
6. `migration/` 与 `middleware/` 现在都有正式 README；后续涉及迁移说明或中间件接入时，先更新对应入口，不要继续把根层散文档直接平铺出来。
