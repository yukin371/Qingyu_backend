# 大纲关联检查工具

## Owner

`internal/app/tools/checkoutline` 是大纲层级与关联字段诊断命令的 owner。它负责 MongoDB 连接配置、游标生命周期和只读输出；`cmd/tools/check_outline/main.go` 只保留兼容入口。

## 边界

- 这是只读诊断工具，不写入业务数据。
- MongoDB URI 和数据库名沿用 CI/部署环境变量，缺省保持本地默认值。
- 查询、解码和游标错误必须返回，不能静默截断检查结果。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/tools/checkoutline ./cmd/tools/check_outline
go vet ./internal/app/tools/checkoutline ./cmd/tools/check_outline
go build ./cmd/tools/check_outline
```
