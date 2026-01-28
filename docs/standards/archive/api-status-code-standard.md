# API状态码规范

## 概述

本文档定义了 Qingyu 项目中 RESTful API 的状态码使用规范，确保所有 API 端点返回一致且符合 HTTP 标准的状态码。

## POST 操作状态码规范

### 1. 201 Created - 资源创建成功

当 POST 请求用于创建新资源时，应返回 `201 Created` 状态码。

**适用场景：**
- 创建新用户
- 创建新权限
- 创建新角色
- 创建新书签
- 发送消息
- 创建@提醒
- 等等

**示例：**
```go
// POST /api/v1/user/auth/register
shared.Success(c, http.StatusCreated, "注册成功", registerResp)

// POST /api/v1/admin/permissions
shared.Success(c, http.StatusCreated, "创建成功", nil)

// POST /api/v1/reader/books/{bookId}/bookmarks
shared.Success(c, http.StatusCreated, "创建成功", bookmark)
```

### 2. 200 OK - 非 POST 创建操作

当 POST 请求不用于创建资源时（如登录、批量操作、审核等），应返回 `200 OK` 状态码。

**适用场景：**
- 用户登录
- 重置密码
- 批量更新状态
- 批量删除
- 分配权限/角色
- 审核操作
- 统计数据更新
- 等等

**示例：**
```go
// POST /api/v1/user/auth/login
shared.Success(c, http.StatusOK, "登录成功", loginResp)

// POST /api/v1/admin/users/{id}/reset-password
shared.Success(c, http.StatusOK, "重置成功", result)

// POST /api/v1/admin/users/batch-update-status
shared.Success(c, http.StatusOK, "批量更新成功", nil)

// POST /api/v1/admin/audit/{id}/review
shared.Success(c, http.StatusOK, "审核已通过", nil)
```

### 3. 202 Accepted - 异步任务处理中

当 POST 请求触发异步操作时，应返回 `202 Accepted` 状态码。

**适用场景：**
- 异步发布项目
- 异步批量操作
- 长时间运行的任务

**示例：**
```go
// POST /api/v1/writer/projects/{id}/publish
shared.Success(c, http.StatusAccepted, "发布任务已创建", record)

// POST /api/v1/writer/documents/{id}/publish
shared.Success(c, http.StatusAccepted, "发布任务已创建", record)

// POST /api/v1/writer/projects/{projectId}/documents/batch-publish
shared.Success(c, http.StatusAccepted, "批量发布任务已创建", result)
```

## 其他常用状态码

### 成功响应

| 状态码 | 含义 | 使用场景 |
|--------|------|----------|
| 200 OK | 请求成功 | GET、PUT、DELETE 及非创建的 POST 操作 |
| 201 Created | 资源创建成功 | POST 创建新资源 |
| 202 Accepted | 请求已接受 | POST 异步任务 |
| 204 No Content | 成功无返回内容 | DELETE 操作 |

### 客户端错误

| 状态码 | 含义 | 使用场景 |
|--------|------|----------|
| 400 Bad Request | 请求参数错误 | 参数验证失败 |
| 401 Unauthorized | 未授权 | 未登录或 Token 无效 |
| 403 Forbidden | 禁止访问 | 权限不足 |
| 404 Not Found | 资源不存在 | 请求的资源未找到 |
| 409 Conflict | 资源冲突 | 资源已存在或状态冲突 |
| 422 Unprocessable Entity | 无法处理 | 语义错误（如业务规则验证失败） |
| 429 Too Many Requests | 请求过多 | 超过速率限制 |

### 服务端错误

| 状态码 | 含义 | 使用场景 |
|--------|------|----------|
| 500 Internal Server Error | 服务器内部错误 | 未预期的服务器错误 |
| 502 Bad Gateway | 网关错误 | 上游服务错误 |
| 503 Service Unavailable | 服务不可用 | 服务维护或过载 |

## 实施指南

### 1. 使用 shared.Success() 方法

```go
// 创建资源
shared.Success(c, http.StatusCreated, "创建成功", data)

// 非创建操作
shared.Success(c, http.StatusOK, "操作成功", data)

// 异步任务
shared.Success(c, http.StatusAccepted, "任务已创建", data)
```

### 2. 错误处理

```go
// 参数错误
shared.BadRequest(c, "参数错误", "用户ID不能为空")
shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())

// 认证错误
shared.Unauthorized(c, "未认证")
shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")

// 权限错误
shared.Error(c, http.StatusForbidden, "权限不足", err.Error())

// 资源不存在
shared.NotFound(c, "用户不存在")
shared.Error(c, http.StatusNotFound, "资源不存在", err.Error())

// 服务器错误
shared.InternalError(c, "操作失败", err)
shared.Error(c, http.StatusInternalServerError, "操作失败", err.Error())
```

### 3. Swagger 注释规范

在编写 API 处理器时，确保 Swagger 注释中的 `@Success` 状态码与实际返回的状态码一致：

```go
// CreatePermission 创建权限
//
//	@Summary		创建权限
//	@Description	管理员创建新权限
//	@Tags			Admin-Permission
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.Permission	true	"权限信息"
//	@Success		201		{object}	shared.APIResponse  // 注意这里是 201
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		403		{object}	shared.APIResponse
//	@Router			/api/v1/admin/permissions [post]
func (api *PermissionAPI) CreatePermission(c *gin.Context) {
	// ...
	shared.Success(c, http.StatusCreated, "创建成功", nil)
}
```

## 项目当前实施情况

### 已正确实施的端点

#### 创建资源 (201 Created)
- ✅ POST /api/v1/user/auth/register
- ✅ POST /api/v1/admin/permissions (已修复)
- ✅ POST /api/v1/admin/roles (已修复)
- ✅ POST /api/v1/social/messages
- ✅ POST /api/v1/social/mentions
- ✅ POST /api/v1/reader/books/{bookId}/bookmarks (已修复)

#### 非创建操作 (200 OK)
- ✅ POST /api/v1/user/auth/login
- ✅ POST /api/v1/admin/users/{id}/reset-password
- ✅ POST /api/v1/admin/users/batch-update-status
- ✅ POST /api/v1/admin/users/batch-delete
- ✅ POST /api/v1/admin/roles/{id}/permissions/{permissionCode}
- ✅ POST /api/v1/admin/users/{userId}/roles
- ✅ POST /api/v1/admin/audit/{id}/review
- ✅ POST /api/v1/admin/audit/{id}/appeal/review
- ✅ POST /api/v1/bookstore/books/{id}/view
- ✅ POST /api/v1/bookstore/banners/{id}/click
- ✅ POST /api/v1/writer/projects/{id}/unpublish

#### 异步任务 (202 Accepted)
- ✅ POST /api/v1/writer/projects/{id}/publish
- ✅ POST /api/v1/writer/documents/{id}/publish
- ✅ POST /api/v1/writer/projects/{projectId}/documents/batch-publish

## 参考资料

- [MDN: HTTP response status codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status)
- [RFC 7231: Status Codes](https://tools.ietf.org/html/rfc7231#section-6)
- [RESTful API Design: Best Practices](https://restfulapi.net/)

## 更新日志

- 2026-01-25: 初始版本，统一所有 POST 操作的状态码
- 修复 permission_api.go 中 CreatePermission 和 CreateRole 返回 201
- 修复 bookmark_api.go 中 CreateBookmark 返回 201
