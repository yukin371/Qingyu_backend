# 权限系统测试环境使用指南

## 快速开始

### 方法一：使用本地MongoDB/Redis

```bash
# 1. 确保MongoDB和Redis已启动
# Windows
net start MongoDB
net start Redis

# Linux/Mac
sudo systemctl start mongod
sudo systemctl start redis

# 2. 运行测试环境准备脚本
# Windows
scripts\test\permission-test-setup.bat

# Linux/Mac
bash scripts/test/permission-test-setup.sh

# 3. 设置环境变量并启动服务器
# Windows
set QINGYU_DATABASE_NAME=qingyu_permission_test
godotenv -f .env.test go run cmd/server/main.go

# Linux/Mac
export QINGYU_DATABASE_NAME=qingyu_permission_test
source .env.test
go run cmd/server/main.go

# 4. 运行测试
go test ./internal/middleware/auth/... -v
```

### 方法二：使用Docker

```bash
# 1. 启动Docker测试环境
docker-compose -f docker-compose.test.yml up -d

# 2. 等待服务启动（约10秒）
docker-compose -f docker-compose.test.yml ps

# 3. 运行测试环境准备脚本
bash scripts/test/permission-test-setup.sh

# 4. 启动服务器
godotenv -f .env.test go run cmd/server/main.go

# 5. 运行测试
go test ./internal/middleware/auth/... -v

# 6. 完成后清理
docker-compose -f docker-compose.test.yml down
```

## 测试账号

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin@test.com | Admin@123 | admin | 管理员，所有权限 |
| author@test.com | Author@123 | author | 作者，作品管理权限 |
| reader@test.com | Reader@123 | reader | 读者，只读权限 |
| editor@test.com | Editor@123 | editor | 编辑，审核权限 |
| limited@test.com | Limited@123 | limited_user | 受限用户 |
| author_reader@test.com | MultiRole@123 | author, reader | 多角色测试 |

## 快速测试

### 测试登录

```bash
# 管理员登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin@test.com", "password": "Admin@123"}'

# 读者登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "reader@test.com", "password": "Reader@123"}'
```

### 测试权限

```bash
# 获取书籍列表（所有角色都可以）
TOKEN="your_token_here"
curl -X GET http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer $TOKEN"

# 创建书籍（仅admin和author可以）
curl -X POST http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "测试书籍", "author": "测试作者"}'

# 删除书籍（reader应该失败）
curl -X DELETE http://localhost:8080/api/v1/books/123 \
  -H "Authorization: Bearer $TOKEN"
```

## 脚本选项

### permission-test-setup.sh / .bat

```bash
# 完整设置（包括数据）
permission-test-setup.sh

# 仅准备数据库
permission-test-setup.sh --db-only

# 仅设置数据库，不填充数据
permission-test-setup.sh --skip-data

# 强制重建数据库
permission-test-setup.sh --force

# 显示帮助
permission-test-setup.sh --help
```

### permission-test-data.go

```bash
# 使用默认数据库
go run scripts/test/permission-test-data.go

# 指定数据库
go run scripts/test/permission-test-data.go --db=my_test_db

# 指定MongoDB URI
go run scripts/test/permission-test-data.go --uri=mongodb://localhost:27018

# 详细输出
go run scripts/test/permission-test-data.go -v
```

## 验证测试环境

```bash
# 检查MongoDB数据
mongosh qingyu_permission_test --eval "
  print('角色数量:', db.roles.countDocuments());
  print('用户数量:', db.users.countDocuments());
  print('角色列表:', db.roles.distinct('name'));
"

# 检查Redis缓存
redis-cli KEYS "user:permissions:*"

# 运行单元测试
go test ./internal/middleware/auth/... -v

# 运行集成测试
go test ./test/integration/... -v
```

## 清理测试环境

```bash
# 删除测试数据库
mongosh --eval "db.getSiblingDB('qingyu_permission_test').dropDatabase()"

# 清理Redis缓存
redis-cli KEYS "user:permissions:*" | xargs redis-cli DEL

# 停止Docker服务
docker-compose -f docker-compose.test.yml down -v

# 完整清理
mongosh qingyu_permission_test --eval "
  db.roles.deleteMany({});
  db.users.deleteMany({});
"
redis-cli FLUSHDB
```

## 故障排查

### MongoDB连接失败

```bash
# 检查MongoDB是否运行
mongosh --eval "db.version()"

# 启动MongoDB
# Windows
net start MongoDB

# Linux
sudo systemctl start mongod

# Mac
brew services start mongodb-community

# Docker
docker run -d -p 27017:27017 mongo:7
```

### Redis连接失败

```bash
# 检查Redis是否运行
redis-cli ping

# 启动Redis
# Windows
net start Redis

# Linux
sudo systemctl start redis

# Mac
brew services start redis

# Docker
docker run -d -p 6379:6379 redis:7-alpine
```

### 权限检查失败

```bash
# 1. 检查用户角色
mongosh qingyu_permission_test --eval "
  db.users.findOne({username: 'admin@test.com'}, {roles: 1})
"

# 2. 检查角色权限
mongosh qingyu_permission_test --eval "
  db.roles.findOne({name: 'admin'})
"

# 3. 重新加载权限
mongosh qingyu_permission_test --eval "
  db.roles.find().forEach(printjson)
"
```

## 相关文档

- [测试环境准备完整指南](../../../docs/testing/permission-test-setup.md)
- [权限系统集成指南](../../../docs/middleware/permission-integration-guide.md)
- [中间件配置](../../../configs/middleware.yaml)
- [权限配置](../../../configs/permissions.yaml)
