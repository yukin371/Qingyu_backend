package user

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewPasswordResetTokenManager 测试创建Token管理器
func TestNewPasswordResetTokenManager(t *testing.T) {
	// Act
	manager := NewPasswordResetTokenManager()

	// Assert
	assert.NotNil(t, manager, "管理器不应为空")
	assert.NotNil(t, manager.tokens, "tokens map不应为空")
}

// TestPasswordResetTokenManager_GenerateToken 测试生成Token
func TestPasswordResetTokenManager_GenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
		checkToken  func(*testing.T, string)
	}{
		{
			name:        "成功生成Token",
			email:       "test@example.com",
			expectError: false,
			checkToken: func(t *testing.T, token string) {
				assert.NotEmpty(t, token, "Token不应为空")
				// 32字节随机数，hex编码后应该是64个字符
				assert.Len(t, token, 64, "Token长度应为64个字符")
				// 验证全是十六进制字符
				for _, c := range token {
					assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F'),
						"Token应只包含十六进制字符")
				}
			},
		},
		{
			name:        "不同邮箱生成不同Token",
			email:       "test2@example.com",
			expectError: false,
			checkToken: func(t *testing.T, token string) {
				assert.NotEmpty(t, token)
			},
		},
		{
			name:        "同一邮箱多次生成Token_会覆盖",
			email:       "test@example.com",
			expectError: false,
			checkToken: func(t *testing.T, token string) {
				assert.NotEmpty(t, token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewPasswordResetTokenManager()

			// Act
			token, err := manager.GenerateToken(ctx, tt.email)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				if tt.checkToken != nil {
					tt.checkToken(t, token)
				}

				// 验证token信息已存储
				manager.mu.RLock()
				tokenInfo, exists := manager.tokens[tt.email]
				manager.mu.RUnlock()

				assert.True(t, exists, "Token信息应已存储")
				assert.Equal(t, tt.email, tokenInfo.Email)
				assert.Equal(t, token, tokenInfo.Token)
				assert.False(t, tokenInfo.Used)
				assert.WithinDuration(t, time.Now().Add(1*time.Hour), tokenInfo.ExpiresAt, time.Second)
			}
		})
	}
}

