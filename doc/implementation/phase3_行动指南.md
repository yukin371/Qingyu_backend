# Phase3 v2.0 - 下一步行动指南

> **当前状态**: 阶段1（基础架构）100% 代码完成，Docker验证待完成  
> **最后更新**: 2025-10-28 18:45

---

## ✅ 已完成

### 阶段 1.1: Python 微服务项目搭建 ✓
- ✅ 完整的 FastAPI 应用骨架
- ✅ 配置管理（Pydantic Settings）
- ✅ 结构化日志（Structlog）
- ✅ 自定义异常体系
- ✅ 健康检查 API
- ✅ Docker 支持
- ✅ 测试框架

### 阶段 1.2: gRPC 通信协议 ✓
- ✅ Protobuf 协议定义（6个 RPC 方法）
- ✅ Go gRPC 客户端
- ✅ Python gRPC 服务端骨架
- ✅ 构建脚本（Makefile）
- ✅ gRPC 通信验证成功（2025-10-28）

### 阶段 1.3: Milvus 向量数据库部署 ✓
- ✅ MilvusClient 核心功能实现（~250行）
- ✅ EmbeddingService 完整实现（~135行）
- ✅ 集成测试用例（~180行）
- ✅ Docker Compose 配置完成
- ✅ 中文实施文档（完整）
- ⏳ Docker 服务验证待完成（镜像拉取问题）

---

## 📚 实施报告

### 阶段 1 实施报告
- [阶段 1.3：最终总结报告](./00进度指导/阶段1.3最终总结_2025-10-28.md) ✨ **最新**
- [阶段 1.3：Milvus 向量数据库部署实施报告](./00进度指导/阶段1.3_Milvus向量数据库部署实施报告_2025-10-28.md)
- [阶段 1.3：Docker部署问题说明](./00进度指导/Docker部署问题说明_2025-10-28.md)
- [gRPC 通信验证成功报告](./GRPC_SUCCESS_REPORT.md)
- [配置整合总结](./CONFIGURATION_INTEGRATION_SUMMARY.md)

---

## 🚀 立即执行（5个步骤）

### 步骤 1: 生成 Protobuf 代码

#### 检查 protoc 是否已安装

```bash
# 检查版本
protoc --version
```

