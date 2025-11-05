package service_test

import (
	user2 "Qingyu_backend/service/interfaces/user"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/users"
	"Qingyu_backend/repository/interfaces/infrastructure"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	"Qingyu_backend/service/user"
	"Qingyu_backend/test/testutil"
)

// MockUserRepository 模拟UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *users.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*users.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id string, ip string) error {
	args := m.Called(ctx, id, ip)
	return args.Error(0)
}

func (m *MockUserRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// 以下是实现UserRepository接口所需的其他方法

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*users.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	args := m.Called(ctx, phone)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateStatus(ctx context.Context, id string, status users.UserStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserRepository) GetActiveUsers(ctx context.Context, limit int64) ([]*users.User, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersByRole(ctx context.Context, role string, limit int64) ([]*users.User, error) {
	args := m.Called(ctx, role, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

func (m *MockUserRepository) SetPhoneVerified(ctx context.Context, id string, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

func (m *MockUserRepository) BatchUpdateStatus(ctx context.Context, ids []string, status users.UserStatus) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockUserRepository) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockUserRepository) FindWithFilter(ctx context.Context, filter *users.UserFilter) ([]*users.User, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*users.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) SearchUsers(ctx context.Context, keyword string, limit int) ([]*users.User, error) {
	args := m.Called(ctx, keyword, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	args := m.Called(ctx, role)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountByStatus(ctx context.Context, status users.UserStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Transaction(ctx context.Context, u *users.User, fn func(context.Context, repoInterfaces.UserRepository) error) error {
	args := m.Called(ctx, u, fn)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// 测试活跃用户可以成功登录
func TestUserService_LoginUser_ActiveStatus_Success(t *testing.T) {
	// Arrange
	// 初始化JWT配置（测试需要）
	testutil.InitTestConfig()

	mockRepo := new(MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	// 创建活跃用户（密码123456的hash）
	activeUser := testutil.CreateTestUser(
		testutil.WithUsername("activeuser"),
		testutil.WithEmail("active@example.com"),
		testutil.WithStatus("active"),
	)
	// 设置密码hash (password123)
	activeUser.SetPassword("password123")

	// Mock期望
	mockRepo.On("GetByUsername", ctx, "activeuser").Return(activeUser, nil)
	mockRepo.On("UpdateLastLogin", ctx, activeUser.ID, mock.Anything).Return(nil)

	// Act
	resp, err := service.LoginUser(ctx, &user2.LoginUserRequest{
		Username: "activeuser",
		Password: "password123",
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "activeuser", resp.User.Username)
	mockRepo.AssertExpectations(t)
}

// 测试未激活用户无法登录
func TestUserService_LoginUser_InactiveStatus_Rejected(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	// 创建未激活用户
	inactiveUser := testutil.CreateTestUser(
		testutil.WithUsername("inactiveuser"),
		testutil.WithEmail("inactive@example.com"),
		testutil.WithStatus("inactive"),
	)
	inactiveUser.SetPassword("password123")

	// Mock期望
	mockRepo.On("GetByUsername", ctx, "inactiveuser").Return(inactiveUser, nil)

	// Act
	resp, err := service.LoginUser(ctx, &user2.LoginUserRequest{
		Username: "inactiveuser",
		Password: "password123",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "账号未激活")
	assert.Contains(t, err.Error(), "验证邮箱")
	mockRepo.AssertExpectations(t)
}

// 测试已封禁用户无法登录
func TestUserService_LoginUser_BannedStatus_Rejected(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	// 创建已封禁用户
	bannedUser := testutil.CreateTestUser(
		testutil.WithUsername("banneduser"),
		testutil.WithEmail("banned@example.com"),
		testutil.WithStatus("banned"),
	)
	bannedUser.SetPassword("password123")

	// Mock期望
	mockRepo.On("GetByUsername", ctx, "banneduser").Return(bannedUser, nil)

	// Act
	resp, err := service.LoginUser(ctx, &user2.LoginUserRequest{
		Username: "banneduser",
		Password: "password123",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "已被封禁")
	assert.Contains(t, err.Error(), "联系管理员")
	mockRepo.AssertExpectations(t)
}

// 测试已删除用户无法登录
func TestUserService_LoginUser_DeletedStatus_Rejected(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	// 创建已删除用户
	deletedUser := testutil.CreateTestUser(
		testutil.WithUsername("deleteduser"),
		testutil.WithEmail("deleted@example.com"),
		testutil.WithStatus("deleted"),
	)
	deletedUser.SetPassword("password123")

	// Mock期望
	mockRepo.On("GetByUsername", ctx, "deleteduser").Return(deletedUser, nil)

	// Act
	resp, err := service.LoginUser(ctx, &user2.LoginUserRequest{
		Username: "deleteduser",
		Password: "password123",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "账号已删除")
	mockRepo.AssertExpectations(t)
}

// 测试密码错误的情况
func TestUserService_LoginUser_WrongPassword(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	// 创建正常用户
	activeUser := testutil.CreateTestUser(
		testutil.WithUsername("normaluser"),
		testutil.WithStatus("active"),
	)
	activeUser.SetPassword("correct_password")

	// Mock期望
	mockRepo.On("GetByUsername", ctx, "normaluser").Return(activeUser, nil)

	// Act - 使用错误的密码
	resp, err := service.LoginUser(ctx, &user2.LoginUserRequest{
		Username: "normaluser",
		Password: "wrong_password",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "密码错误")
	mockRepo.AssertExpectations(t)
}

// 测试用户不存在的情况
func TestUserService_LoginUser_UserNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := user.NewUserService(mockRepo)
	ctx := context.Background()

	// Mock期望 - 返回NotFoundError
	notFoundErr := repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "用户不存在", nil)
	mockRepo.On("GetByUsername", ctx, "nonexistent").Return((*users.User)(nil), notFoundErr)

	// Act
	resp, err := service.LoginUser(ctx, &user2.LoginUserRequest{
		Username: "nonexistent",
		Password: "password123",
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")
	mockRepo.AssertExpectations(t)
}

// 测试表驱动：批量测试所有状态
func TestUserService_LoginUser_AllStatuses(t *testing.T) {
	tests := []struct {
		name          string
		userStatus    users.UserStatus
		expectSuccess bool
		errorContains string
	}{
		{
			name:          "活跃用户可以登录",
			userStatus:    users.UserStatusActive,
			expectSuccess: true,
			errorContains: "",
		},
		{
			name:          "未激活用户被拒绝",
			userStatus:    users.UserStatusInactive,
			expectSuccess: false,
			errorContains: "账号未激活",
		},
		{
			name:          "已封禁用户被拒绝",
			userStatus:    users.UserStatusBanned,
			expectSuccess: false,
			errorContains: "已被封禁",
		},
		{
			name:          "已删除用户被拒绝",
			userStatus:    users.UserStatusDeleted,
			expectSuccess: false,
			errorContains: "账号已删除",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockUserRepository)
			service := user.NewUserService(mockRepo)
			ctx := context.Background()

			// 创建测试用户
			testUser := &users.User{
				ID:       primitive.NewObjectID().Hex(),
				Username: "testuser",
				Email:    "test@example.com",
				Role:     "user",
				Status:   tt.userStatus,
			}
			testUser.SetPassword("password123")

			// Mock期望
			mockRepo.On("GetByUsername", ctx, "testuser").Return(testUser, nil)
			if tt.expectSuccess {
				mockRepo.On("UpdateLastLogin", ctx, testUser.ID, mock.Anything).Return(nil)
			}

			// Act
			resp, err := service.LoginUser(ctx, &user2.LoginUserRequest{
				Username: "testuser",
				Password: "password123",
			})

			// Assert
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.Token)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), tt.errorContains)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
