package social

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"Qingyu_backend/models/social"
)

// BenchmarkGetRatingStats_CacheHit 基准测试：缓存命中场景
// 目标：测试从Redis缓存读取评分统计的性能
// 预期：<1ms/op，最小化内存分配
func BenchmarkGetRatingStats_CacheHit(b *testing.B) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 创建预期的评分统计（模拟缓存数据）
	expectedStats := &social.RatingStats{
		TargetID:      "benchmark123",
		TargetType:    "comment",
		AverageRating: 4.5,
		TotalRatings:  1,
		Distribution:  map[int]int64{5: 1},
		UpdatedAt:     time.Now(),
	}

	// 序列化为JSON（模拟Redis存储）
	serializedData := `{"targetId":"benchmark123","targetType":"comment","averageRating":4.5,"totalRatings":1,"distribution":{"5":1},"updatedAt":"` + expectedStats.UpdatedAt.Format(time.RFC3339Nano) + `"}`

	// 设置预期行为 - Redis返回缓存数据
	mockRedis.On("Get", mock.Anything, "rating:stats:comment:benchmark123").Return(serializedData, nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 重置计时器，排除setup时间
	b.ResetTimer()

	// 运行benchmark
	for i := 0; i < b.N; i++ {
		_, err := service.GetRatingStats(context.Background(), "comment", "benchmark123")
		if err != nil {
			b.Fatalf("GetRatingStats failed: %v", err)
		}
	}
}

// BenchmarkGetRatingStats_CacheMiss 基准测试：缓存未命中场景
// 目标：测试从数据库聚合评分并写入缓存的性能
// 预期：<100ms/op，符合性能指标要求
func BenchmarkGetRatingStats_CacheMiss(b *testing.B) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Redis返回空（缓存未命中）
	mockRedis.On("Get", mock.Anything, "rating:stats:comment:benchmark456").Return("", nil)

	// 设置Comment Repository返回评论数据
	comment := &social.Comment{
		AuthorID: "user_benchmark",
		Rating:   5,
		BookID:   "book_benchmark",
		Content:  "Great book!",
	}
	mockCommentRepo.On("GetByID", mock.Anything, "benchmark456").Return(comment, nil)

	// 设置Redis写入缓存
	mockRedis.On("Set", mock.Anything, "rating:stats:comment:benchmark456", mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 重置计时器，排除setup时间
	b.ResetTimer()

	// 运行benchmark
	for i := 0; i < b.N; i++ {
		_, err := service.GetRatingStats(context.Background(), "comment", "benchmark456")
		if err != nil {
			b.Fatalf("GetRatingStats failed: %v", err)
		}
	}
}

// BenchmarkGetRatingStats_Book_CacheMiss 基准测试：书籍评分聚合（缓存未命中）
// 目标：测试书籍级别评分聚合的性能（涉及更复杂的聚合逻辑）
// 预期：<100ms/op，符合性能指标要求
func BenchmarkGetRatingStats_Book_CacheMiss(b *testing.B) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Redis返回空（缓存未命中）
	mockRedis.On("Get", mock.Anything, "rating:stats:book:book_benchmark").Return("", nil)

	// 设置Comment Repository返回书籍评分统计
	statsMap := map[string]interface{}{
		"average_rating": 4.3,
		"total_ratings":  int64(100),
		"distribution": map[string]int64{
			"1": 5,
			"2": 10,
			"3": 20,
			"4": 35,
			"5": 30,
		},
	}
	mockCommentRepo.On("GetBookRatingStats", mock.Anything, "book_benchmark").Return(statsMap, nil)

	// 设置Redis写入缓存
	mockRedis.On("Set", mock.Anything, "rating:stats:book:book_benchmark", mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 重置计时器，排除setup时间
	b.ResetTimer()

	// 运行benchmark
	for i := 0; i < b.N; i++ {
		_, err := service.GetRatingStats(context.Background(), "book", "book_benchmark")
		if err != nil {
			b.Fatalf("GetRatingStats failed: %v", err)
		}
	}
}

