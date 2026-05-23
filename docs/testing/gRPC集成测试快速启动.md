# gRPC集成测试 - 快速启动指南

> **快速上手** | **10分钟完成测试** | **零配置烦恼**

---

## 🚀 一键启动（最简单）

### Windows

```bash
# 1. 设置 API Key（仅首次需要）
set GOOGLE_API_KEY=your_api_key_here

# 2. 运行集成测试
scripts\testing\test_grpc_integration.bat
```

### Linux/Mac

```bash
# 1. 设置 API Key（仅首次需要）
export GOOGLE_API_KEY=your_api_key_here

# 2. 运行集成测试（需要先创建对应的 .sh 脚本）
# 或手动执行以下步骤
```

**完成！** 脚本会自动：
- ✅ 检查依赖
- ✅ 启动 Python AI 服务
- ✅ 运行 Python 测试
- ✅ 运行 Go 测试
- ✅ 生成测试报告

---

## 📝 手动测试（分步骤）

### 第一步：准备环境

#### 1. 安装依赖

```bash
# Go 依赖
go mod download

# Python 依赖
cd Qingyu-Ai-Service
pip install -r requirements.txt
cd ..
```

#### 2. 设置 API Key

```bash
# Windows
set GOOGLE_API_KEY=your_google_api_key

# Linux/Mac
export GOOGLE_API_KEY=your_google_api_key
```

#### 3. 验证环境

```bash
# 检查 Go
go version

# 检查 Python
python --version

# 检查 API Key
echo %GOOGLE_API_KEY%  # Windows
echo $GOOGLE_API_KEY    # Linux/Mac
```

### 第二步：启动 Python AI 服务

打开一个**新的终端窗口**：

```bash
cd Qingyu-Ai-Service
python run_grpc_server.py
```

**看到以下输出表示成功**：

```
========================================
Phase3 gRPC Server Startup
========================================

API Key: AIzaSyB1234567890...

[1/2] Checking dependencies...
Dependencies: OK

[2/2] Starting gRPC server...

========================================
Server will listen on 0.0.0.0:50051
Press Ctrl+C to stop
========================================

🚀 gRPC服务器启动 - 监听地址: 0.0.0.0:50051
📋 可用服务:
  - ExecuteCreativeWorkflow: 完整创作工作流
  - GenerateOutline: 大纲生成
  - GenerateCharacters: 角色生成
  - GeneratePlot: 情节生成
  - HealthCheck: 健康检查
✅ gRPC服务器就绪，等待请求...
```

**保持此窗口运行！**

### 第三步：测试 Python 客户端

打开**另一个新终端窗口**：

```bash
cd Qingyu-Ai-Service
python tests\test_grpc_phase3.py
```

**预期输出**：

```
========================================
🚀 Phase3 gRPC服务测试
========================================

========================================
🏥 测试健康检查
========================================
✅ 健康状态: healthy
📋 检查结果:
  - agents: ok
  - llm: ok

========================================
📖 测试大纲生成
========================================
📝 任务: 创作修仙小说大纲...

✅ 大纲生成成功!
📖 标题: 逆天修仙路
🎭 类型: 修仙
📚 章节数: 5
⏱️  耗时: 9.83秒

... (更多输出)

========================================
✅ 所有测试通过!
========================================
```

### 第四步：测试 Go 客户端

打开**第三个终端窗口**（或关闭 Python 测试窗口后使用）：

#### 方式 A：运行单个命令行工具

```bash
# 在项目根目录
go run cmd\test_phase3_grpc\main.go --addr localhost:50051
```

#### 方式 B：运行完整测试套件

```bash
# 在项目根目录
go test -v -timeout 300s ./test/integration -run TestGRPC
```

**预期输出**：

```
========================================
Phase3 gRPC客户端测试
========================================

连接到gRPC服务器: localhost:50051
✅ 连接成功

1️⃣  健康检查...
   状态: healthy
   组件状态:
     - agents: ok
     - llm: ok

2️⃣  生成大纲...
   任务: 创作一个修仙小说大纲，主角是天才少年
   ✅ 成功! 耗时: 9.78秒
   📖 标题: 天才少年修仙录
   🎭 类型: 修仙
   📚 章节数: 5

... (更多输出)

========================================
✅ 测试完成
========================================
```

