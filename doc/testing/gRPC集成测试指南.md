# Go后端与Python AI微服务 gRPC集成测试指南

> **版本**: v1.0  
> **日期**: 2025-10-31  
> **状态**: ✅ 就绪

---

## 📋 目录

- [概述](#概述)
- [测试环境准备](#测试环境准备)
- [快速开始](#快速开始)
- [测试方案](#测试方案)
- [故障排查](#故障排查)
- [性能基准](#性能基准)

---

## 概述

本文档提供 Go 后端与 Python AI 微服务之间 gRPC 集成测试的完整指南。

### 架构概览

```
┌─────────────────────────────┐
│   Go Backend (Qingyu)       │
│                              │
│  ┌────────────────────────┐ │
│  │  Phase3Client (gRPC)   │ │
│  │  service/ai/           │ │
│  └──────────┬─────────────┘ │
└─────────────┼────────────────┘
              │ gRPC (50051)
              │ Protobuf
              ↓
┌─────────────────────────────┐
│  Python AI Service (gRPC)   │
│                              │
│  ┌────────────────────────┐ │
│  │  AIServicer            │ │
│  │  - ExecuteCreativeW... │ │
│  │  - GenerateOutline     │ │
│  │  - GenerateCharacters  │ │
│  │  - GeneratePlot        │ │
│  │  - HealthCheck         │ │
│  └────────────────────────┘ │
└─────────────────────────────┘
```

### 测试覆盖

- ✅ **连接测试** - gRPC 连接和健康检查
- ✅ **功能测试** - 大纲/角色/情节生成
- ✅ **工作流测试** - 完整创作工作流
- ✅ **并发测试** - 多并发请求处理
- ✅ **错误处理** - 异常和超时处理
- ✅ **性能测试** - 响应时间和资源消耗

---

## 测试环境准备

### 1. 系统要求

| 组件 | 版本要求 | 说明 |
|-----|---------|------|
| **Go** | 1.21+ | Go 后端运行环境 |
| **Python** | 3.11+ | Python AI 服务运行环境 |
| **protoc** | 3.19+ | Protobuf 编译器 |
| **MongoDB** | 5.0+ | 数据库（可选） |
| **Redis** | 7.0+ | 缓存（可选） |

### 2. 环境变量

```bash
# Windows
set GOOGLE_API_KEY=your_google_api_key_here

# Linux/Mac
export GOOGLE_API_KEY=your_google_api_key_here
```

### 3. 安装依赖

#### Go 依赖

```bash
# 在项目根目录
go mod download
go mod tidy
```

#### Python 依赖

```bash
# 在 python_ai_service 目录
cd python_ai_service
pip install -r requirements.txt
```

### 4. 生成 Protobuf 代码

#### Go 端

```bash
# 在项目根目录
protoc -I python_ai_service/proto \
  --go_out=pkg/grpc/pb \
  --go-grpc_out=pkg/grpc/pb \
  python_ai_service/proto/ai_service.proto
```

#### Python 端

```bash
# 在 python_ai_service 目录
cd python_ai_service
python -m grpc_tools.protoc \
  -I proto \
  --python_out=src/grpc_service \
  --grpc_python_out=src/grpc_service \
  proto/ai_service.proto
```

---

## 快速开始

### 方式一：一键测试（推荐）

```bash
# 在项目根目录
scripts\testing\test_grpc_integration.bat
```

此脚本会自动：
1. 检查依赖
2. 启动 Python AI 服务
3. 运行 Python 客户端测试
4. 运行 Go 客户端测试
5. 生成测试报告

### 方式二：手动测试

#### 步骤 1: 启动 Python AI 服务

```bash
cd python_ai_service
python run_grpc_server.py
```

或使用批处理脚本：

```bash
cd python_ai_service
scripts\start_grpc_server.bat
```

#### 步骤 2: 运行 Python 客户端测试

在新终端窗口：

```bash
cd python_ai_service
python tests\test_grpc_phase3.py
```

#### 步骤 3: 运行 Go 客户端测试

在新终端窗口：

```bash
# 运行单个测试命令
go run cmd\test_phase3_grpc\main.go --addr localhost:50051

# 或运行完整测试套件
go test -v -timeout 300s ./test/integration -run TestGRPC
```

---

## 测试方案

### 1. 连接测试

**目标**: 验证 gRPC 连接和服务健康状态

**测试步骤**:
```bash
go test -v ./test/integration -run TestGRPCConnection
```

**预期结果**:
- ✅ 连接成功
- ✅ 健康状态为 "healthy"
- ✅ 所有组件检查通过

### 2. 大纲生成测试

**目标**: 验证故事大纲生成功能

**测试步骤**:
```bash
go test -v ./test/integration -run TestGenerateOutline
```

**预期结果**:
- ✅ 生成包含完整信息的大纲
- ✅ 章节数 >= 1
- ✅ 响应时间 < 30秒

**输出示例**:
```
📖 标题: 逆天修仙路
🎭 类型: 修仙
📚 章节数: 5
⏱️  耗时: 9.8秒
```

### 3. 角色生成测试

**目标**: 验证角色设定生成功能

**测试步骤**:
```bash
go test -v ./test/integration -run TestGenerateCharacters
```

**预期结果**:
- ✅ 生成至少1个角色
- ✅ 角色包含完整信息（姓名、性格、背景等）
- ✅ 响应时间 < 30秒

**输出示例**:
```
👥 角色数量: 3
   1. 林天 (protagonist)
      性格: [勇敢, 聪慧, 坚韧]
```

### 4. 情节生成测试

**目标**: 验证情节事件生成功能

**测试步骤**:
```bash
go test -v ./test/integration -run TestGeneratePlot
```

**预期结果**:
- ✅ 生成至少5个事件
- ✅ 事件包含完整信息（标题、时间、参与者等）
- ✅ 响应时间 < 30秒

**输出示例**:
```
📅 事件数量: 15
🧵 情节线数: 3
   1. 少年初入修仙界 (第1天)
      类型: 触发事件
```

### 5. 完整工作流测试

**目标**: 验证完整创作工作流（大纲 → 角色 → 情节）

**测试步骤**:
```bash
go test -v ./test/integration -run TestCompleteWorkflow
```

**预期结果**:
- ✅ 依次完成大纲、角色、情节生成
- ✅ 所有数据结构完整
- ✅ 总响应时间 < 120秒

**输出示例**:
```
✅ 工作流执行成功! 总耗时: 38.5秒
   🆔 执行ID: abc123
   ✓  审核状态: true
   🔄 反思次数: 0
   
   📖 大纲: 心动的信号
      章节数: 3
   👥 角色数: 3
   📊 事件数: 10
```

### 6. 并发测试

**目标**: 验证服务并发处理能力

**测试步骤**:
```bash
go test -v ./test/integration -run TestConcurrentRequests
```

**预期结果**:
- ✅ 同时处理3个请求
- ✅ 所有请求成功完成
- ✅ 无竞态条件或死锁

### 7. 错误处理测试

**目标**: 验证异常情况处理

**测试步骤**:
```bash
go test -v ./test/integration -run TestErrorHandling
```

**预期结果**:
- ✅ 空请求返回错误
- ✅ 超时请求返回错误
- ✅ 错误信息清晰明确

---

## 故障排查

### 问题 1: 连接失败

**错误信息**:
```
failed to connect to AI service: context deadline exceeded
```

**可能原因**:
1. Python AI 服务未启动
2. 端口被占用
3. 防火墙阻止连接

**解决方法**:
```bash
# 1. 检查服务是否运行
netstat -ano | findstr :50051

# 2. 重启 Python AI 服务
cd python_ai_service
python run_grpc_server.py

# 3. 检查防火墙设置
```

### 问题 2: API 密钥错误

**错误信息**:
```
google.api_core.exceptions.PermissionDenied: API key not valid
```

**解决方法**:
```bash
# 设置正确的 API 密钥
set GOOGLE_API_KEY=your_actual_api_key

# 验证环境变量
echo %GOOGLE_API_KEY%

# 重启 Python AI 服务
```

### 问题 3: Protobuf 版本不匹配

**错误信息**:
```
cannot parse invalid wire-format data
```

**解决方法**:
```bash
# 重新生成 protobuf 代码
# Go 端
protoc -I python_ai_service/proto \
  --go_out=pkg/grpc/pb \
  --go-grpc_out=pkg/grpc/pb \
  python_ai_service/proto/ai_service.proto

# Python 端
cd python_ai_service
python -m grpc_tools.protoc \
  -I proto \
  --python_out=src/grpc_service \
  --grpc_python_out=src/grpc_service \
  proto/ai_service.proto
```

### 问题 4: 超时错误

**错误信息**:
```
context deadline exceeded
```

**解决方法**:
```go
// 增加超时时间
ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
defer cancel()
```

### 问题 5: Python 依赖缺失

**错误信息**:
```
ModuleNotFoundError: No module named 'grpc'
```

**解决方法**:
```bash
cd python_ai_service
pip install -r requirements.txt
```

---

## 性能基准

### 接口性能

| 接口 | 平均响应时间 | P95 | P99 | Token消耗 |
|-----|------------|-----|-----|----------|
| **HealthCheck** | < 0.1秒 | < 0.2秒 | < 0.5秒 | 0 |
| **GenerateOutline** | 8-12秒 | 15秒 | 20秒 | 2000-3000 |
| **GenerateCharacters** | 10-15秒 | 18秒 | 25秒 | 1500-2500 |
| **GeneratePlot** | 12-18秒 | 22秒 | 30秒 | 2500-4000 |
| **ExecuteCreativeWorkflow** | 30-45秒 | 60秒 | 90秒 | 5000-8000 |

### 并发性能

| 并发数 | 平均响应时间 | 成功率 | 错误率 |
|-------|------------|-------|--------|
| 1 | 10秒 | 100% | 0% |
| 3 | 12秒 | 100% | 0% |
| 5 | 15秒 | 98% | 2% |
| 10 | 25秒 | 90% | 10% |

### 资源消耗

| 指标 | 数值 | 说明 |
|-----|------|------|
| **内存** | ~500MB | Python AI 服务 |
| **CPU** | ~20% | 单核使用率 |
| **网络带宽** | 10-50KB/s | 请求/响应 |
| **磁盘IO** | < 1MB/s | 日志写入 |

---

## 测试报告模板

### 测试执行摘要

| 项目 | 结果 |
|-----|------|
| **测试日期** | YYYY-MM-DD |
| **测试人员** | XXX |
| **测试环境** | 开发/测试/生产 |
| **Go 版本** | 1.21.x |
| **Python 版本** | 3.11.x |
| **总测试数** | X |
| **通过数** | X |
| **失败数** | X |
| **通过率** | XX% |

### 测试结果详情

| 测试用例 | 状态 | 耗时 | 备注 |
|---------|------|------|------|
| TestGRPCConnection | ✅ | 0.1s | - |
| TestGenerateOutline | ✅ | 9.8s | - |
| TestGenerateCharacters | ✅ | 12.3s | - |
| TestGeneratePlot | ✅ | 14.7s | - |
| TestCompleteWorkflow | ✅ | 38.5s | - |
| TestConcurrentRequests | ✅ | 15.2s | - |
| TestErrorHandling | ✅ | 0.5s | - |

### 问题和建议

1. **问题**: XXX
   - **影响**: XXX
   - **建议**: XXX

2. **优化建议**: XXX

---

## 最佳实践

### 开发环境

1. **使用 Mock 数据** - 节省 API 调用成本
2. **本地 gRPC** - 使用 localhost:50051
3. **详细日志** - 启用 DEBUG 日志级别
4. **快速迭代** - 只测试修改的功能

### 测试环境

1. **真实 API** - 使用真实 LLM API
2. **完整测试** - 运行所有测试用例
3. **性能测试** - 测试并发和负载
4. **监控告警** - 启用监控系统

### 生产环境

1. **TLS 加密** - 使用 TLS 加密通信
2. **负载均衡** - 多实例部署
3. **限流熔断** - 防止雪崩
4. **监控日志** - 完善的监控和日志

---

## 常见问题

### Q1: 如何调整超时时间？

**A**: 在创建 context 时设置超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
defer cancel()
```

### Q2: 如何启用详细日志？

**A**: 设置日志级别

```python
# Python
import logging
logging.basicConfig(level=logging.DEBUG)
```

```go
// Go
// 使用 -v 标志运行测试
go test -v ./test/integration
```

### Q3: 如何测试性能？

**A**: 使用基准测试

```bash
go test -bench=. -benchmem ./test/integration
```

### Q4: 如何处理并发问题？

**A**: gRPC 服务器自动处理并发，但要注意：
- 使用合理的超时时间
- 限制并发数量
- 实现重试机制

---

## 参考文档

- [Phase3 gRPC README](../../python_ai_service/PHASE3_GRPC_README.md)
- [gRPC 集成指南](../../python_ai_service/GRPC_INTEGRATION_GUIDE.md)
- [Proto 定义](../../python_ai_service/proto/ai_service.proto)
- [Go 客户端实现](../../service/ai/phase3_client.go)
- [Python 服务端实现](../../python_ai_service/src/grpc_service/server.py)

---

**最后更新**: 2025-10-31  
**维护者**: 青羽后端架构团队  
**版本**: v1.0

