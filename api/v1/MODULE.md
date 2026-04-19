# API v1

> 最后更新：2026-03-29

## 职责

HTTP API 处理层，接收请求、参数校验、调用 Service、格式化响应。所有 API 注册在 `/api/v1` 路由组下。不包含业务逻辑。

## 数据流

```
Gin Router → Middleware（Auth/CORS/RateLimit） → Handler → Service → Repository
                                                      ↓
                                              pkg/response（统一响应格式）
```

## 约定 & 陷阱

- **响应格式强制**：必须使用 `pkg/response` 包（`response.Success`/`response.BadRequest` 等），禁止直接 `c.JSON()`
- **4 位错误码**：业务错误码为 4 位数字（1001 参数错误、2001 用户不存在等），禁止用 HTTP 状态码作为业务码
- **前端前缀自动添加**：前端 HTTP 拦截器自动加 `/api/v1`，后端路由必须注册在此前缀下
- **字段名转换**：后端返回 `snake_case`，前端拦截器自动转 `camelCase`，后端无需处理
- **shared/ 公共层**：`api/v1/shared/` 包含通用的请求验证、响应构建、认证处理，新 API 模块应复用而非重写；其中标准认证/OAuth HTTP owner 已收敛到 `api/v1/shared/auth_api.go` 与 `api/v1/shared/oauth_api.go`
- **Swagger 注解**：每个 API 端点必须有 Swagger 注解，用于自动生成文档和 Orval 前端类型

## 辅助函数使用规范（强制）

所有 handler **必须**使用 `api/v1/shared` 包的辅助函数，禁止内联重复逻辑：

| 场景 | 禁止写法 | 必须使用 |
|------|----------|----------|
| 获取用户ID（必需） | `c.Get("user_id")` + 类型断言 + 错误响应 | `shared.GetUserID(c)` |
| 获取用户ID（可选） | `c.Get("user_id")` + 静默返回 | `shared.GetUserIDOptional(c)` |
| 获取用户名 | `c.Get("username")` + 类型断言 | `shared.GetUserName(c)` |
| 获取用户角色 | `c.Get("roles")` + 类型断言 | `shared.GetUserRoles(c)` |
| JSON 绑定 | `c.ShouldBindJSON` + err 响应 | `shared.BindJSON(c, &req)` |
| 兼容 JSON 绑定文案 | 手写 `请求参数错误: ...` | `shared.BindJSONWithMessage(c, &req, prefix)` |
| Bearer Token（必需） | `c.GetHeader("Authorization")` + 手动裁剪 | `shared.GetBearerToken(c)` |
| Bearer Token（可选） | `c.GetHeader("Authorization")` + 静默返回 | `shared.GetBearerTokenOptional(c)` |
| 路径参数 | `c.Param` + 空值校验 | `shared.GetRequiredParam(c, key, name)` |
| 分页参数 | 手动 `strconv.Atoi` | `shared.GetPaginationParamsStandard(c)` |
| 传递 userID 到 service | `context.WithValue(ctx, "userId", ...)` | `shared.AddUserIDToContext(c)` |

## Auth 边界

- 标准认证/OAuth 路由 owner：`api/v1/shared/*`
- 用户域兼容认证入口：`api/v1/user/handler/auth_handler.go`
- `api/v1/auth/` 历史包已退场，不应再新增或恢复第二套 auth API owner

### 当前已固化的兼容差异

- `POST /api/v1/shared/auth/register`：标准入口，启用邮箱验证码后必须提供并通过 `verification_code`
- `POST /api/v1/shared/auth/register`：用户名/邮箱重复时会返回明确的冲突响应，而不是兜底 500
- `POST /api/v1/user/auth/register`：兼容入口，当前仍保持“不要求 `verification_code`”的旧语义
- `POST /api/v1/shared/auth/login`：标准入口，统一返回 `登录失败: ...` 风格的未授权提示
- `POST /api/v1/shared/auth/login`：请求体字段缺失等校验错误会在进入 service 前直接返回 400 `请求参数错误: ...`
- `POST /api/v1/user/auth/login`：兼容入口，对账号不存在/密码错误统一返回 `用户名或密码错误`，并保留旧错误映射
- `POST /api/v1/shared/auth/logout`：缺少 Bearer Token 时返回未认证错误
- `POST /api/v1/user/auth/logout`：缺少 Bearer Token 时保持幂等成功，避免破坏旧前端/旧客户端流程

以上差异已由 `test/api/auth_api_test.go` 与 `api/v1/user/handler/auth_handler_logout_test.go` 覆盖；后续若要继续收薄，必须先统一外部契约，再调整实现。

### Context Key 统一

- **gin.Context 层**：`"user_id"`（由 JWT 中间件设置）、`"username"`、`"roles"`
- **context.Context 层**（传给 service）：`"userId"`
- 禁止使用 `"userID"`、`"userId"` 在 gin.Context 中，或 `"user_id"` 在 context.Context 中
