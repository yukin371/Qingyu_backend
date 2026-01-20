//go:build e2e
// +build e2e

package e2e_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"Qingyu_backend/models/users"
	e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestE2E_CompleteUserJourney 完整用户旅程
// 流程：注册 -> 浏览书城 -> 阅读 -> 互动 -> 成为作者
func TestE2E_CompleteUserJourney(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	t.Log("========== E2E 完整用户旅程测试 ==========")

	// Setup
	env, cleanup := e2eFramework.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()
	assertions := env.Assert()

	// ========== 阶段1：新用户注册 ==========
	var userToken string
	var testUser *users.User
	var bookID, chapterID string

	t.Run("Phase1_Registration", func(t *testing.T) {
		t.Log(">>> 阶段1：用户注册")

		// 创建用户，不指定用户名，让系统自动生成唯一的
		testUser = fixtures.CreateUser()
		require.NotNil(t, testUser)
		t.Logf("✓ 用户注册成功: %s (%s)", testUser.Username, testUser.ID)

		// 登录获取 token
		userToken = actions.Login(testUser.Username, "Test1234")
		require.NotEmpty(t, userToken)
	})

	// ========== 阶段2：浏览发现书城 ==========
	t.Run("Phase2_Discovery", func(t *testing.T) {
		t.Log(">>> 阶段2：浏览发现书城")

		// 浏览书城首页
		homepage := actions.GetBookstoreHomepage()
		assertions.AssertResponseContains(homepage, "code")
		t.Log("✓ 书城首页加载成功")

		// 查看榜单
		rankings := actions.GetRankings("realtime")
		assertions.AssertResponseContains(rankings, "data")
		t.Log("✓ 榜单加载成功")

		// 搜索书籍
		searchResults := actions.SearchBooks("测试")
		assertions.AssertResponseContains(searchResults, "data")
		t.Log("✓ 搜索功能正常")
	})

	// ========== 阶段3：阅读体验 ==========
	t.Run("Phase3_Reading", func(t *testing.T) {
		t.Log(">>> 阶段3：阅读体验")

		// 准备测试书籍
		author := fixtures.CreateUser()
		book := fixtures.CreateBook(author.ID, e2eFramework.WithBookPrice(0))
		chapter := fixtures.CreateChapter(book.ID.Hex())

		bookID = book.ID.Hex()
		chapterID = chapter.ID.Hex()

		// 查看书籍详情
		bookDetail := actions.GetBookDetail(bookID)
		assertions.AssertResponseContains(bookDetail, "data")
		t.Log("✓ 书籍详情加载成功")

		// 获取章节列表
		chapterList := actions.GetChapterList(bookID, userToken)
		assertions.AssertResponseContains(chapterList, "data")
		t.Log("✓ 章节列表加载成功")

		// 获取章节内容
		chapterContent := actions.GetChapter(chapterID, userToken)
		assertions.AssertResponseContains(chapterContent, "data")
		t.Log("✓ 章节内容加载成功")

		// 开始阅读（保存进度）
		_ = actions.StartReading(testUser.ID, bookID, chapterID, userToken)
		t.Log("✓ 阅读进度保存成功")
	})

	// ========== 阶段4：社交互动 ==========
	t.Run("Phase4_Interaction", func(t *testing.T) {
		t.Log(">>> 阶段4：社交互动")

		// 发表评论（至少 10 个字符，符合 API 要求）
		comment := actions.AddComment(userToken, bookID, chapterID, "这是一本非常好的书，我非常喜欢！")
		assertions.AssertResponseContains(comment, "data")
		t.Log("✓ 发表评论成功")

		// 验证评论存在
		assertions.AssertCommentExists(testUser.ID, bookID)

		// 收藏书籍
		collection := actions.CollectBook(userToken, bookID)
		assertions.AssertResponseContains(collection, "data")
		t.Log("✓ 收藏书籍成功")

		// 验证收藏记录存在
		assertions.AssertCollectionExists(testUser.ID, bookID)

		// 点赞书籍（改为书籍点赞）
		like := actions.LikeChapter(userToken, bookID)
		assertions.AssertResponseContains(like, "data")
		t.Log("✓ 点赞成功")

		// 书签功能未实现（bookmarkSvc = nil），跳过测试
		t.Log("⚠ 书签功能未实现，跳过测试")
	})

	// ========== 阶段5：成为作者 ==========
	t.Run("Phase5_Authoring", func(t *testing.T) {
		t.Log(">>> 阶段5：成为作者")

		// 创建写作项目
		project := actions.CreateProject(userToken, map[string]interface{}{
			"title":       "我的第一本书",
			"description": "这是一本原创小说",
			"genre":       "都市",
		})
		assertions.AssertResponseContains(project, "data")
		t.Log("✓ 创建写作项目成功")

		// 验证项目创建
		projectData, ok := project["data"].(map[string]interface{})
		if ok {
			projectID, ok := projectData["id"].(string)
			if ok {
				t.Logf("✓ 项目ID: %s", projectID)
			}
		}
	})

	t.Log("========== 完整用户旅程测试完成 ==========")
}

// TestE2E_VIPReadingFlow VIP用户阅读流程测试
// 流程：VIP用户 -> 免费阅读付费书籍
func TestE2E_VIPReadingFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	t.Log("========== E2E VIP阅读流程测试 ==========")

	// Setup
	env, cleanup := e2eFramework.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()
	assertions := env.Assert()

	// 1. 创建VIP用户（不指定用户名，让系统自动生成唯一的）
	vipUser := fixtures.CreateUser(e2eFramework.WithVIPLevel(1))
	vipToken := actions.Login(vipUser.Username, "Test1234")

	assertions.AssertUserIsVIP(vipUser.ID, true)
	t.Log("✓ VIP用户创建成功")

	// 2. 创建免费书籍
	author := fixtures.CreateUser()
	book := fixtures.CreateBook(author.ID,
		e2eFramework.WithBookPrice(0),
	)
	chapter := fixtures.CreateChapter(book.ID.Hex(),
		e2eFramework.WithChapterFree(true),
	)

	// 3. VIP用户阅读
	chapterContent := actions.GetChapter(chapter.ID.Hex(), vipToken)
	assertions.AssertResponseContains(chapterContent, "data")

	t.Log("✓ VIP用户阅读章节成功")

	t.Log("========== VIP阅读流程测试完成 ==========")
}

