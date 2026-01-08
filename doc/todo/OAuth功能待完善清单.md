# OAuth2.0 第三方登录功能 - 待完善清单

**创建日期**: 2025-01-07
**状态**: 🚧 暂时搁置
**当前进度**: 核心功能已完成 80%

## 概述

OAuth2.0 第三方登录功能已完成核心代码实现，但尚未完全集成到系统主流程中。本文档记录当前完成状态和待完成任务。

## ✅ 已完成的工作

### 1. 数据模型层 (100%)
- ✅ `models/auth/oauth.go` - OAuth 数据模型定义
  - `OAuthProvider` 类型
  - `OAuthAccount` 结构
  - `OAuthSession` 结构
  - `UserIdentity` 结构

### 2. 仓储层 (100%)
- ✅ `repository/interfaces/auth/oauth_repository.go` - OAuth 仓储接口
- ✅ `repository/mongodb/auth/oauth_repository_mongo.go` - MongoDB 实现
  - 账号管理（CRUD）
  - 会话管理
  - 令牌管理
  - 索引创建方法

### 3. 服务层 (90%)
- ✅ `service/shared/auth/oauth_service.go` - OAuth 核心服务
  - `GetAuthURL()` - 获取授权URL
  - `ExchangeCode()` - 交换授权码
  - `GetUserInfo()` - 获取用户信息（Google、GitHub、QQ）
  - `LinkAccount()` - 绑定账号
  - `UnlinkAccount()` - 解绑账号
  - `RefreshToken()` - 刷新令牌
- ✅ `service/shared/auth/auth_service.go` - AuthService 集成
  - 添加 `oauthRepo` 依赖
  - 实现 `OAuthLogin()` 方法
  - 添加辅助函数

### 4. API 层 (100%)
- ✅ `api/v1/shared/oauth_api.go` - OAuth API 处理器
  - `GetAuthorizeURL()` - 获取授权URL
  - `HandleCallback()` - 处理OAuth回调
  - `GetLinkedAccounts()` - 获取绑定账号
  - `UnlinkAccount()` - 解绑账号
  - `SetPrimaryAccount()` - 设置主账号

### 5. 配置管理 (100%)
- ✅ `config/oauth_config.go` - OAuth 配置管理器
  - `LoadFromEnv()` - 从环境变量加载
  - `LoadFromConfig()` - 从配置文件加载
  - `GetConfig()` / `IsProviderEnabled()` 等方法
- ✅ `config/config.go` - 添加 OAuth 配置字段

### 6. 路由注册 (90%)
- ✅ `router/shared/shared_router.go` - OAuth 路由定义
  - 更新 `RegisterRoutes()` 签名
  - 更新 `RegisterAuthRoutes()` 添加 OAuth 路由

### 7. 文档 (100%)
- ✅ `doc/design/auth/第三方登录OAuth设计文档.md` - 完整设计文档
- ✅ `doc/design/auth/OAuth集成完成报告.md` - 集成完成报告

## ❌ 待完成的工作

### 1. 服务容器集成 (0%)

**优先级**: 🔴 高

**文件**: `service/container/service_container.go`

**任务**:
```go
// 在 ServiceContainer 结构体中添加
oauthService *auth.OAuthService

// 在 Initialize() 方法中初始化
func (c *ServiceContainer) Initialize(ctx context.Context, cfg *config.Config) error {
    // ... 现有代码 ...

    // 初始化 OAuth 配置管理器
    oauthConfigMgr := config.NewOAuthConfigManager()
    oauthConfigMgr.LoadFromEnv()
    oauthConfigMgr.LoadFromConfig(cfg)

    // 创建 OAuth 仓储
    oauthRepo := mongoAuth.NewMongoOAuthRepository(c.mongoDB)

    // 创建 OAuth 服务
    c.oauthService = auth.NewOAuthService(
        global.Logger,
        oauthRepo,
        oauthConfigMgr.GetConfigs(),
    )

    // ... 现有代码 ...
}

// 添加获取方法
func (c *ServiceContainer) GetOAuthService() (*auth.OAuthService, error) {
    if c.oauthService == nil {
        return nil, fmt.Errorf("OAuthService未初始化")
    }
    return c.oauthService, nil
}
```

### 2. 主路由注册更新 (0%)

**优先级**: 🔴 高

**文件**: `router/enter.go`

**任务**:
```go
// 在 RegisterAllRoutes() 函数中更新
func RegisterAllRoutes(r *gin.Engine, container *container.ServiceContainer) {
    // 获取 OAuth 服务
    oauthService, err := container.GetOAuthService()
    if err != nil {
        global.Logger.Warn("OAuth服务未初始化", zap.Error(err))
        oauthService = nil
    }

    // 注册共享路由（传递 OAuth 服务）
    shared.RegisterRoutes(
        r.Group("/api/v1/shared"),
        authService,
        oauthService,  // 新增参数
        global.Logger,
        walletService,
        storageService,
        multipartService,
        imageProcessor,
    )
}
```

### 3. 数据库索引创建 (0%)

**优先级**: 🟡 中

**任务**:
```go
// 在服务初始化时调用
oauthRepo := mongoAuth.NewMongoOAuthRepository(db)
if err := oauthRepo.EnsureIndexes(ctx); err != nil {
    global.Logger.Error("创建OAuth索引失败", zap.Error(err))
}
```

**需要创建的索引**:
```javascript
// oauth_accounts 集合
db.oauth_accounts.createIndex(
  { provider: 1, provider_user_id: 1 },
  { unique: true }
)
db.oauth_accounts.createIndex({ user_id: 1 })
db.oauth_accounts.createIndex({ user_id: 1, is_primary: 1 })

// oauth_sessions 集合
db.oauth_sessions.createIndex(
  { state: 1 },
  { unique: true, expireAfterSeconds: 600 }
)
db.oauth_sessions.createIndex({ expires_at: 1 })
```

