# API 参考文档

## 概述

青羽写作平台后端API参考文档，提供所有可用API端点的详细说明。

## 基础信息

- **Base URL**: `/api/v1`
- **认证方式**: Bearer Token (JWT)
- **内容类型**: `application/json`
- **字符编码**: UTF-8

## 认证

大多数API端点需要JWT认证。在请求头中添加：

```
Authorization: Bearer <your-jwt-token>
```

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": { ... },
  "timestamp": 1234567890,
  "request_id": "uuid"
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误描述",
  "details": "详细错误信息",
  "timestamp": 1234567890,
  "request_id": "uuid"
}
```

---

## 管理员API (api/v1/admin/)

### 用户管理

#### 获取用户列表
```http
GET /api/v1/admin/users?page=1&limit=20&status=active&keyword=
```

**查询参数**:
- `page` (int, 可选): 页码，默认1
- `limit` (int, 可选): 每页数量，默认20，最大100
- `status` (string, 可选): 用户状态筛选
- `keyword` (string, 可选): 搜索关键词

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "users": [
      {
        "id": "user123",
        "username": "testuser",
        "email": "test@example.com",
        "nickname": "测试用户",
        "status": "active",
        "role": "user",
        "created_at": "2026-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

#### 创建用户
```http
POST /api/v1/admin/users
```

**请求体**:
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "SecurePass123!",
  "nickname": "新用户",
  "role": "user"
}
```

#### 更新用户信息
```http
PUT /api/v1/admin/users/{userId}
```

**路径参数**:
- `userId` (string, 必需): 用户ID

**请求体**:
```json
{
  "nickname": "更新昵称",
  "status": "active",
  "role": "author"
}
```

#### 删除用户
```http
DELETE /api/v1/admin/users/{userId}
```

### 权限管理

#### 获取所有权限
```http
GET /api/v1/admin/permissions
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "code": "user.read",
      "name": "读取用户信息",
      "description": "允许读取用户基本信息",
      "category": "user"
    }
  ]
}
```

#### 获取权限详情
```http
GET /api/v1/admin/permissions/{code}
```

#### 创建权限
```http
POST /api/v1/admin/permissions
```

**请求体**:
```json
{
  "code": "content.write",
  "name": "写入内容",
  "description": "允许创建和编辑内容",
  "category": "content"
}
```

#### 更新权限
```http
PUT /api/v1/admin/permissions/{code}
```

#### 删除权限
```http
DELETE /api/v1/admin/permissions/{code}
```

### 权限模板管理

#### 获取权限模板列表
```http
GET /api/v1/admin/permission-templates
```

**查询参数**:
- `category` (string, 可选): 模板分类
- `is_system` (boolean, 可选): 是否为系统模板

#### 创建权限模板
```http
POST /api/v1/admin/permission-templates
```

**请求体**:
```json
{
  "name": "作者模板",
  "code": "author_template",
  "description": "适用于作者角色的权限模板",
  "permissions": ["content.write", "content.read"],
  "category": "role"
}
```

#### 应用权限模板到角色
```http
POST /api/v1/admin/permission-templates/{templateId}/apply
```

**请求体**:
```json
{
  "roleId": "role123"
}
```

### 审计日志

#### 获取审计日志
```http
GET /api/v1/admin/audit/trail
```

**查询参数**:
- `admin_id` (string, 可选): 管理员ID
- `operation` (string, 可选): 操作类型
- `resource_type` (string, 可选): 资源类型
- `resource_id` (string, 可选): 资源ID
- `start_date` (string, 可选): 开始日期 (YYYY-MM-DD)
- `end_date` (string, 可选): 结束日期 (YYYY-MM-DD)
- `page` (int, 可选): 页码，默认1
- `size` (int, 可选): 每页数量，默认20

### 统计分析

#### 获取用户增长趋势
```http
GET /api/v1/admin/analytics/user-growth
```

