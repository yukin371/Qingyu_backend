# 用户社区功能API设计文档

**版本**：v1.0  
**创建时间**：2025-10-26  
**最后更新**：2025-10-26

---

## 文档概述

本文档描述青羽平台用户社区功能的后端API设计，包括用户主页、评论系统扩展等核心功能。

---

## 1. 用户主页API

### 1.1 获取用户公开信息

**接口描述**：获取指定用户的公开信息，用于展示用户主页。

**请求信息**：
- **URL**：`GET /api/v1/users/:userId/profile`
- **认证**：可选（登录后可查看更多信息）
- **权限**：公开接口

**路径参数**：
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | string | 是 | 用户ID |

**响应示例**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "user_id": "user123",
    "username": "author_name",
    "nickname": "作者昵称",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "这是用户的个人简介",
    "role": "writer",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**字段说明**：
- `user_id`：用户唯一标识
- `username`：用户名（不可修改）
- `nickname`：昵称（可修改）
- `avatar`：头像URL
- `bio`：个人简介（最多500字）
- `role`：用户角色（reader/writer/admin）
- `created_at`：注册时间

**注意事项**：
- 此接口不返回敏感信息（邮箱、手机号等）
- 区别于 `/api/v1/users/profile`（获取当前用户详细信息）

---

### 1.2 获取用户作品列表

**接口描述**：获取指定用户的已发布作品列表（仅作者角色）。

**请求信息**：
- **URL**：`GET /api/v1/users/:userId/books`
- **认证**：可选
- **权限**：公开接口

**路径参数**：
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | string | 是 | 用户ID |

**查询参数**：
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| size | int | 否 | 20 | 每页数量（最大100） |
| status | string | 否 | - | 状态筛选（published/completed） |

**响应示例**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "books": [
      {
        "book_id": "book123",
        "title": "示例作品",
        "cover": "https://example.com/cover.jpg",
        "description": "这是一部精彩的作品",
        "category": "玄幻",
        "status": "published",
        "word_count": 100000,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 20
  }
}
```

**业务规则**：
- 只返回已发布的作品（status=published）
- 按创建时间倒序排列
- 不返回草稿或下架的作品

**TODO**：
- 需要在 BookService 中实现 `GetBooksByAuthor` 方法
- 需要添加作品统计信息（阅读量、收藏数等）

---

## 2. 评论系统扩展API

### 2.1 获取评论完整线程

**接口描述**：获取评论及其所有回复的完整树状结构。

**请求信息**：
- **URL**：`GET /api/v1/reader/comments/:id/thread`
- **认证**：可选
- **权限**：公开接口

**路径参数**：
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 评论ID |

**响应示例**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": "comment123",
    "content": "这是一条评论",
    "user": {
      "user_id": "user123",
      "username": "reader",
      "nickname": "读者昵称",
      "avatar": "https://example.com/avatar.jpg",
      "role": "reader"
    },
    "like_count": 10,
    "reply_count": 5,
    "is_liked": false,
    "created_at": "2024-01-01T00:00:00Z",
    "replies": [
      {
        "id": "reply123",
        "content": "这是一条回复",
        "user": {...},
        "like_count": 2,
        "reply_count": 1,
        "is_liked": false,
        "created_at": "2024-01-01T01:00:00Z",
        "replies": [...]
      }
    ]
  }
}
```

**数据结构说明**：
- `replies`：嵌套的回复数组，支持无限层级
- `is_liked`：当前用户是否已点赞（未登录时为false）

**性能考虑**：
- 递归查询可能影响性能
- 建议限制最大层级（如3层）
- 考虑使用缓存优化热门评论

**TODO**：
- 需要在 CommentService 中实现 `GetCommentThread` 方法
- 需要优化递归查询性能

---

### 2.2 获取热门评论

**接口描述**：获取指定书籍的热门评论列表（按点赞数、回复数排序）。

**请求信息**：
- **URL**：`GET /api/v1/reader/comments/top`
- **认证**：可选
- **权限**：公开接口

**查询参数**：
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| book_id | string | 是 | - | 书籍ID |
| limit | int | 否 | 10 | 返回数量（最大50） |

