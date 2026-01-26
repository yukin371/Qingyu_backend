package user

import (
	"context"
	"testing"

	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/service/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestNewPasswordService 测试创建密码服务
func TestNewPasswordService(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	verificationService := NewVerificationService(mockUserRepo, mockAuthRepo, nil)
	service := NewPasswordService(verificationService, mockUserRepo)

	assert.NotNil(t, service, "密码服务创建失败")
	assert.NotNil(t, service.verificationService, "验证服务未正确设置")
	assert.NotNil(t, service.userRepo, "用户仓库未正确设置")
}

// TestPasswordService_SendResetCode 测试发送重置验证码
func TestPasswordService_SendResetCode(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	verificationService := NewVerificationService(mockUserRepo, mockAuthRepo, nil)
	service := NewPasswordService(verificationService, mockUserRepo)

	ctx := context.Background()
	email := "test@example.com"

	// 设置mock预期 - 用户存在
	userID := primitive.NewObjectID()
	mockUser := &usersModel.User{
		Email: email,
	}
	mockUser.ID = userID
	// GetByEmail会被调用两次：
	// 1. PasswordService.SendResetCode中检查邮箱是否存在
	// 2. VerificationService.SendEmailCode中检查邮箱是否存在（当purpose是reset_password时）
	mockUserRepo.On("GetByEmail", ctx, email).Return(mockUser, nil).Times(2)

	// 执行测试
	err := service.SendResetCode(ctx, email)
	assert.NoError(t, err, "发送重置验证码应该成功")

	mockUserRepo.AssertExpectations(t)
}

// TestPasswordService_ResetPassword 测试重置密码
func TestPasswordService_ResetPassword(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	verificationService := NewVerificationService(mockUserRepo, mockAuthRepo, nil)
	service := NewPasswordService(verificationService, mockUserRepo)

	ctx := context.Background()
	email := "test@example.com"
	userID := primitive.NewObjectID()

	// 设置mock预期
	mockUser := &usersModel.User{
		Email: email,
	}
	mockUser.ID = userID

	// 测试用例1：成功重置密码
	t.Run("Success", func(t *testing.T) {
		// GetByEmail会被调用两次：
		// 1. PasswordService.SendResetCode中检查邮箱是否存在
		// 2. VerificationService.SendEmailCode中检查邮箱是否存在（当purpose是reset_password时）
		mockUserRepo.On("GetByEmail", ctx, email).Return(mockUser, nil).Times(2)

		// 发送重置验证码
		err := service.SendResetCode(ctx, email)
		assert.NoError(t, err, "发送重置验证码应该成功")

		mockUserRepo.AssertExpectations(t)
	})

	// 测试用例2：邮箱不存在
	t.Run("EmailNotExists", func(t *testing.T) {
		// 当邮箱不存在时，GetByEmail只会在PasswordService.SendResetCode中被调用1次
		// SendEmailCode会因为邮箱不存在而提前返回，不会再调用GetByEmail
		mockUserRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, context.Canceled).Once()

		// 发送重置验证码到不存在的邮箱
		err := service.SendResetCode(ctx, "nonexistent@example.com")
		assert.Error(t, err, "应该返回邮箱不存在错误")

		mockUserRepo.AssertExpectations(t)
	})
}

// TestPasswordService_UpdatePassword 测试修改密码
func TestPasswordService_UpdatePassword(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	verificationService := NewVerificationService(mockUserRepo, mockAuthRepo, nil)
	service := NewPasswordService(verificationService, mockUserRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// 创建测试用户并设置密码
	mockUser := &usersModel.User{}
	mockUser.ID = userID
	err := mockUser.SetPassword("password")
	assert.NoError(t, err, "设置密码应该成功")

	// 测试用例1：使用正确的旧密码
	t.Run("Success", func(t *testing.T) {
		mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()
		mockUserRepo.On("UpdatePassword", ctx, userID.Hex(), mock.AnythingOfType("string")).Return(nil).Once()

		err = service.UpdatePassword(ctx, userID.Hex(), "password", "NewPassword123")
		assert.NoError(t, err, "修改密码应该成功")

		mockUserRepo.AssertExpectations(t)
	})

	// 测试用例2：使用错误的旧密码
	t.Run("WrongOldPassword", func(t *testing.T) {
		mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()

		err = service.UpdatePassword(ctx, userID.Hex(), "WrongPassword", "NewPassword123")
		assert.Error(t, err, "应该返回旧密码错误")
		assert.Equal(t, ErrOldPasswordMismatch, err, "应该返回旧密码错误")

		mockUserRepo.AssertExpectations(t)
	})

	// 测试用例3：用户不存在
	t.Run("UserNotFound", func(t *testing.T) {
		// 使用一个有效的ObjectID格式，但GetByID返回错误
		fakeUserID := primitive.NewObjectID().Hex()
		mockUserRepo.On("GetByID", ctx, fakeUserID).Return(nil, context.Canceled).Once()

		err = service.UpdatePassword(ctx, fakeUserID, "password", "NewPassword123")
		assert.Error(t, err, "应该返回用户不存在错误")

		mockUserRepo.AssertExpectations(t)
	})
}

// TestPasswordService_GetUserByEmail 测试根据邮箱获取用户
func TestPasswordService_GetUserByEmail(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	verificationService := NewVerificationService(mockUserRepo, mockAuthRepo, nil)
	service := NewPasswordService(verificationService, mockUserRepo)

	ctx := context.Background()
	email := "test@example.com"

	// 设置mock预期
	userID := primitive.NewObjectID()
	mockUser := &usersModel.User{
		Email: email,
	}
	mockUser.ID = userID
	mockUserRepo.On("GetByEmail", ctx, email).Return(mockUser, nil).Once()

	// 获取存在的用户
	user, err := service.GetUserByEmail(ctx, email)
	assert.NoError(t, err, "获取用户应该成功")
	assert.Equal(t, email, user.Email, "获取的用户邮箱不匹配")

	// 重置mock
	mockUserRepo.ExpectedCalls = nil

	// 设置mock预期 - 用户不存在
	mockUserRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, context.Canceled).Once()

	// 获取不存在的用户
	_, err = service.GetUserByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err, "应该返回用户不存在错误")

	mockUserRepo.AssertExpectations(t)
}
