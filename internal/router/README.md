# 路由配置文档

## 概述

本目录包含 Qingyu_backend 项目的路由配置文件。路由是应用程序的入口点，负责将HTTP请求分发到对应的处理器。

## 当前状态 (2026-02-27)

### 已完成

1. **废弃模块清理**
   - `api/v1/shared/notification_api.go` 已删除
   - 通知功能已迁移到 `api/v1/notifications/` 和 `router/notifications/`
   - 所有路由引用已更新

2. **路由架构优化**
   - 主路由入口: `router/enter.go`
   - 模块化路由: 各功能模块的路由独立在 `router/<module>/` 目录下

## 路由结构

### 主路由文件

| 文件 | 说明 |
|------|------|
| `router/enter.go` | 主路由注册入口，负责注册所有模块的路由 |
| `internal/router/setup.go` | 新架构的路由配置（中间件配置） |

### 模块路由目录

| 目录 | 路由前缀 | 说明 |
|------|---------|------|
| `router/announcements/` | `/api/v1/announcements` | 公告模块路由 |
| `router/admin/` | `/api/v1/admin` | 管理员路由 |
| `router/ai/` | `/api/v1/ai` | AI服务路由 |
| `router/bookstore/` | `/api/v1/bookstore` | 书店路由 |
| `router/finance/` | `/api/v1/finance` | 财务路由 |
| `router/notifications/` | `/api/v1/notifications` | 通知路由 |
| `router/reader/` | `/api/v1/reader` | 阅读器路由 |
| `router/reading-stats/` | `/api/v1/reading-stats` | 阅读统计路由 |
| `router/recommendation/` | `/api/v1/recommendation` | 推荐系统路由 |
| `router/shared/` | `/api/v1/shared` | 共享服务路由（认证、存储） |
| `router/social/` | `/api/v1/social` | 社交功能路由 |
| `router/system/` | `/api/v1/system` | 系统监控路由 |
| `router/user/` | `/api/v1/user` | 用户管理路由 |
| `router/writer/` | `/api/v1/writer` | 写作端路由 |

## 路由注册流程

### 1. 主入口

`router/enter.go` 的 `RegisterRoutes()` 函数是路由注册的主入口:

```go
func RegisterRoutes(r *gin.Engine) {
    // 初始化日志
    logger := initRouterLogger()

    // 创建API版本组
    v1 := r.Group("/api/v1")

    // 获取服务容器
    serviceContainer := service.GetServiceContainer()

    // 注册各模块路由
    // ...
}
```

### 2. 模块路由注册

每个模块都有自己的注册函数，例如:

```go
// 通知路由
notificationsRouter.RegisterRoutes(v1, notificationAPI)
notificationsRouter.RegisterUserManagementRoutes(v1, notificationAPI)

// 社交路由
socialRouter.RegisterSocialRoutes(v1, relationAPI, commentAPI, ...)

// 公告路由
announcementsRouter.RegisterAnnouncementRoutes(v1, announcementSvc)
```

## 路由列表

### 认证路由 (`/api/v1/shared/auth`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/send-verification-code` | 发送验证码 | 否 |
| POST | `/register` | 用户注册 | 否 |
| POST | `/login` | 用户登录 | 否 |
| POST | `/logout` | 用户登出 | 是 |
| POST | `/refresh` | 刷新令牌 | 是 |
| GET | `/permissions` | 获取用户权限 | 是 |
| GET | `/roles` | 获取用户角色 | 是 |

### 通知路由 (`/api/v1/notifications`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `` | 获取通知列表 | 是 |
| GET | `/unread-count` | 获取未读数量 | 是 |
| GET | `/stats` | 获取通知统计 | 是 |
| GET | `/:id` | 获取通知详情 | 是 |
| POST | `/:id/read` | 标记为已读 | 是 |
| POST | `/batch-read` | 批量标记已读 | 是 |
| POST | `/read-all` | 全部标记已读 | 是 |
| POST | `/clear-read` | 清除已读通知 | 是 |
| POST | `/:id/resend` | 重新发送通知 | 是 |
| GET | `/ws-endpoint` | 获取WebSocket端点 | 是 |
| DELETE | `/:id` | 删除通知 | 是 |
| POST | `/batch-delete` | 批量删除 | 是 |
| DELETE | `/delete-all` | 删除所有 | 是 |
| GET | `/preferences` | 获取通知偏好 | 是 |
| PUT | `/preferences` | 更新通知偏好 | 是 |
| POST | `/preferences/reset` | 重置通知偏好 | 是 |
| POST | `/push/register` | 注册推送设备 | 是 |
| DELETE | `/push/unregister/:deviceId` | 取消注册推送设备 | 是 |
| GET | `/push/devices` | 获取推送设备列表 | 是 |

### 公告路由 (`/api/v1/announcements`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/effective` | 获取有效公告 | 否 |
| GET | `/:id` | 获取公告详情 | 否 |
| POST | `/:id/view` | 增加查看次数 | 否 |

### 社交路由 (`/api/v1/social`)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/follow/:userId` | 关注用户 | 是 |
| DELETE | `/follow/:userId` | 取消关注 | 是 |
| GET | `/follow/:userId/status` | 检查关注状态 | 是 |
| GET | `/users/:userId/followers` | 获取粉丝列表 | 是 |
| GET | `/users/:userId/following` | 获取关注列表 | 是 |
| POST | `/comments` | 创建评论 | 是 |
| GET | `/comments` | 获取评论列表 | 是 |
| POST | `/books/:bookId/like` | 点赞书籍 | 是 |
| DELETE | `/books/:bookId/like` | 取消点赞 | 是 |
| POST | `/collections` | 添加收藏 | 是 |
| GET | `/collections` | 获取收藏列表 | 是 |

## 废弃路由变更记录

### 2026-02-27 - 通知路由迁移

**废弃**: `api/v1/shared/notification_api.go`

**迁移到**:
- API: `api/v1/notifications/notification_api.go`
- 路由: `router/notifications/notification_router.go`

**影响**:
- 通知相关API从 `/api/v1/shared/notifications/*` 迁移到 `/api/v1/notifications/*`
- 用户管理相关通知设置保留 `/api/v1/user-management/*` 路径
- 推送设备管理迁移到 `/api/v1/notifications/push/*`

**兼容性**:
- 旧路由已不再可用
- 客户端需要更新到新路由路径

## 待处理事项

### 高优先级

1. **Messaging模块整合** (后续P2任务)
   - 当前保留旧版 `api/v1/messages/message_api.go`
   - 计划统一到 `/api/v1/social/messages/*`

### 中优先级

1. **路由文档完善**
   - 补充所有模块的完整路由列表
   - 添加请求/响应示例

2. **路由分组优化**
   - 审查路由分组是否合理
   - 统一路由命名规范

## 开发指南

### 添加新路由

1. 在对应模块的路由目录创建或修改路由文件
2. 实现 `Register*Routes()` 函数
3. 在 `router/enter.go` 中调用注册函数
4. 更新本文档

### 路由命名规范

- 使用RESTful风格
- 复数名词表示资源 (如 `/books`, `/comments`)
- 使用连字符分隔多词路径 (如 `/book-details`)
- HTTP方法语义化:
  - `GET`: 查询
  - `POST`: 创建
  - `PUT/PATCH`: 更新
  - `DELETE`: 删除

## 参考文档

- [Gin路由文档](https://gin-gonic.com/docs/examples/routes-grouping/)
- [RESTful API设计规范](https://restfulapi.net/)
- [项目API设计规范](../../docs/api/API设计规范.md)

---

**最后更新**: 2026-02-27
**维护者**: Qingyu Backend Team
