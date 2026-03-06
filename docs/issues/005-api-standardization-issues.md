# Issue #005: API 标准化问题

**优先级**: 高 (P0)
**类型**: 架构问题
**状态**: 部分修复
**创建日期**: 2026-03-05
**来源报告**: [后端综合审计报告](../reports/archived/backend-comprehensive-audit-summary-2026-01-26.md)、[后端 API 分析](../reports/archived/backend-api-analysis-2026-01-26.md)

---

## 问题描述

RESTful API 存在多个标准化问题，导致前后端对接困难、API 管理混乱。

### 具体问题

#### 1. 响应码系统严重不一致 🔴 P0

**问题**: API 响应中的 `code` 字段值不一致。

```go
// 有些地方使用 200 表示成功
{"code": 200, "message": "success", "data": {...}}

// 有些地方使用 0 表示成功
{"code": 0, "message": "success", "data": {...}}

// 有些地方使用 1 表示成功
{"code": 1, "message": "success", "data": {...}}
```

**影响**:
- 前端需要处理多种成功码
- 错误判断逻辑复杂
- API 响应中间件难以统一实现

#### 2. URL 前缀不统一 🔴 P0

**问题**: 12 个端点使用 `/system/` 前缀，与其他端点不一致。

```
正常模式: /api/v1/{module}/{action}
异常模式: /system/{action}

异常示例:
- /system/health
- /system/info
- /system/metrics
```

**影响**:
- API 网关配置复杂
- 权限控制不一致
- API 文档混乱

#### 3. 完全没有使用 PATCH 方法 🔴 P0

**问题**: 591 个 API 端点中，PATCH 方法使用数为 0。

```
GET:    298个 (50.4%)
POST:   168个 (28.4%)
PUT:     58个 (9.8%)
DELETE:  55个 (9.3%)
PATCH:    0个 ⚠️
```

**影响**:
- 部分更新需要 PUT 整个资源
- 带宽浪费
- 不符合 RESTful 最佳实践

#### 4. 分页响应格式不符合规范 🟡 P1

**问题**: 分页响应格式不统一。

```go
// 格式 1: 使用 total 和 items
{
  "total": 100,
  "items": [...]
}

// 格式 2: 使用 data 和 pagination
{
  "data": [...],
  "pagination": {
    "total": 100,
    "page": 1,
    "size": 20
  }
}
```

#### 5. 重复的 API 端点 🟡 P1

**问题**: 存在功能重复的 API 端点。

```
示例:
- /api/v1/books/:id
- /api/v1/bookstore/books/:id
- /api/v1/reader/books/:id
```

---

## 解决方案

### 1. 统一响应码规范

```go
// pkg/response/types.go
package response

const (
    CodeSuccess = 0     // 成功统一使用 0
    CodeError   = -1    // 错误统一使用 -1
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) *Response {
    return &Response{
        Code:    CodeSuccess,
        Message: "success",
        Data:    data,
    }
}

func Error(code int, message string) *Response {
    return &Response{
        Code:    CodeError,
        Message: message,
    }
}
```

### 2. 统一 URL 前缀

```go
// 将 /system/* 端点迁移到 /api/v1/system/*
// /system/health  →  /api/v1/system/health
// /system/info    →  /api/v1/system/info
```

### 3. 添加 PATCH 方法支持

```go
// 使用 PATCH 实现部分更新
func (h *BookHandler) PatchBook(c *gin.Context) {
    bookID := c.Param("id")
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        response.Error(c, err)
        return
    }

    // 只更新提供的字段
    if err := h.service.PatchBook(c, bookID, updates); err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, nil)
}
```

### 4. 统一分页响应格式

```go
// pkg/response/pagination.go
type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    Total int64 `json:"total"`
    Page  int   `json:"page"`
    Size  int   `json:"size"`
    Pages int   `json:"pages"`
}
```

---

## 实施计划

### Phase 1: 响应码统一（1-2 周）

1. 创建统一的响应类型定义
2. 更新所有 API Handler 使用新响应类型
3. 更新前端适配新的响应格式
4. 添加响应中间件确保一致性

**需要修改的文件**:
- `pkg/response/` - 新建响应包
- `api/v1/**/*.go` - 所有 API Handler
- 前端 API 调用代码

### Phase 2: URL 前缀统一（2-3 天）

1. 列出所有 `/system/*` 端点
2. 修改路由定义
3. 更新 API 文档
4. 更新前端调用

### Phase 3: 添加 PATCH 支持（1 周）

1. 为需要部分更新的资源添加 PATCH 端点
2. 实现 `PatchXxx` Service 方法
3. 添加测试
4. 更新 API 文档

### Phase 4: 分页格式统一（3-5 天）

1. 定义统一的分页响应格式
2. 更新所有分页 API
3. 更新前端适配

---

## 检查清单

### 响应码统一
- [x] 定义统一的响应常量
- [x] 创建响应工具函数
- [ ] 更新所有 API Handler
- [ ] 前端适配完成
- [ ] 测试验证

### URL 前缀统一
- [ ] 列出所有需要迁移的端点
- [ ] 修改路由定义
- [ ] 更新 API 文档
- [ ] 前端调用更新

### PATCH 方法支持
- [ ] 识别需要部分更新的资源
- [ ] 实现 PATCH Handler
- [ ] 实现 Service 层方法
- [ ] 添加测试

---

## 当前进展（2026-03-06）

已完成一批低风险、可独立合并的 API 标准化修复：

1. `api/v1/admin/analytics_api.go`
   - 从裸 `c.JSON` 切换到 `pkg/response`
   - 成功响应统一为 `code=0`
   - 参数错误统一走 `response.BadRequest`

2. `api/v1/admin/audit_api.go`
   - 列表接口改为统一 `pagination` 响应结构
   - 成功响应统一为 `code=0`

3. `api/v1/content/project_api.go`
4. `api/v1/content/document_api.go`
5. `api/v1/content/progress_api.go`
   - 分页接口从旧 `shared.Paginated` 收敛到 `pkg/response.Paginated`

未完成：

- `/system/*` 前缀迁移
- PATCH 方法补齐
- 其他仍使用旧 `shared`/裸 `c.JSON` 的 handler 收口

注意：

- 当前分支基线存在与本 issue 无关的 `ObjectID` 迁移编译错误，导致无法在该 worktree 上完成 `go test ./api/v1/admin ./api/v1/content/...` 的整体验证。
- 本轮修改已完成代码收敛和测试断言调整，但需在基线修复后再做完整测试。

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [后端 API 分析](../reports/archived/backend-api-analysis-2026-01-26.md) | 详细 API 问题分析 |
| [设计审查 - API 标准化](../reports/archived/design-review-block7-api-standardization-20260127.md) | API 标准化设计审查 |
| [前后端 API 对齐报告](../reports/archived/2026-01-25-frontend-backend-api-alignment-report.md) | 前后端 API 对比分析 |

---

## 相关Issue

### 依赖Issue（必须先处理）
- 无

### 相关Issue（联合处理）
- [#001: 统一模型层 ID 字段类型](./001-unify-id-type-in-models.md) - ID类型转换影响API响应格式
- [#011: 前后端数据类型不一致](./011-frontend-backend-data-type-inconsistency.md) - 响应拦截器处理不一致与API响应格式相关
- [#006: 数据库索引问题](./006-database-index-issues.md) - API性能优化需要索引支持
