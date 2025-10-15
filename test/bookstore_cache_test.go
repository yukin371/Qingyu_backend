package test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockRedisClient 模拟Redis客户端
type MockRedisClient struct {
	mock.Mock
	data map[string]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx, "get", key)

	if value, exists := m.data[key]; exists {
		cmd.SetVal(value)
	} else {
		cmd.SetErr(redis.Nil)
	}

	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.Called(ctx, key, value, expiration)

	if str, ok := value.(string); ok {
		m.data[key] = str
	} else if bytes, ok := value.([]byte); ok {
		m.data[key] = string(bytes)
	}

	cmd := redis.NewStatusCmd(ctx, "set", key, value)
	cmd.SetVal("OK")
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	m.Called(ctx, keys)

	deleted := int64(0)
	for _, key := range keys {
		if _, exists := m.data[key]; exists {
			delete(m.data, key)
			deleted++
		}
	}

	cmd := redis.NewIntCmd(ctx, "del", keys)
	cmd.SetVal(deleted)
	return cmd
}

func (m *MockRedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	m.Called(ctx, pattern)
	cmd := redis.NewStringSliceCmd(ctx, "keys", pattern)

	var keys []string
	for key := range m.data {
		if pattern == "*" || key == pattern {
			keys = append(keys, key)
		}
	}

	cmd.SetVal(keys)
	return cmd
}

// TestRedisCacheService_GetHomepageData 测试获取首页数据缓存
func TestRedisCacheService_GetHomepageData(t *testing.T) {
	mockClient := NewMockRedisClient()
	cacheService := bookstoreService.NewRedisCacheService(mockClient, "test")

	// 测试缓存不存在的情况
	mockClient.On("Get", mock.Anything, "test:bookstore:homepage").Return(nil)

	ctx := context.Background()
	result, err := cacheService.GetHomepageData(ctx)

	assert.NoError(t, err)
	assert.Nil(t, result)

	// 测试缓存存在的情况
	homepageData := &bookstoreService.HomepageData{
		Banners: []*bookstore.Banner{
			{
				ID:    primitive.NewObjectID(),
				Title: "Test Banner",
			},
		},
		Rankings: map[string][]*bookstore.RankingItem{
			"realtime": {
				{
					ID:   primitive.NewObjectID(),
					Rank: 1,
				},
			},
		},
	}

	jsonData, _ := json.Marshal(homepageData)
	mockClient.data["test:bookstore:homepage"] = string(jsonData)

	result, err = cacheService.GetHomepageData(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Banners, 1)
	assert.Equal(t, "Test Banner", result.Banners[0].Title)
}

// TestRedisCacheService_SetHomepageData 测试设置首页数据缓存
func TestRedisCacheService_SetHomepageData(t *testing.T) {
	mockClient := NewMockRedisClient()
	cacheService := bookstoreService.NewRedisCacheService(mockClient, "test")

	homepageData := &bookstoreService.HomepageData{
		Banners: []*bookstore.Banner{
			{
				ID:    primitive.NewObjectID(),
				Title: "Test Banner",
			},
		},
	}

	mockClient.On("Set", mock.Anything, "test:bookstore:homepage", mock.Anything, 5*time.Minute).Return(nil)

	ctx := context.Background()
	err := cacheService.SetHomepageData(ctx, homepageData, 5*time.Minute)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// TestRedisCacheService_GetRanking 测试获取榜单缓存
func TestRedisCacheService_GetRanking(t *testing.T) {
	mockClient := NewMockRedisClient()
	cacheService := bookstoreService.NewRedisCacheService(mockClient, "test")

	// 准备测试数据
	rankingItems := []*bookstore.RankingItem{
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore.RankingTypeRealtime,
			Rank:      1,
			Score:     100.0,
			ViewCount: 1000,
			Period:    "2024-01-15",
		},
	}

	jsonData, _ := json.Marshal(rankingItems)
	key := "test:bookstore:ranking:realtime:2024-01-15"
	mockClient.data[key] = string(jsonData)

	mockClient.On("Get", mock.Anything, key).Return(nil)

	ctx := context.Background()
	result, err := cacheService.GetRanking(ctx, bookstore.RankingTypeRealtime, "2024-01-15")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, result[0].Rank)
}