**响应示例**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "comments": [
      {
        "id": "comment123",
        "content": "这是一条热门评论",
        "user": {...},
        "like_count": 100,
        "reply_count": 50,
        "is_liked": false,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 10
  }
}
```

**排序规则**：
- 优先级：点赞数 > 回复数 > 创建时间
- 综合权重计算：`score = like_count * 2 + reply_count * 1`
- 时间衰减：较新的评论有加权

**临时实现**：
- 当前使用常规评论列表，按 "hot" 排序
- 待 Service 层实现完整的热门评论算法

**TODO**：
- 需要在 CommentService 中实现 `GetTopComments` 方法
- 实现热门度计算算法
- 添加缓存机制

---

### 2.3 获取评论回复列表

**接口描述**：获取指定评论的直接回复列表（分页）。

**请求信息**：
- **URL**：`GET /api/v1/reader/comments/:id/replies`
- **认证**：可选
- **权限**：公开接口

**路径参数**：
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 评论ID |

**查询参数**：
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| size | int | 否 | 20 | 每页数量（最大50） |

**响应示例**：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "replies": [
      {
        "id": "reply123",
        "content": "这是一条回复",
        "user": {...},
        "like_count": 5,
        "reply_count": 2,
        "is_liked": false,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "size": 20
  }
}
```

**业务规则**：
- 只返回直接回复（parent_id = :id）
- 不返回嵌套的子回复
- 按时间倒序排列

**临时实现**：
- 当前返回空列表
- 前端可以正常渲染，不会出错

**TODO**：
- 需要在 CommentService 中实现 `GetCommentReplies` 方法
- 支持排序选项（最新/最热）

---

## 3. 统一响应格式

### 3.1 成功响应

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {...}
}
```

### 3.2 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "error": "详细错误信息"
}
```

**常见错误码**：
- `400`：参数错误
- `401`：未授权（需要登录）
- `403`：权限不足
- `404`：资源不存在
- `500`：服务器内部错误

---

## 4. 认证和授权

### 4.1 认证机制

- **认证方式**：JWT Token
- **Token传递**：`Authorization: Bearer {token}`
- **Token有效期**：7天

### 4.2 可选认证接口

以下接口支持可选认证（登录后获得更多信息）：
- 用户主页查询：显示关注状态
- 评论查询：显示点赞状态

### 4.3 公开接口

无需认证即可访问：
- 获取用户公开信息
- 获取用户作品列表
- 获取评论列表和详情

---

## 5. 数据库设计

### 5.1 用户表（users）

已存在，无需修改。

### 5.2 评论表（comments）

**关键字段**：
- `id`：评论ID
- `book_id`：书籍ID
- `chapter_id`：章节ID（可选）
- `user_id`：用户ID
- `parent_id`：父评论ID（用于回复）
- `content`：评论内容
- `like_count`：点赞数
- `reply_count`：回复数
- `created_at`：创建时间
- `updated_at`：更新时间

**索引**：
- `book_id` + `created_at`（查询书籍评论）
- `parent_id` + `created_at`（查询回复）
- `user_id` + `created_at`（查询用户评论）

---

## 6. 实施计划

### 6.1 已完成

- ✅ API接口定义和注释
- ✅ 基础的请求处理和参数验证
- ✅ 统一的错误处理
- ✅ Swagger文档注释

### 6.2 待实施（优先级：高）

1. **CommentService 扩展**
   - [ ] 实现 `GetCommentThread` 方法
   - [ ] 实现 `GetTopComments` 方法
   - [ ] 实现 `GetCommentReplies` 方法

2. **BookService 扩展**
   - [ ] 实现 `GetBooksByAuthor` 方法
   - [ ] 添加作品统计信息

3. **性能优化**
   - [ ] 添加Redis缓存（热门评论、用户信息）
   - [ ] 优化递归查询（评论线程）

### 6.3 待实施（优先级：中）

1. **功能增强**
   - [ ] 关注/粉丝功能
   - [ ] 用户统计信息（作品数、粉丝数）
   - [ ] 评论举报功能

2. **数据分析**
   - [ ] 用户活跃度统计
   - [ ] 评论热度计算

---

## 7. 测试建议

### 7.1 单元测试

```go
// 测试获取用户公开信息
func TestGetUserProfile(t *testing.T) {
    // 测试正常用户
    // 测试不存在的用户
    // 测试各种角色的用户
}
```

### 7.2 集成测试

- 测试完整的用户主页访问流程
- 测试评论线程的递归查询
- 测试分页功能

### 7.3 性能测试

- 并发请求测试（1000 QPS）
- 递归查询深度测试
- 大数据量分页测试

---

## 8. 安全考虑

### 8.1 数据脱敏

- 用户主页不返回敏感信息（邮箱、手机）
- 只返回必要的公开信息

### 8.2 防刷保护

- 添加访问频率限制（Rate Limiting）
- 防止爬虫抓取用户信息

### 8.3 SQL注入防护

- 使用参数化查询
- 输入验证和清理

---

## 9. 版本历史

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|----------|------|
| v1.0 | 2025-10-26 | 初始版本，包含用户主页和评论扩展API | 开发团队 |

---

**文档维护者**：青羽后端开发团队  
**联系方式**：backend@qingyu.com


