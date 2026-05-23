# Phase3 gRPC集成指南

## Page Role

- phase3-legacy-guide
- current-owner: `design/ai/`
- current-bounded: Phase3 阶段内的 gRPC 集成操作指南，只用于历史集成路径回看

## Recommended Read Path

1. 先读 `README.md`。
2. 再读 `phase3/README.md`。
3. 需要看当时的 gRPC 集成步骤时，再读本文件。

## Boundary

- 这是 Phase3 历史集成指南，不是 today owner。
- 当前 AI 设计判断应回到总入口和主线文档。
- 文中的启动方式、目录结构和 Agent 清单带有明显时间性。

## Quick Section Map

- 概述
- 架构
- 文件结构
- 集成步骤
- 验证方式

## Quick Takeaways

- 本文回答“Phase3 当时如何把 Python AI 服务接进 gRPC”。
- 它不负责定义当前的跨服务集成边界。

## Skip Guide

- 只看当前主线：跳过本文件。
- 只看当时接线步骤：重点看架构、文件结构和集成步骤。

## 📋 概述

本指南介绍如何将Phase3专业Agent集成到gRPC服务中，以便Go后端可以通过gRPC调用AI服务。

## 🏗️ 架构

```
Go后端 (Qingyu_backend)
    ↓ gRPC调用
Python AI服务 (python_ai_service)
    ↓ 调用
Phase3 Agents (OutlineAgent, CharacterAgent, PlotAgent)
    ↓ 调用
Gemini 2.0 Flash API
```

## 📁 文件结构

```
python_ai_service/
├── proto/
│   └── ai_service.proto          # 更新后的protobuf定义
├── src/
│   ├── agents/
│   │   └── specialized/          # Phase3专业Agents
│   │       ├── outline_agent.py
│   │       ├── character_agent.py
│   │       └── plot_agent.py
│   └── grpc_service/             # 新增：gRPC服务模块
│       ├── __init__.py
│       ├── converters.py         # 数据转换工具
│       ├── ai_servicer.py        # gRPC服务实现
│       ├── server.py             # 服务器启动脚本
│       ├── ai_service_pb2.py     # 生成的protobuf代码
│       └── ai_service_pb2_grpc.py # 生成的gRPC代码
└── scripts/
    └── generate_grpc_proto.bat   # 代码生成脚本
```

## 🚀 快速开始

### 1. 生成Protobuf代码

**Windows:**

```bash
cd python_ai_service
scripts\generate_grpc_proto.bat
```

**Linux/Mac:**

```bash
cd python_ai_service
python -m grpc_tools.protoc \
    -I proto \
    --python_out=src/grpc_service \
    --grpc_python_out=src/grpc_service \
    proto/ai_service.proto
```

### 2. 启动gRPC服务器

```bash
cd python_ai_service
set GOOGLE_API_KEY=your_api_key_here
python src/grpc_service/server.py --host 0.0.0.0 --port 50051
```

### 3. 测试服务

创建测试脚本 `test_grpc_client.py`:

```python
import grpc
import asyncio
from grpc_service import ai_service_pb2, ai_service_pb2_grpc


async def test_outline_generation():
    """测试大纲生成"""
    async with grpc.aio.insecure_channel('localhost:50051') as channel:
        stub = ai_service_pb2_grpc.AIServiceStub(channel)
        
        request = ai_service_pb2.OutlineRequest(
            task="创作一个修仙小说大纲，主角是天才少年",
            user_id="test_user",
            project_id="test_project"
        )
        
        response = await stub.GenerateOutline(request)
        
        print(f"✅ 大纲生成成功")
        print(f"📖 标题: {response.outline.title}")
        print(f"📚 章节数: {len(response.outline.chapters)}")
        print(f"⏱️  耗时: {response.execution_time:.2f}秒")


async def test_creative_workflow():
    """测试完整创作工作流"""
    async with grpc.aio.insecure_channel('localhost:50051') as channel:
        stub = ai_service_pb2_grpc.AIServiceStub(channel)
        
        request = ai_service_pb2.CreativeWorkflowRequest(
            task="创作一个现代都市爱情小说的完整设定",
            user_id="test_user",
            project_id="test_project",
            max_reflections=3,
            enable_human_review=False
        )
        
        response = await stub.ExecuteCreativeWorkflow(request)
        
        print(f"✅ 工作流执行成功")
        print(f"📖 大纲: {response.outline.title}")
        print(f"👥 角色数: {len(response.characters.characters)}")
        print(f"📊 情节事件数: {len(response.plot.timeline_events)}")
        print(f"⏱️  总耗时: {sum(response.execution_times.values()):.2f}秒")


if __name__ == "__main__":
    asyncio.run(test_outline_generation())
    asyncio.run(test_creative_workflow())
```

