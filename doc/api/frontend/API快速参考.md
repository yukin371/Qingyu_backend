# API 快速参考

> **一页纸速查表** - 适合打印或快速查阅

## 基础信息

**Base URL**: `http://localhost:8080/api/v1`  
**认证方式**: `Authorization: Bearer <token>`  
**Content-Type**: `application/json`

---

## 认证相关

| 功能 | 方法 | 路径 | 需认证 | 请求体 |
|------|------|------|--------|--------|
| 用户注册 | POST | `/register` | ❌ | `{username, email, password}` |
| 用户登录 | POST | `/login` | ❌ | `{username, password}` |
| 登出 | POST | `/shared/auth/logout` | ✅ | - |
| 刷新Token | POST | `/shared/auth/refresh` | ✅ | - |
| 获取权限 | GET | `/shared/auth/permissions` | ✅ | - |
| 获取角色 | GET | `/shared/auth/roles` | ✅ | - |

---

## 用户相关

| 功能 | 方法 | 路径 | 需认证 | 说明 |
|------|------|------|--------|------|
| 获取个人信息 | GET | `/users/profile` | ✅ | - |
| 更新个人信息 | PUT | `/users/profile` | ✅ | `{nickname?, bio?, avatar?, phone?}` |
| 修改密码 | PUT | `/users/password` | ✅ | `{oldPassword, newPassword}` |
| 获取用户列表 | GET | `/admin/users` | ✅🔑 | 管理员权限 |
| 获取指定用户 | GET | `/admin/users/:id` | ✅🔑 | 管理员权限 |
| 更新用户 | PUT | `/admin/users/:id` | ✅🔑 | 管理员权限 |
| 删除用户 | DELETE | `/admin/users/:id` | ✅🔑 | 管理员权限 |

---

## 书城相关

| 功能 | 方法 | 路径 | 参数 |
|------|------|------|------|
| 获取首页数据 | GET | `/bookstore/homepage` | - |
| 获取书籍详情 | GET | `/bookstore/books/:id` | - |
| 根据分类获取书籍 | GET | `/bookstore/categories/:categoryId/books` | `page, size` |
| 获取推荐书籍 | GET | `/bookstore/books/recommended` | `page, size` |
| 获取精选书籍 | GET | `/bookstore/books/featured` | `page, size` |
| 搜索书籍 | GET | `/bookstore/books/search` | `keyword, author, categoryId, page, size` |
| 获取分类树 | GET | `/bookstore/categories/tree` | - |
| 获取分类详情 | GET | `/bookstore/categories/:id` | - |
| 获取Banner列表 | GET | `/bookstore/banners` | `limit` |
| 增加浏览量 | POST | `/bookstore/books/:id/view` | - |
| 增加Banner点击 | POST | `/bookstore/banners/:id/click` | - |

### 榜单

| 功能 | 方法 | 路径 | 参数 |
|------|------|------|------|
| 实时榜 | GET | `/bookstore/rankings/realtime` | `limit` |
| 周榜 | GET | `/bookstore/rankings/weekly` | `period, limit` |
| 月榜 | GET | `/bookstore/rankings/monthly` | `period, limit` |
| 新人榜 | GET | `/bookstore/rankings/newbie` | `period, limit` |
| 按类型获取榜单 | GET | `/bookstore/rankings/:type` | `period, limit` |

---

## 阅读器相关

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 获取章节信息 | GET | `/reader/chapters/:id` | ❌ | - |
| 获取章节内容 | GET | `/reader/chapters/:id/content` | ✅ | - |
| 获取书籍章节列表 | GET | `/reader/chapters` | ❌ | `bookId, page, size` |
| 获取阅读设置 | GET | `/reader/settings` | ✅ | - |
| 保存阅读设置 | POST | `/reader/settings` | ✅ | Settings对象 |
| 更新阅读设置 | PUT | `/reader/settings` | ✅ | 部分字段 |

---

## 推荐系统

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 获取个性化推荐 | GET | `/recommendation/personalized` | ✅ | `limit` |
| 获取相似物品 | GET | `/recommendation/similar` | ❌ | `itemId, limit` |
| 记录用户行为 | POST | `/recommendation/behavior` | ✅ | 行为数据 |
| 获取首页推荐 | GET | `/recommendation/homepage` | ❌ | `limit` |
| 获取热门推荐 | GET | `/recommendation/hot` | ❌ | `limit, days` |
| 获取分类推荐 | GET | `/recommendation/category` | ❌ | `category, limit` |

