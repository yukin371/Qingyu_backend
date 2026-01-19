# 测试层级规范

## 概述

本目录包含青羽写作平台各测试层级的详细规范。每个层级都有明确的测试策略、工具使用和最佳实践。

## 测试层级概览

```
┌─────────────────────────────────────┐
│         E2E测试 (10%)                │  完整业务流程
│      真实环境 + 真实数据              │  test/e2e/
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│        API测试 (30%)                 │  HTTP集成测试
│   HTTP请求 + 真实数据库               │  test/api/
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│       Service测试 (60%)              │  业务逻辑测试
│     Mock Repository + 业务逻辑        │  service/*/
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│     Repository测试 (100%)            │  数据访问测试
│      真实MongoDB + CRUD操作          │  repository/mongodb/*/
└─────────────────────────────────────┘
```

## 测试层级对比

| 层级 | 目的 | 数据源 | 执行速度 | 测试方式 | 文档 |
|------|------|--------|----------|----------|------|
| **Repository** | 验证数据库操作 | 真实MongoDB | 慢 (秒级) | 真实DB操作 | [规范](./repository_层测试规范.md) |
| **Service** | 验证业务逻辑 | Mock Repository | 快 (毫秒级) | 单元测试 | [规范](./service_层测试规范.md) |
| **API** | 验证HTTP接口 | 真实数据库 | 中 (百毫秒级) | 集成测试 | [规范](./api_层测试规范.md) |
| **E2E** | 验证业务流程 | 真实环境 | 慢 (秒级) | 端到端测试 | [规范](./e2e_测试规范.md) |

## 快速导航

### Repository层测试规范

**适用人群**：后端开发、数据工程师
**测试重点**：数据库操作、CRUD完整性、数据验证
**关键原则**：
- ✅ 必须使用真实MongoDB
- ❌ 严格禁止Mock Repository
- 使用`testutil.SetupTestDB()`进行数据库setup
- 覆盖率目标：≥80%

