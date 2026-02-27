# pkg/errors - 统一错误处理系统

## 概述

本包提供了 Qingyu_backend 项目的统一错误处理系统，包括：

- **UnifiedError**: 统一的错误结构体
- **ErrorFactory**: 预定义的错误工厂
- **ErrorCode**: 标准化的错误码定义
- **ErrorConverter**: 旧错误类型转换器
- **HTTP响应构建器**: 标准化的HTTP错误响应

## 目录结构

```
pkg/errors/
├── README.md                   # 本文档
├── unified_error.go            # 统一错误结构体
├── codes.go                    # 错误码定义
├── module_codes.go             # 模块专属错误码
├── error_factory.go            # 错误工厂
├── error_converter.go          # 错误转换器
├── middleware_funcs.go         # 错误处理中间件
├── ai_errors.go               # AI服务错误
├── adapter_errors.go          # 适配器错误
├── layer_errors.go            # 分层错误
├── examples/                  # 使用示例
│   └── user_service_example.go
└── unified_error_system_test.go
```

## 快速开始

### 1. 使用预定义的错误工厂

```go
import "qingyu_backend/pkg/errors"

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*User, error) {
    if userID == "" {
        return nil, errors.UserServiceFactory.ValidationError(
            "1008",
            "用户ID不能为空",
            "field: user_id",
        )
    }

    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        return nil, errors.UserServiceFactory.InternalError(
            "5001",
            "查询用户失败",
            err,
        )
    }

    if user == nil {
        return nil, errors.UserServiceFactory.NotFoundError("用户", userID)
    }

    return user, nil
}
```

### 2. 使用错误构建器

```go
import "qingyu_backend/pkg/errors"

err := errors.NewErrorBuilder().
    WithCode("1001").
    WithCategory(errors.CategoryValidation).
    WithLevel(errors.LevelError).
    WithMessage("参数无效").
    WithDetails("用户名长度必须在3-20个字符之间").
    WithHTTPStatus(400).
    WithService("user-service", "GetUserByID").
    Build()
```

### 3. 转换旧错误类型

```go
import "qingyu_backend/pkg/errors"

func (h *UserHandler) handleError(c *gin.Context, err error) {
    unifiedErr := errors.ToUnifiedError("user-service", err)
    
    statusCode, response := errors.ToHTTPResponse(
        unifiedErr,
        c.GetString("request_id"),
        c.GetString("trace_id"),
    )
    
    c.JSON(statusCode, response)
}
```

## 错误码体系

项目使用4位数字错误码，按类别划分：

### 错误类别

| 类别 | 范围 | 说明 |
|------|------|------|
| 通用客户端错误 | 1000-1099 | 参数验证、格式错误等 |
| 用户相关错误 | 2000-2999 | 认证、授权、用户操作 |
| 业务逻辑错误 | 3000-3999 | 书籍、章节、内容等业务错误 |
| 限流配额错误 | 4000-4099 | 频率限制、配额限制 |
| 服务器内部错误 | 5000-5099 | 系统、数据库、外部服务错误 |

### 模块专属错误码

| 模块 | 范围 | 说明 |
|------|------|------|
| Writer | 3300-3399 | 写作模块专属错误 |
| Reader | 3400-3499 | 阅读模块专属错误 |
| AIService | 3500-3599 | AI服务专属错误 |
| Social | 3600-3699 | 社交功能专属错误 |
| Messaging | 3700-3799 | 消息功能专属错误 |
| Admin | 3800-3899 | 管理功能专属错误 |

详细错误码列表请参考：
- [错误码标准](../../docs/standards/error_code_standard.md)
- [错误处理指南](../../docs/standards/error_handling_guide.md)

## 错误分类

### ErrorCategory

```go
const (
    CategoryValidation ErrorCategory = "validation"  // 验证错误
    CategoryBusiness   ErrorCategory = "business"    // 业务错误
    CategorySystem     ErrorCategory = "system"      // 系统错误
    CategoryExternal   ErrorCategory = "external"    // 外部服务错误
    CategoryNetwork    ErrorCategory = "network"     // 网络错误
    CategoryAuth       ErrorCategory = "auth"        // 认证授权错误
    CategoryDatabase   ErrorCategory = "database"    // 数据库错误
    CategoryCache      ErrorCategory = "cache"       // 缓存错误
)
```

### ErrorLevel

```go
const (
    LevelInfo     ErrorLevel = "info"      // 信息性错误
    LevelWarning  ErrorLevel = "warning"   // 警告性错误
    LevelError    ErrorLevel = "error"     // 错误
    LevelCritical ErrorLevel = "critical"  // 严重错误
)
```

## 预定义错误工厂

项目为每个服务模块预定义了错误工厂：

```go
var (
    AIServiceFactory        = NewErrorFactory("ai-service")
    UserServiceFactory      = NewErrorFactory("user-service")
    DocumentServiceFactory  = NewErrorFactory("document-service")
    ProjectFactory          = NewErrorFactory("project-service")
    BookstoreServiceFactory = NewErrorFactory("bookstore-service")
    ReaderServiceFactory    = NewErrorFactory("reader-service")
    WriterServiceFactory    = NewErrorFactory("writer-service")
)
```

### 工厂方法

每个工厂提供以下方法：

