package redis

import (
	"context"
	"fmt"
	"time"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklistRepositoryRedis JWT Token 黑名单 Repository 的 Redis 实现
type TokenBlacklistRepositoryRedis struct {
	client *redis.Client
	prefix string
}

// NewTokenBlacklistRepository 创建 Token 黑名单 Repository
func NewTokenBlacklistRepository(client *redis.Client) sharedRepo.TokenBlacklistRepository {
	return &TokenBlacklistRepositoryRedis{
		client: client,
		prefix: "token:blacklist:",
	}
}

// AddToBlacklist 将 Token 加入黑名单
func (r *TokenBlacklistRepositoryRedis) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := r.prefix + token

	// 使用 SET 命令设置键值，并设置过期时间
	err := r.client.Set(ctx, key, "1", expiration).Err()
	if err != nil {
		return fmt.Errorf("添加Token到黑名单失败: %w", err)
	}

	return nil
}

// IsBlacklisted 检查 Token 是否在黑名单中
func (r *TokenBlacklistRepositoryRedis) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := r.prefix + token

	// 检查键是否存在
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查Token黑名单状态失败: %w", err)
	}

	return exists > 0, nil
}

// RemoveFromBlacklist 从黑名单中移除 Token
func (r *TokenBlacklistRepositoryRedis) RemoveFromBlacklist(ctx context.Context, token string) error {
	key := r.prefix + token

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("从黑名单移除Token失败: %w", err)
	}

	return nil
}

// ClearExpiredTokens 清理过期的 Token
// Redis 会自动清理过期键，此方法为空实现
func (r *TokenBlacklistRepositoryRedis) ClearExpiredTokens(ctx context.Context) error {
	// Redis 自动处理过期键，无需手动清理
	return nil
}

// Health 健康检查
func (r *TokenBlacklistRepositoryRedis) Health(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
