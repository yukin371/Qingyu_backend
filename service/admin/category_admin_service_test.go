package admin

import (
	"context"
	"errors"
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

// 测试创建分类 - 带父分类
func TestCategoryAdminService_CreateCategory_WithParent(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	parentID := "parent123"
	req := &CreateCategoryRequest{
		Name:      "玄幻",
		ParentID:  &parentID,
		SortOrder: 1,
	}

	parentCategory := &bookstore.Category{
		ID:     parentID,
		Name:   "小说",
		Level:  0,
		ParentID: nil,
	}

	mockRepo.On("NameExistsAtLevel", ctx, &parentID, "玄幻", "").Return(false, nil)
	mockRepo.On("GetByID", ctx, parentID).Return(parentCategory, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*bookstore.Category")).Return(nil)

	category, err := service.CreateCategory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "玄幻", category.Name)
	assert.Equal(t, 1, category.Level)
	assert.Equal(t, &parentID, category.ParentID)
	mockRepo.AssertExpectations(t)
}

// 测试创建分类 - 父分类不存在
func TestCategoryAdminService_CreateCategory_ParentNotFound(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	parentID := "nonexistent"
	req := &CreateCategoryRequest{
		Name:      "玄幻",
		ParentID:  &parentID,
		SortOrder: 1,
	}

	mockRepo.On("NameExistsAtLevel", ctx, &parentID, "玄幻", "").Return(false, nil)
	mockRepo.On("GetByID", ctx, parentID).Return(nil, errors.New("not found"))

	category, err := service.CreateCategory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "父分类不存在")
	mockRepo.AssertExpectations(t)
}

// 测试创建分类 - 层级超过限制
func TestCategoryAdminService_CreateCategory_LevelExceeded(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	parentID := "parent123"
	req := &CreateCategoryRequest{
		Name:      "玄幻",
		ParentID:  &parentID,
		SortOrder: 1,
	}

	parentCategory := &bookstore.Category{
		ID:       parentID,
		Name:     "子分类",
		Level:    2,
		ParentID: nil,
	}

	mockRepo.On("NameExistsAtLevel", ctx, &parentID, "玄幻", "").Return(false, nil)
	mockRepo.On("GetByID", ctx, parentID).Return(parentCategory, nil)

	category, err := service.CreateCategory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "分类层级超过限制")
	mockRepo.AssertExpectations(t)
}

// 测试更新分类成功
func TestCategoryAdminService_UpdateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"
	newName := "玄幻小说"
	newDesc := "玄幻类小说"
	newOrder := 5

	req := &UpdateCategoryRequest{
		Name:        &newName,
		Description: &newDesc,
		SortOrder:   &newOrder,
	}

	existingCategory := &bookstore.Category{
		ID:          id,
		Name:        "玄幻",
		Description: "旧描述",
		SortOrder:   1,
		ParentID:    nil,
	}

	mockRepo.On("GetByID", ctx, id).Return(existingCategory, nil)
	mockRepo.On("NameExistsAtLevel", ctx, existingCategory.ParentID, newName, id).Return(false, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*bookstore.Category")).Return(nil)

	category, err := service.UpdateCategory(ctx, id, req)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, newName, category.Name)
	assert.Equal(t, newDesc, category.Description)
	assert.Equal(t, newOrder, category.SortOrder)
	mockRepo.AssertExpectations(t)
}

// 测试更新分类 - 分类不存在
func TestCategoryAdminService_UpdateCategory_NotFound(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "nonexistent"
	newName := "新名称"

	req := &UpdateCategoryRequest{
		Name: &newName,
	}

	mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	category, err := service.UpdateCategory(ctx, id, req)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "分类不存在")
	mockRepo.AssertExpectations(t)
}

// 测试更新分类 - 名称重复
func TestCategoryAdminService_UpdateCategory_DuplicateName(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"
	newName := "武侠"

	req := &UpdateCategoryRequest{
		Name: &newName,
	}

	existingCategory := &bookstore.Category{
		ID:       id,
		Name:     "玄幻",
		ParentID: nil,
	}

	mockRepo.On("GetByID", ctx, id).Return(existingCategory, nil)
	mockRepo.On("NameExistsAtLevel", ctx, existingCategory.ParentID, newName, id).Return(true, nil)

	category, err := service.UpdateCategory(ctx, id, req)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "分类名称已存在")
	mockRepo.AssertExpectations(t)
}

