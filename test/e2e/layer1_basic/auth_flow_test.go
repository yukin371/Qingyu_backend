//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"

	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestAuthFlow 测试认证流程
// 流程: 注册 -> 登录 -> 获取用户信息 -> 登出
func TestAuthFlow(t *testing.T) {
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

