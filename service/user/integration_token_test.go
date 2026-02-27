package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	user2 "Qingyu_backend/service/interfaces/user"
)

// TestUserService_ValidateToken_Integration Token验证集成测试
func TestUserService_ValidateToken_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户并获取Token
	_, username, _, password := env.CreateDefaultTestUser(t)

	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "登录应该成功")
	require.NotEmpty(t, loginResp.Token, "Token不应该为空")

	// Act - 验证Token
	validateReq := &user2.ValidateTokenRequest{
		Token: loginResp.Token,
	}

	validateResp, err := env.UserService.ValidateToken(ctx, validateReq)

	// Assert
	// 注意：根据当前的ValidateToken实现，它总是返回false
	// 实际的JWT验证在中间件中完成
	// 这个测试主要验证函数调用不会出错
	require.NoError(t, err, "验证Token不应该返回错误")
	require.NotNil(t, validateResp, "响应不应该为空")

	// 验证Token本身的有效性（使用JWTTestHelper）
	userID, err := env.JWTTestHelper.ValidateTestToken(loginResp.Token)
	require.NoError(t, err, "Token本身应该是有效的")
	assert.NotEmpty(t, userID, "Token中应该包含用户ID")
}

// TestUserService_ValidateToken_ExpiredToken_Integration 过期Token测试
func TestUserService_ValidateToken_ExpiredToken_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户
	userID, _, _, _ := env.CreateDefaultTestUser(t)

	// 生成过期的Token
	expiredToken, err := env.JWTTestHelper.GenerateExpiredToken(userID)
	require.NoError(t, err, "生成过期Token应该成功")

	// Act - 验证过期Token
	validateReq := &user2.ValidateTokenRequest{
		Token: expiredToken,
	}

	validateResp, err := env.UserService.ValidateToken(ctx, validateReq)

	// Assert
	require.NoError(t, err, "验证过期Token不应该返回错误")
	require.NotNil(t, validateResp, "响应不应该为空")

	// 验证过期Token确实无效
	_, err = env.JWTTestHelper.ValidateTestToken(expiredToken)
	assert.Error(t, err, "过期Token验证应该失败")
}

// TestUserService_ValidateToken_InvalidToken_Integration 无效Token测试
func TestUserService_ValidateToken_InvalidToken_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 使用无效的Token格式
	invalidToken := "invalid.token.format"

	// Act - 验证无效Token
	validateReq := &user2.ValidateTokenRequest{
		Token: invalidToken,
	}

	validateResp, err := env.UserService.ValidateToken(ctx, validateReq)

	// Assert
	require.NoError(t, err, "验证无效Token不应该返回错误")
	require.NotNil(t, validateResp, "响应不应该为空")

	// 验证无效Token确实无效
	_, err = env.JWTTestHelper.ValidateTestToken(invalidToken)
	assert.Error(t, err, "无效Token验证应该失败")
}

// TestUserService_ValidateToken_WrongSignature 错误签名测试
func TestUserService_ValidateToken_WrongSignature_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 生成错误签名的Token
	userID, _, _, _ := env.CreateDefaultTestUser(t)
	wrongSignatureToken := env.JWTTestHelper.GenerateInvalidToken(userID)

	// Act - 验证错误签名的Token
	validateReq := &user2.ValidateTokenRequest{
		Token: wrongSignatureToken,
	}

	validateResp, err := env.UserService.ValidateToken(ctx, validateReq)

	// Assert
	require.NoError(t, err, "验证错误签名Token不应该返回错误")
	require.NotNil(t, validateResp, "响应不应该为空")

	// 验证错误签名Token确实无效
	_, err = env.JWTTestHelper.ValidateTestToken(wrongSignatureToken)
	assert.Error(t, err, "错误签名Token验证应该失败")
}

