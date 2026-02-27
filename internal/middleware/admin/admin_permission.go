package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authService "Qingyu_backend/service/auth"
)

// AdminPermissionConfig 管理员权限中间件配置
type AdminPermissionConfig struct {
	PermissionService authService.PermissionService
	SuperAdminRole     string // 超级管理员角色名称，默认为 "super_admin"
}

// AdminPermissionOption 配置选项函数
type AdminPermissionOption func(*AdminPermissionConfig)

// WithPermissionService 设置权限服务
func WithPermissionService(service authService.PermissionService) AdminPermissionOption {
	return func(c *AdminPermissionConfig) {
		c.PermissionService = service
	}
}

// WithSuperAdminRole 设置超级管理员角色
func WithSuperAdminRole(role string) AdminPermissionOption {
	return func(c *AdminPermissionConfig) {
		c.SuperAdminRole = role
	}
}

// RequireAdminPermission 要求特定管理权限的中间件
// 用法: router.Use(RequireAdminPermission("user.read"))
func RequireAdminPermission(requiredPermission string, opts ...AdminPermissionOption) gin.HandlerFunc {
	config := &AdminPermissionConfig{
		SuperAdminRole: "super_admin",
	}

	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	return func(c *gin.Context) {
		// 1. 从上下文获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "UNAUTHORIZED",
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "INVALID_USER_ID",
				"message": "无效的用户ID",
			})
			c.Abort()
			return
		}

		// 2. 检查是否是超级管理员
		if config.PermissionService != nil {
			isSuperAdmin, err := config.PermissionService.HasRole(c.Request.Context(), userIDStr, config.SuperAdminRole)
			if err == nil && isSuperAdmin {
				// 超级管理员绕过权限检查
				c.Next()
				return
			}
		}

		// 3. 检查用户是否有指定权限
		if config.PermissionService == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "SERVICE_UNAVAILABLE",
				"message": "权限服务不可用",
			})
			c.Abort()
			return
		}

		hasPermission, err := config.PermissionService.CheckPermission(c.Request.Context(), userIDStr, requiredPermission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "PERMISSION_CHECK_FAILED",
				"message": "权限检查失败",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    "FORBIDDEN",
				"message": "权限不足",
				"details": gin.H{
					"required": requiredPermission,
				},
			})
			c.Abort()
			return
		}

		// 4. 权限检查通过，继续处理请求
		c.Next()
	}
}

// RequireAnyAdminPermission 要求任意一个权限的中间件
// 用法: router.Use(RequireAnyAdminPermission(opts)("user.read", "user.write"))
func RequireAnyAdminPermission(opts ...AdminPermissionOption) func(...string) gin.HandlerFunc {
	config := &AdminPermissionConfig{
		SuperAdminRole: "super_admin",
	}

	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	return func(permissions ...string) gin.HandlerFunc {
		return func(c *gin.Context) {
			// 1. 从上下文获取用户ID
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"code":    "UNAUTHORIZED",
					"message": "用户未认证",
				})
				c.Abort()
				return
			}

			userIDStr, ok := userID.(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"code":    "INVALID_USER_ID",
					"message": "无效的用户ID",
				})
				c.Abort()
				return
			}

			// 2. 检查是否是超级管理员
			if config.PermissionService != nil {
				isSuperAdmin, err := config.PermissionService.HasRole(c.Request.Context(), userIDStr, config.SuperAdminRole)
				if err == nil && isSuperAdmin {
					// 超级管理员绕过权限检查
					c.Next()
					return
				}
			}

			// 3. 检查是否有任意一个权限
			if config.PermissionService == nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"code":    "SERVICE_UNAVAILABLE",
					"message": "权限服务不可用",
				})
				c.Abort()
				return
			}

			for _, permission := range permissions {
				hasPermission, err := config.PermissionService.CheckPermission(c.Request.Context(), userIDStr, permission)
				if err == nil && hasPermission {
					// 有任意一个权限即可
					c.Next()
					return
				}
			}

			// 4. 没有任何权限
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    "FORBIDDEN",
				"message": "权限不足",
				"details": gin.H{
					"required_any": permissions,
				},
			})
			c.Abort()
		}
	}
}

// RequireResourceOwnerPermission 要求资源所有者权限或管理员权限
// 用法: router.Use(RequireResourceOwnerPermission("user", "user_id"))
func RequireResourceOwnerPermission(resourceType, resourceIDParam string, opts ...AdminPermissionOption) gin.HandlerFunc {
	config := &AdminPermissionConfig{
		SuperAdminRole: "super_admin",
	}

	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	return func(c *gin.Context) {
		// 1. 从上下文获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "UNAUTHORIZED",
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "INVALID_USER_ID",
				"message": "无效的用户ID",
			})
			c.Abort()
			return
		}

		// 2. 检查是否是超级管理员
		if config.PermissionService != nil {
			isSuperAdmin, err := config.PermissionService.HasRole(c.Request.Context(), userIDStr, config.SuperAdminRole)
			if err == nil && isSuperAdmin {
				// 超级管理员绕过权限检查
				c.Next()
				return
			}
		}

		// 3. 检查是否有管理权限（可以管理所有资源）
		adminPermission := resourceType + ".manage"
		if config.PermissionService != nil {
			hasAdminPermission, err := config.PermissionService.CheckPermission(c.Request.Context(), userIDStr, adminPermission)
			if err == nil && hasAdminPermission {
				c.Next()
				return
			}
		}

		// 4. 检查是否是资源所有者
		resourceID := c.Param(resourceIDParam)
		if resourceID == "" {
			resourceID = c.Query(resourceIDParam)
		}

		if resourceID == userIDStr {
			// 是资源所有者
			c.Next()
			return
		}

		// 5. 权限不足
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "FORBIDDEN",
			"message": "权限不足",
			"details": gin.H{
				"reason":       "not_resource_owner",
				"resource_type": resourceType,
			},
		})
		c.Abort()
	}
}

// MatchWildcardPermission 通配符权限匹配检查
// 检查用户权限是否匹配通配符模式
func MatchWildcardPermission(userPermissions []string, requiredPermission string) bool {
	for _, permission := range userPermissions {
		// 精确匹配
		if permission == requiredPermission || permission == "*" {
			return true
		}

		// 通配符匹配 (如 user.* 匹配 user.read)
		if strings.HasSuffix(permission, ".*") {
			prefix := strings.TrimSuffix(permission, ".*")
			if strings.HasPrefix(requiredPermission, prefix+".") {
				return true
			}
		}
	}
	return false
}
