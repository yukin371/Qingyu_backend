# 阅读器系统API文档

> **文档版本**: v1.0
> **创建日期**: 2025-10-09
> **最后更新**: 2025-10-09

## Page Role

- current-reference
- current-owner: `docs/api/reader/`
- current-bounded: 当前阅读器系统 API 详细参考文档，负责阅读器侧接口说明

## Recommended Read Path

1. 先读 `README.md`。
2. 再读 `../usage_guide.md`。
3. 最后读本页的详细接口定义。

## Boundary

- 本页是阅读器系统详细参考手册，不是 API 总入口。
- 当前阅读器专题 owner 仍在 `docs/api/reader/README.md`。

## Quick Section Map

- 目录
- 概述
- 通用说明
- 章节API
- 阅读进度API
- 标注API
- 阅读设置API
- 错误码说明

## Quick Takeaways

- 这是当前阅读器系统 API 的详细参考文档。

## Skip Guide

- 只看阅读器入口：回 `README.md`。

## 📋 目录

- [1. 概述](#1-概述)
- [2. 通用说明](#2-通用说明)
- [3. 章节API](#3-章节api)
- [4. 阅读进度API](#4-阅读进度api)
- [5. 标注API](#5-标注api)
- [6. 阅读设置API](#6-阅读设置api)
- [7. 错误码说明](#7-错误码说明)

---

## 1. 概述

### 1.1 接口基础信息

- **Base URL**: `http://your-domain/api/v1`
- **协议**: HTTP/HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8
- **认证方式**: JWT Token（大部分接口需要）

### 1.2 功能模块

```
阅读器系统API
├── 章节模块 - 章节内容获取、导航
├── 进度模块 - 阅读进度记录、统计
├── 标注模块 - 书签、高亮、笔记管理
└── 设置模块 - 阅读偏好设置
```

### 1.3 API概览

| 模块           | 接口数量     | 主要功能             | 认证要求 |
| -------------- | ------------ | -------------------- | -------- |
| 章节管理       | 6            | 章节获取、内容、导航 | 部分需要 |
| 阅读进度       | 8            | 进度记录、统计、历史 | 需要     |
| 标注管理       | 13           | 书签、高亮、笔记     | 需要     |
| 阅读设置       | 3            | 个性化设置           | 需要     |
| **总计** | **30** | -                    | -        |

---

## 2. 通用说明

### 2.1 统一响应格式

#### 成功响应

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 具体数据
  }
}
```

#### 分页响应

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "错误详细描述",
  "error": "具体错误信息"
}
```

### 2.2 认证说明

大部分接口需要用户认证，请在请求头中携带JWT Token：

```
Authorization: Bearer {your_token}
```

Token获取方式：通过用户登录接口获取。

### 2.3 数据模型

#### Chapter（章节）

```json
{
  "id": "507f1f77bcf86cd799439011",
  "bookId": "507f1f77bcf86cd799439012",
  "chapterNum": 1,
  "title": "第一章 开始",
  "content": "章节内容...",
  "wordCount": 3000,
  "isVip": false,
  "isPublished": true,
  "publishedAt": "2024-01-01T00:00:00Z",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### ReadingProgress（阅读进度）

```json
{
  "id": "507f1f77bcf86cd799439013",
  "userId": "507f1f77bcf86cd799439014",
  "bookId": "507f1f77bcf86cd799439012",
  "chapterId": "507f1f77bcf86cd799439011",
  "progress": 0.65,
  "readingTime": 3600,
  "lastReadAt": "2024-01-01T12:00:00Z",
  "isFinished": false,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```

#### Annotation（标注）

```json
{
  "id": "507f1f77bcf86cd799439015",
  "userId": "507f1f77bcf86cd799439014",
  "bookId": "507f1f77bcf86cd799439012",
  "chapterId": "507f1f77bcf86cd799439011",
  "type": "highlight",
  "text": "标注的文本内容",
  "note": "我的笔记",
  "range": "100-200",
  "isPublic": false,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**标注类型（type）**：

- `bookmark` - 书签
- `highlight` - 高亮
- `note` - 笔记

#### ReadingSettings（阅读设置）

```json
{
  "id": "507f1f77bcf86cd799439016",
  "userId": "507f1f77bcf86cd799439014",
  "fontSize": 16,
  "fontFamily": "Microsoft YaHei",
  "lineHeight": 1.8,
  "backgroundColor": "#FFFFFF",
  "textColor": "#000000",
  "pageMode": "scroll",
  "autoSave": true,
  "showProgress": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

**翻页模式（pageMode）**：

- `scroll` - 滚动模式
- `paginate` - 分页模式

---

## 3. 章节API

### 3.1 获取章节信息

**接口说明**

根据章节ID获取章节的基本信息（不包含内容）。

**接口详情**

- **URL**: `/api/v1/reader/chapters/{id}`
- **Method**: `GET`
- **认证**: 不需要

**路径参数**

| 参数 | 类型   | 必填 | 说明   |
| ---- | ------ | ---- | ------ |
| id   | string | 是   | 章节ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "bookId": "507f1f77bcf86cd799439012",
    "chapterNum": 1,
    "title": "第一章 开始",
    "wordCount": 3000,
    "isVip": false,
    "publishedAt": "2024-01-01T00:00:00Z"
  }
}
```

**错误示例**

```json
{
  "code": 404,
  "message": "章节不存在",
  "error": "chapter not found"
}
```

---

### 3.2 获取章节内容

**接口说明**

获取章节的完整内容，支持VIP权限验证。

**接口详情**

- **URL**: `/api/v1/reader/chapters/{id}/content`
- **Method**: `GET`
- **认证**: **需要**

**路径参数**

| 参数 | 类型   | 必填 | 说明   |
| ---- | ------ | ---- | ------ |
| id   | string | 是   | 章节ID |

**请求头**

```
Authorization: Bearer {your_token}
```

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "content": "章节的完整文本内容..."
  }
}
```

**错误示例**

```json
{
  "code": 403,
  "message": "获取章节内容失败",
  "error": "VIP章节需要权限"
}
```

**权限说明**

- VIP章节需要用户具有相应权限
- 如果用户无权限，返回403错误
- 前端应引导用户购买或开通VIP

---

### 3.3 获取书籍章节列表

**接口说明**

获取指定书籍的章节列表，支持分页。

**接口详情**

- **URL**: `/api/v1/reader/chapters`
- **Method**: `GET`
- **认证**: 不需要

**查询参数**

| 参数   | 类型   | 必填 | 默认值 | 说明     |
| ------ | ------ | ---- | ------ | -------- |
| bookId | string | 是   | -      | 书籍ID   |
| page   | int    | 否   | 1      | 页码     |
| size   | int    | 否   | 20     | 每页数量 |

**请求示例**

```
GET /api/v1/reader/chapters?bookId=507f1f77bcf86cd799439012&page=1&size=20
```

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "chapters": [
      {
        "id": "507f1f77bcf86cd799439011",
        "chapterNum": 1,
        "title": "第一章 开始",
        "wordCount": 3000,
        "isVip": false
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

---

### 3.4 获取章节导航

**接口说明**

获取当前章节的上一章和下一章信息，用于章节导航。

**接口详情**

- **URL**: `/api/v1/reader/chapters/navigation`
- **Method**: `GET`
- **认证**: 不需要

**查询参数**

| 参数       | 类型   | 必填 | 说明       |
| ---------- | ------ | ---- | ---------- |
| bookId     | string | 是   | 书籍ID     |
| chapterNum | int    | 是   | 当前章节号 |

**请求示例**

```
GET /api/v1/reader/chapters/navigation?bookId=507f1f77bcf86cd799439012&chapterNum=5
```

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "prevChapter": {
      "id": "507f1f77bcf86cd799439010",
      "chapterNum": 4,
      "title": "第四章"
    },
    "nextChapter": {
      "id": "507f1f77bcf86cd799439012",
      "chapterNum": 6,
      "title": "第六章"
    }
  }
}
```

**说明**

- 如果没有上一章，`prevChapter`为 `null`
- 如果没有下一章，`nextChapter`为 `null`

---

### 3.5 获取第一章

**接口说明**

获取书籍的第一章信息。

**接口详情**

- **URL**: `/api/v1/reader/chapters/first`
- **Method**: `GET`
- **认证**: 不需要

**查询参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "chapterNum": 1,
    "title": "第一章 开始"
  }
}
```

---

### 3.6 获取最后一章

**接口说明**

获取书籍的最后一章信息。

**接口详情**

- **URL**: `/api/v1/reader/chapters/last`
- **Method**: `GET`
- **认证**: 不需要

**查询参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "507f1f77bcf86cd799439099",
    "chapterNum": 100,
    "title": "第一百章 大结局"
  }
}
```

---

## 4. 阅读进度API

### 4.1 获取阅读进度

**接口说明**

获取用户在某本书的阅读进度。

**接口详情**

- **URL**: `/api/v1/reader/progress/{bookId}`
- **Method**: `GET`
- **认证**: **需要**

**路径参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "507f1f77bcf86cd799439013",
    "userId": "507f1f77bcf86cd799439014",
    "bookId": "507f1f77bcf86cd799439012",
    "chapterId": "507f1f77bcf86cd799439011",
    "progress": 0.65,
    "readingTime": 3600,
    "lastReadAt": "2024-01-01T12:00:00Z",
    "isFinished": false
  }
}
```

---

### 4.2 保存阅读进度

**接口说明**

保存用户的阅读进度（章节位置和百分比）。

**接口详情**

- **URL**: `/api/v1/reader/progress`
- **Method**: `POST`
- **认证**: **需要**

**请求体**

```json
{
  "bookId": "507f1f77bcf86cd799439012",
  "chapterId": "507f1f77bcf86cd799439011",
  "progress": 0.65
}
```

**字段说明**

| 字段      | 类型   | 必填 | 说明                  |
| --------- | ------ | ---- | --------------------- |
| bookId    | string | 是   | 书籍ID                |
| chapterId | string | 是   | 章节ID                |
| progress  | float  | 是   | 进度百分比（0-1之间） |

**响应示例**

```json
{
  "code": 200,
  "message": "保存成功",
  "data": null
}
```

---

### 4.3 更新阅读时长

**接口说明**

更新用户在某本书的阅读时长（秒）。

**接口详情**

- **URL**: `/api/v1/reader/progress/reading-time`
- **Method**: `PUT`
- **认证**: **需要**

**请求体**

```json
{
  "bookId": "507f1f77bcf86cd799439012",
  "duration": 300
}
```

**字段说明**

| 字段     | 类型   | 必填 | 说明                    |
| -------- | ------ | ---- | ----------------------- |
| bookId   | string | 是   | 书籍ID                  |
| duration | int    | 是   | 阅读时长（秒），最小值1 |

**响应示例**

```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

**使用建议**

- 建议每5分钟上报一次阅读时长
- 或在用户退出阅读器时上报

---

### 4.4 获取最近阅读记录

**接口说明**

获取用户最近阅读的书籍列表。

**接口详情**

- **URL**: `/api/v1/reader/progress/recent`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数  | 类型 | 必填 | 默认值 | 说明     |
| ----- | ---- | ---- | ------ | -------- |
| limit | int  | 否   | 20     | 数量限制 |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "507f1f77bcf86cd799439013",
      "bookId": "507f1f77bcf86cd799439012",
      "chapterId": "507f1f77bcf86cd799439011",
      "progress": 0.65,
      "lastReadAt": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### 4.5 获取阅读历史

**接口说明**

获取用户的完整阅读历史，支持分页。

**接口详情**

- **URL**: `/api/v1/reader/progress/history`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 说明     |
| ---- | ---- | ---- | ------ | -------- |
| page | int  | 否   | 1      | 页码     |
| size | int  | 否   | 20     | 每页数量 |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "progresses": [...],
    "total": 50,
    "page": 1,
    "size": 20
  }
}
```

---

### 4.6 获取阅读统计

**接口说明**

获取用户的阅读统计数据（总阅读时长、已读/未读书籍数）。

**接口详情**

- **URL**: `/api/v1/reader/progress/stats`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数   | 类型   | 必填 | 默认值 | 说明                           |
| ------ | ------ | ---- | ------ | ------------------------------ |
| period | string | 否   | all    | 统计周期：all/today/week/month |

**period支持的值**

- `all` - 总计
- `today` - 今天
- `week` - 本周
- `month` - 本月

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "totalReadingTime": 36000,
    "unfinishedCount": 5,
    "finishedCount": 10,
    "period": "all"
  }
}
```

