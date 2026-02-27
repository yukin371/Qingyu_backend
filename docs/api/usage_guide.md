# API 使用指南

## 概述

本文档提供 Qingyu 后端 API 的使用指南，包括认证方式、请求/响应格式、错误处理、分页规范等内容。

## 目录

- [认证方式](#认证方式)
- [通用请求格式](#通用请求格式)
- [通用响应格式](#通用响应格式)
- [错误码说明](#错误码说明)
- [分页参数说明](#分页参数说明)
- [请求示例](#请求示例)

## 认证方式

### Bearer Token 认证

大部分 API 需要使用 Bearer Token 进行认证。Token 通过用户登录或刷新获取。

#### 请求头格式

```
Authorization: Bearer <token>
```

#### 完整请求示例

```bash
curl -X GET https://api.example.com/api/v1/user/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Token 获取

#### 用户登录

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}
```

#### 响应

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600
  }
}
```

#### Token 刷新

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 通用请求格式

### HTTP 方法

| 方法 | 说明 | 示例 |
|------|------|------|
| GET | 获取资源 | 获取用户信息 |
| POST | 创建资源 | 创建新文章 |
| PUT | 完整更新资源 | 更新用户资料 |
| PATCH | 部分更新资源 | 更新部分字段 |
| DELETE | 删除资源 | 删除文章 |

### Content-Type

请求体通常使用 JSON 格式：

```
Content-Type: application/json
```

### 请求ID

每个请求会自动分配唯一 ID，用于日志追踪：

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

## 通用响应格式

### 成功响应

```json
{
  "success": true,
  "data": {
    // 响应数据
  },
  "message": "操作成功"
}
```

### 错误响应

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "参数验证失败",
    "details": {
      "email": "邮箱格式不正确"
    }
  },
  "requestId": "550e8400-e29b-41d4-a716-446655440000"
}
```

### HTTP 状态码

| 状态码 | 说明 | 使用场景 |
|--------|------|----------|
| 200 | OK | 请求成功 |
| 201 | Created | 资源创建成功 |
| 204 | No Content | 删除成功 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未认证或认证失败 |
| 403 | Forbidden | 无权限访问 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 资源冲突 |
| 422 | Unprocessable Entity | 参数验证失败 |
| 429 | Too Many Requests | 请求过于频繁 |
| 500 | Internal Server Error | 服务器内部错误 |
| 503 | Service Unavailable | 服务暂时不可用 |

## 错误码说明

### 通用错误码

| 错误码 | HTTP状态 | 说明 |
|--------|----------|------|
| SUCCESS | 200 | 操作成功 |
| VALIDATION_ERROR | 400 | 参数验证失败 |
| UNAUTHORIZED | 401 | 未认证 |
| FORBIDDEN | 403 | 无权限 |
| NOT_FOUND | 404 | 资源不存在 |
| CONFLICT | 409 | 资源冲突 |
| INTERNAL_ERROR | 500 | 内部错误 |
| SERVICE_UNAVAILABLE | 503 | 服务不可用 |

### 认证错误码

| 错误码 | 说明 |
|--------|------|
| INVALID_CREDENTIALS | 用户名或密码错误 |
| TOKEN_EXPIRED | Token 已过期 |
| TOKEN_INVALID | Token 无效 |
| REFRESH_TOKEN_INVALID | 刷新 Token 无效 |
| ACCOUNT_LOCKED | 账户已锁定 |
| ACCOUNT_DISABLED | 账户已禁用 |

### 业务错误码

| 错误码 | 说明 |
|--------|------|
| USER_NOT_FOUND | 用户不存在 |
| DUPLICATE_USERNAME | 用户名已存在 |
| DUPLICATE_EMAIL | 邮箱已存在 |
| BOOK_NOT_FOUND | 书籍不存在 |
| CHAPTER_NOT_FOUND | 章节不存在 |
| INSUFFICIENT_BALANCE | 余额不足 |
| QUOTA_EXCEEDED | 配额超限 |

## 分页参数说明

### 请求参数

```
GET /api/v1/books?limit=20&offset=0
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| limit | int | 20 | 每页数量（最大100） |
| offset | int | 0 | 偏移量 |
| page | int | 1 | 页码（与offset二选一） |
| sortBy | string | created_at | 排序字段 |
| sortOrder | string | desc | 排序方向（asc/desc） |

### 分页响应

```json
{
  "success": true,
  "data": {
    "items": [
      // 数据列表
    ],
    "pagination": {
      "total": 100,
      "limit": 20,
      "offset": 0,
      "hasMore": true
    }
  }
}
```

### 游标分页

对于大数据集，推荐使用游标分页：

```
GET /api/v1/books?cursor=xxx&limit=20
```

```json
{
  "success": true,
  "data": {
    "items": [],
    "pagination": {
      "nextCursor": "eyJpZCI6IjYw...",
      "hasMore": true
    }
  }
}
```

## 请求示例

### cURL 示例

#### GET 请求

```bash
curl -X GET "https://api.example.com/api/v1/books?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### POST 请求

```bash
curl -X POST "https://api.example.com/api/v1/books" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "我的新书",
    "author": "作者名",
    "category": "玄幻"
  }'
```

#### PUT 请求

```bash
curl -X PUT "https://api.example.com/api/v1/books/BOOK_ID" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "更新后的标题"
  }'
```

#### DELETE 请求

```bash
curl -X DELETE "https://api.example.com/api/v1/books/BOOK_ID" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Python 示例

```python
import requests

# 配置
BASE_URL = "https://api.example.com"
TOKEN = "YOUR_TOKEN"

headers = {
    "Authorization": f"Bearer {TOKEN}",
    "Content-Type": "application/json"
}

# GET 请求
response = requests.get(
    f"{BASE_URL}/api/v1/books",
    headers=headers,
    params={"limit": 10, "offset": 0}
)
data = response.json()

# POST 请求
response = requests.post(
    f"{BASE_URL}/api/v1/books",
    headers=headers,
    json={
        "title": "我的新书",
        "author": "作者名"
    }
)
```

### JavaScript 示例

```javascript
// 配置
const BASE_URL = 'https://api.example.com';
const TOKEN = 'YOUR_TOKEN';

// GET 请求
fetch(`${BASE_URL}/api/v1/books?limit=10&offset=0`, {
  headers: {
    'Authorization': `Bearer ${TOKEN}`
  }
})
.then(response => response.json())
.then(data => console.log(data));

// POST 请求
fetch(`${BASE_URL}/api/v1/books`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${TOKEN}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    title: '我的新书',
    author: '作者名'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

### Go 示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func main() {
    baseURL := "https://api.example.com"
    token := "YOUR_TOKEN"

    // 创建 HTTP 客户端
    client := &http.Client{}

    // GET 请求
    req, _ := http.NewRequest("GET", baseURL+"/api/v1/books?limit=10", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    // POST 请求
    data := map[string]string{
        "title": "我的新书",
        "author": "作者名",
    }
    jsonData, _ := json.Marshal(data)
    req, _ = http.NewRequest("POST", baseURL+"/api/v1/books", bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    resp, _ = client.Do(req)
    defer resp.Body.Close()
}
```

## 速率限制

### 默认限制

- 每个用户每分钟最多 100 次请求
- 匿名用户每分钟最多 20 次请求

### 响应头

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

### 超限响应

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "请求过于频繁，请稍后再试"
  }
}
```

## 版本管理

### API 版本

API 通过 URL 路径进行版本控制：

```
/api/v1/...  # 当前版本
/api/v2/...  # 未来版本
```

### 弃用通知

弃用的 API 会在响应头中标注：

```
Deprecation: true
Sunset: Sat, 01 Jan 2025 00:00:00 GMT
Link: </api/v2/resource>; rel="successor-version"
```

## 最佳实践

### 1. 错误处理

始终检查响应状态并处理错误：

```python
response = requests.get(url)
if response.status_code == 200:
    data = response.json()
    if data.get('success'):
        # 处理成功响应
        pass
    else:
        # 处理业务错误
        error = data.get('error')
        print(f"Error: {error['message']}")
else:
    # 处理 HTTP 错误
    print(f"HTTP Error: {response.status_code}")
```

### 2. Token 管理

- 在 Token 过期前自动刷新
- 安全存储 Token
- 处理 Token 失效情况

### 3. 分页处理

```python
def fetch_all_books():
    offset = 0
    limit = 100
    all_books = []

    while True:
        response = requests.get(
            f"{BASE_URL}/api/v1/books",
            params={"limit": limit, "offset": offset}
        )
        data = response.json()
        books = data['data']['items']
        all_books.extend(books)

        if not data['data']['pagination']['hasMore']:
            break

        offset += limit

    return all_books
```

### 4. 请求重试

对于临时性错误（5xx），实现指数退避重试：

```python
import time

def make_request_with_retry(url, max_retries=3):
    for attempt in range(max_retries):
        try:
            response = requests.get(url)
            if response.status_code < 500:
                return response
        except RequestException:
            pass

        if attempt < max_retries - 1:
            time.sleep(2 ** attempt)

    return None
```

## 相关文档

- [API 参考文档](./reference.md)
- [Swagger 文档](https://api.example.com/swagger)
- [错误处理指南](./error-handling.md)
- [模块 API 文档](./README.md)

---

**版本**: v1.0
**更新日期**: 2026-02-27
**维护者**: Backend API Team
