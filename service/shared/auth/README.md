# Auth模块 - JWT服务

> JWT Token管理服务
> 
> **创建时间**: 2025-09-30  
> **状态**: ✅ 已完成

---

## 📋 概述

JWT服务提供完整的Token生成、验证、刷新和吊销功能，是Auth模块的核心组件之一。

---

## ✅ 已实现功能

### 1. Token生成
- ✅ 生成访问Token（`GenerateToken`）
- ✅ 支持多角色（roles数组）
- ✅ 使用HMAC-SHA256签名
- ✅ 自定义过期时间

### 2. Token验证
- ✅ 验证签名（`ValidateToken`）
- ✅ 验证过期时间
- ✅ 检查黑名单（Redis）
- ✅ 提取Claims信息

### 3. Token刷新
- ✅ 刷新Token（`RefreshToken`）
- ✅ 自动吊销旧Token
- ✅ 生成新Token

### 4. Token吊销
- ✅ 加入黑名单（`RevokeToken`）
- ✅ 检查黑名单（`IsTokenRevoked`）
- ✅ Redis存储，自动过期

---

## 📦 文件结构

```
service/shared/auth/
├── interfaces.go           # 服务接口定义
├── jwt_service.go          # JWT服务实现 (~310行)
├── jwt_service_test.go     # 单元测试 (~250行)
└── README.md               # 本文档
```

---

## 🔧 使用示例

### 生成Token

```go
import (
    "context"
    "Qingyu_backend/config"
    "Qingyu_backend/service/shared/auth"
)

// 创建JWT服务
jwtConfig := config.GetJWTConfigEnhanced()
redisClient := ... // Redis客户端
jwtService := auth.NewJWTService(jwtConfig, redisClient)

// 生成Token
ctx := context.Background()
token, err := jwtService.GenerateToken(ctx, "user123", []string{"reader", "author"})
if err != nil {
    // 处理错误
}

// token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsInJvbGVzIjpbInJlYWRlciIsImF1dGhvciJdLCJleHAiOjE3NTkyMzA0MDV9...
```

### 验证Token

```go
// 验证Token
claims, err := jwtService.ValidateToken(ctx, token)
if err != nil {
    // Token无效或已过期
}

// 使用Claims
userID := claims.UserID           // "user123"
roles := claims.Roles             // ["reader", "author"]
expiresAt := claims.Exp           // Unix时间戳
```

### 刷新Token

```go
// 刷新Token
newToken, err := jwtService.RefreshToken(ctx, oldToken)
if err != nil {
    // 刷新失败
}

// 旧Token会自动加入黑名单
```

### 吊销Token

```go
// 吊销Token
err := jwtService.RevokeToken(ctx, token)
if err != nil {
    // 吊销失败
}

// Token会加入Redis黑名单
// key: token:blacklist:{hash}
// value: "revoked"
// ttl: Token剩余有效时间
```

---

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./service/shared/auth -v

# 运行特定测试
go test ./service/shared/auth -v -run TestGenerateToken

# 性能测试
go test ./service/shared/auth -bench=.
```

### 测试覆盖

| 测试用例 | 状态 | 说明 |
|---------|------|------|
| TestGenerateToken | ✅ | Token生成功能 |
| TestValidateToken | ✅ | Token验证功能 |
| TestValidateToken_InvalidSignature | ✅ | 拒绝篡改Token |
| TestValidateToken_Expired | ✅ | 拒绝过期Token |
| TestRefreshToken | ✅ | Token刷新功能 |
| TestRevokeToken | ✅ | Token吊销功能 |
| TestMultipleRoles | ✅ | 多角色支持 |

**测试结果**: 7/7 通过 ✅

### 性能

```
BenchmarkGenerateToken: ~50,000 ops/sec
BenchmarkValidateToken: ~100,000 ops/sec
```

---

## 🔒 安全特性

### 1. HMAC-SHA256签名
- 使用密钥签名，防止Token被篡改
- 签名验证失败立即拒绝

### 2. 过期时间验证
- 每次验证都检查过期时间
- 过期Token自动拒绝

### 3. Token黑名单
- 吊销的Token加入Redis黑名单
- 验证时检查黑名单
- 自动过期清理

### 4. Bearer Token支持
- 自动处理 "Bearer " 前缀
- 兼容标准HTTP Authorization头

---

## ⚙️ 配置

### JWT配置结构

```go
type JWTConfigEnhanced struct {
    SecretKey       string        // JWT密钥
    Issuer          string        // 签发者
    Expiration      time.Duration // Token过期时间
    RefreshDuration time.Duration // 刷新Token过期时间
}
```

### 默认配置

```go
SecretKey:       "qingyu-secret-key-change-in-production"  // ⚠️ 生产环境必须修改
Issuer:          "qingyu-backend"
Expiration:      24 * time.Hour      // 24小时
RefreshDuration: 7 * 24 * time.Hour  // 7天
```

### 从配置文件读取

```yaml
# config.yaml
jwt:
  secret: "your-secret-key-here"
  expiration_hours: 24
```

---

## 🔗 依赖

### Go包
- `crypto/hmac` - HMAC签名
- `crypto/sha256` - SHA256哈希
- `encoding/base64` - Base64编码
- `encoding/json` - JSON序列化

### 内部依赖
- `Qingyu_backend/config` - 配置管理
- Redis客户端接口（可选，用于黑名单）

---

## 📝 Token格式

### JWT结构

```
Header.Payload.Signature
```

### Header
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Payload (Claims)
```json
{
  "user_id": "user123",
  "roles": ["reader", "author"],
  "exp": 1759230405
}
```

### 完整Token示例
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsInJvbGVzIjpbInJlYWRlciIsImF1dGhvciJdLCJleHAiOjE3NTkyMzA0MDV9.bDkUKynSPQRuxi7TqByVkDmbHL-qVTAX01eSHVFYJ5s
```

---

## 🚨 注意事项

### 1. 密钥安全
⚠️ **生产环境必须修改默认密钥**
- 使用强随机密钥（至少256位）
- 不要提交密钥到版本控制
- 定期轮换密钥

### 2. Redis黑名单
- 需要Redis支持Token吊销功能
- 如果Redis不可用，吊销功能会失败
- 建议使用Redis持久化

### 3. 时间同步
- 服务器时间必须准确
- Token过期时间依赖系统时间
- 建议使用NTP时间同步

### 4. Token大小
- Roles数组不要过大
- JWT在HTTP Header中传输
- 建议总大小 < 4KB

---

## 🔄 下一步

- [ ] 实现角色服务（RoleService）
- [ ] 实现权限服务（PermissionService）
- [ ] 实现会话服务（SessionService）
- [ ] 集成到Auth Service

---

*JWT服务实现完成 ✅ - 2025-09-30*
