# E2E 测试框架文档

## 概述

本 E2E 测试框架采用**分层验证架构**，通过三个独立的测试层级全面验证系统功能：

- **Layer 1: 基础流程测试** (2-3分钟) - 验证核心业务流程的正确性
- **Layer 2: 数据一致性测试** (3-5分钟) - 验证跨模块数据一致性
- **Layer 3: 边界场景测试** (5-8分钟) - 验证边界条件和并发场景

**总测试时间**: 10-15分钟

## 目录结构

```
test/e2e/
├── framework/                      # E2E 测试框架
│   ├── environment.go             # 测试环境框架
│   ├── fixtures.go                # 测试数据夹具
│   ├── actions.go                 # 业务操作辅助
│   └── assertions.go              # E2E 断言辅助
├── data/                          # 数据工厂和验证器
│   ├── factory.go                 # 测试数据工厂
│   ├── scenarios.go               # 测试场景构建器
│   └── consistency_validator.go   # 数据一致性验证器
├── layer1_basic/                  # Layer 1: 基础流程测试
│   ├── auth_flow_test.go          # 认证流程测试
│   ├── reading_flow_test.go       # 阅读流程测试
│   ├── social_flow_test.go        # 社交流程测试
│   └── writing_flow_test.go       # 写作流程测试
├── layer2_consistency/            # Layer 2: 数据一致性测试
│   ├── user_reading_consistency_test.go      # 用户阅读一致性
│   ├── book_chapter_consistency_test.go      # 书籍章节一致性
│   ├── social_interaction_consistency_test.go # 社交互动一致性
│   └── helper.go                  # 测试辅助函数
├── layer3_boundary/               # Layer 3: 边界场景测试
│   └── concurrent_reading_test.go # 并发和边界测试
├── scenarios/                     # 完整场景测试（可选）
└── suite_test.go                  # 测试套件入口
```

## 快速开始

### 运行所有 E2E 测试

```bash
# 使用 Makefile（推荐）
make test-e2e

# 或直接使用 go test
go test -v -timeout 20m ./test/e2e/...
```

### 运行特定层级的测试

```bash
# 快速测试（仅 Layer 1）
make test-e2e-quick

# 标准测试（Layer 1 + Layer 2）
make test-e2e-standard

# 仅 Layer 1 - 基础流程测试
make test-e2e-layer1

# 仅 Layer 2 - 数据一致性测试
make test-e2e-layer2

# 仅 Layer 3 - 边界场景测试
make test-e2e-layer3
```

### 跳过 E2E 测试

```bash
# 运行单元测试时跳过 E2E
go test ./... -short
```

## 测试层级详解

### Layer 1: 基础流程测试 (2-3分钟)

**目标**: 验证核心业务流程的正确性

**测试用例**:
1. **认证流程** (auth_flow_test.go)
   - 创建测试用户
   - 用户登录
   - 验证用户存在
   - 获取用户详细信息

2. **阅读流程** (reading_flow_test.go)
   - 浏览书城首页
   - 查看书籍详情
   - 获取章节列表
   - 阅读章节内容
   - 保存阅读进度

3. **社交流程** (social_flow_test.go)
   - 发表评论
   - 收藏书籍
   - 点赞书籍
   - 查看收藏列表
   - 查看书籍评论

4. **写作流程** (writing_flow_test.go)
   - 创建写作项目
   - 验证项目存在

### Layer 2: 数据一致性测试 (3-5分钟)

**目标**: 验证跨模块数据一致性

**测试用例**:
1. **用户阅读一致性** (user_reading_consistency_test.go)
   - 验证用户阅读进度在用户模块和阅读器模块中一致
   - 验证阅读历史数据一致性
   - 跨模块数据验证

2. **书籍章节一致性** (book_chapter_consistency_test.go)
   - 验证书籍章节数量一致性
   - 验证章节内容完整性
   - 验证书籍-作者关系

3. **社交互动一致性** (social_interaction_consistency_test.go)
   - 验证评论数据在社交和用户模块中一致
   - 验证收藏数据在社交和阅读器模块中一致
   - 验证点赞数据一致性
   - 验证社交互动时间线

### Layer 3: 边界场景测试 (5-8分钟)

**目标**: 验证边界条件和并发场景

**测试用例**:
1. **并发阅读** (concurrent_reading_test.go)
   - 10个用户同时阅读同一本书
   - 验证并发操作后的数据一致性
   - 验证各用户阅读进度独立存储

2. **并发社交互动** (concurrent_reading_test.go)
   - 5个用户同时发表评论
   - 5个用户同时收藏书籍
   - 验证并发社交操作的数据完整性

3. **边界数据量** (concurrent_reading_test.go)
   - 创建包含100个章节的书籍
   - 测试获取大量章节列表（分页）
   - 验证分页功能的一致性

## 编写新的 E2E 测试

### 选择测试层级

- **Layer 1**: 新功能的核心业务流程
- **Layer 2**: 涉及多个模块的数据操作
- **Layer 3**: 性能敏感、并发场景或边界条件

### 基本结构

```go
package layer1_basic

import (
    "testing"
    e2e "Qingyu_backend/test/e2e/framework"
)

func TestMyNewFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    // 1. 初始化测试环境
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    // 2. 创建测试数据
    user := fixtures.CreateUser()

    // 3. 执行业务操作
    token := actions.Login(user.Username, "Test1234")

    // 4. 验证结果
    assertions.AssertUserExists(user.ID)
}
```

