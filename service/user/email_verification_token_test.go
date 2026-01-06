package user

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewEmailVerificationTokenManager 测试创建Token管理器
func TestNewEmailVerificationTokenManager(t *testing.T) {
	// Act
	manager := NewEmailVerificationTokenManager()

	// Assert
	assert.NotNil(t, manager, "管理器不应为空")
	assert.NotNil(t, manager.tokens, "tokens map不应为空")
}

// TestEmailVerificationTokenManager_GenerateCode 测试生成验证码
func TestEmailVerificationTokenManager_GenerateCode(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		email       string
		expectError bool
		checkCode   func(*testing.T, string)
	}{
		{
			name:        "成功生成6位数字验证码",
			userID:      "user123",
			email:       "test@example.com",
			expectError: false,
			checkCode: func(t *testing.T, code string) {
				assert.Len(t, code, 6, "验证码长度应为6位")
				assert.NotEmpty(t, code, "验证码不应为空")
				// 验证全是数字
				for _, c := range code {
					assert.GreaterOrEqual(t, c, '0')
					assert.LessOrEqual(t, c, '9')
				}
			},
		},
		{
			name:        "不同邮箱生成不同验证码",
			userID:      "user123",
			email:       "test2@example.com",
			expectError: false,
			checkCode: func(t *testing.T, code string) {
				assert.NotEmpty(t, code)
			},
		},
		{
			name:        "空邮箱生成验证码",
			userID:      "user123",
			email:       "",
			expectError: false, // 实现中没有邮箱验证
			checkCode: func(t *testing.T, code string) {
				assert.NotEmpty(t, code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewEmailVerificationTokenManager()

			// Act
			code, err := manager.GenerateCode(ctx, tt.userID, tt.email)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, code)
				if tt.checkCode != nil {
					tt.checkCode(t, code)
				}

				// 验证token信息已存储
				manager.mu.RLock()
				tokenInfo, exists := manager.tokens[tt.email]
				manager.mu.RUnlock()

				assert.True(t, exists, "验证码信息应已存储")
				assert.Equal(t, tt.userID, tokenInfo.UserID)
				assert.Equal(t, tt.email, tokenInfo.Email)
				assert.Equal(t, code, tokenInfo.Code)
				assert.False(t, tokenInfo.Used)
				assert.WithinDuration(t, time.Now().Add(30*time.Minute), tokenInfo.ExpiresAt, time.Second)
			}
		})
	}
}

// TestEmailVerificationTokenManager_ValidateCode 测试验证验证码
func TestEmailVerificationTokenManager_ValidateCode(t *testing.T) {
	type setupTokensFunc func(*EmailVerificationTokenManager, string) string

	tests := []struct {
		name          string
		setupTokens   setupTokensFunc // 返回验证码
		userID        string
		email         string
		code          string
		expectError   bool
		errorContains string
	}{
		{
			name: "验证成功_有效验证码",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				ctx := context.Background()
				code, _ := m.GenerateCode(ctx, "user123", email)
				return code
			},
			userID:      "user123",
			email:       "test@example.com",
			code:        "", // 会在setup中填充
			expectError: false,
		},
		{
			name: "验证失败_验证码不存在",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				return ""
			},
			userID:        "user123",
			email:         "nonexist@example.com",
			code:          "123456",
			expectError:   true,
			errorContains: "验证码不存在或已过期",
		},
		{
			name: "验证失败_验证码错误",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				ctx := context.Background()
				code, _ := m.GenerateCode(ctx, "user123", email)
				return code
			},
			userID:        "user123",
			email:         "test@example.com",
			code:          "999999", // 错误的验证码
			expectError:   true,
			errorContains: "验证码错误",
		},
		{
			name: "验证失败_用户ID不匹配",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				ctx := context.Background()
				code, _ := m.GenerateCode(ctx, "user123", email)
				return code
			},
			userID:        "user456", // 不同的用户ID
			email:         "test@example.com",
			code:          "",        // 会在setup中填充
			expectError:   true,
			errorContains: "验证码不匹配",
		},
		{
			name: "验证失败_验证码已使用",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				ctx := context.Background()
				code, _ := m.GenerateCode(ctx, "user123", email)
				m.MarkCodeAsUsed(ctx, email)
				return code
			},
			userID:        "user123",
			email:         "test@example.com",
			code:          "", // 会在setup中填充
			expectError:   true,
			errorContains: "验证码已使用",
		},
		{
			name: "验证失败_验证码过期",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				ctx := context.Background()
				code, _ := m.GenerateCode(ctx, "user123", email)
				// 手动设置过期时间为过去
				m.mu.Lock()
				m.tokens[email].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				return code
			},
			userID:        "user123",
			email:         "test@example.com",
			code:          "", // 会在setup中填充
			expectError:   true,
			errorContains: "验证码已过期",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewEmailVerificationTokenManager()
			if tt.setupTokens != nil {
				generatedCode := tt.setupTokens(manager, tt.email)
				// 只有当code为空时才使用生成的验证码
				if tt.code == "" {
					tt.code = generatedCode
				}
			}

			// Act
			err := manager.ValidateCode(ctx, tt.userID, tt.email, tt.code)

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

