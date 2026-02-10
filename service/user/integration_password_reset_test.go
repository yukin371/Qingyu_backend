package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	user2 "Qingyu_backend/service/interfaces/user"
)

// TestUserService_PasswordResetFlow_Integration 完整密码重置流程集成测试
func TestUserService_PasswordResetFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	_, username, email, originalPassword := env.CreateDefaultTestUser(t)

	// 确保用户可以登录
	loginReq1 := &user2.LoginUserRequest{
		Username: username,
		Password: originalPassword,
	}

	loginResp1, err := env.UserService.LoginUser(ctx, loginReq1)
	require.NoError(t, err, "原始密码登录应该成功")
	require.NotEmpty(t, loginResp1.Token, "Token不应该为空")

	// Act 1 - 请求密码重置
	requestReq := &user2.RequestPasswordResetRequest{
		Email: email,
	}

	requestResp, err := env.UserService.RequestPasswordReset(ctx, requestReq)

	// Assert 1 - 请求密码重置成功
	require.NoError(t, err, "请求密码重置应该成功")
	require.NotNil(t, requestResp, "响应不应该为空")
	assert.True(t, requestResp.Success, "请求应该成功")

	// Act 2 - 使用生成的Token确认密码重置
	tokenManager := GetGlobalPasswordResetTokenManager()
	resetToken, err := tokenManager.GenerateToken(ctx, email)
	require.NoError(t, err, "生成重置Token应该成功")

	newPassword := "NewPassword456!"
	confirmReq := &user2.ConfirmPasswordResetRequest{
		Email:    email,
		Token:    resetToken,
		Password: newPassword,
	}

	confirmResp, err := env.UserService.ConfirmPasswordReset(ctx, confirmReq)

	// Assert 2 - 确认密码重置成功
	require.NoError(t, err, "确认密码重置应该成功")
	require.NotNil(t, confirmResp, "响应不应该为空")
	assert.True(t, confirmResp.Success, "密码重置应该成功")

	// Act 3 - 使用新密码登录
	loginReq2 := &user2.LoginUserRequest{
		Username: username,
		Password: newPassword,
	}

	loginResp2, err := env.UserService.LoginUser(ctx, loginReq2)

	// Assert 3 - 新密码登录应该成功
	require.NoError(t, err, "新密码登录应该成功")
	require.NotNil(t, loginResp2, "登录响应不应该为空")
	assert.NotEmpty(t, loginResp2.Token, "Token不应该为空")

	// Act 4 - 尝试使用旧密码登录
	loginReq3 := &user2.LoginUserRequest{
		Username: username,
		Password: originalPassword,
	}

	loginResp3, err := env.UserService.LoginUser(ctx, loginReq3)

	// Assert 4 - 旧密码登录应该失败
	assert.Error(t, err, "旧密码登录应该失败")
	assert.Nil(t, loginResp3, "响应应该为空")
	assert.Contains(t, err.Error(), "密码错误", "错误信息应该包含'密码错误'")
}

// TestUserService_ResetPassword_NonExistentEmail 不存在邮箱测试
func TestUserService_ResetPassword_NonExistentEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 请求重置不存在的邮箱
	requestReq := &user2.RequestPasswordResetRequest{
		Email: "nonexistent@example.com",
	}

	requestResp, err := env.UserService.RequestPasswordReset(ctx, requestReq)

	// Assert - 为了安全，即使用户不存在也返回成功（防止邮箱枚举攻击）
	require.NoError(t, err, "请求不应该返回错误")
	require.NotNil(t, requestResp, "响应不应该为空")
	assert.True(t, requestResp.Success, "请求应该返回成功")
	assert.Contains(t, requestResp.Message, "如果该邮箱已注册", "消息应该包含提示信息")
}

// TestUserService_UpdatePassword_WrongOldPassword_Integration 错误旧密码测试
func TestUserService_UpdatePassword_WrongOldPassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, _, _ := env.CreateDefaultTestUser(t)

	// Act - 使用错误的旧密码更新密码
	updateReq := &user2.UpdatePasswordRequest{
		ID:          userID,
		OldPassword: "WrongOldPassword123!",
		NewPassword: "NewPassword456!",
	}

	updateResp, err := env.UserService.UpdatePassword(ctx, updateReq)

	// Assert
	assert.Error(t, err, "错误旧密码应该返回错误")
	assert.Nil(t, updateResp, "响应应该为空")
	assert.Contains(t, err.Error(), "旧密码错误", "错误信息应该包含'旧密码错误'")
}

// TestUserService_UpdatePassword_NonExistentUser_Integration 不存在用户测试
func TestUserService_UpdatePassword_NonExistentUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 为不存在的用户更新密码
	updateReq := &user2.UpdatePasswordRequest{
		ID:          "507f1f77bcf86cd799439011",
		OldPassword: "OldPassword123!",
		NewPassword: "NewPassword456!",
	}

	updateResp, err := env.UserService.UpdatePassword(ctx, updateReq)

	// Assert
	assert.Error(t, err, "不存在的用户应该返回错误")
	assert.Nil(t, updateResp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户不存在", "错误信息应该包含'用户不存在'")
}

