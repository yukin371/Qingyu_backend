# 认证系统 API 参考

> **版本**: v1.0  
> **最后更新**: 2025-10-18  
> **基础路径**: `/api/v1/shared/auth`

---

## 1. 概述

认证系统提供用户注册、登录、Token管理、权限验证等核心安全功能，是整个平台的安全基础设施。

### 1.1 基础信息

- **认证方式**: JWT (JSON Web Token)
- **Token 类型**: Bearer Token
- **Token 有效期**: 24小时（可配置）
- **刷新机制**: 支持 Token 刷新

### 1.2 功能特性

- ✅ 用户注册和登录
- ✅ JWT Token 生成和验证
- ✅ Token 自动刷新
- ✅ 用户登出（Token 失效）
- ✅ 权限获取
- ✅ 角色管理

---

## 2. 认证流程

### 2.1 完整认证流程图

```
┌─────────────┐
│  1. 注册/登录 │
└──────┬──────┘
       │
       ↓
┌─────────────┐
│  2. 获取Token │
└──────┬──────┘
       │
       ↓
┌─────────────────┐
│  3. 存储Token    │
│  (localStorage)  │
└──────┬──────────┘
       │
       ↓
┌───────────────────┐
│  4. 携带Token请求  │
│  (Authorization)   │
└──────┬────────────┘
       │
       ↓
┌─────────────────┐     Token过期?     ┌─────────────┐
│  5. 服务器验证   │ ─────Yes────────→ │  刷新Token   │
└──────┬──────────┘                    └──────┬──────┘
       │                                      │
       No                                     │
       │                                      │
       ↓                                      ↓
┌─────────────┐                      ┌─────────────┐
│  返回数据    │ ←────────────────── │  重新请求    │
└─────────────┘                      └─────────────┘
```

---

## 3. 接口列表

### 3.1 用户注册

**接口说明**: 注册新用户账号

**请求**
```
POST /api/v1/shared/auth/register
```

**认证**: ❌ 无需认证

**请求体**
```json
{
  "username": "test_user",
  "email": "user@example.com",
  "password": "Password123!",
  "confirmPassword": "Password123!"
}
```

**请求参数说明**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（3-50字符） |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码（6-32字符） |
| confirmPassword | string | 是 | 确认密码（需与password一致） |

**密码要求**:
- 长度：6-32 字符
- 至少包含一个大写字母
- 至少包含一个小写字母
- 至少包含一个数字
- 可选：特殊字符

**响应示例**
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "id": "user123",
    "email": "user@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 86400
  }
}
```

**错误响应**
- `400` - 参数验证失败（邮箱格式错误、密码不符合要求等）
- `409` - 用户已存在（邮箱或用户名重复）
- `500` - 服务器错误

**cURL 示例**
```bash
curl -X POST "http://localhost:8080/api/v1/shared/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "email": "user@example.com",
    "password": "Password123!",
    "confirmPassword": "Password123!"
  }'
```

**JavaScript/Axios 示例**
```javascript
const register = async (userData) => {
  try {
    const response = await axios.post('/api/v1/shared/auth/register', {
      username: userData.username,
      email: userData.email,
      password: userData.password,
      confirmPassword: userData.confirmPassword
    });
    
    // 保存 Token
    const { token, expiresIn } = response.data.data;
    localStorage.setItem('token', token);
    localStorage.setItem('tokenExpiry', Date.now() + expiresIn * 1000);
    
    return response.data.data;
  } catch (error) {
    if (error.response?.status === 409) {
      throw new Error('用户已存在');
    } else if (error.response?.status === 400) {
      throw new Error(error.response.data.message || '参数验证失败');
    }
    throw error;
  }
};
```

---

### 3.2 用户登录

**接口说明**: 用户登录获取访问令牌

**请求**
```
POST /api/v1/shared/auth/login
```

**认证**: ❌ 无需认证

**请求体**
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

**请求参数说明**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 邮箱地址或用户名 |
| password | string | 是 | 密码 |

**响应示例**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "userId": "user123",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 86400,
    "user": {
      "id": "user123",
      "username": "test_user",
      "email": "user@example.com",
      "role": "user",
      "avatar": "https://example.com/avatar.jpg"
    }
  }
}
```

