# 用户 Repository 集成测试

## 概述

本目录包含用户管理模块的 Repository 层集成测试，用于验证：
- `UserRepository` - 用户数据访问操作
- `RoleRepository` - 角色数据访问操作

## 环境要求

集成测试需要连接到真实的 MongoDB 数据库。有两种方式运行测试：

### 方式一：使用 Docker（推荐）

1. **启动 Docker 数据库服务**

```bash
# 在项目根目录执行
cd docker
docker-compose -f docker-compose.db-only.yml up -d

# 等待服务就绪
docker-compose -f docker-compose.db-only.yml ps
```

2. **运行集成测试**

```bash
# 返回项目根目录
cd ..

# 运行用户 Repository 测试
go test -v ./test/repository/user/...

# 运行特定测试
go test -v ./test/repository/user -run TestUserRepository_Integration
go test -v ./test/repository/user -run TestRoleRepository_Integration
```

3. **停止 Docker 服务**

```bash
cd docker
docker-compose -f docker-compose.db-only.yml down

# 清理数据（可选）
docker-compose -f docker-compose.db-only.yml down -v
```

### 方式二：使用本地 MongoDB

如果本地已安装 MongoDB，确保：
- MongoDB 运行在 `localhost:27017`
- 配置文件 `config/config.yaml` 正确设置

## 测试内容

### UserRepository 测试

**基础操作测试** (`TestUserRepository_Integration`)
- ✅ 健康检查
- ✅ 创建用户
- ✅ 根据 ID 获取用户
- ✅ 根据 Email 获取用户
- ✅ 根据 Phone 获取用户
- ✅ 检查邮箱/手机号是否存在
- ✅ 更新用户信息
- ✅ 更新最后登录时间和 IP
- ✅ 更新用户状态
- ✅ 设置邮箱验证状态
- ✅ 高级查询（使用 Filter）
- ✅ 搜索用户
- ✅ 按角色/状态统计
- ✅ 删除用户

**批量操作测试** (`TestUserRepository_BatchOperations`)
- ✅ 批量创建用户
- ✅ 批量更新状态
- ✅ 批量删除（软删除）

### RoleRepository 测试

**基础操作测试** (`TestRoleRepository_Integration`)
- ✅ 健康检查
- ✅ 创建角色
- ✅ 根据 ID 获取角色
- ✅ 根据名称获取角色
- ✅ 检查角色名是否存在
- ✅ 更新角色信息
- ✅ 获取角色权限
- ✅ 添加权限
- ✅ 移除权限
- ✅ 更新角色权限列表
- ✅ 列出所有角色
- ✅ 列出默认角色
- ✅ 获取默认角色
- ✅ 按名称统计
- ✅ 删除角色

**默认角色测试** (`TestRoleRepository_DefaultRole`)
- ✅ 创建多个默认角色
- ✅ 获取默认角色
- ✅ 列出所有默认角色

## 跳过测试

如果没有可用的 MongoDB 环境，测试会自动跳过：

```bash
# 使用 -short 标志跳过集成测试
go test -v -short ./test/repository/user/...
```

## 测试数据清理

集成测试会在测试结束时自动清理测试数据。如需手动清理：

```bash
# 连接到 MongoDB
docker exec -it qingyu-mongodb mongosh

# 切换到测试数据库
use Qingyu_writer

# 删除测试用户（可选）
db.users.deleteMany({ username: { $regex: /^testuser_|^batch_user/ } })
db.roles.deleteMany({ name: { $regex: /^test_role_|^default_role/ } })
```

## 常见问题

### Q1: 测试失败：连接数据库超时

**解决方案**：
1. 确认 Docker 服务正在运行：`docker ps`
2. 检查 MongoDB 容器状态：`docker logs qingyu-mongodb`
3. 验证端口映射：确保 27017 端口未被占用

### Q2: 测试失败：配置文件找不到

**解决方案**：
1. 确认配置文件存在：`ls -la config/config.yaml`
2. 检查测试中的配置路径是否正确（相对路径：`../../../config/config.yaml`）

### Q3: 测试失败：重复键错误

**解决方案**：
这通常是因为之前的测试数据未清理。运行：
```bash
# 清理并重启 Docker
cd docker
docker-compose -f docker-compose.db-only.yml down -v
docker-compose -f docker-compose.db-only.yml up -d
```

## 性能基准

每个测试的预期执行时间（在 Docker 环境中）：
- `TestUserRepository_Integration`: ~2-3秒
- `TestUserRepository_BatchOperations`: ~1-2秒
- `TestRoleRepository_Integration`: ~2-3秒
- `TestRoleRepository_DefaultRole`: ~1秒

总测试时间：~6-9秒

## 下一步

完成 Repository 测试后，继续：
1. **Service 层单元测试** - 使用 Mock Repository
2. **API 层集成测试** - 测试完整的请求响应流程
3. **端到端测试** - 测试完整的用户注册登录流程

