# 共享服务API接口文档

## 概述

本文档描述了青羽后端共享服务的RESTful API接口。共享服务包括：
- **认证服务**：用户注册、登录、权限管理
- **钱包服务**：余额查询、充值、消费、转账、提现
- **存储服务**：文件上传、下载、管理
- **管理服务**：内容审核、用户管理、操作日志

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1/shared`
- **认证方式**: JWT Token（Bearer Token）
- **Content-Type**: `application/json`

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {...}
}
```

### 分页响应
```json
{
  "code": 200,
  "message": "查询成功",
  "data": [...],
  "total": 100,
  "page": 1,
  "size": 20
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误描述",
  "error": "详细错误信息"
}
```

## 认证服务API

### 1. 用户注册
**POST** `/auth/register`

#### 请求体
```json
{
  "username": "user123",
  "email": "user@example.com",
  "password": "password123"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "user_id": "user_001",
    "username": "user123",
    "token": "eyJhbGc..."
  }
}
```

### 2. 用户登录
**POST** `/auth/login`

#### 请求体
```json
{
  "username": "user123",
  "password": "password123"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGc...",
    "expires_in": 3600,
    "user_id": "user_001"
  }
}
```

### 3. 用户登出
**POST** `/auth/logout`

**需要认证**

#### Headers
```
Authorization: Bearer {token}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "登出成功"
}
```

### 4. 刷新Token
**POST** `/auth/refresh`

**需要认证**

#### 响应示例
```json
{
  "code": 200,
  "message": "Token刷新成功",
  "data": {
    "token": "eyJhbGc..."
  }
}
```

### 5. 获取用户权限
**GET** `/auth/permissions`

**需要认证**

#### 响应示例
```json
{
  "code": 200,
  "message": "获取权限成功",
  "data": ["read:books", "write:comments"]
}
```

### 6. 获取用户角色
**GET** `/auth/roles`

**需要认证**

#### 响应示例
```json
{
  "code": 200,
  "message": "获取角色成功",
  "data": ["user", "author"]
}
```

---

## 钱包服务API

### 1. 查询余额
**GET** `/wallet/balance`

**需要认证**

#### 响应示例
```json
{
  "code": 200,
  "message": "查询余额成功",
  "data": 1250.50
}
```

### 2. 获取钱包信息
**GET** `/wallet`

**需要认证**

#### 响应示例
```json
{
  "code": 200,
  "message": "获取钱包信息成功",
  "data": {
    "id": "wallet_001",
    "user_id": "user_001",
    "balance": 1250.50,
    "frozen_balance": 50.00,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-10T12:00:00Z"
  }
}
```

### 3. 充值
**POST** `/wallet/recharge`

**需要认证**

#### 请求体
```json
{
  "amount": 100.00,
  "method": "alipay"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "充值成功",
  "data": {
    "id": "txn_001",
    "type": "recharge",
    "amount": 100.00,
    "balance_after": 1350.50,
    "created_at": "2025-01-10T12:00:00Z"
  }
}
```

### 4. 消费
**POST** `/wallet/consume`

**需要认证**

#### 请求体
```json
{
  "amount": 50.00,
  "reason": "购买书籍"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "消费成功",
  "data": {
    "id": "txn_002",
    "type": "consume",
    "amount": 50.00,
    "balance_after": 1300.50,
    "created_at": "2025-01-10T12:05:00Z"
  }
}
```

### 5. 转账
**POST** `/wallet/transfer`

**需要认证**

#### 请求体
```json
{
  "to_user_id": "user_002",
  "amount": 100.00,
  "reason": "打赏"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "转账成功",
  "data": {
    "id": "txn_003",
    "type": "transfer",
    "amount": 100.00,
    "balance_after": 1200.50,
    "created_at": "2025-01-10T12:10:00Z"
  }
}
```

### 6. 查询交易记录
**GET** `/wallet/transactions?page=1&page_size=20&type=recharge`

**需要认证**

#### 查询参数
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：20）
- `type`: 交易类型（可选：recharge, consume, transfer）

#### 响应示例
```json
{
  "code": 200,
  "message": "查询交易记录成功",
  "data": [
    {
      "id": "txn_001",
      "type": "recharge",
      "amount": 100.00,
      "balance_after": 1350.50,
      "created_at": "2025-01-10T12:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "size": 20
}
```

### 7. 申请提现
**POST** `/wallet/withdraw`

**需要认证**

#### 请求体
```json
{
  "amount": 500.00,
  "account": "alipay:user@example.com"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "申请提现成功",
  "data": {
    "id": "wd_001",
    "amount": 500.00,
    "status": "pending",
    "created_at": "2025-01-10T12:15:00Z"
  }
}
```

### 8. 查询提现申请
**GET** `/wallet/withdrawals?page=1&page_size=20&status=pending`

**需要认证**

#### 查询参数
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：20）
- `status`: 状态（可选：pending, approved, rejected）

#### 响应示例
```json
{
  "code": 200,
  "message": "查询提现申请成功",
  "data": [
    {
      "id": "wd_001",
      "amount": 500.00,
      "status": "pending",
      "created_at": "2025-01-10T12:15:00Z"
    }
  ],
  "total": 5,
  "page": 1,
  "size": 20
}
```

