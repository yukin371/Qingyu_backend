# Shared API 模块结构说明

## 文件结构

```
api/v1/shared/
├── auth_api.go            # 认证服务 API（注册/登录/登出/刷新/权限/角色/验证码）
├── auth_helpers.go        # 认证辅助函数（验证码校验、标准登录处理）
├── oauth_api.go           # OAuth 第三方登录 API
├── wallet_api.go          # 钱包服务 API（余额/充值/消费/转账/提现）
├── storage_api.go         # 文件存储 API（上传/下载/分片/缩略图/权限）
├── api_helpers.go         # 通用 API 辅助函数（参数提取/分页/绑定）
├── request_validator.go   # 统一请求校验器
├── response.go            # 统一响应封装（成功/错误/分页/校验错误）
├── types.go               # 跨模块共享响应类型定义
├── user_types.go          # 用户相关 DTO 与公开信息类型
├── user_converter.go      # User Model 与 DTO 双向转换
├── auth_helpers_test.go   # 认证辅助测试
├── response_test.go       # 响应封装测试
├── storage_api_test.go    # 存储 API 测试
├── api_helpers_test.go    # 通用辅助测试
└── README.md              # 本文件
```

## 模块职责划分

### 1. AuthAPI (`auth_api.go`)

**职责**: 标准认证服务 HTTP 入口，提供用户注册、登录、登出、Token 刷新及权限查询。

**核心功能**:
- ✅ 用户注册（支持邮箱验证码校验）
- ✅ 用户登录（标准错误包装）
- ✅ 用户登出（Bearer Token 缺失时返回未认证）
- ✅ Token 刷新
- ✅ 获取当前用户权限列表
- ✅ 获取当前用户角色列表
- ✅ 发送注册验证码

**API 端点**:
```
POST /api/v1/auth/send-verification-code   # 发送注册验证码
POST /api/v1/auth/register                 # 用户注册
POST /api/v1/auth/login                    # 用户登录
POST /api/v1/auth/logout                   # 用户登出
POST /api/v1/auth/refresh                  # 刷新 Token
GET  /api/v1/auth/permissions              # 获取当前用户权限
GET  /api/v1/auth/roles                    # 获取当前用户角色
```

**依赖服务**:
- `auth.AuthService` — 认证核心服务
- `emailcode.Manager` — 邮箱验证码管理器

**兼容关系**: 与 `api/v1/user/handler/auth_handler.go` 存在兼容路由，差异详见下方"Auth 边界"章节。

---

### 2. OAuthAPI (`oauth_api.go`)

**职责**: 第三方 OAuth 登录与账号绑定，支持 Google/GitHub/QQ 等提供商。

**核心功能**:
- ✅ 获取 OAuth 授权 URL（支持登录模式和绑定模式）
- ✅ 处理 OAuth 回调（自动登录/注册或绑定已有账号）
- ✅ 获取已绑定的第三方账号列表
- ✅ 解绑第三方账号
- ✅ 设置主账号

**API 端点**:
```
POST   /api/v1/oauth/:provider/authorize        # 获取 OAuth 授权 URL
POST   /api/v1/oauth/:provider/callback         # 处理 OAuth 回调
GET    /api/v1/oauth/accounts                   # 获取已绑定账号
DELETE /api/v1/oauth/accounts/:accountID        # 解绑账号
PUT    /api/v1/oauth/accounts/:accountID/primary # 设置主账号
```

**依赖服务**:
- `auth.OAuthServiceInterface` — OAuth 核心服务
- `auth.AuthService` — 认证服务（用于回调时签发 Token）

---

### 3. WalletAPI (`wallet_api.go`)

**职责**: 用户钱包服务，提供余额查询、充值、消费、转账和提现功能。

**核心功能**:
- ✅ 查询钱包余额
- ✅ 获取钱包完整信息
- ✅ 充值
- ✅ 消费
- ✅ 转账
- ✅ 查询交易记录（分页）
- ✅ 申请提现
- ✅ 查询提现记录

