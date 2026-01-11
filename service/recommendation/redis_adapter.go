package recommendation

import (
	"context"
	"time"

	"Qingyu_backend/pkg/cache"
)

// RedisAdapter 将 cache.RedisClient 适配为 recommendation 所需的 CacheClient 接口
type RedisAdapter struct {
	client cache.RedisClient
}

// NewRedisAdapter 创建Redis适配器
func NewRedisAdapter(client cache.RedisClient) *RedisAdapter {
	return &RedisAdapter{client: client}
}

// Get 获取缓存值
func (a *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key)
}

// Set 设置缓存值
func (a *RedisAdapter) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return a.client.Set(ctx, key, value, ttl)
}

// Delete 删除缓存键（单个key版本）
func (a *RedisAdapter) Delete(ctx context.Context, key string) error {
	return a.client.Delete(ctx, key)
}