**字段说明**

- `totalReadingTime`: 总阅读时长（秒）
- `unfinishedCount`: 未读完书籍数量
- `finishedCount`: 已读完书籍数量

---

### 4.7 获取未读完的书籍

**接口说明**

获取用户未读完的书籍列表。

**接口详情**

- **URL**: `/api/v1/reader/progress/unfinished`
- **Method**: `GET`
- **认证**: **需要**

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "bookId": "507f1f77bcf86cd799439012",
      "progress": 0.65,
      "lastReadAt": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### 4.8 获取已读完的书籍

**接口说明**

获取用户已读完的书籍列表。

**接口详情**

- **URL**: `/api/v1/reader/progress/finished`
- **Method**: `GET`
- **认证**: **需要**

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "bookId": "507f1f77bcf86cd799439012",
      "progress": 1.0,
      "finishedAt": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

## 5. 标注API

### 5.1 创建标注

**接口说明**

创建书签、高亮或笔记标注。

**接口详情**

- **URL**: `/api/v1/reader/annotations`
- **Method**: `POST`
- **认证**: **需要**

**请求体**

```json
{
  "bookId": "507f1f77bcf86cd799439012",
  "chapterId": "507f1f77bcf86cd799439011",
  "type": "highlight",
  "text": "要标注的文本内容",
  "note": "我的笔记",
  "range": "100-200"
}
```

