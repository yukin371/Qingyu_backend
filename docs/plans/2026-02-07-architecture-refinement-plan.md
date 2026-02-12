# Architecture Refinement Implementation Plan (Revised v3)

**创建日期**: 2026-02-07  
**修订日期**: 2026-02-12  
**目标**: 在不破坏现有可运行路径的前提下，分阶段修复架构问题，优先解决可编译性、可观测性、可回滚性，并显著降低开发/维护成本。

---

## 0. 修订原则

1. 先兼容再替换：新增抽象层时保留旧构造器/旧调用路径，完成全量迁移后再删除旧接口。
2. 以现有真实接口为准，不引入“平行事件系统”。
3. 每个阶段必须满足：
   - 可编译
   - 测试可执行
   - 可回滚（单独 PR）
4. 启动验证统一使用当前入口：`cmd/server/main.go`。

### 0.1 方向修正（2026-02-12）

1. **冻结目录级重构**：暂停 `service/shared/*` 大规模迁移和 `git mv`，先解决接口一致性与编译稳定性。  
2. **统一最小契约**：`quota.Checker` 固定为单一签名，禁止同阶段内变更返回语义。  
3. **先减法再加法**：优先删除重复入口、重复抽象、重复适配器；新增组件必须有明确收益与替换路径。  
4. **测试与实现同提交闭环**：接口变更必须同 PR 内完成 mock / factory / compile gate 修复。  

---

## 阶段总览（v3）

| Phase | 问题 | 优先级 | 预计工期 | 交付形式 |
|---|---|---|---|---|
| 1 | 配额中间件契约统一 | P0 | 0.5-1 天 | 固定接口 + 删除重复路径 |
| 2 | 编译与测试护栏 | P0 | 0.5-1 天 | writer/reader/quota 关键包门禁 |
| 3 | shared 模块降复杂度 | P1 | 2-3 天 | 不迁目录，只做边界收口与 TODO 清偿 |
| 4 | 容器初始化优化（轻量） | P2 | 1-2 天 | 只补验证与日志，不重写初始化框架 |
| 5 | 事件治理增强（可选） | P2 | 2-3 天 | 在现有 persisted bus 上增量增强 |

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

### Task 1.1 固定配额抽象接口（禁止继续漂移）

**Files**
- Create: `pkg/quota/checker.go`

**Implementation**
- 固定 `Checker` 最小契约：`Check(ctx context.Context, userID string, estimatedAmount int) error`。
- 当前阶段不在 `Checker` 中扩展 `Consume/GetInfo`，避免接口过宽导致 mock 与调用点同步成本过高。
- 错误类型统一走 `pkg/quota`，避免双重错误域。

### Task 1.2 精简适配层

**Implementation**
- 若 `QuotaService` 已直接实现 `Checker`，删除冗余 adapter；若未实现，仅保留单一 adapter 路径。
- 清理重复入口函数，避免“同功能多构造器”。

### Task 1.3 重构中间件（兼容但不扩散）

**Files**
- Modify: `pkg/middleware/quota.go`

**Implementation**
1. 新增接口构造器：`NewQuotaMiddlewareWithChecker(checker quota.Checker)`.
2. 保留旧构造器：
   - `NewQuotaMiddleware(quotaService *aiService.QuotaService)` 内部创建 adapter。
3. 兼容现有路由函数签名，但内部必须统一走 `checker.Check(...)`。
4. 不新增新的中间件变体函数；轻量/标准差异通过参数表达。

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

## Phase 2: 编译与测试护栏 (P0)

### 目标
用最小成本建立“防回归护栏”，减少维护人力消耗和返工。

### Task 2.1 关键包门禁

**Commands**
```bash
go test ./service/writer ./service/reader ./pkg/middleware ./pkg/quota -v
```

**Rule**
- 作为每次架构改动前置门禁；未通过不得推进下一步重构。

### Task 2.2 接口变更闭环规则

**Rule**
- 任何 Port 签名变更必须同提交更新对应 mock 与 factory test。
- 禁止“先改接口，后补测试”的跨提交做法。

### Task 2.3 验证

**Commands**
```bash
go test ./service/writer ./service/reader ./pkg/middleware ./pkg/quota -v
go run cmd/server/main.go
```

### Phase 2 验收标准
- [ ] 关键包门禁稳定通过
- [ ] 接口变更不再引入跨包编译断裂
- [ ] 不破坏现有启动链路

---

## Phase 3: shared 模块降复杂度 (P1)

### 目标
在**不迁移目录**前提下，降低 shared 模块维护复杂度。

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

### Task 3.2 暂停目录迁移（冻结策略）

**Decision**
- 冻结 `service/shared/* -> service/*` 迁移。
- 冻结新增多云 adapter（S3/OSS/COS）与大规模策略框架扩展。
- 优先清偿 `stats_service.go` 与 `storage_service.go` 中已识别 TODO。

### Task 3.3 验证

**Commands**
```bash
go test ./service/shared/... -v
go test ./api/... ./router/... ./service/... -run TestNonExistent -count=0
go run cmd/server/main.go
```

### Phase 3 验收标准
- [ ] 在不迁移目录的前提下，shared 关键模块编译与测试稳定
- [ ] TODO 项按优先级持续下降
- [ ] 无新增抽象层重复建设

---

## Phase 4: Writer 服务耦合度优化 (P2)

> v3 说明：本阶段延后，只有在 Phase 1-3 连续稳定后再启动。

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

> v3 说明：本阶段延后，禁止并行启动以免分散维护资源。

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

1. 固定 `Checker` 最小契约，停止接口语义漂移。  
2. Phase 2 调整为“编译与测试护栏”，优先降低返工成本。  
3. Phase 3 调整为“冻结目录迁移 + 清偿 TODO”，避免低收益大重构。  
4. Phase 4/5 改为延后启动，集中资源先完成可编译、可测试、可维护。  
