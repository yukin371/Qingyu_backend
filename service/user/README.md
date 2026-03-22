# User Service 模块架构文档

用户服务模块，提供用户注册、登录、密码管理、邮箱验证等核心用户业务逻辑。

## 架构图

```mermaid
graph TB
    subgraph "User Service Layer"
        UserService["UserServiceImpl<br/>用户核心服务"]
        PasswordService["PasswordService<br/>密码服务"]
        VerificationService["VerificationService<br/>验证服务"]
    end

    subgraph "Support Components"
        PasswordValidator["PasswordValidator<br/>密码强度验证器"]
        EmailTokenManager["EmailVerificationTokenManager<br/>邮箱验证码管理器"]
        PasswordResetManager["PasswordResetTokenManager<br/>密码重置Token管理器"]
        TransactionManager["TransactionManager<br/>事务管理器"]
        CascadeManager["CascadeManager<br/>级联操作管理器"]
        Converter["Converter<br/>DTO转换器"]
        UserValidator["UserValidator<br/>用户验证器"]
    end

    subgraph "External Dependencies"
        UserRepository["UserRepository<br/>用户仓储接口"]
        AuthRepository["AuthRepository<br/>认证仓储接口"]
        EmailService["EmailService<br/>邮件服务"]
    end

    %% UserService 依赖
    UserService --> UserRepository
    UserService --> AuthRepository
    UserService --> Converter
    UserService --> EmailTokenManager
    UserService --> PasswordResetManager

    %% PasswordService 依赖
    PasswordService --> VerificationService
    PasswordService --> UserRepository
    PasswordService --> PasswordValidator

    %% VerificationService 依赖
    VerificationService --> UserRepository
    VerificationService --> AuthRepository
    VerificationService --> EmailService
    VerificationService --> EmailTokenManager

    %% 事务管理
    TransactionManager --> CascadeManager
```

## 核心服务列表

### UserServiceImpl - 用户核心服务

主服务实现，提供完整的用户管理功能。

| 方法 | 职责 |
|------|------|
| `CreateUser` | 创建用户（管理端） |
| `RegisterUser` | 用户注册（公开） |
| `LoginUser` | 用户登录认证 |
| `LogoutUser` | 用户登出 |
| `GetUser` | 获取用户信息 |
| `UpdateUser` | 更新用户信息 |
| `DeleteUser` | 删除用户 |
| `ListUsers` | 获取用户列表 |
| `UpdatePassword` | 修改密码（需旧密码） |
| `ResetPassword` | 重置密码（发送重置邮件） |
| `ConfirmPasswordReset` | 确认密码重置 |
| `SendEmailVerification` | 发送邮箱验证码 |
| `VerifyEmail` | 验证邮箱 |
| `AssignRole` | 分配角色 |
| `RemoveRole` | 移除角色 |
| `GetUserRoles` | 获取用户角色 |
| `GetUserPermissions` | 获取用户权限 |
| `DowngradeRole` | 角色降级 |

### PasswordService - 密码服务

密码管理专用服务。

| 方法 | 职责 |
|------|------|
| `SendResetCode` | 发送密码重置验证码 |
| `ResetPassword` | 通过验证码重置密码 |
| `UpdatePassword` | 修改密码（需验证旧密码） |

### VerificationService - 验证服务

验证码发送和校验服务。

| 方法 | 职责 |
|------|------|
| `SendEmailCode` | 发送邮箱验证码 |
| `SendPhoneCode` | 发送手机验证码 |
| `VerifyCode` | 验证验证码 |
| `MarkCodeAsUsed` | 标记验证码已使用 |
| `SetEmailVerified` | 设置邮箱已验证 |
| `SetPhoneVerified` | 设置手机已验证 |
| `CheckPassword` | 验证密码正确性 |
| `EmailExists` | 检查邮箱是否存在 |
| `PhoneExists` | 检查手机是否存在 |

## 辅助组件

### PasswordValidator - 密码验证器

密码强度验证和评分。

- 最小长度: 8位
- 必须包含: 大写字母、小写字母、数字
- 可选: 特殊字符
- 检测: 常见弱密码、连续字符

### EmailVerificationTokenManager - 邮箱验证码管理器

- 生成6位数字验证码
- 有效期: 30分钟
- 单例模式，自动清理过期Token

### PasswordResetTokenManager - 密码重置Token管理器

- 生成64字符随机Token
- 有效期: 1小时
- 支持一次性使用标记

### TransactionManager - 事务管理器

支持复杂业务场景的事务操作:

- `UserRegistrationTransaction`: 用户注册事务（用户+角色+配置）
- `UserDeletionTransaction`: 用户删除事务（软删除/硬删除）
- `SagaManager`: Saga模式分布式事务

### Converter - DTO转换器

Model 与 DTO 之间的转换:

- `ToUserDTO`: User Model -> UserDTO
- `ToUserDTOs`: 批量转换
- `ToUser`: DTO -> Model（用于更新）
- `ToUserWithoutID`: DTO -> Model（用于创建）

## 依赖关系

```mermaid
graph LR
    subgraph "Service Layer"
        UserService
        PasswordService
        VerificationService
    end

    subgraph "Repository Layer"
        UserRepository["UserRepository"]
        AuthRepository["AuthRepository"]
    end

    subgraph "External Services"
        EmailService["EmailService"]
    end

    UserService --> UserRepository
    UserService --> AuthRepository
    PasswordService --> UserRepository
    VerificationService --> UserRepository
    VerificationService --> EmailService
```

### Repository 依赖

