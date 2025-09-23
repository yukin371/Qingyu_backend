package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// PermissionConfig 权限中间件配置
type PermissionConfig struct {
	RequiredRoles       []string          `json:"required_roles" yaml:"required_roles"`
	RequiredPermissions []string          `json:"required_permissions" yaml:"required_permissions"`
	AllowedRoles        []string          `json:"allowed_roles" yaml:"allowed_roles"`
	AllowedPermissions  []string          `json:"allowed_permissions" yaml:"allowed_permissions"`
	SkipPaths           []string          `json:"skip_paths" yaml:"skip_paths"`
	CheckMode           string            `json:"check_mode" yaml:"check_mode"` // "any", "all"
	Message             string            `json:"message" yaml:"message"`
	StatusCode          int               `json:"status_code" yaml:"status_code"`
	RoleHierarchy       map[string][]string `json:"role_hierarchy" yaml:"role_hierarchy"`
}

// Permission 权限结构
type Permission struct {
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// DefaultPermissionConfig 默认权限配置
func DefaultPermissionConfig() PermissionConfig {
	return PermissionConfig{
		RequiredRoles:       []string{},
		RequiredPermissions: []string{},
		AllowedRoles:        []string{},
		AllowedPermissions:  []string{},
		SkipPaths:           []string{"/health", "/metrics"},
		CheckMode:           "any",
		Message:             "权限不足，无法访问该资源",
		StatusCode:          http.StatusForbidden,
		RoleHierarchy: map[string][]string{
			"admin":     {"user", "guest"},
			"moderator": {"user", "guest"},
			"user":      {"guest"},
		},
	}
}

// RequirePermission 要求特定权限
func RequirePermission(permissions ...string) gin.HandlerFunc {
	config := DefaultPermissionConfig()
	config.RequiredPermissions = permissions
	return PermissionWithConfig(config)
}

// RequireAllRoles 要求所有角色
func RequireAllRoles(roles ...string) gin.HandlerFunc {
	config := DefaultPermissionConfig()
	config.RequiredRoles = roles
	config.CheckMode = "all"
	return PermissionWithConfig(config)
}

// PermissionWithConfig 带配置的权限中间件
func PermissionWithConfig(config PermissionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过权限检查
		if shouldSkipPermissionCheck(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// 获取用户信息
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      40101,
				"message":   "用户未认证",
				"timestamp": time.Now().Unix(),
				"data":      nil,
			})
			c.Abort()
			return
		}

		userCtx, ok := user.(*UserContext)
		if !ok {
			// 尝试从JWT claims获取用户信息
			if claims, exists := c.Get("claims"); exists {
				if jwtClaims, ok := claims.(*JWTClaims); ok {
					userCtx = &UserContext{
						UserID:      jwtClaims.UserID,
						Username:    jwtClaims.Username,
						Roles:       jwtClaims.Roles,
						Permissions: []string{}, // 这里可以根据角色查询权限
						LoginTime:   jwtClaims.IssuedAt.Unix(),
					}
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code":      40102,
						"message":   "用户信息格式错误",
						"timestamp": time.Now().Unix(),
						"data":      nil,
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":      40102,
					"message":   "用户信息格式错误",
					"timestamp": time.Now().Unix(),
					"data":      nil,
				})
				c.Abort()
				return
			}
		}

		// 检查权限
		if !checkUserPermission(userCtx, config) {
			c.JSON(config.StatusCode, gin.H{
				"code":      40301,
				"message":   config.Message,
				"timestamp": time.Now().Unix(),
				"data":      nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// shouldSkipPermissionCheck 检查是否应该跳过权限检查
func shouldSkipPermissionCheck(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// checkUserPermission 检查用户权限
func checkUserPermission(user *UserContext, config PermissionConfig) bool {
	// 检查必需角色
	if len(config.RequiredRoles) > 0 {
		if config.CheckMode == "all" {
			if !hasAllRoles(user.Roles, config.RequiredRoles, config.RoleHierarchy) {
				return false
			}
		} else {
			if !hasAnyRole(user.Roles, config.RequiredRoles, config.RoleHierarchy) {
				return false
			}
		}
	}

	// 检查允许的角色
	if len(config.AllowedRoles) > 0 {
		if !hasAnyRole(user.Roles, config.AllowedRoles, config.RoleHierarchy) {
			return false
		}
	}

	// 检查必需权限
	if len(config.RequiredPermissions) > 0 {
		if config.CheckMode == "all" {
			if !hasAllPermissions(user.Permissions, config.RequiredPermissions) {
				return false
			}
		} else {
			if !hasAnyPermission(user.Permissions, config.RequiredPermissions) {
				return false
			}
		}
	}

	// 检查允许的权限
	if len(config.AllowedPermissions) > 0 {
		if !hasAnyPermission(user.Permissions, config.AllowedPermissions) {
			return false
		}
	}

	return true
}

// hasAllRoles 检查是否拥有所有角色
func hasAllRoles(userRoles, requiredRoles []string, hierarchy map[string][]string) bool {
	for _, requiredRole := range requiredRoles {
		if !hasRole(userRoles, requiredRole, hierarchy) {
			return false
		}
	}
	return true
}

// hasAnyRole 检查是否拥有任意角色
func hasAnyRole(userRoles, allowedRoles []string, hierarchy map[string][]string) bool {
	for _, allowedRole := range allowedRoles {
		if hasRole(userRoles, allowedRole, hierarchy) {
			return true
		}
	}
	return false
}

// hasRole 检查是否拥有特定角色（考虑角色层次）
func hasRole(userRoles []string, targetRole string, hierarchy map[string][]string) bool {
	// 直接匹配
	for _, userRole := range userRoles {
		if userRole == targetRole {
			return true
		}
		
		// 检查角色层次
		if inheritedRoles, exists := hierarchy[userRole]; exists {
			for _, inheritedRole := range inheritedRoles {
				if inheritedRole == targetRole {
					return true
				}
			}
		}
	}
	return false
}

// hasAllPermissions 检查是否拥有所有权限
func hasAllPermissions(userPermissions, requiredPermissions []string) bool {
	for _, requiredPermission := range requiredPermissions {
		if !hasPermission(userPermissions, requiredPermission) {
			return false
		}
	}
	return true
}

// hasAnyPermission 检查是否拥有任意权限
func hasAnyPermission(userPermissions, allowedPermissions []string) bool {
	for _, allowedPermission := range allowedPermissions {
		if hasPermission(userPermissions, allowedPermission) {
			return true
		}
	}
	return false
}

// hasPermission 检查是否拥有特定权限
func hasPermission(userPermissions []string, targetPermission string) bool {
	for _, userPermission := range userPermissions {
		if userPermission == targetPermission {
			return true
		}
		
		// 支持通配符权限，如 "user:*" 匹配 "user:read", "user:write"
		if strings.HasSuffix(userPermission, ":*") {
			prefix := strings.TrimSuffix(userPermission, ":*")
			if strings.HasPrefix(targetPermission, prefix+":") {
				return true
			}
		}
	}
	return false
}

// CreatePermissionMiddleware 创建权限中间件（用于中间件工厂）
func CreatePermissionMiddleware(config map[string]interface{}) (gin.HandlerFunc, error) {
	permissionConfig := DefaultPermissionConfig()
	
	// 解析配置
	if requiredRoles, ok := config["required_roles"].([]string); ok {
		permissionConfig.RequiredRoles = requiredRoles
	}
	if requiredPermissions, ok := config["required_permissions"].([]string); ok {
		permissionConfig.RequiredPermissions = requiredPermissions
	}
	if allowedRoles, ok := config["allowed_roles"].([]string); ok {
		permissionConfig.AllowedRoles = allowedRoles
	}
	if allowedPermissions, ok := config["allowed_permissions"].([]string); ok {
		permissionConfig.AllowedPermissions = allowedPermissions
	}
	if skipPaths, ok := config["skip_paths"].([]string); ok {
		permissionConfig.SkipPaths = skipPaths
	}
	if checkMode, ok := config["check_mode"].(string); ok {
		permissionConfig.CheckMode = checkMode
	}
	if message, ok := config["message"].(string); ok {
		permissionConfig.Message = message
	}
	if statusCode, ok := config["status_code"].(int); ok {
		permissionConfig.StatusCode = statusCode
	}
	if roleHierarchy, ok := config["role_hierarchy"].(map[string][]string); ok {
		permissionConfig.RoleHierarchy = roleHierarchy
	}
	
	return PermissionWithConfig(permissionConfig), nil
}