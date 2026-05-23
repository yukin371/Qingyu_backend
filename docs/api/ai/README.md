# AI API 专题入口

> 最后整理: 2026-05-22

本目录保留写作辅助与 AI 接口专题文档，但其中部分内容带有设计/实施交叉属性，不应替代 `docs/ai/` 或 `docs/architecture/`。

## Page Role

- current-topic-hub
- current-owner: `docs/api/ai/`
- current-bounded: 当前 AI API 专题入口，只负责 AI 相关接口文档导航，不承担 AI 架构主入口职责

## Recommended Read Path

1. 先读 `AI_API_Documentation.md`。
2. 需要架构边界时，再回 `../../architecture/ai_grpc_integration.md`。

## Quick Section Map

- Current Documents
- Boundary

## Current Documents

- [AI_API_Documentation.md](./AI_API_Documentation.md)
- [AI_API_Design_Document.md](./AI_API_Design_Document.md)
- [AI_Module_Design_Document.md](./AI_Module_Design_Document.md)

## Boundary

- `../README.md`：API 总入口 owner。
- `../../ai/README.md`：后端 AI 接入说明 owner。
- `../../architecture/ai_grpc_integration.md`：AI 架构边界 owner。
