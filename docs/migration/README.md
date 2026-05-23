# 后端迁移文档入口

> 最后整理: 2026-05-22

本目录只承载“文档层面的迁移、兼容清理和路径切换说明”，不等同于代码里的 `migration/` 或 `migrations/` 包。

## First Read Path

1. [architecture_migration_guide.md](./architecture_migration_guide.md)
2. [auth-module-migration.md](./auth-module-migration.md)
3. [../standards/documentation-taxonomy.md](../standards/documentation-taxonomy.md)

## Current Documents

| 文档 | 状态 | 用途 |
|------|------|------|
| [architecture_migration_guide.md](./architecture_migration_guide.md) | `current-bounded` | 架构迁移与目录演进说明 |
| [auth-module-migration.md](./auth-module-migration.md) | `current-bounded` | 鉴权模块迁移与兼容边界说明 |
| [ci-dependency-check-update.md](./ci-dependency-check-update.md) | `current-bounded` | CI 依赖检查调整记录 |

## Boundary

- 本目录：解释“为什么迁、迁到哪里、兼容期如何处理、旧路径如何退出”。
- `../architecture/`：当前架构事实与依赖边界 owner。
- `../implementation/`：迁移已经落地后的实施记录 owner。
- 代码目录中的 `migration/` / `migrations/`：运行时迁移逻辑、数据库迁移或脚本 owner，不由这里替代。

## Practical Rules

1. 新增“弃用、迁移、兼容退出”说明时，优先落到本目录。
2. 如果内容已经变成“当前规范”，应回写到 `../architecture/`、`../standards/` 或 `../guides/`，不要长期停在迁移说明里。
3. 如果内容已经变成“具体实施经过”，应放到 `../implementation/`。
4. 不要把数据库 schema migration、运行脚本或代码包路径说明误写进本目录；这里只负责文档和治理视角的迁移边界。
