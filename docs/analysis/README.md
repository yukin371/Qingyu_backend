# 后端分析文档入口

> 最后整理: 2026-05-22  
> 当前状态: `current-bounded`

本目录保存分析、评估、治理盘点与复杂度观察结果，用来回答“现状如何、复杂度在哪、清理成本多大”，不承担“当前规则”或“实施步骤”的 owner。

## Recommended Read Path

1. [README.md](./README.md)
2. [2026-04-07-backend-legacy-cleanup-assessment.md](./2026-04-07-backend-legacy-cleanup-assessment.md)
3. [2026-04-07-backend-legacy-cleanup-phase-report.md](./2026-04-07-backend-legacy-cleanup-phase-report.md)
4. [2026-01-29-writer-migration-analysis.md](./2026-01-29-writer-migration-analysis.md)

## Current Materials

| 文档 | 状态 | 用途 |
|------|------|------|
| [2026-04-07-backend-legacy-cleanup-assessment.md](./2026-04-07-backend-legacy-cleanup-assessment.md) | `analysis` | 历史清理范围、风险和优先级评估 |
| [2026-04-07-backend-legacy-cleanup-phase-report.md](./2026-04-07-backend-legacy-cleanup-phase-report.md) | `analysis` | 阶段治理结果与现状快照 |
| [2026-04-07-backend-legacy-cleanup-phase-2-guide.md](./2026-04-07-backend-legacy-cleanup-phase-2-guide.md) | `analysis` | 第二阶段治理参考说明 |
| [2026-01-29-writer-migration-analysis.md](./2026-01-29-writer-migration-analysis.md) | `analysis` | writer 迁移复杂度与拆解观察 |
| [service_unified_error_migration_effort.md](./service_unified_error_migration_effort.md) | `analysis` | 统一错误处理迁移工作量评估 |

## Non-Markdown Asset

- `writer-complexity-matrix.json`: 分析辅助数据，不是人工主入口。

## Boundary

- 本目录：分析、盘点、复杂度评估、治理成本判断。
- `../architecture/`: 当前结构事实与模块边界 owner。
- `../standards/`: 长期规则与必须遵循的标准 owner。
- `../implementation/`: 已落地实施记录 owner。
- 父仓 `../../../docs/plans/`: 跨仓库长期方案与活跃计划 owner。

## Practical Rules

1. 如果文档在回答“现状问题在哪、复杂度多大、清理顺序如何排”，可以落在本目录。
2. 如果分析结论已经稳定成长期规则，应继续收口到 `../standards/`、`../architecture/` 或父仓计划文档。
3. 如果内容已经变成“怎么改、改了什么、如何验证”，应回写到 `../implementation/`、`../review/` 或父仓 `docs/plans/`，不要长期停在分析目录。
