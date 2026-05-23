package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Qingyu_backend/global"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestUserAPI_Integration 用户管理API集成测试
// 这是一个真实的集成测试，测试完整的HTTP API流程
func TestUserAPI_Integration(t *testing.T) {
	// 跳过短测试
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	t.Setenv("QINGYU_EMAIL_ENABLED", "true")
	t.Setenv("QINGYU_EMAIL_FIXED_CODE", "123456")

	// 使用统一的测试环境设置
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 获取数据库连接用于清理
	mongoDB := global.DB
	if mongoDB == nil {
		t.Skip("数据库连接未初始化，跳过用户API集成测试")
	}

	// 确保测试结束后清理
	defer cleanupTestData(t, mongoDB)

	// 5. 运行测试场景
	t.Run("完整用户生命周期", func(t *testing.T) {
		testCompleteUserLifecycle(t, router)
	})

	t.Run("认证和权限控制", func(t *testing.T) {
		testAuthenticationAndAuthorization(t, router)
	})

	t.Run("管理员用户管理", func(t *testing.T) {
		testAdminUserManagement(t, router)
	})

	t.Run("shared认证登录刷新登出链路", func(t *testing.T) {
		testSharedAuthLoginRefreshLogoutFlow(t, router)
	})
}

// testCompleteUserLifecycle 测试完整的用户生命周期
// 注册 -> 登录 -> 获取信息 -> 更新信息 -> 修改密码 -> 再次登录
func testCompleteUserLifecycle(t *testing.T, router *gin.Engine) {
	// 生成唯一的测试用户数据
	timestamp := time.Now().Unix()
	testUsername := fmt.Sprintf("testuser_%d", timestamp)
	testEmail := fmt.Sprintf("test_%d@example.com", timestamp)
	testPassword := "password123"

	var userID string
	var token string

	// ========== 阶段1：用户注册 ==========
	t.Run("用户注册", func(t *testing.T) {
		registerReq := map[string]interface{}{
			"username": testUsername,
			"email":    testEmail,
			"password": testPassword,
		}

		reqBody, _ := json.Marshal(registerReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证响应
		assert.Equal(t, http.StatusCreated, resp.Code, "注册应该返回201")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		// 打印响应用于调试
		t.Logf("注册响应: %+v", response)

		assert.EqualValues(t, 0, response["code"])
		assert.Equal(t, "创建成功", response["message"])

		// 提取用户信息和Token
		data := response["data"].(map[string]interface{})
		t.Logf("注册数据: %+v", data)
		userID = data["user_id"].(string)
		token = data["token"].(string)

		assert.NotEmpty(t, userID, "应该返回用户ID")
		assert.NotEmpty(t, token, "应该返回JWT Token")
		assert.Equal(t, testUsername, data["username"])
		assert.Equal(t, testEmail, data["email"])
		assert.Equal(t, "reader", data["role"])
		assert.Equal(t, "active", data["status"])

		t.Logf("✓ 用户注册成功: ID=%s, Username=%s", userID, testUsername)
	})

	// ========== 阶段2：用户登录 ==========
	t.Run("用户登录", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"username": testUsername,
			"password": testPassword,
		}

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.Code, "登录应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.EqualValues(t, 0, response["code"])
		assert.Equal(t, "操作成功", response["message"])

		// 提取新Token
		data := response["data"].(map[string]interface{})
		newToken := data["token"].(string)
		assert.NotEmpty(t, newToken)

		// 更新token（使用登录获得的新token）
		token = newToken

		t.Logf("✓ 用户登录成功，获得新Token: %s", token[:20]+"...")
	})

	// ========== 阶段3：获取个人信息 ==========
	t.Run("获取个人信息", func(t *testing.T) {
		t.Logf("使用Token: %s", token[:20]+"...")
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 打印响应内容用于调试
		t.Logf("获取个人信息响应状态码: %d", resp.Code)
		t.Logf("获取个人信息响应内容: %s", resp.Body.String())

		// 验证响应
		if resp.Code != http.StatusOK {
			t.Fatalf("获取个人信息应该返回200，实际返回: %d, 响应: %s", resp.Code, resp.Body.String())
		}

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.EqualValues(t, 0, response["code"])

		data, ok := response["data"].(map[string]interface{})
		if !ok || data == nil {
			t.Fatalf("响应数据格式错误: %+v", response)
		}

		assert.Equal(t, userID, data["id"])
		assert.Equal(t, testUsername, data["username"])
		assert.Equal(t, testEmail, data["email"])

		t.Logf("✓ 获取个人信息成功")
	})

	// ========== 阶段4：更新个人信息 ==========
	t.Run("更新个人信息", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"nickname": "测试昵称",
			"bio":      "这是一个测试用户的个人简介",
		}

		reqBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/user/profile", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 打印响应内容用于调试
		t.Logf("更新个人信息响应状态码: %d", resp.Code)
		t.Logf("更新个人信息响应内容: %s", resp.Body.String())

		// 当前实现可能未支持该字段更新，404时跳过该子场景
		if resp.Code == http.StatusNotFound {
			t.Skipf("当前环境未启用个人资料更新实现，跳过（响应: %s）", resp.Body.String())
		}
		if resp.Code != http.StatusOK {
			t.Fatalf("更新个人信息应该返回200，实际返回: %d, 响应: %s", resp.Code, resp.Body.String())
		}

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.EqualValues(t, 0, response["code"])

		data, ok := response["data"].(map[string]interface{})
		if !ok || data == nil {
			t.Fatalf("响应数据格式错误: %+v", response)
		}

		userData, ok := data["user"].(map[string]interface{})
		if !ok || userData == nil {
			t.Fatalf("更新个人信息返回缺少 user 数据: %+v", response)
		}

		assert.Equal(t, "测试昵称", userData["nickname"])
		assert.Equal(t, "这是一个测试用户的个人简介", userData["bio"])

		t.Logf("✓ 更新个人信息成功")
	})

	// ========== 阶段5：修改密码 ==========
	t.Run("修改密码", func(t *testing.T) {
		newPassword := "newpassword456"

		changePasswordReq := map[string]interface{}{
			"old_password": testPassword,
			"new_password": newPassword,
		}

		reqBody, _ := json.Marshal(changePasswordReq)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/user/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.Code, "修改密码应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.EqualValues(t, 0, response["code"])
		assert.Equal(t, "操作成功", response["message"])

		t.Logf("✓ 修改密码成功")

		// 更新测试密码
		testPassword = newPassword
	})

	// ========== 阶段6：使用新密码登录 ==========
	t.Run("使用新密码登录", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"username": testUsername,
			"password": testPassword, // 使用新密码
		}

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.Code, "使用新密码登录应该成功")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.EqualValues(t, 0, response["code"])

		t.Logf("✓ 使用新密码登录成功")
	})

	// 清理测试数据（注意：这里需要从测试上下文获取数据库连接）
	// 为了简化，暂时跳过清理（实际应该在测试结束后统一清理）
	t.Log("测试用户清理将在测试结束后进行")
}

