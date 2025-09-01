package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 定义JWT密钥和过期时间
var (
	JwtSecret = []byte("qingyu_secret_key") // 在生产环境中应该使用环境变量存储
	TokenExpireDuration = 24 * time.Hour    // Token有效期为24小时
)

// Claims 自定义JWT的声明
type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID, username string) (string, error) {
	now := time.Now()
	expireTime := now.Add(TokenExpireDuration)

	claims := Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "qingyu_backend",
			Subject:   "user_token",
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString(JwtSecret)
	return tokenString, err
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}