---

## 存储服务API

### 1. 上传文件
**POST** `/storage/upload`

**需要认证**

**Content-Type**: `multipart/form-data`

#### 表单参数
- `file`: 文件（必需）
- `path`: 存储路径（可选）

#### 响应示例
```json
{
  "code": 200,
  "message": "上传文件成功",
  "data": {
    "id": "file_001",
    "filename": "document.pdf",
    "size": 1024000,
    "path": "uploads/user_001/document.pdf",
    "url": "https://cdn.example.com/file_001",
    "created_at": "2025-01-10T12:20:00Z"
  }
}
```

### 2. 下载文件
**GET** `/storage/download/{file_id}`

**需要认证**

#### 响应
文件二进制流

### 3. 删除文件
**DELETE** `/storage/files/{file_id}`

**需要认证**

#### 响应示例
```json
{
  "code": 200,
  "message": "删除文件成功"
}
```

### 4. 获取文件信息
**GET** `/storage/files/{file_id}`

**需要认证**

#### 响应示例
```json
{
  "code": 200,
  "message": "获取文件信息成功",
  "data": {
    "id": "file_001",
    "filename": "document.pdf",
    "size": 1024000,
    "content_type": "application/pdf",
    "user_id": "user_001",
    "created_at": "2025-01-10T12:20:00Z"
  }
}
```

### 5. 列出文件
**GET** `/storage/files?page=1&page_size=20&category=avatar`

**需要认证**

#### 查询参数
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：20）
- `category`: 分类（可选）

#### 响应示例
```json
{
  "code": 200,
  "message": "列出文件成功",
  "data": [
    {
      "id": "file_001",
      "filename": "avatar.jpg",
      "size": 102400,
      "category": "avatar",
      "created_at": "2025-01-10T12:20:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "size": 20
}
```

### 6. 获取文件访问URL
**GET** `/storage/files/{file_id}/url?expire=3600`

**需要认证**

#### 查询参数
- `expire`: 过期时间（秒，默认：3600）

#### 响应示例
```json
{
  "code": 200,
  "message": "获取文件URL成功",
  "data": {
    "url": "https://cdn.example.com/file_001?token=xxx&expires=..."
  }
}
```

---

## 管理服务API

### 1. 获取待审核内容
**GET** `/admin/reviews/pending?content_type=book`

**需要认证（管理员）**

#### 查询参数
- `content_type`: 内容类型（可选：book, chapter, comment）

#### 响应示例
```json
{
  "code": 200,
  "message": "获取待审核内容成功",
  "data": [
    {
      "id": "audit_001",
      "content_id": "book_001",
      "content_type": "book",
      "status": "pending",
      "created_at": "2025-01-10T10:00:00Z"
    }
  ]
}
```

### 2. 审核内容
**POST** `/admin/reviews`

**需要认证（管理员）**

#### 请求体
```json
{
  "content_id": "book_001",
  "content_type": "book",
  "action": "approve",
  "reason": "符合规范"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "审核成功"
}
```

### 3. 审核提现
**POST** `/admin/withdraw/review`

**需要认证（管理员）**

#### 请求体
```json
{
  "withdraw_id": "wd_001",
  "approved": true,
  "reason": "审核通过"
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "批准提现成功"
}
```

### 4. 获取用户统计
**GET** `/admin/users/{user_id}/statistics`

**需要认证（管理员）**

#### 响应示例
```json
{
  "code": 200,
  "message": "获取用户统计成功",
  "data": {
    "user_id": "user_001",
    "total_books": 10,
    "total_chapters": 100,
    "total_words": 500000,
    "total_reads": 10000,
    "total_income": 5000.00
  }
}
```

### 5. 获取操作日志
**GET** `/admin/operation-logs?page=1&page_size=20`

**需要认证（管理员）**

#### 查询参数
- `page`: 页码（默认：1）
- `page_size`: 每页数量（默认：20）
- `admin_id`: 管理员ID（可选）
- `operation`: 操作类型（可选）

#### 响应示例
```json
{
  "code": 200,
  "message": "获取操作日志成功",
  "data": [
    {
      "id": "log_001",
      "admin_id": "admin_001",
      "operation": "approve_book",
      "target": "book_001",
      "created_at": "2025-01-10T11:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "size": 20
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200    | 成功 |
| 400    | 请求参数错误 |
| 401    | 未认证或Token无效 |
| 403    | 无权限访问 |
| 404    | 资源不存在 |
| 500    | 服务器内部错误 |

## 认证说明

大多数API需要JWT认证。请在请求头中携带Token：

```
Authorization: Bearer {your_jwt_token}
```

Token可以通过登录接口获取，有效期为1小时。Token过期后可以通过刷新接口获取新Token。

## 速率限制

为保护服务稳定性，API实施了速率限制：
- 普通用户：每分钟100次请求
- 认证用户：每分钟200次请求
- 管理员：每分钟500次请求

超出限制将返回`429 Too Many Requests`错误。

