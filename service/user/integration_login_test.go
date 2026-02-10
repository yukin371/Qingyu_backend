package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	user2 "Qingyu_backend/service/interfaces/user"
	usersModel "Qingyu_backend/models/users"
)

// TestUserService_LoginUser_Integration 完整用户登录集成测试
func TestUserService_LoginUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户
	userID, username, email, password := env.CreateDefaultTestUser(t)

	// Act - 使用正确的用户名和密码登录
	req := &user2.LoginUserRequest{
		Username: username,
		Password: password,
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	require.NoError(t, err, "登录应该成功")
	require.NotNil(t, resp, "响应不应该为空")
	require.NotNil(t, resp.User, "用户信息不应该为空")
	assert.NotEmpty(t, resp.Token, "Token不应该为空")
	assert.NotEmpty(t, resp.RefreshToken, "RefreshToken不应该为空")

	// 验证用户信息
	assert.Equal(t, userID, resp.User.ID, "用户ID应该匹配")
	assert.Equal(t, username, resp.User.Username, "用户名应该匹配")
	assert.Equal(t, email, resp.User.Email, "邮箱应该匹配")

	// 验证Token有效性
	tokenUserID, err := env.JWTTestHelper.ValidateTestToken(resp.Token)
	require.NoError(t, err, "Token验证应该成功")
	assert.Equal(t, userID, tokenUserID, "Token中的用户ID应该匹配")
}

// TestUserService_LoginUser_WrongPassword_Integration 错误密码测试
func TestUserService_LoginUser_WrongPassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户
	_, username, _, _ := env.CreateDefaultTestUser(t)

	// Act - 使用错误的密码登录
	req := &user2.LoginUserRequest{
		Username: username,
		Password: "WrongPassword123!",
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	assert.Error(t, err, "错误密码登录应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "密码错误", "错误信息应该包含'密码错误'")
}

// TestUserService_LoginUser_UserNotFound_Integration 用户不存在测试
func TestUserService_LoginUser_UserNotFound_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 尝试登录不存在的用户
	req := &user2.LoginUserRequest{
		Username: "nonexistentuser",
		Password: "SomePassword123!",
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	assert.Error(t, err, "不存在的用户登录应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户不存在", "错误信息应该包含'用户不存在'")
}

// TestUserService_LoginUser_InactiveAccount 未激活账号测试
func TestUserService_LoginUser_InactiveAccount_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个未激活的用户
	userID := env.CreateTestUserWithStatus(t, usersModel.UserStatusInactive)
	user := env.AssertUserExists(t, userID)

	// Act - 尝试登录未激活的用户
	req := &user2.LoginUserRequest{
		Username: user.Username,
		Password: "TestPassword123!",
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	assert.Error(t, err, "未激活账号登录应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "未激活", "错误信息应该包含'未激活'")
}

// TestUserService_LoginUser_BannedAccount 被封禁账号测试
func TestUserService_LoginUser_BannedAccount_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建一个被封禁的用户
	userID := env.CreateTestUserWithStatus(t, usersModel.UserStatusBanned)
	user := env.AssertUserExists(t, userID)

	// Act - 尝试登录被封禁的用户
	req := &user2.LoginUserRequest{
		Username: user.Username,
		Password: "TestPassword123!",
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	assert.Error(t, err, "被封禁账号登录应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "封禁", "错误信息应该包含'封禁'")
}

// TestUserService_LoginAndLogout 登录登出流程测试
func TestUserService_LoginAndLogout_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册并登录一个用户
	_, username, _, password := env.CreateDefaultTestUser(t)

	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "登录应该成功")
	require.NotEmpty(t, loginResp.Token, "Token不应该为空")

	// Act - 登出用户
	logoutReq := &user2.LogoutUserRequest{
		UserID: loginResp.User.ID,
		Token:  loginResp.Token,
	}

	logoutResp, err := env.UserService.LogoutUser(ctx, logoutReq)

	// Assert
	require.NoError(t, err, "登出应该成功")
	require.NotNil(t, logoutResp, "响应不应该为空")
	assert.True(t, logoutResp.Success, "登出应该成功")

	// 验证Token已失效（可选，取决于系统实现）
	// 有些系统可能会将Token加入黑名单，有些则不会
}

// TestUserService_MultipleLogins 多次登录测试
func TestUserService_MultipleLogins_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户
	_, username, _, password := env.CreateDefaultTestUser(t)

	// Act - 多次登录同一个用户
	var tokens []string
	for i := 0; i < 3; i++ {
		req := &user2.LoginUserRequest{
			Username: username,
			Password: password,
		}

		resp, err := env.UserService.LoginUser(ctx, req)

		// Assert
		require.NoError(t, err, "第%d次登录应该成功", i+1)
		require.NotNil(t, resp, "第%d次登录响应不应该为空", i+1)
		assert.NotEmpty(t, resp.Token, "第%d次登录的Token不应该为空", i+1)

		// 每次登录应该生成不同的Token
		for j, existingToken := range tokens {
			assert.NotEqual(t, existingToken, resp.Token, "第%d次登录的Token应该与第%d次不同", i+1, j+1)
		}

		tokens = append(tokens, resp.Token)
	}

	// 所有Token都应该有效
	for i, token := range tokens {
		userID, err := env.JWTTestHelper.ValidateTestToken(token)
		require.NoError(t, err, "第%d个Token验证应该成功", i+1)
		assert.NotEmpty(t, userID, "第%d个Token中的用户ID不应该为空", i+1)
	}
}

// TestUserService_LoginUser_WithClientIP 测试带客户端IP的登录
func TestUserService_LoginUser_WithClientIP_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户
	_, username, _, password := env.CreateDefaultTestUser(t)

	// Act - 使用客户端IP登录
	req := &user2.LoginUserRequest{
		Username: username,
		Password: password,
		ClientIP: "192.168.1.100",
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	require.NoError(t, err, "带客户端IP的登录应该成功")
	require.NotNil(t, resp, "响应不应该为空")
	assert.NotEmpty(t, resp.Token, "Token不应该为空")

	// 验证用户的最后登录IP已更新（可选，取决于系统实现）
}

// TestUserService_LoginUser_EmptyUsername 测试空用户名
func TestUserService_LoginUser_EmptyUsername_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 使用空用户名登录
	req := &user2.LoginUserRequest{
		Username: "",
		Password: "SomePassword123!",
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	assert.Error(t, err, "空用户名登录应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户名", "错误信息应该包含'用户名'")
}

// TestUserService_LoginUser_EmptyPassword 测试空密码
func TestUserService_LoginUser_EmptyPassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Act - 使用空密码登录
	req := &user2.LoginUserRequest{
		Username: "testuser",
		Password: "",
	}

	resp, err := env.UserService.LoginUser(ctx, req)

	// Assert
	assert.Error(t, err, "空密码登录应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "密码", "错误信息应该包含'密码'")
}
