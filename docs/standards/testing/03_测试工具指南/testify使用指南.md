# testify使用指南

## 概述

testify是Go最流行的测试框架，提供强大的断言库和Mock功能。

## 安装

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/suite
```

## 断言库（assert vs require）

### assert - 失败后继续执行

```go
func TestExample(t *testing.T) {
    // 即使失败，也会继续执行后续断言
    assert.Equal(t, 1, 2, "应该相等")
    assert.NotNil(t, nil, "不应该为nil")  // 仍会执行
}
```

**使用场景**：
- 验证多个独立的条件
- 想要看到所有失败的断言
- 非关键性的验证

### require - 失败后立即停止

```go
func TestExample(t *testing.T) {
    // 失败后立即停止测试
    require.NoError(t, setup(), "Setup失败")  // 如果失败，后面的代码不执行
    require.NotNil(t, user, "用户不应该为nil")

    // 关键验证
    assert.Equal(t, "admin", user.Role)
}
```

**使用场景**：
- Setup阶段的关键验证
- 后续测试依赖的前提条件
- 关键性验证

## 常用断言方法

### 相等性断言

```go
// Equal - 相等
assert.Equal(t, expected, actual, "消息")
assert.Equal(t, 100, result)

// NotEqual - 不相等
assert.NotEqual(t, unexpected, actual)

// Exactly - 完全相等（包括类型）
assert.Exactly(t, int32(100), int32(100))

// NotNil - 不为nil
assert.NotNil(t, actual)

// Nil - 为nil
assert.Nil(t, err)
```

### 布尔断言

```go
// True - 为true
assert.True(t, isValid, "应该有效")

// False - 为false
assert.False(t, hasError, "不应该有错误")
```

### 数值断言

```go
// Greater - 大于
assert.Greater(t, 10, 5)

// GreaterOrEqual - 大于等于
assert.GreaterOrEqual(t, 10, 10)

// Less - 小于
assert.Less(t, 5, 10)

// LessOrEqual - 小于等于
assert.LessOrEqual(t, 5, 5)

// InDelta - 浮点数比较（允许误差）
assert.InDelta(t, 3.14, result, 0.01)

// InEpsilon - 相对误差
assert.InEpsilon(t, 100.0, 101.0, 0.01)  // 允许1%误差
```

### 集合断言

```go
// Contains - 包含元素
assert.Contains(t, []int{1, 2, 3}, 2)
assert.Contains(t, "hello world", "world")

// NotContains - 不包含
assert.NotContains(t, []int{1, 2, 3}, 4)

// Len - 长度
assert.Len(t, []int{1, 2, 3}, 3)

// Empty - 为空
assert.Empty(t, []int{})

// ElementsMatch - 元素匹配（忽略顺序）
assert.ElementsMatch(t, []int{1, 2, 3}, []int{3, 2, 1})

// Subset - 子集
assert.Subset(t, []int{1, 2, 3, 4}, []int{2, 3})
```

### 错误断言

```go
// Error - 有错误
assert.Error(t, err)

// NoError - 没有错误
assert.NoError(t, err)

// ErrorIs - 特定错误
assert.ErrorIs(t, err, ErrNotFound)

// ErrorContains - 错误包含文本
assert.ErrorContains(t, err, "not found")

// Eventually - 最终成功（带重试）
assert.Eventually(t, func() bool {
    return checkCondition()
}, time.Second, 100*time.Millisecond)
```

### 对象断言

```go
// IsType - 类型匹配
assert.IsType(t, &User{}, actual)

// Implements - 实现接口
assert.Implements(t, (*io.Reader)(nil), actual)

// NotSame - 不是同一个对象
a := []int{1, 2}
b := []int{1, 2}
assert.NotSame(t, a, b)
```

## Mock框架使用

### 定义Mock对象

```go
// 1. 定义接口
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
}

// 2. 创建Mock
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}
```

### 设置Mock期望

```go
func TestUserService_GetUser(t *testing.T) {
    // 创建Mock
    mockRepo := new(MockUserRepository)

    // 设置期望 - 简单返回
    mockRepo.On("GetByID", mock.Anything, "user123").Return(&User{
        ID: "user123",
        Username: "testuser",
    }, nil)

    // 设置期望 - 匹配特定参数
    mockRepo.On("GetByID", mock.Anything, "invalid").
        Return(nil, errors.New("not found"))

    // 使用Mock
    service := NewUserService(mockRepo)
    user, err := service.GetUser(context.Background(), "user123")

    // 断言
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)

    // 验证Mock调用
    mockRepo.AssertExpectations(t)
    mockRepo.AssertCalled(t, "GetByID", mock.Anything, "user123")
    mockRepo.AssertNumberOfCalls(t, "GetByID", 1)
}
```

### Mock参数匹配

```go
// Anything - 匹配任何值
mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(user, nil)

// AnythingOfType - 匹配类型
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*User")).Return(nil)

// MatchedBy - 自定义匹配
mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *User) bool {
    return u.Username == "testuser" && u.Email == "test@example.com"
})).Return(nil)

