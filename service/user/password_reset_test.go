package user

import (
	user2 "Qingyu_backend/service/interfaces/user"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	"context"
	"errors"
	"testing"
	"time"

	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/service/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUserService_RequestPasswordReset 测试请求密码重置
func TestUserService_RequestPasswordReset(t *testing.T) {
	tests := []struct {
		name          string
		request       *user2.RequestPasswordResetRequest
		setupMock     func(*mocks.MockUserRepository)
		expectError   bool
		errorContains string
		checkResponse func(*testing.T, *user2.RequestPasswordResetResponse)
	}{
		{
			name: "成功发送重置邮件_用户存在",
			request: &user2.RequestPasswordResetRequest{
				Email: "test@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *user2.RequestPasswordResetResponse) {
				assert.True(t, resp.Success)
				assert.Equal(t, "密码重置邮件已发送", resp.Message)
				assert.Equal(t, 3600, resp.ExpiresIn) // 1小时
			},
		},
		{
			name: "成功_用户不存在但仍返回成功（防止邮箱枚举）",
			request: &user2.RequestPasswordResetRequest{
				Email: "nonexist@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "nonexist@example.com").Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "用户不存在", nil))
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *user2.RequestPasswordResetResponse) {
				assert.True(t, resp.Success)
				assert.Equal(t, "如果该邮箱已注册，您将收到密码重置邮件", resp.Message)
			},
		},
		{
			name: "验证失败_邮箱为空",
			request: &user2.RequestPasswordResetRequest{
				Email: "",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			expectError:   true,
			errorContains: "邮箱不能为空",
		},
		{
			name: "验证失败_数据库错误",
			request: &user2.RequestPasswordResetRequest{
				Email: "test@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("数据库连接失败"))
			},
			expectError:   true,
			errorContains: "检查用户失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

			// Act
			resp, err := service.RequestPasswordReset(ctx, tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_ConfirmPasswordReset 测试确认密码重置
func TestUserService_ConfirmPasswordReset(t *testing.T) {
	tests := []struct {
		name          string
		request       *user2.ConfirmPasswordResetRequest
		setupMock     func(*mocks.MockUserRepository)
		setupTokens   func(*PasswordResetTokenManager, string) string
		expectError   bool
		errorContains string
		checkResponse func(*testing.T, *user2.ConfirmPasswordResetResponse)
	}{
		{
			name: "成功重置密码",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    "", // 在setup中填充
				Password: "NewPassword123!",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
				m.On("UpdatePassword", mock.Anything, "user123", mock.AnythingOfType("string")).Return(nil)
			},
			setupTokens: func(mgr *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := mgr.GenerateToken(ctx, email)
				return token
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *user2.ConfirmPasswordResetResponse) {
				assert.True(t, resp.Success)
				assert.Equal(t, "密码重置成功，请使用新密码登录", resp.Message)
			},
		},
		{
			name: "验证失败_邮箱为空",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "",
				Token:    "token123",
				Password: "NewPassword123!",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			setupTokens:   func(mgr *PasswordResetTokenManager, email string) string { return "" },
			expectError:   true,
			errorContains: "邮箱、Token和新密码不能为空",
		},
		{
			name: "验证失败_Token为空",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    "",
				Password: "NewPassword123!",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			setupTokens:   func(mgr *PasswordResetTokenManager, email string) string { return "" },
			expectError:   true,
			errorContains: "邮箱、Token和新密码不能为空",
		},
		{
			name: "验证失败_密码为空",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    "token123",
				Password: "",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			setupTokens:   func(mgr *PasswordResetTokenManager, email string) string { return "" },
			expectError:   true,
			errorContains: "邮箱、Token和新密码不能为空",
		},
		{
			name: "验证失败_用户不存在",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "nonexist@example.com",
				Token:    "token123",
				Password: "NewPassword123!",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "nonexist@example.com").Return(nil, errors.New("用户不存在"))
			},
			setupTokens:   func(mgr *PasswordResetTokenManager, email string) string { return "" },
			expectError:   true,
			errorContains: "用户不存在",
		},
		{
			name: "验证失败_Token无效",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    "wrongtoken1234567890123456789012345678901234567890123456789012345678",
				Password: "NewPassword123!",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:    "user123",
					Email: "test@example.com",
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			setupTokens: func(mgr *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := mgr.GenerateToken(ctx, email)
				return token // 返回正确的token，但request中使用错误的
			},
			expectError:   true,
			errorContains: "Token验证失败",
		},
		{
			name: "验证失败_Token已使用",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    "", // 在setup中填充
				Password: "NewPassword123!",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:    "user123",
					Email: "test@example.com",
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			setupTokens: func(mgr *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := mgr.GenerateToken(ctx, email)
				mgr.MarkTokenAsUsed(ctx, email)
				return token
			},
			expectError:   true,
			errorContains: "Token验证失败",
		},
		{
			name: "验证失败_Token过期",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    "", // 在setup中填充
				Password: "NewPassword123!",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:    "user123",
					Email: "test@example.com",
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			setupTokens: func(mgr *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := mgr.GenerateToken(ctx, email)
				// 手动设置过期时间
				mgr.mu.Lock()
				mgr.tokens[email].ExpiresAt = mgr.tokens[email].ExpiresAt.Add(-2 * time.Hour)
				mgr.mu.Unlock()
				return token
			},
			expectError:   true,
			errorContains: "Token验证失败",
		},
		{
			name: "验证失败_更新密码失败",
			request: &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    "", // 在setup中填充
				Password: "NewPassword123!",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
				m.On("UpdatePassword", mock.Anything, "user123", mock.AnythingOfType("string")).Return(errors.New("数据库更新失败"))
			},
			setupTokens: func(mgr *PasswordResetTokenManager, email string) string {
				ctx := context.Background()
				token, _ := mgr.GenerateToken(ctx, email)
				return token
			},
			expectError:   true,
			errorContains: "更新密码失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

			// 设置token（如果需要）
			if tt.setupTokens != nil {
				tokenManager := GetGlobalPasswordResetTokenManager()
				token := tt.setupTokens(tokenManager, tt.request.Email)
				if tt.request.Token == "" && token != "" {
					tt.request.Token = token
				}
			}

			// Act
			resp, err := service.ConfirmPasswordReset(ctx, tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.checkResponse != nil {
					tt.checkResponse(t, resp)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_ConfirmPasswordReset_PasswordStrength 测试密码强度
func TestUserService_ConfirmPasswordReset_PasswordStrength(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		shouldSucceed bool
	}{
		{
			name:          "强密码_包含大小写字母、数字和特殊字符",
			password:      "StrongP@ssw0rd!",
			shouldSucceed: true,
		},
		{
			name:          "中等强度密码_包含大小写字母和数字",
			password:      "MediumPass123",
			shouldSucceed: true,
		},
		{
			name:          "弱密码_只有小写字母",
			password:      "weak",
			shouldSucceed: true, // 当前实现没有密码强度验证
		},
		{
			name:          "非常短的密码",
			password:      "123",
			shouldSucceed: true, // 当前实现没有最小长度限制
		},
		{
			name:          "只有数字",
			password:      "12345678",
			shouldSucceed: true,
		},
		{
			name:          "包含特殊字符",
			password:      "!@#$%^&*()",
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockRepo := new(mocks.MockUserRepository)
			user := &usersModel.User{
				ID:       "user123",
				Username: "testuser",
				Email:    "test@example.com",
			}
			mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			mockRepo.On("UpdatePassword", mock.Anything, "user123", mock.AnythingOfType("string")).Return(nil)
			service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

			tokenManager := GetGlobalPasswordResetTokenManager()
			token, _ := tokenManager.GenerateToken(ctx, "test@example.com")

			request := &user2.ConfirmPasswordResetRequest{
				Email:    "test@example.com",
				Token:    token,
				Password: tt.password,
			}

			// Act
			resp, err := service.ConfirmPasswordReset(ctx, request)

			// Assert
			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_PasswordReset_Integration 集成测试：完整的密码重置流程
func TestUserService_PasswordReset_Integration(t *testing.T) {
	// TODO: 需要重构服务以支持依赖注入
	t.Skip("需要重构服务以支持tokenManager依赖注入")

	// Arrange
	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepository)
	service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

	user := &usersModel.User{
		ID:       "user123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	// 步骤1: 请求密码重置
	t.Log("步骤1: 请求密码重置")
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	resetReq := &user2.RequestPasswordResetRequest{
		Email: "test@example.com",
	}

	resetResp, err := service.RequestPasswordReset(ctx, resetReq)
	assert.NoError(t, err)
	assert.True(t, resetResp.Success)

	// 步骤2: 确认密码重置
	t.Log("步骤2: 确认密码重置")
	// 需要获取生成的token，但当前实现不支持
}

// TestUserService_ResetPassword_Deprecated 测试旧的ResetPassword方法
func TestUserService_ResetPassword_Deprecated(t *testing.T) {
	tests := []struct {
		name          string
		request       *user2.ResetPasswordRequest
		setupMock     func(*mocks.MockUserRepository)
		expectError   bool
		errorContains string
	}{
		{
			name: "成功发送重置邮件",
			request: &user2.ResetPasswordRequest{
				Email: "test@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectError: false,
		},
		{
			name: "用户不存在_仍返回成功（安全考虑）",
			request: &user2.ResetPasswordRequest{
				Email: "nonexist@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "nonexist@example.com").Return(nil, errors.New("用户不存在"))
			},
			expectError: false,
		},
		{
			name: "验证失败_邮箱为空",
			request: &user2.ResetPasswordRequest{
				Email: "",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			expectError:   true,
			errorContains: "邮箱不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

			// Act
			resp, err := service.ResetPassword(ctx, tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// Benchmark_RequestPasswordReset 性能测试：请求密码重置
func Benchmark_RequestPasswordReset(b *testing.B) {
	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepository)
	user := &usersModel.User{
		ID:       "user123",
		Username: "testuser",
		Email:    "test@example.com",
	}
	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

	req := &user2.RequestPasswordResetRequest{
		Email: "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.RequestPasswordReset(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
