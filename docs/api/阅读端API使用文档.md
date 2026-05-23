# 阅读端API使用文档

> **版本**: v1.0  
> **最后更新**: 2025-10-16  
> **文档类型**: API使用手册

## Page Role

- current-reference
- current-owner: `docs/api/reader/`
- current-bounded: 当前阅读端 API 使用手册，负责阅读与书城相关接口的整体验证和调用参考

## Recommended Read Path

1. 先读 `reader/README.md` 和 `bookstore/README.md`。
2. 再读本页。

## Boundary

- 本页是使用手册，不是总入口。
- 当前阅读端 owner 仍在 `docs/api/reader/` 与 `docs/api/bookstore/`。

## Quick Section Map

- 文档概述
- 目录
- 基础信息
- 书城系统 API
- 阅读器系统 API
- 认证和权限
- 错误码说明
- 最佳实践

## Quick Takeaways

- 这是当前阅读端整体验证手册。

## Skip Guide

- 只看阅读器专题：去 `reader/README.md`。

---

## 📚 文档概述

本文档详细描述了青羽写作平台阅读端服务的所有API接口，包括书城系统和阅读器系统的完整功能。

### 系统划分

#### 1. 书城系统 (`/api/v1/bookstore`, `/api/v1/books`, `/api/v1/chapters`, `/api/v1/reading`)
- 书籍浏览和搜索
- 分类导航
- 书籍详情
- 章节管理
- 评分统计

#### 2. 阅读器系统 (`/api/v1/reader`)
- 章节阅读
- 阅读进度
- 注记和书签
- 阅读设置

---

## 🔗 目录