### 核心组件

#### TestEnvironment 测试环境

```go
env, cleanup := e2e.SetupTestEnvironment(t)
defer cleanup()

// 存储和获取测试数据
env.SetTestData("key", value)
value := env.GetTestData("key")

// 日志记录
env.LogSuccess("操作成功")
env.LogInfo("信息日志")
env.LogError("错误日志")
```

#### Fixtures 测试数据夹具

```go
fixtures := env.Fixtures()

// 创建用户
user := fixtures.CreateUser()
author := fixtures.CreateAdminUser()

// 创建书籍和章节
book := fixtures.CreateBook(author.ID)
chapter := fixtures.CreateChapter(book.ID.Hex())
```

#### Actions 业务操作辅助

```go
actions := env.Actions()

// 认证
token := actions.Login(username, "Test1234")

// 书城相关
homepage := actions.GetBookstoreHomepage()
rankings := actions.GetRankings("realtime")
bookDetail := actions.GetBookDetail(bookID)
searchResults := actions.SearchBooks("keyword")

// 阅读相关
chapterList := actions.GetChapterList(bookID, token)
chapterContent := actions.GetChapter(chapterID, token)
actions.StartReading(userID, bookID, chapterID, token)
readingProgress := actions.GetReadingProgress(userID, bookID, token)

// 社交互动
actions.AddComment(token, bookID, chapterID, content)
actions.CollectBook(token, bookID)
actions.LikeChapter(token, bookID)
comments := actions.GetBookComments(bookID, token)
collections := actions.GetReaderCollections(userID, token)

// 写作相关
actions.CreateProject(token, reqData)
```

#### ConsistencyValidator 数据一致性验证器

```go
validator := env.ConsistencyValidator()

// 验证用户数据
issues := validator.ValidateUserData(userID)

// 验证书籍数据
issues := validator.ValidateBookData(bookID)

// 过滤特定类型的问题
for _, issue := range issues {
    if issue.Severity == "error" {
        t.Errorf("数据一致性错误: %s", issue.Description)
    }
}
```

#### Assertions 断言辅助

```go
assertions := env.Assert()

// 用户相关
assertions.AssertUserExists(userID)

// 数据库验证
assertions.AssertCommentExists(userID, bookID)
assertions.AssertCollectionExists(userID, bookID)
assertions.AssertReadingProgress(userID, bookID)

// HTTP 响应
assertions.AssertResponseContains(response, "data")
```

## 数据隔离策略

所有测试数据使用 `e2e_test_` 前缀标记：

- 用户名：`e2e_test_user_xxxxx`
- 邮箱：`e2e_test_xxxxx@example.com`
- 书名：`e2e_test_book_xxxxx`

测试结束后自动清理所有带前缀的数据。

## 数据工厂

### TestDataFactory

创建各种测试实体：

```go
import "Qingyu_backend/test/e2e/data"

factory := data.NewTestDataFactory(t)
ctx := context.Background()

// 创建用户（带选项）
user := factory.CreateUser(ctx, data.UserOptions{
    Username: "test_user",
    VIPLevel: 1,
})

// 创建书籍
book := factory.CreateBook(ctx, authorID, data.BookOptions{
    Title: "我的书",
    Price: 100,
})

// 创建章节
chapter := factory.CreateChapter(ctx, bookID.Hex(), data.ChapterOptions{
    Title: "第一章",
})

// 清理测试数据
factory.Cleanup(ctx)
```

### ScenarioBuilder

构建复杂测试场景：

```go
builder := data.NewScenarioBuilder(factory, t)

// 构建阅读场景
scenario := builder.BuildReaderWithProgress()
user := scenario.User
book := scenario.Book
chapter := scenario.Chapter

// 构建社交互动场景
scenario := builder.BuildSocialInteraction()
```

## 最佳实践

1. **使用 `testing.Short()` 跳过**: E2E 测试默认在 `-short` 模式下跳过
2. **选择合适的测试层级**:
   - 新功能 → Layer 1
   - 跨模块操作 → Layer 2
   - 性能敏感场景 → Layer 3
3. **测试隔离**: 每个测试独立运行，不依赖其他测试
4. **完整验证**: 验证完整的数据链路，而不仅仅是 HTTP 响应
5. **使用数据工厂**: 优先使用 `fixtures` 方法，直接操作数据库
6. **清晰日志**: 使用 `env.LogSuccess()` 记录关键步骤
7. **清理数据**: 测试结束后自动清理测试数据

## 依赖项

- `github.com/stretchr/testify` - 测试断言
- `github.com/gin-gonic/gin` - Web 框架
- `go.mongodb.org/mongo-driver` - MongoDB 驱动
- `context` - 上下文管理

## 故障排除

### MongoDB 连接失败

确保 MongoDB 服务正在运行：
```bash
# 检查 MongoDB 状态
systemctl status mongod

# 或使用 Docker
docker ps | grep mongo
```

### Redis 警告

E2E 测试在 Redis 降级模式下运行，Redis 不可用时不会影响测试。

### 测试超时

增加测试超时时间：
```bash
go test -v -timeout 30m ./test/e2e/...
```

## 扩展阅读

- [E2E 测试规范](../../doc/standards/testing/01_测试层级规范/e2e_测试规范.md)
- [集成测试文档](../../test/integration/)
- [单元测试最佳实践](../../doc/standards/testing/)