// BenchmarkGetRatingStats_Concurrent_CacheHit 基准测试：并发缓存命中场景
// 目标：测试并发访问时缓存读取的性能和线程安全性
// 预期：支持高并发，性能稳定
func BenchmarkGetRatingStats_Concurrent_CacheHit(b *testing.B) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 创建预期的评分统计（模拟缓存数据）
	expectedStats := &social.RatingStats{
		TargetID:      "concurrent123",
		TargetType:    "comment",
		AverageRating: 4.5,
		TotalRatings:  1,
		Distribution:  map[int]int64{5: 1},
		UpdatedAt:     time.Now(),
	}

	// 序列化为JSON（模拟Redis存储）
	serializedData := `{"targetId":"concurrent123","targetType":"comment","averageRating":4.5,"totalRatings":1,"distribution":{"5":1},"updatedAt":"` + expectedStats.UpdatedAt.Format(time.RFC3339Nano) + `"}`

	// 设置预期行为 - Redis返回缓存数据
	// 允许并发调用
	mockRedis.On("Get", mock.Anything, "rating:stats:comment:concurrent123").Return(serializedData, nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 重置计时器，排除setup时间
	b.ResetTimer()

	// 运行并发benchmark
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := service.GetRatingStats(context.Background(), "comment", "concurrent123")
			if err != nil {
				b.Fatalf("GetRatingStats failed: %v", err)
			}
		}
	})
}

// BenchmarkSerializeStats 基准测试：评分统计序列化性能
// 目标：测试JSON序列化的性能
// 预期：最小化内存分配，优化序列化速度
func BenchmarkSerializeStats(b *testing.B) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger).(*RatingServiceImplementation)

	// 创建评分统计（包含复杂分布）
	stats := &social.RatingStats{
		TargetID:      "serialize123",
		TargetType:    "book",
		AverageRating: 4.2,
		TotalRatings:  1000,
		Distribution: map[int]int64{
			1: 50,
			2: 100,
			3: 200,
			4: 350,
			5: 300,
		},
		UpdatedAt: time.Now(),
	}

	// 重置计时器
	b.ResetTimer()

	// 运行benchmark
	for i := 0; i < b.N; i++ {
		_, err := service.serializeStats(stats)
		if err != nil {
			b.Fatalf("serializeStats failed: %v", err)
		}
	}
}

// BenchmarkDeserializeStats 基准测试：评分统计反序列化性能
// 目标：测试JSON反序列化的性能
// 预期：最小化内存分配，优化反序列化速度
func BenchmarkDeserializeStats(b *testing.B) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger).(*RatingServiceImplementation)

	// 创建序列化的评分统计
	stats := &social.RatingStats{
		TargetID:      "deserialize123",
		TargetType:    "book",
		AverageRating: 4.2,
		TotalRatings:  1000,
		Distribution: map[int]int64{
			1: 50,
			2: 100,
			3: 200,
			4: 350,
			5: 300,
		},
		UpdatedAt: time.Now(),
	}

	serializedData, _ := service.serializeStats(stats)

	// 重置计时器
	b.ResetTimer()

	// 运行benchmark
	for i := 0; i < b.N; i++ {
		_, err := service.deserializeStats(serializedData)
		if err != nil {
			b.Fatalf("deserializeStats failed: %v", err)
		}
	}
}

// BenchmarkInvalidateCache 基准测试：缓存失效性能
// 目标：测试缓存删除操作的性能
// 预期：<5ms/op，快速失效
func BenchmarkInvalidateCache(b *testing.B) {
	// 准备测试数据
	mockRedis := new(MockRedisClient)
	mockCommentRepo := new(MockCommentRepository)
	mockReviewRepo := new(MockReviewRepository)
	logger := zap.NewNop()

	// 设置Redis删除操作
	mockRedis.On("Delete", mock.Anything, []string{"rating:stats:book:invalidate123"}).Return(nil)

	// 创建service
	service := NewRatingService(mockCommentRepo, mockReviewRepo, mockRedis, logger)

	// 重置计时器
	b.ResetTimer()

	// 运行benchmark
	for i := 0; i < b.N; i++ {
		err := service.InvalidateCache(context.Background(), "book", "invalidate123")
		if err != nil {
			b.Fatalf("InvalidateCache failed: %v", err)
		}
	}
}
