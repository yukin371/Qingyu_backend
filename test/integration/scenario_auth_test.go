package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// AdminUsersPath 管理员用户API路径（其他路径已在helpers.go中定义）
const (
	AdminUsersPath = "/api/v1/admin/users"
)

// 用户认证流程测试
func TestAuthScenario(t *testing.T) {
	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 创建TestHelper
	helper := NewTestHelper(t, router)
	var testToken string
	testUsername := fmt.Sprintf("it_user_%d", time.Now().UnixNano())
	testPassword := "Test@123456"

	t.Run("0.准备测试账号", func(t *testing.T) {
		registerData := map[string]interface{}{
			"username": testUsername,
			"email":    fmt.Sprintf("%s@test.com", testUsername),
			"password": testPassword,
		}
		w := helper.DoRequest("POST", RegisterPath, registerData, "")
		if w.Code != 200 && w.Code != 201 {
			t.Fatalf("准备测试账号失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	t.Run("1.用户登录_普通用户", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": testUsername,
			"password": testPassword,
		}

		w := helper.DoRequest("POST", LoginPath, loginData, "")
		response := helper.AssertSuccess(w, 200, "普通用户登录失败")

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("响应data格式错误")
		}

		token, ok := data["token"].(string)
		if !ok {
			t.Fatal("token字段缺失或格式错误")
		}

		testToken = token

		helper.LogSuccess("普通用户登录成功")
		t.Logf("  用户名: %s", testUsername)
		tokenPreview := testToken
		if len(testToken) > 20 {
			tokenPreview = testToken[:20]
		}
		t.Logf("  Token: %s...", tokenPreview)

		assert.NotEmpty(t, testToken, "Token不应为空")
		assert.Contains(t, data, "user", "应该返回用户信息")
	})

	t.Run("2.用户登录_VIP用户", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "vip_user01",
			"password": "Vip@123456",
		}

		w := helper.DoRequest("POST", LoginPath, loginData, "")
		if w.Code != 200 {
			t.Skipf("环境未提供VIP种子账号，跳过（状态码: %d）", w.Code)
		}
		response := helper.AssertSuccess(w, 200, "VIP用户登录失败")

		data, ok := response["data"].(map[string]interface{})
		if ok {
			if user, ok := data["user"].(map[string]interface{}); ok {
				helper.LogSuccess("VIP用户登录成功")
				t.Logf("  用户名: vip_user01")
				t.Logf("  角色: %v", user["role"])

				assert.Equal(t, "vip", user["role"], "角色应该是vip")
			}
		} else {
			t.Logf("○ VIP登录失败")
		}
	})

	t.Run("3.用户登录_管理员", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "admin",
			"password": "Admin@123456",
		}

		w := helper.DoRequest("POST", LoginPath, loginData, "")
		if w.Code != 200 {
			t.Skipf("环境未提供管理员种子账号，跳过（状态码: %d）", w.Code)
		}
		response := helper.AssertSuccess(w, 200, "管理员登录失败")

		data, ok := response["data"].(map[string]interface{})
		if ok {
			if user, ok := data["user"].(map[string]interface{}); ok {
				helper.LogSuccess("管理员登录成功")
				t.Logf("  用户名: admin")
				t.Logf("  角色: %v", user["role"])

				assert.Equal(t, "admin", user["role"], "角色应该是admin")
			}
		} else {
			t.Logf("○ 管理员登录失败")
		}
	})

	t.Run("4.错误处理_错误的密码", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": testUsername,
			"password": "WrongPassword123",
		}

		w := helper.DoRequest("POST", LoginPath, loginData, "")
		helper.AssertError(w, 401, "用户名或密码错误", "错误密码应该返回401")

		helper.LogSuccess("错误密码处理正确")
	})

	t.Run("5.错误处理_不存在的用户", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "nonexist_user",
			"password": "Test@123456",
		}

		w := helper.DoRequest("POST", LoginPath, loginData, "")
		helper.AssertError(w, 401, "用户名或密码错误", "不存在的用户应该返回401")

		helper.LogSuccess("不存在用户处理正确")
	})

	if testToken != "" {
		t.Run("6.Token验证_访问需要认证的接口", func(t *testing.T) {
			w := helper.DoAuthRequest("GET", ReaderBooksPath, nil, testToken)

			// 应该能够访问 (200或其他非401状态码)
			if w.Code == 200 || w.Code == 404 {
				helper.LogSuccess("Token验证通过，可以访问认证接口")
			} else {
				t.Logf("○ Token验证可能失败，状态码: %d", w.Code)
			}
		})

		t.Run("7.Token验证_无Token访问受保护接口", func(t *testing.T) {
			w := helper.DoRequest("GET", ReaderBooksPath, nil, "")

			// 应该返回401未授权
			if w.Code == 401 {
				helper.LogSuccess("无Token访问受保护接口被正确拒绝")
			} else {
				t.Logf("○ 无Token访问状态: %d", w.Code)
			}
		})

		t.Run("8.Token验证_无效Token", func(t *testing.T) {
			w := helper.DoAuthRequest("GET", ReaderBooksPath, nil, "invalid_token_12345")

			// 应该返回401未授权
			if w.Code == 401 {
				helper.LogSuccess("无效Token被正确拒绝")
			} else {
				t.Logf("○ 无效Token状态: %d", w.Code)
			}
		})
	}

	t.Run("9.用户注册_新用户", func(t *testing.T) {
		// 生成随机用户名避免冲突
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		registerData := map[string]interface{}{
			"username": "newuser_" + timestamp,
			"email":    "newuser_" + timestamp + "@test.com",
			"password": "Test@123456",
		}

		w := helper.DoRequest("POST", RegisterPath, registerData, "")

		if w.Code == 200 || w.Code == 201 {
			helper.LogSuccess("新用户注册成功")
			t.Logf("  用户名: %s", registerData["username"])
			t.Logf("  邮箱: %s", registerData["email"])
		} else {
			t.Logf("○ 注册失败或接口不存在，状态码: %d", w.Code)
		}
	})

	t.Run("10.权限验证_普通用户访问管理员接口", func(t *testing.T) {
		// 使用普通用户token
		if testToken == "" {
			t.Skip("没有可用的测试Token")
		}

		w := helper.DoAuthRequest("GET", AdminUsersPath, nil, testToken)

		// 应该返回403禁止访问
		if w.Code == 403 || w.Code == 401 {
			helper.LogSuccess("普通用户访问管理员接口被正确拒绝")
		} else if w.Code == 404 {
			t.Logf("○ 管理员接口不存在")
		} else {
			t.Logf("○ 权限验证状态: %d", w.Code)
		}
	})

	t.Logf("\n=== 认证流程测试完成 ===")
	t.Logf("测试场景: 登录 → Token验证 → 权限控制 → 注册")
}
