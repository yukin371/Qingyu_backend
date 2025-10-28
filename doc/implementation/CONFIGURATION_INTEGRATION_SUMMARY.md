# Go后端 + Python AI服务配置整合总结

## 实施日期
2025-10-28

## 整合策略
**环境变量优先 + YAML作为默认源**

- Go后端：Viper加载YAML配置（`config/config.yaml`），环境变量（`QINGYU_`前缀）覆盖
- Python服务：Pydantic Settings从环境变量加载，支持`.env`文件
- 共享变量：如`OPENAI_API_KEY`、`REDIS_HOST`等通过环境变量同步

## 已完成工作

### 1. 扩展Go YAML配置 (`config/config.yaml`)
添加了`ai.python_service`配置段：
```yaml
ai:
  api_key: "${OPENAI_API_KEY}"  # 支持环境变量引用
  base_url: "https://api.openai.com/v1"
  model: "gpt-4-turbo-preview"
  python_service:
    host: "localhost"  # Docker: python-ai-service
    grpc_port: 50052
    embedding_model: "BAAI/bge-large-zh-v1.5"
    milvus_host: "localhost"
    milvus_port: 19530
    redis_host: "localhost"
    redis_port: 6379
```

### 2. 更新Go配置结构体 (`config/config.go`)
- 新增`PythonAIServiceConfig`结构体
- 在`AIConfig`中添加`PythonService`字段
- 在`setDefaults()`设置默认值，确保Viper Unmarshal兼容

### 3. 创建共享环境变量模板
创建`.env.example`（根目录），包含：
- AI密钥：`OPENAI_API_KEY`、`ANTHROPIC_API_KEY`
- gRPC连接：`GO_GRPC_HOST`、`GO_GRPC_PORT`、`PYTHON_AI_GRPC_PORT`
- 数据库：`MILVUS_HOST`、`REDIS_HOST`等

### 4. 扩展Docker Compose (`docker/docker-compose.dev.yml`)
- 新增`python-ai-service`服务（ports 8000/50052）
- 新增Milvus栈：milvus、etcd、minio
- Go服务重命名为`go-backend`
- 环境变量注入：`${OPENAI_API_KEY}`从宿主机`.env`加载

### 5. Go AI服务注入Python配置 (`service/ai/ai_service.go`)
- 在`Service`结构体添加`PythonConfig`字段
- 在`GenerateContent`中优先使用gRPC调用Python（fallback到adapter）
- 示例逻辑：
```go
if s.PythonConfig != nil && s.PythonConfig.GrpcPort > 0 {
    conn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.PythonConfig.Host, s.PythonConfig.GrpcPort), ...)
    client := pb.NewAIServiceClient(conn)
    // ... 调用GenerateContent
}
```

## 服务状态

### Python AI服务 (FastAPI)
- ✅ 端口8000已启动
- ✅ 健康检查通过：`/api/v1/health` 返回 `{"status":"healthy"}`
- ✅ 配置加载成功（pydantic-settings）
- ⚠️ gRPC端口50052未监听（需排查）

### Go后端
- ⏳ 配置已更新，支持Python配置读取
- ⏳ gRPC客户端代码已添加
- ⏳ 需验证启动和连接

## gRPC连接问题排查

### 当前问题
Python gRPC服务器（端口50052）未成功监听，Go客户端连接失败：
```
rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:50052: connectex: No connection could be made because the target machine actively refused it."
```

### 排查步骤
1. **Python FastAPI启动检查**：✅ 运行在8000
2. **gRPC服务器启动**：❌ `start_grpc_server()`在FastAPI lifespan中调用，但端口未监听
3. **可能原因**：
   - `asyncio.run(serve())`在已有event loop中调用导致冲突
   - 线程启动逻辑需要同步等待

### 临时解决方案
创建独立测试服务器`python_ai_service/test_grpc_server.py`，直接运行gRPC：
```bash
cd python_ai_service
python test_grpc_server.py  # 监听50052
```

## 使用指南

### 本地开发
1. **复制环境变量**：
   ```bash
   cp .env.example .env
   # 编辑.env，填写OPENAI_API_KEY
   ```

2. **启动Python AI服务**：
   ```bash
   cd python_ai_service
   python test_grpc_server.py  # 独立gRPC
   # 或使用FastAPI（需修复gRPC启动）
   python -m uvicorn src.main:app --host 0.0.0.0 --port 8000
   ```

3. **启动Go后端**：
   ```bash
   go run cmd/server/main.go
   ```

4. **测试gRPC连接**：
   ```bash
   go run test_grpc_connection.go
   ```

### Docker环境
```bash
docker-compose -f docker/docker-compose.dev.yml up -d
docker logs -f qingyu-python-ai  # 查看Python日志
docker logs -f qingyu-backend    # 查看Go日志
```

## 下一步建议

1. **修复gRPC启动**：
   - 使用`asyncio.create_task()`替代`threading.Thread`
   - 或在FastAPI之前独立启动gRPC（separate process）

2. **测试覆盖**：
   - 添加`test/integration/grpc_connection_test.go`
   - Mock Python gRPC服务测试Go客户端

3. **生产配置**：
   - 使用Kubernetes Secrets管理密钥
   - 添加gRPC TLS加密
   - 配置健康检查和监控

4. **文档更新**：
   - 更新`doc/ops/部署指南.md`
   - 添加配置整合示例到`doc/architecture/`

## 技术栈验证
- ✅ Go 1.21+ (Viper + gRPC)
- ✅ Python 3.13.7 (FastAPI + grpcio 1.76.0)
- ✅ Protocol Buffers (ai_service.proto已生成Go/Python代码)
- ✅ Pydantic Settings 2.11.0 (环境变量管理)
- ✅ Milvus/Redis配置（Docker Compose已添加）

## 附件
- 配置文件：`config/config.yaml`
- Go结构体：`config/config.go`
- Docker Compose：`docker/docker-compose.dev.yml`
- gRPC测试：`test_grpc_connection.go`
- Proto定义：`python_ai_service/proto/ai_service.proto`

---
**维护者**：青羽后端架构团队  
**最后更新**：2025-10-28

