package bookstore_test

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
)

// MockBannerRepository Mock实现
type MockBannerRepository struct {
	mock.Mock
	bookstoreRepo.BannerRepository // 嵌入接口
}

func (m *MockBannerRepository) Create(ctx context.Context, banner *bookstore.Banner) error {
	args := m.Called(ctx, banner)
	return args.Error(0)
}

func (m *MockBannerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Banner, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Banner), args.Error(1)
}

func (m *MockBannerRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBannerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBannerRepository) GetActive(ctx context.Context, limit, offset int) ([]*bookstore.Banner, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Banner), args.Error(1)
}

func (m *MockBannerRepository) GetByTargetType(ctx context.Context, targetType string, limit, offset int) ([]*bookstore.Banner, error) {
	args := m.Called(ctx, targetType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Banner), args.Error(1)
}

func (m *MockBannerRepository) IncrementClickCount(ctx context.Context, bannerID primitive.ObjectID) error {
	args := m.Called(ctx, bannerID)
	return args.Error(0)
}

func (m *MockBannerRepository) GetClickStats(ctx context.Context, bannerID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bannerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBannerRepository) BatchUpdateStatus(ctx context.Context, bannerIDs []primitive.ObjectID, isActive bool) error {
	args := m.Called(ctx, bannerIDs, isActive)
	return args.Error(0)
}

// 测试助手函数
func createTestBanner(id primitive.ObjectID, title, targetType string) *bookstore.Banner {
	now := time.Now()
	return &bookstore.Banner{
		ID:          id,
		Title:       title,
		Description: "测试Banner描述",
		Image:       "https://example.com/banner.jpg",
		Target:      "target_id",
		TargetType:  targetType,
		SortOrder:   100,
		IsActive:    true,
		ClickCount:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// TestBannerRepository_Create 测试创建Banner
func TestBannerRepository_Create(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	banner := createTestBanner(primitive.NewObjectID(), "测试Banner", "book")

	// 设置Mock期望
	mockRepo.On("Create", ctx, banner).Return(nil)

	// 执行测试
	err := mockRepo.Create(ctx, banner)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_GetByID_Success 测试根据ID获取Banner成功
func TestBannerRepository_GetByID_Success(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	bannerID := primitive.NewObjectID()
	expectedBanner := createTestBanner(bannerID, "测试Banner", "book")

	// 设置Mock期望
	mockRepo.On("GetByID", ctx, bannerID).Return(expectedBanner, nil)

	// 执行测试
	banner, err := mockRepo.GetByID(ctx, bannerID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, banner)
	assert.Equal(t, expectedBanner.ID, banner.ID)
	assert.Equal(t, expectedBanner.Title, banner.Title)
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_GetByID_NotFound 测试根据ID获取Banner不存在
func TestBannerRepository_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	bannerID := primitive.NewObjectID()

	// 设置Mock期望：Banner不存在
	mockRepo.On("GetByID", ctx, bannerID).Return(nil, errors.New("banner not found"))

	// 执行测试
	banner, err := mockRepo.GetByID(ctx, bannerID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, banner)
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_GetActive 测试获取活跃Banner列表
func TestBannerRepository_GetActive(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	banners := []*bookstore.Banner{
		createTestBanner(primitive.NewObjectID(), "Banner1", "book"),
		createTestBanner(primitive.NewObjectID(), "Banner2", "category"),
	}

	// 设置Mock期望
	mockRepo.On("GetActive", ctx, 10, 0).Return(banners, nil)

	// 执行测试
	result, err := mockRepo.GetActive(ctx, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	for _, banner := range result {
		assert.True(t, banner.IsActive)
	}
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_GetByTargetType 测试根据目标类型获取Banner
func TestBannerRepository_GetByTargetType(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	targetType := "book"
	banners := []*bookstore.Banner{
		createTestBanner(primitive.NewObjectID(), "书籍Banner1", targetType),
		createTestBanner(primitive.NewObjectID(), "书籍Banner2", targetType),
	}

	// 设置Mock期望
	mockRepo.On("GetByTargetType", ctx, targetType, 10, 0).Return(banners, nil)

	// 执行测试
	result, err := mockRepo.GetByTargetType(ctx, targetType, 10, 0)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	for _, banner := range result {
		assert.Equal(t, targetType, banner.TargetType)
	}
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_IncrementClickCount 测试增加Banner点击次数
func TestBannerRepository_IncrementClickCount(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	bannerID := primitive.NewObjectID()

	// 设置Mock期望
	mockRepo.On("IncrementClickCount", ctx, bannerID).Return(nil)

	// 执行测试
	err := mockRepo.IncrementClickCount(ctx, bannerID)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_GetClickStats 测试获取Banner点击统计
func TestBannerRepository_GetClickStats(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	bannerID := primitive.NewObjectID()
	expectedCount := int64(1234)

	// 设置Mock期望
	mockRepo.On("GetClickStats", ctx, bannerID).Return(expectedCount, nil)

	// 执行测试
	count, err := mockRepo.GetClickStats(ctx, bannerID)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_BatchUpdateStatus 测试批量更新Banner状态
func TestBannerRepository_BatchUpdateStatus(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	bannerIDs := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}
	isActive := false

	// 设置Mock期望
	mockRepo.On("BatchUpdateStatus", ctx, bannerIDs, isActive).Return(nil)

	// 执行测试
	err := mockRepo.BatchUpdateStatus(ctx, bannerIDs, isActive)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_Update 测试更新Banner
func TestBannerRepository_Update(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	bannerID := primitive.NewObjectID()
	updates := map[string]interface{}{
		"title":      "更新的标题",
		"is_active":  false,
		"updated_at": time.Now(),
	}

	// 设置Mock期望
	mockRepo.On("Update", ctx, bannerID, updates).Return(nil)

	// 执行测试
	err := mockRepo.Update(ctx, bannerID, updates)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBannerRepository_Delete 测试删除Banner
func TestBannerRepository_Delete(t *testing.T) {
	mockRepo := new(MockBannerRepository)
	ctx := context.Background()

	bannerID := primitive.NewObjectID()

	// 设置Mock期望
	mockRepo.On("Delete", ctx, bannerID).Return(nil)

	// 执行测试
	err := mockRepo.Delete(ctx, bannerID)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
