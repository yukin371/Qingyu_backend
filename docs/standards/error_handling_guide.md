# 错误处理指南

**版本**: v1.0  
**创建日期**: 2026-02-26  
**状态**: ✅ 正式实施  

---

## 一、概述

本指南定义了 Qingyu_backend 项目的错误处理规范，包括错误码定义、错误创建、错误传播和错误响应格式。

---

## 二、错误码标准

### 2.1 错误码格式

项目使用 **4位数字错误码**，格式为 `ABCD`：

```
A - 错误类别 (1-5)
B - 子类别 (0-9)
C-D - 具体错误编号 (00-99)
```

### 2.2 错误类别分类

| 类别 | 范围 | 说明 | HTTP状态码范围 |
|------|------|------|----------------|
| 通用客户端错误 | 1000-1099 | 参数验证、格式错误等 | 400-499 |
| 用户相关错误 | 2000-2999 | 认证、授权、用户操作 | 401-403, 409 |
| 业务逻辑错误 | 3000-3999 | 书籍、章节、内容等业务错误 | 400, 404, 409 |
| 限流配额错误 | 4000-4099 | 频率限制、配额限制 | 429 |
| 服务器内部错误 | 5000-5099 | 系统、数据库、外部服务错误 | 500-503 |

### 2.3 常用错误码列表

#### 通用客户端错误 (1000-1099)

| 错误码 | 名称 | HTTP状态码 | 说明 |
|--------|------|-----------|------|
| 1001 | InvalidParams | 400 | 请求参数无效 |
| 1008 | MissingParam | 400 | 缺少必填参数 |
| 1009 | InvalidFormat | 400 | 参数格式无效 |
| 1015 | ValidationFailed | 400 | 验证失败 |
| 1002 | Unauthorized | 401 | 未授权访问 |
| 1003 | Forbidden | 403 | 禁止访问 |
| 1004 | NotFound | 404 | 资源不存在 |
| 1005 | AlreadyExists | 409 | 资源已存在 |

#### 用户相关错误 (2000-2999)

| 错误码 | 名称 | HTTP状态码 | 说明 |
|--------|------|-----------|------|
| 2001 | UserNotFound | 404 | 用户不存在 |
| 2002 | InvalidCredentials | 401 | 无效凭证 |
| 2003 | UsernameAlreadyUsed | 409 | 用户名已被使用 |
| 2004 | EmailAlreadyUsed | 409 | 邮箱已被使用 |
| 2008 | TokenExpired | 401 | Token过期 |
| 2009 | TokenInvalid | 401 | Token无效 |

#### 业务逻辑错误 (3000-3999)

| 错误码 | 名称 | HTTP状态码 | 说明 |
|--------|------|-----------|------|
| 3001 | BookNotFound | 404 | 书籍不存在 |
| 3002 | ChapterNotFound | 404 | 章节不存在 |
| 3010 | InsufficientQuota | 400 | 配额不足 |
| 3020 | CharacterNotFound | 404 | 角色不存在 |

#### 限流配额错误 (4000-4099)

| 错误码 | 名称 | HTTP状态码 | 说明 |
|--------|------|-----------|------|
| 4000 | RateLimitExceeded | 429 | 频率限制超出 |
| 4001 | DailyLimitExceeded | 429 | 每日限制超出 |

#### 服务器内部错误 (5000-5099)

| 错误码 | 名称 | HTTP状态码 | 说明 |
|--------|------|-----------|------|
| 5000 | InternalError | 500 | 内部错误 |
| 5001 | DatabaseError | 500 | 数据库错误 |
| 5004 | ExternalAPIError | 502 | 外部API错误 |

---

## 三、统一错误系统

### 3.1 使用 UnifiedError

项目统一使用 `pkg/errors.UnifiedError` 作为错误结构：

```go
import "qingyu_backend/pkg/errors"

// 创建错误
err := errors.NewErrorBuilder().
    WithCode("1001").
    WithCategory(errors.CategoryValidation).
    WithLevel(errors.LevelError).
    WithMessage("参数无效").
    WithDetails("用户名长度必须在3-20个字符之间").
    WithHTTPStatus(400).
    Build()
```