// TestCachedBookstoreService 测试带缓存的书城服务
func TestCachedBookstoreService(t *testing.T) {
	// 创建模拟服务
	mockService := new(MockBookstoreServiceCache)
	mockCache := new(MockCacheService)

	cachedService := bookstoreService.NewCachedBookstoreService(mockService, mockCache)

	// 测试获取首页数据（缓存命中）
	expectedData := &bookstoreService.HomepageData{
		Banners: []*bookstore.Banner{
			{ID: primitive.NewObjectID(), Title: "Cached Banner"},
		},
	}

	mockCache.On("GetHomepageData", mock.Anything).Return(expectedData, nil)

	ctx := context.Background()
	result, err := cachedService.GetHomepageData(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)

	// 验证没有调用原服务
	mockService.AssertNotCalled(t, "GetHomepageData")
	mockCache.AssertExpectations(t)
}

// TestCachedBookstoreService_CacheMiss 测试缓存未命中的情况
func TestCachedBookstoreService_CacheMiss(t *testing.T) {
	mockService := new(MockBookstoreServiceCache)
	mockCache := new(MockCacheService)

	cachedService := bookstoreService.NewCachedBookstoreService(mockService, mockCache)

	expectedData := &bookstoreService.HomepageData{
		Banners: []*bookstore.Banner{
			{ID: primitive.NewObjectID(), Title: "Service Banner"},
		},
	}

	// 模拟缓存未命中
	mockCache.On("GetHomepageData", mock.Anything).Return(nil, nil)
	// 模拟从服务获取数据
	mockService.On("GetHomepageData", mock.Anything).Return(expectedData, nil)
	// 模拟设置缓存（异步调用，可能不会立即执行）
	mockCache.On("SetHomepageData", mock.Anything, expectedData, bookstoreService.HomepageCacheExpiration).Return(nil).Maybe()

	ctx := context.Background()
	result, err := cachedService.GetHomepageData(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)

	mockService.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// MockBookstoreServiceCache 模拟书城服务（用于缓存测试）
type MockBookstoreServiceCache struct {
	mock.Mock
}

func (m *MockBookstoreServiceCache) GetHomepageData(ctx context.Context) (*bookstoreService.HomepageData, error) {
	args := m.Called(ctx)
	return args.Get(0).(*bookstoreService.HomepageData), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetActiveBanners(ctx context.Context, limit int) ([]*bookstore.Banner, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.Banner), args.Error(1)
}

// 实现其他必需的接口方法（简化实现）
func (m *MockBookstoreServiceCache) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore.Book, int64, error) {
	args := m.Called(ctx, categoryID, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreServiceCache) GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreServiceCache) SearchBooksWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreServiceCache) GetBookStats(ctx context.Context) (*bookstore.BookStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*bookstore.BookStats), args.Error(1)
}

func (m *MockBookstoreServiceCache) IncrementBookView(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookstoreServiceCache) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore.CategoryTree), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetCategoryByID(ctx context.Context, id string) (*bookstore.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetRootCategories(ctx context.Context) ([]*bookstore.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockBookstoreServiceCache) IncrementBannerClick(ctx context.Context, bannerID string) error {
	args := m.Called(ctx, bannerID)
	return args.Error(0)
}

func (m *MockBookstoreServiceCache) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetRankingByType(ctx context.Context, rankingType bookstore.RankingType, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreServiceCache) UpdateRankings(ctx context.Context, rankingType bookstore.RankingType, period string) error {
	args := m.Called(ctx, rankingType, period)
	return args.Error(0)
}

func (m *MockBookstoreServiceCache) GetFreeBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetHotBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreServiceCache) GetNewReleases(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreServiceCache) SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.Book, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Get(1).(int64), args.Error(2)
}

// MockCacheService 模拟缓存服务
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) GetHomepageData(ctx context.Context) (*bookstoreService.HomepageData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreService.HomepageData), args.Error(1)
}

