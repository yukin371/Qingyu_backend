# JWT身份认证设计文档

## 1. 概述

### 1.1 项目背景
青羽后端系统采用JWT（JSON Web Token）作为用户身份认证和授权的核心机制。JWT提供了无状态的认证方案，适合分布式系统和微服务架构。

### 1.2 设计目标
- **无状态认证**：服务器不需要存储会话信息
- **安全性**：防止token伪造、重放攻击等安全威胁
- **可扩展性**：支持多服务间的认证共享
- **用户体验**：支持token自动刷新和优雅的过期处理

## 2. 技术架构

### 2.1 技术栈
- **JWT库**：`github.com/dgrijalva/jwt-go`
- **密码加密**：`golang.org/x/crypto/bcrypt`
- **Web框架**：Gin
- **数据库**：MongoDB
- **配置管理**：Viper

### 2.2 架构层次
```
┌─────────────────┐
│   Client        │
├─────────────────┤
│   Router        │ ← JWT中间件验证
├─────────────────┤
│   API Layer     │ ← 用户认证API
├─────────────────┤
│   Service Layer │ ← 认证业务逻辑
├─────────────────┤
│   Model Layer   │ ← 用户数据模型
├─────────────────┤
│   Database      │ ← MongoDB存储
└─────────────────┘
```

## 3. 数据模型设计

### 3.1 用户模型
基于现有的 `models/system/sys_user.go`：

```go
type User struct {
    ID        string    `bson:"_id,omitempty" json:"id"`
    Username  string    `bson:"username" json:"username"`
    Email     string    `bson:"email" json:"email"`
    Password  string    `bson:"password" json:"-"`
    CreatedAt time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
```

**字段说明**：
- `ID`：用户唯一标识符
- `Username`：用户名（唯一）
- `Email`：邮箱地址（唯一）
- `Password`：bcrypt加密的密码哈希
- `CreatedAt/UpdatedAt`：时间戳

### 3.2 JWT配置模型
基于现有的 `config/jwt.go`：

```go
type JWTConfig struct {
    Secret          string // JWT签名密钥
    ExpirationHours int    // token过期时间（小时）
}
```

## 4. JWT Token设计

### 4.1 Token结构
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "用户ID",
    "username": "用户名",
    "email": "邮箱",
    "iat": 1640995200,
    "exp": 1641081600
  },
  "signature": "签名"
}
```

### 4.2 Claims定义
```go
type JWTClaims struct {
    UserID   string `json:"sub"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.StandardClaims
}
```

## 5. 业务逻辑设计

### 5.1 用户服务 (UserService)
```go
type UserService struct{}

// 用户注册
func (s *UserService) Register(username, email, password string) (*User, error)

// 用户登录
func (s *UserService) Login(username, password string) (*User, string, error)

// 验证用户
func (s *UserService) ValidateUser(userID string) (*User, error)

// 更新用户信息
func (s *UserService) UpdateProfile(userID string, updates *User) (*User, error)
```

### 5.2 JWT工具服务 (JWTService)
```go
type JWTService struct {
    config *config.JWTConfig
}

// 生成JWT token
func (s *JWTService) GenerateToken(user *User) (string, error)

// 验证JWT token
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error)

// 解析token获取用户信息
func (s *JWTService) ParseToken(tokenString string) (*JWTClaims, error)

// 刷新token
func (s *JWTService) RefreshToken(tokenString string) (string, error)
```

### 5.3 认证中间件 (AuthMiddleware)
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从Header获取token
        // 2. 验证token有效性
        // 3. 解析用户信息
        // 4. 设置用户上下文
        // 5. 继续处理请求
    }
}
```

## 6. API接口设计

### 6.1 认证相关接口

#### 用户注册
```
POST /api/v1/user/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}

Response:
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": "user_id",
    "username": "testuser",
    "email": "test@example.com",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 用户登录
```
POST /api/v1/user/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

Response:
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "user": {
      "id": "user_id",
      "username": "testuser",
      "email": "test@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": "2024-01-02T00:00:00Z"
  }
}
```

