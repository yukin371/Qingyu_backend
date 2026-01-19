# 禁止Mock规范

## 概述

本文档明确列出在测试中**严格禁止**的Mock方案，这些方案会破坏测试的有效性，隐藏真实Bug，增加维护成本。

## 严格禁止的Mock方案

### ❌ 1. 禁止Mock Repository测试Repository层

**问题描述**：使用Mock Repository来测试Repository层代码

```go
// ❌ 严格禁止
func TestUserRepository_CreateUser(t *testing.T) {
    // 危险！这完全绕过了真实的数据库操作
    mockRepo := new(MockUserRepository)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    err := mockRepo.Create(context.Background(), user)
    assert.NoError(t, err)

    // 测试通过了，但真实的数据库操作可能有：
    // - 连接问题
    // - 索引错误
    // - 类型不匹配
    // - 字段名错误
}
```

**为什么禁止**：
- 无法发现数据库连接问题
- 无法验证查询语法是否正确
- 无法发现索引缺失问题
- 无法验证数据类型兼容性
- 测试与生产环境不一致

**正确做法**：
```go
// ✅ 正确：使用真实MongoDB
func TestUserRepository_CreateUser(t *testing.T) {
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := NewUserRepository(db)
    err := repo.Create(context.Background(), user)

    // 验证数据确实写入了数据库
    assert.NoError(t, err)

    // 从数据库查询验证
    found, _ := repo.GetByID(context.Background(), user.ID)
    assert.Equal(t, user.Username, found.Username)
}
```

---

### ❌ 2. 禁止使用内存数据库

**问题描述**：使用SQLite、内存MongoDB等替代真实数据库

```go
// ❌ 严格禁止
func TestBookRepository(t *testing.T) {
    // 使用内存数据库完全改变了查询特性
    db := setupMemoryDB()

    books, err := repo.Query(bson.M{"status": "published"})
    // 内存数据库的查询语法和MongoDB不同
}
```

**为什么禁止**：
- 不同数据库的查询语法差异巨大
- 索引特性完全不同
- 聚合管道无法测试
- 事务行为不一致
- 性能特征完全不同

**正确做法**：
```go
// ✅ 正确：使用真实MongoDB测试数据库
func TestBookRepository(t *testing.T) {
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := NewBookRepository(db)
    // 使用真实的MongoDB
}
```

---

### ❌ 3. 禁止使用环境变量Mock

**问题描述**：在代码中检查TESTING环境变量并返回假数据

```go
// ❌ 严格禁止
func (r *Repository) Query(filter interface{}) ([]*Model, error) {
    if os.Getenv("TESTING") == "true" {
        // 临时mock用于测试
        return mockData, nil
    }

    // 真实逻辑
    cursor, err := r.collection.Find(context.Background(), filter)
    // ...
}
```

**为什么禁止**：
- 可能忘记移除，污染生产代码
- 代码审查难以发现
- Git blame显示污染
- 测试覆盖的代码路径和生产不同
- 增加代码复杂度

**正确做法**：
```go
// ✅ 正确：使用依赖注入
type Repository struct {
    db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Query(filter interface{}) ([]*Model, error) {
    // 只有真实逻辑，没有任何测试分支
    cursor, err := r.collection.Find(context.Background(), filter)
    // ...
}
```

---

### ❌ 4. 禁止硬编码测试返回值

**问题描述**：在函数中硬编码测试专用的返回值

```go
// ❌ 严格禁止
func GetUserData(userID string) (*User, error) {
    // 真实逻辑
    user, err := r.repo.GetByID(userID)

    // 危险的测试hack
    if userID == "test_user_123" {
        return &User{ID: "test_user_123", Username: "测试用户"}, nil
    }

    return user, err
}
```

**为什么禁止**：
- 测试专用代码混入生产代码
- 难以维护和理解
- 可能被滥用
- 增加代码复杂度
- 测试数据和真实数据不一致

**正确做法**：
```go
// ✅ 正确：只有真实逻辑
func GetUserData(userID string) (*User, error) {
    return r.repo.GetByID(userID)
}

// 测试中使用Factory创建测试数据
func TestGetUserData(t *testing.T) {
    factory := fixtures.NewUserFactory()
    user := factory.Create()

    result, err := service.GetUserData(user.ID)
    // ...
}
```

---

### ❌ 5. 禁止使用假数据绕过验证

**问题描述**：使用不符合生产规范的测试数据

```go
// ❌ 严格禁止
func TestEmailValidation(t *testing.T) {
    // 使用假数据绕过验证逻辑
    user := &User{
        Email: "invalid-email",  // 明知无效但测试不管
        Phone: "12345",          // 格式错误
    }

    err := service.CreateUser(user)
    // 测试通过了，但验证逻辑可能没有真正执行
}
```