---

## 共享服务

### 钱包

| 功能 | 方法 | 路径 | 需认证 |
|------|------|------|--------|
| 查询余额 | GET | `/shared/wallet/balance` | ✅ |
| 获取钱包信息 | GET | `/shared/wallet` | ✅ |
| 充值 | POST | `/shared/wallet/recharge` | ✅ |
| 消费 | POST | `/shared/wallet/consume` | ✅ |
| 交易历史 | GET | `/shared/wallet/transactions` | ✅ |

### 存储

| 功能 | 方法 | 路径 | 需认证 |
|------|------|------|--------|
| 上传文件 | POST | `/shared/storage/upload` | ✅ |
| 下载文件 | GET | `/shared/storage/download/:fileId` | ✅ |
| 删除文件 | DELETE | `/shared/storage/:fileId` | ✅ |
| 文件列表 | GET | `/shared/storage/files` | ✅ |

---

## 常用响应码

| 状态码 | 说明 | 处理方式 |
|--------|------|----------|
| 200 | 成功 | 正常处理 |
| 201 | 创建成功 | 正常处理（资源创建） |
| 400 | 参数错误 | 检查请求参数 |
| 401 | 未认证/Token过期 | 跳转登录或刷新Token |
| 403 | 权限不足 | 提示用户权限不够 |
| 404 | 资源不存在 | 提示资源未找到 |
| 500 | 服务器错误 | 提示稍后重试 |

---

## 统一响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 业务数据
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "error": "username is required"
}
```

### 分页响应

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [...],
  "total": 100,
  "page": 1,
  "pageSize": 20
}
```

---

## 快速上手代码片段

### Axios 配置

```javascript
import axios from 'axios';

const request = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 30000
});

// 请求拦截
request.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截
request.interceptors.response.use(
  response => response.data.data,
  error => {
    if (error.response?.status === 401) {
      // 跳转登录
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default request;
```

### 登录示例

```javascript
import request from './request';

// 登录
export async function login(username, password) {
  const response = await request.post('/login', {
    username,
    password
  });
  
  // 保存 Token
  localStorage.setItem('token', response.token);
  return response;
}

// 获取用户信息
export function getUserProfile() {
  return request.get('/users/profile');
}
```

### 分页加载

```javascript
// 获取书籍列表
export function getBookList(page = 1, size = 20, categoryId = '') {
  return request.get('/bookstore/books/search', {
    params: {
      page,
      size,
      categoryId
    }
  });
}

// Vue 组件中使用
const { data: books, total } = await getBookList(1, 20);
```

### 文件上传

```javascript
export function uploadFile(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  return request.post('/shared/storage/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}
```

### 错误处理

```javascript
try {
  const data = await getUserProfile();
  console.log('用户信息:', data);
} catch (error) {
  if (error.response) {
    // 服务器返回错误
    console.error('错误:', error.response.data.message);
  } else if (error.request) {
    // 网络错误
    console.error('网络错误');
  } else {
    // 其他错误
    console.error('未知错误');
  }
}
```

---

## 常见参数说明

### 分页参数

- `page`: 页码（从1开始）
- `pageSize` / `size`: 每页数量（默认20）
- `limit`: 限制数量（不分页）

### 排序参数

- `sortBy`: 排序字段（如 `created_at`, `view_count`）
- `sortOrder`: 排序方向（`asc` 升序, `desc` 降序）

### 过滤参数

- `keyword`: 搜索关键词
- `categoryId`: 分类ID
- `author`: 作者
- `status`: 状态
- `minRating`: 最低评分

---

## 测试环境

**开发环境**: `http://localhost:8080`  
**测试环境**: `http://test-api.qingyu.com`  
**生产环境**: `https://api.qingyu.com`

---

## 相关文档

- [前端集成指南](./前端集成指南.md) - 详细集成步骤
- [用户系统 API](./用户系统API参考.md) - 完整接口文档
- [书城系统 API](./书城API参考.md) - 完整接口文档

---

**最后更新**: 2025-10-18  
**维护者**: 青羽后端团队