#### 获取用户信息
```
GET /api/v1/user/profile
Authorization: Bearer <token>

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "user_id",
    "username": "testuser",
    "email": "test@example.com",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

#### 更新用户信息
```
PUT /api/v1/user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "newusername",
  "email": "newemail@example.com"
}

Response:
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": "user_id",
    "username": "newusername",
    "email": "newemail@example.com",
    "updatedAt": "2024-01-01T12:00:00Z"
  }
}
```

#### Token刷新
```
POST /api/v1/user/refresh
Authorization: Bearer <token>

Response:
{
  "code": 0,
  "message": "刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": "2024-01-02T12:00:00Z"
  }
}
```

### 6.2 错误响应格式
```json
{
  "code": 40001,
  "message": "token已过期",
  "timestamp": 1640995200
}
```

**错误码定义**：
- `40001`: token已过期
- `40002`: token无效
- `40003`: 用户不存在
- `40004`: 密码错误
- `40005`: 用户名已存在
- `40006`: 邮箱已存在

## 7. 数据库设计

### 7.1 用户集合 (users)
```javascript
{
  "_id": "ObjectId或自定义ID",
  "username": "用户名",
  "email": "邮箱地址",
  "password": "bcrypt加密的密码哈希",
  "created_at": "创建时间",
  "updated_at": "更新时间"
}
```

### 7.2 索引设计
```javascript
// 用户名唯一索引
db.users.createIndex({"username": 1}, {"unique": true})

// 邮箱唯一索引
db.users.createIndex({"email": 1}, {"unique": true})

// 创建时间索引（用于排序）
db.users.createIndex({"created_at": -1})
```

## 8. 安全性设计

### 8.1 密码安全
- 使用bcrypt进行密码哈希
- 设置合适的cost参数（默认12）
- 密码复杂度要求：至少8位，包含字母和数字

### 8.2 JWT安全
- 使用强随机密钥（至少256位）
- 设置合理的过期时间（默认24小时）
- 支持token黑名单机制
- 防止XSS攻击：token存储在httpOnly cookie中

### 8.3 API安全
- HTTPS传输加密
- 请求频率限制
- 输入参数验证和过滤
- SQL注入防护

### 8.4 会话管理
```go
// Token黑名单（可选实现）
type TokenBlacklist struct {
    Token     string    `bson:"token"`
    ExpiresAt time.Time `bson:"expires_at"`
    CreatedAt time.Time `bson:"created_at"`
}
```

## 9. 中间件实现

### 9.1 认证中间件
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"code": 40001, "message": "缺少认证token"})
            c.Abort()
            return
        }

        // 解析Bearer token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // 验证token
        jwtService := &JWTService{config: config.LoadJWTConfig()}
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"code": 40002, "message": "token无效"})
            c.Abort()
            return
        }

        // 设置用户上下文
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

### 9.2 可选认证中间件
```go
func OptionalAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader != "" {
            // 尝试解析token，但不强制要求
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            jwtService := &JWTService{config: config.LoadJWTConfig()}
            if claims, err := jwtService.ValidateToken(tokenString); err == nil {
                c.Set("userID", claims.UserID)
                c.Set("username", claims.Username)
            }
        }
        c.Next()
    }
}
```

## 10. 配置管理

### 10.1 环境变量
```bash
# JWT配置
JWT_SECRET=your_super_secret_key_here
JWT_EXPIRATION_HOURS=24

