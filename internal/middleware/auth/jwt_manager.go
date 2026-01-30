package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// DefaultAccessExpiration 默认Access Token过期时间（2小时）
	DefaultAccessExpiration = 2 * time.Hour
	// DefaultRefreshExpiration 默认Refresh Token过期时间（7天）
	DefaultRefreshExpiration = 7 * 24 * time.Hour
	// DefaultIssuer 默认签发者
	DefaultIssuer = "qingyu"
)

// Claims JWT声明结构
type Claims struct {
	UserID    string   `json:"user_id"`
	Username  string   `json:"username,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	TokenType string   `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTManager JWT管理器接口
type JWTManager interface {
	// GenerateTokenPair 生成Access Token和Refresh Token
	GenerateTokenPair(userID string, extraClaims map[string]interface{}) (accessToken, refreshToken string, err error)

	// ValidateToken 验证Token并返回Claims
	ValidateToken(tokenString string) (*Claims, error)

	// ParseToken 解析Token（不验证签名）
	ParseToken(tokenString string) (*Claims, error)

	// RefreshToken 使用Refresh Token获取新的Access Token
	RefreshToken(refreshToken string) (string, error)
}

// jwtManager JWT管理器实现
type jwtManager struct {
	secret           []byte
	issuer           string
	accessExpiration time.Duration
	refreshExpiration time.Duration
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret string, accessExpiration, refreshExpiration time.Duration) (JWTManager, error) {
	if secret == "" {
		return nil, errors.New("secret cannot be empty")
	}

	if accessExpiration <= 0 {
		accessExpiration = DefaultAccessExpiration
	}

	if refreshExpiration <= 0 {
		refreshExpiration = DefaultRefreshExpiration
	}

	return &jwtManager{
		secret:           []byte(secret),
		issuer:           DefaultIssuer,
		accessExpiration: accessExpiration,
		refreshExpiration: refreshExpiration,
	}, nil
}

// GenerateTokenPair 生成Access Token和Refresh Token
func (m *jwtManager) GenerateTokenPair(userID string, extraClaims map[string]interface{}) (string, string, error) {
	now := time.Now()

	// 生成Access Token
	accessClaims := &Claims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// 添加额外的claims
	if username, ok := extraClaims["username"].(string); ok {
		accessClaims.Username = username
	}
	if roles, ok := extraClaims["roles"].([]string); ok {
		accessClaims.Roles = roles
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(m.secret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// 生成Refresh Token
	refreshClaims := &Claims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(m.secret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateToken 验证Token并返回Claims
func (m *jwtManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// ParseToken 解析Token（不验证签名）
func (m *jwtManager) ParseToken(tokenString string) (*Claims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// RefreshToken 使用Refresh Token获取新的Access Token
func (m *jwtManager) RefreshToken(refreshToken string) (string, error) {
	// 验证Refresh Token
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// 检查是否为Refresh Token
	if claims.TokenType != "refresh" {
		return "", errors.New("token is not a refresh token")
	}

	// 生成新的Access Token
	now := time.Now()
	newClaims := &Claims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Roles:     claims.Roles,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   claims.UserID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims).SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign new access token: %w", err)
	}

	return newAccessToken, nil
}
