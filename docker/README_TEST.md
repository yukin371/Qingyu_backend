# Docker 测试环境使用指南

## 概述

`docker-compose.test.yml` 提供了一个专门用于测试的隔离环境，包含：
- MongoDB 6.0（使用内存存储，测试结束自动清理）
- Redis 7（使用内存存储）

## 快速开始

### 使用自动化脚本（推荐）

**Linux/Mac:**
```bash
chmod +x scripts/run_tests_with_docker.sh
./scripts/run_tests_with_docker.sh
```

**Windows:**
```cmd
scripts\run_tests_with_docker.bat
```

### 手动运行

**1. 启动测试环境**
```bash
docker-compose -f docker/docker-compose.test.yml up -d
```

**2. 等待服务就绪**
```bash
# 检查服务状态
docker-compose -f docker/docker-compose.test.yml ps

# 查看日志
docker-compose -f docker/docker-compose.test.yml logs
```

**3. 设置环境变量并运行测试**
```bash
export MONGODB_URI="mongodb://admin:password@localhost:27017"
export MONGODB_DATABASE="qingyu_test"
export REDIS_ADDR="localhost:6379"
export ENVIRONMENT="test"

# 运行单元测试
go test -v -short ./...

# 运行集成测试
go test -v ./test/integration/...

# 运行API测试
go test -v ./test/api/...
```

**4. 清理测试环境**
```bash
docker-compose -f docker/docker-compose.test.yml down -v
```

## 特性

### 内存存储
- 使用 `tmpfs` 将数据库存储在内存中
- 测试速度更快
- 容器关闭时自动清理，不留垃圾数据

### 隔离环境
- 使用独立的网络 `test-network`
- 不影响开发环境的数据库
- 可以并行运行多个测试环境（需修改端口）

### 健康检查
- MongoDB: 每5秒检查一次，最多重试10次
- Redis: 每3秒检查一次，最多重试5次

## 连接信息

### MongoDB
- 地址: `localhost:27017`
- 用户名: `admin`
- 密码: `password`
- 数据库: `qingyu_test`
- 完整URI: `mongodb://admin:password@localhost:27017`

### Redis
- 地址: `localhost:6379`
- 无密码
- 数据库: 0（默认）

## 故障排查

### 查看容器状态
```bash
docker-compose -f docker/docker-compose.test.yml ps
```

### 查看日志
```bash
# 所有服务日志
docker-compose -f docker/docker-compose.test.yml logs

# MongoDB日志
docker-compose -f docker/docker-compose.test.yml logs mongodb-test

# Redis日志
docker-compose -f docker/docker-compose.test.yml logs redis-test
```

### 进入容器调试
```bash
# 进入MongoDB容器
docker exec -it qingyu-mongodb-test bash

# 连接MongoDB
docker exec -it qingyu-mongodb-test mongo -u admin -p password

# 进入Redis容器
docker exec -it qingyu-redis-test sh

# 连接Redis
docker exec -it qingyu-redis-test redis-cli
```

### 完全重置
```bash
# 停止并删除所有容器、网络和卷
docker-compose -f docker/docker-compose.test.yml down -v

# 重新启动
docker-compose -f docker/docker-compose.test.yml up -d
```

## CI/CD 集成

在 `.github/workflows/ci.yml` 中已经配置好使用此测试环境：
- 集成测试使用 Docker Compose
- API测试使用 Docker Compose
- 自动启动和清理

## 性能优化建议

1. **首次运行**: 需要下载Docker镜像，会比较慢
2. **后续运行**: 使用缓存的镜像，启动很快（通常10-15秒）
3. **并行测试**: 可以使用 `go test -parallel` 参数提高测试速度
4. **内存要求**: 建议至少4GB可用内存

## 与开发环境的区别

| 特性 | 测试环境 | 开发环境 |
|------|---------|---------|
| 数据存储 | 内存（tmpfs） | 持久化卷 |
| 端口 | 27017, 6379 | 27017, 6379 |
| 网络 | test-network | qingyu-network |
| 容器名 | qingyu-*-test | qingyu-* |
| 数据保留 | 临时 | 永久 |

## 最佳实践

1. **测试前清理**: 总是使用 `down -v` 清理旧数据
2. **环境隔离**: 不要在测试环境中使用生产数据
3. **并行运行**: 可以同时运行测试和开发环境
4. **资源监控**: 注意内存使用，避免OOM

## 相关文件

- `docker/docker-compose.test.yml` - 测试环境配置
- `scripts/run_tests_with_docker.sh` - Linux/Mac测试脚本
- `scripts/run_tests_with_docker.bat` - Windows测试脚本
- `.github/workflows/ci.yml` - CI配置
- `config/config.test.yaml` - 测试环境配置文件