## 📡 gRPC接口说明

### 1. ExecuteCreativeWorkflow - 完整创作工作流

执行 Outline → Character → Plot 完整流程

**请求:**

```protobuf
message CreativeWorkflowRequest {
  string task = 1;                      // 创作任务描述
  string user_id = 2;                   // 用户ID
  string project_id = 3;                // 项目ID
  int32 max_reflections = 4;            // 最大反思次数
  bool enable_human_review = 5;         // 是否启用人工审核
  map<string, string> workspace_context = 6;  // 工作区上下文
}
```

**响应:**

```protobuf
message CreativeWorkflowResponse {
  string execution_id = 1;              // 执行ID
  bool review_passed = 2;               // 审核是否通过
  int32 reflection_count = 3;           // 反思次数
  OutlineData outline = 4;              // 大纲数据
  CharactersData characters = 5;        // 角色数据
  PlotData plot = 6;                    // 情节数据
  DiagnosticReportData diagnostic_report = 7;  // 诊断报告
  repeated string reasoning = 8;        // 推理链
  map<string, float> execution_times = 9;  // 执行时间
  int32 tokens_used = 10;               // Token使用量
}
```

### 2. GenerateOutline - 生成大纲

生成故事大纲

**请求:**

```protobuf
message OutlineRequest {
  string task = 1;                      // 任务描述
  string user_id = 2;                   // 用户ID
  string project_id = 3;                // 项目ID
  map<string, string> workspace_context = 4;  // 工作区上下文
  string correction_prompt = 5;         // 修正提示（可选）
}
```

**响应:**

```protobuf
message OutlineResponse {
  OutlineData outline = 1;              // 大纲数据
  float execution_time = 2;             // 执行时间
}
```

### 3. GenerateCharacters - 生成角色

基于大纲生成角色

**请求:**

```protobuf
message CharactersRequest {
  string task = 1;                      // 任务描述
  string user_id = 2;                   // 用户ID
  string project_id = 3;                // 项目ID
  OutlineData outline = 4;              // 大纲数据（必需）
  map<string, string> workspace_context = 5;  // 工作区上下文
  string correction_prompt = 6;         // 修正提示（可选）
}
```

**响应:**

```protobuf
message CharactersResponse {
  CharactersData characters = 1;        // 角色数据
  float execution_time = 2;             // 执行时间
}
```

### 4. GeneratePlot - 生成情节

基于大纲和角色生成情节

**请求:**

```protobuf
message PlotRequest {
  string task = 1;                      // 任务描述
  string user_id = 2;                   // 用户ID
  string project_id = 3;                // 项目ID
  OutlineData outline = 4;              // 大纲数据（必需）
  CharactersData characters = 5;        // 角色数据（必需）
  map<string, string> workspace_context = 6;  // 工作区上下文
  string correction_prompt = 7;         // 修正提示（可选）
}
```

**响应:**

```protobuf
message PlotResponse {
  PlotData plot = 1;                    // 情节数据
  float execution_time = 2;             // 执行时间
}
```

### 5. HealthCheck - 健康检查

检查服务健康状态

**请求:**

```protobuf
message HealthCheckRequest {}
```

**响应:**

```protobuf
message HealthCheckResponse {
  string status = 1;                    // healthy/degraded/unhealthy
  map<string, string> checks = 2;       // 各组件检查结果
}
```

## 🔧 Go后端集成

### 1. 生成Go代码

```bash
cd Qingyu_backend
protoc -I python_ai_service/proto \
    --go_out=pkg/grpc/pb \
    --go-grpc_out=pkg/grpc/pb \
    python_ai_service/proto/ai_service.proto
```

### 2. Go客户端示例

