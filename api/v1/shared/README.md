# Shared API 模块结构说明

## 📁 文件结构

```
api/v1/shared/
├── auth_api.go           # 认证API
├── wallet_api.go         # 钱包API
├── storage_api.go        # 存储API
├── response.go           # 统一响应格式
├── request_validator.go  # 请求验证器
├── types.go              # 公共类型定义
└── README.md             # 本文件
```

## 🎯 模块职责

**职责**: 提供跨业务域的通用功能和工具

**核心功能**:
- ✅ 统一响应格式
- ✅ 请求验证工具
- ✅ 认证相关功能
- ✅ 钱包和交易
- ✅ 文件存储

**注意**: 管理员功能已迁移到 `admin` 模块。

---

## 📋 API端点列表

### 认证API（AuthAPI）

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| POST | /api/v1/auth/refresh-token | 刷新Token | AuthAPI.RefreshToken |
| POST | /api/v1/auth/logout | 登出 | AuthAPI.Logout |
| GET | /api/v1/auth/verify | 验证Token | AuthAPI.VerifyToken |

### 钱包API（WalletAPI）

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| GET | /api/v1/wallet/balance | 获取余额 | WalletAPI.GetBalance |
| GET | /api/v1/wallet/transactions | 交易记录 | WalletAPI.GetTransactions |
| POST | /api/v1/wallet/recharge | 充值 | WalletAPI.Recharge |
| POST | /api/v1/wallet/withdraw | 提现 | WalletAPI.Withdraw |
| GET | /api/v1/wallet/income | 收入统计 | WalletAPI.GetIncome |
| GET | /api/v1/wallet/expense | 支出统计 | WalletAPI.GetExpense |

### 存储API（StorageAPI）

| 方法 | 路径 | 说明 | Handler |
|------|------|------|---------|
| POST | /api/v1/storage/upload | 上传文件 | StorageAPI.Upload |
| DELETE | /api/v1/storage/:fileId | 删除文件 | StorageAPI.Delete |
| GET | /api/v1/storage/:fileId/url | 获取文件URL | StorageAPI.GetFileURL |
| POST | /api/v1/storage/batch-upload | 批量上传 | StorageAPI.BatchUpload |

---

## 🔧 通用工具

### 1. 统一响应格式 (`response.go`)

#### Success - 成功响应
```go
func Success(c *gin.Context, code int, message string, data interface{})
```

**示例**:
```go
shared.Success(c, http.StatusOK, "操作成功", result)
```

**响应格式**:
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {...}
}
```

#### Error - 错误响应
```go
func Error(c *gin.Context, code int, message string, detail string)
```

**示例**:
```go
shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
```

**响应格式**:
```json
{
  "code": 400,
  "message": "参数错误",
  "error": "详细错误信息"
}
```

#### Paginated - 分页响应
```go
func Paginated(c *gin.Context, data interface{}, total int64, page int, pageSize int, message string)
```

**示例**:
```go
shared.Paginated(c, users, total, page, pageSize, "获取成功")
```

**响应格式**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

#### 便捷方法

```go
// 200 OK
shared.OK(c, data)

// 201 Created
shared.Created(c, data)

// 400 Bad Request
shared.BadRequest(c, message, detail)

// 401 Unauthorized
shared.Unauthorized(c, message)

// 403 Forbidden
shared.Forbidden(c, message)

// 404 Not Found
shared.NotFound(c, message)

// 500 Internal Server Error
shared.InternalError(c, message, err)
```

---

### 2. 请求验证 (`request_validator.go`)

#### ValidateRequest - 验证请求体
```go
func ValidateRequest(c *gin.Context, req interface{}) bool
```

**使用**:
```go
var req CreateUserRequest
if !shared.ValidateRequest(c, &req) {
    return  // 验证失败会自动返回错误响应
}
```

#### ValidateQueryParams - 验证查询参数
```go
func ValidateQueryParams(c *gin.Context, req interface{}) bool
```

**使用**:
```go
var req ListUsersRequest
if !shared.ValidateQueryParams(c, &req) {
    return
}
```

---

### 3. 公共类型 (`types.go`)

#### APIResponse - 基础响应
```go
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

#### ErrorResponse - 错误响应
```go
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}
```

#### PaginatedResponse - 分页响应
```go
type PaginatedResponse struct {
    Code       int         `json:"code"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}
```

#### Pagination - 分页信息
```go
type Pagination struct {
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}
```

---

## 📊 请求/响应示例

### 刷新Token
```json
POST /api/v1/auth/refresh-token
Authorization: Bearer <old_token>
Content-Type: application/json