### 书城系统
- [首页数据](#1-首页数据)
- [书籍管理](#2-书籍管理)
- [分类管理](#3-分类管理)
- [章节管理](#4-章节管理)
- [评分系统](#5-评分系统)
- [榜单系统](#6-榜单系统)
- [搜索功能](#7-搜索功能)

### 阅读器系统
- [章节阅读](#8-章节阅读)
- [阅读进度](#9-阅读进度)
- [注记功能](#10-注记功能)
- [阅读设置](#11-阅读设置)

---

## 🌐 基础信息

### 服务地址
- **开发环境**: `http://localhost:8080`
- **生产环境**: `https://api.qingyu.com`

### 认证方式
- 大部分公开接口无需认证
- 用户相关接口需要JWT Token认证
- 在Header中添加: `Authorization: Bearer {token}`

### 统一响应格式

#### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 返回数据
  }
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

#### 分页响应
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    // 数据列表
  ],
  "total": 100,
  "page": 1,
  "size": 20
}
```

---

## 📖 书城系统 API

### 1. 首页数据

#### 1.1 获取首页数据

获取书城首页的所有数据，包括Banner、推荐书籍、精选书籍等。

**接口**: `GET /api/v1/bookstore/homepage`  
**认证**: 不需要

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/homepage"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取首页数据成功",
  "data": {
    "banners": [
      {
        "id": "banner_id",
        "title": "热门推荐",
        "image": "https://example.com/banner.jpg",
        "bookId": "book_id",
        "order": 1
      }
    ],
    "recommendedBooks": [
      {
        "id": "book_id",
        "title": "书名",
        "author": "作者",
        "cover": "封面URL",
        "rating": 4.5
      }
    ],
    "featuredBooks": [],
    "categories": []
  }
}
```

---

### 2. 书籍管理

#### 2.1 获取书籍详情

**接口**: `GET /api/v1/bookstore/books/{id}`  
**认证**: 不需要

**路径参数**:
- `id` (string, 必填) - 书籍ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/books/67890"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "67890",
    "title": "斗破苍穹",
    "author": "天蚕土豆",
    "cover": "https://example.com/cover.jpg",
    "description": "书籍简介...",
    "category": "玄幻",
    "status": "completed",
    "wordCount": 1000000,
    "chapterCount": 1000,
    "rating": 4.5,
    "views": 100000,
    "favorites": 5000,
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-12-31T23:59:59Z"
  }
}
```

#### 2.2 获取书籍列表（分页）

**接口**: `GET /api/v1/bookstore/books`  
**认证**: 不需要

**查询参数**:
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量
- `category` (string, 可选) - 分类筛选
- `status` (string, 可选) - 状态筛选（ongoing/completed）
- `sort` (string, 可选) - 排序方式（latest/popular/rating）

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/books?page=1&size=20&category=玄幻&sort=popular"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "book_id",
      "title": "书名",
      "author": "作者",
      "cover": "封面URL",
      "rating": 4.5
    }
  ],
  "total": 100,
  "page": 1,
  "size": 20
}
```

#### 2.3 根据分类获取书籍

**接口**: `GET /api/v1/bookstore/books/category/{category}`  
**认证**: 不需要

**路径参数**:
- `category` (string, 必填) - 分类名称

**查询参数**:
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/books/category/玄幻?page=1&size=20"
```

#### 2.4 根据作者获取书籍

**接口**: `GET /api/v1/bookstore/books/author/{author}`  
**认证**: 不需要

**路径参数**:
- `author` (string, 必填) - 作者名称

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/books/author/天蚕土豆"
```

---

### 3. 分类管理

#### 3.1 获取所有分类

**接口**: `GET /api/v1/bookstore/categories`  
**认证**: 不需要

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/categories"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "category_id",
      "name": "玄幻",
      "description": "玄幻小说",
      "icon": "icon_url",
      "parentId": "",
      "level": 1,
      "order": 1,
      "bookCount": 1000
    }
  ]
}
```

#### 3.2 获取子分类

**接口**: `GET /api/v1/bookstore/categories/{id}/children`  
**认证**: 不需要

**路径参数**:
- `id` (string, 必填) - 父分类ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/categories/parent_id/children"
```

---

### 4. 章节管理

#### 4.1 获取章节详情

**接口**: `GET /api/v1/chapters/{id}`  
**认证**: 不需要

**路径参数**:
- `id` (string, 必填) - 章节ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/chapters/chapter123"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "chapter123",
    "bookId": "book_id",
    "title": "第一章 初入修仙界",
    "chapterNumber": 1,
    "content": "章节内容...",
    "wordCount": 3000,
    "vipRequired": false,
    "price": 0,
    "publishTime": "2023-01-01T00:00:00Z",
    "createdAt": "2023-01-01T00:00:00Z"
  }
}
```

#### 4.2 获取书籍章节列表

**接口**: `GET /api/v1/books/{book_id}/chapters`  
**认证**: 不需要

**路径参数**:
- `book_id` (string, 必填) - 书籍ID

**查询参数**:
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/books/book123/chapters?page=1&size=50"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "chapter_id",
      "title": "第一章 标题",
      "chapterNumber": 1,
      "wordCount": 3000,
      "vipRequired": false,
      "publishTime": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 1000,
  "page": 1,
  "size": 50
}
```

#### 4.3 获取最新章节

**接口**: `GET /api/v1/books/{book_id}/chapters/latest`  
**认证**: 不需要

**路径参数**:
- `book_id` (string, 必填) - 书籍ID

**查询参数**:
- `limit` (int, 可选, 默认10) - 返回数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/books/book123/chapters/latest?limit=5"
```

---

### 5. 评分系统

#### 5.1 获取书籍评分列表

**接口**: `GET /api/v1/reading/books/{book_id}/ratings`  
**认证**: 不需要

**路径参数**:
- `book_id` (string, 必填) - 书籍ID

**查询参数**:
- `page` (int, 可选, 默认1) - 页码
- `limit` (int, 可选, 默认10) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reading/books/book123/ratings?page=1&limit=20"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "rating_id",
      "userId": "user123",
      "userName": "用户昵称",
      "bookId": "book123",
      "rating": 5,
      "review": "非常精彩的小说",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 500,
  "page": 1,
  "limit": 20
}
```

#### 5.2 创建书籍评分

**接口**: `POST /api/v1/reading/ratings`  
**认证**: 需要（JWT Token）

**请求体**:
```json
{
  "bookId": "book123",
  "rating": 5,
  "review": "非常精彩的小说，强烈推荐！"
}
```

**请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/reading/ratings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bookId": "book123",
    "rating": 5,
    "review": "非常精彩的小说"
  }'
