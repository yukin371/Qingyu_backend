# 实时榜演示数据工具

`internal/app/tools/seedrealtimeranking` 是 `cmd/tools/seed_realtime_ranking` 的 owner，用于从书籍数据生成 realtime 排名演示记录。

默认只查询并预览；删除旧记录和插入新记录必须同时使用 `--apply` 与 `QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true`。MongoDB 配置复用 `mongodbconfig`。
