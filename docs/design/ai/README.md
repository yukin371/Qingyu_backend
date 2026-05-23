# AI 模块设计入口

> 目录定位：当前主线、完成态专题、历史参考的统一入口
> 使用原则：只做导航与分层说明，不把这里写成实施流水账

## Page Role

- current-hub
- current-owner: `design/ai/`
- current-bounded: 当前 AI 设计主入口与阅读导航，不承载实现细节的逐步堆叠

## Recommended Read Path

1. 先读 `agent/07.Python_AI_Agent系统架构设计.md`。
2. 再读 `agent/08.Go_AI代理层设计.md` 和 `agent/09.Agent工具调用集成设计.md`。
3. 需要补检索与流式边界时，再读 `rag/10.RAG检索增强系统设计.md`、`rag/11.RAG事件驱动索引设计.md`、`streaming/12.AI流式接口规范.md`。
4. 需要阶段复盘时，再去 `phase3/README.md`。

## Boundary

- 这里是 AI 设计入口层，不是实现清单。
- 目录分层以 `agent/` 为当前主线，`phase3/` 为完成态专题，`archived/` 为历史参考。
- `README 1.md` 仅作历史副本，不作为主入口。

## Quick Section Map

- 目录分层
- 阅读顺序
- 目录说明
- 使用边界
- 快速索引

## Quick Takeaways

- 当前主线看 `agent/`。
- 完成态专题看 `phase3/`。
- 旧方案只看 `archived/`。

## Skip Guide

- 只想找当前设计入口：直接看 `agent/07.Python_AI_Agent系统架构设计.md`。
- 只想看历史演进：去 `archived/README.md`。
- 只想核对完成态专题：去 `phase3/README.md`。

## 目录分层

- 当前主线：`agent/`
- 完成态专题：`phase3/`
- 历史参考：`archived/`
- 补充专题：`rag/`、`streaming/` 以及目录下的迁移说明文档

## 阅读顺序

1. 先读 `agent/07.Python_AI_Agent系统架构设计.md`，把当前主线的 Python AI Agent 架构、工具调用、RAG 和多 Agent 协作边界先看清楚。
2. 再读 `agent/08.Go_AI代理层设计.md` 与 `agent/09.Agent工具调用集成设计.md`，补齐 Go 代理层和工具接入层的协作方式。
3. 需要理解检索和接口约束时，再看 `rag/10.RAG检索增强系统设计.md`、`rag/11.RAG事件驱动索引设计.md`、`streaming/12.AI流式接口规范.md`。
4. 如果你是在复盘一个已经完成的阶段，再转到 `phase3/README.md`。
5. 如果你在查旧方案的演进路径，只把 `archived/README.md` 当历史参考入口。

## 目录说明

### 当前主线

`agent/` 负责当前还在沿用和持续对齐的主线设计，重点是 Python AI Agent 系统、Go 代理层和工具调用边界。这里的文档应优先作为新需求、新实现和新澄清的引用源。

### 完成态专题

`phase3/` 是已经完成的专题索引。这里的文档适合用来确认阶段内设计结论、回看已交付能力和补充实现背景，但不应被当作当前主线继续扩写。

### 历史参考

`archived/` 只保留旧架构和被替代方案，适合看演进史、查旧决策和对照差异，不作为新设计的落点。

## 使用边界

- 本目录是设计文档入口层，不是实现代码入口。
- 新的主线判断默认以 `agent/07.Python_AI_Agent系统架构设计.md` 为准。
- 需要补新设计时，优先补主线文档，再决定是否沉淀为完成态专题或历史参考。
- 目录中另有 `README 1.md`，当前不作为主入口导航。

## 快速索引

| 位置 | 定位 | 说明 |
|---|---|---|
| `agent/` | 当前主线 | 当前仍在承接的 AI 架构、代理层与工具层设计 |
| `phase3/` | 完成态专题 | 16 篇设计文档的完成态索引 |
| `archived/` | 历史参考 | 已被新架构替代的旧文档 |
| `rag/` | 专题补充 | RAG 设计和索引机制 |
| `streaming/` | 专题补充 | 强制流式接口规范 |
