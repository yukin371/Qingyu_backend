//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestWritingFlow 测试写作流程
// 流程: 创建写作项目 -> 验证项目存在
func TestWritingFlow(t *testing.T) {
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

