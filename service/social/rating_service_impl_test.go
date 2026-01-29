package social

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"Qingyu_backend/models/social"
	"Qingyu_backend/pkg/cache"
)

// MockReviewRepository 模拟书评仓库
type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) CreateReview(ctx context.Context, review *social.Review) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockReviewRepository) GetReviewByID(ctx context.Context, reviewID string) (*social.Review, error) {
	args := m.Called(ctx, reviewID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Review), args.Error(1)
}

func (m *MockReviewRepository) GetReviewsByBook(ctx context.Context, bookID string, page, size int) ([]*social.Review, int64, error) {
	args := m.Called(ctx, bookID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) GetReviewsByUser(ctx context.Context, userID string, page, size int) ([]*social.Review, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) GetPublicReviews(ctx context.Context, page, size int) ([]*social.Review, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) GetReviewsByRating(ctx context.Context, bookID string, rating int, page, size int) ([]*social.Review, int64, error) {
	args := m.Called(ctx, bookID, rating, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) UpdateReview(ctx context.Context, reviewID string, updates map[string]interface{}) error {
	args := m.Called(ctx, reviewID, updates)
	return args.Error(0)
}

func (m *MockReviewRepository) DeleteReview(ctx context.Context, reviewID string) error {
	args := m.Called(ctx, reviewID)
	return args.Error(0)
}

func (m *MockReviewRepository) CreateReviewLike(ctx context.Context, reviewLike *social.ReviewLike) error {
	args := m.Called(ctx, reviewLike)
	return args.Error(0)
}

func (m *MockReviewRepository) DeleteReviewLike(ctx context.Context, reviewID, userID string) error {
	args := m.Called(ctx, reviewID, userID)
	return args.Error(0)
}

func (m *MockReviewRepository) GetReviewLike(ctx context.Context, reviewID, userID string) (*social.ReviewLike, error) {
	args := m.Called(ctx, reviewID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.ReviewLike), args.Error(1)
}

func (m *MockReviewRepository) IsReviewLiked(ctx context.Context, reviewID, userID string) (bool, error) {
	args := m.Called(ctx, reviewID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockReviewRepository) GetReviewLikes(ctx context.Context, reviewID string, page, size int) ([]*social.ReviewLike, int64, error) {
	args := m.Called(ctx, reviewID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.ReviewLike), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) IncrementReviewLikeCount(ctx context.Context, reviewID string) error {
	args := m.Called(ctx, reviewID)
	return args.Error(0)
}

func (m *MockReviewRepository) DecrementReviewLikeCount(ctx context.Context, reviewID string) error {
	args := m.Called(ctx, reviewID)
	return args.Error(0)
}

func (m *MockReviewRepository) GetAverageRating(ctx context.Context, bookID string) (float64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockReviewRepository) GetRatingDistribution(ctx context.Context, bookID string) (map[int]int64, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int]int64), args.Error(1)
}

func (m *MockReviewRepository) CountReviews(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReviewRepository) CountUserReviews(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockReviewRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockRedisClient 模拟Redis客户端
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	args := m.Called(ctx, keys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRedisClient) MSet(ctx context.Context, pairs ...interface{}) error {
	args := m.Called(ctx, pairs)
	return args.Error(0)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	args := m.Called(ctx, key, fields)
	return args.Error(0)
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Decr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) GetClient() interface{} {
	args := m.Called()
	return args.Get(0)
}

// TestGetRatingStats_CacheHit 测试缓存命中场景
func TestGetRatingStats_CacheHit(t *testing.T) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 创建预期的评分统计
	expectedStats := &social.RatingStats{
		TargetID:      "123",
		TargetType:    "comment",
		AverageRating: 4.5,
		TotalRatings:  1,
		Distribution:  map[int]int64{5: 1},
		UpdatedAt:     time.Now(),
	}

	// 序列化为JSON
	serializedData := `{"targetId":"123","targetType":"comment","averageRating":4.5,"totalRatings":1,"distribution":{"5":1},"updatedAt":"` + expectedStats.UpdatedAt.Format(time.RFC3339Nano) + `"}`

	// 设置预期行为 - Redis返回缓存数据
	mockRedis.On("Get", mock.Anything, "rating:stats:comment:123").Return(serializedData, nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 执行测试
	stats, err := service.GetRatingStats(context.Background(), "comment", "123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "123", stats.TargetID)
	assert.Equal(t, "comment", stats.TargetType)
	assert.Equal(t, 4.5, stats.AverageRating)
	assert.Equal(t, int64(1), stats.TotalRatings)
	assert.Equal(t, int64(1), stats.Distribution[5])

	// 验证mock调用
	mockRedis.AssertExpectations(t)
}

// TestGetRatingStats_CacheMiss 测试缓存未命中，从DB聚合
func TestGetRatingStats_CacheMiss(t *testing.T) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Redis返回空（缓存未命中）
	mockRedis.On("Get", mock.Anything, "rating:stats:comment:456").Return("", cache.ErrRedisNil)

	// 设置Comment Repository返回评论数据
	comment := &social.Comment{
		AuthorID: "user123",
		Rating:   5,
	}
	mockCommentRepo.On("GetByID", mock.Anything, "456").Return(comment, nil)

	// 设置Redis写入缓存
	mockRedis.On("Set", mock.Anything, "rating:stats:comment:456", mock.Anything, 5*time.Minute).Return(nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 执行测试
	stats, err := service.GetRatingStats(context.Background(), "comment", "456")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "456", stats.TargetID)
	assert.Equal(t, "comment", stats.TargetType)
	assert.Equal(t, 5.0, stats.AverageRating)
	assert.Equal(t, int64(1), stats.TotalRatings)
	assert.Equal(t, int64(1), stats.Distribution[5])

	// 验证mock调用
	mockRedis.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

// TestAggregateRatings_Comment 测试聚合评论评分
func TestAggregateRatings_Comment(t *testing.T) {
	// 准备测试数据
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Comment Repository返回评论数据
	comment := &social.Comment{
		AuthorID: "user789",
		Rating:   4,
	}
	mockCommentRepo.On("GetByID", mock.Anything, "comment789").Return(comment, nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, nil, logger)

	// 执行测试
	stats, err := service.AggregateRatings(context.Background(), "comment", "comment789")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "comment789", stats.TargetID)
	assert.Equal(t, "comment", stats.TargetType)
	assert.Equal(t, 4.0, stats.AverageRating)
	assert.Equal(t, int64(1), stats.TotalRatings)
	assert.Equal(t, int64(1), stats.Distribution[4])

	// 验证mock调用
	mockCommentRepo.AssertExpectations(t)
}

// TestAggregateRatings_Review 测试聚合审核评分
func TestAggregateRatings_Review(t *testing.T) {
	// 准备测试数据
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Review Repository返回书评数据
	review := &social.Review{
		UserID: "user456",
		Rating: 5,
	}
	mockReviewRepo.On("GetReviewByID", mock.Anything, "review123").Return(review, nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, nil, logger)

	// 执行测试
	stats, err := service.AggregateRatings(context.Background(), "review", "review123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "review123", stats.TargetID)
	assert.Equal(t, "review", stats.TargetType)
	assert.Equal(t, 5.0, stats.AverageRating)
	assert.Equal(t, int64(1), stats.TotalRatings)
	assert.Equal(t, int64(1), stats.Distribution[5])

	// 验证mock调用
	mockReviewRepo.AssertExpectations(t)
}

// TestInvalidateCache 测试缓存失效
func TestInvalidateCache(t *testing.T) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Redis删除操作 - Delete接受(ctx, ...keys)
	mockRedis.On("Delete", mock.Anything, []string{"rating:stats:book:book123"}).Return(nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 执行测试
	err := service.InvalidateCache(context.Background(), "book", "book123")

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRedis.AssertExpectations(t)
}

// TestInvalidateCache_NoRedis 测试无Redis客户端时的缓存失效
func TestInvalidateCache_NoRedis(t *testing.T) {
	// 准备测试数据
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 创建service（无Redis客户端）
	service := NewRatingService(mockCommentRepo, mockReviewRepo, nil, logger)

	// 执行测试
	err := service.InvalidateCache(context.Background(), "book", "book123")

	// 验证结果
	assert.NoError(t, err)
}

// TestAggregateRatings_UnsupportedType 测试不支持的目标类型
func TestAggregateRatings_UnsupportedType(t *testing.T) {
	// 准备测试数据
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, nil, logger)

	// 执行测试
	stats, err := service.AggregateRatings(context.Background(), "unsupported", "id123")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "不支持的目标类型")
}

// TestAggregateRatings_Book 测试聚合书籍评分
func TestAggregateRatings_Book(t *testing.T) {
	// 准备测试数据
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Comment Repository返回书籍评分统计
	statsMap := map[string]interface{}{
		"average_rating": 4.3,
		"total_ratings":  int64(10),
		"distribution": map[string]int64{
			"1": 0,
			"2": 1,
			"3": 2,
			"4": 3,
			"5": 4,
		},
	}
	mockCommentRepo.On("GetBookRatingStats", mock.Anything, "book456").Return(statsMap, nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, nil, logger)

	// 执行测试
	stats, err := service.AggregateRatings(context.Background(), "book", "book456")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "book456", stats.TargetID)
	assert.Equal(t, "book", stats.TargetType)
	assert.Equal(t, 4.3, stats.AverageRating)
	assert.Equal(t, int64(10), stats.TotalRatings)
	assert.Equal(t, int64(1), stats.Distribution[2])
	assert.Equal(t, int64(2), stats.Distribution[3])
	assert.Equal(t, int64(3), stats.Distribution[4])
	assert.Equal(t, int64(4), stats.Distribution[5])

	// 验证mock调用
	mockCommentRepo.AssertExpectations(t)
}

