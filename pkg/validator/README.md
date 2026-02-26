# 验证器使用指南

## 概述

`pkg/validator` 提供了统一的请求验证机制，基于 [go-playground/validator](https://github.com/go-playground/validator) 实现。本包采用单例模式管理全局验证器实例，并提供了13种自定义验证规则。

## 核心组件

### 1. 全局验证器 (`validator.go`)

```go
// 获取全局验证器实例（单例模式）
v := validator.GetValidator()

// 验证结构体
err := validator.ValidateStruct(request)

// 验证并返回友好错误
validationErrors := validator.ValidateStructWithErrors(request)
```

### 2. 自定义验证规则 (`custom_validators.go`)

提供了13种自定义验证规则，涵盖金额、文件、字符串、业务等场景。

### 3. 请求验证封装 (`api/v1/shared/request_validator.go`)

为API层提供了便捷的请求验证方法。

## 内置验证规则

### 金额验证

| 规则标签 | 说明 | 参数 | 示例 |
|---------|------|------|------|
| `amount` | 验证金额格式（最多2位小数） | 无 | `validate:"amount"` |
| `positive_amount` | 验证正数金额（> 0） | 无 | `validate:"positive_amount"` |
| `amount_range` | 验证金额范围（0.01 - 1000000） | 无 | `validate:"amount_range"` |

**使用示例：**

```go
type RechargeRequest struct {
    Amount float64 `json:"amount" validate:"positive_amount,amount_range"`
}
```

### 文件验证

| 规则标签 | 说明 | 支持的类型 |
|---------|------|-----------|
| `file_type` | 验证文件类型 | 图片、PDF、Office、文本、ZIP |
| `file_size` | 验证文件大小（最大50MB） | 0 < size <= 50MB |

**支持的文件类型：**

- 图片：`image/jpeg`, `image/png`, `image/gif`, `image/webp`
- 文档：`application/pdf`, `application/msword`, `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- 其他：`text/plain`, `application/zip`

**使用示例：**

```go
type UploadRequest struct {
    FileType string `json:"file_type" validate:"file_type"`
    FileSize int64  `json:"file_size" validate:"file_size"`
}
```

### 字符串验证

| 规则标签 | 说明 | 规则 |
|---------|------|------|
| `username` | 验证用户名 | 3-20个字符，字母数字下划线 |
| `phone` | 验证手机号 | 中国大陆手机号格式（1[3-9]开头） |
| `strong_password` | 验证强密码 | 至少8位，包含大小写字母和数字 |

**使用示例：**

```go
type UserRequest struct {
    Username string `json:"username" validate:"username"`
    Phone    string `json:"phone" validate:"phone"`
    Password string `json:"password" validate:"strong_password"`
}
```

### 业务验证

| 规则标签 | 说明 | 允许的值 |
|---------|------|---------|
| `transaction_type` | 验证交易类型 | `recharge`, `consume`, `transfer`, `refund`, `withdraw` |
| `withdraw_account` | 验证提现账号 | `alipay:xxx`, `wechat:xxx`, `bank:xxx` |
| `content_type` | 验证内容类型 | `book`, `chapter`, `comment`, `review` |

**使用示例：**

```go
type TransactionRequest struct {
    Type    string  `json:"type" validate:"transaction_type"`
    Amount  float64 `json:"amount" validate:"positive_amount,amount_range"`
    Account string  `json:"account" validate:"withdraw_account"`
}
```

## 使用方式

### 1. API层验证（推荐）

在API Handler中使用 `request_validator` 包提供的工具：

```go
package user

import (
    "github.com/gin-gonic/gin"
    appValidator "Qingyu_backend/pkg/validator"
    "Qingyu_backend/api/v1/shared"
)

type RegisterRequest struct {
    Username string `json:"username" binding:"required" validate:"username"`
    Email    string `json:"email" binding:"required,email" validate:"required,email"`
    Password string `json:"password" binding:"required" validate:"strong_password"`
}

func RegisterHandler(c *gin.Context) {
    var req RegisterRequest

    // 使用共享验证器
    if !shared.ValidateRequest(c, &req) {
        return // 错误响应已在 ValidateRequest 中处理
    }

    // 处理业务逻辑...
}
```

### 2. Service层验证

对于需要数据库查询的验证（如唯一性检查），在Service层进行：

```go
package user

type UserValidator struct {
    repo UserRepository
}

func (v *UserValidator) ValidateCreate(ctx context.Context, user *User) error {
    var errors ValidationErrors

    // 1. 基础字段验证
    if err := v.validateUsername(user.Username); err != nil {
        errors = append(errors, *err)
    }

    // 2. 唯一性验证
    if exists, _ := v.repo.ExistsByUsername(ctx, user.Username); exists {
        errors = append(errors, ValidationError{
            Field: "username",
            Message: "用户名已存在",
            Code: "DUPLICATE",
        })
    }

    if len(errors) > 0 {
        return errors
    }
    return nil
}
```

### 3. 结构体验证

直接使用全局验证器：

```go
import "Qingyu_backend/pkg/validator"

type CreateProjectRequest struct {
    Name        string `json:"name" binding:"required" validate:"required,min=3,max=50"`
    Description string `json:"description" validate:"max=500"`
    IsPublic    bool   `json:"is_public"`
}

func CreateProject(req *CreateProjectRequest) error {
    // 验证并获取详细错误
    validationErrors := validator.ValidateStructWithErrors(req)
    if len(validationErrors) > 0 {
        // 获取字段级错误
        fieldErrors := validationErrors.GetFieldErrors()
        for field, message := range fieldErrors {
            log.Printf("字段 %s 验证失败: %s", field, message)
        }
        return validationErrors
    }
    return nil
}
```

## 错误处理

### 验证错误结构

```go
// ValidationErrors 验证错误集合
type ValidationErrors []FieldError

// FieldError 单个字段错误
type FieldError struct {
    Field   string // 字段名
    Message string // 错误消息
    Tag     string // 验证规则标签
    Value   interface{} // 实际值
}
```

### 获取字段级错误

```go
validationErrors := validator.ValidateStructWithErrors(request)

// 获取所有字段错误（字段名 -> 错误消息）
fieldErrors := validationErrors.GetFieldErrors()

// 返回给客户端
c.JSON(400, gin.H{
    "code": "VALIDATION_ERROR",
    "message": "请求参数验证失败",
    "errors": fieldErrors,
})
```

### 错误消息示例

```json
{
  "code": 400,
  "message": "请求参数验证失败",
  "errors": {
    "username": "用户名长度不能少于3个字符",
    "email": "邮箱格式不正确",
    "amount": "金额必须在 0.01 到 1000000 之间"
  }
}
```

## 最佳实践

### 1. 验证层级划分

```
┌─────────────────────────────────────────┐
│           API 层（Handler）              │
│  - 格式验证（ShouldBindJSON）            │
│  - 基础规则验证（binding/validate tag）  │
│  - 使用 request_validator 封装           │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         Service 层（Business）           │
│  - 业务逻辑验证                          │
│  - 唯一性验证（需要查库）                │
│  - 权限验证                              │
│  - 使用 Validator 结构体                │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│       Repository 层（Data）              │
│  - 数据完整性验证                        │
│  - 外键引用验证                          │
└─────────────────────────────────────────┘
```

### 2. 命名规范

**验证规则标签：**
- 使用小写和下划线：`positive_amount`, `file_type`
- 语义清晰：`transaction_type` 而非 `tx_type`

**验证器函数：**
- 格式：`validate{RuleName}`
- 示例：`validateAmount`, `validateUsername`

**验证错误：**
- 错误码使用大写下划线：`INVALID_FORMAT`, `DUPLICATE`
- 错误消息使用中文，用户友好

### 3. 性能考虑

**使用 sync.Once 确保单例：**

```go
var once sync.Once
var validate *validator.Validate

func GetValidator() *validator.Validate {
    once.Do(func() {
        validate = validator.New()
        RegisterCustomValidators(validate)
    })
    return validate
}
```

**避免重复验证：**

```go
// 不推荐 - 重复验证
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    if err := validator.ValidateStruct(req); err != nil {
        return err
    }
    // ...又在另一个地方验证
}