**为什么禁止**：
- 验证逻辑没有被测试
- 可能绕过了重要的业务规则
- 测试覆盖虚假
- 无法发现验证Bug

**正确做法**：
```go
// ✅ 正确：使用Factory生成符合规则的数据
func TestEmailValidation(t *testing.T) {
    factory := fixtures.NewUserFactory()

    // 测试有效的邮箱
    user := factory.Create(func(u *User) {
        u.Email = "valid@example.com"
    })

    err := service.CreateUser(user)
    assert.NoError(t, err)

    // 测试无效的邮箱
    invalidUser := factory.Create(func(u *User) {
        u.Email = "invalid-email"
    })

    err = service.CreateUser(invalidUser)
    assert.Error(t, err)
}
```

---

### ❌ 6. 禁止Mock第三方服务

**问题描述**：Mock支付、短信、邮件等外部服务

```go
// ❌ 严格禁止（集成测试中）
func TestPaymentFlow(t *testing.T) {
    // Mock了支付网关，但真实集成测试应该验证支付流程
    mockPaymentGateway := new(MockPaymentGateway)
    mockPaymentGateway.On("Charge", mock.Anything).Return(&PaymentResult{
        Success: true,  // 假数据，无法发现真实的支付问题
    })
}
```

**为什么禁止**（集成测试中）：
- 无法发现API接口变更
- 无法发现网络问题
- 无法发现认证问题
- 无法发现数据格式问题

**正确做法**（集成测试）：
```go
// ✅ 正确：使用测试环境的真实支付服务
func TestPaymentFlow(t *testing.T) {
    // 使用测试环境的支付服务
    paymentService := NewPaymentService(config.TestAPIKey)

    // 测试小额支付
    result, err := paymentService.Charge(context.Background(), &ChargeRequest{
        Amount: 1,  // 测试金额
        Currency: "CNY",
    })

    assert.NoError(t, err)
    assert.True(t, result.Success)
}
```

---

### ❌ 7. 禁止禁用错误检查

**问题描述**：在测试中临时禁用某些检查

```go
// ❌ 严格禁止
func TestComplexFlow(t *testing.T) {
    // 临时禁用权限检查以通过测试
    originalCheck := enableAuthCheck
    enableAuthCheck = false  // 危险的hack
    defer func() { enableAuthCheck = originalCheck }()

    // 执行测试
    result := service.DoSomething()
    assert.NotNil(t, result)
}
```

**为什么禁止**：
- 可能导致安全漏洞
- 测试通过的代码在生产中失败
- 难以追踪临时修改
- 容易忘记恢复

**正确做法**：
```go
// ✅ 正确：使用测试用户权限
func TestComplexFlow(t *testing.T) {
    // 创建有足够权限的测试用户
    factory := fixtures.NewUserFactory()
    admin := factory.CreateAdmin()

    // 使用admin的token进行测试
    token := generateTestToken(admin)
    result := service.DoSomething(token)

    assert.NotNil(t, result)
}
```

---

### ❌ 8. 禁止修改全局状态

**问题描述**：在测试中修改全局变量或单例

```go
// ❌ 严格禁止
var GlobalConfig = &Config{Debug: false}

func TestProductionMode(t *testing.T) {
    // 临时修改全局配置
    oldConfig := GlobalConfig
    GlobalConfig = &Config{Debug: true}  // 污染全局状态
    defer func() { GlobalConfig = oldConfig }()

    // 测试逻辑
}
```

**为什么禁止**：
- 污染全局状态
- 并发测试会失败
- 测试顺序依赖
- 难以调试

**正确做法**：
```go
// ✅ 正确：使用依赖注入
type Service struct {
    config *Config
}

func NewService(config *Config) *Service {
    return &Service{config: config}
}

func TestProductionMode(t *testing.T) {
    // 每个测试独立创建service
    config := &Config{Debug: true}
    service := NewService(config)

    // 测试逻辑，不影响其他测试
}
```

---

## Hack检测清单

代码审查时，检查是否存在以下hack：

- [ ] Repository测试中使用了Mock Repository
- [ ] 使用了内存数据库或假数据库
- [ ] 代码中有`os.Getenv("TESTING")`检查
- [ ] 生产代码中有测试专用的分支
- [ ] 使用了硬编码的测试返回值
- [ ] 使用了不符合规范的测试数据
- [ ] 集成测试中Mock了外部服务
- [ ] 临时禁用了权限或验证检查
- [ ] 修改了全局变量或单例

## 相关文档

- [真实数据库测试规范](真实数据库测试规范.md) - 如何正确使用真实数据库
- [测试数据准备规范](测试数据准备规范.md) - Factory模式创建测试数据
- [数据隔离策略](数据隔离策略.md) - 避免数据污染
