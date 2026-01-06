# gRPC 客户端

Go 端的 gRPC 客户端，用于与 Python AI Service 通信。

## 目录结构

```
pkg/grpc/
├── pb/                 # Protobuf 生成的代码
│   ├── ai_service.pb.go
│   └── ai_service_grpc.pb.go
├── client.go           # gRPC 客户端封装
└── README.md          # 本文件
```

## 生成 Protobuf 代码

从项目根目录执行：

```bash
# 安装 protoc 插件（如果未安装）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码
protoc --go_out=. --go-grpc_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  -I python_ai_service/proto \
  python_ai_service/proto/ai_service.proto
```

## 使用示例

```go
package main

import (
    "context"
    "Qingyu_backend/pkg/grpc"
    "Qingyu_backend/pkg/grpc/pb"
)

func main() {
    // 创建客户端
    client, err := grpc.NewAIClient("localhost:50051")
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 生成内容
    resp, err := client.GenerateContent(context.Background(), &pb.GenerateContentRequest{
        ProjectId: "project-123",
        Prompt: "生成一段开头",
    })
    if err != nil {
        panic(err)
    }
    
    println(resp.Content)
}
```

