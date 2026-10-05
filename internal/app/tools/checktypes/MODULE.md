# 项目类型检查工具

## Owner

`internal/app/tools/checktypes` 是项目关联类型诊断和可选演示数据清理命令的 owner；`cmd/tools/check_types/main.go` 只保留兼容入口。

## 安全边界

- 默认只读检查项目记录。
- 删除非演示项目必须显式传入 `--apply`，且找不到演示项目时拒绝全量删除。
- MongoDB 连接、查询、解码和删除错误必须返回。
- URI/数据库配置复用 `internal/app/tools/mongodbconfig`。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/tools/checktypes ./cmd/tools/check_types
go vet ./internal/app/tools/checktypes ./cmd/tools/check_types
go build ./cmd/tools/check_types
```
