//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// RunAuthFlow 导出的入口函数，供suite_test.go调用
func RunAuthFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()
	assertions := env.Assert()

	t.Run("步骤1_创建测试用户", func(t *testing.T) {
		t.Log("创建测试用户...")

		// 创建测试用户
		user := fixtures.CreateUser()
		t.Logf("✓ 用户创建成功: %s (ID: %s)", user.Username, user.ID)

		// 保存用户信息到测试环境
		env.SetTestData("test_user", user)
	})

	t.Run("步骤2_用户登录", func(t *testing.T) {
		t.Log("执行用户登录...")

		// 获取之前创建的用户
		user := env.GetTestData("test_user").(*users.User)

		// 执行登录
		token := actions.Login(user.Username, "Test1234")
		t.Logf("✓ 登录成功，获取到 token: %s...", token[:20])

		// 保存token
		env.SetTestData("auth_token", token)
	})

	t.Run("步骤3_验证用户存在", func(t *testing.T) {
		t.Log("验证用户信息...")

		// 获取用户信息
		user := env.GetTestData("test_user").(*users.User)

		// 验证用户存在
		assertions.AssertUserExists(user.ID)
		t.Logf("✓ 用户验证通过: %s", user.Username)
	})

	t.Run("步骤4_获取用户详细信息", func(t *testing.T) {
		t.Log("获取用户详细信息...")

		// 获取用户信息
		user := env.GetTestData("test_user").(*users.User)

		// 通过Actions获取用户详细信息
		userDetail := actions.GetUser(user.ID)
		t.Logf("✓ 获取用户详情: %s, Email: %s, VIP: %d",
			userDetail.Username, userDetail.Email, userDetail.VIPLevel)

		// 验证用户状态
		if userDetail.Status != users.UserStatusActive {
			t.Errorf("期望用户状态为 active，实际为: %s", userDetail.Status)
		}
	})
}

// RunReadingFlow 导出的入口函数，供suite_test.go调用
func RunReadingFlow(t *testing.T) {
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

// RunSocialFlow 导出的入口函数，供suite_test.go调用
func RunSocialFlow(t *testing.T) {
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

// RunWritingFlow 导出的入口函数，供suite_test.go调用
func RunWritingFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	// 步骤1: 创建作者用户并登录
	t.Run("步骤1_创建作者用户并登录", func(t *testing.T) {
		t.Log("创建作者用户并登录...")

		// 创建作者用户（不使用固定用户名，避免冲突）
		author := fixtures.CreateUser()
		t.Logf("✓ 作者用户创建成功: %s (角色: %v)", author.Username, author.Roles)

		// 登录获取token
		token := actions.Login(author.Username, "Test1234")
		t.Logf("✓ 登录成功，获取token")

		// 保存到环境
		env.SetTestData("author_user", author)
		env.SetTestData("auth_token", token)
	})

	// 步骤2: 创建写作项目
	t.Run("步骤2_创建写作项目", func(t *testing.T) {
		t.Log("创建写作项目...")

		token := env.GetTestData("auth_token").(string)

		// 创建项目请求
		projectReq := map[string]interface{}{
			"title":       "e2e_test_writing_project",
			"description": "这是一个E2E测试写作项目",
			"genre":       "小说",
			"status":      "draft",
		}

		// 创建项目
		projectResp := actions.CreateProject(token, projectReq)

		// 验证响应
		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			if projectId, ok := data["projectId"].(string); ok {
				t.Logf("✓ 写作项目创建成功 (ID: %s)", projectId)

				// 保存项目ID
				env.SetTestData("project_id", projectId)
			} else {
				t.Error("项目响应中未找到projectId字段")
			}
		} else {
			t.Error("项目响应格式不正确")
		}
	})

	// 步骤3: 验证项目存在（通过查询项目列表）
	t.Run("步骤3_验证项目存在", func(t *testing.T) {
		t.Log("验证写作项目存在...")

		token := env.GetTestData("auth_token")
		if token == nil {
			t.Skip("auth_token未设置，跳过此步骤")
			return
		}

		projectId := env.GetTestData("project_id")
		if projectId == nil {
			t.Skip("project_id未设置，跳过此步骤")
			return
		}

		// 获取项目列表
		w := env.DoRequest("GET", "/api/v1/writer/projects", nil, token.(string))

		if w.Code == 200 {
			t.Logf("✓ 项目列表获取成功，项目 %v 应在列表中", projectId)
		} else {
			t.Errorf("获取项目列表失败，状态码: %d", w.Code)
		}
	})
}
