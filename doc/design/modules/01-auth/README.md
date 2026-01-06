# 01 - 认证授权模块

> **模块编号**: 01
> **模块名称**: Authentication & Authorization
> **负责功能**: 用户身份验证、授权和账户管理
> **完成度**: 🟢 95%

## 📋 目录结构

```
认证授权模块/
├── api/v1/
│   ├── user/                      # 用户API
│   │   ├── user_api.go           # 用户基本操作
│   │   ├── profile_api.go        # 用户资料
│   │   └── password_api.go       # 密码管理
│   └── shared/                    # 共享认证API
│       └── auth_api.go           # 登录/注册/令牌刷新
├── service/user/                   # 用户服务层
│   ├── user_service.go           # 用户业务逻辑
│   └── password_service.go       # 密码服务
├── repository/interfaces/user/    # 用户仓储接口
├── repository/mongodb/user/        # MongoDB仓储实现
│   ├── user_repository_mongo.go  # 用户CRUD
│   └── password_reset_repository_mongo.go  # 密码重置
└── models/users/                   # 数据模型
    ├── user.go                   # 用户实体
    ├── user_stats.go             # 用户统计
    └── password_reset.go         # 密码重置
```

## 🎯 核心功能

### 1. 用户认证

- **注册**: 邮箱/手机注册，验证码验证
- **登录**: 用户名/邮箱 + 密码登录
- **登出**: 清除会话，加入令牌黑名单
- **令牌刷新**: 访问令牌过期时使用刷新令牌获取新令牌

### 2. 用户管理

- **资料管理**: 昵称、头像、简介
- **密码管理**: 修改密码、重置密码
- **状态管理**: 活跃/禁用/删除状态
- **邮箱验证**: 邮箱验证流程

### 3. 权限控制

- **角色管理**: user/author/editor/admin/superadmin
- **权限检查**: 基于角色的访问控制(RBAC)
- **API权限**: 中间件级别的权限验证

## 🔑 安全特性

### 密码安全

- 使用 bcrypt 哈希算法（cost factor: 10）
- 密码强度验证（最小长度、复杂度要求）
- 密码重置令牌有效期控制

### 令牌安全

- JWT 签名密钥存储在环境变量
- 访问令牌有效期：2小时
- 刷新令牌有效期：7天
- 令牌黑名单机制（Redis存储）

### 账户安全

- 登录失败次数限制
- 邮箱/手机验证要求
- 可疑登录检测
- 账户封禁/删除状态检查

## 📊 数据模型

### User (用户实体)

```go
type User struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username        string             `bson:"username" json:"username"`
    Email           string             `bson:"email" json:"email"`
    Phone           string             `bson:"phone,omitempty" json:"phone,omitempty"`
    PasswordHash    string             `bson:"password_hash" json:"-"`
    Nickname        string             `bson:"nickname" json:"nickname"`
    Avatar          string             `bson:"avatar" json:"avatar"`
    Bio             string             `bson:"bio" json:"bio"`
    Status          UserStatus         `bson:"status" json:"status"`
    Roles           []string           `bson:"roles" json:"roles"`
    EmailVerified   bool               `bson:"email_verified" json:"emailVerified"`
    PhoneVerified   bool               `bson:"phone_verified" json:"phoneVerified"`
    CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
    LastLoginAt     *time.Time         `bson:"last_login_at,omitempty" json:"lastLoginAt,omitempty"`
}
```

### PasswordReset (密码重置)

```go
type PasswordReset struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email     string             `bson:"email" json:"email"`
    Token     string             `bson:"token" json:"token"`
    ExpiresAt time.Time          `bson:"expires_at" json:"expiresAt"`
    Used      bool               `bson:"used" json:"used"`
    CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}
```

## 🔄 核心流程

### 注册流程

