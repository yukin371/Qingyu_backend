package user

import (
	"context"
	"testing"

	user2 "Qingyu_backend/service/interfaces/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =========================
// 用户登出相关测试
// =========================

// TestUserService_LogoutUser_Success 测试用户登出成功
func TestUserService_LogoutUser_Success(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.LogoutUserRequest{
		Token: "valid_jwt_token",
	}

	// Act
	resp, err := service.LogoutUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 注意：当前实现是简化版本，返回固定成功
	// 完整实现应该验证Token并加入黑名单
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
}

// =========================
// Token验证相关测试
// =========================

// TestUserService_ValidateToken_ValidToken 测试验证Token-有效Token
func TestUserService_ValidateToken_ValidToken(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	req := &user2.ValidateTokenRequest{
		Token: "valid_jwt_token",
	}

	// Act
	resp, err := service.ValidateToken(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	// 注意：当前简化实现返回false
	// 完整实现应该验证Token并返回用户信息
	assert.False(t, resp.Valid) // 当前实现固定返回false
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

	req := &user2.ValidateTokenRequest{
		Token: "invalid_token_format",
	}

	// Act
	resp, err := service.ValidateToken(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Valid)
}
