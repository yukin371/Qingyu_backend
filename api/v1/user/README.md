# User API 模块结构说明

## 文件结构

```
api/v1/user/
├── password_api.go                  # 密码管理API（重置、修改）
├── security_api.go                  # 用户安全API（邮箱验证、密码重置流程）
├── verification_api.go              # 验证码API（邮箱/手机验证、设备管理）
├── user_api_test_disabled.go        # 用户API集成测试（已禁用）
├── password_api_test.go             # 密码API单元测试
├── verification_api_test.go         # 验证码API单元测试
├── dto/
│   ├── auth_dto.go                  # 认证相关 DTO（注册/登录请求响应）
│   ├── password_dto.go              # 密码相关 DTO（重置/修改请求响应）
│   ├── user_dto.go                  # 用户相关 DTO（资料更新/头像上传/公开信息）
│   └── verification_dto.go          # 验证相关 DTO（验证码/解绑/设备请求响应）
├── handler/
│   ├── auth_handler.go              # 认证处理器（注册/登录/登出）
│   ├── auth_handler_logout_test.go  # 登出处理器测试
│   ├── profile_handler.go           # 个人信息处理器（资料CRUD/头像上传/角色降级）
│   ├── profile_handler_test.go      # 个人信息处理器测试
│   ├── public_user_handler.go       # 公开用户信息处理器（用户主页/作品列表/批量查询）
│   ├── stats_handler.go             # 用户统计处理器（统计/内容统计/活跃度/收益）
│   └── stats_handler_test.go        # 统计处理器测试
└── README.md                        # 本文件
```

## 模块职责划分

### 1. PasswordAPI (`password_api.go`)

**职责**: 密码管理功能

**核心功能**:
- ✅ 发送密码重置验证码（向邮箱发送6位验证码，有效期5分钟）
- ✅ 重置密码（通过邮箱+验证码重置密码）
- ✅ 修改密码（已登录用户验证旧密码后修改）

**API端点**:
```
POST /api/v1/users/password/reset/send    # 发送密码重置验证码
POST /api/v1/users/password/reset/verify   # 验证重置码并重置密码
PUT  /api/v1/users/password                # 修改密码（需JWT认证）
```

**依赖服务**:
- `user.PasswordService` - 密码管理服务

---

### 2. SecurityAPI (`security_api.go`)

**职责**: 用户安全相关操作（邮箱验证、密码重置流程）

**核心功能**:
- ✅ 发送邮箱验证码
- ✅ 验证邮箱（使用6位数字验证码）
- ✅ 请求密码重置（发送含Token的重置链接）
- ✅ 确认密码重置（使用Token完成重置）

**API端点**:
```
POST /api/v1/user/email/send-code          # 发送邮箱验证码
POST /api/v1/user/email/verify             # 验证邮箱
POST /api/v1/user/password/reset-request   # 请求密码重置（发送重置链接）
POST /api/v1/user/password/reset           # 确认密码重置（Token验证）
```

**依赖服务**:
- `user.UserService` - 用户服务接口

**备注**: 此API通过 `user.UserService` 接口（基于service/interfaces/user）实现，与PasswordAPI使用不同的服务层接口。密码重置流程与PasswordAPI存在功能交叉，但实现路径不同——SecurityAPI走接口化服务层，PasswordAPI直接依赖具体服务实现。

---

### 3. VerificationAPI (`verification_api.go`)

**职责**: 验证码发送/验证、账号绑定管理、设备管理

**核心功能**:
- ✅ 发送邮箱验证码（6位数字验证码，有效期30分钟）
- ✅ 发送手机验证码（模拟实现，控制台打印）
- ✅ 验证邮箱（验证码校验 + 标记已验证 + 防重复使用）
- ✅ 解绑邮箱（需验证密码）
- ✅ 解绑手机（模拟实现）
- ✅ 删除设备（需验证密码）

