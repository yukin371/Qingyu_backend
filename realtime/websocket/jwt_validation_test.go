package websocket

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/service/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJWTService 模拟JWT服务用于测试
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(ctx context.Context, userID string, roles []string) (string, error) {
	args := m.Called(ctx, userID, roles)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenClaims), args.Error(1)
}

func (m *MockJWTService) RefreshToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) RevokeToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockJWTService) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

// TestNotificationHubValidateToken 测试通知Hub的token验证
func TestNotificationHubValidateToken(t *testing.T) {
	mockJWT := new(MockJWTService)
	hub := NewWSHub(mockJWT)
	ctx := context.Background()

	t.Run("有效token", func(t *testing.T) {
		expectedClaims := &auth.TokenClaims{
			UserID: "user123",
			Roles:  []string{"reader"},
			Exp:    time.Now().Add(1 * time.Hour).Unix(),
			Iat:    time.Now().Unix(),
		}

		mockJWT.On("ValidateToken", ctx, "valid_token").Return(expectedClaims, nil).Once()

		userID, err := hub.validateToken(ctx, "valid_token")

		assert.NoError(t, err)
		assert.Equal(t, "user123", userID)
	})

	t.Run("无效token", func(t *testing.T) {
		mockJWT.On("ValidateToken", ctx, "invalid_token").Return(nil, assert.AnError).Once()

		userID, err := hub.validateToken(ctx, "invalid_token")

		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Contains(t, err.Error(), "token验证失败")
	})

	t.Run("JWT服务未初始化", func(t *testing.T) {
		hubWithoutJWT := NewWSHub(nil)

		userID, err := hubWithoutJWT.validateToken(ctx, "any_token")

		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Contains(t, err.Error(), "JWT服务未初始化")
	})
}

// TestMessagingHubValidateToken 测试消息Hub的token验证
func TestMessagingHubValidateToken(t *testing.T) {
	mockJWT := new(MockJWTService)
	hub := NewMessagingWSHub(mockJWT)
	ctx := context.Background()

	t.Run("有效token", func(t *testing.T) {
		expectedClaims := &auth.TokenClaims{
			UserID: "user456",
			Roles:  []string{"reader", "writer"},
			Exp:    time.Now().Add(1 * time.Hour).Unix(),
			Iat:    time.Now().Unix(),
		}

		mockJWT.On("ValidateToken", ctx, "valid_token").Return(expectedClaims, nil).Once()

		userID, err := hub.validateMessagingToken(ctx, "valid_token")

		assert.NoError(t, err)
		assert.Equal(t, "user456", userID)
	})

	t.Run("无效token", func(t *testing.T) {
		mockJWT.On("ValidateToken", ctx, "invalid_token").Return(nil, assert.AnError).Once()

		userID, err := hub.validateMessagingToken(ctx, "invalid_token")

		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Contains(t, err.Error(), "token验证失败")
	})

	t.Run("JWT服务未初始化", func(t *testing.T) {
		hubWithoutJWT := NewMessagingWSHub(nil)

		userID, err := hubWithoutJWT.validateMessagingToken(ctx, "any_token")

		assert.Error(t, err)
		assert.Empty(t, userID)
		assert.Contains(t, err.Error(), "JWT服务未初始化")
	})

	t.Run("过期token", func(t *testing.T) {
		expiredClaims := &auth.TokenClaims{
			UserID: "user789",
			Roles:  []string{"reader"},
			Exp:    time.Now().Add(-1 * time.Hour).Unix(), // 已过期
			Iat:    time.Now().Add(-2 * time.Hour).Unix(),
		}

		mockJWT.On("ValidateToken", ctx, "expired_token").Return(expiredClaims, assert.AnError).Once()

		userID, err := hub.validateMessagingToken(ctx, "expired_token")

		assert.Error(t, err)
		assert.Empty(t, userID)
	})
}
