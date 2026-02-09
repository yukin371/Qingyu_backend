# User模块TDD迁移实施计划

> **计划日期**: 2026-02-10  
> **制定者**: 架构重构女仆  
> **预期工期**: 10-16个工作日

---

## 📋 总体策略

### TDD方法论
采用严格的**Red-Green-Refactor**循环：

```
┌─────────────────────────────────────────────────────────┐
│  RED (编写失败的测试)                                     │
│  → 编写测试用例，描述预期行为                              │
│  → 运行测试，确认失败（红灯）                               │
│  → 记录失败的测试场景                                      │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  GREEN (编写最小可行代码)                                 │
│  → 编写最简单的代码使测试通过                              │
│  → 运行测试，确认通过（绿灯）                               │
│  → 不考虑代码质量和设计                                    │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  REFACTOR (重构优化代码)                                  │
│  → 重构代码，改进设计                                      │
│  → 保持测试通过                                           │
│  → 应用设计模式和最佳实践                                  │
└─────────────────────────────────────────────────────────┘
```

### 测试金字塔

```
         ┌──────────┐
         │  E2E测试  │  15% (端到端场景)
         │   (15)   │
         ├──────────┤
         │ 集成测试  │  30% (组件间交互)
         │   (30)   │
         ├──────────┤
         │ 单元测试  │  55% (独立功能)
         │   (55)   │
         └──────────┘
         总计: ~100个测试用例
```

---

## 🚀 阶段1: 基础组件迁移 (Day 1-2)

### 目标
建立坚实的测试基础，迁移低风险的基础组件。

### 迁移清单

#### 1.1 Constants常量 (0.5天)

**文件**: `service/user/constants.go`

**迁移内容**:
```go
// 验证码相关常量
- VerificationCodeLength
- VerificationCodeExpiry
- VerificationRateLimitCount

// 密码重置相关
- PasswordResetTokenExpiry

// Token相关
- TokenDefaultExpiry

// 验证目的常量
- VerificationPurposeEmail
- VerificationPurposePhone
- VerificationPurposeReset
```

**测试用例** (无，纯常量):
- ✅ 编译验证
- ✅ 常量值检查

**验收标准**:
- [ ] 常量定义完整
- [ ] 编译无错误
- [ ] 文档注释完整

---

#### 1.2 UserValidator用户验证器 (0.5天)

**文件**: `service/user/user_validator.go`

**测试用例** (15个):

| # | 测试场景 | 预期结果 |
|---|---------|---------|
| 1 | 用户名为空 | 返回"用户名不能为空"错误 |
| 2 | 用户名过短(<3字符) | 返回"长度不能少于3个字符" |
| 3 | 用户名过长(>30字符) | 返回"长度不能超过30个字符" |
| 4 | 用户名包含非法字符 | 返回"只能包含字母、数字和下划线" |
| 5 | 用户名以数字开头 | 返回"不能以数字开头" |
| 6 | 用户名为保留名 | 返回"系统保留，不能使用" |
| 7 | 邮箱格式错误 | 返回"邮箱格式不正确" |
| 8 | 密码过短(<8字符) | 返回"长度不能少于8个字符" |
| 9 | 密码无字母 | 返回"必须包含字母" |
| 10 | 密码无数字 | 返回"必须包含数字" |
| 11 | 密码为弱密码 | 返回"密码过于简单" |
| 12 | 用户名唯一性检查 | 检查数据库中的唯一性 |
| 13 | 邮箱唯一性检查 | 检查数据库中的唯一性 |
| 14 | 业务规则验证 | 用户名与邮箱前缀不能相同 |
| 15 | 创建时间验证 | 创建时间不能是未来时间 |

**TDD执行**:

```go
// RED: 编写失败的测试
func TestUserValidator_ValidateUsername_Empty(t *testing.T) {
    validator := NewUserValidator(nil)
    user := &usersModel.User{Username: ""}
    
    err := validator.validateUsername(user.Username)
    
    assert.Error(t, err)
    assert.Equal(t, "REQUIRED", err.Code)
}

// GREEN: 实现最小代码
func (v *UserValidator) validateUsername(username string) *ValidationError {
    if username == "" {
        return &ValidationError{Field: "username", Code: "REQUIRED"}
    }
    return nil
}

// REFACTOR: 重构优化
func (v *UserValidator) validateUsername(username string) *ValidationError {
    // 完整的实现...
}
```