**API 端点**:
```
GET  /api/v1/wallet/balance       # 查询余额
GET  /api/v1/wallet               # 获取钱包概览
GET  /api/v1/wallet/transactions  # 获取交易记录
GET  /api/v1/wallet/withdrawals   # 获取提现记录
POST /api/v1/wallet/recharge      # 充值
POST /api/v1/wallet/consume       # 消费
POST /api/v1/wallet/transfer      # 转账
POST /api/v1/wallet/withdraw      # 提现
```

**依赖服务**:
- `wallet.WalletService` — 钱包核心服务

---

### 4. StorageAPI (`storage_api.go`)

**职责**: 文件存储服务，提供文件上传/下载、分片上传、缩略图生成及文件访问权限管理。

**核心功能**:
- ✅ 小文件上传（<50MB，multipart/form-data）
- ✅ 文件下载
- ✅ 获取文件信息
- ✅ 删除文件
- ✅ 文件列表查询（分页）
- ✅ 获取下载 URL（临时签名链接）
- ✅ 分片上传（初始化/上传分片/完成/中止/进度查询）
- ✅ 缩略图生成
- ✅ 文件访问授权（Grant/Revoke）

**API 端点**:
```
# 基础文件操作
POST   /api/v1/storage/upload                # 上传文件
GET    /api/v1/storage/download/:file_id     # 下载文件
GET    /api/v1/storage/files/:file_id        # 获取文件信息
DELETE /api/v1/storage/files/:file_id        # 删除文件
GET    /api/v1/storage/files                 # 获取文件列表
GET    /api/v1/storage/files/:file_id/url    # 获取下载 URL

# 分片上传
POST   /api/v1/storage/multipart/initiate    # 初始化分片上传
PUT    /api/v1/storage/multipart/chunk       # 上传分片
POST   /api/v1/storage/multipart/complete    # 完成分片上传
DELETE /api/v1/storage/multipart/:upload_id  # 中止分片上传
GET    /api/v1/storage/multipart/:upload_id/progress  # 查询上传进度

# 图片与权限
POST /api/v1/storage/files/:file_id/thumbnail  # 生成缩略图
POST /api/v1/storage/files/:file_id/access     # 授权访问
DELETE /api/v1/storage/files/:file_id/access   # 撤销访问
```

**依赖服务**:
- `storage.StorageService` — 存储核心服务
- `storage.MultipartUploadManager` — 分片上传管理
- `storage.ImageProcessorService` — 图片处理（缩略图）

---

### 5. api_helpers (`api_helpers.go`)

**职责**: 通用 API 辅助函数库，为所有 handler 提供参数提取、分页、请求绑定等公共能力。

**核心功能**:
- ✅ 用户 ID 提取（必需/可选/自定义消息）
- ✅ Bearer Token 提取（必需/可选）
- ✅ 路径参数提取（`GetRequiredParam`、`GetFirstParam`）
- ✅ 查询参数提取（`GetRequiredQuery`、`GetIntQueryInRange`）
- ✅ 分页参数解析（标准/大页/小页预设）
- ✅ 请求体绑定（JSON/Query/Params，含错误处理）
- ✅ 分页响应构造（`RespondWithPaginated`）
- ✅ 用户角色/用户名提取
- ✅ 批量 ID 校验

**依赖**:
- `pkg/response` — 统一响应包

---

### 6. auth_helpers (`auth_helpers.go`)

**职责**: 认证相关辅助函数，提供验证码校验、验证码发送、标准登录流程封装。

**核心功能**:
- ✅ 注册邮箱验证码校验（支持验证码功能未启用时自动跳过）
- ✅ 注册邮箱验证码发送
- ✅ 标准登录流程封装（绑定 JSON + 调用 AuthService + 统一错误响应）

**依赖**:
- `pkg/emailcode.Manager` — 邮箱验证码管理器
- `service/auth.AuthService` — 认证服务

---

### 7. request_validator (`request_validator.go`)

**职责**: 统一请求校验器，封装 validator 库并提供字段级错误响应。

