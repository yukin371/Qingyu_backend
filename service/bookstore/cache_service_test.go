package bookstore

import (
	"context"
	"errors"
	"testing"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// setupTestRedis 设置测试用的Redis miniredis
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

// TestNewRedisCacheService 测试创建缓存服务
func TestNewRedisCacheService(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test")

	assert.NotNil(t, service)
	redisService, ok := service.(*RedisCacheService)
	assert.True(t, ok)
	assert.Equal(t, "test", redisService.prefix)
}

// TestHomepageDataCache 测试首页数据缓存
func TestHomepageDataCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	data := &HomepageData{
		FeaturedBooks: []*bookstoreModel.Book{
			{Title: "推荐书籍1"},
			{Title: "推荐书籍2"},
		},
	}

	// 测试设置缓存
	err := service.SetHomepageData(ctx, data, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取缓存
	retrieved, err := service.GetHomepageData(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, 2, len(retrieved.FeaturedBooks))
	assert.Equal(t, "推荐书籍1", retrieved.FeaturedBooks[0].Title)

	// 测试清除缓存
	err = service.InvalidateHomepageCache(ctx)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetHomepageData(ctx)
	assert.NoError(t, err)
	assert.Nil(t, retrieved) // 缓存不存在
}

// TestBookDetailCache 测试书籍详情缓存
func TestBookDetailCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	bookID := primitive.NewObjectID().Hex()
	bookDetail := &bookstoreModel.BookDetail{
		Title:  "测试书籍",
		Author: "测试作者",
	}

	// 测试设置缓存
	err := service.SetBookDetail(ctx, bookID, bookDetail, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取缓存
	retrieved, err := service.GetBookDetail(ctx, bookID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "测试书籍", retrieved.Title)
	assert.Equal(t, "测试作者", retrieved.Author)

	// 测试清除缓存
	err = service.InvalidateBookDetailCache(ctx, bookID)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetBookDetail(ctx, bookID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestBookCache 测试书籍缓存
func TestBookCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	bookID := primitive.NewObjectID().Hex()
	book := &bookstoreModel.Book{
		Title:  "测试书籍",
		Author: "测试作者",
	}

	// 测试设置缓存
	err := service.SetBook(ctx, bookID, book, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取缓存
	retrieved, err := service.GetBook(ctx, bookID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "测试书籍", retrieved.Title)

	// 测试清除缓存
	err = service.InvalidateBookCache(ctx, bookID)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetBook(ctx, bookID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestChapterCache 测试章节缓存
func TestChapterCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	chapterID := primitive.NewObjectID().Hex()
	chapter := &bookstoreModel.Chapter{
		Title:   "第一章",
		Content: "章节内容",
	}

	// 测试设置缓存
	err := service.SetChapter(ctx, chapterID, chapter, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取缓存
	retrieved, err := service.GetChapter(ctx, chapterID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "第一章", retrieved.Title)
	assert.Equal(t, "章节内容", retrieved.Content)

	// 测试清除缓存
	err = service.InvalidateChapterCache(ctx, chapterID)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetChapter(ctx, chapterID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestRankingCache 测试榜单缓存
func TestRankingCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	rankingType := bookstoreModel.RankingTypeDaily
	period := "2024-01-01"
	items := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
		{BookID: primitive.NewObjectID(), Rank: 2},
	}

	// 测试设置缓存
	err := service.SetRanking(ctx, rankingType, period, items, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取缓存
	retrieved, err := service.GetRanking(ctx, rankingType, period)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, 2, len(retrieved))
	assert.Equal(t, 1, retrieved[0].Rank)

	// 测试清除缓存
	err = service.InvalidateRankingCache(ctx, rankingType, period)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetRanking(ctx, rankingType, period)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestBannerCache 测试Banner缓存
func TestBannerCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	banners := []*bookstoreModel.Banner{
		{Title: "Banner1", ImageURL: "http://example.com/1.jpg"},
		{Title: "Banner2", ImageURL: "http://example.com/2.jpg"},
	}

	// 测试设置缓存
	err := service.SetActiveBanners(ctx, banners, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取缓存
	retrieved, err := service.GetActiveBanners(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, 2, len(retrieved))
	assert.Equal(t, "Banner1", retrieved[0].Title)

	// 测试清除缓存
	err = service.InvalidateBannerCache(ctx)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetActiveBanners(ctx)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestCategoryTreeCache 测试分类树缓存
func TestCategoryTreeCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	tree := []*bookstoreModel.CategoryTree{
		{
			Name: "玄幻",
			Children: []*bookstoreModel.CategoryTree{
				{Name: "东方玄幻"},
				{Name: "西方玄幻"},
			},
		},
		{
			Name: "都市",
			Children: []*bookstoreModel.CategoryTree{
				{Name: "都市生活"},
				{Name: "都市异能"},
			},
		},
	}

	// 测试设置缓存
	err := service.SetCategoryTree(ctx, tree, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取缓存
	retrieved, err := service.GetCategoryTree(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, 2, len(retrieved))
	assert.Equal(t, "玄幻", retrieved[0].Name)
	assert.Equal(t, 2, len(retrieved[0].Children))

	// 测试清除缓存
	err = service.InvalidateCategoryCache(ctx, "玄幻")
	assert.NoError(t, err)
}

// TestBookRatingCache 测试书籍评分缓存
func TestBookRatingCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	bookID := primitive.NewObjectID().Hex()
	rating := &bookstoreModel.BookRating{
		BookID: bookID,
		Rating: 4.5,
	}

	// 测试设置评分缓存
	err := service.SetBookRating(ctx, bookID, rating, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取评分缓存
	retrieved, err := service.GetBookRating(ctx, bookID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, 4.5, retrieved.Rating)

	// 测试设置平均评分缓存
	err = service.SetBookAverageRating(ctx, bookID, 4.5, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取平均评分缓存
	avgRating, err := service.GetBookAverageRating(ctx, bookID)
	assert.NoError(t, err)
	assert.Equal(t, 4.5, avgRating)

	// 测试清除评分缓存
	err = service.InvalidateBookRatingCache(ctx, bookID)
	assert.NoError(t, err)
}

// TestBookStatisticsCache 测试书籍统计缓存
func TestBookStatisticsCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	bookID := primitive.NewObjectID().Hex()
	stats := &bookstoreModel.BookStatistics{
		ViewCount:   1000,
		LikeCount:   100,
		CommentRate: 50,
	}

	// 测试设置统计缓存
	err := service.SetBookStatistics(ctx, bookID, stats, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取统计缓存
	retrieved, err := service.GetBookStatistics(ctx, bookID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, int64(1000), retrieved.ViewCount)
	assert.Equal(t, int64(100), retrieved.LikeCount)

	// 测试清除统计缓存
	err = service.InvalidateBookStatisticsCache(ctx, bookID)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetBookStatistics(ctx, bookID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestTopBooksCache 测试热门书籍缓存
func TestTopBooksCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	books := []*bookstoreModel.Book{
		{Title: "热门书籍1", ViewCount: 1000},
		{Title: "热门书籍2", ViewCount: 900},
	}

	// 测试设置热门浏览书籍缓存
	err := service.SetTopViewedBooks(ctx, books, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取热门浏览书籍缓存
	retrieved, err := service.GetTopViewedBooks(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, 2, len(retrieved))

	// 测试清除热门浏览书籍缓存
	err = service.InvalidateTopViewedBooksCache(ctx)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetTopViewedBooks(ctx)
	assert.Error(t, err) // 缓存不存在，返回错误
}

// TestAggregatedStatisticsCache 测试聚合统计缓存
func TestAggregatedStatisticsCache(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	stats := map[string]interface{}{
		"total_books":      1000,
		"total_authors":    500,
		"total_categories": 50,
	}

	// 测试设置聚合统计缓存
	err := service.SetAggregatedStatistics(ctx, stats, 5*time.Minute)
	assert.NoError(t, err)

	// 测试获取聚合统计缓存
	retrieved, err := service.GetAggregatedStatistics(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, int64(1000), retrieved["total_books"])
	assert.Equal(t, int64(500), retrieved["total_authors"])

	// 测试清除聚合统计缓存
	err = service.InvalidateAggregatedStatisticsCache(ctx)
	assert.NoError(t, err)

	// 验证缓存已清除
	retrieved, err = service.GetAggregatedStatistics(ctx)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestCacheMiss 测试缓存未命中
func TestCacheMiss(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	bookID := primitive.NewObjectID().Hex()

	// 测试获取不存在的书籍详情
	bookDetail, err := service.GetBookDetail(ctx, bookID)
	assert.NoError(t, err)
	assert.Nil(t, bookDetail) // 缓存不存在，返回nil

	// 测试获取不存在的书籍
	book, err := service.GetBook(ctx, bookID)
	assert.NoError(t, err)
	assert.Nil(t, book)

	// 测试获取不存在的章节
	chapter, err := service.GetChapter(ctx, bookID)
	assert.NoError(t, err)
	assert.Nil(t, chapter)
}

// TestCacheExpiration 测试缓存过期
func TestCacheExpiration(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	bookID := primitive.NewObjectID().Hex()
	bookDetail := &bookstoreModel.BookDetail{
		Title:  "测试书籍",
		Author: "测试作者",
	}

	// 设置短过期时间的缓存
	err := service.SetBookDetail(ctx, bookID, bookDetail, 100*time.Millisecond)
	assert.NoError(t, err)

	// 立即获取，应该存在
	retrieved, err := service.GetBookDetail(ctx, bookID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 再次获取，应该不存在
	retrieved, err = service.GetBookDetail(ctx, bookID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

// TestCacheKeyGeneration 测试缓存键生成
func TestCacheKeyGeneration(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)

	// 测试getKey方法
	key1 := service.getKey("homepage")
	assert.Equal(t, "test:bookstore:homepage", key1)

	key2 := service.getKey("book:123456")
	assert.Equal(t, "test:bookstore:book:123456", key2)

	key3 := service.getKey("ranking:daily:2024-01-01")
	assert.Equal(t, "test:bookstore:ranking:daily:2024-01-01", key3)
}

// TestConcurrentCacheAccess 测试并发缓存访问
func TestConcurrentCacheAccess(t *testing.T) {
	mr, client := setupTestRedis(t)
	defer mr.Close()

	service := NewRedisCacheService(client, "test").(*RedisCacheService)
	ctx := context.Background()

	bookID := primitive.NewObjectID().Hex()
	bookDetail := &bookstoreModel.BookDetail{
		Title:  "测试书籍",
		Author: "测试作者",
	}

	// 并发写入
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			err := service.SetBookDetail(ctx, bookID, bookDetail, 5*time.Minute)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证缓存存在
	retrieved, err := service.GetBookDetail(ctx, bookID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
}
