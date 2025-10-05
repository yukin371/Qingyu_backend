# 青羽共享服务 API 使用指南

## 文档概述

本文档提供青羽后端共享服务的完整API使用指南，帮助开发者快速集成和使用各项服务。

**版本**: v1.0  
**基础URL**: `http://localhost:8080/api/v1/shared`  
**文档更新**: 2025-10-04

---

## 目录

- [快速开始](#快速开始)
- [通用规范](#通用规范)
- [认证服务 API](#认证服务-api)
- [钱包服务 API](#钱包服务-api)
- [存储服务 API](#存储服务-api)
- [管理服务 API](#管理服务-api)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)

---

## 快速开始

### 1. 基础配置

```bash
# 克隆项目
git clone https://github.com/yourusername/Qingyu_backend.git
cd Qingyu_backend

# 启动服务
go run main.go
```

### 2. 第一个API请求

```bash
# 注册新用户
curl -X POST http://localhost:8080/api/v1/shared/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123456"
  }'
```

### 3. 认证流程

```bash
# 1. 登录获取Token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/shared/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"Test123456"}' \
  | jq -r '.data.token')

# 2. 使用Token访问受保护的API
curl -X GET http://localhost:8080/api/v1/shared/wallet/balance \
  -H "Authorization: Bearer $TOKEN"
```

---

## 通用规范

### 请求格式

#### 请求头
```
Content-Type: application/json
Authorization: Bearer <token>  # 需要认证的接口
X-Request-ID: <request-id>     # 可选，用于追踪
```

#### 请求体
所有POST/PUT请求使用JSON格式：
```json
{
  "field1": "value1",
  "field2": 123
}
```

### 响应格式

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 业务数据
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "abc123"
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "字段username为必填项",
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "abc123"
}
```

#### 分页响应
```json
{
  "code": 200,
  "message": "success",
  "data": [
    // 数据列表
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_pages": 10
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "abc123"
}
```

### 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或Token无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁（限流） |
| 500 | 服务器内部错误 |

### 速率限制

不同服务有不同的限流策略：

| 服务 | 限流策略 | 说明 |
|------|---------|------|
| 认证（公开） | 10次/分钟 | 防止暴力破解 |
| 认证（已登录） | 30次/分钟 | 正常使用 |
| 钱包 | 50次/分钟 | 支持频繁查询 |
| 存储 | 20次/分钟 | 文件操作限制 |
| 管理 | 100次/分钟 | 管理员权限 |

**超出限流响应**：
```json
{
  "code": 429,
  "message": "请求过于频繁，每60秒最多10次请求"
}
```

---

## 认证服务 API

### 1. 用户注册

**端点**: `POST /auth/register`  
**认证**: 不需要  
**限流**: 10次/分钟

#### 请求参数
```json
{
  "username": "testuser",      // 必填，3-20字符，字母数字下划线
  "email": "test@example.com", // 必填，有效邮箱格式
  "password": "Test123456"     // 必填，8-20字符，需包含大小写字母和数字
}
```

#### 参数说明
| 参数 | 类型 | 必填 | 验证规则 |
|------|------|------|---------|
| username | string | 是 | valid_username（3-20字符） |
| email | string | 是 | valid_email（邮箱格式） |
| password | string | 是 | valid_password（8-20字符，含大小写字母数字） |

#### 响应示例
```json
{
  "code": 201,
  "message": "用户注册成功",
  "data": {
    "user_id": "user_123456",
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2025-10-04T12:00:00Z"
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

#### 错误示例
```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "用户名必须为3-20个字符，只能包含字母、数字和下划线",
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

#### cURL示例
```bash
curl -X POST http://localhost:8080/api/v1/shared/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123456"
  }'
```

---

### 2. 用户登录

**端点**: `POST /auth/login`  
**认证**: 不需要  
**限流**: 10次/分钟

#### 请求参数
```json
{
  "username": "testuser",    // 必填
  "password": "Test123456"   // 必填
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user_123456",
      "username": "testuser",
      "email": "test@example.com",
      "role": "user"
    },
    "expires_at": "2025-10-05T12:00:00Z"
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

#### cURL示例
```bash
curl -X POST http://localhost:8080/api/v1/shared/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456"
  }'
```

---

### 3. 用户登出

**端点**: `POST /auth/logout`  
**认证**: 需要（JWT Token）  
**限流**: 30次/分钟

#### 请求头
```
Authorization: Bearer <token>
```

#### 响应示例
```json
{
  "code": 200,
  "message": "登出成功",
  "data": null,
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

#### cURL示例
```bash
curl -X POST http://localhost:8080/api/v1/shared/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

---

### 4. 刷新Token

**端点**: `POST /auth/refresh`  
**认证**: 需要（JWT Token）  
**限流**: 30次/分钟

#### 请求参数
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."  // 必填，旧Token
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "Token刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-10-05T12:00:00Z"
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

---

### 5. 获取用户权限

**端点**: `GET /auth/permissions`  
**认证**: 需要（JWT Token）  
**限流**: 30次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "permissions": [
      "user:read",
      "user:update",
      "wallet:read",
      "wallet:write"
    ]
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

---

### 6. 获取用户角色

**端点**: `GET /auth/roles`  
**认证**: 需要（JWT Token）  
**限流**: 30次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "roles": ["user"]
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

---

## 钱包服务 API

### 1. 查询余额

**端点**: `GET /wallet/balance`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "balance": 1000.50
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

#### cURL示例
```bash
curl -X GET http://localhost:8080/api/v1/shared/wallet/balance \
  -H "Authorization: Bearer $TOKEN"
```

---

### 2. 获取钱包信息

**端点**: `GET /wallet`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "wallet_123456",
    "user_id": "user_123456",
    "balance": 1000.50,
    "frozen_amount": 0,
    "status": "active",
    "created_at": "2025-10-01T10:00:00Z",
    "updated_at": "2025-10-04T12:00:00Z"
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

---

### 3. 充值

**端点**: `POST /wallet/recharge`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 请求参数
```json
{
  "amount": 100.00,           // 必填，正数，最多2位小数
  "payment_method": "alipay", // 必填，支付方式：alipay/wechat/bank
  "description": "充值100元"  // 可选，描述信息
}
```

#### 参数说明
| 参数 | 类型 | 必填 | 验证规则 |
|------|------|------|---------|
| amount | float64 | 是 | positive_amount（正数，最多2位小数） |
| payment_method | string | 是 | 支付方式枚举 |
| description | string | 否 | 最大200字符 |

#### 响应示例
```json
{
  "code": 200,
  "message": "充值成功",
  "data": {
    "transaction_id": "tx_123456",
    "amount": 100.00,
    "balance": 1100.50,
    "created_at": "2025-10-04T12:00:00Z"
  },
  "timestamp": "2025-10-04T12:00:00Z",
  "request_id": "req_abc123"
}
```

#### cURL示例
```bash
curl -X POST http://localhost:8080/api/v1/shared/wallet/recharge \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "payment_method": "alipay",
    "description": "充值100元"
  }'
```

---

### 4. 消费

**端点**: `POST /wallet/consume`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 请求参数
```json
{
  "amount": 50.00,            // 必填，正数，最多2位小数
  "service_type": "reading",  // 必填，服务类型
  "description": "购买小说章节" // 可选，描述信息
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "消费成功",
  "data": {
    "transaction_id": "tx_123457",
    "amount": 50.00,
    "balance": 1050.50,
    "created_at": "2025-10-04T12:05:00Z"
  },
  "timestamp": "2025-10-04T12:05:00Z",
  "request_id": "req_abc124"
}
```

#### 错误示例（余额不足）
```json
{
  "code": 400,
  "message": "余额不足",
  "error": "当前余额: 10.00，需要: 50.00",
  "timestamp": "2025-10-04T12:05:00Z",
  "request_id": "req_abc124"
}
```

---

### 5. 转账

**端点**: `POST /wallet/transfer`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 请求参数
```json
{
  "to_user_id": "user_789",   // 必填，接收方用户ID
  "amount": 30.00,            // 必填，正数，最多2位小数
  "description": "还款"       // 可选，描述信息
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "转账成功",
  "data": {
    "transaction_id": "tx_123458",
    "from_user_id": "user_123456",
    "to_user_id": "user_789",
    "amount": 30.00,
    "balance": 1020.50,
    "created_at": "2025-10-04T12:10:00Z"
  },
  "timestamp": "2025-10-04T12:10:00Z",
  "request_id": "req_abc125"
}
```

---

### 6. 获取交易记录

**端点**: `GET /wallet/transactions`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 查询参数
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10，最大100 |
| type | string | 否 | 交易类型：recharge/consume/transfer |
| start_time | string | 否 | 开始时间（RFC3339格式） |
| end_time | string | 否 | 结束时间（RFC3339格式） |

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "tx_123458",
      "type": "transfer",
      "amount": 30.00,
      "balance_after": 1020.50,
      "description": "还款",
      "created_at": "2025-10-04T12:10:00Z"
    },
    {
      "id": "tx_123457",
      "type": "consume",
      "amount": 50.00,
      "balance_after": 1050.50,
      "description": "购买小说章节",
      "created_at": "2025-10-04T12:05:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 25,
    "total_pages": 3
  },
  "timestamp": "2025-10-04T12:15:00Z",
  "request_id": "req_abc126"
}
```

#### cURL示例
```bash
# 获取第1页
curl -X GET "http://localhost:8080/api/v1/shared/wallet/transactions?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 按类型筛选
curl -X GET "http://localhost:8080/api/v1/shared/wallet/transactions?type=recharge" \
  -H "Authorization: Bearer $TOKEN"
```

---

### 7. 申请提现

**端点**: `POST /wallet/withdraw`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 请求参数
```json
{
  "amount": 500.00,           // 必填，正数，最多2位小数
  "account_type": "alipay",   // 必填，账户类型：alipay/wechat/bank
  "account": "13800138000",   // 必填，提现账号
  "account_name": "张三",      // 必填，账户名称
  "description": "提现500元"   // 可选，描述信息
}
```

#### 参数说明
| 参数 | 类型 | 必填 | 验证规则 |
|------|------|------|---------|
| amount | float64 | 是 | positive_amount（正数） |
| account_type | string | 是 | 账户类型枚举 |
| account | string | 是 | withdraw_account（根据类型验证） |
| account_name | string | 是 | 2-20字符 |

#### 响应示例
```json
{
  "code": 200,
  "message": "提现申请已提交",
  "data": {
    "request_id": "wd_123456",
    "amount": 500.00,
    "status": "pending",
    "created_at": "2025-10-04T12:20:00Z"
  },
  "timestamp": "2025-10-04T12:20:00Z",
  "request_id": "req_abc127"
}
```

---

### 8. 获取提现记录

**端点**: `GET /wallet/withdrawals`  
**认证**: 需要（JWT Token）  
**限流**: 50次/分钟

#### 查询参数
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| status | string | 否 | 状态：pending/approved/rejected |

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "wd_123456",
      "amount": 500.00,
      "status": "pending",
      "account_type": "alipay",
      "created_at": "2025-10-04T12:20:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1,
    "total_pages": 1
  },
  "timestamp": "2025-10-04T12:25:00Z",
  "request_id": "req_abc128"
}
```

---

## 存储服务 API

### 1. 上传文件

**端点**: `POST /storage/upload`  
**认证**: 需要（JWT Token）  
**限流**: 20次/分钟

#### 请求参数（multipart/form-data）
```
file: <binary>              // 必填，文件内容
category: "avatar"          // 必填，文件分类：avatar/document/image/video
description: "用户头像"     // 可选，文件描述
```

#### 响应示例
```json
{
  "code": 200,
  "message": "文件上传成功",
  "data": {
    "file_id": "file_123456",
    "filename": "avatar.jpg",
    "size": 102400,
    "url": "http://localhost:8080/api/v1/shared/storage/download/file_123456",
    "uploaded_at": "2025-10-04T12:30:00Z"
  },
  "timestamp": "2025-10-04T12:30:00Z",
  "request_id": "req_abc129"
}
```

#### cURL示例
```bash
curl -X POST http://localhost:8080/api/v1/shared/storage/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/avatar.jpg" \
  -F "category=avatar" \
  -F "description=用户头像"
```

---

### 2. 下载文件

**端点**: `GET /storage/download/:file_id`  
**认证**: 需要（JWT Token）  
**限流**: 20次/分钟

#### 响应
- 成功：返回文件二进制数据
- 失败：返回JSON错误信息

#### cURL示例
```bash
# 下载文件
curl -X GET http://localhost:8080/api/v1/shared/storage/download/file_123456 \
  -H "Authorization: Bearer $TOKEN" \
  -o avatar.jpg
```

---

### 3. 删除文件

**端点**: `DELETE /storage/files/:file_id`  
**认证**: 需要（JWT Token）  
**限流**: 20次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "文件删除成功",
  "data": null,
  "timestamp": "2025-10-04T12:35:00Z",
  "request_id": "req_abc130"
}
```

---

### 4. 获取文件信息

**端点**: `GET /storage/files/:file_id`  
**认证**: 需要（JWT Token）  
**限流**: 20次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "file_123456",
    "filename": "avatar.jpg",
    "size": 102400,
    "category": "avatar",
    "owner_id": "user_123456",
    "description": "用户头像",
    "uploaded_at": "2025-10-04T12:30:00Z"
  },
  "timestamp": "2025-10-04T12:40:00Z",
  "request_id": "req_abc131"
}
```

---

### 5. 列出文件

**端点**: `GET /storage/files`  
**认证**: 需要（JWT Token）  
**限流**: 20次/分钟

#### 查询参数
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| category | string | 否 | 文件分类筛选 |

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "file_123456",
      "filename": "avatar.jpg",
      "size": 102400,
      "category": "avatar",
      "uploaded_at": "2025-10-04T12:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 5,
    "total_pages": 1
  },
  "timestamp": "2025-10-04T12:45:00Z",
  "request_id": "req_abc132"
}
```

---

### 6. 获取文件访问URL

**端点**: `GET /storage/files/:file_id/url`  
**认证**: 需要（JWT Token）  
**限流**: 20次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "url": "http://localhost:8080/api/v1/shared/storage/download/file_123456",
    "expires_at": "2025-10-04T13:00:00Z"
  },
  "timestamp": "2025-10-04T12:50:00Z",
  "request_id": "req_abc133"
}
```

---

## 管理服务 API

### 1. 获取待审核内容

**端点**: `GET /admin/reviews/pending`  
**认证**: 需要（JWT Token + 管理员权限）  
**限流**: 100次/分钟

#### 查询参数
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content_type | string | 否 | 内容类型：article/comment/novel |

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "content_id": "content_123",
      "content_type": "article",
      "status": "pending",
      "submitted_at": "2025-10-04T12:00:00Z"
    }
  ],
  "timestamp": "2025-10-04T12:55:00Z",
  "request_id": "req_abc134"
}
```

---

### 2. 审核内容

**端点**: `POST /admin/reviews`  
**认证**: 需要（JWT Token + 管理员权限）  
**限流**: 100次/分钟

#### 请求参数
```json
{
  "content_id": "content_123",  // 必填，内容ID
  "content_type": "article",    // 必填，内容类型
  "action": "approve",          // 必填，审核动作：approve/reject
  "reason": "内容违规"          // 拒绝时必填，原因说明
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "审核完成",
  "data": null,
  "timestamp": "2025-10-04T13:00:00Z",
  "request_id": "req_abc135"
}
```

---

### 3. 审核提现

**端点**: `POST /admin/withdraw/review`  
**认证**: 需要（JWT Token + 管理员权限）  
**限流**: 100次/分钟

#### 请求参数
```json
{
  "request_id": "wd_123456",   // 必填，提现申请ID
  "action": "approve",          // 必填，审核动作：approve/reject
  "reason": "银行账号错误"      // 拒绝时必填，原因说明
}
```

#### 响应示例
```json
{
  "code": 200,
  "message": "提现审核完成",
  "data": null,
  "timestamp": "2025-10-04T13:05:00Z",
  "request_id": "req_abc136"
}
```

---

### 4. 获取用户统计

**端点**: `GET /admin/users/:user_id/statistics`  
**认证**: 需要（JWT Token + 管理员权限）  
**限流**: 100次/分钟

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_id": "user_123456",
    "total_transactions": 50,
    "total_amount": 5000.00,
    "total_posts": 100,
    "total_comments": 200
  },
  "timestamp": "2025-10-04T13:10:00Z",
  "request_id": "req_abc137"
}
```

