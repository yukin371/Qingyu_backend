package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authsvc "Qingyu_backend/service/auth"
	serviceInterfaces "Qingyu_backend/service/interfaces/base"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAuthHandlerService struct {
	mock.Mock
}

func (m *MockAuthHandlerService) Register(ctx context.Context, req *authsvc.RegisterRequest) (*authsvc.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authsvc.RegisterResponse), args.Error(1)
}

func (m *MockAuthHandlerService) Login(ctx context.Context, req *authsvc.LoginRequest) (*authsvc.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authsvc.LoginResponse), args.Error(1)
}

func (m *MockAuthHandlerService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestAuthHandler_Logout_WithTokenCallsService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthHandlerService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer token-123")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("Logout", mock.Anything, "token-123").Return(nil).Once()

	handler.Logout(c)

	require.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Logout_WithoutTokenIsIdempotent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthHandlerService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth/logout", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Logout_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthHandlerService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer token-123")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("Logout", mock.Anything, "token-123").Return(serviceInterfaces.NewServiceError("AuthService", serviceInterfaces.ErrorTypeInternal, "登出失败", nil)).Once()

	handler.Logout(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthHandlerService)
	handler := NewAuthHandler(mockService)

	body := map[string]string{
		"username": "tester",
		"email":    "tester@example.com",
		"password": "Password@123",
	}
	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("Register", mock.Anything, &authsvc.RegisterRequest{
		Username: "tester",
		Email:    "tester@example.com",
		Password: "Password@123",
	}).Return(&authsvc.RegisterResponse{
		User: &authsvc.UserInfo{
			ID:       "user-1",
			Username: "tester",
			Email:    "tester@example.com",
			Roles:    []string{"reader"},
			Status:   "active",
		},
		Token: "token-123",
	}, nil).Once()

	handler.Register(c)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "创建成功", resp["message"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "user-1", data["user_id"])
	assert.Equal(t, "active", data["status"])
	assert.Equal(t, "token-123", data["token"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAuthHandlerService)
	handler := NewAuthHandler(mockService)

	body := map[string]string{
		"username": "tester",
		"password": "Password@123",
	}
	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("Login", mock.Anything, &authsvc.LoginRequest{
		Username: "tester",
		Password: "Password@123",
		ClientIP: "1.2.3.4",
	}).Return(&authsvc.LoginResponse{
		User: &authsvc.UserInfo{
			ID:       "user-1",
			Username: "tester",
			Email:    "tester@example.com",
			Roles:    []string{"reader"},
			Status:   "active",
		},
		Token: "token-123",
	}, nil).Once()

	handler.Login(c)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "操作成功", resp["message"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "token-123", data["token"])
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "user-1", user["user_id"])
	mockService.AssertExpectations(t)
}
