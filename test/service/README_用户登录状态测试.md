# 用户登录状态测试文档

**文件**: `test/service/user_login_status_test.go`  
**创建日期**: 2025-10-17  
**关联修复**: [登录状态检查修复](../../doc/implementation/登录状态检查修复完成_2025-10-17.md)

---

## 📖 概述

本测试文件针对用户登录逻辑的状态检查功能进行全面测试，确保只有 `active` 状态的用户可以登录，其他状态（`inactive`、`banned`、`deleted`）的用户会被正确拒绝。

---

## 🎯 测试覆盖

### 测试场景列表

| 测试函数 | 测试场景 | 预期结果 |
|---------|---------|---------|
| `TestUserService_LoginUser_ActiveStatus_Success` | 活跃用户登录 | ✅ 成功，返回 token |
| `TestUserService_LoginUser_InactiveStatus_Rejected` | 未激活用户登录 | ❌ 拒绝，提示"账号未激活" |
| `TestUserService_LoginUser_BannedStatus_Rejected` | 已封禁用户登录 | ❌ 拒绝，提示"已被封禁" |
| `TestUserService_LoginUser_DeletedStatus_Rejected` | 已删除用户登录 | ❌ 拒绝，提示"账号已删除" |
| `TestUserService_LoginUser_WrongPassword` | 密码错误 | ❌ 拒绝，提示"密码错误" |
| `TestUserService_LoginUser_UserNotFound` | 用户不存在 | ❌ 拒绝，提示"用户不存在" |
| `TestUserService_LoginUser_AllStatuses` | 表驱动：所有状态 | 批量验证所有状态 |

---

## 📝 测试设计

### 遵循的最佳实践

根据 [测试最佳实践文档](../../doc/testing/测试最佳实践.md)，本测试实现了以下最佳实践：

#### 1. ✅ AAA 测试模式

```go
func TestUserService_LoginUser_ActiveStatus_Success(t *testing.T) {
    // ===== Arrange (准备) =====
    mockRepo := new(MockUserRepository)
    service := user.NewUserService(mockRepo)
    activeUser := testutil.CreateTestUser(...)
    mockRepo.On("GetByUsername", ...).Return(activeUser, nil)
    
    // ===== Act (执行) =====
    resp, err := service.LoginUser(ctx, &interfaces.LoginUserRequest{...})
    
    // ===== Assert (断言) =====
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    mockRepo.AssertExpectations(t)
}
```

#### 2. ✅ 使用 Mock 接口

```go
// MockUserRepository 实现完整的 UserRepository 接口
type MockUserRepository struct {
    mock.Mock
}

// 实现所有接口方法
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*users.User, error) {
    args := m.Called(ctx, username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*users.User), args.Error(1)
}
```

#### 3. ✅ 表驱动测试

```go
func TestUserService_LoginUser_AllStatuses(t *testing.T) {
    tests := []struct {
        name          string
        userStatus    users.UserStatus
        expectSuccess bool
        errorContains string
    }{
        {
            name:          "活跃用户可以登录",
            userStatus:    users.UserStatusActive,
            expectSuccess: true,
        },
        {
            name:          "未激活用户被拒绝",
            userStatus:    users.UserStatusInactive,
            expectSuccess: false,
            errorContains: "账号未激活",
        },
        // ... 更多测试用例
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试实现
        })
    }
}
```

#### 4. ✅ 清晰的测试命名

- 格式：`TestServiceName_MethodName_Scenario_ExpectedResult`
- 示例：`TestUserService_LoginUser_BannedStatus_Rejected`
- 每个测试名称都清楚地描述了测试场景和预期结果

#### 5. ✅ 测试隔离

- 每个测试创建独立的 Mock 实例
- 使用 `testutil.CreateTestUser()` 生成独立的测试数据
- 测试之间互不影响

---

## 🔧 Mock 实现

### MockUserRepository

完整实现了 `UserRepository` 接口的所有方法：

**基础 CRUD**:
- `Create()`
- `GetByID()`
- `GetByUsername()`
- `GetByEmail()`
- `GetByPhone()`
- `Update()`
- `Delete()`
- `List()`

**状态管理**:
- `ExistsByUsername()`
- `ExistsByEmail()`
- `ExistsByPhone()`
- `UpdateLastLogin()`
- `UpdatePassword()`
- `UpdateStatus()`
- `SetEmailVerified()`
- `SetPhoneVerified()`

**高级功能**:
- `GetActiveUsers()`
- `GetUsersByRole()`
- `BatchUpdateStatus()`
- `BatchDelete()`
- `FindWithFilter()`
- `SearchUsers()`
- `CountByRole()`
- `CountByStatus()`
- `Transaction()`
- `Count()`
- `Health()`

---

## 🧪 测试详细说明

### 1. 活跃用户登录测试

```go
func TestUserService_LoginUser_ActiveStatus_Success(t *testing.T)
```

**目的**: 验证正常用户可以成功登录

**步骤**:
1. 创建 `active` 状态的测试用户
2. Mock `GetByUsername` 返回该用户
3. Mock `UpdateLastLogin` 成功
4. 调用 `LoginUser` 方法
5. 验证返回 token 且无错误

**断言**:
- `assert.NoError(t, err)`
- `assert.NotNil(t, resp)`
- `assert.NotEmpty(t, resp.Token)`
- `assert.Equal(t, "activeuser", resp.User.Username)`

---

### 2. 未激活用户登录测试

```go
func TestUserService_LoginUser_InactiveStatus_Rejected(t *testing.T)
```

**目的**: 验证未激活用户无法登录

**预期行为**:
- 密码验证通过
- 状态检查失败
- 返回错误：`"账号未激活，请先验证邮箱"`

