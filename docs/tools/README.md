# 工具文档入口

> 最后整理: 2026-05-22  
> 当前状态: `current-bounded`

本目录保存构建工具、迁移工具和脚本使用说明。当前入口以本页为准，不再只把它当成单一的 Protobuf 构建说明。

## Recommended Read Path

1. [README.md](./README.md)
2. [QUICKSTART.md](./QUICKSTART.md)
3. [README_MIGRATION.md](./README_MIGRATION.md)
4. [使用说明.md](./使用说明.md)

## Current Materials

- [QUICKSTART.md](./QUICKSTART.md): 工具快速开始
- [README_MIGRATION.md](./README_MIGRATION.md): 迁移相关工具说明
- [使用说明.md](./使用说明.md): Scripts 目录使用说明

## Boundary

- 本目录：构建与工具使用说明
- `Qingyu_backend/scripts/`: 实际脚本 owner
- `Qingyu_backend/docs/testing/README.md`: 测试文档入口
- `Qingyu_backend/docs/migration/README.md`: 文档层迁移说明入口

## Legacy Detailed Manual

以下内容保留为较早一版“构建脚本说明”长文正文，可继续回看，但本页顶部现在才是目录入口层。

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
