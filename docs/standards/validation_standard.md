# 请求验证规范

## 概述

本文档定义了 Qingyu_backend 项目中请求验证的统一标准和最佳实践。遵循这些规范可以确保验证逻辑的一致性、可维护性和可测试性。

## 验证层级

### 三层验证架构

```
┌──────────────────────────────────────────────────────────────┐
│                       API 层 (Handler)                        │
│  职责：格式验证、基础规则验证                                  │
│  工具：binding/validate tag, request_validator                │
│  示例：required, email, min, max, 自定义tag                   │
└──────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────┐
│                     Service 层 (Business)                     │
│  职责：业务逻辑验证、唯一性验证、权限验证                      │
│  工具：Validator 结构体、Repository 查询                      │
│  示例：用户名唯一、余额充足、权限检查                          │
└──────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────┐
│                   Repository 层 (Data Integrity)              │
│  职责：数据完整性验证、外键引用验证                            │
│  工具：ReferenceValidator                                    │
│  示例：外键存在性、关联关系有效性                              │
└──────────────────────────────────────────────────────────────┘
```

### 各层职责详解

#### 1. API 层验证

**适用场景：**
- JSON 格式验证
- 必填字段检查
- 字段长度限制
- 数据格式检查（email、URL等）
- 数值范围检查
- 枚举值验证

**验证方式：**
```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required" validate:"username"`
    Email    string `json:"email" binding:"required,email" validate:"required,email"`
    Password string `json:"password" binding:"required" validate:"strong_password"`
}
```

**最佳实践：**
- 使用 `binding` 标签进行基础验证
- 使用 `validate` 标签进行自定义规则验证
- 使用 `shared.ValidateRequest` 统一处理
- 不要在 API 层进行数据库查询

#### 2. Service 层验证

**适用场景：**
- 业务规则验证
- 唯一性约束检查
- 权限验证
- 状态转换验证
- 依赖关系验证

**验证方式：**
```go
type UserValidator struct {
    repo UserRepository
}

func (v *UserValidator) ValidateCreate(ctx context.Context, user *User) error {
    // 唯一性验证
    if exists, _ := v.repo.ExistsByUsername(ctx, user.Username); exists {
        return ValidationError{
            Field: "username",
            Message: "用户名已存在",
            Code: "DUPLICATE",
        }
    }

    // 业务规则验证
    if user.Age < 18 {
        return ValidationError{
            Field: "age",
            Message: "用户必须年满18岁",
            Code: "BUSINESS_RULE_VIOLATION",
        }
    }

    return nil
}
```

**最佳实践：**
- 为每个 Service 创建独立的 Validator
- 使用结构化的 ValidationError
- 区分字段错误和业务规则错误
- 提供清晰的错误码

#### 3. Repository 层验证

**适用场景：**
- 外键引用完整性
- 关联数据存在性
- 数据一致性检查

**验证方式：**
```go
// 使用 ReferenceValidator 验证引用
func (s *Service) CreateComment(ctx context.Context, req *CreateCommentRequest) error {
    // 验证用户和书籍存在
    if err := s.refValidator.ValidateCommentReference(ctx, req.AuthorID, req.BookID); err != nil {
        return err
    }

    // 继续处理...
}
```

## 验证规则

### 命名规范

#### 验证规则标签

**格式：** 小写，单词间用下划线分隔

```go
// 正确示例
"positive_amount"
"file_type"
"transaction_type"
"withdraw_account"

// 错误示例
"positiveAmount"  // 使用了驼峰命名
"pos_amt"        // 缩写不清晰
"pa"              // 过于简短
```

#### 验证器函数

**格式：** `validate{RuleName}`

```go
// 正确示例
func validateAmount(fl validator.FieldLevel) bool
func validatePositiveAmount(fl validator.FieldLevel) bool
func validateFileType(fl validator.FieldLevel) bool

// 错误示例
func checkAmount(...)    // 应该使用 validate 前缀
func amountValidate(...) // 动词在后
func ValidateAmount(...) // 导出函数不需要（仅内部使用）
```

#### 验证错误码

**格式：** 大写，单词间用下划线分隔

```go
// 通用错误码
"REQUIRED"           // 必填字段
"INVALID_FORMAT"     // 格式无效
"MIN_LENGTH"         // 长度不足
"MAX_LENGTH"         // 长度超限
"DUPLICATE"          // 重复值

// 业务错误码
"WEAK_PASSWORD"      // 弱密码
"RESERVED_NAME"      // 保留名称
"INSUFFICIENT_BALANCE" // 余额不足
```

### 验证规则分类

#### 1. 基础验证规则

使用 Gin 框架提供的内置验证器：

```go
// 必填
binding:"required"

// 字符串
binding:"email"              // 邮箱格式
binding:"url"                // URL格式
binding:"min=3"              // 最小长度
binding:"max=50"             // 最大长度
binding:"len=10"             // 固定长度
binding:"eq=admin"           // 等于
binding:"ne=admin"           // 不等于
binding:"oneof=red green"    // 枚举值