**API端点**:
```
POST   /api/v1/users/verify/email/send       # 发送邮箱验证码
POST   /api/v1/users/verify/phone/send       # 发送手机验证码（模拟）
POST   /api/v1/users/email/verify            # 验证邮箱
DELETE /api/v1/users/email/unbind            # 解绑邮箱
DELETE /api/v1/users/phone/unbind            # 解绑手机
DELETE /api/v1/users/devices/{deviceId}      # 删除设备
```

**依赖服务**:
- `user.VerificationService` - 验证码服务（发送、校验、标记使用、密码验证）
- `user.UserService` - 用户服务接口（解绑邮箱/手机、删除设备）

---

### 4. AuthHandler (`handler/auth_handler.go`)

**职责**: 用户认证（注册、登录、登出）

**核心功能**:
- ✅ 用户注册（用户名+邮箱+密码，注册后返回Token）
- ✅ 用户登录（用户名+密码，记录客户端IP，返回Token和角色信息）
- ✅ 用户登出（幂等操作，使Token失效）

**API端点**:
```
POST /api/v1/user/auth/register   # 用户注册
POST /api/v1/user/auth/login      # 用户登录
POST /api/v1/user/auth/logout     # 用户登出（需JWT认证）
```

**依赖服务**:
- `auth.AuthService`（通过 `authHandlerService` 接口解耦） - 认证服务

**与 `api/v1/shared/auth_api.go` 的关系**:

两个模块均提供注册/登录/登出功能，存在路由重叠：

| 功能 | `user/handler/AuthHandler` | `shared/AuthAPI` |
|------|---------------------------|-------------------|
| 注册 | `/api/v1/user/auth/register` | `/api/v1/shared/auth/register` |
| 登录 | `/api/v1/user/auth/login` | `/api/v1/shared/auth/login` |
| 登出 | `/api/v1/user/auth/logout` | `/api/v1/shared/auth/logout` |
| 刷新Token | - | `/api/v1/shared/auth/refresh` |
| 权限/角色 | - | `/api/v1/shared/auth/permissions`, `/api/v1/shared/auth/roles` |
| 验证码 | - | `/api/v1/shared/auth/send-verification-code` |

- `shared/AuthAPI` 是旧版认证入口，注册流程包含邮箱验证码校验（`VerificationCode` 字段），还提供Token刷新、权限/角色查询等额外能力
- `user/handler/AuthHandler` 是新版认证入口，注册流程更简洁（无验证码步骤），返回结构包含完整角色列表（`Roles`字段）
- 两者底层共享 `auth.AuthService`，路由注册由 Router 层决定实际生效路径

---

### 5. ProfileHandler (`handler/profile_handler.go`)

**职责**: 已登录用户的个人信息管理

**核心功能**:
- ✅ 获取当前用户信息
- ✅ 更新个人信息（昵称、简介、头像URL、手机、性别、生日、地区、网站）
- ✅ 修改密码（验证旧密码）
- ✅ 上传头像（支持JPG/PNG/GIF，最大5MB，通过StorageService存储）
- ✅ 降级用户角色（admin->author/reader, author->reader，需确认）

**API端点**:
```
GET  /api/v1/user/profile            # 获取当前用户信息（需JWT认证）
PUT  /api/v1/user/profile            # 更新个人信息（需JWT认证）
PUT  /api/v1/user/password           # 修改密码（需JWT认证）
POST /api/v1/user/avatar             # 上传头像（需JWT认证，multipart/form-data）
POST /api/v1/user/role/downgrade     # 降级角色（需JWT认证）
```

**依赖服务**:
- `user.UserService` - 用户服务接口
- `storage.StorageService` - 文件存储服务（可选依赖，用于头像上传）

**备注**: UpdatePassword 与 PasswordAPI.UpdatePassword 功能重叠，两者路径不同（`/api/v1/user/password` vs `/api/v1/users/password`），由 Router 层决定实际路由映射。

---

### 6. PublicUserHandler (`handler/public_user_handler.go`)

**职责**: 公开用户信息查询（无需认证）

**核心功能**:
- ✅ 获取用户公开信息（不含敏感字段）
- ✅ 获取用户详细资料（公开主页展示）
- ✅ 获取用户作品列表（分页查询，依赖BookstoreService）
- ✅ 批量获取用户信息（逗号分隔ID列表，上限50个）

