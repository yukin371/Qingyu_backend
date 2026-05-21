package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/config"
	internalAuth "Qingyu_backend/internal/middleware/auth"
)

// JWTService接口已在interfaces.go中定义，这里直接实现

// JWTServiceImpl JWT服务实现
type JWTServiceImpl struct {
	config      *config.JWTConfigEnhanced
	redisClient RedisClient // Redis客户端接口（用于黑名单）
}

// RedisClient Redis客户端接口（简化版）
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
}

// NewJWTService 创建JWT服务
func NewJWTService(cfg *config.JWTConfigEnhanced, redisClient RedisClient) JWTService {
	if cfg == nil {
		cfg = config.GetJWTConfigEnhanced()
	}
	return &JWTServiceImpl{
		config:      cfg,
		redisClient: redisClient,
	}
}

func (s *JWTServiceImpl) newJWTManager() (internalAuth.JWTManager, error) {
	cfg := s.config
	if config.GlobalConfig != nil && config.GlobalConfig.JWT != nil && config.GlobalConfig.JWT.Secret != "" {
		liveCfg := config.GetJWTConfigEnhanced()
		if liveCfg != nil && liveCfg.SecretKey != "" {
			cfg = liveCfg
			s.config = liveCfg
		}
	}

	return internalAuth.NewJWTManager(cfg.SecretKey, cfg.Expiration, cfg.RefreshDuration)
}

func (s *JWTServiceImpl) buildExtraClaims(roles []string) map[string]interface{} {
	if len(roles) == 0 {
		return nil
	}
	return map[string]interface{}{"roles": roles}
}

func (s *JWTServiceImpl) buildTokenClaims(claims *internalAuth.Claims) *TokenClaims {
	if claims == nil {
		return nil
	}

	result := &TokenClaims{
		UserID: claims.UserID,
		Roles:  append([]string(nil), claims.Roles...),
	}
	if claims.ExpiresAt != nil {
		result.Exp = claims.ExpiresAt.Unix()
	}
	if claims.IssuedAt != nil {
		result.Iat = claims.IssuedAt.Unix()
	}
	return result
}

// ============ Token生成 ============

// GenerateToken 生成访问Token（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) GenerateToken(ctx context.Context, userID string, roles []string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id不能为空")
	}

	manager, err := s.newJWTManager()
	if err != nil {
		return "", err
	}

	accessToken, _, err := manager.GenerateTokenPair(userID, s.buildExtraClaims(roles))
	if err != nil {
		return "", fmt.Errorf("生成访问Token失败: %w", err)
	}
	return accessToken, nil
}

// GenerateTokenWithUsername 生成带用户名的访问Token（扩展方法）
func (s *JWTServiceImpl) GenerateTokenWithUsername(ctx context.Context, userID string, username string, roles []string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id不能为空")
	}

	manager, err := s.newJWTManager()
	if err != nil {
		return "", err
	}

	extraClaims := s.buildExtraClaims(roles)
	if extraClaims == nil {
		extraClaims = make(map[string]interface{})
	}
	extraClaims["username"] = username

	accessToken, _, err := manager.GenerateTokenPair(userID, extraClaims)
	if err != nil {
		return "", fmt.Errorf("生成访问Token失败: %w", err)
	}
	return accessToken, nil
}

// GenerateTokenPair 生成Token对（访问Token + 刷新Token）- 扩展方法
func (s *JWTServiceImpl) GenerateTokenPair(ctx context.Context, userID string, roles []string) (accessToken, refreshToken string, err error) {
	manager, err := s.newJWTManager()
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = manager.GenerateTokenPair(userID, s.buildExtraClaims(roles))
	if err != nil {
		return "", "", fmt.Errorf("生成Token对失败: %w", err)
	}
	return accessToken, refreshToken, nil
}

// ============ Token验证 ============

// ValidateToken 验证Token并返回Claims（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token为空")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	manager, err := s.newJWTManager()
	if err != nil {
		return nil, err
	}

	claims, err := manager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("token解析失败: %w", err)
	}

	if s.redisClient != nil {
		revoked, revokeErr := s.IsTokenRevoked(ctx, tokenString)
		if revokeErr == nil && revoked {
			return nil, errors.New("token已被吊销")
		}
	}

	return s.buildTokenClaims(claims), nil
}

// ============ Token刷新 ============

// RefreshToken 使用刷新Token获取新的Token（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) RefreshToken(ctx context.Context, token string) (string, error) {
	claims, err := s.ValidateToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("刷新Token验证失败: %w", err)
	}

	newToken, err := s.GenerateToken(ctx, claims.UserID, claims.Roles)
	if err != nil {
		return "", fmt.Errorf("生成新Token失败: %w", err)
	}

	return newToken, nil
}

// ============ Token吊销（黑名单） ============

// RevokeToken 吊销Token（加入黑名单）- 匹配interfaces.go中的定义
func (s *JWTServiceImpl) RevokeToken(ctx context.Context, token string) error {
	if s.redisClient == nil {
		return errors.New("Redis客户端未配置，无法吊销Token")
	}

	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	claims, err := s.ParseTokenClaims(token)
	if err != nil {
		return nil
	}

	var remainingTime time.Duration
	if claims.Exp > 1e15 {
		remainingTime = time.Until(time.Unix(0, claims.Exp))
	} else {
		remainingTime = time.Until(time.Unix(claims.Exp, 0))
	}

	if remainingTime <= 0 {
		return nil
	}

	key := s.getBlacklistKey(token)
	if err := s.redisClient.Set(ctx, key, "revoked", remainingTime); err != nil {
		return fmt.Errorf("存储黑名单记录失败: %w", err)
	}

	return nil
}

// IsTokenRevoked 检查Token是否已被吊销
func (s *JWTServiceImpl) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	if s.redisClient == nil {
		return false, nil
	}

	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	key := s.getBlacklistKey(token)
	exists, err := s.redisClient.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("检查黑名单失败: %w", err)
	}

	return exists > 0, nil
}

// getBlacklistKey 获取黑名单Redis key
func (s *JWTServiceImpl) getBlacklistKey(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("token:blacklist:%s", tokenHash[:32])
}

// ============ 辅助方法 ============

// ParseTokenClaims 解析Token的Claims（不验证签名）- 内部方法
func (s *JWTServiceImpl) ParseTokenClaims(token string) (*TokenClaims, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	manager, err := s.newJWTManager()
	if err != nil {
		return nil, err
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf("解析claims失败: %w", err)
	}

	return s.buildTokenClaims(claims), nil
}