```

**响应示例**:
```json
{
  "code": 200,
  "message": "评分成功",
  "data": {
    "id": "rating_id",
    "userId": "user123",
    "bookId": "book123",
    "rating": 5,
    "review": "非常精彩的小说",
    "createdAt": "2023-01-01T00:00:00Z"
  }
}
```

#### 5.3 更新评分

**接口**: `PUT /api/v1/reading/ratings/{id}`  
**认证**: 需要（JWT Token）

**路径参数**:
- `id` (string, 必填) - 评分ID

**请求体**:
```json
{
  "rating": 4,
  "review": "更新后的评价"
}
```

**请求示例**:
```bash
curl -X PUT "http://localhost:8080/api/v1/reading/ratings/rating123" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 4,
    "review": "更新后的评价"
  }'
```

#### 5.4 删除评分

**接口**: `DELETE /api/v1/reading/ratings/{id}`  
**认证**: 需要（JWT Token）

**路径参数**:
- `id` (string, 必填) - 评分ID

**请求示例**:
```bash
curl -X DELETE "http://localhost:8080/api/v1/reading/ratings/rating123" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 5.5 获取书籍评分统计

**接口**: `GET /api/v1/reading/books/{book_id}/ratings/stats`  
**认证**: 不需要

**路径参数**:
- `book_id` (string, 必填) - 书籍ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reading/books/book123/ratings/stats"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "bookId": "book123",
    "averageRating": 4.5,
    "totalRatings": 500,
    "ratingDistribution": {
      "5": 300,
      "4": 150,
      "3": 30,
      "2": 15,
      "1": 5
    }
  }
}
```

---

### 6. 榜单系统

#### 6.1 获取热门榜单

**接口**: `GET /api/v1/bookstore/rankings/hot`  
**认证**: 不需要

**查询参数**:
- `limit` (int, 可选, 默认20) - 返回数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/bookstore/rankings/hot?limit=50"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "rank": 1,
      "bookId": "book123",
      "title": "书名",
      "author": "作者",
      "cover": "封面URL",
      "score": 95.8,
      "views": 1000000,
      "favorites": 50000
    }
  ]
}
```

#### 6.2 获取新书榜

**接口**: `GET /api/v1/bookstore/rankings/new`  
**认证**: 不需要

#### 6.3 获取完结榜

**接口**: `GET /api/v1/bookstore/rankings/completed`  
**认证**: 不需要

#### 6.4 获取收藏榜

**接口**: `GET /api/v1/bookstore/rankings/favorites`  
**认证**: 不需要

---

### 7. 搜索功能

#### 7.1 书籍搜索

**接口**: `GET /api/v1/books/search`  
**认证**: 不需要

**查询参数**:
- `q` (string, 必填) - 搜索关键词
- `type` (string, 可选) - 搜索类型（title/author/all）
- `category` (string, 可选) - 分类筛选
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/books/search?q=斗破&type=title&page=1&size=20"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "搜索成功",
  "data": [
    {
      "id": "book_id",
      "title": "斗破苍穹",
      "author": "天蚕土豆",
      "cover": "封面URL",
      "rating": 4.5,
      "highlight": "斗破苍穹"
    }
  ],
  "total": 100,
  "page": 1,
  "size": 20
}
```

#### 7.2 根据标题搜索

**接口**: `GET /api/v1/books/search/title`  
**认证**: 不需要

**查询参数**:
- `title` (string, 必填) - 书籍标题
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/books/search/title?title=斗破"
```

#### 7.3 根据标签搜索

**接口**: `GET /api/v1/books/search/tags`  
**认证**: 不需要

**查询参数**:
- `tags` (string, 必填) - 标签列表（逗号分隔）
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/books/search/tags?tags=玄幻,热血"
```

---

## 📱 阅读器系统 API

### 8. 章节阅读

#### 8.1 获取章节信息

**接口**: `GET /api/v1/reader/chapters/{id}`  
**认证**: 不需要

**路径参数**:
- `id` (string, 必填) - 章节ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters/chapter123"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "chapter123",
    "bookId": "book123",
    "title": "第一章 初入修仙界",
    "chapterNumber": 1,
    "wordCount": 3000,
    "vipRequired": false
  }
}
```

#### 8.2 获取章节内容

**接口**: `GET /api/v1/reader/chapters/{id}/content`  
**认证**: 需要（JWT Token）

