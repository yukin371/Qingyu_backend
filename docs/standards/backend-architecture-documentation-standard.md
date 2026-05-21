# 后端架构与文档标准

> 更新日期: 2026-04-07
> 适用范围: `Qingyu_backend` 的架构总览、模块文档、Mermaid 图、审查报告

## 1. 标准目标

这份标准解决两个问题：

1. 如何让人类开发者快速建立后端上下文
2. 如何让 AI 在最少输入下得到接近真实运行结构的理解

## 2. Source of Truth 规则

### 2.1 当前一线主入口

以下文档优先作为当前后端 source of truth：

- `docs/architecture/README.md`
- `docs/architecture/2026-04-07-backend-runtime-architecture.md`
- `docs/architecture/2026-04-07-backend-module-map.md`
- `docs/guides/backend-context-quickstart.md`
- `docs/review/2026-04-07-backend-architecture-risk-review.md`

### 2.2 二线参考文档

这些文档可引用，但不应单独承担“总入口”职责：

- `docs/architecture/system_architecture.md`
- `docs/architecture/data_model.md`
- `docs/architecture/api_architecture.md`
- `docs/standards/layer-*.md`
- `docs/issues/*.md`

### 2.3 历史文档处理规则

以下类型默认视为历史或阶段性资料：

- 单次阶段交接
- MVP 审查阶段产物
- 特定模块迁移分析
- 已经过期的 TODO / completion / handoff

历史资料可以保留，但不能替代当前架构总览。

## 3. 架构文档最小集合

后端至少要维护这 4 类文档：

1. **运行时架构**
   - 启动链路
   - 中间件顺序
   - 请求主路径
   - 全局状态 / 容器
2. **模块地图**
   - 业务模块
   - 基础设施模块
   - 层级对齐情况
   - 特殊命名和包袱
3. **标准**
   - 文档边界
   - 命名和 Mermaid 规则
   - source of truth 规则
4. **审查**
   - 架构坑点
   - 反直觉点
   - 高风险区域

## 4. Mermaid 图标准

### 4.1 必须画的图

每个主架构包至少维护下面 3 类 Mermaid 图中的 2 类：

- 启动 / bootstrap flow
- 请求流 / runtime request flow
- 模块依赖图

### 4.2 Mermaid 图要求

1. 节点名称应尽量映射真实目录或真实对象名。
2. 必须显式标出运行时中心对象，例如 `ServiceContainer`。
3. 如果存在兼容路径、fallback 或 optional dependency，图中要标出来。
4. 不要只画理想分层图，必须允许出现过渡态结构。

### 4.3 应避免的图

- 只展示“理论层级”但完全不反映真实启动过程
- 省略 `router` 与 `api/v1` 的双层职责
- 省略容器、事件总线、搜索初始化等运行时关键节点

## 5. 模块文档标准

如果为某个后端模块建立说明，至少应包含：

1. 模块职责
2. 对应目录
3. 上游入口
4. 下游依赖
5. 关键模型 / 关键服务
6. 已知坑点

推荐结构：

```text
模块名
- 职责
- 入口: router/api
- 核心服务: service/*
- 数据入口: repository/*
- 关键模型: models/*
- 风险与例外
```

## 6. 评审与问题文档标准

### 6.1 `review/`

适合放：

- 阶段性审查
- 架构风险报告
- 文档与代码不一致盘点

### 6.2 `issues/`

适合放：

- 长期存在的结构问题
- 可被继续跟踪和拆解的整改项

### 6.3 关系

- `review` 负责“这一轮发现了什么”
- `issues` 负责“这个问题如何持续跟踪”

## 7. 命名与路径约定

1. 新总览类文件优先使用英文文件名或日期前缀英文文件名。
2. README 只负责索引和边界，不承担完整专题正文。
3. 文档里的路径示例优先使用当前真实 `docs/` 路径，不再写旧 `doc/` 路径。
4. 模块名在表格中可以用业务名，但正文第一次出现时应给出真实目录。

## 8. 何时必须更新文档

出现以下变化时，必须同步更新至少一份主入口文档：

1. 启动链路变化
2. 中间件顺序变化
3. 路由入口或模块分组变化
4. 服务容器/ProviderRegistry 行为变化
5. 新增跨模块核心能力，例如搜索、事件、审查、AI 工作流

## 9. 文档检查清单

- [ ] 运行时主链是否还能映射到当前代码
- [ ] Mermaid 图是否仍与真实目录一致
- [ ] README 是否说明了当前推荐阅读顺序
- [ ] 是否明确标出了过渡态和兼容路径
- [ ] 是否把新发现的坑点沉淀到 `review/` 或 `issues/`
