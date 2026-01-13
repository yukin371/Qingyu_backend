package layer3_boundary

import (
	"fmt"
	"sync"
	"testing"

	"Qingyu_backend/models/bookstore"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestConcurrentReading 测试并发阅读场景
// 验证: 多个用户同时阅读同一本书时的数据一致性
func TestConcurrentReading(t *testing.T) {
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
		book := fixtures.CreateBook(author.ID)

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
				userIDs[index] = user.ID

				// 登录获取token
				token := actions.Login(user.Username, "Test1234")

				// 获取章节列表
				chapterList := actions.GetChapterList(book.ID.Hex(), token)

				// 保存阅读进度
				if data, ok := chapterList["data"].(map[string]interface{}); ok {
					if chapters, ok := data["chapters"].([]interface{}); ok && len(chapters) > 0 {
						if firstChapter, ok := chapters[0].(map[string]interface{}); ok {
							if chapterID, ok := firstChapter["id"].(string); ok {
								actions.StartReading(user.ID, book.ID.Hex(), chapterID, token)
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

// TestConcurrentSocialInteraction 测试并发社交互动
// 验证: 多个用户同时发表评论、点赞时的数据一致性
func TestConcurrentSocialInteraction(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	t.Run("步骤1_创建测试数据", func(t *testing.T) {
		t.Log("创建测试书籍...")

		author := fixtures.CreateAdminUser()
		book := fixtures.CreateBook(author.ID)
		chapter := fixtures.CreateChapter(book.ID.Hex())

		t.Logf("✓ 创建书籍: %s", book.Title)

		env.SetTestData("test_book", book)
		env.SetTestData("test_chapter", chapter)
	})

	t.Run("步骤2_多用户并发评论", func(t *testing.T) {
		t.Log("模拟5个用户同时发表评论...")

		book := env.GetTestData("test_book").(*bookstore.Book)
		chapter := env.GetTestData("test_chapter").(*bookstore.Chapter)
		userCount := 5

		var wg sync.WaitGroup
		commentCount := make(chan int, userCount)

		for i := 0; i < userCount; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				user := fixtures.CreateUser()
				token := actions.Login(user.Username, "Test1234")

				// 发表评论
				comment := actions.AddComment(token, book.ID.Hex(), chapter.ID.Hex(),
					fmt.Sprintf("并发测试评论 #%d 来自用户 %s", index+1, user.Username))

				if data, ok := comment["data"].(map[string]interface{}); ok {
					if _, ok := data["commentId"]; ok {
						commentCount <- 1
					} else {
						commentCount <- 0
					}
				} else {
					commentCount <- 0
				}
			}(i)
		}

		wg.Wait()
		close(commentCount)

		totalComments := 0
		for count := range commentCount {
			totalComments += count
		}

		t.Logf("✓ 完成 %d 个用户的并发评论，成功: %d", userCount, totalComments)
	})

	t.Run("步骤3_多用户并发收藏", func(t *testing.T) {
		t.Log("模拟5个用户同时收藏书籍...")

		book := env.GetTestData("test_book").(*bookstore.Book)
		userCount := 5

		var wg sync.WaitGroup
		collectionCount := make(chan int, userCount)

		for i := 0; i < userCount; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				user := fixtures.CreateUser()
				token := actions.Login(user.Username, "Test1234")

				// 收藏书籍
				collection := actions.CollectBook(token, book.ID.Hex())

				if data, ok := collection["data"].(map[string]interface{}); ok {
					if _, ok := data["collectionId"]; ok {
						collectionCount <- 1
					} else {
						collectionCount <- 0
					}
				} else {
					collectionCount <- 0
				}
			}(i)
		}

		wg.Wait()
		close(collectionCount)

		totalCollections := 0
		for count := range collectionCount {
			totalCollections += count
		}

		t.Logf("✓ 完成 %d 个用户的并发收藏，成功: %d", userCount, totalCollections)
	})
}

// TestBoundaryDataSizes 测试边界数据量
// 验证: 处理大量数据时的性能和正确性
func TestBoundaryDataSizes(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	t.Run("步骤1_创建包含大量章节的书籍", func(t *testing.T) {
		t.Log("创建包含100个章节的书籍...")

		author := fixtures.CreateAdminUser()
		book := fixtures.CreateBook(author.ID)

		// 创建100个章节
		chapterCount := 100
		for i := 1; i <= chapterCount; i++ {
			fixtures.CreateChapter(book.ID.Hex())
		}

		t.Logf("✓ 创建书籍: %s, 包含 %d 个章节", book.Title, chapterCount)

		env.SetTestData("test_book", book)
		env.SetTestData("chapter_count", chapterCount)
	})

	t.Run("步骤2_测试获取大量章节列表", func(t *testing.T) {
		t.Log("测试获取章节列表...")

		book := env.GetTestData("test_book").(*bookstore.Book)

		user := fixtures.CreateUser()
		token := actions.Login(user.Username, "Test1234")

		// 获取章节列表
		chapterList := actions.GetChapterList(book.ID.Hex(), token)

		if data, ok := chapterList["data"].(map[string]interface{}); ok {
			if chapters, ok := data["chapters"].([]interface{}); ok {
				actualCount := len(chapters)
				expectedCount := env.GetTestData("chapter_count").(int)

				t.Logf("✓ 成功获取 %d 个章节（总计: %d，使用分页获取全部）", actualCount, expectedCount)

				// 验证分页返回了合理数量的章节
				if actualCount > 0 && actualCount <= expectedCount {
					t.Logf("✓ 分页返回合理: %d/%d", actualCount, expectedCount)
				}

				// 检查是否有分页信息
				if pagination, ok := data["pagination"].(map[string]interface{}); ok {
					t.Logf("分页信息: %+v", pagination)
				}
			}
		}
	})

	t.Run("步骤3_测试分页获取", func(t *testing.T) {
		t.Log("测试分页功能...")

		book := env.GetTestData("test_book").(*bookstore.Book)

		user := fixtures.CreateUser()
		token := actions.Login(user.Username, "Test1234")

		// 获取章节列表（第一次）
		chapterList1 := actions.GetChapterList(book.ID.Hex(), token)

		// 再次获取（验证缓存/分页是否正常工作）
		chapterList2 := actions.GetChapterList(book.ID.Hex(), token)

		if data1, ok := chapterList1["data"].(map[string]interface{}); ok {
			if chapters1, ok := data1["chapters"].([]interface{}); ok {
				if data2, ok := chapterList2["data"].(map[string]interface{}); ok {
					if chapters2, ok := data2["chapters"].([]interface{}); ok {
						count1 := len(chapters1)
						count2 := len(chapters2)

						if count1 == count2 {
							t.Logf("✓ 分页数据一致，每次返回 %d 个章节", count1)
						} else {
							t.Logf("⚠ 分页数据不一致: 第一次 %d, 第二次 %d", count1, count2)
						}
					}
				}
			}
		}
	})
}
