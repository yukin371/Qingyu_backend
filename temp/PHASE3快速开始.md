# 🚀 Phase3 Go客户端快速开始

> **5分钟完成测试！**

---

## ✅ 前提条件

- [x] Go 1.21+ 已安装
- [x] Python 3.11+ 已安装
- [x] Gemini API Key 已获取

---

## 📝 第1步: 启动Python gRPC服务器（终端1）

```powershell
cd python_ai_service

# 设置API密钥
$env:GOOGLE_API_KEY="AIzaSyD-q07WgZPd8mw4f1hVKVw44yXNUWrAuOk"

# 启动服务器
python run_grpc_server.py
```

**期望看到**:
```
========================================
Phase3 gRPC Server Startup
========================================

API Key is set
...
✅ gRPC服务器就绪，等待请求...
```

---

## 🧪 第2步: 运行Go测试（终端2，新开）

### 方式1: 使用编译好的程序

```powershell
cd E:\Github\Qingyu\Qingyu_backend

# 运行测试
.\test_phase3_grpc.exe
```

### 方式2: 直接运行Go代码

```powershell
go run cmd/test_phase3_grpc/main.go
```

### 方式3: 测试完整工作流

```powershell
.\test_phase3_grpc.exe -workflow
```

---

## 📊 期望输出

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

---

## 🧪 运行Go单元测试

```powershell
# 测试健康检查
go test ./service/ai -run TestPhase3Client_HealthCheck -v

# 测试大纲生成  
go test ./service/ai -run TestPhase3Client_GenerateOutline -v

# 测试完整工作流（需要30-60秒）
go test ./service/ai -run TestPhase3Client_ExecuteCreativeWorkflow -v
```

---

## 🎯 核心文件

| 文件 | 说明 |
|-----|------|
| `service/ai/phase3_client.go` | gRPC客户端 |
| `service/ai/phase3_client_test.go` | 单元测试 |
| `cmd/test_phase3_grpc/main.go` | 命令行工具 |
| `pkg/grpc/pb/*.go` | Protobuf生成代码 |

---

## 🐛 常见问题

### Q: 连接失败

```
connection error: desc = "transport: Error while dialing"
```

**A**: 确保Python gRPC服务器正在运行（终端1）

### Q: API密钥错误

**A**: 检查环境变量
```powershell
echo $env:GOOGLE_API_KEY
```

### Q: 编译错误

**A**: 重新生成Protobuf代码
```powershell
protoc -I python_ai_service/proto --go_out=. --go-grpc_out=. python_ai_service/proto/ai_service.proto
```

---

## 📚 详细文档

- [完整集成总结](doc/implementation/00进度指导/Phase3_Go集成完成总结_2025-10-30.md)
- [Python gRPC指南](python_ai_service/GRPC_INTEGRATION_GUIDE.md)
- [测试工具README](cmd/test_phase3_grpc/README.md)

---

**就这么简单！享受Phase3 AI能力吧！** 🎉

