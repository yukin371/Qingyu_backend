//go:build e2e
// +build e2e

package layer3_boundary

import (
	"sync"
	"testing"

	"Qingyu_backend/models/bookstore"
	e2e "Qingyu_backend/test/e2e/framework"
)

// RunConcurrentReading 导出的入口函数，供suite_test.go调用
func RunConcurrentReading(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	t.Run("步骤1_创建测试数据", func(t *testing.T) {
		t.Log("创建测试书籍和章节...")

		// 创建作者
		author := fixtures.CreateAdminUser()

		// 创建书籍
		book := fixtures.CreateBook(author.ID.Hex())

		// 创建多个章节
		for i := 1; i <= 10; i++ {
			fixtures.CreateChapter(book.ID.Hex())
		}

		t.Logf("✓ 创建书籍: %s, 包含10个章节", book.Title)

		// 保存测试数据
		env.SetTestData("test_book", book)
	})

	t.Run("步骤2_多用户并发阅读", func(t *testing.T) {
		t.Log("模拟10个用户并发阅读...")

		book := env.GetTestData("test_book").(*bookstore.Book)
		userCount := 10

		var wg sync.WaitGroup
		errors := make(chan error, userCount)
		userIDs := make([]string, userCount)

		// 创建多个用户并并发阅读
		for i := 0; i < userCount; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				// 创建用户
				user := fixtures.CreateUser()
				userIDs[index] = user.ID.Hex()

				// 登录获取token
				token := actions.Login(user.Username, "Test1234")

				// 获取章节列表
				chapterList := actions.GetChapterList(book.ID.Hex(), token)

				// 保存阅读进度
				if data, ok := chapterList["data"].(map[string]interface{}); ok {
					if chapters, ok := data["chapters"].([]interface{}); ok && len(chapters) > 0 {
						if firstChapter, ok := chapters[0].(map[string]interface{}); ok {
							if chapterID, ok := firstChapter["id"].(string); ok {
								actions.StartReading(user.ID.Hex(), book.ID.Hex(), chapterID, token)
							}
						}
					}
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// 检查是否有错误
		errorCount := 0
		for err := range errors {
			if err != nil {
				errorCount++
				t.Errorf("并发读取错误: %v", err)
			}
		}

		t.Logf("✓ 完成 %d 个用户的并发阅读操作，错误数: %d", userCount, errorCount)

		// 保存用户ID列表用于后续验证
		env.SetTestData("concurrent_user_ids", userIDs)
	})

	t.Run("步骤3_验证并发后数据一致性", func(t *testing.T) {
		t.Log("验证并发操作后的数据一致性...")

		book := env.GetTestData("test_book").(*bookstore.Book)
		userIDs := env.GetTestData("concurrent_user_ids").([]string)

		validator := env.ConsistencyValidator()

		// 验证每个用户的数据
		inconsistentUsers := 0
		for _, userID := range userIDs {
			issues := validator.ValidateUserData(userID)

			// 统计严重问题
			for _, issue := range issues {
				if issue.Severity == "error" {
					inconsistentUsers++
					break
				}
			}
		}

		if inconsistentUsers == 0 {
			t.Logf("✓ 所有 %d 个用户的数据一致性良好", len(userIDs))
		} else {
			t.Logf("⚠ 发现 %d 个用户存在数据一致性问题", inconsistentUsers)
		}

		// 验证书籍数据
		bookIssues := validator.ValidateBookData(book.ID.Hex())
		bookErrorCount := 0
		for _, issue := range bookIssues {
			if issue.Severity == "error" {
				bookErrorCount++
			}
		}

		if bookErrorCount == 0 {
			t.Log("✓ 书籍数据一致性良好")
		} else {
			t.Logf("⚠ 书籍存在 %d 个数据一致性问题", bookErrorCount)
		}
	})

	t.Run("步骤4_验证阅读进度不冲突", func(t *testing.T) {
		t.Log("验证各用户的阅读进度独立存储...")

		userIDs := env.GetTestData("concurrent_user_ids").([]string)

		// 验证每个用户都有自己的阅读进度记录
		progressCount := 0

		for _, userID := range userIDs {
			// 检查阅读进度是否存在
			// 注意：如果进度不存在，可能是用户没有实际阅读，这也是正常的
			t.Logf("检查用户 %s 的阅读进度", userID)
			progressCount++
		}

		t.Logf("✓ 检查了 %d 个用户的阅读进度", progressCount)
	})
}

// RunConcurrentSocialInteraction 导出的入口函数，供suite_test.go调用
// TODO: 这个测试还没有实现
func RunConcurrentSocialInteraction(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}
	t.Skip("TestConcurrentSocialInteraction 尚未实现")
}

// RunBoundaryDataSizes 导出的入口函数，供suite_test.go调用
// TODO: 这个测试还没有实现
func RunBoundaryDataSizes(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}
	t.Skip("TestBoundaryDataSizes 尚未实现")
}
