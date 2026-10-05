# 基准测试装配模块

## Owner

`internal/app/benchmark` 是性能基准命令的装配 owner。它负责命令参数解析、测试场景组装、超时控制和结果文件输出；HTTP 基准实现仍由 `benchmark` 包持有。

`cmd/benchmark/main.go` 只保留兼容进程入口，因此现有 `go run cmd/benchmark/main.go` 和性能脚本路径不变。

## 边界

- 不依赖业务数据库、服务容器或 HTTP 路由注册。
- 不把 benchmark 结果写入业务数据源。
- 参数默认值和 JSON 结果结构需要保持兼容。

## 验证

在 `Qingyu_backend` 执行：

```text
go test ./internal/app/benchmark ./cmd/benchmark
go vet ./internal/app/benchmark ./cmd/benchmark
go build ./cmd/benchmark
```
