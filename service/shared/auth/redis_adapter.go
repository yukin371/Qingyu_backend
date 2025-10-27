package auth

import (
	"context"
	"time"

	"Qingyu_backend/pkg/cache"
)

// RedisAdapter 将 cache.RedisClient 适配为 auth 所需的接口
type RedisAdapter struct {
	client cache.RedisClient
}

// NewRedisAdapter 创建Redis适配器
func NewRedisAdapter(client cache.RedisClient) *RedisAdapter {
	return &RedisAdapter{client: client}
}

// ============ CacheClient 接口实现 ============

// Get 获取缓存值
func (a *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key)
}

// Set 设置缓存值
func (a *RedisAdapter) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return a.client.Set(ctx, key, value, ttl)
}

// Delete 删除缓存键
func (a *RedisAdapter) Delete(ctx context.Context, key string) error {
	return a.client.Delete(ctx, key)
}

// ============ RedisClient 接口实现 ============

// Exists 检查键是否存在
func (a *RedisAdapter) Exists(ctx context.Context, keys ...string) (int64, error) {
	return a.client.Exists(ctx, keys...)
}

// GetTTL 获取键的过期时间
func (a *RedisAdapter) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return a.client.TTL(ctx, key)
}

// SetWithExpire 设置键值对并指定过期时间
func (a *RedisAdapter) SetWithExpire(ctx context.Context, key, value string, expire time.Duration) error {
	return a.client.Set(ctx, key, value, expire)
}

// SAdd 添加集合成员
func (a *RedisAdapter) SAdd(ctx context.Context, key string, members ...string) error {
	// 将[]string转换为[]interface{}
	memberInterfaces := make([]interface{}, len(members))
	for i, m := range members {
		memberInterfaces[i] = m
	}
	return a.client.SAdd(ctx, key, memberInterfaces...)
}

// SMembers 获取集合所有成员
func (a *RedisAdapter) SMembers(ctx context.Context, key string) ([]string, error) {
	return a.client.SMembers(ctx, key)
}

// SRem 删除集合成员
func (a *RedisAdapter) SRem(ctx context.Context, key string, members ...string) error {
	// 将[]string转换为[]interface{}
	memberInterfaces := make([]interface{}, len(members))
	for i, m := range members {
		memberInterfaces[i] = m
	}
	return a.client.SRem(ctx, key, memberInterfaces...)
}

// Expire 设置键的过期时间
func (a *RedisAdapter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return a.client.Expire(ctx, key, ttl)
}
