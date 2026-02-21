# Shared 架构简化 PR 草案

**日期**: 2026-02-13  
**范围**: `Qingyu_backend`（架构简化优先，功能实现暂缓）  
**目标**: 降低维护成本，收敛接口依赖，减少重复抽象与入口噪音。

## 建议 PR 标题

`refactor(shared): converge storage/stats ports and simplify router/container wiring`

## 变更摘要

1. `storage` 依赖倒置收敛
- API 层改为依赖 `service/shared/storage/interfaces.go` 导出的端口接口。
- 容器层移除 `StorageServiceImpl` 具体实现入口，统一使用端口接口注入。
- `router/shared` 删除未使用兼容壳函数，避免重复入口。

2. `stats` 端口命名与聚合接口收敛
- 统一使用 `StatsPort`，移除 `PlatformStatsService` 接口别名依赖链。
- 聚合端口拆分为 `UserAggregatorPort` / `ContentAggregatorPort`，并补齐契约测试。

3. `user` 路由链路去弱类型
- `RegisterUserRoutes` 从 `interface{}` 入参改为明确端口类型。
- `PublicUserHandler` 去运行时类型断言，直接依赖最小 `BookstoreService` 端口。
- `ProfileHandler` 显式注入可选 `StorageService`，去除注释式依赖。

4. 入口文件与容器可读性降噪
- `router/enter.go` 统一存储相关变量命名、日志口径与注释风格。
- 清理过时 TODO 和 `nolint` 占位噪音（保持行为不变）。
- `service/container/service_container.go` 统一 storage 初始化段局部变量语义命名与日志文案。

## 关键文件（建议纳入本 PR）

- `service/shared/storage/interfaces.go`
- `service/shared/storage/storage_service.go`
- `service/shared/storage/multipart_upload_service.go`
- `service/shared/storage/mock/mocks.go`
- `service/shared/storage/storage_service_test.go`
- `service/shared/storage/multipart_upload_service_test.go`
- `api/v1/shared/storage_api.go`
- `api/v1/shared/storage_api_test.go`
- `service/shared/stats/ports.go`
- `service/shared/stats/stats_service.go`
- `service/shared/stats/aggregators/user_aggregator.go`
- `service/shared/stats/aggregators/content_aggregator.go`
- `service/shared/stats/aggregators/aggregators_test.go`
- `service/shared/cache/redis_cache_service.go`
- `service/shared/cache/redis_cache_service_strategy_test.go`
- `service/shared/cache/strategies/cache_strategy.go`
- `service/shared/cache/strategies/ttl_strategy.go`
- `service/shared/cache/strategies/strategies_test.go`
- `service/container/service_container.go`
- `router/shared/shared_router.go`
- `router/user/user_router.go`
- `router/enter.go`
- `api/v1/user/handler/public_user_handler.go`
- `api/v1/user/handler/profile_handler.go`
- `api/v1/user/handler/stats_handler.go`
- `docs/plans/2026-02-12-shared-module-p1-refactor.md`

## 明确不在本 PR（避免范围扩散）

- 断点续传功能完整实现。
- 图片水印/文本水印功能。
- 多云存储适配器（S3/OSS/COS）功能实现。
- stats 聚合真实业务指标补齐（当前以端口与骨架收敛为主）。

## 回归验证（已执行）

- `go test ./service/shared/storage -count=1`
- `go test ./api/v1/shared ./service/shared/storage ./service/shared/cache/strategies ./service/shared/stats/... -count=1`
- `go test ./api/v1/user/handler ./router/user ./router -count=1`
- `go test ./service/container -run Test__CompileOnly__ -count=1`

## 风险与回滚

1. 风险点
- 主要风险是依赖注入链路签名调整导致编译期断裂；运行时行为变化风险低。

2. 回滚策略
- 若出现链路回归，可按模块回滚：
  - `router/user` 与 `api/v1/user/handler`（用户路由注入链）
  - `router/shared` 与 `service/container`（shared 路由注册链）
  - `service/shared/stats`（端口别名收敛）

## 推荐提交拆分

1. `storage` + `api/v1/shared` + storage tests  
2. `stats` ports/aggregators + stats tests  
3. `cache` strategy manager 接入 + tests  
4. `container` + `router/shared` + `router/enter`（依赖注入与入口降噪）  
5. `router/user` + `api/v1/user/handler`（去弱类型/命名收敛）  
6. `docs/plans` 更新

