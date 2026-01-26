package user

import (
	"context"
	"errors"
	"testing"
	"time"

	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/service/user/mocks"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestEmailVerificationTokenManager_GenerateCode 测试验证码生成
func TestEmailVerificationTokenManager_GenerateCode(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 测试生成验证码
	code, err := manager.GenerateCode(ctx, "user123", "test@example.com")
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 验证码应该是6位数字
	if len(code) != 6 {
		t.Errorf("验证码长度应该是6位，实际为: %d", len(code))
	}

	// 验证码应该全是数字
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("验证码应该只包含数字，发现: %c", c)
		}
	}
}

// TestEmailVerificationTokenManager_ValidateCode 测试验证码验证
func TestEmailVerificationTokenManager_ValidateCode(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 生成验证码
	userID := "user123"
	email := "test@example.com"
	code, err := manager.GenerateCode(ctx, userID, email)
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 验证正确的验证码
	err = manager.ValidateCode(ctx, userID, email, code)
	if err != nil {
		t.Errorf("验证正确验证码失败: %v", err)
	}

	// 验证错误的验证码
	err = manager.ValidateCode(ctx, userID, email, "000000")
	if err == nil {
		t.Error("应该返回验证码错误")
	}

	// 验证错误用户ID的验证码
	err = manager.ValidateCode(ctx, "wrong_user", email, code)
	if err == nil {
		t.Error("应该返回用户不匹配错误")
	}

	// 验证错误邮箱的验证码
	err = manager.ValidateCode(ctx, userID, "wrong@example.com", code)
	if err == nil {
		t.Error("应该返回验证码不存在错误")
	}
}

// TestEmailVerificationTokenManager_MarkCodeAsUsed 测试标记验证码已使用
func TestEmailVerificationTokenManager_MarkCodeAsUsed(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 生成验证码
	userID := "user123"
	email := "test@example.com"
	code, err := manager.GenerateCode(ctx, userID, email)
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 验证验证码应该成功
	err = manager.ValidateCode(ctx, userID, email, code)
	if err != nil {
		t.Errorf("验证正确验证码失败: %v", err)
	}

	// 标记验证码已使用
	err = manager.MarkCodeAsUsed(ctx, email)
	if err != nil {
		t.Fatalf("标记验证码已使用失败: %v", err)
	}

	// 验证已使用的验证码应该失败
	err = manager.ValidateCode(ctx, userID, email, code)
	if err == nil {
		t.Error("已使用的验证码应该验证失败")
	}
}

// TestEmailVerificationTokenManager_CleanExpiredCodes 测试清理过期验证码
func TestEmailVerificationTokenManager_CleanExpiredCodes(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 生成验证码
	userID := "user123"
	_, err := manager.GenerateCode(ctx, userID, "test1@example.com")
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 手动添加一个过期的验证码
	manager.tokens["expired@example.com"] = &VerificationTokenInfo{
		Code:      "123456",
		Email:     "expired@example.com",
		UserID:    userID,
		ExpiresAt: time.Now().Add(-time.Hour), // 1小时前过期
		Used:      false,
	}

	// 记录清理前的数量
	beforeCount := len(manager.tokens)

	// 清理过期验证码
	manager.CleanExpiredCodes(ctx)

	// 记录清理后的数量
	afterCount := len(manager.tokens)

	// 应该清理了至少1个过期验证码
	if afterCount >= beforeCount {
		t.Error("应该清理至少1个过期验证码")
	}

	// 验证过期的验证码已被清理
	_, exists := manager.tokens["expired@example.com"]
	if exists {
		t.Error("过期验证码应该被清理")
	}
}

