# 后端中间件专题入口

> 最后整理: 2026-05-22

本目录目前是一个 `current-bounded` 专题区，只保留中间件集成类说明，不再扩张成第二个架构或标准总入口。

## First Read Path

1. [permission-integration-guide.md](./permission-integration-guide.md)
2. [../standards/layer-middleware.md](../standards/layer-middleware.md)
3. [../architecture/README.md](../architecture/README.md)

## Current Documents

| 文档 | 状态 | 用途 |
|------|------|------|
| [permission-integration-guide.md](./permission-integration-guide.md) | `current-bounded` | 权限中间件集成、接线和排障说明 |

## Boundary

- 本目录：中间件接入、替换、排障、兼容期说明。
- `../standards/layer-middleware.md`：中间件层的现行规范 owner。
- `../architecture/`：中间件在整体请求链中的位置和依赖边界 owner。
- `../implementation/infrastructure/`：中间件改造落地记录 owner。

## Practical Rules

1. 如果文档在讲“中间件应该怎么设计”，优先更新 `../standards/` 或 `../architecture/`。
2. 如果文档在讲“某个中间件如何接入或替换”，可以落在本目录。
3. 如果本目录后续不再只有单一专题，应优先并入更稳定的 owner，而不是继续在根层平铺新的中间件专题树。
