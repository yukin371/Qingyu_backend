# 测试数据清理工具

## Owner

`internal/app/tools/cleartestdata` 是测试数据集合清理命令的 owner；`cmd/tools/clear_test_data/main.go` 只保留兼容入口。

## 安全边界

- 默认只显示提示，不删除任何数据。
- 必须同时传入 `--apply` 并设置 `QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true`。
- MongoDB 连接、Ping 和任一集合删除失败都返回错误，不继续静默清理。
- URI/数据库配置复用 `internal/app/tools/mongodbconfig`。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/tools/cleartestdata ./cmd/tools/clear_test_data
go vet ./internal/app/tools/cleartestdata ./cmd/tools/clear_test_data
go build ./cmd/tools/clear_test_data
```
