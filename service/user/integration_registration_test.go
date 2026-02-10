package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	user2 "Qingyu_backend/service/interfaces/user"
)

// TestUserService_RegisterUser_Integration 完整用户注册集成测试
func TestUserService_RegisterUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange
	testUser := env.GenerateUniqueTestUser(t)

	req := &user2.RegisterUserRequest{
		Username: testUser.Username,
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	// Act
	resp, err := env.UserService.RegisterUser(ctx, req)

	// Assert
	require.NoError(t, err, "注册用户应该成功")
	require.NotNil(t, resp, "响应不应该为空")
	require.NotNil(t, resp.User, "用户信息不应该为空")

	// 验证用户基本信息
	assert.Equal(t, testUser.Username, resp.User.Username, "用户名应该匹配")
	assert.Equal(t, testUser.Email, resp.User.Email, "邮箱应该匹配")
	assert.NotEmpty(t, resp.User.ID, "用户ID不应该为空")
	assert.NotEmpty(t, resp.Token, "Token不应该为空")

	// 验证默认值
	assert.Equal(t, "reader", resp.User.Roles[0], "默认角色应该是reader")
	assert.Equal(t, "active", string(resp.User.Status), "默认状态应该是active")
	assert.False(t, resp.User.EmailVerified, "邮箱验证状态应该是false")

	// 验证Token有效性
	userID, err := env.JWTTestHelper.ValidateTestToken(resp.Token)
	require.NoError(t, err, "Token验证应该成功")
	assert.Equal(t, resp.User.ID, userID, "Token中的用户ID应该匹配")

	// 验证用户已在数据库中创建
	dbUser := env.AssertUserExists(t, resp.User.ID)
	assert.Equal(t, testUser.Username, dbUser.Username, "数据库中的用户名应该匹配")
	assert.Equal(t, testUser.Email, dbUser.Email, "数据库中的邮箱应该匹配")
	assert.NotEqual(t, testUser.Password, dbUser.Password, "密码应该被加密")
	assert.True(t, dbUser.CheckPassword(testUser.Password), "密码验证应该成功")
}

