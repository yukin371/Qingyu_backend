# Service层测试规范

## 概述

Service层负责业务逻辑处理，位于Repository层和API层之间。本规范定义了Service层测试的详细要求和最佳实践。

## 核心原则

### ✅ 必须使用Mock Repository

```go
// ✅ 正确示例
func TestAuthService_Login(t *testing.T) {
    // Setup - 创建Mock Repository
    mockRepo := new(MockAuthRepository)
    mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(&users.User{
            ID:       "user123",
            Username: "testuser",
            Password: "$2a$10$hashedpassword", // bcrypt hash
        }, nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act - 执行业务逻辑
    token, err := service.Login(ctx, "test@example.com", "password123")

    // Assert - 验证结果
    require.NoError(t, err)
    assert.NotEmpty(t, token)
    mockRepo.AssertExpectations(t)
}
```

### ❌ 严格禁止

```go
// ❌ 错误：Service层测试中使用真实Repository
func TestUserService_CreateUser(t *testing.T) {
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := NewUserRepository(db)  // 危险！这变成了集成测试
    service := NewUserService(repo)

    // 问题：测试失败时无法确定是Service逻辑问题还是Repository问题
}
```

## 测试组织结构

### 文件位置

```
service/{module}/*_test.go
```

示例：
```
service/shared/auth/jwt_service_test.go
service/bookstore/bookstore_service_test.go
service/reader/reading_progress_service_test.go
```

### 包命名

与被测试代码使用相同的包名：

```go
package auth

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)
```

## Mock Repository创建

### 定义Mock对象

```go
// service/shared/auth/auth_service_test.go
package auth

import "github.com/stretchr/testify/mock"

// MockAuthRepository Mock仓储接口
type MockAuthRepository struct {
    mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *users.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockAuthRepository) GetUserByID(ctx context.Context, id string) (*users.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*users.User), args.Error(1)
}
```

### Mock使用最佳实践

```go
func TestAuthService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)

    // 设置期望 - 成功场景
    mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *users.User) bool {
        return u.Username == "testuser" && u.Email == "test@example.com"
    })).Return(nil)

    // 设置期望 - 邮箱已存在场景
    mockRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").
        Return(&users.User{ID: "existing"}, nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act
    user := &users.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    err := service.CreateUser(ctx, user)

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, user.ID)
    mockRepo.AssertExpectations(t)
}
```

## AAA测试模式

所有Service测试必须遵循AAA模式（Arrange-Act-Assert）：

### Arrange（准备）

```go
// Arrange - 准备Mock和测试数据
mockRepo := new(MockAuthRepository)
mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
    Return(&users.User{
        ID:       "user123",
        Username: "testuser",
        Password: "$2a$10$hashedpassword",
    }, nil)

service := auth.NewAuthService(mockRepo)
ctx := context.Background()
```

### Act（执行）

```go
// Act - 执行业务逻辑
token, err := service.Login(ctx, "test@example.com", "password123")
```

### Assert（断言）

```go
// Assert - 验证结果
require.NoError(t, err)
assert.NotEmpty(t, token)

// 验证Mock调用
mockRepo.AssertCalled(t, "GetUserByEmail", mock.Anything, "test@example.com")
mockRepo.AssertNumberOfCalls(t, "GetUserByEmail", 1)
```

## 测试用例设计

### 1. 业务逻辑测试

#### 成功场景测试

```go
func TestAuthService_Login_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    mockUser := &users.User{
        ID:       "user123",
        Username: "testuser",
        Email:    "test@example.com",
        Password: string(hashedPassword),
    }

    mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(mockUser, nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act
    token, err := service.Login(ctx, "test@example.com", "password123")

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, token)

    // 验证token格式
    parts := strings.Split(token, ".")
    assert.Len(t, parts, 3, "JWT应该有3个部分")

    mockRepo.AssertExpectations(t)
}
```

#### 失败场景测试

```go
func TestAuthService_Login_UserNotFound(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)
    mockRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").
        Return(nil, errors.New("用户不存在"))

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act
    token, err := service.Login(ctx, "notfound@example.com", "password123")

    // Assert
    require.Error(t, err)
    assert.Empty(t, token)
    assert.Contains(t, err.Error(), "用户不存在")

    mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    mockUser := &users.User{
        ID:       "user123",
        Email:    "test@example.com",
        Password: string(hashedPassword),
    }

    mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(mockUser, nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act - 使用错误密码
    token, err := service.Login(ctx, "test@example.com", "wrongpassword")

    // Assert
    require.Error(t, err)
    assert.Empty(t, token)
    assert.Contains(t, err.Error(), "密码错误")

    mockRepo.AssertExpectations(t)
}
```

