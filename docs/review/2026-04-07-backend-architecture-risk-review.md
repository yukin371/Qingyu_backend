# 2026-04-07 Backend Architecture Risk Review

> 审查范围：`router`、`api/v1`、`service`、`models`、`docs/issues`、历史 `docs/review`  
> 审查目标：识别当前架构隐患，给出可执行结论，而非阶段性口号  
> 审查日期：2026-04-07

## 1. 执行结论

当前后端并非“分层失效”，而是处于“分层骨架存在 + 编排层膨胀 + 命名与契约漂移并存”的过渡态。  
最高风险不在单点 bug，而在结构可理解性和演进可控性：运行时入口过重、跨域能力边界不清、前后端契约仍有历史兼容负担。

### 风险分级总览

| 风险ID | 风险主题 | 等级 | 当前状态 |
|---|---|---|---|
| R-01 | `router/enter.go` 过重，承担运行时编排 | P0 | 持续存在 |
| R-02 | `shared` 职责黑洞风险 | P1 | 持续存在 |
| R-03 | `writer` 复合子域过大 | P1 | 持续存在 |
| R-04 | `search/stats/auth/audit` 非标准 vertical slice | P1 | 持续存在 |
| R-05 | 命名漂移与前后端契约风险 | P0 | 部分收敛，未闭环 |
| R-06 | 历史 issue 关闭与未关闭项混杂，治理断层风险 | P1 | 持续存在 |

---

## 2. 详细风险项

## R-01 `router/enter.go` 过重，承担运行时编排（P0）

**现象**

- `router/enter.go` 文件达到 1029 行，显著超出“路由注册层”常见规模。
- 同一文件中同时做了：
  - 服务可用性探测与分支注册
  - 搜索服务初始化（`initSearchService`）
  - 事件订阅（`project.published`）
  - 兼容路由保留与日志输出

**证据**

- `Qingyu_backend/router/enter.go:79` `RegisterRoutes`
- `Qingyu_backend/router/enter.go:191` `initSearchService(...)`
- `Qingyu_backend/router/enter.go:197` `Subscribe("project.published", ...)`
- `Qingyu_backend/router/enter.go:906` `func initSearchService(...)`

**影响**

- 路由层与 bootstrap 层边界被打散，问题定位成本上升。
- 启动成功不再代表关键能力完整注册，容易出现“部分可用”隐性故障。
- 对 AI 和新成员不友好，容易误解运行时真实装配路径。

**当前状态**

- 问题长期存在，尚无结构化拆分动作。

**建议动作**

1. 将运行时装配逻辑从 `router/enter.go` 抽离为独立 bootstrap/registrar 层。
2. 保留 `RegisterRoutes` 仅做路由绑定，不再承担服务初始化。
3. 增加“关键能力注册完成度”启动检查，避免 silent degrade。

---

## R-02 `shared` 职责黑洞风险（P1）

**现象**

- `service/shared` 承载了过多横切能力：`auth`、`cache`、`metrics`、`stats`、`storage`、`permission` 等。
- `shared` 既像基础设施层，又承接业务协作能力，边界不清。

**证据**

- 目录结构：`Qingyu_backend/service/shared/*`
- 历史审查已明确同类问题：`Qingyu_backend/docs/architecture/architecture_review_report.md`（“shared模块职责过重”）

**影响**

- 继续演进时容易把新需求“默认塞进 shared”，形成长期技术债。
- 模块复用与模块污染共存，测试范围扩大且隔离难度上升。

**当前状态**

- 未见成体系拆分计划，风险持续。

**建议动作**

1. 给 `shared` 增加明确准入规则：仅保留跨域基础能力，不承载业务编排。
2. 建立“shared 候选能力评审清单”，不满足条件则落入独立模块。
3. 对现有子域分批迁出：优先识别最接近业务域的组件。

---

## R-03 `writer` 复合子域过大（P1）

**现象**

- `writer` 在 API、router、service 层均表现为“模块内平台”：
  - `project`、`document`、`outline`、`story_harness`、`publish`、`stats`、`keyword`、`location`、`encyclopedia` 等并存。
- 单模块承担过多上下游协作职责。

**证据**

- 目录结构：`Qingyu_backend/router/writer/*`、`Qingyu_backend/api/v1/writer/*`、`Qingyu_backend/service/writer/*`
- 历史 issue 背景：`docs/issues/010-repository-layer-business-logic-permeation.md` 中 writer 域存在多处边界问题。

**影响**

- 变更影响面过大，难以做到低风险迭代。
- 单测和回归范围膨胀，CI 成本增加。
- 领域边界模糊导致重复能力和职责漂移。

**当前状态**

- 高复杂度现状持续，尚未完成子域治理。

**建议动作**

1. 将 `writer` 拆分为“创作核心子域 + 支撑子域”两级边界。
2. 为每个子域定义稳定契约，禁止跨子域直接调用内部实现。
3. 优先从 `publish`、`story_harness`、`stats` 三个变化频繁子域开始分离。

---

## R-04 `search/stats/auth/audit` 非标准 vertical slice（P1）

**现象**

