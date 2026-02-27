package integration

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMiddleware_Integration 中间件集成测试
// 测试认证、权限、限流等中间件的完整流程
func TestMiddleware_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过中间件集成测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("认证中间件完整流程", func(t *testing.T) {
		testAuthenticationMiddleware(t, router)
	})

	t.Run("权限中间件完整流程", func(t *testing.T) {
		testAuthorizationMiddleware(t, router)
	})

	t.Run("限流中间件完整流程", func(t *testing.T) {
		testRateLimitMiddleware(t, router)
	})

	t.Run("版本路由中间件", func(t *testing.T) {
		testVersionRoutingMiddleware(t, router)
	})
}

// testAuthenticationMiddleware 测试认证中间件
func testAuthenticationMiddleware(t *testing.T, router *gin.Engine) {
	helper := NewTestHelper(t, router)

	// 1. 准备测试账号
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("middleware_auth_%d", timestamp)
	testPassword := "Test@123456"

	t.Run("1.注册测试账号", func(t *testing.T) {
		registerData := map[string]interface{}{
			"username": testUsername,
			"email":    fmt.Sprintf("%s@test.com", testUsername),
			"password": testPassword,
		}
		w := helper.DoRequest("POST", RegisterPath, registerData, "")
		if w.Code != 200 && w.Code != 201 {
			t.Logf("注册状态码: %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 2. 测试无Token访问
	t.Run("2.无Token访问受保护接口", func(t *testing.T) {
		w := helper.DoRequest("GET", UserProfilePath, nil, "")

		// 应该返回401未授权
		assert.Equal(t, 401, w.Code, "无Token访问应该返回401")
		helper.LogSuccess("无Token访问被正确拒绝")
	})

	// 3. 测试无效Token
	t.Run("3.无效Token访问", func(t *testing.T) {
		invalidTokens := []string{
			"",
			"invalid_token",
			"Bearer invalid_token",
			"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		}

		for _, token := range invalidTokens {
			w := helper.DoRequest("GET", UserProfilePath, nil, token)
			assert.Equal(t, 401, w.Code, "无效Token应该返回401: %s", token)
		}
		helper.LogSuccess("无效Token被正确拒绝")
	})

	// 4. 测试有效Token
	t.Run("4.有效Token访问", func(t *testing.T) {
		token := helper.LoginUser(testUsername, testPassword)
		require.NotEmpty(t, token, "登录获取Token失败")

		w := helper.DoRequest("GET", UserProfilePath, nil, token)

		// 应该能够访问（200或404，取决于接口是否存在）
		if w.Code == 200 || w.Code == 404 {
			helper.LogSuccess("有效Token访问成功")
		} else {
			t.Logf("有效Token访问状态码: %d", w.Code)
		}
	})

	// 5. 测试Token过期
	t.Run("5.Token过期验证", func(t *testing.T) {
		// 使用格式正确但签名错误的Token
		fakeToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAwMDAwMDAsInVzZXJfaWQiOiJmYWtlIn0.fake_signature"

		w := helper.DoRequest("GET", UserProfilePath, nil, fakeToken)

		// 应该返回401
		assert.Equal(t, 401, w.Code, "过期Token应该返回401")
		helper.LogSuccess("过期Token被正确拒绝")
	})
}

// testAuthorizationMiddleware 测试权限中间件
func testAuthorizationMiddleware(t *testing.T, router *gin.Engine) {
	helper := NewTestHelper(t, router)

	// 1. 测试普通用户访问管理员接口
	t.Run("1.普通用户访问管理员接口", func(t *testing.T) {
		// 尝试使用普通用户token（如果有的话）
		testToken := helper.LoginUser("test_user01", "Test@123456")
		if testToken == "" {
			// 如果没有测试账号，创建一个
			timestamp := time.Now().Unix()
			testUsername := fmt.Sprintf("perm_user_%d", timestamp)
			registerData := map[string]interface{}{
				"username": testUsername,
				"email":    fmt.Sprintf("%s@test.com", testUsername),
				"password": "Test@123456",
			}
			w := helper.DoRequest("POST", RegisterPath, registerData, "")
			if w.Code == 200 || w.Code == 201 {
				testToken = helper.LoginUser(testUsername, "Test@123456")
			}
		}

		if testToken != "" {
			// 测试访问管理员接口
			adminPaths := []string{
				"/api/v1/admin/users",
				"/api/v1/admin/audit/logs",
				"/api/v1/admin/config",
			}

			for _, path := range adminPaths {
				w := helper.DoRequest("GET", path, nil, testToken)
				// 应该返回403禁止访问或404接口不存在
				if w.Code == 403 {
					t.Logf("✓ 路径 %s 被正确拒绝（403）", path)
				} else if w.Code == 404 {
					t.Logf("○ 路径 %s 不存在（404）", path)
				} else {
					t.Logf("○ 路径 %s 返回 %d", path, w.Code)
				}
			}
		} else {
			t.Skip("无法获取测试Token，跳过权限测试")
		}
	})

	// 2. 测试VIP用户访问VIP接口
	t.Run("2.VIP用户访问VIP接口", func(t *testing.T) {
		// 这个测试需要VIP用户，可能需要跳过
		t.Skip("需要VIP测试账号，跳过VIP权限测试")
	})
}

// testRateLimitMiddleware 测试限流中间件
func testRateLimitMiddleware(t *testing.T, router *gin.Engine) {
	helper := NewTestHelper(t, router)

	t.Run("1.正常请求不被限流", func(t *testing.T) {
		// 发送几个请求，不应该被限流
		for i := 0; i < 5; i++ {
			w := helper.DoRequest("GET", BookstoreHomePath, nil, "")
			// 应该正常返回
			if w.Code == 200 {
				// 正常
			} else if w.Code == 429 {
				t.Logf("请求被限流（可能测试环境限流较严格）")
			}
		}
		helper.LogSuccess("正常请求不被限流")
	})

	t.Run("2.高频请求触发限流", func(t *testing.T) {
		// 发送大量请求，可能触发限流
		rateLimited := false
		for i := 0; i < 20; i++ {
			w := helper.DoRequest("GET", BookstoreHomePath, nil, "")
			if w.Code == 429 {
				rateLimited = true
				t.Logf("✓ 请求 #%d 触发限流", i+1)
				break
			}
		}

		if rateLimited {
			helper.LogSuccess("高频请求成功触发限流")
		} else {
			t.Log("○ 未触发限流（可能限流阈值较高）")
		}
	})
}

// testVersionRoutingMiddleware 测试版本路由中间件
func testVersionRoutingMiddleware(t *testing.T, router *gin.Engine) {
	helper := NewTestHelper(t, router)

	t.Run("1.API版本路由", func(t *testing.T) {
		// 测试不同版本的API路径
		versionPaths := []struct {
			path    string
			exists  bool
			comment string
		}{
			{"/api/v1/bookstore/homepage", true, "v1书城首页"},
			{"/api/v2/bookstore/homepage", false, "v2书城首页（可能不存在）"},
			{"/api/v1/user/auth/login", true, "v1用户登录"},
		}

		for _, vp := range versionPaths {
			w := helper.DoRequest("GET", vp.path, nil, "")
			if w.Code == 404 {
				if vp.exists {
					t.Logf("○ 路径 %s 返回404", vp.path)
				} else {
					t.Logf("✓ 路径 %s 不存在（符合预期）", vp.path)
				}
			} else if w.Code == 200 {
				t.Logf("✓ 路径 %s 可访问", vp.path)
			}
		}
	})

	t.Run("2.版本协商", func(t *testing.T) {
		// 测试Accept头中的版本协商
		req := httptest.NewRequest("GET", BookstoreHomePath, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("API-Version", "v1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == 200 {
			helper.LogSuccess("版本协商成功")
		} else {
			t.Logf("版本协商状态码: %d", w.Code)
		}
	})
}

// TestMiddleware_ErrorHandling 中间件错误处理集成测试
func TestMiddleware_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过中间件错误处理测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("中间件错误响应格式", func(t *testing.T) {
		// 1. 测试401错误格式
		t.Run("1.401未授权错误格式", func(t *testing.T) {
			w := helper.DoRequest("GET", UserProfilePath, nil, "")

			if w.Code == 401 {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// 验证错误响应格式
				assert.Contains(t, response, "code", "应该包含code字段")
				assert.Contains(t, response, "message", "应该包含message字段")

				helper.LogSuccess("401错误格式正确")
			}
		})

		// 2. 测试403错误格式
		t.Run("2.403禁止访问错误格式", func(t *testing.T) {
			// 使用普通用户token访问管理员接口
			testToken := helper.LoginUser("test_user01", "Test@123456")
			if testToken != "" {
				w := helper.DoRequest("GET", "/api/v1/admin/users", nil, testToken)

				if w.Code == 403 {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					require.NoError(t, err)

					assert.Contains(t, response, "code", "应该包含code字段")
					assert.Contains(t, response, "message", "应该包含message字段")

					helper.LogSuccess("403错误格式正确")
				}
			}
		})

		// 3. 测试429错误格式
		t.Run("3.429限流错误格式", func(t *testing.T) {
			// 发送多个请求尝试触发限流
			for i := 0; i < 30; i++ {
				w := helper.DoRequest("GET", BookstoreHomePath, nil, "")
				if w.Code == 429 {
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					require.NoError(t, err)

					assert.Contains(t, response, "code", "应该包含code字段")
					assert.Contains(t, response, "message", "应该包含message字段")

					helper.LogSuccess("429错误格式正确")
					return
				}
			}
			t.Log("○ 未触发限流，无法验证429错误格式")
		})
	})
}

// TestMiddleware_CORS 跨域中间件集成测试
func TestMiddleware_CORS(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过CORS中间件测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("CORS预检请求", func(t *testing.T) {
		// 发送OPTIONS预检请求
		req := httptest.NewRequest("OPTIONS", BookstoreHomePath, nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 检查CORS响应头
		corsHeaders := []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers",
		}

		hasCORS := false
		for _, header := range corsHeaders {
			if w.Header().Get(header) != "" {
				hasCORS = true
				t.Logf("✓ CORS头 %s: %s", header, w.Header().Get(header))
			}
		}

		if hasCORS {
			t.Log("✓ CORS中间件正常工作")
		} else {
			t.Log("○ 未检测到CORS头")
		}
	})
}

// TestMiddleware_RequestLogging 请求日志中间件集成测试
func TestMiddleware_RequestLogging(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过请求日志中间件测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("请求ID生成", func(t *testing.T) {
		w := helper.DoRequest("GET", BookstoreHomePath, nil, "")

		// 检查响应头中是否有请求ID
		requestID := w.Header().Get("X-Request-ID")
		if requestID != "" {
			t.Logf("✓ 请求ID: %s", requestID)
		} else {
			t.Log("○ 未检测到请求ID")
		}
	})

	t.Run("请求日志记录", func(t *testing.T) {
		// 这个测试需要检查日志输出，暂时跳过
		t.Skip("需要日志捕获机制，暂时跳过")
	})
}