**路径参数**:
- `id` (string, 必填) - 章节ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters/chapter123/content" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "content": "章节正文内容..."
  }
}
```

**注意事项**:
- 需要登录才能获取章节内容
- VIP章节需要相应的VIP权限
- 内容会进行缓存以提高性能

#### 8.3 获取书籍章节列表

**接口**: `GET /api/v1/reader/chapters`  
**认证**: 不需要

**查询参数**:
- `bookId` (string, 必填) - 书籍ID
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters?bookId=book123&page=1&size=50"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "chapters": [
      {
        "id": "chapter1",
        "title": "第一章",
        "chapterNumber": 1,
        "vipRequired": false
      }
    ],
    "total": 1000,
    "page": 1,
    "size": 50
  }
}
```

#### 8.4 获取章节导航

**接口**: `GET /api/v1/reader/chapters/navigation`  
**认证**: 不需要

**查询参数**:
- `bookId` (string, 必填) - 书籍ID
- `chapterNum` (int, 必填) - 当前章节号

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters/navigation?bookId=book123&chapterNum=10"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "prevChapter": {
      "id": "chapter9",
      "title": "第九章",
      "chapterNumber": 9
    },
    "nextChapter": {
      "id": "chapter11",
      "title": "第十一章",
      "chapterNumber": 11
    }
  }
}
```

#### 8.5 获取第一章

**接口**: `GET /api/v1/reader/chapters/first`  
**认证**: 不需要

**查询参数**:
- `bookId` (string, 必填) - 书籍ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters/first?bookId=book123"
```

#### 8.6 获取最后一章

**接口**: `GET /api/v1/reader/chapters/last`  
**认证**: 不需要

**查询参数**:
- `bookId` (string, 必填) - 书籍ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/chapters/last?bookId=book123"
```

---

### 9. 阅读进度

#### 9.1 获取阅读进度

**接口**: `GET /api/v1/reader/progress/{bookId}`  
**认证**: 需要（JWT Token）

**路径参数**:
- `bookId` (string, 必填) - 书籍ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/progress/book123" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "userId": "user123",
    "bookId": "book123",
    "chapterId": "chapter50",
    "progress": 0.65,
    "lastReadAt": "2023-12-31T23:59:59Z",
    "totalReadTime": 7200
  }
}
```

#### 9.2 保存阅读进度

**接口**: `POST /api/v1/reader/progress`  
**认证**: 需要（JWT Token）

**请求体**:
```json
{
  "bookId": "book123",
  "chapterId": "chapter50",
  "progress": 0.65
}
```

**请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/reader/progress" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bookId": "book123",
    "chapterId": "chapter50",
    "progress": 0.65
  }'
```

**响应示例**:
```json
{
  "code": 200,
  "message": "保存成功",
  "data": null
}
```

**参数说明**:
- `progress`: 章节内阅读进度，取值范围 0.0-1.0

#### 9.3 更新阅读时长

**接口**: `PUT /api/v1/reader/progress/time`  
**认证**: 需要（JWT Token）

**请求体**:
```json
{
  "bookId": "book123",
  "duration": 1800
}
```

**请求示例**:
```bash
curl -X PUT "http://localhost:8080/api/v1/reader/progress/time" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bookId": "book123",
    "duration": 1800
  }'
```

**参数说明**:
- `duration`: 本次阅读时长（秒）

#### 9.4 获取阅读历史

**接口**: `GET /api/v1/reader/progress/history`  
**认证**: 需要（JWT Token）

**查询参数**:
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/progress/history?page=1&size=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "history": [
      {
        "bookId": "book123",
        "bookTitle": "斗破苍穹",
        "bookCover": "封面URL",
        "chapterId": "chapter50",
        "chapterTitle": "第五十章",
        "progress": 0.65,
        "lastReadAt": "2023-12-31T23:59:59Z"
      }
    ],
    "total": 50,
    "page": 1,
    "size": 20
  }
}
```

#### 9.5 获取总阅读时长

