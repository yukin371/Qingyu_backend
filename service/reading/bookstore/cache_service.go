package bookstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"Qingyu_backend/models/reading/bookstore"
)

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
	
	// 分类缓存
	GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error)
	SetCategoryTree(ctx context.Context, tree []*bookstore.CategoryTree, expiration time.Duration) error
	
	// 缓存清理
	InvalidateHomepageCache(ctx context.Context) error
	InvalidateRankingCache(ctx context.Context, rankingType bookstore.RankingType, period string) error
	InvalidateBannerCache(ctx context.Context) error
	InvalidateBookCache(ctx context.Context, bookID string) error
	InvalidateCategoryCache(ctx context.Context) error
}

// RedisCacheService Redis缓存服务实现
type RedisCacheService struct {
	client *redis.Client
	prefix string
}

// NewRedisCacheService 创建Redis缓存服务
func NewRedisCacheService(client *redis.Client, prefix string) CacheService {
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
func (c *RedisCacheService) InvalidateCategoryCache(ctx context.Context) error {
	key := c.getKey("categories:tree")
	return c.client.Del(ctx, key).Err()
}