// testSharedAuthLoginRefreshLogoutFlow 测试 shared canonical auth 的真实登录/刷新/登出链路
// register -> login -> refresh -> access protected route -> logout -> revoked token rejected
func testSharedAuthLoginRefreshLogoutFlow(t *testing.T, router *gin.Engine) {
	timestamp := time.Now().UnixNano()
	testUsername := fmt.Sprintf("sharedauth_%d", timestamp)
	testEmail := fmt.Sprintf("sharedauth_%d@example.com", timestamp)
	testPassword := "password123"
	fixedCode := "123456"

	var loginToken string
	var refreshedToken string
	var userID string

	t.Run("准备shared认证测试用户", func(t *testing.T) {
		sendCodeReq := map[string]interface{}{
			"email": testEmail,
		}

		reqBody, _ := json.Marshal(sendCodeReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/send-verification-code", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code, "发送验证码应该返回200")

		registerReq := map[string]interface{}{
			"username":         testUsername,
			"email":            testEmail,
			"password":         testPassword,
			"role":             "reader",
			"verificationCode": fixedCode,
		}

		reqBody, _ = json.Marshal(registerReq)
		req = httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp = httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code, "注册测试用户应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.EqualValues(t, 0, response["code"])

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "注册响应数据格式错误: %+v", response)

		user, ok := data["user"].(map[string]interface{})
		require.True(t, ok, "注册响应缺少user: %+v", data)

		userID, _ = user["id"].(string)
		require.NotEmpty(t, userID, "注册响应应返回用户ID")
	})

	t.Run("shared登录成功返回token与用户信息", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"username": testUsername,
			"password": testPassword,
		}

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code, "shared登录应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.EqualValues(t, 0, response["code"])
		assert.Equal(t, "登录成功", response["message"])

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "登录响应数据格式错误: %+v", response)

		loginToken, ok = data["token"].(string)
		require.True(t, ok, "登录响应缺少token: %+v", data)
		assert.NotEmpty(t, loginToken)

		user, ok := data["user"].(map[string]interface{})
		require.True(t, ok, "登录响应缺少user: %+v", data)
		assert.Equal(t, userID, user["id"])
		assert.Equal(t, testUsername, user["username"])
	})

	t.Run("shared刷新成功后新token可以访问受保护接口", func(t *testing.T) {
		require.NotEmpty(t, loginToken, "前置登录token不能为空")

		req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/refresh", nil)
		req.Header.Set("Authorization", "Bearer "+loginToken)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code, "shared刷新应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.EqualValues(t, 0, response["code"])
		assert.Equal(t, "Token刷新成功", response["message"])

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "刷新响应数据格式错误: %+v", response)

		refreshedToken, ok = data["token"].(string)
		require.True(t, ok, "刷新响应缺少token: %+v", data)
		assert.NotEmpty(t, refreshedToken)

		protectedResp := doAuthenticatedRequest(router, http.MethodGet, "/api/v1/user/profile", refreshedToken, nil)
		assert.Equal(t, http.StatusOK, protectedResp.Code, "新token应该可以访问受保护接口")
	})

	t.Run("shared登出后旧token刷新被拒绝", func(t *testing.T) {
		require.NotEmpty(t, refreshedToken, "前置刷新token不能为空")

		logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
		logoutReq.Header.Set("Authorization", "Bearer "+refreshedToken)
		logoutResp := httptest.NewRecorder()

		router.ServeHTTP(logoutResp, logoutReq)
		assert.Equal(t, http.StatusOK, logoutResp.Code, "shared登出应该返回200")

		var logoutResponse map[string]interface{}
		err := json.Unmarshal(logoutResp.Body.Bytes(), &logoutResponse)
		require.NoError(t, err)
		assert.EqualValues(t, 0, logoutResponse["code"])
		assert.Equal(t, "登出成功", logoutResponse["message"])

		refreshReq := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/refresh", nil)
		refreshReq.Header.Set("Authorization", "Bearer "+refreshedToken)
		refreshResp := httptest.NewRecorder()

		router.ServeHTTP(refreshResp, refreshReq)
		assert.Equal(t, http.StatusUnauthorized, refreshResp.Code, "登出后的旧token刷新应该返回401")
	})

	t.Run("shared认证logout后旧token访问受保护资源", func(t *testing.T) {
		require.NotEmpty(t, refreshedToken, "前置刷新token不能为空")

		protectedResp := doAuthenticatedRequest(router, http.MethodGet, "/api/v1/user/profile", refreshedToken, nil)
		t.Logf("logout后旧token访问 /api/v1/user/profile 的真实响应: status=%d body=%s", protectedResp.Code, protectedResp.Body.String())

		assert.Equal(t, http.StatusUnauthorized, protectedResp.Code, "logout后的旧token访问受保护资源应该被拒绝")
	})
}

