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

// === Mock AuthService ===

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

func (m *MockAuthService) RefreshToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *MockAuthService) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	args := m.Called(ctx, userID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) HasRole(ctx context.Context, userID, role string) (bool, error) {
	args := m.Called(ctx, userID, role)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) CreateRole(ctx context.Context, req *auth.CreateRoleRequest) (*auth.Role, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Role), args.Error(1)
}

func (m *MockAuthService) UpdateRole(ctx context.Context, roleID string, req *auth.UpdateRoleRequest) error {
	args := m.Called(ctx, roleID, req)
	return args.Error(0)
}

func (m *MockAuthService) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockAuthService) AssignRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthService) RemoveRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthService) CreateSession(ctx context.Context, userID string) (*auth.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Session), args.Error(1)
}

func (m *MockAuthService) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Session), args.Error(1)
}

func (m *MockAuthService) DestroySession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockAuthService) RefreshSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockAuthService) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// === 测试辅助函数 ===

func setupAuthTestRouter(authService auth.AuthService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置user_id（用于需要认证的端点）
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := sharedAPI.NewAuthAPI(authService)

	v1 := r.Group("/api/v1/shared/auth")
	{
		v1.POST("/register", api.Register)
		v1.POST("/login", api.Login)
		v1.POST("/logout", api.Logout)
		v1.POST("/refresh", api.RefreshToken)
		v1.GET("/permissions", api.GetUserPermissions)
		v1.GET("/roles", api.GetUserRoles)
	}

	return r
}

// === 测试用例 ===

func TestAuthAPI_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    auth.RegisterRequest
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "成功注册",
			requestBody: auth.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(service *MockAuthService) {
				resp := &auth.RegisterResponse{
					User: &auth.UserInfo{
						ID:       "user123",
						Username: "testuser",
						Email:    "test@example.com",
						Roles:    []string{"user"},
					},
					Token: "test_token_123",
				}
				service.On("Register", mock.Anything, mock.AnythingOfType("*auth.RegisterRequest")).Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "注册成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["token"])
				user := data["user"].(map[string]interface{})
				assert.Equal(t, "user123", user["id"])
				assert.Equal(t, "testuser", user["username"])
			},
		},
		{
			name: "注册失败-服务错误",
			requestBody: auth.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(service *MockAuthService) {
				service.On("Register", mock.Anything, mock.AnythingOfType("*auth.RegisterRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(500), resp["code"])
				assert.Contains(t, resp["message"], "注册失败")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			router := setupAuthTestRouter(mockService, "")

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/shared/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthAPI_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    auth.LoginRequest
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "成功登录",
			requestBody: auth.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			setupMock: func(service *MockAuthService) {
				resp := &auth.LoginResponse{
					User: &auth.UserInfo{
						ID:       "user123",
						Username: "testuser",
						Email:    "test@example.com",
						Roles:    []string{"user"},
					},
					Token: "test_token_123",
				}
				service.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginRequest")).Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "登录成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "test_token_123", data["token"])
				user := data["user"].(map[string]interface{})
				assert.Equal(t, "user123", user["id"])
			},
		},
		{
			name: "登录失败-密码错误",
			requestBody: auth.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setupMock: func(service *MockAuthService) {
				service.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Contains(t, resp["message"], "登录失败")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			router := setupAuthTestRouter(mockService, "")

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/shared/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthAPI_Logout(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:  "成功登出",
			token: "Bearer test_token_123",
			setupMock: func(service *MockAuthService) {
				service.On("Logout", mock.Anything, "test_token_123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "登出成功", resp["message"])
			},
		},
		{
			name:  "未提供Token",
			token: "",
			setupMock: func(service *MockAuthService) {
				// 不设置Mock，因为不会调用Service
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未提供Token", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			router := setupAuthTestRouter(mockService, "")

			req := httptest.NewRequest("POST", "/api/v1/shared/auth/logout", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthAPI_RefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:  "成功刷新Token",
			token: "Bearer old_token_123",
			setupMock: func(service *MockAuthService) {
				service.On("RefreshToken", mock.Anything, "old_token_123").Return("new_token_456", nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "Token刷新成功", resp["message"])
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "new_token_456", data["token"])
			},
		},
		{
			name:  "Token无效",
			token: "Bearer invalid_token",
			setupMock: func(service *MockAuthService) {
				service.On("RefreshToken", mock.Anything, "invalid_token").Return("", assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Contains(t, resp["message"], "Token刷新失败")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			router := setupAuthTestRouter(mockService, "")

			req := httptest.NewRequest("POST", "/api/v1/shared/auth/refresh", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthAPI_GetUserPermissions(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功获取权限",
			userID: "user123",
			setupMock: func(service *MockAuthService) {
				permissions := []string{"read", "write", "delete"}
				service.On("GetUserPermissions", mock.Anything, "user123").Return(permissions, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "获取权限成功", resp["message"])
				data := resp["data"].([]interface{})
				assert.Len(t, data, 3)
				assert.Contains(t, data, "read")
			},
		},
		{
			name:   "未认证",
			userID: "",
			setupMock: func(service *MockAuthService) {
				// 不设置Mock，因为不会调用Service
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未认证", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			router := setupAuthTestRouter(mockService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/shared/auth/permissions", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthAPI_GetUserRoles(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "成功获取角色",
			userID: "user123",
			setupMock: func(service *MockAuthService) {
				roles := []string{"admin", "editor"}
				service.On("GetUserRoles", mock.Anything, "user123").Return(roles, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(200), resp["code"])
				assert.Equal(t, "获取角色成功", resp["message"])
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
				assert.Contains(t, data, "admin")
			},
		},
		{
			name:   "未认证",
			userID: "",
			setupMock: func(service *MockAuthService) {
				// 不设置Mock，因为不会调用Service
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Equal(t, float64(401), resp["code"])
				assert.Equal(t, "未认证", resp["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			router := setupAuthTestRouter(mockService, tt.userID)

			req := httptest.NewRequest("GET", "/api/v1/shared/auth/roles", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			tt.checkResponse(t, resp)
			mockService.AssertExpectations(t)
		})
	}
}
