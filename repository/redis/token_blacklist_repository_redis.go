package redis

import (
	"context"
	"fmt"
	"time"

	sharedRepo "Qingyu_backend/repository/interfaces/shared"

	"github.com/redis/go-redis/v9"
)

const (
	// DefaultTokenBlacklistPrefix 默认的 Token 黑名单键前缀
	DefaultTokenBlacklistPrefix = "token:blacklist:"

	// TokenBlacklistValue Token 黑名单键的值（仅用于标识存在）
	TokenBlacklistValue = "1"
)

// TokenBlacklistConfig Token 黑名单配置
type TokenBlacklistConfig struct {
	// Prefix Redis 键前缀
	Prefix string
}

// DefaultTokenBlacklistConfig 返回默认配置
func DefaultTokenBlacklistConfig() *TokenBlacklistConfig {
	return &TokenBlacklistConfig{
		Prefix: DefaultTokenBlacklistPrefix,
	}
}

// TokenBlacklistRepositoryRedis JWT Token 黑名单 Repository 的 Redis 实现
type TokenBlacklistRepositoryRedis struct {
	client *redis.Client
	config *TokenBlacklistConfig
}

// NewTokenBlacklistRepository 创建 Token 黑名单 Repository（使用默认配置）
func NewTokenBlacklistRepository(client *redis.Client) sharedRepo.TokenBlacklistRepository {
	return NewTokenBlacklistRepositoryWithConfig(client, DefaultTokenBlacklistConfig())
}

// NewTokenBlacklistRepositoryWithConfig 创建 Token 黑名单 Repository（使用自定义配置）
func NewTokenBlacklistRepositoryWithConfig(client *redis.Client, config *TokenBlacklistConfig) sharedRepo.TokenBlacklistRepository {
	if config == nil {
		config = DefaultTokenBlacklistConfig()
	}
	if config.Prefix == "" {
		config.Prefix = DefaultTokenBlacklistPrefix
	}

	return &TokenBlacklistRepositoryRedis{
		client: client,
		config: config,
	}
}

// AddToBlacklist 将 Token 加入黑名单
func (r *TokenBlacklistRepositoryRedis) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	if token == "" {
		return fmt.Errorf("token不能为空")
	}

	key := r.config.Prefix + token

	// 使用 SET 命令设置键值，并设置过期时间
	err := r.client.Set(ctx, key, TokenBlacklistValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("添加Token到黑名单失败: %w", err)
	}

	return nil
}

// IsBlacklisted 检查 Token 是否在黑名单中
func (r *TokenBlacklistRepositoryRedis) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("token不能为空")
	}

	key := r.config.Prefix + token

	// 检查键是否存在
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查Token黑名单状态失败: %w", err)
	}

	return exists > 0, nil
}

// RemoveFromBlacklist 从黑名单中移除 Token
func (r *TokenBlacklistRepositoryRedis) RemoveFromBlacklist(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("token不能为空")
	}

	key := r.config.Prefix + token

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
	if r.client == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis健康检查失败: %w", err)
	}
	return nil
}