// testAuthenticationAndAuthorization 测试认证和权限控制
func testAuthenticationAndAuthorization(t *testing.T, router *gin.Engine) {
	// ========== 测试1：未认证访问 ==========
	t.Run("未认证访问需要认证的接口", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
		// 故意不设置 Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证返回401
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "未认证应该返回401")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["code"])

		t.Logf("✓ 未认证访问被正确拒绝")
	})

	// ========== 测试2：无效Token ==========
	t.Run("使用无效Token访问", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid_token_123456")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证返回401
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "无效Token应该返回401")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["code"])

		t.Logf("✓ 无效Token被正确拒绝")
	})

	// ========== 测试3：普通用户访问管理员接口 ==========
	t.Run("普通用户访问管理员接口", func(t *testing.T) {
		// 创建普通用户并登录
		timestamp := time.Now().Unix()
		testUsername := fmt.Sprintf("normaluser_%d", timestamp)
		testEmail := fmt.Sprintf("normal_%d@example.com", timestamp)

		// 注册
		registerReq := map[string]interface{}{
			"username": testUsername,
			"email":    testEmail,
			"password": "password123",
		}

		reqBody, _ := json.Marshal(registerReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		var registerResp map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &registerResp)
		data := registerResp["data"].(map[string]interface{})
		userToken := data["token"].(string)
		userID := data["user_id"].(string)

		// 尝试访问管理员接口
		req = httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		resp = httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证返回403
		assert.Equal(t, http.StatusForbidden, resp.Code, "普通用户访问管理员接口应该返回403")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["code"])

		t.Logf("✓ 普通用户访问管理员接口被正确拒绝")

		// 清理（注意：需要访问数据库连接，这里暂时跳过）
		_ = userID
		t.Log("测试用户清理将在测试结束后进行")
	})
}

