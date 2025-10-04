package bookstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Qingyu_backend/models/reading/bookstore"

	"github.com/redis/go-redis/v9"
)

// RedisClient 适配接口，便于在测试中使用 MockRedisClient
// 只包含缓存服务需要的最小方法集合
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
}

// CacheService 缓存服务接口
type CacheService interface {
	// 首页数据缓存
	GetHomepageData(ctx context.Context) (*HomepageData, error)
	SetHomepageData(ctx context.Context, data *HomepageData, expiration time.Duration) error

	// 榜单缓存
	GetRanking(ctx context.Context, rankingType bookstore.RankingType, period string) ([]*bookstore.RankingItem, error)
	SetRanking(ctx context.Context, rankingType bookstore.RankingType, period string, items []*bookstore.RankingItem, expiration time.Duration) error

	// Banner缓存
	GetActiveBanners(ctx context.Context) ([]*bookstore.Banner, error)
	SetActiveBanners(ctx context.Context, banners []*bookstore.Banner, expiration time.Duration) error

	// 书籍缓存
	GetBook(ctx context.Context, bookID string) (*bookstore.Book, error)
	SetBook(ctx context.Context, bookID string, book *bookstore.Book, expiration time.Duration) error

	// 书籍详情缓存
	GetBookDetail(ctx context.Context, bookID string) (*bookstore.BookDetail, error)
	SetBookDetail(ctx context.Context, bookID string, bookDetail *bookstore.BookDetail, expiration time.Duration) error

	// 章节缓存
	GetChapter(ctx context.Context, chapterID string) (*bookstore.Chapter, error)
	SetChapter(ctx context.Context, chapterID string, chapter *bookstore.Chapter, expiration time.Duration) error

	// 分类缓存
	GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error)
	SetCategoryTree(ctx context.Context, tree []*bookstore.CategoryTree, expiration time.Duration) error

	// 书籍评分缓存
	GetBookRating(ctx context.Context, key string) (*bookstore.BookRating, error)
	SetBookRating(ctx context.Context, key string, rating *bookstore.BookRating, expiration time.Duration) error
	GetBookAverageRating(ctx context.Context, bookID string) (float64, error)
	SetBookAverageRating(ctx context.Context, bookID string, rating float64, expiration time.Duration) error
	InvalidateBookRatingCache(ctx context.Context, bookID string) error
	InvalidateBookRatingsCache(ctx context.Context, bookID string) error
	InvalidateUserRatingsCache(ctx context.Context, userID string) error
	InvalidateBookAverageRatingCache(ctx context.Context, bookID string) error
	InvalidateTopViewedBooksCache(ctx context.Context) error
	InvalidateHottestBooksCache(ctx context.Context) error
	InvalidateTopFavoritedBooksCache(ctx context.Context) error

	// 书籍统计缓存
	GetBookStatistics(ctx context.Context, bookID string) (*bookstore.BookStatistics, error)
	SetBookStatistics(ctx context.Context, bookID string, stats *bookstore.BookStatistics, expiration time.Duration) error
	InvalidateBookStatisticsCache(ctx context.Context, bookID string) error
	GetTopViewedBooks(ctx context.Context) ([]*bookstore.Book, error)
	SetTopViewedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error
	GetTopFavoritedBooks(ctx context.Context) ([]*bookstore.Book, error)
	SetTopFavoritedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error
	GetTopRatedBooks(ctx context.Context) ([]*bookstore.Book, error)
	SetTopRatedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error
	GetHottestBooks(ctx context.Context) ([]*bookstore.Book, error)
	SetHottestBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error

	// 缓存清理
	InvalidateHomepageCache(ctx context.Context) error
	InvalidateRankingCache(ctx context.Context, rankingType bookstore.RankingType, period string) error
	InvalidateBannerCache(ctx context.Context) error
	InvalidateBookCache(ctx context.Context, bookID string) error
	InvalidateBookDetailCache(ctx context.Context, bookID string) error
	InvalidateChapterCache(ctx context.Context, chapterID string) error
	InvalidateBookChaptersCache(ctx context.Context, bookID string) error
	InvalidateCategoryCache(ctx context.Context, category string) error
	InvalidateAuthorCache(ctx context.Context, author string) error
	InvalidateTopRatedBooksCache(ctx context.Context) error
	InvalidateAggregatedStatisticsCache(ctx context.Context) error

	// 聚合统计缓存
	GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error)
	SetAggregatedStatistics(ctx context.Context, stats map[string]interface{}, expiration time.Duration) error
}

