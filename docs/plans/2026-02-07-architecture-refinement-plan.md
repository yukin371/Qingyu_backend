# Architecture Refinement Implementation Plan (Revised v2)

**创建日期**: 2026-02-07  
**修订日期**: 2026-02-07  
**目标**: 在不破坏现有可运行路径的前提下，分阶段修复架构问题，优先解决可编译性、可观测性、可回滚性。

---

## 0. 修订原则

1. 先兼容再替换：新增抽象层时保留旧构造器/旧调用路径，完成全量迁移后再删除旧接口。
2. 以现有真实接口为准，不引入“平行事件系统”。
3. 每个阶段必须满足：
   - 可编译
   - 测试可执行
   - 可回滚（单独 PR）
4. 启动验证统一使用当前入口：`cmd/server/main.go`。

---

## 阶段总览

| Phase | 问题 | 优先级 | 预计工期 | 交付形式 |
|---|---|---|---|---|
| 1 | 中间件跨层依赖 | P0 | 1-2 天 | 小步兼容重构 |
| 2 | 服务初始化顺序隐式依赖 | P1 | 2-3 天 | 显式初始化计划 + 依赖校验 |
| 3 | shared 模块职责过重 | P1 | 3-5 天 | 渐进拆分（先接口再迁移） |
| 4 | Writer 服务耦合度高 | P2 | 3-4 天 | 基于现有 EventBus 的事件化 |
| 5 | 事件总线持久化治理 | P2 | 2-3 天 | 复用已有 persisted bus 能力 |

---

## Phase 1: 中间件跨层依赖 (P0)

### 目标
让 `pkg/middleware/quota.go` 依赖抽象接口，而不是直接依赖 `*ai.QuotaService`，但保持现有路由不改即可运行。

### 现状约束（必须对齐）
- 现有中间件直接依赖 `*aiService.QuotaService`（`pkg/middleware/quota.go`）。
- 现有 `QuotaService` 方法签名：
  - `CheckQuota(ctx, userID, amount)`
  - `ConsumeQuota(ctx, userID, amount, service, model, requestID)`
  - `GetQuotaInfo(ctx, userID)`（不是 `GetUserQuota`）

### Task 1.1 新增配额抽象接口（与现有签名一致）

**Files**
- Create: `pkg/quota/checker.go`

**Implementation**
- 定义 `Checker`：
  - `Check(ctx context.Context, userID string, estimatedAmount int) error`
  - `Consume(ctx context.Context, userID string, amount int, service, model, requestID string) error`
  - `GetInfo(ctx context.Context, userID string) (*ai.UserQuota, error)`
- 错误不重复定义，沿用 `models/ai` 中已有错误（避免双重错误域）。

### Task 1.2 实现 QuotaService 适配器

**Files**
- Create: `service/ai/quota_checker_adapter.go`
- Create: `service/ai/quota_checker_adapter_test.go`

**Implementation**
- `type QuotaCheckerAdapter struct { quotaService *QuotaService }`
- `Check` -> `quotaService.CheckQuota`
- `Consume` -> `quotaService.ConsumeQuota`
- `GetInfo` -> `quotaService.GetQuotaInfo`
- 测试使用可注入接口或 stub，不把 `*MockQuotaService` 强行当 `*QuotaService` 传入。

### Task 1.3 重构中间件（保留兼容入口）

**Files**
- Modify: `pkg/middleware/quota.go`

**Implementation**
1. 新增接口构造器：`NewQuotaMiddlewareWithChecker(checker quota.Checker)`.
2. 保留旧构造器：
   - `NewQuotaMiddleware(quotaService *aiService.QuotaService)` 内部创建 adapter。
3. 兼容现有路由函数签名：
   - `QuotaCheckMiddleware(quotaService *aiService.QuotaService)`
   - `LightQuotaCheckMiddleware(quotaService *aiService.QuotaService)`
   - `HeavyQuotaCheckMiddleware(quotaService *aiService.QuotaService)`
4. 中间件内部改为统一走 `checker.Check(...)`。

### Task 1.4 验证

**Commands**
```bash
go test ./pkg/middleware/... ./service/ai/... -v
go run cmd/server/main.go
```

### Phase 1 验收标准
- [ ] 中间件核心逻辑不再直接调用 `QuotaService` 方法
- [ ] `router/ai/ai_router.go` 不改动也可正常编译运行
- [ ] 单测与集成测试通过

---

## Phase 2: 服务初始化顺序隐式依赖 (P1)

### 目标
把“隐式依赖 + 手工顺序”变成“显式依赖 + 可验证顺序”，避免一次性重写整个容器。

