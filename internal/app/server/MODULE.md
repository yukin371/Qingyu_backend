# 服务启动装配模块

## Owner

`internal/app/server` 是后端 HTTP 服务的启动装配 owner。它负责配置路径解析、配置热重载注册、数据库初始化和 `core` 服务启动顺序。

`cmd/server/main.go` 只保留进程入口和 Swagger 元数据，Air、Docker、CI 继续通过原路径启动。

## 边界

- 可以依赖 `config`、`core` 以及内部基础设施包。
- 不承载业务 handler、service 或 repository 实现。
- 不改变监听端口、环境变量名称和配置文件兼容路径。
- 新增 worker 或迁移入口应建立各自的装配模块，不直接复制本模块的启动逻辑。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/server ./cmd/server
go vet ./internal/app/server ./cmd/server
go build ./cmd/server
```
