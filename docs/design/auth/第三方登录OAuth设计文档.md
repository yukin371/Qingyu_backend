# 第三方登录 OAuth2.0 设计文档

**版本**: v1.0
**创建日期**: 2025-01-07
**状态**: 设计完成

## 概述

本文档描述了青羽写作平台的第三方登录功能设计，支持用户通过 Google、GitHub、QQ 等第三方平台进行登录和账号绑定。

## 目标

1. **简化注册流程**：用户可以通过第三方账号快速注册和登录
2. **提升用户体验**：减少用户记忆密码的负担
3. **账号关联**：支持将多个第三方账号绑定到同一个平台账号
4. **安全性**：遵循 OAuth2.0 标准，确保用户信息安全

## 技术选型

### OAuth 2.0 协议

使用标准 OAuth 2.0 授权码流程 (Authorization Code Flow)：

```
用户 → 第三方平台 → 授权码 → 后端 → Access Token → 用户信息
```

**优势**：
- 行业标准，安全性高
- 支持主流平台
- 用户授权过程透明

### 支持的提供商

| 提供商 | 状态 | 说明 |
|--------|------|------|
| Google | ✅ | 使用 golang.org/x/oauth2/google |
| GitHub | ✅ | 使用 golang.org/x/oauth2/github |
| QQ | ✅ | 使用自定义端点 |
| 微信 | 🚧 | 预留接口，待实现 |
| 微博 | 🚧 | 预留接口，待实现 |

## 系统架构

### 整体架构图

```
┌──────────────────────────────────────────────────────────────┐
│                         前端应用                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ 登录按钮     │  │ 授权页面     │  │ 账号管理     │        │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘        │
└─────────┼──────────────────┼──────────────────┼───────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌──────────────────────────────────────────────────────────────┐
│                      API层 (Gin)                              │
│  ┌────────────────────────────────────────────────────┐      │
│  │ POST /api/v1/shared/oauth/{provider}/authorize     │      │
│  │ POST /api/v1/shared/oauth/{provider}/callback      │      │
│  │ GET  /api/v1/shared/oauth/accounts                 │      │
│  │ DELETE /api/v1/shared/oauth/accounts/{id}          │      │
│  └────────────────────────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                    服务层 (Service)                           │
│  ┌────────────────────────────────────────────────────┐      │
│  │           OAuthService                              │      │
│  │  - GetAuthURL()         获取授权URL                  │      │
│  │  - ExchangeCode()       交换授权码                   │      │
│  │  - GetUserInfo()        获取用户信息                 │      │
│  │  - LinkAccount()        绑定账号                     │      │
│  │  - UnlinkAccount()      解绑账号                     │      │
│  │  - RefreshToken()       刷新令牌                     │      │
│  └────────────────────────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                  仓储层 (Repository)                          │
│  ┌────────────────────────────────────────────────────┐      │
│  │           OAuthRepository                            │      │
│  │  - FindByProviderAndProviderID()                    │      │
│  │  - FindByUserID()                                   │      │
│  │  - Create() / Update() / Delete()                   │      │
│  │  - UpdateTokens()                                   │      │
│  └────────────────────────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                  数据层 (MongoDB)                             │
│  ┌──────────────────┐  ┌──────────────────┐                 │
│  │ oauth_accounts   │  │ oauth_sessions   │                 │
│  └──────────────────┘  └──────────────────┘                 │
└──────────────────────────────────────────────────────────────┘
```

### 数据流

#### 登录流程

```
1. 用户点击 "使用 Google 登录"
   ↓
2. 前端调用 POST /api/v1/shared/oauth/google/authorize
   ↓
3. 后端生成授权 URL 并返回
   ↓
4. 前端重定向用户到 Google 授权页面
   ↓
5. 用户在 Google 页面授权
   ↓
6. Google 重定向回前端，携带授权码
   ↓
7. 前端调用 POST /api/v1/shared/oauth/google/callback
   ↓
8. 后端交换授权码获取 Access Token
   ↓
9. 后端使用 Access Token 获取用户信息
   ↓
10. 后端检查用户是否存在：
    - 存在：直接登录，返回 JWT Token
    - 不存在：自动注册，返回 JWT Token
   ↓
11. 前端使用 JWT Token 登录系统
```

