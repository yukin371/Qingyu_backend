package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	authModel "Qingyu_backend/models/auth"
)

// OAuthRepository OAuth仓储接口
type OAuthRepository interface {
	// ==================== OAuth账号管理 ====================

	// FindByProviderAndProviderID 根据提供商和提供商用户ID查找OAuth账号
	FindByProviderAndProviderID(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error)

	// FindByUserID 根据用户ID查找所有OAuth账号
	FindByUserID(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error)

	// FindByID 根据ID查找OAuth账号
	FindByID(ctx context.Context, id string) (*authModel.OAuthAccount, error)

	// Create 创建OAuth账号
	Create(ctx context.Context, account *authModel.OAuthAccount) error

	// Update 更新OAuth账号
	Update(ctx context.Context, account *authModel.OAuthAccount) error

	// Delete 删除OAuth账号
	Delete(ctx context.Context, id string) error

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, id string) error

	// UpdateTokens 更新访问令牌
	UpdateTokens(ctx context.Context, id string, accessToken, refreshToken string, expiresAt primitive.DateTime) error

	// SetPrimaryAccount 设置主账号
	SetPrimaryAccount(ctx context.Context, userID string, accountID string) error

	// GetPrimaryAccount 获取用户的主账号
	GetPrimaryAccount(ctx context.Context, userID string) (*authModel.OAuthAccount, error)

	// CountByUserID 统计用户的OAuth账号数量
	CountByUserID(ctx context.Context, userID string) (int64, error)

	// ==================== OAuth会话管理 ====================

	// CreateSession 创建OAuth会话
	CreateSession(ctx context.Context, session *authModel.OAuthSession) error

	// FindSessionByID 根据ID查找OAuth会话
	FindSessionByID(ctx context.Context, id string) (*authModel.OAuthSession, error)

	// FindSessionByState 根据state查找OAuth会话
	FindSessionByState(ctx context.Context, state string) (*authModel.OAuthSession, error)

	// DeleteSession 删除OAuth会话
	DeleteSession(ctx context.Context, id string) error

	// CleanupExpiredSessions 清理过期的OAuth会话
	CleanupExpiredSessions(ctx context.Context) (int64, error)
}
