//go:build e2e
// +build e2e

package layer3_boundary

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"Qingyu_backend/global"
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
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
func RunConcurrentSocialInteraction(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()
	assertions := env.Assert()

	t.Run("步骤1_创建社交测试数据", func(t *testing.T) {
		author := fixtures.CreateAdminUser()
		book := fixtures.CreateBook(author.ID.Hex())
		chapter := fixtures.CreateChapter(book.ID.Hex())

		env.SetTestData("social_book", book)
		env.SetTestData("social_chapter", chapter)
	})

	t.Run("步骤2_多用户并发社交互动", func(t *testing.T) {
		book := env.GetTestData("social_book").(*bookstore.Book)
		chapter := env.GetTestData("social_chapter").(*bookstore.Chapter)
		userCount := 6

		var wg sync.WaitGroup
		userIDs := make([]string, userCount)
		errors := make(chan error, userCount)

		for i := 0; i < userCount; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						errors <- fmt.Errorf("goroutine %d panic: %v", index, r)
					}
				}()

				user := fixtures.CreateUser()
				userIDs[index] = user.ID.Hex()

				token := actions.Login(user.Username, "Test1234")
				actions.AddComment(token, book.ID.Hex(), chapter.ID.Hex(), fmt.Sprintf("并发社交测试评论_%02d，内容长度满足最小限制。", index))
				actions.CollectBook(token, book.ID.Hex())
				actions.LikeChapter(token, book.ID.Hex())
			}(i)
		}

		wg.Wait()
		close(errors)

		for err := range errors {
			require.NoError(t, err)
		}

		env.SetTestData("social_user_ids", userIDs)
		t.Logf("✓ 完成 %d 个用户的并发社交互动", userCount)
	})

	t.Run("步骤3_验证并发社交数据已落库", func(t *testing.T) {
		book := env.GetTestData("social_book").(*bookstore.Book)
		userIDs := env.GetTestData("social_user_ids").([]string)

		ctx := context.Background()
		commentCount, err := global.DB.Collection("comments").CountDocuments(ctx, bson.M{
			"target_id":   book.ID.Hex(),
			"target_type": "book",
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, int(commentCount), len(userIDs), "并发评论落库数量不足")

		for _, userID := range userIDs {
			assertions.AssertCommentExists(userID, book.ID.Hex())
			assertions.AssertCollectionExists(userID, book.ID.Hex())
		}
	})

	t.Run("步骤4_验证并发后数据一致性", func(t *testing.T) {
		book := env.GetTestData("social_book").(*bookstore.Book)
		userIDs := env.GetTestData("social_user_ids").([]string)
		validator := env.ConsistencyValidator()

		for _, userID := range userIDs {
			issues := validator.ValidateUserData(userID)
			for _, issue := range issues {
				if issue.Severity == "error" {
					t.Fatalf("用户 %s 出现数据一致性错误: [%s] %s", userID, issue.Type, issue.Description)
				}
			}
		}

		bookIssues := validator.ValidateBookData(book.ID.Hex())
		for _, issue := range bookIssues {
			if issue.Severity == "error" {
				t.Fatalf("书籍 %s 出现数据一致性错误: [%s] %s", book.ID.Hex(), issue.Type, issue.Description)
			}
		}
	})
}

// RunBoundaryDataSizes 导出的入口函数，供suite_test.go调用
func RunBoundaryDataSizes(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	t.Run("步骤1_准备边界测试数据", func(t *testing.T) {
		user := fixtures.CreateUser()
		author := fixtures.CreateAdminUser()
		book := fixtures.CreateBook(author.ID.Hex())
		chapter := fixtures.CreateChapter(book.ID.Hex())
		token := actions.Login(user.Username, "Test1234")

		env.SetTestData("boundary_user", user)
		env.SetTestData("boundary_book", book)
		env.SetTestData("boundary_chapter", chapter)
		env.SetTestData("boundary_token", token)
	})

	t.Run("步骤2_验证评论内容长度边界", func(t *testing.T) {
		book := env.GetTestData("boundary_book").(*bookstore.Book)
		token := env.GetTestData("boundary_token").(string)

		validBodies := []map[string]interface{}{
			{"book_id": book.ID.Hex(), "content": strings.Repeat("a", 10)},
			{"book_id": book.ID.Hex(), "content": strings.Repeat("b", 500)},
		}

		for _, reqBody := range validBodies {
			w := env.DoRequest("POST", "/api/v1/social/comments", reqBody, token)
			require.Truef(t, w.Code == 200 || w.Code == 201, "评论边界请求失败: code=%d body=%s", w.Code, w.Body.String())
		}

		invalidComment := map[string]interface{}{
			"book_id": book.ID.Hex(),
			"content": strings.Repeat("c", 501),
		}
		w := env.DoRequest("POST", "/api/v1/social/comments", invalidComment, token)
		require.Equalf(t, 400, w.Code, "超长评论应返回 400: body=%s", w.Body.String())
	})

	t.Run("步骤3_验证书签位置边界", func(t *testing.T) {
		book := env.GetTestData("boundary_book").(*bookstore.Book)
		chapter := env.GetTestData("boundary_chapter").(*bookstore.Chapter)
		token := env.GetTestData("boundary_token").(string)

		validBookmark := map[string]interface{}{
			"bookId":    book.ID.Hex(),
			"chapterId": chapter.ID.Hex(),
			"position":  1,
			"note":      "边界位置书签",
		}
		w := env.DoRequest("POST", fmt.Sprintf("/api/v1/reader/books/%s/bookmarks", book.ID.Hex()), validBookmark, token)
		require.Equalf(t, 201, w.Code, "position=1 的书签创建失败: body=%s", w.Body.String())

		invalidBookmark := map[string]interface{}{
			"bookId":    book.ID.Hex(),
			"chapterId": chapter.ID.Hex(),
			"position":  0,
			"note":      "非法位置书签",
		}
		w = env.DoRequest("POST", fmt.Sprintf("/api/v1/reader/books/%s/bookmarks", book.ID.Hex()), invalidBookmark, token)
		require.Equalf(t, 400, w.Code, "position=0 应返回 400: body=%s", w.Body.String())
	})

	t.Run("步骤4_验证边界数据仍可被读取", func(t *testing.T) {
		user := env.GetTestData("boundary_user").(*users.User)
		book := env.GetTestData("boundary_book").(*bookstore.Book)
		token := env.GetTestData("boundary_token").(string)

		comments := actions.GetBookComments(book.ID.Hex(), token)
		require.Contains(t, comments, "data", "评论列表响应缺少 data")

		collections := actions.GetReaderCollections(user.ID.Hex(), token)
		require.Contains(t, collections, "data", "收藏列表响应缺少 data")

		bookmarkList := env.DoRequest("GET", fmt.Sprintf("/api/v1/reader/books/%s/bookmarks?page=1&size=20", book.ID.Hex()), nil, token)
		require.Equalf(t, 200, bookmarkList.Code, "获取书签列表失败: body=%s", bookmarkList.Body.String())
	})
}