#### 账号绑定流程

```
1. 用户已登录，访问账号设置页面
   ↓
2. 用户点击 "绑定 GitHub 账号"
   ↓
3. 前端调用 POST /api/v1/shared/oauth/github/authorize
   （携带用户 JWT Token）
   ↓
4. 后端生成授权 URL，标记为绑定模式
   ↓
5. 前端重定向用户到 GitHub 授权页面
   ↓
6. 用户在 GitHub 页面授权
   ↓
7. GitHub 重定向回前端，携带授权码
   ↓
8. 前端调用 POST /api/v1/shared/oauth/github/callback
   ↓
9. 后端检查会话状态，发现是绑定模式
   ↓
10. 后端将 GitHub 账号绑定到当前用户
   ↓
11. 返回绑定成功结果
```

## 数据模型

### OAuthAccount

```go
type OAuthAccount struct {
    ID              string        `bson:"_id"`
    UserID          string        `bson:"user_id"`                    // 关联的用户ID
    Provider        OAuthProvider `bson:"provider"`                   // 提供商
    ProviderUserID  string        `bson:"provider_user_id"`           // 提供商用户ID
    Email           string        `bson:"email"`                      // 邮箱
    Username        string        `bson:"username"`                   // 用户名
    Avatar          string        `bson:"avatar"`                     // 头像URL
    AccessToken     string        `bson:"access_token"`               // 访问令牌（加密）
    RefreshToken    string        `bson:"refresh_token"`              // 刷新令牌（加密）
    ExpiresAt       time.Time     `bson:"expires_at"`                 // 过期时间
    TokenExpiresAt  time.Time     `bson:"token_expires_at"`           // 令牌过期时间
    Scope           string        `bson:"scope"`                      // 授权范围
    IsPrimary       bool          `bson:"is_primary"`                 // 是否为主账号
    LastLoginAt     time.Time     `bson:"last_login_at"`              // 最后登录时间
    Metadata        map[string]interface{} `bson:"metadata"`          // 额外信息
    CreatedAt       time.Time     `bson:"created_at"`                 // 创建时间
    UpdatedAt       time.Time     `bson:"updated_at"`                 // 更新时间
}
```

### OAuthSession

```go
type OAuthSession struct {
    ID              string        `bson:"_id"`
    State           string        `bson:"state"`                      // OAuth状态参数
    Provider        OAuthProvider `bson:"provider"`
    RedirectURI     string        `bson:"redirect_uri"`
    Scope           string        `bson:"scope"`
    UserID          string        `bson:"user_id,omitempty"`          // 已登录用户ID（绑定模式）
    LinkMode        bool          `bson:"link_mode"`                  // 是否为绑定模式
    ExpiresAt       time.Time     `bson:"expires_at"`                 // 过期时间（10分钟）
    CreatedAt       time.Time     `bson:"created_at"`
}
```

### UserIdentity

```go
type UserIdentity struct {
    Provider       OAuthProvider `json:"provider"`
    ProviderID     string        `json:"provider_id"`
    Email          string        `json:"email"`
    EmailVerified  bool          `json:"email_verified"`
    Name           string        `json:"name"`
    Avatar         string        `json:"avatar"`
    Username       string        `json:"username,omitempty"`
    Locale         string        `json:"locale,omitempty"`
    Gender         string        `json:"gender,omitempty"`
    Birthday       string        `json:"birthday,omitempty"`
}
```

## API 接口设计

### 1. 获取授权 URL

**端点**: `POST /api/v1/shared/oauth/{provider}/authorize`

**请求参数**:
```json
{
  "redirect_uri": "https://yourapp.com/oauth/callback",
  "state": "random_state_string"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "获取授权URL成功",
  "data": {
    "authorize_url": "https://accounts.google.com/o/oauth2/v2/auth?...",
    "provider": "google"
  }
}
```