**接口**: `GET /api/v1/reader/progress/total-time`  
**认证**: 需要（JWT Token）

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/progress/total-time" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "totalTime": 360000,
    "todayTime": 3600,
    "weekTime": 25200
  }
}
```

**说明**:
- `totalTime`: 总阅读时长（秒）
- `todayTime`: 今日阅读时长（秒）
- `weekTime`: 本周阅读时长（秒）

---

### 10. 注记功能

#### 10.1 创建注记

**接口**: `POST /api/v1/reader/annotations`  
**认证**: 需要（JWT Token）

**请求体**:
```json
{
  "bookId": "book123",
  "chapterId": "chapter50",
  "type": "highlight",
  "text": "选中的文本",
  "note": "我的笔记",
  "range": "100-150"
}
```

**请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/reader/annotations" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bookId": "book123",
    "chapterId": "chapter50",
    "type": "highlight",
    "text": "选中的文本",
    "note": "我的笔记",
    "range": "100-150"
  }'
```

**响应示例**:
```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "id": "annotation123",
    "userId": "user123",
    "bookId": "book123",
    "chapterId": "chapter50",
    "type": "highlight",
    "text": "选中的文本",
    "note": "我的笔记",
    "range": "100-150",
    "createdAt": "2023-12-31T23:59:59Z"
  }
}
```

**注记类型说明**:
- `bookmark`: 书签
- `highlight`: 高亮标注
- `note`: 笔记

#### 10.2 更新注记

**接口**: `PUT /api/v1/reader/annotations/{id}`  
**认证**: 需要（JWT Token）

**路径参数**:
- `id` (string, 必填) - 注记ID

**请求体**:
```json
{
  "text": "更新后的文本",
  "note": "更新后的笔记",
  "range": "100-160"
}
```

**请求示例**:
```bash
curl -X PUT "http://localhost:8080/api/v1/reader/annotations/annotation123" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "note": "更新后的笔记"
  }'
```

#### 10.3 删除注记

**接口**: `DELETE /api/v1/reader/annotations/{id}`  
**认证**: 需要（JWT Token）

**路径参数**:
- `id` (string, 必填) - 注记ID

**请求示例**:
```bash
curl -X DELETE "http://localhost:8080/api/v1/reader/annotations/annotation123" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 10.4 获取书籍注记列表

**接口**: `GET /api/v1/reader/annotations/book/{bookId}`  
**认证**: 需要（JWT Token）

**路径参数**:
- `bookId` (string, 必填) - 书籍ID

**查询参数**:
- `type` (string, 可选) - 注记类型筛选
- `page` (int, 可选, 默认1) - 页码
- `size` (int, 可选, 默认20) - 每页数量

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/annotations/book/book123?type=highlight" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "annotations": [
      {
        "id": "annotation123",
        "chapterId": "chapter50",
        "chapterTitle": "第五十章",
        "type": "highlight",
        "text": "选中的文本",
        "note": "我的笔记",
        "createdAt": "2023-12-31T23:59:59Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

#### 10.5 获取章节注记列表

**接口**: `GET /api/v1/reader/annotations/chapter/{chapterId}`  
**认证**: 需要（JWT Token）

**路径参数**:
- `chapterId` (string, 必填) - 章节ID

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/annotations/chapter/chapter50" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 10.6 获取注记统计

**接口**: `GET /api/v1/reader/annotations/stats`  
**认证**: 需要（JWT Token）

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/annotations/stats" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "totalCount": 500,
    "bookmarkCount": 100,
    "highlightCount": 300,
    "noteCount": 100
  }
}
```

#### 10.7 批量创建注记

**接口**: `POST /api/v1/reader/annotations/batch`  
**认证**: 需要（JWT Token）

**请求体**:
```json
{
  "annotations": [
    {
      "bookId": "book123",
      "chapterId": "chapter50",
      "type": "highlight",
      "text": "文本1",
      "range": "100-150"
    },
    {
      "bookId": "book123",
      "chapterId": "chapter50",
      "type": "note",
      "text": "文本2",
      "note": "笔记2",
      "range": "200-250"
    }
  ]
}
```

**请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/reader/annotations/batch" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "annotations": [...]
  }'
