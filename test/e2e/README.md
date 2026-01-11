# E2E 测试框架文档

## 概述

本 E2E 测试框架遵循 `doc/standards/testing/01_测试层级规范/e2e_测试规范.md` 规范，用于测试完整的端到端业务流程。

## 目录结构

```
test/e2e/
├── framework/                      # E2E 测试框架
│   ├── environment.go             # 测试环境框架
│   ├── fixtures.go                # 测试数据夹具
│   ├── actions.go                 # 业务操作辅助
│   └── assertions.go              # E2E 断言辅助
└── scenarios/                      # 测试场景
    └── complete_user_journey_test.go  # 完整用户旅程测试
```

## 运行 E2E 测试

```bash
# 运行所有 E2E 测试
go test ./test/e2e/... -v

# 运行完整用户旅程测试
go test ./test/e2e/scenarios/complete_user_journey_test.go -v

# 跳过 E2E 测试
go test ./... -short
```

## 编写新的 E2E 测试

### 基本结构

```go
package e2e_test

import (
    "testing"
    "github.com/stretchr/testify/require"

    "Qingyu_backend/test/e2e/framework"
)

func TestE2E_MyScenario(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    // 1. 初始化测试环境
    env, cleanup := framework.SetupTestEnvironment(t)
    defer cleanup()

    // 2. 获取辅助工具
    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    // 3. 创建测试数据
    user := fixtures.CreateUser(
        framework.WithUsername("test_user"),
        framework.WithBalance(1000),
    )

    // 4. 执行业务操作
    token := actions.Login(user.Email, "password123")

    // 5. 验证结果
    assertions.AssertUserBalance(user.ID, 1000)
}
```

### 核心组件

#### TestEnvironment 测试环境

`SetupTestEnvironment(t)` 初始化完整的测试环境：
- 加载配置
- 初始化服务容器
- 设置 HTTP 路由
- 创建清理函数

```go
env, cleanup := framework.SetupTestEnvironment(t)
defer cleanup()
```

#### Fixtures 测试数据夹具

创建测试数据，自动添加 `e2e_test_` 前缀：

```go
fixtures := env.Fixtures()

// 创建用户
user := fixtures.CreateUser(
    framework.WithUsername("test_user"),
    framework.WithBalance(1000),
    framework.WithVIP(true),
)

// 创建书籍
book := fixtures.CreateBook(author.ID,
    framework.WithBookTitle("我的书"),
    framework.WithBookPrice(100),
)

// 创建章节
chapter := fixtures.CreateChapter(book.ID,
    framework.WithChapterTitle("第一章"),
    framework.WithChapterFree(false),
)
```

#### Actions 业务操作辅助

执行业务操作，通过 API 或直接操作数据库：

```go
actions := env.Actions()

// 认证
token := actions.Login(email, password)

// 书城相关
homepage := actions.GetBookstoreHomepage()
rankings := actions.GetRankings("realtime")
bookDetail := actions.GetBookDetail(bookID)
chapterContent := actions.GetChapter(chapterID, token)

// 阅读相关
actions.StartReading(userID, bookID, chapterID, token)
history := actions.GetReadingHistory(userID)

// 社交互动
actions.AddComment(token, bookID, chapterID, "好书！")
actions.CollectBook(token, bookID)
actions.LikeChapter(token, chapterID)
actions.AddBookmark(token, bookID, chapterID, 100)

// 写作相关
actions.CreateProject(token, reqData)
```

#### Assertions 断言辅助

验证测试结果：

```go
assertions := env.Assert()

// 用户相关
assertions.AssertUserBalance(userID, 1000)
assertions.AssertUserVIP(userID, true)

// HTTP 响应
assertions.AssertHTTPSuccess(statusCode, 200)
assertions.AssertResponseContains(response, "data")

// 数据库
assertions.AssertCollectionCount("users", 10)
assertions.AssertDocumentExists("purchases", filter)

// 业务流程
assertions.AssertPurchaseRecord(userID, bookID)
assertions.AssertReadingProgress(userID, bookID)
```

## 数据隔离策略

所有测试数据使用 `e2e_test_` 前缀标记：

- 用户名：`e2e_test_user_xxxxx`
- 邮箱：`e2e_test_xxxxx@example.com`
- 书名：`e2e_test_book_xxxxx`
- 章节名：`e2e_test_chapter_xxxxx`

测试结束后自动清理所有带前缀的数据。

## 测试场景

### 完整用户旅程 (TestE2E_CompleteUserJourney)

流程：注册 → 浏览书城 → 充值 → 阅读 → 互动 → 成为作者

### 书籍购买流程 (TestE2E_BookPurchaseFlow)

流程：作者发布书籍 → 读者充值 → 购买书籍 → 阅读付费内容

### VIP阅读流程 (TestE2E_VIPReadingFlow)

流程：用户订阅VIP → 免费阅读付费书籍

## 最佳实践

1. **使用 `testing.Short()` 跳过**：E2E 测试默认在 `-short` 模式下跳过
2. **命名规范**：所有 E2E 测试用例以 `TestE2E_` 开头
3. **测试隔离**：每个测试独立运行，不依赖其他测试
4. **完整验证**：验证完整的数据链路，而不仅仅是 HTTP 响应
5. **清晰日志**：使用 `env.LogSuccess()` 记录关键步骤

## 依赖项

- `github.com/stretchr/testify` - 测试断言
- `github.com/gin-gonic/gin` - Web 框架
- `go.mongodb.org/mongo-driver` - MongoDB 驱动

## 扩展阅读

- [E2E 测试规范](../../doc/standards/testing/01_测试层级规范/e2e_测试规范.md)
- [集成测试文档](../../test/integration/)
