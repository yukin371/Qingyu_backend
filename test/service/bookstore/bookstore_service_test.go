package bookstore_test

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// ============ Mock Repositories ============

// MockBookRepository Mock书籍Repository - 使用接口嵌入
type MockBookRepository struct {
	mock.Mock
	BookstoreRepo.BookRepository // 嵌入接口，避免实现所有方法
}

// 只实现测试中实际使用的方法
func (m *MockBookRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetByCategory(ctx context.Context, categoryID primitive.ObjectID, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, categoryID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) CountByCategory(ctx context.Context, categoryID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) SearchWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Book), args.Error(1)
}

func (m *MockBookRepository) GetStats(ctx context.Context) (*bookstore2.BookStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.BookStats), args.Error(1)
}

func (m *MockBookRepository) IncrementViewCount(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRepository) CountByFilter(ctx context.Context, filter *bookstore2.BookFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// MockCategoryRepository Mock分类Repository
type MockCategoryRepository struct {
	mock.Mock
	BookstoreRepo.CategoryRepository // 嵌入接口
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context) ([]*bookstore2.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context) ([]*bookstore2.CategoryTree, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.CategoryTree), args.Error(1)
}

// MockBannerRepository Mock Banner Repository
type MockBannerRepository struct {
	mock.Mock
	BookstoreRepo.BannerRepository // 嵌入接口
}

func (m *MockBannerRepository) GetActive(ctx context.Context, limit, offset int) ([]*bookstore2.Banner, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.Banner), args.Error(1)
}

func (m *MockBannerRepository) IncrementClickCount(ctx context.Context, bannerID primitive.ObjectID) error {
	args := m.Called(ctx, bannerID)
	return args.Error(0)
}

func (m *MockBannerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore2.Banner, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore2.Banner), args.Error(1)
}

// MockRankingRepository Mock排行榜Repository
type MockRankingRepository struct {
	mock.Mock
	BookstoreRepo.RankingRepository // 嵌入接口
}

func (m *MockRankingRepository) GetByType(ctx context.Context, rankingType bookstore2.RankingType, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) GetByTypeWithBooks(ctx context.Context, rankingType bookstore2.RankingType, period string, limit, offset int) ([]*bookstore2.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore2.RankingItem), args.Error(1)
}

// ============ Test Helpers ============

func createTestBook(id, title string, status bookstore2.BookStatus) *bookstore2.Book {
	objID, _ := primitive.ObjectIDFromHex(id)
	return &bookstore2.Book{
		ID:     objID,
		Title:  title,
		Status: status,
		Author: "测试作者",
	}
}

func createTestCategory(id, name string) *bookstore2.Category {
	objID, _ := primitive.ObjectIDFromHex(id)
	return &bookstore2.Category{
		ID:   objID,
		Name: name,
	}
}

func createTestBanner(id, title string) *bookstore2.Banner {
	objID, _ := primitive.ObjectIDFromHex(id)
	return &bookstore2.Banner{
		ID:       objID,
		Title:    title,
		IsActive: true,
	}
}

// ============ GetBookByID Tests ============

