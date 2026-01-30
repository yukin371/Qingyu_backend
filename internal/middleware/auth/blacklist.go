package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Blacklist Token黑名单接口
type Blacklist interface {
	// Add 添加Token到黑名单
	Add(ctx context.Context, token string, ttl time.Duration) error

	// IsBlacklisted 检查Token是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)

	// Remove 从黑名单移除Token
	Remove(ctx context.Context, token string) error
}

// redisBlacklist Redis实现的Token黑名单
type redisBlacklist struct {
	client *redis.Client
	prefix string
}

// NewRedisBlacklist 创建Redis黑名单
func NewRedisBlacklist(client *redis.Client, prefix string) Blacklist {
	if prefix == "" {
		prefix = "blacklist:"
	}

	return &redisBlacklist{
		client: client,
		prefix: prefix,
	}
}

// Add 添加Token到黑名单
func (b *redisBlacklist) Add(ctx context.Context, token string, ttl time.Duration) error {
	key := b.prefix + token
	return b.client.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted 检查Token是否在黑名单中
func (b *redisBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := b.prefix + token
	result, err := b.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}
	return result > 0, nil
}

// Remove 从黑名单移除Token
func (b *redisBlacklist) Remove(ctx context.Context, token string) error {
	key := b.prefix + token
	return b.client.Del(ctx, key).Err()
}

// mockBlacklist 用于测试的模拟黑名单
type mockBlacklist struct {
	tokens map[string]time.Time
}

// NewMockBlacklist 创建模拟黑名单
func NewMockBlacklist() Blacklist {
	return &mockBlacklist{
		tokens: make(map[string]time.Time),
	}
}

// Add 添加Token到黑名单
func (b *mockBlacklist) Add(ctx context.Context, token string, ttl time.Duration) error {
	b.tokens[token] = time.Now().Add(ttl)
	return nil
}

// IsBlacklisted 检查Token是否在黑名单中
func (b *mockBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	_, exists := b.tokens[token]
	return exists, nil
}

// Remove 从黑名单移除Token
func (b *mockBlacklist) Remove(ctx context.Context, token string) error {
	delete(b.tokens, token)
	return nil
}
