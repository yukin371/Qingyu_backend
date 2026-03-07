# 用户管理 API 使用指南

**版本**: v1.0  
**最后更新**: 2025-10-13  
**基础路径**: `/api/v1`

---

## 📋 目录

1. [认证说明](#认证说明)
2. [公开接口](#公开接口)
   - [用户注册](#用户注册)
   - [用户登录](#用户登录)
3. [需要认证的接口](#需要认证的接口)
   - [获取个人信息](#获取个人信息)
   - [更新个人信息](#更新个人信息)
   - [修改密码](#修改密码)
4. [管理员接口](#管理员接口)
   - [获取用户列表](#获取用户列表)
   - [获取指定用户](#获取指定用户)
   - [更新用户信息](#更新用户信息)
   - [删除用户](#删除用户)
5. [错误码说明](#错误码说明)
6. [数据结构](#数据结构)

---

## 认证说明

### JWT Token 认证

所有需要认证的接口都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <your_jwt_token>
```

### 获取 Token

通过 [用户注册](#用户注册) 或 [用户登录](#用户登录) 接口获取 Token。

### Token 有效期

- 默认有效期：24 小时
- Token 过期后需要重新登录

### 权限说明

| 角色 | 权限 |
|------|------|
| user | 普通用户，可访问个人信息相关接口 |
| author | 作者，拥有用户权限 + 写作权限 |
| admin | 管理员，拥有所有权限 |

---

## 公开接口

### 用户注册

创建新用户账号。

**接口**: `POST /register`  
**需要认证**: ❌

#### 请求参数

```json
{
  "username": "testuser",      // 必填，3-50字符
  "email": "test@example.com", // 必填，有效邮箱
  "password": "password123"    // 必填，6-100字符
}
```

#### 成功响应

```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "user_id": "670abcdef123456789",
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "status": "active",
    "email_verified": false,
    "phone_verified": false,
    "created_at": "2025-10-13T10:00:00Z",
    "updated_at": "2025-10-13T10:00:00Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": 1697203200
}
```

#### 错误响应

**400 - 参数错误**:
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "用户名长度必须在3-50之间"
}
```

**409 - 邮箱已存在**:
```json
{
  "code": 409,
  "message": "注册失败",
  "error": "邮箱已被注册"
}
```

#### cURL 示例

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

---

### 用户登录

用户登录获取 Token。

**接口**: `POST /login`  
**需要认证**: ❌

#### 请求参数

```json
{
  "username": "testuser",   // 必填，用户名或邮箱
  "password": "password123" // 必填
}
```

#### 成功响应

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user_id": "670abcdef123456789",
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "status": "active",
    "last_login_at": "2025-10-13T10:30:00Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": 1697203200
}
```

#### 错误响应

**401 - 用户名或密码错误**:
```json
{
  "code": 401,
  "message": "登录失败",
  "error": "用户名或密码错误"
}
```

**403 - 账号被禁用**:
```json
{
  "code": 403,
  "message": "登录失败",
  "error": "账号已被禁用"
}
```

#### cURL 示例

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

---

## 需要认证的接口

### 获取个人信息

获取当前登录用户的详细信息。

**接口**: `GET /users/profile`  
**需要认证**: ✅ (JWT Token)

#### 请求参数

无

#### 成功响应

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "user_id": "670abcdef123456789",
    "username": "testuser",
    "email": "test@example.com",
    "phone": "",
    "role": "user",
    "status": "active",
    "avatar": "",
    "nickname": "",
    "bio": "",
    "email_verified": false,
    "phone_verified": false,
    "last_login_at": "2025-10-13T10:30:00Z",
    "last_login_ip": "192.168.1.100",
    "created_at": "2025-10-13T10:00:00Z",
    "updated_at": "2025-10-13T10:30:00Z"
  },
  "timestamp": 1697203200
}
```

#### 错误响应

**401 - 未认证**:
```json
{
  "code": 40101,
  "message": "未提供认证令牌",
  "data": null
}
```

**401 - Token 无效**:
```json
{
  "code": 40103,
  "message": "无效的认证令牌",
  "data": null
}
```

#### cURL 示例

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### 更新个人信息

更新当前用户的个人信息。

**接口**: `PUT /users/profile`  
**需要认证**: ✅ (JWT Token)

#### 请求参数

```json
{
  "nickname": "我的昵称",        // 可选，0-50字符
  "avatar": "https://...",     // 可选，头像URL
  "bio": "这是我的个人简介",    // 可选，0-500字符
  "phone": "13800138000"       // 可选，手机号
}
```

#### 成功响应

```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "user_id": "670abcdef123456789",
    "username": "testuser",
    "nickname": "我的昵称",
    "avatar": "https://...",
    "bio": "这是我的个人简介",
    "phone": "13800138000",
    "updated_at": "2025-10-13T11:00:00Z"
  },
  "timestamp": 1697207200
}
```

#### 错误响应

**400 - 参数错误**:
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "昵称长度不能超过50字符"
}
```

#### cURL 示例

```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "我的昵称",
    "bio": "这是我的个人简介"
  }'
```

---

### 修改密码

修改当前用户的登录密码。

**接口**: `PUT /users/password`  
**需要认证**: ✅ (JWT Token)

#### 请求参数

```json
{
  "old_password": "oldpassword123", // 必填，当前密码
  "new_password": "newpassword456"  // 必填，6-100字符
}
```

#### 成功响应

```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null,
  "timestamp": 1697207200
}
```

#### 错误响应

**400 - 旧密码错误**:
```json
{
  "code": 400,
  "message": "修改密码失败",
  "error": "旧密码错误"
}
```

**400 - 新密码不符合要求**:
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "新密码长度不能少于6位"
}
```

#### cURL 示例

```bash
curl -X PUT http://localhost:8080/api/v1/users/password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "oldpassword123",
    "new_password": "newpassword456"
  }'
```

---

## 管理员接口

### 获取用户列表

获取系统中的用户列表（分页）。

**接口**: `GET /admin/users`  
**需要认证**: ✅ (JWT Token + Admin 角色)

#### 请求参数（Query）

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| page | int | 否 | 页码（从1开始） | 1 |
| page_size | int | 否 | 每页数量 | 10 |
| role | string | 否 | 筛选角色 | - |
| status | string | 否 | 筛选状态 | - |
| keyword | string | 否 | 搜索关键词 | - |

#### 成功响应

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "users": [
      {
        "user_id": "670abcdef123456789",
        "username": "testuser",
        "email": "test@example.com",
        "role": "user",
        "status": "active",
        "created_at": "2025-10-13T10:00:00Z"
      },
      {
        "user_id": "670abcdef987654321",
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin",
        "status": "active",
        "created_at": "2025-10-01T08:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10
  },
  "timestamp": 1697207200
}
```

#### 错误响应

**403 - 权限不足**:
```json
{
  "code": 40301,
  "message": "权限不足",
  "data": null
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

# 搜索关键词
curl -X GET "http://localhost:8080/api/v1/admin/users?keyword=test" \
  -H "Authorization: Bearer <admin_token>"
```

---

### 获取指定用户

获取指定用户的详细信息。

**接口**: `GET /admin/users/:id`  
**需要认证**: ✅ (JWT Token + Admin 角色)

#### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 用户ID |

#### 成功响应

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "user_id": "670abcdef123456789",
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13800138000",
    "role": "user",
    "status": "active",
    "avatar": "https://...",
    "nickname": "测试用户",
    "bio": "这是测试用户",
    "email_verified": true,
    "phone_verified": false,
    "last_login_at": "2025-10-13T10:30:00Z",
    "last_login_ip": "192.168.1.100",
    "created_at": "2025-10-13T10:00:00Z",
    "updated_at": "2025-10-13T10:30:00Z"
  },
  "timestamp": 1697207200
}
```

#### 错误响应

**404 - 用户不存在**:
```json
{
  "code": 404,
  "message": "用户不存在",
  "error": "未找到指定用户"
}
```

#### cURL 示例

```bash
curl -X GET http://localhost:8080/api/v1/admin/users/670abcdef123456789 \
  -H "Authorization: Bearer <admin_token>"
```

---

### 更新用户信息

管理员更新指定用户的信息。

**接口**: `PUT /admin/users/:id`  
**需要认证**: ✅ (JWT Token + Admin 角色)

#### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 用户ID |

#### 请求参数

```json
{
  "role": "author",      // 可选，用户角色
  "status": "banned",    // 可选，用户状态
  "nickname": "新昵称",   // 可选
  "avatar": "https://...", // 可选
  "bio": "新简介"         // 可选
}
```

#### 成功响应

```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "user_id": "670abcdef123456789",
    "username": "testuser",
    "role": "author",
    "status": "banned",
    "nickname": "新昵称",
    "updated_at": "2025-10-13T11:30:00Z"
  },
  "timestamp": 1697209800
}
```

#### 错误响应

**400 - 参数错误**:
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "无效的角色类型"
}
```

#### cURL 示例

```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/670abcdef123456789 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "author",
    "status": "active"
  }'
```

---

### 删除用户

删除指定用户（软删除）。

**接口**: `DELETE /admin/users/:id`  
**需要认证**: ✅ (JWT Token + Admin 角色)

#### 路径参数

| 参数 | 类型 | 说明 |
|------|------|------|
| id | string | 用户ID |

#### 成功响应

```json
{
  "code": 200,
  "message": "删除成功",
  "data": null,
  "timestamp": 1697209800
}
```

#### 错误响应

**404 - 用户不存在**:
```json
{
  "code": 404,
  "message": "用户不存在",
  "error": "未找到指定用户"
}
```

#### cURL 示例

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/670abcdef123456789 \
  -H "Authorization: Bearer <admin_token>"
```

---

## 错误码说明

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 参数错误 |
| 401 | 未认证 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器错误 |

### 业务错误码

| 错误码 | 说明 |
|--------|------|
| 40101 | 未提供认证令牌 |
| 40102 | 无效的认证令牌格式 |
| 40103 | 无效的认证令牌 |
| 40301 | 权限不足 |

---

## 数据结构

### User（用户对象）

```typescript
interface User {
  user_id: string;           // 用户ID
  username: string;          // 用户名
  email: string;             // 邮箱
  phone?: string;            // 手机号（可选）
  role: string;              // 角色：user/author/admin
  status: string;            // 状态：active/inactive/banned/deleted
  avatar?: string;           // 头像URL（可选）
  nickname?: string;         // 昵称（可选）
  bio?: string;              // 个人简介（可选）
  email_verified: boolean;   // 邮箱是否验证
  phone_verified: boolean;   // 手机是否验证
  last_login_at?: string;    // 最后登录时间（可选）
  last_login_ip?: string;    // 最后登录IP（可选）
  created_at: string;        // 创建时间
  updated_at: string;        // 更新时间
}
```

### 角色说明

| 角色 | 值 | 说明 |
|------|-----|------|
| 普通用户 | user | 基础权限 |
| 作者 | author | 用户权限 + 写作权限 |
| 管理员 | admin | 所有权限 |

### 用户状态

| 状态 | 值 | 说明 |
|------|-----|------|
| 正常 | active | 正常使用 |
| 未激活 | inactive | 未激活，不能登录 |
| 封禁 | banned | 已封禁，不能登录 |
| 已删除 | deleted | 已删除（软删除） |

---

## 前端集成示例

### Vue 3 + Axios

```javascript
// api/user.js
import axios from 'axios';

const BASE_URL = 'http://localhost:8080/api/v1';

// 设置 Token
export function setAuthToken(token) {
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  } else {
    delete axios.defaults.headers.common['Authorization'];
  }
}

// 用户注册
export async function register(data) {
  const response = await axios.post(`${BASE_URL}/register`, data);
  const token = response.data.data.token;
  setAuthToken(token);
  return response.data;
}

// 用户登录
export async function login(data) {
  const response = await axios.post(`${BASE_URL}/login`, data);
  const token = response.data.data.token;
  setAuthToken(token);
  return response.data;
}

// 获取个人信息
export async function getProfile() {
  const response = await axios.get(`${BASE_URL}/users/profile`);
  return response.data;
}

// 更新个人信息
export async function updateProfile(data) {
  const response = await axios.put(`${BASE_URL}/users/profile`, data);
  return response.data;
}

// 修改密码
export async function changePassword(data) {
  const response = await axios.put(`${BASE_URL}/users/password`, data);
  return response.data;
}

// 管理员 - 获取用户列表
export async function getUserList(params) {
  const response = await axios.get(`${BASE_URL}/admin/users`, { params });
  return response.data;
}
```

### React + Fetch

```javascript
// utils/api.js
const BASE_URL = 'http://localhost:8080/api/v1';

// 获取 Token
export function getToken() {
  return localStorage.getItem('token');
}

// 保存 Token
export function saveToken(token) {
  localStorage.setItem('token', token);
}

// 删除 Token
export function removeToken() {
  localStorage.removeItem('token');
}

// 统一请求方法
async function request(url, options = {}) {
  const token = getToken();
  
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  const response = await fetch(`${BASE_URL}${url}`, {
    ...options,
    headers,
  });
  
  const data = await response.json();
  
  if (!response.ok) {
    throw new Error(data.message || data.error);
  }
  
  return data;
}

// 用户注册
export async function register(data) {
  const result = await request('/register', {
    method: 'POST',
    body: JSON.stringify(data),
  });
  saveToken(result.data.token);
  return result;
}

// 用户登录
export async function login(data) {
  const result = await request('/login', {
    method: 'POST',
    body: JSON.stringify(data),
  });
  saveToken(result.data.token);
  return result;
}

// 获取个人信息
export async function getProfile() {
  return request('/users/profile');
}
```

---

## 测试建议

### Postman Collection

建议创建 Postman Collection，包含以下环境变量：

```json
{
  "baseUrl": "http://localhost:8080/api/v1",
  "token": "",
  "userId": ""
}
```

### 测试流程

1. **注册新用户** → 保存 Token
2. **登录** → 验证 Token
3. **获取个人信息** → 验证认证
4. **更新个人信息** → 验证更新
5. **修改密码** → 验证密码更新
6. **管理员登录** → 保存 Admin Token
7. **获取用户列表** → 验证权限
8. **更新其他用户** → 验证管理员权限

---

## 注意事项

1. **密码安全**:
   - 密码在传输前不需要加密（HTTPS 负责传输安全）
   - 密码存储时使用 bcrypt 加密
   - 建议密码长度至少 8 位，包含字母和数字

2. **Token 管理**:
   - Token 应存储在 localStorage 或内存中
   - 不要在 URL 中传递 Token
   - Token 过期后需要重新登录

3. **CORS 配置**:
   - 确保后端已配置 CORS
   - 允许前端域名访问

4. **错误处理**:
   - 401 错误：清除 Token，跳转到登录页
   - 403 错误：提示权限不足
   - 其他错误：显示友好的错误信息

---

**文档维护**: 如有问题或建议，请联系开发团队。

