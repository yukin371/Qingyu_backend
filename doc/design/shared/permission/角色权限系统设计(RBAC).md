# 角色权限系统设计（RBAC）

> **版本**: v1.0  
> **创建日期**: 2025-10-21  
> **最后更新**: 2025-10-21  
> **状态**: ⚠️ 部分实现，补充设计文档

---

## 1. 为什么需要这个设计

### 1.1 业务背景

角色权限系统（RBAC - Role-Based Access Control）是平台安全的核心组件，负责：
- 用户角色管理（普通用户、作者、VIP、管理员等）
- 权限细粒度控制（操作级别、资源级别）
- 灵活的权限分配和回收
- 权限检查和验证
- 支持多角色和角色继承

### 1.2 实际实现情况

当前已实现部分功能：
- **基础Role模型**: `models/shared/auth/role.go`, `models/users/role.go`
- **权限中间件**: `middleware/permission_middleware.go`
- **RoleService**: `service/shared/auth/role_service.go`
- **PermissionService**: `service/shared/auth/permission_service.go`

**缺失部分**：
- 完整的RBAC数据模型设计
- 权限继承机制
- 资源级权限控制
- 权限缓存策略

**实现文件**：
- `models/shared/auth/role.go`
- `models/users/role.go`
- `service/shared/auth/role_service.go`
- `service/shared/auth/permission_service.go`
- `middleware/permission_middleware.go`
- `middleware/admin_permission.go`
- `middleware/vip_permission.go`

---

## 2. RBAC模型设计

### 2.1 RBAC核心概念

```
User (用户) ←→ Role (角色) ←→ Permission (权限) ←→ Resource (资源)
```

**四个核心实体**：
1. **User**: 用户
2. **Role**: 角色
3. **Permission**: 权限
4. **Resource**: 资源

**关系**：
- 一个用户可以拥有多个角色
- 一个角色包含多个权限
- 一个权限可以应用于多个资源
- 用户通过角色获得权限

---

## 3. 数据模型设计

### 3.1 Role（角色）

#### 3.1.1 数据结构

```go
type Role struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    Name        string    `bson:"name" json:"name" validate:"required"`           // 角色名称
    Code        string    `bson:"code" json:"code" validate:"required"`           // 角色代码（唯一）
    Description string    `bson:"description" json:"description"`                 // 角色描述
    Level       int       `bson:"level" json:"level"`                             // 角色等级（用于权限继承）
    ParentID    string    `bson:"parentId,omitempty" json:"parentId,omitempty"`   // 父角色ID（支持角色继承）
    Permissions []string  `bson:"permissions" json:"permissions"`                 // 权限列表（权限代码）
    IsSystem    bool      `bson:"isSystem" json:"isSystem"`                       // 是否系统角色（不可删除）
    IsActive    bool      `bson:"isActive" json:"isActive"`                       // 是否启用
    CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
    CreatedBy   string    `bson:"createdBy" json:"createdBy"`
}
```

#### 3.1.2 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ID | string | 否 | 角色唯一标识 |
| Name | string | 是 | 角色名称（如：管理员、作者、VIP用户） |
| Code | string | 是 | 角色代码（如：admin, author, vip），全局唯一 |
| Description | string | 否 | 角色描述 |
| Level | int | 否 | 角色等级，用于权限继承（数字越大权限越高） |
| ParentID | string | 否 | 父角色ID，支持角色继承 |
| Permissions | []string | 否 | 权限代码列表 |
| IsSystem | bool | 否 | 是否系统角色（系统角色不可删除） |
| IsActive | bool | 否 | 是否启用 |
| CreatedAt | time.Time | 是 | 创建时间 |
| UpdatedAt | time.Time | 是 | 更新时间 |
| CreatedBy | string | 否 | 创建人ID |

#### 3.1.3 预定义系统角色

