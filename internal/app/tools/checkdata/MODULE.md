# 数据关联检查工具

## Owner

`internal/app/tools/checkdata` 是项目关联字段诊断命令的 owner。它负责 MongoDB 连接配置、健康检查和代表性集合读取；`cmd/tools/check_data/main.go` 只保留兼容入口。

## 边界

- 这是只读诊断工具，不写入业务数据。
- MongoDB URI 和数据库名沿用 CI/部署环境变量，缺省保持本地默认值。
- 任一连接或读取失败必须返回错误，不能静默输出空字段。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/tools/checkdata ./cmd/tools/check_data
go vet ./internal/app/tools/checkdata ./cmd/tools/check_data
go build ./cmd/tools/check_data
```
