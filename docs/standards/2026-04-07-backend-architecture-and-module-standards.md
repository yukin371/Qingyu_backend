# Backend Architecture And Module Standards (2026-04-07)

> 范围: `Qingyu_backend`  
> 角色: 项目级总则（不替代 `layer-*.md` 的分层细则）  
> 状态: 现行

## 1. 目标

本标准解决三个长期问题：

1. 分层规则存在，但模块边界和命名长期漂移。
2. `router/enter.go`、`service/shared` 等中心点持续膨胀。
3. 文档和代码容易出现“看起来分层清晰，实际运行链路复杂”的认知落差。

## 2. 分层边界总则

分层方向保持单向依赖：

`router -> api -> service -> repository -> models`

允许的基础设施例外：

- `service/container` 作为运行时装配中心
- `service/events` 作为事件总线能力
- `pkg/*` 作为工具层，不承载业务编排

禁止项：

1. API 直接访问 repository。
2. repository 直接实现业务规则（计算策略、流程编排、状态机）。
3. middleware 直接承载核心业务流程。

具体层内细则以 `layer-*.md` 为准。

## 3. 模块分组标准

模块必须按三组表达，而不是仅按目录平铺：

1. 业务域模块  
   例如 `bookstore / reader / writer / social / finance / admin / recommendation`
2. 横切能力模块  
   例如 `auth / audit / notification / announcements / search / stats`
3. 平台基础设施模块  
   例如 `service/container / service/events / service/shared / repository/querybuilder / repository/cache`

在架构文档、评审文档、设计文档中，必须优先使用这三组叙事。

## 4. 依赖规则

### 4.1 容器与入口

- `service/container` 是服务生命周期和依赖注入中心。
- `router/enter.go` 只应承担“路由聚合与装配编排”，不得继续吞并业务规则。

### 4.2 横切能力

- 横切能力（auth/audit/search/stats 等）优先通过 service 接口暴露，不直接嵌入 router 细节逻辑。
- `shared` 中新增能力必须先过边界评估，避免把 `shared` 当默认兜底目录。

### 4.3 事件和异步

- 事件订阅/发布逻辑可在入口装配，但业务语义要落在 service 层，不落在 router 初始化细节里。

## 5. 命名与对齐约束

以下命名漂移是当前明确问题，必须纳入治理基线：

1. `notification` vs `notifications`
2. `user` vs `users`
3. `stats` vs `reading-stats`

约束规则：

1. 一个业务概念只允许一个 canonical 名称，其他名称只作为兼容别名。
2. canonical 名称必须写入 API/Router/Service/Model 的映射表（在架构文档中维护）。
3. 新增路由或新模块禁止继续扩大命名漂移；如果历史兼容必须保留，文档中要显式标注“兼容路径”。

## 6. 当前重点治理项

### 6.1 `router/enter.go` 过重

标准要求：

1. `router/enter.go` 保持“注册与装配”定位。
2. 搜索、事件订阅、兼容策略等必须逐步沉淀为可复用的初始化组件，而不是在 `enter.go` 增长内联逻辑。

### 6.2 `service/shared` 易膨胀

标准要求：

1. 新能力进入 `shared` 前，必须回答“为何不是独立模块”。
2. `shared` 下模块需要按职责分区（auth/cache/metrics/storage 等），禁止无分类堆叠。
3. 当某子能力出现独立生命周期或复杂策略时，应拆出独立 service 包。

### 6.3 `writer` 复合子域

标准要求：

1. `writer` 视为复合域，不再按“单一模块”治理。
2. 文档中必须显式区分其子域能力：project/document/storyharness/publish/stats 等。

## 7. 文档维护规则

### 7.1 文档角色分离

- `docs/standards/*`: 规则与约束
- `docs/architecture/*`: 现状结构和运行链路
- `docs/review/*`: 风险与审查结论
- `docs/issues/*`: 持续跟踪的问题项

### 7.2 更新触发条件

以下变化发生时，必须至少更新一个 standards 文档：

1. 分层边界变更
2. 模块分组或命名变更
3. 入口装配方式变化（尤其 `router/enter.go` 与容器初始化）
4. 新增横切能力或重大兼容路径

### 7.3 路径规范

文档统一使用 `docs/standards/...` 路径，不再使用旧写法 `doc/standards/...`。

## 8. 与 `layer-*.md` 的关系

本文件只定义项目级约束，不重写分层细则。  
具体实现规范继续由以下文档承接：

- `layer-api.md`
- `layer-router.md`
- `layer-service.md`
- `layer-repository.md`
- `layer-models.md`
- `layer-dto.md`
- `layer-middleware.md`
- `layer-config.md`
- `layer-pkg.md`

## 9. 执行检查清单

- [ ] 新增模块是否标注了三组归属（业务域/横切/平台）
- [ ] 是否引入了新的命名漂移
- [ ] 是否扩大了 `router/enter.go` 的职责
- [ ] 是否把本应独立的能力继续塞进 `service/shared`
- [ ] 是否同步更新了对应 `layer-*.md` 或关联架构/审查文档
