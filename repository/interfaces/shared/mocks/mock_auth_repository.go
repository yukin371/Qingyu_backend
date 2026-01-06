package mocks

import (
	authModel "Qingyu_backend/models/auth"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockAuthRepository AuthRepository 的 Mock 实现
type MockAuthRepository struct {
	mock.Mock
}

// CreateRole 创建角色
func (m *MockAuthRepository) CreateRole(ctx context.Context, role *authModel.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

// GetRole 获取角色
func (m *MockAuthRepository) GetRole(ctx context.Context, roleID string) (*authModel.Role, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

// GetRoleByName 根据名称获取角色
func (m *MockAuthRepository) GetRoleByName(ctx context.Context, name string) (*authModel.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

// UpdateRole 更新角色
func (m *MockAuthRepository) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error {
	args := m.Called(ctx, roleID, updates)
	return args.Error(0)
}

// DeleteRole 删除角色
func (m *MockAuthRepository) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

// ListRoles 列出所有角色
func (m *MockAuthRepository) ListRoles(ctx context.Context) ([]*authModel.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Role), args.Error(1)
}

// AssignUserRole 分配用户角色
func (m *MockAuthRepository) AssignUserRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// RemoveUserRole 移除用户角色
func (m *MockAuthRepository) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// GetUserRoles 获取用户角色
func (m *MockAuthRepository) GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Role), args.Error(1)
}

// HasUserRole 检查用户是否有某角色
func (m *MockAuthRepository) HasUserRole(ctx context.Context, userID, roleID string) (bool, error) {
	args := m.Called(ctx, userID, roleID)
	return args.Bool(0), args.Error(1)
}

// GetRolePermissions 获取角色权限
func (m *MockAuthRepository) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// GetUserPermissions 获取用户权限
func (m *MockAuthRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// Health 健康检查
func (m *MockAuthRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