func TestBookstoreService_GetBookByID_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)
	mockRankingRepo := new(MockRankingRepository)

	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		mockCategoryRepo,
		mockBannerRepo,
		mockRankingRepo,
	)

	ctx := context.Background()
	bookID := "507f1f77bcf86cd799439011"
	objID, _ := primitive.ObjectIDFromHex(bookID)

	expectedBook := createTestBook(bookID, "测试书籍", bookstore2.BookStatusPublished)

	mockBookRepo.On("GetByID", ctx, objID).Return(expectedBook, nil)

	// Act
	result, err := service.GetBookByID(ctx, bookID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBook.Title, result.Title)
	assert.Equal(t, expectedBook.Status, result.Status)
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_GetBookByID_InvalidID(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	invalidID := "invalid-id"

	// Act
	result, err := service.GetBookByID(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid book ID")
}

func TestBookstoreService_GetBookByID_NotFound(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	bookID := "507f1f77bcf86cd799439011"
	objID, _ := primitive.ObjectIDFromHex(bookID)

	mockBookRepo.On("GetByID", ctx, objID).Return(nil, nil)

	// Act
	result, err := service.GetBookByID(ctx, bookID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "book not found")
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_GetBookByID_NotPublished(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	bookID := "507f1f77bcf86cd799439011"
	objID, _ := primitive.ObjectIDFromHex(bookID)

	unpublishedBook := createTestBook(bookID, "未发布书籍", bookstore2.BookStatusDraft)
	mockBookRepo.On("GetByID", ctx, objID).Return(unpublishedBook, nil)

	// Act
	result, err := service.GetBookByID(ctx, bookID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "book not available")
	mockBookRepo.AssertExpectations(t)
}

// ============ GetBooksByCategory Tests ============

func TestBookstoreService_GetBooksByCategory_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	categoryID := "507f1f77bcf86cd799439011"
	objID, _ := primitive.ObjectIDFromHex(categoryID)
	page := 1
	pageSize := 10

	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439012", "书籍1", bookstore2.BookStatusPublished),
		createTestBook("507f1f77bcf86cd799439013", "书籍2", bookstore2.BookStatusPublished),
	}

	mockBookRepo.On("GetByCategory", ctx, objID, pageSize, 0).Return(expectedBooks, nil)
	mockBookRepo.On("CountByCategory", ctx, objID).Return(int64(2), nil)

	// Act
	result, total, err := service.GetBooksByCategory(ctx, categoryID, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_GetBooksByCategory_EmptyResult(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	categoryID := "507f1f77bcf86cd799439011"
	objID, _ := primitive.ObjectIDFromHex(categoryID)

	mockBookRepo.On("GetByCategory", ctx, objID, 10, 0).Return([]*bookstore2.Book{}, nil)
	mockBookRepo.On("CountByCategory", ctx, objID).Return(int64(0), nil)

	// Act
	result, total, err := service.GetBooksByCategory(ctx, categoryID, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	mockBookRepo.AssertExpectations(t)
}

// ============ GetRecommendedBooks Tests ============

func TestBookstoreService_GetRecommendedBooks_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439011", "推荐书籍1", bookstore2.BookStatusPublished),
		createTestBook("507f1f77bcf86cd799439012", "推荐书籍2", bookstore2.BookStatusPublished),
	}

	mockBookRepo.On("GetRecommended", ctx, 10, 0).Return(expectedBooks, nil)

	// Act
	result, err := service.GetRecommendedBooks(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockBookRepo.AssertExpectations(t)
}

// ============ GetFeaturedBooks Tests ============

func TestBookstoreService_GetFeaturedBooks_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439011", "精选书籍1", bookstore2.BookStatusPublished),
	}

	mockBookRepo.On("GetFeatured", ctx, 10, 0).Return(expectedBooks, nil)

	// Act
	result, err := service.GetFeaturedBooks(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockBookRepo.AssertExpectations(t)
}

// ============ GetHotBooks Tests ============

func TestBookstoreService_GetHotBooks_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439011", "热门书籍1", bookstore2.BookStatusPublished),
		createTestBook("507f1f77bcf86cd799439012", "热门书籍2", bookstore2.BookStatusPublished),
		createTestBook("507f1f77bcf86cd799439013", "热门书籍3", bookstore2.BookStatusPublished),
	}

	mockBookRepo.On("GetHotBooks", ctx, 10, 0).Return(expectedBooks, nil)

	// Act
	result, err := service.GetHotBooks(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "热门书籍1", result[0].Title)
	mockBookRepo.AssertExpectations(t)
}

// ============ GetNewReleases Tests ============

func TestBookstoreService_GetNewReleases_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439011", "新书1", bookstore2.BookStatusPublished),
	}

	mockBookRepo.On("GetNewReleases", ctx, 10, 0).Return(expectedBooks, nil)

	// Act
	result, err := service.GetNewReleases(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockBookRepo.AssertExpectations(t)
}

// ============ GetFreeBooks Tests ============

func TestBookstoreService_GetFreeBooks_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439011", "免费书籍1", bookstore2.BookStatusPublished),
		createTestBook("507f1f77bcf86cd799439012", "免费书籍2", bookstore2.BookStatusPublished),
	}

	mockBookRepo.On("GetFreeBooks", ctx, 10, 0).Return(expectedBooks, nil)

	// Act
	result, err := service.GetFreeBooks(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockBookRepo.AssertExpectations(t)
}

// ============ SearchBooks Tests ============

