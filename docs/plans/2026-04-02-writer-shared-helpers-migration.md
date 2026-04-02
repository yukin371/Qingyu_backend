# Writer API 共享辅助函数迁移设计

> 创建日期：2026-04-02
> 分支：refactor/writer-shared-helpers
> 状态：实施中

## 背景

writer 模块 19 个 handler 文件中存在 ~400 处重复的 context 处理逻辑，完全不使用已有的 `api/v1/shared/` 包辅助函数。同时还存在 context key 不一致和错误使用的问题。

## 问题清单

### 重复模式统计

| 模式 | 出现次数 | shared 替代函数 |
|---|---|---|
| `c.Get("user_id")` + 类型断言 + 错误响应 | ~80 | `shared.GetUserID(c)` |
| `c.Get("user_id")` + 静默返回空字符串 | ~15 | `shared.GetUserIDOptional(c)` |
| `c.Param("id")` + 空值校验 | ~60 | `shared.GetRequiredParam(c, ...)` |
| `c.ShouldBindJSON` + err 响应 | ~100 | `shared.BindJSON(c, &req)` |
| 手动分页解析 | ~20 | `shared.GetPaginationParamsStandard(c)` |
| `context.WithValue(ctx, "userID", ...)` | ~15 | `shared.AddUserIDToContext(c)` |

### Context Key 错误（Bug）

JWT 中间件设置的 key（`internal/middleware/auth/jwt.go:166-167`）：
- `c.Set("user_id", claims.UserID)`
- `c.Set("username", claims.Username)`
- `c.Set("roles", claims.Roles)`

| 文件 | 错误 key | 正确 key | 影响 |
|---|---|---|---|
| `audit_api.go` (6处) | `"userID"` | `"user_id"` | 永远拿不到用户 |
| `import_export_api.go` (2处) | `"userId"` | `"user_id"` | 永远拿不到用户 |
| `comment_api.go`, `lock_api.go` | `"userName"` | `"username"` | 永远拿不到用户名 |
| `lock_api.go` | `"userRole"` | `"roles"` | 永远拿不到角色 |

### Context Key 不一致（API→Service 传递）

service 层从 `context.Context` 读取 userID 的 key 也不统一：

| Service 文件 | context key |
|---|---|
| `service/writer/document/auth_helper.go` (9处) | `"userID"` |
| `service/writer/document/batch_operation_service.go` (3处) | `"userID"` |
| `service/writer/project/project_service.go` (5处) | `"userId"` |
| `service/search/search.go` (2处) | `"userId"` + `"user_id"` fallback |

**统一方案**：全部统一为 `"userId"`，与 `shared.AddUserIDToContext` 一致。

## 解决方案

### 1. shared 包新增函数

在 `api/v1/shared/api_helpers.go` 中新增：

```go
// GetUserName 获取用户名（可选，不存在返回空字符串）
func GetUserName(c *gin.Context) string

// GetUserRoles 获取用户角色列表（可选，不存在返回 nil）
func GetUserRoles(c *gin.Context) []string
```

### 2. writer handler 迁移（4 批次）

#### 第1批：简单迁移（user_id 提取 → shared.GetUserID）

| 文件 | 迁移点 |
|---|---|
| `character_api.go` | 5处 user_id 提取 |
| `location_api.go` | 4处 user_id 提取 |
| `outline_api.go` | 4处 user_id 提取 |
| `export_api.go` | 5处 user_id 提取 |
| `keyword_api.go` | 少量 |

#### 第2批：context 创建 + 分页迁移

| 文件 | 迁移点 |
|---|---|
| `document_api.go` | 删除本地 `getContextWithUserID`，5处 context 创建迁移 |
| `editor_api.go` | 删除本地 `getAuthUserID`，context 创建迁移 |
| `project_api.go` | 6处 context 创建迁移 |
| `search_api.go` | 分页迁移 |
| `stats_api.go` | 分页迁移 |
| `template_api.go` | 分页迁移 |

#### 第3批：修复错误 context key（Bug 修复）

| 文件 | 修复内容 |
|---|---|
| `audit_api.go` | `"userID"` → `shared.GetUserID(c)` |
| `import_export_api.go` | `"userId"` → `shared.GetUserID(c)` |
| `comment_api.go` | `"userName"` → `shared.GetUserName(c)` |
| `lock_api.go` | `"userName"` → `shared.GetUserName(c)`, `"userRole"` → `shared.GetUserRoles(c)` |

#### 第4批：其余文件

| 文件 | 迁移点 |
|---|---|
| `publish_api.go` | 5处 user_id 提取 |
| `batch_operation_api.go` | 3处 user_id 提取 |
| `timeline_api.go` | 少量 |
| `version_api.go` | 少量 |
| `encyclopedia_api.go` | 检查 |
| `writer_stats_aggregate_api.go` | 1处 |

### 3. service 层 context key 统一

`service/writer/document/auth_helper.go` + `batch_operation_service.go`：
- `"userID"` → `"userId"`（共 12 处）

`service/search/search.go`：
- 删除 `"user_id"` fallback，统一用 `"userId"`

### 4. MODULE.md 更新

在 `api/v1/MODULE.md` 新增"辅助函数使用规范"段落。

## 验收标准

1. 所有 writer handler 使用 shared 包辅助函数，无内联 context 提取
2. 本地 `getContextWithUserID` 和 `getAuthUserID` 函数已删除
3. context key 全部统一（gin: `user_id`/`username`/`roles`，context: `userId`）
4. `go build ./...` 编译通过
5. 现有测试通过
6. MODULE.md 已更新
