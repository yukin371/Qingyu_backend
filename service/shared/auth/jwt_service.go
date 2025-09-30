package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/config"
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

// ============ Token生成 ============

// GenerateToken 生成访问Token（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) GenerateToken(ctx context.Context, userID string, roles []string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id不能为空")
	}

	now := time.Now()
	claims := &TokenClaims{
		UserID: userID,
		Roles:  roles,
		Exp:    now.Add(s.config.Expiration).Unix(),
	}

	return s.generateJWT(claims)
}

// GenerateTokenWithUsername 生成带用户名的访问Token（扩展方法）
func (s *JWTServiceImpl) GenerateTokenWithUsername(ctx context.Context, userID string, username string, roles []string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id不能为空")
	}

	now := time.Now()
	claims := &TokenClaims{
		UserID: userID,
		Roles:  roles,
		Exp:    now.Add(s.config.Expiration).Unix(),
	}

	return s.generateJWT(claims)
}

// GenerateTokenPair 生成Token对（访问Token + 刷新Token）- 扩展方法
func (s *JWTServiceImpl) GenerateTokenPair(ctx context.Context, userID string, roles []string) (accessToken, refreshToken string, err error) {
	// 1. 生成访问Token
	accessToken, err = s.GenerateToken(ctx, userID, roles)
	if err != nil {
		return "", "", fmt.Errorf("生成访问Token失败: %w", err)
	}

	// 2. 生成刷新Token
	now := time.Now()
	refreshClaims := &TokenClaims{
		UserID: userID,
		Roles:  roles,
		Exp:    now.Add(s.config.RefreshDuration).Unix(),
	}

	refreshToken, err = s.generateJWT(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("生成刷新Token失败: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateJWT 内部方法：生成JWT
func (s *JWTServiceImpl) generateJWT(claims *TokenClaims) (string, error) {
	// 1. 创建Header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// 2. 创建Payload
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("序列化claims失败: %w", err)
	}
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// 3. 创建签名
	message := headerBase64 + "." + payloadBase64
	signature := s.createSignature(message)

	// 4. 组合JWT
	token := message + "." + signature

	return token, nil
}

// createSignature 创建HMAC-SHA256签名
func (s *JWTServiceImpl) createSignature(message string) string {
	h := hmac.New(sha256.New, []byte(s.config.SecretKey))
	h.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return signature
}

// ============ Token验证 ============

// ValidateToken 验证Token并返回Claims（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	// 1. 检查Token格式
	if token == "" {
		return nil, errors.New("token为空")
	}

	// 移除Bearer前缀
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	// 2. 解析Token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("token格式错误")
	}

	headerBase64 := parts[0]
	payloadBase64 := parts[1]
	signatureBase64 := parts[2]

	// 3. 验证签名
	message := headerBase64 + "." + payloadBase64
	expectedSignature := s.createSignature(message)
	if signatureBase64 != expectedSignature {
		return nil, errors.New("token签名验证失败")
	}

	// 4. 解析Claims
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return nil, fmt.Errorf("解码payload失败: %w", err)
	}

	var claims TokenClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("解析claims失败: %w", err)
	}

	// 5. 验证过期时间
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token已过期")
	}

	// 7. 检查Token是否在黑名单中
	if s.redisClient != nil {
		revoked, err := s.IsTokenRevoked(ctx, token)
		if err == nil && revoked {
			return nil, errors.New("token已被吊销")
		}
	}

	return &claims, nil
}

// ============ Token刷新 ============

// RefreshToken 使用刷新Token获取新的Token（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// 1. 验证刷新Token
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("刷新Token验证失败: %w", err)
	}

	// 2. 生成新的访问Token
	newToken, err := s.GenerateToken(ctx, claims.UserID, claims.Roles)
	if err != nil {
		return "", fmt.Errorf("生成新Token失败: %w", err)
	}

	// 3. 可选：将旧的刷新Token加入黑名单
	if s.redisClient != nil {
		remainingTime := time.Until(time.Unix(claims.Exp, 0))
		if remainingTime > 0 {
			_ = s.RevokeToken(ctx, refreshToken)
		}
	}

	return newToken, nil
}

// ============ Token吊销（黑名单） ============

// RevokeToken 吊销Token（加入黑名单）- 匹配interfaces.go中的定义
func (s *JWTServiceImpl) RevokeToken(ctx context.Context, token string) error {
	if s.redisClient == nil {
		return errors.New("Redis客户端未配置，无法吊销Token")
	}

	// 移除Bearer前缀
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	// 解析Token获取过期时间
	claims, err := s.ParseTokenClaims(token)
	if err != nil {
		return fmt.Errorf("解析Token失败: %w", err)
	}

	// 计算剩余时间
	remainingTime := time.Until(time.Unix(claims.Exp, 0))
	if remainingTime <= 0 {
		return nil // Token已过期，无需加入黑名单
	}

	// 存入Redis，设置过期时间为Token剩余时间
	key := s.getBlacklistKey(token)
	if err := s.redisClient.Set(ctx, key, "revoked", remainingTime); err != nil {
		return fmt.Errorf("存储黑名单记录失败: %w", err)
	}

	return nil
}

// IsTokenRevoked 检查Token是否已被吊销
func (s *JWTServiceImpl) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	if s.redisClient == nil {
		return false, nil // 没有Redis，认为未吊销
	}

	// 移除Bearer前缀
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
	// 使用Token的哈希值作为key（避免key过长）
	h := sha256.New()
	h.Write([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("token:blacklist:%s", tokenHash[:32]) // 取前32位
}

// ============ 辅助方法 ============

// ParseTokenClaims 解析Token的Claims（不验证签名）- 内部方法
func (s *JWTServiceImpl) ParseTokenClaims(token string) (*TokenClaims, error) {
	// 移除Bearer前缀
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("token格式错误")
	}

	payloadBase64 := parts[1]
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return nil, fmt.Errorf("解码payload失败: %w", err)
	}

	var claims TokenClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("解析claims失败: %w", err)
	}

	return &claims, nil
}