// 数值
binding:"min=1"              // 最小值
binding:"max=100"            // 最大值
binding:"gt=0"               // 大于
binding:"gte=0"              // 大于等于
binding:"lt=100"             // 小于
binding:"lte=100"            // 小于等于

// 特殊
binding:"-"                  // 跳过验证
binding:"omitempty"          // 空值时跳过
```

#### 2. 自定义验证规则

在 `pkg/validator/custom_validators.go` 中定义：

```go
// 金额类
amount           // 金额格式（最多2位小数）
positive_amount  // 正数金额
amount_range     // 金额范围（0.01-1000000）

// 文件类
file_type        // 文件类型
file_size        // 文件大小（最大50MB）

// 字符串类
username         // 用户名（3-20字符，字母数字下划线）
phone            // 手机号（中国大陆）
strong_password  // 强密码（8位以上，大小写字母+数字）

// 业务类
transaction_type // 交易类型
withdraw_account // 提现账号
content_type     // 内容类型
```

### 注册规范

所有自定义验证器必须在 `RegisterCustomValidators` 中注册：

```go
func RegisterCustomValidators(v *validator.Validate) {
    validations := []struct {
        tag string
        fn  validator.Func
    }{
        // 金额验证
        {"amount", validateAmount},
        {"positive_amount", validatePositiveAmount},
        {"amount_range", validateAmountRange},

        // ... 其他验证器
    }

    for _, validation := range validations {
        if err := v.RegisterValidation(validation.tag, validation.fn); err != nil {
            log.Printf("Warning: failed to register validation '%s': %v", validation.tag, err)
        }
    }
}
```

**注意事项：**
- 必须处理注册失败的情况（使用日志警告）
- 不要在注册失败时终止程序
- 注册失败应在监控中被发现

## 错误处理规范

### 错误响应格式

#### 单个字段错误

```json
{
  "code": 400,
  "message": "请求参数验证失败",
  "errors": {
    "username": "用户名长度不能少于3个字符",
    "email": "邮箱格式不正确"
  }
}
```

#### 多个字段错误

```json
{
  "code": 400,
  "message": "请求参数验证失败",
  "errors": {
    "username": "用户名长度不能少于3个字符",
    "email": "邮箱格式不正确",
    "password": "密码必须包含至少一个字母和一个数字",
    "amount": "金额必须在 0.01 到 1000000 之间"
  }
}
```

### 错误消息规范

**格式：** 清晰、用户友好、可操作

```go
// 好的错误消息
"用户名长度不能少于3个字符"
"邮箱格式不正确"
"密码必须包含至少一个字母和一个数字"

// 不好的错误消息
"验证失败"              // 太笼统
"Invalid username"      // 不中文化
"username: invalid"     // 不友好
"Err: username format"  // 技术性太强
```

### 错误码规范

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| `REQUIRED` | 必填字段缺失 | 400 |
| `INVALID_FORMAT` | 格式不正确 | 400 |
| `MIN_LENGTH` | 长度不足 | 400 |
| `MAX_LENGTH` | 长度超限 | 400 |
| `INVALID_EMAIL` | 邮箱格式错误 | 400 |
| `DUPLICATE` | 值重复 | 409 |
| `NOT_FOUND` | 资源不存在 | 404 |
| `BUSINESS_RULE_VIOLATION` | 违反业务规则 | 400 |
| `INSUFFICIENT_BALANCE` | 余额不足 | 400 |
| `WEAK_PASSWORD` | 弱密码 | 400 |
| `RESERVED_NAME` | 保留名称 | 400 |

## 常用验证场景

### 1. 金额验证

```go
type RechargeRequest struct {
    Amount float64 `json:"amount" validate:"positive_amount,amount_range"`
}

// 验证规则：
// - 金额必须大于 0
// - 金额范围：0.01 - 1000000.00
// - 最多2位小数
```

### 2. 文件验证

```go
type FileUploadRequest struct {
    FileName string `json:"file_name" binding:"required"`
    FileType string `json:"file_type" validate:"file_type"`
    FileSize int64  `json:"file_size" validate:"file_size"`
}

// 验证规则：
// - 文件类型：图片、PDF、Office、文本、ZIP
// - 文件大小：0 < size <= 50MB
```

### 3. 字符串验证

```go
type UserProfileRequest struct {
    Username string `json:"username" validate:"username"`
    Phone    string `json:"phone" validate:"phone"`
    Password string `json:"password" validate:"strong_password"`
}

// 验证规则：
// - 用户名：3-20字符，字母数字下划线
// - 手机号：1[3-9]开头的11位数字
// - 密码：至少8位，包含大小写字母和数字
```

### 4. 枚举验证

```go
type TransactionRequest struct {
    Type string `json:"type" validate:"transaction_type"`
}

