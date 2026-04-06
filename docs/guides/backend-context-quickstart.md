# Backend Context Quickstart

> 目标: 让第一次接手 Qingyu Backend 的人类或 AI 在 20 分钟内建立足够深的项目上下文

## 1. 先建立什么认知

不要一上来扫完整个仓库。先建立这四个事实：

1. 服务启动链路在哪里
2. 路由如何进入业务服务
3. 哪些模块是核心业务域
4. 当前哪些文档是主入口，哪些只是历史背景

## 2. 推荐阅读顺序

### 第 1 步: 看入口

- `cmd/server/main.go`
- `core/server.go`
- `core/init_db.go`
- `service/enter.go`

这一步只需要回答：

- 配置在哪里加载
- 服务容器在哪里创建
- Gin 在哪里初始化
- 哪些中间件最先执行

### 第 2 步: 看运行时主轴

读 [后端运行时架构](../architecture/2026-04-07-backend-runtime-architecture.md)。

重点关注：

- `core.InitDB()` 已经是兼容保留 no-op
- 真正的 Mongo/Redis/EventBus 初始化在 `ServiceContainer`
- 路由层采用渐进式注册，不是一次性全量失败

### 第 3 步: 看模块地图

读 [后端模块地图](../architecture/2026-04-07-backend-module-map.md)。

重点关注：

- `writer / bookstore / reader` 是主业务三角
- `api/v1` 和 `router/*` 是两层，不是重复层
- `search / events / shared / internalapi` 是横切能力，不是纯业务域

### 第 4 步: 看标准

读 [后端架构与文档标准](../standards/backend-architecture-documentation-standard.md)。

重点关注：

- 哪些文档被视为 source of truth
- Mermaid 图该怎么画
- 模块文档最少需要写什么

### 第 5 步: 看风险

读 [后端架构风险审查](../review/2026-04-07-backend-architecture-risk-review.md)。

这份文档会告诉你：

- 哪些地方最容易被旧文档误导
- 哪些结构虽然能跑，但维护成本高
- 哪些是后续架构整治的优先点

## 3. 调试或继续深入时怎么找

### 如果你在追启动问题

从这里开始：

- `cmd/server/main.go`
- `core/server.go`
- `service/enter.go`
- `service/container/service_container.go`

### 如果你在追某个 HTTP 接口

从这里开始：

- `router/<domain>`
- `api/v1/<domain>`
- `service/<domain>`
- `repository/interfaces` / `repository/mongodb`

### 如果你在追共享能力

优先看：

- `service/container`
- `service/events`
- `service/shared`
- `internal/middleware`
- `pkg`

### 如果你在追写作链路

优先看：

- `router/writer`
- `api/v1/writer`
- `service/writer`
- `models/writer`
- `repository/mongodb/writer`

## 4. 当前最容易踩的认知坑

1. 看到 `InitDB` 就以为数据库在这里初始化。
2. 看到 `router` 和 `api/v1` 以为有重复实现。
3. 看到 `ProviderRegistry` 就以为容器已经完全 provider 化。
4. 看到 Gin 启动成功就以为所有业务路由都已可用。
5. 直接把 `docs/design` 或旧 `review` 当作当前 source of truth。

## 5. 给 AI 的最短提示词骨架

如果后续要让 AI 快速理解后端，建议先喂这几个文件：

1. `Qingyu_backend/docs/guides/backend-context-quickstart.md`
2. `Qingyu_backend/docs/architecture/2026-04-07-backend-runtime-architecture.md`
3. `Qingyu_backend/docs/architecture/2026-04-07-backend-module-map.md`
4. `Qingyu_backend/docs/review/2026-04-07-backend-architecture-risk-review.md`

推荐要求 AI 先回答：

- 真实启动链路
- 当前运行时中心对象
- 模块分层和不一致命名
- 可能误导维护者的 3 个点
