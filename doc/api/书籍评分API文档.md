# 书籍评分API文档

> **版本**: v1.0  
> **创建日期**: 2025-10-21

---

## 1. API概述

书籍评分API提供评分、评论、点赞等功能。

**Base URL**: `/api/v1/ratings`

---

## 2. 评分操作

### 2.1 评分书籍

**接口**: `POST /ratings/books/:bookId`

**请求**：
```json
{
  "rating": 4.5,
  "review": "非常精彩的作品！",
  "tags": ["剧情精彩", "文笔优美"]
}
```

**响应**：
```json
{
  "code": 201,
  "message": "评分成功",
  "data": {
    "ratingId": "rating_123",
    "rating": 4.5,
    "createdAt": "2025-10-21T10:00:00Z"
  }
}
```

**限制**：
- 评分范围：0-5星（支持0.5间隔）
- 每本书只能评分一次
- 24小时内可修改

### 2.2 修改评分

**接口**: `PUT /ratings/:ratingId`

**请求**：
```json
{
  "rating": 5.0,
  "review": "更新后的评论"
}
```

### 2.3 删除评分

**接口**: `DELETE /ratings/:ratingId`

---

## 3. 评分查询

### 3.1 获取书籍评分列表

**接口**: `GET /ratings/books/:bookId`

**参数**：
- `sortBy` - helpful/recent/rating
- `starFilter` - 1-5（星级筛选）
- `page` - 页码
- `pageSize` - 每页数量

**响应**：
```json
{
  "code": 200,
  "data": {
    "ratings": [
      {
        "id": "rating_123",
        "userId": "user_456",
        "username": "testuser",
        "rating": 4.5,
        "review": "非常精彩的作品！",
        "tags": ["剧情精彩"],
        "isHelpful": 10,
        "createdAt": "2025-10-21T10:00:00Z"
      }
    ],
    "total": 300,
    "averageRating": 4.3,
    "distribution": {
      "5": 100,
      "4": 80,
      "3": 50,
      "2": 30,
      "1": 40
    }
  }
}
```

### 3.2 获取用户评分历史

**接口**: `GET /ratings/users/:userId`

**响应**：
```json
{
  "code": 200,
  "data": {
    "ratings": [
      {
        "bookId": "book_123",
        "bookTitle": "测试书籍",
        "rating": 4.5,
        "createdAt": "2025-10-21T10:00:00Z"
      }
    ],
    "total": 20
  }
}
```

---

## 4. 评价互动

### 4.1 标记评价有用

**接口**: `POST /ratings/:ratingId/helpful`

**响应**：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    "helpfulCount": 11
  }
}
```

### 4.2 取消标记

**接口**: `DELETE /ratings/:ratingId/helpful`

---

## 5. 评价统计

### 5.1 获取评分分布

**接口**: `GET /ratings/books/:bookId/distribution`

**响应**：
```json
{
  "code": 200,
  "data": {
    "averageRating": 4.3,
    "totalRatings": 300,
    "distribution": {
      "5": {"count": 100, "percent": 33.3},
      "4": {"count": 80, "percent": 26.7},
      "3": {"count": 50, "percent": 16.7},
      "2": {"count": 30, "percent": 10.0},
      "1": {"count": 40, "percent": 13.3}
    }
  }
}
```

---

## 6. 评价标签

### 6.1 获取热门标签

**接口**: `GET /ratings/books/:bookId/tags`

**响应**：
```json
{
  "code": 200,
  "data": {
    "tags": [
      {"tag": "剧情精彩", "count": 150},
      {"tag": "文笔优美", "count": 120},
      {"tag": "人物丰满", "count": 80}
    ]
  }
}
```

---

**文档状态**: ✅ 已完成