```
1. 用户提交注册信息（用户名/邮箱/密码）
   ↓
2. 后端验证信息格式和唯一性
   ↓
3. 创建用户记录（状态：inactive，未验证）
   ↓
4. 发送验证邮件/短信
   ↓
5. 用户点击验证链接或输入验证码
   ↓
6. 激活账户（状态：active）
   ↓
7. 返回JWT令牌
```

### 登录流程

```
1. 用户提交用户名/邮箱和密码
   ↓
2. 后端查询用户记录
   ↓
3. 验证密码（bcrypt比较）
   ↓
4. 检查账户状态（是否被封禁/删除）
   ↓
5. 生成JWT访问令牌和刷新令牌
   ↓
6. 更新最后登录时间
   ↓
7. 返回令牌和用户基本信息
```

### 令牌刷新流程

```
1. 客户端提交刷新令牌
   ↓
2. 后端验证刷新令牌有效性
   ↓
3. 检查令牌是否在黑名单中
   ↓
4. 生成新的访问令牌和刷新令牌
   ↓
5. 将旧刷新令牌加入黑名单
   ↓
6. 返回新令牌
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/shared/auth/register | 用户注册 | 否 |
| POST | /api/v1/shared/auth/login | 用户登录 | 否 |
| POST | /api/v1/shared/auth/logout | 用户登出 | 是 |
| POST | /api/v1/shared/auth/refresh | 刷新令牌 | 否 |
| GET | /api/v1/shared/auth/permissions | 获取权限 | 是 |
| GET | /api/v1/shared/auth/roles | 获取角色 | 是 |
| GET | /api/v1/users/profile | 获取个人资料 | 是 |
| PUT | /api/v1/users/profile | 更新个人资料 | 是 |
| PUT | /api/v1/users/password | 修改密码 | 是 |
| GET | /api/v1/users/:userId/books | 获取用户作品 | 是 |
| POST | /api/v1/users/password/reset | 请求密码重置 | 否 |
| POST | /api/v1/users/password/reset/confirm | 确认密码重置 | 否 |

## 🔧 依赖关系

### 依赖的模块
- 无（认证模块是基础模块）

### 被依赖的模块
- 所有其他模块（写作、阅读、社交等）都依赖认证模块进行用户身份验证

### 外部服务
- **邮件服务**: 用于发送验证邮件和密码重置邮件
- **短信服务**: 用于发送手机验证码
- **Redis**: 存储令牌黑名单和验证码

## ⚙️ 配置项

```yaml
auth:
  jwt:
    secret: ${JWT_SECRET}
    access_token_duration: 2h
    refresh_token_duration: 168h  # 7天
  bcrypt:
    cost: 10
  email:
    from: noreply@qingyu.com
    verification_url: https://qingyu.com/verify
    reset_url: https://qingyu.com/reset-password
  password:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_digit: true
    require_special_char: true
```

## 📈 扩展点

1. **第三方登录集成**
   - 可添加 OAuth2.0 支持（微信、QQ、GitHub等）
   - 统一身份认证（SSO）

2. **多因素认证（MFA）**
   - TOTP（基于时间的一次性密码）
   - 短信验证码二次确认

3. **单点登录（SSO）**
   - 支持多个子系统的统一登录
   - CAS 或 OAuth2.0 协议支持

4. **审计日志**
   - 记录所有认证和授权操作
   - 异常登录检测和告警

## 🚀 性能优化

1. **缓存策略**
   - Redis 缓存用户基本信息
   - 令牌黑名单使用 Redis 存储

2. **数据库优化**
   - 用户名、邮箱、手机号建立唯一索引
   - 复合索引优化查询性能

3. **异步处理**
   - 邮件发送异步处理
   - 短信发送异步处理

## 📊 监控指标

- 注册成功率
- 登录成功率
- 令牌刷新频率
- 密码重置请求量
- 活跃用户数
- 封禁用户数

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/user/`, `../../Qingyu_backend/api/v1/shared/`