{
  "refresh_token": "refresh_token_string"
}

Response:
{
  "code": 200,
  "message": "刷新成功",
  "data": {
    "access_token": "new_access_token",
    "refresh_token": "new_refresh_token",
    "expires_in": 3600
  }
}
```

### 获取钱包余额
```json
GET /api/v1/wallet/balance
Authorization: Bearer <token>

Response:
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "balance": 1000.50,
    "currency": "CNY",
    "frozen": 100.00,
    "available": 900.50
  }
}
```

### 上传文件
```json
POST /api/v1/storage/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <binary_data>
type: "image"
folder: "avatars"

Response:
{
  "code": 201,
  "message": "上传成功",
  "data": {
    "file_id": "file_123",
    "filename": "avatar.jpg",
    "url": "https://cdn.example.com/avatars/avatar.jpg",
    "size": 102400,
    "mime_type": "image/jpeg",
    "uploaded_at": "2025-10-24T10:00:00Z"
  }
}
```

### 查询交易记录
```json
GET /api/v1/wallet/transactions?page=1&page_size=20&type=income
Authorization: Bearer <token>

Response:
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": "txn_123",
      "type": "income",
      "amount": 50.00,
      "balance_after": 1000.50,
      "description": "订阅收入",
      "created_at": "2025-10-24T09:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## 🔧 设计原则

### 1. 统一性
所有API使用统一的响应格式和错误处理。

### 2. 可复用性
提供通用工具函数，减少重复代码。

### 3. 易用性
简洁的API设计，便于其他模块调用。

### 4. 扩展性
易于添加新的通用功能。

---

## 📝 开发规范

### 1. 使用统一响应

**正确**:
```go
// 使用shared包的响应函数
shared.Success(c, http.StatusOK, "操作成功", data)
```

**错误**:
```go
// 不要直接使用gin.JSON
c.JSON(http.StatusOK, gin.H{"data": data})
```

### 2. 错误处理

```go
if err != nil {
    shared.InternalError(c, "操作失败", err)
    return
}
```

### 3. 参数验证

```go
var req CreateRequest
if !shared.ValidateRequest(c, &req) {
    return  // 自动返回400错误
}
```

---

## 🚀 扩展建议

### 未来可添加的功能

1. **通知服务**
   - 站内通知
   - 邮件通知
   - 短信通知
   - 推送通知

2. **搜索服务**
   - 全文搜索
   - 标签搜索
   - 智能搜索建议

3. **缓存服务**
   - 统一缓存接口
   - 多级缓存
   - 缓存预热

4. **日志服务**
   - 操作日志
   - 审计日志
   - 性能日志

5. **配置服务**
   - 动态配置
   - 配置热更新
   - 配置版本管理

---

## 🔄 与其他模块的关系

### Shared模块的作用

```
┌─────────────────────────────────────┐
│           Shared Module             │
│  ┌──────────────────────────────┐  │
│  │  Response Helper             │  │
│  │  Request Validator           │  │
│  │  Common Types                │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
         ↑      ↑      ↑      ↑
         │      │      │      │
    ┌────┴──┐ ┌┴────┐ ┌┴────┐ ┌┴────┐
    │ User  │ │ AI  │ │Writer│ │Admin│
    │ API   │ │ API │ │ API  │ │ API │
    └───────┘ └─────┘ └──────┘ └─────┘
```

所有业务模块都依赖Shared模块提供的通用功能。

---

## 📚 相关文档

- [API设计规范](../../../doc/api/API设计规范.md)
- [错误处理规范](../../../doc/engineering/软件工程规范_v2.0.md)
- [响应格式说明](../../../doc/api/API接口总览.md)

---

## 💡 最佳实践

### 1. 响应格式

**始终使用统一的响应格式**:
```go
// 成功
shared.Success(c, http.StatusOK, "操作成功", data)

// 失败
shared.Error(c, http.StatusBadRequest, "操作失败", err.Error())
```

### 2. 错误处理

**明确的错误消息**:
```go
// Good
shared.BadRequest(c, "用户名不能为空", "")

// Bad
shared.Error(c, 400, "error", "")
```

### 3. 数据验证

**使用binding标签**:
```go
type Request struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
}
```

### 4. HTTP状态码

```
200 OK          - 请求成功
201 Created     - 创建成功
400 Bad Request - 请求参数错误
401 Unauthorized - 未认证
403 Forbidden   - 无权限
404 Not Found   - 资源不存在
500 Internal Server Error - 服务器错误
```

---

**版本**: v1.0  
**创建日期**: 2025-10-24  
**维护者**: Shared模块开发组

