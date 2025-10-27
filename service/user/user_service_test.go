package user

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/user"
	"context"
	"errors"
	"testing"

	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/service/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============ 测试用例 ============

// TestNewUserService 测试服务创建
func TestNewUserService(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)

	assert.NotNil(t, service, "服务不应为空")
	assert.Equal(t, "UserService", service.GetServiceName())
	assert.Equal(t, "1.0.0", service.GetVersion())
}

// TestUserService_Initialize 测试服务初始化
func TestUserService_Initialize(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mocks.MockUserRepository)
		expectError   bool
		errorContains string
	}{
		{
			name: "初始化成功",
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("Health", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "初始化失败_数据库连接失败",
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("Health", mock.Anything).Return(errors.New("连接失败"))
			},
			expectError:   true,
			errorContains: "连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := NewUserService(mockRepo)

			// Act
			err := service.Initialize(context.Background())

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_CreateUser 测试创建用户
func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		request       *serviceInterfaces.CreateUserRequest
		setupMock     func(*mocks.MockUserRepository)
		expectError   bool
		errorContains string
	}{
		{
			name: "创建成功",
			request: &serviceInterfaces.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("ExistsByUsername", mock.Anything, "testuser").Return(false, nil)
				m.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
				m.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "用户名为空",
			request: &serviceInterfaces.CreateUserRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			expectError:   true,
			errorContains: "用户名不能为空",
		},
		{
			name: "邮箱为空",
			request: &serviceInterfaces.CreateUserRequest{
				Username: "testuser",
				Email:    "",
				Password: "password123",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			expectError:   true,
			errorContains: "邮箱不能为空",
		},
		{
			name: "用户名已存在",
			request: &serviceInterfaces.CreateUserRequest{
				Username: "existuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("ExistsByUsername", mock.Anything, "existuser").Return(true, nil)
			},
			expectError:   true,
			errorContains: "用户名已存在",
		},
		{
			name: "邮箱已存在",
			request: &serviceInterfaces.CreateUserRequest{
				Username: "testuser",
				Email:    "exist@example.com",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("ExistsByUsername", mock.Anything, "testuser").Return(false, nil)
				m.On("ExistsByEmail", mock.Anything, "exist@example.com").Return(true, nil)
			},
			expectError:   true,
			errorContains: "邮箱已存在",
		},
		{
			name: "数据库创建失败",
			request: &serviceInterfaces.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("ExistsByUsername", mock.Anything, "testuser").Return(false, nil)
				m.On("ExistsByEmail", mock.Anything, "test@example.com").Return(false, nil)
				m.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).Return(errors.New("数据库错误"))
			},
			expectError:   true,
			errorContains: "创建用户失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := NewUserService(mockRepo)

			// Act
			resp, err := service.CreateUser(context.Background(), tt.request)

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
				assert.NotNil(t, resp.User)
				assert.Equal(t, tt.request.Username, resp.User.Username)
				assert.Equal(t, tt.request.Email, resp.User.Email)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_GetUser 测试获取用户
func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name          string
		request       *serviceInterfaces.GetUserRequest
		setupMock     func(*mocks.MockUserRepository)
		expectError   bool
		errorContains string
	}{
		{
			name: "获取成功",
			request: &serviceInterfaces.GetUserRequest{
				ID: "user123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				user := &usersModel.User{
					ID:       "user123",
					Username: "testuser",
					Email:    "test@example.com",
				}
				m.On("GetByID", mock.Anything, "user123").Return(user, nil)
			},
			expectError: false,
		},
		{
			name: "ID为空",
			request: &serviceInterfaces.GetUserRequest{
				ID: "",
			},
			setupMock:     func(m *mocks.MockUserRepository) {},
			expectError:   true,
			errorContains: "ID不能为空",
		},
		{
			name: "数据库查询失败",
			request: &serviceInterfaces.GetUserRequest{
				ID: "user123",
			},
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.Anything, "user123").Return(nil, errors.New("数据库错误"))
			},
			expectError:   true,
			errorContains: "获取用户失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := NewUserService(mockRepo)

			// Act
			resp, err := service.GetUser(context.Background(), tt.request)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)
				assert.Equal(t, tt.request.ID, resp.User.ID)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserService_Health 测试健康检查
func TestUserService_Health(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*mocks.MockUserRepository)
		expectError   bool
		errorContains string
	}{
		{
			name: "健康检查通过",
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("Health", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "健康检查失败",
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("Health", mock.Anything).Return(errors.New("数据库连接失败"))
			},
			expectError:   true,
			errorContains: "数据库连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockUserRepository)
			tt.setupMock(mockRepo)
			service := NewUserService(mockRepo)

			// Act
			err := service.Health(context.Background())

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
