package recommendation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RecommendationCacheService 推荐结果缓存服务接口
type RecommendationCacheService interface {
	// 个性化推荐缓存
	GetPersonalizedCache(ctx context.Context, userID string) ([]string, error)
	SetPersonalizedCache(ctx context.Context, userID string, items []string, expiration time.Duration) error
	InvalidatePersonalizedCache(ctx context.Context, userID string) error

	// 相似物品缓存
	GetSimilarItemsCache(ctx context.Context, itemID string) ([]string, error)
	SetSimilarItemsCache(ctx context.Context, itemID string, items []string, expiration time.Duration) error
	InvalidateSimilarItemsCache(ctx context.Context, itemID string) error
}

// RedisRecommendationCacheService Redis缓存实现
type RedisRecommendationCacheService struct {
	client *redis.Client
	prefix string
}

// NewRedisRecommendationCacheService 创建Redis推荐缓存服务
func NewRedisRecommendationCacheService(client *redis.Client, prefix string) RecommendationCacheService {
	if prefix == "" {
		prefix = "qingyu"
	}
	return &RedisRecommendationCacheService{
		client: client,
		prefix: prefix,
	}
}

func (c *RedisRecommendationCacheService) getKey(key string) string {
	return fmt.Sprintf("%s:reco:%s", c.prefix, key)
}

// GetPersonalizedCache 获取个性化推荐缓存
func (c *RedisRecommendationCacheService) GetPersonalizedCache(ctx context.Context, userID string) ([]string, error) {
	key := c.getKey(fmt.Sprintf("personalized:%s", userID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在
		}
		return nil, fmt.Errorf("failed to get personalized cache: %w", err)
	}

	var items []string
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return items, nil
}

// SetPersonalizedCache 设置个性化推荐缓存
func (c *RedisRecommendationCacheService) SetPersonalizedCache(ctx context.Context, userID string, items []string, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("personalized:%s", userID))
	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// InvalidatePersonalizedCache 清除个性化推荐缓存
func (c *RedisRecommendationCacheService) InvalidatePersonalizedCache(ctx context.Context, userID string) error {
	key := c.getKey(fmt.Sprintf("personalized:%s", userID))
	return c.client.Del(ctx, key).Err()
}

// GetSimilarItemsCache 获取相似物品缓存
func (c *RedisRecommendationCacheService) GetSimilarItemsCache(ctx context.Context, itemID string) ([]string, error) {
	key := c.getKey(fmt.Sprintf("similar:%s", itemID))
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get similar items cache: %w", err)
	}

	var items []string
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return items, nil
}

// SetSimilarItemsCache 设置相似物品缓存
func (c *RedisRecommendationCacheService) SetSimilarItemsCache(ctx context.Context, itemID string, items []string, expiration time.Duration) error {
	key := c.getKey(fmt.Sprintf("similar:%s", itemID))
	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// InvalidateSimilarItemsCache 清除相似物品缓存
func (c *RedisRecommendationCacheService) InvalidateSimilarItemsCache(ctx context.Context, itemID string) error {
	key := c.getKey(fmt.Sprintf("similar:%s", itemID))
	return c.client.Del(ctx, key).Err()
}
