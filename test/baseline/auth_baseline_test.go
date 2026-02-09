package baseline

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAuthService Mock认证服务
type MockAuthService struct {
	ValidateTokenFunc    func(ctx context.Context, token string) (string, error)
	CheckPermissionFunc  func(ctx context.Context, userID, permission string) (bool, error)
	CreateSessionFunc    func(ctx context.Context, userID string) (string, error)
	DestroySessionFunc    func(ctx context.Context, sessionID string) error
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(ctx, token)
	}
	return "user-123", nil
}

func (m *MockAuthService) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	if m.CheckPermissionFunc != nil {
		return m.CheckPermissionFunc(ctx, userID, permission)
	}
	return true, nil
}

func (m *MockAuthService) CreateSession(ctx context.Context, userID string) (string, error) {
	if m.CreateSessionFunc != nil {
		return m.CreateSessionFunc(ctx, userID)
	}
	return "session-123", nil
}

func (m *MockAuthService) DestroySession(ctx context.Context, sessionID string) error {
	if m.DestroySessionFunc != nil {
		return m.DestroySessionFunc(ctx, sessionID)
	}
	return nil
}

// AuthServicePort 认证服务端口接口（简化版，用于测试）
type AuthServicePort interface {
	ValidateToken(ctx context.Context, token string) (string, error)
	CheckPermission(ctx context.Context, userID, permission string) (bool, error)
}

// TestAuthServicePortInterface 测试认证服务端口接口
func TestAuthServicePortInterface(t *testing.T) {
	t.Run("接口应该包含核心认证方法", func(t *testing.T) {
		var port AuthServicePort = &MockAuthService{}

		// 测试ValidateToken
		userID, err := port.ValidateToken(context.Background(), "valid-token")
		assert.NoError(t, err)
		assert.Equal(t, "user-123", userID)

		// 测试CheckPermission
		hasPermission, err := port.CheckPermission(context.Background(), "user-123", "read:books")
		assert.NoError(t, err)
		assert.True(t, hasPermission)
	})

	t.Run("应该支持自定义行为", func(t *testing.T) {
		customValidateCalled := false
		mock := &MockAuthService{
			ValidateTokenFunc: func(ctx context.Context, token string) (string, error) {
				customValidateCalled = true
				return "custom-user", nil
			},
		}

		var port AuthServicePort = mock
		userID, _ := port.ValidateToken(context.Background(), "token")

		assert.True(t, customValidateCalled)
		assert.Equal(t, "custom-user", userID)
	})
}

// TestTokenOperations 测试令牌操作基线
func TestTokenOperations(t *testing.T) {
	t.Run("应该能够验证有效的令牌", func(t *testing.T) {
		mock := &MockAuthService{
			ValidateTokenFunc: func(ctx context.Context, token string) (string, error) {
				if token == "valid-token" {
					return "user-123", nil
				}
				return "", assert.AnError
			},
		}

		var port AuthServicePort = mock

		// 测试有效令牌
		userID, err := port.ValidateToken(context.Background(), "valid-token")
		assert.NoError(t, err)
		assert.Equal(t, "user-123", userID)

		// 测试无效令牌
		_, err = port.ValidateToken(context.Background(), "invalid-token")
		assert.Error(t, err)
	})

	t.Run("应该能够检查权限", func(t *testing.T) {
		mock := &MockAuthService{
			CheckPermissionFunc: func(ctx context.Context, userID, permission string) (bool, error) {
				// 只允许read权限
				if permission[:4] == "read" {
					return true, nil
				}
				return false, nil
			},
		}

		var port AuthServicePort = mock

		// 测试有权限
		hasRead, _ := port.CheckPermission(context.Background(), "user-123", "read:books")
		assert.True(t, hasRead)

		// 测试无权限
		hasWrite, _ := port.CheckPermission(context.Background(), "user-123", "write:books")
		assert.False(t, hasWrite)
	})
}

// BenchmarkTokenOperations 性能基线测试
func BenchmarkTokenOperations(b *testing.B) {
	mock := &MockAuthService{
		ValidateTokenFunc: func(ctx context.Context, token string) (string, error) {
			return "user-123", nil
		},
		CheckPermissionFunc: func(ctx context.Context, userID, permission string) (bool, error) {
			return true, nil
		},
	}

	var port AuthServicePort = mock
	ctx := context.Background()

	b.Run("ValidateToken", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = port.ValidateToken(ctx, "test-token")
		}
	})

	b.Run("CheckPermission", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = port.CheckPermission(ctx, "user-123", "read:books")
		}
	})
}
