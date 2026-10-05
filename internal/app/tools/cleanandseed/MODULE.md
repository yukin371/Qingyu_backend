# 演示数据清理与生成工具

`internal/app/tools/cleanandseed` 是 `cmd/tools/clean_and_seed` 的 owner，用于清理演示集合并重建完整示例项目。

默认只预览；执行必须同时使用 `--apply` 和 `QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true`。MongoDB 配置复用 `mongodbconfig`，清理、查询、插入和验证错误都会中止命令并返回。