```go
// 验证错误
ValidationError(code, message string, details ...string) *UnifiedError

// 业务错误
BusinessError(code, message string, details ...string) *UnifiedError

// 资源不存在
NotFoundError(resource, id string) *UnifiedError

// 认证错误
AuthError(code, message string) *UnifiedError

// 禁止访问
ForbiddenError(code, message string) *UnifiedError

// 内部错误
InternalError(code, message string, cause error) *UnifiedError

// 外部服务错误
ExternalError(code, message string, retryable bool) *UnifiedError

// 网络错误
NetworkError(message string) *UnifiedError

// 超时错误
TimeoutError(operation string) *UnifiedError

// 频率限制
RateLimitError(limit int) *UnifiedError

// 数据库错误
DatabaseError(operation string, cause error) *UnifiedError

// 缓存错误
CacheError(operation string, cause error) *UnifiedError
```

## HTTP响应

### 标准错误响应格式

```json
{
  "code": "2003",
  "message": "用户名已被使用",
  "details": "用户名 'testuser' 已被使用，请选择其他用户名",
  "timestamp": "2026-02-26T10:30:00Z",
  "request_id": "req-123456",
  "trace_id": "trace-789"
}
```

### 验证错误响应格式

```json
{
  "code": "1001",
  "message": "请求参数无效",
  "details": {
    "fields": [
      {
        "field": "username",
        "message": "用户名长度必须在3-20个字符之间"
      },
      {
        "field": "email",
        "message": "邮箱格式无效"
      }
    ]
  },
  "timestamp": "2026-02-26T10:30:00Z",
  "request_id": "req-123456"
}
```

## 最佳实践

### 1. Service层错误处理

```go
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    // 参数验证
    if req.Username == "" {
        return s.factory.ValidationError(
            "1008",
            "用户名不能为空",
            "field: username",
        )
    }

    // 业务逻辑检查
    exists, err := s.repo.ExistsByUsername(ctx, req.Username)
    if err != nil {
        return s.factory.InternalError("5001", "检查用户名失败", err)
    }
    if exists {
        return s.factory.BusinessError(
            "2003",
            "用户名已被使用",
        )
    }

    // 数据库操作
    if err := s.repo.Create(ctx, user); err != nil {
        return s.factory.InternalError("5001", "创建用户失败", err)
    }

    return nil
}
```

### 2. API层错误处理

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "code":    "1001",
            "message": "请求参数无效",
            "details": err.Error(),
        })
        return
    }

    if err := h.service.CreateUser(c, &req); err != nil {
        statusCode, response := errors.ToHTTPResponseWithError(
            "user-service",
            err,
            c.GetString("request_id"),
            c.GetString("trace_id"),
        )
        c.JSON(statusCode, response)
        return
    }

    c.JSON(200, gin.H{"message": "成功"})
}
```

### 3. Repository层错误处理

```go
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    
    if err == mongo.ErrNoDocuments {
        return nil, nil  // 用户不存在，返回 nil
    }
    
    if err != nil {
        return nil, err  // 返回原始错误，由Service层处理
    }

    return &user, nil
}
```

## 迁移指南

### 从旧错误类型迁移

如果你有旧的错误类型（如 UserError, ReaderError），可以使用错误转换器：

```go
// 旧代码
func (s *UserService) GetUser(id string) (*User, error) {
    user, err := s.repo.Get(id)
    if err != nil {
        return nil, NotFound("用户", id)  // 旧的UserError
    }
    return user, nil
}

// 新代码
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, s.factory.InternalError("5001", "查询用户失败", err)
    }
    if user == nil {
        return nil, s.factory.NotFoundError("用户", id)
    }
    return user, nil
}
```

### 使用转换器兼容旧代码

```go
import "qingyu_backend/pkg/errors"

func (h *Handler) handleError(c *gin.Context, err error) {
    // 转换为统一错误
    unifiedErr := errors.ToUnifiedError("user-service", err)
    
    // 或者使用特定转换器
    converter := errors.NewLegacyErrorConverter("user-service")
    if userErr, ok := err.(*UserError); ok {
        unifiedErr = converter.ConvertFromUserError(userErr)
    }
    
    // 返回HTTP响应
    statusCode, response := errors.ToHTTPResponse(unifiedErr, ...)
    c.JSON(statusCode, response)
}
```

## 错误日志

建议在Service层记录错误日志：

```go
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*User, error) {
    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        s.logger.Error("查询用户失败",
            "user_id", userID,
            "error", err,
        )
        return nil, s.factory.InternalError("5001", "查询用户失败", err)
    }

    if user == nil {
        s.logger.Info("用户不存在", "user_id", userID)
        return nil, s.factory.NotFoundError("用户", userID)
    }

    return user, nil
}
```

## 测试

```go
func TestUserService_CreateUser(t *testing.T) {
    service := NewUserService(mockRepo)

    tests := []struct {
        name    string
        req     *CreateUserRequest
        wantErr error
    }{
        {
            name: "用户名为空",
            req:  &CreateUserRequest{Username: ""},
            wantErr: errors.UserServiceFactory.ValidationError(
                "1008",
                "用户名不能为空",
                "",
            ),
        },
        // ... 更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := service.CreateUser(context.Background(), tt.req)
            
            if tt.wantErr != nil {
                assert.Error(t, err)
                if unifiedErr, ok := err.(*errors.UnifiedError); ok {
                    assert.Equal(t, tt.wantErr.Code, unifiedErr.Code)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 相关文档

- [错误码标准](../../docs/standards/error_code_standard.md)
- [错误处理指南](../../docs/standards/error_handling_guide.md)
- [API设计规范](../../docs/standards/api/API设计规范.md)

## 维护者

如有问题或建议，请联系架构团队。