func (m *MockCacheService) SetHomepageData(ctx context.Context, data *bookstoreService.HomepageData, expiration time.Duration) error {
	args := m.Called(ctx, data, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetRanking(ctx context.Context, rankingType bookstore.RankingType, period string) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, rankingType, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockCacheService) SetRanking(ctx context.Context, rankingType bookstore.RankingType, period string, items []*bookstore.RankingItem, expiration time.Duration) error {
	args := m.Called(ctx, rankingType, period, items, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetActiveBanners(ctx context.Context) ([]*bookstore.Banner, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Banner), args.Error(1)
}

func (m *MockCacheService) SetActiveBanners(ctx context.Context, banners []*bookstore.Banner, expiration time.Duration) error {
	args := m.Called(ctx, banners, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBook(ctx context.Context, bookID string) (*bookstore.Book, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Book), args.Error(1)
}

func (m *MockCacheService) SetBook(ctx context.Context, bookID string, book *bookstore.Book, expiration time.Duration) error {
	args := m.Called(ctx, bookID, book, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.CategoryTree), args.Error(1)
}

func (m *MockCacheService) SetCategoryTree(ctx context.Context, tree []*bookstore.CategoryTree, expiration time.Duration) error {
	args := m.Called(ctx, tree, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateHomepageCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateRankingCache(ctx context.Context, rankingType bookstore.RankingType, period string) error {
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

func (m *MockCacheService) InvalidateCategoryCache(ctx context.Context, categoryID string) error {
	args := m.Called(ctx, categoryID)
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

func (m *MockCacheService) GetBookAverageRating(ctx context.Context, bookID string) (float64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCacheService) SetBookAverageRating(ctx context.Context, bookID string, rating float64, expiration time.Duration) error {
	args := m.Called(ctx, bookID, rating, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBookDetail(ctx context.Context, bookID string) (*bookstore.BookDetail, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookDetail), args.Error(1)
}

func (m *MockCacheService) SetBookDetail(ctx context.Context, bookID string, book *bookstore.BookDetail, expiration time.Duration) error {
	args := m.Called(ctx, bookID, book, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBookRating(ctx context.Context, bookID string) (*bookstore.BookRating, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookRating), args.Error(1)
}

func (m *MockCacheService) SetBookRating(ctx context.Context, bookID string, rating *bookstore.BookRating, expiration time.Duration) error {
	args := m.Called(ctx, bookID, rating, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetBookStatistics(ctx context.Context, bookID string) (*bookstore.BookStatistics, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookStatistics), args.Error(1)
}

func (m *MockCacheService) SetBookStatistics(ctx context.Context, bookID string, stats *bookstore.BookStatistics, expiration time.Duration) error {
	args := m.Called(ctx, bookID, stats, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetChapter(ctx context.Context, chapterID string) (*bookstore.Chapter, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockCacheService) SetChapter(ctx context.Context, chapterID string, chapter *bookstore.Chapter, expiration time.Duration) error {
	args := m.Called(ctx, chapterID, chapter, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetHottestBooks(ctx context.Context) ([]*bookstore.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockCacheService) SetHottestBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetTopFavoritedBooks(ctx context.Context) ([]*bookstore.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockCacheService) SetTopFavoritedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetTopRatedBooks(ctx context.Context) ([]*bookstore.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockCacheService) SetTopRatedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetTopViewedBooks(ctx context.Context) ([]*bookstore.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockCacheService) SetTopViewedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	args := m.Called(ctx, books, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateAggregatedStatisticsCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateAuthorCache(ctx context.Context, authorID string) error {
	args := m.Called(ctx, authorID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookAverageRatingCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookChaptersCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookDetailCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
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

func (m *MockCacheService) InvalidateBookStatisticsCache(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateChapterCache(ctx context.Context, chapterID string) error {
	args := m.Called(ctx, chapterID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateTopRatedBooksCache(ctx context.Context) error {
	args := m.Called(ctx)
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

func (m *MockCacheService) InvalidateUserRatingsCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
