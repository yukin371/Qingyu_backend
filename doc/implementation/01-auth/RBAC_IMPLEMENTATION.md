# RBAC 权限控制实施文档

## 概述

本文档记录了青羽平台基于角色的访问控制 (RBAC) 系统的实施过程。

**实施日期**: 2025-12-20
**状态**: ✅ 已完成

## 什么是 RBAC

RBAC (Role-Based Access Control) 是一种基于角色的访问控制机制，通过将权限分配给角色，再将角色分配给用户，实现权限管理。

### 核心概念

- **用户 (User)**: 系统的使用者
- **角色 (Role)**: 权限的集合
- **权限 (Permission)**: 对资源的操作许可
- **资源 (Resource)**: 系统中的实体（API、数据等）

## 架构设计

### RBAC 模型

```
┌─────────┐         ┌─────────┐         ┌────────────┐
│  User   │ ──────> │  Role   │ ──────> │ Permission │
└─────────┘         └─────────┘         └────────────┘
                          │
                          ↓
                   ┌──────────────┐
                   │    Resource  │
                   └──────────────┘
```

### 预定义角色

| 角色 | 说明 | 权限范围 |
|------|------|----------|
| `admin` | 系统管理员 | 所有权限 |
| `author` | 作者 | 写作、发布、作品管理 |
| `reader` | 读者 | 阅读、评论、收藏 |
| `vip` | VIP会员 | 读者权限 + 专属内容 |

### 权限矩阵

| 操作 | admin | author | reader | vip |
|------|-------|--------|--------|-----|
| 用户管理 | ✅ | ❌ | ❌ | ❌ |
| 内容审核 | ✅ | ❌ | ❌ | ❌ |
| 写作/发布 | ✅ | ✅ | ❌ | ❌ |
| 阅读 | ✅ | ✅ | ✅ | ✅ |
| 评论 | ✅ | ✅ | ✅ | ✅ |
| 专属内容 | ✅ | ❌ | ❌ | ✅ |
| 数据查看 | ✅ | 自己的 | 自己的 | 自己的 |

## 实施步骤

**提交**: `11ad486 feat(rbac): 实现RBAC权限控制系统`

### 步骤 1: 定义权限模型

```go
// models/rbac/permission.go

type Permission struct {
    ID          string
    Name        string
    Description string
    Resource    string
    Action      string
}

type Role struct {
    ID          string
    Name        string
    Description string
    Permissions []string
}

type UserRole struct {
    UserID    string
    RoleID    string
    CreatedAt time.Time
}
```

### 步骤 2: 实现权限检查中间件

```go
// middleware/rbac.go

func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取当前用户
        userID := c.GetString("user_id")
        if userID == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "code":    401,
                "message": "Unauthorized",
            })
            c.Abort()
            return
        }

        // 获取用户角色
        userRoles, err := getUserRoles(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "code":    500,
                "message": "Failed to get user roles",
            })
            c.Abort()
            return
        }

        // 检查角色
        hasRole := false
        for _, userRole := range userRoles {
            for _, requiredRole := range roles {
                if userRole == requiredRole {
                    hasRole = true
                    break
                }
            }
            if hasRole {
                break
            }
        }

        if !hasRole {
            c.JSON(http.StatusForbidden, gin.H{
                "code":    403,
                "message": "Forbidden: Insufficient permissions",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

func RequirePermission(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")

        // 检查用户是否有指定权限
        hasPermission, err := checkPermission(userID, resource, action)
        if err != nil || !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{
                "code":    403,
                "message": "Forbidden: Missing required permission",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 步骤 3: 创建权限服务

```go
// service/rbac/permission_service.go

type PermissionService struct {
    permissionRepo repository.PermissionRepository
    roleRepo       repository.RoleRepository
    userRoleRepo   repository.UserRoleRepository
}

func (s *PermissionService) CheckPermission(userID, resource, action string) (bool, error) {
    // 获取用户角色
    userRoles, err := s.userRoleRepo.FindByUserID(userID)
    if err != nil {
        return false, err
    }

    // 检查每个角色的权限
    for _, userRole := range userRoles {
        role, err := s.roleRepo.FindByID(userRole.RoleID)
        if err != nil {
            continue
        }

        for _, permissionID := range role.Permissions {
            permission, err := s.permissionRepo.FindByID(permissionID)
            if err != nil {
                continue
            }

            if permission.Resource == resource && permission.Action == action {
                return true, nil
            }
        }
    }

    return false, nil
}

