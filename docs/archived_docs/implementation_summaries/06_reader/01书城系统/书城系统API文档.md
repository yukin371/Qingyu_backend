# 书城系统API文档

> **模块**: 书城系统  
> **版本**: v1.0  
> **最后更新**: 2025-10-08

## 📋 目录

1. [API概览](#1-api概览)
2. [首页API](#2-首页api)
3. [书籍API](#3-书籍api)
4. [分类API](#4-分类api)
5. [BannerAPI](#5-bannerapi)
6. [搜索API](#6-搜索api)
7. [错误码](#7-错误码)
8. [使用示例](#8-使用示例)

## 1. API概览

### 1.1 基础信息

- **Base URL**: `/api/v1/bookstore`
- **请求格式**: JSON
- **响应格式**: JSON
- **字符编码**: UTF-8
- **认证方式**: JWT Token（部分接口需要）

### 1.2 API列表

| 方法 | 路径 | 功能 | 认证 |
|-----|------|------|------|
| GET | `/homepage` | 获取首页数据 | 否 |
| GET | `/books` | 获取书籍列表 | 否 |
| GET | `/books/:id` | 获取书籍详情 | 否 |
| GET | `/books/recommended` | 获取推荐书籍 | 否 |
| GET | `/books/featured` | 获取精选书籍 | 否 |
| GET | `/books/hot` | 获取热门书籍 | 否 |
| GET | `/books/search` | 搜索书籍 | 否 |
| POST | `/books/:id/view` | 增加浏览量 | 可选 |
| GET | `/categories` | 获取分类列表 | 否 |
| GET | `/categories/:id` | 获取分类详情 | 否 |
| GET | `/categories/:id/books` | 获取分类下的书籍 | 否 |
| GET | `/categories/tree` | 获取分类树 | 否 |
| GET | `/banners` | 获取Banner列表 | 否 |
| POST | `/banners/:id/click` | 增加点击量 | 可选 |

### 1.3 统一响应格式

#### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 具体数据
  }
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "error": "详细错误信息"
}
```

## 2. 首页API

### 2.1 获取首页数据

获取书城首页所有数据，包括Banner、推荐书籍、精选书籍等。

**请求**

```
GET /api/v1/bookstore/homepage
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "banners": [
      {
        "id": "65f1234567890abcdef12345",
        "title": "新书推荐",
        "description": "本月最火新书",
        "image": "https://example.com/banner1.jpg",
        "target": "65f1234567890abcdef12346",
        "targetType": "book",
        "sortOrder": 1,
        "clickCount": 1250
      }
    ],
    "recommendedBooks": [
      {
        "id": "65f1234567890abcdef12346",
        "title": "修真世界",
        "author": "方想",
        "cover": "https://example.com/cover1.jpg",
        "introduction": "这是一个修真的世界...",
        "categoryIds": ["65f1234567890abcdef12347"],
        "tags": ["玄幻", "修真"],
        "rating": 4.8,
        "viewCount": 125000,
        "likeCount": 3500,
        "isRecommended": true,
        "isFeatured": false
      }
    ],
    "featuredBooks": [],
    "hotBooks": [],
    "categories": [
      {
        "id": "65f1234567890abcdef12347",
        "name": "玄幻",
        "icon": "https://example.com/icon1.png",
        "bookCount": 5420,
        "sortOrder": 1
      }
    ],
    "stats": {
      "totalBooks": 50000,
      "publishedBooks": 45000
    }
  }
}
```

**字段说明**

| 字段 | 类型 | 说明 |
|-----|------|------|
| banners | Array | Banner列表，最多5个 |
| recommendedBooks | Array | 推荐书籍，最多20本 |
| featuredBooks | Array | 精选书籍，最多10本 |
| hotBooks | Array | 热门书籍，最多10本 |
| categories | Array | 分类列表 |
| stats | Object | 统计数据 |

## 3. 书籍API

### 3.1 获取书籍列表

获取书籍列表，支持分页。

**请求**

```
GET /api/v1/bookstore/books?page=1&size=20&status=published
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| size | int | 否 | 每页数量 | 20 |
| status | string | 否 | 状态 | published |
| categoryId | string | 否 | 分类ID | - |
| author | string | 否 | 作者 | - |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "books": [
      {
        "id": "65f1234567890abcdef12346",
        "title": "修真世界",
        "author": "方想",
        "cover": "https://example.com/cover1.jpg",
        "rating": 4.8,
        "viewCount": 125000,
        "chapterCount": 2000
      }
    ],
    "total": 5000,
    "page": 1,
    "size": 20,
    "totalPages": 250
  }
}
```

### 3.2 获取书籍详情

获取单本书籍的详细信息。

**请求**

```
GET /api/v1/bookstore/books/:id
```

**路径参数**

| 参数 | 类型 | 说明 |
|-----|------|------|
| id | string | 书籍ID |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "65f1234567890abcdef12346",
    "title": "修真世界",
    "author": "方想",
    "authorId": "65f1234567890abcdef12348",
    "introduction": "这是一个修真的世界，方运带着刚刚研究出的计算机技术来到这里...",
    "cover": "https://example.com/cover1.jpg",
    "categoryIds": ["65f1234567890abcdef12347"],
    "tags": ["玄幻", "修真", "东方玄幻"],
    "status": "published",
    "wordCount": 5000000,
    "chapterCount": 2000,
    "viewCount": 125000,
    "likeCount": 3500,
    "commentCount": 1200,
    "rating": 4.8,
    "ratingCount": 5600,
    "isRecommended": true,
    "isFeatured": false,
    "publishedAt": "2024-01-01T00:00:00Z",
    "createdAt": "2023-12-01T00:00:00Z",
    "updatedAt": "2024-10-08T00:00:00Z"
  }
}
```

### 3.3 获取推荐书籍

获取系统推荐的书籍列表。

**请求**

```
GET /api/v1/bookstore/books/recommended?limit=20
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| limit | int | 否 | 数量限制 | 20 |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "books": [
      {
        "id": "65f1234567890abcdef12346",
        "title": "修真世界",
        "author": "方想",
        "cover": "https://example.com/cover1.jpg",
        "rating": 4.8,
        "viewCount": 125000
      }
    ],
    "total": 20
  }
}
```

### 3.4 获取精选书籍

获取编辑精选的书籍列表。

**请求**

```
GET /api/v1/bookstore/books/featured?limit=10
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| limit | int | 否 | 数量限制 | 10 |

**响应格式同推荐书籍**

### 3.5 获取热门书籍

获取热度最高的书籍列表。

**请求**

```
GET /api/v1/bookstore/books/hot?limit=10
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| limit | int | 否 | 数量限制 | 10 |

**响应格式同推荐书籍**

### 3.6 增加浏览量

记录用户浏览书籍行为，增加浏览量统计。

**请求**

```
POST /api/v1/bookstore/books/:id/view
```

**Headers**

```
Authorization: Bearer <token>  // 可选
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "viewCount": 125001
  }
}
```

**说明**:
- 支持匿名访问
- 带Token时记录用户ID
- 每个用户每本书每天只统计一次

## 4. 分类API

### 4.1 获取分类列表

获取所有分类列表。

**请求**

```
GET /api/v1/bookstore/categories?level=0
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| level | int | 否 | 分类层级 | - |
| parentId | string | 否 | 父分类ID | - |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "categories": [
      {
        "id": "65f1234567890abcdef12347",
        "name": "玄幻",
        "description": "玄幻小说分类",
        "icon": "https://example.com/icon1.png",
        "parentId": null,
        "level": 0,
        "sortOrder": 1,
        "bookCount": 5420,
        "isActive": true
      }
    ],
    "total": 12
  }
}
```

### 4.2 获取分类详情

获取单个分类的详细信息。

**请求**

```
GET /api/v1/bookstore/categories/:id
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "65f1234567890abcdef12347",
    "name": "玄幻",
    "description": "玄幻小说分类，包括东方玄幻、异世大陆等",
    "icon": "https://example.com/icon1.png",
    "parentId": null,
    "level": 0,
    "sortOrder": 1,
    "bookCount": 5420,
    "isActive": true,
    "children": [
      {
        "id": "65f1234567890abcdef12349",
        "name": "东方玄幻",
        "bookCount": 2100
      }
    ]
  }
}
```

### 4.3 获取分类下的书籍

获取指定分类下的书籍列表。

**请求**

```
GET /api/v1/bookstore/categories/:id/books?page=1&size=20
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| size | int | 否 | 每页数量 | 20 |
| sortBy | string | 否 | 排序字段 | viewCount |
| sortOrder | string | 否 | 排序方向 | desc |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "category": {
      "id": "65f1234567890abcdef12347",
      "name": "玄幻"
    },
    "books": [
      {
        "id": "65f1234567890abcdef12346",
        "title": "修真世界",
        "author": "方想",
        "cover": "https://example.com/cover1.jpg",
        "rating": 4.8
      }
    ],
    "total": 5420,
    "page": 1,
    "size": 20
  }
}
```

### 4.4 获取分类树

获取完整的分类层级树结构。

**请求**

```
GET /api/v1/bookstore/categories/tree
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "tree": [
      {
        "id": "65f1234567890abcdef12347",
        "name": "玄幻",
        "level": 0,
        "bookCount": 5420,
        "children": [
          {
            "id": "65f1234567890abcdef12349",
            "name": "东方玄幻",
            "level": 1,
            "bookCount": 2100,
            "children": []
          }
        ]
      }
    ]
  }
}
```

## 5. BannerAPI

### 5.1 获取Banner列表

获取活动中的Banner列表。

**请求**

```
GET /api/v1/bookstore/banners?limit=5
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 默认值 |
|-----|------|------|------|--------|
| limit | int | 否 | 数量限制 | 5 |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "banners": [
      {
        "id": "65f1234567890abcdef12345",
        "title": "新书推荐",
        "description": "本月最火新书",
        "image": "https://example.com/banner1.jpg",
        "target": "65f1234567890abcdef12346",
        "targetType": "book",
        "sortOrder": 1,
        "clickCount": 1250
      }
    ]
  }
}
```

### 5.2 增加点击量

记录Banner点击行为。

**请求**

```
POST /api/v1/bookstore/banners/:id/click
```

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "clickCount": 1251
  }
}
```

## 6. 搜索API

### 6.1 搜索书籍

搜索书籍，支持多条件筛选。

**请求**

```
POST /api/v1/bookstore/books/search
```

**请求体**

```json
{
  "keyword": "修真",
  "categoryIds": ["65f1234567890abcdef12347"],
  "author": "方想",
  "tags": ["玄幻"],
  "minRating": 4.0,
  "sortBy": "viewCount",
  "sortOrder": "desc",
  "page": 1,
  "size": 20
}
```

**参数说明**

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| keyword | string | 否 | 关键词（书名/作者） |
| categoryIds | Array | 否 | 分类ID列表 |
| author | string | 否 | 作者名 |
| tags | Array | 否 | 标签列表 |
| minRating | float | 否 | 最低评分 |
| sortBy | string | 否 | 排序字段 |
| sortOrder | string | 否 | 排序方向 |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**响应**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "books": [
      {
        "id": "65f1234567890abcdef12346",
        "title": "修真世界",
        "author": "方想",
        "cover": "https://example.com/cover1.jpg",
        "rating": 4.8,
        "viewCount": 125000,
        "matchScore": 0.95
      }
    ],
    "total": 156,
    "page": 1,
    "size": 20,
    "keyword": "修真"
  }
}
```