### 2. 处理 OAuth 回调

**端点**: `POST /api/v1/shared/oauth/{provider}/callback`

**请求参数**:
```json
{
  "code": "4/0AeanS0J...",
  "state": "random_state_string"
}
```

**响应（登录模式）**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user": {
      "id": "user_id",
      "username": "john_doe",
      "email": "john@example.com"
    },
    "token": {
      "access_token": "jwt_token",
      "refresh_token": "refresh_token",
      "expires_in": 3600
    },
    "provider": "google"
  }
}
```

**响应（绑定模式）**:
```json
{
  "code": 200,
  "message": "账号绑定成功",
  "data": {
    "account": {
      "id": "account_id",
      "provider": "github",
      "username": "github_user",
      "is_primary": false
    },
    "provider": "github"
  }
}
```

### 3. 获取绑定的账号列表

**端点**: `GET /api/v1/shared/oauth/accounts`

**需要认证**: ✅

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "account_id_1",
      "provider": "google",
      "username": "john@gmail.com",
      "avatar": "https://...",
      "is_primary": true,
      "last_login_at": "2025-01-07T10:30:00Z"
    },
    {
      "id": "account_id_2",
      "provider": "github",
      "username": "github_user",
      "avatar": "https://...",
      "is_primary": false,
      "last_login_at": "2025-01-06T15:20:00Z"
    }
  ]
}
```

### 4. 解绑账号

**端点**: `DELETE /api/v1/shared/oauth/accounts/{accountID}`

**需要认证**: ✅

**响应**:
```json
{
  "code": 200,
  "message": "解绑成功"
}
```

### 5. 设置主账号

**端点**: `PUT /api/v1/shared/oauth/accounts/{accountID}/primary`

**需要认证**: ✅

**响应**:
```json
{
  "code": 200,
  "message": "设置成功"
}
```

## 安全考虑

### 1. 状态参数 (State Parameter)

- 生成随机的 `state` 参数，防止 CSRF 攻击
- 验证回调时返回的 `state` 是否匹配

### 2. 令牌加密

- Access Token 和 Refresh Token 加密后存储
- 使用 AES-256 加密算法

### 3. 会话管理

- OAuth 会话 10 分钟后自动过期
- 会话完成后立即删除

### 4. HTTPS

- 所有 OAuth 通信必须使用 HTTPS
- 生产环境禁用 HTTP

### 5. 权限范围

- 只请求必要的权限
- Google: `openid`, `email`, `profile`
- GitHub: `read:user`, `user:email`
- QQ: `get_user_info`

## 配置管理

### 环境变量

```bash
# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URI=https://yourapp.com/oauth/google/callback

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URI=https://yourapp.com/oauth/github/callback

# QQ OAuth
QQ_CLIENT_ID=your_qq_client_id
QQ_CLIENT_SECRET=your_qq_client_secret
QQ_REDIRECT_URI=https://yourapp.com/oauth/qq/callback
QQ_AUTH_URL=https://graph.qq.com/oauth2.0/authorize
QQ_TOKEN_URL=https://graph.qq.com/oauth2.0/token
```

### 配置结构

```go
type OAuthConfig struct {
    Provider     OAuthProvider `json:"provider"`
    ClientID     string        `json:"client_id"`
    ClientSecret string        `json:"client_secret"`
    AuthURL      string        `json:"auth_url"`
    TokenURL     string        `json:"token_url"`
    UserInfoURL  string        `json:"user_info_url"`
    Scopes       string        `json:"scopes"`
    Enabled      bool          `json:"enabled"`
}
```

## 错误处理

### 常见错误码

