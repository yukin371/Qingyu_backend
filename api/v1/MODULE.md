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
- **shared/ 公共层**：`api/v1/shared/` 包含通用的请求验证、响应构建、认证处理，新 API 模块应复用而非重写
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
| 路径参数 | `c.Param` + 空值校验 | `shared.GetRequiredParam(c, key, name)` |
| 分页参数 | 手动 `strconv.Atoi` | `shared.GetPaginationParamsStandard(c)` |
| 传递 userID 到 service | `context.WithValue(ctx, "userId", ...)` | `shared.AddUserIDToContext(c)` |

### Context Key 统一

- **gin.Context 层**：`"user_id"`（由 JWT 中间件设置）、`"username"`、`"roles"`
- **context.Context 层**（传给 service）：`"userId"`
- 禁止使用 `"userID"`、`"userId"` 在 gin.Context 中，或 `"user_id"` 在 context.Context 中
