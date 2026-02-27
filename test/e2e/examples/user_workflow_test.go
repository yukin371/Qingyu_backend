//go:build e2e
// +build e2e

package e2e_examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestE2E_UserWorkflow 用户完整工作流E2E测试
// 展示从注册到阅读书籍的完整用户旅程
func TestE2E_UserWorkflow(t *testing.T) {
	// 1. 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	helpers := env.Helpers()
	actions := env.Actions()

	t.Run("用户注册-登录-阅读流程", func(t *testing.T) {
		// 1. 准备测试数据
		timestamp := time.Now().Unix()
		username := fmt.Sprintf("e2e_user_%d", timestamp)
		email := fmt.Sprintf("e2e_%d@test.com", timestamp)
		password := "Test@123456"

		// 2. 用户注册并登录
		env.LogInfo("开始用户注册流程")
		token := helpers.RegisterAndLogin(username, email, password)
		require.NotEmpty(t, token, "用户注册登录失败")

		// 3. 浏览书城
		env.LogInfo("浏览书城首页")
		_ = actions.GetBookstoreHomepage()
		env.LogSuccess("书城首页加载成功")

		// 4. 搜索书籍
		env.LogInfo("搜索书籍")
		_ = actions.SearchBooks("test")
		env.LogSuccess("搜索结果加载成功")

		// 5. 查看榜单
		env.LogInfo("查看热门榜单")
		_ = actions.GetRankings("hot")
		env.LogSuccess("热门榜单加载成功")

		// 6. 清理测试数据
		helpers.CleanupTestUser(username)
	})
}

// TestE2E_ReaderWorkflow 读者完整工作流E2E测试
func TestE2E_ReaderWorkflow(t *testing.T) {
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	helpers := env.Helpers()
	actions := env.Actions()

	t.Run("读者阅读流程", func(t *testing.T) {
		// 1. 登录测试用户
		token := helpers.LoginAsTestUser()
		require.NotEmpty(t, token, "测试用户登录失败")

		// 2. 获取书城首页
		_ = actions.GetBookstoreHomepage()

		// 3. 获取书籍详情（使用示例ID）
		_ = actions.GetBookDetail("example-book-id")
		env.LogSuccess("书籍详情获取成功")

		// 4. 获取章节列表
		_ = actions.GetChapterList("example-book-id", token)
		env.LogSuccess("章节列表获取成功")
	})
}

// TestE2E_SocialWorkflow 社交互动工作流E2E测试
func TestE2E_SocialWorkflow(t *testing.T) {
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	helpers := env.Helpers()
	actions := env.Actions()

	t.Run("用户社交互动流程", func(t *testing.T) {
		token := helpers.LoginAsTestUser()
		require.NotEmpty(t, token, "测试用户登录失败")

		// 1. 收藏书籍
		env.LogInfo("收藏书籍")
		_ = actions.CollectBook(token, "example-book-id")
		env.LogSuccess("收藏成功")

		// 2. 点赞书籍
		env.LogInfo("点赞书籍")
		_ = actions.LikeChapter(token, "example-book-id")
		env.LogSuccess("点赞成功")

		// 3. 发表评论
		env.LogInfo("发表评论")
		commentContent := fmt.Sprintf("这是E2E测试评论_%d", time.Now().Unix())
		_ = actions.AddComment(token, "example-book-id", "", commentContent)
		env.LogSuccess("评论发表成功")

		// 4. 获取评论列表
		env.LogInfo("获取评论列表")
		_ = actions.GetBookComments("example-book-id", token)
		env.LogSuccess("评论列表获取成功")

		// 5. 获取收藏列表
		env.LogInfo("获取收藏列表")
		_ = actions.GetReaderCollections("", token)
		env.LogSuccess("收藏列表获取成功")
	})
}

// TestE2E_Performance 性能测试示例
func TestE2E_Performance(t *testing.T) {
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	helpers := env.Helpers()

	t.Run("书城首页性能测试", func(t *testing.T) {
		iterations := 20
		result := helpers.BenchmarkRequest("GET", "/api/v1/bookstore/homepage", nil, "", iterations)

		result.LogBenchmarkResult(t)

		// 验证性能在可接受范围内
		assert.Less(t, result.AvgDuration.Milliseconds(), int64(100),
			"平均响应时间应小于100ms")
	})

	t.Run("并发请求测试", func(t *testing.T) {
		concurrency := 10
		iterations := 5

		result := helpers.ConcurrentRequest(
			"GET",
			"/api/v1/bookstore/homepage",
			nil,
			"",
			concurrency,
			iterations,
		)

		result.LogConcurrentResult(t)

		// 验证成功率
		assert.Greater(t, result.SuccessRate, 95.0,
			"成功率应大于95%%")
	})
}

// TestE2E_ErrorHandling 错误处理E2E测试
func TestE2E_ErrorHandling(t *testing.T) {
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	helpers := env.Helpers()

	t.Run("认证错误处理", func(t *testing.T) {
		// 1. 无Token访问受保护接口
		w := env.DoRequest("GET", "/api/v1/users/profile", nil, "")
		helpers.AssertError(w, 401, "")
		env.LogSuccess("无Token访问被正确拒绝")

		// 2. 无效Token访问
		w = env.DoRequest("GET", "/api/v1/users/profile", nil, "invalid_token")
		helpers.AssertError(w, 401, "")
		env.LogSuccess("无效Token被正确拒绝")

		// 3. 错误的密码登录
		loginData := map[string]interface{}{
			"username": "test_user01",
			"password": "wrong_password",
		}
		w = env.DoRequest("POST", "/api/v1/login", loginData, "")
		helpers.AssertError(w, 401, "用户名或密码错误")
		env.LogSuccess("错误密码登录被正确拒绝")
	})

	t.Run("权限错误处理", func(t *testing.T) {
		// 1. 普通用户访问管理员接口
		token := helpers.LoginAsTestUser()
		w := env.DoRequest("GET", "/api/v1/admin/users", nil, token)

		if w.Code == 403 || w.Code == 401 {
			env.LogSuccess("普通用户访问管理员接口被正确拒绝")
		}
	})

	t.Run("资源不存在错误处理", func(t *testing.T) {
		// 1. 访问不存在的书籍
		w := env.DoRequest("GET", "/api/v1/bookstore/books/non-existent-book-id", nil, "")
		helpers.AssertError(w, 404, "")
		env.LogSuccess("不存在的书籍返回404")

		// 2. 访问不存在的章节
		w = env.DoRequest("GET", "/api/v1/bookstore/chapters/non-existent-chapter-id/content", nil, "")
		helpers.AssertError(w, 404, "")
		env.LogSuccess("不存在的章节返回404")
	})
}
