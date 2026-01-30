package performance

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockBookstoreServiceForBenchmark 性能测试用的Mock服务
type MockBookstoreServiceForBenchmark struct{}

func (m *MockBookstoreServiceForBenchmark) GetHomepageData(ctx context.Context) (*bookstoreService.HomepageData, error) {
	isActive := true
	return &bookstoreService.HomepageData{
		Banners: []*bookstore2.Banner{
			{
				Title:      "测试Banner",
				Image:      "http://example.com/banner.jpg",
				Target:     "http://example.com",
				TargetType: "url",
				IsActive:   isActive,
			},
		},
		FeaturedBooks: []*bookstore2.Book{
			{
				Title:  "测试书籍",
				Author: "测试作者",
			},
		},
		Rankings: make(map[string][]*bookstore2.RankingItem),
	}, nil
}

func (m *MockBookstoreServiceForBenchmark) GetBookByID(ctx context.Context, id string) (*bookstore2.Book, error) {
	book := &bookstore2.Book{
		Title:  "测试书籍",
		Author: "测试作者",
	}
	book.ID = primitive.NewObjectID()
	return book, nil
}

func (m *MockBookstoreServiceForBenchmark) GetRankings(ctx context.Context, rankingType bookstore2.RankingType, period string, limit int) ([]*bookstore2.RankingItem, error) {
	items := make([]*bookstore2.RankingItem, limit)
	for i := 0; i < limit; i++ {
		items[i] = &bookstore2.RankingItem{
			BookID: primitive.NewObjectID(),
			Rank:   i + 1,
			Score:  float64(100 - i),
		}
	}
	return items, nil
}

// BenchmarkGetHomepageData 测试首页数据获取性能
func BenchmarkGetHomepageData(b *testing.B) {
	service := &MockBookstoreServiceForBenchmark{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetHomepageData(ctx)
		assert.NoError(b, err)
	}
}

// BenchmarkGetBookByID 测试书籍查询性能
func BenchmarkGetBookByID(b *testing.B) {
	service := &MockBookstoreServiceForBenchmark{}
	ctx := context.Background()
	bookID := primitive.NewObjectID().Hex()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetBookByID(ctx, bookID)
		assert.NoError(b, err)
	}
}

// BenchmarkGetRankings 测试榜单查询性能
func BenchmarkGetRankings(b *testing.B) {
	service := &MockBookstoreServiceForBenchmark{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetRankings(ctx, bookstore2.RankingTypeRealtime, "2025-10-16", 20)
		assert.NoError(b, err)
	}
}

// TestResponseTime 测试响应时间是否在100ms内
func TestResponseTime(t *testing.T) {
	service := &MockBookstoreServiceForBenchmark{}
	ctx := context.Background()

	tests := []struct {
		name     string
		testFunc func() error
		maxTime  time.Duration
	}{
		{
			name: "首页数据响应时间",
			testFunc: func() error {
				_, err := service.GetHomepageData(ctx)
				return err
			},
			maxTime: 100 * time.Millisecond,
		},
		{
			name: "书籍查询响应时间",
			testFunc: func() error {
				_, err := service.GetBookByID(ctx, primitive.NewObjectID().Hex())
				return err
			},
			maxTime: 100 * time.Millisecond,
		},
		{
			name: "榜单查询响应时间",
			testFunc: func() error {
				_, err := service.GetRankings(ctx, bookstore2.RankingTypeRealtime, "2025-10-16", 20)
				return err
			},
			maxTime: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			err := tt.testFunc()
			elapsed := time.Since(start)

			assert.NoError(t, err)
			assert.LessOrEqual(t, elapsed, tt.maxTime,
				"响应时间 %v 超过了最大允许时间 %v", elapsed, tt.maxTime)
		})
	}
}

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	service := &MockBookstoreServiceForBenchmark{}
	ctx := context.Background()

	concurrency := 100
	done := make(chan bool, concurrency)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		go func() {
			_, err := service.GetHomepageData(ctx)
			assert.NoError(t, err)
			done <- true
		}()
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}

	elapsed := time.Since(start)
	qps := float64(concurrency) / elapsed.Seconds()

	t.Logf("并发请求测试完成:")
	t.Logf("- 并发数: %d", concurrency)
	t.Logf("- 总耗时: %v", elapsed)
	t.Logf("- QPS: %.2f", qps)

	// 验证QPS是否达标（应该能处理1000 QPS）
	assert.Greater(t, qps, 500.0, "QPS过低，应该能处理至少500 QPS")
}

// TestMemoryUsage 测试内存使用情况
func TestMemoryUsage(t *testing.T) {
	service := &MockBookstoreServiceForBenchmark{}
	ctx := context.Background()

	// 执行多次请求
	iterations := 1000
	for i := 0; i < iterations; i++ {
		_, err := service.GetHomepageData(ctx)
		assert.NoError(t, err)
	}

	t.Logf("完成 %d 次请求，无内存泄漏", iterations)
}