| 服务 | Repository | 用途 |
|------|------------|------|
| UserServiceImpl | `UserRepository` | 用户CRUD、状态管理 |
| UserServiceImpl | `AuthRepository` | 角色、权限管理 |
| PasswordService | `UserRepository` | 密码更新、用户查询 |
| VerificationService | `UserRepository` | 用户信息查询、验证状态更新 |
| VerificationService | `AuthRepository` | 认证相关操作 |

## 核心流程说明

### 用户注册流程

```mermaid
sequenceDiagram
    participant Client
    participant UserService
    participant UserRepository
    participant TokenManager

    Client->>UserService: RegisterUser(username, email, password)
    UserService->>UserService: 验证请求参数
    UserService->>UserRepository: ExistsByUsername()
    UserService->>UserRepository: ExistsByEmail()
    UserService->>UserService: 创建User对象
    UserService->>UserService: SetPassword() (bcrypt加密)
    UserService->>UserRepository: Create(user)

    alt 创建成功
        UserService->>TokenManager: GenerateToken(userID, roles)
        TokenManager-->>UserService: JWT Token
        UserService-->>Client: RegisterUserResponse(user, token)
    else 并发冲突
        UserService->>UserService: 重试机制(最多3次)
    end
```

### 用户登录流程

```mermaid
sequenceDiagram
    participant Client
    participant UserService
    participant UserRepository
    participant TokenManager

    Client->>UserService: LoginUser(username, password)
    UserService->>UserService: 验证参数非空
    UserService->>UserRepository: GetByUsername(username)

    alt 用户不存在
        UserRepository-->>UserService: NotFoundError
        UserService-->>Client: 用户不存在
    else 用户存在
        UserRepository-->>UserService: User
        UserService->>UserService: ValidatePassword()

        alt 密码错误
            UserService-->>Client: 密码错误
        else 密码正确
            UserService->>UserService: 检查用户状态

            alt 状态异常
                UserService-->>Client: 账号未激活/已封禁/已删除
            else 状态正常
                UserService->>UserRepository: UpdateLastLogin()
                UserService->>TokenManager: GenerateToken()
                UserService-->>Client: LoginUserResponse(user, token)
            end
        end
    end
```

### 密码重置流程

```mermaid
sequenceDiagram
    participant Client
    participant PasswordService
    participant VerificationService
    participant UserRepository
    participant EmailService

    Client->>PasswordService: SendResetCode(email)
    PasswordService->>UserRepository: GetByEmail(email)

    alt 用户存在
        PasswordService->>VerificationService: SendEmailCode(email, "reset_password")
        VerificationService->>EmailService: 发送验证码邮件
        PasswordService-->>Client: 发送成功
    else 用户不存在
        PasswordService-->>Client: 邮箱不存在
    end

    Client->>PasswordService: ResetPassword(email, code, newPassword)
    PasswordService->>VerificationService: VerifyCode(email, code, "reset_password")

    alt 验证码有效
        PasswordService->>VerificationService: MarkCodeAsUsed(email)
        PasswordService->>PasswordService: bcrypt加密新密码
        PasswordService->>UserRepository: UpdatePasswordByEmail()
        PasswordService-->>Client: 重置成功
    else 验证码无效
        PasswordService-->>Client: 验证码无效或已过期
    end
```

### 邮箱验证流程

```mermaid
sequenceDiagram
    participant Client
    participant UserService
    participant UserRepository
    participant TokenManager
    participant EmailService

    Client->>UserService: SendEmailVerification(userID, email)
    UserService->>UserRepository: GetByID(userID)

    alt 用户存在且邮箱匹配
        UserService->>TokenManager: GenerateCode(userID, email)
        TokenManager-->>UserService: 6位验证码
        UserService->>EmailService: 发送验证邮件
        UserService-->>Client: 发送成功
    else 邮箱不匹配或已验证
        UserService-->>Client: 相应错误
    end

    Client->>UserService: VerifyEmail(userID, code)
    UserService->>UserRepository: GetByID(userID)
    UserService->>TokenManager: ValidateCode(userID, email, code)

    alt 验证码有效
        UserService->>TokenManager: MarkCodeAsUsed(email)
        UserService->>UserRepository: Update(userID, {email_verified: true, status: active})
        UserService-->>Client: 验证成功
    else 验证码无效
        UserService-->>Client: 验证码无效或已过期
    end
```

## 错误处理

模块使用统一的错误码体系:

| 错误码 | 说明 | HTTP状态码 |
|--------|------|------------|
| 40401 | 用户不存在 | 404 |
| 40001 | 邮箱格式无效 | 400 |
| 40002 | 密码格式无效 | 400 |
| 40901 | 用户已存在 | 409 |
| 40101 | 令牌无效 | 401 |
| 40102 | 令牌过期 | 401 |
| 40301 | 权限不足 | 403 |
| 50001 | 内部错误 | 500 |

## 文件结构

```
service/user/
├── user_service.go              # 用户核心服务实现
├── password_service.go          # 密码服务
├── verification_service.go      # 验证服务
├── password_validator.go        # 密码强度验证器
├── email_verification_token.go  # 邮箱验证码管理器
├── password_reset_token.go      # 密码重置Token管理器
├── transaction_manager.go       # 事务管理器
├── converter.go                 # DTO转换器
├── user_validator.go            # 用户验证器
├── errors.go                    # 错误定义
├── constants.go                 # 常量定义
└── mocks/                       # Mock文件（测试用）
    ├── mock_user_repository.go
    └── mock_auth_repository.go
```
