# 青羽写作平台 - 论文答辩演示版本

## 简介

这是青羽写作平台的论文答辩演示版本，连接真实MongoDB数据库，支持完整的CRUD操作。

### 特点

- **完整CRUD操作**: 支持所有增删改查功能
- **真实数据库**: 连接MongoDB，数据持久化
- **内置演示数据**: 可通过seeder快速注入演示数据
- **完整API**: 保留所有核心API接口
- **Swagger文档**: 自动生成的API文档

## 环境要求

- Go 1.24+
- MongoDB 7.0+（必须）
- Redis 7.0+（可选，用于缓存）

## 快速开始

### 1. 启动MongoDB

```bash
# Windows（如果已安装MongoDB服务）
net start MongoDB

# 或使用Docker
docker run -d -p 27017:27017 --name qingyu-mongo mongo:7.0

# 或使用已有的MongoDB实例
```

### 2. 初始化演示数据（可选但推荐）

```bash
cd Qingyu_backend
go run cmd/seeder/main.go
```

### 3. 启动演示服务器

```bash
# 方式一：直接运行
go run cmd/demo_server/main.go

# 方式二：使用配置文件
go run cmd/demo_server/main.go

# 方式三：指定MongoDB连接
set MONGO_URI=mongodb://localhost:27017
set DB_NAME=qingyu
go run cmd/demo_server/main.go
```

### 4. 访问服务

- API地址: http://localhost:9090
- 健康检查: http://localhost:9090/health
- Swagger文档: http://localhost:9090/swagger/index.html
- 演示信息: http://localhost:9090/demo/info

## 演示数据

### 演示账户

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| demo | demo123 | 读者 | 普通读者账户 |
| author | author123 | 作者 | 作者账户，可发布作品 |
| admin | admin123 | 管理员 | 系统管理员 |

### 演示书籍

运行seeder后会自动创建：
- 《星际迷途》- 科幻类
- 《剑破苍穹》- 玄幻类
- 《都市仙尊》- 都市类
- 《青玉案》- 古言类

## 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| DEMO_PORT | 9090 | 服务端口 |
| MONGO_URI | mongodb://localhost:27017 | MongoDB连接字符串 |
| DB_NAME | qingyu | 数据库名称 |
| SEED_DATA | true | 是否自动检查演示数据 |
| CONFIG_FILE | ./config | 配置文件路径 |

## 演示端点

### 演示信息
```
GET /demo/info          # 获取演示版本信息
POST /demo/reset        # 重置演示数据
POST /demo/create-user  # 快速创建用户
```

### 核心API

所有完整的API接口都可以通过 http://localhost:9090/swagger/index.html 查看。

主要包括：
- 用户认证（登录、注册、登出）
- 书城管理（书籍CRUD、分类、排行榜）
- 阅读器（章节阅读、进度同步）
- 社交互动（评论、点赞、收藏）
- AI辅助写作（需要Python AI服务）
- 后台管理（用户管理、数据统计）

## 演示流程建议

### 1. 用户注册与登录
```bash
# 注册新用户
curl -X POST http://localhost:9090/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123","nickname":"测试用户"}'

# 登录
curl -X POST http://localhost:9090/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo123"}'
```

### 2. 浏览书城
```bash
# 获取书籍列表
curl http://localhost:9090/api/v1/bookstore/books

# 搜索书籍
curl "http://localhost:9090/api/v1/bookstore/search?keyword=星际"

# 获取分类
curl http://localhost:9090/api/v1/bookstore/categories
```

### 3. 阅读功能
```bash
# 获取书籍章节
curl http://localhost:9090/api/v1/bookstore/books/{book_id}/chapters

# 阅读章节
curl http://localhost:9090/api/v1/bookstore/chapters/{chapter_id}
```

## 注意事项

1. **需要MongoDB**: 演示版本需要MongoDB服务运行中
2. **数据持久化**: 数据存储在MongoDB中，重启不会丢失
3. **AI功能**: 默认跳过AI服务，如需完整功能请启动Python AI服务

## 目录结构

```
cmd/demo_server/
├── main.go           # 主入口（简化配置，连接MongoDB）
├── memory_store.go   # 内存存储（备用，当前版本未使用）
├── handlers.go       # API处理函数（备用，当前版本未使用）
└── README.md         # 本文档

config/
└── config.demo.yaml  # 演示配置文件
```

## 故障排除

### MongoDB连接失败
```bash
# 检查MongoDB是否运行
# Windows
sc query MongoDB

# Linux/Mac
systemctl status mongod

# 或使用Docker
docker ps | grep mongo
```

### 端口被占用
```bash
# 修改端口
set DEMO_PORT=8080
go run cmd/demo_server/main.go
```

### 数据库为空
```bash
# 运行数据填充
go run cmd/seeder/main.go
```

---

青羽写作平台 - 论文答辩演示版本
