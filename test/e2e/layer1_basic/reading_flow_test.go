//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestReadingFlow 测试阅读流程
// 流程: 浏览书城 -> 查看书籍详情 -> 获取章节列表 -> 阅读章节内容 -> 保存阅读进度
func TestReadingFlow(t *testing.T) {
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

		// 获取作者（使用admin角色）
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

	// 步骤3: 浏览书城首页
	t.Run("步骤3_浏览书城首页", func(t *testing.T) {
		t.Log("获取书城首页数据...")

		homepage := actions.GetBookstoreHomepage()

		// 验证响应包含数据
		assertions.AssertResponseContains(homepage, "data")
		t.Logf("✓ 书城首页数据获取成功")
	})

	// 步骤4: 查看书籍详情
	t.Run("步骤4_查看书籍详情", func(t *testing.T) {
		t.Log("获取书籍详情...")

		book := env.GetTestData("test_book").(*bookstore.Book)

		// 获取书籍详情
		bookDetail := actions.GetBookDetail(book.ID.Hex())

		// 验证响应
		assertions.AssertResponseContains(bookDetail, "data")
		t.Logf("✓ 书籍详情获取成功: %s", book.Title)
	})

	// 步骤5: 获取章节列表
	t.Run("步骤5_获取章节列表", func(t *testing.T) {
		t.Log("获取书籍章节列表...")

		book := env.GetTestData("test_book").(*bookstore.Book)
		token := env.GetTestData("auth_token").(string)

		// 获取章节列表
		chapterList := actions.GetChapterList(book.ID.Hex(), token)

		// 验证响应
		assertions.AssertResponseContains(chapterList, "data")
		t.Logf("✓ 章节列表获取成功")
	})

	// 步骤6: 阅读章节内容
	t.Run("步骤6_阅读章节内容", func(t *testing.T) {
		t.Log("获取章节内容...")

		chapter := env.GetTestData("test_chapter").(*bookstore.Chapter)
		token := env.GetTestData("auth_token").(string)

		// 获取章节内容
		chapterContent := actions.GetChapter(chapter.ID.Hex(), token)

		// 验证响应
		assertions.AssertResponseContains(chapterContent, "data")
		t.Logf("✓ 章节内容获取成功: %s", chapter.Title)
	})

	// 步骤7: 保存阅读进度
	t.Run("步骤7_保存阅读进度", func(t *testing.T) {
		t.Log("保存阅读进度...")

		user := env.GetTestData("test_user").(*users.User)
		book := env.GetTestData("test_book").(*bookstore.Book)
		chapter := env.GetTestData("test_chapter").(*bookstore.Chapter)
		token := env.GetTestData("auth_token").(string)

		// 保存阅读进度
		progress := actions.StartReading(user.ID, book.ID.Hex(), chapter.ID.Hex(), token)

		// 验证响应包含code和message字段（data可能为null）
		assertions.AssertResponseContains(progress, "code")
		assertions.AssertResponseContains(progress, "message")
		t.Logf("✓ 阅读进度保存成功")

		// 验证阅读进度已保存
		assertions.AssertReadingProgress(user.ID, book.ID.Hex())
	})
}