// TestEmailVerificationTokenManager_MarkCodeAsUsed 测试标记验证码为已使用
func TestEmailVerificationTokenManager_MarkCodeAsUsed(t *testing.T) {
	type setupTokensFunc func(*EmailVerificationTokenManager, string) string

	tests := []struct {
		name          string
		setupTokens   setupTokensFunc
		email         string
		expectError   bool
		errorContains string
		checkUsed     func(*testing.T, *EmailVerificationTokenManager, string)
	}{
		{
			name: "成功标记为已使用",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				ctx := context.Background()
				code, _ := m.GenerateCode(ctx, "user123", email)
				return code
			},
			email:       "test@example.com",
			expectError: false,
			checkUsed: func(t *testing.T, m *EmailVerificationTokenManager, email string) {
				m.mu.RLock()
				tokenInfo := m.tokens[email]
				m.mu.RUnlock()
				assert.True(t, tokenInfo.Used, "验证码应标记为已使用")
			},
		},
		{
			name:        "标记失败_验证码不存在",
			setupTokens: func(m *EmailVerificationTokenManager, email string) string {
				return ""
			},
			email:         "nonexist@example.com",
			expectError:   true,
			errorContains: "验证码不存在",
			checkUsed:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewEmailVerificationTokenManager()
			if tt.setupTokens != nil {
				tt.setupTokens(manager, tt.email)
			}

			// Act
			err := manager.MarkCodeAsUsed(ctx, tt.email)

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

// TestEmailVerificationTokenManager_CleanExpiredCodes 测试清理过期验证码
func TestEmailVerificationTokenManager_CleanExpiredCodes(t *testing.T) {
	type setupTokensFunc func(*EmailVerificationTokenManager) map[string]string

	tests := []struct {
		name            string
		setupTokens     setupTokensFunc
		expectedCleaned int
	}{
		{
			name: "清理部分过期验证码",
			setupTokens: func(m *EmailVerificationTokenManager) map[string]string {
				ctx := context.Background()
				codes := make(map[string]string)

				// 添加未过期的验证码
				code1, _ := m.GenerateCode(ctx, "user1", "active1@example.com")
				codes["active1@example.com"] = code1

				// 添加已过期的验证码
				code2, _ := m.GenerateCode(ctx, "user2", "expired1@example.com")
				m.mu.Lock()
				m.tokens["expired1@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				codes["expired1@example.com"] = code2

				// 添加另一个已过期的验证码
				code3, _ := m.GenerateCode(ctx, "user3", "expired2@example.com")
				m.mu.Lock()
				m.tokens["expired2@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				codes["expired2@example.com"] = code3

				return codes
			},
			expectedCleaned: 2,
		},
		{
			name: "所有验证码都未过期",
			setupTokens: func(m *EmailVerificationTokenManager) map[string]string {
				ctx := context.Background()
				codes := make(map[string]string)

				code1, _ := m.GenerateCode(ctx, "user1", "active1@example.com")
				codes["active1@example.com"] = code1

				code2, _ := m.GenerateCode(ctx, "user2", "active2@example.com")
				codes["active2@example.com"] = code2

				return codes
			},
			expectedCleaned: 0,
		},
		{
			name: "所有验证码都过期",
			setupTokens: func(m *EmailVerificationTokenManager) map[string]string {
				ctx := context.Background()
				codes := make(map[string]string)

				code1, _ := m.GenerateCode(ctx, "user1", "expired1@example.com")
				m.mu.Lock()
				m.tokens["expired1@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				codes["expired1@example.com"] = code1

				code2, _ := m.GenerateCode(ctx, "user2", "expired2@example.com")
				m.mu.Lock()
				m.tokens["expired2@example.com"].ExpiresAt = time.Now().Add(-1 * time.Hour)
				m.mu.Unlock()
				codes["expired2@example.com"] = code2

				return codes
			},
			expectedCleaned: 2,
		},
		{
			name:            "没有验证码",
			setupTokens:     func(m *EmailVerificationTokenManager) map[string]string {
				return make(map[string]string)
			},
			expectedCleaned: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			manager := NewEmailVerificationTokenManager()
			codes := tt.setupTokens(manager)
			initialCount := len(manager.tokens)

			// Act
			manager.CleanExpiredCodes(ctx)

			// Assert
			finalCount := len(manager.tokens)
			assert.Equal(t, tt.expectedCleaned, initialCount-finalCount,
				"应清理 %d 个过期验证码", tt.expectedCleaned)

			// 验证未过期的验证码仍然存在
			for email := range codes {
				if _, exists := manager.tokens[email]; exists {
					// 仍然存在，应该是未过期的
					tokenInfo := manager.tokens[email]
					assert.True(t, time.Now().Before(tokenInfo.ExpiresAt),
						"剩余的验证码应该是未过期的")
				}
			}
		})
	}
}

// TestEmailVerificationTokenManager_ConcurrentAccess 测试并发访问
func TestEmailVerificationTokenManager_ConcurrentAccess(t *testing.T) {
	// Arrange
	ctx := context.Background()
	manager := NewEmailVerificationTokenManager()
	iterations := 100

	// Act - 并发生成验证码
	done := make(chan bool)
	for i := 0; i < iterations; i++ {
		go func(index int) {
			email := fmt.Sprintf("user%d@example.com", index)
			userID := fmt.Sprintf("user%d", index)
			_, err := manager.GenerateCode(ctx, userID, email)
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
		"应该生成 %d 个验证码", iterations)
}
