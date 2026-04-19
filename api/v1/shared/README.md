# Shared API 模块说明

## 职责

`api/v1/shared` 是标准共享 HTTP owner，负责放置跨业务域但 owner 明确的 shared API：

- 认证与 OAuth 标准入口
- 钱包 API
- 存储 API
- 通用请求校验、响应封装、认证 helper

它不是“历史兼容层”，也不是新增第二套业务 owner 的落点。

## Auth 边界

- 标准认证/OAuth HTTP owner：`api/v1/shared/auth_api.go`、`api/v1/shared/oauth_api.go`
- 用户域兼容入口：`api/v1/user/handler/auth_handler.go`
- `api/v1/auth/` 历史包已删除，不应恢复

如果后续需要收薄认证入口，必须先统一 `shared` 与 `user` 兼容入口的响应/错误语义，不能机械合并。

当前已确认的关键差异：
- `shared/auth/register` 在邮箱验证码启用时要求 `verification_code`
- `shared/auth/register` 在用户名/邮箱重复时返回明确的冲突错误，而不是内部错误
- `user/auth/register` 仍保留不要求 `verification_code` 的兼容语义
- `shared/auth/login` 走标准错误包装
- `shared/auth/login` 的请求字段校验错误会先返回 400 `请求参数错误: ...`，不会进入 auth service
- `user/auth/login` 对账号不存在/密码错误保留统一的兼容提示
- `shared/auth/logout` 在缺少 Bearer Token 时返回未认证
- `user/auth/logout` 在缺少 Bearer Token 时返回 200 幂等成功

## 主要文件

| 文件 | 作用 |
|------|------|
| `auth_api.go` | 标准认证入口 |
| `oauth_api.go` | 标准 OAuth 入口 |
| `wallet_api.go` | 钱包接口 |
| `storage_api.go` | 文件存储接口 |
| `api_helpers.go` | 通用 API helper |
| `auth_helpers.go` | Bearer Token、登录处理、验证码相关 helper |
| `request_validator.go` | 统一请求校验 |
| `response.go` | 统一响应工具 |

## 当前路由

以下路由由 [`shared_router.go`](/E:/Github/Qingyu/Qingyu_backend/.worktrees/backend-phase4b-auth-user-boundary/router/shared/shared_router.go) 注册。

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/send-verification-code` | 发送注册验证码 |
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/logout` | 用户登出 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |
| GET | `/api/v1/auth/permissions` | 获取当前用户权限 |
| GET | `/api/v1/auth/roles` | 获取当前用户角色 |

### OAuth

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/oauth/:provider/authorize` | 获取 OAuth 授权 URL |
| POST | `/api/v1/oauth/:provider/callback` | 处理 OAuth 回调 |
| GET | `/api/v1/oauth/accounts` | 获取已绑定账号 |
| DELETE | `/api/v1/oauth/accounts/:accountID` | 解绑账号 |
| PUT | `/api/v1/oauth/accounts/:accountID/primary` | 设置主账号 |

### 钱包

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/wallet/balance` | 获取余额 |
| GET | `/api/v1/wallet` | 获取钱包概览 |
| GET | `/api/v1/wallet/transactions` | 获取交易记录 |
| GET | `/api/v1/wallet/withdrawals` | 获取提现记录 |
| POST | `/api/v1/wallet/recharge` | 充值 |
| POST | `/api/v1/wallet/consume` | 消费 |
| POST | `/api/v1/wallet/transfer` | 转账 |
| POST | `/api/v1/wallet/withdraw` | 提现 |

### 存储

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/storage/upload` | 上传文件 |
| GET | `/api/v1/storage/download/:file_id` | 下载文件 |
| DELETE | `/api/v1/storage/files/:file_id` | 删除文件 |
| GET | `/api/v1/storage/files/:file_id` | 获取文件信息 |
| GET | `/api/v1/storage/files` | 获取文件列表 |
| GET | `/api/v1/storage/files/:file_id/url` | 获取下载 URL |

## 约束

- 新的 shared API 只能在 owner 明确时落到本模块，禁止把模块当成“临时收容站”
- handler 必须优先复用 `api_helpers.go` / `auth_helpers.go`
- 保持统一响应语义，禁止直接 `c.JSON()`
