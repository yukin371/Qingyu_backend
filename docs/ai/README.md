# 后端 AI 集成文档入口

> 最后整理: 2026-05-22

本目录只承载 `Qingyu_backend` 侧的 AI 接入说明，不再兼任 Python AI 服务总览或第二套架构入口。

## First Read Path

1. [../architecture/ai_grpc_integration.md](../architecture/ai_grpc_integration.md)
2. [GRPC_INTEGRATION_GUIDE.md](./GRPC_INTEGRATION_GUIDE.md)
3. [PHASE3_QUICKSTART.md](./PHASE3_QUICKSTART.md)
4. [../implementation/README.md](../implementation/README.md)

## Owner Boundary

- `../architecture/ai_grpc_integration.md`：当前后端与 AI 服务的架构边界、依赖关系和运行链路 owner。
- `GRPC_INTEGRATION_GUIDE.md`：后端接入 AI gRPC 服务的操作型主指南。
- `Qingyu-Ai-Service/README.md`：Python AI 服务启动、目录结构、部署和运行细节 owner。
- 父仓 `docs/plans/`：跨仓库 AI 规划、阶段方案和长期设计 owner。

## Current Documents

| 文档 | 状态 | 用途 |
|------|------|------|
| [GRPC_INTEGRATION_GUIDE.md](./GRPC_INTEGRATION_GUIDE.md) | `current-owner` | Go 后端接入 AI gRPC 服务的主指南 |
| [PHASE3_QUICKSTART.md](./PHASE3_QUICKSTART.md) | `current-bounded` | Phase 3 相关快速启动与联调入口 |
| [PHASE3_GRPC_README.md](./PHASE3_GRPC_README.md) | `current-bounded` | gRPC 目录和联调背景说明 |
| [PHASE3_GRPC_INTEGRATION_COMPLETE.md](./PHASE3_GRPC_INTEGRATION_COMPLETE.md) | `current-bounded` | 一轮集成完成总结 |

## Historical Supplements

以下文档保留用于历史追溯，不应继续被当作当前总入口：

- [LANGCHAIN_1.0_IMPLEMENTATION_SUMMARY.md](./LANGCHAIN_1.0_IMPLEMENTATION_SUMMARY.md)
- [LANGCHAIN_1.0_REFACTOR_PROGRESS.md](./LANGCHAIN_1.0_REFACTOR_PROGRESS.md)
- [LANGCHAIN_1.0_REFACTOR_COMPLETE.md](./LANGCHAIN_1.0_REFACTOR_COMPLETE.md)
- [LangChain1.0迁移指南_2025-1105.md](./LangChain1.0迁移指南_2025-1105.md)
- [Service层快速参考_2025-1028.md](./Service层快速参考_2025-1028.md)
- [Service层完善实施总结_2025-1028.md](./Service层完善实施总结_2025-1028.md)

## Practical Rules

1. 如果文档在讲“当前架构边界”，优先更新 `../architecture/ai_grpc_integration.md`，不要在本目录复制第二份架构总览。
2. 如果文档在讲“如何联调、如何接入、如何排查 gRPC”，优先落在本目录。
3. 如果文档在讲“AI 服务本身如何运行”，回到 `Qingyu-Ai-Service` 仓库，不在这里复制 Python 运行手册。
4. 如果需要追溯更老的 Phase 3 设计稿，当前记为 `TBD`。
   确认路径：父仓 `docs/plans/submodules/backend/`、`docs/plans/` 的历史计划，以及 `Qingyu-Ai-Service` 仓库中的设计/迁移文档。
