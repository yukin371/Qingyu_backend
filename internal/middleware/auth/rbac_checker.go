package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// RBACChecker 基于角色的权限检查器
//
// 实现了基于角色（Role-Based Access Control）的权限检查
type RBACChecker struct {
	// rolePerms 角色到权限的映射
	// key: role name
	// value: permission set (key: permission string, value: true)
	rolePerms map[string]map[string]bool

	// userRoles 用户到角色的映射
	// key: user ID
	// value: list of roles
	userRoles map[string][]string

	// mu 读写锁
	mu sync.RWMutex
}

// NewRBACChecker 创建RBAC检查器
func NewRBACChecker(config *CheckerConfig) (Checker, error) {
	checker := &RBACChecker{
		rolePerms: make(map[string]map[string]bool),
		userRoles: make(map[string][]string),
	}

	// 如果提供了配置文件路径，从配置文件加载
	if config != nil && config.ConfigPath != "" {
		// TODO: 实现从配置文件加载
		// 这里暂时返回空检查器
	}

	return checker, nil
}

// Name 返回检查器名称
func (c *RBACChecker) Name() string {
	return "rbac"
}

// Check 检查权限
//
// 检查流程：
// 1. 获取用户的所有角色
// 2. 检查任一角色是否有权限
// 3. 支持通配符权限（如 "*:*" 匹配所有权限）
func (c *RBACChecker) Check(ctx context.Context, subject string, perm Permission) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 1. 获取用户的所有角色
	roles, exists := c.userRoles[subject]
	if !exists || len(roles) == 0 {
		return false, nil
	}

	// 2. 构建权限键
	permKey := perm.String()

	// 3. 检查任一角色是否有权限
	for _, role := range roles {
		perms, ok := c.rolePerms[role]
		if !ok {
			continue
		}

		// 检查通配符权限 "*:*"
		if perms["*:*"] {
			return true, nil
		}

		// 检查资源通配符 "resource:*"
		resourceWildcard := fmt.Sprintf("%s:*", perm.Resource)
		if perms[resourceWildcard] {
			return true, nil
		}

		// 检查精确权限
		if perms[permKey] {
			return true, nil
		}
	}

	return false, nil
}

// BatchCheck 批量检查权限
func (c *RBACChecker) BatchCheck(ctx context.Context, subject string, perms []Permission) ([]bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	results := make([]bool, len(perms))

	// 获取用户的所有角色
	roles, exists := c.userRoles[subject]
	if !exists || len(roles) == 0 {
		// 没有角色，所有权限都是false
		return results, nil
	}

	// 收集用户的所有权限
	userPerms := make(map[string]bool)
	for _, role := range roles {
		if perms, ok := c.rolePerms[role]; ok {
			for perm := range perms {
				userPerms[perm] = true
			}
		}
	}

	// 检查每个权限
	for i, perm := range perms {
		permKey := perm.String()

		// 检查通配符权限
		if userPerms["*:*"] {
			results[i] = true
			continue
		}

		// 检查资源通配符
		resourceWildcard := fmt.Sprintf("%s:*", perm.Resource)
		if userPerms[resourceWildcard] {
			results[i] = true
			continue
		}

		// 检查精确权限
		results[i] = userPerms[permKey]
	}

	return results, nil
}

// Close 关闭检查器
func (c *RBACChecker) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清空数据
	c.rolePerms = make(map[string]map[string]bool)
	c.userRoles = make(map[string][]string)

	return nil
}

// ========== 管理方法 ==========

// AssignRole 为用户分配角色
//
// 如果用户已有该角色，不会重复添加
func (c *RBACChecker) AssignRole(userID string, role string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	roles := c.userRoles[userID]
	for _, r := range roles {
		if r == role {
			return // 已有该角色
		}
	}

	c.userRoles[userID] = append(roles, role)
}

// RevokeRole 撤销用户的角色
func (c *RBACChecker) RevokeRole(userID string, role string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	roles := c.userRoles[userID]
	newRoles := make([]string, 0, len(roles))

	for _, r := range roles {
		if r != role {
			newRoles = append(newRoles, r)
		}
	}

	if len(newRoles) == 0 {
		delete(c.userRoles, userID)
	} else {
		c.userRoles[userID] = newRoles
	}
}

// GetUserRoles 获取用户的所有角色
func (c *RBACChecker) GetUserRoles(userID string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	roles := c.userRoles[userID]
	if roles == nil {
		return []string{}
	}

	// 返回副本
	result := make([]string, len(roles))
	copy(result, roles)
	return result
}

// GrantPermission 为角色授予权限
//
// perm格式: "resource:action" 或 "resource:action:resource_id"
func (c *RBACChecker) GrantPermission(role string, perm string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rolePerms[role] == nil {
		c.rolePerms[role] = make(map[string]bool)
	}

	c.rolePerms[role][perm] = true
}

// RevokePermission 撤销角色的权限
func (c *RBACChecker) RevokePermission(role string, perm string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if perms, ok := c.rolePerms[role]; ok {
		delete(perms, perm)
	}
}

