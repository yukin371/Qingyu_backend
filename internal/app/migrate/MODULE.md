# 数据迁移装配模块

## Owner

`internal/app/migrate` 是 MongoDB schema/index migration 命令的装配 owner。它负责参数校验、环境到数据库的映射、Mongo 生命周期和迁移注册；具体迁移实现仍归 `migration` 与 `migration/mongodb`。

`cmd/migrate/main.go` 只保留兼容命令入口。原有 `-command`、`-name`、`-env`、`-force` 参数和 `go run cmd/migrate/main.go` 用法保持不变。

## 边界

- 不承载业务 service、handler 或 repository 逻辑。
- 不读取 HTTP 请求或全局路由状态。
- 新迁移必须显式注册，并通过 `migration.SimpleMigration` 执行。
- 生产环境必须保留显式确认或 `--force` 保护。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/migrate ./cmd/migrate
go vet ./internal/app/migrate ./cmd/migrate
go build ./cmd/migrate
```