```

**说明**:
- 最多支持50个注记批量创建
- 适用于多端同步场景

---

### 11. 阅读设置

#### 11.1 获取阅读设置

**接口**: `GET /api/v1/reader/settings`  
**认证**: 需要（JWT Token）

**请求示例**:
```bash
curl -X GET "http://localhost:8080/api/v1/reader/settings" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "userId": "user123",
    "fontSize": 16,
    "fontFamily": "宋体",
    "lineHeight": 1.8,
    "backgroundColor": "#FFFFFF",
    "textColor": "#333333",
    "pageMode": "scroll",
    "autoSave": true,
    "showProgress": true,
    "theme": "default"
  }
}
```

#### 11.2 保存阅读设置

**接口**: `POST /api/v1/reader/settings`  
**认证**: 需要（JWT Token）

**请求体**:
```json
{
  "fontSize": 18,
  "fontFamily": "黑体",
  "lineHeight": 2.0,
  "backgroundColor": "#F5F5F5",
  "textColor": "#000000",
  "pageMode": "paginate",
  "autoSave": true,
  "showProgress": true,
  "theme": "night"
}
```

**请求示例**:
```bash
curl -X POST "http://localhost:8080/api/v1/reader/settings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fontSize": 18,
    "theme": "night"
  }'
```

**响应示例**:
```json
{
  "code": 200,
  "message": "保存成功",
  "data": null
}
```

#### 11.3 更新阅读设置

**接口**: `PUT /api/v1/reader/settings`  
**认证**: 需要（JWT Token）

**请求体**:
```json
{
  "fontSize": 20,
  "theme": "sepia"
}
```

**请求示例**:
```bash
curl -X PUT "http://localhost:8080/api/v1/reader/settings" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fontSize": 20,
    "theme": "sepia"
  }'
```

**说明**:
- 支持部分字段更新
- 未提供的字段保持原值不变

**设置字段说明**:
- `fontSize`: 字体大小（12-24）
- `fontFamily`: 字体类型
- `lineHeight`: 行高（1.0-3.0）
- `backgroundColor`: 背景颜色
- `textColor`: 文本颜色
- `pageMode`: 翻页模式（scroll/paginate）
- `autoSave`: 自动保存进度
- `showProgress`: 显示阅读进度
- `theme`: 主题（default/night/sepia）

---

## 🔐 认证和权限

### JWT Token 获取

首先需要通过用户登录接口获取Token：

```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

响应示例：
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "userId": "user123",
    "expiresAt": "2024-01-01T00:00:00Z"
  }
}
```

### 使用Token

在需要认证的接口中，将Token添加到Header：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### VIP权限验证

部分VIP章节需要相应的VIP等级：

- VIP等级0：普通用户
- VIP等级1：基础VIP
- VIP等级2：VIP Plus
- VIP等级3：VIP Pro
- VIP等级4：VIP Ultra

VIP权限由中间件自动验证，无需额外操作。

---

## 📊 错误码说明

| 错误码 | 说明 | 常见原因 |
|-------|------|---------|
| 200 | 成功 | 操作成功完成 |
| 201 | 创建成功 | 资源创建成功 |
| 400 | 请求参数错误 | 参数缺失或格式错误 |
| 401 | 未授权 | 未登录或Token无效 |
| 403 | 禁止访问 | 权限不足（如VIP权限） |
| 404 | 资源不存在 | 请求的资源不存在 |
| 409 | 冲突 | 资源已存在（如重复评分） |
| 500 | 服务器错误 | 服务器内部错误 |

---

## 💡 最佳实践

### 1. 分页请求
- 默认分页大小：20
- 建议最大分页大小：100
- 使用合理的分页大小以提高性能

### 2. 缓存策略
- 书籍详情、章节列表等数据会被缓存
- 缓存时间：5-30分钟
- 可以通过`Cache-Control: no-cache`头强制刷新

### 3. 请求频率限制
- 普通用户：100次/小时
- VIP用户：根据等级提升限额
- 超出限制返回429错误

### 4. 数据同步
- 阅读进度建议每30秒同步一次
- 注记数据支持批量同步
- 使用批量接口可以减少请求次数

### 5. 错误处理
```javascript
try {
  const response = await fetch(url, options);
  const data = await response.json();
  
  if (data.code !== 200) {
    // 处理业务错误
    console.error(data.message);
  }
  
  return data.data;
} catch (error) {
  // 处理网络错误
  console.error('请求失败', error);
}
```

---

## 📞 技术支持

- **API文档**: https://api.qingyu.com/docs

---

**文档版本**: v1.0  
**最后更新**: 2025-10-16  
**维护者**: 青羽后端团队