**错误响应**
- `400` - 参数错误
- `401` - 邮箱或密码错误
- `403` - 账号已被禁用
- `500` - 服务器错误

**JavaScript/Axios 示例**
```javascript
const login = async (email, password) => {
  try {
    const response = await axios.post('/api/v1/shared/auth/login', {
      email,
      password
    });
    
    const { token, expiresIn, user } = response.data.data;
    
    // 保存 Token 和用户信息
    localStorage.setItem('token', token);
    localStorage.setItem('tokenExpiry', Date.now() + expiresIn * 1000);
    localStorage.setItem('user', JSON.stringify(user));
    
    return response.data.data;
  } catch (error) {
    if (error.response?.status === 401) {
      throw new Error('邮箱或密码错误');
    } else if (error.response?.status === 403) {
      throw new Error('账号已被禁用');
    }
    throw error;
  }
};
```

---

### 3.3 用户登出

**接口说明**: 用户登出，使当前Token失效

**请求**
```
POST /api/v1/shared/auth/logout
```

**认证**: 🔒 需要 JWT Token

**请求头**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**
```json
{
  "code": 200,
  "message": "登出成功"
}
```

**错误响应**
- `401` - 未提供Token或Token无效
- `500` - 服务器错误

**JavaScript/Axios 示例**
```javascript
const logout = async () => {
  try {
    await axios.post('/api/v1/shared/auth/logout');
    
    // 清除本地存储
    localStorage.removeItem('token');
    localStorage.removeItem('tokenExpiry');
    localStorage.removeItem('user');
    
    // 跳转到登录页
    router.push('/login');
  } catch (error) {
    console.error('登出失败:', error);
    // 即使失败也清除本地数据
    localStorage.clear();
    router.push('/login');
  }
};
```

---

### 3.4 刷新Token

**接口说明**: 使用当前Token获取新的Token，延长登录时间

**请求**
```
POST /api/v1/shared/auth/refresh
```

**认证**: 🔒 需要 JWT Token（可以是即将过期的Token）

**请求头**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**
```json
{
  "code": 200,
  "message": "Token刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...(新Token)",
    "expiresIn": 86400
  }
}
```

**错误响应**
- `401` - Token无效或已失效
- `500` - 服务器错误

**Token刷新策略**:
- 建议在Token过期前5分钟刷新
- 如果Token已完全过期，需要重新登录
- 刷新后的Token有效期重新计算

**JavaScript/Axios 示例**
```javascript
const refreshToken = async () => {
  try {
    const response = await axios.post('/api/v1/shared/auth/refresh');
    const { token, expiresIn } = response.data.data;
    
    // 更新存储的Token
    localStorage.setItem('token', token);
    localStorage.setItem('tokenExpiry', Date.now() + expiresIn * 1000);
    
    return token;
  } catch (error) {
    if (error.response?.status === 401) {
      // Token完全失效，需要重新登录
      localStorage.clear();
      router.push('/login');
      throw new Error('登录已过期，请重新登录');
    }
    throw error;
  }
};

// 自动刷新Token（在过期前5分钟刷新）
const setupTokenRefresh = () => {
  const checkAndRefresh = async () => {
    const tokenExpiry = localStorage.getItem('tokenExpiry');
    if (!tokenExpiry) return;
    
    const expiryTime = parseInt(tokenExpiry);
    const now = Date.now();
    const timeUntilExpiry = expiryTime - now;
    
    // 如果Token将在5分钟内过期，刷新它
    if (timeUntilExpiry < 5 * 60 * 1000 && timeUntilExpiry > 0) {
      try {
        await refreshToken();
      } catch (error) {
        console.error('自动刷新Token失败:', error);
      }
    }
  };
  
  // 每分钟检查一次
  setInterval(checkAndRefresh, 60 * 1000);
  
  // 立即执行一次
  checkAndRefresh();
};
```

