package mocks

import (
	"context"

	usersModel "Qingyu_backend/models/users"
	"Qingyu_backend/repository/interfaces/infrastructure"
	userRepo "Qingyu_backend/repository/interfaces/user"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository Mock实现
type MockUserRepository struct {
	mock.Mock
}

// 实现CRUDRepository接口
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

func (m *MockUserRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*usersModel.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	if len(args) > 0 && args.Get(0) != nil {
		return args.Get(0).(int64), args.Error(1)
	}
	return 0, args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// 用户特定方法
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

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*usersModel.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usersModel.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	args := m.Called(ctx, phone)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id string, ip string) error {
	args := m.Called(ctx, id, ip)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
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
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*usersModel.User), args.Get(1).(int64), args.Error(2)
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

// Transaction - 现在可以正确使用UserRepository类型
func (m *MockUserRepository) Transaction(ctx context.Context, user *usersModel.User, fn func(context.Context, userRepo.UserRepository) error) error {
	args := m.Called(ctx, user, fn)
	return args.Error(0)
}

func (m *MockUserRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
