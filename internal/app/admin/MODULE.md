# 管理员引导装配模块

## Owner

`internal/app/admin` 是管理员引导命令的应用装配 owner。它负责参数校验、配置发现、MongoDB 生命周期和管理员账号初始化用例编排。

`cmd/admin/main.go` 只负责命令分发、帮助文本和退出码；`admin bootstrap` 的命令行用法保持不变。

## 边界

- 可以依赖 `config`、管理员/用户模型和 MongoDB driver。
- 不承载 HTTP handler、路由或普通用户业务 API。
- 密码只能从指定环境变量读取，不接受命令行明文密码。
- 管理员创建/重置必须保留已有账号冲突检查、角色合并和操作日志。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/admin ./cmd/admin
go vet ./internal/app/admin ./cmd/admin
go build ./cmd/admin
go run ./cmd/admin help
```
