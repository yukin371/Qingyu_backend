# 后端 Docker 配置

后端项目的 Docker 配置文件，包含 MongoDB、Redis 和后端服务。

## 📁 目录结构

```
Qingyu_backend/docker/
├── Dockerfile.dev           # 开发环境Dockerfile
├── Dockerfile.prod          # 生产环境Dockerfile
├── .dockerignore            # Docker忽略文件
├── .air.toml                # Air热重载配置
├── docker-compose.dev.yml   # 开发环境编排（含数据库）
├── docker-compose.prod.yml  # 生产环境编排
├── dev.bat                  # 开发环境启动脚本
├── stop.bat                 # 停止服务脚本
└── README.md                # 本文件
```

## 🚀 快速开始

### 开发环境

#### 使用脚本（推荐）
```bash
# 在 Qingyu_backend 目录下
cd docker
dev.bat
```

#### 使用 docker-compose
```bash
cd Qingyu_backend/docker
docker-compose -f docker-compose.dev.yml up -d
```

这将启动：
- MongoDB（数据库）
- Redis（缓存）
- Backend（Go服务，支持热重载）

### 生产环境

```bash
cd Qingyu_backend/docker
docker-compose -f docker-compose.prod.yml up -d --build
```

## 📋 服务说明

### MongoDB
- **端口**: 27017
- **数据库**: Qingyu_writer
- **数据持久化**: Docker Volume

### Redis
- **端口**: 6379
- **数据持久化**: Docker Volume

### Backend
- **端口**: 8080
- **热重载**: Air工具
- **框架**: Gin

## 🔧 配置说明

### 开发环境特性
- ✅ Air热重载（代码修改自动重启）
- ✅ 源代码实时挂载
- ✅ MongoDB + Redis
- ✅ 健康检查

### 生产环境特性
- ✅ 多阶段构建优化
- ✅ 二进制文件优化（-ldflags）
- ✅ 密码保护（MongoDB、Redis）
- ✅ 自动重启策略

### 环境变量

生产环境需要设置：
- `MONGO_PASSWORD` - MongoDB密码
- `REDIS_PASSWORD` - Redis密码

创建 `.env` 文件：
```env
MONGO_PASSWORD=your_secure_password
REDIS_PASSWORD=your_secure_password
```

## 📝 常用命令

### 启动服务
```bash
# 开发环境（含数据库）
docker-compose -f docker-compose.dev.yml up -d

# 生产环境
docker-compose -f docker-compose.prod.yml up -d
```

### 停止服务
```bash
docker-compose -f docker-compose.dev.yml down
```

### 查看日志
```bash
# 所有服务
docker-compose -f docker-compose.dev.yml logs -f

# 特定服务
docker-compose -f docker-compose.dev.yml logs -f backend
docker-compose -f docker-compose.dev.yml logs -f mongodb
```

### 进入容器
```bash
# 后端容器
docker-compose -f docker-compose.dev.yml exec backend sh

# MongoDB
docker-compose -f docker-compose.dev.yml exec mongodb mongosh

# Redis
docker-compose -f docker-compose.dev.yml exec redis redis-cli
```

### 重建服务
```bash
docker-compose -f docker-compose.dev.yml up -d --build
```

## 🌐 访问地址

- **后端API**: http://localhost:8080
- **MongoDB**: localhost:27017
- **Redis**: localhost:6379

## 🔗 网络配置

后端服务会创建并使用 `qingyu-network` 网络，前端服务可以通过加入此网络与后端通信。

## 🔍 故障排除

### 端口冲突
修改 `docker-compose.dev.yml` 中的端口映射：
```yaml
ports:
  - "8081:8080"  # 改为其他端口
```

### 数据库连接失败
1. 检查健康检查状态
2. 等待数据库完全启动（约30秒）
3. 查看日志排查问题

### 热重载不工作
1. 检查 `.air.toml` 配置
2. 查看容器日志
3. 重启容器

## 📚 相关文档

- [主项目文档](../../README.md)
- [前端Docker配置](../../Qingyu/docker/README.md)
- [Docker使用指南](../../README.Docker.md)