func (s *PermissionService) AssignRole(userID, roleID string) error {
    userRole := &rbac.UserRole{
        UserID:    userID,
        RoleID:    roleID,
        CreatedAt: time.Now(),
    }

    return s.userRoleRepo.Create(userRole)
}

func (s *PermissionService) RevokeRole(userID, roleID string) error {
    return s.userRoleRepo.Delete(userID, roleID)
}
```

### 步骤 4: 应用到路由

```go
// router/admin/admin_router.go

func RegisterAdminRoutes(r *gin.RouterGroup, adminSvc *admin.Service) {
    adminGroup := r.Group("/admin")
    adminGroup.Use(middleware.JWTAuth())
    adminGroup.Use(middleware.RequireRole("admin"))
    {
        // 用户管理
        adminGroup.GET("/users", adminSvc.ListUsers)
        adminGroup.GET("/users/:id", adminSvc.GetUser)
        adminGroup.PUT("/users/:id", adminSvc.UpdateUser)
        adminGroup.DELETE("/users/:id", adminSvc.DeleteUser)

        // 角色管理
        adminGroup.GET("/roles", adminSvc.ListRoles)
        adminGroup.POST("/roles", adminSvc.CreateRole)
        adminGroup.PUT("/roles/:id", adminSvc.UpdateRole)
        adminGroup.DELETE("/roles/:id", adminSvc.DeleteRole)

        // 权限管理
        adminGroup.GET("/permissions", adminSvc.ListPermissions)
        adminGroup.POST("/permissions", adminSvc.CreatePermission)
    }
}
```

## API 端点

### 角色管理 (管理员)

#### 获取角色列表
```
GET /api/v1/admin/roles
Authorization: Bearer {admin_token}
```

#### 创建角色
```
POST /api/v1/admin/roles
Authorization: Bearer {admin_token}
Body: {
  "name": "moderator",
  "description": "内容审核员",
  "permissions": ["content:read", "content:review"]
}
```

#### 更新角色
```
PUT /api/v1/admin/roles/:id
Authorization: Bearer {admin_token}
```

#### 删除角色
```
DELETE /api/v1/admin/roles/:id
Authorization: Bearer {admin_token}
```

### 用户角色管理 (管理员)

#### 分配角色
```
POST /api/v1/admin/users/:id/roles
Authorization: Bearer {admin_token}
Body: {
  "role_id": "author"
}
```

#### 撤销角色
```
DELETE /api/v1/admin/users/:id/roles/:roleId
Authorization: Bearer {admin_token}
```

#### 获取用户角色
```
GET /api/v1/admin/users/:id/roles
Authorization: Bearer {admin_token}
```

### 权限检查 (内部使用)

#### 检查权限
```
GET /api/v1/rbac/check?resource=book&action=write
Authorization: Bearer {token}
```

## 数据模型

### Permission (权限)
```go
type Permission struct {
    ID          string
    Name        string  // 如: "book:write"
    Description string
    Resource    string  // 资源: "book"
    Action      string  // 操作: "write"
}
```

### Role (角色)
```go
type Role struct {
    ID          string
    Name        string  // 如: "author"
    Description string
    Permissions []string  // 权限 ID 列表
}
```

### UserRole (用户角色关联)
```go
type UserRole struct {
    UserID    string
    RoleID    string
    CreatedAt time.Time
}
```

## 使用示例

### 示例 1: 创建作者角色

```go
permissionSvc := serviceContainer.GetPermissionService()

// 创建权限
writePerm, _ := permissionSvc.CreatePermission("book:write", "写书权限", "book", "write")
publishPerm, _ := permissionSvc.CreatePermission("book:publish", "发布权限", "book", "publish")

// 创建角色
authorRole, _ := permissionSvc.CreateRole("author", "作者", []string{writePerm.ID, publishPerm.ID})

// 分配给用户
permissionSvc.AssignRole("user123", authorRole.ID)
```

### 示例 2: 检查权限

```go
hasPermission, err := permissionSvc.CheckPermission("user123", "book", "write")
if err != nil {
    // 处理错误
}

if !hasPermission {
    return errors.New("没有写书权限")
}
```

### 示例 3: 保护路由

```go
// 只允许作者访问
r.POST("/api/v1/writer/books",
    middleware.JWTAuth(),
    middleware.RequireRole("author"),
    writerAPI.CreateBook,
)