---

### 5. 获取操作日志

**端点**: `GET /admin/operation-logs`  
**认证**: 需要（JWT Token + 管理员权限）  
**限流**: 100次/分钟

#### 查询参数
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| user_id | string | 否 | 用户ID筛选 |
| action | string | 否 | 操作类型筛选 |

#### 响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "log_123",
      "user_id": "user_123456",
      "action": "login",
      "ip": "192.168.1.1",
      "created_at": "2025-10-04T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1000,
    "total_pages": 100
  },
  "timestamp": "2025-10-04T13:15:00Z",
  "request_id": "req_abc138"
}
```

---

## 错误处理

### 常见错误码

| 错误码 | 说明 | 处理建议 |
|--------|------|---------|
| 400 | 请求参数错误 | 检查请求参数格式和内容 |
| 401 | 未认证 | 重新登录获取Token |
| 403 | 权限不足 | 联系管理员获取权限 |
| 404 | 资源不存在 | 检查资源ID是否正确 |
| 429 | 请求过于频繁 | 降低请求频率或等待 |
| 500 | 服务器错误 | 联系技术支持 |

### 错误响应格式

```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "详细错误信息",
  "timestamp": "2025-10-04T13:20:00Z",
  "request_id": "req_abc139"
}
```

### 错误处理最佳实践

```go
// Go示例
resp, err := client.Do(req)
if err != nil {
    log.Printf("请求失败: %v", err)
    return
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    var errorResp ErrorResponse
    json.NewDecoder(resp.Body).Decode(&errorResp)
    log.Printf("API错误[%d]: %s (RequestID: %s)", 
        errorResp.Code, errorResp.Message, errorResp.RequestID)
    return
}
```

```javascript
// JavaScript示例
try {
  const response = await fetch(url, options);
  const data = await response.json();
  
  if (data.code !== 200) {
    console.error(`API错误[${data.code}]: ${data.message}`);
    console.error(`RequestID: ${data.request_id}`);
    return;
  }
  
  // 处理成功响应
  console.log(data.data);
} catch (error) {
  console.error('请求失败:', error);
}
```

---

## 最佳实践

### 1. Token管理

#### Token刷新策略
```javascript
// JavaScript示例
let token = localStorage.getItem('token');
let tokenExpiresAt = localStorage.getItem('tokenExpiresAt');

