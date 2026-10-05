# 诊断工具 MongoDB 配置

## Owner

`internal/app/tools/mongodbconfig` 是后端只读诊断工具的 MongoDB URI 和数据库名环境解析 owner。

## 规则

- 优先使用 `QINGYU_DATABASE_PRIMARY_MONGODB_URI` 与 `QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE`。
- 兼容 `MONGODB_URI` 与 `MONGODB_DATABASE`。
- 未设置时使用本地默认 `mongodb://localhost:27017/?directConnection=true` 和 `qingyu`。
- 不负责建立连接、Ping 或关闭 client；生命周期由具体工具 owner 管理。