// TestUserService_RegisterUser_DuplicateUsername_Integration 重复用户名测试
func TestUserService_RegisterUser_DuplicateUsername_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户
	testUser := env.GenerateUniqueTestUser(t)
	req1 := &user2.RegisterUserRequest{
		Username: testUser.Username,
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	_, err := env.UserService.RegisterUser(ctx, req1)
	require.NoError(t, err, "第一次注册应该成功")

	// Act - 尝试用相同用户名注册另一个用户
	duplicateUser := env.GenerateUniqueTestUserWithPrefix(t, "duplicate")
	duplicateUser.Username = testUser.Username // 使用相同的用户名

	req2 := &user2.RegisterUserRequest{
		Username: duplicateUser.Username,
		Email:    duplicateUser.Email,
		Password: duplicateUser.Password,
	}

	// Assert
	resp, err := env.UserService.RegisterUser(ctx, req2)
	assert.Error(t, err, "重复用户名注册应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "用户名已存在", "错误信息应该包含'用户名已存在'")
}

// TestUserService_RegisterUser_DuplicateEmail_Integration 重复邮箱测试
func TestUserService_RegisterUser_DuplicateEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange - 先注册一个用户
	testUser := env.GenerateUniqueTestUser(t)
	req1 := &user2.RegisterUserRequest{
		Username: testUser.Username,
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	_, err := env.UserService.RegisterUser(ctx, req1)
	require.NoError(t, err, "第一次注册应该成功")

	// Act - 尝试用相同邮箱注册另一个用户
	duplicateUser := env.GenerateUniqueTestUserWithPrefix(t, "duplicate")
	duplicateUser.Email = testUser.Email // 使用相同的邮箱

	req2 := &user2.RegisterUserRequest{
		Username: duplicateUser.Username,
		Email:    duplicateUser.Email,
		Password: duplicateUser.Password,
	}

	// Assert
	resp, err := env.UserService.RegisterUser(ctx, req2)
	assert.Error(t, err, "重复邮箱注册应该失败")
	assert.Nil(t, resp, "响应应该为空")
	assert.Contains(t, err.Error(), "邮箱已存在", "错误信息应该包含'邮箱已存在'")
}

// TestUserService_RegisterUser_TokenExpiration 测试Token过期时间
func TestUserService_RegisterUser_TokenExpiration_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange
	testUser := env.GenerateUniqueTestUser(t)
	req := &user2.RegisterUserRequest{
		Username: testUser.Username,
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	// Act
	resp, err := env.UserService.RegisterUser(ctx, req)

	// Assert
	require.NoError(t, err, "注册用户应该成功")
	require.NotEmpty(t, resp.Token, "Token不应该为空")

	// 验证Token包含正确的过期时间
	claims, err := env.JWTTestHelper.ParseTokenWithoutValidation(resp.Token)
	require.NoError(t, err, "解析Token应该成功")

	exp, ok := (*claims)["exp"].(float64)
	require.True(t, ok, "Token应该包含exp字段")
	expTime := time.Unix(int64(exp), 0)

	// Token应该在24小时后过期
	expectedExp := time.Now().Add(env.TestConfig.JWTExpiration)
	timeDiff := expTime.Sub(expectedExp)

	assert.True(t, timeDiff < time.Minute, "Token过期时间应该接近24小时")
	assert.True(t, expTime.After(time.Now()), "Token应该在未来过期")
}

// TestUserService_RegisterUser_DefaultValues 测试默认值设置
func TestUserService_RegisterUser_DefaultValues_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange
	testUser := env.GenerateUniqueTestUser(t)
	req := &user2.RegisterUserRequest{
		Username: testUser.Username,
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	// Act
	resp, err := env.UserService.RegisterUser(ctx, req)

	// Assert
	require.NoError(t, err, "注册用户应该成功")

	// 验证默认角色
	assert.NotEmpty(t, resp.User.Roles, "角色列表不应该为空")
	assert.Contains(t, resp.User.Roles, "reader", "应该包含reader角色")

	// 验证默认状态
	assert.Equal(t, "active", string(resp.User.Status), "默认状态应该是active")

	// 验证邮箱验证状态
	assert.False(t, resp.User.EmailVerified, "邮箱验证状态应该是false")
	assert.Nil(t, resp.User.EmailVerifiedAt, "邮箱验证时间应该为空")

	// 验证手机验证状态
	assert.False(t, resp.User.PhoneVerified, "手机验证状态应该是false")

	// 验证时间戳
	assert.NotZero(t, resp.User.CreatedAt, "创建时间不应该为空")
	assert.NotZero(t, resp.User.UpdatedAt, "更新时间不应该为空")
}

// TestUserService_RegisterUser_InvalidEmail 测试无效邮箱
func TestUserService_RegisterUser_InvalidEmail_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange
	testUser := env.GenerateUniqueTestUser(t)

	testCases := []struct {
		name  string
		email string
	}{
		{"无效邮箱格式-无@", "invalidemail"},
		{"无效邮箱格式-无域名", "invalid@"},
		{"无效邮箱格式-无用户名", "@example.com"},
		{"无效邮箱格式-无点", "user@examplecom"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &user2.RegisterUserRequest{
				Username: testUser.Username,
				Email:    tc.email,
				Password: testUser.Password,
			}

			// Act
			resp, err := env.UserService.RegisterUser(ctx, req)

			// Assert
			assert.Error(t, err, "无效邮箱应该返回错误")
			assert.Nil(t, resp, "响应应该为空")
			assert.Contains(t, err.Error(), "邮箱", "错误信息应该包含'邮箱'")
		})
	}
}

// TestUserService_RegisterUser_WeakPassword 测试弱密码
func TestUserService_RegisterUser_WeakPassword_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange
	testUser := env.GenerateUniqueTestUser(t)

	testCases := []struct {
		name     string
		password string
	}{
		{"密码太短", "12345"},
		{"密码为空", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &user2.RegisterUserRequest{
				Username: testUser.Username,
				Email:    testUser.Email,
				Password: tc.password,
			}

			// Act
			resp, err := env.UserService.RegisterUser(ctx, req)

			// Assert
			assert.Error(t, err, "弱密码应该返回错误")
			assert.Nil(t, resp, "响应应该为空")
		})
	}
}

// TestUserService_RegisterUser_InvalidUsername 测试无效用户名
func TestUserService_RegisterUser_InvalidUsername_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// Setup
	env := SetupIntegrationTestEnvironment(t)
	defer env.CleanupFunc()

	ctx := context.Background()

	// Arrange
	testUser := env.GenerateUniqueTestUser(t)

	testCases := []struct {
		name     string
		username string
	}{
		{"用户名太短", "ab"},
		{"用户名为空", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &user2.RegisterUserRequest{
				Username: tc.username,
				Email:    testUser.Email,
				Password: testUser.Password,
			}

			// Act
			resp, err := env.UserService.RegisterUser(ctx, req)

			// Assert
			assert.Error(t, err, "无效用户名应该返回错误")
			assert.Nil(t, resp, "响应应该为空")
		})
	}
}
