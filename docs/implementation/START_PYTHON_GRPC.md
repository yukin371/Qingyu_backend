# Python gRPC服务器启动指南

## 问题诊断

Python gRPC服务器在后台启动时没有正确运行。这是因为PowerShell的后台任务机制可能导致进程立即退出。

## 手动启动步骤（推荐）

### 方法1：使用新的PowerShell窗口

打开**新的PowerShell窗口**，运行：

```powershell
cd E:\Github\Qingyu\Qingyu_backend\python_ai_service
python quick_test_grpc.py
```

**预期输出**：
```
============================================================
Starting gRPC Server Test
============================================================
✅ grpc imported
✅ ai_service_pb2 imported
✅ ai_service_pb2_grpc imported

🚀 Starting server on port 50052...
✅ Server is RUNNING on port 50052
🔗 Test with: go run test_grpc_connection.go
```

保持这个窗口打开！服务器会一直运行。

### 方法2：使用Windows Terminal（如果已安装）

```powershell
wt -d E:\Github\Qingyu\Qingyu_backend\python_ai_service python quick_test_grpc.py
```

### 方法3：使用批处理文件

创建`python_ai_service/start_server.bat`：
```batch
@echo off
cd /d E:\Github\Qingyu\Qingyu_backend\python_ai_service
python quick_test_grpc.py
pause
```

双击运行此批处理文件。

## 验证服务器运行

在另一个PowerShell窗口：

```powershell
# 检查端口
netstat -ano | findstr ":50052"
# 应显示：TCP    0.0.0.0:50052   ...   LISTENING

# 运行Go测试
cd E:\Github\Qingyu\Qingyu_backend
go run test_grpc_connection.go
```

## 如果仍然失败

### 检查Python版本
```powershell
python --version
# 应为Python 3.10+
```

### 检查依赖
```powershell
pip show grpcio grpcio-reflection
```

### 防火墙设置
Windows Defender可能阻止端口50052，添加入站规则：
```powershell
# 以管理员身份运行
netsh advfirewall firewall add rule name="Python gRPC" dir=in action=allow protocol=TCP localport=50052
```

### 查看完整错误
如果服务器启动后立即退出，查看错误信息：
```powershell
cd python_ai_service
python quick_test_grpc.py 2>&1 | Tee-Object -FilePath grpc_error.log
```

## 成功标志

✅ Python窗口显示"Server is RUNNING"
✅ netstat显示端口50052在LISTENING
✅ Go测试输出"gRPC连接成功"

## 下一步

一旦gRPC通信成功，您可以：
1. 集成FastAPI和gRPC（同时运行）
2. 实现完整的AI功能（RAG、Agent）
3. 使用Docker Compose部署

