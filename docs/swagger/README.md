# Swagger 生成说明

## 生成方式

推荐直接使用仓库脚本：

Windows PowerShell:

```powershell
.\scripts\docs\generate_swagger.ps1
```

Linux / macOS:

```bash
./scripts/docs/generate_swagger.sh
```

也可以使用 `make`：

```bash
make swagger
```

校验产物是否需要更新：

```bash
make swagger-check
```

严格校验产物是否已经提交：

```bash
make swagger-verify-committed
```

## 当前生成参数

Swagger 入口文件：

```text
api/v1/swagger.go
```

扫描目录：

```text
api/v1
pkg/response
models
models/dto
service/interfaces
service/ai/dto
service/shared/storage
service/shared/stats
```

输出目录：

```text
docs
```

## 产物

- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

## 说明

- 当前方案使用 `--parseDependency=false`，避免 `swaggo/swag` 在本项目上递归解析依赖时崩溃。
- 如果新增注解引用了新的外部 DTO 包，需要把对应目录加入脚本里的 `-d` 列表。
- `make swagger-check` 只校验 Swagger 产物是否可以成功生成，适合本地快速检查和 CI。
- `make swagger-verify-committed` 会在生成后对 `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml` 做 `git diff` 校验，适合提交前人工确认，不建议直接放进 CI。
- 仓库还提供了独立工作流 `Swagger Artifact Sync`，仅在 `api/v1/**` 发生变更时触发，用于校验产物是否已同步提交。
