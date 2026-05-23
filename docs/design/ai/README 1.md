# AI 设计历史副本

> Python 微服务：AI Agent 工作流、RAG 系统、LangGraph 编排
> 最后整理：2026-05-22

## Page Role

- legacy-snapshot
- current-owner: `design/ai/`
- current-bounded: 仅保留一份与主入口相邻的历史副本，不作为新读者入口

## Recommended Read Path

1. 先读 `README.md`。
2. 再读 `agent/07.Python_AI_Agent系统架构设计.md`。
3. 如需历史副本，再回看本文件。

## Boundary

- 本文件是主入口的历史副本。
- 不应作为新文档导航的默认落点。
- 若要删改入口层，优先改 `README.md`，不是这里。

## Quick Section Map

- 项目概述
- Phase 3 MVP 状态
- 技术栈
- 快速开始
- 项目结构
- 开发规范
- 测试
- 部署
- 监控
- 文档
- 许可证

## Quick Takeaways

- 这是副本，不是主入口。
- 正式导航看 `README.md`。

## Skip Guide

- 只看当前入口：跳过本文件。
- 只做历史比对：只看项目概述和文档段落即可。

## 📋 项目概述

本服务实现了 Qingyu 写作系统的 AI 能力提升 Phase3 v2.0，包括：

- ✅ **Creative Agent 工作流**：理解 → RAG检索 → 生成 → 审核 → 最终化（带重试循环）
- ✅ **LangChain Tools**：RAGTool、CharacterTool、OutlineTool
- ✅ **RAG 系统**：向量检索 + 元数据过滤
- ✅ **gRPC 通信**：与 Go 后端高性能通信
- ⏳ **上下文感知**：WorkspaceContextTool（待实现）
- ⏳ **A2A 创作流水线**：大纲 → 角色 → 情节（待实现）

## 🎯 Phase 3 MVP 状态

**当前版本**: MVP v1.0  
**完成度**: 核心功能 100%  
**详细报告**: 见 [phase3/README.md](./phase3/README.md)

## 技术栈

- **框架**: FastAPI 0.109+
- **Agent**: LangChain + LangGraph
- **向量数据库**: Milvus 2.3+
- **向量模型**: BAAI/bge-large-zh-v1.5
- **通信协议**: gRPC
- **日志**: structlog

## 快速开始

### 1. 安装依赖

```bash
# 使用 Poetry
poetry install

# 或使用 pip
pip install -r requirements.txt
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 配置 API Keys 和服务地址
```

### 3. 启动服务

```bash
# 开发模式（热重载）
poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000

# 生产模式
poetry run uvicorn src.main:app --host 0.0.0.0 --port 8000 --workers 4
```

### 4. 访问 API 文档

- FastAPI 文档: http://localhost:8000/docs
- ReDoc 文档: http://localhost:8000/redoc

## 📁 项目结构

```
python_ai_service/
├── src/
│   ├── core/                      # 核心模块
│   │   ├── config.py              # 配置管理
│   │   ├── logger.py              # 日志系统
│   │   ├── exceptions.py          # 异常定义
│   │   └── tools/                 # Tool基础框架 ✅
│   │       ├── base.py            # BaseTool基类
│   │       ├── registry.py        # ToolRegistry
│   │       └── langchain/         # LangChain Tools
│   │           ├── rag_tool.py        # RAG检索工具 ✅
│   │           ├── character_tool.py  # 角色工具 ✅
│   │           └── outline_tool.py    # 大纲工具 ✅
│   ├── agents/                    # Agent系统 ✅
│   │   ├── states/                # 状态定义
│   │   │   ├── base_state.py      # 基础状态 ✅
│   │   │   └── creative_state.py  # 创作状态 ✅
│   │   ├── nodes/                 # 工作流节点
│   │   │   ├── understanding.py   # 理解任务 ✅
│   │   │   ├── retrieval.py       # RAG检索 ✅
│   │   │   ├── generation.py      # 内容生成 ✅
│   │   │   ├── review.py          # 审核评估 ✅
│   │   │   └── finalize.py        # 最终化 ✅
│   │   └── workflows/             # 工作流编排
│   │       ├── creative.py        # 创作工作流 ✅
│   │       └── routers.py         # 路由函数 ✅
│   ├── services/                  # Service层 ✅
│   │   ├── agent_service.py       # Agent服务 ✅
│   │   ├── tool_service.py        # Tool服务 ✅
│   │   └── rag_service.py         # RAG服务 ✅
│   ├── infrastructure/            # 基础设施 ✅
│   │   └── go_api/                # Go API客户端
│   │       └── http_client.py     # HTTP客户端 ✅
│   ├── rag/                       # RAG系统
│   │   ├── milvus_client.py
│   │   ├── embedding_service.py
│   │   └── rag_pipeline.py
│   ├── grpc_server/               # gRPC服务端
│   │   └── servicer.py            # gRPC实现 ✅
│   ├── api/                       # FastAPI路由
│   │   └── health.py
│   └── main.py                    # FastAPI入口
├── proto/                         # Protobuf定义
├── tests/                         # 测试
├── pyproject.toml                 # Poetry配置
├── requirements.txt               # Pip依赖 ✅
├── PHASE3_MVP_IMPLEMENTATION.md   # MVP实施总结 ✅
└── README.md                      # 本文件
```

**✅ 已完成** | **⏳ 进行中** | **⏸️ 待实现**

## 开发规范

### 代码风格

- 使用 Black 格式化（行长 100）
- 使用 isort 排序导入
- 使用 mypy 类型检查
- 遵循 PEP 8

### 提交规范

- feat: 新功能
- fix: 修复 Bug
- refactor: 重构
- docs: 文档更新
- test: 测试相关

## 测试

```bash
# 运行所有测试
poetry run pytest

# 运行单个测试文件
poetry run pytest tests/test_api.py

# 生成覆盖率报告
poetry run pytest --cov=src --cov-report=html
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t qingyu-ai-service:latest .

# 运行容器
docker run -p 8000:8000 qingyu-ai-service:latest
```

### Docker Compose 部署

```bash
cd ../docker
docker-compose up -d
```

## 监控

- Prometheus 指标: http://localhost:8000/metrics
- Grafana 仪表盘: http://localhost:3000

## 文档

- [架构设计](./phase3/README.md)
- [API 文档](./phase3/14.Python_AI_Service_API设计.md)
- [当前 AI 设计入口](./README.md)

## 许可证

Copyright © 2025 Qingyu Team

