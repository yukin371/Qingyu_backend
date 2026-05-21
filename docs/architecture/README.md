# Qingyu Backend Architecture Hub

本目录是后端架构主入口，目标是让人类和 AI 在 10 分钟内建立当前可用认知，而不是在历史文档中盲目搜索。

## First Read Path

1. [system_architecture.md](./system_architecture.md)
2. [2026-04-07-backend-runtime-flow.md](./2026-04-07-backend-runtime-flow.md)
3. [2026-04-07-backend-module-map.md](./2026-04-07-backend-module-map.md)
4. [dependency-rules.md](./dependency-rules.md)

## Current Source Of Truth

- [system_architecture.md](./system_architecture.md): 当前分层边界、启动流程、请求链路主图
- [2026-04-07-backend-runtime-flow.md](./2026-04-07-backend-runtime-flow.md): `main -> core -> container -> router` 真实运行时链路
- [2026-04-07-backend-module-map.md](./2026-04-07-backend-module-map.md): 业务域/横切能力/平台基础设施三层模块地图
- [dependency-rules.md](./dependency-rules.md): 依赖约束与边界规则
- [data_model.md](./data_model.md): 模型与数据层结构

## Specialized But Still Active

- [ai_grpc_integration.md](./ai_grpc_integration.md): AI 服务专项集成路径
- [api_architecture.md](./api_architecture.md): API 分层和接口视角

## Historical Reference Boundary

以下文档保留用于历史追溯，不作为当前架构主依据：

- [项目开发规则.md](./项目开发规则.md): 分层描述已过时（缺 Repository，示例与现状不一致）
- [项目概述.md](./项目概述.md): 包含大量愿景和规划性叙述，不是当前运行时事实
- [architecture_review_report.md](./architecture_review_report.md): 历史审查结论，需二次核验后再采纳
- [component_analysis.md](./component_analysis.md): 组件分析可参考，但模块关系需按当前代码刷新
- [module-dependency-analysis.md](./module-dependency-analysis.md): 主题较窄，不能替代全局模块地图
- [2026-03-04-writer-editor-v2-*.md](./2026-03-04-writer-editor-v2-merge-status.md): 阶段性交接文档
- [2026-04-07-backend-runtime-architecture.md](./2026-04-07-backend-runtime-architecture.md): 迁移过程中的中间稿，已由 `backend-runtime-flow` 取代

## Maintenance Rules

1. 新增“主入口级”架构结论，优先更新本目录这三份文件：`README`、`system_architecture`、`backend-runtime-flow`。
2. Mermaid 图必须映射真实代码路径，禁止只画理想态。
3. 发现文档与代码不一致，先在文档中显式标注偏差，再进入重构或治理任务。
4. 阶段性交付文档可保留，但必须明确标注为历史资料，避免与当前主入口混淆。
