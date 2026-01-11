# 真实数据测试规范

## ⚠️ 核心原则

### 禁止使用Hack，必须使用真实数据

青羽写作平台严格要求在某些测试场景中使用**真实数据库**和**真实数据**，禁止使用Mock或其他临时hack方案。

## 为什么必须使用真实数据？

### ❌ 使用Mock/Hack的危害

1. **隐藏真实Bug**
   ```go
   // ❌ 错误示例：使用Mock Repository测试Repository层
   func TestUserRepository_Create(t *testing.T) {
       mockRepo := new(MockUserRepository)  // 危险！
       mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

       // 测试通过了，但真实的数据库操作可能有Bug
   }
   ```
   **问题**：无法发现数据库连接问题、索引错误、类型不匹配等

2. **测试与生产环境不一致**
   ```go
   // ❌ 错误示例：使用内存数据库
   func TestBookRepository_Query(t *testing.T) {
       db := setupMemoryDB()  // 使用SQLite内存数据库

       // 测试通过了，但MongoDB的查询语法不同
   }
   ```
   **问题**：MongoDB的聚合管道、索引特性无法在内存数据库中验证

3. **临时Hack导致维护困难**
   ```go
   // ❌ 错误示例：临时修改代码适配测试
   // if os.Getenv("TESTING") == "true" {
   //     return mockData  // 临时hack
   // }
   ```
   **问题**：可能忘记移除，污染生产代码

### ✅ 使用真实数据的好处

1. **发现真实Bug**
   - 数据库连接问题
   - 索引缺失导致的性能问题
   - 数据类型不匹配
   - 事务处理错误

2. **测试即文档**
   - 测试代码展示了如何正确使用数据库
   - 新人可以通过测试了解数据模型

3. **重构信心**
   - 修改Repository代码时，真实测试能立即发现问题
   - 确保数据库操作的正确性

## 必须使用真实数据的场景

### ✅ 1. Repository层所有测试

**原因**：Repository层的职责就是与数据库交互

```go
// ✅ 正确示例：使用真实MongoDB
func TestAuthRepository_CreateRole(t *testing.T) {
    // Setup - 使用真实测试数据库
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := shared.NewAuthRepository(db)
    ctx := context.Background()

    role := &authModel.Role{
        Name:        "test_role",
        Permissions: []string{"test.read", "test.write"},
        IsSystem:    false,
    }

    // Act - 真实数据库操作
    err := repo.CreateRole(ctx, role)

    // Assert - 验证数据库中的结果
    require.NoError(t, err)

    // 从数据库查询验证
    found, err := repo.GetRole(ctx, role.ID)
    require.NoError(t, err)
    assert.Equal(t, role.Name, found.Name)
}
```

**位置**：`repository/mongodb/**/*_test.go`

---

### ✅ 2. 集成测试

**原因**：验证多个组件协作时的正确性

```go
// ✅ 正确示例：完整集成测试
func TestBookPurchaseFlow(t *testing.T) {
    // Setup - 真实数据库和HTTP服务器
    router, cleanup := setupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 1. 创建书籍（真实数据库）
    bookID := helper.CreateTestBook()

    // 2. 用户登录
    token := helper.LoginTestUser()

    // 3. 购买章节（真实HTTP请求 + 真实数据库）
    reqBody := map[string]interface{}{
        "book_id": bookID,
        "chapter_ids": []string{"chapter1"},
    }
    w := helper.DoAuthRequest("POST", "/api/v1/bookstore/purchase", reqBody, token)
    helper.AssertSuccess(w, 200, "购买失败")

    // 4. 验证数据库中的购买记录
}
```

**位置**：`test/integration/*.go`

---

### ✅ 3. 性能测试

**原因**：只有真实数据库才能反映真实性能

```go
// ✅ 正确示例：真实性能测试
func BenchmarkBookRepository_GetByID(b *testing.B) {
    // Setup - 真实数据库
    db, cleanup := testutil.SetupTestDB(&testing.T{})
    defer cleanup()

    repo := NewBookRepository(db)

    // 准备测试数据
    book := &Book{Title: "测试书籍"}
    repo.Create(context.Background(), book)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 真实数据库查询
        repo.GetByID(context.Background(), book.ID)
    }
}
```

**位置**：`test/performance/*.go`

---

## 禁止的Hack方案

### ❌ 1. 禁止Mock Repository测试Repository

```go
// ❌ 错误
func TestUserRepository_Create(t *testing.T) {
    mockRepo := new(MockUserRepository)  // 危险！
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    // ...
}

// ✅ 正确
func TestUserRepository_Create(t *testing.T) {
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := NewUserRepository(db)
    // 使用真实数据库
}
```

