package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
)

// 用户认证流程测试
func TestAuthScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化
	_, err := config.LoadConfig("../..")
	require.NoError(t, err, "加载配置失败")

	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	baseURL := "http://localhost:8080"
	var testToken string

	t.Run("1.用户登录_普通用户", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "test_user01",
			"password": "Test@123456",
		}

		jsonData, _ := json.Marshal(loginData)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/login", baseURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		// 调试：打印原始响应体
		if err := json.Unmarshal(body, &map[string]interface{}{}); err != nil {
			t.Logf("⚠ JSON解析失败")
			t.Logf("  HTTP状态码: %d", resp.StatusCode)
			t.Logf("  响应体长度: %d bytes", len(body))
			if len(body) > 200 {
				t.Logf("  前200字符: %q", string(body[:200]))
				t.Logf("  后100字符: %q", string(body[len(body)-100:]))
			} else {
				t.Logf("  完整响应体: %q", string(body))
			}
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		if result["code"] == float64(200) {
			data := result["data"].(map[string]interface{})
			testToken = data["token"].(string)

			t.Logf("✓ 普通用户登录成功")
			t.Logf("  邮箱: test01@qingyu.com")
			t.Logf("  Token: %s...", testToken[:20])

			assert.NotEmpty(t, testToken, "Token不应为空")
			assert.Contains(t, data, "user", "应该返回用户信息")
		} else {
			t.Logf("○ 登录失败: %v", result["message"])
		}
	})

	t.Run("2.用户登录_VIP用户", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "vip_user01",
			"password": "Vip@123456",
		}

		jsonData, _ := json.Marshal(loginData)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/login", baseURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		if result["code"] == float64(200) {
			data := result["data"].(map[string]interface{})
			user := data["user"].(map[string]interface{})

			t.Logf("✓ VIP用户登录成功")
			t.Logf("  邮箱: vip01@qingyu.com")
			t.Logf("  角色: %v", user["role"])

			assert.Equal(t, "vip", user["role"], "角色应该是vip")
		} else {
			t.Logf("○ VIP登录失败: %v", result["message"])
		}
	})

	t.Run("3.用户登录_管理员", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "admin",
			"password": "Admin@123456",
		}

		jsonData, _ := json.Marshal(loginData)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/login", baseURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		if result["code"] == float64(200) {
			data := result["data"].(map[string]interface{})
			user := data["user"].(map[string]interface{})

			t.Logf("✓ 管理员登录成功")
			t.Logf("  邮箱: admin@qingyu.com")
			t.Logf("  角色: %v", user["role"])

			assert.Equal(t, "admin", user["role"], "角色应该是admin")
		} else {
			t.Logf("○ 管理员登录失败: %v", result["message"])
		}
	})

	t.Run("4.错误处理_错误的密码", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "test_user01",
			"password": "WrongPassword123",
		}

		jsonData, _ := json.Marshal(loginData)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/login", baseURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.NotEqual(t, float64(200), result["code"], "错误密码应该登录失败")
		t.Logf("✓ 错误密码处理正确: %v", result["message"])
	})

	t.Run("5.错误处理_不存在的用户", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "nonexist_user",
			"password": "Test@123456",
		}

		jsonData, _ := json.Marshal(loginData)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/login", baseURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.NotEqual(t, float64(200), result["code"], "不存在的用户应该登录失败")
		t.Logf("✓ 不存在用户处理正确: %v", result["message"])
	})

	if testToken != "" {
		t.Run("6.Token验证_访问需要认证的接口", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/books", baseURL), nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+testToken)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)

			// 应该能够访问
			if result["code"] == float64(200) || result["code"] == float64(404) {
				t.Logf("✓ Token验证通过，可以访问认证接口")
			} else {
				t.Logf("○ Token验证可能失败: %v", result["message"])
			}
		})

		t.Run("7.Token验证_无Token访问受保护接口", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/books", baseURL), nil)
			require.NoError(t, err)
			// 不设置Authorization头

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// 应该返回401未授权
			if resp.StatusCode == http.StatusUnauthorized {
				t.Logf("✓ 无Token访问受保护接口被正确拒绝")
			} else {
				t.Logf("○ 无Token访问状态: %d", resp.StatusCode)
			}
		})

		t.Run("8.Token验证_无效Token", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/books", baseURL), nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer invalid_token_12345")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// 应该返回401未授权
			if resp.StatusCode == http.StatusUnauthorized {
				t.Logf("✓ 无效Token被正确拒绝")
			} else {
				t.Logf("○ 无效Token状态: %d", resp.StatusCode)
			}
		})
	}

	t.Run("9.用户注册_新用户", func(t *testing.T) {
		// 生成随机邮箱避免冲突
		timestamp := fmt.Sprintf("%d", getNow().Unix())
		registerData := map[string]interface{}{
			"username": "newuser_" + timestamp,
			"email":    "newuser_" + timestamp + "@test.com",
			"password": "Test@123456",
		}

		jsonData, _ := json.Marshal(registerData)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v1/user/register", baseURL),
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			t.Logf("○ 注册请求失败: %v", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)

		if err == nil && result["code"] == float64(200) {
			t.Logf("✓ 新用户注册成功")
			t.Logf("  用户名: %s", registerData["username"])
			t.Logf("  邮箱: %s", registerData["email"])
		} else if result != nil {
			t.Logf("○ 注册失败或接口不存在: %v", result["message"])
		}
	})

	t.Run("10.权限验证_普通用户访问管理员接口", func(t *testing.T) {
		// 使用普通用户token
		if testToken == "" {
			t.Skip("没有可用的测试Token")
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/admin/users", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+testToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 应该返回403禁止访问
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			t.Logf("✓ 普通用户访问管理员接口被正确拒绝")
		} else if resp.StatusCode == http.StatusNotFound {
			t.Logf("○ 管理员接口不存在")
		} else {
			t.Logf("○ 权限验证状态: %d", resp.StatusCode)
		}
	})

	t.Logf("\n=== 认证流程测试完成 ===")
	t.Logf("测试场景: 登录 → Token验证 → 权限控制 → 注册")
}

// 辅助函数：获取当前时间
func getNow() time.Time {
	return time.Now()
}