### 现状约束
- `ServiceContainer.Initialize()` 与 `SetupDefaultServices()` 两段式初始化（`service/enter.go`）。
- `SetupDefaultServices()` 内存在真实顺序依赖（例如 `ProjectService` 在 `AIService` 前）。
- 不宜一次性替换为全量工厂注册系统（风险过大）。

### Task 2.1 新增轻量初始化计划结构

**Files**
- Create: `service/container/init_plan.go`
- Create: `service/container/init_plan_test.go`

**Implementation**
- 定义 `InitStep`：
  - `Name string`
  - `DependsOn []string`
  - `Required bool`
  - `Build func(*ServiceContainer) error`
- 提供：
  - `ValidatePlan(steps []InitStep) error`（依赖存在、循环依赖检测）
  - `ResolveOrder(steps []InitStep) ([]InitStep, error)`（稳定拓扑排序）
- 不使用“优先级分组后再拓扑”的算法，避免误判循环。

### Task 2.2 在 SetupDefaultServices 中接入初始化计划

**Files**
- Modify: `service/container/service_container.go`

**Implementation**
1. 把现有大段初始化拆成多个 `BuildXxx` 私有函数。
2. 在 `SetupDefaultServices()` 构建 `[]InitStep`，先验证、再排序、再执行。
3. `Required=false` 的步骤失败仅告警，不阻塞（用于可选服务）。
4. 增加初始化日志：输出最终执行顺序与失败项。

### Task 2.3 启动时依赖图导出（仅调试用途）

**Files**
- Modify: `service/container/service_container.go`

**Implementation**
- 增加 `DumpInitPlan() string`，用于测试与排障。
- 暂不在 `cmd/server/main.go` 增加新 flag，先通过日志与单测验证。

### Task 2.4 验证

**Commands**
```bash
go test ./service/container/... -v
go run cmd/server/main.go
```

### Phase 2 验收标准
- [ ] 初始化顺序由可验证的数据结构驱动
- [ ] 可检测循环依赖和缺失依赖
- [ ] 不破坏现有 `service.InitializeServices()` 启动链路

---

## Phase 3: shared 模块职责过重 (P1)

### 目标
渐进拆分 `service/shared/*`，先解决边界和依赖，再做目录迁移，避免大规模 `git mv` 引发冲突。

### 现状约束
- `service/shared/auth`、`service/shared/storage`、`service/shared/messaging` 被大量引用。
- 仓库已存在 `service/messaging`，不可直接把 `service/shared/messaging` 强行迁入同路径。

### Task 3.1 先做“包边界收口”，不移动目录

**Files**
- Create: `service/interfaces/shared/{auth.go,storage.go,messaging.go}`
- Modify: `service/container/service_container.go`

**Implementation**
1. 提炼最小接口到 `service/interfaces/shared/*`。
2. 容器字段尽量依赖接口而非具体实现。
3. API 层逐步改为依赖这些接口包，减少直接耦合到 `service/shared/*` 具体实现。

### Task 3.2 按模块分批迁移（每批单独 PR）

**Batch A（低风险）**
- `service/shared/cache` -> `service/cache`
- 更新导入并回归测试

**Batch B（中风险）**
- `service/shared/stats` -> `service/metrics`（或 `service/stats`，二选一并统一）

**Batch C（高风险，最后做）**
- `service/shared/auth`、`service/shared/storage`、`service/shared/messaging`
- 先建立新包 re-export/适配层，再迁移调用，最后删旧路径

### Task 3.3 删除旧目录（仅在全量引用清零后）

**硬门槛**
```bash
rg -n "Qingyu_backend/service/shared/" -g "*.go"
```
结果必须为 0，才允许删旧目录。

### Task 3.4 验证

**Commands**
```bash
go test ./... -run TestNonExistent -count=0
go test ./api/... ./router/... ./service/... -v
go run cmd/server/main.go
```

### Phase 3 验收标准
- [ ] 模块依赖先收口到接口
- [ ] 迁移按批次完成且每批可回滚
- [ ] 不出现与现有 `service/messaging` 的目录冲突

---

## Phase 4: Writer 服务耦合度优化 (P2)

### 目标
在不改写全局事件模型的前提下，让 Writer 相关跨模块动作通过现有 `EventBus` 事件驱动。

### 现状约束（必须对齐）
- 项目已有标准 `Event` / `EventHandler` / `EventBus` 接口（`service/interfaces/base/base_service.go`）。
- 项目已有 `service/events/writer_events.go`，应优先扩展而不是重复建事件体系。

