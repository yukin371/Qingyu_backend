package cache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

// SearchCache 搜索缓存管理
type SearchCache struct {
	// TODO: 注入 Redis client
	defaultTTL time.Duration
	hotTTL     time.Duration
}

// NewSearchCache 创建搜索缓存
func NewSearchCache(defaultTTL, hotTTL time.Duration) (*SearchCache, error) {
	return &SearchCache{
		defaultTTL: defaultTTL,
		hotTTL:     hotTTL,
	}, nil
}

// Get 获取缓存
func (c *SearchCache) Get(ctx context.Context, key string) ([]byte, error) {
	// TODO: 实现 Redis 获取
	return nil, fmt.Errorf("cache get not implemented")
}

// Set 设置缓存
func (c *SearchCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	// TODO: 实现 Redis 设置
	return fmt.Errorf("cache set not implemented")
}

// Delete 删除缓存
func (c *SearchCache) Delete(ctx context.Context, key string) error {
	// TODO: 实现 Redis 删除
	return fmt.Errorf("cache delete not implemented")
}

// GenerateCacheKey 生成缓存键
func GenerateCacheKey(searchType, query string, page, pageSize int, filter map[string]interface{}) string {
	// 将 filter 序列化为 JSON
	filterJSON, _ := json.Marshal(filter)
	filterHash := md5.Sum(filterJSON)

	// 生成缓存键
	return fmt.Sprintf("search:%s:%s:%d:%d:%x", searchType, query, page, pageSize, filterHash)
}
