package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"Qingyu_backend/config"
)

// MockRedisClient Mock Redis客户端
type MockRedisClient struct {
	data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.data[key] = value.(string)
	return nil
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	count := int64(0)
	for _, key := range keys {
		if _, ok := m.data[key]; ok {
			count++
		}
	}
	return count, nil
}

// TestJWTService_GenerateToken_TableDriven 表格驱动测试 - 生成Token
func TestJWTService_GenerateToken_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		roles     []string
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "正常生成Token",
			userID:  "user_123",
			roles:   []string{"reader", "author"},
			wantErr: false,
		},
		{
			name:    "生成管理员Token",
			userID:  "admin_001",
			roles:   []string{"admin", "editor", "reader"},
			wantErr: false,
		},
		{
			name:    "无角色生成Token",
			userID:  "user_456",
			roles:   []string{},
			wantErr: false,
		},
		{
			name:    "空用户ID",
			userID:  "",
			roles:   []string{"reader"},
			wantErr: true,
			errMsg:  "不能为空",
		},
		{
			name:    "多角色Token",
			userID:  "super_user",
			roles:   []string{"reader", "author", "editor", "admin", "reviewer"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cache := NewMockRedisClient()
			cfg := &config.JWTConfigEnhanced{
				SecretKey:       "test-secret-key-12345678",
				Issuer:          "qingyu-test",
				Expiration:      1 * time.Hour,
				RefreshDuration: 24 * time.Hour,
			}
			service := NewJWTService(cfg, cache)
			ctx := context.Background()

			// Act
			token, err := service.GenerateToken(ctx, tt.userID, tt.roles)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// 验证Token可以解析
				claims, err := service.ValidateToken(ctx, token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.roles, claims.Roles)
			}
		})
	}
}

// TestJWTService_ValidateToken_TableDriven 表格驱动测试 - 验证Token
func TestJWTService_ValidateToken_TableDriven(t *testing.T) {
	cache := NewMockRedisClient()
	cfg := &config.JWTConfigEnhanced{
		Secret:      "test-secret-key-12345678",
		Expiration:  1 * time.Hour,
	}
	service := NewJWTService(cfg, cache)
	ctx := context.Background()

	// 预生成有效Token
	validToken, _ := service.GenerateToken(ctx, "user_123", []string{"reader"})

	tests := []struct {
		name      string
		token     string
		setup     func()
		wantErr   bool
		errMsg    string
		checkErr  string
	}{
		{
			name:    "有效Token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "空Token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "格式错误的Token",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:  "已吊销的Token",
			token: validToken,
			setup: func() {
				service.RevokeToken(ctx, validToken)
			},
			wantErr:  true,
		},
		{
			name:    "伪造的Token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fake.signature",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			if tt.setup != nil {
				tt.setup()
			}

			// Act
			claims, err := service.ValidateToken(ctx, tt.token)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.NotEmpty(t, claims.UserID)
			}
		})
	}
}

