# Qingyu AI Service - Phase3 v2.0

> Python 微服务：AI Agent 工作流、RAG 系统、LangGraph 编排

## 项目概述

本服务实现了 Qingyu 写作系统的 AI 能力提升 Phase3 v2.0，包括：

- **A2A 创作流水线**：大纲 → 角色 → 情节 → 审核（带反思循环）
- **RAG 系统**：结构化向量检索 + 元数据增强
- **上下文感知**：WorkspaceContextTool（借鉴 Cursor）
- **反思循环**：增强审核 Agent + 元调度器
- **gRPC 通信**：与 Go 后端高性能通信

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

## 项目结构

```
python_ai_service/
├── src/
│   ├── core/           # 核心模块（配置、日志、异常）
│   ├── api/            # FastAPI 路由
│   ├── agents/         # Agent 实现
│   │   ├── nodes/      # LangGraph 节点
│   │   ├── states/     # 状态 Schema
│   │   └── workflows/  # 工作流编排
│   ├── tools/          # LangChain Tools
│   ├── rag/            # RAG 系统
│   │   ├── milvus_client.py
│   │   ├── embedding_service.py
│   │   └── hybrid_retriever.py
│   ├── grpc_server/    # gRPC 服务端
│   └── main.py         # FastAPI 入口
├── proto/              # Protobuf 定义
├── tests/              # 测试
├── pyproject.toml      # Poetry 配置
└── README.md           # 本文件
```

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

- [架构设计](../doc/design/ai/phase3/README_v2.0升级指南.md)
- [API 文档](../doc/design/ai/phase3/14.Python_AI_Service_API设计.md)
- [开发指南](./docs/development.md)

## 许可证

Copyright © 2025 Qingyu Team