// 推荐 - 在API层验证
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if !shared.ValidateRequest(c, &req) {
        return
    }
    // Service层直接使用
}
```

### 4. 自定义验证规则开发

**添加新的自定义验证器：**

1. 在 `custom_validators.go` 中添加验证函数：

```go
// validateISBN 验证ISBN格式
func validateISBN(fl validator.FieldLevel) bool {
    isbn := fl.Field().String()
    // 实现ISBN验证逻辑
    matched, _ := regexp.MatchString(`^\d{3}-\d{10}$`, isbn)
    return matched
}
```

2. 在 `RegisterCustomValidators` 中注册：

```go
validations := []struct {
    tag string
    fn  validator.Func
}{
    // ...现有验证器
    {"isbn", validateISBN}, // 新增
}
```

3. 使用新验证器：

```go
type BookRequest struct {
    ISBN string `json:"isbn" validate:"isbn"`
}
```

### 5. 测试验证器

**单元测试示例：**

```go
func TestValidateAmount(t *testing.T) {
    tests := []struct {
        name    string
        amount  float64
        wantErr bool
    }{
        {"valid amount", 100.50, false},
        {"negative amount", -50.00, true},
        {"too many decimals", 100.505, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            field := level.Field{}
            field.SetReflectedValue(reflect.ValueOf(tt.amount))

            err := validateAmount(field)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateAmount() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 常见问题

### Q1: 为什么验证器注册失败使用日志警告而不是返回错误？

**A:** 为了保证系统启动的健壮性。如果某个验证器注册失败，系统仍然可以启动，只是该特定规则不可用。在生产环境中，这些警告应该被监控。

### Q2: 什么时候使用 tag 验证，什么时候使用 Service 层验证？

**A:**
- **Tag 验证**：格式检查、范围检查、枚举值等不需要数据库的验证
- **Service 验证**：唯一性检查、权限检查、业务规则等需要数据库的验证

### Q3: 如何自定义验证错误消息？

**A:** 可以使用 `msg` 参数：

```go
type Request struct {
    Name string `validate:"required,min=3,max=50" msg:"名称长度必须在3-50个字符之间"`
}
```

或者通过 `TranslateError` 进行自定义翻译。

## 相关文档

- [验证规范文档](../../../docs/standards/validation_standard.md)
- [API架构文档](../../../architecture/api_architecture.md)
- [错误处理规范](../errors/README.md)

## 维护者

- 创建日期：2026-02-26
- 维护者：Backend Team
- 版本：1.0.0