### 3.2 使用错误工厂

推荐使用预定义的错误工厂：

```go
import "qingyu_backend/pkg/errors"

// 使用预定义工厂
factory := errors.UserServiceFactory

// 验证错误
err := factory.ValidationError(
    "1001",
    "用户名格式无效",
    "用户名长度必须在3-20个字符之间",
)

// 业务错误
err := factory.BusinessError(
    "2003",
    "用户名已被使用",
)

// 不存在错误
err := factory.NotFoundError("用户", userID)

// 认证错误
err := factory.AuthError("2008", "登录已过期")

// 内部错误
err := factory.InternalError("5000", "数据库查询失败", causeErr)
```

---

## 四、错误处理最佳实践

### 4.1 Service层错误处理

```go
// service/user/user_service.go

import (
    "qingyu_backend/pkg/errors"
)

type UserService struct {
    factory *errors.ErrorFactory
}

func NewUserService() *UserService {
    return &UserService{
        factory: errors.UserServiceFactory,
    }
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*User, error) {
    // 参数验证
    if userID == "" {
        return nil, s.factory.ValidationError(
            "1008",
            "用户ID不能为空",
        )
    }

    // 查询用户
    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        return nil, s.factory.InternalError("5001", "查询用户失败", err)
    }

    // 检查用户是否存在
    if user == nil {
        return nil, s.factory.NotFoundError("用户", userID)
    }

    return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    // 验证用户名格式
    if len(req.Username) < 3 || len(req.Username) > 20 {
        return s.factory.ValidationError(
            "1001",
            "用户名长度必须在3-20个字符之间",
            "username: " + req.Username,
        )
    }

    // 检查用户名是否存在
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

    // 创建用户
    if err := s.repo.Create(ctx, user); err != nil {
        return s.factory.InternalError("5001", "创建用户失败", err)
    }

    return nil
}
```

### 4.2 API层错误处理

```go
// api/v1/user/user_api.go

import (
    "qingyu_backend/pkg/errors"
)

func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")

    user, err := h.service.GetUserByID(c, userID)
    if err != nil {
        // 统一错误处理
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
    // 类型断言，检查是否是 UnifiedError
    if unifiedErr, ok := err.(*errors.UnifiedError); ok {
        c.JSON(unifiedErr.GetHTTPStatus(), gin.H{
            "code":    unifiedErr.Code,
            "message": unifiedErr.Message,
            "details": unifiedErr.Details,
        })
        return
    }

    // 其他错误类型，返回内部错误
    c.JSON(http.StatusInternalServerError, gin.H{
        "code":    "5000",
        "message": "内部错误",
    })
}
```

### 4.3 Repository层错误处理

```go
// repository/user_repository.go

import (
    "qingyu_backend/pkg/errors"
)

func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    
    if err == mongo.ErrNoDocuments {
        // 用户不存在，返回 nil 而不是错误
        return nil, nil
    }
    
    if err != nil {
        // 数据库错误
        return nil, &errors.RepositoryError{
            Code:    "DATABASE_ERROR",
            Message: "查询用户失败",
            Err:     err,
        }
    }

    return &user, nil
}
```

---

## 五、错误传播规则

### 5.1 错误传播层级

```
┌─────────────┐
│   API层     │ → 统一错误响应格式
├─────────────┤
│  Service层  │ → 使用UnifiedError，添加业务上下文
├─────────────┤
│ Repository层│ → 返回RepositoryError或nil
└─────────────┘
```

### 5.2 错误转换规则

| 层级 | 接收的错误类型 | 返回的错误类型 |
|------|---------------|----------------|
| Repository | MongoDB错误 | RepositoryError |
| Service | RepositoryError | UnifiedError |
| API | UnifiedError | HTTP响应 |

---

## 六、错误响应格式

### 6.1 标准错误响应

```json
{
  "code": "2003",
  "message": "用户名已被使用",
  "details": "用户名 'testuser' 已被使用，请选择其他用户名",
  "timestamp": "2026-02-26T10:30:00Z",
  "request_id": "req-123456"
}
```

### 6.2 验证错误响应

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

### 6.3 内部错误响应

