# API 层测试说明

> **创建日期**: 2025-10-18  
> **状态**: 测试框架已创建，待完善

---

## 📋 概述

本目录包含 API 层的集成测试，重点测试 HTTP 请求/响应处理逻辑。

---

## 📁 测试文件

### 已创建的测试框架

| 文件 | 模块 | 状态 | 说明 |
|------|------|------|------|
| `user_api_test.go` | 用户系统 | 🚧 框架完成 | 需要完善 Mock |
| `auth_api_test.go` | 认证系统 | 🚧 框架完成 | 需要完善 Mock |
| `reader_api_integration_test.go` | 阅读器 | 🚧 框架完成 | 需要 Repository Mock |
| `bookstore_api_test.go` | 书城 | ✅ 部分完成 | 已有基础测试 |

### Build Tags

所有新增的测试文件使用 `integration` build tag：

```go
//go:build integration
// +build integration
```

这意味着：
- ✅ 不会影响正常编译
- ✅ 可以单独运行集成测试
- ✅ 作为测试框架保留

---

## 🚀 运行测试

### 运行所有测试

```bash
go test ./test/api/... -tags=integration -v
```

### 运行特定测试

```bash
# 用户 API 测试
go test ./test/api/ -tags=integration -v -run TestUser

# 认证 API 测试
go test ./test/api/ -tags=integration -v -run TestAuth

# 书城 API 测试（不需要 tag）
go test ./test/api/ -v -run TestBookstore
```

---

## 📝 测试覆盖

### 用户系统 API (user_api_test.go)

**测试场景**: 12+

- ✅ 用户注册（成功、失败、参数验证）
- ✅ 用户登录（成功、失败）
- ✅ 获取个人信息
- ✅ 更新个人信息
- ✅ 修改密码
- ✅ 管理员获取用户列表
- ✅ 管理员获取用户详情
- ✅ 管理员更新用户
- ✅ 管理员删除用户
- ✅ 未认证访问测试

**依赖**: MockUserRepository

### 认证系统 API (auth_api_test.go)

**测试场景**: 10+

- ✅ 注册
- ✅ 登录
- ✅ 登出
- ✅ Token 刷新
- ✅ 获取用户权限
- ✅ 获取用户角色
- ✅ 参数验证
- ✅ 认证检查

**依赖**: MockAuthService

### 阅读器 API (reader_api_integration_test.go)

**测试场景**: 8+（框架）

- 📝 获取章节信息
- 📝 获取章节内容
- 📝 获取章节列表
- 📝 获取阅读设置
- 📝 保存阅读设置
- 📝 参数验证

**依赖**: Repository Mocks（待实现）

---

## 🛠️ 待完善工作

### 1. Mock 完善

需要完善以下 Mock 实现：

#### MockUserRepository
- ❌ `BatchDelete` 方法

#### MockAuthService
- ❌ `AssignRole` 方法
- ❌ `CreateRole` 方法
- ❌ `UpdateRole` 方法
- ❌ `DeleteRole` 方法
- ❌ `RemoveRole` 方法
- ❌ 其他角色管理方法

### 2. 类型适配

#### Auth API 测试
- 需要适配 `RegisterResponse` 和 `LoginResponse` 结构
  - 当前: `{UserID, Username, Email, Token}`
  - 实际: `{User: *UserInfo, Token}`

#### Reader API 测试
- 需要适配 `ReadingSettings` 结构
- 需要实现 Repository Mocks

### 3. 集成测试环境

需要：
- 测试数据库连接
- 测试数据准备
- 清理机制

---

## 💡 使用建议

### 当前阶段

1. **参考测试框架**
   - 测试文件提供了完整的测试结构
   - 可以作为编写其他 API 测试的模板

2. **运行现有测试**
   - 书城 API 测试可以正常运行
   - 其他测试需要完善后运行

### 完善步骤

1. **完善 Mock 实现**
   ```bash
   # 参考现有的 Service Mock
   # 实现缺失的方法
   ```

2. **适配类型**
   ```bash
   # 检查实际的 Response 结构
   # 更新测试代码中的断言
   ```

3. **运行测试**
   ```bash
   go test ./test/api/ -tags=integration -v
   ```

4. **修复失败的测试**
   - 检查 Mock 行为
   - 验证请求/响应格式
   - 确保中间件正常工作

---

## 📚 参考资料

### 现有测试

- `test/service/` - Service 层测试（包含 Mock 示例）
- `test/repository/` - Repository 层测试
- `test/api/bookstore_api_test.go` - 可运行的 API 测试示例

### 测试工具

- **testify/mock**: Mock 框架
- **testify/assert**: 断言库
- **httptest**: HTTP 测试工具

### 最佳实践

1. **测试结构**
   - Setup: 创建测试环境
   - Execute: 执行 HTTP 请求
   - Assert: 验证响应

2. **Mock 使用**
   - 使用 `mock.On()` 设置期望
   - 使用 `AssertExpectations()` 验证调用

3. **HTTP 测试**
   - 使用 `httptest.NewRecorder()` 记录响应
   - 验证状态码、响应体、Headers

---

## 🎯 下一步

### 短期（1周内）

1. 完善 MockUserRepository 和 MockAuthService
2. 适配 Auth API 测试的类型
3. 运行并修复所有测试

### 中期（2-4周）

1. 创建 Reader API 的 Repository Mocks
2. 补充更多测试场景
3. 达到 80% 测试覆盖率

### 长期（持续）

1. 添加性能测试
2. 添加并发测试
3. 集成到 CI/CD

---

**维护者**: 青羽后端团队  
**最后更新**: 2025-10-18