### 2. 边缘案例测试

```go
func TestAuthService_CreateUser_EmptyUsername(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)
    mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act - 用户名为空
    user := &users.User{
        Username: "",
        Email:    "test@example.com",
        Password: "password123",
    }
    err := service.CreateUser(ctx, user)

    // Assert
    require.Error(t, err)
    assert.Contains(t, err.Error(), "用户名不能为空")

    // 验证不应该调用CreateUser
    mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

func TestAuthService_CreateUser_DuplicateEmail(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)

    // 模拟邮箱已存在
    mockRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").
        Return(&users.User{ID: "existing_user"}, nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act - 创建重复邮箱的用户
    user := &users.User{
        Username: "newuser",
        Email:    "existing@example.com",
        Password: "password123",
    }
    err := service.CreateUser(ctx, user)

    // Assert
    require.Error(t, err)
    assert.Contains(t, err.Error(), "邮箱已被使用")

    // 验证不会尝试创建用户
    mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
}
```

### 3. 业务规则测试

```go
func TestAuthService_DeleteUser_CannotDeleteSystemUser(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)

    // 模拟系统用户
    systemUser := &users.User{
        ID:       "system_user",
        Username: "system",
        Email:    "system@qingyu.com",
        IsSystem: true,
    }

    mockRepo.On("GetUserByID", mock.Anything, "system_user").
        Return(systemUser, nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act - 尝试删除系统用户
    err := service.DeleteUser(ctx, "system_user")

    // Assert
    require.Error(t, err)
    assert.Contains(t, err.Error(), "不能删除系统用户")

    // 验证不会调用删除
    mockRepo.AssertNotCalled(t, "DeleteUser", mock.Anything, mock.Anything)
}

func TestBookstoreService_PurchaseChapter_InsufficientBalance(t *testing.T) {
    // Arrange
    mockUserRepo := new(MockUserRepository)
    mockBookRepo := new(MockBookRepository)
    mockTxRepo := new(MockTransactionRepository)

    // 用户余额不足
    mockUserRepo.On("GetUserByID", mock.Anything, "user123").
        Return(&users.User{
            ID:      "user123",
            Balance: 50, // 余额50元
        }, nil)

    // 章节价格100元
    mockBookRepo.On("GetChapterByID", mock.Anything, "chapter123").
        Return(&bookstore.Chapter{
            ID:    "chapter123",
            Price: 100,
        }, nil)

    service := bookstore.NewPurchaseService(mockUserRepo, mockBookRepo, mockTxRepo)
    ctx := context.Background()

    // Act - 尝试购买
    err := service.PurchaseChapter(ctx, "user123", "chapter123")

    // Assert
    require.Error(t, err)
    assert.Contains(t, err.Error(), "余额不足")

    // 验证不会创建交易
    mockTxRepo.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything)
}
```

### 4. 复杂业务流程测试

```go
func TestBookstoreService_PurchaseBook_SuccessfulFlow(t *testing.T) {
    // Arrange
    mockUserRepo := new(MockUserRepository)
    mockBookRepo := new(MockBookRepository)
    mockTxRepo := new(MockTransactionRepository)
    mockPurchaseRepo := new(MockPurchaseRepository)

    // 用户有足够余额
    user := &users.User{
        ID:      "user123",
        Balance: 1000,
    }
    mockUserRepo.On("GetUserByID", mock.Anything, "user123").Return(user, nil)

    // 书籍信息
    book := &bookstore.Book{
        ID:    "book123",
        Title: "测试书籍",
        Price: 100,
    }
    mockBookRepo.On("GetBookByID", mock.Anything, "book123").Return(book, nil)

    // 创建交易记录
    mockTxRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx *Transaction) bool {
        return tx.UserID == "user123" && tx.Amount == 100
    })).Return(nil)

    // 更新用户余额
    mockUserRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *users.User) bool {
        return u.ID == "user123" && u.Balance == 900
    })).Return(nil)

    // 记录购买
    mockPurchaseRepo.On("CreatePurchase", mock.Anything, mock.Anything).Return(nil)

    service := bookstore.NewPurchaseService(mockUserRepo, mockBookRepo, mockTxRepo, mockPurchaseRepo)
    ctx := context.Background()

    // Act
    purchaseID, err := service.PurchaseBook(ctx, "user123", "book123")

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, purchaseID)

    // 验证完整的调用链
    mockUserRepo.AssertNumberOfCalls(t, "GetUserByID", 1)
    mockBookRepo.AssertNumberOfCalls(t, "GetBookByID", 1)
    mockTxRepo.AssertNumberOfCalls(t, "CreateTransaction", 1)
    mockUserRepo.AssertNumberOfCalls(t, "UpdateUser", 1)
    mockPurchaseRepo.AssertNumberOfCalls(t, "CreatePurchase", 1)
}
```