// Arguments - 具体参数
mockRepo.On("GetByID", mock.Anything, "user123").Return(user, nil)
```

### Mock高级用法

```go
// Once - 只调用一次
mockRepo.On("GetByID", mock.Anything, "user123").Return(user, nil).Once()

// Times - 调用n次
mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(3)

// Return - 多次返回不同值
mockRepo.On("GetByID", mock.Anything, "user123").
    Return(user1, nil).
    Once().
    On("GetByID", mock.Anything, "user123").
    Return(user2, nil).
    Once()

// Run - 执行回调
mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *User) bool {
    return true
})).Return(nil).Run(func(args mock.Arguments) {
    user := args.Get(1).(*User)
    user.ID = "generated_id"  // 模拟生成ID
})

// Maybe - 可选调用
mockRepo.On("OptionalMethod", mock.Anything).Maybe().Return(nil)
```

## 测试套件（testify/suite）

### 创建测试套件

```go
type UserServiceTestSuite struct {
    suite.Suite
    mockRepo *MockUserRepository
    service  *UserService
}

// SetupSuite - 所有测试前执行一次
func (s *UserServiceTestSuite) SetupSuite() {
    // 初始化共享资源
}

// TearDownSuite - 所有测试后执行一次
func (s *UserServiceTestSuite) TearDownSuite() {
    // 清理共享资源
}

// SetupTest - 每个测试前执行
func (s *UserServiceTestSuite) SetupTest() {
    s.mockRepo = new(MockUserRepository)
    s.service = NewUserService(s.mockRepo)
}

// TearDownTest - 每个测试后执行
func (s *UserServiceTestSuite) TearDownTest() {
    // 清理测试数据
}

// 测试方法
func (s *UserServiceTestSuite) TestCreateUser_Success() {
    // Arrange
    user := &User{Username: "testuser"}
    s.mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    // Act
    err := s.service.CreateUser(context.Background(), user)

    // Assert
    s.NoError(err)
    s.mockRepo.AssertExpectations(s.T())
}

// 运行测试套件
func TestUserServiceTestSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

## 最佳实践

### 1. 使用Table-Driven测试

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
        {"missing domain", "test@", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 2. 使用描述性消息

```go
// ✅ 好
assert.NoError(t, err, "创建用户失败")
assert.Equal(t, expected, actual, "用户名不匹配")

// ❌ 差
assert.NoError(t, err)
assert.Equal(t, expected, actual)
```

### 3. 使用require处理前置条件

```go
func TestExample(t *testing.T) {
    // Setup阶段 - 使用require
    db, err := setupTestDB()
    require.NoError(t, err, "Setup数据库失败")
    require.NotNil(t, db, "数据库不应该为nil")

    // 测试阶段 - 使用assert
    result, err := service.Query()
    assert.NoError(t, err, "查询失败")
    assert.NotNil(t, result)
}
```

### 4. 使用自定义断言

```go
// 自定义断言函数
func assertUserEqual(t *testing.T, expected, actual *User) {
    t.Helper()

    assert.Equal(t, expected.ID, actual.ID, "ID不匹配")
    assert.Equal(t, expected.Username, actual.Username, "用户名不匹配")
    assert.Equal(t, expected.Email, actual.Email, "邮箱不匹配")
}

// 使用
func TestExample(t *testing.T) {
    expected := &User{ID: "1", Username: "test"}
    actual := getUser()
    assertUserEqual(t, expected, actual)
}
```

## 完整示例

```go
package service_test

import (
    "context"
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestUserService_GetUser(t *testing.T) {
    tests := []struct {
        name      string
        userID    string
        setup     func(*MockUserRepository)
        wantErr   bool
        errMsg    string
        wantUser  *User
    }{
        {
            name:   "成功获取用户",
            userID: "user123",
            setup: func(m *MockUserRepository) {
                m.On("GetByID", mock.Anything, "user123").
                    Return(&User{ID: "user123", Username: "test"}, nil)
            },
            wantErr:  false,
            wantUser: &User{ID: "user123", Username: "test"},
        },
        {
            name:   "用户不存在",
            userID: "invalid",
            setup: func(m *MockUserRepository) {
                m.On("GetByID", mock.Anything, "invalid").
                    Return(nil, errors.New("not found"))
            },
            wantErr: true,
            errMsg:  "not found",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(MockUserRepository)
            tt.setup(mockRepo)

            service := NewUserService(mockRepo)
            ctx := context.Background()

            // Act
            user, err := service.GetUser(ctx, tt.userID)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
                assert.Nil(t, user)
            } else {
                require.NoError(t, err)
                require.NotNil(t, user)
                assert.Equal(t, tt.wantUser.ID, user.ID)
                assert.Equal(t, tt.wantUser.Username, user.Username)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

## 参考资源

- [testify GitHub](https://github.com/stretchr/testify)
- [assert文档](https://pkg.go.dev/github.com/stretchr/testify/assert)
- [mock文档](https://pkg.go.dev/github.com/stretchr/testify/mock)
- [测试模板集](../../06_快速参考/测试模板集.md)
