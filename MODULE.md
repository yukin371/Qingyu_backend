# Qingyu_backend 模块上下文

## 模块职责
- 提供 Go + Gin 实现的 API 层，聚合 10 大业务面向的路由与请求封装。
- 维护 Service 层业务逻辑与 Repository 层数据访问之间的事务与一致性。
- 作为主业务真相的主要持有者，负责与 MongoDB/Redis/Milvus 等核心数据组件的交互并保持一致的错误/响应语义。
- 支撑中台能力（配置、鉴权、日志、依赖注入容器）并为 Python AI 服务、前端与子系统提供统一的边界。

## Owns
- `router/enter.go` 及所有公共路由注册、中间件链（鉴权、限流、日志）的编排。
- `service/` 中的业务实现，确保多模块复用时的依赖准确；`repository/` 中的存储接口与缓存协作。
- `core/` 与 `container/` 负责启动顺序、依赖注入与共享能力。
- 构建跨模块命名空间（如 error code、response schema）与共享错误处理、统一日志格式。

## Must not own
- 书籍/阅读等业务主视图渲染，应留给前端；仅提供可消费的 API，无直接 UI 逻辑。
- AI 运行时的细节（LangChain/LangGraph Agent）由 `Qingyu-Ai-Service` 负责，后端只负责 gRPC/HTTP 网关。
- 不在子模块 `docs/plans/` 中写跨仓库治理内容；这些计划应在父仓库 `docs/plans/submodules/backend/` 中归档、收敛。

## 关键依赖
- MongoDB/Redis/Milvus 等持久化组件与 `config` 中的配置管理。
- 描述到 AI 的 gRPC 通道（proto、`service/ai` 适配器）和 Python 服务的 health-check。
- `service/shared` 等复用层，任何修改需先确认不会造成能力重复或命名漂移。
- 与前端 Editor V3 相关的 route 路径与 DTO 约定（`api/v1/writer` 等）。

## 不变量
- API 响应 `code/message/data` 结构必须一致，错误码在 `docs/standards/error_code_standard.md` 中维护。
- `router/enter.go` 的中间件顺序（如鉴权 > RBAC > 日志）不得随意调换，任何改动需要同步文档。
- 共享服务（如 `service/shared`）必须保持单一 owner，不允许不同模块私自 fork 实现。
- 命名约定不宜出现 `notification/notifications` / `user/users` 等重复层级；若需要调整需先列出命名影响范围。

## 常见坑
1. `router/enter.go` 责任过重，增加新的全局中间件前必须评估对现有路由的影响并补充测试。
2. `service/shared` 是职责黑洞，复用前必须确认现有实现满足需求，否则会造成隐性依赖循环。
3. 命名漂移（`notification/notifications`、`user/users`、`stats/reading-stats` 等）依然存在，改名需同步 API 文档与前端 contract。
4. 缺乏明确的 ownership 时容易重复实现 shared helper，改动前必须回答 “本实现是否已经存在” 这类强制问题。

## 文档同步触发条件
- 新增/调整 API 时同步更新 `docs/api`、`docs/architecture` 与 swagger 产物。
- 任何涉及 `router/enter.go`、`service/shared`、命名空间调整的改动需在 `docs/standards` 与 `docs/review` 下写 review/risk 说明。
- 关键模块重构（例如 writer、reader、ai）的边界或验证方式变化时更新 `Qingyu_backend/MODULE.md` 并在父仓库 `docs/plans/submodules/backend/` 创建 plan 梳理。