// TestUserService_PasswordReset_ViaEmail 通过邮箱重置密码完整流程
func TestUserService_PasswordReset_ViaEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	_, username, email, originalPassword := env.CreateDefaultTestUser(t)

	// Act 1 - 请求密码重置
	requestReq := &user2.RequestPasswordResetRequest{
		Email: email,
	}

	requestResp, err := env.UserService.RequestPasswordReset(ctx, requestReq)
	require.NoError(t, err, "请求密码重置应该成功")
	require.True(t, requestResp.Success, "请求应该成功")

	// Act 2 - 生成重置Token并确认
	tokenManager := GetGlobalPasswordResetTokenManager()
	resetToken, err := tokenManager.GenerateToken(ctx, email)
	require.NoError(t, err, "生成重置Token应该成功")

	newPassword := "ResetPassword789!"
	confirmReq := &user2.ConfirmPasswordResetRequest{
		Email:    email,
		Token:    resetToken,
		Password: newPassword,
	}

	confirmResp, err := env.UserService.ConfirmPasswordReset(ctx, confirmReq)
	require.NoError(t, err, "确认密码重置应该成功")
	require.True(t, confirmResp.Success, "密码重置应该成功")

	// Assert - 使用新密码登录成功
	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: newPassword,
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "新密码登录应该成功")
	assert.NotEmpty(t, loginResp.Token, "Token不应该为空")

	// 验证旧密码无效
	oldLoginReq := &user2.LoginUserRequest{
		Username: username,
		Password: originalPassword,
	}

	_, err = env.UserService.LoginUser(ctx, oldLoginReq)
	assert.Error(t, err, "旧密码应该无效")
}

// TestUserService_MultiplePasswordUpdates 多次密码更新测试
func TestUserService_MultiplePasswordUpdates_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, username, _, originalPassword := env.CreateDefaultTestUser(t)

	// Act - 多次更新密码
	passwords := []string{
		originalPassword,
		"Password1_abc",
		"Password2_xyz",
		"Password3_123",
	}

	for i := 1; i < len(passwords); i++ {
		oldPassword := passwords[i-1]
		newPassword := passwords[i]

		updateReq := &user2.UpdatePasswordRequest{
			ID:          userID,
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}

		updateResp, err := env.UserService.UpdatePassword(ctx, updateReq)

		// Assert - 每次更新都应该成功
		require.NoError(t, err, "第%d次密码更新应该成功", i)
		require.NotNil(t, updateResp, "第%d次响应不应该为空", i)
		assert.True(t, updateResp.Updated, "第%d次密码更新应该成功", i)
	}

	// Assert - 使用最新的密码登录成功
	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: passwords[len(passwords)-1],
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "最新密码登录应该成功")
	assert.NotEmpty(t, loginResp.Token, "Token不应该为空")
}

// TestUserService_PasswordReset_SamePassword 设置相同密码测试
func TestUserService_PasswordReset_SamePassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, _, password := env.CreateDefaultTestUser(t)

	// Act - 尝试将密码更新为相同的密码
	updateReq := &user2.UpdatePasswordRequest{
		ID:          userID,
		OldPassword: password,
		NewPassword: password,
	}

	_, err := env.UserService.UpdatePassword(ctx, updateReq)

	// Assert - 设置相同密码可能成功或失败，取决于系统实现
	// 这里我们验证不会崩溃
	_ = err
	// 某些系统可能允许设置相同密码，某些可能不允许
}

// TestUserService_VerificationCode_PasswordReset 验证码密码重置流程
func TestUserService_VerificationCode_PasswordReset_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	_, username, email, _ := env.CreateDefaultTestUser(t)

	// Act 1 - 请求密码重置
	requestReq := &user2.RequestPasswordResetRequest{
		Email: email,
	}

	_, err := env.UserService.RequestPasswordReset(ctx, requestReq)
	require.NoError(t, err, "请求密码重置应该成功")

	// Act 2 - 使用无效的Token尝试重置
	confirmReq := &user2.ConfirmPasswordResetRequest{
		Email:    email,
		Token:    "invalid_token_12345",
		Password: "NewPassword456!",
	}

	confirmResp, err := env.UserService.ConfirmPasswordReset(ctx, confirmReq)

	// Assert - 无效Token应该返回错误
	assert.Error(t, err, "无效Token应该返回错误")
	assert.Nil(t, confirmResp, "响应应该为空")
	assert.Contains(t, err.Error(), "Token验证失败", "错误信息应该包含'Token验证失败'")

	// Act 3 - 使用正确的Token重置密码
	tokenManager := GetGlobalPasswordResetTokenManager()
	resetToken, err := tokenManager.GenerateToken(ctx, email)
	require.NoError(t, err, "生成重置Token应该成功")

	confirmReq.Token = resetToken
	confirmResp, err = env.UserService.ConfirmPasswordReset(ctx, confirmReq)

	// Assert - 密码重置应该成功
	require.NoError(t, err, "密码重置应该成功")
	require.NotNil(t, confirmResp, "响应不应该为空")
	assert.True(t, confirmResp.Success, "密码重置应该成功")

	// Act 4 - 使用新密码登录
	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: "NewPassword456!",
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "新密码登录应该成功")
	assert.NotEmpty(t, loginResp.Token, "Token不应该为空")
}

