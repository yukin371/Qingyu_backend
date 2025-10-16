# 集成测试说明

本目录包含用户管理模块的集成测试，测试真实的HTTP API流程。

---

## 📋 测试文件

### 1. `user_api_integration_test.go`

完整的用户管理API集成测试，包括：

- **完整用户生命周期测试**
  - 用户注册
  - 用户登录
  - 获取个人信息
  - 更新个人信息
  - 修改密码
  - 使用新密码登录

- **认证和权限控制测试**
  - 未认证访问
  - 无效Token访问
  - 普通用户访问管理员接口

- **管理员用户管理测试**（待实现）
  - 获取用户列表
  - 更新用户信息
  - 删除用户

---

## 🚀 运行测试

### 前提条件

1. **MongoDB 运行中**
   ```bash
   # 使用 Docker Compose 启动测试数据库
   cd ../../
   docker-compose -f docker-compose-db.yaml up -d
   ```

2. **配置文件**
   - 确保 `config/config.yaml` 配置正确
   - MongoDB 连接信息正确

### 运行所有集成测试

```bash
# 从项目根目录运行
go test ./test/integration/ -v

# 或指定特定测试
go test ./test/integration/ -v -run TestUserAPI_Integration
```

### 跳过集成测试（仅运行单元测试）

```bash
go test ./test/integration/ -v -short
```

---

## 🧪 测试场景

### 场景1：完整用户生命周期

```
1. 注册用户
   POST /api/v1/register
   → 返回用户信息和Token

2. 登录
   POST /api/v1/login
   → 返回新Token

3. 获取个人信息
   GET /api/v1/users/profile
   Header: Authorization: Bearer <token>
   → 返回用户详细信息

4. 更新个人信息
   PUT /api/v1/users/profile
   Body: { "nickname": "...", "bio": "..." }
   → 更新成功

5. 修改密码
   PUT /api/v1/users/password
   Body: { "old_password": "...", "new_password": "..." }
   → 修改成功

6. 使用新密码登录
   POST /api/v1/login
   → 登录成功
```

### 场景2：认证和权限控制

```
1. 未认证访问
   GET /api/v1/users/profile
   （不带Authorization header）
   → 401 Unauthorized

2. 无效Token访问
   GET /api/v1/users/profile
   Header: Authorization: Bearer invalid_token
   → 401 Unauthorized

3. 普通用户访问管理员接口
   GET /api/v1/admin/users
   Header: Authorization: Bearer <user_token>
   → 403 Forbidden
```

---

## 📊 测试输出示例

```
=== RUN   TestUserAPI_Integration
=== RUN   TestUserAPI_Integration/完整用户生命周期
=== RUN   TestUserAPI_Integration/完整用户生命周期/用户注册
    user_api_integration_test.go:95: ✓ 用户注册成功: ID=670abcdef123456789, Username=testuser_1697203200
=== RUN   TestUserAPI_Integration/完整用户生命周期/用户登录
    user_api_integration_test.go:127: ✓ 用户登录成功，获得新Token
=== RUN   TestUserAPI_Integration/完整用户生命周期/获取个人信息
    user_api_integration_test.go:149: ✓ 获取个人信息成功
=== RUN   TestUserAPI_Integration/完整用户生命周期/更新个人信息
    user_api_integration_test.go:174: ✓ 更新个人信息成功
=== RUN   TestUserAPI_Integration/完整用户生命周期/修改密码
    user_api_integration_test.go:201: ✓ 修改密码成功
=== RUN   TestUserAPI_Integration/完整用户生命周期/使用新密码登录
    user_api_integration_test.go:220: ✓ 使用新密码登录成功
=== RUN   TestUserAPI_Integration/认证和权限控制
=== RUN   TestUserAPI_Integration/认证和权限控制/未认证访问需要认证的接口
    user_api_integration_test.go:245: ✓ 未认证访问被正确拒绝
=== RUN   TestUserAPI_Integration/认证和权限控制/使用无效Token访问
    user_api_integration_test.go:262: ✓ 无效Token被正确拒绝
=== RUN   TestUserAPI_Integration/认证和权限控制/普通用户访问管理员接口
    user_api_integration_test.go:310: ✓ 普通用户访问管理员接口被正确拒绝
--- PASS: TestUserAPI_Integration (2.34s)
    --- PASS: TestUserAPI_Integration/完整用户生命周期 (1.82s)
        --- PASS: TestUserAPI_Integration/完整用户生命周期/用户注册 (0.31s)
        --- PASS: TestUserAPI_Integration/完整用户生命周期/用户登录 (0.25s)
        --- PASS: TestUserAPI_Integration/完整用户生命周期/获取个人信息 (0.18s)
        --- PASS: TestUserAPI_Integration/完整用户生命周期/更新个人信息 (0.21s)
        --- PASS: TestUserAPI_Integration/完整用户生命周期/修改密码 (0.23s)
        --- PASS: TestUserAPI_Integration/完整用户生命周期/使用新密码登录 (0.24s)
    --- PASS: TestUserAPI_Integration/认证和权限控制 (0.52s)
        --- PASS: TestUserAPI_Integration/认证和权限控制/未认证访问需要认证的接口 (0.05s)
        --- PASS: TestUserAPI_Integration/认证和权限控制/使用无效Token访问 (0.04s)
        --- PASS: TestUserAPI_Integration/认证和权限控制/普通用户访问管理员接口 (0.43s)
PASS
ok      Qingyu_backend/test/integration 2.345s
```

---

## 🔍 故障排查

### 问题1：无法连接到MongoDB

**错误**: `failed to connect to MongoDB`

**解决方案**:
```bash
# 检查MongoDB是否运行
docker ps | grep mongo

# 启动MongoDB
docker-compose -f docker-compose-db.yaml up -d
```

### 问题2：配置文件找不到

**错误**: `cannot find config file`

**解决方案**:
```bash
# 确保从项目根目录运行测试
cd /path/to/Qingyu_backend
go test ./test/integration/ -v
```

### 问题3：测试数据未清理

**解决方案**:
```bash
# 手动清理测试数据
mongo qingyu_test --eval "db.users.deleteMany({username: /^(testuser_|normaluser_|admin_)/})"
```

---

## 📝 注意事项

1. **测试数据清理**
   - 测试会自动清理以 `testuser_`、`normaluser_`、`admin_` 开头的测试用户
   - 测试结束后会自动执行清理

2. **测试隔离**
   - 每个测试使用唯一的时间戳作为用户名
   - 确保测试之间不会互相干扰

3. **数据库选择**
   - 使用配置文件中指定的数据库
   - 建议使用专门的测试数据库

4. **测试跳过**
   - 使用 `-short` 标志跳过集成测试
   - 适用于不想启动数据库的快速测试

---

## 🎯 下一步

- [ ] 实现管理员用户管理测试
- [ ] 添加更多边界条件测试
- [ ] 添加性能测试
- [ ] 添加并发测试

---

**更新时间**: 2025-10-13  
**维护者**: AI Assistant


