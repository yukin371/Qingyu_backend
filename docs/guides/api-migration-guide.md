# API迁移指南

> **版本**: v1.0
> **更新日期**: 2026-01-29
> **适用范围**: 从shared包迁移到response包的API规范化

## 📋 目录

1. [迁移概述](#迁移概述)
2. [准备工作](#准备工作)
3. [迁移步骤](#迁移步骤)
4. [错误码映射](#错误码映射)
5. [响应函数对照](#响应函数对照)
6. [常见问题](#常见问题)
7. [最佳实践](#最佳实践)
8. [检查清单](#检查清单)

---

## 迁移概述

### 目标

将API处理器从旧的`shared`包迁移到新的统一`response`包，实现：
- ✅ 统一的响应格式
- ✅ 4位错误码规范
- ✅ 简化的API调用
- ✅ 毫秒级时间戳

### 迁移收益

| 方面 | 迁移前 | 迁移后 | 改善 |
|------|--------|--------|------|
| 响应调用 | 4个参数 | 2个参数 | 简化50% |
| 错误码 | 6位 | 4位 | 更规范 |
| 代码行数 | 基准 | -2~3行/文件 | 更简洁 |
| 依赖 | shared+http | response | 依赖减少 |

### Block 7成果参考

- **迁移文件**: 11个Reader模块API
- **响应调用**: 213次成功迁移
- **测试覆盖**: 174/174测试通过（100%）
- **参考文档**: [Block 7进展报告](../../../docs/plans/submodules/backend/api-governance/2026-01-28-block7-api-standardization-progress.md)

---

## 准备工作

### 1. 环境准备

#### 创建feature分支
```bash
git checkout -b feature/block8-writer-migration
```

#### 验证基线测试
```bash
# 运行response包测试
cd Qingyu_backend/pkg/response
go test -v

# 运行Writer模块测试
cd Qingyu_backend/api/v1/writer
go test -v
```

#### 备份当前状态
```bash
# 创建备份分支
git branch backup-before-block8-migration
```

### 2. 理解response包

#### 响应结构
```go
type APIResponse struct {
    Code      int         `json:"code"`       // 0=成功, 4位错误码
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp int64       `json:"timestamp"`  // 毫秒级时间戳
    RequestID string      `json:"request_id"`
}
```

#### 可用函数
```go
// 成功响应
response.Success(c, data)                    // 200 OK
response.Created(c, data)                    // 201 Created
response.NoContent(c)                        // 204 No Content
response.Paginated(c, data, total, page, size, message) // 分页

// 错误响应
response.BadRequest(c, message, details)     // 400
response.Unauthorized(c, message)            // 401
response.Forbidden(c, message)               // 403
response.NotFound(c, message)                // 404
response.Conflict(c, message, details)       // 409
response.InternalError(c, err)               // 500
```

---

## 迁移步骤

### TDD流程：Red-Green-Refactor-Integration

```
┌─────────────────────────────────────────────────────────┐
│ 1. RED - 编写失败的测试（如果需要新测试）                │
└─────────────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────────────┐
│ 2. GREEN - 迁移代码使测试通过                           │
│    ├─ 替换shared.Error调用                              │
│    ├─ 替换shared.Success调用                            │
│    ├─ 替换shared.ValidationError调用                    │
│    └─ 更新Swagger注释                                   │
└─────────────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────────────┐
│ 3. REFACTOR - 重构优化代码                              │
│    ├─ 清理导入依赖                                      │
│    ├─ 提取helper函数                                    │
│    └─ 优化代码结构                                      │
└─────────────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────────────┐
│ 4. INTEGRATION - 集成验证                               │
│    ├─ 运行单元测试                                      │
│    ├─ 运行集成测试                                      │
│    ├─ 编译验证                                          │
│    └─ Git提交                                           │
└─────────────────────────────────────────────────────────┘
```

### 单个文件迁移流程

#### Step 1: 分析文件

```bash
# 统计响应调用次数
grep -E "shared\.(Error|Success|ValidationError)" api/v1/writer/xxx_api.go | wc -l

# 检查特殊场景
grep -E "(WebSocket|c.FileAttachment)" api/v1/writer/xxx_api.go
```

#### Step 2: 替换响应调用

**基本模式**:
```go
// ─────────────────────────────────────────────
// 旧代码 → 新代码
// ─────────────────────────────────────────────

// 1. 成功响应 (200 OK)
shared.Success(c, http.StatusOK, "获取成功", data)
→ response.Success(c, data)

// 2. 创建成功 (201 Created)
shared.Success(c, http.StatusCreated, "创建成功", data)
→ response.Created(c, data)

// 3. 参数错误 (400 Bad Request)
shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
→ response.BadRequest(c, "参数错误", err.Error())

// 4. 参数验证错误
shared.ValidationError(c, err)
→ response.BadRequest(c, "参数错误", err.Error())

// 5. 未授权 (401 Unauthorized)
shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
→ response.Unauthorized(c, "请先登录")

// 6. 禁止访问 (403 Forbidden)
shared.Error(c, http.StatusForbidden, "禁止访问", "无权限")
→ response.Forbidden(c, "无权限")

// 7. 资源不存在 (404 Not Found)
shared.Error(c, http.StatusNotFound, "未找到", "资源不存在")
→ response.NotFound(c, "资源不存在")

// 8. 版本冲突 (409 Conflict)
shared.Error(c, http.StatusConflict, "版本冲突", "文档已被修改")
→ response.Conflict(c, "版本冲突", "文档已被修改")

// 9. 服务器错误 (500 Internal Server Error)
shared.Error(c, http.StatusInternalServerError, "服务器错误", err.Error())
→ response.InternalError(c, err)
```

#### Step 3: 清理导入

```go
// ─────────────────────────────────────────────
// 移除的导入
// ─────────────────────────────────────────────
import (
    "net/http"      // 如果没有WebSocket，移除
    "Qingyu_backend/api/v1/shared"  // 移除
)

// ─────────────────────────────────────────────
// 保留的导入
// ─────────────────────────────────────────────
import (
    "Qingyu_backend/pkg/response"  // 添加
)
```

**注意**: 如果使用了WebSocket，保留`net/http`导入：
```go
import (
    "net/http"  // 保留，WebSocket需要
    "Qingyu_backend/pkg/response"
)
```

#### Step 4: 更新Swagger注释

```go
// ─────────────────────────────────────────────
// 旧注释
// ─────────────────────────────────────────────
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse

// ─────────────────────────────────────────────
// 新注释
// ─────────────────────────────────────────────
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
```

#### Step 5: 验证和测试

```bash
# 编译检查
cd Qingyu_backend
go build ./api/v1/writer/xxx_api.go

# 运行测试
cd api/v1/writer
go test -v -run TestXxx

# 运行完整测试
go test -v
```

#### Step 6: Git提交

```bash
git add api/v1/writer/xxx_api.go
git commit -m "feat(api): migrate xxx_api to new response package

- Replace all shared.Error calls with response functions
- Replace all shared.Success calls with response functions
- Remove HTTP status code parameters
- Update Swagger annotations
- Clean up imports (remove shared, net/http)"
```

---

## 错误码映射

### 6位错误码 → 4位错误码

| 旧错误码 | 新错误码 | 常量名 | 含义 |
|---------|---------|--------|------|
| 0 | 0 | CodeSuccess | 成功 |
| 100001 | 1001 | CodeParamError | 参数错误 |
| 100403 | 1003 | CodeForbidden | 禁止访问 |
| 100404 | 1004 | CodeNotFound | 资源不存在 |
| 100409 | 1006 | CodeConflict | 资源冲突 |
| 100500 | 5000 | CodeInternalError | 服务器内部错误 |
| 100601 | 1002 | CodeUnauthorized | 未授权 |

### 错误码分类

```go
// 0xxx - 成功
0 = CodeSuccess

// 1xxx - 客户端错误
1001 = CodeParamError       // 参数错误
1002 = CodeUnauthorized     // 未授权
1003 = CodeForbidden        // 禁止访问
1004 = CodeNotFound         // 资源不存在
1005 = CodeMethodNotAllowed // 方法不允许
1006 = CodeConflict         // 资源冲突

// 2xxx - 用户相关错误
2001 = CodeUserNotFound     // 用户不存在
2002 = CodeUserDisabled     // 用户被禁用
2003 = CodeInvalidPassword  // 密码错误

// 3xxx - 业务逻辑错误
3001 = CodeBusinessError    // 业务错误
3002 = CodePermissionDenied // 权限不足

// 4xxx - 限流相关
4001 = CodeRateLimitExceeded // 超出限流

// 5xxx - 服务器错误
5000 = CodeInternalError     // 服务器内部错误
5001 = CodeDatabaseError     // 数据库错误
5002 = CodeServiceUnavailable // 服务不可用
```

---

## 响应函数对照

### 完整对照表

| HTTP状态码 | 旧函数 | 新函数 | 参数变化 |
|-----------|--------|--------|---------|
| 200 OK | `shared.Success(c, http.StatusOK, msg, data)` | `response.Success(c, data)` | 4参数→2参数 |
| 201 Created | `shared.Success(c, http.StatusCreated, msg, data)` | `response.Created(c, data)` | 4参数→2参数 |
| 204 No Content | `shared.Success(c, http.StatusNoContent, msg, nil)` | `response.NoContent(c)` | 4参数→1参数 |
| 400 Bad Request | `shared.Error(c, http.StatusBadRequest, msg, details)` | `response.BadRequest(c, msg, details)` | 4参数→3参数 |
| 401 Unauthorized | `shared.Error(c, http.StatusUnauthorized, msg, details)` | `response.Unauthorized(c, msg)` | 4参数→2参数 |
| 403 Forbidden | `shared.Error(c, http.StatusForbidden, msg, details)` | `response.Forbidden(c, msg)` | 4参数→2参数 |
| 404 Not Found | `shared.Error(c, http.StatusNotFound, msg, details)` | `response.NotFound(c, msg)` | 4参数→2参数 |
| 409 Conflict | `shared.Error(c, http.StatusConflict, msg, details)` | `response.Conflict(c, msg, details)` | 4参数→3参数 |
| 500 Internal Error | `shared.Error(c, http.StatusInternalServerError, msg, err.Error())` | `response.InternalError(c, err)` | 4参数→2参数 |

### 分页响应

```go
// ─────────────────────────────────────────────
// 旧代码
// ─────────────────────────────────────────────
shared.Success(c, http.StatusOK, "获取成功", gin.H{
    "list": data,
    "total": total,
    "page": page,
    "pageSize": pageSize,
})

// ─────────────────────────────────────────────
// 新代码（推荐）
// ─────────────────────────────────────────────
response.Paginated(c, data, total, page, pageSize, "获取成功")

// 或者使用Success返回自定义结构
response.Success(c, gin.H{
    "list": data,
    "total": total,
    "page": page,
    "pageSize": pageSize,
})
```

---

## 常见问题

### Q1: 如何处理版本冲突？

**问题**: 编辑器API需要特殊处理版本冲突

**解决方案**:
```go
// 检查错误类型
if err.Error() == "版本冲突" {
    response.Conflict(c, "版本冲突", "文档已被其他用户修改，请刷新后重试")
    return
}
```

### Q2: 如何保留WebSocket支持？

**问题**: WebSocket需要`net/http`导入

**解决方案**:
```go
import (
    "net/http"  // 保留，WebSocket需要
    "Qingyu_backend/pkg/response"
)

// WebSocket升级不需要修改
upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
```

### Q3: 如何处理文件下载？

**问题**: 文件下载使用`c.FileAttachment`

**解决方案**:
```go
// 文件下载不需要修改
c.FileAttachment(filePath, fileName)

// 但错误处理需要迁移
if err != nil {
    response.InternalError(c, err)
    return
}
```

### Q4: 如何处理批量操作的异步响应？

**问题**: 批量操作异步执行，立即返回

**解决方案**:
```go
// 提交批量操作
response.Success(c, gin.H{
    "batchId": batchOp.ID.Hex(),
    "status": "submitted",
})

// 异步执行
go func() {
    api.batchOpSvc.Execute(ctx, batchId)
}()
```

### Q5: Swagger注释引用了shared.APIResponse怎么办？

**问题**: Swagger注释中的`shared.APIResponse`需要更新

**解决方案**:
```go
// 批量替换
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse

// 改为
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse

// 或者使用具体的数据类型
// @Success 200 {object} response.APIResponse{data=DocumentResponse}
```

### Q6: 如何处理复杂的错误场景？

**问题**: 需要根据错误类型返回不同的响应

**解决方案**:
```go
// 使用errors.Is检查特定错误
if errors.Is(err, ErrNotFound) {
    response.NotFound(c, "文档不存在")
    return
}

if errors.Is(err, ErrUnauthorized) {
    response.Unauthorized(c, "无权访问")
    return
}

// 默认处理
response.InternalError(c, err)
```

### Q7: 测试失败怎么办？

**问题**: 迁移后测试失败

**解决方案**:
1. 检查响应格式是否匹配
2. 检查错误码是否正确
3. 检查时间戳格式（毫秒级）
4. 查看测试输出，定位具体问题

```bash
# 运行详细测试
go test -v -run TestFailingTest

# 查看覆盖率
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Q8: 如何验证迁移完整性？

**问题**: 确保没有遗漏的shared调用

**解决方案**:
```bash
# 搜索所有shared调用
grep -r "shared\.(Error|Success|ValidationError)" api/v1/writer

# 应该没有输出（除了注释）

# 检查导入
grep -r "Qingyu_backend/api/v1/shared" api/v1/writer

# 应该没有输出（除了Swagger注释）
```

### Q9: 如何处理第三方库的错误？

**问题**: 第三方库返回的错误需要转换

**解决方案**:
```go
// 包装第三方错误
if err := thirdPartyCall(); err != nil {
    response.InternalError(c, fmt.Errorf("第三方服务错误: %w", err))
    return
}

// 或者转换为业务错误
if err := thirdPartyCall(); err != nil {
    response.BadRequest(c, "第三方服务不可用", err.Error())
    return
}
```

### Q10: 如何处理自定义响应格式？

**问题**: 需要返回自定义的响应结构

**解决方案**:
```go
// 使用Success返回自定义结构
response.Success(c, MyCustomResponse{
    Field1: value1,
    Field2: value2,
    Nested: NestedStruct{
        Field3: value3,
    },
})

// 或使用gin.H
response.Success(c, gin.H{
    "customField": customValue,
    "data": data,
})
```

---

## 最佳实践

### 1. 遵循TDD流程

```
Red → Green → Refactor → Integration
```

- **Red**: 先写测试（如果需要）
- **Green**: 快速通过测试
- **Refactor**: 优化代码结构
- **Integration**: 确保所有测试通过

### 2. 小步快跑，频繁提交

```bash
# 每迁移1-2个函数就提交一次
git add xxx_api.go
git commit -m "feat(api): migrate xxx function"

# 而不是迁移完整个文件才提交
```

### 3. 保持测试覆盖

```bash
# 运行测试确保覆盖
go test -v -cover

# 目标：每个API至少有单元测试
```

### 4. 文档同步更新

```go
// 更新Swagger注释
// @Success 200 {object} response.APIResponse

// 更新注释说明
// GetDocument 获取文档详情
// 返回文档的完整信息，包括内容和元数据
```

### 5. 错误处理一致性

```go
// 统一的错误处理模式
if err != nil {
    response.InternalError(c, err)
    return
}

// 而不是
if err != nil {
    log.Error(err)
    c.JSON(500, gin.H{"error": err.Error()})
    return
}
```

### 6. 提取helper减少重复

```go
// 提取getUserID helper
func getUserID(c *gin.Context) (string, error) {
    userID, exists := c.Get("userId")
    if !exists {
        return "", errors.New("用户未登录")
    }
    return userID.(string), nil
}

// 使用
userID, err := getUserID(c)
if err != nil {
    response.Unauthorized(c, "请先登录")
    return
}
```

### 7. 验证参数统一模式

```go
// 统一的参数验证模式
var req CreateRequest
if err := c.ShouldBindJSON(&req); err != nil {
    response.BadRequest(c, "参数错误", err.Error())
    return
}

// 验证必填字段
if req.Name == "" {
    response.BadRequest(c, "参数错误", "名称不能为空")
    return
}
```

### 8. 分页响应标准化

```go
// 标准分页响应
response.Paginated(c, list, total, page, pageSize, "获取成功")

// 而不是自定义结构
response.Success(c, gin.H{
    "list": list,
    "total": total,
    "page": page,
    "pageSize": pageSize,
})
```

---

## 检查清单

### 迁移前检查

- [ ] 已创建feature分支
- [ ] 已备份当前代码
- [ ] 已运行基线测试
- [ ] 已了解response包API
- [ ] 已阅读本文档

### 迁移中检查（每个文件）

- [ ] 替换所有`shared.Error`调用
- [ ] 替换所有`shared.Success`调用
- [ ] 替换所有`shared.ValidationError`调用
- [ ] 移除HTTP状态码参数
- [ ] 更新错误码（6位→4位）
- [ ] 清理导入依赖
- [ ] 更新Swagger注释
- [ ] 代码编译通过

### 迁移后检查（每个文件）

- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] 无shared包残留
- [ ] Swagger注释正确
- [ ] Git提交成功

### 整体验收

- [ ] 所有文件迁移完成
- [ ] 所有测试通过
- [ ] 无编译错误
- [ ] 无shared包残留
- [ ] Swagger文档完整
- [ ] 代码审查通过
- [ ] PR创建成功

---

## 参考资料

### Block 7参考文档

- [Block 7 API规范化试点 - 进展报告](../../../docs/plans/submodules/backend/api-governance/2026-01-28-block7-api-standardization-progress.md)
- [Block 7 全面回归测试报告](../reports/block7-p2-regression-test-report.md)

### 相关代码

- [Response包实现](../../pkg/response/writer.go)
- [错误码定义](../../pkg/response/codes.go)
- [Reader模块示例](../../api/v1/reader/)

### Writer模块分析

- [Writer模块迁移预分析报告](../analysis/2026-01-29-writer-migration-analysis.md)
- [Writer模块复杂度矩阵](../analysis/writer-complexity-matrix.json)

---

**文档版本**: v1.0
**最后更新**: 2026-01-29
**维护者**: Backend Team