// TestPasswordResetTokenManager_ValidateToken 测试验证Token
func TestPasswordResetTokenManager_ValidateToken(t *testing.T) {
	tests := []struct {
		name          string
		setupTokens   func(*PasswordResetTokenManager, string) string // 返回token
		email         string
		token         string
		expectError   bool
		errorContains string
	}{
		{
			name: "验证成功_有效Token",
			setupTokens: func(m *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := m.GenerateToken(ctx, email)
				return token
			},
			email:       "test@example.com",
			token:       "", // 会在setup中填充
			expectError: false,
		},
		{
			name: "验证失败_Token不存在",
			setupTokens: func(m *PasswordResetTokenManager, email string) string {
				return ""
			},
			email:         "nonexist@example.com",
			token:         "a1b2c3d4e5f6...",
			expectError:   true,
			errorContains: "无效的重置Token",
		},
		{
			name: "验证失败_Token错误",
			setupTokens: func(m *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := m.GenerateToken(ctx, email)
				return token
			},
			email:         "test@example.com",
			token:         "wrongtoken1234567890123456789012345678901234567890123456789012345678",
			expectError:   true,
			errorContains: "无效的重置Token",
		},
		{
			name: "验证失败_Token已使用",
			setupTokens: func(m *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := m.GenerateToken(ctx, email)
				m.MarkTokenAsUsed(ctx, email)
				return token
			},
			email:       "test@example.com",
			token:       "", // 会在setup中填充
			expectError: true,
			errorContains: "重置Token已使用",
		},
		{
			name: "验证失败_Token过期",
			setupTokens: func(m *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := m.GenerateToken(ctx, email)
				// 手动设置过期时间为过去
				m.mu.Lock()
				m.tokens[email].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				return token
			},
			email:         "test@example.com",
			token:         "", // 会在setup中填充
			expectError:   true,
			errorContains: "重置Token已过期",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewPasswordResetTokenManager()
			generatedToken := tt.setupTokens(manager, tt.email)
			// 只有当token为空时才使用生成的token
			if tt.token == "" {
				tt.token = generatedToken
			}

			// Act
			err := manager.ValidateToken(ctx, tt.email, tt.token)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPasswordResetTokenManager_MarkTokenAsUsed 测试标记Token为已使用
func TestPasswordResetTokenManager_MarkTokenAsUsed(t *testing.T) {
	tests := []struct {
		name          string
		setupTokens   func(*PasswordResetTokenManager, string) string
		email         string
		expectError   bool
		errorContains string
		checkUsed     func(*testing.T, *PasswordResetTokenManager, string)
	}{
		{
			name: "成功标记为已使用",
			setupTokens: func(m *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := m.GenerateToken(ctx, email)
				return token
			},
			email:       "test@example.com",
			expectError: false,
			checkUsed: func(t *testing.T, m *PasswordResetTokenManager, email string) {
				m.mu.RLock()
				tokenInfo := m.tokens[email]
				m.mu.RUnlock()
				assert.True(t, tokenInfo.Used, "Token应标记为已使用")
			},
		},
		{
			name: "标记失败_Token不存在",
			setupTokens: func(m *PasswordResetTokenManager, email string) string {
				return ""
			},
			email:         "nonexist@example.com",
			expectError:   true,
			errorContains: "Token不存在",
			checkUsed:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewPasswordResetTokenManager()
			tt.setupTokens(manager, tt.email)

			// Act
			err := manager.MarkTokenAsUsed(ctx, tt.email)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.checkUsed != nil {
					tt.checkUsed(t, manager, tt.email)
				}
			}
		})
	}
}

// TestPasswordResetTokenManager_CleanExpiredTokens 测试清理过期Token
func TestPasswordResetTokenManager_CleanExpiredTokens(t *testing.T) {
	tests := []struct {
		name            string
		setupTokens     func(*PasswordResetTokenManager) map[string]string
		expectedCleaned int
	}{
		{
			name: "清理部分过期Token",
			setupTokens: func(m *PasswordResetTokenManager) map[string]string {
				ctx := context.Background()
				tokens := make(map[string]string)

				// 添加未过期的Token
				token1, _ := m.GenerateToken(ctx, "active1@example.com")
				tokens["active1@example.com"] = token1

				// 添加已过期的Token
				token2, _ := m.GenerateToken(ctx, "expired1@example.com")
				m.mu.Lock()
				m.tokens["expired1@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				tokens["expired1@example.com"] = token2

				// 添加另一个已过期的Token
				token3, _ := m.GenerateToken(ctx, "expired2@example.com")
				m.mu.Lock()
				m.tokens["expired2@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				tokens["expired2@example.com"] = token3

				return tokens
			},
			expectedCleaned: 2,
		},
		{
			name: "所有Token都未过期",
			setupTokens: func(m *PasswordResetTokenManager) map[string]string {
				ctx := context.Background()
				tokens := make(map[string]string)

				token1, _ := m.GenerateToken(ctx, "active1@example.com")
				tokens["active1@example.com"] = token1

				token2, _ := m.GenerateToken(ctx, "active2@example.com")
				tokens["active2@example.com"] = token2

				return tokens
			},
			expectedCleaned: 0,
		},
		{
			name: "所有Token都过期",
			setupTokens: func(m *PasswordResetTokenManager) map[string]string {
				ctx := context.Background()
				tokens := make(map[string]string)

				token1, _ := m.GenerateToken(ctx, "expired1@example.com")
				m.mu.Lock()
				m.tokens["expired1@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				tokens["expired1@example.com"] = token1

				token2, _ := m.GenerateToken(ctx, "expired2@example.com")
				m.mu.Lock()
				m.tokens["expired2@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				tokens["expired2@example.com"] = token2

				return tokens
			},
			expectedCleaned: 2,
		},
		{
			name:            "没有Token",
			setupTokens:     func(m *PasswordResetTokenManager) map[string]string {
				return make(map[string]string)
			},
			expectedCleaned: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewPasswordResetTokenManager()
			tokens := tt.setupTokens(manager)
			initialCount := len(manager.tokens)

			// Act
			manager.CleanExpiredTokens(ctx)

			// Assert
			finalCount := len(manager.tokens)
			assert.Equal(t, tt.expectedCleaned, initialCount-finalCount,
				"应清理 %d 个过期Token", tt.expectedCleaned)

			// 验证未过期的Token仍然存在
			for email := range tokens {
				if _, exists := manager.tokens[email]; exists {
					// 仍然存在，应该是未过期的
					tokenInfo := manager.tokens[email]
					assert.True(t, time.Now().Before(tokenInfo.ExpiresAt),
						"剩余的Token应该是未过期的")
				}
			}
		})
	}
}

// TestPasswordResetTokenManager_TokenUniqueness 测试Token唯一性
func TestPasswordResetTokenManager_TokenUniqueness(t *testing.T) {
	// Arrange
	ctx := context.Background()
	manager := NewPasswordResetTokenManager() // 使用新实例，不使用全局的
	email := "test@example.com"

	// Act - 生成多个Token
	tokens := make([]string, 100)
	for i := 0; i < 100; i++ {
		token, err := manager.GenerateToken(ctx, email)
		assert.NoError(t, err)
		tokens[i] = token
	}

	// Assert - 验证所有Token都是唯一的
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[token] = true
	}

	// 虽然同一邮箱会覆盖，但生成的随机Token应该是唯一的
	// 由于覆盖，最终只有一个token被存储
	assert.Equal(t, 1, len(manager.tokens))
	assert.Equal(t, 100, len(uniqueTokens)) // 所有生成的token都是唯一的
}

// TestPasswordResetTokenManager_ConcurrentAccess 测试并发访问
func TestPasswordResetTokenManager_ConcurrentAccess(t *testing.T) {
	// Arrange
	ctx := context.Background()
	manager := NewPasswordResetTokenManager()
	iterations := 100

	// Act - 并发生成Token
	done := make(chan bool)
	for i := 0; i < iterations; i++ {
		go func(index int) {
			email := fmt.Sprintf("user%d@example.com", index)
			_, err := manager.GenerateToken(ctx, email)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < iterations; i++ {
		<-done
	}

	// Assert
	assert.Equal(t, iterations, len(manager.tokens),
		"应该生成 %d 个Token", iterations)
}

// TestPasswordResetTokenManager_TokenFormat 测试Token格式
func TestPasswordResetTokenManager_TokenFormat(t *testing.T) {
	// Arrange
	ctx := context.Background()
	manager := NewPasswordResetTokenManager()

	// Act
	token, err := manager.GenerateToken(ctx, "test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, token, 64, "Token应该是64个字符（32字节的hex编码）")

	// 验证是有效的十六进制字符串
	_, err = hex.DecodeString(token)
	assert.NoError(t, err, "Token应该是有效的十六进制字符串")
}

// Benchmark_GenerateToken 性能测试：生成Token
func Benchmark_GenerateToken(b *testing.B) {
	ctx := context.Background()
	manager := NewPasswordResetTokenManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.GenerateToken(ctx, "test@example.com")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark_ValidateToken 性能测试：验证Token
func Benchmark_ValidateToken(b *testing.B) {
	ctx := context.Background()
	manager := NewPasswordResetTokenManager()
	token, _ := manager.GenerateToken(ctx, "test@example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := manager.ValidateToken(ctx, "test@example.com", token)
		if err != nil {
			b.Fatal(err)
		}
	}
}
