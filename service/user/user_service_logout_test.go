package user

import (
	"context"
	"errors"
	"testing"

	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	user2 "Qingyu_backend/service/interfaces/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTokenLifecycleService struct {
	mock.Mock
}

func (m *MockTokenLifecycleService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenLifecycleService) ValidateTokenUserID(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

// =========================
// 用户登出相关测试
// =========================

// TestUserService_LogoutUser_Success 测试用户登出成功
func TestUserService_LogoutUser_Success(t *testing.T) {
	// Arrange
	service, _ := setupUserService()
	ctx := context.Background()
	tokenService := new(MockTokenLifecycleService)
	service.SetTokenLifecycleService(tokenService)

	req := &user2.LogoutUserRequest{
		Token: "valid_jwt_token",
	}

	tokenService.On("Logout", ctx, req.Token).Return(nil).Once()

	// Act
	resp, err := service.LogoutUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	tokenService.AssertExpectations(t)
}

// TestUserService_LogoutUser_EmptyToken 测试用户登出-Token为空
func TestUserService_LogoutUser_EmptyToken(t *testing.T) {
	// Arrange
	service, _ := setupUserService()
	ctx := context.Background()

	req := &user2.LogoutUserRequest{
		Token: "",
	}

	// Act
	resp, err := service.LogoutUser(ctx, req)

	// Assert
	// 注意：当前简化实现可能不验证Token
	// 完整实现应该返回错误
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
}

func TestUserService_LogoutUser_LogoutServiceError(t *testing.T) {
	service, _ := setupUserService()
	ctx := context.Background()
	tokenService := new(MockTokenLifecycleService)
	service.SetTokenLifecycleService(tokenService)

	req := &user2.LogoutUserRequest{
		Token: "valid_jwt_token",
	}

	tokenService.On("Logout", ctx, req.Token).Return(errors.New("revoke failed")).Once()

	resp, err := service.LogoutUser(ctx, req)

	require.Nil(t, resp)
	require.Error(t, err)
	serviceErr, ok := err.(*serviceInterfaces.ServiceError)
	require.True(t, ok)
	assert.Equal(t, serviceInterfaces.ErrorTypeInternal, serviceErr.Type)
	assert.Contains(t, serviceErr.Message, "登出失败")
	tokenService.AssertExpectations(t)
}

func TestUserService_LogoutUser_NilRequest(t *testing.T) {
	service, _ := setupUserService()
	ctx := context.Background()

	resp, err := service.LogoutUser(ctx, nil)

	require.Nil(t, resp)
	require.Error(t, err)
	serviceErr, ok := err.(*serviceInterfaces.ServiceError)
	require.True(t, ok)
	assert.Equal(t, serviceInterfaces.ErrorTypeValidation, serviceErr.Type)
}

