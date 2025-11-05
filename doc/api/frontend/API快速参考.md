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

### 章节阅读

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 获取章节信息 | GET | `/reader/chapters/:id` | ❌ | - |
| 获取章节内容 | GET | `/reader/chapters/:id/content` | ✅ | - |
| 获取书籍章节列表 | GET | `/reader/chapters` | ❌ | `bookId, page, size` |
| 获取阅读设置 | GET | `/reader/settings` | ✅ | - |
| 保存阅读设置 | POST | `/reader/settings` | ✅ | Settings对象 |
| 更新阅读设置 | PUT | `/reader/settings` | ✅ | 部分字段 |

### 评论功能 ⭐️新增

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 发表评论 | POST | `/reader/comments` | ✅ | `{book_id, chapter_id?, content, rating?}` |
| 获取评论列表 | GET | `/reader/comments` | ❌ | `book_id, sortBy, page, size` |
| 获取评论详情 | GET | `/reader/comments/:id` | ❌ | - |
| 更新评论 | PUT | `/reader/comments/:id` | ✅ | `{content}` |
| 删除评论 | DELETE | `/reader/comments/:id` | ✅ | - |
| 回复评论 | POST | `/reader/comments/:id/reply` | ✅ | `{content}` |
| **点赞评论** ⭐️ | POST | `/reader/comments/:id/like` | ✅ | - |
| **取消点赞** ⭐️ | DELETE | `/reader/comments/:id/like` | ✅ | - |

### 阅读历史

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 记录阅读 | POST | `/reader/reading-history` | ✅ | `{book_id, chapter_id, start_time, end_time, progress}` |
| 获取历史列表 | GET | `/reader/reading-history` | ✅ | `page, page_size, book_id?` |
| 获取阅读统计 | GET | `/reader/reading-history/stats` | ✅ | `days` |
| 删除历史记录 | DELETE | `/reader/reading-history/:id` | ✅ | - |
| 清空历史记录 | DELETE | `/reader/reading-history` | ✅ | - |

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

## 写作系统

### 项目管理

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 创建项目 | POST | `/projects` | ✅ | `{title, description, category}` |
| 获取项目列表 | GET | `/projects` | ✅ | `page, pageSize, status, category` |
| 获取项目详情 | GET | `/projects/:id` | ✅ | - |
| 更新项目 | PUT | `/projects/:id` | ✅ | `{title?, description?, status?}` |
| 删除项目 | DELETE | `/projects/:id` | ✅ | - |
| 更新项目统计 | PUT | `/projects/:id/statistics` | ✅ | - |

### 文档管理

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 创建文档 | POST | `/projects/:projectId/documents` | ✅ | `{title, content, parentId}` |
| 获取文档列表 | GET | `/projects/:projectId/documents` | ✅ | `page, pageSize` |
| 获取文档树 | GET | `/projects/:projectId/documents/tree` | ✅ | - |
| 获取文档详情 | GET | `/documents/:id` | ✅ | - |
| 更新文档 | PUT | `/documents/:id` | ✅ | `{title?, content?, status?}` |
| 删除文档 | DELETE | `/documents/:id` | ✅ | - |
| 移动文档 | PUT | `/documents/:id/move` | ✅ | `{newParentId, newOrder}` |
| 重新排序 | PUT | `/projects/:projectId/documents/reorder` | ✅ | `{orders}` |

### 编辑器

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 自动保存 | POST | `/documents/:id/autosave` | ✅ | `{content, version}` |
| 获取保存状态 | GET | `/documents/:id/save-status` | ✅ | - |
| 获取文档内容 | GET | `/documents/:id/content` | ✅ | - |
| 更新文档内容 | PUT | `/documents/:id/content` | ✅ | `{content}` |
| 计算字数 | POST | `/documents/:id/word-count` | ✅ | `{content, filterMarkdown}` |
| 获取快捷键配置 | GET | `/user/shortcuts` | ✅ | - |
| 更新快捷键 | PUT | `/user/shortcuts` | ✅ | `{shortcuts}` |
| 重置快捷键 | POST | `/user/shortcuts/reset` | ✅ | - |