```go
const (
    RoleCodeSuperAdmin    = "super_admin"    // 超级管理员
    RoleCodeAdmin         = "admin"          // 管理员
    RoleCodeModerator     = "moderator"      // 版主
    RoleCodeAuthor        = "author"         // 作者
    RoleCodeVIP           = "vip"            // VIP用户
    RoleCodeUser          = "user"           // 普通用户
    RoleCodeGuest         = "guest"          // 访客
)

// 系统角色定义
var SystemRoles = []Role{
    {
        Code:        RoleCodeSuperAdmin,
        Name:        "超级管理员",
        Description: "拥有所有权限",
        Level:       100,
        Permissions: []string{"*"},  // 所有权限
        IsSystem:    true,
        IsActive:    true,
    },
    {
        Code:        RoleCodeAdmin,
        Name:        "管理员",
        Description: "平台管理权限",
        Level:       90,
        ParentID:    "",  // 无父角色
        Permissions: []string{
            "user:manage",
            "content:manage",
            "audit:manage",
            "system:config",
        },
        IsSystem: true,
        IsActive: true,
    },
    {
        Code:        RoleCodeModerator,
        Name:        "版主",
        Description: "内容审核和管理",
        Level:       70,
        Permissions: []string{
            "content:review",
            "content:delete",
            "user:view",
        },
        IsSystem: true,
        IsActive: true,
    },
    {
        Code:        RoleCodeAuthor,
        Name:        "作者",
        Description: "内容创作者",
        Level:       50,
        Permissions: []string{
            "content:create",
            "content:update",
            "content:publish",
            "stats:view",
        },
        IsSystem: true,
        IsActive: true,
    },
    {
        Code:        RoleCodeVIP,
        Name:        "VIP用户",
        Description: "付费会员",
        Level:       40,
        Permissions: []string{
            "chapter:unlock",
            "book:favorite",
            "reading:ad_free",
        },
        IsSystem: true,
        IsActive: true,
    },
    {
        Code:        RoleCodeUser,
        Name:        "普通用户",
        Description: "注册用户",
        Level:       20,
        Permissions: []string{
            "book:read",
            "book:comment",
            "book:favorite",
        },
        IsSystem: true,
        IsActive: true,
    },
    {
        Code:        RoleCodeGuest,
        Name:        "访客",
        Description: "未登录用户",
        Level:       10,
        Permissions: []string{
            "book:browse",
            "book:search",
        },
        IsSystem: true,
        IsActive: true,
    },
}
```

#### 3.1.4 索引策略

```javascript
// MongoDB索引
db.roles.createIndex({ "code": 1 }, { unique: true })   // 角色代码唯一
db.roles.createIndex({ "level": -1 })                    // 按等级查询
db.roles.createIndex({ "parentId": 1 })                  // 查询子角色
db.roles.createIndex({ "isSystem": 1, "isActive": 1 })   // 查询系统角色
```

---

### 3.2 Permission（权限）

#### 3.2.1 数据结构

```go
type Permission struct {
    ID          string    `bson:"_id,omitempty" json:"id"`
    Code        string    `bson:"code" json:"code" validate:"required"`           // 权限代码（唯一）
    Name        string    `bson:"name" json:"name" validate:"required"`           // 权限名称
    Description string    `bson:"description" json:"description"`                 // 权限描述
    Resource    string    `bson:"resource" json:"resource"`                       // 资源类型
    Action      string    `bson:"action" json:"action"`                           // 操作类型
    Category    string    `bson:"category" json:"category"`                       // 权限分类
    IsSystem    bool      `bson:"isSystem" json:"isSystem"`                       // 是否系统权限
    CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}
```

#### 3.2.2 权限命名规范

```
格式: <resource>:<action>:<scope>

resource: 资源类型（user, book, chapter, comment等）
action:   操作类型（create, read, update, delete, manage等）
scope:    作用域（own, all, system等，可选）

示例:
- user:create:all       - 创建所有用户
- book:update:own       - 更新自己的书籍
- chapter:delete:all    - 删除所有章节
- comment:read:all      - 查看所有评论
- system:config:all     - 系统配置管理
```

