package user

import (
	"context"
	"testing"

	authModel "Qingyu_backend/models/auth"
	usersModel "Qingyu_backend/models/users"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	repoInterfaces "Qingyu_backend/repository/interfaces/user"
	user2 "Qingyu_backend/service/interfaces/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =========================
// Mock Repository实现
// =========================

// MockUserRepository Mock用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *usersModel.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*usersModel.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filter base.Filter) ([]*usersModel.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*usersModel.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*usersModel.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id string, ip string) error {
	args := m.Called(ctx, id, ip)
	return args.Error(0)
}

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*usersModel.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	args := m.Called(ctx, phone)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateStatus(ctx context.Context, id string, status usersModel.UserStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserRepository) GetActiveUsers(ctx context.Context, limit int64) ([]*usersModel.User, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersByRole(ctx context.Context, role string, limit int64) ([]*usersModel.User, error) {
	args := m.Called(ctx, role, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

func (m *MockUserRepository) SetPhoneVerified(ctx context.Context, id string, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

func (m *MockUserRepository) BatchUpdateStatus(ctx context.Context, ids []string, status usersModel.UserStatus) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockUserRepository) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockUserRepository) FindWithFilter(ctx context.Context, filter *usersModel.UserFilter) ([]*usersModel.User, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(1)
	}
	return args.Get(0).([]*usersModel.User), args.Get(1).(int64), args.Error(1)
}

func (m *MockUserRepository) SearchUsers(ctx context.Context, keyword string, limit int) ([]*usersModel.User, error) {
	args := m.Called(ctx, keyword, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	args := m.Called(ctx, role)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountByStatus(ctx context.Context, status usersModel.UserStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Transaction(ctx context.Context, user *usersModel.User, fn func(ctx context.Context, repo repoInterfaces.UserRepository) error) error {
	args := m.Called(ctx, user, fn)
	return args.Error(0)
}

func (m *MockUserRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockAuthRepository Mock认证仓储
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateRole(ctx context.Context, role *authModel.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockAuthRepository) GetRole(ctx context.Context, roleID string) (*authModel.Role, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

func (m *MockAuthRepository) GetRoleByName(ctx context.Context, name string) (*authModel.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

func (m *MockAuthRepository) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error {
	args := m.Called(ctx, roleID, updates)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockAuthRepository) ListRoles(ctx context.Context) ([]*authModel.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Role), args.Error(1)
}

func (m *MockAuthRepository) AssignUserRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthRepository) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Role), args.Error(1)
}

func (m *MockAuthRepository) HasUserRole(ctx context.Context, userID, roleID string) (bool, error) {
	args := m.Called(ctx, userID, roleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// =========================
// 测试辅助函数
// =========================

// setupUserService 创建测试用的UserService实例
func setupUserService() (*UserServiceImpl, *MockUserRepository, *MockAuthRepository) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)

	service := NewUserService(mockUserRepo, mockAuthRepo)

	return service.(*UserServiceImpl), mockUserRepo, mockAuthRepo
}

// =========================
// 创建用户相关测试
// =========================

// TestUserService_CreateUser_Success 测试创建用户成功
func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByUsername", ctx, req.Username).Return(false, nil)
	mockUserRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)
	mockUserRepo.On("Create", ctx, mock.MatchedBy(func(u *usersModel.User) bool {
		return u.Username == req.Username && u.Email == req.Email
	})).Return(nil)

	// Act
	resp, err := service.CreateUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Username, resp.User.Username)
	assert.Equal(t, req.Email, resp.User.Email)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_CreateUser_DuplicateUsername 测试创建用户-用户名已存在
func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.CreateUserRequest{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByUsername", ctx, req.Username).Return(true, nil)

	// Act
	resp, err := service.CreateUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户名已存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_CreateUser_DuplicateEmail 测试创建用户-邮箱已存在
func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.CreateUserRequest{
		Username: "testuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByUsername", ctx, req.Username).Return(false, nil)
	mockUserRepo.On("ExistsByEmail", ctx, req.Email).Return(true, nil)

	// Act
	resp, err := service.CreateUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱已存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_CreateUser_ValidationFailed 测试创建用户-验证失败
func TestUserService_CreateUser_ValidationFailed(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *user2.CreateUserRequest
		wantErr string
	}{
		{
			name: "空用户名",
			req: &user2.CreateUserRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: "请求数据验证失败",
		},
		{
			name: "空邮箱",
			req: &user2.CreateUserRequest{
				Username: "testuser",
				Email:    "",
				Password: "password123",
			},
			wantErr: "请求数据验证失败",
		},
		{
			name: "空密码",
			req: &user2.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: "请求数据验证失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			resp, err := service.CreateUser(ctx, tt.req)

			// Assert
			require.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// =========================
// 获取用户相关测试
// =========================

// TestUserService_GetUser_Success 测试获取用户成功
func TestUserService_GetUser_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	expectedUser := &usersModel.User{
		Username: "testuser",
		Email:    "test@example.com",
	}
	expectedUser.ID = primitive.NewObjectID()
	userID := expectedUser.ID.Hex()

	mockUserRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	// Act
	req := &user2.GetUserRequest{ID: userID}
	resp, err := service.GetUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedUser, resp.User)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetUser_NotFound 测试获取用户-不存在
func TestUserService_GetUser_NotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	req := &user2.GetUserRequest{ID: userID}
	resp, err := service.GetUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetUser_EmptyID 测试获取用户-空ID
func TestUserService_GetUser_EmptyID(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	// Act
	req := &user2.GetUserRequest{ID: ""}
	resp, err := service.GetUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

// =========================
// 更新用户相关测试
// =========================

// TestUserService_UpdateUser_Success 测试更新用户成功
func TestUserService_UpdateUser_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	updatedUser := &usersModel.User{
		Username: "updateduser",
		Email:    "updated@example.com",
	}
	updatedUser.ID = primitive.NewObjectID()
	userID := updatedUser.ID.Hex()

	updates := map[string]interface{}{
		"username": "updateduser",
	}

	mockUserRepo.On("Exists", ctx, userID).Return(true, nil)
	mockUserRepo.On("Update", ctx, userID, updates).Return(nil)
	mockUserRepo.On("GetByID", ctx, userID).Return(updatedUser, nil)

	// Act
	req := &user2.UpdateUserRequest{
		ID:      userID,
		Updates: updates,
	}
	resp, err := service.UpdateUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, *updatedUser, resp.User)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_UpdateUser_NotFound 测试更新用户-不存在
func TestUserService_UpdateUser_NotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()

	mockUserRepo.On("Exists", ctx, userID).Return(false, nil)

	// Act
	req := &user2.UpdateUserRequest{
		ID:      userID,
		Updates: map[string]interface{}{"username": "newname"},
	}
	resp, err := service.UpdateUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_UpdateUser_EmptyUpdates 测试更新用户-空更新数据
func TestUserService_UpdateUser_EmptyUpdates(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	// Act
	req := &user2.UpdateUserRequest{
		ID:      primitive.NewObjectID().Hex(),
		Updates: map[string]interface{}{},
	}
	resp, err := service.UpdateUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "更新数据不能为空")
}

// =========================
// 删除用户相关测试
// =========================

// TestUserService_DeleteUser_Success 测试删除用户成功
func TestUserService_DeleteUser_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()

	mockUserRepo.On("Exists", ctx, userID).Return(true, nil)
	mockUserRepo.On("Delete", ctx, userID).Return(nil)

	// Act
	req := &user2.DeleteUserRequest{ID: userID}
	resp, err := service.DeleteUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Deleted)
	assert.False(t, resp.DeletedAt.IsZero())

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_DeleteUser_NotFound 测试删除用户-不存在
func TestUserService_DeleteUser_NotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()

	mockUserRepo.On("Exists", ctx, userID).Return(false, nil)

	// Act
	req := &user2.DeleteUserRequest{ID: userID}
	resp, err := service.DeleteUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// =========================
// 列出用户相关测试
// =========================

// TestUserService_ListUsers_Success 测试列出用户成功
func TestUserService_ListUsers_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	expectedUsers := []*usersModel.User{
		func() *usersModel.User {
			u := &usersModel.User{Username: "user1"}
			u.ID = primitive.NewObjectID()
			return u
		}(),
		func() *usersModel.User {
			u := &usersModel.User{Username: "user2"}
			u.ID = primitive.NewObjectID()
			return u
		}(),
	}

	req := &user2.ListUsersRequest{
		Page:     1,
		PageSize: 20,
	}

	mockUserRepo.On("List", ctx, mock.Anything).Return(expectedUsers, nil)
	mockUserRepo.On("Count", ctx, mock.Anything).Return(int64(2), nil)

	// Act
	resp, err := service.ListUsers(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.Len(t, resp.Users, 2)
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.PageSize)
	assert.Equal(t, 1, resp.TotalPages)

	mockUserRepo.AssertExpectations(t)
}