// TestGetUserRating_Book 测试获取用户对书籍的评分
func TestGetUserRating_Book(t *testing.T) {
	// 准备测试数据
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Comment Repository返回书籍的评论列表
	comments := []*social.Comment{
		{
			AuthorID: "user123",
			Rating:   5,
		},
		{
			AuthorID: "user456",
			Rating:   4,
		},
	}
	mockCommentRepo.On("GetCommentsByBookID", mock.Anything, "book789", 1, 100).Return(comments, int64(2), nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, nil, logger)

	// 执行测试 - 用户有评分
	rating, err := service.GetUserRating(context.Background(), "user123", "book", "book789")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 5, rating)

	// 执行测试 - 用户无评分
	rating, err = service.GetUserRating(context.Background(), "user999", "book", "book789")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 0, rating)

	// 验证mock调用
	mockCommentRepo.AssertExpectations(t)
}

// TestGetRatingStats_CacheDeserializeError 测试缓存反序列化失败场景
func TestGetRatingStats_CacheDeserializeError(t *testing.T) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Redis返回无效的JSON数据
	mockRedis.On("Get", mock.Anything, "rating:stats:comment:999").Return("invalid json", nil)

	// 设置Comment Repository返回评论数据
	comment := &social.Comment{
		AuthorID: "user999",
		Rating:   3,
	}
	mockCommentRepo.On("GetByID", mock.Anything, "999").Return(comment, nil)

	// 设置Redis写入缓存
	mockRedis.On("Set", mock.Anything, "rating:stats:comment:999", mock.Anything, 5*time.Minute).Return(nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 执行测试
	stats, err := service.GetRatingStats(context.Background(), "comment", "999")

	// 验证结果 - 应该从数据库重新聚合
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "999", stats.TargetID)
	assert.Equal(t, 3.0, stats.AverageRating)

	// 验证mock调用
	mockRedis.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

// TestInvalidateCache_DeleteError 测试缓存失效失败场景
func TestInvalidateCache_DeleteError(t *testing.T) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Redis删除操作返回错误 - Delete接受(ctx, ...keys)
	mockRedis.On("Delete", mock.Anything, []string{"rating:stats:book:error123"}).Return(errors.New("redis connection error"))

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 执行测试
	err := service.InvalidateCache(context.Background(), "book", "error123")

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "缓存失效失败")

	// 验证mock调用
	mockRedis.AssertExpectations(t)
}