**核心功能**:
- ✅ 请求体自动绑定 + 结构体校验（向后兼容旧代码未显式绑定的情况）
- ✅ 查询参数绑定 + 校验
- ✅ 字段级验证错误响应（`ValidationErrorResponse`）
- ✅ 通用验证错误处理（`HandleValidationError`）

**依赖**:
- `pkg/validator` — 全局验证器实例

---

### 8. response (`response.go`)

**职责**: 统一 HTTP 响应封装，为所有 handler 提供一致的响应格式。

**核心功能**:
- ✅ 成功响应（`Success`、`SuccessData`、`SuccessResponse`）
- ✅ 错误响应（`Error`、`ErrorResponseWithCode`）
- ✅ 分页响应（`Paginated`、`PaginatedResponseHelper`）
- ✅ 校验错误响应（`ValidationError`）
- ✅ HTTP 语义响应（`Unauthorized`/`Forbidden`/`NotFound`/`InternalError`/`BadRequest`）
- ✅ 服务层错误统一处理（`HandleServiceError`、`HandleError`）
- ✅ 统一错误封装（`WrapServiceError`）

**响应结构体**:
- `APIResponse` — 标准响应（code/message/data/timestamp/request_id）
- `PaginatedResponse` — 分页响应（含 Pagination 元信息）
- `ErrorResponse` — 错误响应（含 debug 字段，仅开发环境）

---

### 9. types (`types.go`)

**职责**: 跨模块共享的响应类型定义。

**包含类型**:
- `WalletBalanceResponse` — 钱包余额
- `TransactionResponse` — 交易记录
- `WithdrawResponse` — 提现记录
- `LoginResponse` — 登录响应（含 Token 和用户信息）
- `FileUploadResponse` — 文件上传结果
- `StorageQuotaResponse` — 存储配额
- `AdminStatsResponse` — 管理员统计概览

---

### 10. user_types (`user_types.go`)

**职责**: 用户相关 DTO 与公开信息类型定义。

**包含类型**:
- `UserProfileResponse` — 用户完整信息（含敏感字段，用于自身/管理员查看）
- `PublicUserProfileResponse` — 用户公开信息（去敏感，用于公开场景）
- `UserBasicInfo` — 基本用户信息（登录/注册响应、简单引用）
- `UserDTO` — 重新导出 `models/dto.UserDTO`（向后兼容）

---

### 11. user_converter (`user_converter.go`)

**职责**: User Model 与 DTO 的双向转换，处理 ID/时间戳等字段格式映射。

**核心功能**:
- ✅ `ToUserDTO` — Model 转 DTO（用于 API 层返回）
- ✅ `ToUserDTOs` — 批量 Model 转 DTO
- ✅ `ToUserModel` — DTO 转 Model（用于更新，保留 ID）
- ✅ `ToUserModelWithoutID` — DTO 转 Model（用于创建，ID 由数据库生成）

**依赖**:
- `models/shared/types.DTOConverter` — 通用类型转换器

---

## Auth 边界说明

`shared/auth_api.go` 与 `api/v1/user/handler/auth_handler.go` 存在兼容路由关系，关键差异：

| 行为 | shared（标准） | user（兼容） |
|------|---------------|-------------|
| 注册验证码 | 邮箱验证码启用时必填 | 不要求验证码 |
| 用户名/邮箱冲突 | 返回 409 + 明确冲突码 | 返回通用错误 |
| 登录参数错误 | 先返回 400 再进入 service | 直接进入 service |
| 登录账号不存在 | 标准错误包装 | 统一兼容提示 |
| 登出缺少 Token | 返回未认证 | 返回 200 幂等成功 |

> 历史包 `api/v1/auth/` 已删除，不应恢复。如需收薄认证入口，必须先统一两套入口的响应/错误语义。

## 约束

- 新的 shared API 只能在 owner 明确时落到本模块，禁止作为"临时收容站"
- handler 必须优先复用 `api_helpers.go` / `auth_helpers.go`
- 保持统一响应语义，禁止直接 `c.JSON()`