## 7. 错误码

### 7.1 HTTP状态码

| 状态码 | 说明 |
|-------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用 |

### 7.2 业务错误码

| 错误码 | 说明 |
|-------|------|
| 10001 | 参数验证失败 |
| 10002 | 书籍不存在 |
| 10003 | 分类不存在 |
| 10004 | Banner不存在 |
| 10005 | 搜索条件为空 |
| 10006 | 数据库操作失败 |
| 10007 | 缓存操作失败 |

## 8. 使用示例

### 8.1 JavaScript示例

```javascript
// 获取首页数据
async function getHomepage() {
  try {
    const response = await fetch('http://api.example.com/api/v1/bookstore/homepage');
    const data = await response.json();
    
    if (data.code === 200) {
      console.log('首页数据:', data.data);
      return data.data;
    } else {
      console.error('获取失败:', data.message);
    }
  } catch (error) {
    console.error('请求错误:', error);
  }
}

// 搜索书籍
async function searchBooks(keyword) {
  try {
    const response = await fetch('http://api.example.com/api/v1/bookstore/books/search', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        keyword: keyword,
        page: 1,
        size: 20
      })
    });
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('搜索错误:', error);
  }
}

// 增加浏览量
async function incrementViewCount(bookId, token) {
  try {
    const headers = {
      'Content-Type': 'application/json'
    };
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    
    const response = await fetch(`http://api.example.com/api/v1/bookstore/books/${bookId}/view`, {
      method: 'POST',
      headers: headers
    });
    
    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('请求错误:', error);
  }
}
```

### 8.2 Python示例

```python
import requests

