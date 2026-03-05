# Issue #012: 401认证错误和权限配置问题

**优先级**: 高 (P1)
**类型**: 安全/认证问题
**状态**: 待处理
**创建日期**: 2026-03-05
**来源报告**: [401认证错误诊断报告](../reports/archived/2026-02-08-auth-401-issue-diagnosis.md)

---

## 问题描述

用户登录后访问 `/api/v1/writer/projects` 返回 401 Unauthorized 错误，涉及 JWT 认证和权限配置问题。

---

## 具体问题

### 1. 权限配置不完整 🟡 P1

**问题**: `author` 角色没有明确的 `project:*` 权限定义。

```yaml
# permissions.yaml (第56-70行)
# 作者权限
author:
  - "book:read"
  - "book:create"
  - "book:update"
  - "book:delete"
  - "chapter:read"
  - "chapter:create"
  - "chapter:update"
  - "chapter:delete"
  - "ai:generate"
  - "ai:chat"
  - "document:read"
  - "document:create"
  # ❌ 缺少 project:* 相关权限
```

**影响**:
- `testauthor001` 账号无法访问 Project 相关接口
- `testadmin001` 账号可以访问（因为有 `*:*` 权限）

**修复方案**: 在 `permissions.yaml` 中为 `author` 角色添加 project 权限

```yaml
author:
  - "book:read"
  - "book:create"
  - "book:update"
  - "book:delete"
  - "chapter:read"
  - "chapter:create"
  - "chapter:update"
  - "chapter:delete"
  - "ai:generate"
  - "ai:chat"
  - "document:read"
  - "document:create"
  - "document:update"
  - "document:delete"
  - "project:read"      # 新增
  - "project:create"    # 新增
  - "project:update"    # 新增
  - "project:delete"    # 新增
```

---

### 2. 权限中间件未正确启用 🟡 P1

**问题**: `internal/router/setup.go` 中的 `SetupAuthMiddleware` 返回一个空的中间件。

```go
// internal/router/setup.go (第54-61行)
func SetupAuthMiddleware() gin.HandlerFunc {
    // ❌ 返回空中间件，不做任何权限检查
    return func(c *gin.Context) {
        c.Next()
    }
}
```

**影响**: 权限配置虽然存在，但没有实际生效

**修复方案**: 实现真正的权限检查中间件

```go
func SetupAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取用户信息（从 JWT 中解析）
        user, exists := c.Get("user")
        if !exists {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // 2. 获取所需权限
        requiredPermissions := c.GetString("permissions")
        if requiredPermissions == "" {
            c.Next()
            return
        }

        // 3. 检查用户是否有所需权限
        userInfo := user.(*users.User)
        if !hasPermission(userInfo, requiredPermissions) {
            c.JSON(403, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

---

### 3. 权限检查分散在 Handler 层 🔵 P2

**问题**: 部分权限检查逻辑硬编码在 Handler 层，缺少统一的权限验证器。

**影响**:
- 权限逻辑难以维护
- 容易出现安全漏洞
- 代码重复

**修复方案**: 实现统一的权限验证器

```go
// pkg/auth/permission_checker.go
type PermissionChecker struct {
    permissions *PermissionsConfig
}

func (pc *PermissionChecker) HasPermission(user *users.User, requiredPermission string) bool {
    // 1. 检查用户角色
    // 2. 检查角色是否有该权限
    // 3. 返回检查结果
}

// 使用示例
func (h *ProjectHandler) GetProjects(c *gin.Context) {
    // 权限检查
    if !h.permissionChecker.HasPermission(getUser(c), "project:read") {
        c.JSON(403, gin.H{"error": "Permission denied"})
        return
    }

    // 业务逻辑
    // ...
}
```

---

### 4. 角色和权限关系不清晰 🔵 P2

**问题**:
- 角色定义在用户模型中
- 权限定义在配置文件中
- 两者关系不清晰

**当前模型**:
```go
// models/users/user.go
type UserRole string

const (
    RoleReader UserRole = "reader"
    RoleAuthor UserRole = "author"
    RoleAdmin  UserRole = "admin"
)

type User struct {
    Roles []string `bson:"roles" json:"roles"`
}
```

**建议改进**:
1. 明确角色和权限的层次关系
2. 支持角色的继承
3. 支持自定义权限组合

---

## 实施计划

### Phase 1: 权限配置修复（1天）

1. **更新 permissions.yaml**
   - 为 `author` 角色添加 `project:*` 权限
   - 检查其他角色是否有缺失权限
   - 重启后端服务

2. **验证权限配置**
   - 使用 `testauthor001` 账号登录
   - 访问 `/api/v1/writer/projects`
   - 确认不再返回 401

### Phase 2: 权限中间件实现（2-3天）

1. **实现权限检查中间件**
   - 从 JWT 解析用户信息
   - 加载权限配置
   - 检查用户是否有所需权限

2. **集成到路由**
   - 更新 `SetupAuthMiddleware` 实现
   - 在路由中添加权限要求
   - 测试权限检查

### Phase 3: 权限系统优化（1-2周）

1. **实现统一的权限验证器**
   - 创建 `PermissionChecker` 类型
   - 提供权限检查 API
   - 支持权限组合

2. **优化权限模型**
   - 明确角色和权限关系
   - 支持角色继承
   - 支持动态权限

3. **添加权限管理 API**
   - 创建权限
   - 分配权限给角色
   - 查询用户权限

---

## 安全检查清单

### 权限配置
- [ ] 所有角色都有明确的权限定义
- [ ] 敏感操作有权限保护
- [ ] 超级管理员权限正确配置

### JWT 认证
- [ ] JWT secret 一致性检查
- [ ] Token 过期时间合理
- [ ] Token 刷新机制正常

### 权限检查
- [ ] 所有 API 端点都有权限检查
- [ ] 权限检查在中间件层统一处理
- [ ] 权限错误返回正确的 HTTP 状态码

### 测试验证
- [ ] 测试各角色的访问权限
- [ ] 测试权限边界情况
- [ ] 测试未认证访问

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [401认证错误诊断报告](../reports/archived/2026-02-08-auth-401-issue-diagnosis.md) | 详细问题分析 |
| [JWT 配置](../../config/config.yaml) | JWT 密钥配置 |
| [权限配置](../../config/permissions.yaml) | 角色权限定义 |
| [认证中间件](../../internal/router/setup.go) | 中间件实现 |

---

## 相关 Issue

- [#008: 中间件架构问题](./008-middleware-architecture-issues.md)
- [#005: API 标准化问题](./005-api-standardization-issues.md)
