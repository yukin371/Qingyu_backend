# Service层测试改进快速指南

## 📋 本指南概述

本指南提供了快速完善Service层测试的实用工具和模板，帮助开发者高效地为各个Service编写高质量的单元测试。

---

## 🚀 快速开始

### 第一步：选择要测试的Service

从优先级列表中选择一个Service：

1. **P0（本周）**：UserService, ProjectService, DocumentService, AuthService
2. **P1（本周末）**：BookstoreService系列, ReadingHistoryService, StorageService, AutoSaveService
3. **P2（下周）**：Writer模块, 缓存服务, 统计服务

### 第二步：复制并修改模板

使用下面的标准模板，替换`[ServiceName]`, `[MethodName]`, `[Repo]`等占位符。

### 第三步：运行并验证

```bash
go test ./test/service/[module]/... -v
go test ./test/service/[module]/... -cover
```

---

## 📝 标准测试模板

### 1. 基础Service测试模板

```go
package [module]_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pkgErrors "Qingyu_backend/pkg/errors"
	"Qingyu_backend/service/[module]"
	"Qingyu_backend/service/base"
)

// Mock[Repository] 实现[Repository]接口
type Mock[Repository] struct {
	mock.Mock
}

// 实现接口方法...
func (m *Mock[Repository]) Method1(ctx context.Context, param string) (string, error) {
	args := m.Called(ctx, param)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

// MockEventBus 实现EventBus接口
type MockEventBus struct {
	publishedEvents []base.Event
	mock.Mock
}

func (m *MockEventBus) Subscribe(handler base.EventHandler) error {
	args := m.Called(handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(handler base.EventHandler) error {
	args := m.Called(handler)
	return args.Error(0)
}

func (m *MockEventBus) Publish(event base.Event) error {
	args := m.Called(event)
	m.publishedEvents = append(m.publishedEvents, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(event base.Event, callback func(error)) error {
	args := m.Called(event, callback)
	m.publishedEvents = append(m.publishedEvents, event)
	return args.Error(0)
}

// ============ 测试用例 ============

// Test[Service]_[MethodName] 测试方法
func Test[Service]_[MethodName](t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// 1. 准备
		mockRepo := new(Mock[Repository])
		mockEventBus := new(MockEventBus)
		service := [module].New[Service](mockRepo, mockEventBus)
		ctx := context.Background()

		// 2. 设置Mock期望
		mockRepo.On("Method1", ctx, "test").Return("result", nil).Once()
		mockEventBus.On("Publish", mock.MatchedBy(func(e base.Event) bool {
			return e.GetEventType() == "[EventType]"
		})).Return(nil).Once()

		// 3. 执行
		result, err := service.[MethodName](ctx, "test")

		// 4. 断言
		assert.NoError(t, err)
		assert.Equal(t, "expected", result)
		mockRepo.AssertExpectations(t)
		mockEventBus.AssertExpectations(t)

		t.Logf("✓ 成功流程测试通过")
	})

	t.Run("InvalidInput", func(t *testing.T) {
		mockRepo := new(Mock[Repository])
		mockEventBus := new(MockEventBus)
		service := [module].New[Service](mockRepo, mockEventBus)
		ctx := context.Background()

		// 参数验证失败应返回ValidationError
		_, err := service.[MethodName](ctx, "")
		
		assert.Error(t, err)
		assert.True(t, pkgErrors.IsValidationError(err))
		
		t.Logf("✓ 参数验证测试通过")
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(Mock[Repository])
		mockEventBus := new(MockEventBus)
		service := [module].New[Service](mockRepo, mockEventBus)
		ctx := context.Background()

		// 设置Mock返回NotFound错误
		mockRepo.On("Method1", ctx, mock.Anything).
			Return(nil, pkgErrors.NewNotFoundError("[Entity]", "not found")).Once()

		_, err := service.[MethodName](ctx, "notexist")

		assert.Error(t, err)
		assert.True(t, pkgErrors.IsNotFoundError(err))
		
		t.Logf("✓ 资源不存在测试通过")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(Mock[Repository])
		mockEventBus := new(MockEventBus)
		service := [module].New[Service](mockRepo, mockEventBus)
		ctx := context.Background()

		// 数据库错误
		mockRepo.On("Method1", ctx, mock.Anything).
			Return(nil, pkgErrors.NewInternalError("database error")).Once()

		_, err := service.[MethodName](ctx, "test")

		assert.Error(t, err)
		
		t.Logf("✓ 数据库错误测试通过")
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		mockRepo := new(Mock[Repository])
		mockEventBus := new(MockEventBus)
		service := [module].New[Service](mockRepo, mockEventBus)
		ctx := context.Background()

		// 并发调用
		mockRepo.On("Method1", ctx, mock.Anything).Return("result", nil)
		mockEventBus.On("Publish", mock.Anything).Return(nil)

		// 并发执行
		done := make(chan error, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_, err := service.[MethodName](ctx, "test")
				done <- err
			}()
		}

		// 验证所有并发请求都成功
		for i := 0; i < 10; i++ {
			err := <-done
			assert.NoError(t, err)
		}

		t.Logf("✓ 并发访问测试通过")
	})
}
```

