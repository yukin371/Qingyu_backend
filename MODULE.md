# Qingyu_backend 模块上下文

## 模块职责
- 提供 Go + Gin 实现的 API 层，聚合 10 大业务面向的路由与请求封装。
- 维护 Service 层业务逻辑与 Repository 层数据访问之间的事务与一致性。
- 作为主业务真相的主要持有者，负责与 MongoDB/Redis/Milvus 等核心数据组件的交互并保持一致的错误/响应语义。
- 支撑中台能力（配置、鉴权、日志、依赖注入容器）并为 Python AI 服务、前端与子系统提供统一的边界。

> 当前目录仍处于历史结构与目标模块结构并存的迁移期。目标边界以父仓库
> `docs/architecture/目标架构与模块边界.md` 和 `docs/standards/开发治理与变更流程.md` 为准；本文件描述当前 owner 和迁移注意事项，不代表所有目录已完成收敛。

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
- 运行时主链路不得使用临时 stdout/debug dump 作为诊断手段；生产可见诊断应走既有 logger / middleware 链路，并在风险清理计划中记录例外。

## 常见坑
1. `router/enter.go` 责任过重，增加新的全局中间件前必须评估对现有路由的影响并补充测试。
2. `service/shared` 是职责黑洞，复用前必须确认现有实现满足需求，否则会造成隐性依赖循环。
3. 命名漂移（`notification/notifications`、`user/users`、`stats/reading-stats` 等）依然存在，改名需同步 API 文档与前端 contract。
4. 缺乏明确的 ownership 时容易重复实现 shared helper，改动前必须回答 “本实现是否已经存在” 这类强制问题。
5. Phase 5 已完成 writer factory / compat / port adapter 核心退场，删除前备份归档在 `backup/writer-compat-predelete-20260430` 与 `backup/ai-adapter-predelete-20260503` 两个后端 tag；后续不要把这些旧文件名当作当前运行态 owner。
6. `pkg/sync` 的阅读进度同步只依赖 `ReaderProgressPort`，不得重新直接导入 `service/reader`；本次迁移存档为 `migration/phase3-t1-progress-port-20261005`。
7. 邮箱验证码管理器已归入 `service/channels`，认证 API 只消费 `channels.Manager`；不得恢复 `pkg/emailcode` 兼容入口。本次迁移存档为 `migration/phase3-t2-emailcode-predelete-20261005`。
8. writer Mongo 仓储直接使用 `models/dto` 的 `ExportTask`、`PublicationRecord`；`service/interfaces` 只保留导出任务兼容别名，仓储不得重新导入 service 包。本次迁移存档为 `migration/phase3-t2-writer-repository-dto-20261005`。
9. admin 内容导出 handler 只依赖 `service/interfaces.AdminContentBookPort` 与 `AdminContentChapterPort`；repository 的 `infrastructure.Filter` 只能在组合根 adapter 中映射，不能重新泄漏到 `api/v1/admin` 或 service port。生产路由由 `router/admin.RegisterContentExportRoutes` 独立注册，并沿用 JWT + `admin` 角色边界；Mongo 或读取端口不可用时跳过该组路由。
10. `internal/middleware/admin` 只依赖本包的 `PermissionChecker` 最小 port（`HasRole`、`CheckPermission`）；不要重新导入 `service/auth.PermissionService`，也不要把缓存/RBAC 装载方法泄漏到 HTTP 中间件边界。
11. JWT 签名、解析、验证、刷新和 wire-level claims 的唯一 owner 是 `internal/security/jwt`；`internal/middleware/auth` 只保留迁移期兼容 façade，`service/auth` 新代码直接依赖 `security/jwt.Manager`。
12. `internal/middleware/setup.RegisterAuthFactory` 直接消费 `internal/security/jwt.Manager`；不得把 JWT 类型重新从 `internal/middleware/auth` façade 暴露给组合根。
13. finance 路由的生产注册通过组合根注入 `RegisterFinanceRoutesWithAuth`；`RegisterFinanceRoutes` 与 `auth.JWTAuth` 仅作迁移期兼容入口。后续路由迁移应复用 `internal/security/jwt.Manager`，不得新增直接构造 middleware façade 的生产调用。
14. notifications 路由的生产注册通过组合根注入 `RegisterRoutesWithAuth` 与 `RegisterUserManagementRoutesWithAuth`；旧注册函数仅作迁移期兼容入口。通知路由迁移不得重新直接构造 `auth.JWTAuth()`。
15. reading-stats 路由的生产注册通过组合根注入 `RegisterReadingStatsRoutesWithAuth`；旧 `RegisterReadingStatsRoutes` 仅作迁移期兼容入口，统计 API 不得重新直接构造 JWT façade。
16. recommendation 路由的认证与管理员子组通过组合根注入 `RegisterRecommendationRoutesWithAuth`；公开推荐路径不应引入鉴权，旧 `RegisterRecommendationRoutes` 仅作迁移期兼容入口。
17. community 路由的认证写操作通过组合根注入 `RegisterCommunityRoutesWithAuth`；公开动态读取路径不得误加鉴权，旧 `RegisterCommunityRoutes` 仅作迁移期兼容入口。
18. search 只有灰度管理组需要鉴权，生产注册通过组合根注入 `RegisterSearchRoutesWithAuth`；公开搜索和健康路由不得引入 JWT，旧 `RegisterSearchRoutes` 仅作迁移期兼容入口。
19. admin 主路由和 content export 路由通过组合根注入 `RegisterAdminRoutesWithAuth`/`RegisterContentExportRoutesWithAuth`；必须保持 JWT 认证先于 `RequireRole("admin")`，旧注册函数仅作迁移期兼容入口。
20. AI 主路由和 creative 路由通过组合根注入 `InitAIRouterWithAuth`；认证必须保持在 quota middleware 之前，旧 `InitAIRouter`/`InitCreativeRoutes` 仅作迁移期兼容入口。

## 文档同步触发条件
- 新增/调整 API 时同步更新 `docs/api`、`docs/architecture` 与 swagger 产物。
- 任何涉及 `router/enter.go`、`service/shared`、命名空间调整的改动需在 `docs/standards` 与 `docs/review` 下写 review/risk 说明。
- 关键模块重构（例如 writer、reader、ai）的边界、备份归档或验证方式变化时更新 `Qingyu_backend/MODULE.md` 并在父仓库 `docs/plans/submodules/backend/` 创建或回写 plan 梳理。
- 后端新增业务能力必须明确所属领域和 handler -> service -> repository 依赖方向；不得以新增顶层 shared helper 规避 owner 决策。