**查询参数**:
- `start_date` (string, 必需): 开始日期 (YYYY-MM-DD)
- `end_date` (string, 必需): 结束日期 (YYYY-MM-DD)
- `interval` (string, 必需): 间隔类型 (daily/weekly/monthly)

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "trend": [
      {
        "date": "2026-01-01",
        "count": 150,
        "growth_rate": 5.2
      }
    ]
  }
}
```

#### 获取内容统计
```http
GET /api/v1/admin/analytics/content-statistics
```

**查询参数**:
- `start_date` (string, 可选): 开始日期
- `end_date` (string, 可选): 结束日期

#### 获取收入报告
```http
GET /api/v1/admin/analytics/revenue-report
```

**查询参数**:
- `start_date` (string, 必需): 开始日期
- `end_date` (string, 必需): 结束日期
- `interval` (string, 必需): 间隔类型 (daily/weekly/monthly)

#### 获取活跃用户报告
```http
GET /api/v1/admin/analytics/active-users
```

**查询参数**:
- `start_date` (string, 必需): 开始日期
- `end_date` (string, 必需): 结束日期
- `type` (string, 必需): 报告类型 (dau/wau/mau)

#### 获取系统概览
```http
GET /api/v1/admin/analytics/system-overview
```

#### 导出统计分析报告
```http
GET /api/v1/admin/analytics/export
```

**查询参数**:
- `report_type` (string, 必需): 报告类型 (user-growth/content-statistics/revenue-report/active-users)
- `format` (string, 可选): 导出格式，默认csv
- 其他查询参数根据report_type而定

### 配置管理

#### 获取系统配置
```http
GET /api/v1/admin/config
```

#### 更新系统配置
```http
PUT /api/v1/admin/config
```

**请求体**:
```json
{
  "key": "site.name",
  "value": "青羽写作平台"
}
```

### 公告管理

#### 创建公告
```http
POST /api/v1/admin/announcements
```

**请求体**:
```json
{
  "title": "系统维护通知",
  "content": "系统将于今晚进行维护",
  "type": "maintenance",
  "priority": "high",
  "is_pinned": true,
  "start_time": "2026-01-01T00:00:00Z",
  "end_time": "2026-01-02T00:00:00Z"
}
```

#### 获取公告列表
```http
GET /api/v1/admin/announcements
```

#### 更新公告
```http
PUT /api/v1/admin/announcements/{announcementId}
```

#### 删除公告
```http
DELETE /api/v1/admin/announcements/{announcementId}
```

### Banner管理

#### 创建Banner
```http
POST /api/v1/admin/banners
```

**请求体**:
```json
{
  "title": "新春活动",
  "image_url": "https://example.com/banner.jpg",
  "link_url": "https://example.com/activity",
  "position": "home",
  "start_time": "2026-01-01T00:00:00Z",
  "end_time": "2026-01-31T23:59:59Z",
  "is_active": true
}
```

#### 获取Banner列表
```http
GET /api/v1/admin/banners
```

**查询参数**:
- `position` (string, 可选): Banner位置
- `is_active` (boolean, 可选): 是否激活

### 内容导出

#### 导出书籍数据
```http
GET /api/v1/admin/content/books/export
```

**查询参数**:
- `format` (string, 可选): 导出格式 (csv/excel)，默认csv
- `status` (string, 可选): 状态筛选
- `author` (string, 可选): 作者筛选
- `start_date` (string, 可选): 开始日期
- `end_date` (string, 可选): 结束日期

#### 获取导出历史
```http
GET /api/v1/admin/content/export-history
```

---

## 公告API (api/v1/announcements/)

### 获取公告列表
```http
GET /api/v1/announcements
```

**查询参数**:
- `type` (string, 可选): 公告类型
- `is_pinned` (boolean, 可选): 是否置顶
- `page` (int, 可选): 页码
- `limit` (int, 可选): 每页数量

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "announcements": [
      {
        "id": "anno123",
        "title": "系统维护通知",
        "content": "系统将于今晚进行维护",
        "type": "maintenance",
        "priority": "high",
        "is_pinned": true,
        "created_at": "2026-01-01T00:00:00Z"
      }
    ],
    "total": 10
  }
}
```