---

### 3.5 获取用户权限

**接口说明**: 获取当前登录用户的权限列表

**请求**
```
GET /api/v1/shared/auth/permissions
```

**认证**: 🔒 需要 JWT Token

**响应示例**
```json
{
  "code": 200,
  "message": "获取权限成功",
  "data": [
    "user.read",
    "user.write",
    "book.read",
    "comment.write"
  ]
}
```

**权限命名规范**:
- 格式：`resource.action`
- 资源：user, book, comment, admin 等
- 操作：read, write, delete, manage 等

**JavaScript/Axios 示例**
```javascript
const getUserPermissions = async () => {
  try {
    const response = await axios.get('/api/v1/shared/auth/permissions');
    return response.data.data;
  } catch (error) {
    console.error('获取权限失败:', error);
    return [];
  }
};

// 检查权限
const hasPermission = (permissions, required) => {
  return permissions.includes(required);
};

// 使用示例
const permissions = await getUserPermissions();
if (hasPermission(permissions, 'admin.manage')) {
  // 显示管理员功能
}
```

---

### 3.6 获取用户角色

**接口说明**: 获取当前登录用户的角色列表

**请求**
```
GET /api/v1/shared/auth/roles
```

**认证**: 🔒 需要 JWT Token

**响应示例**
```json
{
  "code": 200,
  "message": "获取角色成功",
  "data": ["user", "vip"]
}
```

**角色类型**:
| 角色 | 说明 | 权限范围 |
|------|------|----------|
| user | 普通用户 | 基础阅读权限 |
| vip | VIP用户 | VIP章节访问 |
| author | 作者 | 写作和发布权限 |
| admin | 管理员 | 平台管理权限 |
| super_admin | 超级管理员 | 完全管理权限 |

**JavaScript/Axios 示例**
```javascript
const getUserRoles = async () => {
  try {
    const response = await axios.get('/api/v1/shared/auth/roles');
    return response.data.data;
  } catch (error) {
    console.error('获取角色失败:', error);
    return [];
  }
};

// 检查角色
const hasRole = (roles, required) => {
  return roles.includes(required);
};

// 使用示例
const roles = await getUserRoles();
if (hasRole(roles, 'vip')) {
  // 显示VIP章节
}
```

---

## 4. Axios 拦截器配置

### 4.1 请求拦截器（自动添加Token）

```javascript
// axios-config.js
import axios from 'axios';
import { useRouter } from 'vue-router';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 10000
});

// 请求拦截器 - 自动添加Token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default apiClient;
```

### 4.2 响应拦截器（Token刷新和错误处理）

```javascript
// 响应拦截器 - 处理Token过期和错误
let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;
    
    // Token过期处理
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // 如果正在刷新，将请求加入队列
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers['Authorization'] = `Bearer ${token}`;
          return apiClient(originalRequest);
        }).catch(err => {
          return Promise.reject(err);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // 尝试刷新Token
        const response = await axios.post('/api/v1/shared/auth/refresh', {}, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        
        const { token, expiresIn } = response.data.data;
        
        // 保存新Token
        localStorage.setItem('token', token);
        localStorage.setItem('tokenExpiry', Date.now() + expiresIn * 1000);
        
        // 更新原请求的Token
        originalRequest.headers['Authorization'] = `Bearer ${token}`;
        
        // 处理队列中的请求
        processQueue(null, token);
        
        // 重试原请求
        return apiClient(originalRequest);
      } catch (refreshError) {
        // 刷新失败，清除Token并跳转登录页
        processQueue(refreshError, null);
        localStorage.clear();
        
        const router = useRouter();
        router.push('/login');
        
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }
    
    // 其他错误处理
    if (error.response?.status === 403) {
      console.error('权限不足');
      // 可以显示提示或跳转
    } else if (error.response?.status >= 500) {
      console.error('服务器错误');
      // 显示错误提示
    }
    
    return Promise.reject(error);
  }
);
```