- 这些模块在五层中的呈现不一致：
  - `search`: API/Service 有，路由无独立目录，初始化依赖 `router/enter.go`
  - `stats`: `api/v1/stats` 存在，但 service 分散于 `service/reader/stats` 与 `service/shared/stats`，路由是 `reading-stats`
  - `auth`/`audit`: API 与服务存在，但路由不独立，挂在 shared 或其他域

**证据**

- `Qingyu_backend/router/*` 与 `Qingyu_backend/api/v1/*` 目录对照
- `Qingyu_backend/service/*` 与 `Qingyu_backend/models/*` 目录对照
- `Qingyu_backend/router/enter.go` 中对搜索模块的特殊初始化路径

**影响**

- 架构理解成本高，模块职责依赖经验知识而非结构自解释。
- 新增功能容易沿用历史旁路，进一步加重结构不一致。

**当前状态**

- 持续存在，暂无统一对齐计划。

**建议动作**

1. 为 `search/stats/auth/audit` 各自定义“目标层级形态”。
2. 在文档层明确“当前临时形态”和“目标形态”，避免误判。
3. 先做命名和入口对齐，再做实现层重构。

---

## R-05 命名漂移与前后端契约风险（P0）

**现象**

- 目录和模块命名存在漂移：
  - `notification` vs `notifications`
  - `models/users` vs `service/user` / `api/v1/user`
  - `stats` vs `reading-stats`
- 历史契约兼容仍在持续：
  - `BookStatus` 历史值兼容
  - `categoryId` 与 `categoryIds` 双字段兼容
  - `updateTime/publishTime` 与 `updatedAt/publishedAt` 共存

**证据**

- `docs/issues/011-frontend-backend-data-type-inconsistency.md`（状态：部分存在问题）
- `docs/issues/005-api-standardization-issues.md`（Phase 1 完成，但仍有后续 TODO）

**影响**

- 前后端联调和自动化契约验证复杂度上升。
- 同义字段长期并存，导致数据语义歧义与消费侧误判。

**当前状态**

- 已有阶段性收敛，但仍未完成契约闭环。

**建议动作**

1. 建立“命名与字段契约冻结清单”，明确目标命名和弃用窗口。
2. 兼容字段全部打上退役计划，按版本逐步移除。
3. 将 API 标准化与类型统一纳入同一治理看板，禁止分散推进。

---

## R-06 历史 issue 关闭与未关闭项混杂，治理断层风险（P1）

**现象**

- `docs/issues/ISSUE_RELATIONSHIPS.md` 显示多个 archived issue 已关闭，但当前活跃 issue 仍有待处理与部分修复并存。
- 一些 issue 已写“已完成”，但正文仍保留“问题确认存在/后续 TODO”，状态语义不统一（例如 #010）。

**证据**

- `docs/issues/ISSUE_RELATIONSHIPS.md`
- `docs/issues/010-repository-layer-business-logic-permeation.md`
- `docs/issues/003-test-infrastructure-improvements.md`（部分修复）
- `docs/issues/004-code-quality-improvements.md`（待处理）
- `docs/issues/009-test-coverage-issues.md`（待处理）
- `docs/issues/012-service-layer-id-conversion-refactor.md`（待处理）

**影响**

- 管理层可能误判“风险已收敛”，实际仍有关键尾项未闭环。
- 执行优先级容易反复摇摆，影响长期架构治理节奏。

**当前状态**

- 文档已沉淀大量信息，但状态表达和执行闭环机制不足。

**建议动作**

1. 为活跃 issue 增加统一状态机：`open -> in_progress -> partial -> resolved -> archived`。
2. 对“已完成但仍有 TODO”的 issue 拆分后续子任务，避免主状态失真。
3. 在每次架构审查报告中同步“风险项与 issue 映射表”。

---

## 3. 历史 Issue 状态盘点（2026-04-07）

### 已归档且核心问题已解决

- #001 统一模型层 ID 字段类型
- #002 Repository Create 方法未回设 ID
- #006 数据库索引问题（核心）
- #007 Service 层事务管理缺失（核心）
- #008 中间件架构问题（核心）
- #013 测试用户种子数据 ID 未设置

### 部分解决，仍需继续跟踪

- #003 测试基础设施改进（部分修复）
- #005 API 标准化问题（Phase 1 完成，仍有后续补齐）
- #010 Repository 业务逻辑渗透（子项有完成，也有未收口）
- #011 前后端数据类型不一致（明确为部分存在）

### 尚未进入收敛阶段

- #004 代码质量改进（待处理）
- #009 测试覆盖率不足（待处理）
- #012 Service 层 ID 转换重构（待处理）

---

## 4. 建议优先级（下一轮）

1. **P0 首要跟踪**：R-01（入口编排过重）与 R-05（命名/契约漂移）。
2. **P1 结构治理**：R-02（shared 边界）与 R-04（非 vertical slice 模块对齐）。
3. **P1 执行治理**：R-06（issue 状态与任务闭环机制）。

---

## 5. 最值得优先跟踪的 3 个坑

1. `router/enter.go` 过重导致运行时边界混乱（R-01）。
2. 命名漂移与契约兼容长期并存导致对接风险持续（R-05）。
3. `shared` 作为职责黑洞继续膨胀（R-02）。