func TestBookstoreService_SearchBooks_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	keyword := "测试"
	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439011", "测试书籍", bookstore2.BookStatusPublished),
	}

	mockBookRepo.On("Search", ctx, keyword, 10, 0).Return(expectedBooks, nil)
	mockBookRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookFilter")).Return(int64(1), nil)

	// Act
	result, total, err := service.SearchBooks(ctx, keyword, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	assert.Contains(t, result[0].Title, "测试")
	mockBookRepo.AssertExpectations(t)
}

func TestBookstoreService_SearchBooks_NoResults(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	keyword := "不存在的关键词"

	mockBookRepo.On("Search", ctx, keyword, 10, 0).Return([]*bookstore2.Book{}, nil)
	mockBookRepo.On("CountByFilter", ctx, mock.AnythingOfType("*bookstore.BookFilter")).Return(int64(0), nil)

	// Act
	result, total, err := service.SearchBooks(ctx, keyword, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	mockBookRepo.AssertExpectations(t)
}

// ============ GetCategoryTree Tests ============

func TestBookstoreService_GetCategoryTree_Success(t *testing.T) {
	// Arrange
	mockCategoryRepo := new(MockCategoryRepository)
	service := bookstoreService.NewBookstoreService(
		new(MockBookRepository),
		mockCategoryRepo,
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	expectedTree := []*bookstore2.CategoryTree{
		{
			Category: *createTestCategory("507f1f77bcf86cd799439011", "分类1"),
			Children: []*bookstore2.CategoryTree{},
		},
	}

	mockCategoryRepo.On("GetCategoryTree", ctx).Return(expectedTree, nil)

	// Act
	result, err := service.GetCategoryTree(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "分类1", result[0].Category.Name)
	mockCategoryRepo.AssertExpectations(t)
}

// ============ GetActiveBanners Tests ============

func TestBookstoreService_GetActiveBanners_Success(t *testing.T) {
	// Arrange
	mockBannerRepo := new(MockBannerRepository)
	service := bookstoreService.NewBookstoreService(
		new(MockBookRepository),
		new(MockCategoryRepository),
		mockBannerRepo,
		new(MockRankingRepository),
	)

	ctx := context.Background()
	expectedBanners := []*bookstore2.Banner{
		createTestBanner("507f1f77bcf86cd799439011", "Banner 1"),
		createTestBanner("507f1f77bcf86cd799439012", "Banner 2"),
	}

	mockBannerRepo.On("GetActive", ctx, 5, 0).Return(expectedBanners, nil)

	// Act
	result, err := service.GetActiveBanners(ctx, 5)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockBannerRepo.AssertExpectations(t)
}

// ============ IncrementBannerClick Tests ============

func TestBookstoreService_IncrementBannerClick_Success(t *testing.T) {
	// Arrange
	mockBannerRepo := new(MockBannerRepository)
	service := bookstoreService.NewBookstoreService(
		new(MockBookRepository),
		new(MockCategoryRepository),
		mockBannerRepo,
		new(MockRankingRepository),
	)

	ctx := context.Background()
	bannerID := "507f1f77bcf86cd799439011"
	objID, _ := primitive.ObjectIDFromHex(bannerID)

	// 需要先Mock GetByID，因为服务会先检查Banner是否存在
	banner := createTestBanner(bannerID, "Test Banner")
	mockBannerRepo.On("GetByID", ctx, objID).Return(banner, nil)
	mockBannerRepo.On("IncrementClickCount", ctx, objID).Return(nil)

	// Act
	err := service.IncrementBannerClick(ctx, bannerID)

	// Assert
	assert.NoError(t, err)
	mockBannerRepo.AssertExpectations(t)
}

func TestBookstoreService_IncrementBannerClick_InvalidID(t *testing.T) {
	// Arrange
	service := bookstoreService.NewBookstoreService(
		new(MockBookRepository),
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	invalidID := "invalid-id"

	// Act
	err := service.IncrementBannerClick(ctx, invalidID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid banner ID")
}

// ============ GetHomepageData Tests ============

func TestBookstoreService_GetHomepageData_Success(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)
	mockRankingRepo := new(MockRankingRepository)

	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		mockCategoryRepo,
		mockBannerRepo,
		mockRankingRepo,
	)

	ctx := context.Background()

	// Mock所有需要的数据
	expectedBanners := []*bookstore2.Banner{
		createTestBanner("507f1f77bcf86cd799439011", "Banner 1"),
	}
	expectedBooks := []*bookstore2.Book{
		createTestBook("507f1f77bcf86cd799439011", "推荐书籍", bookstore2.BookStatusPublished),
	}
	expectedStats := &bookstore2.BookStats{
		TotalBooks: 100,
	}
	expectedCategories := []*bookstore2.Category{
		createTestCategory("507f1f77bcf86cd799439011", "分类1"),
	}
	expectedRankings := []*bookstore2.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
	}

	mockBannerRepo.On("GetActive", ctx, 5, 0).Return(expectedBanners, nil)
	mockBookRepo.On("GetRecommended", ctx, 10, 0).Return(expectedBooks, nil)
	mockBookRepo.On("GetFeatured", ctx, 10, 0).Return(expectedBooks, nil)
	mockCategoryRepo.On("GetRootCategories", ctx).Return(expectedCategories, nil)
	mockBookRepo.On("GetStats", ctx).Return(expectedStats, nil)

	// GetHomepageData使用GetByTypeWithBooks而不是GetByType
	mockRankingRepo.On("GetByTypeWithBooks", ctx, bookstore2.RankingTypeRealtime, mock.Anything, 10, 0).Return(expectedRankings, nil)
	mockRankingRepo.On("GetByTypeWithBooks", ctx, bookstore2.RankingTypeWeekly, mock.Anything, 10, 0).Return(expectedRankings, nil)
	mockRankingRepo.On("GetByTypeWithBooks", ctx, bookstore2.RankingTypeMonthly, mock.Anything, 10, 0).Return(expectedRankings, nil)

	// Act
	result, err := service.GetHomepageData(ctx)

	// Assert
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Banners, 1)
	assert.Len(t, result.RecommendedBooks, 1)
	assert.Len(t, result.Categories, 1)
	assert.NotNil(t, result.Stats)

	mockBannerRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
	mockRankingRepo.AssertExpectations(t)
}

