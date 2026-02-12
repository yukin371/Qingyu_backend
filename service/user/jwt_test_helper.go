package user

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTTestHelper JWT测试辅助工具
type JWTTestHelper struct {
	secret string
	issuer string
}

// NewJWTTestHelper 创建JWT测试辅助工具
func NewJWTTestHelper(secret string) *JWTTestHelper {
	if secret == "" {
		secret = "test-secret-key-for-integration-testing-only"
	}
	return &JWTTestHelper{
		secret: secret,
		issuer: "qingyu-test",
	}
}

// GenerateTestToken 生成测试Token
func (h *JWTTestHelper) GenerateTestToken(userID string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": "access",
		"exp":        time.Now().Add(expiration).Unix(),
		"iat":        time.Now().Unix(),
		"iss":        h.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}

// GenerateTestTokenWithRoles 生成带角色的测试Token
func (h *JWTTestHelper) GenerateTestTokenWithRoles(userID string, roles []string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"roles":      roles,
		"token_type": "access",
		"exp":        time.Now().Add(expiration).Unix(),
		"iat":        time.Now().Unix(),
		"iss":        h.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}

// ValidateTestToken 验证测试Token
func (h *JWTTestHelper) ValidateTestToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
	}

	return "", jwt.ErrSignatureInvalid
}

// GenerateExpiredToken 生成过期的测试Token
func (h *JWTTestHelper) GenerateExpiredToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"token_type": "access",
		"exp":        time.Now().Add(-time.Hour).Unix(), // 1小时前过期
		"iat":        time.Now().Add(-2 * time.Hour).Unix(),
		"iss":        h.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}

// GenerateInvalidToken 生成无效签名Token
func (h *JWTTestHelper) GenerateInvalidToken(userID string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用错误的密钥签名
	tokenString, _ := token.SignedString([]byte("wrong-secret"))
	return tokenString
}

// ParseTokenWithoutValidation 解析Token但不验证签名
func (h *JWTTestHelper) ParseTokenWithoutValidation(tokenString string) (*jwt.MapClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