**字段说明**

| 字段      | 类型   | 必填 | 说明                              |
| --------- | ------ | ---- | --------------------------------- |
| bookId    | string | 是   | 书籍ID                            |
| chapterId | string | 是   | 章节ID                            |
| type      | string | 是   | 标注类型：bookmark/highlight/note |
| text      | string | 否   | 标注的文本内容                    |
| note      | string | 否   | 注释/笔记内容                     |
| range     | string | 否   | 标注范围（如：100-200）           |

**响应示例**

```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "userId": "507f1f77bcf86cd799439014",
    "bookId": "507f1f77bcf86cd799439012",
    "chapterId": "507f1f77bcf86cd799439011",
    "type": "highlight",
    "text": "要标注的文本内容",
    "note": "我的笔记",
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

---

### 5.2 更新标注

**接口说明**

更新已有标注的内容。

**接口详情**

- **URL**: `/api/v1/reader/annotations/{id}`
- **Method**: `PUT`
- **认证**: **需要**

**路径参数**

| 参数 | 类型   | 必填 | 说明   |
| ---- | ------ | ---- | ------ |
| id   | string | 是   | 标注ID |

**请求体**

```json
{
  "text": "更新后的文本",
  "note": "更新后的笔记",
  "range": "100-250"
}
```

**响应示例**

```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