// ============ Table-Driven Tests ============

func TestBookstoreService_Pagination(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		offset   int
	}{
		{"第1页", 1, 10, 0},
		{"第2页", 2, 10, 10},
		{"第3页", 3, 20, 40},
		{"大页面", 1, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockBookRepo := new(MockBookRepository)
			service := bookstoreService.NewBookstoreService(
				mockBookRepo,
				new(MockCategoryRepository),
				new(MockBannerRepository),
				new(MockRankingRepository),
			)

			ctx := context.Background()
			expectedBooks := []*bookstore2.Book{
				createTestBook("507f1f77bcf86cd799439011", "书籍", bookstore2.BookStatusPublished),
			}

			mockBookRepo.On("GetRecommended", ctx, tt.pageSize, tt.offset).Return(expectedBooks, nil)

			// Act
			result, err := service.GetRecommendedBooks(ctx, tt.page, tt.pageSize)

			// Assert
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			mockBookRepo.AssertExpectations(t)
		})
	}
}

// ============ Error Handling Tests ============

func TestBookstoreService_ErrorHandling(t *testing.T) {
	t.Run("Repository错误传播", func(t *testing.T) {
		// Arrange
		mockBookRepo := new(MockBookRepository)
		service := bookstoreService.NewBookstoreService(
			mockBookRepo,
			new(MockCategoryRepository),
			new(MockBannerRepository),
			new(MockRankingRepository),
		)

		ctx := context.Background()
		expectedError := errors.New("database connection failed")

		mockBookRepo.On("GetRecommended", ctx, 10, 0).Return(nil, expectedError)

		// Act
		result, err := service.GetRecommendedBooks(ctx, 1, 10)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get recommended books")
		mockBookRepo.AssertExpectations(t)
	})
}

// ============ Benchmark Tests ============

func BenchmarkBookstoreService_GetBookByID(b *testing.B) {
	mockBookRepo := new(MockBookRepository)
	service := bookstoreService.NewBookstoreService(
		mockBookRepo,
		new(MockCategoryRepository),
		new(MockBannerRepository),
		new(MockRankingRepository),
	)

	ctx := context.Background()
	bookID := "507f1f77bcf86cd799439011"
	objID, _ := primitive.ObjectIDFromHex(bookID)
	expectedBook := createTestBook(bookID, "测试书籍", bookstore2.BookStatusPublished)

	mockBookRepo.On("GetByID", ctx, objID).Return(expectedBook, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetBookByID(ctx, bookID)
	}
}