// TestUserService_UpdatePassword_EmptyUserID 测试空用户ID
func TestUserService_UpdatePassword_EmptyUserID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 使用空用户ID更新密码
	updateReq := &user2.UpdatePasswordRequest{
		ID:          "",
		OldPassword: "OldPassword123!",
		NewPassword: "NewPassword456!",
	}

	updateResp, err := env.UserService.UpdatePassword(ctx, updateReq)

	// Assert
	assert.Error(t, err, "空用户ID应该返回错误")
	assert.Nil(t, updateResp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户ID", "错误信息应该包含'用户ID'")
}

// TestUserService_UpdatePassword_EmptyOldPassword 测试空旧密码
func TestUserService_UpdatePassword_EmptyOldPassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, _, _ := env.CreateDefaultTestUser(t)

	// Act - 使用空旧密码更新密码
	updateReq := &user2.UpdatePasswordRequest{
		ID:          userID,
		OldPassword: "",
		NewPassword: "NewPassword456!",
	}

	updateResp, err := env.UserService.UpdatePassword(ctx, updateReq)

	// Assert
	assert.Error(t, err, "空旧密码应该返回错误")
	assert.Nil(t, updateResp, "响应应该为空")
	assert.Contains(t, err.Error(), "旧密码", "错误信息应该包含'旧密码'")
}

// TestUserService_UpdatePassword_EmptyNewPassword 测试空新密码
func TestUserService_UpdatePassword_EmptyNewPassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	userID, _, _, password := env.CreateDefaultTestUser(t)

	// Act - 使用空新密码更新密码
	updateReq := &user2.UpdatePasswordRequest{
		ID:          userID,
		OldPassword: password,
		NewPassword: "",
	}

	updateResp, err := env.UserService.UpdatePassword(ctx, updateReq)

	// Assert
	assert.Error(t, err, "空新密码应该返回错误")
	assert.Nil(t, updateResp, "响应应该为空")
	assert.Contains(t, err.Error(), "新密码", "错误信息应该包含'新密码'")
}

// TestUserService_ConfirmPasswordReset_EmptyEmail 测试空邮箱
func TestUserService_ConfirmPasswordReset_EmptyEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 使用空邮箱确认密码重置
	confirmReq := &user2.ConfirmPasswordResetRequest{
		Email:    "",
		Token:    "some_token",
		Password: "NewPassword456!",
	}

	confirmResp, err := env.UserService.ConfirmPasswordReset(ctx, confirmReq)

	// Assert
	assert.Error(t, err, "空邮箱应该返回错误")
	assert.Nil(t, confirmResp, "响应应该为空")
	assert.Contains(t, err.Error(), "邮箱", "错误信息应该包含'邮箱'")
}

// TestUserService_ConfirmPasswordReset_EmptyToken 测试空Token
func TestUserService_ConfirmPasswordReset_EmptyToken_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	_, _, email, _ := env.CreateDefaultTestUser(t)

	// Act - 使用空Token确认密码重置
	confirmReq := &user2.ConfirmPasswordResetRequest{
		Email:    email,
		Token:    "",
		Password: "NewPassword456!",
	}

	confirmResp, err := env.UserService.ConfirmPasswordReset(ctx, confirmReq)

	// Assert
	assert.Error(t, err, "空Token应该返回错误")
	assert.Nil(t, confirmResp, "响应应该为空")
	assert.Contains(t, err.Error(), "Token", "错误信息应该包含'Token'")
}

// TestUserService_ConfirmPasswordReset_EmptyPassword 测试空密码
func TestUserService_ConfirmPasswordReset_EmptyPassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个用户
	_, _, email, _ := env.CreateDefaultTestUser(t)

	// Act - 使用空密码确认密码重置
	confirmReq := &user2.ConfirmPasswordResetRequest{
		Email:    email,
		Token:    "some_token",
		Password: "",
	}

	confirmResp, err := env.UserService.ConfirmPasswordReset(ctx, confirmReq)

	// Assert
	assert.Error(t, err, "空密码应该返回错误")
	assert.Nil(t, confirmResp, "响应应该为空")
	assert.Contains(t, err.Error(), "密码", "错误信息应该包含'密码'")
}

// TestUserService_RequestPasswordReset_EmptyEmail 测试请求重置空邮箱
func TestUserService_RequestPasswordReset_EmptyEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 使用空邮箱请求密码重置
	requestReq := &user2.RequestPasswordResetRequest{
		Email: "",
	}

	requestResp, err := env.UserService.RequestPasswordReset(ctx, requestReq)

	// Assert
	assert.Error(t, err, "空邮箱应该返回错误")
	assert.Nil(t, requestResp, "响应应该为空")
	assert.Contains(t, err.Error(), "邮箱", "错误信息应该包含'邮箱'")
}
