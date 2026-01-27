# 权限系统测试环境准备指南

> Qingyu Backend 权限系统测试环境配置文档
>
> 版本: v1.0 | 更新时间: 2026-01-27

## 目录

- [快速开始](#快速开始)
- [环境要求](#环境要求)
- [安装步骤](#安装步骤)
- [配置说明](#配置说明)
- [测试数据](#测试数据)
- [验证方法](#验证方法)
- [常见问题](#常见问题)
- [清理环境](#清理环境)

---

## 快速开始

### 一键准备测试环境

**Windows:**
```bash
# 运行准备脚本
scripts\test\permission-test-setup.bat

# 仅准备数据库
scripts\test\permission-test-setup.bat --db-only

# 跳过数据填充
scripts\test\permission-test-setup.bat --skip-data
```

**Linux/Mac:**
```bash
# 运行准备脚本
bash scripts/test/permission-test-setup.sh

# 仅准备数据库
bash scripts/test/permission-test-setup.sh --db-only

# 跳过数据填充
bash scripts/test/permission-test-setup.sh --skip-data
```

### 启动测试服务器

**Windows:**
```batch
# 设置环境变量
set QINGYU_DATABASE_NAME=qingyu_permission_test

# 启动服务器
godotenv -f .env.test go run cmd/server/main.go
```

**Linux/Mac:**
```bash
# 设置环境变量
export QINGYU_DATABASE_NAME=qingyu_permission_test

# 启动服务器
source .env.test
go run cmd/server/main.go
```

### 运行权限测试

```bash
# 运行所有权限测试
go test ./internal/middleware/auth/... -v

# 运行特定测试
go test ./internal/middleware/auth/ -run TestRBACChecker -v

# 运行集成测试
go test ./test/integration/... -v
```

---

## 环境要求

### 必需组件

| 组件 | 最低版本 | 推荐版本 | 下载地址 |
|------|---------|---------|---------|
| Go | 1.21+ | 1.22+ | https://golang.org/dl/ |
| MongoDB | 5.0+ | 7.0+ | https://www.mongodb.com/try/download |
| Redis | 6.0+ | 7.0+ | https://redis.io/download |

### 可选组件

| 组件 | 用途 | 下载地址 |
|------|------|---------|
| Docker | 快速启动Redis | https://www.docker.com/ |
| godotenv | 加载环境变量 | `go install github.com/joho/godotenv/cmd/godotenv@latest` |

### 检查环境

```bash
# 检查Go版本
go version

# 检查MongoDB
mongosh --version  # 或 mongo --version

# 检查Redis
redis-cli --version

# 检查Docker（可选）
docker --version
```

---

## 安装步骤

### 步骤 1: 启动MongoDB

**Windows:**
```batch
# 方法1: 使用Windows服务
net start MongoDB

# 方法2: 手动启动
"C:\Program Files\MongoDB\Server\7.0\bin\mongod.exe" --dbpath C:\data\db
```

**Linux:**
```bash
sudo systemctl start mongod
sudo systemctl enable mongod
```

**Mac:**
```bash
brew services start mongodb-community
```

**Docker:**
```bash
docker run -d \
  --name mongodb-test \
  -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME= \
  -e MONGO_INITDB_ROOT_PASSWORD= \
  mongo:7
```

### 步骤 2: 启动Redis

**Windows:**
```batch
# 方法1: 使用Windows服务
net start Redis

# 方法2: 手动启动
redis-server
```

**Linux:**
```bash
sudo systemctl start redis
sudo systemctl enable redis
```

**Mac:**
```bash
brew services start redis
```

**Docker:**
```bash
docker run -d \
  --name redis-test \
  -p 6379:6379 \
  redis:7-alpine
```

### 步骤 3: 验证连接

```bash
# 验证MongoDB
mongosh --eval "db.version()"

# 验证Redis
redis-cli ping
# 应该返回: PONG
```

### 步骤 4: 运行准备脚本

按照上面的[快速开始](#快速开始)中的说明运行准备脚本。

---

## 配置说明

### 环境变量配置

测试环境配置位于 `.env.test` 文件中：

```bash
# 数据库配置
QINGYU_DATABASE_NAME=qingyu_permission_test
QINGYU_DATABASE_URI=mongodb://localhost:27017

# Redis配置
QINGYU_REDIS_ADDR=localhost:6379
QINGYU_REDIS_DB=1

# JWT配置
QINGYU_JWT_SECRET=test_jwt_secret_key_for_permission_testing_12345

# 权限配置
PERMISSION_ENABLED=true
PERMISSION_CACHE_ENABLED=true
PERMISSION_LOAD_FROM_DB=true
```

### 权限配置文件

权限定义位于 `configs/permissions.yaml`：

```yaml
permissions:
  roles:
    admin:
      name: "管理员"
      description: "系统管理员，拥有所有权限"

    author:
      name: "作者"
      description: "内容创作者，可以管理自己的作品"

    reader:
      name: "读者"
      description: "普通读者，可以阅读内容"

  role_permissions:
    admin:
      - "*:*"  # 所有权限

    author:
      - "book:read"
      - "book:create"
      - "book:update"
      - "book:delete"

    reader:
      - "book:read"
      - "chapter:read"
```

### 中间件配置

权限中间件配置位于 `configs/middleware.yaml`：

```yaml
middleware:
  permission:
    enabled: true
    strategy: "rbac"
    config_path: "configs/permissions.yaml"

    skip_paths:
      - "/health"
      - "/api/v1/auth/login"

    rbac:
      load_from_db: true
      cache:
        enabled: true
        ttl: 5m
```

---

## 测试数据

### 角色列表

| 角色名 | 描述 | 权限数量 | 示例权限 |
|--------|------|---------|---------|
| admin | 系统管理员 | 1 | `*:*` (所有权限) |
| author | 作者 | 16 | `book:*`, `chapter:*`, `ai:generate` |
| reader | 读者 | 5 | `book:read`, `chapter:read` |
| editor | 编辑 | 11 | `book:review`, `chapter:review` |
| limited_user | 受限用户 | 1 | `book:read` |

### 测试账号

| 用户名 | 密码 | 角色 | VIP | 用途 |
|--------|------|------|-----|------|
| admin@test.com | Admin@123 | admin | 是 | 管理员测试 |
| author@test.com | Author@123 | author | 是 | 作者测试 |
| reader@test.com | Reader@123 | reader | 否 | 读者测试 |
| editor@test.com | Editor@123 | editor | 否 | 编辑测试 |
| limited@test.com | Limited@123 | limited_user | 否 | 受限用户测试 |
| author_reader@test.com | MultiRole@123 | author, reader | 是 | 多角色测试 |

### 权限测试场景

#### 1. 管理员场景
- ✓ 可以访问所有API
- ✓ 可以创建、读取、更新、删除任何资源
- ✓ 拥有通配符权限 `*:*`

#### 2. 作者场景
- ✓ 可以创建书籍 (`book:create`)
- ✓ 可以更新自己的书籍 (`book:update`)
- ✓ 可以删除自己的书籍 (`book:delete`)
- ✗ 不能审核书籍 (`book:review`)
- ✗ 不能管理用户 (`user:*`)

#### 3. 读者场景
- ✓ 可以阅读书籍 (`book:read`)
- ✓ 可以阅读章节 (`chapter:read`)
- ✗ 不能创建书籍 (`book:create`)
- ✗ 不能删除书籍 (`book:delete`)

#### 4. 编辑场景
- ✓ 可以审核内容 (`*:review`)
- ✓ 可以更新内容 (`*:update`)
- ✗ 不能删除内容 (`*:delete`)
- ✗ 不能创建内容 (`*:create`)

#### 5. 受限用户场景
- ✓ 可以读取书籍 (`book:read`)
- ✗ 不能读取章节 (`chapter:read`)
- ✗ 不能执行其他操作

#### 6. 多角色场景
- ✓ 拥有author和reader的所有权限
- ✓ 权限是累加的

---

## 验证方法

### 1. 验证数据库连接

```bash
# MongoDB
mongosh qingyu_permission_test --eval "db.roles.countDocuments()"
# 应该返回: 5

mongosh qingyu_permission_test --eval "db.users.countDocuments()"
# 应该返回: 6
```

### 2. 验证Redis连接

```bash
redis-cli
> KEYS user:permissions:*
# 应该返回权限缓存键列表
```

### 3. 验证API访问

#### 测试管理员登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@test.com",
    "password": "Admin@123"
  }'
```

预期响应:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "username": "admin@test.com",
      "roles": ["admin"]
    }
  }
}
```

#### 测试权限检查

```bash
# 使用admin token访问受保护资源
curl -X GET http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer <admin_token>"

# 应该返回200（有权限）

# 使用reader token尝试删除书籍
curl -X DELETE http://localhost:8080/api/v1/books/123 \
  -H "Authorization: Bearer <reader_token>"

# 应该返回403（无权限）
```

### 4. 验证中间件日志

启动服务器后，检查日志输出：

```
[INFO] RBACChecker已创建
[INFO] 从数据库加载权限到RBACChecker
[INFO] 加载角色权限 role=admin permissions=1
[INFO] 加载角色权限 role=author permissions=16
[INFO] 加载角色权限 role=reader permissions=5
[INFO] 权限加载完成 roles=5
```

### 5. 运行单元测试

```bash
# 测试RBAC检查器
go test ./internal/middleware/auth/ -run TestRBACChecker -v

# 测试权限中间件
go test ./internal/middleware/auth/ -run TestPermissionMiddleware -v

# 测试权限服务
go test ./service/shared/auth/ -run TestPermissionService -v

# 运行所有测试
go test ./internal/middleware/auth/... ./service/shared/auth/... -v
```

### 6. 运行集成测试

```bash
# API权限集成测试
go test ./test/api/ -run TestPermissionAPI -v

# 端到端测试
go test ./test/e2e/ -run TestAuthFlow -v
```

---

## 常见问题

### Q1: MongoDB连接失败

**错误信息:**
```
MongoDB连接失败: connection refused
```

**解决方法:**
1. 检查MongoDB是否启动:
   ```bash
   # Windows
   net start MongoDB

   # Linux
   sudo systemctl status mongod
   ```

2. 检查端口是否被占用:
   ```bash
   # 检查27017端口
   netstat -ano | findstr :27017
   ```

3. 修改MongoDB URI:
   ```bash
   export QINGYU_DATABASE_URI=mongodb://localhost:27018
   ```

### Q2: Redis连接失败

**错误信息:**
```
Redis连接失败: dial tcp: connect: connection refused
```

**解决方法:**
1. 使用Docker快速启动Redis:
   ```bash
   docker run -d -p 6379:6379 redis:7-alpine
   ```

2. 禁用Redis缓存:
   ```bash
   export PERMISSION_CACHE_ENABLED=false
   ```

### Q3: 权限检查总是失败

**可能原因:**
1. 用户未分配角色
2. 角色未分配权限
3. 权限格式不匹配

**排查步骤:**

```bash
# 1. 检查用户角色
mongosh qingyu_permission_test --eval '
  db.users.findOne({username: "admin@test.com"}, {roles: 1})
'

# 2. 检查角色权限
mongosh qingyu_permission_test --eval '
  db.roles.findOne({name: "admin"})
'

# 3. 检查权限格式
# 确保使用 ":" 分隔符，例如 "book:read" 而不是 "book.read"
```

### Q4: 测试数据填充失败

**错误信息:**
```
插入角色失败: duplicate key error
```

**解决方法:**
```bash
# 强制重建数据库
bash scripts/test/permission-test-setup.sh --force

# 或手动清理
mongosh qingyu_permission_test --eval '
  db.roles.deleteMany({})
  db.users.deleteMany({})
'
```

### Q5: 环境变量未生效

**解决方法:**

**Windows:**
```batch
# 使用godotenv工具
go install github.com/joho/godotenv/cmd/godotenv@latest
godotenv -f .env.test go run cmd/server/main.go
```

**Linux/Mac:**
```bash
# 方法1: 直接导出
export $(cat .env.test | grep -v '^#' | xargs)

# 方法2: 使用source
set -a
source .env.test
set +a
```

---

## 清理环境

### 清理测试数据

```bash
# 删除测试数据库
mongosh --eval "db.getSiblingDB('qingyu_permission_test').dropDatabase()"

# 或使用脚本
mongosh qingyu_permission_test --eval '
  db.roles.deleteMany({})
  db.users.deleteMany({})
'
```

### 清理Redis缓存

```bash
# 清除特定缓存
redis-cli KEYS "user:permissions:*" | xargs redis-cli DEL

# 清除所有缓存（谨慎）
redis-cli FLUSHDB
```

### 完全清理

```bash
# 1. 停止服务
# Ctrl+C 或 kill 进程

# 2. 清理数据库
mongosh --eval "db.getSiblingDB('qingyu_permission_test').dropDatabase()"

# 3. 清理Redis
redis-cli FLUSHDB

# 4. 清理日志
rm -f logs/*.log

# 5. 清理临时文件
rm -rf /tmp/qingyu_test_*
```

---

## 附录

### A. 测试场景清单

- [ ] 管理员访问所有API
- [ ] 作者创建书籍
- [ ] 作者更新自己的书籍
- [ ] 作者删除自己的书籍
- [ ] 读者阅读书籍
- [ ] 读者无法创建书籍
- [ ] 编辑审核内容
- [ ] 受限用户只有基本权限
- [ ] 多角色用户权限累加
- [ ] 通配符权限 `*:*`
- [ ] 通配符权限 `book:*`
- [ ] 权限缓存生效
- [ ] 权限缓存刷新

### B. 测试API端点

| 方法 | 路径 | 所需权限 | 测试用户 |
|------|------|---------|---------|
| POST | /api/v1/auth/login | 无 | 所有 |
| GET | /api/v1/books | book:read | admin, author, reader, editor |
| POST | /api/v1/books | book:create | admin, author |
| PUT | /api/v1/books/:id | book:update | admin, author |
| DELETE | /api/v1/books/:id | book:delete | admin, author |
| GET | /api/v1/books/:id/reviews | book:review | admin, editor |
| POST | /api/v1/admin/users | user:create | admin |

### C. 相关文档

- [权限系统集成指南](../middleware/permission-integration-guide.md)
- [中间件架构设计](../middleware/architecture.md)
- [RBAC模型说明](../models/auth/role.md)
- [API文档](../api/README.md)

### D. 故障排查命令

```bash
# 检查MongoDB连接
mongosh --eval "db.version()"

# 检查Redis连接
redis-cli ping

# 检查测试数据库
mongosh qingyu_permission_test --eval "db.getCollectionNames()"

# 检查角色数据
mongosh qingyu_permission_test --eval "db.roles.find().pretty()"

# 检查用户数据
mongosh qingyu_permission_test --eval "db.users.find({}, {username: 1, roles: 1}).pretty()"

# 检查Redis缓存
redis-cli KEYS "user:permissions:*"

# 查看服务器日志
tail -f logs/app.log

# 查看中间件日志
tail -f logs/middleware.log
```

---

**更新日志:**
- 2026-01-27: v1.0 初始版本