BASE_URL = 'http://api.example.com/api/v1/bookstore'

# 获取首页数据
def get_homepage():
    response = requests.get(f'{BASE_URL}/homepage')
    if response.status_code == 200:
        data = response.json()
        if data['code'] == 200:
            return data['data']
    return None

# 搜索书籍
def search_books(keyword, page=1, size=20):
    payload = {
        'keyword': keyword,
        'page': page,
        'size': size
    }
    response = requests.post(f'{BASE_URL}/books/search', json=payload)
    if response.status_code == 200:
        return response.json()['data']
    return None

# 获取书籍详情
def get_book_detail(book_id):
    response = requests.get(f'{BASE_URL}/books/{book_id}')
    if response.status_code == 200:
        return response.json()['data']
    return None
```

### 8.3 Go示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const baseURL = "http://api.example.com/api/v1/bookstore"

// 获取首页数据
func GetHomepage() (*HomepageResponse, error) {
    resp, err := http.Get(baseURL + "/homepage")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Code    int                `json:"code"`
        Message string             `json:"message"`
        Data    *HomepageResponse  `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    if result.Code != 200 {
        return nil, fmt.Errorf(result.Message)
    }
    
    return result.Data, nil
}

// 搜索书籍
func SearchBooks(keyword string, page, size int) (*SearchBooksResponse, error) {
    payload := map[string]interface{}{
        "keyword": keyword,
        "page":    page,
        "size":    size,
    }
    
    data, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.Post(
        baseURL+"/books/search",
        "application/json",
        bytes.NewBuffer(data),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Code    int                   `json:"code"`
        Message string                `json:"message"`
        Data    *SearchBooksResponse  `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Data, nil
}
```

---

**文档维护**: 青羽后端团队
**最后更新**: 2025-10-08