```json
{
  "code": "5000",
  "message": "服务器内部错误",
  "timestamp": "2026-02-26T10:30:00Z",
  "request_id": "req-123456"
}
```

**注意**：内部错误不要暴露具体错误详情给客户端

---

## 七、错误日志记录

### 7.1 日志级别

| 错误级别 | 说明 | 示例 |
|---------|------|------|
| LevelInfo | 信息性错误 | 用户不存在（可能的情况） |
| LevelWarning | 警告性错误 | 外部服务降级 |
| LevelError | 错误 | 业务逻辑错误 |
| LevelCritical | 严重错误 | 数据库连接失败 |

### 7.2 日志记录示例

```go
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*User, error) {
    user, err := s.repo.GetByID(ctx, userID)
    if err != nil {
        // 记录错误日志
        s.logger.Error("查询用户失败",
            "user_id", userID,
            "error", err,
        )
        return nil, s.factory.InternalError("5001", "查询用户失败", err)
    }

    if user == nil {
        // 记录信息日志
        s.logger.Info("用户不存在", "user_id", userID)
        return nil, s.factory.NotFoundError("用户", userID)
    }

    return user, nil
}
```

---

## 八、常见场景示例

### 8.1 参数验证

```go
func ValidateCreateUserRequest(req *CreateUserRequest) error {
    factory := errors.UserServiceFactory

    if req.Username == "" {
        return factory.ValidationError(
            "1008",
            "用户名不能为空",
            "field: username",
        )
    }

    if len(req.Username) < 3 || len(req.Username) > 20 {
        return factory.ValidationError(
            "1001",
            "用户名长度无效",
            "username长度必须在3-20个字符之间",
        )
    }

    if !emailRegex.MatchString(req.Email) {
        return factory.ValidationError(
            "1009",
            "邮箱格式无效",
            "email: " + req.Email,
        )
    }

    return nil
}
```

### 8.2 资源不存在

```go
func (s *BookService) GetBook(ctx context.Context, bookID string) (*Book, error) {
    book, err := s.repo.GetByID(ctx, bookID)
    if err != nil {
        return nil, s.factory.InternalError("5001", "查询书籍失败", err)
    }

    if book == nil {
        return nil, s.factory.NotFoundError("书籍", bookID)
    }

    return book, nil
}
```

### 8.3 权限检查

```go
func (s *BookService) DeleteBook(ctx context.Context, userID, bookID string) error {
    book, err := s.GetBook(ctx, bookID)
    if err != nil {
        return err
    }

    if book.AuthorID != userID {
        return s.factory.ForbiddenError(
            "1003",
            "只有作者才能删除书籍",
        )
    }

    return s.repo.Delete(ctx, bookID)
}
```

### 8.4 外部服务调用

```go
func (s *AIService) GenerateContent(ctx context.Context, req *GenerateRequest) (*Content, error) {
    // 调用外部AI服务
    response, err := s.aiClient.Generate(ctx, req)
    if err != nil {
        // 判断是否可重试
        if isRetryable(err) {
            return nil, s.factory.ExternalError(
                "5004",
                "AI服务暂时不可用",
                true, // retryable
            )
        }
        return nil, s.factory.ExternalError(
            "5004",
            "AI服务调用失败",
            false,
        )
    }

    return response, nil
}
```

---

## 九、测试规范

### 9.1 错误创建测试

```go
func TestUserService_CreateUser(t *testing.T) {
    service := NewUserService()

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
        {
            name: "用户名过短",
            req:  &CreateUserRequest{Username: "ab"},
            wantErr: errors.UserServiceFactory.ValidationError(
                "1001",
                "用户名长度无效",
                "",
            ),
        },
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

---

## 十、检查清单

在代码审查时，请检查以下项目：

- [ ] 使用统一的错误码格式
- [ ] 使用预定义的错误工厂
- [ ] 错误消息清晰明确
- [ ] 内部错误不暴露敏感信息
- [ ] 错误日志记录完整
- [ ] 可重试错误正确标记
- [ ] HTTP状态码映射正确

---

## 十一、相关文档

- [错误码标准](./error_code_standard.md)
- [pkg/errors包文档](../../pkg/errors/README.md)
- [API设计规范](./api/API设计规范.md)