### ❌ 2. 禁止使用内存数据库

```go
// ❌ 错误
func TestBookRepository(t *testing.T) {
    db := setupMemoryDB()  // 使用SQLite或其他内存DB
    // ...
}

// ✅ 正确
func TestBookRepository(t *testing.T) {
    db, cleanup := testutil.SetupTestDB(t)  // 真实MongoDB
    defer cleanup()
    // ...
}
```

### ❌ 3. 禁止临时修改代码

```go
// ❌ 错误：在代码中添加测试标志
func (r *Repository) Query(filter interface{}) ([]*Model, error) {
    if os.Getenv("TESTING") == "true" {
        return mockData  // 临时hack
    }
    // 真实逻辑
}

// ✅ 正确：使用依赖注入
type Repository struct {
    db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
    return &Repository{db: db}
}
```

### ❌ 4. 禁止使用假数据绕过验证

```go
// ❌ 错误：使用无效测试数据
func TestValidation(t *testing.T) {
    user := &User{
        Email: "invalid-email",  // 无效邮箱，但测试通过了
    }
    // ...
}

// ✅ 正确：使用真实格式的测试数据
func TestValidation(t *testing.T) {
    factory := fixtures.NewUserFactory()
    user := factory.Create()  // 生成符合规则的数据
    // ...
}
```

## 测试数据准备规范

### 使用Factory模式创建测试数据

```go
// ✅ 使用测试数据工厂
func TestUserService_CreateUser(t *testing.T) {
    // 创建工厂
    userFactory := fixtures.NewUserFactory()

    // 创建标准测试用户
    user := userFactory.Create()

    // 创建特定类型的用户
    admin := userFactory.CreateAdmin()
    author := userFactory.CreateAuthor()

    // 批量创建
    users := userFactory.CreateBatch(10)

    // 自定义属性
    customUser := userFactory.Create(func(u *users.User) {
        u.Username = "custom_name"
        u.Email = "custom@example.com"
    })
}
```

### 测试数据隔离

```go
// ✅ 每个测试独立的数据
func TestExample(t *testing.T) {
    // Setup - 创建独立的测试数据
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()  // Cleanup会清理所有测试数据

    // 使用factory创建数据
    factory := fixtures.NewUserFactory()
    user := factory.Create()

    // 测试逻辑...
}
```

## 测试数据库配置

### 测试数据库Setup

```go
// test/testutil/database.go
func SetupTestDB(t *testing.T) (*mongo.Database, func()) {
    t.Helper()

    // 加载配置
    cfg, err := config.LoadConfig("config/config.yaml")
    if err != nil {
        t.Fatalf("加载配置失败: %v", err)
    }

    // 初始化服务容器
    c := container.NewServiceContainer()
    err = c.Initialize(context.Background())
    if err != nil {
        t.Fatalf("初始化失败: %v", err)
    }

    // 获取测试数据库
    db := c.GetMongoDB()

    // 返回清理函数
    cleanup := func() {
        // 清理所有测试集合
        collections := []string{
            "users", "books", "projects", "documents",
            "roles", "permissions", "transactions",
        }
        ctx := context.Background()
        for _, coll := range collections {
            _ = db.Collection(coll).Drop(ctx)
        }
        _ = c.Close(ctx)
    }

    return db, cleanup
}
```

### 配置文件（config/config.yaml）

```yaml
database:
  mongodb:
    # 测试数据库配置
    uri: "mongodb://localhost:27017"
    database: "qingyu_test"  # 使用独立的测试数据库
    options:
      maxPoolSize: 10
      minPoolSize: 1
```

## 数据隔离策略

### 1. 使用独立测试数据库

```yaml
# 测试环境使用独立数据库
database:
  mongodb:
    database: "qingyu_test"  # 不是qingyu_prod
```

### 2. 每个测试独立清理

```go
func TestExample(t *testing.T) {
    // Setup
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()  // 测试结束后清理

    // 测试逻辑...
}
```

### 3. 使用测试用户前缀

```go
// 测试数据使用特定前缀，便于识别和清理
const (
    TestUserPrefix = "test_user_"
    TestBookPrefix = "test_book_"
)

func CreateTestUser(t *testing.T) *User {
    return &User{
        Username: TestUserPrefix + primitive.NewObjectID().Hex(),
        Email: TestUserPrefix + "@test.com",
    }
}
```

## 相关文档

- [测试数据准备规范](测试数据准备规范.md) - Factory模式详细使用
- [数据隔离策略](数据隔离策略.md) - 避免数据污染的完整方案
- [MongoDB测试配置指南](../03_测试工具指南/mongodb_测试配置指南.md) - 测试数据库配置
