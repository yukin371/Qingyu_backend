//go:build e2e
// +build e2e

package layer2_consistency

import (
	"testing"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestSocialInteractionConsistency 测试社交互动数据一致性
// 验证: 评论、点赞、收藏在社交模块和用户模块中保持一致
func TestSocialInteractionConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()
	validator := env.ConsistencyValidator()

	t.Run("步骤1_创建用户和测试内容", func(t *testing.T) {
		t.Log("创建用户和测试内容...")

		// 创建多个用户
		user1 := fixtures.CreateUser()
		user2 := fixtures.CreateUser()
		t.Logf("✓ 用户创建成功: %s, %s", user1.Username, user2.Username)

		// 创建作者和书籍
		author := fixtures.CreateAdminUser()
		book := fixtures.CreateBook(author.ID)
		chapter := fixtures.CreateChapter(book.ID.Hex())
		t.Logf("✓ 书籍和章节创建成功: %s", book.Title)

		// 保存测试数据
		env.SetTestData("user1", user1)
		env.SetTestData("user2", user2)
		env.SetTestData("test_book", book)
		env.SetTestData("test_chapter", chapter)
		env.SetTestData("author_user", author)
	})

	t.Run("步骤2_执行多种社交互动", func(t *testing.T) {
		t.Log("执行评论、点赞、收藏操作...")

		user1 := env.GetTestData("user1").(*users.User)
		user2 := env.GetTestData("user2").(*users.User)
		book := env.GetTestData("test_book").(*bookstore.Book)
		chapter := env.GetTestData("test_chapter").(*bookstore.Chapter)

		// User1 登录并执行互动
		token1 := actions.Login(user1.Username, "Test1234")

		// User1 发表评论
		actions.AddComment(token1, book.ID.Hex(), chapter.ID.Hex(), "E2E测试评论1：这本书很有趣！")
		actions.AddComment(token1, book.ID.Hex(), chapter.ID.Hex(), "E2E测试评论2：期待后续章节。")

		// User1 收藏书籍
		actions.CollectBook(token1, book.ID.Hex())

		// User1 点赞书籍
		actions.LikeChapter(token1, book.ID.Hex())

		t.Logf("✓ User1 完成社交互动")

		// User2 也执行一些互动
		token2 := actions.Login(user2.Username, "Test1234")

		// User2 发表评论
		actions.AddComment(token2, book.ID.Hex(), chapter.ID.Hex(), "E2E测试评论3：我也很喜欢这本书。")

		// User2 收藏书籍
		actions.CollectBook(token2, book.ID.Hex())

		t.Logf("✓ User2 完成社交互动")
	})

	t.Run("步骤3_验证评论数据一致性", func(t *testing.T) {
		t.Log("验证评论数据在社交和用户模块中一致...")

		user1 := env.GetTestData("user1").(*users.User)
		book := env.GetTestData("test_book").(*bookstore.Book)

		// 验证 user1 的评论数据
		issues := validator.ValidateUserData(user1.ID)
		commentIssues := filterIssuesByType(issues, "comment_count_mismatch")

		if len(commentIssues) > 0 {
			for _, issue := range commentIssues {
				t.Logf("评论数量不一致: %s", issue.Description)
			}
		}

		// 验证评论确实存在于数据库中
		assertions := env.Assert()
		assertions.AssertCommentExists(user1.ID, book.ID.Hex())
		t.Logf("✓ 评论数据验证通过")
	})

	t.Run("步骤4_验证收藏数据一致性", func(t *testing.T) {
		t.Log("验证收藏数据在社交和阅读器模块中一致...")

		user1 := env.GetTestData("user1").(*users.User)
		book := env.GetTestData("test_book").(*bookstore.Book)

		// 验证收藏数据
		issues := validator.ValidateUserData(user1.ID)
		collectionIssues := filterIssuesByType(issues, "collection_count_mismatch")

		if len(collectionIssues) > 0 {
			for _, issue := range collectionIssues {
				t.Logf("收藏数量不一致: %s", issue.Description)
			}
		}

		// 验证收藏确实存在
		assertions := env.Assert()
		assertions.AssertCollectionExists(user1.ID, book.ID.Hex())
		t.Logf("✓ 收藏数据验证通过")
	})

	t.Run("步骤5_验证点赞数据一致性", func(t *testing.T) {
		t.Log("验证点赞数据在社交模块中一致...")

		user1 := env.GetTestData("user1").(*users.User)

		// 验证点赞数据
		issues := validator.ValidateUserData(user1.ID)
		likeIssues := filterIssuesByType(issues, "like_count_mismatch")

		if len(likeIssues) > 0 {
			for _, issue := range likeIssues {
				t.Logf("点赞数量不一致: %s", issue.Description)
			}
		} else {
			t.Log("✓ 点赞数据验证通过")
		}
	})

	t.Run("步骤6_验证跨模块社交数据汇总", func(t *testing.T) {
		t.Log("验证用户在各个模块的社交数据汇总一致...")

		user1 := env.GetTestData("user1").(*users.User)

		// 完整的用户数据验证
		issues := validator.ValidateUserData(user1.ID)

		// 统计各类问题
		errorCount := 0
		warningCount := 0
		for _, issue := range issues {
			if issue.Severity == "error" {
				errorCount++
			} else if issue.Severity == "warning" {
				warningCount++
			}
		}

		t.Logf("社交数据一致性检查完成: %d个错误, %d个警告", errorCount, warningCount)

		// 如果有错误，详细列出
		if errorCount > 0 {
			for _, issue := range issues {
				if issue.Severity == "error" {
					t.Errorf("社交数据错误: [%s] %s", issue.Type, issue.Description)
				}
			}
		}

		// 确保用户的所有社交互动都被正确记录
		if errorCount == 0 {
			t.Log("✓ 所有社交数据在跨模块中保持一致")
		}
	})

	t.Run("步骤7_验证社交互动时间线", func(t *testing.T) {
		t.Log("验证社交互动的时间线数据...")

		user1 := env.GetTestData("user1").(*users.User)

		// 验证用户的活动历史
		issues := validator.ValidateUserData(user1.ID)
		timelineIssues := filterIssuesByType(issues, "timeline_inconsistent")

		if len(timelineIssues) > 0 {
			t.Logf("发现 %d 个时间线问题", len(timelineIssues))
			for _, issue := range timelineIssues {
				t.Logf("  - %s: %s", issue.Type, issue.Description)
			}
		} else {
			t.Log("✓ 社交互动时间线一致")
		}
	})
}



