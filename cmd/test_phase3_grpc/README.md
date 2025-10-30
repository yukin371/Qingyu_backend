# Phase3 gRPC客户端测试工具

这是一个用于测试Phase3 AI服务gRPC连接的命令行工具。

## 📋 前提条件

1. **Python AI服务已启动**

```bash
cd python_ai_service
$env:GOOGLE_API_KEY="your_api_key"
python run_grpc_server.py
```

2. **Go依赖已安装**

```bash
go mod tidy
```

## 🚀 使用方法

### 方式1: 测试单个Agent（快速）

```bash
# 从项目根目录运行
go run cmd/test_phase3_grpc/main.go -task "创作一个修仙小说大纲"
```

**预期输出**:
```
========================================
Phase3 gRPC客户端测试
========================================

连接到gRPC服务器: localhost:50051
✅ 连接成功

1️⃣  健康检查...
   状态: healthy
   组件状态:
     - outline_agent: healthy
     - character_agent: healthy
     - plot_agent: healthy

2️⃣  生成大纲...
   任务: 创作一个修仙小说大纲
   ✅ 成功! 耗时: 9.8秒
   📖 标题: 天命之子
   🎭 类型: 修仙
   📚 章节数: 5
   ...
```

### 方式2: 测试完整工作流（慢，30-60秒）

```bash
go run cmd/test_phase3_grpc/main.go -workflow -task "创作都市爱情小说设定"
```

### 方式3: 自定义服务器地址

```bash
go run cmd/test_phase3_grpc/main.go -addr "192.168.1.100:50051" -task "你的任务"
```

## 🧪 运行Go单元测试

```bash
# 测试客户端连接
go test ./service/ai -run TestPhase3Client_HealthCheck -v

# 测试大纲生成
go test ./service/ai -run TestPhase3Client_GenerateOutline -v

# 测试完整工作流（长测试，跳过）
go test ./service/ai -run TestPhase3Client_ExecuteCreativeWorkflow -v

# 运行所有测试（跳过长测试）
go test ./service/ai -short -v

# 性能测试
go test ./service/ai -bench=. -benchmem
```

## 📊 命令行参数

| 参数 | 默认值 | 说明 |
|-----|--------|------|
| `-addr` | `localhost:50051` | gRPC服务器地址 |
| `-task` | `"创作一个修仙小说大纲..."` | 创作任务描述 |
| `-workflow` | `false` | 是否执行完整工作流 |

## 🎯 测试场景

### 场景1: 健康检查

```bash
# 编译并运行
go build -o test_phase3_grpc.exe cmd/test_phase3_grpc/main.go
./test_phase3_grpc.exe
```

### 场景2: 大纲生成

```bash
go run cmd/test_phase3_grpc/main.go \
  -task "创作一个科幻小说大纲，主角是星际探险家，5章"
```

### 场景3: 完整工作流

```bash
go run cmd/test_phase3_grpc/main.go \
  -workflow \
  -task "创作一个悬疑推理小说的完整设定，3章"
```

## 🐛 故障排查

### 问题1: 连接失败

```
❌ 连接失败: context deadline exceeded
```

**解决方案**:
1. 检查Python AI服务是否启动
2. 检查端口50051是否正确
3. 检查防火墙设置

### 问题2: 健康检查失败

```
❌ 健康检查失败: rpc error
```

**解决方案**:
1. 确认API密钥已设置
2. 查看Python服务日志
3. 重启Python服务

### 问题3: 超时错误

```
❌ 大纲生成失败: context deadline exceeded
```

**解决方案**:
1. 增加超时时间（修改`phase3_client.go`）
2. 检查网络连接
3. 查看Python服务日志

## 📝 示例输出

### 单Agent测试输出

```
========================================
Phase3 gRPC客户端测试
========================================

连接到gRPC服务器: localhost:50051
✅ 连接成功

1️⃣  健康检查...
   状态: healthy

2️⃣  生成大纲...
   ✅ 成功! 耗时: 9.80秒
   📖 标题: 天命之子
   📚 章节数: 5

3️⃣  生成角色...
   ✅ 成功! 耗时: 12.30秒
   👥 角色数: 3

4️⃣  生成情节...
   ✅ 成功! 耗时: 14.70秒
   📅 事件数: 18

========================================
✅ 测试完成
========================================
```

### 完整工作流输出

```
🎨 执行完整创作工作流...
   ⏳ 这可能需要30-60秒...

============================================================
✅ 工作流执行成功!
============================================================
🆔 执行ID: uuid-xxx
✓  审核状态: true
🔄 反思次数: 0
⏱️  总耗时: 38.50秒

📖 大纲:
   标题: 心动的信号
   章节数: 3

👥 角色:
   角色数: 2

📊 情节:
   事件数: 15

⏱️  执行时间分析:
   outline: 9.80秒
   character: 12.30秒
   plot: 14.70秒
   总计: 36.80秒
```

## 🔗 相关文档

- [Phase3 gRPC集成指南](../../python_ai_service/GRPC_INTEGRATION_GUIDE.md)
- [Phase3完成报告](../../doc/implementation/00进度指导/Phase3_gRPC集成完成报告_2025-10-30.md)
- [API文档](../../python_ai_service/PHASE3_GRPC_README.md)