# 数据库配置
MONGO_URI=mongodb://localhost:27017
DB_NAME=qingyu_backend
```

### 10.2 配置加载
基于现有的配置系统：
```go
func LoadJWTConfig() *JWTConfig {
    return &JWTConfig{
        Secret:          getEnv("JWT_SECRET", "qingyu_secret_key"),
        ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
    }
}
```

## 11. 错误处理

### 11.1 认证错误类型
```go
var (
    ErrTokenExpired     = errors.New("token已过期")
    ErrTokenInvalid     = errors.New("token无效")
    ErrUserNotFound     = errors.New("用户不存在")
    ErrPasswordInvalid  = errors.New("密码错误")
    ErrUsernameExists   = errors.New("用户名已存在")
    ErrEmailExists      = errors.New("邮箱已存在")
)
```

### 11.2 统一错误响应
```go
func HandleAuthError(c *gin.Context, err error) {
    var code int
    var message string
    
    switch err {
    case ErrTokenExpired:
        code, message = 40001, "token已过期"
    case ErrTokenInvalid:
        code, message = 40002, "token无效"
    case ErrUserNotFound:
        code, message = 40003, "用户不存在"
    case ErrPasswordInvalid:
        code, message = 40004, "密码错误"
    default:
        code, message = 50000, "内部服务器错误"
    }
    
    c.JSON(http.StatusUnauthorized, gin.H{
        "code":      code,
        "message":   message,
        "timestamp": time.Now().Unix(),
    })
}
```

## 12. 测试策略

### 12.1 单元测试
- JWT token生成和验证
- 密码加密和验证
- 用户模型方法测试

### 12.2 集成测试
- 用户注册流程
- 用户登录流程
- 认证中间件测试
- API接口测试

### 12.3 安全测试
- token伪造测试
- 过期token测试
- 权限绕过测试
- 密码暴力破解测试

## 13. 部署和监控

### 13.1 部署配置
- 生产环境使用强随机JWT密钥
- 启用HTTPS
- 配置反向代理
- 设置适当的CORS策略

### 13.2 监控指标
- 登录成功/失败率
- token验证成功/失败率
- API响应时间
- 异常登录检测

### 13.3 日志记录
```go
// 认证相关日志
log.Info("用户登录成功", 
    zap.String("userID", user.ID),
    zap.String("username", user.Username),
    zap.String("ip", c.ClientIP()))

log.Warn("token验证失败",
    zap.String("token", tokenString),
    zap.String("error", err.Error()),
    zap.String("ip", c.ClientIP()))
```

## 14. 扩展性设计

### 14.1 多租户支持
```go
type JWTClaims struct {
    UserID   string `json:"sub"`
    Username string `json:"username"`
    Email    string `json:"email"`
    TenantID string `json:"tenant_id"` // 租户ID
    jwt.StandardClaims
}
```

### 14.2 角色权限系统
```go
type User struct {
    ID        string   `bson:"_id,omitempty" json:"id"`
    Username  string   `bson:"username" json:"username"`
    Email     string   `bson:"email" json:"email"`
    Password  string   `bson:"password" json:"-"`
    Roles     []string `bson:"roles" json:"roles"`     // 用户角色
    CreatedAt time.Time `bson:"created_at" json:"createdAt"`
    UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
```

### 14.3 OAuth2集成
- 支持第三方登录（Google、GitHub等）
- 统一的认证接口
- token互换机制

## 15. 最佳实践

### 15.1 开发规范
- 遵循项目分层架构
- 统一的错误处理
- 完善的日志记录
- 充分的单元测试

### 15.2 安全规范
- 定期更换JWT密钥
- 监控异常登录行为
- 实施账户锁定策略
- 定期安全审计

### 15.3 性能优化
- token缓存策略
- 数据库连接池优化
- 异步日志记录
- 合理的索引设计

---

## 关联文档
- 软件需求规格说明书(SRS) ../软件需求规格说明书(SRS).md
- 架构设计说明书 ../架构设计说明书.md
- API 接口总览 ../API接口总览.md
- 数据库设计说明书 ../数据库设计说明书.md
- 测试计划与用例 ../测试计划与用例.md
- 部署与运维指南 ../部署与运维指南.md
- 安全设计与威胁建模 ../安全设计与威胁建模.md
- 日志与监控 ../日志与监控.md
- 需求追踪矩阵 ../需求追踪矩阵.md


本设计文档基于青羽后端项目的现有架构和开发规范，提供了完整的JWT身份认证解决方案。文档将随着项目发展持续更新和完善。