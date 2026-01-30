package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuth 提供简单的JWT认证中间件
// 这是一个便捷函数，使用默认配置创建JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	// 创建logger
	logger, _ := zap.NewDevelopment()

	// 创建默认的JWT管理器和黑名单
	// 注意：这里使用默认配置，实际使用时应该从服务容器获取
	jwtManager, err := NewJWTManager("default-secret-change-in-production", DefaultAccessExpiration, DefaultRefreshExpiration)
	if err != nil {
		logger.Fatal("Failed to create JWT manager", zap.Error(err))
	}

	middleware := NewJWTAuthMiddleware(jwtManager, nil, logger)
	return middleware.Handler()
}

// GenerateToken 生成JWT令牌
// 这是一个便捷函数，用于生成JWT令牌
// 注意：这个函数使用默认配置，生产环境应该使用服务容器中的JWT管理器
func GenerateToken(userID, username string, roles []string) (string, error) {
	// 创建默认的JWT管理器
	jwtManager, err := NewJWTManager("default-secret-change-in-production", DefaultAccessExpiration, DefaultRefreshExpiration)
	if err != nil {
		return "", err
	}

	// 构建额外的claims
	extraClaims := make(map[string]interface{})
	if username != "" {
		extraClaims["username"] = username
	}
	if len(roles) > 0 {
		extraClaims["roles"] = roles
	}

	// 生成token（这里只返回access token）
	accessToken, _, err := jwtManager.GenerateTokenPair(userID, extraClaims)
	return accessToken, err
}

// RequireRole 要求特定角色的中间件
// 这是一个便捷函数，用于检查用户是否具有指定角色
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(401, gin.H{
				"code":    40101,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(500, gin.H{
				"code":    50001,
				"message": "用户信息格式错误",
			})
			c.Abort()
			return
		}

		// 检查是否有所需角色
		hasRole := false
		for _, r := range userRoles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(403, gin.H{
				"code":    40301,
				"message": "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole 要求任意角色的中间件
// 检查用户是否具有指定角色中的任意一个
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		userRoles, exists := c.Get("roles")
		if !exists {
			c.JSON(401, gin.H{
				"code":    40101,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		userRoleList, ok := userRoles.([]string)
		if !ok {
			c.JSON(500, gin.H{
				"code":    50001,
				"message": "用户信息格式错误",
			})
			c.Abort()
			return
		}

		// 检查是否有任意所需角色
		hasAnyRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoleList {
				if userRole == requiredRole {
					hasAnyRole = true
					break
				}
			}
			if hasAnyRole {
				break
			}
		}

		if !hasAnyRole {
			c.JSON(403, gin.H{
				"code":    40301,
				"message": "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllRoles 要求所有角色的中间件
// 检查用户是否具有所有指定的角色
func RequireAllRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		userRoles, exists := c.Get("roles")
		if !exists {
			c.JSON(401, gin.H{
				"code":    40101,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		userRoleList, ok := userRoles.([]string)
		if !ok {
			c.JSON(500, gin.H{
				"code":    50001,
				"message": "用户信息格式错误",
			})
			c.Abort()
			return
		}

		// 检查是否拥有所有角色
		roleMap := make(map[string]bool)
		for _, r := range userRoleList {
			roleMap[r] = true
		}

		for _, requiredRole := range roles {
			if !roleMap[requiredRole] {
				c.JSON(403, gin.H{
					"code":    40301,
					"message": "权限不足",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
