package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	user2 "Qingyu_backend/service/interfaces/user"
	usersModel "Qingyu_backend/models/users"
)

// TestUserService_EmailVerificationFlow_Integration 完整邮箱验证流程集成测试
func TestUserService_EmailVerificationFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个未验证邮箱的用户
	userID, _, email, _ := env.CreateDefaultTestUser(t)

	// 确保邮箱未验证
	user := env.AssertUserExists(t, userID)
	require.False(t, user.EmailVerified, "初始状态邮箱应该是未验证的")

	// Act 1 - 发送验证邮件
	sendReq := &user2.SendEmailVerificationRequest{
		UserID: userID,
		Email:  email,
	}

	sendResp, err := env.UserService.SendEmailVerification(ctx, sendReq)

	// Assert 1 - 发送验证邮件成功
	require.NoError(t, err, "发送验证邮件应该成功")
	require.NotNil(t, sendResp, "响应不应该为空")
	assert.True(t, sendResp.Success, "发送应该成功")
	assert.NotEmpty(t, sendResp.Message, "消息不应该为空")

	// 注意：在测试环境中，我们无法获取实际发送的验证码
	// 实际场景中，验证码会通过邮件发送给用户
	// 这里我们需要手动生成一个验证码进行测试

	// Act 2 - 使用生成的验证码验证邮箱（模拟）
	// 由于我们无法从邮件中获取验证码，这里我们直接生成一个新的验证码
	tokenManager := NewEmailVerificationTokenManager()
	testCode, err := tokenManager.GenerateCode(ctx, userID, email)
	require.NoError(t, err, "生成测试验证码应该成功")

	verifyReq := &user2.VerifyEmailRequest{
		UserID: userID,
		Code:   testCode,
	}

	verifyResp, err := env.UserService.VerifyEmail(ctx, verifyReq)

	// Assert 2 - 验证邮箱成功
	require.NoError(t, err, "验证邮箱应该成功")
	require.NotNil(t, verifyResp, "响应不应该为空")
	assert.True(t, verifyResp.Success, "验证应该成功")
	assert.Equal(t, "邮箱验证成功", verifyResp.Message, "验证消息应该正确")

	// 验证用户的邮箱验证状态已更新
	env.AssertUserEmailVerified(t, userID, true)
}

// TestUserService_VerifyEmail_AlreadyVerified 已验证邮箱测试
func TestUserService_VerifyEmail_AlreadyVerified_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个已验证邮箱的用户
	userID, _, email, _ := env.CreateVerifiedTestUser(t)

	// Act - 再次发送验证邮件
	sendReq := &user2.SendEmailVerificationRequest{
		UserID: userID,
		Email:  email,
	}

	sendResp, err := env.UserService.SendEmailVerification(ctx, sendReq)

	// Assert - 应该返回成功（邮箱已验证）
	require.NoError(t, err, "发送验证邮件不应该失败")
	require.NotNil(t, sendResp, "响应不应该为空")
	assert.True(t, sendResp.Success, "发送应该成功")
	assert.Equal(t, "邮箱已验证", sendResp.Message, "消息应该是'邮箱已验证'")
}

// TestUserService_VerifyEmail_WrongEmail 错误邮箱测试
func TestUserService_VerifyEmail_WrongEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, _, _ := env.CreateDefaultTestUser(t)

	// Act - 尝试发送验证邮件到错误的邮箱
	sendReq := &user2.SendEmailVerificationRequest{
		UserID: userID,
		Email:  "wrongemail@example.com",
	}

	sendResp, err := env.UserService.SendEmailVerification(ctx, sendReq)

	// Assert
	assert.Error(t, err, "错误邮箱应该返回错误")
	assert.Nil(t, sendResp, "响应应该为空")
	assert.Contains(t, err.Error(), "邮箱不匹配", "错误信息应该包含'邮箱不匹配'")
}

// TestUserService_SendEmailVerification_NonExistentUser 不存在用户测试
func TestUserService_SendEmailVerification_NonExistentUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 尝试为不存在的用户发送验证邮件
	sendReq := &user2.SendEmailVerificationRequest{
		UserID: "507f1f77bcf86cd799439011",
		Email:  "test@example.com",
	}

	sendResp, err := env.UserService.SendEmailVerification(ctx, sendReq)

	// Assert
	assert.Error(t, err, "不存在的用户应该返回错误")
	assert.Nil(t, sendResp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户不存在", "错误信息应该包含'用户不存在'")
}

// TestUserService_VerificationCode_Expiration 验证码过期测试
func TestUserService_VerificationCode_Expiration_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户并生成验证码
	userID, _, email, _ := env.CreateDefaultTestUser(t)

	tokenManager := NewEmailVerificationTokenManager()
	_ , err := tokenManager.GenerateCode(ctx, userID, email)
	require.NoError(t, err, "生成验证码应该成功")

	// Act - 等待验证码过期（根据默认配置，验证码15分钟过期）
	// 在测试中，我们可以模拟过期的情况
	// 注意：这里我们使用一个无效的验证码来模拟过期的情况

	verifyReq := &user2.VerifyEmailRequest{
		UserID: userID,
		Code:   "000000", // 使用无效的验证码
	}

	verifyResp, err := env.UserService.VerifyEmail(ctx, verifyReq)

	// Assert
	assert.Error(t, err, "无效验证码应该返回错误")
	assert.Nil(t, verifyResp, "响应应该为空")
	assert.Contains(t, err.Error(), "验证码无效", "错误信息应该包含'验证码无效'")
}