**验收标准**:
- [ ] 15个测试全部通过
- [ ] 代码覆盖率 > 90%
- [ ] Mock repository正确使用

---

#### 1.3 PasswordValidator密码验证器 (0.5天)

**文件**: `service/user/password_validator.go`

**测试用例** (10个):

| # | 测试场景 | 预期结果 |
|---|---------|---------|
| 1 | 密码过短(<8字符) | 返回false, "长度不能少于8位" |
| 2 | 密码无大写字母 | 返回false, "必须包含大写字母" |
| 3 | 密码无小写字母 | 返回false, "必须包含小写字母" |
| 4 | 密码无数字 | 返回false, "必须包含数字" |
| 5 | 密码为常见弱密码 | 返回false, "密码过于常见" |
| 6 | 密码包含连续字符(123) | 返回false, "不能包含连续字符" |
| 7 | 密码包含连续字符(abc) | 返回false, "不能包含连续字符" |
| 8 | 强密码评分 | 返回评分 ≥ 80 |
| 9 | 中等密码评分 | 返回评分 60-79 |
| 10 | 弱密码评分 | 返回评分 < 40 |

**验收标准**:
- [ ] 10个测试全部通过
- [ ] 密码强度评分算法正确
- [ ] 弱密码字典完整

---

#### 1.4 Converter转换器 (0.5天)

**文件**: `service/user/converter.go`

**测试用例** (8个):

| # | 测试场景 | 预期结果 |
|---|---------|---------|
| 1 | User → UserDTO转换 | 字段完整映射 |
| 2 | UserDTO → User转换 | 字段完整映射 |
| 3 | 批量User → DTO转换 | 列表完整转换 |
| 4 | 时间格式转换 | ISO8601格式正确 |
| 5 | ObjectID转换 | hex字符串正确 |
| 6 | 状态枚举转换 | 枚举值正确 |
| 7 | 空值处理 | nil安全处理 |
| 8 | 边界值处理 | 极限值正确处理 |

**验收标准**:
- [ ] 8个测试全部通过
- [ ] 时间转换无精度丢失
- [ ] ID转换无错误

---

### 阶段1验收
- [ ] 所有基础组件迁移完成
- [ ] 测试覆盖率 > 85%
- [ ] 编译无错误
- [ ] 代码审查通过
- [ ] 更新Serena记忆

---

## 🚀 阶段2: 核心UserService迁移 (Day 3-6)

### 目标
迁移最复杂的UserService，采用分阶段策略。

### Phase 2.1: 基础CRUD (Day 3)

**迁移方法**:
1. CreateUser
2. GetUser
3. UpdateUser
4. DeleteUser
5. ListUsers

**测试用例** (15个):

| # | 测试场景 | 测试方法 |
|---|---------|---------|
| 1 | 创建用户成功 | TestUserService_CreateUser_Success |
| 2 | 用户名已存在 | TestUserService_CreateUser_DuplicateUsername |
| 3 | 邮箱已存在 | TestUserService_CreateUser_DuplicateEmail |
| 4 | 空用户名 | TestUserService_CreateUser_EmptyUsername |
| 5 | 空邮箱 | TestUserService_CreateUser_EmptyEmail |
| 6 | 获取用户成功 | TestUserService_GetUser_Success |
| 7 | 用户不存在 | TestUserService_GetUser_NotFound |
| 8 | 空ID | TestUserService_GetUser_EmptyID |
| 9 | 更新用户成功 | TestUserService_UpdateUser_Success |
| 10 | 用户不存在 | TestUserService_UpdateUser_NotFound |
| 11 | 空更新数据 | TestUserService_UpdateUser_EmptyUpdates |
| 12 | 删除用户成功 | TestUserService_DeleteUser_Success |
| 13 | 用户不存在 | TestUserService_DeleteUser_NotFound |
| 14 | 列出用户成功 | TestUserService_ListUsers_Success |
| 15 | 分页正确 | TestUserService_ListUsers_Pagination |

**Mock设置**:
```go
func setupUserService() (*UserServiceImpl, *MockUserRepository, *MockAuthRepository) {
    mockUserRepo := new(MockUserRepository)
    mockAuthRepo := new(MockAuthRepository)
    service := NewUserService(mockUserRepo, mockAuthRepo)
    return service.(*UserServiceImpl), mockUserRepo, mockAuthRepo
}
```

