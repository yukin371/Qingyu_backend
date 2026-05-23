# 后端本地 plans 迁移说明

> 最后整理: 2026-05-22  
> 当前状态: `legacy-readonly`

`Qingyu_backend/docs/plans/` 已退出“长期计划 owner”角色，当前只保留迁移说明，避免后端子模块重新长出一棵本地 plan 树。

## Recommended Read Path

1. [../../../docs/plans/submodules/backend/README.md](../../../docs/plans/submodules/backend/README.md)
2. [../architecture/README.md](../architecture/README.md)
3. [../implementation/README.md](../implementation/README.md)

## Current Boundary

- 本目录：仅保留“计划已迁移到哪里”的说明。
- 父仓 [../../../docs/plans/submodules/backend/README.md](../../../docs/plans/submodules/backend/README.md): 后端长期计划、专题设计、实施方案与阶段计划的唯一入口。
- `../architecture/`: 当前架构事实与边界 owner。
- `../implementation/`: 已落地实施记录 owner。

## Parent Taxonomy

后端新计划应统一写入父仓以下专题，而不是在本地重建平行目录：

- `architecture/`: 当前架构、模型与版本治理设计
- `api-governance/`: API 标准化、接口治理、错误处理治理
- `publication/`: 发布与审核工作流设计
- `shared-and-layering/`: shared 模块与分层重构
- `testing-and-quality/`: 测试策略与质量治理
- `legacy-phases/`: 历史 rollout 计划与完成报告

## Practical Rules

1. 不要在 `Qingyu_backend/docs/plans/` 新增 backend 计划、设计或长期治理文档。
2. 如果文档在讲“未来怎么改、分几步改、风险与验收是什么”，应写到父仓 `docs/plans/submodules/backend/`。
3. 如果文档已经变成“当前事实”，应回写到 `../architecture/`、`../standards/` 或 `../guides/`。
4. 如果文档已经变成“实施完成记录”，应回写到 `../implementation/` 或 `../review/`。
