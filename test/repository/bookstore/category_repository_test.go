package bookstore_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// MockCategoryRepository Mock实现
type MockCategoryRepository struct {
	mock.Mock
	bookstoreRepo.CategoryRepository // 嵌入接口
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *bookstore.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByName(ctx context.Context, name string) (*bookstore.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByParent(ctx context.Context, parentID primitive.ObjectID, limit, offset int) ([]*bookstore.Category, error) {
	args := m.Called(ctx, parentID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context) ([]*bookstore.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.CategoryTree), args.Error(1)
}

func (m *MockCategoryRepository) GetChildren(ctx context.Context, parentID primitive.ObjectID) ([]*bookstore.Category, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockCategoryRepository) UpdateBookCount(ctx context.Context, categoryID primitive.ObjectID, count int64) error {
	args := m.Called(ctx, categoryID, count)
	return args.Error(0)
}

func (m *MockCategoryRepository) BatchUpdateStatus(ctx context.Context, categoryIDs []primitive.ObjectID, isActive bool) error {
	args := m.Called(ctx, categoryIDs, isActive)
	return args.Error(0)
}

// 测试助手函数
func createTestCategory(id primitive.ObjectID, name string, level int) *bookstore.Category {
	return &bookstore.Category{
		ID:          id,
		Name:        name,
		Description: "测试分类",
		Level:       level,
		SortOrder:   100,
		BookCount:   0,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// TestCategoryRepository_Create 测试创建分类
func TestCategoryRepository_Create(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	category := createTestCategory(primitive.NewObjectID(), "测试分类", 0)

	// 设置Mock期望
	mockRepo.On("Create", ctx, category).Return(nil)

	// 执行测试
	err := mockRepo.Create(ctx, category)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetByID_Success 测试根据ID获取分类成功
func TestCategoryRepository_GetByID_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	categoryID := primitive.NewObjectID()
	expectedCategory := createTestCategory(categoryID, "玄幻", 0)

	// 设置Mock期望
	mockRepo.On("GetByID", ctx, categoryID).Return(expectedCategory, nil)

	// 执行测试
	category, err := mockRepo.GetByID(ctx, categoryID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, expectedCategory.ID, category.ID)
	assert.Equal(t, expectedCategory.Name, category.Name)
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetByID_NotFound 测试根据ID获取分类不存在
func TestCategoryRepository_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	categoryID := primitive.NewObjectID()

	// 设置Mock期望：分类不存在
	mockRepo.On("GetByID", ctx, categoryID).Return(nil, errors.New("category not found"))

	// 执行测试
	category, err := mockRepo.GetByID(ctx, categoryID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, category)
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetByName 测试根据名称获取分类
func TestCategoryRepository_GetByName(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	categoryName := "玄幻"
	expectedCategory := createTestCategory(primitive.NewObjectID(), categoryName, 0)

	// 设置Mock期望
	mockRepo.On("GetByName", ctx, categoryName).Return(expectedCategory, nil)

	// 执行测试
	category, err := mockRepo.GetByName(ctx, categoryName)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, categoryName, category.Name)
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetByParent 测试获取子分类列表
func TestCategoryRepository_GetByParent(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	parentID := primitive.NewObjectID()
	childCategories := []*bookstore.Category{
		createTestCategory(primitive.NewObjectID(), "东方玄幻", 1),
		createTestCategory(primitive.NewObjectID(), "异界大陆", 1),
	}

	// 设置Mock期望
	mockRepo.On("GetByParent", ctx, parentID, 10, 0).Return(childCategories, nil)

	// 执行测试
	categories, err := mockRepo.GetByParent(ctx, parentID, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Equal(t, 2, len(categories))
	assert.Equal(t, 1, categories[0].Level)
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetRootCategories 测试获取根分类列表
func TestCategoryRepository_GetRootCategories(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	rootCategories := []*bookstore.Category{
		createTestCategory(primitive.NewObjectID(), "玄幻", 0),
		createTestCategory(primitive.NewObjectID(), "都市", 0),
		createTestCategory(primitive.NewObjectID(), "武侠", 0),
	}

	// 设置Mock期望
	mockRepo.On("GetRootCategories", ctx).Return(rootCategories, nil)

	// 执行测试
	categories, err := mockRepo.GetRootCategories(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Equal(t, 3, len(categories))
	for _, cat := range categories {
		assert.Equal(t, 0, cat.Level)
	}
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetCategoryTree 测试获取分类树
func TestCategoryRepository_GetCategoryTree(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	// 创建分类树结构
	rootCategory := createTestCategory(primitive.NewObjectID(), "玄幻", 0)
	childCategory1 := createTestCategory(primitive.NewObjectID(), "东方玄幻", 1)
	childCategory2 := createTestCategory(primitive.NewObjectID(), "异界大陆", 1)

	categoryTree := []*bookstore.CategoryTree{
		{
			Category: *rootCategory,
			Children: []*bookstore.CategoryTree{
				{Category: *childCategory1},
				{Category: *childCategory2},
			},
		},
	}

	// 设置Mock期望
	mockRepo.On("GetCategoryTree", ctx).Return(categoryTree, nil)

	// 执行测试
	tree, err := mockRepo.GetCategoryTree(ctx)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 1, len(tree))
	assert.Equal(t, 2, len(tree[0].Children))
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetChildren 测试获取直接子分类
func TestCategoryRepository_GetChildren(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	parentID := primitive.NewObjectID()
	children := []*bookstore.Category{
		createTestCategory(primitive.NewObjectID(), "子分类1", 1),
		createTestCategory(primitive.NewObjectID(), "子分类2", 1),
	}

	// 设置Mock期望
	mockRepo.On("GetChildren", ctx, parentID).Return(children, nil)

	// 执行测试
	result, err := mockRepo.GetChildren(ctx, parentID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_UpdateBookCount 测试更新分类书籍数量
func TestCategoryRepository_UpdateBookCount(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	categoryID := primitive.NewObjectID()
	newCount := int64(100)

	// 设置Mock期望
	mockRepo.On("UpdateBookCount", ctx, categoryID, newCount).Return(nil)

	// 执行测试
	err := mockRepo.UpdateBookCount(ctx, categoryID, newCount)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_BatchUpdateStatus 测试批量更新分类状态
func TestCategoryRepository_BatchUpdateStatus(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	ctx := context.Background()

	categoryIDs := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}
	isActive := false

	// 设置Mock期望
	mockRepo.On("BatchUpdateStatus", ctx, categoryIDs, isActive).Return(nil)

	// 执行测试
	err := mockRepo.BatchUpdateStatus(ctx, categoryIDs, isActive)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
