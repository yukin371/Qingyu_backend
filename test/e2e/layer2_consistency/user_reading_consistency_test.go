//go:build e2e
// +build e2e

package layer2_consistency

import (
	"testing"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestUserReadingConsistency 测试用户阅读数据一致性
// 验证: 用户的阅读进度在用户模块和阅读器模块中保持一致
func TestUserReadingConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()
	validator := env.ConsistencyValidator()

	t.Run("步骤1_创建用户和阅读数据", func(t *testing.T) {
		t.Log("创建测试用户和阅读数据...")

		// 创建用户
		user := fixtures.CreateUser()
		t.Logf("✓ 用户创建成功: %s", user.Username)

		// 创建作者和书籍
		author := fixtures.CreateAdminUser()
		book := fixtures.CreateBook(author.ID)
		chapter := fixtures.CreateChapter(book.ID.Hex())

		t.Logf("✓ 书籍和章节创建成功: %s", book.Title)

		// 保存测试数据
		env.SetTestData("test_user", user)
		env.SetTestData("test_book", book)
		env.SetTestData("test_chapter", chapter)
		env.SetTestData("author_user", author)
	})

	t.Run("步骤2_执行阅读操作并保存进度", func(t *testing.T) {
		t.Log("执行阅读操作...")

		user := env.GetTestData("test_user").(*users.User)
		book := env.GetTestData("test_book").(*bookstore.Book)
		chapter := env.GetTestData("test_chapter").(*bookstore.Chapter)

		token := actions.Login(user.Username, "Test1234")

		// 模拟多次阅读操作
		for i := 1; i <= 3; i++ {
			// 更新阅读进度
			actions.StartReading(user.ID, book.ID.Hex(), chapter.ID.Hex(), token)

			// 获取章节内容（模拟阅读）
			actions.GetChapter(chapter.ID.Hex(), token)
		}

		t.Logf("✓ 完成3次阅读操作")
	})

	t.Run("步骤3_验证阅读进度数据一致性", func(t *testing.T) {
		t.Log("验证阅读进度数据一致性...")

		user := env.GetTestData("test_user").(*users.User)

		// 使用一致性验证器检查用户数据
		issues := validator.ValidateUserData(user.ID)

		// 统计问题数量
		errorCount := 0
		warningCount := 0
		for _, issue := range issues {
			if issue.Severity == "error" {
				errorCount++
			} else if issue.Severity == "warning" {
				warningCount++
			}
		}

		t.Logf("数据一致性检查完成: %d个错误, %d个警告", errorCount, warningCount)

		// 对于这个测试，我们不应该有错误
		if errorCount > 0 {
			for _, issue := range issues {
				if issue.Severity == "error" {
					t.Errorf("数据一致性错误: %s - %s", issue.Type, issue.Description)
				}
			}
		}

		// 确认阅读进度存在
		readingProgressIssues := filterIssuesByType(issues, "reading_progress_missing")
		if len(readingProgressIssues) > 0 {
			t.Error("发现阅读进度缺失问题，这可能表明数据同步失败")
		}
	})

	t.Run("步骤4_验证阅读历史数据一致性", func(t *testing.T) {
		t.Log("验证阅读历史数据一致性...")

		user := env.GetTestData("test_user").(*users.User)

		// 验证阅读历史
		issues := validator.ValidateUserData(user.ID)
		historyIssues := filterIssuesByType(issues, "reading_history_inconsistent")

		if len(historyIssues) > 0 {
			t.Logf("阅读历史问题: %d个", len(historyIssues))
			for _, issue := range historyIssues {
				t.Logf("  - %s: %s", issue.Type, issue.Description)
			}
		} else {
			t.Log("✓ 阅读历史数据一致")
		}
	})

	t.Run("步骤5_跨模块数据验证", func(t *testing.T) {
		t.Log("执行跨模块数据验证...")

		user := env.GetTestData("test_user").(*users.User)

		// 获取用户详细信息（从用户模块）
		userDetail := actions.GetUser(user.ID)

		// 使用一致性验证器进行跨模块验证
		issues := validator.ValidateUserData(user.ID)

		// 验证没有严重的数据不一致问题
		errorCount := 0
		for _, issue := range issues {
			if issue.Severity == "error" {
				errorCount++
			}
		}

		if errorCount == 0 {
			t.Logf("✓ 跨模块数据验证完成，用户状态: %s", userDetail.Status)
		} else {
			t.Errorf("发现 %d 个数据一致性问题", errorCount)
		}
	})
}



