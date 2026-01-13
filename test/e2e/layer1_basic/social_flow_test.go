package layer1_basic

import (
	"testing"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestSocialFlow 测试社交流程
// 流程: 发表评论 -> 收藏书籍 -> 点赞 -> 查看互动记录
func TestSocialFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()
	assertions := env.Assert()

	// 步骤1: 创建用户并登录
	t.Run("步骤1_创建用户并登录", func(t *testing.T) {
		t.Log("创建测试用户并登录...")

		// 创建测试用户
		user := fixtures.CreateUser()
		t.Logf("✓ 用户创建成功: %s", user.Username)

		// 登录获取token
		token := actions.Login(user.Username, "Test1234")
		t.Logf("✓ 登录成功，获取token")

		// 保存到环境
		env.SetTestData("test_user", user)
		env.SetTestData("auth_token", token)
	})

	// 步骤2: 创建测试书籍和章节
	t.Run("步骤2_创建测试书籍和章节", func(t *testing.T) {
		t.Log("创建测试书籍和章节...")

		// 创建作者
		author := fixtures.CreateAdminUser()
		t.Logf("✓ 作者创建成功: %s", author.Username)

		// 创建书籍
		book := fixtures.CreateBook(author.ID)
		t.Logf("✓ 书籍创建成功: %s (ID: %s)", book.Title, book.ID.Hex())

		// 创建章节
		chapter := fixtures.CreateChapter(book.ID.Hex())
		t.Logf("✓ 章节创建成功: %s (ID: %s)", chapter.Title, chapter.ID.Hex())

		// 保存到环境
		env.SetTestData("test_book", book)
		env.SetTestData("test_chapter", chapter)
	})

	// 步骤3: 发表评论
	t.Run("步骤3_发表评论", func(t *testing.T) {
		t.Log("发表书籍评论...")

		user := env.GetTestData("test_user").(*users.User)
		book := env.GetTestData("test_book").(*bookstore.Book)
		token := env.GetTestData("auth_token").(string)

		// 发表评论
		comment := actions.AddComment(token, book.ID.Hex(), "", "这是一条E2E测试评论，测试评论功能是否正常工作。")

		// 验证响应
		assertions.AssertResponseContains(comment, "data")
		t.Logf("✓ 评论发表成功")

		// 验证评论存在
		assertions.AssertCommentExists(user.ID, book.ID.Hex())
	})

	// 步骤4: 收藏书籍
	t.Run("步骤4_收藏书籍", func(t *testing.T) {
		t.Log("收藏书籍...")

		user := env.GetTestData("test_user").(*users.User)
		book := env.GetTestData("test_book").(*bookstore.Book)
		token := env.GetTestData("auth_token").(string)

		// 收藏书籍
		collection := actions.CollectBook(token, book.ID.Hex())

		// 验证响应
		assertions.AssertResponseContains(collection, "data")
		t.Logf("✓ 书籍收藏成功: %s", book.Title)

		// 验证收藏记录存在
		assertions.AssertCollectionExists(user.ID, book.ID.Hex())
	})

	// 步骤5: 点赞书籍
	t.Run("步骤5_点赞书籍", func(t *testing.T) {
		t.Log("点赞书籍...")

		book := env.GetTestData("test_book").(*bookstore.Book)
		token := env.GetTestData("auth_token").(string)

		// 点赞书籍
		likeResp := actions.LikeChapter(token, book.ID.Hex())

		// 验证响应
		assertions.AssertResponseContains(likeResp, "data")
		t.Logf("✓ 书籍点赞成功: %s", book.Title)
	})

	// 步骤6: 查看收藏列表
	t.Run("步骤6_查看收藏列表", func(t *testing.T) {
		t.Log("查看用户收藏列表...")

		user := env.GetTestData("test_user").(*users.User)
		token := env.GetTestData("auth_token").(string)

		// 获取收藏列表
		collections := actions.GetReaderCollections(user.ID, token)

		// 验证响应
		assertions.AssertResponseContains(collections, "data")
		t.Logf("✓ 收藏列表获取成功")
	})

	// 步骤7: 查看书籍评论
	t.Run("步骤7_查看书籍评论", func(t *testing.T) {
		t.Log("查看书籍评论列表...")

		book := env.GetTestData("test_book").(*bookstore.Book)
		token := env.GetTestData("auth_token").(string)

		// 获取书籍评论
		comments := actions.GetBookComments(book.ID.Hex(), token)

		// 验证响应
		assertions.AssertResponseContains(comments, "data")
		t.Logf("✓ 书籍评论列表获取成功")
	})
}
