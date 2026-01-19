# Phase3 v2.0 实施指南

> **实施状态**: 🚧 进行中  
> **当前阶段**: 阶段1 - 基础架构搭建  
> **完成度**: 60%

---

## 📋 文档索引

### 计划文档
- [实施计划](../../../phase3-v2-0-implementation.plan.md) - 完整实施计划
- [实施进度](./实施进度_2025-10-28.md) - 当前进度报告

### 设计文档
- [v2.0 升级指南](../../../../design/ai/phase3/README_v2.0升级指南.md) - 架构设计
- [A2A 流水线设计](../../../../design/ai/phase3/05.A2A创作流水线Agent设计_v2.0_智能协作生态.md) - 详细设计

---

## 🚀 快速开始

### 前置条件

1. **Python 环境**
   - Python 3.10+
   - Poetry（推荐）或 pip

2. **Go 环境**
   - Go 1.21+
   - Protobuf 编译器

3. **Docker**
   - Docker 20.10+
   - Docker Compose 2.0+

### 步骤 1: 克隆项目

```bash
cd Qingyu_backend
```

### 步骤 2: 生成 Protobuf 代码

```bash
# 确保安装了 protoc
# macOS: brew install protobuf
# Linux: sudo apt-get install protobuf-compiler

# 生成所有 Protobuf 代码
make proto

# 或分别生成
make proto-go      # 生成 Go 代码
make proto-python  # 生成 Python 代码
```

### 步骤 3: 安装 Python 依赖

```bash
cd python_ai_service

# 使用 Poetry（推荐）
poetry install

# 或使用 pip
pip install -r requirements.txt
```

### 步骤 4: 配置环境变量

```bash
# 复制示例配置
cp .env.example .env

# 编辑配置（填写 API Keys）
vim .env
```

### 步骤 5: 启动服务

```bash
# 使用 Poetry
poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000

# 或使用快速启动脚本
./run.sh  # Linux/macOS
run.bat   # Windows
```

### 步骤 6: 验证服务

```bash
# 访问 FastAPI 文档
open http://localhost:8000/docs

# 健康检查
curl http://localhost:8000/api/v1/health
```

---

## 📦 Docker 部署

### 构建镜像

```bash
cd python_ai_service
docker build -t qingyu-ai-service:latest .
```

### 运行容器

```bash
docker run -d \
  --name qingyu-ai-service \
  -p 8000:8000 \
  -e OPENAI_API_KEY=your_key \
  qingyu-ai-service:latest
```

### Docker Compose（完整栈）

```bash
cd docker
docker-compose -f docker-compose.dev.yml up -d
```

---

## 🧪 测试

### 运行测试

```bash
cd python_ai_service

# 运行所有测试
poetry run pytest tests/ -v

# 运行特定测试
poetry run pytest tests/test_health.py -v

# 生成覆盖率报告
poetry run pytest --cov=src --cov-report=html
```

### 手动测试 API

```bash
# 健康检查
curl http://localhost:8000/api/v1/health

# 就绪检查
curl http://localhost:8000/api/v1/health/ready

# 存活检查
curl http://localhost:8000/api/v1/health/live
```

---

## 🔧 开发工作流

### 1. 创建新功能分支

```bash
git checkout -b feature/stage-1-3-milvus
```

### 2. 实现功能

参考实施计划，按照模块划分实现：
- RAG 系统 → `src/rag/`
- Agent 节点 → `src/agents/nodes/`
- 工具层 → `src/tools/`

### 3. 编写测试

```python
# tests/test_new_feature.py
import pytest

def test_new_feature(client):
    response = client.get("/api/v1/new-feature")
    assert response.status_code == 200
```

### 4. 运行测试和 Linter

```bash
# 运行测试
poetry run pytest

# 代码格式化
poetry run black src/ tests/
poetry run isort src/ tests/

# 类型检查
poetry run mypy src/
```

### 5. 提交代码

```bash
git add .
git commit -m "feat: implement stage 1.3 - milvus integration"
git push origin feature/stage-1-3-milvus
```

---

## 📂 项目结构

