# outline 查询诊断工具

`internal/app/tools/debugoutlinequery` 是 `cmd/tools/debug_outline_query` 的运行 owner，用于比较 ObjectID/string 查询并检查 BSON 解码形态。

工具只读，项目通过 `--project-id` 指定；MongoDB 配置复用 `mongodbconfig`。