**API端点**:
```
GET /api/v1/user/users/{id}            # 获取用户公开信息
GET /api/v1/user/users/{id}/profile    # 获取用户详细资料
GET /api/v1/user/users/{id}/books      # 获取用户作品列表
GET /api/v1/user/users/batch           # 批量获取用户信息（?ids=id1,id2,...）
```

**依赖服务**:
- `user.UserService` - 用户服务接口
- `BookstoreService`（可选依赖，通过 `SetBookstoreService` 注入） - 书店服务（用于查询用户作品）

---

### 7. StatsHandler (`handler/stats_handler.go`)

**职责**: 用户统计数据查询

**核心功能**:
- ✅ 获取我的统计（基础用户统计）
- ✅ 获取我的内容统计（创作内容相关统计）
- 开发中：活跃度统计（7天/自定义天数）
- 开发中：收益统计（按日期范围查询）

**API端点**:
```
GET /api/v1/user/stats/my           # 获取我的统计（需JWT认证）
GET /api/v1/user/stats/my/content   # 获取内容统计（需JWT认证）
GET /api/v1/user/stats/my/activity  # 获取活跃度统计（开发中，需JWT认证）
GET /api/v1/user/stats/my/revenue   # 获取收益统计（开发中，需JWT认证）
```

**依赖服务**:
- `UserStatsService` - 用户统计服务接口（切片B seam）
- `ContentStatsService` - 内容统计服务接口（切片C seam）

**备注**: 活跃度统计和收益统计尚未实现，调用时返回 501 Not Implemented。

---

## DTO 层说明

| 文件 | 包含结构体 | 说明 |
|------|-----------|------|
| `dto/auth_dto.go` | `RegisterRequest`, `RegisterResponse`, `LoginRequest`, `LoginResponse` | 认证相关请求/响应 |
| `dto/password_dto.go` | `SendPasswordResetRequest`, `ResetPasswordRequest`, `UpdatePasswordRequest` 等 | 密码管理相关请求/响应 |
| `dto/user_dto.go` | `UpdateProfileRequest`, `GetUserBooksResponse`, `UploadAvatarResponse` 等 | 用户资料相关请求/响应，复用 `shared` 层公共类型 |
| `dto/verification_dto.go` | `SendEmailCodeRequest`, `VerifyEmailRequest`, `UnbindEmailRequest` 等 | 验证码与绑定管理相关请求/响应 |

---

## 设计原则

### 1. 接口解耦
Handler 通过接口（`authHandlerService`、`UserStatsService`、`ContentStatsService`、`BookstoreService`）与具体服务实现解耦，便于单元测试和依赖管理。

### 2. 可选依赖注入
`ProfileHandler.storageService` 和 `PublicUserHandler.bookstoreService` 通过 Setter 方法注入，未注入时仍可正常工作（头像上传和作品查询功能降级处理）。

### 3. 统一响应格式
所有 API 使用 `response.Success`、`response.BadRequest`、`response.NotFound`、`response.Unauthorized`、`response.InternalError` 统一响应格式。

### 4. 错误类型路由
Handler 层根据 `serviceInterfaces.ServiceError.Type`（NotFound/Unauthorized/Validation/Business）返回对应 HTTP 状态码。

---

## 与 `shared/auth` 的兼容关系

`api/v1/user/` 和 `api/v1/shared/auth_api.go` 在注册/登录/登出功能上存在路由重叠。当前两套路由并存：

- **`shared/auth` (旧版)**: 路径前缀 `/api/v1/shared/auth/`，注册需要邮箱验证码，额外提供Token刷新、权限/角色查询
- **`user/handler` (新版)**: 路径前缀 `/api/v1/user/auth/`，注册更简洁，返回结构更丰富

底层共享 `auth.AuthService` 服务，Router 层决定实际注册的路由。建议后续统一到 `user/handler` 版本，将 `shared/auth` 的Token刷新和权限查询能力迁移过来后移除旧路由。

---

**版本**: v1.0
**更新日期**: 2026-05-21
**维护者**: yukin371
