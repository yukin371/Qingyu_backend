package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"Qingyu_backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims 定义JWT的声明结构
type JWTClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// UserContext 用户上下文信息
type UserContext struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	LoginTime   int64    `json:"login_time"`
}

// AuthConfig JWT认证中间件配置
type AuthConfig struct {
	SkipPaths    []string `json:"skip_paths"`
	RequiredRole string   `json:"required_role"`
	AllowedRoles []string `json:"allowed_roles"`
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return JWTAuthWithConfig(AuthConfig{})
}

// JWTAuthWithConfig 带配置的JWT认证中间件
func JWTAuthWithConfig(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过认证
		if shouldSkipAuth(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40101,
				"message": "未提供认证令牌",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 解析token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40102,
				"message": "无效的认证令牌格式",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 验证token
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40103,
				"message": "无效的认证令牌",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 检查权限
		if !checkPermission(claims, config) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    40301,
				"message": "权限不足",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		userContext := UserContext{
			UserID:   claims.UserID,
			Username: claims.Username,
			Roles:    claims.Roles,
		}
		c.Set("user", userContext)
		// 同时设置两种命名方式以保持向后兼容
		// TODO: 移除驼峰命名方式
		c.Set("userId", claims.UserID)  // 驼峰命名（reader, ai, writer模块使用）
		c.Set("user_id", claims.UserID) // 下划线命名（user, admin, shared模块使用）
		c.Set("username", claims.Username)
		c.Set("userRoles", claims.Roles)

		c.Next()
	}
}

// RequireRole 要求特定角色的中间件
func RequireRole(role string) gin.HandlerFunc {
	return JWTAuthWithConfig(AuthConfig{
		RequiredRole: role,
	})
}

// RequireAnyRole 要求任意角色的中间件
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return JWTAuthWithConfig(AuthConfig{
		AllowedRoles: roles,
	})
}

// shouldSkipAuth 检查是否应该跳过认证
func shouldSkipAuth(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// checkPermission 检查用户权限
func checkPermission(claims *JWTClaims, config AuthConfig) bool {
	// 如果没有配置权限要求，则通过
	if config.RequiredRole == "" && len(config.AllowedRoles) == 0 {
		return true
	}

	// 检查必需角色
	if config.RequiredRole != "" {
		for _, role := range claims.Roles {
			if role == config.RequiredRole {
				return true
			}
		}
		return false
	}

	// 检查允许的角色
	if len(config.AllowedRoles) > 0 {
		for _, userRole := range claims.Roles {
			for _, allowedRole := range config.AllowedRoles {
				if userRole == allowedRole {
					return true
				}
			}
		}
		return false
	}

	return true
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID, username string, roles []string) (string, error) {
	// 获取JWT配置
	cfg := config.GlobalConfig.JWT
	if cfg == nil {
		return "", errors.New("JWT configuration is missing")
	}

	// 创建声明
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Qingyu",
			Subject:   userID,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString([]byte(cfg.Secret))
}

// GenerateTokenCompat 兼容旧版本的令牌生成函数
func GenerateTokenCompat(userID, role string) (string, error) {
	return GenerateToken(userID, "", []string{role})
}

// RefreshToken 刷新JWT令牌
func RefreshToken(tokenString string) (string, error) {
	// 解析旧令牌
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// 检查令牌是否即将过期（在过期前30分钟内可以刷新）
	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return "", errors.New("token is not eligible for refresh")
	}

	// 生成新令牌
	return GenerateToken(claims.UserID, claims.Username, claims.Roles)
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 获取JWT配置
	cfg := config.GlobalConfig.JWT
	if cfg == nil {
		return nil, errors.New("JWT configuration is missing")
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// 验证令牌类型和有效性
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// 额外的验证
		if claims.UserID == "" {
			return nil, errors.New("invalid token: missing user ID")
		}

		// 检查令牌是否过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// GetUserFromContext 从Gin上下文中获取用户信息
func GetUserFromContext(c *gin.Context) (*UserContext, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	userContext, ok := user.(UserContext)
	if !ok {
		return nil, errors.New("invalid user context type")
	}

	return &userContext, nil
}

// HasRole 检查用户是否具有指定角色
func HasRole(c *gin.Context, role string) bool {
	user, err := GetUserFromContext(c)
	if err != nil {
		return false
	}

	for _, userRole := range user.Roles {
		if userRole == role {
			return true
		}
	}

	return false
}

// HasAnyRole 检查用户是否具有任意指定角色
func HasAnyRole(c *gin.Context, roles ...string) bool {
	user, err := GetUserFromContext(c)
	if err != nil {
		return false
	}

	for _, userRole := range user.Roles {
		for _, role := range roles {
			if userRole == role {
				return true
			}
		}
	}

	return false
}