### 获取公告详情
```http
GET /api/v1/announcements/{announcementId}
```

---

## 通知API (api/v1/notifications/)

### 获取通知列表
```http
GET /api/v1/notifications
```

**查询参数**:
- `type` (string, 可选): 通知类型
- `is_read` (boolean, 可选): 是否已读
- `page` (int, 可选): 页码
- `limit` (int, 可选): 每页数量

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "notifications": [
      {
        "id": "notif123",
        "title": "新消息提醒",
        "content": "您有一条新消息",
        "type": "message",
        "is_read": false,
        "created_at": "2026-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "unread_count": 10
  }
}
```

### 标记通知为已读
```http
PUT /api/v1/notifications/{notificationId}/read
```

### 标记所有通知为已读
```http
PUT /api/v1/notifications/read-all
```

### 删除通知
```http
DELETE /api/v1/notifications/{notificationId}
```

---

## 社交API (api/v1/social/)

### 关注

#### 关注用户
```http
POST /api/v1/social/follow/{userId}
```

#### 取消关注
```http
DELETE /api/v1/social/follow/{userId}
```

#### 获取关注列表
```http
GET /api/v1/social/following
```

**查询参数**:
- `page` (int, 可选): 页码
- `limit` (int, 可选): 每页数量

#### 获取粉丝列表
```http
GET /api/v1/social/followers
```

### 收藏

#### 添加收藏
```http
POST /api/v1/social/collections
```

**请求体**:
```json
{
  "target_type": "book",
  "target_id": "book123",
  "remark": "推荐阅读"
}
```

#### 获取收藏列表
```http
GET /api/v1/social/collections
```

**查询参数**:
- `target_type` (string, 可选): 目标类型
- `page` (int, 可选): 页码
- `limit` (int, 可选): 每页数量

#### 取消收藏
```http
DELETE /api/v1/social/collections/{collectionId}
```

### 评论

#### 发表评论
```http
POST /api/v1/social/comments
```

**请求体**:
```json
{
  "target_type": "book",
  "target_id": "book123",
  "content": "这本书很棒！",
  "parent_id": null
}
```

#### 获取评论列表
```http
GET /api/v1/social/comments
```

**查询参数**:
- `target_type` (string, 必需): 目标类型
- `target_id` (string, 必需): 目标ID
- `parent_id` (string, 可选): 父评论ID
- `page` (int, 可选): 页码
- `limit` (int, 可选): 每页数量

### 书单

#### 创建书单
```http
POST /api/v1/social/booklists
```

**请求体**:
```json
{
  "title": "我的推荐书单",
  "description": "这些书都很值得阅读",
  "is_public": true,
  "cover_url": "https://example.com/cover.jpg"
}
```

#### 获取书单列表
```http
GET /api/v1/social/booklists
```

**查询参数**:
- `is_public` (boolean, 可选): 是否公开
- `user_id` (string, 可选): 用户ID
- `page` (int, 可选): 页码
- `limit` (int, 可选): 每页数量

#### 添加书籍到书单
```http
POST /api/v1/social/booklists/{booklistId}/books
```

**请求体**:
```json
{
  "book_id": "book123",
  "remark": "强烈推荐"
}
```

---

## 错误代码

| 代码 | 描述 |
|------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用 |

---

## 分页参数

所有支持分页的API都遵循以下规范：

- `page`: 页码，从1开始
- `limit` 或 `size`: 每页数量，默认20，最大100
- 响应中包含 `total`（总数）和 `total_pages`（总页数）

## 速率限制

- 未认证用户: 100请求/小时
- 已认证用户: 1000请求/小时
- 管理员: 无限制

超过限制将返回429状态码。

## 版本管理

API版本通过URL路径指定，当前版本为 `v1`。

## 更多信息

- [使用指南](usage_guide.md)
- [错误处理](error_handling.md)
- [最佳实践](best_practices.md)