---

### 5.3 删除标注

**接口说明**

删除指定的标注。

**接口详情**

- **URL**: `/api/v1/reader/annotations/{id}`
- **Method**: `DELETE`
- **认证**: **需要**

**路径参数**

| 参数 | 类型   | 必填 | 说明   |
| ---- | ------ | ---- | ------ |
| id   | string | 是   | 标注ID |

**响应示例**

```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

---

### 5.4 获取章节标注

**接口说明**

获取指定章节的所有标注（当前用户）。

**接口详情**

- **URL**: `/api/v1/reader/annotations/chapter`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数      | 类型   | 必填 | 说明   |
| --------- | ------ | ---- | ------ |
| bookId    | string | 是   | 书籍ID |
| chapterId | string | 是   | 章节ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "507f1f77bcf86cd799439015",
      "type": "highlight",
      "text": "标注文本",
      "note": "笔记内容",
      "range": "100-200"
    }
  ]
}
```

---

### 5.5 获取书籍标注

**接口说明**

获取整本书的所有标注（当前用户）。

**接口详情**

- **URL**: `/api/v1/reader/annotations/book`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [...]
}
```

---

### 5.6 获取笔记列表

**接口说明**

获取用户在某本书的所有笔记。

**接口详情**

- **URL**: `/api/v1/reader/annotations/notes`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "507f1f77bcf86cd799439015",
      "chapterId": "507f1f77bcf86cd799439011",
      "text": "原文内容",
      "note": "我的笔记",
      "createdAt": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### 5.7 搜索笔记

**接口说明**

在用户的所有笔记中搜索关键词。

**接口详情**

- **URL**: `/api/v1/reader/annotations/notes/search`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数    | 类型   | 必填 | 说明       |
| ------- | ------ | ---- | ---------- |
| keyword | string | 是   | 搜索关键词 |

**响应示例**

```json
{
  "code": 200,
  "message": "搜索成功",
  "data": [...]
}
```

---

### 5.8 获取书签列表

**接口说明**

获取用户在某本书的所有书签。

**接口详情**

- **URL**: `/api/v1/reader/annotations/bookmarks`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "507f1f77bcf86cd799439015",
      "chapterId": "507f1f77bcf86cd799439011",
      "chapterTitle": "第一章",
      "createdAt": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### 5.9 获取最新书签

**接口说明**

获取用户在某本书的最新书签（用于快速跳转）。

**接口详情**

- **URL**: `/api/v1/reader/annotations/bookmarks/latest`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "507f1f77bcf86cd799439015",
    "chapterId": "507f1f77bcf86cd799439011",
    "chapterTitle": "第一章",
    "createdAt": "2024-01-01T12:00:00Z"
  }
}
```