## Table-Driven测试

对于有多个输入场景的测试，使用表格驱动模式：

```go
func TestAuthService_ValidateEmail_TableDriven(t *testing.T) {
    service := auth.NewAuthService(nil)

    tests := []struct {
        name      string
        email     string
        wantErr   bool
        errMsg    string
    }{
        {
            name:    "有效邮箱",
            email:   "test@example.com",
            wantErr: false,
        },
        {
            name:    "空邮箱",
            email:   "",
            wantErr: true,
            errMsg:  "邮箱不能为空",
        },
        {
            name:    "缺少@符号",
            email:   "invalidemail.com",
            wantErr: true,
            errMsg:  "邮箱格式无效",
        },
        {
            name:    "缺少域名",
            email:   "test@",
            wantErr: true,
            errMsg:  "邮箱格式无效",
        },
        {
            name:    "缺少用户名",
            email:   "@example.com",
            wantErr: true,
            errMsg:  "邮箱格式无效",
        },
        {
            name:    "多个@符号",
            email:   "test@test@example.com",
            wantErr: true,
            errMsg:  "邮箱格式无效",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            err := service.ValidateEmail(tt.email)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                if tt.errMsg != "" {
                    assert.Contains(t, err.Error(), tt.errMsg)
                }
            } else {
                require.NoError(t, err)
            }
        })
    }
}

func TestUserService_CreateUser_TableDriven(t *testing.T) {
    tests := []struct {
        name      string
        user      *users.User
        setupMock func(*MockAuthRepository)
        wantErr   bool
        errMsg    string
    }{
        {
            name: "成功创建用户",
            user: &users.User{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            setupMock: func(m *MockAuthRepository) {
                m.On("GetUserByEmail", mock.Anything, "test@example.com").
                    Return(nil, nil)
                m.On("CreateUser", mock.Anything, mock.Anything).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name: "邮箱已存在",
            user: &users.User{
                Username: "testuser",
                Email:    "existing@example.com",
                Password: "password123",
            },
            setupMock: func(m *MockAuthRepository) {
                m.On("GetUserByEmail", mock.Anything, "existing@example.com").
                    Return(&users.User{ID: "existing"}, nil)
            },
            wantErr: true,
            errMsg:  "邮箱已被使用",
        },
        {
            name: "用户名为空",
            user: &users.User{
                Username: "",
                Email:    "test@example.com",
                Password: "password123",
            },
            setupMock: func(m *MockAuthRepository) {},
            wantErr:   true,
            errMsg:    "用户名不能为空",
        },
        {
            name: "密码太短",
            user: &users.User{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "12345",
            },
            setupMock: func(m *MockAuthRepository) {},
            wantErr:   true,
            errMsg:    "密码长度至少8位",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(MockAuthRepository)
            tt.setupMock(mockRepo)

            service := auth.NewAuthService(mockRepo)
            ctx := context.Background()

            // Act
            err := service.CreateUser(ctx, tt.user)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
                assert.NotEmpty(t, tt.user.ID)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

## 错误处理测试

### Repository错误传播

```go
func TestAuthService_GetUser_RepositoryError(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)
    mockRepo.On("GetUserByID", mock.Anything, "user123").
        Return(nil, errors.New("database connection failed"))

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act
    user, err := service.GetUser(ctx, "user123")

    // Assert
    require.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "database connection failed")

    mockRepo.AssertExpectations(t)
}
```

### 超时处理

```go
func TestSlowOperation_Timeout(t *testing.T) {
    // Arrange
    mockRepo := new(MockSlowRepository)
    mockRepo.On("SlowQuery", mock.Anything).
        WaitUntil(time.After(5 * time.Second)).
        Return("result", nil)

    service := NewSlowService(mockRepo)
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    // Act
    result, err := service.SlowOperation(ctx)

    // Assert
    require.Error(t, err)
    assert.Empty(t, result)
    assert.Equal(t, context.DeadlineExceeded, err)
}
```

## 测试覆盖率目标

### 必须达到的覆盖率

| 方法类型 | 覆盖率目标 | 说明 |
|---------|-----------|------|
| 公开方法 | 100% | 对外接口必须完整测试 |
| 私有方法 | ≥80% | 重要私有方法需要测试 |
| 错误分支 | 100% | 所有错误路径必须覆盖 |
| 业务规则 | 100% | 核心业务规则必须完整覆盖 |
| 边缘案例 | ≥90% | 边界条件充分测试 |

### 覆盖率检查命令

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./service/shared/auth

# 查看覆盖率
go tool cover -func=coverage.out | grep "auth_service.go"

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html
```

