package middleware

import (
	"context"
	"strings"

	"Qingyu_backend/service/shared/auth"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限中间件
type PermissionMiddleware struct {
	authService auth.AuthService
}

// NewPermissionMiddleware 创建权限中间件
func NewPermissionMiddleware(authService auth.AuthService) *PermissionMiddleware {
	return &PermissionMiddleware{
		authService: authService,
	}
}

// RequirePermission 需要特定权限中间件
func (m *PermissionMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 2. 检查权限
		ctx := context.Background()
		has, err := m.authService.CheckPermission(ctx, userID.(string), permission)
		if err != nil {
			c.JSON(500, gin.H{"error": "权限检查失败: " + err.Error()})
			c.Abort()
			return
		}

		if !has {
			c.JSON(403, gin.H{"error": "权限不足: 需要" + permission + "权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 需要任一权限中间件
func (m *PermissionMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 2. 检查是否有任一权限
		ctx := context.Background()
		for _, perm := range permissions {
			has, err := m.authService.CheckPermission(ctx, userID.(string), perm)
			if err == nil && has {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "权限不足: 需要以下权限之一: " + strings.Join(permissions, ", ")})
		c.Abort()
	}
}

// RequireAllPermissions 需要所有权限中间件
func (m *PermissionMiddleware) RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 2. 检查是否有所有权限
		ctx := context.Background()
		for _, perm := range permissions {
			has, err := m.authService.CheckPermission(ctx, userID.(string), perm)
			if err != nil || !has {
				c.JSON(403, gin.H{"error": "权限不足: 需要" + perm + "权限"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// CheckResourcePermission 检查资源权限（动态权限）
func (m *PermissionMiddleware) CheckResourcePermission(resourceType, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 2. 构建权限字符串
		permission := resourceType + "." + action

		// 3. 检查权限
		ctx := context.Background()
		has, err := m.authService.CheckPermission(ctx, userID.(string), permission)
		if err != nil {
			c.JSON(500, gin.H{"error": "权限检查失败: " + err.Error()})
			c.Abort()
			return
		}

		if !has {
			c.JSON(403, gin.H{
				"error":    "权限不足",
				"resource": resourceType,
				"action":   action,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ============ 权限检查函数（非中间件） ============

// CheckPermission 检查权限（辅助函数）
func CheckPermission(c *gin.Context, authService auth.AuthService, permission string) bool {
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}

	ctx := context.Background()
	has, err := authService.CheckPermission(ctx, userID.(string), permission)
	return err == nil && has
}

// GetUserPermissions 获取用户权限列表
func GetUserPermissions(c *gin.Context, authService auth.AuthService) ([]string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, nil
	}

	ctx := context.Background()
	return authService.GetUserPermissions(ctx, userID.(string))
}
