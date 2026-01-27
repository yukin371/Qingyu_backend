package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	authModel "Qingyu_backend/models/auth"
	middlewareAuth "Qingyu_backend/internal/middleware/auth"
)

// MockAuthRepository Mock仓储实现
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

func (m *MockAuthRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ========== NewPermissionService 测试 ==========

func TestNewPermissionService(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	assert.NotNil(t, service)
}

// ========== CheckPermission 测试 ==========

func TestCheckPermission_WithWildcard(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	// Mock返回通配符权限
	mockRepo.On("GetUserPermissions", mock.Anything, "user1").Return([]string{"*:*"}, nil)

	service := NewPermissionService(mockRepo, nil, logger)

	allowed, err := service.CheckPermission(context.Background(), "user1", "any:permission")

	assert.NoError(t, err)
	assert.True(t, allowed)
	mockRepo.AssertExpectations(t)
}

func TestCheckPermission_ExactMatch(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	// Mock返回具体权限
	mockRepo.On("GetUserPermissions", mock.Anything, "user1").Return([]string{
		"user:read",
		"book:read",
		"book:write",
	}, nil)

	service := NewPermissionService(mockRepo, nil, logger)

	// 测试精确匹配
	allowed, err := service.CheckPermission(context.Background(), "user1", "book:read")

	assert.NoError(t, err)
	assert.True(t, allowed)
	mockRepo.AssertExpectations(t)
}

func TestCheckPermission_NoPermission(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	// Mock返回空权限
	mockRepo.On("GetUserPermissions", mock.Anything, "user1").Return([]string{}, nil)

	service := NewPermissionService(mockRepo, nil, logger)

	allowed, err := service.CheckPermission(context.Background(), "user1", "book:read")

	assert.NoError(t, err)
	assert.False(t, allowed)
	mockRepo.AssertExpectations(t)
}

// ========== HasRole 测试 ==========

func TestHasRole_HasRole(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	// Mock返回角色列表
	mockRepo.On("GetUserRoles", mock.Anything, "user1").Return([]*authModel.Role{
		{Name: "admin"},
		{Name: "author"},
	}, nil)

	service := NewPermissionService(mockRepo, nil, logger)

	hasRole, err := service.HasRole(context.Background(), "user1", "admin")

	assert.NoError(t, err)
	assert.True(t, hasRole)
	mockRepo.AssertExpectations(t)
}

func TestHasRole_NoRole(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	// Mock返回角色列表
	mockRepo.On("GetUserRoles", mock.Anything, "user1").Return([]*authModel.Role{
		{Name: "author"},
	}, nil)

	service := NewPermissionService(mockRepo, nil, logger)

	hasRole, err := service.HasRole(context.Background(), "user1", "admin")

	assert.NoError(t, err)
	assert.False(t, hasRole)
	mockRepo.AssertExpectations(t)
}

// ========== RBAC集成测试 ==========

func TestSetChecker(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	// 创建checker
	checker, _ := middlewareAuth.NewRBACChecker(nil)

	// 设置checker
	service.SetChecker(checker)

	// 验证设置成功（通过调用LoadPermissionsToChecker验证）
	mockRepo.On("ListRoles", mock.Anything).Return([]*authModel.Role{
		{Name: "admin", Permissions: []string{"user.read", "book.write"}},
	}, nil)

	err := service.LoadPermissionsToChecker(context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLoadPermissionsToChecker(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	// 创建checker
	checker, _ := middlewareAuth.NewRBACChecker(nil)
	rbacChecker := checker.(*middlewareAuth.RBACChecker)
	service.SetChecker(rbacChecker)

	// Mock返回角色列表
	mockRepo.On("ListRoles", mock.Anything).Return([]*authModel.Role{
		{
			Name: "admin",
			Permissions: []string{
				"user.read",
				"book.write",
				"book.delete",
			},
		},
		{
			Name: "reader",
			Permissions: []string{
				"book.read",
				"chapter.read",
			},
		},
	}, nil)

	// 加载权限
	err := service.LoadPermissionsToChecker(context.Background())

	assert.NoError(t, err)

	// 验证checker中的权限
	adminPerms := rbacChecker.GetRolePermissions("admin")
	assert.Contains(t, adminPerms, "user:read")
	assert.Contains(t, adminPerms, "book:write")

	readerPerms := rbacChecker.GetRolePermissions("reader")
	assert.Contains(t, readerPerms, "book:read")

	mockRepo.AssertExpectations(t)
}

func TestLoadUserRolesToChecker(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	// 创建checker并设置权限
	checker, _ := middlewareAuth.NewRBACChecker(nil)
	rbacChecker := checker.(*middlewareAuth.RBACChecker)
	rbacChecker.GrantPermission("admin", "*:*")
	service.SetChecker(rbacChecker)

	// Mock返回用户角色
	mockRepo.On("GetUserRoles", mock.Anything, "user1").Return([]*authModel.Role{
		{Name: "admin"},
		{Name: "author"},
	}, nil)

	// 加载用户角色
	err := service.LoadUserRolesToChecker(context.Background(), "user1")

	assert.NoError(t, err)

	// 验证用户角色已设置到checker
	roles := rbacChecker.GetUserRoles("user1")
	assert.Contains(t, roles, "admin")
	assert.Contains(t, roles, "author")

	mockRepo.AssertExpectations(t)
}

func TestReloadAllFromDatabase(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	// 创建checker
	checker, _ := middlewareAuth.NewRBACChecker(nil)
	rbacChecker := checker.(*middlewareAuth.RBACChecker)
	service.SetChecker(rbacChecker)

	// Mock返回角色列表
	mockRepo.On("ListRoles", mock.Anything).Return([]*authModel.Role{
		{Name: "admin", Permissions: []string{"*:*"}},
	}, nil)

	// 重新加载
	err := service.ReloadAllFromDatabase(context.Background())

	assert.NoError(t, err)

	// 验证权限已加载
	perms := rbacChecker.GetRolePermissions("admin")
	assert.Contains(t, perms, "*:*")

	mockRepo.AssertExpectations(t)
}

func TestSetChecker_NotSet(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	// 不设置checker，直接调用LoadPermissionsToChecker
	err := service.LoadPermissionsToChecker(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RBACChecker未设置")
}

// ========== convertPermissions 测试 ==========

func TestConvertPermissions(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	input := []string{
		"user.read",
		"book.write",
		"admin.manage",
	}

	expected := []string{
		"user:read",
		"book:write",
		"admin:manage",
	}

	result := service.(*PermissionServiceImpl).convertPermissions(input)

	assert.Equal(t, expected, result)
}

// ========== 缓存测试 ==========

func TestGetRoleFromCache(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	logger := zap.NewNop()

	service := NewPermissionService(mockRepo, nil, logger)

	// 缓存中没有角色
	role, exists := service.(*PermissionServiceImpl).getRoleFromCache("admin")

	assert.Nil(t, role)
	assert.False(t, exists)
}
