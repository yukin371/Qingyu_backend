package shared

import (
	"context"
	"time"
)

// TokenBlacklistRepository JWT Token 黑名单 Repository
type TokenBlacklistRepository interface {
	// AddToBlacklist 将 Token 加入黑名单
	AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error

	// IsBlacklisted 检查 Token 是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)

	// RemoveFromBlacklist 从黑名单中移除 Token（一般不需要，Token 会自动过期）
	RemoveFromBlacklist(ctx context.Context, token string) error

	// ClearExpiredTokens 清理过期的 Token（Redis 自动过期，此方法可选）
	ClearExpiredTokens(ctx context.Context) error

	// Health 健康检查
	Health(ctx context.Context) error
}
