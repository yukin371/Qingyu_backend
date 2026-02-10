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
// 发送邮箱验证相关测试
// =========================

// TestUserService_SendEmailVerification_Success 测试发送邮箱验证成功
func TestUserService_SendEmailVerification_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"

	testUser := &usersModel.User{
		Username:      "testuser",
		Email:         testEmail,
		EmailVerified: false,
		Status:        usersModel.UserStatusActive,
	}
	testUser.ID, _ = primitive.ObjectIDFromHex(testUserID)

	req := &user2.SendEmailVerificationRequest{
		UserID: testUserID,
		Email:  testEmail,
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(testUser, nil)

	// Act
	resp, err := service.SendEmailVerification(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Message)
	assert.Equal(t, 1800, resp.ExpiresIn) // 30分钟 = 1800秒

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_SendEmailVerification_UserNotFound 测试发送邮箱验证-用户不存在
func TestUserService_SendEmailVerification_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"

	req := &user2.SendEmailVerificationRequest{
		UserID: testUserID,
		Email:  testEmail,
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.SendEmailVerification(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_SendEmailVerification_EmptyUserID 测试发送邮箱验证-用户ID为空
func TestUserService_SendEmailVerification_EmptyUserID(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.SendEmailVerificationRequest{
		UserID: "",
		Email:  "test@example.com",
	}

	// Act
	resp, err := service.SendEmailVerification(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID和邮箱不能为空")
}

// TestUserService_SendEmailVerification_EmptyEmail 测试发送邮箱验证-邮箱为空
func TestUserService_SendEmailVerification_EmptyEmail(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.SendEmailVerificationRequest{
		UserID: primitive.NewObjectID().Hex(),
		Email:  "",
	}

	// Act
	resp, err := service.SendEmailVerification(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID和邮箱不能为空")
}

// TestUserService_SendEmailVerification_EmailMismatch 测试发送邮箱验证-邮箱不匹配
func TestUserService_SendEmailVerification_EmailMismatch(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"
	differentEmail := "different@example.com"

	testUser := &usersModel.User{
		Username:      "testuser",
		Email:         differentEmail,
		EmailVerified: false,
		Status:        usersModel.UserStatusActive,
	}
	testUser.ID, _ = primitive.ObjectIDFromHex(testUserID)

	req := &user2.SendEmailVerificationRequest{
		UserID: testUserID,
		Email:  testEmail,
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(testUser, nil)

	// Act
	resp, err := service.SendEmailVerification(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱不匹配")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_SendEmailVerification_AlreadyVerified 测试发送邮箱验证-已验证过
func TestUserService_SendEmailVerification_AlreadyVerified(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"

	testUser := &usersModel.User{
		Username:      "testuser",
		Email:         testEmail,
		EmailVerified: true,
		Status:        usersModel.UserStatusActive,
	}
	testUser.ID, _ = primitive.ObjectIDFromHex(testUserID)

	req := &user2.SendEmailVerificationRequest{
		UserID: testUserID,
		Email:  testEmail,
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(testUser, nil)

	// Act
	resp, err := service.SendEmailVerification(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "邮箱已验证", resp.Message)
	assert.Equal(t, 0, resp.ExpiresIn)

	mockUserRepo.AssertExpectations(t)
}

// =========================
// 验证邮箱相关测试
// =========================

// TestUserService_VerifyEmail_Success 测试验证邮箱成功
func TestUserService_VerifyEmail_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"

	testUser := &usersModel.User{
		Username:      "testuser",
		Email:         testEmail,
		EmailVerified: false,
		Status:        usersModel.UserStatusInactive,
	}
	testUser.ID, _ = primitive.ObjectIDFromHex(testUserID)

	// 生成验证码
	tokenManager := NewEmailVerificationTokenManager()
	verificationCode, _ := tokenManager.GenerateCode(ctx, testUserID, testEmail)

	req := &user2.VerifyEmailRequest{
		UserID: testUserID,
		Code:   verificationCode,
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(testUser, nil)
	mockUserRepo.On("Update", ctx, testUserID, mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// Act
	resp, err := service.VerifyEmail(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.Message)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_VerifyEmail_EmptyUserID 测试验证邮箱-用户ID为空
func TestUserService_VerifyEmail_EmptyUserID(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.VerifyEmailRequest{
		UserID: "",
		Code:   "123456",
	}

	// Act
	resp, err := service.VerifyEmail(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID和验证码不能为空")
}

// TestUserService_VerifyEmail_EmptyCode 测试验证邮箱-验证码为空
func TestUserService_VerifyEmail_EmptyCode(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.VerifyEmailRequest{
		UserID: primitive.NewObjectID().Hex(),
		Code:   "",
	}

	// Act
	resp, err := service.VerifyEmail(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID和验证码不能为空")
}

// TestUserService_VerifyEmail_UserNotFound 测试验证邮箱-用户不存在
func TestUserService_VerifyEmail_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	req := &user2.VerifyEmailRequest{
		UserID: testUserID,
		Code:   "123456",
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.VerifyEmail(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_VerifyEmail_InvalidCode 测试验证邮箱-验证码无效
func TestUserService_VerifyEmail_InvalidCode(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"

	testUser := &usersModel.User{
		Username:      "testuser",
		Email:         testEmail,
		EmailVerified: false,
		Status:        usersModel.UserStatusInactive,
	}
	testUser.ID, _ = primitive.ObjectIDFromHex(testUserID)

	req := &user2.VerifyEmailRequest{
		UserID: testUserID,
		Code:   "invalid_code",
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(testUser, nil)

	// Act
	resp, err := service.VerifyEmail(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "验证码无效")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_VerifyEmail_CodeExpired 测试验证邮箱-验证码过期
func TestUserService_VerifyEmail_CodeExpired(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"

	testUser := &usersModel.User{
		Username:      "testuser",
		Email:         testEmail,
		EmailVerified: false,
		Status:        usersModel.UserStatusInactive,
	}
	testUser.ID, _ = primitive.ObjectIDFromHex(testUserID)

	// 使用过期的验证码
	req := &user2.VerifyEmailRequest{
		UserID: testUserID,
		Code:   "expired_code",
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(testUser, nil)

	// Act
	resp, err := service.VerifyEmail(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "验证码无效")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_VerifyEmail_UpdateError 测试验证邮箱-更新用户状态失败
func TestUserService_VerifyEmail_UpdateError(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testEmail := "test@example.com"

	testUser := &usersModel.User{
		Username:      "testuser",
		Email:         testEmail,
		EmailVerified: false,
		Status:        usersModel.UserStatusInactive,
	}
	testUser.ID, _ = primitive.ObjectIDFromHex(testUserID)

	// 生成验证码
	tokenManager := NewEmailVerificationTokenManager()
	verificationCode, _ := tokenManager.GenerateCode(ctx, testUserID, testEmail)

	req := &user2.VerifyEmailRequest{
		UserID: testUserID,
		Code:   verificationCode,
	}

	mockUserRepo.On("GetByID", ctx, testUserID).Return(testUser, nil)
	mockUserRepo.On("Update", ctx, testUserID, mock.AnythingOfType("map[string]interface {}")).Return(assert.AnError)

	// Act
	resp, err := service.VerifyEmail(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "更新用户状态失败")

	mockUserRepo.AssertExpectations(t)
}
