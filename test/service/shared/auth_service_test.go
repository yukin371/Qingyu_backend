package shared

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUser 模拟用户模型
type MockUser struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	Phone        string             `bson:"phone" json:"phone"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Salt         string             `bson:"salt" json:"-"`
	Status       string             `bson:"status" json:"status"`
	Roles        []string           `bson:"roles" json:"roles"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockRole 模拟角色模型
type MockRole struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	DisplayName string             `bson:"display_name" json:"displayName"`
	Permissions []string           `bson:"permissions" json:"permissions"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *MockUser) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*MockUser, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockUser), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*MockUser, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockUser), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*MockUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockUser), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRoleRepository 模拟角色仓储
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*MockRole, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockRole), args.Error(1)
}

func (m *MockRoleRepository) GetByNames(ctx context.Context, names []string) ([]*MockRole, error) {
	args := m.Called(ctx, names)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockRole), args.Error(1)
}

// MockAuthService 模拟认证服务
type MockAuthService struct {
	userRepo *MockUserRepository
	roleRepo *MockRoleRepository
}

func NewMockAuthService(userRepo *MockUserRepository, roleRepo *MockRoleRepository) *MockAuthService {
	return &MockAuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// Register 用户注册
func (s *MockAuthService) Register(ctx context.Context, username, email, password string) (*MockUser, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingEmail, _ := s.userRepo.GetByEmail(ctx, email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// 创建新用户
	user := &MockUser{
		ID:           primitive.NewObjectID(),
		Username:     username,
		Email:        email,
		PasswordHash: hashPassword(password),
		Salt:         generateSalt(),
		Status:       "active",
		Roles:        []string{"user"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *MockAuthService) Login(ctx context.Context, username, password string) (*MockUser, string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if user.Status != "active" {
		return nil, "", errors.New("user account is not active")
	}

	if !verifyPassword(password, user.PasswordHash, user.Salt) {
		return nil, "", errors.New("invalid credentials")
	}

	// 生成JWT Token
	token := generateJWTToken(user)

	return user, token, nil
}

// GetUserPermissions 获取用户权限
func (s *MockAuthService) GetUserPermissions(ctx context.Context, userID primitive.ObjectID) ([]string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles, err := s.roleRepo.GetByNames(ctx, user.Roles)
	if err != nil {
		return nil, err
	}

	permissionSet := make(map[string]bool)
	for _, role := range roles {
		for _, permission := range role.Permissions {
			permissionSet[permission] = true
		}
	}

	permissions := make([]string, 0, len(permissionSet))
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// 辅助函数
func hashPassword(password string) string {
	// 模拟密码哈希
	return "hashed_" + password
}

func generateSalt() string {
	// 模拟盐值生成
	return "salt_" + time.Now().Format("20060102150405")
}

func verifyPassword(password, hash, salt string) bool {
	// 模拟密码验证
	return hash == "hashed_"+password
}

func generateJWTToken(user *MockUser) string {
	// 模拟JWT Token生成
	return "jwt_token_" + user.Username
}

// 测试用例

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	roleRepo := new(MockRoleRepository)
	authService := NewMockAuthService(userRepo, roleRepo)

	ctx := context.Background()
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// Mock 设置
	userRepo.On("GetByUsername", ctx, username).Return(nil, errors.New("user not found"))
	userRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("email not found"))
	userRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockUser")).Return(nil)

	// 执行测试
	user, err := authService.Register(ctx, username, email, password)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "active", user.Status)
	assert.Contains(t, user.Roles, "user")

	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_UsernameExists(t *testing.T) {
	userRepo := new(MockUserRepository)
	roleRepo := new(MockRoleRepository)
	authService := NewMockAuthService(userRepo, roleRepo)

	ctx := context.Background()
	username := "existinguser"
	email := "test@example.com"
	password := "password123"

	existingUser := &MockUser{
		ID:       primitive.NewObjectID(),
		Username: username,
		Email:    "existing@example.com",
	}

	// Mock 设置
	userRepo.On("GetByUsername", ctx, username).Return(existingUser, nil)

	// 执行测试
	user, err := authService.Register(ctx, username, email, password)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "username already exists", err.Error())

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	roleRepo := new(MockRoleRepository)
	authService := NewMockAuthService(userRepo, roleRepo)

	ctx := context.Background()
	username := "testuser"
	password := "password123"

	existingUser := &MockUser{
		ID:           primitive.NewObjectID(),
		Username:     username,
		Email:        "test@example.com",
		PasswordHash: "hashed_" + password,
		Salt:         "salt_123",
		Status:       "active",
		Roles:        []string{"user"},
	}

	// Mock 设置
	userRepo.On("GetByUsername", ctx, username).Return(existingUser, nil)

	// 执行测试
	user, token, err := authService.Login(ctx, username, password)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, "jwt_token_"+username, token)

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	userRepo := new(MockUserRepository)
	roleRepo := new(MockRoleRepository)
	authService := NewMockAuthService(userRepo, roleRepo)

	ctx := context.Background()
	username := "testuser"
	password := "wrongpassword"

	// Mock 设置
	userRepo.On("GetByUsername", ctx, username).Return(nil, errors.New("user not found"))

	// 执行测试
	user, token, err := authService.Login(ctx, username, password)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())

	userRepo.AssertExpectations(t)
}

func TestAuthService_GetUserPermissions_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	roleRepo := new(MockRoleRepository)
	authService := NewMockAuthService(userRepo, roleRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	user := &MockUser{
		ID:       userID,
		Username: "testuser",
		Roles:    []string{"admin", "editor"},
	}

	roles := []*MockRole{
		{
			Name:        "admin",
			Permissions: []string{"user:read", "user:write", "content:read", "content:write"},
		},
		{
			Name:        "editor",
			Permissions: []string{"content:read", "content:write", "content:publish"},
		},
	}

	// Mock 设置
	userRepo.On("GetByID", ctx, userID).Return(user, nil)
	roleRepo.On("GetByNames", ctx, []string{"admin", "editor"}).Return(roles, nil)

	// 执行测试
	permissions, err := authService.GetUserPermissions(ctx, userID)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, permissions)
	assert.Contains(t, permissions, "user:read")
	assert.Contains(t, permissions, "user:write")
	assert.Contains(t, permissions, "content:read")
	assert.Contains(t, permissions, "content:write")
	assert.Contains(t, permissions, "content:publish")

	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}