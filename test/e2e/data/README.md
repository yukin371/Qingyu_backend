# E2E测试数据工厂

这个包提供了用于E2E测试的数据生成工具，包括测试数据工厂和场景构建器。

## 功能特性

### 数据工厂 (TestDataFactory)

提供创建各种测试实体的方法：

- **CreateUser**: 创建测试用户
  - 支持自定义用户名、邮箱、VIP等级、余额、角色
  - 自动生成密码 (Test1234)
  - 自动设置默认值

- **CreateUsers**: 批量创建用户
  - 支持批量创建多个用户
  - 自动生成唯一的用户名和邮箱

- **CreateBook**: 创建测试书籍
  - 支持自定义标题、作者ID、价格、是否免费、分类等
  - 自动设置默认值

- **CreateChapter**: 创建测试章节
  - 支持自定义书籍ID、章节号、是否免费
  - 自动创建章节内容

- **CreateComment**: 创建测试评论
  - 支持自定义作者ID、目标ID、目标类型、内容
  - 自动设置默认值

- **CreateCollection**: 创建测试收藏
  - 支持自定义用户ID和书籍ID

- **Cleanup**: 清理测试数据
  - 根据前缀清理测试数据

### 场景构建器 (ScenarioBuilder)

提供预定义的测试场景：

- **BuildReaderWithProgress**: 构建阅读进度场景
  - 创建读者、作者、书籍和章节
  - 模拟真实的阅读进度

- **BuildSocialInteraction**: 构建社交互动场景
  - 创建作者、书籍和多个互动用户
  - 可用于测试评论、点赞等功能

- **BuildPaidContent**: 构建付费内容场景
  - 创建作者、免费用户、VIP用户
  - 创建包含免费章节和付费章节的书籍

## 使用方法

### 基本用法

```go
package data_test

import (
    "testing"
    "Qingyu_backend/test/e2e/data"
)

func TestExample(t *testing.T) {
    // 1. 设置测试环境
    data.SetupTestEnvironment(t)

    // 2. 创建数据工厂
    factory := data.NewTestDataFactory(t)

    // 3. 创建用户
    user := factory.CreateUser(data.UserOptions{
        Username: "test_user",
        VIPLevel: 1,
    })

    // 4. 创建书籍
    book := factory.CreateBook(data.BookOptions{
        Title:    "测试书籍",
        AuthorID: user.ID,
        Price:    100,
        IsFree:   false,
    })

    // 5. 创建章节
    chapter := factory.CreateChapter(book.ID.Hex(), 1, true)

    // 6. 清理测试数据
    factory.Cleanup("test_")
}
```

### 使用场景构建器

```go
func TestScenarioExample(t *testing.T) {
    data.SetupTestEnvironment(t)

    // 创建场景构建器
    builder := data.NewScenarioBuilder(t)

    // 构建付费内容场景
    scenario := builder.BuildPaidContent()

    // 使用场景中的数据
    freeUser := scenario.FreeUser
    vipUser := scenario.VIPUser
    paidBook := scenario.PaidBook

    // 执行测试逻辑...
}
```

## 数据前缀规则

所有测试数据都使用 `e2e_` 前缀，便于识别和清理：

- 用户名: `e2e_user_*` 或自定义前缀
- 邮箱: `e2e_*@example.com`
- 书籍标题: `e2e_book_*` 或自定义

## 注意事项

1. **必须先调用 SetupTestEnvironment**: 在所有测试前必须调用此函数来初始化数据库连接和服务
2. **数据隔离**: 每个测试应该使用唯一的用户名和标题，避免冲突
3. **清理数据**: 测试结束后应该调用 Cleanup 方法清理测试数据
4. **数据库要求**: 需要运行MongoDB服务（测试配置在 config/config.test.yaml）

## 示例测试

查看 `factory_test.go` 文件以获取更多使用示例。