#### 3.2.3 预定义权限列表

```go
// 用户权限
const (
    PermUserCreate        = "user:create:all"         // 创建用户
    PermUserRead          = "user:read:all"           // 查看用户
    PermUserUpdate        = "user:update:all"         // 更新用户
    PermUserDelete        = "user:delete:all"         // 删除用户
    PermUserManage        = "user:manage:all"         // 管理用户（包含所有用户操作）
    PermUserUpdateOwn     = "user:update:own"         // 更新自己的信息
)

// 内容权限
const (
    PermContentCreate     = "content:create:own"      // 创建内容
    PermContentRead       = "content:read:all"        // 阅读内容
    PermContentUpdate     = "content:update:own"      // 更新自己的内容
    PermContentDelete     = "content:delete:own"      // 删除自己的内容
    PermContentPublish    = "content:publish:own"     // 发布内容
    PermContentReview     = "content:review:all"      // 审核内容
    PermContentManage     = "content:manage:all"      // 管理所有内容
)

// 书籍权限
const (
    PermBookCreate        = "book:create:own"         // 创建书籍
    PermBookRead          = "book:read:all"           // 阅读书籍
    PermBookUpdate        = "book:update:own"         // 更新自己的书籍
    PermBookDelete        = "book:delete:own"         // 删除自己的书籍
    PermBookFavorite      = "book:favorite:own"       // 收藏书籍
    PermBookComment       = "book:comment:all"        // 评论书籍
)

// 章节权限
const (
    PermChapterCreate     = "chapter:create:own"      // 创建章节
    PermChapterRead       = "chapter:read:all"        // 阅读章节
    PermChapterUpdate     = "chapter:update:own"      // 更新自己的章节
    PermChapterDelete     = "chapter:delete:own"      // 删除自己的章节
    PermChapterUnlock     = "chapter:unlock:all"      // 解锁付费章节（VIP）
)

// 审核权限
const (
    PermAuditView         = "audit:view:all"          // 查看审核记录
    PermAuditReview       = "audit:review:all"        // 进行人工审核
    PermAuditManage       = "audit:manage:all"        // 管理审核系统
)

// 统计权限
const (
    PermStatsView         = "stats:view:own"          // 查看自己的统计
    PermStatsViewAll      = "stats:view:all"          // 查看所有统计
)

// 系统权限
const (
    PermSystemConfig      = "system:config:all"       // 系统配置
    PermSystemLogs        = "system:logs:view"        // 查看系统日志
    PermSystemBackup      = "system:backup:all"       // 系统备份
)

// 特殊权限
const (
    PermAllAccess         = "*"                       // 所有权限（超级管理员）
    PermReadingAdFree     = "reading:ad_free:all"     // 无广告阅读（VIP）
)
```

#### 3.2.4 权限分类

```go
const (
    CategoryUser    = "user"      // 用户管理
    CategoryContent = "content"   // 内容管理
    CategoryBook    = "book"      // 书籍管理
    CategoryChapter = "chapter"   // 章节管理
    CategoryAudit   = "audit"     // 审核管理
    CategoryStats   = "stats"     // 统计管理
    CategorySystem  = "system"    // 系统管理
    CategorySpecial = "special"   // 特殊权限
)
```

---

### 3.3 UserRole（用户角色关联）

#### 3.3.1 数据结构

```go
type UserRole struct {
    ID        string     `bson:"_id,omitempty" json:"id"`
    UserID    string     `bson:"userId" json:"userId" validate:"required"`
    RoleID    string     `bson:"roleId" json:"roleId" validate:"required"`
    RoleCode  string     `bson:"roleCode" json:"roleCode"`                     // 冗余字段，便于查询
    ExpiresAt *time.Time `bson:"expiresAt,omitempty" json:"expiresAt"`         // 过期时间（VIP等）
    IsActive  bool       `bson:"isActive" json:"isActive"`
    CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`
    GrantedBy string     `bson:"grantedBy" json:"grantedBy"`                   // 授予人ID
}
```

#### 3.3.2 业务方法

```go
// IsExpired 是否已过期
func (ur *UserRole) IsExpired() bool {
    if ur.ExpiresAt == nil {
        return false  // 永久角色
    }
    return time.Now().After(*ur.ExpiresAt)
}