// TestJWTService_TokenExpiration 表格驱动测试 - Token过期
func TestJWTService_TokenExpiration(t *testing.T) {
	tests := []struct {
		name         string
		expiration   time.Duration
		validateWait time.Duration
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "Token未过期",
			expiration:   1 * time.Hour,
			validateWait: 0,
			wantErr:      false,
		},
		{
			name:         "Token已过期",
			expiration:   10 * time.Millisecond,
			validateWait: 20 * time.Millisecond,
			wantErr:      true,
		},
		{
			name:         "即将过期但仍然有效",
			expiration:   100 * time.Millisecond,
			validateWait: 50 * time.Millisecond,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cache := NewMockRedisClient()
			cfg := &config.JWTConfigEnhanced{
				Secret:      "test-secret-key-12345678",
				Expiration:  tt.expiration,
			}
			service := NewJWTService(cfg, cache)
			ctx := context.Background()

			token, err := service.GenerateToken(ctx, "user_123", []string{"reader"})
			require.NoError(t, err)

			if tt.validateWait > 0 {
				time.Sleep(tt.validateWait)
			}

			// Act
			claims, err := service.ValidateToken(ctx, token)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

// TestJWTService_RefreshToken 表格驱动测试 - 刷新Token
func TestJWTService_RefreshToken_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(JWTService, context.Context) string
		wantErr      bool
		errMsg       string
		checkRevoked bool
	}{
		{
			name: "正常刷新Token",
			setup: func(s JWTService, ctx context.Context) string {
				token, _ := s.GenerateToken(ctx, "user_123", []string{"reader"})
				return token
			},
			wantErr: false,
		},
		{
			name: "刷新已吊销的Token",
			setup: func(s JWTService, ctx context.Context) string {
				token, _ := s.GenerateToken(ctx, "user_123", []string{"reader"})
				s.RevokeToken(ctx, token)
				return token
			},
			wantErr:      true,
			checkRevoked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cache := NewMockRedisClient()
			cfg := &config.JWTConfigEnhanced{
				SecretKey:       "test-secret-key-12345678",
				Issuer:          "qingyu-test",
				Expiration:      1 * time.Hour,
				RefreshDuration: 24 * time.Hour,
			}
			service := NewJWTService(cfg, cache)
			ctx := context.Background()

			oldToken := tt.setup(service, ctx)

			// Act
			newToken, err := service.RefreshToken(ctx, oldToken)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, newToken)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, newToken)
				assert.NotEqual(t, oldToken, newToken)

				// 验证新Token有效
				claims, err := service.ValidateToken(ctx, newToken)
				assert.NoError(t, err)
				assert.Equal(t, "user_123", claims.UserID)
			}
		})
	}
}

// TestJWTService_RevokeToken 表格驱动测试 - 吊销Token
func TestJWTService_RevokeToken_TableDriven(t *testing.T) {
	cache := NewMockRedisClient()
	cfg := &config.JWTConfigEnhanced{
		Secret:      "test-secret-key-12345678",
		Expiration:  1 * time.Hour,
	}
	service := NewJWTService(cfg, cache)
	ctx := context.Background()

	tests := []struct {
		name    string
		token   func() string
		wantErr bool
		errMsg  string
	}{
		{
			name: "吊销有效Token",
			token: func() string {
				t, _ := service.GenerateToken(ctx, "user_123", []string{"reader"})
				return t
			},
			wantErr: false,
		},
		{
			name: "吊销已吊销的Token（幂等）",
			token: func() string {
				t, _ := service.GenerateToken(ctx, "user_123", []string{"reader"})
				service.RevokeToken(ctx, t)
				return t
			},
			wantErr: false,
		},
		{
			name:  "吊销无效Token",
			token: func() string { return "invalid_token" },
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			token := tt.token()

			// Act
			err := service.RevokeToken(ctx, token)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// BenchmarkJWTService_GenerateToken 基准测试 - 生成Token
func BenchmarkJWTService_GenerateToken(b *testing.B) {
	cache := NewMockRedisClient()
	cfg := &config.JWTConfigEnhanced{
		Secret:      "test-secret-key-12345678",
		Expiration:  1 * time.Hour,
	}
	service := NewJWTService(cfg, cache)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GenerateToken(ctx, "user_123", []string{"reader", "author"})
	}
}

// BenchmarkJWTService_ValidateToken 基准测试 - 验证Token
func BenchmarkJWTService_ValidateToken(b *testing.B) {
	cache := NewMockRedisClient()
	cfg := &config.JWTConfigEnhanced{
		Secret:      "test-secret-key-12345678",
		Expiration:  1 * time.Hour,
	}
	service := NewJWTService(cfg, cache)
	ctx := context.Background()

	token, _ := service.GenerateToken(ctx, "user_123", []string{"reader"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateToken(ctx, token)
	}
}