// 使用 Gin 内置验证
type TransactionRequest struct {
    Type string `json:"type" binding:"required,oneof=recharge consume transfer refund withdraw"`
}
```

### 5. 嵌套验证

```go
type CreateProjectRequest struct {
    Name        string           `json:"name" binding:"required" validate:"required,min=3,max=50"`
    Description string           `json:"description" validate:"max=500"`
    Settings    ProjectSettings  `json:"settings" validate:"dive"`
}

type ProjectSettings struct {
    IsPublic    bool `json:"is_public"`
    AllowComment bool `json:"allow_comment"`
}

// 使用 dive 验证嵌套结构
type CreateProjectRequest struct {
    Members []Member `json:"members" validate:"dive"`
}
```

### 6. 条件验证

```go
type WithdrawRequest struct {
    Amount  float64 `json:"amount" validate:"positive_amount,amount_range"`
    Method  string  `json:"method" binding:"required,oneof=alipay wechat bank"`
    Account string  `json:"account" binding:"required"`
    // Account 的格式验证需要在 Service 层根据 Method 进行
}
```

## 最佳实践

### 1. 验证优先级

```
1. 必填字段检查 (required)
2. 格式验证 (email, url, 等)
3. 长度验证 (min, max)
4. 范围验证 (数值范围)
5. 业务规则验证 (Service 层)
6. 唯一性验证 (Service 层)
```

### 2. 避免过度验证

```go
// 不推荐 - 验证太多细节
type Request struct {
    Name string `validate:"required,min=3,max=50,ascii,printascii"`
}

// 推荐 - 只验证必要的
type Request struct {
    Name string `validate:"required,min=3,max=50"`
    // 其他格式化细节可以在使用时处理
}
```

### 3. 提供验证反馈

```go
// 前端可以使用验证错误来提供实时反馈
{
  "errors": {
    "password": "密码必须包含至少一个字母和一个数字",
    "confirm_password": "两次输入的密码不一致"
  }
}
```

### 4. 统一错误处理

使用 `shared.ValidateRequest` 统一处理 API 层验证：

```go
func Handler(c *gin.Context) {
    var req Request
    if !shared.ValidateRequest(c, &req) {
        return // 错误已在 ValidateRequest 中处理
    }
    // 继续处理...
}
```

### 5. 验证器测试

为每个自定义验证器编写测试：

```go
func TestValidateAmount(t *testing.T) {
    tests := []struct {
        name    string
        amount  float64
        wantErr bool
    }{
        {"valid", 100.50, false},
        {"negative", -50.00, true},
        {"too many decimals", 100.505, true},
    }
    // 测试实现...
}
```

## 性能考虑

### 1. 避免重复验证

```go
// 不推荐 - 在多层重复验证
func (h *Handler) Create(c *gin.Context) {
    var req Request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // Service 层又验证了一遍
    h.service.Create(c, req)
}

// 推荐 - 在 API 层验证
func (h *Handler) Create(c *gin.Context) {
    var req Request
    if !shared.ValidateRequest(c, &req) {
        return
    }
    // Service 层直接使用
    h.service.Create(c, req)
}
```

### 2. 使用单例验证器

全局验证器使用 `sync.Once` 确保只初始化一次：

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

### 3. 避免复杂验证

复杂的业务规则验证应该在 Service 层进行，而不是通过 tag：

```go
// 不推荐 - 复杂的 tag 验证
type Request struct {
    Balance float64 `validate:"balance_check"`
}

// 推荐 - Service 层验证
func (s *Service) Withdraw(ctx context.Context, req *WithdrawRequest) error {
    if req.Amount > s.GetBalance(ctx, req.UserID) {
        return errors.New("余额不足")
    }
    // ...
}
```

## 安全考虑

### 1. 输入长度限制

```go
type Request struct {
    Name  string `validate:"max=100"`
    Email string `validate:"max=255"`
    Note  string `validate:"max=5000"`
}
```

### 2. 特殊字符过滤

```go
// 在 Service 层进行更严格的验证
func (v *Validator) validateHTMLInjection(content string) error {
    dangerous := []string{"<script>", "javascript:", "onerror="}
    for _, d := range dangerous {
        if strings.Contains(strings.ToLower(content), d) {
            return errors.New("内容包含危险字符")
        }
    }
    return nil
}
```

### 3. 防止验证绕过

```go
// 不要只依赖前端验证
// 始终在后端进行完整的验证
func (h *Handler) Create(c *gin.Context) {
    var req Request
    // 始终验证，即使前端已验证
    if !shared.ValidateRequest(c, &req) {
        return
    }
    // ...
}
```

## 相关文档

- [验证器使用指南](../../pkg/validator/README.md)
- [API架构文档](../architecture/api_architecture.md)
- [错误处理规范](./error_handling_standard.md)

## 变更历史

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|---------|------|
| 1.0.0 | 2026-02-26 | 初始版本 | Backend Team |

## 维护者

- 创建日期：2026-02-26
- 维护者：Backend Team
- 审核者：架构师