## 并发测试

### 测试并发安全性

```go
func TestAuthService_CreateUser_Concurrent(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuthRepository)
    mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    // Act - 并发创建用户
    const concurrency = 10
    var wg sync.WaitGroup
    errors := make(chan error, concurrency)

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            user := &users.User{
                Username: fmt.Sprintf("concurrent_user_%d", index),
                Email:    fmt.Sprintf("concurrent_%d@example.com", index),
                Password: "password123",
            }
            if err := service.CreateUser(ctx, user); err != nil {
                errors <- err
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Assert - 验证没有错误
    for err := range errors {
        t.Errorf("并发创建用户失败: %v", err)
    }
}
```

### 竞态检测

```bash
# 使用-race标志检测竞态条件
go test -race ./service/shared/auth
```

## 性能测试

### Benchmark测试

```go
func BenchmarkAuthService_Login(b *testing.B) {
    mockRepo := new(MockAuthRepository)
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

    mockUser := &users.User{
        ID:       "user123",
        Email:    "test@example.com",
        Password: string(hashedPassword),
    }

    mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").
        Return(mockUser, nil)

    service := auth.NewAuthService(mockRepo)
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.Login(ctx, "test@example.com", "password123")
    }
}

func BenchmarkAuthService_VerifyToken(b *testing.B) {
    mockRepo := new(MockAuthRepository)
    service := auth.NewAuthService(mockRepo)

    // 生成测试token
    token, _ := service.GenerateToken(context.Background(), "user123")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.VerifyToken(context.Background(), token)
    }
}
```

### 运行性能测试

```bash
# 运行Service层的性能测试
go test -bench=. -benchmem ./service/shared/auth

# 运行特定的benchmark测试
go test -bench=BenchmarkAuthService -benchmem ./service/shared/auth
```

## Mock高级用法

### 参数匹配器

```go
func TestUserService_UpdateUser_AdvancedMatching(t *testing.T) {
    mockRepo := new(MockUserRepository)

    // 使用MatchedBy进行复杂匹配
    mockRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *users.User) bool {
        return u.ID == "user123" &&
               u.Username == "updated_user" &&
               len(u.Roles) > 0
    })).Return(nil)

    service := NewUserService(mockRepo)
    ctx := context.Background()

    user := &users.User{
        ID:       "user123",
        Username: "updated_user",
        Roles:    []string{"user"},
    }
    err := service.UpdateUser(ctx, user)

    require.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### 多次调用不同返回值

```go
func TestRetryLogic(t *testing.T) {
    mockRepo := new(MockRepository)

    // 前两次失败，第三次成功
    mockRepo.On("GetData", mock.Anything).
        Return(nil, errors.New("temporary error")).Once()
    mockRepo.On("GetData", mock.Anything).
        Return(nil, errors.New("temporary error")).Once()
    mockRepo.On("GetData", mock.Anything).
        Return("success", nil)

    service := NewServiceWithRetry(mockRepo, 3)
    result, err := service.GetDataWithRetry(context.Background())

    require.NoError(t, err)
    assert.Equal(t, "success", result)
    mockRepo.AssertNumberOfCalls(t, "GetData", 3)
}
```

### Run回调

```go
func TestUserService_CreateUser_GeneratesID(t *testing.T) {
    mockRepo := new(MockUserRepository)

    // 使用Run修改传入的参数
    mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *users.User) bool {
        return u.ID == "" // 创建时没有ID
    })).Return(nil).Run(func(args mock.Arguments) {
        user := args.Get(1).(*users.User)
        user.ID = "generated_id_123" // 模拟数据库生成ID
    })

    service := NewUserService(mockRepo)
    ctx := context.Background()

    user := &users.User{
        Username: "testuser",
        Email:    "test@example.com",
    }
    err := service.CreateUser(ctx, user)

    require.NoError(t, err)
    assert.Equal(t, "generated_id_123", user.ID)
}
```

## 常见问题

### Q1: 为什么Service测试必须使用Mock Repository？

**A**: Service层的职责是业务逻辑处理，使用Mock Repository可以：
- **隔离测试**：只测试Service层的业务逻辑
- **快速执行**：避免数据库操作，测试运行更快
- **可控环境**：精确控制Repository的返回值和错误
- **明确职责**：Repository层由自己的测试覆盖

### Q2: 何时应该测试私有方法？

**A**: 私有方法测试策略：
- **不直接测试**：通过公开方法间接测试私有方法
- **复杂逻辑**：如果私有方法包含复杂业务逻辑，值得测试
- **提取公开**：考虑将复杂的私有方法提取为独立的公开Service

```go
// ✅ 好：通过公开方法测试私有方法
func TestUserService_CalculateDiscount_Indirect(t *testing.T) {
    // 通过公开的Purchase方法间接测试CalculateDiscount
    service := setupService()
    result, err := service.Purchase(ctx, userID, bookID)
    // 验证折扣计算结果
}