如果未安装，请参考 [安装指南](#q1-protoc-命令找不到)。

#### 安装 Go 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### 生成代码

**Linux / macOS**:
```bash
make proto
```

**Windows (PowerShell)**:
```powershell
.\scripts\generate_proto_all.ps1
```

**手动生成（所有平台）**:
```bash
# Go 代码
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  -I python_ai_service/proto \
  python_ai_service/proto/ai_service.proto

# Python 代码
cd python_ai_service
python -m grpc_tools.protoc -I proto \
  --python_out=src/grpc_server \
  --grpc_python_out=src/grpc_server \
  proto/ai_service.proto
```

**预期输出**:
- `pkg/grpc/pb/ai_service.pb.go`
- `pkg/grpc/pb/ai_service_grpc.pb.go`
- `python_ai_service/src/grpc_server/ai_service_pb2.py`
- `python_ai_service/src/grpc_server/ai_service_pb2_grpc.py`

---

### 步骤 2: 安装 Python 依赖

```bash
cd python_ai_service

# 使用 Poetry（推荐）
poetry install

# 或使用 pip
pip install -r requirements.txt
```

---

### 步骤 3: 配置环境变量

```bash
# 复制示例配置
cp python_ai_service/.env.example python_ai_service/.env

# 编辑配置，至少需要设置：
# - OPENAI_API_KEY 或 ANTHROPIC_API_KEY
# - MILVUS_HOST（如果使用 Docker 则为 milvus）
# - REDIS_HOST（如果使用 Docker 则为 redis）
```

**关键配置项**:
```env
# AI 提供商（至少配置一个）
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_anthropic_key

# Milvus（阶段1.3后需要）
MILVUS_HOST=localhost
MILVUS_PORT=19530

# Go gRPC（用于 Python 调用 Go）
GO_GRPC_HOST=localhost
GO_GRPC_PORT=50051

# Embedding 模型
EMBEDDING_MODEL_NAME=BAAI/bge-large-zh-v1.5
EMBEDDING_MODEL_DEVICE=cpu  # 或 cuda
```

---

### 步骤 4: 测试 Python 服务

```bash
# 进入 Python 服务目录
cd python_ai_service

# 运行测试
poetry run pytest tests/ -v

# 启动服务（开发模式）
poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000

# 或使用快速启动脚本
./run.sh  # Linux/macOS
# 或
run.bat   # Windows
```

**验证**:
```bash
# 访问 API 文档
open http://localhost:8000/docs

# 健康检查
curl http://localhost:8000/api/v1/health

# 应该看到：
# {
#   "status": "healthy",
#   "service": "qingyu-ai-service",
#   "timestamp": "2025-10-28T...",
#   "version": "0.1.0"
# }
```

---

### 步骤 5: 开始阶段 1.3（部署 Milvus）

```bash
# 进入 docker 目录
cd docker

# 编辑 docker-compose.dev.yml，添加 Milvus 服务
# （参考下面的配置）

# 启动服务
docker-compose -f docker-compose.dev.yml up -d milvus etcd minio
```

**Docker Compose 配置示例**（添加到 `docker-compose.dev.yml`）:

```yaml
services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.5
    environment:
      - ETCD_AUTO_COMPACTION_MODE=revision
      - ETCD_AUTO_COMPACTION_RETENTION=1000
      - ETCD_QUOTA_BACKEND_BYTES=4294967296
      - ETCD_SNAPSHOT_COUNT=50000
    volumes:
      - etcd:/etcd
    command: etcd -advertise-client-urls=http://127.0.0.1:2379 -listen-client-urls http://0.0.0.0:2379 --data-dir /etcd
    healthcheck:
      test: ["CMD", "etcdctl", "endpoint", "health"]
      interval: 30s
      timeout: 20s
      retries: 3

  minio:
    image: minio/minio:RELEASE.2023-03-20T20-16-18Z
    environment:
      MINIO_ACCESS_KEY: minioadmin
      MINIO_SECRET_KEY: minioadmin
    volumes:
      - minio:/minio_data
    command: minio server /minio_data --console-address ":9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  milvus:
    image: milvusdb/milvus:v2.3.4
    command: ["milvus", "run", "standalone"]
    environment:
      ETCD_ENDPOINTS: etcd:2379
      MINIO_ADDRESS: minio:9000
    volumes:
      - milvus:/var/lib/milvus
    ports:
      - "19530:19530"
      - "9091:9091"
    depends_on:
      - etcd
      - minio
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9091/healthz"]
      interval: 30s
      start_period: 90s
      timeout: 20s
      retries: 3

volumes:
  etcd:
  minio:
  milvus:
```

---

## 📋 后续阶段（Week 2-10）

### Week 2: 向量化和 RAG（阶段2）
- 实现 Milvus 客户端完整功能
- 加载 BGE 向量模型
- 实现文本向量化
- 实现混合检索引擎

### Week 3: 上下文感知（阶段3）
- 实现 WorkspaceContextTool
- 任务类型识别
- 结构化上下文构建

### Week 4-6: 反思循环（阶段4）
- 增强审核 Agent
- 元调度器
- 修正策略
- A/B 测试

### Week 6-8: Agent 工作流（阶段5）
- LangGraph 工作流
- 4 个专业 Agent
- 7 个核心工具
- 端到端测试

---

## 🐛 常见问题

### Q1: protoc 命令找不到

**解决**:
```bash
# macOS
brew install protobuf

# Linux (Ubuntu/Debian)
sudo apt-get install -y protobuf-compiler

# Windows
# 从 https://github.com/protocolbuffers/protobuf/releases 下载
```

### Q2: Poetry 安装失败

**解决**:
```bash
# 使用 pip 安装 Poetry
pip install poetry

# 或使用官方安装脚本
curl -sSL https://install.python-poetry.org | python3 -
```

### Q3: 端口 8000 被占用

**解决**:
```bash
# 查看占用进程
lsof -i :8000  # Linux/macOS
netstat -ano | findstr :8000  # Windows

# 修改端口（编辑 .env）
SERVICE_PORT=8001
```

### Q4: gRPC 生成代码失败

**解决**:
```bash
# 检查 Python gRPC 工具
pip install grpcio-tools

# 手动生成 Python 代码
cd python_ai_service
python -m grpc_tools.protoc -I proto \
  --python_out=src/grpc_server \
  --grpc_python_out=src/grpc_server \
  proto/ai_service.proto
```

---

## 📚 参考文档

### 实施相关
- [实施计划](doc/implementation/00进度指导/计划/phase3-v2-0-implementation.plan.md)
- [实施进度](doc/implementation/00进度指导/计划/Phase3-v2.0/实施进度_2025-10-28.md)
- [实施总结](python_ai_service/IMPLEMENTATION_SUMMARY.md)

### 设计相关
- [v2.0 升级指南](doc/design/ai/phase3/README_v2.0升级指南.md)
- [A2A 流水线设计](doc/design/ai/phase3/05.A2A创作流水线Agent设计_v2.0_智能协作生态.md)

### Python 项目
- [Python 服务 README](python_ai_service/README.md)

---

## ✅ 检查清单

开始下一阶段前，确认：

- [ ] Protobuf 代码已生成（`make proto`）
- [ ] Python 依赖已安装（`poetry install`）
- [ ] 环境变量已配置（`.env` 文件）
- [ ] Python 服务可以启动
- [ ] 健康检查 API 正常
- [ ] 测试通过（`pytest tests/`）
- [ ] Docker Compose 配置已更新（Milvus）

---

## 💡 建议

1. **按阶段实施**: 不要跳过阶段，每个阶段都是后续的基础
2. **及时测试**: 每完成一个模块就测试，不要等到最后
3. **文档先行**: 先理解设计文档，再开始编码
4. **增量开发**: 小步快跑，频繁提交
5. **代码审查**: 关键模块需要 Review

---

**准备好了吗？** 开始执行步骤 1 → 生成 Protobuf 代码！

```bash
make proto
```

祝顺利！🚀

