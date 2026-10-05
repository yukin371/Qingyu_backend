# 重复大纲检查工具

## Owner

`internal/app/tools/checkduplicates` 是 arc 大纲重复标题诊断命令的 owner；`cmd/tools/check_duplicates/main.go` 只保留兼容入口。

## 边界

- 工具只读，不执行删除；清理由独立 `clean_duplicates` 命令负责并要求显式安全闸门。
- 可通过 `--project-id` 指定目标项目，默认值仅为历史诊断项目 ID。
- MongoDB URI/数据库配置复用 `internal/app/tools/mongodbconfig`。
- BSON 类型不匹配时输出空/零值并保持报告可读，不触发 panic。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/tools/checkduplicates ./cmd/tools/check_duplicates
go vet ./internal/app/tools/checkduplicates ./cmd/tools/check_duplicates
go build ./cmd/tools/check_duplicates
```