### Task 4.1 复用并扩展现有 writer 事件定义

**Files**
- Modify: `service/events/writer_events.go`

**Implementation**
- 新增确有业务价值的事件（例如 `ChapterPublishedEvent`），实现 `base.Event`。
- 事件字段使用稳定 schema，避免 `map[string]interface{}`。

### Task 4.2 新增标准事件处理器（实现 EventHandler 接口）

**Files**
- Create: `service/events/handlers/writer_notification_handler.go`
- Create: `service/events/handlers/writer_finance_handler.go`

**Implementation**
- 每个 handler 实现：
  - `Handle(ctx context.Context, event base.Event) error`
  - `GetHandlerName() string`
  - `GetSupportedEventTypes() []string`
- 不再使用 `Subscribe(func(interface{}) error)` 这种非标准签名。

### Task 4.3 在容器注册订阅

**Files**
- Modify: `service/container/service_container.go`

**Implementation**
- 在 `SetupDefaultServices()` 末尾注册 writer 相关 handlers 到 `c.eventBus`。
- 统一处理订阅失败日志。

### Task 4.4 验证

**Commands**
```bash
go test ./service/events/... ./service/writer/... -v
go run cmd/server/main.go
```

### Phase 4 验收标准
- [ ] Writer 相关跨模块流程通过标准 EventBus 驱动
- [ ] 事件定义和处理器都符合现有接口约束
- [ ] 无新增“第二套事件协议”

---

## Phase 5: 事件总线持久化治理 (P2)

### 目标
基于现有 `service/events/persisted_event_bus.go` 增强可观测与运维能力，而不是重复造轮子。

### 现状约束
- 已有 `PersistedEventBus`、`RetryableEventBus`、`EventStore`。
- 新功能应围绕已有能力补齐（重放、管理 API、指标）。

### Task 5.1 补齐 replay 能力到现有 persisted bus 体系

**Files**
- Modify: `service/events/event_store.go`
- Modify: `service/events/mongo_event_store.go`
- Modify: `service/events/persisted_event_bus.go`

**Implementation**
- 在 `EventStore` 增加可分页读取历史事件接口。
- 在 `PersistedEventBus` 增加 `Replay(ctx, filter)`（按类型、时间窗、limit）。
- Replay 期间支持 dry-run 模式（只统计不分发）。

### Task 5.2 管理 API（仅 admin）

**Files**
- Create: `api/v1/admin/events_api.go`
- Modify: `router/...`（admin 路由接入）

**Implementation**
- `POST /admin/events/replay` 使用 JSON body：
  - `event_type`
  - `from`
  - `to`
  - `limit`
  - `dry_run`
- 参数标签使用 `json`，与 `ShouldBindJSON` 对齐。

### Task 5.3 可观测性

**Files**
- Modify: `service/events/metrics.go`
- Modify: `service/events/logging.go`

**Implementation**
- 增加指标：
  - replay_count
  - replay_failed_count
  - replay_duration_ms
- 关键日志包含 trace id / event type / range / result。

### Task 5.4 验证

**Commands**
```bash
go test ./service/events/... ./api/v1/admin/... -v
go run cmd/server/main.go
```

### Phase 5 验收标准
- [ ] 持久化事件支持可控 replay
- [ ] 管理 API 可用且参数校验正确
- [ ] replay 过程可观测

---

## 统一提交与回滚策略

1. 每个 Task 独立提交，禁止跨 Phase 混合提交。
2. 每个 Phase 至少包含：
   - 一个测试提交
   - 一个文档更新提交
3. 回滚策略：
   - 若线上风险高，优先回滚“接入点改动”提交
   - 保留抽象层代码（不影响运行）

---

## 统一验证清单

### 功能验证
- [ ] 关键 API 可正常调用（auth / ai / writer / bookstore）
- [ ] 事件发布与处理链路可观测

### 质量验证
- [ ] 相关模块单测通过
- [ ] 至少一次完整服务启动验证通过（`go run cmd/server/main.go`）
- [ ] 无新增编译警告和循环依赖

### 文档验证
- [ ] 更新架构文档（模块边界、事件流）
- [ ] 更新迁移说明（尤其 Phase 3）

---

## 本版相对 v1 的关键修正

1. 对齐现有 `QuotaService`/`EventBus` 真实接口，移除不可编译示例。
2. Phase 2 改为轻量可验证方案，避免一次性替换整个容器体系。
3. Phase 3 改为“先收口再迁移”，规避 `service/messaging` 路径冲突。
4. Phase 5 改为增强现有 `persisted_event_bus`，不重复建设并行实现。