---

## 🎯 特定Service的测试模板

### UserService测试
```go
// 需要测试的关键方法
func TestUserService_Register(t *testing.T) {
	t.Run("SuccessfulRegistration", func(t *testing.T) {
		// 验证用户成功注册
	})
	t.Run("DuplicateEmail", func(t *testing.T) {
		// 验证邮箱已存在时的错误处理
	})
	t.Run("InvalidPassword", func(t *testing.T) {
		// 验证密码规则验证
	})
	t.Run("PasswordHashingCorrect", func(t *testing.T) {
		// 验证密码被正确哈希
	})
}

func TestUserService_Login(t *testing.T) {
	t.Run("SuccessfulLogin", func(t *testing.T) {
		// 验证用户成功登录
	})
	t.Run("WrongPassword", func(t *testing.T) {
		// 验证密码错误
	})
	t.Run("UserNotFound", func(t *testing.T) {
		// 验证用户不存在
	})
}

func TestUserService_GetUserProfile(t *testing.T) {
	t.Run("SuccessfulRetrieval", func(t *testing.T) {
		// 验证获取用户资料
	})
	t.Run("Unauthorized", func(t *testing.T) {
		// 验证权限检查
	})
}
```

### ProjectService测试
```go
func TestProjectService_CreateProject(t *testing.T) {
	t.Run("SuccessfulCreation", func(t *testing.T) {
		// 验证项目创建
	})
	t.Run("InvalidProjectName", func(t *testing.T) {
		// 验证项目名称验证
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	t.Run("SuccessfulUpdate", func(t *testing.T) {
		// 验证项目更新
	})
	t.Run("UnauthorizedUpdate", func(t *testing.T) {
		// 验证权限检查
	})
	t.Run("VersionConflict", func(t *testing.T) {
		// 验证并发冲突处理
	})
}
```

### DocumentService测试
```go
func TestDocumentService_CreateDocument(t *testing.T) {
	t.Run("SuccessfulCreation", func(t *testing.T) {
		// 验证文档创建
	})
	t.Run("ProjectNotFound", func(t *testing.T) {
		// 验证项目权限检查
	})
	t.Run("DuplicateName", func(t *testing.T) {
		// 验证文档名称唯一性
	})
}

func TestDocumentService_UpdateContent(t *testing.T) {
	t.Run("SuccessfulUpdate", func(t *testing.T) {
		// 验证内容更新
	})
	t.Run("AutoSaveTrigger", func(t *testing.T) {
		// 验证自动保存事件触发
	})
}
```

---

## ✅ 测试检查清单

### 每个测试函数应包含

- [ ] **Happy Path（正常流程）**
  - 验证方法在正常输入下返回预期结果
  - 验证事件被正确发布（如适用）

- [ ] **Input Validation（参数验证）**
  - 空值处理
  - 无效格式
  - 超长字符串
  - 负数/零值
  - 返回 `ValidationError`

- [ ] **Business Rules（业务规则）**
  - 重复检查（邮箱、项目名等）
  - 权限检查
  - 业务约束（如余额检查）
  - 返回 `BusinessError`

- [ ] **Error Handling（错误处理）**
  - Repository 返回 `NotFound`
  - Repository 返回其他错误
  - Context 超时/取消
  - 返回正确的错误类型

- [ ] **Permissions（权限检查）**
  - 用户权限验证
  - 所有权检查
  - 返回 `AuthError`

- [ ] **Concurrent Access（并发访问）**
  - 幂等性检查
  - 竞态条件处理
  - 使用goroutines进行测试

- [ ] **Cache Behavior（缓存行为，如适用）**
  - 缓存命中
  - 缓存失效
  - 缓存更新

- [ ] **Event Publishing（事件发布）**
  - 验证发布了正确的事件类型
  - 验证事件数据正确
  - 验证事件在适当时机发布

### Mock设计检查