**验收标准**:
- [ ] 15个测试全部通过
- [ ] Mock正确设置
- [ ] 错误处理完整

---

### Phase 2.2: 认证相关 (Day 4)

**迁移方法**:
1. RegisterUser
2. LoginUser
3. LogoutUser
4. ValidateToken
5. UpdatePassword

**测试用例** (10个):

| # | 测试场景 | 测试方法 |
|---|---------|---------|
| 1 | 注册成功 | TestUserService_RegisterUser_Success |
| 2 | 用户名已存在 | TestUserService_RegisterUser_DuplicateUsername |
| 3 | 登录成功 | TestUserService_LoginUser_Success |
| 4 | 用户不存在 | TestUserService_LoginUser_UserNotFound |
| 5 | 密码错误 | TestUserService_LoginUser_WrongPassword |
| 6 | 账号未激活 | TestUserService_LoginUser_AccountInactive |
| 7 | 账号被封禁 | TestUserService_LoginUser_AccountBanned |
| 8 | 更新密码成功 | TestUserService_UpdatePassword_Success |
| 9 | 旧密码错误 | TestUserService_UpdatePassword_WrongOldPassword |
| 10 | 用户不存在 | TestUserService_UpdatePassword_UserNotFound |

**JWT相关**:
```go
// 注意: 需要JWT配置
// 跳过需要JWT的测试，在集成测试中运行
t.Skip("需要JWT配置，集成测试中运行")
```

**验收标准**:
- [ ] 10个测试全部通过（可跳过JWT相关）
- [ ] 密码加密正确
- [ ] Token生成逻辑准备

---

### Phase 2.3: 角色权限 (Day 5-6)

**迁移方法**:
1. AssignRole
2. RemoveRole
3. GetUserRoles
4. GetUserPermissions
5. DowngradeRole

**测试用例** (10个):

| # | 测试场景 | 测试方法 |
|---|---------|---------|
| 1 | 分配角色成功 | TestUserService_AssignRole_Success |
| 2 | 用户不存在 | TestUserService_AssignRole_UserNotFound |
| 3 | 角色不存在 | TestUserService_AssignRole_RoleNotFound |
| 4 | 移除角色成功 | TestUserService_RemoveRole_Success |
| 5 | 用户不存在 | TestUserService_RemoveRole_UserNotFound |
| 6 | 获取角色列表 | TestUserService_GetUserRoles_Success |
| 7 | 空角色列表 | TestUserService_GetUserRoles_Empty |
| 8 | 获取权限列表 | TestUserService_GetUserPermissions_Success |
| 9 | 角色降级成功 | TestUserService_DowngradeRole_Success |
| 10 | 降级未确认 | TestUserService_DowngradeRole_NotConfirmed |

**AuthRepository集成**:
```go
mockAuthRepo := new(MockAuthRepository)
mockAuthRepo.On("GetRole", ctx, roleID).Return(role, nil)
mockAuthRepo.On("AssignUserRole", ctx, userID, roleID).Return(nil)
```

**验收标准**:
- [ ] 10个测试全部通过
- [ ] 与AuthRepository集成正确
- [ ] 权限检查逻辑正确

---

### 阶段2验收
- [ ] UserService核心方法迁移完成
- [ ] 现有35个测试全部通过
- [ ] 新增测试通过
- [ ] 与auth模块集成验证通过
- [ ] 代码覆盖率 > 75%

---

## 🚀 阶段3: 依赖服务迁移 (Day 7-9)

### 3.1 VerificationService迁移 (Day 7-8)

**文件**: `service/user/verification_service.go`

**测试用例** (20个 - 现有已覆盖):

