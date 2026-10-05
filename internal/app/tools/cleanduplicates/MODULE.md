# 重复大纲清理工具

## Owner

`internal/app/tools/cleanduplicates` 是重复 arc 大纲诊断与清理命令的 owner；`cmd/tools/clean_duplicates/main.go` 只保留兼容入口。

## 安全边界

- 默认只预览重复项，不删除数据。
- 必须同时传入 `--apply` 和设置 `QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true`。
- 可通过 `--project-id` 指定目标项目，默认值仅为历史诊断项目 ID。
- 查询、解码、删除和清理后统计失败均返回错误；无效字段不会触发强制类型断言。
- URI/数据库配置复用 `internal/app/tools/mongodbconfig`。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/tools/cleanduplicates ./cmd/tools/clean_duplicates
go vet ./internal/app/tools/cleanduplicates ./cmd/tools/clean_duplicates
go build ./cmd/tools/clean_duplicates
```