// 要求特定权限
r.POST("/api/v1/admin/settings",
    middleware.JWTAuth(),
    middleware.RequirePermission("system", "write"),
    adminAPI.UpdateSettings,
)
```

## 中间件使用

### RequireRole 中间件

```go
// 单个角色
r.GET("/admin/users", middleware.RequireRole("admin"), handler)

// 多个角色（满足任一即可）
r.DELETE("/admin/users/:id", middleware.RequireRole("admin", "moderator"), handler)
```

### RequirePermission 中间件

```go
// 需要特定权限
r.POST("/writer/books", middleware.RequirePermission("book", "write"), handler)
```

### 组合使用

```go
r.POST("/api/v1/writer/books/:id/publish",
    middleware.JWTAuth(),              // 必须登录
    middleware.RequireRole("author"),  // 必须是作者
    middleware.RequirePermission("book", "publish"),  // 必须有发布权限
    writerAPI.PublishBook,
)
```

## 最佳实践

### 1. 权限命名规范

```
{resource}:{action}
```

示例:
- `book:read` - 阅读书籍
- `book:write` - 写作
- `book:publish` - 发布
- `book:delete` - 删除
- `user:manage` - 用户管理
- `system:config` - 系统配置

### 2. 角色设计原则

- **最小权限原则**: 角色只包含必需的权限
- **职责分离**: 避免一个角色拥有过多权限
- **可组合性**: 通过角色组合实现复杂权限

### 3. 缓存策略

```go
// 缓存用户角色（5分钟）
func getUserRolesWithCache(userID string) ([]string, error) {
    cacheKey := fmt.Sprintf("user:roles:%s", userID)

    // 尝试从缓存获取
    if roles, found := cache.Get(cacheKey); found {
        return roles.([]string), nil
    }

    // 从数据库获取
    roles, err := getUserRoles(userID)
    if err != nil {
        return nil, err
    }

    // 存入缓存
    cache.Set(cacheKey, roles, 5*time.Minute)

    return roles, nil
}
```

### 4. 审计日志

记录所有权限相关的操作：

```go
func (s *PermissionService) AssignRole(userID, roleID string) error {
    err := s.userRoleRepo.Create(&UserRole{
        UserID: userID,
        RoleID: roleID,
    })

    // 记录审计日志
    s.auditLogger.Log(AuditEvent{
        Action:    "role.assign",
        UserID:    userID,
        RoleID:    roleID,
        Timestamp: time.Now(),
    })

    return err
}
```

## 安全考虑

### 1. 防止权限提升

- 普通用户不能提升自己的权限
- 修改权限需要二次确认
- 记录所有权限变更

### 2. 默认拒绝

```go
// 默认拒绝访问
func checkPermission(userID, resource, action string) (bool, error) {
    // 明确授权才允许
    // 未授权的默认拒绝
}
```

### 3. 定期审查

- 定期审查用户权限
- 删除未使用的角色
- 撤销不活跃用户的高级权限

## 测试

### 单元测试

```go
func TestCheckPermission(t *testing.T) {
    service := setupTestService()

    // 测试有权限的用户
    hasPermission, _ := service.CheckPermission("user1", "book", "write")
    assert.True(t, hasPermission)

    // 测试无权限的用户
    hasPermission, _ = service.CheckPermission("user2", "book", "write")
    assert.False(t, hasPermission)
}
```

### 集成测试

```go
func TestRoleBasedAccess(t *testing.T) {
    // 测试管理员可以访问
    w := performRequest(router, "GET", "/api/v1/admin/users", adminToken)
    assert.Equal(t, 200, w.Code)

    // 测试普通用户不能访问
    w = performRequest(router, "GET", "/api/v1/admin/users", userToken)
    assert.Equal(t, 403, w.Code)
}
```

## 后续优化

1. **动态权限**: 支持基于资源的动态权限
2. **权限继承**: 支持角色继承
3. **临时权限**: 支持有时效的权限授予
4. **权限委托**: 支持权限委托给其他用户
5. **批量操作**: 支持批量权限管理

## 相关文档

- [P0 中间件集成](MIDDLEWARE_INTEGRATION.md)
- [路由修复报告](ROUTER_FIX_REPORT.md)
- [项目结构总结](../docs/项目结构总结.md)

## 提交历史

```
11ad486 - feat(rbac): 实现RBAC权限控制系统
```
