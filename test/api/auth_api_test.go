//go:build integration
// +build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	sharedAPI "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/shared/auth"
)

// MockAuthService Auth Service Mock
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	args := m.Called(ctx, oldToken)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// AuthAPITestSuite 认证API测试套件
type AuthAPITestSuite struct {
	authAPI     *sharedAPI.AuthAPI
	router      *gin.Engine
	mockService *MockAuthService
}

// setupAuthAPITest 设置认证API测试环境
func setupAuthAPITest(t *testing.T) *AuthAPITestSuite {
	gin.SetMode(gin.TestMode)

	// 创建Mock Service
	mockService := new(MockAuthService)

	// 创建AuthAPI
	authAPI := sharedAPI.NewAuthAPI(mockService)

	// 设置路由
	router := gin.New()
	router.Use(gin.Recovery())

	// 公开路由
	auth := router.Group("/api/v1/shared/auth")
	{
		auth.POST("/register", authAPI.Register)
		auth.POST("/login", authAPI.Login)
		auth.POST("/logout", authAPI.Logout)
		auth.POST("/refresh", authAPI.RefreshToken)

		// 需要认证的路由
		authenticated := auth.Group("")
		authenticated.Use(mockAuthMiddleware())
		{
			authenticated.GET("/permissions", authAPI.GetUserPermissions)
			authenticated.GET("/roles", authAPI.GetUserRoles)
		}
	}

	return &AuthAPITestSuite{
		authAPI:     authAPI,
		router:      router,
		mockService: mockService,
	}
}

// TestAuthRegister_Success 测试注册成功
func TestAuthRegister_Success(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 设置Mock期望
	expectedResp := &auth.RegisterResponse{
		UserID:   "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Token:    "test-token-123",
	}
	suite.mockService.On("Register", mock.Anything, mock.AnythingOfType("*auth.RegisterRequest")).Return(expectedResp, nil)

	// 准备请求数据
	reqBody := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "注册成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, expectedResp.Username, data["username"])
	assert.Equal(t, expectedResp.Email, data["email"])
	assert.NotEmpty(t, data["token"])

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestAuthLogin_Success 测试登录成功
func TestAuthLogin_Success(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 设置Mock期望
	expectedResp := &auth.LoginResponse{
		UserID:   "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Token:    "test-token-123",
	}
	suite.mockService.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginRequest")).Return(expectedResp, nil)

	// 准备请求数据
	reqBody := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "登录成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, expectedResp.Username, data["username"])
	assert.NotEmpty(t, data["token"])

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestAuthLogin_InvalidCredentials 测试登录失败 - 无效凭证
func TestAuthLogin_InvalidCredentials(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 设置Mock期望 - 返回错误
	suite.mockService.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginRequest")).
		Return(nil, assert.AnError)

	// 准备请求数据
	reqBody := map[string]interface{}{
		"username": "testuser",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应 - 应该返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(401), response["code"])

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestAuthLogout_Success 测试登出成功
func TestAuthLogout_Success(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 设置Mock期望
	suite.mockService.On("Logout", mock.Anything, "test-token").Return(nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "登出成功", response["message"])

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestAuthLogout_NoToken 测试登出失败 - 缺少Token
func TestAuthLogout_NoToken(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 创建请求 - 不带Authorization header
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应 - 应该返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(401), response["code"])
}

// TestRefreshToken_Success 测试刷新Token成功
func TestRefreshToken_Success(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 设置Mock期望
	newToken := "new-token-456"
	suite.mockService.On("RefreshToken", mock.Anything, "old-token").Return(newToken, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer old-token")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "Token刷新成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, newToken, data["token"])

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestRefreshToken_ExpiredToken 测试刷新Token失败 - Token过期
func TestRefreshToken_ExpiredToken(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 设置Mock期望 - 返回错误
	suite.mockService.On("RefreshToken", mock.Anything, "expired-token").
		Return("", assert.AnError)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应 - 应该返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(401), response["code"])

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestGetUserPermissions_Success 测试获取用户权限成功
func TestGetUserPermissions_Success(t *testing.T) {
	suite := setupAuthAPITest(t)

	testUserID := "user-123"
	expectedPermissions := []string{
		"read:books",
		"write:articles",
		"manage:profile",
	}

	// 设置Mock期望
	suite.mockService.On("GetUserPermissions", mock.Anything, testUserID).
		Return(expectedPermissions, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/shared/auth/permissions", nil)
	req.Header.Set("X-Test-User-ID", testUserID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取权限成功", response["message"])

	// 验证权限列表
	permissions := response["data"].([]interface{})
	assert.Len(t, permissions, len(expectedPermissions))

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestGetUserRoles_Success 测试获取用户角色成功
func TestGetUserRoles_Success(t *testing.T) {
	suite := setupAuthAPITest(t)

	testUserID := "user-123"
	expectedRoles := []string{"user", "author"}

	// 设置Mock期望
	suite.mockService.On("GetUserRoles", mock.Anything, testUserID).
		Return(expectedRoles, nil)

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/v1/shared/auth/roles", nil)
	req.Header.Set("X-Test-User-ID", testUserID)
	w := httptest.NewRecorder()

	// 执行请求
	suite.router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取角色成功", response["message"])

	// 验证角色列表
	roles := response["data"].([]interface{})
	assert.Len(t, roles, len(expectedRoles))

	// 验证Mock调用
	suite.mockService.AssertExpectations(t)
}

// TestAuth_ValidationErrors 测试参数验证错误
func TestAuth_ValidationErrors(t *testing.T) {
	suite := setupAuthAPITest(t)

	tests := []struct {
		name     string
		endpoint string
		method   string
		body     map[string]interface{}
	}{
		{
			name:     "注册 - 缺少用户名",
			endpoint: "/api/v1/shared/auth/register",
			method:   http.MethodPost,
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
		},
		{
			name:     "注册 - 无效的邮箱",
			endpoint: "/api/v1/shared/auth/register",
			method:   http.MethodPost,
			body: map[string]interface{}{
				"username": "testuser",
				"email":    "invalid-email",
				"password": "password123",
			},
		},
		{
			name:     "登录 - 缺少密码",
			endpoint: "/api/v1/shared/auth/login",
			method:   http.MethodPost,
			body: map[string]interface{}{
				"username": "testuser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, tt.endpoint, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.router.ServeHTTP(w, req)

			// 应该返回400 Bad Request
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

// TestAuth_RequiresAuthentication 测试需要认证的接口
func TestAuth_RequiresAuthentication(t *testing.T) {
	suite := setupAuthAPITest(t)

	// 创建不使用认证中间件的路由
	router := gin.New()
	router.GET("/api/v1/shared/auth/permissions", suite.authAPI.GetUserPermissions)

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "获取权限 - 未认证",
			url:  "/api/v1/shared/auth/permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// 应该返回401 Unauthorized
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}
