# ID 错误处理指南

## 概述

本文档定义了项目中 ID 相关错误的统一处理策略，确保从 Repository → Service → API 的错误翻译链路清晰一致。

## 错误类型定义

```go
// repository/errors.go
var (
    // ErrEmptyID 表示ID为空字符串
    ErrEmptyID = errors.New("ID cannot be empty")

    // ErrInvalidIDFormat 表示ID格式无效（不是有效的ObjectID）
    ErrInvalidIDFormat = errors.New("invalid ID format")
)
```

## 错误判断工具

```go
// repository/id_converter.go

// IsIDError 判断是否为ID相关错误
func IsIDError(err error) bool {
    return errors.Is(err, ErrEmptyID) || errors.Is(err, ErrInvalidIDFormat)
}
```

## 分层处理策略

### 1. Repository 层

**职责**：返回原始的 ID 错误，不做业务翻译

```go
// 使用统一工具
func (r *SomeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Model, error) {
    // 查询逻辑...
    // 如果查询失败，返回原始错误
}
```

### 2. Service 层

**职责**：将 ID 错误翻译为业务错误

```go
func (s *SomeService) DoSomething(ctx context.Context, id string) error {
    // 解析 ID
    oid, err := repository.ParseID(id)
    if err != nil {
        // 统一翻译为业务错误
        if errors.Is(err, repository.ErrEmptyID) {
            return serviceerrors.ErrMissingParameter  // 或自定义错误
        }
        return serviceerrors.ErrInvalidID
    }

    // 调用 Repository
    result, err := s.repo.GetByID(ctx, oid)
    if err != nil {
        return err  // 透传其他错误
    }
    // ...
}
```

### 3. API 层

**职责**：保留快速失败机制，统一错误响应格式

#### 模式 A：快速失败（推荐用于热路径）

```go
func (api *SomeAPI) GetSomething(c *gin.Context) {
    id := c.Param("id")

    // 快速失败 - 明显错误在API层拦截
    if id == "" {
        response.BadRequest(c, "参数错误", "ID不能为空")
        return
    }

    // 格式校验可以在这里做，也可以交给 Service
    // 如果 Service 已经做了翻译，这里可以省略

    // 调用 Service
    result, err := api.service.GetSomething(ctx, id)
    if err != nil {
        // 统一错误翻译
        if repository.IsIDError(err) {
            response.BadRequest(c, "参数错误", "无效的ID格式")
            return
        }
        // 其他错误处理...
    }

    response.Success(c, result)
}
```

#### 模式 B：委托 Service（简化版）

```go
func (api *SomeAPI) GetSomething(c *gin.Context) {
    id := c.Param("id")

    // 只做空值快速检查
    if id == "" {
        response.BadRequest(c, "参数错误", "ID不能为空")
        return
    }

    // 格式校验交给 Service
    result, err := api.service.GetSomething(ctx, id)
    if err != nil {
        if repository.IsIDError(err) {
            response.BadRequest(c, "参数错误", "无效的ID格式")
            return
        }
        // 其他错误处理...
    }

    response.Success(c, result)
}
```

## 错误响应格式

所有 ID 相关错误统一使用以下响应格式：

```json
{
    "code": 400,
    "message": "参数错误",
    "error": "无效的ID格式",
    "timestamp": 1234567890123
}
```

**调用方式**：
```go
response.BadRequest(c, "参数错误", "无效的ID格式")
response.BadRequest(c, "参数错误", "ID不能为空")
```

## 最佳实践

### ✅ 推荐做法

1. **Service 层使用 `repository.ParseID`**
   ```go
   oid, err := repository.ParseID(id)
   if err != nil {
       return nil, err  // 统一错误语义
   }
   ```

2. **API 层使用 `repository.IsIDError` 判断**
   ```go
   if repository.IsIDError(err) {
       response.BadRequest(c, "参数错误", "无效的ID格式")
       return
   }
   ```

3. **保留 API 层快速失败**
   ```go
   if id == "" {
       response.BadRequest(c, "参数错误", "ID不能为空")
       return
   }
   ```

### ❌ 避免的做法

1. **直接使用 `primitive.ObjectIDFromHex`**
   ```go
   // 不推荐
   oid, err := primitive.ObjectIDFromHex(id)
   ```

2. **自定义错误消息不一致**
   ```go
   // 不推荐 - 消息不统一
   response.BadRequest(c, "错误", "id不对")
   ```

3. **跳过 API 层快速失败**
   ```go
   // 不推荐 - 空值应该在 API 层拦截
   result, err := api.service.GetSomething(ctx, c.Param("id"))
   ```

## 错误翻译对照表

| 原始错误 | Service 层翻译 | API 层响应 |
|---------|---------------|-----------|
| `repository.ErrEmptyID` | `serviceerrors.ErrMissingParameter` | `BadRequest("参数错误", "ID不能为空")` |
| `repository.ErrInvalidIDFormat` | `serviceerrors.ErrInvalidID` | `BadRequest("参数错误", "无效的ID格式")` |

## 相关文件

- `repository/errors.go` - ID 错误定义
- `repository/id_converter.go` - ID 转换工具
- `api/v1/shared/response.go` - 响应工具
- `service/shared/errors/service_errors.go` - 业务错误定义