### 内容审核

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 实时检测 | POST | `/audit/check` | ✅ | `{content}` |
| 全文审核 | POST | `/documents/:id/audit` | ✅ | `{content}` |
| 获取审核结果 | GET | `/documents/:id/audit-result` | ✅ | `targetType` |
| 提交申诉 | POST | `/audit/:id/appeal` | ✅ | `{reason}` |
| 获取违规记录 | GET | `/users/:userId/violations` | ✅ | - |
| 获取违规统计 | GET | `/users/:userId/violation-summary` | ✅ | - |

### 数据统计

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 获取作品统计 | GET | `/writer/books/:book_id/stats` | ✅ | - |
| 获取章节统计 | GET | `/writer/chapters/:chapter_id/stats` | ✅ | - |
| 获取阅读热力图 | GET | `/writer/books/:book_id/heatmap` | ✅ | - |
| 获取收入统计 | GET | `/writer/books/:book_id/revenue` | ✅ | `start_date, end_date` |
| 获取热门章节 | GET | `/writer/books/:book_id/top-chapters` | ✅ | - |
| 获取每日统计 | GET | `/writer/books/:book_id/daily-stats` | ✅ | `days` |
| 获取跳出点分析 | GET | `/writer/books/:book_id/drop-off-points` | ✅ | - |
| 获取留存率 | GET | `/writer/books/:book_id/retention` | ✅ | `days` |

### 版本管理

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 获取版本历史 | GET | `/documents/:documentId/versions` | ✅ | `page, pageSize` |
| 获取特定版本 | GET | `/documents/:documentId/versions/:versionId` | ✅ | - |
| 比较版本 | GET | `/documents/:documentId/versions/compare` | ✅ | `fromVersion, toVersion` |
| 恢复版本 | POST | `/documents/:documentId/versions/:versionId/restore` | ✅ | - |

---

## 共享服务

### 钱包

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 查询余额 | GET | `/shared/wallet/balance` | ✅ | - |
| 获取钱包信息 | GET | `/shared/wallet` | ✅ | - |
| 充值 | POST | `/shared/wallet/recharge` | ✅ | `{amount, method}` |
| 消费 | POST | `/shared/wallet/consume` | ✅ | `{amount, reason}` |
| 转账 | POST | `/shared/wallet/transfer` | ✅ | `{to_user_id, amount, reason}` |
| 交易历史 | GET | `/shared/wallet/transactions` | ✅ | `page, page_size, type` |
| 申请提现 | POST | `/shared/wallet/withdraw` | ✅ | `{amount, account}` |
| 查询提现申请 | GET | `/shared/wallet/withdrawals` | ✅ | `page, page_size, status` |

### 存储

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 上传文件 | POST | `/shared/storage/upload` | ✅ | `file, path` |
| 下载文件 | GET | `/shared/storage/download/:fileId` | ✅ | - |
| 删除文件 | DELETE | `/shared/storage/files/:fileId` | ✅ | - |
| 获取文件信息 | GET | `/shared/storage/files/:fileId` | ✅ | - |
| 文件列表 | GET | `/shared/storage/files` | ✅ | `page, page_size, category` |
| 获取文件URL | GET | `/shared/storage/files/:fileId/url` | ✅ | `expire` |

### 管理员

| 功能 | 方法 | 路径 | 需认证 | 参数 |
|------|------|------|--------|------|
| 获取待审核内容 | GET | `/shared/admin/reviews/pending` | ✅🔑 | `content_type` |
| 审核内容 | POST | `/shared/admin/reviews` | ✅🔑 | `{content_id, action, reason}` |
| 审核提现 | POST | `/shared/admin/withdraw/review` | ✅🔑 | `{withdraw_id, approved, reason}` |
| 获取用户统计 | GET | `/shared/admin/users/:user_id/statistics` | ✅🔑 | - |
| 获取操作日志 | GET | `/shared/admin/operation-logs` | ✅🔑 | `page, page_size, admin_id, operation` |

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

## 统一响应格式 ⭐️已更新

### 成功响应

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 业务数据
  },
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde"
}
```

> **新增字段说明**:
> - `timestamp`: Unix时间戳，服务器响应时间
> - `request_id`: 请求追踪ID，便于调试和日志追踪（可选）

### 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "error": "username is required",
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde"
}
```