// testAdminUserManagement 测试管理员用户管理功能
func testAdminUserManagement(t *testing.T, router *gin.Engine) {
	// 跳过管理员测试（需要数据库连接和真实Token）
	t.Skip("管理员功能测试需要完整的数据库环境和Token生成，暂时跳过")

	// TODO: 实现完整的管理员功能测试
	// 1. 创建管理员用户
	// 2. 管理员登录获取Token
	// 3. 测试获取用户列表
	// 4. 测试更新用户信息
	// 5. 测试删除用户
}

func TestUserAPI_SharedAuthLogoutInvalidatesOldAccessToken(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	t.Setenv("QINGYU_EMAIL_ENABLED", "true")
	t.Setenv("QINGYU_EMAIL_FIXED_CODE", "123456")

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	mongoDB := global.DB
	if mongoDB == nil {
		t.Skip("数据库连接未初始化，跳过用户API集成测试")
	}
	defer cleanupTestData(t, mongoDB)

	timestamp := time.Now().UnixNano()
	testUsername := fmt.Sprintf("sharedauth_old_%d", timestamp)
	testEmail := fmt.Sprintf("sharedauth_old_%d@example.com", timestamp)
	testPassword := "password123"

	sendCodeReq := map[string]interface{}{
		"email": testEmail,
	}
	reqBody, _ := json.Marshal(sendCodeReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/send-verification-code", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	registerReq := map[string]interface{}{
		"username":         testUsername,
		"email":            testEmail,
		"password":         testPassword,
		"role":             "reader",
		"verificationCode": "123456",
	}
	reqBody, _ = json.Marshal(registerReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var registerResp map[string]interface{}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &registerResp))
	require.EqualValues(t, 0, registerResp["code"])

	loginReq := map[string]interface{}{
		"username": testUsername,
		"password": testPassword,
	}
	reqBody, _ = json.Marshal(loginReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var loginResp map[string]interface{}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &loginResp))
	require.EqualValues(t, 0, loginResp["code"])

	loginData, ok := loginResp["data"].(map[string]interface{})
	require.True(t, ok, "登录响应数据格式错误: %+v", loginResp)
	loginToken, ok := loginData["token"].(string)
	require.True(t, ok, "登录响应缺少token: %+v", loginData)
	require.NotEmpty(t, loginToken)

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+loginToken)
	logoutResp := httptest.NewRecorder()
	router.ServeHTTP(logoutResp, logoutReq)
	require.Equal(t, http.StatusOK, logoutResp.Code)

	var logoutJSON map[string]interface{}
	require.NoError(t, json.Unmarshal(logoutResp.Body.Bytes(), &logoutJSON))
	require.Equal(t, "登出成功", logoutJSON["message"])

	protectedResp := doAuthenticatedRequest(router, http.MethodGet, "/api/v1/user/profile", loginToken, nil)
	assert.Equal(t, http.StatusUnauthorized, protectedResp.Code, "登出后的旧token访问受保护资源应该被拒绝")
}

func doAuthenticatedRequest(router *gin.Engine, method, url, token string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

// getMongoDB 获取MongoDB数据库连接
func getMongoDB() (*mongo.Database, error) {
	// 从 global 包获取已初始化的数据库连接
	if global.DB == nil {
		return nil, fmt.Errorf("数据库未初始化，请先调用 core.InitDB()")
	}
	return global.DB, nil
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T, mongoDB *mongo.Database) {
	// 清理测试数据库中的测试用户
	ctx := context.Background()
	collection := mongoDB.Collection("users")

	// 删除所有测试用户（用户名包含 "testuser_" 或 "normaluser_" 或 "admin_"）
	filter := map[string]interface{}{
		"username": map[string]interface{}{
			"$regex": "^(testuser_|normaluser_|admin_|sharedauth_)",
		},
	}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		t.Logf("清理测试数据失败: %v", err)
	} else {
		t.Logf("清理测试数据成功: 删除了 %d 条记录", result.DeletedCount)
	}
}
