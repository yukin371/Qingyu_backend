package cache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

// SearchCache 搜索缓存管理
type SearchCache struct {
	redis      *redis.Client
	defaultTTL time.Duration
	hotTTL     time.Duration
	keyPrefix  string
	enabled    bool
	stats      struct {
		hits   int64
		misses int64
	}
}

// NewSearchCache 创建搜索缓存
func NewSearchCache(redisClient *redis.Client, defaultTTL, hotTTL time.Duration) (*SearchCache, error) {
	return &SearchCache{
		redis:      redisClient,
		defaultTTL: defaultTTL,
		hotTTL:     hotTTL,
		keyPrefix:  "search:cache:",
		enabled:    true,
	}, nil
}

// Get 获取缓存
func (c *SearchCache) Get(ctx context.Context, key string) ([]byte, error) {
	if !c.enabled || c.redis == nil {
		atomic.AddInt64(&c.stats.misses, 1)
		return nil, fmt.Errorf("cache disabled")
	}

	fullKey := c.keyPrefix + key
	val, err := c.redis.Get(ctx, fullKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			atomic.AddInt64(&c.stats.misses, 1)
			return nil, nil // 缓存未命中，返回 nil 而不是错误
		}
		return nil, err
	}

	atomic.AddInt64(&c.stats.hits, 1)
	return val, nil
}

// Set 设置缓存
func (c *SearchCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if !c.enabled || c.redis == nil {
		return nil
	}

	fullKey := c.keyPrefix + key
	return c.redis.Set(ctx, fullKey, value, ttl).Err()
}

// Delete 删除缓存
func (c *SearchCache) Delete(ctx context.Context, key string) error {
	if !c.enabled || c.redis == nil {
		return nil
	}

	fullKey := c.keyPrefix + key
	return c.redis.Del(ctx, fullKey).Err()
}

// DeletePattern 批量删除匹配模式的缓存
func (c *SearchCache) DeletePattern(ctx context.Context, pattern string) error {
	if !c.enabled || c.redis == nil {
		return nil
	}

	fullPattern := c.keyPrefix + pattern
	iter := c.redis.Scan(ctx, 0, fullPattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := c.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Exists 检查缓存是否存在
func (c *SearchCache) Exists(ctx context.Context, key string) (bool, error) {
	if !c.enabled || c.redis == nil {
		return false, nil
	}

	fullKey := c.keyPrefix + key
	n, err := c.redis.Exists(ctx, fullKey).Result()
	return n > 0, err
}

// Clear 清空所有搜索缓存
func (c *SearchCache) Clear(ctx context.Context) error {
	if !c.enabled || c.redis == nil {
		return nil
	}

	fullPattern := c.keyPrefix + "*"
	iter := c.redis.Scan(ctx, 0, fullPattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := c.redis.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Ping 健康检查
func (c *SearchCache) Ping(ctx context.Context) error {
	if c.redis == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.redis.Ping(ctx).Err()
}

// Close 关闭连接
func (c *SearchCache) Close() error {
	// Redis client 由外部管理，不在这里关闭
	return nil
}

// Stats 获取缓存统计信息
func (c *SearchCache) Stats(ctx context.Context) (*CacheStats, error) {
	hits := atomic.LoadInt64(&c.stats.hits)
	misses := atomic.LoadInt64(&c.stats.misses)
	total := hits + misses

	hitRate := 0.0
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	// 获取缓存键数量
	var keys int64
	if c.redis != nil {
		fullPattern := c.keyPrefix + "*"
		iter := c.redis.Scan(ctx, 0, fullPattern, 100).Iterator()
		for iter.Next(ctx) {
			keys++
		}
	}

	return &CacheStats{
		Hits:       hits,
		Misses:     misses,
		Keys:       keys,
		HitRate:    hitRate,
		MemoryUsage: 0, // Redis 内存使用需要通过 INFO 命令获取
	}, nil
}

// Reset 重置统计信息
func (c *SearchCache) Reset(ctx context.Context) error {
	atomic.StoreInt64(&c.stats.hits, 0)
	atomic.StoreInt64(&c.stats.misses, 0)
	return nil
}

// GetTTLForPage 根据页码获取 TTL
func (c *SearchCache) GetTTLForPage(page int) time.Duration {
	// 第一页使用较长的 TTL
	if page == 1 {
		return c.hotTTL
	}
	return c.defaultTTL
}

// GenerateCacheKey 生成缓存键
func GenerateCacheKey(searchType, query string, page, pageSize int, filter map[string]interface{}) string {
	// 将 filter 序列化为 JSON
	filterJSON, _ := json.Marshal(filter)
	filterHash := md5.Sum(filterJSON)

	// 生成缓存键
	return fmt.Sprintf("search:%s:%s:%d:%d:%x", searchType, query, page, pageSize, filterHash)
}

// IsEnabled 检查缓存是否启用
func (c *SearchCache) IsEnabled() bool {
	return c.enabled
}

// SetEnabled 设置缓存启用状态
func (c *SearchCache) SetEnabled(enabled bool) {
	c.enabled = enabled
}