// ✅ 也可以：提取为独立的Service
type DiscountService struct{}

func (s *DiscountService) CalculateDiscount(user *User, book *Book) float64 {
    // 独立的折扣计算逻辑
}
```

### Q3: 如何测试事务处理？

**A**: 事务处理测试策略：
```go
func TestPurchaseService_PurchaseBook_Transaction(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    mockBookRepo := new(MockBookRepository)
    mockTxRepo := new(MockTransactionRepository)

    // 模拟事务成功
    mockTxRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
    mockTxRepo.On("CreateTransaction", mock.Anything, mock.Anything).Return(nil)
    mockUserRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
    mockTxRepo.On("Commit", mock.Anything).Return(nil)

    service := NewPurchaseService(mockUserRepo, mockBookRepo, mockTxRepo)
    ctx := context.Background()

    err := service.PurchaseBook(ctx, userID, bookID)

    require.NoError(t, err)

    // 验证事务完整流程
    mockTxRepo.AssertCalled(t, "BeginTx")
    mockTxRepo.AssertCalled(t, "CreateTransaction")
    mockTxRepo.AssertCalled(t, "Commit")
}

func TestPurchaseService_PurchaseBook_TransactionRollback(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    mockBookRepo := new(MockBookRepository)
    mockTxRepo := new(MockTransactionRepository)

    // 模拟事务失败
    mockTxRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
    mockTxRepo.On("CreateTransaction", mock.Anything, mock.Anything).
        Return(errors.New("transaction failed"))
    mockTxRepo.On("Rollback", mock.Anything).Return(nil)

    service := NewPurchaseService(mockUserRepo, mockBookRepo, mockTxRepo)
    ctx := context.Background()

    err := service.PurchaseBook(ctx, userID, bookID)

    require.Error(t, err)

    // 验证回滚
    mockTxRepo.AssertCalled(t, "Rollback")
    mockTxRepo.AssertNotCalled(t, "Commit")
}
```

### Q4: 如何测试第三方依赖？

**A**: 第三方依赖应该通过接口抽象和Mock：
```go
// 定义接口
type EmailService interface {
    SendEmail(ctx context.Context, to, subject, body string) error
}

// Service使用接口
type UserService struct {
    userRepo  UserRepository
    emailSvc  EmailService
}

// Mock
type MockEmailService struct {
    mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
    args := m.Called(ctx, to, subject, body)
    return args.Error(0)
}

// 测试
func TestUserService_SendWelcomeEmail(t *testing.T) {
    mockUserRepo := new(MockUserRepository)
    mockEmailSvc := new(MockEmailService)

    mockEmailSvc.On("SendEmail", mock.Anything, "test@example.com", "欢迎", mock.Anything).
        Return(nil)

    service := NewUserService(mockUserRepo, mockEmailSvc)
    // ...
}
```

## 参考文档

- [testify使用指南](../03_测试工具指南/testify使用指南.md)
- [Mock框架使用指南](../03_测试工具指南/mock框架使用指南.md)
- [Repository层测试规范](./repository_层测试规范.md)
- [API层测试规范](./api_层测试规范.md)