**断言**:
- `assert.Error(t, err)`
- `assert.Nil(t, resp)`
- `assert.Contains(t, err.Error(), "账号未激活")`
- `assert.Contains(t, err.Error(), "验证邮箱")`

---

### 3. 已封禁用户登录测试

```go
func TestUserService_LoginUser_BannedStatus_Rejected(t *testing.T)
```

**目的**: 验证已封禁用户无法登录

**预期行为**:
- 密码验证通过
- 状态检查失败
- 返回错误：`"账号已被封禁，请联系管理员"`

**断言**:
- `assert.Error(t, err)`
- `assert.Contains(t, err.Error(), "已被封禁")`
- `assert.Contains(t, err.Error(), "联系管理员")`

---

### 4. 已删除用户登录测试

```go
func TestUserService_LoginUser_DeletedStatus_Rejected(t *testing.T)
```

**目的**: 验证已删除用户无法登录

**预期行为**:
- 密码验证通过
- 状态检查失败
- 返回错误：`"账号已删除"`

**断言**:
- `assert.Error(t, err)`
- `assert.Contains(t, err.Error(), "账号已删除")`

---

### 5. 密码错误测试

```go
func TestUserService_LoginUser_WrongPassword(t *testing.T)
```

**目的**: 验证密码验证在状态检查之前

**预期行为**:
- 密码验证失败（在状态检查之前）
- 返回错误：`"密码错误"`

**断言**:
- `assert.Error(t, err)`
- `assert.Contains(t, err.Error(), "密码错误")`

---

### 6. 用户不存在测试

```go
func TestUserService_LoginUser_UserNotFound(t *testing.T)
```

**目的**: 验证用户不存在的情况

**预期行为**:
- Repository 返回 `NotFoundError`
- 返回错误：`"用户不存在"`

**断言**:
- `assert.Error(t, err)`
- `assert.Contains(t, err.Error(), "用户不存在")`

---

### 7. 表驱动测试：所有状态

```go
func TestUserService_LoginUser_AllStatuses(t *testing.T)
```

**目的**: 批量验证所有用户状态的行为

**测试用例**:
| 状态 | 应该成功？ | 错误信息 |
|------|----------|---------|
| `active` | ✅ 是 | - |
| `inactive` | ❌ 否 | "账号未激活" |
| `banned` | ❌ 否 | "已被封禁" |
| `deleted` | ❌ 否 | "账号已删除" |

---

## 🏃 运行测试

### 运行所有登录状态测试

```bash
# 运行整个测试文件
go test -v ./test/service/user_login_status_test.go

# 运行特定测试
go test -v ./test/service -run TestUserService_LoginUser_ActiveStatus

# 运行表驱动测试
go test -v ./test/service -run TestUserService_LoginUser_AllStatuses

# 带覆盖率运行
go test -v -cover ./test/service/user_login_status_test.go
```

### 运行所有 Service 测试

```bash
go test -v ./test/service/...
```

---

## 📊 预期测试结果

### 成功输出示例

```
=== RUN   TestUserService_LoginUser_ActiveStatus_Success
--- PASS: TestUserService_LoginUser_ActiveStatus_Success (0.00s)
=== RUN   TestUserService_LoginUser_InactiveStatus_Rejected
--- PASS: TestUserService_LoginUser_InactiveStatus_Rejected (0.00s)
=== RUN   TestUserService_LoginUser_BannedStatus_Rejected
--- PASS: TestUserService_LoginUser_BannedStatus_Rejected (0.00s)
=== RUN   TestUserService_LoginUser_DeletedStatus_Rejected
--- PASS: TestUserService_LoginUser_DeletedStatus_Rejected (0.00s)
=== RUN   TestUserService_LoginUser_WrongPassword
--- PASS: TestUserService_LoginUser_WrongPassword (0.00s)
=== RUN   TestUserService_LoginUser_UserNotFound
--- PASS: TestUserService_LoginUser_UserNotFound (0.00s)
=== RUN   TestUserService_LoginUser_AllStatuses
=== RUN   TestUserService_LoginUser_AllStatuses/活跃用户可以登录
=== RUN   TestUserService_LoginUser_AllStatuses/未激活用户被拒绝
=== RUN   TestUserService_LoginUser_AllStatuses/已封禁用户被拒绝
=== RUN   TestUserService_LoginUser_AllStatuses/已删除用户被拒绝
--- PASS: TestUserService_LoginUser_AllStatuses (0.00s)
PASS
ok      Qingyu_backend/test/service     0.123s
```

---

## 🔗 相关文档

- [登录状态检查修复详细说明](../../doc/implementation/用户登录状态检查修复_2025-10-17.md)
- [登录状态检查修复完成报告](../../doc/implementation/登录状态检查修复完成_2025-10-17.md)
- [测试最佳实践](../../doc/testing/测试最佳实践.md)
- [测试组织规范](../../doc/testing/测试组织规范.md)
- [用户模型定义](../../models/users/user.go)
- [用户服务实现](../../service/user/user_service.go)

---

## ✅ 测试完成清单

- [x] 创建完整的 MockUserRepository
- [x] 测试活跃用户成功登录
- [x] 测试未激活用户被拒绝
- [x] 测试已封禁用户被拒绝
- [x] 测试已删除用户被拒绝
- [x] 测试密码错误场景
- [x] 测试用户不存在场景
- [x] 实现表驱动测试
- [x] 遵循 AAA 模式
- [x] 使用测试工具函数（testutil）
- [x] 验证 Mock 期望
- [x] 通过 Linter 检查

---

**创建时间**: 2025-10-17  
**维护者**: 青羽后端团队  
**状态**: ✅ 已完成