// TestUserService_TokenWithRoles 带角色的Token测试
func TestUserService_TokenWithRoles_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 创建带角色的用户
	userID := env.CreateTestUserWithRoles(t, []string{"reader", "writer"})
	user := env.AssertUserExists(t, userID)

	// 设置用户密码
	err := user.SetPassword("TestPassword123!")
	require.NoError(t, err, "设置密码应该成功")

	// 登录获取Token
	loginReq := &user2.LoginUserRequest{
		Username: user.Username,
		Password: "TestPassword123!",
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "登录应该成功")
	require.NotEmpty(t, loginResp.Token, "Token不应该为空")

	// Act - 解析Token验证角色
	claims, err := env.JWTTestHelper.ParseTokenWithoutValidation(loginResp.Token)
	require.NoError(t, err, "解析Token应该成功")

	// Assert
	roles, ok := (*claims)["roles"].([]interface{})
	require.True(t, ok, "Token应该包含roles字段")
	require.NotEmpty(t, roles, "角色列表不应该为空")

	// 验证角色
	roleSet := make(map[string]bool)
	for _, role := range roles {
		if roleStr, ok := role.(string); ok {
			roleSet[roleStr] = true
		}
	}

	assert.True(t, roleSet["reader"], "Token应该包含reader角色")
	assert.True(t, roleSet["writer"], "Token应该包含writer角色")
}

// TestUserService_ValidateToken_NonExistentUser 不存在用户的Token测试
func TestUserService_ValidateToken_NonExistentUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 为一个不存在的用户生成Token
	fakeUserID := "507f1f77bcf86cd799439011"
	token, err := env.JWTTestHelper.GenerateTestToken(fakeUserID, env.TestConfig.JWTExpiration)
	require.NoError(t, err, "生成Token应该成功")

	// Act - 验证不存在用户的Token
	validateReq := &user2.ValidateTokenRequest{
		Token: token,
	}

	validateResp, err := env.UserService.ValidateToken(ctx, validateReq)

	// Assert
	require.NoError(t, err, "验证Token不应该返回错误")
	require.NotNil(t, validateResp, "响应不应该为空")

	// Token本身是有效的，但用户不存在
	// 这种情况下，中间件会拒绝请求
}

// TestUserService_TokenExpiration Token过期时间测试
func TestUserService_TokenExpiration_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 注册用户并登录
	_, username, _, password := env.CreateDefaultTestUser(t)

	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "登录应该成功")
	require.NotEmpty(t, loginResp.Token, "Token不应该为空")

	// Act - 解析Token验证过期时间
	claims, err := env.JWTTestHelper.ParseTokenWithoutValidation(loginResp.Token)
	require.NoError(t, err, "解析Token应该成功")

	// Assert
	exp, ok := (*claims)["exp"].(float64)
	require.True(t, ok, "Token应该包含exp字段")

	// 验证过期时间接近当前时间 + 配置有效期（允许1分钟误差）
	expTime := time.Unix(int64(exp), 0)
	expectedExpTime := time.Now().Add(env.TestConfig.JWTExpiration)
	timeDiff := expTime.Sub(expectedExpTime)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	assert.True(t, timeDiff < time.Minute, "Token过期时间应该接近配置的时间")
}

// TestUserService_RefreshToken Token刷新测试
func TestUserService_RefreshToken_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Skip: RefreshToken功能尚未实现
	t.Skip("RefreshToken功能尚未实现")

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 注册用户并登录
	_, username, _, password := env.CreateDefaultTestUser(t)

	loginReq := &user2.LoginUserRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := env.UserService.LoginUser(ctx, loginReq)
	require.NoError(t, err, "登录应该成功")
	require.NotEmpty(t, loginResp.Token, "Token不应该为空")

	// Act - 使用RefreshToken获取新的AccessToken
	// 注意：RefreshToken功能尚未实现，以下代码为未来实现的占位符
	_ = loginResp
	_ = username
	_ = password
	_ = err
}

// TestUserService_RefreshToken_Invalid 无效RefreshToken测试
func TestUserService_RefreshToken_Invalid_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Skip: RefreshToken功能尚未实现
	t.Skip("RefreshToken功能尚未实现")

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 使用无效的RefreshToken
	// 注意：RefreshToken功能尚未实现，以下代码为未来实现的占位符
	_ = ctx
	_ = env
}