// 测试删除分类成功
func TestCategoryAdminService_DeleteCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"

	category := &bookstore.Category{
		ID:        id,
		Name:      "玄幻",
		BookCount: 0,
	}

	mockRepo.On("HasChildren", ctx, id).Return(false, nil)
	mockRepo.On("GetByID", ctx, id).Return(category, nil)
	mockRepo.On("Delete", ctx, id).Return(nil)

	err := service.DeleteCategory(ctx, id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// 测试删除分类 - 有子分类
func TestCategoryAdminService_DeleteCategory_HasChildren(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"

	mockRepo.On("HasChildren", ctx, id).Return(true, nil)

	err := service.DeleteCategory(ctx, id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "存在子分类")
	mockRepo.AssertExpectations(t)
}

// 测试删除分类 - 有关联作品
func TestCategoryAdminService_DeleteCategory_HasBooks(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"

	category := &bookstore.Category{
		ID:        id,
		Name:      "玄幻",
		BookCount: 10,
	}

	mockRepo.On("HasChildren", ctx, id).Return(false, nil)
	mockRepo.On("GetByID", ctx, id).Return(category, nil)

	err := service.DeleteCategory(ctx, id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "存在关联作品")
	mockRepo.AssertExpectations(t)
}

// 测试删除分类 - 分类不存在
func TestCategoryAdminService_DeleteCategory_NotFound(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "nonexistent"

	mockRepo.On("HasChildren", ctx, id).Return(false, nil)
	mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	err := service.DeleteCategory(ctx, id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "分类不存在")
	mockRepo.AssertExpectations(t)
}

// 测试获取分类树成功
func TestCategoryAdminService_GetCategoryTree_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()

	categories := []*bookstore.Category{
		{ID: "1", Name: "小说", ParentID: nil, Level: 0},
		{ID: "2", Name: "玄幻", ParentID: strPtr("1"), Level: 1},
		{ID: "3", Name: "武侠", ParentID: strPtr("1"), Level: 1},
	}

	mockRepo.On("List", ctx, mock.Anything).Return(categories, nil)

	tree, err := service.GetCategoryTree(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Len(t, tree, 1)
	assert.Equal(t, "小说", tree[0].Name)
	assert.Len(t, tree[0].Children, 2)
	mockRepo.AssertExpectations(t)
}

// 测试获取分类列表成功
func TestCategoryAdminService_GetCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()

	categories := []*bookstore.Category{
		{ID: "1", Name: "玄幻", ParentID: nil},
		{ID: "2", Name: "武侠", ParentID: nil},
	}

	mockRepo.On("List", ctx, mock.AnythingOfType("map[string]interface {}")).Return(categories, nil)

	filter := &CategoryFilter{}
	result, err := service.GetCategories(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// 测试获取分类列表 - 带筛选条件
func TestCategoryAdminService_GetCategories_WithFilter(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()

	parentID := "parent123"
	isActive := true

	categories := []*bookstore.Category{
		{ID: "1", Name: "玄幻", ParentID: &parentID},
	}

	mockRepo.On("List", ctx, mock.AnythingOfType("map[string]interface {}")).Return(categories, nil)

	filter := &CategoryFilter{
		ParentID: &parentID,
		IsActive: &isActive,
	}
	result, err := service.GetCategories(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

// 测试获取分类详情成功
func TestCategoryAdminService_GetCategoryByID_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"

	category := &bookstore.Category{
		ID:   id,
		Name: "玄幻",
	}

	mockRepo.On("GetByID", ctx, id).Return(category, nil)

	result, err := service.GetCategoryByID(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "玄幻", result.Name)
	mockRepo.AssertExpectations(t)
}

// 测试移动分类成功
func TestCategoryAdminService_MoveCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"
	newParentID := "parent456"

	req := &MoveCategoryRequest{
		ParentID: &newParentID,
	}

	category := &bookstore.Category{
		ID:     id,
		Name:   "玄幻",
		Level:  0,
		ParentID: nil,
	}

	newParent := &bookstore.Category{
		ID:       newParentID,
		Name:     "小说",
		Level:    0,
		ParentID: nil,
	}

	mockRepo.On("GetByID", ctx, id).Return(category, nil)
	mockRepo.On("GetByID", ctx, newParentID).Return(newParent, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*bookstore.Category")).Return(nil)

	err := service.MoveCategory(ctx, id, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// 测试移动分类 - 分类不存在
func TestCategoryAdminService_MoveCategory_CategoryNotFound(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "nonexistent"
	newParentID := "parent456"

	req := &MoveCategoryRequest{
		ParentID: &newParentID,
	}

	mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	err := service.MoveCategory(ctx, id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "分类不存在")
	mockRepo.AssertExpectations(t)
}

// 测试移动分类 - 循环引用
func TestCategoryAdminService_MoveCategory_CircularReference(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"
	newParentID := "child456"

	req := &MoveCategoryRequest{
		ParentID: &newParentID,
	}

	category := &bookstore.Category{
		ID:       id,
		Name:     "玄幻",
		Level:    0,
		ParentID: nil,
	}

	childCategory := &bookstore.Category{
		ID:       newParentID,
		Name:     "子分类",
		Level:    1,
		ParentID: &id,
	}

	mockRepo.On("GetByID", ctx, id).Return(category, nil)
	mockRepo.On("GetByID", ctx, newParentID).Return(childCategory, nil)
	mockRepo.On("GetByID", ctx, newParentID).Return(childCategory, nil)

	err := service.MoveCategory(ctx, id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "循环引用")
	mockRepo.AssertExpectations(t)
}

// 测试移动分类 - 层级超过限制
func TestCategoryAdminService_MoveCategory_LevelExceeded(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"
	newParentID := "parent456"

	req := &MoveCategoryRequest{
		ParentID: &newParentID,
	}

	category := &bookstore.Category{
		ID:     id,
		Name:   "玄幻",
		Level:  0,
		ParentID: nil,
	}

	newParent := &bookstore.Category{
		ID:       newParentID,
		Name:     "子分类",
		Level:    2,
		ParentID: nil,
	}

	mockRepo.On("GetByID", ctx, id).Return(category, nil)
	mockRepo.On("GetByID", ctx, newParentID).Return(newParent, nil)

	err := service.MoveCategory(ctx, id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "分类层级超过限制")
	mockRepo.AssertExpectations(t)
}

// 测试移动分类到自己下面
func TestCategoryAdminService_MoveCategory_MoveToSelf(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"

	req := &MoveCategoryRequest{
		ParentID: &id,
	}

	category := &bookstore.Category{
		ID:       id,
		Name:     "玄幻",
		Level:    0,
		ParentID: nil,
	}

	mockRepo.On("GetByID", ctx, id).Return(category, nil)

	err := service.MoveCategory(ctx, id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不能将分类移动到自己下面")
	mockRepo.AssertExpectations(t)
}

// 测试调整排序成功
func TestCategoryAdminService_SortCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "category123"
	newSortOrder := 10

	category := &bookstore.Category{
		ID:        id,
		Name:      "玄幻",
		SortOrder: 1,
	}

	mockRepo.On("GetByID", ctx, id).Return(category, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*bookstore.Category")).Return(nil)

	err := service.SortCategory(ctx, id, newSortOrder)

	assert.NoError(t, err)
	assert.Equal(t, newSortOrder, category.SortOrder)
	mockRepo.AssertExpectations(t)
}

// 测试调整排序 - 分类不存在
func TestCategoryAdminService_SortCategory_NotFound(t *testing.T) {
	mockRepo := new(MockCategoryAdminRepository)
	service := NewCategoryAdminService(mockRepo)

	ctx := context.Background()
	id := "nonexistent"

	mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	err := service.SortCategory(ctx, id, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "分类不存在")
	mockRepo.AssertExpectations(t)
}

// 辅助函数
func strPtr(s string) *string {
	return &s
}