| # | 测试场景 | 测试方法 |
|---|---------|---------|
| 1 | 生成验证码 | TestEmailVerificationTokenManager_GenerateCode |
| 2 | 验证验证码成功 | TestEmailVerificationTokenManager_ValidateCode |
| 3 | 验证码错误 | TestEmailVerificationTokenManager_ValidateCode_Invalid |
| 4 | 标记已使用 | TestEmailVerificationTokenManager_MarkCodeAsUsed |
| 5 | 清理过期验证码 | TestEmailVerificationTokenManager_CleanExpiredCodes |
| 6 | 并发访问 | TestEmailVerificationTokenManager_ConcurrentAccess |
| 7 | 发送邮箱验证码 | TestVerificationService_SendEmailCode_Success |
| 8 | 密码重置-用户存在 | TestVerificationService_SendEmailCode_ResetPassword_UserExists |
| 9 | 密码重置-用户不存在 | TestVerificationService_SendEmailCode_ResetPassword_UserNotExists |
| 10 | 发送手机验证码 | TestVerificationService_SendPhoneCode_Success |
| 11 | 验证码验证成功 | TestVerificationService_VerifyCode_Success |
| 12 | 验证码错误 | TestVerificationService_VerifyCode_InvalidCode |
| 13 | 验证码过期 | TestVerificationService_VerifyCode_CodeExpired |
| 14 | 验证码已使用 | TestVerificationService_VerifyCode_CodeUsed |
| 15 | 标记已使用成功 | TestVerificationService_MarkCodeAsUsed_Success |
| 16 | 设置邮箱验证 | TestVerificationService_SetEmailVerified_Success |
| 17 | 邮箱不匹配 | TestVerificationService_SetEmailVerified_EmailMismatch |
| 18 | 用户不存在 | TestVerificationService_SetEmailVerified_UserNotFound |
| 19 | 设置手机验证 | TestVerificationService_SetPhoneVerified_Success |
| 20 | 检查密码 | TestVerificationService_CheckPassword_Success |

**EmailService集成**:
```go
// 需要与channels模块的EmailService集成
// 验证邮件发送功能
```

**验收标准**:
- [ ] 20个测试全部通过
- [ ] 并发安全验证
- [ ] 与EmailService集成通过

---

### 3.2 PasswordService迁移 (Day 9)

**文件**: `service/user/password_service.go`

**测试用例** (10个 - 新增5个):

| # | 测试场景 | 测试方法 |
|---|---------|---------|
| 1 | 创建服务成功 | TestNewPasswordService |
| 2 | 发送重置码 | TestPasswordService_SendResetCode |
| 3 | 邮箱不存在 | TestPasswordService_SendResetCode_EmailNotExists |
| 4 | 重置密码成功 | TestPasswordService_ResetPassword_Success |
| 5 | 更新密码成功 | TestPasswordService_UpdatePassword_Success |
| 6 | 旧密码错误 | TestPasswordService_UpdatePassword_WrongOldPassword |
| 7 | 用户不存在 | TestPasswordService_UpdatePassword_UserNotFound |
| 8 | 根据邮箱获取用户 | TestPasswordService_GetUserByEmail |
| 9 | 邮箱不存在 | TestPasswordService_GetUserByEmail_NotExists |
| 10 | 根据ID获取用户 | TestPasswordService_GetUserByID |

**验收标准**:
- [ ] 10个测试全部通过
- [ ] 与VerificationService集成正确
- [ ] 密码重置流程完整

---

### 阶段3验收
- [ ] VerificationService迁移完成
- [ ] PasswordService迁移完成
- [ ] 与channels模块集成通过
- [ ] 验证码端到端测试通过
- [ ] 代码覆盖率 > 70%

---

## 🚀 阶段4: 高级特性迁移 (Day 10-12)

### 4.1 TransactionManager迁移 (Day 10-11)

**文件**: `service/user/transaction_manager.go`

**测试用例** (15个):

| # | 测试场景 | 测试方法 |
|---|---------|---------|
| 1 | 创建事务管理器 | TestNewTransactionManager |
| 2 | 用户注册事务成功 | TestUserRegistrationTransaction_Execute_Success |
| 3 | 创建用户失败 | TestUserRegistrationTransaction_Execute_UserCreateFailed |
| 4 | 分配角色失败 | TestUserRegistrationTransaction_Execute_AssignRoleFailed |
| 5 | 初始化配置失败 | TestUserRegistrationTransaction_Execute_InitConfigFailed |
| 6 | 用户软删除成功 | TestUserDeletionTransaction_Execute_SoftDelete_Success |
| 7 | 用户硬删除成功 | TestUserDeletionTransaction_Execute_HardDelete_Success |
| 8 | 软删除归档项目 | TestUserDeletionTransaction_executeSoftDelete_ArchiveProjects |
| 9 | 软删除禁用会话 | TestUserDeletionTransaction_executeSoftDelete_DisableSessions |
| 10 | 硬删除删除角色 | TestUserDeletionTransaction_executeHardDelete_DeleteRoles |
| 11 | 硬删除删除配置 | TestUserDeletionTransaction_executeHardDelete_DeleteConfigs |
| 12 | 硬删除处理项目 | TestUserDeletionTransaction_executeHardDelete_HandleProjects |
| 13 | Saga执行成功 | TestSagaManager_ExecuteSaga_Success |
| 14 | Saga失败补偿 | TestSagaManager_ExecuteSaga_Failure_Compensate |
| 15 | 引用完整性验证 | TestReferenceIntegrityManager_ValidateReferences |