// RedisCacheService Redis缓存服务实现
type RedisCacheService struct {
	client RedisClient
	prefix string
}

// NewRedisCacheService 创建Redis缓存服务
func NewRedisCacheService(client RedisClient, prefix string) CacheService {
	return &RedisCacheService{
		client: client,
		prefix: prefix,
	}
}

// 缓存键生成
func (c *RedisCacheService) getKey(key string) string {
	return fmt.Sprintf("%s:bookstore:%s", c.prefix, key)
}

// GetHomepageData 获取首页数据缓存
func (c *RedisCacheService) GetHomepageData(ctx context.Context) (*HomepageData, error) {
	key := c.getKey("homepage")
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get homepage cache: %w", err)
	}

	var homepage HomepageData
	if err := json.Unmarshal([]byte(data), &homepage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal homepage data: %w", err)
	}

	return &homepage, nil
}

// SetHomepageData 设置首页数据缓存
func (c *RedisCacheService) SetHomepageData(ctx context.Context, data *HomepageData, expiration time.Duration) error {
	key := c.getKey("homepage")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal homepage data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// GetRanking 获取榜单缓存
func (c *RedisCacheService) GetRanking(ctx context.Context, rankingType bookstore.RankingType, period string) ([]*bookstore.RankingItem, error) {
	key := c.getKey(fmt.Sprintf("ranking:%s:%s", rankingType, period))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get ranking cache: %w", err)
	}

	var items []*bookstore.RankingItem
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ranking data: %w", err)
	}

	return items, nil
}

// SetRanking 设置榜单缓存
func (c *RedisCacheService) SetRanking(ctx context.Context, rankingType bookstore.RankingType, period string, items []*bookstore.RankingItem, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("ranking:%s:%s", rankingType, period))
	jsonData, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal ranking data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// GetActiveBanners 获取Banner缓存
func (c *RedisCacheService) GetActiveBanners(ctx context.Context) ([]*bookstore.Banner, error) {
	key := c.getKey("banners:active")
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get banners cache: %w", err)
	}

	var banners []*bookstore.Banner
	if err := json.Unmarshal([]byte(data), &banners); err != nil {
		return nil, fmt.Errorf("failed to unmarshal banners data: %w", err)
	}

	return banners, nil
}