```
python_ai_service/
├── src/
│   ├── core/              # 核心模块
│   │   ├── config.py      # 配置管理
│   │   ├── logger.py      # 日志系统
│   │   └── exceptions.py  # 异常定义
│   ├── api/               # FastAPI 路由
│   │   └── health.py      # 健康检查
│   ├── grpc_server/       # gRPC 服务端
│   │   ├── servicer.py    # Servicer 实现
│   │   └── server.py      # Server 启动
│   ├── agents/            # Agent 实现
│   │   ├── nodes/         # LangGraph 节点
│   │   ├── states/        # 状态 Schema
│   │   └── workflows/     # 工作流编排
│   ├── tools/             # LangChain Tools
│   ├── rag/               # RAG 系统
│   └── main.py            # 应用入口
├── proto/                 # Protobuf 定义
├── tests/                 # 测试
├── pyproject.toml         # Poetry 配置
└── README.md              # 项目说明
```

---

## 🐛 故障排查

### 问题 1: Protobuf 代码生成失败

**症状**: `make proto` 失败

**解决方案**:
```bash
# 检查 protoc 是否安装
protoc --version

# 安装 Python gRPC 工具
pip install grpcio-tools

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 问题 2: 依赖安装失败

**症状**: `poetry install` 失败

**解决方案**:
```bash
# 清理缓存
poetry cache clear pypi --all

# 重新安装
poetry install --no-cache

# 或使用 pip
pip install -r requirements.txt
```

### 问题 3: 服务启动失败

**症状**: Uvicorn 启动错误

**解决方案**:
```bash
# 检查端口是否被占用
lsof -i :8000

# 检查环境变量
cat .env

# 查看详细日志
poetry run uvicorn src.main:app --log-level debug
```

---

## 📊 实施里程碑

| 阶段 | 任务 | 状态 | 预计完成 |
|------|------|------|---------|
| 1.1 | Python 微服务项目搭建 | ✅ 完成 | 2025-10-28 |
| 1.2 | gRPC 通信协议实现 | 🚧 进行中 | 2025-10-28 |
| 1.3 | Milvus 向量数据库部署 | ⏳ 待开始 | 2025-10-29 |
| 2.1 | 向量化引擎实现 | ⏳ 待开始 | 2025-11-01 |
| 2.2 | 结构化 RAG 实现 | ⏳ 待开始 | 2025-11-05 |
| 2.3 | 事件驱动索引更新 | ⏳ 待开始 | 2025-11-08 |
| 3.1 | WorkspaceContextTool | ⏳ 待开始 | 2025-11-12 |
| 3.2 | 集成到 Agent Prompt | ⏳ 待开始 | 2025-11-15 |
| 4.1 | 增强审核 Agent | ⏳ 待开始 | 2025-11-20 |
| 4.2 | 元调度器 | ⏳ 待开始 | 2025-11-25 |
| 4.3 | 反思循环集成测试 | ⏳ 待开始 | 2025-11-28 |
| 5.1 | LangGraph 工作流搭建 | ⏳ 待开始 | 2025-12-05 |
| 5.2 | 专业 Agent 实现 | ⏳ 待开始 | 2025-12-12 |
| 5.3 | 工具层实现 | ⏳ 待开始 | 2025-12-15 |
| 5.4 | A2A 流水线集成测试 | ⏳ 待开始 | 2025-12-20 |

---

## 🔗 相关资源

### 官方文档
- [FastAPI](https://fastapi.tiangolo.com/)
- [LangChain](https://python.langchain.com/)
- [LangGraph](https://langchain-ai.github.io/langgraph/)
- [Milvus](https://milvus.io/docs)
- [gRPC Python](https://grpc.io/docs/languages/python/)

### 项目文档
- [后端 API 文档](../../../api/)
- [架构设计](../../../architecture/)
- [测试文档](../../../testing/)

---

## 💬 联系和支持

- **技术问题**: 查看实施计划和设计文档
- **Bug 报告**: 创建 Issue
- **功能建议**: 参考 v2.0 升级指南

---

**最后更新**: 2025-10-28  
**维护者**: Qingyu AI Team