// IsValid 是否有效
func (ur *UserRole) IsValid() bool {
    return ur.IsActive && !ur.IsExpired()
}
```

#### 3.3.3 索引策略

```javascript
// MongoDB索引
db.user_roles.createIndex({ "userId": 1, "roleId": 1 }, { unique: true })  // 用户角色唯一
db.user_roles.createIndex({ "userId": 1, "isActive": 1 })                   // 查询用户的有效角色
db.user_roles.createIndex({ "roleCode": 1 })                                // 按角色代码查询
db.user_roles.createIndex({ "expiresAt": 1 })                               // 过期检查
```

---

### 3.4 ResourcePermission（资源权限）

用于更细粒度的权限控制。

#### 3.4.1 数据结构

```go
type ResourcePermission struct {
    ID           string    `bson:"_id,omitempty" json:"id"`
    ResourceType string    `bson:"resourceType" json:"resourceType" validate:"required"`  // 资源类型
    ResourceID   string    `bson:"resourceId" json:"resourceId" validate:"required"`      // 资源ID
    UserID       string    `bson:"userId,omitempty" json:"userId,omitempty"`              // 用户ID（可选）
    RoleID       string    `bson:"roleId,omitempty" json:"roleId,omitempty"`              // 角色ID（可选）
    Permissions  []string  `bson:"permissions" json:"permissions"`                        // 权限列表
    IsAllow      bool      `bson:"isAllow" json:"isAllow"`                                // 允许/拒绝
    CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
    UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`
}
```

#### 3.4.2 使用场景

```go
// 示例：授予用户对特定项目的编辑权限
resourcePerm := &ResourcePermission{
    ResourceType: "project",
    ResourceID:   "project_123",
    UserID:       "user_456",
    Permissions:  []string{"project:update", "project:delete"},
    IsAllow:      true,
}
```

---

## 4. 权限检查逻辑

### 4.1 权限检查流程

```
1. 获取用户的所有角色（包括已过期检查）
   ↓
2. 收集所有角色的权限（包括继承的权限）
   ↓
3. 检查是否有特殊权限（*）
   ↓
4. 检查是否有目标权限
   ↓
5. 检查资源级权限（如果需要）
   ↓
6. 返回检查结果
```

### 4.2 实现示例

```go
// HasPermission 检查用户是否拥有指定权限
func (s *PermissionService) HasPermission(ctx context.Context, userID string, permission string) (bool, error) {
    // 1. 获取用户的所有有效角色
    userRoles, err := s.GetUserRoles(ctx, userID)
    if err != nil {
        return false, err
    }

    // 2. 收集所有权限
    permissions := make(map[string]bool)
    for _, userRole := range userRoles {
        if !userRole.IsValid() {
            continue
        }

        role, err := s.roleRepo.GetByID(ctx, userRole.RoleID)
        if err != nil || role == nil {
            continue
        }

        // 添加角色的所有权限
        for _, perm := range role.Permissions {
            permissions[perm] = true
        }

        // 递归获取父角色的权限
        parentPerms, _ := s.getParentRolePermissions(ctx, role.ParentID)
        for _, perm := range parentPerms {
            permissions[perm] = true
        }
    }

    // 3. 检查特殊权限（*）
    if permissions["*"] {
        return true, nil
    }

    // 4. 检查目标权限
    if permissions[permission] {
        return true, nil
    }

    // 5. 通配符匹配（resource:* 或 resource:action:*）
    if s.matchWildcardPermission(permissions, permission) {
        return true, nil
    }

    return false, nil
}

// matchWildcardPermission 通配符权限匹配
func (s *PermissionService) matchWildcardPermission(userPerms map[string]bool, targetPerm string) bool {
    parts := strings.Split(targetPerm, ":")
    if len(parts) == 0 {
        return false
    }

    // 检查 resource:*
    if userPerms[parts[0]+":*"] {
        return true
    }

    // 检查 resource:action:*
    if len(parts) >= 2 && userPerms[parts[0]+":"+parts[1]+":*"] {
        return true
    }

    return false
}
```

### 4.3 资源级权限检查

```go
// HasResourcePermission 检查用户对特定资源的权限
func (s *PermissionService) HasResourcePermission(
    ctx context.Context,
    userID string,
    resourceType string,
    resourceID string,
    permission string,
) (bool, error) {
    // 1. 先检查用户级别的全局权限
    hasGlobal, err := s.HasPermission(ctx, userID, permission)
    if err != nil {
        return false, err
    }
    if hasGlobal {
        return true, nil
    }

    // 2. 检查资源级权限
    resourcePerm, err := s.resourcePermRepo.Get(ctx, resourceType, resourceID, userID)
    if err != nil {
        return false, err
    }

    if resourcePerm != nil && resourcePerm.IsAllow {
        for _, perm := range resourcePerm.Permissions {
            if perm == permission {
                return true, nil
            }
        }
    }

    // 3. 检查所有者权限（own scope）
    if strings.HasSuffix(permission, ":own") {
        isOwner, err := s.checkOwnership(ctx, userID, resourceType, resourceID)
        if err != nil {
            return false, err
        }
        return isOwner, nil
    }

    return false, nil
}
```

---

## 5. 权限继承

### 5.1 角色继承设计

```
SuperAdmin (Level 100)
  ↓
Admin (Level 90)
  ↓
Moderator (Level 70)
  ↓
Author (Level 50)
  ↓
VIP (Level 40)
  ↓
User (Level 20)
  ↓
Guest (Level 10)
```

### 5.2 继承实现

```go
// getParentRolePermissions 递归获取父角色的所有权限
func (s *PermissionService) getParentRolePermissions(ctx context.Context, parentID string) ([]string, error) {
    if parentID == "" {
        return []string{}, nil
    }

    parentRole, err := s.roleRepo.GetByID(ctx, parentID)
    if err != nil || parentRole == nil {
        return []string{}, err
    }

    permissions := make([]string, len(parentRole.Permissions))
    copy(permissions, parentRole.Permissions)

    // 递归获取父级的父级
    if parentRole.ParentID != "" {
        grandParentPerms, _ := s.getParentRolePermissions(ctx, parentRole.ParentID)
        permissions = append(permissions, grandParentPerms...)
    }

    return permissions, nil
}
```

---

## 6. 权限缓存策略

### 6.1 用户权限缓存

```go
type PermissionCache struct {
    UserID      string
    Permissions map[string]bool
    Roles       []string
    CachedAt    time.Time
    TTL         time.Duration
}

// GetUserPermissionsWithCache 带缓存的权限获取
func (s *PermissionService) GetUserPermissionsWithCache(ctx context.Context, userID string) (map[string]bool, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:permissions:%s", userID)
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil && cached != nil {
        return cached.(map[string]bool), nil
    }

    // 2. 从数据库加载
    permissions, err := s.GetUserPermissions(ctx, userID)
    if err != nil {
        return nil, err
    }

    // 3. 缓存结果（5分钟TTL）
    s.cache.Set(ctx, cacheKey, permissions, 5*time.Minute)

    return permissions, nil
}

// InvalidateUserPermissionCache 清除用户权限缓存
func (s *PermissionService) InvalidateUserPermissionCache(ctx context.Context, userID string) error {
    cacheKey := fmt.Sprintf("user:permissions:%s", userID)
    return s.cache.Delete(ctx, cacheKey)
}
```

### 6.2 角色权限缓存

```go
// GetRolePermissionsWithCache 带缓存的角色权限获取
func (s *PermissionService) GetRolePermissionsWithCache(ctx context.Context, roleID string) ([]string, error) {
    cacheKey := fmt.Sprintf("role:permissions:%s", roleID)
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil && cached != nil {
        return cached.([]string), nil
    }

    role, err := s.roleRepo.GetByID(ctx, roleID)
    if err != nil {
        return nil, err
    }

    permissions := role.Permissions

    // 包含继承的权限
    if role.ParentID != "" {
        parentPerms, _ := s.getParentRolePermissions(ctx, role.ParentID)
        permissions = append(permissions, parentPerms...)
    }

    // 缓存结果（30分钟TTL，角色权限变化较少）
    s.cache.Set(ctx, cacheKey, permissions, 30*time.Minute)

    return permissions, nil
}
```

---

## 7. 中间件集成

### 7.1 权限检查中间件

```go
// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(permissionService PermissionService, requiredPermission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从上下文获取用户ID
        userID, exists := c.Get("userId")
        if !exists {
            response.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
            c.Abort()
            return
        }

        // 2. 检查权限
        hasPermission, err := permissionService.HasPermission(c.Request.Context(), userID.(string), requiredPermission)
        if err != nil {
            response.Error(c, http.StatusInternalServerError, "权限检查失败", err.Error())
            c.Abort()
            return
        }

        if !hasPermission {
            response.Error(c, http.StatusForbidden, "权限不足", fmt.Sprintf("需要权限: %s", requiredPermission))
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 7.2 角色检查中间件

```go
// RoleMiddleware 角色检查中间件
func RoleMiddleware(permissionService PermissionService, requiredRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userId")
        if !exists {
            response.Error(c, http.StatusUnauthorized, "未授权", "无法获取用户信息")
            c.Abort()
            return
        }

        // 获取用户角色
        userRoles, err := permissionService.GetUserRoles(c.Request.Context(), userID.(string))
        if err != nil {
            response.Error(c, http.StatusInternalServerError, "角色检查失败", err.Error())
            c.Abort()
            return
        }

        // 检查是否拥有任一所需角色
        hasRole := false
        for _, userRole := range userRoles {
            if !userRole.IsValid() {
                continue
            }
            for _, requiredRole := range requiredRoles {
                if userRole.RoleCode == requiredRole {
                    hasRole = true
                    break
                }
            }
            if hasRole {
                break
            }
        }

        if !hasRole {
            response.Error(c, http.StatusForbidden, "权限不足", fmt.Sprintf("需要角色: %v", requiredRoles))
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 7.3 路由使用示例

```go
// 需要管理员权限
adminGroup := r.Group("/admin")
adminGroup.Use(middleware.JWTAuth())
adminGroup.Use(PermissionMiddleware(permissionService, PermUserManage))
{
    adminGroup.GET("/users", userApi.ListUsers)
    adminGroup.DELETE("/users/:id", userApi.DeleteUser)
}

// 需要作者角色
writerGroup := r.Group("/writer")
writerGroup.Use(middleware.JWTAuth())
writerGroup.Use(RoleMiddleware(permissionService, RoleCodeAuthor))
{
    writerGroup.POST("/books", bookApi.CreateBook)
    writerGroup.PUT("/books/:id", bookApi.UpdateBook)
}

// 需要VIP权限
vipChapters := r.Group("/chapters/premium")
vipChapters.Use(middleware.JWTAuth())
vipChapters.Use(PermissionMiddleware(permissionService, PermChapterUnlock))
{
    vipChapters.GET("/:id", chapterApi.ReadPremiumChapter)
}
```

---

## 8. 服务接口设计

### 8.1 RoleService接口

```go
type RoleService interface {
    // CreateRole 创建角色
    CreateRole(ctx context.Context, role *Role) error
    
    // GetRole 获取角色
    GetRole(ctx context.Context, id string) (*Role, error)
    
    // GetRoleByCode 根据代码获取角色
    GetRoleByCode(ctx context.Context, code string) (*Role, error)
    
    // UpdateRole 更新角色
    UpdateRole(ctx context.Context, id string, updates map[string]interface{}) error
    
    // DeleteRole 删除角色（非系统角色）
    DeleteRole(ctx context.Context, id string) error
    
    // ListRoles 列出所有角色
    ListRoles(ctx context.Context, filter Filter) ([]*Role, error)
    
    // AddPermissionToRole 为角色添加权限
    AddPermissionToRole(ctx context.Context, roleID string, permission string) error
    
    // RemovePermissionFromRole 从角色移除权限
    RemovePermissionFromRole(ctx context.Context, roleID string, permission string) error
}
```

### 8.2 PermissionService接口

```go
type PermissionService interface {
    // AssignRoleToUser 为用户分配角色
    AssignRoleToUser(ctx context.Context, userID string, roleID string, expiresAt *time.Time) error
    
    // RemoveRoleFromUser 移除用户角色
    RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error
    
    // GetUserRoles 获取用户的所有角色
    GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error)
    
    // GetUserPermissions 获取用户的所有权限
    GetUserPermissions(ctx context.Context, userID string) (map[string]bool, error)
    
    // HasPermission 检查用户是否拥有指定权限
    HasPermission(ctx context.Context, userID string, permission string) (bool, error)
    
    // HasAnyPermission 检查用户是否拥有任一指定权限
    HasAnyPermission(ctx context.Context, userID string, permissions []string) (bool, error)
    
    // HasAllPermissions 检查用户是否拥有所有指定权限
    HasAllPermissions(ctx context.Context, userID string, permissions []string) (bool, error)
    
    // HasRole 检查用户是否拥有指定角色
    HasRole(ctx context.Context, userID string, roleCode string) (bool, error)
    
    // HasResourcePermission 检查用户对特定资源的权限
    HasResourcePermission(ctx context.Context, userID string, resourceType string, resourceID string, permission string) (bool, error)
    
    // InvalidateUserPermissionCache 清除用户权限缓存
    InvalidateUserPermissionCache(ctx context.Context, userID string) error
}
```

---

## 9. 与v2.1架构的关系

### 9.1 在整体架构中的位置

```
Shared Module (共享模块)
  └─ Auth Module (认证授权模块)
      ├─ Models (数据模型层)
      │   ├─ Role
      │   ├─ Permission
      │   ├─ UserRole
      │   └─ ResourcePermission
      ├─ Repository (数据访问层)
      │   ├─ RoleRepository
      │   ├─ UserRoleRepository
      │   └─ ResourcePermissionRepository
      ├─ Service (业务逻辑层)
      │   ├─ RoleService
      │   └─ PermissionService
      └─ Middleware (中间件层)
          ├─ PermissionMiddleware
          └─ RoleMiddleware
```

### 9.2 事件发布

```go
// 角色分配事件
type RoleAssignedEvent struct {
    UserID    string
    RoleID    string
    RoleCode  string
    ExpiresAt *time.Time
    GrantedBy string
}

// 角色移除事件
type RoleRemovedEvent struct {
    UserID   string
    RoleID   string
    RoleCode string
    RemovedBy string
}

// 权限变更事件
type PermissionChangedEvent struct {
    RoleID      string
    RoleCode    string
    Permission  string
    Action      string  // added/removed
}
```

---

## 10. 实现参考

### 10.1 代码文件路径

- **数据模型**: 
  - `models/shared/auth/role.go`
  - `models/users/role.go`
  
- **Service层**:
  - `service/shared/auth/role_service.go`
  - `service/shared/auth/permission_service.go`
  
- **中间件**:
  - `middleware/permission_middleware.go`
  - `middleware/admin_permission.go`
  - `middleware/vip_permission.go`

### 10.2 需要补充的实现

1. ✅ 完整的Permission模型
2. ✅ UserRole关联模型
3. ✅ ResourcePermission模型
4. ✅ 权限继承逻辑
5. ✅ 权限缓存策略
6. ✅ 资源级权限检查

---

**文档状态**: ✅ 已完成  
**最后审核**: 2025-10-21  
**下一次审核**: 随代码变更同步更新