// SetActiveBanners 设置Banner缓存
func (c *RedisCacheService) SetActiveBanners(ctx context.Context, banners []*bookstore.Banner, expiration time.Duration) error {
	key := c.getKey("banners:active")
	jsonData, err := json.Marshal(banners)
	if err != nil {
		return fmt.Errorf("failed to marshal banners data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// GetBook 获取书籍缓存
func (c *RedisCacheService) GetBook(ctx context.Context, bookID string) (*bookstore.Book, error) {
	key := c.getKey(fmt.Sprintf("book:%s", bookID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get book cache: %w", err)
	}

	var book bookstore.Book
	if err := json.Unmarshal([]byte(data), &book); err != nil {
		return nil, fmt.Errorf("failed to unmarshal book data: %w", err)
	}

	return &book, nil
}

// SetBook 设置书籍缓存
func (c *RedisCacheService) SetBook(ctx context.Context, bookID string, book *bookstore.Book, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("book:%s", bookID))
	jsonData, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("failed to marshal book data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// GetCategoryTree 获取分类树缓存
func (c *RedisCacheService) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	key := c.getKey("categories:tree")
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get category tree cache: %w", err)
	}

	var tree []*bookstore.CategoryTree
	if err := json.Unmarshal([]byte(data), &tree); err != nil {
		return nil, fmt.Errorf("failed to unmarshal category tree data: %w", err)
	}

	return tree, nil
}

// SetCategoryTree 设置分类树缓存
func (c *RedisCacheService) SetCategoryTree(ctx context.Context, tree []*bookstore.CategoryTree, expiration time.Duration) error {
	key := c.getKey("categories:tree")
	jsonData, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("failed to marshal category tree data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// InvalidateHomepageCache 清除首页缓存
func (c *RedisCacheService) InvalidateHomepageCache(ctx context.Context) error {
	key := c.getKey("homepage")
	return c.client.Del(ctx, key).Err()
}

// InvalidateRankingCache 清除榜单缓存
func (c *RedisCacheService) InvalidateRankingCache(ctx context.Context, rankingType bookstore.RankingType, period string) error {
	key := c.getKey(fmt.Sprintf("ranking:%s:%s", rankingType, period))
	return c.client.Del(ctx, key).Err()
}

// InvalidateBannerCache 清除Banner缓存
func (c *RedisCacheService) InvalidateBannerCache(ctx context.Context) error {
	key := c.getKey("banners:active")
	return c.client.Del(ctx, key).Err()
}

// InvalidateBookCache 清除书籍缓存
func (c *RedisCacheService) InvalidateBookCache(ctx context.Context, bookID string) error {
	key := c.getKey(fmt.Sprintf("book:%s", bookID))
	return c.client.Del(ctx, key).Err()
}

// InvalidateCategoryCache 清除分类缓存
func (c *RedisCacheService) InvalidateCategoryCache(ctx context.Context, category string) error {
	key := c.getKey(fmt.Sprintf("category:%s", category))
	return c.client.Del(ctx, key).Err()
}

// GetBookDetail 获取书籍详情缓存
func (c *RedisCacheService) GetBookDetail(ctx context.Context, bookID string) (*bookstore.BookDetail, error) {
	key := c.getKey(fmt.Sprintf("book_detail:%s", bookID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get book detail cache: %w", err)
	}

	var bookDetail bookstore.BookDetail
	if err := json.Unmarshal([]byte(data), &bookDetail); err != nil {
		return nil, fmt.Errorf("failed to unmarshal book detail data: %w", err)
	}

	return &bookDetail, nil
}

// SetBookDetail 设置书籍详情缓存
func (c *RedisCacheService) SetBookDetail(ctx context.Context, bookID string, bookDetail *bookstore.BookDetail, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("book_detail:%s", bookID))
	jsonData, err := json.Marshal(bookDetail)
	if err != nil {
		return fmt.Errorf("failed to marshal book detail data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// GetChapter 获取章节缓存
func (c *RedisCacheService) GetChapter(ctx context.Context, chapterID string) (*bookstore.Chapter, error) {
	key := c.getKey(fmt.Sprintf("chapter:%s", chapterID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get chapter cache: %w", err)
	}

	var chapter bookstore.Chapter
	if err := json.Unmarshal([]byte(data), &chapter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chapter data: %w", err)
	}

	return &chapter, nil
}

// SetChapter 设置章节缓存
func (c *RedisCacheService) SetChapter(ctx context.Context, chapterID string, chapter *bookstore.Chapter, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("chapter:%s", chapterID))
	jsonData, err := json.Marshal(chapter)
	if err != nil {
		return fmt.Errorf("failed to marshal chapter data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, expiration).Err()
}

// InvalidateBookDetailCache 清除书籍详情缓存
func (c *RedisCacheService) InvalidateBookDetailCache(ctx context.Context, bookID string) error {
	key := c.getKey(fmt.Sprintf("book_detail:%s", bookID))
	return c.client.Del(ctx, key).Err()
}

// InvalidateChapterCache 清除章节缓存
func (c *RedisCacheService) InvalidateChapterCache(ctx context.Context, chapterID string) error {
	key := c.getKey(fmt.Sprintf("chapter:%s", chapterID))
	return c.client.Del(ctx, key).Err()
}

// InvalidateBookChaptersCache 清除书籍章节列表缓存
func (c *RedisCacheService) InvalidateBookChaptersCache(ctx context.Context, bookID string) error {
	key := c.getKey(fmt.Sprintf("book_chapters:%s", bookID))
	return c.client.Del(ctx, key).Err()
}

// InvalidateAuthorCache 清除作者相关缓存
func (c *RedisCacheService) InvalidateAuthorCache(ctx context.Context, author string) error {
	key := c.getKey(fmt.Sprintf("author:%s", author))
	return c.client.Del(ctx, key).Err()
}

// GetBookRating 获取书籍评分缓存
func (r *RedisCacheService) GetBookRating(ctx context.Context, key string) (*bookstore.BookRating, error) {
	cacheKey := fmt.Sprintf("book_rating:%s", key)
	data, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var rating bookstore.BookRating
	if err := json.Unmarshal([]byte(data), &rating); err != nil {
		return nil, err
	}

	return &rating, nil
}

// SetBookRating 设置书籍评分缓存
func (r *RedisCacheService) SetBookRating(ctx context.Context, key string, rating *bookstore.BookRating, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("book_rating:%s", key)
	data, err := json.Marshal(rating)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, expiration).Err()
}

// GetBookAverageRating 获取书籍平均评分缓存
func (r *RedisCacheService) GetBookAverageRating(ctx context.Context, bookID string) (float64, error) {
	cacheKey := fmt.Sprintf("book_avg_rating:%s", bookID)
	result, err := r.client.Get(ctx, cacheKey).Float64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return result, nil
}

// SetBookAverageRating 设置书籍平均评分缓存
func (r *RedisCacheService) SetBookAverageRating(ctx context.Context, bookID string, rating float64, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("book_avg_rating:%s", bookID)
	return r.client.Set(ctx, cacheKey, rating, expiration).Err()
}

// InvalidateBookRatingCache 清除书籍评分缓存
func (r *RedisCacheService) InvalidateBookRatingCache(ctx context.Context, bookID string) error {
	pattern := fmt.Sprintf("book_rating:*%s*", bookID)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	// 同时清除平均评分缓存
	avgRatingKey := fmt.Sprintf("book_avg_rating:%s", bookID)
	return r.client.Del(ctx, avgRatingKey).Err()
}

// GetBookStatistics 获取书籍统计缓存
func (r *RedisCacheService) GetBookStatistics(ctx context.Context, bookID string) (*bookstore.BookStatistics, error) {
	cacheKey := fmt.Sprintf("book_statistics:%s", bookID)
	data, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var stats bookstore.BookStatistics
	if err := json.Unmarshal([]byte(data), &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// SetBookStatistics 设置书籍统计缓存
func (r *RedisCacheService) SetBookStatistics(ctx context.Context, bookID string, stats *bookstore.BookStatistics, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("book_statistics:%s", bookID)
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, expiration).Err()
}

// InvalidateBookStatisticsCache 清除书籍统计缓存
func (r *RedisCacheService) InvalidateBookStatisticsCache(ctx context.Context, bookID string) error {
	cacheKey := fmt.Sprintf("book_statistics:%s", bookID)
	return r.client.Del(ctx, cacheKey).Err()
}

// InvalidateBookRatingsCache 清除书籍评分缓存
func (r *RedisCacheService) InvalidateBookRatingsCache(ctx context.Context, bookID string) error {
	cacheKey := fmt.Sprintf("book_ratings:%s", bookID)
	return r.client.Del(ctx, cacheKey).Err()
}

// InvalidateUserRatingsCache 清除用户评分缓存
func (r *RedisCacheService) InvalidateUserRatingsCache(ctx context.Context, userID string) error {
	cacheKey := fmt.Sprintf("user_ratings:%s", userID)
	return r.client.Del(ctx, cacheKey).Err()
}

// InvalidateBookAverageRatingCache 清除书籍平均评分缓存
func (r *RedisCacheService) InvalidateBookAverageRatingCache(ctx context.Context, bookID string) error {
	cacheKey := fmt.Sprintf("book_average_rating:%s", bookID)
	return r.client.Del(ctx, cacheKey).Err()
}

// GetTopViewedBooks 获取热门浏览书籍缓存
func (r *RedisCacheService) GetTopViewedBooks(ctx context.Context) ([]*bookstore.Book, error) {
	cacheKey := "top_viewed_books"
	result := r.client.Get(ctx, cacheKey)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var books []*bookstore.Book
	if err := json.Unmarshal([]byte(result.Val()), &books); err != nil {
		return nil, err
	}

	return books, nil
}

// SetTopViewedBooks 设置热门浏览书籍缓存
func (r *RedisCacheService) SetTopViewedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	cacheKey := "top_viewed_books"
	data, err := json.Marshal(books)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, expiration).Err()
}

// GetTopFavoritedBooks 获取热门收藏书籍缓存
func (r *RedisCacheService) GetTopFavoritedBooks(ctx context.Context) ([]*bookstore.Book, error) {
	cacheKey := "top_favorited_books"
	result := r.client.Get(ctx, cacheKey)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var books []*bookstore.Book
	if err := json.Unmarshal([]byte(result.Val()), &books); err != nil {
		return nil, err
	}

	return books, nil
}

// SetTopFavoritedBooks 设置热门收藏书籍缓存
func (r *RedisCacheService) SetTopFavoritedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	cacheKey := "top_favorited_books"
	data, err := json.Marshal(books)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, expiration).Err()
}

// GetTopRatedBooks 获取热门评分书籍缓存
func (r *RedisCacheService) GetTopRatedBooks(ctx context.Context) ([]*bookstore.Book, error) {
	cacheKey := "top_rated_books"
	result := r.client.Get(ctx, cacheKey)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var books []*bookstore.Book
	if err := json.Unmarshal([]byte(result.Val()), &books); err != nil {
		return nil, err
	}

	return books, nil
}

// SetTopRatedBooks 设置热门评分书籍缓存
func (r *RedisCacheService) SetTopRatedBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	cacheKey := "top_rated_books"
	data, err := json.Marshal(books)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, expiration).Err()
}

// GetHottestBooks 获取最热门书籍缓存
func (r *RedisCacheService) GetHottestBooks(ctx context.Context) ([]*bookstore.Book, error) {
	cacheKey := "hottest_books"
	result := r.client.Get(ctx, cacheKey)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var books []*bookstore.Book
	if err := json.Unmarshal([]byte(result.Val()), &books); err != nil {
		return nil, err
	}

	return books, nil
}

// SetHottestBooks 设置最热门书籍缓存
func (r *RedisCacheService) SetHottestBooks(ctx context.Context, books []*bookstore.Book, expiration time.Duration) error {
	cacheKey := "hottest_books"
	data, err := json.Marshal(books)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, expiration).Err()
}

// InvalidateTopViewedBooksCache 清除热门浏览书籍缓存
func (r *RedisCacheService) InvalidateTopViewedBooksCache(ctx context.Context) error {
	cacheKey := "top_viewed_books"
	return r.client.Del(ctx, cacheKey).Err()
}

// InvalidateHottestBooksCache 清除最热门书籍缓存
func (r *RedisCacheService) InvalidateHottestBooksCache(ctx context.Context) error {
	cacheKey := "hottest_books"
	return r.client.Del(ctx, cacheKey).Err()
}

// InvalidateTopFavoritedBooksCache 清除热门收藏书籍缓存
func (r *RedisCacheService) InvalidateTopFavoritedBooksCache(ctx context.Context) error {
	cacheKey := "top_favorited_books"
	return r.client.Del(ctx, cacheKey).Err()
}

// InvalidateTopRatedBooksCache 清除热门评分书籍缓存
func (r *RedisCacheService) InvalidateTopRatedBooksCache(ctx context.Context) error {
	cacheKey := "top_rated_books"
	return r.client.Del(ctx, cacheKey).Err()
}

// InvalidateAggregatedStatisticsCache 清除聚合统计缓存
func (r *RedisCacheService) InvalidateAggregatedStatisticsCache(ctx context.Context) error {
	cacheKey := "aggregated_statistics"
	return r.client.Del(ctx, cacheKey).Err()
}

// GetAggregatedStatistics 获取聚合统计缓存
func (r *RedisCacheService) GetAggregatedStatistics(ctx context.Context) (map[string]interface{}, error) {
	cacheKey := "aggregated_statistics"
	result := r.client.Get(ctx, cacheKey)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, nil
		}
		return nil, result.Err()
	}

	var stats map[string]interface{}
	if err := json.Unmarshal([]byte(result.Val()), &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// SetAggregatedStatistics 设置聚合统计缓存
func (r *RedisCacheService) SetAggregatedStatistics(ctx context.Context, stats map[string]interface{}, expiration time.Duration) error {
	cacheKey := "aggregated_statistics"
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, expiration).Err()
}