// TestEmailVerificationTokenManager_ConcurrentAccess 测试并发访问
func TestEmailVerificationTokenManager_ConcurrentAccess(t *testing.T) {
	manager := NewEmailVerificationTokenManager()
	ctx := context.Background()

	// 并发生成验证码
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(index int) {
			email := "user" + string(rune('0'+index)) + "@example.com"
			userID := "user" + string(rune('0'+index))
			_, err := manager.GenerateCode(ctx, userID, email)
			if err != nil {
				t.Errorf("并发生成验证码失败: %v", err)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// ==================== VerificationService 业务逻辑测试 ====================

// TestVerificationService_SendEmailCode_Success 测试成功发送邮箱验证码
func TestVerificationService_SendEmailCode_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	purpose := "verify_email"

	// 执行测试
	err := service.SendEmailCode(ctx, email, purpose)

	// 验证结果
	assert.NoError(t, err, "发送邮箱验证码应该成功")
}

// TestVerificationService_SendEmailCode_ResetPassword_UserExists 测试重置密码时用户存在
func TestVerificationService_SendEmailCode_ResetPassword_UserExists(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	purpose := "reset_password"

	// 设置mock预期 - 用户存在
	mockUser := &usersModel.User{
		Email: email,
	}
	mockUser.ID = primitive.NewObjectID()
	mockUserRepo.On("GetByEmail", ctx, email).Return(mockUser, nil).Once()

	// 执行测试
	err := service.SendEmailCode(ctx, email, purpose)

	// 验证结果
	assert.NoError(t, err, "发送邮箱验证码应该成功")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_SendEmailCode_ResetPassword_UserNotExists 测试重置密码时用户不存在
func TestVerificationService_SendEmailCode_ResetPassword_UserNotExists(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "nonexistent@example.com"
	purpose := "reset_password"

	// 设置mock预期 - 用户不存在
	mockUserRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("用户不存在")).Once()

	// 执行测试
	err := service.SendEmailCode(ctx, email, purpose)

	// 验证结果
	assert.Error(t, err, "用户不存在时应该返回错误")
	assert.Contains(t, err.Error(), "邮箱不存在", "错误信息应该包含'邮箱不存在'")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_SendPhoneCode_Success 测试成功发送手机验证码
func TestVerificationService_SendPhoneCode_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	phone := "13800138000"
	purpose := "verify_phone"

	// 执行测试
	err := service.SendPhoneCode(ctx, phone, purpose)

	// 验证结果
	assert.NoError(t, err, "发送手机验证码应该成功")
}

// TestVerificationService_VerifyCode_Success 测试验证码验证成功
func TestVerificationService_VerifyCode_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	userID := primitive.NewObjectID()

	// 先生成验证码
	code, err := service.GetVerificationTokenManager().GenerateCode(ctx, userID.Hex(), email)
	assert.NoError(t, err, "生成验证码应该成功")

	// 执行测试 - 直接使用 tokenManager 验证
	err = service.GetVerificationTokenManager().ValidateCode(ctx, userID.Hex(), email, code)

	// 验证结果
	assert.NoError(t, err, "验证码验证应该成功")
}

// TestVerificationService_VerifyCode_InvalidCode 测试验证码错误
func TestVerificationService_VerifyCode_InvalidCode(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	userID := primitive.NewObjectID()
	purpose := "verify_email"

	// 先生成验证码
	_, err := service.GetVerificationTokenManager().GenerateCode(ctx, userID.Hex(), email)
	assert.NoError(t, err, "生成验证码应该成功")

	// 使用错误的验证码
	err = service.VerifyCode(ctx, email, "000000", purpose)

	// 验证结果
	assert.Error(t, err, "使用错误验证码应该返回错误")
	assert.Contains(t, err.Error(), "验证码无效或已过期", "错误信息应该包含'验证码无效或已过期'")
}

// TestVerificationService_VerifyCode_CodeExpired 测试验证码过期
func TestVerificationService_VerifyCode_CodeExpired(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	userID := primitive.NewObjectID()
	purpose := "verify_email"

	// 手动添加一个过期的验证码
	tokenManager := service.GetVerificationTokenManager()
	tokenManager.tokens[email] = &VerificationTokenInfo{
		Code:      "123456",
		Email:     email,
		UserID:    userID.Hex(),
		ExpiresAt: time.Now().Add(-time.Hour), // 1小时前过期
		Used:      false,
	}

	// 执行测试
	err := service.VerifyCode(ctx, email, "123456", purpose)

	// 验证结果
	assert.Error(t, err, "过期验证码应该返回错误")
	assert.Contains(t, err.Error(), "验证码无效或已过期", "错误信息应该包含'验证码无效或已过期'")
}

// TestVerificationService_VerifyCode_CodeUsed 测试验证码已使用
func TestVerificationService_VerifyCode_CodeUsed(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	userID := primitive.NewObjectID()
	purpose := "verify_email"

	// 生成验证码
	code, err := service.GetVerificationTokenManager().GenerateCode(ctx, userID.Hex(), email)
	assert.NoError(t, err, "生成验证码应该成功")

	// 标记验证码为已使用
	err = service.MarkCodeAsUsed(ctx, email)
	assert.NoError(t, err, "标记验证码已使用应该成功")

	// 尝试验证已使用的验证码
	err = service.VerifyCode(ctx, email, code, purpose)

	// 验证结果
	assert.Error(t, err, "已使用的验证码应该返回错误")
	assert.Contains(t, err.Error(), "验证码无效或已过期", "错误信息应该包含'验证码无效或已过期'")
}

// TestVerificationService_MarkCodeAsUsed_Success 测试标记验证码为已使用成功
func TestVerificationService_MarkCodeAsUsed_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	userID := primitive.NewObjectID()

	// 生成验证码
	_, err := service.GetVerificationTokenManager().GenerateCode(ctx, userID.Hex(), email)
	assert.NoError(t, err, "生成验证码应该成功")

	// 执行测试
	err = service.MarkCodeAsUsed(ctx, email)

	// 验证结果
	assert.NoError(t, err, "标记验证码已使用应该成功")

	// 验证验证码已被标记
	tokenManager := service.GetVerificationTokenManager()
	tokenInfo := tokenManager.tokens[email]
	assert.True(t, tokenInfo.Used, "验证码应该被标记为已使用")
}