| 错误码 | 说明 | 处理方式 |
|--------|------|----------|
| `OAUTH_INVALID_PROVIDER` | 不支持的提供商 | 返回400，提示用户选择其他提供商 |
| `OAUTH_INVALID_CODE` | 授权码无效 | 返回400，提示用户重新授权 |
| `OAUTH_STATE_MISMATCH` | 状态参数不匹配 | 返回400，可能的CSRF攻击 |
| `OAUTH_ACCOUNT_ALREADY_LINKED` | 账号已绑定 | 返回400，提示用户账号已被占用 |
| `OAUTH_TOKEN_EXPIRED` | 令牌过期 | 尝试刷新令牌或重新授权 |

### 错误响应示例

```json
{
  "code": 400,
  "message": "账号已绑定",
  "error": "OAUTH_ACCOUNT_ALREADY_LINKED",
  "details": "该 Google 账号已绑定到其他用户"
}
```

## 部署清单

### 1. 第三方平台配置

#### Google OAuth

1. 访问 [Google Cloud Console](https://console.cloud.google.com/)
2. 创建项目或选择现有项目
3. 启用 Google+ API
4. 创建 OAuth 2.0 客户端ID
5. 配置授权重定向 URI

#### GitHub OAuth

1. 访问 [GitHub Developer Settings](https://github.com/settings/developers)
2. 创建 OAuth App
3. 配置 Authorization callback URL

#### QQ OAuth

1. 访问 [QQ互联平台](https://connect.qq.com/)
2. 创建网站应用
3. 配置回调地址

### 2. 数据库索引

```javascript
// oauth_accounts 集合
db.oauth_accounts.createIndex(
  { provider: 1, provider_user_id: 1 },
  { unique: true }
)

db.oauth_accounts.createIndex({ user_id: 1 })

db.oauth_accounts.createIndex(
  { user_id: 1, is_primary: 1 }
)

// oauth_sessions 集合
db.oauth_sessions.createIndex(
  { state: 1 },
  { unique: true, expireAfterSeconds: 600 }
)

db.oauth_sessions.createIndex({ expires_at: 1 })
```

## 测试计划

### 单元测试

- [ ] OAuthService.GetAuthURL()
- [ ] OAuthService.ExchangeCode()
- [ ] OAuthService.GetUserInfo()
- [ ] OAuthService.LinkAccount()
- [ ] OAuthService.UnlinkAccount()

### 集成测试

- [ ] Google OAuth 完整流程
- [ ] GitHub OAuth 完整流程
- [ ] QQ OAuth 完整流程
- [ ] 账号绑定流程
- [ ] 账号解绑流程

### 端到端测试

- [ ] 用户使用 Google 注册
- [ ] 用户使用 GitHub 登录
- [ ] 用户绑定多个账号
- [ ] 用户解绑账号
- [ ] 令牌刷新

## 未来扩展

### 计划中的提供商

- [ ] 微信登录
- [ ] 微博登录
- [ ] Facebook 登录
- [ ] Apple 登录

### 增强功能

- [ ] 社交账号合并
- [ ] 账号转移
- [ ] OAuth 审计日志
- [ ] 多因素认证集成

## 参考资料

- [OAuth 2.0 规范 (RFC 6749)](https://tools.ietf.org/html/rfc6749)
- [Google OAuth 2.0 文档](https://developers.google.com/identity/protocols/oauth2)
- [GitHub OAuth 文档](https://developer.github.com/apps/building-oauth-apps/)
- [QQ 互联文档](http://wiki.connect.qq.com/)

## 附录

### 文件清单

| 文件路径 | 说明 |
|----------|------|
| `models/auth/oauth.go` | OAuth 数据模型 |
| `service/shared/auth/oauth_service.go` | OAuth 服务实现 |
| `repository/interfaces/auth/oauth_repository.go` | OAuth 仓储接口 |
| `repository/mongodb/auth/oauth_repository_mongo.go` | OAuth MongoDB 实现 |
| `api/v1/shared/oauth_api.go` | OAuth API 处理器 |
| `doc/design/auth/第三方登录OAuth设计文档.md` | 本文档 |

### 变更日志

| 日期 | 版本 | 说明 |
|------|------|------|
| 2025-01-07 | v1.0 | 初始版本 |

---

**文档维护者**: AI Assistant
**最后更新**: 2025-01-07