---

## 🔍 验证测试结果

### 成功标志

所有测试应该看到：

- ✅ Python 服务启动成功
- ✅ Python 客户端测试全部通过
- ✅ Go 客户端测试全部通过
- ✅ 无报错和异常

### 测试覆盖

| 测试项 | Python 客户端 | Go 客户端 |
|-------|-------------|----------|
| 健康检查 | ✅ | ✅ |
| 大纲生成 | ✅ | ✅ |
| 角色生成 | ✅ | ✅ |
| 情节生成 | ✅ | ✅ |
| 完整工作流 | ✅ | ✅ |

---

## 🎯 只想快速验证连接？

最简单的测试 - 仅验证 gRPC 连接：

```bash
# 1. 启动 Python 服务（终端1）
cd Qingyu-Ai-Service
python run_grpc_server.py

# 2. 测试连接（终端2）
cd Qingyu-Ai-Service
python -c "
import grpc
from grpc_service import ai_service_pb2, ai_service_pb2_grpc
import asyncio

async def test():
    async with grpc.aio.insecure_channel('localhost:50051') as channel:
        stub = ai_service_pb2_grpc.AIServiceStub(channel)
        resp = await stub.HealthCheck(ai_service_pb2.HealthCheckRequest())
        print(f'Status: {resp.status}')

asyncio.run(test())
"
```

---

## ❌ 常见问题

### 问题 1: "GOOGLE_API_KEY 未设置"

```
ERROR: GOOGLE_API_KEY environment variable is not set
```

**解决**：
```bash
set GOOGLE_API_KEY=your_actual_api_key
```

### 问题 2: "连接被拒绝"

```
grpc._channel._InactiveRpcError: failed to connect to all addresses
```

**原因**: Python AI 服务未启动

**解决**: 先启动 Python 服务
```bash
cd Qingyu-Ai-Service
python run_grpc_server.py
```

### 问题 3: "端口已被占用"

```
OSError: [WinError 10048] 通常每个套接字地址(协议/网络地址/端口)只允许使用一次
```

**解决**:
```bash
# 查找占用端口的进程
netstat -ano | findstr :50051

# 结束进程（记下 PID）
taskkill /PID <PID> /F

# 重新启动服务
```

### 问题 4: "Python 依赖缺失"

```
ModuleNotFoundError: No module named 'grpc'
```

**解决**:
```bash
cd Qingyu-Ai-Service
pip install -r requirements.txt
```

### 问题 5: "Go 编译错误"

```
cannot find package "google.golang.org/grpc"
```

**解决**:
```bash
go mod download
go mod tidy
```

---

## 📊 性能预期

| 操作 | 预期时间 | 说明 |
|-----|---------|------|
| 服务启动 | 5-10秒 | Python AI 服务初始化 |
| 健康检查 | < 0.1秒 | 连接验证 |
| 大纲生成 | 8-15秒 | 调用 LLM API |
| 角色生成 | 10-18秒 | 调用 LLM API |
| 情节生成 | 12-20秒 | 调用 LLM API |
| 完整工作流 | 30-60秒 | 完整流程 |

---

## 📖 下一步

测试通过后，您可以：

1. **集成到 Go 服务** - 在 Service 层集成 Phase3Client
2. **添加业务逻辑** - 实现具体的业务场景
3. **性能优化** - 添加缓存和并发控制
4. **生产部署** - 使用 Docker 和 K8s 部署

详细文档：
- [gRPC 集成测试指南](./gRPC集成测试指南.md)
- [Phase3 gRPC README](../../../Qingyu-Ai-Service/PHASE3_GRPC_README.md)

---

## 🤝 需要帮助？

如果遇到问题：

1. 查看 [故障排查指南](./gRPC集成测试指南.md#故障排查)
2. 检查 Python 服务日志
3. 查看 [完整测试指南](./gRPC集成测试指南.md)

---

**最后更新**: 2025-10-31  
**维护者**: 青羽后端架构团队

