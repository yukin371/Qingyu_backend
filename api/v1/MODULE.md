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
