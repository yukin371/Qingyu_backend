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

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
	"Qingyu_backend/repository/mongodb/user"
	userRouter "Qingyu_backend/router/user"
	userService "Qingyu_backend/service/user"

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

	// 1. 初始化配置和数据库
	_, err := config.LoadConfig("../../config/config.yaml")
	require.NoError(t, err, "加载配置失败")

	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	// 获取数据库连接
	mongoDB, err := getMongoDB()
	require.NoError(t, err, "获取数据库连接失败")

	// 确保测试结束后清理
	defer cleanupTestData(t, mongoDB)

	// 2. 创建Repository和Service
	userRepo := user.NewMongoUserRepository(mongoDB)
	userSvc := userService.NewUserService(userRepo)

	// 3. 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 4. 注册路由
	apiV1 := router.Group("/api/v1")
	userRouter.RegisterUserRoutes(apiV1, userSvc)

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

		assert.Equal(t, float64(201), response["code"])
		assert.Equal(t, "注册成功", response["message"])

		// 提取用户信息和Token
		data := response["data"].(map[string]interface{})
		t.Logf("注册数据: %+v", data)
		userID = data["user_id"].(string)
		token = data["token"].(string)

		assert.NotEmpty(t, userID, "应该返回用户ID")
		assert.NotEmpty(t, token, "应该返回JWT Token")
		assert.Equal(t, testUsername, data["username"])
		assert.Equal(t, testEmail, data["email"])
		assert.Equal(t, "user", data["role"])
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

		assert.Equal(t, float64(200), response["code"])
		assert.Equal(t, "登录成功", response["message"])

		// 提取新Token
		data := response["data"].(map[string]interface{})
		newToken := data["token"].(string)
		assert.NotEmpty(t, newToken)

		// 更新token（使用登录获得的新token）
		token = newToken

		t.Logf("✓ 用户登录成功，获得新Token")
	})

	// ========== 阶段3：获取个人信息 ==========
	t.Run("获取个人信息", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.Code, "获取个人信息应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, float64(200), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, userID, data["user_id"])
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
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.Code, "更新个人信息应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, float64(200), response["code"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "测试昵称", data["nickname"])
		assert.Equal(t, "这是一个测试用户的个人简介", data["bio"])

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
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/password", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.Code, "修改密码应该返回200")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, float64(200), response["code"])
		assert.Equal(t, "密码修改成功", response["message"])

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

		assert.Equal(t, float64(200), response["code"])

		t.Logf("✓ 使用新密码登录成功")
	})

	// 清理测试数据（注意：这里需要从测试上下文获取数据库连接）
	// 为了简化，暂时跳过清理（实际应该在测试结束后统一清理）
	t.Log("测试用户清理将在测试结束后进行")
}

// testAuthenticationAndAuthorization 测试认证和权限控制
func testAuthenticationAndAuthorization(t *testing.T, router *gin.Engine) {
	// ========== 测试1：未认证访问 ==========
	t.Run("未认证访问需要认证的接口", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		// 故意不设置 Authorization header
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证返回401
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "未认证应该返回401")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, float64(40101), response["code"])
		assert.Equal(t, "未提供认证令牌", response["message"])

		t.Logf("✓ 未认证访问被正确拒绝")
	})

	// ========== 测试2：无效Token ==========
	t.Run("使用无效Token访问", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid_token_123456")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		// 验证返回401
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "无效Token应该返回401")

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, float64(40103), response["code"])

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

		assert.Equal(t, float64(40301), response["code"])

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
			"$regex": "^(testuser_|normaluser_|admin_)",
		},
	}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		t.Logf("清理测试数据失败: %v", err)
	} else {
		t.Logf("清理测试数据成功: 删除了 %d 条记录", result.DeletedCount)
	}
}
