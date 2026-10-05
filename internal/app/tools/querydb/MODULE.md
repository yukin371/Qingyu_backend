# Mongo 查询诊断工具

`internal/app/tools/querydb` 是 `cmd/tools/query_db` 的运行 owner，按自然顺序输出最近的大纲和文档记录。

工具只读；连接、Ping、查询、游标关闭前的解码错误必须返回。MongoDB 配置复用 `mongodbconfig`。
