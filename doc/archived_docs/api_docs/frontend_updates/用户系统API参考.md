# 用户系统 API 参考

> **版本**: v1.0  
> **最后更新**: 2025-10-18  
> **模块**: 用户系统（User System）

## 📋 目录

- [概述](#概述)
- [认证接口](#认证接口)
- [用户接口](#用户接口)
- [管理员接口](#管理员接口)
- [数据结构](#数据结构)
- [错误码说明](#错误码说明)
- [示例代码](#示例代码)

---

## 概述

### 基础信息

- **Base URL**: `/api/v1`
- **认证方式**: JWT Token (Bearer认证)
- **Content-Type**: `application/json`

### 功能特性

✅ 用户注册与登录  
✅ JWT Token 认证  
✅ 用户信息管理  
✅ 密码管理  
✅ 基于角色的权限控制（RBAC）  
✅ 用户状态管理

---

## 认证接口

### 用户注册

创建新用户账号并自动登录。

**请求**

- **方法**: `POST`
- **路径**: `/register`
- **需要认证**: ❌ 否

**请求体**

```json
{
  "username": "string (required, 3-50字符)",
  "email": "string (required, email格式)",
  "password": "string (required, 最少6字符)"
}
```

**响应 - 成功 (201)**

```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "userID": "string",
    "username": "string",
    "email": "string",
    "role": "user",
    "status": "active",
    "token": "string (JWT Token)"
  }
}
```

**响应 - 失败**

```json
{
  "code": 400,
  "message": "注册失败",
  "error": "用户名已存在"
}
```

**cURL 示例**

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**JavaScript (Axios) 示例**

```javascript
import axios from 'axios';

async function register(username, email, password) {
  try {
    const response = await axios.post('/api/v1/register', {
      username,
      email,
      password
    });
    
    // 保存 Token
    localStorage.setItem('token', response.data.data.token);
    
    return response.data;
  } catch (error) {
    console.error('注册失败:', error.response?.data);
    throw error;
  }
}
```

---

### 用户登录

用户登录并获取 JWT Token。

**请求**

- **方法**: `POST`
- **路径**: `/login`
- **需要认证**: ❌ 否

**请求体**

```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "userID": "string",
    "username": "string",
    "email": "string",
    "token": "string (JWT Token)"
  }
}
```

**响应 - 失败**

```json
{
  "code": 401,
  "message": "用户名或密码错误"
}
```

**JavaScript (Fetch) 示例**

```javascript
async function login(username, password) {
  const response = await fetch('/api/v1/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });
  
  if (!response.ok) {
    throw new Error('登录失败');
  }
  
  const data = await response.json();
  
  // 保存 Token
  localStorage.setItem('token', data.data.token);
  localStorage.setItem('userId', data.data.userID);
  
  return data;
}
```

---

## 用户接口

### 获取个人信息

获取当前登录用户的详细信息。

**请求**

- **方法**: `GET`
- **路径**: `/users/profile`
- **需要认证**: ✅ 是

**Headers**

```
Authorization: Bearer <your_jwt_token>
```

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "userID": "string",
    "username": "string",
    "email": "string",
    "phone": "string",
    "role": "user",
    "status": "active",
    "avatar": "string (URL)",
    "nickname": "string",
    "bio": "string",
    "emailVerified": boolean,
    "phoneVerified": boolean,
    "lastLoginAt": "2025-10-18T10:30:00Z",
    "lastLoginIP": "192.168.1.1",
    "createdAt": "2025-01-01T00:00:00Z",
    "updatedAt": "2025-10-18T10:30:00Z"
  }
}
```

**响应 - 失败**

```json
{
  "code": 401,
  "message": "未认证"
}
```

**JavaScript 示例**

```javascript
async function getUserProfile() {
  const token = localStorage.getItem('token');
  
  const response = await axios.get('/api/v1/users/profile', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return response.data.data;
}
```

---

### 更新个人信息

更新当前登录用户的个人信息。

**请求**

- **方法**: `PUT`
- **路径**: `/users/profile`
- **需要认证**: ✅ 是

**请求体**（所有字段可选）

```json
{
  "nickname": "string (optional)",
  "bio": "string (optional)",
  "avatar": "string (optional, URL)",
  "phone": "string (optional)"
}
```

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "更新成功"
}
```

**响应 - 失败**

```json
{
  "code": 400,
  "message": "更新失败",
  "error": "手机号格式无效"
}
```

**JavaScript 示例**

```javascript
async function updateProfile(updates) {
  const token = localStorage.getItem('token');
  
  const response = await axios.put('/api/v1/users/profile', updates, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return response.data;
}

// 使用示例
await updateProfile({
  nickname: "青羽用户",
  bio: "热爱阅读和写作"
});
```

---

### 修改密码

修改当前登录用户的密码。

**请求**

- **方法**: `PUT`
- **路径**: `/users/password`
- **需要认证**: ✅ 是

**请求体**

```json
{
  "oldPassword": "string (required)",
  "newPassword": "string (required, 最少6字符)"
}
```

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "密码修改成功"
}
```

**响应 - 失败**

```json
{
  "code": 401,
  "message": "旧密码错误"
}
```

**JavaScript 示例**

```javascript
async function changePassword(oldPassword, newPassword) {
  const token = localStorage.getItem('token');
  
  const response = await axios.put(
    '/api/v1/users/password',
    {
      oldPassword,
      newPassword
    },
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  return response.data;
}
```

---

## 管理员接口

> **注意**: 以下接口需要管理员权限

### 获取用户列表

管理员获取用户列表，支持分页和筛选。

**请求**

- **方法**: `GET`
- **路径**: `/admin/users`
- **需要认证**: ✅ 是（管理员）

**Query 参数**

```
page        (int, optional)    页码，默认1
page_size   (int, optional)    每页数量，默认10
username    (string, optional)  用户名筛选
email       (string, optional)  邮箱筛选
role        (string, optional)  角色筛选 (user/author/admin)
status      (string, optional)  状态筛选 (active/inactive/banned)
```

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "userID": "string",
      "username": "string",
      "email": "string",
      "role": "user",
      "status": "active",
      "createdAt": "2025-01-01T00:00:00Z",
      "updatedAt": "2025-10-18T10:30:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "pageSize": 10
}
```

**JavaScript 示例**

```javascript
async function listUsers(page = 1, pageSize = 10, filters = {}) {
  const token = localStorage.getItem('token');
  
  const params = {
    page,
    page_size: pageSize,
    ...filters
  };
  
  const response = await axios.get('/api/v1/admin/users', {
    headers: {
      'Authorization': `Bearer ${token}`
    },
    params
  });
  
  return response.data;
}

// 使用示例
const { data: users, total } = await listUsers(1, 20, {
  role: 'user',
  status: 'active'
});
```

---

### 获取指定用户信息

管理员获取指定用户的详细信息。

**请求**

- **方法**: `GET`
- **路径**: `/admin/users/:id`
- **需要认证**: ✅ 是（管理员）

**Path 参数**

```
id (string, required)    用户ID
```

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "userID": "string",
    "username": "string",
    "email": "string",
    "phone": "string",
    "role": "user",
    "status": "active",
    "avatar": "string",
    "nickname": "string",
    "bio": "string",
    "emailVerified": boolean,
    "phoneVerified": boolean,
    "lastLoginAt": "2025-10-18T10:30:00Z",
    "lastLoginIP": "192.168.1.1",
    "createdAt": "2025-01-01T00:00:00Z",
    "updatedAt": "2025-10-18T10:30:00Z"
  }
}
```

**JavaScript 示例**

```javascript
async function getUserById(userId) {
  const token = localStorage.getItem('token');
  
  const response = await axios.get(`/api/v1/admin/users/${userId}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return response.data.data;
}
```

---

### 更新用户信息

管理员更新指定用户的信息。

**请求**

- **方法**: `PUT`
- **路径**: `/admin/users/:id`
- **需要认证**: ✅ 是（管理员）

**请求体**（所有字段可选）

```json
{
  "nickname": "string (optional)",
  "bio": "string (optional)",
  "avatar": "string (optional)",
  "phone": "string (optional)",
  "role": "string (optional, user/author/admin)",
  "status": "string (optional, active/inactive/banned)",
  "emailVerified": boolean (optional),
  "phoneVerified": boolean (optional)
}
```

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "更新成功"
}
```

**JavaScript 示例**

```javascript
async function updateUser(userId, updates) {
  const token = localStorage.getItem('token');
  
  const response = await axios.put(
    `/api/v1/admin/users/${userId}`,
    updates,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  
  return response.data;
}

// 使用示例 - 禁用用户
await updateUser('user-123', {
  status: 'banned'
});
```

---

### 删除用户

管理员删除指定用户。

**请求**

- **方法**: `DELETE`
- **路径**: `/admin/users/:id`
- **需要认证**: ✅ 是（管理员）

**响应 - 成功 (200)**

```json
{
  "code": 200,
  "message": "删除成功"
}
```

**响应 - 失败**

```json
{
  "code": 404,
  "message": "用户不存在"
}
```

**JavaScript 示例**

```javascript
async function deleteUser(userId) {
  const token = localStorage.getItem('token');
  
  const confirmed = confirm('确定要删除该用户吗？');
  if (!confirmed) return;
  
  const response = await axios.delete(`/api/v1/admin/users/${userId}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return response.data;
}
```

---

## 数据结构

### User 对象

```typescript
interface User {
  userID: string;           // 用户ID
  username: string;         // 用户名
  email: string;            // 邮箱
  phone?: string;           // 手机号
  role: 'user' | 'author' | 'admin';  // 角色
  status: 'active' | 'inactive' | 'banned';  // 状态
  avatar?: string;          // 头像URL
  nickname?: string;        // 昵称
  bio?: string;             // 个人简介
  emailVerified: boolean;   // 邮箱是否验证
  phoneVerified: boolean;   // 手机号是否验证
  lastLoginAt?: string;     // 最后登录时间 (ISO 8601)
  lastLoginIP?: string;     // 最后登录IP
  createdAt: string;        // 创建时间 (ISO 8601)
  updatedAt: string;        // 更新时间 (ISO 8601)
}
```

### UserRole 枚举

| 值 | 说明 | 权限 |
|----|------|------|
| `user` | 普通用户 | 基础阅读权限 |
| `author` | 作者 | 用户权限 + 写作权限 |
| `admin` | 管理员 | 所有权限 |

### UserStatus 枚举

| 值 | 说明 |
|----|------|
| `active` | 正常状态 |
| `inactive` | 未激活 |
| `banned` | 已封禁 |

---

## 错误码说明

### 通用错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 参数错误 | 检查请求参数格式和必填项 |
| 401 | 未认证 | 提供有效的 JWT Token |
| 403 | 权限不足 | 需要更高权限（如管理员） |
| 404 | 资源不存在 | 确认用户ID是否正确 |
| 409 | 资源冲突 | 用户名或邮箱已存在 |
| 500 | 服务器错误 | 稍后重试或联系管理员 |

### 业务错误码

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| "用户名已存在" | 注册时用户名重复 | 更换用户名 |
| "邮箱已被注册" | 注册时邮箱重复 | 更换邮箱或找回密码 |
| "用户名或密码错误" | 登录凭证无效 | 检查输入是否正确 |
| "旧密码错误" | 修改密码时原密码不正确 | 确认原密码 |
| "Token过期" | JWT Token 已过期 | 重新登录或刷新 Token |

---

## 示例代码

### 完整的认证流程

```javascript
// auth.js
import axios from 'axios';

const API_BASE = '/api/v1';

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: API_BASE,
  timeout: 30000
});

// 请求拦截器 - 添加 Token
apiClient.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器 - 处理错误
apiClient.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      // Token 过期，清除并跳转登录
      localStorage.clear();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 注册
export async function register(username, email, password) {
  const data = await apiClient.post('/register', {
    username,
    email,
    password
  });
  
  // 保存 Token
  localStorage.setItem('token', data.data.token);
  localStorage.setItem('userId', data.data.userID);
  
  return data;
}

// 登录
export async function login(username, password) {
  const data = await apiClient.post('/login', {
    username,
    password
  });
  
  // 保存 Token
  localStorage.setItem('token', data.data.token);
  localStorage.setItem('userId', data.data.userID);
  
  return data;
}

// 登出
export async function logout() {
  await apiClient.post('/shared/auth/logout');
  localStorage.clear();
}

// 获取用户信息
export function getUserProfile() {
  return apiClient.get('/users/profile');
}

// 更新用户信息
export function updateUserProfile(updates) {
  return apiClient.put('/users/profile', updates);
}

// 修改密码
export function changePassword(oldPassword, newPassword) {
  return apiClient.put('/users/password', {
    oldPassword,
    newPassword
  });
}

export default apiClient;
```

### Vue 3 组件示例

```vue
<!-- LoginPage.vue -->
<template>
  <div class="login-page">
    <h1>用户登录</h1>
    <form @submit.prevent="handleLogin">
      <input
        v-model="form.username"
        type="text"
        placeholder="用户名"
        required
      />
      <input
        v-model="form.password"
        type="password"
        placeholder="密码"
        required
      />
      <button type="submit" :disabled="loading">
        {{ loading ? '登录中...' : '登录' }}
      </button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { login } from '@/api/auth';

const router = useRouter();
const form = ref({
  username: '',
  password: ''
});
const loading = ref(false);
const error = ref('');

const handleLogin = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    await login(form.value.username, form.value.password);
    router.push('/');
  } catch (err) {
    error.value = err.response?.data?.message || '登录失败';
  } finally {
    loading.value = false;
  }
};
</script>
```

---

## 相关文档

- [前端集成指南](./前端集成指南.md) - 完整集成教程
- [API 快速参考](./API快速参考.md) - 一页纸速查表
- [认证 API 参考](./认证API参考.md) - 共享认证服务

---

**维护者**: 青羽后端团队  
**联系方式**: backend@qingyu.com

