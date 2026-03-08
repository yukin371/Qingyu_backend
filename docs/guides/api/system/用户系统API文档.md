# 用户系统 API 文档

**版本**: v1.0  
**最后更新**: 2025-10-15  
**基础路径**: `/api/v1`  
**模块**: 用户系统 (System - User)

---

## 📋 目录

1. [概述](#概述)
2. [认证说明](#认证说明)
3. [公开接口](#公开接口)
   - [用户注册](#用户注册)
   - [用户登录](#用户登录)
4. [用户接口（需要认证）](#用户接口需要认证)
   - [获取个人信息](#获取个人信息)
   - [更新个人信息](#更新个人信息)
   - [修改密码](#修改密码)
5. [管理员接口](#管理员接口)
   - [获取用户列表](#获取用户列表)
   - [获取指定用户信息](#获取指定用户信息)
   - [更新用户信息](#更新用户信息)
   - [删除用户](#删除用户)
6. [数据结构](#数据结构)
7. [错误码说明](#错误码说明)
8. [前端集成示例](#前端集成示例)
9. [测试建议](#测试建议)

---

## 概述

用户系统提供完整的用户管理功能，包括：
- 用户注册与登录
- 个人信息管理
- 密码管理
- 用户权限管理（管理员）

### 功能特性

- ✅ JWT Token 认证
- ✅ 基于角色的权限控制（RBAC）
- ✅ 密码加密存储（bcrypt）
- ✅ 邮箱和手机号验证
- ✅ 用户状态管理
- ✅ 分页查询支持

### 技术架构

- **API层**: `api/v1/system/sys_user.go`
- **DTO层**: `api/v1/system/user_dto.go`
- **Service层**: `service/interfaces` (用户服务接口)
- **认证中间件**: JWT Token

---

## 认证说明

### JWT Token 认证

所有需要认证的接口都需要在请求头中携带 JWT Token：

```http
Authorization: Bearer <your_jwt_token>
```

### 获取 Token

通过以下方式获取 Token：
1. [用户注册](#用户注册) - 注册成功后自动返回 Token
2. [用户登录](#用户登录) - 登录成功后返回 Token

### Token 说明

- **有效期**: 默认 24 小时
- **存储位置**: 建议存储在 localStorage 或内存中
- **刷新机制**: Token 过期后需要重新登录
- **安全性**: 不要在 URL 中传递 Token

### 权限角色

| 角色 | 值 | 权限说明 |
|------|-----|----------|
| 普通用户 | user | 基础权限，可访问个人信息相关接口 |
| 作者 | author | 拥有用户权限 + 写作相关权限 |
| 管理员 | admin | 拥有所有权限，可管理用户 |

---

## 公开接口

### 用户注册

创建新用户账号。

**接口地址**: `POST /api/v1/register`  
**需要认证**: ❌ 否

#### 请求参数

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| username | string | ✅ | 3-50字符 | 用户名 |
| email | string | ✅ | 有效邮箱格式 | 邮箱地址 |
| password | string | ✅ | 最少6字符 | 登录密码 |

#### 请求示例

```json
{
  "username": "zhangsan",
  "email": "zhangsan@example.com",
  "password": "password123"
}
```

#### 成功响应 (201)

```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "user_id": "670abcdef123456789abcdef",
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "role": "user",
    "status": "active",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjcwYWJjZGVmMTIzNDU2Nzg5YWJjZGVmIiwiZXhwIjoxNzI5MDQ4MDAwfQ.xxx"
  },
  "timestamp": 1697203200
}
```

#### 错误响应

**400 - 参数验证失败**:
```json
{
  "code": 400,
  "message": "注册失败",
  "error": "用户名长度必须在3-50之间"
}
```

**400 - 邮箱已存在**:
```json
{
  "code": 400,
  "message": "注册失败",
  "error": "邮箱已被注册"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "注册失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "password": "password123"
  }'
```

#### 业务逻辑

1. 验证请求参数格式
2. 检查邮箱是否已注册
3. 检查用户名是否已存在
4. 使用 bcrypt 加密密码
5. 创建用户记录
6. 生成 JWT Token
7. 返回用户信息和 Token

---

### 用户登录

用户登录获取 Token。

**接口地址**: `POST /api/v1/login`  
**需要认证**: ❌ 否

#### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | ✅ | 用户名或邮箱 |
| password | string | ✅ | 登录密码 |

#### 请求示例

```json
{
  "username": "zhangsan",
  "password": "password123"
}
```

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user_id": "670abcdef123456789abcdef",
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjcwYWJjZGVmMTIzNDU2Nzg5YWJjZGVmIiwiZXhwIjoxNzI5MDQ4MDAwfQ.xxx"
  },
  "timestamp": 1697203200
}
```

#### 错误响应

**401 - 用户名或密码错误**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "用户名或密码错误"
}
```

**401 - 账号未激活** ⭐ 新增 (2025-10-17):
```json
{
  "code": 401,
  "message": "未认证",
  "error": "账号未激活，请先验证邮箱"
}
```

**401 - 账号已被封禁** ⭐ 新增 (2025-10-17):
```json
{
  "code": 401,
  "message": "未认证",
  "error": "账号已被封禁，请联系管理员"
}
```

**401 - 账号已删除** ⭐ 新增 (2025-10-17):
```json
{
  "code": 401,
  "message": "未认证",
  "error": "账号已删除"
}
```

**400 - 参数验证失败**:
```json
{
  "code": 400,
  "message": "登录失败",
  "error": "用户名和密码不能为空"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "登录失败",
  "error": "服务器内部错误"
}
```

> **📝 重要更新 (2025-10-17)**:  
> **登录现已增加用户状态检查**，系统会验证账号状态后才允许登录：
> - `active` (活跃): 正常用户，可以登录 ✅
> - `inactive` (未激活): 需要先验证邮箱才能登录 ❌
> - `banned` (已封禁): 被管理员封禁，无法登录 ❌
> - `deleted` (已删除): 账号已删除，无法登录 ❌
> 
> 详见：[登录状态检查修复说明](../../implementation/登录状态检查修复完成_2025-10-17.md)

#### cURL 示例

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "password": "password123"
  }'
```

#### 业务逻辑

1. 验证请求参数
2. 根据用户名查找用户
3. 验证密码（bcrypt 比对）
4. 更新最后登录时间和IP
5. 生成 JWT Token
6. 返回用户信息和 Token

---

## 用户接口（需要认证）

### 获取个人信息

获取当前登录用户的详细信息。

**接口地址**: `GET /api/v1/users/profile`  
**需要认证**: ✅ 是 (JWT Token)

#### 请求参数

无

#### 请求头

```http
Authorization: Bearer <your_jwt_token>
```

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "user_id": "670abcdef123456789abcdef",
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "role": "user",
    "status": "active",
    "avatar": "https://example.com/avatars/zhangsan.jpg",
    "nickname": "张三",
    "bio": "这是我的个人简介",
    "email_verified": true,
    "phone_verified": false,
    "last_login_at": "2025-10-15T10:30:00Z",
    "last_login_ip": "192.168.1.100",
    "created_at": "2025-10-01T08:00:00Z",
    "updated_at": "2025-10-15T10:30:00Z"
  },
  "timestamp": 1697203200
}
```

#### 错误响应

**401 - 未认证**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "未提供认证令牌"
}
```

**404 - 用户不存在**:
```json
{
  "code": 404,
  "message": "未找到",
  "error": "用户不存在"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "获取用户信息失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### 更新个人信息

更新当前登录用户的个人信息。

**接口地址**: `PUT /api/v1/users/profile`  
**需要认证**: ✅ 是 (JWT Token)

#### 请求参数

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| nickname | string | ❌ | 最多50字符 | 昵称 |
| bio | string | ❌ | 最多500字符 | 个人简介 |
| avatar | string | ❌ | 有效URL | 头像URL |
| phone | string | ❌ | E.164格式 | 手机号 |

**注意**: 所有字段都是可选的，只更新提供的字段。

#### 请求示例

```json
{
  "nickname": "张三同学",
  "bio": "热爱阅读和写作",
  "avatar": "https://example.com/avatars/new-avatar.jpg",
  "phone": "+8613800138000"
}
```

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "更新成功",
  "data": null,
  "timestamp": 1697203200
}
```

#### 错误响应

**400 - 参数验证失败**:
```json
{
  "code": 400,
  "message": "更新失败",
  "error": "昵称长度不能超过50字符"
}
```

**401 - 未认证**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "未提供认证令牌"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "更新失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "张三同学",
    "bio": "热爱阅读和写作"
  }'
```

---

### 修改密码

修改当前用户的登录密码。

**接口地址**: `PUT /api/v1/users/password`  
**需要认证**: ✅ 是 (JWT Token)

#### 请求参数

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| old_password | string | ✅ | - | 当前密码 |
| new_password | string | ✅ | 最少6字符 | 新密码 |

#### 请求示例

```json
{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null,
  "timestamp": 1697203200
}
```

#### 错误响应

**401 - 旧密码错误**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "旧密码错误"
}
```

**400 - 新密码不符合要求**:
```json
{
  "code": 400,
  "message": "修改密码失败",
  "error": "新密码长度不能少于6位"
}
```

**404 - 用户不存在**:
```json
{
  "code": 404,
  "message": "未找到",
  "error": "用户不存在"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "修改密码失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
curl -X PUT http://localhost:8080/api/v1/users/password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }'
```

#### 安全建议

1. 修改密码后建议重新登录
2. 密码应包含字母、数字和特殊字符
3. 不要使用常见密码
4. 定期更换密码

---

## 管理员接口

以下接口仅限管理员（admin 角色）访问。

### 获取用户列表

获取系统中的用户列表，支持分页和筛选。

**接口地址**: `GET /api/v1/admin/users`  
**需要认证**: ✅ 是 (JWT Token + Admin 角色)

#### 请求参数（Query）

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | ❌ | 1 | 页码（从1开始） |
| page_size | int | ❌ | 10 | 每页数量（最大100） |
| username | string | ❌ | - | 用户名筛选（模糊匹配） |
| email | string | ❌ | - | 邮箱筛选（模糊匹配） |
| role | string | ❌ | - | 角色筛选（精确匹配） |
| status | string | ❌ | - | 状态筛选（精确匹配） |

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "user_id": "670abcdef123456789abcdef",
      "username": "zhangsan",
      "email": "zhangsan@example.com",
      "phone": "13800138000",
      "role": "user",
      "status": "active",
      "avatar": "https://example.com/avatars/zhangsan.jpg",
      "nickname": "张三",
      "bio": "这是我的个人简介",
      "email_verified": true,
      "phone_verified": false,
      "last_login_at": "2025-10-15T10:30:00Z",
      "last_login_ip": "192.168.1.100",
      "created_at": "2025-10-01T08:00:00Z",
      "updated_at": "2025-10-15T10:30:00Z"
    },
    {
      "user_id": "670abcdef987654321fedcba",
      "username": "lisi",
      "email": "lisi@example.com",
      "phone": "",
      "role": "author",
      "status": "active",
      "avatar": "",
      "nickname": "李四",
      "bio": "",
      "email_verified": false,
      "phone_verified": false,
      "last_login_at": "2025-10-14T15:20:00Z",
      "last_login_ip": "192.168.1.101",
      "created_at": "2025-09-15T12:00:00Z",
      "updated_at": "2025-10-14T15:20:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "page_size": 10,
  "timestamp": 1697203200
}
```

#### 错误响应

**400 - 参数错误**:
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "页码必须大于0"
}
```

**401 - 未认证**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "未提供认证令牌"
}
```

**403 - 权限不足**:
```json
{
  "code": 403,
  "message": "禁止访问",
  "error": "权限不足"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "获取用户列表失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
# 获取第1页，每页10条
curl -X GET "http://localhost:8080/api/v1/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer <admin_token>"

# 筛选角色为 author 的用户
curl -X GET "http://localhost:8080/api/v1/admin/users?role=author" \
  -H "Authorization: Bearer <admin_token>"

# 搜索用户名包含 "zhang" 的用户
curl -X GET "http://localhost:8080/api/v1/admin/users?username=zhang" \
  -H "Authorization: Bearer <admin_token>"

# 组合筛选：激活状态的作者，第2页
curl -X GET "http://localhost:8080/api/v1/admin/users?page=2&role=author&status=active" \
  -H "Authorization: Bearer <admin_token>"
```

---

### 获取指定用户信息

管理员获取指定用户的详细信息。

**接口地址**: `GET /api/v1/admin/users/{id}`  
**需要认证**: ✅ 是 (JWT Token + Admin 角色)

#### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | ✅ | 用户ID |

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "user_id": "670abcdef123456789abcdef",
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "role": "user",
    "status": "active",
    "avatar": "https://example.com/avatars/zhangsan.jpg",
    "nickname": "张三",
    "bio": "这是我的个人简介",
    "email_verified": true,
    "phone_verified": false,
    "last_login_at": "2025-10-15T10:30:00Z",
    "last_login_ip": "192.168.1.100",
    "created_at": "2025-10-01T08:00:00Z",
    "updated_at": "2025-10-15T10:30:00Z"
  },
  "timestamp": 1697203200
}
```

#### 错误响应

**400 - 参数错误**:
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "用户ID不能为空"
}
```

**401 - 未认证**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "未提供认证令牌"
}
```

**403 - 权限不足**:
```json
{
  "code": 403,
  "message": "禁止访问",
  "error": "权限不足"
}
```

**404 - 用户不存在**:
```json
{
  "code": 404,
  "message": "未找到",
  "error": "用户不存在"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "获取用户信息失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
curl -X GET http://localhost:8080/api/v1/admin/users/670abcdef123456789abcdef \
  -H "Authorization: Bearer <admin_token>"
```

---

### 更新用户信息

管理员更新指定用户的信息（可以更新角色、状态等）。

**接口地址**: `PUT /api/v1/admin/users/{id}`  
**需要认证**: ✅ 是 (JWT Token + Admin 角色)

#### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | ✅ | 用户ID |

#### 请求参数

| 参数名 | 类型 | 必填 | 校验规则 | 说明 |
|--------|------|------|----------|------|
| nickname | string | ❌ | 最多50字符 | 昵称 |
| bio | string | ❌ | 最多500字符 | 个人简介 |
| avatar | string | ❌ | 有效URL | 头像URL |
| phone | string | ❌ | E.164格式 | 手机号 |
| role | string | ❌ | user/author/admin | 用户角色 |
| status | string | ❌ | active/inactive/banned/deleted | 用户状态 |
| email_verified | bool | ❌ | - | 邮箱验证状态 |
| phone_verified | bool | ❌ | - | 手机验证状态 |

**注意**: 所有字段都是可选的，只更新提供的字段。

#### 请求示例

```json
{
  "role": "author",
  "status": "active",
  "nickname": "张三作家",
  "email_verified": true
}
```

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "更新成功",
  "data": null,
  "timestamp": 1697203200
}
```

#### 错误响应

**400 - 参数验证失败**:
```json
{
  "code": 400,
  "message": "更新失败",
  "error": "无效的角色类型"
}
```

**401 - 未认证**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "未提供认证令牌"
}
```

**403 - 权限不足**:
```json
{
  "code": 403,
  "message": "禁止访问",
  "error": "权限不足"
}
```

**404 - 用户不存在**:
```json
{
  "code": 404,
  "message": "未找到",
  "error": "用户不存在"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "更新失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/670abcdef123456789abcdef \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "author",
    "status": "active",
    "email_verified": true
  }'
```

---

### 删除用户

管理员删除指定用户（软删除）。

**接口地址**: `DELETE /api/v1/admin/users/{id}`  
**需要认证**: ✅ 是 (JWT Token + Admin 角色)

#### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | ✅ | 用户ID |

#### 成功响应 (200)

```json
{
  "code": 200,
  "message": "删除成功",
  "data": null,
  "timestamp": 1697203200
}
```

#### 错误响应

**400 - 参数错误**:
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "用户ID不能为空"
}
```

**401 - 未认证**:
```json
{
  "code": 401,
  "message": "未认证",
  "error": "未提供认证令牌"
}
```

**403 - 权限不足**:
```json
{
  "code": 403,
  "message": "禁止访问",
  "error": "权限不足"
}
```

**404 - 用户不存在**:
```json
{
  "code": 404,
  "message": "未找到",
  "error": "用户不存在"
}
```

**500 - 服务器错误**:
```json
{
  "code": 500,
  "message": "删除失败",
  "error": "服务器内部错误"
}
```

#### cURL 示例

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/670abcdef123456789abcdef \
  -H "Authorization: Bearer <admin_token>"
```

#### 注意事项

- 删除操作为软删除，用户状态将被标记为 `deleted`
- 软删除的用户数据保留在数据库中，但无法登录
- 不建议删除管理员账户

---

## 数据结构

### UserStatus - 用户状态枚举

| 状态值 | 说明 | 可登录 |
|--------|------|--------|
| active | 正常激活 | ✅ |
| inactive | 未激活 | ❌ |
| banned | 已封禁 | ❌ |
| deleted | 已删除 | ❌ |

### RegisterRequest - 注册请求

```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### RegisterResponse - 注册响应

```go
type RegisterResponse struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Status   string `json:"status"`
    Token    string `json:"token"`
}
```

### LoginRequest - 登录请求

```go
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

### LoginResponse - 登录响应

```go
type LoginResponse struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Token    string `json:"token"`
}
```

### UserProfileResponse - 用户信息响应

```go
type UserProfileResponse struct {
    UserID        string    `json:"user_id"`
    Username      string    `json:"username"`
    Email         string    `json:"email"`
    Phone         string    `json:"phone,omitempty"`
    Role          string    `json:"role"`
    Status        string    `json:"status"`
    Avatar        string    `json:"avatar,omitempty"`
    Nickname      string    `json:"nickname,omitempty"`
    Bio           string    `json:"bio,omitempty"`
    EmailVerified bool      `json:"email_verified"`
    PhoneVerified bool      `json:"phone_verified"`
    LastLoginAt   time.Time `json:"last_login_at,omitempty"`
    LastLoginIP   string    `json:"last_login_ip,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

### UpdateProfileRequest - 更新个人信息请求

```go
type UpdateProfileRequest struct {
    Nickname *string `json:"nickname,omitempty" validate:"omitempty,max=50"`
    Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
    Avatar   *string `json:"avatar,omitempty" validate:"omitempty,url"`
    Phone    *string `json:"phone,omitempty" validate:"omitempty,e164"`
}
```

### ChangePasswordRequest - 修改密码请求

```go
type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=6"`
}
```

### ListUsersRequest - 用户列表请求

```go
type ListUsersRequest struct {
    Page     int        `form:"page" validate:"omitempty,min=1"`
    PageSize int        `form:"page_size" validate:"omitempty,min=1,max=100"`
    Username string     `form:"username" validate:"omitempty"`
    Email    string     `form:"email" validate:"omitempty,email"`
    Role     string     `form:"role" validate:"omitempty"`
    Status   UserStatus `form:"status" validate:"omitempty"`
}
```

### AdminUpdateUserRequest - 管理员更新用户请求

```go
type AdminUpdateUserRequest struct {
    Nickname      *string     `json:"nickname,omitempty" validate:"omitempty,max=50"`
    Bio           *string     `json:"bio,omitempty" validate:"omitempty,max=500"`
    Avatar        *string     `json:"avatar,omitempty" validate:"omitempty,url"`
    Phone         *string     `json:"phone,omitempty" validate:"omitempty,e164"`
    Role          *string     `json:"role,omitempty" validate:"omitempty,oneof=user author admin"`
    Status        *UserStatus `json:"status,omitempty" validate:"omitempty"`
    EmailVerified *bool       `json:"email_verified,omitempty"`
    PhoneVerified *bool       `json:"phone_verified,omitempty"`
}
```

---

## 错误码说明

### HTTP 状态码

| 状态码 | 说明 | 使用场景 |
|--------|------|----------|
| 200 | OK - 成功 | 请求成功 |
| 201 | Created - 创建成功 | 注册成功 |
| 400 | Bad Request - 参数错误 | 请求参数验证失败 |
| 401 | Unauthorized - 未认证 | Token 无效或未提供 |
| 403 | Forbidden - 权限不足 | 没有访问权限 |
| 404 | Not Found - 资源不存在 | 用户不存在 |
| 500 | Internal Server Error - 服务器错误 | 服务器内部错误 |

### Service 层错误类型

| 错误类型 | 说明 | 对应HTTP状态码 |
|---------|------|---------------|
| ErrorTypeValidation | 验证错误 | 400 |
| ErrorTypeBusiness | 业务逻辑错误 | 400 |
| ErrorTypeNotFound | 资源不存在 | 404 |
| ErrorTypeUnauthorized | 认证失败 | 401 |
| ErrorTypeForbidden | 权限不足 | 403 |
| ErrorTypeInternal | 内部错误 | 500 |

---

## 前端集成示例

### Vue 3 + Axios

#### 1. 创建 API 服务文件

```javascript
// src/api/user.js
import request from '@/utils/request'

/**
 * 用户注册
 */
export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

/**
 * 用户登录
 */
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

/**
 * 获取个人信息
 */
export function getProfile() {
  return request({
    url: '/users/profile',
    method: 'get'
  })
}

/**
 * 更新个人信息
 */
export function updateProfile(data) {
  return request({
    url: '/users/profile',
    method: 'put',
    data
  })
}

/**
 * 修改密码
 */
export function changePassword(data) {
  return request({
    url: '/users/password',
    method: 'put',
    data
  })
}

/**
 * 获取用户列表（管理员）
 */
export function getUserList(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

/**
 * 获取指定用户（管理员）
 */
export function getUser(id) {
  return request({
    url: `/admin/users/${id}`,
    method: 'get'
  })
}

/**
 * 更新用户（管理员）
 */
export function updateUser(id, data) {
  return request({
    url: `/admin/users/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除用户（管理员）
 */
export function deleteUser(id) {
  return request({
    url: `/admin/users/${id}`,
    method: 'delete'
  })
}
```

#### 2. Axios 请求拦截器配置

```javascript
// src/utils/request.js
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建 axios 实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    const token = userStore.token
    
    // 如果存在 token，添加到请求头
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 如果响应码不是 200/201，视为错误
    if (res.code !== 200 && res.code !== 201) {
      ElMessage.error(res.message || res.error || '请求失败')
      
      // 401: Token 过期或无效
      if (res.code === 401) {
        const userStore = useUserStore()
        userStore.logout()
        // 跳转到登录页
        window.location.href = '/login'
      }
      
      return Promise.reject(new Error(res.message || res.error || '请求失败'))
    } else {
      return res
    }
  },
  error => {
    console.error('响应错误:', error)
    
    let message = '网络错误'
    if (error.response) {
      const status = error.response.status
      switch (status) {
        case 400:
          message = '请求参数错误'
          break
        case 401:
          message = '未授权，请重新登录'
          const userStore = useUserStore()
          userStore.logout()
          window.location.href = '/login'
          break
        case 403:
          message = '权限不足'
          break
        case 404:
          message = '请求的资源不存在'
          break
        case 500:
          message = '服务器错误'
          break
        default:
          message = error.response.data?.message || '请求失败'
      }
    }
    
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default service
```

#### 3. Pinia 用户状态管理

```javascript
// src/stores/user.js
import { defineStore } from 'pinia'
import { login, register, getProfile } from '@/api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.userInfo?.role === 'admin',
    isAuthor: (state) => state.userInfo?.role === 'author' || state.userInfo?.role === 'admin'
  },
  
  actions: {
    // 登录
    async login(loginForm) {
      try {
        const res = await login(loginForm)
        this.token = res.data.token
        localStorage.setItem('token', res.data.token)
        return res
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    // 注册
    async register(registerForm) {
      try {
        const res = await register(registerForm)
        this.token = res.data.token
        localStorage.setItem('token', res.data.token)
        return res
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    // 获取用户信息
    async getUserInfo() {
      try {
        const res = await getProfile()
        this.userInfo = res.data
        return res
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    // 登出
    logout() {
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
    }
  }
})
```

#### 4. 组件使用示例

```vue
<!-- src/views/Login.vue -->
<template>
  <div class="login-container">
    <el-form :model="loginForm" :rules="rules" ref="loginFormRef">
      <el-form-item prop="username">
        <el-input v-model="loginForm.username" placeholder="用户名或邮箱" />
      </el-form-item>
      <el-form-item prop="password">
        <el-input v-model="loginForm.password" type="password" placeholder="密码" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleLogin" :loading="loading">登录</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  const valid = await loginFormRef.value.validate()
  if (!valid) return
  
  loading.value = true
  try {
    await userStore.login(loginForm)
    await userStore.getUserInfo()
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}
</script>
```

### React + Fetch

#### 1. API 服务文件

```javascript
// src/api/user.js
const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// 统一请求方法
async function request(url, options = {}) {
  const token = localStorage.getItem('token')
  
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  
  const response = await fetch(`${BASE_URL}${url}`, {
    ...options,
    headers
  })
  
  const data = await response.json()
  
  if (!response.ok) {
    // 401 未授权
    if (response.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    throw new Error(data.message || data.error || '请求失败')
  }
  
  return data
}

// 用户注册
export async function register(data) {
  const result = await request('/register', {
    method: 'POST',
    body: JSON.stringify(data)
  })
  if (result.data.token) {
    localStorage.setItem('token', result.data.token)
  }
  return result
}

// 用户登录
export async function login(data) {
  const result = await request('/login', {
    method: 'POST',
    body: JSON.stringify(data)
  })
  if (result.data.token) {
    localStorage.setItem('token', result.data.token)
  }
  return result
}

// 获取个人信息
export async function getProfile() {
  return request('/users/profile')
}

// 更新个人信息
export async function updateProfile(data) {
  return request('/users/profile', {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

// 修改密码
export async function changePassword(data) {
  return request('/users/password', {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

// 获取用户列表（管理员）
export async function getUserList(params) {
  const queryString = new URLSearchParams(params).toString()
  return request(`/admin/users?${queryString}`)
}

// 获取指定用户（管理员）
export async function getUser(id) {
  return request(`/admin/users/${id}`)
}

// 更新用户（管理员）
export async function updateUser(id, data) {
  return request(`/admin/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  })
}

// 删除用户（管理员）
export async function deleteUser(id) {
  return request(`/admin/users/${id}`, {
    method: 'DELETE'
  })
}
```

#### 2. 自定义 Hook

```javascript
// src/hooks/useAuth.js
import { useState, useEffect } from 'react'
import { login, register, getProfile } from '../api/user'

export function useAuth() {
  const [user, setUser] = useState(null)
  const [token, setToken] = useState(localStorage.getItem('token'))
  const [loading, setLoading] = useState(false)
  
  useEffect(() => {
    if (token) {
      fetchUserInfo()
    }
  }, [token])
  
  const fetchUserInfo = async () => {
    try {
      setLoading(true)
      const res = await getProfile()
      setUser(res.data)
    } catch (error) {
      console.error('获取用户信息失败:', error)
      logout()
    } finally {
      setLoading(false)
    }
  }
  
  const handleLogin = async (credentials) => {
    setLoading(true)
    try {
      const res = await login(credentials)
      setToken(res.data.token)
      setUser(res.data)
      return res
    } catch (error) {
      throw error
    } finally {
      setLoading(false)
    }
  }
  
  const handleRegister = async (data) => {
    setLoading(true)
    try {
      const res = await register(data)
      setToken(res.data.token)
      setUser(res.data)
      return res
    } catch (error) {
      throw error
    } finally {
      setLoading(false)
    }
  }
  
  const logout = () => {
    setToken(null)
    setUser(null)
    localStorage.removeItem('token')
  }
  
  return {
    user,
    token,
    loading,
    isLoggedIn: !!token,
    isAdmin: user?.role === 'admin',
    isAuthor: user?.role === 'author' || user?.role === 'admin',
    login: handleLogin,
    register: handleRegister,
    logout
  }
}
```

---

## 测试建议

### Postman 测试

#### 1. 创建环境变量

```json
{
  "baseUrl": "http://localhost:8080/api/v1",
  "token": "",
  "adminToken": "",
  "userId": ""
}
```

#### 2. 测试流程

**步骤 1: 注册新用户**
```
POST {{baseUrl}}/register
Body:
{
  "username": "testuser",
  "email": "testuser@example.com",
  "password": "password123"
}

保存返回的 token 到环境变量
```

**步骤 2: 登录**
```
POST {{baseUrl}}/login
Body:
{
  "username": "testuser",
  "password": "password123"
}
```

**步骤 3: 获取个人信息**
```
GET {{baseUrl}}/users/profile
Headers:
Authorization: Bearer {{token}}
```

**步骤 4: 更新个人信息**
```
PUT {{baseUrl}}/users/profile
Headers:
Authorization: Bearer {{token}}
Body:
{
  "nickname": "测试用户",
  "bio": "这是测试账号"
}
```

**步骤 5: 修改密码**
```
PUT {{baseUrl}}/users/password
Headers:
Authorization: Bearer {{token}}
Body:
{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

**步骤 6: 管理员登录**
```
POST {{baseUrl}}/login
Body:
{
  "username": "admin",
  "password": "admin123"
}

保存返回的 token 到 adminToken
```

**步骤 7: 获取用户列表**
```
GET {{baseUrl}}/admin/users?page=1&page_size=10
Headers:
Authorization: Bearer {{adminToken}}
```

**步骤 8: 更新用户角色**
```
PUT {{baseUrl}}/admin/users/{{userId}}
Headers:
Authorization: Bearer {{adminToken}}
Body:
{
  "role": "author",
  "status": "active"
}
```

### 单元测试示例

```go
// api/v1/system/sys_user_test.go
package system

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestUserAPI_Register(t *testing.T) {
    // 设置 Gin 为测试模式
    gin.SetMode(gin.TestMode)
    
    tests := []struct {
        name           string
        requestBody    RegisterRequest
        expectedStatus int
        expectedMsg    string
    }{
        {
            name: "成功注册",
            requestBody: RegisterRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            expectedStatus: http.StatusCreated,
            expectedMsg:    "注册成功",
        },
        {
            name: "用户名太短",
            requestBody: RegisterRequest{
                Username: "ab",
                Email:    "test@example.com",
                Password: "password123",
            },
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "邮箱格式错误",
            requestBody: RegisterRequest{
                Username: "testuser",
                Email:    "invalid-email",
                Password: "password123",
            },
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "密码太短",
            requestBody: RegisterRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "12345",
            },
            expectedStatus: http.StatusBadRequest,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 创建测试请求
            body, _ := json.Marshal(tt.requestBody)
            req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            // 创建响应记录器
            w := httptest.NewRecorder()
            
            // 执行请求
            router := gin.Default()
            router.POST("/api/v1/register", api.Register)
            router.ServeHTTP(w, req)
            
            // 验证响应
            assert.Equal(t, tt.expectedStatus, w.Code)
            
            if tt.expectedMsg != "" {
                var response map[string]interface{}
                json.Unmarshal(w.Body.Bytes(), &response)
                assert.Equal(t, tt.expectedMsg, response["message"])
            }
        })
    }
}
```

---

## 注意事项

### 安全性

1. **密码安全**
   - 密码使用 bcrypt 加密存储
   - 密码传输时确保使用 HTTPS
   - 建议密码长度至少 8 位，包含字母、数字和特殊字符
   - 不要在日志中记录密码

2. **Token 安全**
   - Token 应存储在 localStorage 或内存中
   - 不要在 URL 中传递 Token
   - Token 过期后需要重新登录
   - 建议定期刷新 Token

3. **输入验证**
   - 所有输入都经过严格验证
   - 防止 SQL 注入
   - 防止 XSS 攻击
   - 限制请求频率

### 性能优化

1. **数据库查询**
   - 用户列表查询使用分页
   - 合理使用索引
   - 避免 N+1 查询问题

2. **缓存策略**
   - 用户信息可以缓存
   - Token 验证可以使用缓存
   - 设置合理的缓存过期时间

### 错误处理

1. **客户端错误处理**
   - 401 错误：清除 Token，跳转到登录页
   - 403 错误：提示权限不足
   - 其他错误：显示友好的错误信息

2. **服务端错误处理**
   - 统一的错误响应格式
   - 记录详细的错误日志
   - 不向客户端暴露敏感信息

---

## 变更日志

### v1.0 (2025-10-15)
- ✅ 初始版本
- ✅ 实现用户注册和登录
- ✅ 实现个人信息管理
- ✅ 实现密码修改
- ✅ 实现管理员用户管理功能
- ✅ JWT Token 认证
- ✅ 基于角色的权限控制

---

**文档维护**: 如有问题或建议，请联系开发团队。  
**相关文档**: 
- [API设计规范](../API设计规范.md)
- [用户管理API使用指南](用户管理API使用指南.md)

