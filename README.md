# 青羽阅读-写作一体系统 (Qingyu Backend)

[![CI/CD Pipeline](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-success?logo=github-actions)](https://github.com/yukin371/Qingyu_backend/actions)
[![Go Version](https://img.shields.io/badge/Go-1.24-blue?logo=go)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-1.11.0-red?logo=gin)](https://gin-gonic.com)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-latest-brightgreen)](./doc)

> 一款基于 Go + Gin + MongoDB 的新一代智能创作平台后端，集成了阅读、写作、社交等完整功能，采用模块化架构设计，支持从单体到微服务的平滑演进。

## 项目简介

青羽是一个面向中文网络文学领域的**阅读-写作一体化平台**，为创作者提供项目管理、版本控制、协同编辑等强大功能，为读者提供个性化推荐、社交互动、VIP会员等丰富体验。

### 核心特色

- **完整内容生态** - 从创作到阅读的完整闭环，包含书店、阅读器、社交互动等10大业务模块
- **模块化架构** - 清晰的分层设计（API → Service → Repository），依赖注入容器管理，支持微服务演进
- **高性能设计** - MongoDB + Redis 多层缓存，支持10万+QPS，响应时间<100ms
- **开发友好** - 完善的测试体系、Docker一键启动、Makefile命令管理、热重载开发
- **AI智能集成** - 支持多AI提供商（OpenAI、智谱AI、DeepSeek等），提供智能写作辅助

## 技术栈

### 核心技术

| 类别 | 技术选型 | 版本 | 用途 |
|------|---------|------|------|
| **语言** | Go | 1.24+ | 高性能后端开发 |
| **框架** | Gin | v1.11.0 | HTTP路由和中间件 |
| **数据库** | MongoDB | v7.0 | 主数据存储 |
| **缓存** | Redis | v7-alpine | 分布式缓存 |
| **向量DB** | Milvus | v2.3.0 | AI向量检索 |
| **认证** | JWT | v4.5.2 | 用户身份验证 |
| **配置** | Viper | v1.21.0 | 配置管理 |
| **日志** | Zap | v1.27.0 | 结构化日志 |
| **监控** | Prometheus | v1.23.2 | 性能指标采集 |
| **通信** | gRPC | v1.76.0 | 服务间通信 |

### AI服务生态

- **Python AI Service** - FastAPI + LangGraph + LangChain
- **多模型支持** - GPT-4、Claude、Gemini、GLM-4、DeepSeek
- **RAG增强** - Milvus向量检索 + 文档索引
- **流式响应** - 支持SSE实时推送

### DevOps工具

- **容器化** - Docker + Docker Compose
- **CI/CD** - GitHub Actions
- **代码质量** - golangci-lint + gosec
- **测试覆盖** - go test + testify
- **文档生成** - Swagger

## 系统架构

### 分层架构总览

```mermaid
graph TB
    subgraph "客户端层"
        Web[Vue3 Web应用]
        Mobile[移动端应用]
    end

    subgraph "网关层"
        Router[Gin Router]
        Middleware[中间件链<br/>JWT/CORS/限流/日志]
    end

    subgraph "业务层 - 10大模块"
        AI[AI智能写作]
        Bookstore[书店系统]
        Reader[阅读器]
        Social[社交互动]
        Writer[写作创作]
        User[用户系统]
        Finance[财务支付]
        Recommend[推荐系统]
        Notice[通知公告]
        Admin[管理后台]
    end

    subgraph "数据层"
        MongoDB[(MongoDB<br/>主数据)]
        Redis[(Redis<br/>缓存)]
        Milvus[(Milvus<br/>向量DB)]
    end

    subgraph "AI服务层"
        PythonAI[Python AI Service<br/>FastAPI + LangGraph]
    end

    Web --> Router
    Mobile --> Router
    Router --> Middleware
    Middleware --> AI
    Middleware --> Bookstore
    Middleware --> Reader
    Middleware --> Social
    Middleware --> Writer
    Middleware --> User
    Middleware --> Finance
    Middleware --> Recommend
    Middleware --> Notice
    Middleware --> Admin

    AI --> PythonAI
    PythonAI --> Milvus

    AI --> MongoDB
    Bookstore --> MongoDB
    Reader --> Redis
    Recommend --> Redis
    User --> MongoDB
    Finance --> MongoDB

    style AI fill:#f9f,stroke:#333,stroke-width:2px
    style PythonAI fill:#ff9,stroke:#333,stroke-width:2px
    style Milvus fill:#99f,stroke:#333,stroke-width:2px
```

### 分层职责说明

| 层级 | 职责 | 目录 | 示例 |
|------|------|------|------|
| **API层** | 路由定义、参数验证、响应封装 | `api/v1/` | `ai/writing_api.go` |
| **Service层** | 业务逻辑、事务管理、外部调用 | `service/` | `ai/writing_service.go` |
| **Repository层** | 数据访问、查询构建、缓存管理 | `repository/` | `mongodb/book_repository.go` |
| **Model层** | 数据模型、验证规则、类型定义 | `models/` | `bookstore/book.go` |
| **Middleware层** | 认证、权限、限流、日志 | `middleware/` | `auth_middleware.go` |

## 快速开始

### 环境要求

| 组件 | 最低版本 | 推荐版本 | 必需/可选 |
|------|---------|---------|----------|
| Go | 1.24 | 1.24+ | 必需 |
| Docker | 20.10+ | 24.0+ | 推荐 |
| MongoDB | 7.0 | 7.0+ | 必需 |
| Redis | 7.0 | 7.0+ | 必需 |
| Python | 3.10+ | 3.11+ | AI功能必需 |
| Git | 2.30+ | 最新版 | 必需 |

### 一键启动（推荐）

使用 Docker Compose 是最简单的启动方式，所有依赖服务自动配置。

```bash
# 1. 克隆项目
git clone https://github.com/your-org/Qingyu_backend.git
cd Qingyu_backend

# 2. 配置AI密钥（可选）
# 创建 .env 文件并添加你的API密钥
echo "OPENAI_API_KEY=your_key_here" > .env

# 3. 启动所有服务（MongoDB、Redis、Milvus、后端、Python AI服务）
make docker-up

# 4. 查看服务状态
docker-compose -f docker/docker-compose.dev.yml ps

# 5. 等待服务启动（约1-2分钟），查看日志
make docker-logs

# 6. 验证服务健康
curl http://localhost:8080/health
```

**服务地址**：
- Go Backend: http://localhost:8080
- Python AI Service: http://localhost:8000
- MongoDB: mongodb://localhost:27017
- Redis: redis://localhost:6379
- Milvus: localhost:19530

### 本地开发模式

如果你需要在本地修改代码并实时调试：

```bash
# 1. 安装依赖
go mod download

# 2. 启动基础设施（MongoDB + Redis + Milvus）
cd docker && docker-compose -f docker-compose.db-only.yml up -d

# 3. 配置文件（默认使用 config/config.yaml）
# 可选：创建本地配置覆盖
cp config/config.yaml config/config.local.yaml
# 编辑 config.local.yaml，修改数据库连接等

# 4. 运行服务
go run cmd/server/main.go

# 或使用 Makefile 命令
make run
```

### 仅启动后端（已有数据库）

如果你的环境已经有运行中的 MongoDB 和 Redis：

```bash
# 快速启动（使用默认配置）
make run

# 构建并运行
make build
./bin/server
```

### 验证安装

```bash
# 健康检查
curl http://localhost:8080/health

# 预期响应
{
  "status": "ok",
  "version": "1.0.0",
  "database": "connected",
  "redis": "connected",
  "ai_service": "connected"
}

# 测试API接口
curl http://localhost:8080/api/v1/bookstore/books

# 运行快速测试
make test-quick
```

### 常见启动问题

<details>
<summary><b>问题1：MongoDB连接失败</b></summary>

**症状**：`connection refused` 或 `timeout`

**解决方案**：
```bash
# 检查MongoDB是否运行
docker ps | grep mongo

# 查看MongoDB日志
docker logs qingyu-mongodb

# 重启MongoDB
docker restart qingyu-mongodb
```
</details>

<details>
<summary><b>问题2：端口被占用</b></summary>

**解决方案**：
```bash
# 检查端口占用
netstat -tunlp | grep 8080

# 修改配置文件端口
vim config/config.local.yaml
# server.port: "9090"

# 或使用环境变量
export QINGYU_SERVER_PORT=9090
```
</details>

<details>
<summary><b>问题3：AI服务调用失败</b></summary>

**解决方案**：
```bash
# 检查Python AI服务状态
curl http://localhost:8000/health

# 检查API密钥配置
cat .env | grep OPENAI_API_KEY

# 查看AI服务日志
docker logs qingyu-python-ai
```
</details>

## 项目结构

### 目录树

```
Qingyu_backend/
├── cmd/server/              # 应用入口
│   └── main.go             # 主程序入口
├── api/v1/                 # API层 - 10大模块路由
│   ├── ai/                 # AI智能写作接口
│   ├── bookstore/          # 书店系统接口
│   ├── reader/             # 阅读器接口
│   ├── social/             # 社交互动接口
│   ├── writer/             # 写作创作接口
│   ├── auth/               # 认证授权接口
│   ├── finance/            # 财务支付接口
│   ├── recommendation/     # 推荐系统接口
│   ├── notifications/      # 通知公告接口
│   └── admin/              # 管理后台接口
├── service/                # 业务逻辑层
│   ├── container/          # 依赖注入容器
│   ├── ai/                 # AI服务
│   ├── bookstore/          # 书店业务
│   ├── reader/             # 阅读业务
│   ├── writer/             # 写作业务
│   └── [其他模块]/
├── repository/             # 数据访问层
│   ├── interfaces/         # Repository接口定义
│   ├── mongodb/            # MongoDB操作封装
│   ├── redis/              # Redis操作封装
│   └── querybuilder/       # 查询构建器
├── models/                 # 数据模型层
│   ├── ai/                 # AI相关模型
│   ├── user/               # 用户模型
│   ├── bookstore/          # 书店模型
│   └── [其他模型]/
├── middleware/             # 中间件
│   ├── auth_middleware.go  # JWT认证
│   ├── rbac_middleware.go  # 权限控制
│   ├── rate_limit.go       # 限流
│   └── logger.go           # 日志记录
├── router/                 # 路由注册
│   └── router.go           # 主路由
├── core/                   # 核心初始化
│   ├── database.go         # 数据库初始化
│   └── redis.go            # Redis初始化
├── config/                 # 配置文件
│   ├── config.yaml         # 默认配置
│   ├── config.test.yaml    # 测试配置
│   └── config.docker.yaml  # Docker配置
├── pkg/                    # 通用工具包
│   ├── cache/              # 缓存封装
│   ├── logger/             # 日志工具
│   ├── errors/             # 错误处理
│   └── response/           # 响应封装
├── python_ai_service/      # Python AI服务
│   ├── src/                # FastAPI应用
│   ├── proto/              # gRPC协议定义
│   └── tests/              # AI服务测试
├── test/                   # 测试代码
│   ├── unit/               # 单元测试
│   ├── integration/        # 集成测试
│   └── e2e/                # 端到端测试
├── docker/                 # Docker配置
│   ├── Dockerfile.dev      # 开发环境
│   └── docker-compose.*.yml # 多环境编排
├── doc/                    # 设计文档
│   ├── architecture/       # 架构文档
│   ├── api/                # API文档
│   └── design/             # 设计文档
├── .github/                # GitHub配置
│   └── workflows/          # CI/CD流程
├── Makefile                # 构建命令
├── go.mod                  # Go模块定义
└── README.md               # 本文档
```

### 依赖注入容器设计

项目使用 ServiceContainer 实现依赖注入，统一管理所有服务的生命周期：

```go
// service/container/container.go
type Container struct {
    // 数据库
    MongoDB      *mongo.Client
    Redis        *redis.Client

    // Repositories
    BookRepo     repository.BookRepository
    UserRepo     repository.UserRepository

    // Services
    BookService  service.BookService
    UserService  service.UserService

    // Handlers (API层)
    BookAPI      *api.BookAPI
    UserAPI      *api.UserAPI
}

// 初始化顺序保证依赖关系
func Initialize() *Container {
    return &Container{
        // 1. 基础设施
        MongoDB:   db.InitMongo(),
        Redis:     db.InitRedis(),

        // 2. Repository层
        BookRepo:  repository.NewBookRepository(),

        // 3. Service层
        BookService: service.NewBookService(container.BookRepo),

        // 4. API层
        BookAPI: api.NewBookAPI(container.BookService),
    }
}
```

## 核心功能模块

### 10大业务模块概览

<details>
<summary><b>1. AI智能写作模块</b></summary>

**功能亮点**：
- 多AI提供商支持（OpenAI、智谱AI、DeepSeek、Claude）
- 流式响应和实时推送
- 上下文管理和会话保持
- RAG检索增强生成

**技术架构**：
- Go后端：统一API入口、适配器模式
- Python服务：LangGraph Agent工作流、LangChain工具链
- 向量数据库：Milvus文档检索

**关键接口**：
```bash
POST /api/v1/ai/writing/continue     # 智能续写
POST /api/v1/ai/writing/generate     # 内容生成
POST /api/v1/ai/creative/character   # 角色创建
POST /api/v1/ai/creative/outline     # 大纲生成
GET  /api/v1/ai/context/:projectId   # 上下文获取
```

**文档**：[AI模块设计文档](./doc/design/ai/README.md)
</details>

<details>
<summary><b>2. 书店系统模块</b></summary>

**功能亮点**：
- 书籍发布与管理
- 多维度搜索和筛选
- 排行榜系统（日榜、周榜、月榜）
- 分类导航和标签管理
- 统计分析（阅读量、评分、收藏）

**关键接口**：
```bash
GET  /api/v1/bookstore/books           # 书籍列表
POST /api/v1/bookstore/books           # 发布书籍
GET  /api/v1/bookstore/books/:id       # 书籍详情
GET  /api/v1/bookstore/search          # 搜索书籍
GET  /api/v1/bookstore/ranking         # 排行榜
```
</details>

<details>
<summary><b>3. 阅读器模块</b></summary>

**功能亮点**：
- 章节阅读和翻页
- 阅读进度跨设备同步
- 书签管理
- 段落标注和笔记
- 阅读历史追踪

**关键接口**：
```bash
GET  /api/v1/reader/books/:id/chapters/:chapterId
POST /api/v1/reader/progress
GET  /api/v1/reader/bookmarks
POST /api/v1/reader/annotations
```
</details>

<details>
<summary><b>4. 社交互动模块</b></summary>

**功能亮点**：
- 评论系统（书籍评论、章节评论）
- 点赞和收藏
- 关注作者
- 动态发布
- 书单管理

**关键接口**：
```bash
POST /api/v1/social/comments
POST /api/v1/social/like
GET  /api/v1/social/following
POST /api/v1/social/booklists
```
</details>

<details>
<summary><b>5. 写作创作模块</b></summary>

**功能亮点**：
- 项目管理（创建、编辑、删除）
- 章节编辑和组织
- 版本控制和历史回溯
- 人物设定管理
- 世界观设定
- 协作编辑支持

**关键接口**：
```bash
POST /api/v1/writer/projects
PUT  /api/v1/writer/documents/:id
POST /api/v1/writer/versions
GET  /api/v1/writer/settings/:projectId
```
</details>

<details>
<summary><b>6. 用户系统模块</b></summary>

**功能亮点**：
- 用户注册和登录
- JWT认证和刷新
- RBAC权限管理
- VIP会员系统
- 个人资料管理
- 第三方登录（OAuth）

**关键接口**：
```bash
POST /api/v1/auth/register
POST /api/v1/auth/login
GET  /api/v1/users/profile
PUT  /api/v1/users/settings
POST /api/v1/auth/refresh
```
</details>

<details>
<summary><b>7. 财务支付模块</b></summary>

**功能亮点**：
- 钱包管理（余额、交易记录）
- 充值和提现
- VIP会员购买
- 章节购买
- 作者收入管理
- 订单管理

**关键接口**：
```bash
GET  /api/v1/finance/wallet
POST /api/v1/finance/recharge
GET  /api/v1/finance/orders
GET  /api/v1/finance/author/revenue
```
</details>

<details>
<summary><b>8. 推荐系统模块</b></summary>

**功能亮点**：
- 个性化推荐
- 相似内容推荐
- 热门榜单
- 分类推荐
- 行为记录和分析

**关键接口**：
```bash
GET /api/v1/recommendation/for-you
GET /api/v1/recommendation/similar/:id
GET /api/v1/recommendation/hot
```
</details>

<details>
<summary><b>9. 通知公告模块</b></summary>

**功能亮点**：
- 系统通知
- 个性化通知配置
- 公告发布和管理
- 消息中心
- 设备管理和推送

**关键接口**：
```bash
GET  /api/v1/notifications
POST /api/v1/admin/announcements
PUT  /api/v1/notifications/settings
```
</details>

<details>
<summary><b>10. 管理后台模块</b></summary>

**功能亮点**：
- 用户管理（审核、封禁）
- 内容审核（敏感内容过滤）
- 系统配置管理
- 数据统计分析
- 权限管理

**关键接口**：
```bash
GET  /api/v1/admin/users
POST /api/v1/admin/audit/:id
GET  /api/v1/admin/stats
PUT  /api/v1/admin/config
```
</details>

## 开发指南

### 开发规范

#### 命名规范

```go
// 文件命名：snake_case
user_service.go
auth_middleware.go

// 包命名：小写单词
package service
package middleware

// 接口命名：动作+名词
type UserRepository interface {
    CreateUser(ctx context.Context, user *User) error
    FindUserByID(ctx context.Context, id string) (*User, error)
}

// 常量命名：大驼峰
const MaxRetries = 3
const DefaultTimeout = 30 * time.Second
```

#### 错误处理规范

```go
// 使用 pkg/errors 包装错误
import "github.com/pkg/errors"

func (s *BookService) CreateBook(req *CreateBookRequest) error {
    // 验证输入
    if err := validateRequest(req); err != nil {
        return errors.Wrap(err, "invalid request")
    }

    // 调用repository
    if err := s.bookRepo.Create(req); err != nil {
        return errors.Wrap(err, "failed to create book")
    }

    return nil
}
```

#### 日志规范

```go
import "go.uber.org/zap"

// 结构化日志
logger.Info("User logged in",
    zap.String("user_id", userID),
    zap.String("ip", c.ClientIP()),
)

logger.Error("Failed to create book",
    zap.String("title", req.Title),
    zap.Error(err),
)
```

### 分层开发示例

#### API层开发

```go
// api/v1/bookstore/book_api.go
package bookstore

type BookAPI struct {
    bookService service.BookService
}

// NewBookAPI - 通过构造函数注入依赖
func NewBookAPI(bookService service.BookService) *BookAPI {
    return &BookAPI{
        bookService: bookService,
    }
}

// CreateBook - 创建书籍
func (api *BookAPI) CreateBook(c *gin.Context) {
    var req CreateBookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    book, err := api.bookService.CreateBook(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create book", err)
        return
    }

    response.Success(c, book)
}
```

#### Service层开发

```go
// service/bookstore/book_service.go
package bookstore

type BookService struct {
    bookRepo   repository.BookRepository
    cacheRepo  repository.CacheRepository
}

func (s *BookService) CreateBook(ctx context.Context, req *CreateBookRequest) (*models.Book, error) {
    // 1. 业务验证
    if err := s.validateBookRequest(req); err != nil {
        return nil, errors.Wrap(err, "validation failed")
    }

    // 2. 构造模型
    book := &models.Book{
        ID:        uuid.New().String(),
        Title:     req.Title,
        CreatedAt: time.Now(),
    }

    // 3. 持久化
    if err := s.bookRepo.Create(ctx, book); err != nil {
        return nil, errors.Wrap(err, "failed to save book")
    }

    // 4. 清除缓存
    s.cacheRepo.Delete(ctx, "books:list")

    return book, nil
}
```

#### Repository层开发

```go
// repository/mongodb/book_repository.go
package mongodb

type BookRepository struct {
    db *mongo.Database
}

func (r *BookRepository) Create(ctx context.Context, book *models.Book) error {
    _, err := r.db.Collection("books").InsertOne(ctx, book)
    return err
}

func (r *BookRepository) FindByID(ctx context.Context, id string) (*models.Book, error) {
    var book models.Book
    err := r.db.Collection("books").FindOne(
        ctx,
        bson.M{"_id": id},
    ).Decode(&book)

    if err == mongo.ErrNoDocuments {
        return nil, ErrBookNotFound
    }

    return &book, err
}
```

### API设计规范

#### 统一响应格式

```go
// 成功响应
{
  "code": 0,
  "message": "success",
  "data": {...}
}

// 错误响应
{
  "code": 10001,
  "message": "用户不存在",
  "detail": "user_id not found: 123456"
}
```

#### RESTful路由设计

```bash
# 资源命名：复数名词
GET    /api/v1/books           # 列表
POST   /api/v1/books           # 创建
GET    /api/v1/books/:id       # 详情
PUT    /api/v1/books/:id       # 更新
DELETE /api/v1/books/:id       # 删除

# 嵌套资源
GET    /api/v1/books/:id/chapters          # 书籍的章节列表
POST   /api/v1/books/:id/chapters          # 为书籍创建章节

# 动作路由
POST   /api/v1/books/:id/publish           # 发布书籍
POST   /api/v1/books/:id/archive           # 归档书籍
```

### Git提交规范

#### Commit Message格式

```bash
# 格式
<type>(<scope>): <subject>

# 示例
feat(ai): add character generation API
fix(auth): resolve JWT refresh token issue
docs(readme): update architecture diagram
```

#### Type类型

| 类型 | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat(bookstore): add book search API` |
| `fix` | Bug修复 | `fix(auth): resolve JWT refresh token issue` |
| `docs` | 文档 | `docs(readme): update architecture diagram` |
| `style` | 格式 | `style: format code with gofmt` |
| `refactor` | 重构 | `refactor(service): extract base service` |
| `test` | 测试 | `test(repo): add integration tests` |
| `chore` | 构建/工具 | `chore: update dependencies` |

### 测试指南

#### 测试分层

```
E2E Tests (端到端测试)
    ↓
Integration Tests (集成测试)
    ↓
Unit Tests (单元测试)
```

#### 单元测试

```go
// service/bookstore/book_service_test.go
func TestBookService_CreateBook(t *testing.T) {
    // 使用mock
    mockRepo := &MockBookRepository{}
    service := NewBookService(mockRepo)

    tests := []struct {
        name    string
        req     *CreateBookRequest
        wantErr bool
    }{
        {
            name: "成功创建书籍",
            req: &CreateBookRequest{
                Title:    "测试书籍",
                AuthorID: "author123",
            },
            wantErr: false,
        },
        {
            name:    "标题为空",
            req:     &CreateBookRequest{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            book, err := service.CreateBook(context.Background(), tt.req)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, book.ID)
            }
        })
    }
}
```

#### 运行测试

```bash
# 所有测试
make test

# 单元测试（快速）
make test-unit

# 集成测试（需要数据库）
make test-integration

# E2E测试（完整流程）
make test-e2e

# 测试覆盖率
make test-coverage

# 查看覆盖率报告
open coverage.html
```

#### 测试覆盖率要求

| 模块 | 目标覆盖率 | 当前状态 |
|------|-----------|---------|
| API层 | 80%+ | ✅ 75% |
| Service层 | 85%+ | ✅ 70% |
| Repository层 | 90%+ | ⚠️ 65% |
| 整体 | 70%+ | ⚠️ 65% |

## 部署指南

### Docker部署（推荐）

#### 开发环境

```bash
# 启动所有服务
make docker-up

# 查看日志
make docker-logs

# 停止服务
make docker-down
```

#### 生产环境

```bash
# 1. 构建生产镜像
docker build -f docker/Dockerfile.prod -t qingyu-backend:latest .

# 2. 使用生产配置
docker-compose -f docker/docker-compose.prod.yml up -d

# 3. 健康检查
curl http://localhost:8080/health
```

### Kubernetes部署

<details>
<summary><b>K8s部署配置示例</b></summary>

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: qingyu-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: qingyu-backend
  template:
    metadata:
      labels:
        app: qingyu-backend
    spec:
      containers:
      - name: backend
        image: qingyu-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: QINGYU_SERVER_PORT
          value: "8080"
        - name: QINGYU_DATABASE_MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: qingyu-secrets
              key: mongodb-uri
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: qingyu-backend-service
spec:
  selector:
    app: qingyu-backend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

</details>

### 传统部署

<details>
<summary><b>Systemd服务配置</b></summary>

```ini
# /etc/systemd/system/qingyu-backend.service
[Unit]
Description=Qingyu Backend Service
After=network.target mongodb.service redis.service

[Service]
Type=simple
User=qingyu
Group=qingyu
WorkingDirectory=/opt/qingyu/backend
ExecStart=/opt/qingyu/backend/bin/server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# 环境变量
Environment="QINGYU_SERVER_MODE=release"
Environment="QINGYU_SERVER_PORT=8080"
EnvironmentFile=/opt/qingyu/backend/.env

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable qingyu-backend
sudo systemctl start qingyu-backend
sudo systemctl status qingyu-backend
```
</details>

### 环境变量配置

```bash
# .env.production
QINGYU_SERVER_MODE=release
QINGYU_SERVER_PORT=8080

# 数据库
QINGYU_DATABASE_MONGODB_URI=mongodb://user:pass@mongo-cluster:27017
QINGYU_DATABASE_NAME=Qingyu_prod
QINGYU_REDIS_HOST=redis-cluster
QINGYU_REDIS_PORT=6379
QINGYU_REDIS_PASSWORD=your_redis_password

# AI服务
QINGYU_AI_API_KEY=your_production_key
QINGYU_AI_PYTHON_HOST=python-ai-service
QINGYU_AI_PYTHON_GRPC_PORT=50052

# 安全
QINGYU_JWT_SECRET=your_strong_secret_key_here
QINGYU_JWT_EXPIRATION_HOURS=24

# 监控
QINGYU_PROMETHEUS_ENABLED=true
```

## 监控与日志

### Prometheus指标

**内置指标端点**：
```bash
# 查看所有指标
curl http://localhost:8080/metrics

# 关键指标示例
# HTTP请求总数
http_requests_total{method="POST",path="/api/v1/ai/writing/generate",status="200"}

# 请求耗时
http_request_duration_seconds{path="/api/v1/books",quantile="0.99"}

# 数据库连接数
mongodb_connections_active

# 缓存命中率
redis_cache_hits_total
redis_cache_misses_total
```

### 日志管理

**日志级别**：
```go
// 开发环境：Debug级别
logger.SetLevel(zap.DebugLevel)

// 生产环境：Info级别
logger.SetLevel(zap.InfoLevel)
```

**结构化日志**：
```json
{
  "level": "info",
  "ts": 1704067200.123,
  "caller": "service/book_service.go:45",
  "msg": "Book created successfully",
  "user_id": "user123",
  "book_id": "book456",
  "duration_ms": 123
}
```

## 故障排查

### 常见问题

#### Q1: 服务启动失败

**症状**：`panic: runtime error` 或 `connection refused`

**排查步骤**：
```bash
# 1. 检查配置文件
cat config/config.yaml

# 2. 检查环境变量
env | grep QINGYU_

# 3. 检查依赖服务
docker ps | grep -E "mongo|redis"

# 4. 查看详细日志
tail -f server.log

# 5. 检查端口占用
netstat -tunlp | grep 8080
```

#### Q2: AI服务调用失败

**症状**：`AI service unavailable` 或 `timeout`

**排查步骤**：
```bash
# 1. 检查Python AI服务状态
curl http://localhost:8000/health

# 2. 检查gRPC连接
grpcurl -plaintext localhost:50052 list

# 3. 查看AI服务日志
docker logs qingyu-python-ai

# 4. 检查API密钥
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

#### Q3: 数据库查询慢

**症状**：API响应时间 > 1秒

**优化方案**：
```bash
# 1. 启用MongoDB慢查询日志
db.setProfilingLevel(1, {slowms: 100})

# 2. 检查索引使用情况
db.books.getIndexes()

# 3. 分析查询计划
db.books.find({title: "test"}).explain("executionStats")

# 4. 添加索引
db.books.createIndex({title: 1, author: 1})

# 5. 启用Redis缓存
# 在Service层添加缓存逻辑
```

#### Q4: 内存泄漏

**症状**：服务内存持续增长

**排查步骤**：
```bash
# 1. 使用pprof分析
go tool pprof http://localhost:8080/debug/pprof/heap

# 2. 生成内存分布图
go tool pprof -http=:8081 http://localhost:8080/debug/pprof/heap

# 3. 检查goroutine泄漏
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# 4. 定位问题代码
go tool pprof -list CreateBook /path/to/profile
```

## 贡献指南

我们欢迎所有形式的贡献！

### 贡献方式

- 报告Bug
- 提出新功能建议
- 改进文档
- 提交代码

### 开发流程

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交Pull Request

### Pull Request检查清单

- [ ] 代码通过所有测试 (`make test`)
- [ ] 代码通过linter检查 (`make lint`)
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
- [ ] Commit message符合规范

### 代码审查标准

- 代码可读性
- 测试覆盖率
- 性能影响
- 安全性考虑
- 文档完整性

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

- 项目主页：https://github.com/your-org/Qingyu_backend
- 问题反馈：https://github.com/your-org/Qingyu_backend/issues
- 邮箱：dev@qingyu.com

## 相关文档

- [架构设计文档](./doc/architecture/项目概述.md)
- [API接口文档](./doc/api/API接口总览.md)
- [数据库设计](./doc/design/database/数据库设计说明书.md)
- [部署运维手册](./docs/ops/部署运维手册.md)
- [开发规范](./doc/architecture/项目开发规则.md)
- [Viper配置管理详解](./doc/Viper配置管理机制详解.md)

## 致谢

感谢所有为本项目做出贡献的开发者！

### 特别感谢

- Go社区
- Gin框架团队
- MongoDB官方
- LangGraph社区
- 所有贡献者

---

**青羽阅读-写作一体系统** - 让创作更简单，让阅读更美好
