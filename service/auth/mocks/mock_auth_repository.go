package mocks

import (
	"context"

	authModel "Qingyu_backend/models/auth"

	"github.com/stretchr/testify/mock"
)

// MockAuthRepository 是认证模块测试使用的角色仓储替身。
type MockAuthRepository struct {
	mock.Mock
}

// NewMockAuthRepository 创建认证模块测试替身。
func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{}
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
