package user

import (
	user2 "Qingyu_backend/service/interfaces/user"
	"context"
	"errors"
	"testing"

	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/service/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUserService_SendEmailVerification 测试发送邮箱验证码
func TestUserService_SendEmailVerification(t *testing.T) {
	tests := []struct {
		name          string
		request       *user2.SendEmailVerificationRequest
		setupMock     func(*mocks.MockUserRepository)
		expectError   bool
		errorContains string
		checkResponse func(*testing.T, *user2.SendEmailVerificationResponse)
	}{
		{
			name: "成功发送验证码_邮箱未验证",
			request: &user2.SendEmailVerificationRequest{
				UserID: "user123",
				Email:  "test@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:            "user123",
					Username:      "testuser",
					Email:         "test@example.com",
					EmailVerified: false,
				}
				m.On("GetByID", mock.Anything, "user123").Return(user, nil)
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *user2.SendEmailVerificationResponse) {
				assert.True(t, resp.Success)
				assert.Equal(t, "验证码已发送到您的邮箱", resp.Message)
				assert.Equal(t, 1800, resp.ExpiresIn) // 30分钟
			},
		},
		{
			name: "邮箱已验证_直接返回成功",
			request: &user2.SendEmailVerificationRequest{
				UserID: "user123",
				Email:  "test@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:            "user123",
					Username:      "testuser",
					Email:         "test@example.com",
					EmailVerified: true, // 已验证
				}
				m.On("GetByID", mock.Anything, "user123").Return(user, nil)
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *user2.SendEmailVerificationResponse) {
				assert.True(t, resp.Success)
				assert.Equal(t, "邮箱已验证", resp.Message)
				assert.Equal(t, 0, resp.ExpiresIn)
			},
		},
		{
			name: "验证失败_用户ID为空",
			request: &user2.SendEmailVerificationRequest{
				UserID: "",
				Email:  "test@example.com",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			expectError:   true,
			errorContains: "用户ID和邮箱不能为空",
		},
		{
			name: "验证失败_邮箱为空",
			request: &user2.SendEmailVerificationRequest{
				UserID: "user123",
				Email:  "",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			expectError:   true,
			errorContains: "用户ID和邮箱不能为空",
		},
		{
			name: "验证失败_用户不存在",
			request: &user2.SendEmailVerificationRequest{
				UserID: "nonexistent",
				Email:  "test@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.Anything, "nonexistent").Return(nil, errors.New("用户不存在"))
			},
			expectError:   true,
			errorContains: "用户不存在",
		},
		{
			name: "验证失败_邮箱不匹配",
			request: &user2.SendEmailVerificationRequest{
				UserID: "user123",
				Email:  "different@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:    "user123",
					Email: "test@example.com", // 不同的邮箱
				}
				m.On("GetByID", mock.Anything, "user123").Return(user, nil)
			},
			expectError:   true,
			errorContains: "邮箱不匹配",
		},
		{
			name: "验证失败_数据库错误",
			request: &user2.SendEmailVerificationRequest{
				UserID: "user123",
				Email:  "test@example.com",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.Anything, "user123").Return(nil, errors.New("数据库连接失败"))
			},
			expectError:   true,
			errorContains: "获取用户失败",
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
			resp, err := service.SendEmailVerification(ctx, tt.request)

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

// TestUserService_VerifyEmail 测试验证邮箱
func TestUserService_VerifyEmail(t *testing.T) {
	type setupTokensFunc func(*EmailVerificationTokenManager, string) string

	tests := []struct {
		name          string
		request       *user2.VerifyEmailRequest
		setupMock     func(*mocks.MockUserRepository)
		setupTokens   setupTokensFunc
		expectError   bool
		errorContains string
		checkResponse func(*testing.T, *user2.VerifyEmailResponse)
	}{
		{
			name: "验证成功",
			request: &user2.VerifyEmailRequest{
				UserID: "user123",
				Code:   "", // 在setup中填充
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:            "user123",
					Username:      "testuser",
					Email:         "test@example.com",
					EmailVerified: false,
					Status:        usersModel.UserStatusInactive,
				}
				m.On("GetByID", mock.Anything, "user123").Return(user, nil)
				m.On("Update", mock.Anything, "user123", mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
			setupTokens: func(mgr *EmailVerificationTokenManager, email string) string {
				ctx := context.Background()
				code, _ := mgr.GenerateCode(ctx, "user123", email)
				return code
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *user2.VerifyEmailResponse) {
				assert.True(t, resp.Success)
				assert.Equal(t, "邮箱验证成功", resp.Message)
			},
		},
		{
			name: "验证失败_用户ID为空",
			request: &user2.VerifyEmailRequest{
				UserID: "",
				Code:   "123456",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			setupTokens:   func(mgr *EmailVerificationTokenManager, email string) string { return "" },
			expectError:   true,
			errorContains: "用户ID和验证码不能为空",
		},
		{
			name: "验证失败_验证码为空",
			request: &user2.VerifyEmailRequest{
				UserID: "user123",
				Code:   "",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			setupTokens:   func(mgr *EmailVerificationTokenManager, email string) string { return "" },
			expectError:   true,
			errorContains: "用户ID和验证码不能为空",
		},
		{
			name: "验证失败_用户不存在",
			request: &user2.VerifyEmailRequest{
				UserID: "nonexistent",
				Code:   "123456",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.Anything, "nonexistent").Return(nil, errors.New("用户不存在"))
			},
			setupTokens:   func(mgr *EmailVerificationTokenManager, email string) string { return "" },
			expectError:   true,
			errorContains: "用户不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

			// 设置验证码（如果需要）
			if tt.setupTokens != nil {
				tokenManager := NewEmailVerificationTokenManager()
				// 获取用户邮箱
				if tt.request.UserID != "" {
					if user, err := mockRepo.GetByID(ctx, tt.request.UserID); err == nil && user != nil {
						code := tt.setupTokens(tokenManager, user.Email)
						if tt.request.Code == "" {
							tt.request.Code = code
						}
					}
				}
			}

			// Act
			resp, err := service.VerifyEmail(ctx, tt.request)

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

// TestUserService_VerifyEmail_Integration 集成测试：完整的验证流程
func TestUserService_VerifyEmail_Integration(t *testing.T) {
	// TODO: 重构服务以支持依赖注入，然后完成集成测试
	t.Skip("需要重构服务以支持tokenManager依赖注入")
}

// Benchmark_SendEmailVerification 性能测试：发送验证码
func Benchmark_SendEmailVerification(b *testing.B) {
	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepository)
	user := &usersModel.User{
		ID:            "user123",
		Username:      "testuser",
		Email:         "test@example.com",
		EmailVerified: false,
	}
	mockRepo.On("GetByID", mock.Anything, "user123").Return(user, nil)
	service := &UserServiceImpl{userRepo: mockRepo, name: "UserService"}

	req := &user2.SendEmailVerificationRequest{
		UserID: "user123",
		Email:  "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.SendEmailVerification(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
