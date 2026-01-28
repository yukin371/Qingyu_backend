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

---

# 数据库迁移与索引验证指南

本文档说明如何执行数据库迁移并验证索引是否正确创建。

## 前置要求

1. **MongoDB 测试环境**
   - 确保 MongoDB 已启动
   - 默认连接: `mongodb://localhost:27017`

2. **Go 环境**
   - Go >= 1.22
   - 已安装项目依赖

## 快速开始

### 1. 设置环境变量

```bash
# Windows (PowerShell)
$env:MONGO_URI="mongodb://localhost:27017"

# Linux/macOS
export MONGO_URI="mongodb://localhost:27017"
```

### 2. 执行迁移到测试环境

```bash
# 方法1: 使用迁移工具（推荐）
go run cmd/migrate/main.go --env=test

# 方法2: 手动执行单个迁移
go run migration/mongodb/002_create_users_indexes.go
go run migration/mongodb/003_create_books_indexes_p0.go
go run migration/mongodb/004_create_chapters_indexes.go
go run migration/mongodb/005_create_reading_progress_indexes.go
```

### 3. 运行索引验证测试

```bash
# 运行所有索引验证测试
go test -v ./scripts/verify_indexes_test.go

# 跳过集成测试（快速检查编译）
go test -short -v ./scripts/verify_indexes_test.go

# 列出所有索引（调试用）
go test -v -run TestListAllIndexes ./scripts/verify_indexes_test.go
```

## 手动验证索引

### 使用 mongosh 验证

#### 1. 连接到 MongoDB

```bash
mongosh "mongodb://localhost:27017"
```

#### 2. 切换到测试数据库

```javascript
use qingyu_test
```

#### 3. 验证 Users 集合索引

```javascript
db.users.getIndexes()

// 预期索引:
// - _id_ (默认)
// - status_1_created_at_-1
// - roles_1
// - last_login_at_-1
```

#### 4. 验证 Books 集合索引

```javascript
db.books.getIndexes()

// 预期索引:
// - _id_ (默认)
// - status_1_created_at_-1
// - status_1_rating_-1
// - author_id_1_status_1_created_at_-1
// - category_ids_1_rating_-1
// - is_completed_1_status_1
```

#### 5. 验证 Chapters 集合索引

```javascript
db.chapters.getIndexes()

// 预期索引:
// - _id_ (默认)
// - book_id_1_chapter_num_1
// - book_id_1_status_1_chapter_num_1
```

#### 6. 验证 ReadingProgress 集合索引

```javascript
db.reading_progress.getIndexes()

// 预期索引:
// - _id_ (默认)
// - user_id_1_book_id_1
// - user_id_1_last_read_at_-1
// - book_id_1
```

## 索引说明

### Users 集合

| 索引名称 | 字段 | 用途 |
|---------|------|------|
| `status_1_created_at_-1` | status, created_at | 按状态和创建时间查询用户 |
| `roles_1` | roles | 按角色查询用户 |
| `last_login_at_-1` | last_login_at | 查询最近登录用户 |

### Books 集合 (P0)

| 索引名称 | 字段 | 用途 |
|---------|------|------|
| `status_1_created_at_-1` | status, created_at | 书籍列表查询 |
| `status_1_rating_-1` | status, rating | 按评分排序书籍 |
| `author_id_1_status_1_created_at_-1` | author_id, status, created_at | 作者书籍查询 |
| `category_ids_1_rating_-1` | category_ids, rating | 分类书籍查询 |
| `is_completed_1_status_1` | is_completed, status | 已完结书籍查询 |

### Chapters 集合

| 索引名称 | 字段 | 用途 |
|---------|------|------|
| `book_id_1_chapter_num_1` | book_id, chapter_num | 书籍章节查询 |
| `book_id_1_status_1_chapter_num_1` | book_id, status, chapter_num | 有效章节查询 |

### ReadingProgress 集合

| 索引名称 | 字段 | 用途 |
|---------|------|------|
| `user_id_1_book_id_1` | user_id, book_id | 用户阅读进度查询 |
| `user_id_1_last_read_at_-1` | user_id, last_read_at | 用户最近阅读记录 |
| `book_id_1` | book_id | 书籍阅读统计 |

## 故障排查

### 测试跳过

如果测试被跳过，检查：

1. MongoDB 是否正在运行
   ```bash
   # Windows
   Get-Process mongo*

   # Linux/macOS
   ps aux | grep mongod
   ```

2. 环境变量是否设置
   ```bash
   # Windows
   echo $env:MONGO_URI

   # Linux/macOS
   echo $MONGO_URI
   ```

### 索引不存在

如果测试报告索引不存在：

1. 确认迁移已执行
2. 检查迁移日志是否有错误
3. 手动运行迁移并查看输出
4. 使用 mongosh 验证索引实际状态

### 连接失败

如果无法连接到 MongoDB：

1. 检查 MongoDB 是否启动
2. 验证连接字符串格式
3. 检查防火墙设置
4. 确认端口号（默认 27017）

## 清理测试数据

```javascript
// 删除测试数据库
use qingyu_test
db.dropDatabase()

// 或者删除特定集合的索引
db.users.dropIndexes()
db.books.dropIndexes()
db.chapters.dropIndexes()
db.reading_progress.dropIndexes()
```

## 相关文档

- [索引设计文档](../../docs/plans/2026-01-26-block3-database-optimization-design.md)
- [迁移执行器](../cmd/migrate/main.go)
- [索引迁移文件](../migration/mongodb/)