---

### 5.10 获取高亮列表

**接口说明**

获取用户在某本书的所有高亮标注。

**接口详情**

- **URL**: `/api/v1/reader/annotations/highlights`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数   | 类型   | 必填 | 说明   |
| ------ | ------ | ---- | ------ |
| bookId | string | 是   | 书籍ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "507f1f77bcf86cd799439015",
      "chapterId": "507f1f77bcf86cd799439011",
      "text": "高亮的文本",
      "range": "100-200",
      "createdAt": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

### 5.11 获取最近标注

**接口说明**

获取用户最近创建的标注列表。

**接口详情**

- **URL**: `/api/v1/reader/annotations/recent`
- **Method**: `GET`
- **认证**: **需要**

**查询参数**

| 参数  | 类型 | 必填 | 默认值 | 说明     |
| ----- | ---- | ---- | ------ | -------- |
| limit | int  | 否   | 20     | 数量限制 |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [...]
}
```

---

### 5.12 获取公开标注

**接口说明**

获取某章节的公开标注（其他用户分享的标注）。

**接口详情**

- **URL**: `/api/v1/reader/annotations/public`
- **Method**: `GET`
- **认证**: 不需要

**查询参数**

| 参数      | 类型   | 必填 | 说明   |
| --------- | ------ | ---- | ------ |
| bookId    | string | 是   | 书籍ID |
| chapterId | string | 是   | 章节ID |

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "507f1f77bcf86cd799439015",
      "userId": "507f1f77bcf86cd799439014",
      "userName": "张三",
      "text": "标注文本",
      "note": "笔记内容",
      "likeCount": 10,
      "createdAt": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

## 6. 阅读设置API

### 6.1 获取阅读设置

**接口说明**

获取用户的阅读偏好设置。

**接口详情**

- **URL**: `/api/v1/reader/settings`
- **Method**: `GET`
- **认证**: **需要**

**响应示例**

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "507f1f77bcf86cd799439016",
    "userId": "507f1f77bcf86cd799439014",
    "fontSize": 16,
    "fontFamily": "Microsoft YaHei",
    "lineHeight": 1.8,
    "backgroundColor": "#FFFFFF",
    "textColor": "#000000",
    "pageMode": "scroll",
    "autoSave": true,
    "showProgress": true
  }
}
```

---

### 6.2 保存阅读设置

**接口说明**

保存完整的阅读设置（覆盖式保存）。

**接口详情**

- **URL**: `/api/v1/reader/settings`
- **Method**: `POST`
- **认证**: **需要**

**请求体**

```json
{
  "fontSize": 16,
  "fontFamily": "Microsoft YaHei",
  "lineHeight": 1.8,
  "backgroundColor": "#FFFFFF",
  "textColor": "#000000",
  "pageMode": "scroll",
  "autoSave": true,
  "showProgress": true
}
```

**响应示例**

```json
{
  "code": 200,
  "message": "保存成功",
  "data": null
}
```

---

### 6.3 更新阅读设置

**接口说明**

部分更新阅读设置（只更新提供的字段）。

**接口详情**

- **URL**: `/api/v1/reader/settings`
- **Method**: `PUT`
- **认证**: **需要**

**请求体**

```json
{
  "fontSize": 18,
  "lineHeight": 2.0
}
```

**字段说明**（所有字段都是可选的）

| 字段            | 类型   | 说明                      |
| --------------- | ------ | ------------------------- |
| fontSize        | int    | 字体大小（像素）          |
| fontFamily      | string | 字体名称                  |
| lineHeight      | float  | 行高                      |
| backgroundColor | string | 背景颜色（HEX）           |
| textColor       | string | 文字颜色（HEX）           |
| pageMode        | string | 翻页模式：scroll/paginate |
| autoSave        | bool   | 是否自动保存进度          |
| showProgress    | bool   | 是否显示进度              |

**响应示例**