// GetRolePermissions 获取角色的所有权限
func (c *RBACChecker) GetRolePermissions(role string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	perms := c.rolePerms[role]
	if perms == nil {
		return []string{}
	}

	result := make([]string, 0, len(perms))
	for perm := range perms {
		result = append(result, perm)
	}
	return result
}

// ========== 角色层次支持（可选） ==========

// RoleHierarchy 角色层次定义
//
// 例如：
//
//	map[string][]string{
//	    "admin":     {"user", "guest"},  // admin 拥有 user 和 guest 的权限
//	    "moderator": {"user", "guest"},  // moderator 拥有 user 和 guest 的权限
//	    "user":      {"guest"},          // user 拥有 guest 的权限
//	}
type RoleHierarchy map[string][]string

// SetRoleHierarchy 设置角色层次
//
// 注意：这需要在检查权限时考虑角色继承
func (c *RBACChecker) SetRoleHierarchy(hierarchy RoleHierarchy) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 为每个角色添加继承角色的权限
	for role, inheritedRoles := range hierarchy {
		for _, inheritedRole := range inheritedRoles {
			if perms, ok := c.rolePerms[inheritedRole]; ok {
				if c.rolePerms[role] == nil {
					c.rolePerms[role] = make(map[string]bool)
				}
				for perm := range perms {
					c.rolePerms[role][perm] = true
				}
			}
		}
	}
}

// ========== 批量操作方法 ==========

// BatchAssignRoles 批量为用户分配角色
func (c *RBACChecker) BatchAssignRoles(userID string, roles []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	existingRoles := c.userRoles[userID]
	roleSet := make(map[string]bool)

	// 添加现有角色
	for _, role := range existingRoles {
		roleSet[role] = true
	}

	// 添加新角色
	for _, role := range roles {
		if !roleSet[role] {
			existingRoles = append(existingRoles, role)
			roleSet[role] = true
		}
	}

	c.userRoles[userID] = existingRoles
}

// BatchGrantPermissions 批量为角色授予权限
func (c *RBACChecker) BatchGrantPermissions(role string, perms []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rolePerms[role] == nil {
		c.rolePerms[role] = make(map[string]bool)
	}

	for _, perm := range perms {
		c.rolePerms[role][perm] = true
	}
}

// ========== 配置加载方法 ==========

// LoadFromMap 从map加载配置
//
// 用于测试和简单配置
func (c *RBACChecker) LoadFromMap(config map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 加载角色权限
	if rolePerms, ok := config["role_permissions"].(map[string][]string); ok {
		for role, perms := range rolePerms {
			if c.rolePerms[role] == nil {
				c.rolePerms[role] = make(map[string]bool)
			}
			for _, perm := range perms {
				c.rolePerms[role][perm] = true
			}
		}
	}

	// 加载用户角色
	if userRoles, ok := config["user_roles"].(map[string][]string); ok {
		for userID, roles := range userRoles {
			c.userRoles[userID] = roles
		}
	}

	return nil
}

// ========== 辅助方法 ==========

// HasRole 检查用户是否有指定角色
func (c *RBACChecker) HasRole(userID string, role string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	roles := c.userRoles[userID]
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否有任一角色
func (c *RBACChecker) HasAnyRole(userID string, roles []string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userRoles := c.userRoles[userID]
	for _, userRole := range userRoles {
		for _, role := range roles {
			if userRole == role {
				return true
			}
		}
	}
	return false
}

// HasAllRoles 检查用户是否有所有角色
func (c *RBACChecker) HasAllRoles(userID string, roles []string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	userRoles := c.userRoles[userID]
	userRoleSet := make(map[string]bool)
	for _, role := range userRoles {
		userRoleSet[role] = true
	}

	for _, role := range roles {
		if !userRoleSet[role] {
			return false
		}
	}
	return true
}

// Stats 返回检查器统计信息
func (c *RBACChecker) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 计算权限总数
	totalPerms := 0
	for _, perms := range c.rolePerms {
		totalPerms += len(perms)
	}

	return map[string]interface{}{
		"total_roles":     len(c.rolePerms),
		"total_users":     len(c.userRoles),
		"total_permissions": totalPerms,
	}
}

// ========== 初始化注册 ==========

func init() {
	// 注册RBAC检查器
	RegisterChecker("rbac", NewRBACChecker)
}

// ========== 权限字符串辅助函数 ==========

// matchPermission 检查权限是否匹配
//
// 支持通配符：
// - "*:*" 匹配所有权限
// - "resource:*" 匹配资源的所有操作
// - "*:action" 匹配所有资源的该操作
func matchPermission(permKey string, requiredKey string) bool {
	// 精确匹配
	if permKey == requiredKey {
		return true
	}

	permParts := strings.Split(permKey, ":")
	requiredParts := strings.Split(requiredKey, ":")

	// "*:*" 匹配所有
	if permParts[0] == "*" && permParts[1] == "*" {
		return true
	}

	// "resource:*" 匹配
	if len(permParts) == 2 && permParts[1] == "*" {
		if permParts[0] == requiredParts[0] {
			return true
		}
	}

	// "*:action" 匹配
	if len(permParts) == 2 && permParts[0] == "*" {
		if permParts[1] == requiredParts[1] {
			return true
		}
	}

	return false
}