---

## 5. 完整示例

### 5.1 登录页面组件（Vue 3）

```vue
<template>
  <div class="login-page">
    <el-card class="login-card">
      <h2>用户登录</h2>
      <el-form :model="loginForm" :rules="loginRules" ref="formRef">
        <el-form-item prop="email">
          <el-input
            v-model="loginForm.email"
            placeholder="邮箱或用户名"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <el-link @click="goToRegister">注册账号</el-link>
        <el-link @click="goToForgotPassword">忘记密码</el-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { ElMessage } from 'element-plus';
import axios from '@/utils/axios';

const router = useRouter();
const userStore = useUserStore();
const formRef = ref();
const loading = ref(false);

const loginForm = reactive({
  email: '',
  password: ''
});

const loginRules = {
  email: [
    { required: true, message: '请输入邮箱或用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
};

const handleLogin = async () => {
  if (!formRef.value) return;
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    
    loading.value = true;
    try {
      const response = await axios.post('/api/v1/shared/auth/login', {
        email: loginForm.email,
        password: loginForm.password
      });
      
      const { token, expiresIn, user } = response.data.data;
      
      // 保存Token和用户信息
      localStorage.setItem('token', token);
      localStorage.setItem('tokenExpiry', Date.now() + expiresIn * 1000);
      
      // 更新Pinia store
      userStore.setUser(user);
      userStore.setToken(token);
      
      ElMessage.success('登录成功');
      
      // 跳转到首页或之前访问的页面
      const redirect = router.currentRoute.value.query.redirect || '/';
      router.push(redirect);
    } catch (error) {
      if (error.response?.status === 401) {
        ElMessage.error('邮箱或密码错误');
      } else if (error.response?.status === 403) {
        ElMessage.error('账号已被禁用，请联系管理员');
      } else {
        ElMessage.error('登录失败，请稍后重试');
      }
    } finally {
      loading.value = false;
    }
  });
};

const goToRegister = () => {
  router.push('/register');
};

const goToForgotPassword = () => {
  router.push('/forgot-password');
};
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  padding: 40px;
}

.login-card h2 {
  text-align: center;
  margin-bottom: 30px;
}

.login-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}
</style>
```

### 5.2 用户Store（Pinia）

```javascript
// stores/user.js
import { defineStore } from 'pinia';
import axios from '@/utils/axios';

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: null,
    permissions: [],
    roles: []
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.roles.includes('admin'),
    isVip: (state) => state.roles.includes('vip'),
    hasPermission: (state) => (permission) => {
      return state.permissions.includes(permission);
    }
  },
  
  actions: {
    setUser(user) {
      this.user = user;
    },
    
    setToken(token) {
      this.token = token;
    },
    
    async loadUserInfo() {
      if (!this.token) return;
      
      try {
        // 加载权限和角色
        const [permissionsRes, rolesRes] = await Promise.all([
          axios.get('/api/v1/shared/auth/permissions'),
          axios.get('/api/v1/shared/auth/roles')
        ]);
        
        this.permissions = permissionsRes.data.data;
        this.roles = rolesRes.data.data;
      } catch (error) {
        console.error('加载用户信息失败:', error);
      }
    },
    
    async logout() {
      try {
        await axios.post('/api/v1/shared/auth/logout');
      } catch (error) {
        console.error('登出失败:', error);
      } finally {
        this.user = null;
        this.token = null;
        this.permissions = [];
        this.roles = [];
        localStorage.clear();
      }
    },
    
    async initFromStorage() {
      const token = localStorage.getItem('token');
      const userStr = localStorage.getItem('user');
      
      if (token && userStr) {
        this.token = token;
        this.user = JSON.parse(userStr);
        await this.loadUserInfo();
      }
    }
  }
});
```

