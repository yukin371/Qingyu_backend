# 后端审查文档入口

> 最后整理: 2026-05-22  
> 当前状态: `current-bounded`

`docs/review/` 用于沉淀“审查结论与风险判断”，回答“问题严重吗、证据是什么、需要跟进什么”，不作为设计规范或实现说明的替代。

## Recommended Read Path

1. [README.md](./README.md)
2. [2026-04-07-backend-architecture-risk-review.md](./2026-04-07-backend-architecture-risk-review.md)
3. [../issues/README.md](../issues/README.md)

## Current Materials

| 文档 | 状态 | 用途 |
|------|------|------|
| [2026-04-07-backend-architecture-risk-review.md](./2026-04-07-backend-architecture-risk-review.md) | `current-bounded` | 当前后端架构风险主入口，覆盖运行时入口、模块边界、命名契约和历史 issue 状态 |

## Archive Reference

- [../archive/legacy-2026-05/review/2025-10-batch/](../archive/legacy-2026-05/review/2025-10-batch/): 2025-10 批次阶段性 review 归档，仅用于回溯阶段判断，不再作为当前现状入口。

## Boundary

- 本目录：审查结论、风险判断、证据汇总、整改建议。
- `../issues/`: 可执行跟踪项与问题台账 owner。
- `../architecture/`: 当前结构事实 owner。
- `../implementation/`: 已实施改动与修复经过 owner。

## Practical Rules

1. 新一轮审查新增独立报告，不覆盖历史报告正文。
2. 审查报告负责“现象、证据、影响、状态、动作”，不重复写实现细节或长实施步骤。
3. 可执行跟踪项必须同步到 `../issues/`，避免审查结论孤立。
