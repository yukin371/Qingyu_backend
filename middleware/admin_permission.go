package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminPermissionMiddleware 管理员权限验证中间件
func AdminPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Context中获取用户角色数组（由JWTAuth中间件设置）
		roles, exists := c.Get("userRoles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：无法获取用户角色",
			})
			c.Abort()
			return
		}

		// 转换为字符串切片
		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：角色格式错误",
			})
			c.Abort()
			return
		}

		// 检查是否为管理员
		hasAdminRole := false
		for _, role := range userRoles {
			if role == "admin" || role == "super_admin" {
				hasAdminRole = true
				break
			}
		}

		if !hasAdminRole {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SuperAdminPermissionMiddleware 超级管理员权限验证中间件
func SuperAdminPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Context中获取用户角色数组（由JWTAuth中间件设置）
		roles, exists := c.Get("userRoles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：无法获取用户角色",
			})
			c.Abort()
			return
		}

		// 转换为字符串切片
		userRoles, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：角色格式错误",
			})
			c.Abort()
			return
		}

		// 只允许超级管理员
		hasSuperAdminRole := false
		for _, role := range userRoles {
			if role == "super_admin" {
				hasSuperAdminRole = true
				break
			}
		}

		if !hasSuperAdminRole {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：需要超级管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermissionsMiddleware 通用权限验证中间件
func RequirePermissionsMiddleware(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Context中获取用户权限列表
		permissions, exists := c.Get("user_permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：无法获取用户权限",
			})
			c.Abort()
			return
		}

		// 转换为字符串切片
		userPerms, ok := permissions.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足：权限格式错误",
			})
			c.Abort()
			return
		}

		// 检查是否拥有所需权限
		for _, required := range requiredPermissions {
			if !checkAdminPermission(userPerms, required) {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "权限不足：缺少权限 " + required,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// checkAdminPermission 检查用户是否拥有指定权限
func checkAdminPermission(userPermissions []string, required string) bool {
	for _, perm := range userPermissions {
		if perm == required || perm == "*" { // * 表示所有权限
			return true
		}
	}
	return false
}