### 5.3 路由守卫

```javascript
// router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores/user';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/views/Home.vue')
    },
    {
      path: '/login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/profile',
      component: () => import('@/views/Profile.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      component: () => import('@/views/Admin.vue'),
      meta: { requiresAuth: true, requiresRole: 'admin' }
    }
  ]
});

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore();
  
  // 初始化用户信息（从localStorage）
  if (!userStore.user && localStorage.getItem('token')) {
    await userStore.initFromStorage();
  }
  
  // 检查是否需要认证
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    });
    return;
  }
  
  // 检查角色权限
  if (to.meta.requiresRole) {
    const requiredRole = to.meta.requiresRole;
    if (!userStore.roles.includes(requiredRole)) {
      next({ path: '/403' }); // 权限不足页面
      return;
    }
  }
  
  // 检查特定权限
  if (to.meta.requiresPermission) {
    const requiredPermission = to.meta.requiresPermission;
    if (!userStore.hasPermission(requiredPermission)) {
      next({ path: '/403' });
      return;
    }
  }
  
  next();
});

export default router;
```

---

## 6. 最佳实践

### 6.1 Token安全存储

```javascript
// 使用 sessionStorage（更安全，关闭标签页后失效）
const useSecureStorage = () => {
  const storage = sessionStorage; // 或 localStorage
  
  return {
    setToken(token) {
      // 可以考虑加密存储
      storage.setItem('token', btoa(token)); // 简单Base64编码
    },
    
    getToken() {
      const encoded = storage.getItem('token');
      return encoded ? atob(encoded) : null;
    },
    
    removeToken() {
      storage.removeItem('token');
    }
  };
};
```

### 6.2 防止CSRF攻击

```javascript
// 添加CSRF Token
apiClient.interceptors.request.use((config) => {
  const csrfToken = document.querySelector('meta[name="csrf-token"]')?.content;
  if (csrfToken) {
    config.headers['X-CSRF-Token'] = csrfToken;
  }
  return config;
});
```

### 6.3 请求重试机制

```javascript
// 添加请求重试
const retryRequest = (error, retries = 3, delay = 1000) => {
  return new Promise((resolve, reject) => {
    if (retries === 0) {
      return reject(error);
    }
    
    setTimeout(() => {
      apiClient(error.config)
        .then(resolve)
        .catch((err) => {
          retryRequest(err, retries - 1, delay * 2)
            .then(resolve)
            .catch(reject);
        });
    }, delay);
  });
};
```

---

## 7. 常见问题

### Q1: Token存储在localStorage还是sessionStorage？

**A**: 
- **localStorage**: 持久化存储，关闭浏览器后仍然有效，适合"记住我"功能
- **sessionStorage**: 关闭标签页后失效，更安全，适合敏感应用

**推荐做法**:
```javascript
const storage = rememberMe ? localStorage : sessionStorage;
storage.setItem('token', token);
```

### Q2: 如何处理多标签页Token同步？

**A**: 使用 `storage` 事件监听：

```javascript
window.addEventListener('storage', (e) => {
  if (e.key === 'token') {
    if (!e.newValue) {
      // Token被删除，用户在其他标签页登出
      window.location.href = '/login';
    } else {
      // Token更新，刷新页面或更新store
      userStore.setToken(e.newValue);
    }
  }
});
```

### Q3: 如何实现"记住我"功能？

**A**: 
```javascript
const login = async (email, password, rememberMe) => {
  const response = await axios.post('/api/v1/shared/auth/login', {
    email,
    password
  });
  
  const { token, expiresIn } = response.data.data;
  
  // 根据rememberMe选择存储方式
  const storage = rememberMe ? localStorage : sessionStorage;
  storage.setItem('token', token);
  storage.setItem('tokenExpiry', Date.now() + expiresIn * 1000);
};
```

---

**文档版本**: v1.0  
**最后更新**: 2025-10-18  
**维护者**: 青羽后端团队

