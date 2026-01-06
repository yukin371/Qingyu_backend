# Windows 快速开始指南

> **针对 Windows 用户的 Phase3 v2.0 实施指南**

---

## ✅ 已解决的问题

Windows 系统没有 `make` 命令，我们提供了 PowerShell 脚本来替代。

---

## 🚀 快速开始（Windows）

### 步骤 1: 生成 Protobuf 代码

您已经安装了 `protoc`（版本 33.0-rc2），接下来：

#### 1.1 安装 Go 插件

```powershell
# 在项目根目录执行
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### 1.2 生成所有 Protobuf 代码

```powershell
# 使用 PowerShell 脚本（推荐）
.\scripts\generate_proto_all.ps1
```

**或者分别生成**：

```powershell
# 生成 Go 代码
.\scripts\generate_proto_go.ps1

# 生成 Python 代码
.\scripts\generate_proto_python.ps1
```

**预期输出**：
```
=== Generating All Protobuf Code ===

[1/2] Generating Go protobuf code...
✓ Go protobuf code generated successfully in pkg\grpc\pb\

[2/2] Generating Python protobuf code...
✓ Python protobuf code generated successfully in src\grpc_server\
✓ Import paths fixed

=== All protobuf code generated successfully! ===
```

---

### 步骤 2: 安装 Python 依赖

```powershell
cd python_ai_service

# 检查是否安装了 Poetry
poetry --version

# 如果没有，安装 Poetry
pip install poetry

# 安装依赖
poetry install
```

---

### 步骤 3: 配置环境变量

```powershell
# 复制示例配置
Copy-Item .env.example .env

# 使用记事本编辑（或 VSCode）
notepad .env
# 或
code .env
```

**最低配置**：
```env
# 至少配置一个 AI 提供商
OPENAI_API_KEY=your_key_here
# 或
ANTHROPIC_API_KEY=your_key_here

# 其他保持默认即可
```

---

### 步骤 4: 启动 Python 服务

```powershell
# 方式 1: 使用 Poetry（推荐）
poetry run uvicorn src.main:app --reload --host 0.0.0.0 --port 8000

# 方式 2: 使用批处理脚本
.\run.bat
```

**验证服务**：

打开浏览器访问：
- API 文档: http://localhost:8000/docs
- 健康检查: http://localhost:8000/api/v1/health

或使用 PowerShell：
```powershell
Invoke-WebRequest -Uri http://localhost:8000/api/v1/health | Select-Object -ExpandProperty Content
```

---

### 步骤 5: 运行测试

```powershell
# 在 python_ai_service 目录下
poetry run pytest tests/ -v
```

---

## 📂 生成的文件位置

生成成功后，您应该看到以下文件：

### Go 代码
```
pkg\grpc\pb\
├── ai_service.pb.go
└── ai_service_grpc.pb.go
```

### Python 代码
```
python_ai_service\src\grpc_server\
├── ai_service_pb2.py
└── ai_service_pb2_grpc.py
```

---

## 🐛 Windows 特定问题

### Q1: PowerShell 提示"无法运行脚本"

**错误信息**：
```
.\scripts\generate_proto_all.ps1 : 无法加载文件，因为在此系统上禁止运行脚本。
```

**解决方案**：
```powershell
# 临时允许运行脚本（仅当前会话）
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process

# 然后再运行脚本
.\scripts\generate_proto_all.ps1
```

或者永久设置（需要管理员权限）：
```powershell
# 以管理员身份运行 PowerShell
Set-ExecutionPolicy RemoteSigned
```

---

### Q2: protoc 找不到

**错误信息**：
```
protoc : 无法将"protoc"项识别为 cmdlet、函数、脚本文件或可运行程序的名称。
```

**解决方案**：

1. 从 [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases) 下载最新版本
2. 解压到某个目录，例如 `C:\protoc`
3. 将 `C:\protoc\bin` 添加到系统 PATH
4. 重启 PowerShell

---

### Q3: Go 插件找不到

**错误信息**：
```
'protoc-gen-go' 不是内部或外部命令
```

**解决方案**：

确保 `%USERPROFILE%\go\bin` 在 PATH 中：

```powershell
# 查看当前 PATH
$env:Path

# 临时添加（仅当前会话）
$env:Path += ";$env:USERPROFILE\go\bin"

# 永久添加：
# 控制面板 → 系统 → 高级系统设置 → 环境变量
# 在 Path 中添加：%USERPROFILE%\go\bin
```

---

### Q4: Python 找不到

**错误信息**：
```
python : 无法将"python"项识别为 cmdlet
```

**解决方案**：

安装 Python 3.10+ 并确保添加到 PATH：
- 下载：https://www.python.org/downloads/
- 安装时勾选 "Add Python to PATH"

---

### Q5: Poetry 安装失败

**解决方案**：

```powershell
# 使用 pip 安装
pip install poetry

# 或使用官方安装脚本
(Invoke-WebRequest -Uri https://install.python-poetry.org -UseBasicParsing).Content | py -
```

---

## 📝 Windows vs Linux/macOS 命令对照

| 任务 | Linux/macOS | Windows (PowerShell) |
|-----|-------------|---------------------|
| 生成 Protobuf | `make proto` | `.\scripts\generate_proto_all.ps1` |
| 复制文件 | `cp .env.example .env` | `Copy-Item .env.example .env` |
| 查看文件 | `cat .env` | `Get-Content .env` |
| 编辑文件 | `vim .env` | `notepad .env` 或 `code .env` |
| 测试连接 | `curl http://localhost:8000` | `Invoke-WebRequest http://localhost:8000` |
| 查找进程 | `lsof -i :8000` | `netstat -ano \| findstr :8000` |

---

## ✅ 验证清单

在继续下一步之前，确认：

- [x] protoc 已安装（`protoc --version`）
- [x] Go 插件已安装（`protoc-gen-go` 和 `protoc-gen-go-grpc`）
- [ ] Protobuf 代码已生成（检查 `pkg\grpc\pb\` 和 `src\grpc_server\`）
- [ ] Python 依赖已安装（`poetry install`）
- [ ] 环境变量已配置（`.env` 文件）
- [ ] Python 服务可以启动
- [ ] 健康检查 API 正常

---

## 🎯 下一步

完成上述步骤后，继续阅读 [`NEXT_STEPS_PHASE3.md`](NEXT_STEPS_PHASE3.md) 中的**步骤 5**（部署 Milvus）。

---

**祝顺利！** 如有问题，请参考 [`scripts/README.md`](scripts/README.md) 获取更多帮助。

---

**提示**：Windows 用户建议使用 **Windows Terminal** 或 **VSCode 集成终端**，体验更好！

