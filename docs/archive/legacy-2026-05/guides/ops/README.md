# Ops Guide Mirror Notice

`Qingyu_backend/docs/guides/ops` 当前不是运维文档主入口。

## Current Source Of Truth

- 当前运维文档 canonical owner 请看 [../../ops/README.md](../../ops/README.md)。
- 部署、监控、CI/CD、性能优化等内容，后续优先写入 `Qingyu_backend/docs/ops/`。

## Status Of This Directory

- 本目录与 `docs/ops/` 存在大量同名文件镜像。
- 其中至少有部分同名文件已经发生内容分叉，不能直接当作“完全副本”处理。
- 为避免误判最新版本，本目录后续只保留迁移提示和必要的历史兼容。

## Rule

1. 新增或更新运维正文，优先改 `Qingyu_backend/docs/ops/`。
2. 本目录如需改动，仅允许补迁移说明、引用修正和历史提示。
3. 后续若逐篇完成差异合并，再考虑把本目录收口为纯迁移提示目录。