```go
package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "Qingyu_backend/pkg/grpc/pb"
)

func main() {
    // 连接gRPC服务
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer conn.Close()

    client := pb.NewAIServiceClient(conn)

    // 调用大纲生成
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    request := &pb.OutlineRequest{
        Task:      "创作一个科幻小说大纲",
        UserId:    "user123",
        ProjectId: "project456",
    }

    response, err := client.GenerateOutline(ctx, request)
    if err != nil {
        log.Fatalf("调用失败: %v", err)
    }

    log.Printf("✅ 大纲生成成功")
    log.Printf("📖 标题: %s", response.Outline.Title)
    log.Printf("📚 章节数: %d", len(response.Outline.Chapters))
    log.Printf("⏱️  耗时: %.2f秒", response.ExecutionTime)
}
```

### 3. 集成到Go Service层

在 `service/ai/` 中创建 `phase3_client.go`:

```go
package ai

import (
    "context"
    "fmt"

    pb "Qingyu_backend/pkg/grpc/pb"
    "google.golang.org/grpc"
)

// Phase3Client Phase3 AI服务客户端
type Phase3Client struct {
    client pb.AIServiceClient
    conn   *grpc.ClientConn
}

// NewPhase3Client 创建Phase3客户端
func NewPhase3Client(address string) (*Phase3Client, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, fmt.Errorf("连接AI服务失败: %w", err)
    }

    return &Phase3Client{
        client: pb.NewAIServiceClient(conn),
        conn:   conn,
    }, nil
}

// GenerateOutline 生成大纲
func (c *Phase3Client) GenerateOutline(ctx context.Context, task, userID, projectID string) (*pb.OutlineResponse, error) {
    request := &pb.OutlineRequest{
        Task:      task,
        UserId:    userID,
        ProjectId: projectID,
    }

    return c.client.GenerateOutline(ctx, request)
}

// ExecuteCreativeWorkflow 执行完整创作工作流
func (c *Phase3Client) ExecuteCreativeWorkflow(ctx context.Context, req *pb.CreativeWorkflowRequest) (*pb.CreativeWorkflowResponse, error) {
    return c.client.ExecuteCreativeWorkflow(ctx, req)
}

// Close 关闭连接
func (c *Phase3Client) Close() error {
    return c.conn.Close()
}
```

## 🐛 故障排查

### 问题1: protobuf代码生成失败

**解决方案:**

```bash
pip install --upgrade grpcio-tools
```

### 问题2: 导入路径错误

**解决方案:**

修改生成的 `ai_service_pb2_grpc.py`:

```python
# 将
import ai_service_pb2 as ai__service__pb2

# 改为
from . import ai_service_pb2 as ai__service__pb2
```

### 问题3: gRPC服务启动失败

**检查:**

1. 端口是否被占用: `netstat -ano | findstr 50051`
2. API密钥是否设置: `echo %GOOGLE_API_KEY%`
3. 依赖是否安装: `pip list | findstr grpc`

## 📊 性能指标

| 接口 | 平均耗时 | 备注 |
|-----|---------|------|
| GenerateOutline | 8-12秒 | 生成5-10章大纲 |
| GenerateCharacters | 10-15秒 | 生成3-5个主要角色 |
| GeneratePlot | 12-18秒 | 生成15-25个情节事件 |
| ExecuteCreativeWorkflow | 30-45秒 | 完整流程 |

## 🔐 安全建议

1. **生产环境使用TLS**:

```python
# server.py
server_credentials = grpc.ssl_server_credentials(...)
server.add_secure_port(server_address, server_credentials)
```

2. **API密钥管理**:

使用环境变量或密钥管理服务，不要硬编码

3. **访问控制**:

添加认证中间件验证请求来源

## 📚 参考资料

- [gRPC Python文档](https://grpc.io/docs/languages/python/)
- [Protocol Buffers文档](https://developers.google.com/protocol-buffers)
- [Phase3 Agent设计文档](./phase3/README.md)

## 🎯 下一步

- [ ] 添加反思循环集成
- [ ] 实现流式响应（用于长文本生成）
- [ ] 添加Redis缓存
- [ ] 集成监控和日志
- [ ] 性能优化和压测

---

**最后更新**: 2025-10-30
**维护者**: 青羽后端架构团队

