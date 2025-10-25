# 青羽阅读-写作一体系统 (Qingyu Backend)

[![CI](https://github.com/yukin371/Qingyu_backend/workflows/Simple%20CI/badge.svg)](https://github.com/yukin371/Qingyu_backend/actions/workflows/ci.yml)
[![CodeQL](https://github.com/yukin371/Qingyu_backend/workflows/CodeQL%20Analysis/badge.svg)](https://github.com/yukin371/Qingyu_backend/actions/workflows/codeql.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## ⚠️ 重要说明

**关于 `.env` 文件**：本项目使用 **Viper + YAML** 进行配置管理，**不依赖 `.env` 文件**。所有配置通过 YAML 配置文件和环境变量管理。

📖 **配置文档**：
- [Viper配置管理机制详解](doc/Viper配置管理机制详解.md)
- [配置文件加载顺序和优先级说明](doc/配置文件加载顺序和优先级说明.md)
- [配置管理工具使用指南](doc/usage/配置管理工具使用指南.md)

## 项目简介

青语智能写作系统是一个基于Go语言开发的智能写作平台后端服务，集成了先进的AI技术，为用户提供智能写作辅助、内容生成、文本分析等功能。

## 主要功能

### 🤖 AI智能写作
- **智能内容生成**: 基于上下文和提示词生成高质量文本内容
- **智能续写**: 根据现有内容智能续写，保持风格连贯
- **文本分析**: 提供情节、角色、风格等多维度文本分析
- **内容优化**: 自动优化语法、风格和流畅度
- **大纲生成**: 基于主题自动生成详细故事大纲

### 📝 项目管理
- 项目创建和管理
- 章节组织和编辑
- 版本控制和历史记录
- 协作功能支持

### 👥 用户系统
- 用户注册和认证
- 权限管理
- 个人设置和偏好

### 🎯 角色与世界观
- 角色信息管理
- 世界观设定
- 情节线索追踪

## 技术栈

- **后端框架**: Go + Gin
- **数据库**: PostgreSQL / MySQL
- **认证**: JWT
- **AI服务**: OpenAI API / 其他AI服务商
- **配置管理**: Viper
- **日志**: Logrus
- **API文档**: Swagger

## 快速开始

### 环境要求

- Go 1.19+
- PostgreSQL 12+ 或 MySQL 8.0+
- Redis (可选，用于缓存)

### 安装步骤

1. **克隆项目**
```bash
git clone https://github.com/your-org/Qingyu_backend.git
cd Qingyu_backend
```

2. **安装依赖**
```bash
go mod download
```

3. **配置系统**
```bash
# 项目使用 YAML 配置文件，不需要 .env 文件
# 开发环境：使用 config/config.yaml（默认）
# 测试环境：使用 config/config.test.yaml（自动检测）
# 生产环境：通过 CONFIG_FILE 环境变量指定

# 查看配置文件
cat config/config.yaml

# 如需自定义，可复制配置模板
cp config/config.yaml config/config.local.yaml
# 编辑本地配置（不会提交到Git）
nano config/config.local.yaml
```

4. **配置数据库**
```bash
# 创建数据库
createdb qingyu_db

# 运行数据库迁移
go run cmd/migrate/main.go
```

5. **启动服务**
```bash
# 开发模式
go run main.go

# 或者构建后运行
go build -o qingyu_backend
./qingyu_backend
```

服务将在 `http://localhost:8080` 启动

## 配置说明

### YAML配置文件

本项目使用 **Viper + YAML** 配置，配置文件优先级：

```
1. config.test.yaml    (测试配置，自动检测优先)
2. config.yaml         (默认配置)
3. config.local.yaml   (本地配置，不提交)
4. config.docker.yaml  (Docker环境)
5. config.prod.yaml    (生产环境，不提交)
```

**配置文件示例** (`config/config.yaml`):

```yaml
# 服务器配置
server:
  port: "8080"
  mode: "debug"  # debug/release

# 数据库配置
mongodb:
  uri: "mongodb://localhost:27017"
  database: "Qingyu_backend"

# JWT配置
jwt:
  secret: "qingyu_secret_key"  # 生产环境必须修改
  expiration_hours: 24

# AI配置
ai:
  api_key: ""  # 你的AI API密钥
  base_url: "https://generativelanguage.googleapis.com/v1beta"
  model: "gemini-1.5-flash"
  max_tokens: 2000
  temperature: 0.7
```

### 环境变量覆盖

可以通过环境变量覆盖配置文件（优先级最高）：

```bash
# 环境变量命名规则：QINGYU_ + 配置路径（.替换为_）
export QINGYU_SERVER_PORT="9090"
export QINGYU_JWT_SECRET="production-secret-key"
export QINGYU_DATABASE_PRIMARY_MONGODB_URI="mongodb://prod-host:27017"
export QINGYU_AI_API_KEY="your-api-key"
```

### AI服务配置

系统支持多种AI服务提供商：

- **OpenAI**: 设置 `AI_PROVIDER=openai`
- **Azure OpenAI**: 设置 `AI_PROVIDER=azure`
- **其他兼容服务**: 设置相应的提供商标识

## API文档

### 主要接口

#### AI服务接口
- `POST /api/v1/ai/generate` - 生成内容
- `POST /api/v1/ai/continue` - 续写内容
- `POST /api/v1/ai/analyze` - 分析文本
- `POST /api/v1/ai/optimize` - 优化文本
- `POST /api/v1/ai/outline` - 生成大纲
- `GET /api/v1/ai/context/:projectId/:chapterId` - 获取上下文

#### 项目管理接口
- `GET /api/v1/projects` - 获取项目列表
- `POST /api/v1/projects` - 创建项目
- `GET /api/v1/projects/:id` - 获取项目详情
- `PUT /api/v1/projects/:id` - 更新项目
- `DELETE /api/v1/projects/:id` - 删除项目

#### 用户管理接口
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新Token
- `GET /api/v1/users/profile` - 获取用户信息

详细的API文档请参考：
- [AI API文档](doc/api/AI_API_Documentation.md)
- [AI服务使用指南](doc/usage/AI_Service_Usage_Guide.md)

## 项目结构

```
Qingyu_backend/
├── cmd/                    # 命令行工具
├── config/                 # 配置管理
│   ├── config.go          # 主配置
│   ├── ai.go              # AI服务配置
│   ├── database.go        # 数据库配置
│   └── jwt.go             # JWT配置
├── controllers/            # 控制器
├── middleware/             # 中间件
├── models/                 # 数据模型
│   ├── ai/                # AI相关模型
│   ├── project/           # 项目模型
│   └── user/              # 用户模型
├── router/                 # 路由配置
│   ├── ai/                # AI路由
│   └── api/               # API路由
├── service/                # 业务逻辑
│   ├── ai/                # AI服务
│   ├── project/           # 项目服务
│   └── user/              # 用户服务
├── utils/                  # 工具函数
├── doc/                    # 文档
│   ├── api/               # API文档
│   └── usage/             # 使用指南
├── config/                # 配置文件目录
│   ├── config.yaml       # 默认配置
│   ├── config.test.yaml  # 测试配置
│   └── config.docker.yaml # Docker配置
├── go.mod                 # Go模块文件
├── go.sum                 # 依赖校验文件
├── main.go                # 程序入口
└── README.md              # 项目说明
```

## 开发指南

### 代码规范

- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 添加必要的注释和文档
- 编写单元测试

### 开发工具安装

```bash
# 安装所有开发工具
make install-tools

# 安装 golangci-lint
make install-lint

# 初始化开发环境
make init
```

### 提交规范

本项目遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 功能开发
git commit -m "feat: 添加AI内容生成功能"

# 问题修复
git commit -m "fix: 修复用户认证问题"

# 文档更新
git commit -m "docs: 更新API文档"

# 性能优化
git commit -m "perf: 优化数据库查询性能"

# 重构
git commit -m "refactor: 重构用户服务层"

# 测试
git commit -m "test: 添加AI服务单元测试"

# 构建/CI
git commit -m "ci: 更新GitHub Actions配置"
```

### 本地开发流程

```bash
# 1. 启动开发服务器
make run

# 2. 启动热重载模式
make dev

# 3. 代码检查
make check

# 4. 运行测试
make test

# 5. 查看测试覆盖率
make test-coverage

# 6. 提交前检查
make pr-check
```

### 测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 运行API测试
make test-api

# 生成测试覆盖率报告
make test-coverage

# 快速测试（跳过慢速测试）
make test-quick

# 检查覆盖率是否达标（>=60%）
make test-coverage-check
```

### 代码质量检查

```bash
# 运行 linter
make lint

# 代码格式化
make fmt

# 安全扫描
make security

# 依赖漏洞检查
make vuln-check

# 代码复杂度检查
make complexity

# 完整的 CI 检查（本地模拟）
make ci-local
```

### CI/CD

项目使用简化的 GitHub Actions 进行持续集成：

- **代码检查**: golangci-lint 代码质量检查
- **安全扫描**: gosec 和 govulncheck 安全检查
- **单元测试**: 快速单元测试（不需要外部依赖）
- **集成测试**: 完整的集成测试（使用 MongoDB）
- **API 测试**: API 端点测试
- **CodeQL**: 自动化代码安全分析

详细信息请参考：
- [CI 修复总结](CI_FIX_SUMMARY.md)
- [快速参考 - CI/CD 命令](doc/ops/快速参考-CI_CD命令.md)
- [GitHub Actions 工作流说明](.github/workflows/README.md)

## 部署

### Docker部署

```bash
# 构建镜像
docker build -t qingyu-backend .

# 运行容器
docker run -d \
  --name qingyu-backend \
  -p 8080:8080 \
  --env-file .env \
  qingyu-backend
```

### 生产环境部署

1. **构建生产版本**
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qingyu_backend .
```

2. **配置反向代理** (Nginx示例)
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

3. **配置系统服务** (systemd示例)
```ini
[Unit]
Description=Qingyu Backend Service
After=network.target

[Service]
Type=simple
User=qingyu
WorkingDirectory=/opt/qingyu
ExecStart=/opt/qingyu/qingyu_backend
Restart=always

[Install]
WantedBy=multi-user.target
```

## 监控和日志

### 日志配置

系统使用结构化日志，支持多种输出格式：

```go
// 配置日志级别
log.SetLevel(log.InfoLevel)

// 设置日志格式
log.SetFormatter(&log.JSONFormatter{})
```

### 性能监控

- API响应时间监控
- 数据库查询性能
- AI服务调用统计
- 错误率和成功率

## 常见问题

### Q: AI服务调用失败怎么办？
A: 
1. 检查API密钥是否正确
2. 确认网络连接正常
3. 查看是否触发限流
4. 检查AI服务商状态

### Q: 数据库连接失败？
A:
1. 确认数据库服务正在运行
2. 检查连接参数是否正确
3. 验证用户权限
4. 查看防火墙设置

### Q: 如何优化性能？
A:
1. 启用Redis缓存
2. 优化数据库查询
3. 使用连接池
4. 实现请求限流
