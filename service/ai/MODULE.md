# service/ai

> 最后更新：2026-04-17

## 职责

承接 Qingyu 的 AI 领域服务编排，当前除通用 gRPC/聊天能力外，还负责 AI 配额治理闭环。
本轮新增的 quota 子能力覆盖管理员批量操作、策略管理、预警查询、仪表盘聚合与定时刷新，不负责直接定义前端展示结构。

## 数据流

管理员路由启动时：
`router.RegisterRoutes` 创建 Mongo quota/policy/alert repositories，再组装 `QuotaAdminService`、`QuotaDashboardService`、`QuotaPolicyService`、`QuotaAlertService`，最后注入 `router/admin.RegisterAdminRoutes`。

运行时主链：
`QuotaService` 负责用户配额读写与消费事务，`QuotaPolicyService` 提供默认配额策略，`QuotaDashboardService` 从 Mongo 聚合并写 Redis 缓存，`QuotaAlertService` 负责告警状态流转，`QuotaScheduler` 周期执行异常检测与 dashboard 缓存刷新。

## 约定 & 陷阱

- `QuotaService` 仍是 AI 配额的主 owner；管理员相关 service 只做后台管理能力，不应绕开 `QuotaService` 直接改用户消费逻辑。
- `QuotaDashboardService` 使用 Redis key `quota:dashboard`；手动刷新和定时刷新都走 `RefreshDashboardCache`，不要在别处重复写缓存键。
- `QuotaAnomalyDetector` 读取 Redis 速率键时依赖原始 `*redis.Client`，启动期必须通过 `cache.RedisClient.GetClient()` 传入，而不是直接把包装接口强转。
- 跨服务一致性检测当前只有调度入口，没有可用的 AI 侧 quota 统计 API；实现前应先补 `Phase3Client` 或 AI gRPC 协议的对账查询能力。
- 若后续新增 quota 子模块文档，`MODULE.md` 只保留业务边界和陷阱，不复制 README 的接口清单。
