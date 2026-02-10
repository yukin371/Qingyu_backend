package user

import (
	"context"
	"testing"

	usersModel "Qingyu_backend/models/users"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	user2 "Qingyu_backend/service/interfaces/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =========================
// 密码重置相关测试
// =========================

// TestUserService_ResetPassword_Success 测试重置密码成功
func TestUserService_ResetPassword_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "test@example.com"
	testUser := &usersModel.User{
		Username: "testuser",
		Email:    testEmail,
		Status:   usersModel.UserStatusActive,
	}
	testUser.ID = primitive.NewObjectID()

	req := &user2.ResetPasswordRequest{
		Email: testEmail,
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(testUser, nil)

	// Act
	resp, err := service.ResetPassword(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_ResetPassword_UserNotFound 测试重置密码-用户不存在（为了安全应返回成功）
func TestUserService_ResetPassword_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "nonexistent@example.com"

	req := &user2.ResetPasswordRequest{
		Email: testEmail,
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.ResetPassword(ctx, req)

	// Assert
	// 为了安全，即使用户不存在也返回成功（防止邮箱枚举攻击）
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_ResetPassword_EmptyEmail 测试重置密码-邮箱为空
func TestUserService_ResetPassword_EmptyEmail(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.ResetPasswordRequest{
		Email: "",
	}

	// Act
	resp, err := service.ResetPassword(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱不能为空")
}

// TestUserService_ResetPassword_RepositoryError 测试重置密码-仓库错误
func TestUserService_ResetPassword_RepositoryError(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "test@example.com"

	req := &user2.ResetPasswordRequest{
		Email: testEmail,
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(nil, assert.AnError)

	// Act
	resp, err := service.ResetPassword(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "检查用户失败")

	mockUserRepo.AssertExpectations(t)
}

// =========================
// 请求密码重置相关测试
// =========================

// TestUserService_RequestPasswordReset_Success 测试请求密码重置成功
func TestUserService_RequestPasswordReset_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "test@example.com"
	testUser := &usersModel.User{
		Username: "testuser",
		Email:    testEmail,
		Status:   usersModel.UserStatusActive,
	}
	testUser.ID = primitive.NewObjectID()

	req := &user2.RequestPasswordResetRequest{
		Email: testEmail,
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(testUser, nil)

	// Act
	resp, err := service.RequestPasswordReset(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Message)
	assert.Equal(t, 3600, resp.ExpiresIn) // 1小时

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_RequestPasswordReset_UserNotFound 测试请求密码重置-用户不存在（为了安全应返回成功）
func TestUserService_RequestPasswordReset_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "nonexistent@example.com"

	req := &user2.RequestPasswordResetRequest{
		Email: testEmail,
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.RequestPasswordReset(ctx, req)

	// Assert
	// 为了安全，即使用户不存在也返回成功（防止邮箱枚举攻击）
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Message, "如果该邮箱已注册")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_RequestPasswordReset_EmptyEmail 测试请求密码重置-邮箱为空
func TestUserService_RequestPasswordReset_EmptyEmail(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.RequestPasswordResetRequest{
		Email: "",
	}

	// Act
	resp, err := service.RequestPasswordReset(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱不能为空")
}

// =========================
// 确认密码重置相关测试
// =========================

// TestUserService_ConfirmPasswordReset_Success 测试确认密码重置成功
func TestUserService_ConfirmPasswordReset_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "test@example.com"
	testPassword := "newPassword123"

	testUser := &usersModel.User{
		Username: "testuser",
		Email:    testEmail,
		Status:   usersModel.UserStatusActive,
	}
	testUser.ID = primitive.NewObjectID()

	// 使用全局Token管理器生成一个有效的Token
	tokenManager := GetGlobalPasswordResetTokenManager()
	validToken, _ := tokenManager.GenerateToken(ctx, testEmail)

	req := &user2.ConfirmPasswordResetRequest{
		Email:    testEmail,
		Token:    validToken,
		Password: testPassword,
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(testUser, nil)
	mockUserRepo.On("UpdatePassword", ctx, testUser.ID.Hex(), mock.AnythingOfType("string")).Return(nil)

	// Act
	resp, err := service.ConfirmPasswordReset(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Message)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_ConfirmPasswordReset_EmptyEmail 测试确认密码重置-邮箱为空
func TestUserService_ConfirmPasswordReset_EmptyEmail(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.ConfirmPasswordResetRequest{
		Email:    "",
		Token:    "some_token",
		Password: "newPassword123",
	}

	// Act
	resp, err := service.ConfirmPasswordReset(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱、Token和新密码不能为空")
}

// TestUserService_ConfirmPasswordReset_EmptyToken 测试确认密码重置-Token为空
func TestUserService_ConfirmPasswordReset_EmptyToken(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.ConfirmPasswordResetRequest{
		Email:    "test@example.com",
		Token:    "",
		Password: "newPassword123",
	}

	// Act
	resp, err := service.ConfirmPasswordReset(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱、Token和新密码不能为空")
}

// TestUserService_ConfirmPasswordReset_EmptyPassword 测试确认密码重置-密码为空
func TestUserService_ConfirmPasswordReset_EmptyPassword(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.ConfirmPasswordResetRequest{
		Email:    "test@example.com",
		Token:    "some_token",
		Password: "",
	}

	// Act
	resp, err := service.ConfirmPasswordReset(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱、Token和新密码不能为空")
}

// TestUserService_ConfirmPasswordReset_UserNotFound 测试确认密码重置-用户不存在
func TestUserService_ConfirmPasswordReset_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "nonexistent@example.com"

	req := &user2.ConfirmPasswordResetRequest{
		Email:    testEmail,
		Token:    "some_token",
		Password: "newPassword123",
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.ConfirmPasswordReset(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_ConfirmPasswordReset_InvalidToken 测试确认密码重置-Token无效
func TestUserService_ConfirmPasswordReset_InvalidToken(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "test@example.com"
	testUser := &usersModel.User{
		Username: "testuser",
		Email:    testEmail,
		Status:   usersModel.UserStatusActive,
	}
	testUser.ID = primitive.NewObjectID()

	req := &user2.ConfirmPasswordResetRequest{
		Email:    testEmail,
		Token:    "invalid_token",
		Password: "newPassword123",
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(testUser, nil)

	// Act
	resp, err := service.ConfirmPasswordReset(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Token验证失败")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_ConfirmPasswordReset_TokenExpired 测试确认密码重置-Token过期
func TestUserService_ConfirmPasswordReset_TokenExpired(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testEmail := "test@example.com"
	testUser := &usersModel.User{
		Username: "testuser",
		Email:    testEmail,
		Status:   usersModel.UserStatusActive,
	}
	testUser.ID = primitive.NewObjectID()

	// 使用过期的Token（这里需要手动创建一个过期场景）
	req := &user2.ConfirmPasswordResetRequest{
		Email:    testEmail,
		Token:    "expired_token",
		Password: "newPassword123",
	}

	mockUserRepo.On("GetByEmail", ctx, testEmail).Return(testUser, nil)

	// Act
	resp, err := service.ConfirmPasswordReset(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Token验证失败")

	mockUserRepo.AssertExpectations(t)
}
