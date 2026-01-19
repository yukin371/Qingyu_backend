# 构建脚本说明

本目录包含项目的构建和工具脚本。

## Protobuf 代码生成

### Linux / macOS

使用 Makefile：

```bash
# 生成所有 Protobuf 代码
make proto

# 或分别生成
make proto-go      # 生成 Go 代码
make proto-python  # 生成 Python 代码
```

### Windows (PowerShell)

使用 PowerShell 脚本：

```powershell
# 生成所有 Protobuf 代码
.\scripts\generate_proto_all.ps1

# 或分别生成
.\scripts\generate_proto_go.ps1      # 生成 Go 代码
.\scripts\generate_proto_python.ps1  # 生成 Python 代码
```

### 手动生成（跨平台）

#### Go 代码

```bash
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  -I python_ai_service/proto \
  python_ai_service/proto/ai_service.proto
```

#### Python 代码

```bash
cd python_ai_service
python -m grpc_tools.protoc -I proto \
  --python_out=src/grpc_server \
  --grpc_python_out=src/grpc_server \
  proto/ai_service.proto
```

## 前置条件

### 1. Protocol Buffers 编译器

**Linux (Ubuntu/Debian)**:
```bash
sudo apt-get install -y protobuf-compiler
```

**macOS**:
```bash
brew install protobuf
```

**Windows**:
从 [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases) 下载并安装

### 2. Go 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 3. Python gRPC 工具

```bash
pip install grpcio-tools
```

## 验证生成结果

生成后应该看到以下文件：

```
pkg/grpc/pb/
├── ai_service.pb.go
└── ai_service_grpc.pb.go

python_ai_service/src/grpc_server/
├── ai_service_pb2.py
└── ai_service_pb2_grpc.py
```

## 常见问题

### Q: Windows 提示 "protoc 不是内部或外部命令"

**A**: 需要将 protoc.exe 添加到 PATH 环境变量中。

### Q: Go 插件找不到

**A**: 确保 `$GOPATH/bin` (或 `$HOME/go/bin`) 在 PATH 中：
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Q: Python 生成失败

**A**: 确保安装了 grpcio-tools：
```bash
pip install --upgrade grpcio-tools
```
