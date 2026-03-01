package admin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/bookstore"
)

// MockCategoryAdminRepository Mock仓储
type MockCategoryAdminRepository struct {
	mock.Mock
}

func (m *MockCategoryAdminRepository) Create(ctx context.Context, category *bookstore.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryAdminRepository) GetByID(ctx context.Context, id string) (*bookstore.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockCategoryAdminRepository) Update(ctx context.Context, category *bookstore.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryAdminRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryAdminRepository) List(ctx context.Context, filter interface{}, opts ...interface{}) ([]*bookstore.Category, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockCategoryAdminRepository) GetTree(ctx context.Context) ([]*bookstore.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockCategoryAdminRepository) UpdateBookCount(ctx context.Context, categoryID string, count int64) error {
	args := m.Called(ctx, categoryID, count)
	return args.Error(0)
}

func (m *MockCategoryAdminRepository) BatchUpdateStatus(ctx context.Context, categoryIDs []string, isActive bool) error {
	args := m.Called(ctx, categoryIDs, isActive)
	return args.Error(0)
}

func (m *MockCategoryAdminRepository) GetDescendantIDs(ctx context.Context, categoryID string) ([]string, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCategoryAdminRepository) HasChildren(ctx context.Context, categoryID string) (bool, error) {
	args := m.Called(ctx, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCategoryAdminRepository) NameExistsAtLevel(ctx context.Context, parentID *string, name string, excludeID string) (bool, error) {
	args := m.Called(ctx, parentID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

// 测试创建分类成功
func TestCategoryAdminService_CreateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	req := &CreateCategoryRequest{
		Name:      "玄幻",
		ParentID:  nil,
		SortOrder: 1,
	}

	mockRepo.On("NameExistsAtLevel", ctx, (*string)(nil), "玄幻", "").Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*bookstore.Category")).Return(nil)

	category, err := service.CreateCategory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "玄幻", category.Name)
	assert.Equal(t, 0, category.Level)
	mockRepo.AssertExpectations(t)
}

// 测试创建分类 - 名称重复
func TestCategoryAdminService_CreateCategory_DuplicateName(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	req := &CreateCategoryRequest{
		Name:      "玄幻",
		ParentID:  nil,
		SortOrder: 1,
	}

	mockRepo.On("NameExistsAtLevel", ctx, (*string)(nil), "玄幻", "").Return(true, nil)

	category, err := service.CreateCategory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "分类名称已存在")
	mockRepo.AssertExpectations(t)
}
