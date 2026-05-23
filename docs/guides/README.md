# Guides Index

本目录只应承载“上手路径、协作说明、操作指南”，不应继续充当 `api/`、`ops/`、`usage/` 的镜像容器。

## First Read Path

- [2026-04-07-backend-human-ai-onboarding.md](./2026-04-07-backend-human-ai-onboarding.md)
- [backend-context-quickstart.md](./backend-context-quickstart.md)
- [../README.md](../README.md)

## Current Boundary

- `guides/`：保留 onboarding、协作路径、开发/运维/排障类指南。
- `../api/`：API 参考与 Swagger 相关文档的 canonical owner。
- `../ops/`：当前运维文档 canonical owner。
- `../usage/`：仅保留历史 usage 迁移说明；真实旧文档已归档。
- `../testing/`：当前测试文档 canonical owner。

## Legacy Mirror Notice

- [api/README.md](./api/README.md)：历史 API 镜像已归档，请优先看 `../api/README.md`。
- [ops/README.md](./ops/README.md)：历史运维镜像已归档，请优先看 `../ops/README.md`。
- [usage/README.md](./usage/README.md)：历史 usage 镜像已归档，请优先看 `../usage/README.md`。

## Topic Directories

- [开发/README.md](./开发/README.md)：开发协作类操作指南
- [运维/README.md](./运维/README.md)：运维使用类指南补充
- [数据库/README.md](./数据库/README.md)：数据库优化与使用指南补充
- [测试/README.md](./测试/README.md)：历史测试报告与旧指南只读入口

## Rule

1. 新增 API 参考，不写到 `guides/api`。
2. 新增运维文档，不写到 `guides/ops`。
3. 新增使用说明，不在 `guides/usage` 和 `usage` 双写。
4. `guides/测试` 当前视为历史只读区；新的测试指南与报告分别回到 `../testing/` 或 `../standards/testing/`。
