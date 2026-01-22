package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/config"

	"github.com/golang-jwt/jwt/v4"
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
	// 使用秒级时间戳以匹配标准JWT格式（middleware/jwt.go期望秒级）
	claims := &TokenClaims{
		UserID: userID,
		Roles:  roles,
		Exp:    now.Add(s.config.Expiration).Unix(), // 使用秒级时间戳（标准JWT格式）
		Iat:    now.Unix(),                          // 使用秒级时间戳
	}

	return s.generateJWT(claims)
}

// GenerateTokenWithUsername 生成带用户名的访问Token（扩展方法）
func (s *JWTServiceImpl) GenerateTokenWithUsername(ctx context.Context, userID string, username string, roles []string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id不能为空")
	}

	now := time.Now()
	// 使用秒级时间戳以匹配标准JWT格式（middleware/jwt.go期望秒级）
	claims := &TokenClaims{
		UserID: userID,
		Roles:  roles,
		Exp:    now.Add(s.config.Expiration).Unix(), // 使用秒级时间戳（标准JWT格式）
		Iat:    now.Unix(),                          // 使用秒级时间戳
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
		Exp:    now.Add(s.config.RefreshDuration).Unix(), // 使用秒级时间戳
		Iat:    now.Unix(),                               // 使用秒级时间戳
	}

	refreshToken, err = s.generateJWT(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("生成刷新Token失败: %w", err)
	}

	return accessToken, refreshToken, nil
}

// generateJWT 内部方法：使用标准JWT库生成token
func (s *JWTServiceImpl) generateJWT(tokenClaims *TokenClaims) (string, error) {
	// 创建标准JWT Claims结构
	claims := jwt.MapClaims{
		"user_id": tokenClaims.UserID,
		"roles":   tokenClaims.Roles,
		"exp":     float64(tokenClaims.Exp), // 标准JWT使用float64
		"iat":     float64(tokenClaims.Iat),
		"iss":     "Qingyu",
		"sub":     tokenClaims.UserID,
	}

	// 使用标准JWT库生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// ============ Token验证 ============

// ValidateToken 验证Token并返回Claims（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	// 1. 检查Token格式
	if tokenString == "" {
		return nil, errors.New("token为空")
	}

	// 移除Bearer前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	// 2. 使用标准JWT库解析并验证token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token解析失败: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token无效")
	}

	// 3. 提取claims并转换为TokenClaims结构
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, _ := claims["user_id"].(string)
		exp, _ := claims["exp"].(float64)
		iat, _ := claims["iat"].(float64)

		// 提取roles
		var roles []string
		if rolesRaw, ok := claims["roles"].([]interface{}); ok {
			roles = make([]string, len(rolesRaw))
			for i, r := range rolesRaw {
				roles[i] = r.(string)
			}
		}

		tokenClaims := &TokenClaims{
			UserID: userID,
			Roles:  roles,
			Exp:    int64(exp),
			Iat:    int64(iat),
		}

		// 4. 检查Token是否在黑名单中
		if s.redisClient != nil {
			revoked, err := s.IsTokenRevoked(ctx, tokenString)
			if err == nil && revoked {
				return nil, errors.New("token已被吊销")
			}
		}

		return tokenClaims, nil
	}

	return nil, errors.New("无法解析token claims")
}

// ============ Token刷新 ============

// RefreshToken 使用刷新Token获取新的Token（匹配interfaces.go中的定义）
func (s *JWTServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// 1. 验证刷新Token
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("刷新Token验证失败: %w", err)
	}

	// 2. 生成新的访问Token，确保与旧token不同
	var newToken string
	maxAttempts := 10 // 最多尝试10次
	for i := 0; i < maxAttempts; i++ {
		now := time.Now()
		newClaims := &TokenClaims{
			UserID: claims.UserID,
			Roles:  claims.Roles,
			Exp:    now.Add(s.config.Expiration).UnixNano(), // 使用纳秒精度支持短过期时间
			Iat:    now.UnixNano(),                          // 纳秒精度确保唯一性
		}
		newToken, err = s.generateJWT(newClaims)
		if err != nil {
			return "", fmt.Errorf("生成新Token失败: %w", err)
		}

		// 确保新token与旧token不同
		if newToken != refreshToken {
			break
		}

		// 如果相同，等待一小段时间再试
		time.Sleep(time.Nanosecond)
	}

	// 3. 不将refreshToken加入黑名单，允许正常刷新
	// 只有明确调用RevokeToken时才会将token加入黑名单

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
		// 无效的token不需要吊销，直接返回成功（幂等性）
		return nil
	}

	// 计算剩余时间（支持秒级和纳秒级Exp）
	var remainingTime time.Duration
	if claims.Exp > 1e15 {
		// 纳秒级时间戳
		remainingTime = time.Until(time.Unix(0, claims.Exp))
	} else {
		// 秒级时间戳
		remainingTime = time.Until(time.Unix(claims.Exp, 0))
	}

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
