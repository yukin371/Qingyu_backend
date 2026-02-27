# API v1 层文档

> **版本**: v1.0
> **最后更新**: 2026-02-26
> **API版本**: v1

---

## 目录

- [目录结构说明](#目录结构说明)
- [设计原则](#设计原则)
- [依赖关系](#依赖关系)
- [BaseAPI使用说明](#baseapi使用说明)
- [中间件使用说明](#中间件使用说明)
- [错误处理规范](#错误处理规范)
- [代码示例](#代码示例)
- [最佳实践](#最佳实践)

---

## 目录结构说明

### 整体结构

```
api/v1/
├── admin/              # 管理员功能模块
├── ai/                 # AI服务模块
├── announcements/      # 公告模块
├── audit/              # 审计模块
├── auth/               # 认证模块
├── bookstore/          # 书店模块
├── examples/           # 示例代码
├── finance/            # 财务模块
├── messages/           # 消息模块
├── notifications/      # 通知模块
├── reader/             # 阅读器模块 (待迁移)
├── recommendation/     # 推荐模块
├── search/             # 搜索模块
├── shared/             # 共享工具和类型
├── social/             # 社交模块
├── stats/              # 统计模块
├── system/             # 系统模块
├── user/               # 用户模块 (待迁移)
├── writer/             # 作者模块 (待迁移)
└── version_api.go      # 版本API
```

### 模块说明

#### 按功能划分的模块 (推荐模式)

| 模块 | 职责 | 状态 |
|------|------|------|
| `admin` | 管理员功能 | ✅ 活跃 |
| `ai` | AI服务集成 | ✅ 活跃 |
| `announcements` | 公告管理 | ✅ 活跃 |
| `audit` | 审计日志 | ✅ 活跃 |
| `auth` | 认证授权 | ✅ 活跃 |
| `bookstore` | 书店功能 | ✅ 活跃 |
| `finance` | 财务管理 | ✅ 活跃 |
| `messages` | 消息通信 | ✅ 活跃 |
| `notifications` | 通知服务 | ✅ 活跃 |
| `recommendation` | 推荐服务 | ✅ 活跃 |
| `search` | 搜索服务 | ✅ 活跃 |
| `shared` | 共享工具 | ✅ 活跃 |
| `social` | 社交功能 | ✅ 活跃 |
| `stats` | 统计分析 | ✅ 活跃 |
| `system` | 系统管理 | ✅ 活跃 |

#### 按角色划分的模块 (待重构)

| 模块 | 职责 | 状态 | 迁移计划 |
|------|------|------|----------|
| `reader` | 阅读相关功能 | ⚠️ 待迁移 | 迁移到 content-management |
| `user` | 用户相关功能 | ⚠️ 待迁移 | 已迁移到 user-management |
| `writer` | 作者相关功能 | ⚠️ 待迁移 | 迁移到 content-management |

### 标准模块结构

每个模块应遵循以下标准结构：

```
module_name/
├── handler/            # 处理器 (按功能细分)
│   ├── auth_handler.go
│   ├── profile_handler.go
│   └── admin_handler.go
├── dto/                # 数据传输对象
│   ├── auth_dto.go
│   ├── user_dto.go
│   └── dto_test.go
├── converter.go        # 数据转换器 (可选)
├── models.go           # 请求/响应模型 (可选)
├── routes.go           # 路由注册
├── middleware.go       # 模块特定中间件 (可选)
├── validator.go        # 自定义验证器 (可选)
└── README.md           # 模块文档
```

---

## 设计原则

### 1. 按功能领域划分

**原则**: 模块应按业务功能领域组织，而非按用户角色划分。

**示例**:
- ✅ 正确: `user-management`, `content-management`, `social`
- ❌ 错误: `reader`, `writer`, `admin`

### 2. 单一职责原则

**原则**: 每个模块只负责一个明确的业务领域。

**要求**:
- 模块内部按handler进一步细分
- 每个handler只处理一类相关功能
- 避免功能交叉和重复

### 3. 统一DTO管理

**原则**: 所有跨模块共享的DTO统一定义在 `shared` 包。

**规则**:
- 模块私有DTO放在模块的 `dto/` 目录
- 跨模块共享DTO放在 `shared/` 目录
- 避免重复定义相同的DTO

### 4. 依赖方向

**原则**: 模块间的依赖应该是单向的，避免循环依赖。

```
business-modules → shared → internal/pkg
```

**规则**:
- 业务模块可以依赖 shared
- shared 可以依赖 internal/pkg
- 业务模块之间避免直接依赖
- 通过 shared 或 service 层交互

### 5. RESTful设计

**原则**: 遵循RESTful API设计规范。

**要求**:
- 使用合适的HTTP方法 (GET, POST, PUT, DELETE)
- 使用资源命名的URL路径
- 统一的响应格式
- 合理的HTTP状态码

---

## 依赖关系

### 模块依赖图

```
┌─────────────────────────────────────────────────────────┐
│                      业务模块层                          │
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐     │
│  │ admin│  │  ai  │  │social│  │content│  │finance│     │
│  └──┬───┘  └──┬───┘  └──┬───┘  └──┬───┘  └──┬───┘     │
│     │         │         │         │         │           │
└─────┼─────────┼─────────┼─────────┼─────────┼───────────┘
      │         │         │         │         │
      ▼         ▼         ▼         ▼         ▼
┌─────────────────────────────────────────────────────────┐
│                      共享层                              │
│              ┌─────────────────────────┐                │
│              │       shared/           │                │
│              │  - DTO                  │                │
│              │  - Response Helper      │                │
│              │  - Validator            │                │
│              │  - Common Types         │                │
│              └─────────────────────────┘                │
└─────────────────────────────────────────────────────────┘
      │         │         │         │         │
      ▼         ▼         ▼         ▼         ▼
┌─────────────────────────────────────────────────────────┐
│                    内部层                                │
│  ┌──────────────┐        ┌──────────────┐              │
│  │ internal/    │        │   pkg/       │              │
│  │ - service    │  ◄────►│ - middleware │              │
│  │ - repository │        │ - errors     │              │
│  │ - model      │        │ - utils      │              │
│  └──────────────┘        └──────────────┘              │
└─────────────────────────────────────────────────────────┘
```

### 具体依赖说明

#### shared 依赖

```go
// 所有业务模块都依赖 shared
import "github.com/yukin371/Qingyu/backend/api/v1/shared"
```

**shared 提供**:
- 统一响应函数 (`response.go`)
- 请求验证器 (`request_validator.go`)
- 公共类型定义 (`types.go`, `user_types.go`)
- 数据转换器 (`converter.go`)

#### internal/service 依赖

```go
// 业务模块依赖服务层
import "github.com/yukin371/Qingyu/backend/internal/service"
```

**service 提供**:
- 业务逻辑实现
- 跨模块协调
- 事务管理

#### pkg/middleware 依赖

```go
// 业务模块使用中间件
import "github.com/yukin371/Qingyu/backend/pkg/middleware"
```

**middleware 提供**:
- 认证中间件
- 权限检查中间件
- 限流中间件
- 日志中间件

---

## BaseAPI使用说明

### 概述

BaseAPI 是所有API处理器的基础结构，提供通用功能和工具方法。

### 使用方式

#### 方式1: 直接使用shared包函数 (推荐)

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/yukin371/Qingyu/backend/api/v1/shared"
    "net/http"
)

type AuthHandler struct {
    // 依赖注入
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if !shared.ValidateRequest(c, &req) {
        return
    }

    // 业务逻辑...

    // 成功响应
    shared.Success(c, http.StatusOK, "登录成功", data)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
    // 业务逻辑...

    // 错误处理
    if err != nil {
        shared.NotFound(c, "用户不存在")
        return
    }

    shared.OK(c, data)
}
```

#### 方式2: 创建模块特定Base结构

```go
package handler

import (
    "github.com/gin-gonic/gin"
)

type BaseHandler struct {
    // 公共字段
}

func (b *BaseHandler) Success(c *gin.Context, data interface{}) {
    shared.Success(c, http.StatusOK, "操作成功", data)
}

func (b *BaseHandler) Error(c *gin.Context, code int, message string) {
    shared.Error(c, code, message, "")
}

type UserHandler struct {
    BaseHandler
    // 特定字段
}
```

### 响应函数

#### Success - 成功响应

```go
func Success(c *gin.Context, code int, message string, data interface{})
```

**示例**:
```go
shared.Success(c, http.StatusOK, "获取成功", user)
shared.Success(c, http.StatusCreated, "创建成功", book)
```

#### 便捷方法

```go
shared.OK(c, data)              // 200 OK
shared.Created(c, data)         // 201 Created
shared.BadRequest(c, msg, det)  // 400 Bad Request
shared.Unauthorized(c, msg)     // 401 Unauthorized
shared.Forbidden(c, msg)        // 403 Forbidden
shared.NotFound(c, msg)         // 404 Not Found
shared.InternalError(c, msg, err) // 500 Internal Error
```

#### Paginated - 分页响应

```go
func Paginated(c *gin.Context, data interface{}, total int64, page int, pageSize int, message string)
```

**示例**:
```go
shared.Paginated(c, users, total, page, pageSize, "获取成功")
```

### 请求验证

#### ValidateRequest - 验证请求体

```go
func ValidateRequest(c *gin.Context, req interface{}) bool
```

**示例**:
```go
var req CreateUserRequest
if !shared.ValidateRequest(c, &req) {
    return  // 验证失败会自动返回400错误
}
```

#### ValidateQueryParams - 验证查询参数

```go
func ValidateQueryParams(c *gin.Context, req interface{}) bool
```

**示例**:
```go
var req ListUsersRequest
if !shared.ValidateQueryParams(c, &req) {
    return
}
```

---

## 中间件使用说明

### 可用中间件

| 中间件 | 位置 | 用途 |
|--------|------|------|
| Auth | `internal/middleware/auth` | JWT认证 |
| Permission | `internal/middleware/core` | 权限检查 |
| RateLimit | `internal/middleware/ratelimit` | 请求限流 |
| CORS | `pkg/middleware/cors` | 跨域支持 |
| Logger | `pkg/middleware/logger` | 请求日志 |
| Recovery | `pkg/middleware/recovery` | 恢复panic |
| Validation | `internal/middleware/validation` | 参数验证 |
| Version | `internal/middleware/version_routing` | 版本路由 |

### 使用示例

#### 在路由中使用

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/yukin371/Qingyu/backend/internal/middleware/auth"
    "github.com/yukin371/Qingyu/backend/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine) {
    // 公开路由
    public := r.Group("/api/v1/public")
    {
        public.POST("/login", handler.Login)
    }

    // 需要认证的路由
    authenticated := r.Group("/api/v1")
    authenticated.Use(auth.AuthMiddleware())
    {
        authenticated.GET("/profile", handler.GetProfile)
    }

    // 需要管理员权限的路由
    admin := r.Group("/api/v1/admin")
    admin.Use(
        auth.AuthMiddleware(),
        middleware.RequirePermission("admin"),
    )
    {
        admin.GET("/users", handler.ListUsers)
    }

    // 使用限流
    limited := r.Group("/api/v1")
    limited.Use(middleware.RateLimitMiddleware(100, 60))
    {
        limited.GET("/search", handler.Search)
    }
}
```

#### 创建模块特定中间件

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "strings"
)

func ContentTypeMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查Content-Type
        contentType := c.GetHeader("Content-Type")
        if !strings.Contains(contentType, "application/json") {
            c.AbortWithStatusJSON(415, gin.H{
                "code":    415,
                "message": "Unsupported Media Type",
            })
            return
        }
        c.Next()
    }
}
```

---

## 错误处理规范

> **参考文档**：主仓库 `docs/standards/api-status-code-standard.md`

### 错误响应格式

所有API响应必须包含完整的响应体：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "request_id": "01HR4XM2K9Y5P3Q7R6T8W0N1V2",
  "timestamp": 1737792000
}
```

### HTTP状态码使用

**重要原则**：
- 统一使用 `200 OK` 表示成功（除创建场景外）
- **DELETE操作返回 200 OK**（不使用 204 No Content）
- 所有成功响应必须包含完整响应体

| 状态码 | 场景 | 示例 | shared方法 |
|--------|------|------|-----------|
| 200 | 成功 | GET、PUT、PATCH、DELETE | `shared.Success(c, http.StatusOK, ...)` |
| 201 | 创建成功 | POST创建新资源 | `shared.Success(c, http.StatusCreated, ...)` |
| 202 | 异步任务 | POST异步操作 | `shared.Success(c, http.StatusAccepted, ...)` |
| 400 | 请求参数错误 | 参数验证失败 | `shared.BadRequest(c, "参数错误", msg)` |
| 401 | 未认证 | Token无效或过期 | `shared.Unauthorized(c, "未认证")` |
| 403 | 禁止访问 | 权限不足 | `shared.Forbidden(c, "权限不足")` |
| 404 | 资源不存在 | 资源未找到 | `shared.NotFound(c, "资源不存在")` |
| 409 | 资源冲突 | 资源已存在 | `shared.Error(c, http.StatusConflict, ...)` |
| 422 | 无法处理 | 业务规则验证失败 | `shared.Error(c, http.StatusUnprocessableEntity, ...)` |
| 429 | 请求过多 | 超过限流阈值 | `shared.Error(c, http.StatusTooManyRequests, ...)` |
| 500 | 服务器错误 | 未预期的错误 | `shared.InternalError(c, "操作失败", err)` |
| 503 | 服务不可用 | 服务维护中 | `shared.Error(c, http.StatusServiceUnavailable, ...)` |

### 状态码规范细节

#### 1. 成功响应状态码

**200 OK** - 适用于：
- GET、PUT、PATCH、DELETE 操作
- 非创建的 POST 操作（登录、重置密码、批量操作等）

**201 Created** - 适用于：
- POST 创建新资源（必须有Location响应头）
- 创建用户、权限、角色、书签、消息等

**202 Accepted** - 适用于：
- POST 异步任务操作
- 异步发布、批量操作、长时间运行的任务

#### 2. DELETE操作特殊说明

```go
// ✅ 正确：DELETE返回 200 OK
shared.Success(c, http.StatusOK, "删除成功", nil)

// ❌ 错误：不要使用 204
// shared.NoContent(c) // 删除操作不应使用
```

**原因**：项目使用统一响应格式，需要返回 `code`, `message`, `request_id`, `timestamp` 字段

### 错误处理模式

根据主仓库 `docs/standards/api-status-code-standard.md` 规范，使用 `shared` 包的方法进行统一错误处理。

#### 1. 参数验证错误

```go
// ✅ 使用 binding 标签自动验证
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
}

// 手动验证
if id == "" {
    shared.BadRequest(c, "参数错误", "ID不能为空")
    return
}
```

#### 2. 业务逻辑错误

```go
user, err := h.service.GetUserByID(id)
if err != nil {
    // 根据Service层返回的错误类型处理
    if strings.Contains(err.Error(), "not found") {
        shared.NotFound(c, "用户不存在")
        return
    }
    shared.InternalError(c, "获取用户失败", err)
    return
}

// ✅ 推荐：使用Service层返回的标准错误
shared.Success(c, http.StatusOK, "获取成功", user)
```

#### 3. 权限错误

```go
// ✅ 使用shared的权限错误方法
if !hasPermission {
    shared.Forbidden(c, "权限不足")
    return
}

// 或使用 HTTP 状态码
shared.Error(c, http.StatusForbidden, "权限不足", "需要管理员权限")
```

#### 4. 资源冲突

```go
// ✅ 409 Conflict
shared.Error(c, http.StatusConflict, "用户名已存在", "")
```

#### 5. 统一响应方法

```go
// 成功响应
shared.Success(c, http.StatusOK, "操作成功", data)      // 一般成功
shared.Success(c, http.StatusCreated, "创建成功", data)    // 创建资源
shared.Success(c, http.StatusAccepted, "任务已创建", data) // 异步任务

// 错误响应
shared.BadRequest(c, "参数错误", "用户ID不能为空")
shared.Unauthorized(c, "未认证")
shared.Forbidden(c, "权限不足")
shared.NotFound(c, "资源不存在")
shared.Error(c, http.StatusConflict, "资源已存在", "")
shared.Error(c, http.StatusUnprocessableEntity, "业务规则验证失败", "")
shared.InternalError(c, "操作失败", err)
```

### 预定义错误

错误定义在 `pkg/errors` 包中，Service层返回标准错误类型，API层通过错误匹配进行响应处理。

```go
// Service层返回错误
return errors.ErrUserNotFound
return errors.ErrInvalidInput
return errors.ErrDuplicateEntry

// API层处理
if errors.Is(err, errors.ErrUserNotFound) {
    shared.NotFound(c, "用户不存在")
    return
}
```
)
```

---

## 代码示例

### 完整的Handler示例

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/yukin371/Qingyu/backend/api/v1/shared"
    "github.com/yukin371/Qingyu/backend/internal/service"
    "net/http"
)

// UserHandler 用户处理器
type UserHandler struct {
    userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// LoginRequest 登录请求
type LoginRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
}

// Login 登录
// @Summary 用户登录
// @Description 使用用户名和密码登录
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
    // 验证请求
    var req LoginRequest
    if !shared.ValidateRequest(c, &req) {
        return
    }

    // 业务逻辑
    token, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
    if err != nil {
        shared.Unauthorized(c, "用户名或密码错误")
        return
    }

    // 返回结果
    shared.Success(c, http.StatusOK, "登录成功", gin.H{
        "token": token,
        "expires_in": 3600,
    })
}

// GetProfile 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags user
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /api/v1/profile [get]
// @Security Bearer
func (h *UserHandler) GetProfile(c *gin.Context) {
    // 从上下文获取用户ID (由认证中间件设置)
    userID, exists := c.Get("user_id")
    if !exists {
        shared.Unauthorized(c, "未认证")
        return
    }

    // 获取用户信息
    user, err := h.userService.GetUserByID(c.Request.Context(), userID.(string))
    if err != nil {
        shared.InternalError(c, "获取用户信息失败", err)
        return
    }

    shared.OK(c, user)
}

// ListUsersRequest 用户列表请求
type ListUsersRequest struct {
    Page     int    `form:"page,default=1" binding:"min=1"`
    PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
    Keyword  string `form:"keyword"`
    Role     string `form:"role"`
}

// ListUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表，支持关键词和角色筛选
// @Tags admin
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "搜索关键词"
// @Param role query string false "角色筛选"
// @Success 200 {object} shared.PaginatedResponse
// @Router /api/v1/admin/users [get]
// @Security Bearer
func (h *UserHandler) ListUsers(c *gin.Context) {
    // 验证查询参数
    var req ListUsersRequest
    if !shared.ValidateQueryParams(c, &req) {
        return
    }

    // 获取用户列表
    users, total, err := h.userService.ListUsers(
        c.Request.Context(),
        req.Page,
        req.PageSize,
        req.Keyword,
        req.Role,
    )
    if err != nil {
        shared.InternalError(c, "获取用户列表失败", err)
        return
    }

    // 返回分页结果
    shared.Paginated(c, users, total, req.Page, req.PageSize, "获取成功")
}
```

### 路由注册示例

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/yukin371/Qingyu/backend/api/v1/handler"
    "github.com/yukin371/Qingyu/backend/internal/middleware/auth"
    "github.com/yukin371/Qingyu/backend/internal/middleware/core"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
    // 公开路由
    public := r.Group("/api/v1/auth")
    {
        public.POST("/login", userHandler.Login)
        public.POST("/register", userHandler.Register)
    }

    // 认证路由
    authenticated := r.Group("/api/v1")
    authenticated.Use(auth.AuthMiddleware())
    {
        authenticated.GET("/profile", userHandler.GetProfile)
        authenticated.PUT("/profile", userHandler.UpdateProfile)
        authenticated.PUT("/password", userHandler.ChangePassword)
    }

    // 管理员路由
    admin := r.Group("/api/v1/admin/users")
    admin.Use(
        auth.AuthMiddleware(),
        core.RequirePermission("admin"),
    )
    {
        admin.GET("", userHandler.ListUsers)
        admin.GET("/:id", userHandler.GetUser)
        admin.PUT("/:id", userHandler.UpdateUser)
        admin.DELETE("/:id", userHandler.DeleteUser)
        admin.POST("/:id/ban", userHandler.BanUser)
        admin.POST("/:id/unban", userHandler.UnbanUser)
    }
}
```

---

## 最佳实践

### 1. 命名规范

#### 文件命名
- handler文件: `*_handler.go` 或 `*_api.go`
- DTO文件: `*_dto.go`
- 路由文件: `routes.go`

#### 函数命名
- Handler方法: 大写开头 (导出) + HTTP方法
  - `GetUser`, `CreateBook`, `UpdateProfile`
- 私有函数: 小写开头
  - `validateUser`, `formatResponse`

#### 变量命名
- 请求DTO: `*Request` 后缀
  - `CreateUserRequest`, `LoginRequest`
- 响应DTO: `*Response` 后缀
  - `UserResponse`, `BookResponse`
- Handler实例: `*Handler` 后缀
  - `UserHandler`, `BookHandler`

### 2. 结构组织

#### Handler结构
```go
// ✅ 推荐: 使用依赖注入
type UserHandler struct {
    userService *service.UserService
    authService *service.AuthService
}

func NewUserHandler(
    userService *service.UserService,
    authService *service.AuthService,
) *UserHandler {
    return &UserHandler{
        userService: userService,
        authService: authService,
    }
}
```

#### DTO组织
```go
// ✅ 推荐: 按功能组织DTO文件
// auth_dto.go
type LoginRequest struct { ... }
type RegisterRequest struct { ... }

// user_dto.go
type UserResponse struct { ... }
type UpdateProfileRequest struct { ... }
```

### 3. 错误处理

```go
// ✅ 推荐: 明确的错误处理
user, err := h.service.GetUser(id)
if err != nil {
    if errors.Is(err, service.ErrNotFound) {
        shared.NotFound(c, "用户不存在")
        return
    }
    shared.InternalError(c, "获取用户失败", err)
    return
}
shared.OK(c, user)
```

### 4. 参数验证

```go
// ✅ 推荐: 使用binding标签
type Request struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// ✅ 推荐: 自定义验证
func (r *RegisterRequest) Validate() error {
    if !isValidPhone(r.Phone) {
        return errors.New("手机号格式错误")
    }
    return nil
}
```

### 5. 响应格式

```go
// ✅ 推荐: 使用统一响应函数
shared.OK(c, data)
shared.Success(c, http.StatusCreated, "创建成功", data)
shared.Paginated(c, items, total, page, size, "获取成功")

// ❌ 避免: 直接使用c.JSON
c.JSON(200, gin.H{"data": data})
```

### 6. 注释规范

```go
// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录系统
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} shared.APIResponse
// @Failure 401 {object} shared.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
    // ...
}
```

---

## 相关文档

- [API设计规范](../../docs/STANDARDS.md)
- [共享模块说明](./shared/README.md)
- [中间件文档](../../docs/middleware/README.md)
- [错误处理规范](../../docs/engineering/软件工程规范_v2.0.md)
- [重构实施计划](../../docs/api_refactor_plan.md)

---

## 版本历史

| 版本 | 日期 | 变更内容 |
|------|------|----------|
| v1.0 | 2026-02-26 | 初始版本，整合API层规范 |

---

**维护者**: Backend Team
**最后更新**: 2026-02-26
