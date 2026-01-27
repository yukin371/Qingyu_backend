# 权限系统集成指南

> Qingyu Backend 权限中间件使用文档
>
> 版本: v2.0 | 更新时间: 2026-01-27

## 目录

- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [配置说明](#配置说明)
- [集成示例](#集成示例)
- [高级功能](#高级功能)
- [API参考](#api参考)
- [最佳实践](#最佳实践)
- [故障排查](#故障排查)

---

## 快速开始

### 1. 基本配置

在 `configs/middleware.yaml` 中配置权限中间件：

```yaml
middleware:
  permission:
    enabled: true
    strategy: "rbac"
    config_path: "configs/permissions.yaml"
    skip_paths:
      - "/health"
      - "/metrics"
```

### 2. 路由集成

```go
import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "Qingyu_backend/internal/middleware/auth"
    "Qingyu_backend/internal/router"
)

func SetupRoutes(router *gin.Engine, logger *zap.Logger) error {
    // 1. 设置权限中间件
    permMW, err := router.SetupPermissionMiddleware(
        "configs/permissions.yaml",
        logger,
    )
    if err != nil {
        return err
    }

    // 2. 应用到路由组
    v1 := router.Group("/api/v1")
    v1.Use(permMW.Handler())

    return nil
}
```

### 3. 运行验证

```bash
# 启动服务
go run cmd/server/main.go

# 测试权限检查
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/api/v1/books
```

---

## 核心概念

### RBAC模型

权限系统基于**角色-权限-资源**模型：

```
用户 (User)
  ↓ 拥有
角色 (Role)
  ↓ 包含
权限 (Permission)
  ↓ 定义为
资源:操作 (Resource:Action)
```

**示例**：
- 用户 "alice" 拥有角色 "author"
- 角色 "author" 包含权限 "book:read", "book:write"
- 权限 "book:write" 定义为对书籍资源的写操作

### 权限格式

权限字符串格式：`<resource>:<action>[:<resource_id>]`

| 格式 | 示例 | 说明 |
|------|------|------|
| `resource:action` | `book:read` | 读取书籍 |
| `*:*` | `*:*` | 所有权限（管理员） |
| `resource:*` | `book:*` | 书籍所有操作 |
| `*:action` | `*:read` | 所有资源的读取 |
| `resource:action:id` | `book:update:123` | 更新特定书籍 |

### 中间件优先级

权限中间件优先级为 `10`，执行顺序：

```
RequestID(1) → Recovery(2) → ErrorHandler(3) → Security(4)
→ CORS(5) → Auth(8) → Permission(10) → Logger(11)
```

---

## 配置说明

### middleware.yaml 完整配置

```yaml
middleware:
  permission:
    # 基础配置
    enabled: ${PERMISSION_ENABLED:true}
    strategy: "rbac"  # rbac | casbin | noop
    config_path: "configs/permissions.yaml"

    # 跳过权限检查的路径
    skip_paths:
      - "/health"
      - "/metrics"
      - "/api/v1/auth/login"
      - "/api/v1/auth/register"
      - "/api/v1/public"

    # 权限不足时的响应
    message: "权限不足，无法访问该资源"
    status_code: 403
    error_code: 40301

    # RBAC特定配置
    rbac:
      # 从数据库加载权限
      load_from_db: ${PERMISSION_LOAD_FROM_DB:true}

      # 缓存配置
      cache:
        enabled: ${PERMISSION_CACHE_ENABLED:true}
        ttl: 5m  # 缓存过期时间
        refresh_interval: 1m  # 自动刷新间隔

      # 权限格式转换
      format_conversion:
        enabled: true  # user.read -> user:read

    # 会话超时
    session_timeout: 30m

    # 热更新
    hot_reload:
      enabled: ${PERMISSION_HOT_RELOAD:false}
      watch_config: true
      watch_db: true
```

### permissions.yaml 权限定义

```yaml
# 角色定义
roles:
  admin:
    name: "管理员"
    description: "系统管理员，拥有所有权限"
    permissions:
      - "*:*"

  author:
    name: "作者"
    description: "可以创作和管理自己的作品"
    permissions:
      - "book:read"
      - "book:create"
      - "book:update"
      - "book:delete"
      - "chapter:read"
      - "chapter:create"
      - "chapter:update"

  reader:
    name: "读者"
    description: "可以阅读作品"
    permissions:
      - "book:read"
      - "chapter:read"

# 用户角色分配
user_roles:
  user_123:
    - admin
  user_456:
    - author
  user_789:
    - reader

# 角色继承
role_hierarchy:
  author:
    inherits:
      - reader  # 作者继承读者权限
```

---

## 集成示例

### 示例1：基本路由保护

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "Qingyu_backend/internal/router"
)

func main() {
    router := gin.Default()
    logger := zap.NewExample()

    // 设置路由（包含权限中间件）
    if err := router.SetupRoutes(router, logger); err != nil {
        logger.Fatal("Failed to setup routes", zap.Error(err))
    }

    router.Run(":8080")
}
```

### 示例2：手动权限检查

```go
package handlers

import (
    "github.com/gin-gonic/gin"
    "Qingyu_backend/internal/middleware/auth"
    "Qingyu_backend/service/shared/auth"
)

type BookHandler struct {
    permService auth.PermissionService
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
    userID := c.GetString("user_id")
    bookID := c.Param("id")

    // 手动检查权限
    allowed, err := h.permService.CheckPermission(
        c.Request.Context(),
        userID,
        "book:update",
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    if !allowed {
        c.JSON(403, gin.H{"error": "权限不足"})
        return
    }

    // 执行业务逻辑
    // ...
}
```

### 示例3：数据库权限集成

```go
package main

import (
    "context"
    "go.uber.org/zap"

    "Qingyu_backend/internal/middleware/auth"
    "Qingyu_backend/repository/mongodb/auth"
    authService "Qingyu_backend/service/shared/auth"
)

func setupPermissionService() {
    logger := zap.NewExample()

    // 1. 创建仓储
    roleRepo := mongodb_auth.NewRoleRepository(mongoClient)

    // 2. 创建权限服务
    permService := authService.NewPermissionService(roleRepo, nil, logger)

    // 3. 创建RBAC检查器
    checker, _ := auth.NewRBACChecker(nil)

    // 4. 设置检查器到服务
    permService.SetChecker(checker)

    // 5. 从数据库加载权限
    if err := permService.LoadPermissionsToChecker(context.Background()); err != nil {
        logger.Error("Failed to load permissions", zap.Error(err))
    }

    // 6. 使用checker创建权限中间件
    permConfig := &auth.PermissionConfig{
        Enabled:  true,
        Strategy: "rbac",
        Checker:  checker,
    }
    permMW, _ := auth.NewPermissionMiddleware(permConfig, logger)

    // 7. 应用到路由
    router.Use(permMW.Handler())
}
```

### 示例4：动态权限加载

```go
package api

import (
    "github.com/gin-gonic/gin"
    "Qingyu_backend/service/shared/auth"
)

type PermissionAPI struct {
    permService auth.PermissionService
}

// ReloadPermissions 重新加载权限
func (api *PermissionAPI) ReloadPermissions(c *gin.Context) {
    // 从数据库重新加载所有权限
    if err := api.permService.ReloadAllFromDatabase(c.Request.Context()); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "权限已重新加载"})
}

// InvalidateUserCache 清除用户权限缓存
func (api *PermissionAPI) InvalidateUserCache(c *gin.Context) {
    userID := c.Param("userId")

    if err := api.permService.InvalidateUserPermissionsCache(
        c.Request.Context(),
        userID,
    ); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "缓存已清除"})
}
```

---

## 高级功能

### 1. 权限继承

```yaml
# permissions.yaml
role_hierarchy:
  admin:
    inherits: []  # 管理员不继承任何人

  editor:
    inherits:
      - author  # 编辑继承作者权限

  author:
    inherits:
      - reader  # 作者继承读者权限

  reader:
    inherits: []  # 读者是最底层角色
```

### 2. 细粒度权限控制

```go
// 只能删除自己创建的文档
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
    userID := c.GetString("user_id")
    docID := c.Param("id")

    // 1. 检查基本权限
    allowed, _ := h.permService.CheckPermission(
        c.Request.Context(),
        userID,
        "document:delete",
    )
    if !allowed {
        c.JSON(403, gin.H{"error": "权限不足"})
        return
    }

    // 2. 检查资源所有权
    doc, err := h.docService.Get(c.Request.Context(), docID)
    if err != nil {
        c.JSON(404, gin.H{"error": "文档不存在"})
        return
    }

    if doc.CreatedBy != userID {
        // 3. 检查是否有管理员权限
        hasAdmin, _ := h.permService.CheckPermission(
            c.Request.Context(),
            userID,
            "document:delete:*",
        )
        if !hasAdmin {
            c.JSON(403, gin.H{"error": "只能删除自己创建的文档"})
            return
        }
    }

    // 执行删除
    h.docService.Delete(c.Request.Context(), docID)
}
```

### 3. 条件权限

```go
// 基于时间的权限检查
func (h *PromotionHandler) ViewPromotion(c *gin.Context) {
    userID := c.GetString("user_id")

    // 检查会员权限
    hasMembership, _ := h.permService.CheckPermission(
        c.Request.Context(),
        userID,
        "promotion:view:vip",
    )

    if !hasMembership {
        // 检查是否在活动期间
        promotion, _ := h.promoService.GetCurrent()
        if promotion == nil || !promotion.IsActive() {
            c.JSON(403, gin.H{"error": "活动未开始或已结束"})
            return
        }
    }

    // 返回促销信息
}
```

### 4. 批量权限检查

```go
// 一次性检查多个权限
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
    userID := c.GetString("user_id")

    permissions := []auth.Permission{
        {Resource: "book", Action: "read"},
        {Resource: "user", Action: "read"},
        {Resource: "stats", Action: "read"},
    }

    results, err := h.checker.BatchCheck(
        c.Request.Context(),
        userID,
        permissions,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    dashboard := gin.H{}
    if results[0] { // book:read
        dashboard["books"] = h.bookService.GetRecent()
    }
    if results[1] { // user:read
        dashboard["users"] = h.userService.GetStats()
    }
    if results[2] { // stats:read
        dashboard["stats"] = h.statsService.GetOverview()
    }

    c.JSON(200, dashboard)
}
```

---

## API参考

### PermissionService接口

```go
type PermissionService interface {
    // 权限检查
    CheckPermission(ctx context.Context, userID, permission string) (bool, error)
    GetUserPermissions(ctx context.Context, userID string) ([]string, error)
    GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
    HasRole(ctx context.Context, userID, role string) (bool, error)

    // 缓存管理
    InvalidateUserPermissionsCache(ctx context.Context, userID string) error

    // RBAC集成
    SetChecker(checker interface{})
    LoadPermissionsToChecker(ctx context.Context) error
    LoadUserRolesToChecker(ctx context.Context, userID string) error
    ReloadAllFromDatabase(ctx context.Context) error
}
```

### RBACChecker接口

```go
type RBACChecker struct {
    // 角色管理
    AssignRole(userID, roleName string)
    RevokeRole(userID, roleName string)

    // 权限管理
    GrantPermission(roleName, permission string)
    BatchGrantPermissions(roleName string, permissions []string)
    RevokePermission(roleName, permission string)

    // 权限检查
    Check(ctx context.Context, userID string, perm Permission) (bool, error)
    BatchCheck(ctx context.Context, userID string, perms []Permission) ([]bool, error)

    // 查询方法
    GetUserRoles(userID string) []string
    GetRolePermissions(roleName string) []string
    HasRole(userID, roleName string) bool
    HasAnyRole(userID string, roleNames []string) bool
    HasAllRoles(userID string, roleNames []string) bool

    // 配置加载
    LoadFromMap(config map[string]interface{}) error
    LoadFromYAML(path string) error

    // 统计信息
    Stats() map[string]interface{}

    // 关闭
    Close() error
}
```

### PermissionMiddleware

```go
type PermissionMiddleware struct {
    // 中间件接口
    Name() string
    Priority() int
    Handler() gin.HandlerFunc

    // 可配置中间件接口
    LoadConfig(config map[string]interface{}) error
    ValidateConfig() error
    Reload(config map[string]interface{}) error

    // 热更新中间件接口
    EnableHotReload() error
    DisableHotReload() error
}

// 创建权限中间件
func NewPermissionMiddleware(
    config *PermissionConfig,
    logger *zap.Logger,
) (*PermissionMiddleware, error)

// 权限配置
type PermissionConfig struct {
    Enabled    bool
    Strategy   string  // rbac | casbin | noop
    ConfigPath string
    SkipPaths  []string
    Message    string
    StatusCode int
}
```

---

## 最佳实践

### 1. 权限设计原则

**最小权限原则**：
- 默认拒绝所有访问
- 只授予必要的权限
- 定期审查和清理权限

**职责分离**：
- 创建角色时考虑职责边界
- 避免角色权限过于宽泛
- 使用多个角色组合而非大角色

**示例**：
```yaml
# ❌ 不推荐：超级角色
super_admin:
  permissions:
    - "*:*"

# ✅ 推荐：细分角色
content_admin:
  permissions:
    - "book:*"
    - "chapter:*"

system_admin:
  permissions:
    - "user:*"
    - "system:*"
```

### 2. 性能优化

**启用缓存**：
```yaml
rbac:
  cache:
    enabled: true
    ttl: 5m
    refresh_interval: 1m
```

**批量检查**：
```go
// ❌ 不推荐：多次单独检查
for _, resource := range resources {
    allowed, _ := checker.Check(ctx, userID, resource)
    if !allowed {
        return
    }
}

// ✅ 推荐：批量检查
results, _ := checker.BatchCheck(ctx, userID, permissions)
for i, allowed := range results {
    if !allowed {
        return
    }
}
```

**跳过不必要的路径**：
```yaml
skip_paths:
  - "/health"      # 健康检查
  - "/metrics"     # 监控指标
  - "/public"      # 公开资源
  - "/api/v1/auth" # 认证端点
```

### 3. 错误处理

**明确的错误信息**：
```go
// ❌ 不推荐：模糊错误
if !allowed {
    return errors.New("access denied")
}

// ✅ 推荐：详细错误
if !allowed {
    return errors.New("需要 book:write 权限才能执行此操作")
}
```

**日志记录**：
```go
func (m *PermissionMiddleware) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        perm := getPermissionFromContext(c)

        allowed, err := m.checker.Check(c, userID, perm)
        if !allowed {
            // 记录拒绝的权限检查
            m.logger.Warn("Permission denied",
                zap.String("user", userID),
                zap.String("permission", perm.String()),
                zap.String("path", c.Request.URL.Path),
                zap.Error(err),
            )

            c.JSON(403, gin.H{
                "code":    40301,
                "message": m.config.Message,
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 4. 测试策略

**单元测试**：
```go
func TestPermissionChecker(t *testing.T) {
    checker := NewTestChecker()

    // 设置测试数据
    checker.AssignRole("user1", "author")
    checker.GrantPermission("author", "book:write")

    // 测试权限检查
    allowed, _ := checker.Check(
        context.Background(),
        "user1",
        Permission{Resource: "book", Action: "write"},
    )

    assert.True(t, allowed)
}
```

**集成测试**：
```go
func TestPermissionMiddleware(t *testing.T) {
    router := setupTestRouter()

    // 测试有权限的请求
    req := httptest.NewRequest("GET", "/api/v1/books", nil)
    req.Header.Set("Authorization", "Bearer valid_token")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)

    // 测试无权限的请求
    req = httptest.NewRequest("DELETE", "/api/v1/books/123", nil)
    req.Header.Set("Authorization", "Bearer reader_token")
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 403, w.Code)
}
```

---

## 故障排查

### 问题1：权限检查总是返回false

**可能原因**：
1. 用户未分配角色
2. 角色未分配权限
3. 权限格式不匹配
4. 配置文件未加载

**排查步骤**：
```go
// 1. 检查用户角色
roles, _ := permService.GetUserRoles(ctx, userID)
fmt.Printf("User roles: %v\n", roles) // 应该非空

// 2. 检查角色权限
for _, role := range roles {
    perms, _ := permService.GetRolePermissions(ctx, role)
    fmt.Printf("Role %s permissions: %v\n", role, perms)
}

// 3. 检查权限格式
perm := "book:read"  // 确保使用 ":" 分隔符
allowed, _ := permService.CheckPermission(ctx, userID, perm)

// 4. 检查配置
stats := checker.Stats()
fmt.Printf("Checker stats: %v\n", stats)
```

### 问题2：缓存导致权限更新不生效

**解决方案**：
```go
// 1. 清除特定用户缓存
permService.InvalidateUserPermissionsCache(ctx, userID)

// 2. 重新加载所有权限
permService.ReloadAllFromDatabase(ctx)

// 3. 如果使用Redis，直接删除缓存键
cacheClient.Delete(ctx, fmt.Sprintf("user:permissions:%s", userID))
```

### 问题3：性能问题

**优化措施**：

1. **启用缓存**：
```yaml
rbac:
  cache:
    enabled: true
    ttl: 5m
```

2. **批量检查**：
```go
results, _ := checker.BatchCheck(ctx, userID, permissions)
```

3. **减少数据库查询**：
```go
// ❌ 每次检查都查询数据库
func CheckPermission(userID, perm string) bool {
    roles := db.GetUserRoles(userID)  // DB查询
    for _, role := range roles {
        perms := db.GetRolePermissions(role)  // DB查询
        if contains(perms, perm) {
            return true
        }
    }
    return false
}

// ✅ 使用缓存
type PermissionService struct {
    cache CacheClient
    repo  Repository
}

func (s *PermissionService) CheckPermission(ctx, userID, perm string) bool {
    // 先查缓存
    if cached, _ := s.cache.Get(ctx, userID); cached != "" {
        return contains(strings.Split(cached, ","), perm)
    }

    // 缓存未命中，查数据库
    perms, _ := s.repo.GetUserPermissions(ctx, userID)
    s.cache.Set(ctx, userID, strings.Join(perms, ","), 5*time.Minute)

    return contains(perms, perm)
}
```

### 问题4：权限格式不一致

**统一格式**：
```go
// 格式化权限字符串
func normalizePermission(perm string) string {
    // 转换 . 到 :
    perm = strings.ReplaceAll(perm, ".", ":")

    // 转换 :: 到 :
    perm = strings.ReplaceAll(perm, "::", ":")

    // 转换为小写
    perm = strings.ToLower(perm)

    // 去除空格
    perm = strings.TrimSpace(perm)

    return perm
}

// 使用示例
perm := normalizePermission("Book . Read")  // "book:read"
```

---

## 总结

本指南涵盖了Qingyu Backend权限系统的完整集成流程：

- ✅ **快速配置**：5分钟即可完成基本集成
- ✅ **灵活扩展**：支持RBAC、Casbin等多种策略
- ✅ **高性能**：内置缓存和批量检查
- ✅ **易维护**：清晰的接口和配置结构

**下一步**：
1. 根据实际业务需求配置角色和权限
2. 编写单元测试和集成测试
3. 监控权限检查性能和日志
4. 定期审查和优化权限设计

**相关文档**：
- [中间件架构设计](./architecture.md)
- [权限模型定义](../../models/auth/role.go)
- [API文档](../../api/README.md)
