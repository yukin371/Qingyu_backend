# Backend Human/AI Onboarding (10min)

> 目标: 让新加入的人类开发者和 AI agent 在 10 分钟内建立后端“可工作级”上下文，而不是泛读历史文档。
> 范围: `Qingyu_backend` 当前代码结构与当前有效文档入口。

## 1. 10 分钟阅读路径

### 0-2 分钟: 先看真实入口

1. `cmd/server/main.go`
2. `core/server.go`

你只需要确认三件事:

- 配置和热重载在哪做
- Gin 和中间件在哪初始化
- 路由注册入口在哪里

### 2-5 分钟: 看运行时主链路

读 [后端运行时链路](../architecture/2026-04-07-backend-runtime-flow.md)。

重点抓住:

- `core.InitDB()` 当前是兼容保留入口，真实初始化在 `ServiceContainer.Initialize()`
- `service/container` 是运行时装配中心
- `router/enter.go` 不只是“挂路由”，还做了部分服务接线与初始化

### 5-8 分钟: 看模块地图

读 [后端模块地图](../architecture/2026-04-07-backend-module-map.md)。

重点抓住:

- 业务域 vs 横切能力 vs 平台基础设施
- 哪些模块对齐良好，哪些存在命名漂移和边界混叠

### 8-10 分钟: 看标准与风险

1. [后端标准索引](../standards/README.md)
2. [后端架构风险审查](../review/2026-04-07-backend-architecture-risk-review.md)

目的:

- 知道“按什么规则改”
- 知道“哪里最容易踩坑”

## 2. 代码目录到模块概念映射

| 目录 | 概念角色 | 你应该怎么理解 |
|---|---|---|
| `cmd/` | 进程入口与运维工具 | `cmd/server` 是主服务入口，其他多为迁移/seed/工具 |
| `core/` | 启动与装配编排 | 把配置、日志、服务初始化、Gin 启动串起来 |
| `router/` | HTTP 路由组装层 | 路由分组 + 部分能力接线；`enter.go` 是总入口 |
| `api/v1/` | Handler/API Facade 层 | 处理请求参数、响应封装，调用 service |
| `service/` | 业务逻辑与平台能力 | 包含业务域服务、容器、事件总线、shared 能力 |
| `repository/` | 数据访问层 | interfaces + mongodb/redis/search 等实现 |
| `models/` | 领域模型与 DTO | 业务模型、共享类型、DTO |
| `internal/middleware/` | 请求横切逻辑 | 认证、日志、恢复、限流等 |
| `pkg/` | 通用基础设施 | cache/logger/metrics/transaction 等复用能力 |

## 3. 启动链路简述

启动主链路可简化为:

`cmd/server/main.go -> config.LoadConfig -> core.InitServer -> core.InitServices -> service.InitializeServices -> ServiceContainer.Initialize -> gin.New + middlewares -> router.RegisterRoutes -> core.RunServer`

对排障最关键的两个事实:

1. 业务服务实例主要在 `ServiceContainer` 内部构建和持有。
2. 路由不是纯静态注册，部分服务会在 `router/enter.go` 阶段做可用性判断与接线。

## 4. 当前最重要的 5 个模块/入口

1. `service/container/service_container.go`
   - 后端装配中心，决定服务生命周期和依赖关系。
2. `router/enter.go`
   - 路由总入口，也是当前部分初始化与兼容接线中心。
3. `service/writer/`
   - 复合子域，包含 project/document/storyharness/publish/stats 等多能力。
4. `service/shared/`
   - 横切能力集中区，复用价值高但也是职责膨胀风险点。
5. `service/ai/`
   - 同时连接 quota/chat/context/外部 AI，属于平台化关键能力。

## 5. 常见误解与坑点

1. 误解: `core.InitDB()` 完成数据库初始化。
   - 现实: 真实 Mongo/Redis/EventBus 初始化主链路在 `ServiceContainer.Initialize()`。
2. 误解: `router/*` 和 `api/v1/*` 是重复层。
   - 现实: `router` 负责分组与装配，`api/v1` 负责 handler/facade。
3. 误解: `shared` 可以继续承接任何“先放这里”的能力。
   - 现实: `shared` 已接近职责黑洞，新增能力要先判断是否应独立模块化。
4. 误解: 模块命名天然一致。
   - 现实: 存在 `notification/notifications`、`user/users`、`stats/reading-stats` 等漂移。
5. 误解: 旧设计与旧 review 可以直接当现行事实。
   - 现实: 先看 `architecture/README.md`、`standards/README.md`、`review` 最新报告，再回看历史材料。

## 6. 深入阅读索引

### Architecture

- [架构主入口](../architecture/README.md)
- [后端运行时链路](../architecture/2026-04-07-backend-runtime-flow.md)
- [后端模块地图](../architecture/2026-04-07-backend-module-map.md)
- [系统架构总览](../architecture/system_architecture.md)

### Standards

- [标准主入口](../standards/README.md)
- [后端架构与模块标准](../standards/2026-04-07-backend-architecture-and-module-standards.md)
- [Service 层标准](../standards/layer-service.md)
- [Repository 层标准](../standards/layer-repository.md)

### Review

- [后端架构风险审查](../review/2026-04-07-backend-architecture-risk-review.md)
- [Review 区索引](../review/README.md)

## 7. 给 AI Agent 的最小上下文包

如果你要让 AI 快速进入“可执行分析”状态，优先提供这 4 份文档:

1. 本文档
2. `../architecture/2026-04-07-backend-runtime-flow.md`
3. `../architecture/2026-04-07-backend-module-map.md`
4. `../review/2026-04-07-backend-architecture-risk-review.md`

并要求 AI 先回答:

- 真实启动链路与运行时装配点
- 当前 5 个最关键模块/入口
- 3 个最高风险架构坑点
