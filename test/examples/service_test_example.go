package examples

import (
	"context"
	"errors"
	"testing"

	"Qingyu_backend/models/users"
	"Qingyu_backend/test/fixtures"
	"Qingyu_backend/test/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/*
本文件提供Service层测试的完整示例
展示Table-Driven测试模式、Mock使用、测试数据工厂等最佳实践
*/

// ============ Mock定义 ============

// MockUserRepository Mock用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *users.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
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

// ============ 示例Service ============

// UserService 示例用户服务
type UserService struct {
	userRepo *MockUserRepository
}

func NewUserService(userRepo *MockUserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *users.User) error {
	// 检查邮箱是否已存在
	existing, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && err.Error() != "not found" {
		return err
	}
	if existing != nil {
		return errors.New("email already exists")
	}

	// 创建用户
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*users.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新用户
	return s.userRepo.Update(ctx, id, updates)
}

// ============ Table-Driven测试示例 ============

func TestuserserviceCreateuserTabledriven(t *testing.T) {
	// 定义测试用例表
	tests := []struct {
		name    string                    // 测试用例名称
		user    *users.User               // 输入用户
		setup   func(*MockUserRepository) // Mock设置
		wantErr bool                      // 是否期望错误
		errMsg  string                    // 错误消息
	}{
		{
			name: "成功创建用户",
			user: testutil.CreateTestUser(
				testutil.WithUsername("newuser"),
				testutil.WithEmail("new@test.com"),
			),
			setup: func(m *MockUserRepository) {
				// Mock: 邮箱不存在
				m.On("GetByEmail", mock.Anything, "new@test.com").
					Return(nil, errors.New("not found"))
				// Mock: 创建成功
				m.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "邮箱已存在_返回错误",
			user: testutil.CreateTestUser(
				testutil.WithEmail("existing@test.com"),
			),
			setup: func(m *MockUserRepository) {
				// Mock: 邮箱已存在
				existingUser := testutil.CreateTestUser(
					testutil.WithEmail("existing@test.com"),
				)
				m.On("GetByEmail", mock.Anything, "existing@test.com").
					Return(existingUser, nil)
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
		{
			name: "数据库查询错误_返回错误",
			user: testutil.CreateTestUser(),
			setup: func(m *MockUserRepository) {
				// Mock: 数据库错误
				m.On("GetByEmail", mock.Anything, mock.Anything).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name: "创建失败_返回错误",
			user: testutil.CreateTestUser(),
			setup: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, mock.Anything).
					Return(nil, errors.New("not found"))
				m.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).
					Return(errors.New("create failed"))
			},
			wantErr: true,
			errMsg:  "create failed",
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: 准备Mock和Service
			mockRepo := new(MockUserRepository)
			tt.setup(mockRepo)
			service := NewUserService(mockRepo)
			ctx := testutil.CreateTestContext()

			// Act: 执行被测试的方法
			err := service.CreateUser(ctx, tt.user)

			// Assert: 验证结果
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			// 验证Mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// ============ 使用工厂模式的测试示例 ============

func TestuserserviceGetuserWithfactory(t *testing.T) {
	// 使用工厂创建测试数据
	userFactory := fixtures.NewUserFactory()

	tests := []struct {
		name    string
		userID  string
		setup   func(*MockUserRepository, *users.User)
		wantErr bool
	}{
		{
			name:   "成功获取用户",
			userID: "user123",
			setup: func(m *MockUserRepository, user *users.User) {
				m.On("GetByID", mock.Anything, "user123").
					Return(user, nil)
			},
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: "nonexistent",
			setup: func(m *MockUserRepository, user *users.User) {
				m.On("GetByID", mock.Anything, "nonexistent").
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockUserRepository)
			expectedUser := userFactory.Create()
			tt.setup(mockRepo, expectedUser)
			service := NewUserService(mockRepo)
			ctx := testutil.CreateTestContext()

			// Act
			user, err := service.GetUser(ctx, tt.userID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				testutil.AssertUserEqual(t, expectedUser, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// ============ 子测试示例 ============

func TestuserserviceUpdateuser(t *testing.T) {
	// 准备通用数据
	userFactory := fixtures.NewUserFactory()
	user := userFactory.Create()
	ctx := testutil.CreateTestContext()

	t.Run("成功更新用户名", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		mockRepo.On("Update", mock.Anything, user.ID, mock.MatchedBy(func(updates map[string]interface{}) bool {
			return updates["username"] == "newname"
		})).Return(nil)

		service := NewUserService(mockRepo)
		updates := map[string]interface{}{"username": "newname"}

		// Act
		err := service.UpdateUser(ctx, user.ID, updates)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在_返回错误", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		mockRepo.On("GetByID", mock.Anything, "nonexistent").
			Return(nil, errors.New("not found"))

		service := NewUserService(mockRepo)
		updates := map[string]interface{}{"username": "newname"}

		// Act
		err := service.UpdateUser(ctx, "nonexistent", updates)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})
}

// ============ 使用testutil助手的示例 ============

func TestuserserviceWithhelpers(t *testing.T) {
	// 创建测试上下文
	ctx := testutil.CreateTestContext()

	// 创建测试用户
	user1 := testutil.CreateTestUser()
	user2 := testutil.CreateTestUser(
		testutil.WithUsername("custom"),
		testutil.WithRole("admin"),
	)

	// 创建Mock
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByID", mock.Anything, user1.ID).Return(user1, nil)
	mockRepo.On("GetByID", mock.Anything, user2.ID).Return(user2, nil)

	service := NewUserService(mockRepo)

	// 测试获取用户1
	result1, err := service.GetUser(ctx, user1.ID)
	testutil.AssertNoErrorWithMessage(t, err, "获取用户1失败")
	testutil.AssertUserEqual(t, user1, result1)

	// 测试获取用户2
	result2, err := service.GetUser(ctx, user2.ID)
	testutil.AssertNoErrorWithMessage(t, err, "获取用户2失败")
	assert.Contains(t, result2.Roles, "admin")

	mockRepo.AssertExpectations(t)
}

// ============ 性能基准测试示例 ============

func BenchmarkuserserviceCreateuser(b *testing.B) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByEmail", mock.Anything, mock.Anything).
		Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	service := NewUserService(mockRepo)
	ctx := testutil.CreateTestContext()
	user := testutil.CreateTestUser()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.CreateUser(ctx, user)
	}
}

/*
运行本示例测试:
  go test -v ./test/examples/service_test_example.go

运行特定测试:
  go test -run TestUserService_CreateUser_TableDriven

运行基准测试:
  go test -bench=BenchmarkUserService_CreateUser
*/
