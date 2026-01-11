package auth

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

	"Qingyu_backend/service/shared/auth"
)

// MockAuthService 模拟AuthService
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

func (m *MockAuthService) OAuthLogin(ctx context.Context, req *auth.OAuthLoginRequest) (*auth.LoginResponse, error) {
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

// setupAuthTestRouter 设置认证测试路由
func setupAuthTestRouter(authService *MockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := NewAuthAPI(authService)

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

// setupAuthTestRouterWithAuth 设置带认证中间件的测试路由
func setupAuthTestRouterWithAuth(authService *MockAuthService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 模拟认证中间件
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewAuthAPI(authService)

	v1 := r.Group("/api/v1/shared/auth")
	{
		v1.GET("/permissions", api.GetUserPermissions)
		v1.GET("/roles", api.GetUserRoles)
	}

	return r
}

// ==================== Register Tests ====================

func TestAuthAPI_Register_Success(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResp := &auth.RegisterResponse{
		User: &auth.UserInfo{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Roles:    []string{"reader"},
		},
		Token: "test-token",
	}

	mockService.On("Register", mock.Anything, &reqBody).Return(expectedResp, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "注册成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuthAPI_Register_MissingRequiredFields(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := map[string]string{
		"username": "testuser",
		// Missing email and password
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthAPI_Register_InvalidEmail(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := map[string]string{
		"username": "testuser",
		"email":    "invalid-email",
		"password": "password123",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthAPI_Register_ShortPassword(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "12345", // Less than 6 characters
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthAPI_Register_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := auth.RegisterRequest{
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	mockService.On("Register", mock.Anything, &reqBody).Return(nil, assert.AnError)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== Login Tests ====================

func TestAuthAPI_Login_Success(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := auth.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	expectedResp := &auth.LoginResponse{
		User: &auth.UserInfo{
			ID:       "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Roles:    []string{"reader"},
		},
		Token: "test-token",
	}

	mockService.On("Login", mock.Anything, &reqBody).Return(expectedResp, nil)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "登录成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuthAPI_Login_MissingUsername(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := map[string]string{
		"password": "password123",
		// Missing username
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthAPI_Login_MissingPassword(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := map[string]string{
		"username": "testuser",
		// Missing password
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthAPI_Login_InvalidCredentials(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	reqBody := auth.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	mockService.On("Login", mock.Anything, &reqBody).Return(nil, assert.AnError)

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== Logout Tests ====================

func TestAuthAPI_Logout_Success(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	token := "test-token"
	mockService.On("Logout", mock.Anything, token).Return(nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "登出成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuthAPI_Logout_NoToken(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthAPI_Logout_WithoutBearerPrefix(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	token := "test-token"
	mockService.On("Logout", mock.Anything, token).Return(nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/logout", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

func TestAuthAPI_Logout_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	token := "test-token"
	mockService.On("Logout", mock.Anything, token).Return(assert.AnError)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== RefreshToken Tests ====================

func TestAuthAPI_RefreshToken_Success(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	oldToken := "old-token"
	newToken := "new-token"
	mockService.On("RefreshToken", mock.Anything, oldToken).Return(newToken, nil)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+oldToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "Token刷新成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, newToken, data["token"])

	mockService.AssertExpectations(t)
}

func TestAuthAPI_RefreshToken_NoToken(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/refresh", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthAPI_RefreshToken_InvalidToken(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouter(mockService)

	token := "invalid-token"
	mockService.On("RefreshToken", mock.Anything, token).Return("", assert.AnError)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/shared/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetUserPermissions Tests ====================

func TestAuthAPI_GetUserPermissions_Success(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	userID := "user123"
	router := setupAuthTestRouterWithAuth(mockService, userID)

	expectedPermissions := []string{"read:books", "write:books", "manage:users"}
	mockService.On("GetUserPermissions", mock.Anything, userID).Return(expectedPermissions, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/shared/auth/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取权限成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuthAPI_GetUserPermissions_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouterWithAuth(mockService, "") // No user ID

	// When
	req, _ := http.NewRequest("GET", "/api/v1/shared/auth/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthAPI_GetUserPermissions_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	userID := "user123"
	router := setupAuthTestRouterWithAuth(mockService, userID)

	mockService.On("GetUserPermissions", mock.Anything, userID).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/shared/auth/permissions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

// ==================== GetUserRoles Tests ====================

func TestAuthAPI_GetUserRoles_Success(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	userID := "user123"
	router := setupAuthTestRouterWithAuth(mockService, userID)

	expectedRoles := []string{"reader", "author"}
	mockService.On("GetUserRoles", mock.Anything, userID).Return(expectedRoles, nil)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/shared/auth/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取角色成功", response["message"])

	mockService.AssertExpectations(t)
}

func TestAuthAPI_GetUserRoles_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	router := setupAuthTestRouterWithAuth(mockService, "") // No user ID

	// When
	req, _ := http.NewRequest("GET", "/api/v1/shared/auth/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthAPI_GetUserRoles_ServiceError(t *testing.T) {
	// Given
	mockService := new(MockAuthService)
	userID := "user123"
	router := setupAuthTestRouterWithAuth(mockService, userID)

	mockService.On("GetUserRoles", mock.Anything, userID).Return(nil, assert.AnError)

	// When
	req, _ := http.NewRequest("GET", "/api/v1/shared/auth/roles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}