**MongoDB事务测试**:
```go
// 需要MongoDB测试环境
// 使用testcontainers或本地MongoDB
```

**验收标准**:
- [ ] 15个测试全部通过
- [ ] 事务回滚正确
- [ ] Saga补偿机制工作
- [ ] 引用完整性保证

---

### 4.2 Token管理器迁移 (Day 12)

**文件**: 
- `service/user/email_verification_token.go`
- `service/user/password_reset_token.go`

**测试用例** (已在VerificationService中覆盖)

**验收标准**:
- [ ] Token生成正确
- [ ] Token验证正确
- [ ] Token过期机制工作
- [ ] 并发安全保证

---

### 阶段4验收
- [ ] 事务管理迁移完成
- [ ] Token管理迁移完成
- [ ] 并发安全测试通过
- [ ] 性能基准测试通过
- [ ] 代码覆盖率 > 65% (事务测试较难)

---

## 📊 测试覆盖率目标

| 组件 | 目标覆盖率 | 优先级 |
|------|-----------|--------|
| UserValidator | > 90% | 高 |
| PasswordValidator | > 90% | 高 |
| Converter | > 85% | 高 |
| UserService | > 75% | 高 |
| VerificationService | > 80% | 中 |
| PasswordService | > 75% | 中 |
| TransactionManager | > 65% | 低 |
| **总体** | **> 75%** | - |

---

## 🔧 测试工具栈

```go
// 测试框架
import (
    "testing"                    // Go标准测试
    "github.com/stretchr/testify" // 断言和Mock
)

// Mock工具
type MockUserRepository struct {
    mock.Mock
}

// 覆盖率工具
// go test -cover
// go test -coverprofile=coverage.out
// go tool cover -html=coverage.out

// 基准测试
// go test -bench=.
// go test -bench=. -benchmem
```

---

## ✅ 总体验收标准

### 功能验收
- [ ] 所有接口方法实现完成
- [ ] 所有测试用例通过
- [ ] 与auth模块集成验证通过
- [ ] 与channels模块集成验证通过

### 质量验收
- [ ] 代码覆盖率 > 75%
- [ ] 无编译错误
- [ ] 无静态检查错误 (go vet, golangci-lint)
- [ ] 代码审查通过

### 文档验收
- [ ] API文档更新
- [ ] 架构文档更新
- [ ] 测试文档完整

### Git验收
- [ ] 提交信息规范
- [ ] 无敏感信息
- [ ] 分支管理正确

---

## 📅 时间表

```
Week 1 (Day 1-5)
├── Day 1-2: 阶段1 - 基础组件
├── Day 3:   阶段2.1 - 基础CRUD
├── Day 4:   阶段2.2 - 认证相关
└── Day 5:   阶段2.3 - 角色权限(开始)

Week 2 (Day 6-12)
├── Day 6:   阶段2.3 - 角色权限(完成)
├── Day 7-8: 阶段3.1 - VerificationService
├── Day 9:   阶段3.2 - PasswordService
├── Day 10-11: 阶段4.1 - TransactionManager
└── Day 12:  阶段4.2 - Token管理 + 验收

缓冲时间: Day 13-16 (应对意外情况)
```

---

## 🎯 成功标准

✅ **功能完整**: 所有35+个接口方法实现并测试通过  
✅ **质量保证**: 代码覆盖率 > 75%，无静态检查错误  
✅ **集成验证**: 与auth和channels模块集成测试通过  
✅ **文档齐全**: 代码文档、测试文档、架构文档完整  
✅ **代码审查**: 通过Design-Review-Maid和Code-Review-Maid审查  

---

**计划制定完成时间**: 2026-02-10  
**计划版本**: v1.0  
**状态**: ✅ 就绪，等待主人确认执行
