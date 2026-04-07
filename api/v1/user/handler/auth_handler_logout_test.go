package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_Logout_WithTokenCallsService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer token-123")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("LogoutUser", mock.Anything, &userServiceInterface.LogoutUserRequest{
		Token: "token-123",
	}).Return(&userServiceInterface.LogoutUserResponse{Success: true}, nil).Once()

	handler.Logout(c)

	require.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Logout_WithoutTokenIsIdempotent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
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

	mockService := new(MockUserService)
	handler := NewAuthHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer token-123")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("LogoutUser", mock.Anything, &userServiceInterface.LogoutUserRequest{
		Token: "token-123",
	}).Return((*userServiceInterface.LogoutUserResponse)(nil), serviceInterfaces.NewServiceError("UserService", serviceInterfaces.ErrorTypeInternal, "登出失败", nil)).Once()

	handler.Logout(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
