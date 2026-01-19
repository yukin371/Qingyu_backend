# 管理员API文档

> **版本**: v1.0  
> **创建日期**: 2025-10-21

---

## 1. API概述

管理员API提供后台管理功能，包括用户管理、内容审核、系统配置等。

**Base URL**: `/api/v1/admin`

**权限要求**: 管理员角色

---

## 2. 用户管理

### 2.1 获取用户列表

**接口**: `GET /admin/users`

**参数**：
- `page` - 页码
- `pageSize` - 每页数量
- `role` - 角色筛选
- `status` - 状态筛选

**响应**：
```json
{
  "code": 200,
  "data": {
    "users": [
      {
        "id": "user_123",
        "username": "testuser",
        "email": "test@example.com",
        "role": "user",
        "status": "active",
        "createdAt": "2025-10-20T10:00:00Z"
      }
    ],
    "total": 1000,
    "page": 1,
    "pageSize": 20
  }
}
```

### 2.2 封禁用户

**接口**: `POST /admin/users/:userId/ban`

**请求**：
```json
{
  "reason": "违反社区规则",
  "duration": 7,
  "durationUnit": "days"
}
```

### 2.3 解除封禁

**接口**: `POST /admin/users/:userId/unban`

---

## 3. 内容审核

### 3.1 获取待审核内容

**接口**: `GET /admin/audit/pending`

**参数**：
- `targetType` - document/chapter/comment
- `page` - 页码

**响应**：
```json
{
  "code": 200,
  "data": {
    "records": [
      {
        "id": "audit_123",
        "targetType": "document",
        "targetId": "doc_456",
        "authorId": "user_789",
        "riskLevel": 3,
        "violations": [
          {
            "type": "sensitive_word",
            "content": "敏感词示例"
          }
        ],
        "createdAt": "2025-10-21T09:00:00Z"
      }
    ],
    "total": 50
  }
}
```

### 3.2 审核通过

**接口**: `POST /admin/audit/:auditId/approve`

**请求**：
```json
{
  "note": "内容符合规范"
}
```

### 3.3 审核拒绝

**接口**: `POST /admin/audit/:auditId/reject`

**请求**：
```json
{
  "reason": "包含违规内容",
  "penaltyType": "warning",
  "penaltyDuration": 3
}
```

---

## 4. 敏感词管理

### 4.1 添加敏感词

**接口**: `POST /admin/sensitive-words`

**请求**：
```json
{
  "word": "敏感词",
  "category": "政治",
  "level": 5,
  "replacement": "***"
}
```

### 4.2 删除敏感词

**接口**: `DELETE /admin/sensitive-words/:wordId`

### 4.3 批量导入

**接口**: `POST /admin/sensitive-words/import`

**请求**：
```json
{
  "words": [
    {"word": "词1", "category": "政治", "level": 5},
    {"word": "词2", "category": "色情", "level": 5}
  ]
}
```

---

## 5. 系统配置

### 5.1 获取系统配置

**接口**: `GET /admin/config`

**响应**：
```json
{
  "code": 200,
  "data": {
    "allowRegistration": true,
    "requireEmailVerification": true,
    "maxUploadSize": 10485760,
    "enableAudit": true
  }
}
```

### 5.2 更新系统配置

**接口**: `PUT /admin/config`

**请求**：
```json
{
  "allowRegistration": false,
  "maxUploadSize": 20971520
}
```

---

## 6. 数据统计

### 6.1 获取系统统计

**接口**: `GET /admin/stats`

**响应**：
```json
{
  "code": 200,
  "data": {
    "totalUsers": 10000,
    "activeUsers": 3000,
    "totalBooks": 5000,
    "totalRevenue": 100000.00,
    "pendingAudits": 50
  }
}
```

---

## 7. 公告管理

### 7.1 发布公告

**接口**: `POST /admin/announcements`

**请求**：
```json
{
  "title": "系统升级公告",
  "content": "系统将于今晚维护...",
  "type": "system",
  "priority": "high"
}
```

### 7.2 获取公告列表

**接口**: `GET /admin/announcements`

---

**文档状态**: ✅ 已完成