// =========================
// 用户注册相关测试
// =========================

// TestUserService_RegisterUser_Success 测试用户注册成功
func TestUserService_RegisterUser_Success(t *testing.T) {
	t.Skip("需要JWT配置，集成测试中运行")

	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.RegisterUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByUsername", ctx, req.Username).Return(false, nil)
	mockUserRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)
	mockUserRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Act
	resp, err := service.RegisterUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, req.Username, resp.User.Username)
	assert.Equal(t, req.Email, resp.User.Email)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_RegisterUser_DuplicateUsername 测试用户注册-用户名已存在
func TestUserService_RegisterUser_DuplicateUsername(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.RegisterUserRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
	}

	mockUserRepo.On("ExistsByUsername", ctx, req.Username).Return(true, nil)

	// Act
	resp, err := service.RegisterUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户名已存在")

	mockUserRepo.AssertExpectations(t)
}

// =========================
// 用户登录相关测试
// =========================

// TestUserService_LoginUser_Success 测试用户登录成功
func TestUserService_LoginUser_Success(t *testing.T) {
	t.Skip("需要JWT配置，集成测试中运行")

	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.LoginUserRequest{
		Username: "testuser",
		Password: "password123",
		ClientIP: "127.0.0.1",
	}

	user := &usersModel.User{
		Username: "testuser",
		Email:    "test@example.com",
		Status:   usersModel.UserStatusActive,
	}
	user.ID, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	// 设置密码
	err := user.SetPassword(req.Password)
	require.NoError(t, err)

	mockUserRepo.On("GetByUsername", ctx, req.Username).Return(user, nil)
	mockUserRepo.On("UpdateLastLogin", ctx, user.ID.Hex(), req.ClientIP).Return(nil)

	// Act
	resp, err := service.LoginUser(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user, resp.User)
	assert.NotEmpty(t, resp.Token)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_LoginUser_UserNotFound 测试用户登录-用户不存在
func TestUserService_LoginUser_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.LoginUserRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	mockUserRepo.On("GetByUsername", ctx, req.Username).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.LoginUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_LoginUser_WrongPassword 测试用户登录-密码错误
func TestUserService_LoginUser_WrongPassword(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.LoginUserRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	user := &usersModel.User{
		Username: "testuser",
		Status:   usersModel.UserStatusActive,
	}
	user.ID, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	// 设置正确的密码
	err := user.SetPassword("correctpassword")
	require.NoError(t, err)

	mockUserRepo.On("GetByUsername", ctx, req.Username).Return(user, nil)

	// Act
	resp, err := service.LoginUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "密码错误")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_LoginUser_AccountInactive 测试用户登录-账号未激活
func TestUserService_LoginUser_AccountInactive(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.LoginUserRequest{
		Username: "testuser",
		Password: "password123",
	}

	user := &usersModel.User{
		Username: "testuser",
		Status:   usersModel.UserStatusInactive,
	}
	user.ID, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	user.SetPassword(req.Password)

	mockUserRepo.On("GetByUsername", ctx, req.Username).Return(user, nil)

	// Act
	resp, err := service.LoginUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "账号未激活")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_LoginUser_AccountBanned 测试用户登录-账号被封禁
func TestUserService_LoginUser_AccountBanned(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.LoginUserRequest{
		Username: "testuser",
		Password: "password123",
	}

	user := &usersModel.User{
		Username: "testuser",
		Status:   usersModel.UserStatusBanned,
	}
	user.ID, _ = primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	user.SetPassword(req.Password)

	mockUserRepo.On("GetByUsername", ctx, req.Username).Return(user, nil)

	// Act
	resp, err := service.LoginUser(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "账号已被封禁")

	mockUserRepo.AssertExpectations(t)
}

// =========================
// 密码更新相关测试
// =========================

// TestUserService_UpdatePassword_Success 测试更新密码成功
func TestUserService_UpdatePassword_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()
	oldPassword := "oldpassword"
	newPassword := "newpassword123"

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(userID)

	user.SetPassword(oldPassword)

	req := &user2.UpdatePasswordRequest{
		ID:          userID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("UpdatePassword", ctx, userID, mock.MatchedBy(func(hashed string) bool {
		return len(hashed) > 0 && hashed != oldPassword && hashed != newPassword
	})).Return(nil)

	// Act
	resp, err := service.UpdatePassword(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Updated)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_UpdatePassword_WrongOldPassword 测试更新密码-旧密码错误
func TestUserService_UpdatePassword_WrongOldPassword(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(userID)

	user.SetPassword("correctpassword")

	req := &user2.UpdatePasswordRequest{
		ID:          userID,
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Act
	resp, err := service.UpdatePassword(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "旧密码错误")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_UpdatePassword_UserNotFound 测试更新密码-用户不存在
func TestUserService_UpdatePassword_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()

	req := &user2.UpdatePasswordRequest{
		ID:          userID,
		OldPassword: "oldpassword",
		NewPassword: "newpassword123",
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.UpdatePassword(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// =========================
// 角色分配相关测试
// =========================

// TestUserService_AssignRole_Success 测试分配角色成功
func TestUserService_AssignRole_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.AssignRoleRequest{
		UserID: "user123",
		RoleID: "role123",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)
	role := &authModel.Role{ID: req.RoleID}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("GetRole", ctx, req.RoleID).Return(role, nil)
	mockAuthRepo.On("AssignUserRole", ctx, req.UserID, req.RoleID).Return(nil)

	// Act
	resp, err := service.AssignRole(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.True(t, resp.Assigned)

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// TestUserService_AssignRole_UserNotFound 测试分配角色-用户不存在
func TestUserService_AssignRole_UserNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	req := &user2.AssignRoleRequest{
		UserID: "nonexistent",
		RoleID: "role123",
	}

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(nil, repoInterfaces.NewUserRepositoryError(repoInterfaces.ErrorTypeNotFound, "user not found", nil))

	// Act
	resp, err := service.AssignRole(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_AssignRole_RoleNotFound 测试分配角色-角色不存在
func TestUserService_AssignRole_RoleNotFound(t *testing.T) {
	// Arrange
	service, mockUserRepo, mockAuthRepo := setupUserService()
	ctx := context.Background()

	req := &user2.AssignRoleRequest{
		UserID: "user123",
		RoleID: "nonexistent",
	}

	user := &usersModel.User{}
	user.ID, _ = primitive.ObjectIDFromHex(req.UserID)

	mockUserRepo.On("GetByID", ctx, req.UserID).Return(user, nil)
	mockAuthRepo.On("GetRole", ctx, req.RoleID).Return(nil, assert.AnError)

	// Act
	resp, err := service.AssignRole(ctx, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "角色不存在")

	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

// =========================
// BaseService接口测试
// =========================

// TestUserService_Health_Success 测试健康检查成功
func TestUserService_Health_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	mockUserRepo.On("Health", ctx).Return(nil)

	// Act
	err := service.Health(ctx)

	// Assert
	require.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_GetServiceName 测试获取服务名称
func TestUserService_GetServiceName(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()

	// Act
	name := service.GetServiceName()

	// Assert
	assert.Equal(t, "UserService", name)
}

// TestUserService_GetVersion 测试获取服务版本
func TestUserService_GetVersion(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()

	// Act
	version := service.GetVersion()

	// Assert
	assert.Equal(t, "1.0.0", version)
}

// TestUserService_Initialize_Success 测试初始化成功
func TestUserService_Initialize_Success(t *testing.T) {
	// Arrange
	service, mockUserRepo, _ := setupUserService()
	ctx := context.Background()

	mockUserRepo.On("Health", ctx).Return(nil)

	// Act
	err := service.Initialize(ctx)

	// Assert
	require.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
}

// TestUserService_Close_Success 测试关闭成功
func TestUserService_Close_Success(t *testing.T) {
	// Arrange
	service, _, _ := setupUserService()
	ctx := context.Background()

	// Act
	err := service.Close(ctx)

	// Assert
	require.NoError(t, err)
}