- [ ] Mock 实现了接口的所有方法
- [ ] Mock 的方法签名与接口完全一致
- [ ] 使用 `AssertExpectations()` 验证所有期望都满足
- [ ] 使用 `mock.Anything` 处理无需检查的参数
- [ ] 使用 `mock.MatchedBy()` 进行复杂的参数匹配

---

## 🔧 常见错误模式

### ❌ 常见错误1：Mock方法签名不匹配
```go
// ❌ 错误
mockRepo.On("GetByID", "id").Return(nil, nil)  // 缺少context.Context

// ✅ 正确
mockRepo.On("GetByID", mock.Anything, "id").Return(nil, nil)
```

### ❌ 常见错误2：未验证Mock期望
```go
// ❌ 错误 - 没有AssertExpectations
func TestService_Method(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Method", mock.Anything).Return(nil)
	// ... 执行测试 ...
	// 缺少 mockRepo.AssertExpectations(t)
}

// ✅ 正确
func TestService_Method(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Method", mock.Anything).Return(nil).Once()
	// ... 执行测试 ...
	mockRepo.AssertExpectations(t)
}
```

### ❌ 常见错误3：过度Mock
```go
// ❌ 错误 - Mock了太多不必要的调用
mockRepo.On("Method1", mock.Anything).Return(nil)
mockRepo.On("Method2", mock.Anything).Return(nil)
mockRepo.On("Method3", mock.Anything).Return(nil)
// ... 20+ 行Mock设置 ...

// ✅ 正确 - 只Mock相关的调用
mockRepo.On("Method1", mock.Anything).Return(data, nil).Times(1)
```

### ❌ 常见错误4：缺少错误路径测试
```go
// ❌ 错误 - 只测试成功的情况
func TestService_Method(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Method", mock.Anything).Return(data, nil)
	result, err := service.Method(ctx, "test")
	assert.NoError(t, err)
}

// ✅ 正确 - 包含多个子测试
func TestService_Method(t *testing.T) {
	t.Run("Success", func(t *testing.T) { ... })
	t.Run("NotFound", func(t *testing.T) { ... })
	t.Run("RepositoryError", func(t *testing.T) { ... })
	t.Run("ValidationError", func(t *testing.T) { ... })
}
```

---

## 📊 测试覆盖率检查

### 生成覆盖率报告
```bash
# 生成覆盖率数据
go test ./test/service/... -coverprofile=coverage.out

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 查看specific包的覆盖率
go test ./test/service/user/... -cover
```

### 覆盖率目标
- 关键业务逻辑：≥ 85%
- 正常服务：≥ 80%
- 工具类/缓存：≥ 70%

---

## 🔍 Mock完整性检查清单

### 检查Mock是否实现了所有接口方法

```bash
# 步骤1：查看接口定义
cat repository/interfaces/[module]/[interface]_interface.go

# 步骤2：在Mock中实现所有方法
# Mock应该有interface中的每个方法的实现

# 步骤3：运行编译检查
go test ./test/service/[module]/... -v
```

---

## 🚦 PR审核检查清单

提交PR前，确保满足以下条件：

- [ ] 所有新增Service都有对应的测试
- [ ] 所有测试都能成功运行
- [ ] 测试覆盖率 ≥ 80%
- [ ] Mock实现了所有接口方法
- [ ] 所有测试用例都使用`t.Run()`组织
- [ ] 错误路径都有测试
- [ ] 权限检查都有测试（如适用）
- [ ] 并发场景都有测试
- [ ] 使用了 `AssertExpectations()`
- [ ] 测试代码可读性好

---

## 📖 相关文档

- [Service层测试覆盖分析报告](SERVICE_TEST_COVERAGE_REPORT.md)
- [Service层架构规范](../../doc/architecture/架构设计规范.md)
- [软件工程规范](../../doc/engineering/软件工程规范_v2.0.md)
- [项目开发规则](../../doc/architecture/项目开发规则.md)

---

## 💡 建议和最佳实践

### 1. 定期检查覆盖率
每周五运行一次覆盖率检查，跟踪改进进度。

### 2. 测试驱动开发（TDD）
新增Service时，先写测试后写实现。

### 3. Mock管理
- Mock放在各自的子包内
- Mock与接口保持同步
- 定期review Mock实现

### 4. 测试代码质量
- 保持测试代码简洁
- 避免重复代码，提取公共函数
- 使用有意义的变量名
- 添加清晰的注释

### 5. 错误消息
测试失败时输出清晰的错误信息：
```go
assert.Equal(t, expected, actual, "用户名应该是%s", expected)
```

---

**最后更新**：2025-10-31  
**维护者**：后端架构团队