async function ensureValidToken() {
  const now = new Date().getTime();
  const expiresAt = new Date(tokenExpiresAt).getTime();
  
  // Token即将过期（提前5分钟刷新）
  if (expiresAt - now < 5 * 60 * 1000) {
    const newToken = await refreshToken(token);
    token = newToken.token;
    tokenExpiresAt = newToken.expires_at;
    localStorage.setItem('token', token);
    localStorage.setItem('tokenExpiresAt', tokenExpiresAt);
  }
  
  return token;
}
```

### 2. 请求重试

```go
// Go示例 - 指数退避重试
func retryRequest(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // 429限流错误，需要退避
        if isRateLimitError(err) {
            waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
            time.Sleep(waitTime)
            continue
        }
        
        // 其他错误，直接返回
        return err
    }
    return fmt.Errorf("超过最大重试次数")
}
```

### 3. 请求ID追踪

```bash
# 在日志中记录RequestID
curl -X GET http://localhost:8080/api/v1/shared/wallet/balance \
  -H "Authorization: Bearer $TOKEN" \
  -v 2>&1 | grep -i "x-request-id"
```

### 4. 分页处理

```go
// Go示例 - 获取所有交易记录
func getAllTransactions() ([]Transaction, error) {
    var allTransactions []Transaction
    page := 1
    
    for {
        resp, err := getTransactions(page, 100)
        if err != nil {
            return nil, err
        }
        
        allTransactions = append(allTransactions, resp.Data...)
        
        // 没有更多数据
        if page >= resp.Pagination.TotalPages {
            break
        }
        
        page++
    }
    
    return allTransactions, nil
}
```

### 5. 并发控制

```go
// Go示例 - 限制并发数
func uploadFilesWithLimit(files []string, limit int) error {
    sem := make(chan struct{}, limit)
    errChan := make(chan error, len(files))
    
    var wg sync.WaitGroup
    for _, file := range files {
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            sem <- struct{}{} // 获取信号量
            defer func() { <-sem }() // 释放信号量
            
            if err := uploadFile(f); err != nil {
                errChan <- err
            }
        }(file)
    }
    
    wg.Wait()
    close(errChan)
    
    // 检查错误
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 6. 性能优化

