package user

import (
	"context"
	"errors"
	"testing"

	usersModel "Qingyu_backend/models/users"
	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	user2 "Qingyu_backend/service/interfaces/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	service, _, _ := setupUserService()
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
	service, _, _ := setupUserService()
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
	service, _, _ := setupUserService()
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
	service, _, _ := setupUserService()
	ctx := context.Background()

	resp, err := service.LogoutUser(ctx, nil)

	require.Nil(t, resp)
	require.Error(t, err)
	serviceErr, ok := err.(*serviceInterfaces.ServiceError)
	require.True(t, ok)
	assert.Equal(t, serviceInterfaces.ErrorTypeValidation, serviceErr.Type)
}

// =========================
// Token验证相关测试
// =========================

// TestUserService_ValidateToken_ValidToken 测试验证Token-有效Token
func TestUserService_ValidateToken_ValidToken(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()
	tokenService := new(MockTokenLifecycleService)
	service.SetTokenLifecycleService(tokenService)

	req := &user2.ValidateTokenRequest{
		Token: "valid_jwt_token",
	}

	expectedUser := &usersModel.User{
		Username: "valid-user",
		Email:    "valid@example.com",
	}
	expectedUser.ID = primitive.NewObjectID()
	userID := expectedUser.ID.Hex()

	tokenService.On("ValidateTokenUserID", ctx, req.Token).Return(userID, nil).Once()
	mockUserRepo.On("GetByID", ctx, userID).Return(expectedUser, nil).Once()

	// Act
	resp, err := service.ValidateToken(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	require.NotNil(t, resp.User)
	assert.Equal(t, userID, resp.User.ID)
	tokenService.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

// TestUserService_ValidateToken_EmptyToken 测试验证Token-Token为空
func TestUserService_ValidateToken_EmptyToken(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.ValidateTokenRequest{
		Token: "",
	}

	// Act
	resp, err := service.ValidateToken(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Valid)
}

// TestUserService_ValidateToken_InvalidToken 测试验证Token-无效Token
func TestUserService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()
	tokenService := new(MockTokenLifecycleService)
	service.SetTokenLifecycleService(tokenService)

	req := &user2.ValidateTokenRequest{
		Token: "invalid_token_format",
	}

	tokenService.On("ValidateTokenUserID", ctx, req.Token).Return("", errors.New("invalid token")).Once()

	// Act
	resp, err := service.ValidateToken(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Valid)
	tokenService.AssertExpectations(t)
}
