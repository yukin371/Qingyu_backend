package test

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/repository/interfaces/infrastructure"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockBookRepository 模拟书籍仓储
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(ctx context.Context, book *bookstore2.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBookRepository) GetByTitle(ctx context.Context, title string) (*bookstore2.Book, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, author, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetByCategory(ctx context.Context, categoryID primitive.ObjectID, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, categoryID, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetByStatus(ctx context.Context, status bookstore2.BookStatus, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, status, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) SearchWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	args := m.Called(ctx, author)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByStatus(ctx context.Context, status bookstore2.BookStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) GetStats(ctx context.Context) (*bookstore2.BookStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*bookstore2.BookStats), args.Error(1)
}

func (m *MockBookRepository) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRepository) IncrementLikeCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRepository) IncrementCommentCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRepository) UpdateRating(ctx context.Context, bookID primitive.ObjectID, rating float64) error {
	args := m.Called(ctx, bookID, rating)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []primitive.ObjectID, status bookstore2.BookStatus) error {
	args := m.Called(ctx, bookIDs, status)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateCategory(ctx context.Context, bookIDs []primitive.ObjectID, categoryIDs []primitive.ObjectID) error {
	args := m.Called(ctx, bookIDs, categoryIDs)
	return args.Error(0)
}

func (m *MockBookRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateFeatured(ctx context.Context, bookIDs []primitive.ObjectID, isFeatured bool) error {
	args := m.Called(ctx, bookIDs, isFeatured)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateRecommended(ctx context.Context, bookIDs []primitive.ObjectID, isRecommended bool) error {
	args := m.Called(ctx, bookIDs, isRecommended)
	return args.Error(0)
}

func (m *MockBookRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByFilter(ctx context.Context, filter *bookstore2.BookFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) GetHot(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetFree(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, keyword, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockBookRepository) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, authorID, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, minPrice, maxPrice, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

// MockCategoryRepository 模拟分类仓储
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *bookstore2.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByName(ctx context.Context, name string) (*bookstore2.Category, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByParent(ctx context.Context, parentID primitive.ObjectID, limit, offset int) ([]*bookstore2.Category, error) {
	args := m.Called(ctx, parentID, limit, offset)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByLevel(ctx context.Context, level int, limit, offset int) ([]*bookstore2.Category, error) {
	args := m.Called(ctx, level, limit, offset)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context) ([]*bookstore2.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context) ([]*bookstore2.CategoryTree, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore2.CategoryTree), args.Error(1)
}

func (m *MockCategoryRepository) CountByParent(ctx context.Context, parentID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCategoryRepository) UpdateBookCount(ctx context.Context, categoryID primitive.ObjectID, count int64) error {
	args := m.Called(ctx, categoryID, count)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetChildren(ctx context.Context, parentID primitive.ObjectID) ([]*bookstore2.Category, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetAncestors(ctx context.Context, categoryID primitive.ObjectID) ([]*bookstore2.Category, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetDescendants(ctx context.Context, categoryID primitive.ObjectID) ([]*bookstore2.Category, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) BatchUpdateStatus(ctx context.Context, categoryIDs []primitive.ObjectID, isActive bool) error {
	args := m.Called(ctx, categoryIDs, isActive)
	return args.Error(0)
}

func (m *MockCategoryRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockCategoryRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCategoryRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockCategoryRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstore2.Category, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

// MockBannerRepository 模拟Banner仓储
type MockBannerRepository struct {
	mock.Mock
}

func (m *MockBannerRepository) Create(ctx context.Context, banner *bookstore2.Banner) error {
	args := m.Called(ctx, banner)
	return args.Error(0)
}

func (m *MockBannerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.Banner, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.Banner), args.Error(1)
}

func (m *MockBannerRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBannerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBannerRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBannerRepository) GetActive(ctx context.Context, limit, offset int) ([]*bookstore2.Banner, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore2.Banner), args.Error(1)
}

func (m *MockBannerRepository) GetByTargetType(ctx context.Context, targetType string, limit, offset int) ([]*bookstore2.Banner, error) {
	args := m.Called(ctx, targetType, limit, offset)
	return args.Get(0).([]*bookstore2.Banner), args.Error(1)
}

func (m *MockBannerRepository) GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*bookstore2.Banner, error) {
	args := m.Called(ctx, startTime, endTime, limit, offset)
	return args.Get(0).([]*bookstore2.Banner), args.Error(1)
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

func (m *MockBannerRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockBannerRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBannerRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockBannerRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstore2.Banner, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore2.Banner), args.Error(1)
}

// MockRankingRepository 模拟榜单仓储
type MockRankingRepository struct {
	mock.Mock
}

func (m *MockRankingRepository) Create(ctx context.Context, item *bookstore2.RankingItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRankingRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.RankingItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockRankingRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRankingRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRankingRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockRankingRepository) GetByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit, offset)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) GetByTypeWithBooks(ctx context.Context, rankingType bookstore2.RankingType, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit, offset)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID, rankingType bookstore2.RankingType, period string) (*bookstore2.RankingItem, error) {
	args := m.Called(ctx, bookID, rankingType, period)
	return args.Get(0).(*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) GetByPeriod(ctx context.Context, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period, limit, offset)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) GetRankingStats(ctx context.Context, rankingType bookstore2.RankingType, period string) (*bookstore2.RankingStats, error) {
	args := m.Called(ctx, rankingType, period)
	return args.Get(0).(*bookstore2.RankingStats), args.Error(1)
}

func (m *MockRankingRepository) CountByType(ctx context.Context, rankingType bookstore2.RankingType, period string) (int64, error) {
	args := m.Called(ctx, rankingType, period)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRankingRepository) GetTopBooks(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) UpsertRankingItem(ctx context.Context, item *bookstore2.RankingItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRankingRepository) BatchUpsertRankingItems(ctx context.Context, items []*bookstore2.RankingItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockRankingRepository) UpdateRankings(ctx context.Context, rankingType bookstore2.RankingType, period string, items []*bookstore2.RankingItem) error {
	args := m.Called(ctx, rankingType, period, items)
	return args.Error(0)
}

func (m *MockRankingRepository) DeleteByPeriod(ctx context.Context, period string) error {
	args := m.Called(ctx, period)
	return args.Error(0)
}

func (m *MockRankingRepository) DeleteByType(ctx context.Context, rankingType bookstore2.RankingType) error {
	args := m.Called(ctx, rankingType)
	return args.Error(0)
}

func (m *MockRankingRepository) DeleteExpiredRankings(ctx context.Context, beforeDate time.Time) error {
	args := m.Called(ctx, beforeDate)
	return args.Error(0)
}

func (m *MockRankingRepository) CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

// 测试用例

func TestBookstoreService_GetBookByID(t *testing.T) {
	// 准备测试数据
	bookID := primitive.NewObjectID()
	expectedBook := &bookstore2.Book{
		ID:     bookID,
		Title:  "测试书籍",
		Author: "测试作者",
		Status: "published",
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望
	mockBookRepo.On("GetByID", mock.Anything, bookID).Return(expectedBook, nil)

	// 创建 MockRankingRepository (使用 ranking_test.go 中的定义)
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, err := service.GetBookByID(context.Background(), bookID.Hex())

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedBook, result)
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_GetBookByID_NotFound(t *testing.T) {
	// 准备测试数据
	bookID := primitive.NewObjectID()

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望 - 返回nil表示未找到
	mockBookRepo.On("GetByID", mock.Anything, bookID).Return((*bookstore2.Book)(nil), nil)

	// 创建 MockRankingRepository (使用 ranking_test.go 中的定义)
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, err := service.GetBookByID(context.Background(), bookID.Hex())

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "book not found")
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_GetBookByID_InvalidID(t *testing.T) {
	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 创建 MockRankingRepository (使用 ranking_test.go 中的定义)
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试 - 使用无效ID
	result, err := service.GetBookByID(context.Background(), "invalid-id")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid book ID")
}

func TestBookstoreService_GetRecommendedBooks(t *testing.T) {
	// 准备测试数据
	expectedBooks := []*bookstore2.Book{
		{
			ID:            primitive.NewObjectID(),
			Title:         "推荐书籍1",
			Author:        "作者1",
			Status:        "published",
			IsRecommended: true,
		},
		{
			ID:            primitive.NewObjectID(),
			Title:         "推荐书籍2",
			Author:        "作者2",
			Status:        "published",
			IsRecommended: true,
		},
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望
	mockBookRepo.On("GetRecommended", mock.Anything, 20, 0).Return(expectedBooks, nil)

	// 创建 MockRankingRepository
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, total, err := service.GetRecommendedBooks(context.Background(), 1, 20)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedBooks, result)
	assert.Equal(t, int64(2), total)
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_IncrementBookView(t *testing.T) {
	// 准备测试数据
	bookID := primitive.NewObjectID()
	book := &bookstore2.Book{
		ID:     bookID,
		Title:  "测试书籍",
		Status: "published",
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望
	mockBookRepo.On("GetByID", mock.Anything, bookID).Return(book, nil)
	mockBookRepo.On("IncrementViewCount", mock.Anything, bookID).Return(nil)

	// 创建 MockRankingRepository
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	err := service.IncrementBookView(context.Background(), bookID.Hex())

	// 验证结果
	assert.NoError(t, err)
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_IncrementBookView_BookNotPublished(t *testing.T) {
	// 准备测试数据
	bookID := primitive.NewObjectID()
	book := &bookstore2.Book{
		ID:     bookID,
		Title:  "测试书籍",
		Status: "draft", // 未发布状态
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望
	mockBookRepo.On("GetByID", mock.Anything, bookID).Return(book, nil)

	// 创建 MockRankingRepository
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	err := service.IncrementBookView(context.Background(), bookID.Hex())

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book not available")
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_GetCategoryByID(t *testing.T) {
	// 准备测试数据
	categoryID := primitive.NewObjectID()
	expectedCategory := &bookstore2.Category{
		ID:       categoryID,
		Name:     "测试分类",
		IsActive: true,
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望
	mockCategoryRepo.On("GetByID", mock.Anything, categoryID).Return(expectedCategory, nil)

	// 创建 MockRankingRepository
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, err := service.GetCategoryByID(context.Background(), categoryID.Hex())

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, result)
	mockCategoryRepo.AssertExpectations(t)
}

func TestBookstoreService_GetCategoryByID_NotActive(t *testing.T) {
	// 准备测试数据
	categoryID := primitive.NewObjectID()
	category := &bookstore2.Category{
		ID:       categoryID,
		Name:     "测试分类",
		IsActive: false, // 未激活
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望
	mockCategoryRepo.On("GetByID", mock.Anything, categoryID).Return(category, nil)

	// 创建 MockRankingRepository
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, err := service.GetCategoryByID(context.Background(), categoryID.Hex())

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "category not available")
	mockCategoryRepo.AssertExpectations(t)
}

func TestBookstoreService_SearchBooks(t *testing.T) {
	// 准备测试数据
	keyword := "测试"
	expectedBooks := []*bookstore2.Book{
		{
			ID:     primitive.NewObjectID(),
			Title:  "测试书籍",
			Author: "测试作者",
			Status: "published",
		},
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)
	mockRankingRepo := new(MockRankingRepository)

	// 设置Mock期望
	mockBookRepo.On("Search", mock.Anything, keyword, 10, 0).Return(expectedBooks, nil)
	publishedStatus := bookstore2.BookStatusPublished
	keywordPtr := keyword
	mockBookRepo.On("CountByFilter", mock.Anything, mock.MatchedBy(func(f *bookstore2.BookFilter) bool {
		return f.Status != nil && *f.Status == publishedStatus && f.Keyword != nil && *f.Keyword == keywordPtr
	})).Return(int64(1), nil)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, total, err := service.SearchBooks(context.Background(), keyword, 1, 10)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedBooks, result)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_SearchBooks_EmptyKeywordAndFilter(t *testing.T) {
	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 创建 MockRankingRepository (使用 ranking_test.go 中的定义)
	mockRankingRepo := new(MockRankingRepository)

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试 - 空关键词
	result, total, err := service.SearchBooks(context.Background(), "", 1, 10)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	assert.Contains(t, err.Error(), "keyword is required")
}

func TestBookstoreService_GetHomepageData(t *testing.T) {
	// 准备测试数据
	expectedBanners := []*bookstore2.Banner{
		{
			ID:       primitive.NewObjectID(),
			Title:    "测试Banner",
			IsActive: true,
		},
	}
	expectedBooks := []*bookstore2.Book{
		{
			ID:            primitive.NewObjectID(),
			Title:         "测试书籍",
			Status:        "published",
			IsRecommended: true,
		},
	}
	expectedCategories := []*bookstore2.Category{
		{
			ID:       primitive.NewObjectID(),
			Name:     "测试分类",
			IsActive: true,
		},
	}
	expectedStats := &bookstore2.BookStats{
		TotalBooks:     100,
		PublishedBooks: 80,
	}

	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望
	mockBannerRepo.On("GetActive", mock.Anything, 5, 0).Return(expectedBanners, nil)
	mockBookRepo.On("GetRecommended", mock.Anything, 10, 0).Return(expectedBooks, nil)
	mockBookRepo.On("GetFeatured", mock.Anything, 10, 0).Return(expectedBooks, nil)
	mockCategoryRepo.On("GetRootCategories", mock.Anything).Return(expectedCategories, nil)
	mockBookRepo.On("GetStats", mock.Anything).Return(expectedStats, nil)

	// 创建 MockRankingRepository
	mockRankingRepo := new(MockRankingRepository)
	// 设置榜单Mock期望
	mockRankingRepo.On("GetByTypeWithBooks", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(([]*bookstore2.RankingItem)(nil), nil).Maybe()

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, err := service.GetHomepageData(context.Background())

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBanners, result.Banners)
	assert.Equal(t, expectedBooks, result.RecommendedBooks)
	assert.Equal(t, expectedBooks, result.FeaturedBooks)
	assert.Equal(t, expectedCategories, result.Categories)
	assert.Equal(t, expectedStats, result.Stats)

	// 验证所有Mock调用
	mockBannerRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestBookstoreService_GetHomepageData_PartialFailure(t *testing.T) {
	// 创建Mock
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	// 设置Mock期望 - Banner获取失败，但其他方法也需要设置（因为可能并发调用）
	mockBannerRepo.On("GetActive", mock.Anything, 5, 0).Return(([]*bookstore2.Banner)(nil), errors.New("banner error"))
	mockBookRepo.On("GetRecommended", mock.Anything, mock.Anything, mock.Anything).Return(([]*bookstore2.Book)(nil), nil).Maybe()
	mockBookRepo.On("GetFeatured", mock.Anything, mock.Anything, mock.Anything).Return(([]*bookstore2.Book)(nil), nil).Maybe()
	mockCategoryRepo.On("GetRootCategories", mock.Anything).Return(([]*bookstore2.Category)(nil), nil).Maybe()
	mockBookRepo.On("GetStats", mock.Anything).Return((*bookstore2.BookStats)(nil), nil).Maybe()

	// 创建 MockRankingRepository
	mockRankingRepo := new(MockRankingRepository)
	// 设置榜单Mock期望（因为可能并发调用）
	mockRankingRepo.On("GetByTypeWithBooks", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(([]*bookstore2.RankingItem)(nil), nil).Maybe()

	// 创建服务
	service := bookstoreService.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 执行测试
	result, err := service.GetHomepageData(context.Background())

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get banners")
	mockBannerRepo.AssertExpectations(t)
}
