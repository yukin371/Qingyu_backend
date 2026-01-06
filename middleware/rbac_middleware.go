package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	sharedService "Qingyu_backend/service/shared"
)

// RBACMiddleware 基于角色的权限控制中间件
type RBACMiddleware struct {
	permissionService sharedService.PermissionService
}

// NewRBACMiddleware 创建RBAC中间件
func NewRBACMiddleware(permissionService sharedService.PermissionService) *RBACMiddleware {
	return &RBACMiddleware{
		permissionService: permissionService,
	}
}

// RequirePermission 要求特定权限
// 用法: router.GET("/admin/users", middleware.RequirePermission("users.read"))
func (m *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID（从JWT认证中间件设置）
		userID, exists := c.Get("user_id")
		if !exists {
			shared.Error(c, http.StatusUnauthorized, "未认证", "请先登录")
			c.Abort()
			return
		}

		// 2. 检查权限
		ctx := c.Request.Context()
		hasPermission, err := m.permissionService.UserHasPermission(ctx, userID.(string), permission)
		if err != nil {
			shared.Error(c, http.StatusInternalServerError, "权限检查失败", err.Error())
			c.Abort()
			return
		}

		if !hasPermission {
			shared.Error(c, http.StatusForbidden, "权限不足", "需要 "+permission+" 权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 要求拥有任意一个权限
// 用法: router.GET("/content", middleware.RequireAnyPermission("content.read", "content.write"))
func (m *RBACMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			shared.Error(c, http.StatusUnauthorized, "未认证", "请先登录")
			c.Abort()
			return
		}

		// 2. 检查是否有任意一个权限
		ctx := c.Request.Context()
		hasAny, err := m.permissionService.UserHasAnyPermission(ctx, userID.(string), permissions)
		if err != nil {
			shared.Error(c, http.StatusInternalServerError, "权限检查失败", err.Error())
			c.Abort()
			return
		}

		if !hasAny {
			shared.Error(c, http.StatusForbidden, "权限不足", "需要以下权限之一: "+strings.Join(permissions, ", "))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions 要求拥有所有权限
// 用法: router.POST("/admin/users/:id/role", middleware.RequireAllPermissions("users.read", "users.manage"))
func (m *RBACMiddleware) RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			shared.Error(c, http.StatusUnauthorized, "未认证", "请先登录")
			c.Abort()
			return
		}

		// 2. 检查是否有所有权限
		ctx := c.Request.Context()
		hasAll, err := m.permissionService.UserHasAllPermissions(ctx, userID.(string), permissions)
		if err != nil {
			shared.Error(c, http.StatusInternalServerError, "权限检查失败", err.Error())
			c.Abort()
			return
		}

		if !hasAll {
			shared.Error(c, http.StatusForbidden, "权限不足", "需要所有权限: "+strings.Join(permissions, ", "))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireResourcePermission 要求资源级权限
// 用法: router.DELETE("/books/:id", middleware.RequireResourcePermission("books", "delete"))
func (m *RBACMiddleware) RequireResourcePermission(resource, action string) gin.HandlerFunc {
	permission := resource + "." + action
	return m.RequirePermission(permission)
}

// RequireRole 要求特定角色
// 注意：这应该与JWT中间件配合使用
// 用法: router.GET("/admin", middleware.RequireRole("admin"))
func (m *RBACMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			shared.Error(c, http.StatusUnauthorized, "未认证", "请先登录")
			c.Abort()
			return
		}

		// 2. 获取用户角色
		ctx := c.Request.Context()
		userRoles, err := m.permissionService.GetUserRoles(ctx, userID.(string))
		if err != nil {
			shared.Error(c, http.StatusInternalServerError, "角色检查失败", err.Error())
			c.Abort()
			return
		}

		// 3. 检查是否拥有所需角色
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			shared.Error(c, http.StatusForbidden, "权限不足", "需要以下角色之一: "+strings.Join(roles, ", "))
			c.Abort()
			return
		}

		c.Next()
	}
}

// ============ 辅助函数 ============

// GetUserPermissions 获取用户的所有权限
func (m *RBACMiddleware) GetUserPermissions(c *gin.Context) ([]string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, nil
	}

	ctx := c.Request.Context()
	permissions, err := m.permissionService.GetUserPermissions(ctx, userID.(string))
	if err != nil {
		return nil, err
	}

	permNames := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permNames = append(permNames, perm.Name)
	}

	return permNames, nil
}

// HasPermission 检查用户是否有权限（辅助函数，可在处理函数中使用）
func (m *RBACMiddleware) HasPermission(c *gin.Context, permission string) (bool, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return false, nil
	}

	ctx := c.Request.Context()
	return m.permissionService.UserHasPermission(ctx, userID.(string), permission)
}
