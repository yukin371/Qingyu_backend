# Postman API 测试指南

本文档提供使用 Postman 测试青羽后端 API 的详细步骤和示例。

## 目录
- [环境配置](#环境配置)
- [用户系统 API 测试](#用户系统-api-测试)
- [认证流程](#认证流程)
- [书城 API 测试](#书城-api-测试)
- [常见问题](#常见问题)

---

## 环境配置

### 1. 基础设置

**Base URL**: `http://localhost:8080`

建议在 Postman 中创建环境变量：
- 变量名: `base_url`
- 值: `http://localhost:8080`
- 变量名: `token`
- 值: (登录后自动设置)

### 2. 全局请求头

对于需要认证的请求，添加以下请求头：
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

---

## 用户系统 API 测试

### 1. 用户注册

**请求方式**: POST  
**URL**: `{{base_url}}/api/v1/register`  
**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456"
}
```

**成功响应** (201 Created):
```json
{
  "code": 201,
  "message": "用户注册成功",
  "data": {
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "",
      "avatar": "",
      "bio": "",
      "role": "user",
      "status": "active",
      "created_at": "2025-10-15T10:30:00Z",
      "updated_at": "2025-10-15T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应** (400 Bad Request):
```json
{
  "code": 400,
  "message": "请求参数错误",
  "error": "用户名已存在"
}
```

**Postman 测试脚本** (Tests 标签):
```javascript
// 验证状态码
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

// 验证响应结构
pm.test("Response has user and token", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('user');
    pm.expect(jsonData.data).to.have.property('token');
});

// 保存 token 到环境变量
var jsonData = pm.response.json();
if (jsonData.data && jsonData.data.token) {
    pm.environment.set("token", jsonData.data.token);
}
```

---

### 2. 用户登录

**请求方式**: POST  
**URL**: `{{base_url}}/api/v1/login`  
**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "testuser",
  "password": "123456"
}
```

**成功响应** (200 OK):
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "username": "testuser",
      "email": "test@example.com",
      "nickname": "测试用户",
      "avatar": "",
      "bio": "",
      "role": "user",
      "status": "active",
      "created_at": "2025-10-15T10:30:00Z",
      "updated_at": "2025-10-15T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应** (401 Unauthorized):
```json
{
  "code": 401,
  "message": "认证失败",
  "error": "用户名或密码错误"
}
```

**Postman 测试脚本**:
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Login successful", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.code).to.eql(200);
    pm.expect(jsonData.data.token).to.be.a('string');
});

// 保存 token
var jsonData = pm.response.json();
if (jsonData.data && jsonData.data.token) {
    pm.environment.set("token", jsonData.data.token);
    console.log("Token saved:", jsonData.data.token);
}
```

---

### 3. 获取用户信息

**请求方式**: GET  
**URL**: `{{base_url}}/api/v1/users/profile`  
**请求头**:
```
Authorization: Bearer {{token}}
```

**成功响应** (200 OK):
```json
{
  "code": 200,
  "message": "获取用户信息成功",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "这是我的个人简介",
    "role": "user",
    "status": "active",
    "created_at": "2025-10-15T10:30:00Z",
    "updated_at": "2025-10-15T10:30:00Z"
  }
}
```

**Postman 测试脚本**:
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Profile data is complete", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('id');
    pm.expect(jsonData.data).to.have.property('username');
    pm.expect(jsonData.data).to.have.property('email');
});
```

---

### 4. 更新用户信息

**请求方式**: PUT  
**URL**: `{{base_url}}/api/v1/users/profile`  
**请求头**:
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**请求体**:
```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/new-avatar.jpg",
  "bio": "更新后的个人简介"
}
```

**成功响应** (200 OK):
```json
{
  "code": 200,
  "message": "更新用户信息成功",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "新昵称",
    "avatar": "https://example.com/new-avatar.jpg",
    "bio": "更新后的个人简介",
    "role": "user",
    "status": "active",
    "created_at": "2025-10-15T10:30:00Z",
    "updated_at": "2025-10-15T11:00:00Z"
  }
}
```

---

### 5. 修改密码

**请求方式**: PUT  
**URL**: `{{base_url}}/api/v1/users/password`  
**请求头**:
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**请求体**:
```json
{
  "old_password": "123456",
  "new_password": "newpassword123"
}
```

**成功响应** (200 OK):
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

**错误响应** (400 Bad Request):
```json
{
  "code": 400,
  "message": "修改密码失败",
  "error": "原密码错误"
}
```

---

### 6. 管理员：获取用户列表

**请求方式**: GET  
**URL**: `{{base_url}}/api/v1/admin/users?page=1&page_size=10&role=user&status=active`  
**请求头**:
```
Authorization: Bearer {{token}}
```

**查询参数**:
- `page`: 页码 (可选，默认 1)
- `page_size`: 每页数量 (可选，默认 10)
- `role`: 角色过滤 (可选: user, admin)
- `status`: 状态过滤 (可选: active, inactive, banned)
- `keyword`: 搜索关键词 (可选)

**成功响应** (200 OK):
```json
{
  "code": 200,
  "message": "获取用户列表成功",
  "data": {
    "users": [
      {
        "id": "507f1f77bcf86cd799439011",
        "username": "testuser",
        "email": "test@example.com",
        "nickname": "测试用户",
        "role": "user",
        "status": "active",
        "created_at": "2025-10-15T10:30:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_pages": 1,
      "total_items": 1
    }
  }
}
```

---

### 7. 管理员：获取指定用户

**请求方式**: GET  
**URL**: `{{base_url}}/api/v1/admin/users/:userId`  
**请求头**:
```
Authorization: Bearer {{token}}
```

**路径参数**:
- `userId`: 用户ID

**示例**: `{{base_url}}/api/v1/admin/users/507f1f77bcf86cd799439011`

---

### 8. 管理员：更新用户

**请求方式**: PUT  
**URL**: `{{base_url}}/api/v1/admin/users/:userId`  
**请求头**:
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**请求体**:
```json
{
  "role": "admin",
  "status": "active"
}
```

---

### 9. 管理员：删除用户

**请求方式**: DELETE  
**URL**: `{{base_url}}/api/v1/admin/users/:userId`  
**请求头**:
```
Authorization: Bearer {{token}}
```

**成功响应** (200 OK):
```json
{
  "code": 200,
  "message": "删除用户成功",
  "data": null
}
```

---

## 认证流程

### 完整测试流程

1. **注册新用户** → 获得 token
2. **登录** → 获得 token (保存到环境变量)
3. **获取用户信息** → 验证 token 有效
4. **更新用户信息** → 测试修改功能
5. **修改密码** → 测试密码更新
6. **使用新密码登录** → 验证密码已更改

### Token 管理

**自动保存 Token** - 在登录/注册请求的 Tests 标签添加：
```javascript
var jsonData = pm.response.json();
if (jsonData.data && jsonData.data.token) {
    pm.environment.set("token", jsonData.data.token);
}
```

**自动添加 Authorization 头** - 在请求的 Headers 标签添加：
```
Authorization: Bearer {{token}}
```

---

## 书城 API 测试

### 1. 获取推荐书籍

**请求方式**: GET  
**URL**: `{{base_url}}/api/v1/public/bookstore/recommend`

**成功响应** (200 OK):
```json
{
  "code": 200,
  "message": "获取推荐书籍成功",
  "data": {
    "books": [
      {
        "id": "book123",
        "title": "示例书籍",
        "author": "作者",
        "cover": "https://example.com/cover.jpg",
        "description": "书籍描述"
      }
    ]
  }
}
```

### 2. 搜索书籍

**请求方式**: GET  
**URL**: `{{base_url}}/api/v1/public/bookstore/search?keyword=红楼梦&page=1&page_size=10`

**查询参数**:
- `keyword`: 搜索关键词
- `page`: 页码 (可选)
- `page_size`: 每页数量 (可选)

---

## 常见问题

### 1. 401 Unauthorized 错误

**原因**: Token 无效或未提供

**解决方案**:
1. 确保已登录并获取 token
2. 检查 Authorization 头格式: `Bearer {{token}}`
3. 重新登录获取新 token

### 2. 400 Bad Request 错误

**原因**: 请求参数不符合要求

**解决方案**:
1. 检查请求体 JSON 格式是否正确
2. 检查必填字段是否都提供了
3. 检查字段值是否符合验证规则（如密码最少6位）
4. 查看响应中的 `error` 字段获取具体错误信息

### 3. 403 Forbidden 错误

**原因**: 没有权限访问该资源

**解决方案**:
1. 检查用户角色是否足够（如管理员接口需要 admin 角色）
2. 确认用户状态是否为 active

### 4. 500 Internal Server Error

**原因**: 服务器内部错误

**解决方案**:
1. 检查后端服务是否正常运行
2. 查看后端日志获取详细错误信息
3. 检查数据库连接是否正常

---

## Postman Collection 导入

为了方便测试，可以创建一个 Postman Collection，包含以上所有请求。

### Collection 结构建议

```
青羽 API 测试
├── 用户系统
│   ├── 公开接口
│   │   ├── 注册
│   │   └── 登录
│   ├── 认证接口
│   │   ├── 获取用户信息
│   │   ├── 更新用户信息
│   │   └── 修改密码
│   └── 管理员接口
│       ├── 获取用户列表
│       ├── 获取指定用户
│       ├── 更新用户
│       └── 删除用户
├── 书城系统
│   ├── 获取推荐书籍
│   ├── 搜索书籍
│   └── 获取书籍详情
└── 健康检查
    └── Ping
```

### 环境变量设置

```json
{
  "base_url": "http://localhost:8080",
  "token": "",
  "userId": ""
}
```

---

## 测试用例示例

### 用户注册并登录流程

1. **POST** `/api/v1/register` - 注册新用户
   - 验证状态码为 201
   - 验证返回 token
   - 保存 token 到环境变量

2. **POST** `/api/v1/login` - 登录
   - 验证状态码为 200
   - 验证返回用户信息和 token

3. **GET** `/api/v1/users/profile` - 获取个人信息
   - 使用保存的 token
   - 验证返回的用户信息正确

4. **PUT** `/api/v1/users/profile` - 更新个人信息
   - 验证更新成功
   - 验证返回的数据已更新

5. **PUT** `/api/v1/users/password` - 修改密码
   - 验证修改成功

6. **POST** `/api/v1/login` - 使用新密码登录
   - 验证可以成功登录

---

## 附录

### HTTP 状态码说明

- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未授权或认证失败
- `403 Forbidden`: 没有权限访问
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误

### 验证规则

**用户名**:
- 必填
- 长度: 3-50 字符
- 唯一性

**邮箱**:
- 必填
- 必须是有效的邮箱格式
- 唯一性

**密码**:
- 必填
- 最少 6 个字符

**角色**:
- `user`: 普通用户
- `admin`: 管理员

**状态**:
- `active`: 正常
- `inactive`: 未激活
- `banned`: 已封禁


