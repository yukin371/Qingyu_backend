# character BSON 诊断工具

`internal/app/tools/debugcharacter` 是 `cmd/tools/debug_character` 的运行 owner，只读输出一条 character 文档的 BSON 字段类型。

MongoDB 配置统一复用 `mongodbconfig`；连接、Ping 和无数据错误必须返回给命令入口。
