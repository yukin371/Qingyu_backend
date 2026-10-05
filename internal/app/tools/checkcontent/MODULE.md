# document_contents 检查工具

## Owner

`internal/app/tools/checkcontent` 是文档内容计数检查命令的 owner。它负责 MongoDB 连接配置、健康检查和统计错误处理；`cmd/tools/check_content/main.go` 只保留兼容入口。

## 边界

- 这是只读诊断工具，不写入业务数据。
- MongoDB URI 和数据库名优先读取 CI/部署注入的环境变量，缺省保持本地 `localhost:27017` 与 `qingyu`。
- 连接、ping 和查询失败必须返回错误，不能静默输出错误结果。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/tools/checkcontent ./cmd/tools/check_content
go vet ./internal/app/tools/checkcontent ./cmd/tools/check_content
go build ./cmd/tools/check_content
```