// TestVerificationService_MarkCodeAsUsed_CodeNotExists 测试标记不存在的验证码
func TestVerificationService_MarkCodeAsUsed_CodeNotExists(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "nonexistent@example.com"

	// 执行测试
	err := service.MarkCodeAsUsed(ctx, email)

	// 验证结果
	assert.Error(t, err, "标记不存在的验证码应该返回错误")
	assert.Contains(t, err.Error(), "验证码不存在", "错误信息应该包含'验证码不存在'")
}

// TestVerificationService_SetEmailVerified_Success 测试设置邮箱验证状态成功
func TestVerificationService_SetEmailVerified_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	email := "test@example.com"

	// 设置mock预期
	mockUser := &usersModel.User{
		Email:         email,
		EmailVerified: false,
	}
	mockUser.ID = userID
	mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()
	mockUserRepo.On("SetEmailVerified", ctx, userID.Hex(), true).Return(nil).Once()

	// 执行测试
	err := service.SetEmailVerified(ctx, userID.Hex(), email)

	// 验证结果
	assert.NoError(t, err, "设置邮箱验证状态应该成功")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_SetEmailVerified_EmailMismatch 测试邮箱不匹配
func TestVerificationService_SetEmailVerified_EmailMismatch(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	wrongEmail := "wrong@example.com"

	// 设置mock预期 - 用户邮箱与提供的不匹配
	mockUser := &usersModel.User{
		Email: "correct@example.com",
	}
	mockUser.ID = userID
	mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()

	// 执行测试
	err := service.SetEmailVerified(ctx, userID.Hex(), wrongEmail)

	// 验证结果
	assert.Error(t, err, "邮箱不匹配时应该返回错误")
	assert.Contains(t, err.Error(), "邮箱不匹配", "错误信息应该包含'邮箱不匹配'")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_SetEmailVerified_UserNotFound 测试用户不存在
func TestVerificationService_SetEmailVerified_UserNotFound(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := "nonexistent"
	email := "test@example.com"

	// 设置mock预期 - 用户不存在
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, errors.New("用户不存在")).Once()

	// 执行测试
	err := service.SetEmailVerified(ctx, userID, email)

	// 验证结果
	assert.Error(t, err, "用户不存在时应该返回错误")
	assert.Contains(t, err.Error(), "获取用户失败", "错误信息应该包含'获取用户失败'")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_SetPhoneVerified_Success 测试设置手机验证状态成功
func TestVerificationService_SetPhoneVerified_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	phone := "13800138000"

	// 设置mock预期
	mockUser := &usersModel.User{
		Phone:         phone,
		PhoneVerified: false,
	}
	mockUser.ID = userID
	mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()
	mockUserRepo.On("SetPhoneVerified", ctx, userID.Hex(), true).Return(nil).Once()

	// 执行测试
	err := service.SetPhoneVerified(ctx, userID.Hex(), phone)

	// 验证结果
	assert.NoError(t, err, "设置手机验证状态应该成功")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_SetPhoneVerified_PhoneMismatch 测试手机号不匹配
func TestVerificationService_SetPhoneVerified_PhoneMismatch(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	wrongPhone := "13900139000"

	// 设置mock预期 - 用户手机号与提供的不匹配
	mockUser := &usersModel.User{
		Phone: "13800138000",
	}
	mockUser.ID = userID
	mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()

	// 执行测试
	err := service.SetPhoneVerified(ctx, userID.Hex(), wrongPhone)

	// 验证结果
	assert.Error(t, err, "手机号不匹配时应该返回错误")
	assert.Contains(t, err.Error(), "手机号不匹配", "错误信息应该包含'手机号不匹配'")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_SetPhoneVerified_UserNotFound 测试用户不存在
func TestVerificationService_SetPhoneVerified_UserNotFound(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := "nonexistent"
	phone := "13800138000"

	// 设置mock预期 - 用户不存在
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, errors.New("用户不存在")).Once()

	// 执行测试
	err := service.SetPhoneVerified(ctx, userID, phone)

	// 验证结果
	assert.Error(t, err, "用户不存在时应该返回错误")
	assert.Contains(t, err.Error(), "获取用户失败", "错误信息应该包含'获取用户失败'")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_CheckPassword_Success 测试检查密码成功
func TestVerificationService_CheckPassword_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	password := "correctPassword"

	// 创建测试用户并设置密码
	mockUser := &usersModel.User{}
	mockUser.ID = userID
	err := mockUser.SetPassword(password)
	assert.NoError(t, err, "设置密码应该成功")

	// 设置mock预期
	mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()

	// 执行测试
	err = service.CheckPassword(ctx, userID.Hex(), password)

	// 验证结果
	assert.NoError(t, err, "密码验证应该成功")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_CheckPassword_WrongPassword 测试密码错误
func TestVerificationService_CheckPassword_WrongPassword(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	correctPassword := "correctPassword"
	wrongPassword := "wrongPassword"

	// 创建测试用户并设置密码
	mockUser := &usersModel.User{}
	mockUser.ID = userID
	err := mockUser.SetPassword(correctPassword)
	assert.NoError(t, err, "设置密码应该成功")

	// 设置mock预期
	mockUserRepo.On("GetByID", ctx, userID.Hex()).Return(mockUser, nil).Once()

	// 执行测试
	err = service.CheckPassword(ctx, userID.Hex(), wrongPassword)

	// 验证结果
	assert.Error(t, err, "密码错误应该返回错误")
	assert.Contains(t, err.Error(), "密码错误", "错误信息应该包含'密码错误'")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_CheckPassword_UserNotFound 测试用户不存在
func TestVerificationService_CheckPassword_UserNotFound(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	userID := "nonexistent"
	password := "password"

	// 设置mock预期 - 用户不存在
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, errors.New("用户不存在")).Once()

	// 执行测试
	err := service.CheckPassword(ctx, userID, password)

	// 验证结果
	assert.Error(t, err, "用户不存在时应该返回错误")
	assert.Contains(t, err.Error(), "获取用户失败", "错误信息应该包含'获取用户失败'")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_EmailExists_True 测试邮箱存在
func TestVerificationService_EmailExists_True(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"

	// 设置mock预期
	mockUserRepo.On("ExistsByEmail", ctx, email).Return(true, nil).Once()

	// 执行测试
	exists, err := service.EmailExists(ctx, email)

	// 验证结果
	assert.NoError(t, err, "检查邮箱是否存在应该成功")
	assert.True(t, exists, "邮箱应该存在")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_EmailExists_False 测试邮箱不存在
func TestVerificationService_EmailExists_False(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "nonexistent@example.com"

	// 设置mock预期
	mockUserRepo.On("ExistsByEmail", ctx, email).Return(false, nil).Once()

	// 执行测试
	exists, err := service.EmailExists(ctx, email)

	// 验证结果
	assert.NoError(t, err, "检查邮箱是否存在应该成功")
	assert.False(t, exists, "邮箱不应该存在")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_PhoneExists_True 测试手机号存在
func TestVerificationService_PhoneExists_True(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	phone := "13800138000"

	// 设置mock预期
	mockUserRepo.On("ExistsByPhone", ctx, phone).Return(true, nil).Once()

	// 执行测试
	exists, err := service.PhoneExists(ctx, phone)

	// 验证结果
	assert.NoError(t, err, "检查手机号是否存在应该成功")
	assert.True(t, exists, "手机号应该存在")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_PhoneExists_False 测试手机号不存在
func TestVerificationService_PhoneExists_False(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	phone := "13900139000"

	// 设置mock预期
	mockUserRepo.On("ExistsByPhone", ctx, phone).Return(false, nil).Once()

	// 执行测试
	exists, err := service.PhoneExists(ctx, phone)

	// 验证结果
	assert.NoError(t, err, "检查手机号是否存在应该成功")
	assert.False(t, exists, "手机号不应该存在")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_GetUserByEmail_Success 测试根据邮箱获取用户成功
func TestVerificationService_GetUserByEmail_Success(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "test@example.com"
	testID := primitive.NewObjectID()
	expectedUser := &usersModel.User{
		Email: email,
	}
	expectedUser.ID = testID

	// 设置mock预期
	mockUserRepo.On("GetByEmail", ctx, email).Return(expectedUser, nil).Once()

	// 执行测试
	user, err := service.GetUserByEmail(ctx, email)

	// 验证结果
	assert.NoError(t, err, "根据邮箱获取用户应该成功")
	assert.Equal(t, expectedUser, user, "返回的用户应该与预期相同")
	mockUserRepo.AssertExpectations(t)
}

// TestVerificationService_GetUserByEmail_NotFound 测试用户不存在
func TestVerificationService_GetUserByEmail_NotFound(t *testing.T) {
	// 创建mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockAuthRepo := new(mocks.MockAuthRepository)

	// 创建验证服务
	service := NewVerificationService(mockUserRepo, mockAuthRepo, nil)

	ctx := context.Background()
	email := "nonexistent@example.com"

	// 设置mock预期
	mockUserRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("用户不存在")).Once()

	// 执行测试
	user, err := service.GetUserByEmail(ctx, email)

	// 验证结果
	assert.Error(t, err, "用户不存在时应该返回错误")
	assert.Nil(t, user, "返回的用户应该为nil")
	mockUserRepo.AssertExpectations(t)
}