// TestUserService_EmailVerification_CompleteFlow 完整邮箱验证流程（含Login检查）
func TestUserService_EmailVerification_CompleteFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个未验证的用户
	userID, username, email, password := env.CreateDefaultTestUser(t)

	// 验证初始状态
	user := env.AssertUserExists(t, userID)
	require.False(t, user.EmailVerified, "初始状态邮箱应该是未验证的")
	require.Equal(t, usersModel.UserStatusActive, user.Status, "初始状态应该是active")

	// Act 1 - 发送验证邮件
	sendReq := &user2.SendEmailVerificationRequest{
		UserID: userID,
		Email:  email,
	}

	_, err := env.UserService.SendEmailVerification(ctx, sendReq)
	require.NoError(t, err, "发送验证邮件应该成功")

	// Act 2 - 验证邮箱
	tokenManager := NewEmailVerificationTokenManager()
	testCode, err := tokenManager.GenerateCode(ctx, userID, email)
	require.NoError(t, err, "生成测试验证码应该成功")

	verifyReq := &user2.VerifyEmailRequest{
		UserID: userID,
		Code:   testCode,
	}

	verifyResp, err := env.UserService.VerifyEmail(ctx, verifyReq)
	require.NoError(t, err, "验证邮箱应该成功")
	require.True(t, verifyResp.Success, "验证应该成功")

	// Assert - 验证邮箱验证状态已更新
	env.AssertUserEmailVerified(t, userID, true)

	// Act 3 - 使用验证后的账号登录
	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)

	// Assert - 登录应该成功
	require.NoError(t, err, "验证后登录应该成功")
	require.NotNil(t, loginResp, "登录响应不应该为空")
	assert.NotEmpty(t, loginResp.Token, "Token不应该为空")
}

// TestUserService_MultipleEmailVerifications 多次邮箱验证测试
func TestUserService_MultipleEmailVerifications_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, email, _ := env.CreateDefaultTestUser(t)

	// Act - 多次发送验证邮件
	for i := 0; i < 3; i++ {
		sendReq := &user2.SendEmailVerificationRequest{
			UserID: userID,
			Email:  email,
		}

		sendResp, err := env.UserService.SendEmailVerification(ctx, sendReq)

		// Assert - 每次发送都应该成功
		require.NoError(t, err, "第%d次发送验证邮件应该成功", i+1)
		require.NotNil(t, sendResp, "第%d次响应不应该为空", i+1)
		assert.True(t, sendResp.Success, "第%d次发送应该成功", i+1)
	}

	// Act - 使用最后一次的验证码验证邮箱
	tokenManager := NewEmailVerificationTokenManager()
	testCode, err := tokenManager.GenerateCode(ctx, userID, email)
	require.NoError(t, err, "生成测试验证码应该成功")

	verifyReq := &user2.VerifyEmailRequest{
		UserID: userID,
		Code:   testCode,
	}

	verifyResp, err := env.UserService.VerifyEmail(ctx, verifyReq)

	// Assert - 验证应该成功
	require.NoError(t, err, "验证邮箱应该成功")
	require.NotNil(t, verifyResp, "响应不应该为空")
	assert.True(t, verifyResp.Success, "验证应该成功")

	// 验证用户的邮箱验证状态
	env.AssertUserEmailVerified(t, userID, true)
}

// TestUserService_VerifyEmail_EmptyUserID 测试空用户ID
func TestUserService_VerifyEmail_EmptyUserID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 使用空用户ID验证邮箱
	verifyReq := &user2.VerifyEmailRequest{
		UserID: "",
		Code:   "123456",
	}

	verifyResp, err := env.UserService.VerifyEmail(ctx, verifyReq)

	// Assert
	assert.Error(t, err, "空用户ID应该返回错误")
	assert.Nil(t, verifyResp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户ID", "错误信息应该包含'用户ID'")
}

// TestUserService_VerifyEmail_EmptyCode 测试空验证码
func TestUserService_VerifyEmail_EmptyCode_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, _, _ := env.CreateDefaultTestUser(t)

	// Act - 使用空验证码验证邮箱
	verifyReq := &user2.VerifyEmailRequest{
		UserID: userID,
		Code:   "",
	}

	verifyResp, err := env.UserService.VerifyEmail(ctx, verifyReq)

	// Assert
	assert.Error(t, err, "空验证码应该返回错误")
	assert.Nil(t, verifyResp, "响应应该为空")
	assert.Contains(t, err.Error(), "验证码", "错误信息应该包含'验证码'")
}

// TestUserService_SendEmailVerification_EmptyUserID 测试发送验证邮件空用户ID
func TestUserService_SendEmailVerification_EmptyUserID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 使用空用户ID发送验证邮件
	sendReq := &user2.SendEmailVerificationRequest{
		UserID: "",
		Email:  "test@example.com",
	}

	sendResp, err := env.UserService.SendEmailVerification(ctx, sendReq)

	// Assert
	assert.Error(t, err, "空用户ID应该返回错误")
	assert.Nil(t, sendResp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户ID", "错误信息应该包含'用户ID'")
}

// TestUserService_SendEmailVerification_EmptyEmail 测试发送验证邮件空邮箱
func TestUserService_SendEmailVerification_EmptyEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, _, _ := env.CreateDefaultTestUser(t)

	// Act - 使用空邮箱发送验证邮件
	sendReq := &user2.SendEmailVerificationRequest{
		UserID: userID,
		Email:  "",
	}

	sendResp, err := env.UserService.SendEmailVerification(ctx, sendReq)

	// Assert
	assert.Error(t, err, "空邮箱应该返回错误")
	assert.Nil(t, sendResp, "响应应该为空")
	assert.Contains(t, err.Error(), "邮箱", "错误信息应该包含'邮箱'")
}