### 分页响应

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [...],
  "timestamp": 1729875123,
  "request_id": "req-12345-abcde",
  "pagination": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5,
    "has_next": true,
    "has_previous": false
  }
}
```

> **分页字段说明**:
> - `total`: 总记录数
> - `page`: 当前页码
> - `page_size`: 每页大小
> - `total_pages`: 总页数
> - `has_next`: 是否有下一页
> - `has_previous`: 是否有上一页

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

// 响应拦截 ⭐️已更新
request.interceptors.response.use(
  response => {
    const { code, message, data, timestamp, request_id } = response.data;
    
    // 记录请求ID便于追踪
    if (request_id && process.env.NODE_ENV === 'development') {
      console.debug('Request ID:', request_id);
    }
    
    // 返回数据部分
    return data;
  },
  error => {
    const { code, message, error: errorDetail, request_id } = error.response?.data || {};
    
    // 401 未授权
    if (code === 401 || error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    
    // 记录错误追踪ID
    if (request_id) {
      console.error('Error Request ID:', request_id);
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

### 创建项目

```javascript
// 创建写作项目
export async function createProject(title, description, category) {
  const response = await request.post('/projects', {
    title,
    description,
    category,
    type: 'novel'
  });
  return response;
}

// 获取项目列表
export function getProjects(page = 1, pageSize = 10) {
  return request.get('/projects', {
    params: { page, pageSize }
  });
}
```

### 编辑器自动保存

```javascript
// 自动保存（每30秒）
let autoSaveTimer = null;

export function enableAutoSave(documentId, getContent, getVersion) {
  autoSaveTimer = setInterval(async () => {
    try {
      await request.post(`/documents/${documentId}/autosave`, {
        content: getContent(),
        version: getVersion()
      });
      console.log('自动保存成功');
    } catch (error) {
      if (error.response?.status === 409) {
        alert('文档已被修改，请刷新页面');
      }
    }
  }, 30000); // 30秒
}

export function disableAutoSave() {
  if (autoSaveTimer) {
    clearInterval(autoSaveTimer);
  }
}
```

### 钱包操作

```javascript
// 查询余额
export async function getBalance() {
  const response = await request.get('/shared/wallet/balance');
  return response.balance;
}

// 充值
export async function recharge(amount, method = 'alipay') {
  return request.post('/shared/wallet/recharge', {
    amount,
    method
  });
}

// 消费
export async function consume(amount, reason) {
  return request.post('/shared/wallet/consume', {
    amount,
    reason
  });
}
```

### 评论点赞操作 ⭐️新增

```javascript
// 点赞评论
export async function likeComment(commentId) {
  return request.post(`/reader/comments/${commentId}/like`);
}

// 取消点赞
export async function unlikeComment(commentId) {
  return request.delete(`/reader/comments/${commentId}/like`);
}

// Vue组件中使用示例
const handleLike = async (comment) => {
  try {
    if (comment.is_liked) {
      await unlikeComment(comment.id);
      comment.is_liked = false;
      comment.like_count--;
      ElMessage.success('已取消点赞');
    } else {
      await likeComment(comment.id);
      comment.is_liked = true;
      comment.like_count++;
      ElMessage.success('点赞成功');
    }
  } catch (error) {
    ElMessage.error('操作失败');
  }
};
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
- [阅读器 API](./阅读器API参考.md) - 阅读器完整接口
- [推荐系统 API](./推荐系统API参考.md) - 推荐系统完整接口
- [写作系统 API](./写作系统API参考.md) - 写作系统完整接口（项目、文档、编辑器、审核、统计）
- [共享服务 API](./共享服务API参考.md) - 共享服务完整接口（钱包、存储、认证、管理员）

---

---

## AI 服务提供商 ⭐️新增

| 提供商 | 模型 | 说明 | 状态 |
|--------|------|------|------|
| OpenAI | GPT-3.5, GPT-4 | 通用文本生成 | ✅ |
| Claude | Claude-2, Claude-3 | Anthropic AI | ✅ |
| Gemini | Gemini-Pro | Google AI | ✅ |
| **DeepSeek** ⭐️ | deepseek-chat | 兼容OpenAI格式 | ✅ |

**默认提供商**: DeepSeek

**支持功能**:
- ✅ 智能续写
- ✅ 文本改写
- ✅ 文本扩写
- ✅ 文本润色
- ✅ Token使用统计

---

**最后更新**: 2025-10-25  
**维护者**: 青羽后端团队  
**文档版本**: v1.3

