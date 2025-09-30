package middleware

import (
	"context"
	"strings"

	"Qingyu_backend/service/shared/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware Auth认证中间件
type AuthMiddleware struct {
	authService auth.AuthService
}

// NewAuthMiddleware 创建Auth中间件
func NewAuthMiddleware(authService auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth 需要认证中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "未提供认证Token"})
			c.Abort()
			return
		}

		// 2. 解析Bearer Token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(401, gin.H{"error": "Token格式错误"})
			c.Abort()
			return
		}

		// 3. 验证Token
		ctx := context.Background()
		claims, err := m.authService.ValidateToken(ctx, token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Token验证失败: " + err.Error()})
			c.Abort()
			return
		}

		// 4. 将用户信息存入Context
		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Set("token", token)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（不强制要求认证）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 尝试获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 2. 解析Token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.Next()
			return
		}

		// 3. 验证Token
		ctx := context.Background()
		claims, err := m.authService.ValidateToken(ctx, token)
		if err != nil {
			// Token无效，但不阻止请求
			c.Next()
			return
		}

		// 4. 将用户信息存入Context
		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Set("token", token)

		c.Next()
	}
}

// RequireRole 需要特定角色中间件
func (m *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 2. 检查角色
		ctx := context.Background()
		has, err := m.authService.HasRole(ctx, userID.(string), requiredRole)
		if err != nil || !has {
			c.JSON(403, gin.H{"error": "权限不足: 需要" + requiredRole + "角色"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole 需要任一角色中间件
func (m *AuthMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 2. 检查是否有任一角色
		ctx := context.Background()
		for _, role := range roles {
			has, err := m.authService.HasRole(ctx, userID.(string), role)
			if err == nil && has {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "权限不足: 需要以下角色之一: " + strings.Join(roles, ", ")})
		c.Abort()
	}
}

// ============ 辅助函数 ============

// GetUserID 从Context获取用户ID
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserRoles 从Context获取用户角色
func GetUserRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get("roles")
	if !exists {
		return nil, false
	}
	return roles.([]string), true
}

// MustGetUserID 从Context获取用户ID（必须存在）
func MustGetUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		panic("user_id not found in context")
	}
	return userID.(string)
}
