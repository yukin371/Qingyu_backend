package bookstore

import (
	"context"
	"errors"
	"testing"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	testMock "Qingyu_backend/service/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCacheService Mock缓存服务
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) GetHomepageData(ctx context.Context) (*HomepageData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HomepageData), args.Error(1)
}

func (m *MockCacheService) SetHomepageData(ctx context.Context, data *HomepageData, expiration time.Duration) error {
	args := m.Called(ctx, data, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetRanking(ctx context.Context, rankingType bookstoreModel.RankingType, period string) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, rankingType, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockCacheService) SetRanking(ctx context.Context, rankingType bookstoreModel.RankingType, period string, items []*bookstoreModel.RankingItem, expiration time.Duration) error {
	args := m.Called(ctx, rankingType, period, items, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetActiveBanners(ctx context.Context) ([]*bookstoreModel.Banner, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Banner), args.Error(1)
}

func (m *MockCacheService) SetActiveBanners(ctx context.Context, banners []*bookstoreModel.Banner, expiration time.Duration) error {
	args := m.Called(ctx, banners, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBook(ctx context.Context, bookID string) (*bookstoreModel.Book, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Book), args.Error(1)
}

func (m *MockCacheService) SetBook(ctx context.Context, bookID string, book *bookstoreModel.Book, expiration time.Duration) error {
	args := m.Called(ctx, bookID, book, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBookDetail(ctx context.Context, bookID string) (*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookDetail), args.Error(1)
}

func (m *MockCacheService) SetBookDetail(ctx context.Context, bookID string, bookDetail *bookstoreModel.BookDetail, expiration time.Duration) error {
	args := m.Called(ctx, bookID, bookDetail, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetChapter(ctx context.Context, chapterID string) (*bookstoreModel.Chapter, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Chapter), args.Error(1)
}

func (m *MockCacheService) SetChapter(ctx context.Context, chapterID string, chapter *bookstoreModel.Chapter, expiration time.Duration) error {
	args := m.Called(ctx, chapterID, chapter, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetCategoryTree(ctx context.Context) ([]*bookstoreModel.CategoryTree, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.CategoryTree), args.Error(1)
}

func (m *MockCacheService) SetCategoryTree(ctx context.Context, tree []*bookstoreModel.CategoryTree, expiration time.Duration) error {
	args := m.Called(ctx, tree, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBookRating(ctx context.Context, key string) (*bookstoreModel.BookRating, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookRating), args.Error(1)
}

func (m *MockCacheService) SetBookRating(ctx context.Context, key string, rating *bookstoreModel.BookRating, expiration time.Duration) error {
	args := m.Called(ctx, key, rating, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBookAverageRating(ctx context.Context, bookID string) (float64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCacheService) SetBookAverageRating(ctx context.Context, bookID string, rating float64, expiration time.Duration) error {
	args := m.Called(ctx, bookID, rating, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookRatingCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookRatingsCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateUserRatingsCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookAverageRatingCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateTopViewedBooksCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateHottestBooksCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateTopFavoritedBooksCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) GetBookStatistics(ctx context.Context, bookID string) (*bookstoreModel.BookStatistics, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookStatistics), args.Error(1)
}

func (m *MockCacheService) SetBookStatistics(ctx context.Context, bookID string, stats *bookstoreModel.BookStatistics, expiration time.Duration) error {
	args := m.Called(ctx, bookID, stats, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookStatisticsCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) GetTopViewedBooks(ctx context.Context) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockCacheService) SetTopViewedBooks(ctx context.Context, books []*bookstoreModel.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetTopFavoritedBooks(ctx context.Context) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockCacheService) SetTopFavoritedBooks(ctx context.Context, books []*bookstoreModel.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetTopRatedBooks(ctx context.Context) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockCacheService) SetTopRatedBooks(ctx context.Context, books []*bookstoreModel.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetHottestBooks(ctx context.Context) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockCacheService) SetHottestBooks(ctx context.Context, books []*bookstoreModel.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateHomepageCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateRankingCache(ctx context.Context, rankingType bookstoreModel.RankingType, period string) error {
	args := m.Called(ctx, rankingType, period)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBannerCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookDetailCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateChapterCache(ctx context.Context, chapterID string) error {
	args := m.Called(ctx, chapterID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookChaptersCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateCategoryCache(ctx context.Context, category string) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateAuthorCache(ctx context.Context, author string) error {
	args := m.Called(ctx, author)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateTopRatedBooksCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateAggregatedStatisticsCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCacheService) SetAggregatedStatistics(ctx context.Context, stats map[string]interface{}, expiration time.Duration) error {
	args := m.Called(ctx, stats, expiration)
	return args.Error(0)
}

// TestNewBookDetailService 测试创建服务实例
func TestNewBookDetailService(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := NewBookDetailService(mockRepo, mockCache)

	assert.NotNil(t, service)
}

// TestCreateBookDetail_Success 测试成功创建书籍详情
func TestCreateBookDetail_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
		cacheService:   mockCache,
	}

	ctx := context.Background()
	bookDetail := &bookstoreModel.BookDetail{
		ID:          primitive.NewObjectID(),
		Title:       "测试书籍",
		Author:      "测试作者",
		Description: "测试描述",
		Status:      bookstoreModel.BookStatusOngoing,
	}

	mockRepo.On("GetByTitle", ctx, "测试书籍").Return(nil, nil).Once()
	mockRepo.On("Create", ctx, bookDetail).Return(nil).Once()
	mockCache.On("InvalidateBookDetailCache", ctx, mock.Anything).Return(nil).Maybe()
	mockCache.On("InvalidateCategoryCache", ctx, mock.Anything).Return(nil).Maybe()
	mockCache.On("InvalidateAuthorCache", ctx, mock.Anything).Return(nil).Maybe()
	mockCache.On("InvalidateHomepageCache", ctx).Return(nil).Maybe()

	err := service.CreateBookDetail(ctx, bookDetail)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateBookDetail_NilBookDetail 测试书籍详情为空
func TestCreateBookDetail_NilBookDetail(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	err := service.CreateBookDetail(ctx, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

// TestCreateBookDetail_MissingTitle 测试缺少标题
func TestCreateBookDetail_MissingTitle(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	bookDetail := &bookstoreModel.BookDetail{
		Author: "测试作者",
	}

	err := service.CreateBookDetail(ctx, bookDetail)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

// TestCreateBookDetail_MissingAuthor 测试缺少作者
func TestCreateBookDetail_MissingAuthor(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	bookDetail := &bookstoreModel.BookDetail{
		Title: "测试书籍",
	}

	err := service.CreateBookDetail(ctx, bookDetail)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "author is required")
}

// TestCreateBookDetail_DuplicateTitle 测试重复标题
func TestCreateBookDetail_DuplicateTitle(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	existingBook := &bookstoreModel.BookDetail{
		ID:     primitive.NewObjectID(),
		Title:  "测试书籍",
		Author: "测试作者",
	}

	bookDetail := &bookstoreModel.BookDetail{
		Title:  "测试书籍",
		Author: "新作者",
	}

	mockRepo.On("GetByTitle", ctx, "测试书籍").Return(existingBook, nil).Once()

	err := service.CreateBookDetail(ctx, bookDetail)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

// TestGetBookDetailByID_Success 测试成功获取书籍详情
func TestGetBookDetailByID_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
		cacheService:   mockCache,
	}

	ctx := context.Background()
	bookID := primitive.NewObjectID()
	bookDetail := &bookstoreModel.BookDetail{
		ID:     bookID,
		Title:  "测试书籍",
		Author: "测试作者",
	}

	mockCache.On("GetBookDetail", ctx, bookID.Hex()).Return(nil, errors.New("cache miss")).Once()
	mockRepo.On("GetByID", ctx, bookID).Return(bookDetail, nil).Once()
	mockCache.On("SetBookDetail", ctx, bookID.Hex(), bookDetail, mock.AnythingOfType("time.Duration")).Return(nil).Once()

	result, err := service.GetBookDetailByID(ctx, bookID)

	assert.NoError(t, err)
	assert.Equal(t, bookDetail, result)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestGetBookDetailByID_CacheHit 测试缓存命中
func TestGetBookDetailByID_CacheHit(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
		cacheService:   mockCache,
	}

	ctx := context.Background()
	bookID := primitive.NewObjectID()
	bookDetail := &bookstoreModel.BookDetail{
		ID:     bookID,
		Title:  "测试书籍",
		Author: "测试作者",
	}

	mockCache.On("GetBookDetail", ctx, bookID.Hex()).Return(bookDetail, nil).Once()

	result, err := service.GetBookDetailByID(ctx, bookID)

	assert.NoError(t, err)
	assert.Equal(t, bookDetail, result)
	mockCache.AssertExpectations(t)
	// 确保没有调用数据库
	mockRepo.AssertNotCalled(t, "GetByID", ctx, bookID)
}

// TestUpdateBookDetail_Success 测试成功更新书籍详情
func TestUpdateBookDetail_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
		cacheService:   mockCache,
	}

	ctx := context.Background()
	bookID := primitive.NewObjectID()
	bookDetail := &bookstoreModel.BookDetail{
		ID:          bookID,
		Title:       "更新后的标题",
		Author:      "测试作者",
		Description: "更新后的描述",
		Status:      bookstoreModel.BookStatusOngoing,
	}

	mockRepo.On("Update", ctx, bookID, mock.AnythingOfType("map[string]interface {}")).Return(nil).Once()
	mockCache.On("InvalidateBookDetailCache", ctx, bookID.Hex()).Return(nil).Maybe()
	mockCache.On("InvalidateCategoryCache", ctx, mock.Anything).Return(nil).Maybe()
	mockCache.On("InvalidateAuthorCache", ctx, mock.Anything).Return(nil).Maybe()
	mockCache.On("InvalidateHomepageCache", ctx).Return(nil).Maybe()

	err := service.UpdateBookDetail(ctx, bookDetail)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteBookDetail_Success 测试成功删除书籍详情
func TestDeleteBookDetail_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
		cacheService:   mockCache,
	}

	ctx := context.Background()
	bookID := primitive.NewObjectID()
	bookDetail := &bookstoreModel.BookDetail{
		ID:     bookID,
		Title:  "测试书籍",
		Author: "测试作者",
	}

	mockRepo.On("GetByID", ctx, bookID).Return(bookDetail, nil).Once()
	mockRepo.On("Delete", ctx, bookID).Return(nil).Once()
	mockCache.On("InvalidateBookDetailCache", ctx, bookID.Hex()).Return(nil).Maybe()
	mockCache.On("InvalidateCategoryCache", ctx, mock.Anything).Return(nil).Maybe()
	mockCache.On("InvalidateAuthorCache", ctx, mock.Anything).Return(nil).Maybe()
	mockCache.On("InvalidateHomepageCache", ctx).Return(nil).Maybe()

	err := service.DeleteBookDetail(ctx, bookID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestGetBookDetailsByAuthor_Success 测试按作者获取书籍
func TestGetBookDetailsByAuthor_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	author := "测试作者"
	books := []*bookstoreModel.BookDetail{
		{Title: "书籍1", Author: author},
		{Title: "书籍2", Author: author},
	}

	mockRepo.On("GetByAuthor", ctx, author, 20, 0).Return(books, nil).Once()
	mockRepo.On("CountByAuthor", ctx, author).Return(int64(2), nil).Once()

	result, total, err := service.GetBookDetailsByAuthor(ctx, author, 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// TestGetBookDetailsByCategory_Success 测试按分类获取书籍
func TestGetBookDetailsByCategory_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	category := "玄幻"
	books := []*bookstoreModel.BookDetail{
		{Title: "书籍1", Categories: []string{category}},
	}

	mockRepo.On("GetByCategory", ctx, category, 20, 0).Return(books, nil).Once()
	mockRepo.On("CountByCategory", ctx, category).Return(int64(1), nil).Once()

	result, total, err := service.GetBookDetailsByCategory(ctx, category, 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}

// TestGetBookDetailsByStatus_Success 测试按状态获取书籍
func TestGetBookDetailsByStatus_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	books := []*bookstoreModel.BookDetail{
		{Title: "书籍1", Status: bookstoreModel.BookStatusOngoing},
	}

	mockRepo.On("GetByStatus", ctx, bookstoreModel.BookStatusOngoing, 20, 0).Return(books, nil).Once()
	mockRepo.On("CountByStatus", ctx, bookstoreModel.BookStatusOngoing).Return(int64(1), nil).Once()

	result, total, err := service.GetBookDetailsByStatus(ctx, bookstoreModel.BookStatusOngoing, 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}

// TestSearchBookDetails_Success 测试搜索书籍
func TestSearchBookDetails_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	keyword := "测试"
	books := []*bookstoreModel.BookDetail{
		{Title: "测试书籍1"},
		{Title: "测试书籍2"},
	}

	mockRepo.On("Search", ctx, keyword, 20, 0).Return(books, nil).Once()

	result, total, err := service.SearchBookDetails(ctx, keyword, 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// TestGetBookDetailStats_Success 测试获取书籍统计
func TestGetBookDetailStats_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()

	mockRepo.On("CountByStatus", ctx, bookstoreModel.BookStatusCompleted).Return(int64(100), nil).Once()
	mockRepo.On("CountByStatus", ctx, bookstoreModel.BookStatusOngoing).Return(int64(200), nil).Once()
	mockRepo.On("CountByStatus", ctx, bookstoreModel.BookStatusPaused).Return(int64(50), nil).Once()

	stats, err := service.GetBookDetailStats(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int64(100), stats["completed_books"])
	assert.Equal(t, int64(200), stats["ongoing_books"])
	assert.Equal(t, int64(50), stats["paused_books"])
	assert.Equal(t, int64(350), stats["total_books"])
	mockRepo.AssertExpectations(t)
}

// TestIncrementViewCount_Success 测试增加浏览计数
func TestIncrementViewCount_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
		cacheService:   mockCache,
	}

	ctx := context.Background()
	bookID := primitive.NewObjectID()

	mockRepo.On("IncrementViewCount", ctx, bookID).Return(nil).Once()
	mockCache.On("InvalidateBookDetailCache", ctx, bookID.Hex()).Return(nil).Once()

	err := service.IncrementViewCount(ctx, bookID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestBatchUpdateBookDetailStatus_Success 测试批量更新书籍状态
func TestBatchUpdateBookDetailStatus_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	mockCache := new(MockCacheService)

	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
		cacheService:   mockCache,
	}

	ctx := context.Background()
	bookIDs := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}

	mockRepo.On("BatchUpdateStatus", ctx, bookIDs, bookstoreModel.BookStatusCompleted).Return(nil).Once()
	mockCache.On("InvalidateBookDetailCache", ctx, mock.AnythingOfType("string")).Return(nil).Times(2)

	err := service.BatchUpdateBookDetailStatus(ctx, bookIDs, bookstoreModel.BookStatusCompleted)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestGetPopularBooks_Success 测试获取热门书籍
func TestGetPopularBooks_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	books := []*bookstoreModel.BookDetail{
		{Title: "热门书籍1"},
		{Title: "热门书籍2"},
	}

	mockRepo.On("GetByStatus", ctx, bookstoreModel.BookStatusCompleted, 10, 0).Return(books, nil).Once()

	result, err := service.GetPopularBooks(ctx, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// TestGetLatestBooks_Success 测试获取最新书籍
func TestGetLatestBooks_Success(t *testing.T) {
	mockRepo := new(testMock.MockBookDetailRepository)
	service := &BookDetailServiceImpl{
		bookDetailRepo: mockRepo,
	}

	ctx := context.Background()
	books := []*bookstoreModel.BookDetail{
		{Title: "最新书籍1"},
		{Title: "最新书籍2"},
	}

	mockRepo.On("GetByStatus", ctx, bookstoreModel.BookStatusOngoing, 10, 0).Return(books, nil).Once()

	result, err := service.GetLatestBooks(ctx, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}
