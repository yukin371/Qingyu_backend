# Auth模块

## 概述

Auth模块是Qingyu后端项目的核心认证授权服务，提供完整的用户认证、权限控制、会话管理和OAuth第三方登录功能。

## 目录结构

```
service/auth/
├── auth_service.go           # 认证服务主实现
├── jwt_service.go            # JWT令牌服务
├── oauth_service.go          # OAuth第三方登录服务
├── session_service.go        # 会话管理服务
├── permission_service.go     # 权限检查服务
├── role_service.go           # 角色管理服务
├── password_validator.go     # 密码强度验证
├── redis_adapter.go          # Redis存储适配器
├── memory_blacklist.go       # 内存令牌黑名单
├── interfaces.go             # 服务接口定义
├── _migration/               # 迁移兼容层
│   └── shared_compat.go      # 向后兼容导出
└── README.md                 # 本文档
```

**注意**: 所有文件都在同一个包（package auth）中，这是为了保持与原有代码结构的兼容性。

## 核心功能

### 1. 用户认证
- **注册**: 用户注册，支持默认角色分配
- **登录**: 用户名密码登录，支持多端登录限制
- **登出**: 令牌撤销，黑名单机制
- **令牌刷新**: JWT令牌自动刷新
- **令牌验证**: 令牌有效性检查

### 2. OAuth第三方登录
- 支持多个OAuth提供商（Google、GitHub、QQ、微信、微博）
- 自动用户创建
- OAuth账号绑定
- 主账号标识

### 3. 权限管理
- 基于角色的访问控制（RBAC）
- 权限检查API
- 用户权限查询
- 角色权限缓存

### 4. 角色管理
- 角色CRUD操作
- 用户角色分配和移除
- 预定义角色（Admin、Editor、Reader）

### 5. 会话管理
- 多端登录限制
- 设备数量强制执行
- 会话创建和销毁
- 会话刷新机制

### 6. 密码安全
- 密码强度验证
- 可配置密码规则
- 随机密码生成（OAuth用户）

## 使用方法

### 基本使用

```go
import "Qingyu_backend/service/auth"

// 创建服务实例
authService := auth.NewAuthService(
    jwtService,
    roleService,
    permissionService,
    authRepo,
    oauthRepo,
    userService,
    sessionService,
)

// 初始化服务
if err := authService.Initialize(ctx); err != nil {
    log.Fatal("Failed to initialize auth service:", err)
}

// 用户登录
loginResp, err := authService.Login(ctx, &auth.LoginRequest{
    Username: "user@example.com",
    Password: "password123",
})
```

### 权限检查

```go
// 检查用户权限
hasPermission, err := authService.CheckPermission(ctx, userID, "post:write")

// 获取用户所有权限
permissions, err := authService.GetUserPermissions(ctx, userID)

// 检查用户角色
hasRole, err := authService.HasRole(ctx, userID, "admin")
```

### 角色管理

```go
// 创建角色
role, err := authService.CreateRole(ctx, &auth.CreateRoleRequest{
    Name:        "moderator",
    Description: "版主角色",
    Permissions: []string{"post:edit", "post:delete"},
})

// 分配角色
err := authService.AssignRole(ctx, userID, roleID)

// 移除角色
err := authService.RemoveRole(ctx, userID, roleID)
```

## 配置

### JWT配置

```go
type JWTConfig struct {
    SecretKey       string        // JWT密钥
    TokenDuration   time.Duration // 令牌有效期
    RefreshDuration time.Duration // 刷新令牌有效期
    Issuer          string        // 发行者
}
```

### OAuth配置

```go
type OAuthConfig struct {
    Providers map[string]OAuthProviderConfig
}

type OAuthProviderConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
    Scopes       []string
}
```

### 会话配置

```go
type SessionConfig struct {
    MaxDevices    int           // 最大设备数量
    TokenDuration time.Duration // 会话有效期
}
```

## 迁移计划

### 当前状态：第一阶段 - 目录结构创建

- [x] 创建新的service/auth目录结构
- [x] 复制所有源文件到新位置
- [x] 创建迁移兼容层
- [x] 创建模块文档

### 下一步：第二阶段 - 依赖更新

- [ ] 更新service容器注册
- [ ] 更新API层import路径
- [ ] 更新middleware层import路径
- [ ] 更新Port接口适配器
- [ ] 运行测试验证

### 最终阶段：清理

- [ ] 删除旧的service/shared/auth目录
- [ ] 移除迁移兼容层
- [ ] 更新所有文档

## 依赖关系

### 内部依赖
- `service/interfaces/user`: 用户服务接口
- `repository/interfaces/auth`: 认证仓储接口
- `repository/interfaces/shared`: 共享仓储接口
- `models/auth`: 认证数据模型
- `models/users`: 用户数据模型

### 外部依赖
- `go.uber.org/zap`: 日志记录
- `crypto/rand`: 随机数生成
- `encoding/base64`: 编码解码

## 向后兼容性

为了确保平滑迁移，本模块提供了向后兼容层：

```go
// 旧的import路径仍然可用
import "Qingyu_backend/service/shared/auth"

// 实际实现来自新路径
// import "Qingyu_backend/service/auth"
```

兼容层将在所有代码迁移完成后移除（计划版本：v1.2）

## 测试

### 运行单元测试

```bash
# 运行所有auth模块测试
go test ./service/auth/...

# 运行特定服务测试
go test ./service/auth/services/jwt_service_test.go
```

### 测试覆盖率

```bash
# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./service/auth/...
go tool cover -html=coverage.out
```

## 维护者

- 创建日期: 2026-02-09
- 当前版本: v1.0.0
- 状态: 迁移中

## 相关文档

- [迁移计划](../../docs/plan/shared-module-refactor-plan.md)
- [依赖规则](../../docs/architecture/dependency-rules.md)
- [接口契约](../../service/interfaces/shared/README.md)