### 4. 环境变量配置 (0%)

**优先级**: 🟡 中

**文件**: `.env` 或 `config.yaml`

**任务**: 添加以下环境变量
```bash
# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# QQ OAuth
QQ_CLIENT_ID=your_qq_client_id
QQ_CLIENT_SECRET=your_qq_client_secret
```

### 5. 第三方平台应用注册 (0%)

**优先级**: 🟢 低（部署前完成）

#### Google OAuth
1. 访问 [Google Cloud Console](https://console.cloud.google.com/)
2. 创建 OAuth 2.0 客户端ID
3. 配置授权重定向 URI: `https://yourdomain.com/oauth/google/callback`
4. 获取 Client ID 和 Client Secret

#### GitHub OAuth
1. 访问 [GitHub Developer Settings](https://github.com/settings/developers)
2. 创建 OAuth App
3. 配置 Authorization callback URL: `https://yourdomain.com/oauth/github/callback`
4. 获取 Client ID 和 Client Secret

#### QQ 互联
1. 访问 [QQ互联平台](https://connect.qq.com/)
2. 创建网站应用
3. 配置回调地址
4. 获取 App ID 和 App Key

### 6. 单元测试 (0%)

**优先级**: 🟡 中

**需要测试的文件**:
- [ ] `service/shared/auth/oauth_service_test.go`
- [ ] `repository/mongodb/auth/oauth_repository_test.go`
- [ ] `api/v1/shared/oauth_api_test.go`

### 7. 集成测试 (0%)

**优先级**: 🟢 低

**测试场景**:
- [ ] Google 完整登录流程
- [ ] GitHub 完整登录流程
- [ ] QQ 完整登录流程
- [ ] 账号绑定流程
- [ ] 账号解绑流程
- [ ] 错误处理

### 8. 令牌加密实现 (0%)

**优先级**: 🟡 中

**当前状态**: Token 明文存储

**需要实现**:
```go
// 加密方法
func encryptToken(plaintext, key string) (string, error)

// 解密方法
func decryptToken(ciphertext, key string) (string, error)

// 在 OAuthAccount 保存前加密
account.AccessToken = encryptToken(token.AccessToken, encryptionKey)
account.RefreshToken = encryptToken(token.RefreshToken, encryptionKey)

// 在使用时解密
accessToken := decryptToken(account.AccessToken, encryptionKey)
```

### 9. Redis 会话存储 (0%)

**优先级**: 🟢 低

**当前状态**: 使用内存存储（重启丢失）

**需要实现**: 将 OAuth 会话存储到 Redis

### 10. 更多提供商实现 (0%)

**优先级**: 🟢 低

- [ ] 微信登录实现
- [ ] 微博登录实现
- [ ] Facebook 登录实现
- [ ] Apple 登录实现

## 🚀 快速恢复开发指南

### 步骤 1: 更新服务容器

```bash
# 编辑 service/container/service_container.go
# 按照上述"服务容器集成"部分的代码进行修改
```

### 步骤 2: 更新路由注册

```bash
# 编辑 router/enter.go
# 按照上述"主路由注册更新"部分的代码进行修改
```

### 步骤 3: 配置环境变量

```bash
# 在 .env 文件中添加 OAuth 配置
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
# ... 其他提供商
```

### 步骤 4: 创建数据库索引

```bash
# 启动应用后，索引会自动创建
# 或手动运行 MongoDB 命令创建
```

### 步骤 5: 测试

```bash
# 编译并运行
go run cmd/server/main.go

# 测试获取授权URL
curl -X POST http://localhost:8080/api/v1/shared/oauth/google/authorize \
  -H "Content-Type: application/json" \
  -d '{"redirect_uri":"http://localhost:3000/callback","state":"test"}'
```

## 📁 相关文件清单

### 核心文件
```
models/auth/oauth.go
service/shared/auth/oauth_service.go
service/shared/auth/auth_service.go (已更新)
service/shared/auth/interfaces.go (已更新)
repository/interfaces/auth/oauth_repository.go
repository/mongodb/auth/oauth_repository_mongo.go
api/v1/shared/oauth_api.go
config/oauth_config.go
config/config.go (已更新)
router/shared/shared_router.go (已更新)
```

### 待更新文件
```
service/container/service_container.go (需更新)
router/enter.go (需更新)
```

### 文档文件
```
doc/design/auth/第三方登录OAuth设计文档.md
doc/design/auth/OAuth集成完成报告.md
doc/todo/OAuth功能待完善清单.md (本文档)
```

## ⚠️ 注意事项

1. **安全**: 生产环境必须使用 HTTPS
2. **密钥管理**: Client Secret 不应提交到代码仓库
3. **令牌安全**: 考虑实现令牌加密
4. **会话管理**: 生产环境建议使用 Redis 存储会话
5. **速率限制**: 已配置速率限制，可根据实际情况调整

## 🔗 参考链接

- [OAuth 2.0 规范 (RFC 6749)](https://tools.ietf.org/html/rfc6749)
- [Google OAuth 2.0 文档](https://developers.google.com/identity/protocols/oauth2)
- [GitHub OAuth 文档](https://developer.github.com/apps/building-oauth-apps/)
- [QQ 互联文档](http://wiki.connect.qq.com/)

## 📝 变更日志

| 日期 | 状态 | 说明 |
|------|------|------|
| 2025-01-07 | 🚧 暂时搁置 | 核心功能已完成，待系统集成 |

---

**文档维护者**: AI Assistant
**最后更新**: 2025-01-07
**预计恢复时间**: 待定
