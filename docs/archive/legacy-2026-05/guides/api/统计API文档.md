# 统计API文档

> **版本**: v1.0  
> **创建日期**: 2025-10-21

---

## 1. API概述

统计API提供书籍、章节、用户阅读数据统计功能。

**Base URL**: `/api/v1/stats`

---

## 2. 书籍统计

### 2.1 获取书籍统计

**接口**: `GET /stats/books/:bookId`

**响应**：
```json
{
  "code": 200,
  "data": {
    "bookId": "book_123",
    "reads": 10000,
    "favorites": 1500,
    "comments": 500,
    "averageRating": 4.5,
    "totalRatings": 300,
    "revenue": 5000.00
  }
}
```

### 2.2 获取章节热力图

**接口**: `GET /stats/books/:bookId/heatmap`

**响应**：
```json
{
  "code": 200,
  "data": {
    "heatmap": [
      {
        "chapterIndex": 1,
        "readCount": 10000,
        "dropRate": 0.0
      },
      {
        "chapterIndex": 2,
        "readCount": 9500,
        "dropRate": 5.0
      }
    ]
  }
}
```

---

## 3. 聚合统计

### 3.1 时间范围聚合

**接口**: `GET /stats/books/:bookId/aggregate`

**参数**：
- `timeRange` - daily/weekly/monthly

**响应**：
```json
{
  "code": 200,
  "data": {
    "totalReads": 50000,
    "uniqueReaders": 8000,
    "averageDuration": 600,
    "peakReadTime": "20:00-22:00"
  }
}
```

---

## 4. 行为记录

### 4.1 记录阅读行为

**接口**: `POST /stats/behavior`

**请求**：
```json
{
  "bookId": "book_123",
  "chapterId": "chapter_456",
  "action": "read",
  "duration": 600
}
```

---

**文档状态**: ✅ 已完成