**快速开始**：
```go
func TestUserRepository_Create(t *testing.T) {
    db, cleanup := testutil.SetupTestDB(t)
    defer cleanup()

    repo := NewUserRepository(db)
    user := &User{Username: "test"}
    err := repo.Create(context.Background(), user)

    require.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

📖 **详细文档**：[repository_层测试规范.md](./repository_层测试规范.md)

---

### Service层测试规范

**适用人群**：后端开发、业务逻辑开发
**测试重点**：业务规则、数据验证、错误处理
**关键原则**：
- ✅ 必须使用Mock Repository
- ❌ 严格禁止使用真实Repository
- 使用testify/mock进行Mock
- 覆盖率目标：≥80%

**快速开始**：
```go
func TestUserService_Create(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    service := NewUserService(mockRepo)
    err := service.Create(context.Background(), &User{Username: "test"})

    require.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

📖 **详细文档**：[service_层测试规范.md](./service_层测试规范.md)

---

### API层测试规范

**适用人群**：后端开发、API开发、测试工程师
**测试重点**：HTTP接口、请求响应、权限控制
**关键原则**：
- ✅ 必须使用集成测试
- ❌ 严格禁止单独测试handler函数
- 使用TestHelper封装HTTP操作
- 覆盖率目标：≥60%

**快速开始**：
```go
func TestAuthAPI_Login(t *testing.T) {
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)
    user := helper.CreateTestUser(&User{Email: "test@example.com"})

    reqBody := map[string]interface{}{
        "email": "test@example.com",
        "password": "password123",
    }
    w := helper.DoRequest("POST", "/api/v1/auth/login", reqBody, nil)

    helper.AssertSuccess(w, 200, "登录失败")
}
```

📖 **详细文档**：[api_层测试规范.md](./api_层测试规范.md)

---

### E2E测试规范

**适用人群**：测试工程师、全栈开发、产品验证
**测试重点**：完整业务流程、用户体验、数据一致性
**关键原则**：
- ✅ 必须模拟真实用户场景
- ✅ 跨多个API和模块
- ✅ 验证完整数据链路
- 核心流程100%覆盖

**快速开始**：
```go
func TestE2E_PurchaseFlow(t *testing.T) {
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 1. 作者发布书籍
    author := env.CreateUser(&User{Username: "author"})
    book := env.CreateBook(&Book{AuthorID: author.ID})

    // 2. 读者充值并购买
    reader := env.CreateUser(&User{Username: "reader"})
    env.Recharge(reader.ID, 500)
    purchase := env.PurchaseBook(reader.ID, book.ID, readerToken)

    // 3. 验证完整流程
    env.VerifyPurchaseRecord(reader.ID, book.ID)
    env.VerifyReadingProgress(reader.ID, book.ID)
}
```

📖 **详细文档**：[e2e_测试规范.md](./e2e_测试规范.md)

## 测试分层决策树

```
我要测试什么？
│
├─ 数据库操作（CRUD、查询）
│  └─→ Repository层测试
│     - 使用真实MongoDB
│     - 验证数据正确性
│
├─ 业务逻辑（验证、计算、规则）
│  └─→ Service层测试
│     - Mock Repository
│     - 测试业务规则
│
├─ HTTP接口（请求、响应、权限）
│  └─→ API层测试
│     - 真实HTTP请求
│     - 验证接口行为
│
└─ 完整流程（跨多个模块）
   └─→ E2E测试
      - 模拟用户操作
      - 验证端到端流程
```

## 共同原则

### AAA测试模式

所有层级的测试都必须遵循AAA模式：

```go
func TestExample(t *testing.T) {
    // Arrange - 准备测试环境和数据
    // ...

    // Act - 执行被测试的操作
    // ...

    // Assert - 验证结果
    // ...
}
```

### Table-Driven测试

对于有多个场景的测试，使用表格驱动模式：

```go
tests := []struct {
    name    string
    input   InputType
    want    OutputType
    wantErr bool
}{
    {"场景1", input1, output1, false},
    {"场景2", input2, output2, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // 测试逻辑
    })
}
```

### 错误处理测试

每个层级都必须测试错误场景：

- Repository层：数据库错误、重复键、无效ID
- Service层：业务规则违反、权限不足、参数无效
- API层：401、403、404、500等错误响应
- E2E层：异常流程、失败重试、数据一致性

## 测试覆盖率目标

| 层级 | 覆盖率目标 | 重点 |
|------|-----------|------|
| Repository | ≥80% | CRUD操作 |
| Service | ≥80% | 业务逻辑 |
| API | ≥60% | HTTP接口 |
| E2E | 核心流程100% | 关键业务 |

## 检查清单

### Repository层测试
- [ ] 使用真实MongoDB
- [ ] CRUD操作完整覆盖
- [ ] 边缘案例测试（Not found, Invalid ID）
- [ ] 索引和查询验证
- [ ] 并发安全性测试

### Service层测试
- [ ] 使用Mock Repository
- [ ] 业务规则完整测试
- [ ] 成功和失败场景
- [ ] 参数验证测试
- [ ] 错误传播测试

### API层测试
- [ ] 使用集成测试
- [ ] 认证授权测试
- [ ] 请求响应验证
- [ ] 错误状态码测试
- [ ] 输入验证测试

### E2E测试
- [ ] 完整业务流程
- [ ] 数据链路验证
- [ ] 跨模块协作
- [ ] 用户体验验证
- [ ] 性能监控

## 常见问题

### Q1: 何时使用哪种测试？

**A**:
- **数据库操作有问题** → Repository层测试
- **业务逻辑有bug** → Service层测试
- **API接口异常** → API层测试
- **整个流程不通** → E2E测试

### Q2: 测试顺序是什么？

**A**: 自底向上测试
1. 先写Repository测试，确保数据层正确
2. 再写Service测试，验证业务逻辑
3. 然后写API测试，验证接口行为
4. 最后写E2E测试，验证完整流程

### Q3: 可以跳过某一层吗？

**A**: 不建议
- 跳过Repository测试 → 数据库问题难以发现
- 跳过Service测试 → 业务逻辑问题难以定位
- 跳过API测试 → 接口问题影响用户
- 跳过E2E测试 → 集成问题在线上爆发

## 相关文档

- [测试规范总导航](../README.md)
- [真实数据测试规范](../04_真实数据测试规范/README.md)
- [testify使用指南](../03_测试工具指南/testify使用指南.md)
- [集成测试详细规范](../02_测试类型规范/集成测试详细规范.md)

## 更新日志

- **v1.0** (2026-01-09) - 创建四层测试规范体系