```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

---

## 7. 错误码说明

### 7.1 HTTP状态码

| 状态码 | 说明                  | 使用场景         |
| ------ | --------------------- | ---------------- |
| 200    | OK                    | 请求成功         |
| 201    | Created               | 创建成功         |
| 400    | Bad Request           | 请求参数错误     |
| 401    | Unauthorized          | 未授权/Token无效 |
| 403    | Forbidden             | 无权限访问       |
| 404    | Not Found             | 资源不存在       |
| 500    | Internal Server Error | 服务器错误       |

### 7.2 业务错误码

| 错误信息        | 说明             | 解决方案             |
| --------------- | ---------------- | -------------------- |
| 未授权          | Token缺失或无效  | 重新登录获取Token    |
| 章节不存在      | 章节ID无效       | 检查章节ID是否正确   |
| VIP章节需要权限 | 用户无VIP权限    | 引导用户开通VIP      |
| 参数错误        | 请求参数格式错误 | 检查参数格式和必填项 |
| 标注不存在      | 标注ID无效       | 检查标注ID是否正确   |

---

## 附录

### A. 快速开始示例

#### JavaScript/Axios

```javascript
import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://your-domain/api/v1',
  headers: {
    'Authorization': 'Bearer your_token'
  }
});

// 获取章节内容
async function getChapterContent(chapterId) {
  const response = await apiClient.get(`/reader/chapters/${chapterId}/content`);
  return response.data.data.content;
}

// 保存阅读进度
async function saveProgress(bookId, chapterId, progress) {
  await apiClient.post('/reader/progress', {
    bookId,
    chapterId,
    progress
  });
}

// 创建书签
async function createBookmark(bookId, chapterId) {
  await apiClient.post('/reader/annotations', {
    bookId,
    chapterId,
    type: 'bookmark'
  });
}
```

#### Python/Requests

```python
import requests

BASE_URL = 'http://your-domain/api/v1'
TOKEN = 'your_token'

headers = {
    'Authorization': f'Bearer {TOKEN}'
}

# 获取章节内容
def get_chapter_content(chapter_id):
    response = requests.get(
        f'{BASE_URL}/reader/chapters/{chapter_id}/content',
        headers=headers
    )
    return response.json()['data']['content']

# 保存阅读进度
def save_progress(book_id, chapter_id, progress):
    requests.post(
        f'{BASE_URL}/reader/progress',
        json={
            'bookId': book_id,
            'chapterId': chapter_id,
            'progress': progress
        },
        headers=headers
    )
```

#### cURL

```bash
# 获取章节内容
curl -X GET "http://your-domain/api/v1/reader/chapters/{id}/content" \
  -H "Authorization: Bearer your_token"

# 保存阅读进度
curl -X POST "http://your-domain/api/v1/reader/progress" \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "bookId": "507f1f77bcf86cd799439012",
    "chapterId": "507f1f77bcf86cd799439011",
    "progress": 0.65
  }'

# 创建高亮
curl -X POST "http://your-domain/api/v1/reader/annotations" \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "bookId": "507f1f77bcf86cd799439012",
    "chapterId": "507f1f77bcf86cd799439011",
    "type": "highlight",
    "text": "要高亮的文本",
    "range": "100-200"
  }'
```

### B. 常见使用场景

#### 场景1：打开书籍阅读

```javascript
// 1. 获取阅读进度
const progress = await getReadingProgress(bookId);

// 2. 如果有进度，跳转到上次阅读的章节
if (progress && progress.chapterId) {
  loadChapter(progress.chapterId);
} else {
  // 否则跳转到第一章
  const firstChapter = await getFirstChapter(bookId);
  loadChapter(firstChapter.id);
}
```

#### 场景2：阅读时自动保存进度

```javascript
// 每5秒保存一次进度
setInterval(() => {
  const progress = calculateProgress(); // 计算当前进度
  saveProgress(bookId, chapterId, progress);
}, 5000);

// 退出时保存进度和时长
window.addEventListener('beforeunload', () => {
  saveProgress(bookId, chapterId, currentProgress);
  updateReadingTime(bookId, readingDuration);
});
```

#### 场景3：加载章节标注

```javascript
// 加载章节内容和标注
async function loadChapterWithAnnotations(bookId, chapterId) {
  // 并行请求
  const [content, annotations] = await Promise.all([
    getChapterContent(chapterId),
    getChapterAnnotations(bookId, chapterId)
  ]);
  
  // 渲染内容和标注
  renderContent(content, annotations);
}
```

### C. 更新日志

| 版本 | 日期       | 更新内容                    |
| ---- | ---------- | --------------------------- |
| v1.0 | 2025-10-09 | 初始版本，包含30个阅读器API |

---

**最后更新**: 2025-10-09