#### 使用连接池
```go
// Go示例
var httpClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

#### 批量操作
```go
// 不推荐：逐个查询
for _, userID := range userIDs {
    balance, _ := getBalance(userID)
    // 处理balance
}

// 推荐：批量查询（如果API支持）
balances, _ := getBatchBalances(userIDs)
```

---

## 附录

### SDK示例

#### Go SDK
```go
package main

import (
    "fmt"
    "github.com/yourusername/qingyu-go-sdk"
)

func main() {
    // 创建客户端
    client := qingyu.NewClient("http://localhost:8080")
    
    // 登录
    resp, err := client.Auth.Login("testuser", "Test123456")
    if err != nil {
        panic(err)
    }
    
    // 设置Token
    client.SetToken(resp.Token)
    
    // 获取余额
    balance, err := client.Wallet.GetBalance()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("当前余额: %.2f\n", balance)
}
```

#### JavaScript SDK
```javascript
import QingyuClient from 'qingyu-js-sdk';

const client = new QingyuClient('http://localhost:8080');

// 登录
const loginResp = await client.auth.login('testuser', 'Test123456');
client.setToken(loginResp.token);

// 获取余额
const balance = await client.wallet.getBalance();
console.log(`当前余额: ${balance}`);
```

### 常见问题

#### Q: Token过期如何处理？
A: Token过期会返回401错误，此时需要重新登录或使用刷新Token接口获取新Token。

#### Q: 如何提高限流阈值？
A: 普通用户限流阈值是固定的，如需更高阈值请升级为VIP用户或联系管理员。

#### Q: 文件上传大小限制是多少？
A: 默认限制为10MB，如需上传更大文件请使用分片上传或联系管理员。

#### Q: 如何处理429限流错误？
A: 实现指数退避重试策略，或降低请求频率。

#### Q: RequestID的作用是什么？
A: RequestID用于追踪请求链路，遇到问题时提供RequestID可以帮助快速定位。

---

## 联系支持

- **技术文档**: https://docs.qingyu.com
- **GitHub**: https://github.com/yourusername/Qingyu_backend
- **问题反馈**: support@qingyu.com
- **技术社区**: https://community.qingyu.com

---

**文档版本**: v1.0  
**最后更新**: 2025-10-04  
**维护团队**: 青羽后端团队

