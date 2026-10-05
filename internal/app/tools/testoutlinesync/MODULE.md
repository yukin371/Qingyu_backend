# 大纲同步校验工具

`internal/app/tools/testoutlinesync` 是 `cmd/tools/test_outline_sync` 的 owner，只读检查项目中文档与大纲的双向关联及数量。

项目通过 `--project-id` 指定，MongoDB 配置复用 `mongodbconfig`；所有查询和统计错误必须返回。
