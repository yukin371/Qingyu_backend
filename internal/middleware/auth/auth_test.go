package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ========== ParsePermission 测试 ==========

func TestParsePermission(t *testing.T) {
	t.Run("ValidPermission_ResourceAction", func(t *testing.T) {
		perm, err := ParsePermission("project:read")
		assert.NoError(t, err)
		assert.Equal(t, "project", perm.Resource)
		assert.Equal(t, "read", perm.Action)
		assert.Equal(t, "", perm.ResourceID)
	})

	t.Run("ValidPermission_ResourceActionID", func(t *testing.T) {
		perm, err := ParsePermission("project:update:123")
		assert.NoError(t, err)
		assert.Equal(t, "project", perm.Resource)
		assert.Equal(t, "update", perm.Action)
		assert.Equal(t, "123", perm.ResourceID)
	})

	t.Run("ValidPermission_Wildcard", func(t *testing.T) {
		perm, err := ParsePermission("*:*")
		assert.NoError(t, err)
		assert.Equal(t, "*", perm.Resource)
		assert.Equal(t, "*", perm.Action)
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		_, err := ParsePermission("invalid")
		assert.Error(t, err)
	})
}

// ========== Permission.String 测试 ==========

func TestPermission_String(t *testing.T) {
	t.Run("WithoutResourceID", func(t *testing.T) {
		perm := Permission{
			Resource: "project",
			Action:   "read",
		}
		assert.Equal(t, "project:read", perm.String())
	})

	t.Run("WithResourceID", func(t *testing.T) {
		perm := Permission{
			Resource:   "document",
			Action:     "delete",
			ResourceID: "456",
		}
		assert.Equal(t, "document:delete:456", perm.String())
	})
}

// ========== NoOpChecker 测试 ==========

func TestNoOpChecker(t *testing.T) {
	checker := &NoOpChecker{}

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, "noop", checker.Name())
	})

	t.Run("Check", func(t *testing.T) {
		allowed, err := checker.Check(context.Background(), "user1", Permission{
			Resource: "project",
			Action:   "read",
		})
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("BatchCheck", func(t *testing.T) {
		perms := []Permission{
			{Resource: "project", Action: "read"},
			{Resource: "document", Action: "write"},
		}
		results, err := checker.BatchCheck(context.Background(), "user1", perms)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.True(t, results[0])
		assert.True(t, results[1])
	})

	t.Run("Close", func(t *testing.T) {
		err := checker.Close()
		assert.NoError(t, err)
	})
}

// ========== RBACChecker 测试 ==========

// newTestRBACChecker 创建测试用的RBAC检查器
func newTestRBACChecker() *RBACChecker {
	checker, _ := NewRBACChecker(nil)
	return checker.(*RBACChecker)
}

func TestRBACChecker_Name(t *testing.T) {
	checker := newTestRBACChecker()
	assert.Equal(t, "rbac", checker.Name())
}

func TestRBACChecker_AssignRole(t *testing.T) {
	checker := newTestRBACChecker()

	t.Run("AssignNewRole", func(t *testing.T) {
		checker.AssignRole("user1", "admin")
		roles := checker.GetUserRoles("user1")
		assert.Contains(t, roles, "admin")
	})

	t.Run("AssignDuplicateRole", func(t *testing.T) {
		checker.AssignRole("user2", "editor")
		checker.AssignRole("user2", "editor") // 重复分配
		roles := checker.GetUserRoles("user2")
		assert.Len(t, roles, 1) // 只有一个
	})
}

func TestRBACChecker_RevokeRole(t *testing.T) {
	checker := newTestRBACChecker()
	checker.AssignRole("user1", "admin")
	checker.AssignRole("user1", "editor")

	t.Run("RevokeExistingRole", func(t *testing.T) {
		checker.RevokeRole("user1", "admin")
		roles := checker.GetUserRoles("user1")
		assert.NotContains(t, roles, "admin")
		assert.Contains(t, roles, "editor")
	})

	t.Run("RevokeAllRoles", func(t *testing.T) {
		checker.RevokeRole("user1", "editor")
		roles := checker.GetUserRoles("user1")
		assert.Empty(t, roles)
	})
}

func TestRBACChecker_GrantPermission(t *testing.T) {
	checker := newTestRBACChecker()

	t.Run("GrantSinglePermission", func(t *testing.T) {
		checker.GrantPermission("admin", "project:read")
		perms := checker.GetRolePermissions("admin")
		assert.Contains(t, perms, "project:read")
	})

	t.Run("GrantMultiplePermissions", func(t *testing.T) {
		checker.BatchGrantPermissions("editor", []string{
			"document:read",
			"document:write",
		})
		perms := checker.GetRolePermissions("editor")
		assert.Contains(t, perms, "document:read")
		assert.Contains(t, perms, "document:write")
	})
}

func TestRBACChecker_RevokePermission(t *testing.T) {
	checker := newTestRBACChecker()
	checker.GrantPermission("admin", "project:read")
	checker.GrantPermission("admin", "project:write")

	checker.RevokePermission("admin", "project:read")
	perms := checker.GetRolePermissions("admin")
	assert.NotContains(t, perms, "project:read")
	assert.Contains(t, perms, "project:write")
}

func TestRBACChecker_Check(t *testing.T) {
	checker := newTestRBACChecker()

	// 设置初始数据
	checker.GrantPermission("admin", "*:*")
	checker.GrantPermission("editor", "project:read")
	checker.GrantPermission("editor", "project:*") // 资源通配符
	checker.AssignRole("user1", "admin")
	checker.AssignRole("user2", "editor")
	checker.AssignRole("user3", "viewer") // 没有权限

	t.Run("Check_WithWildcardPermission", func(t *testing.T) {
		allowed, err := checker.Check(context.Background(), "user1", Permission{
			Resource: "any",
			Action:   "any",
		})
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Check_WithResourceWildcard", func(t *testing.T) {
		allowed, err := checker.Check(context.Background(), "user2", Permission{
			Resource: "project",
			Action:   "write",
		})
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Check_WithExactPermission", func(t *testing.T) {
		allowed, err := checker.Check(context.Background(), "user2", Permission{
			Resource: "project",
			Action:   "read",
		})
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Check_NoPermission", func(t *testing.T) {
		allowed, err := checker.Check(context.Background(), "user3", Permission{
			Resource: "project",
			Action:   "read",
		})
		assert.NoError(t, err)
		assert.False(t, allowed)
	})

	t.Run("Check_NoRole", func(t *testing.T) {
		allowed, err := checker.Check(context.Background(), "user999", Permission{
			Resource: "project",
			Action:   "read",
		})
		assert.NoError(t, err)
		assert.False(t, allowed)
	})
}

func TestRBACChecker_BatchCheck(t *testing.T) {
	checker := newTestRBACChecker()

	// 设置初始数据
	checker.GrantPermission("editor", "project:read")
	checker.GrantPermission("editor", "document:write")
	checker.AssignRole("user1", "editor")

	perms := []Permission{
		{Resource: "project", Action: "read"},
		{Resource: "document", Action: "write"},
		{Resource: "book", Action: "delete"}, // 没有权限
	}

	results, err := checker.BatchCheck(context.Background(), "user1", perms)
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	assert.True(t, results[0])
	assert.True(t, results[1])
	assert.False(t, results[2])
}

func TestRBACChecker_HasRole(t *testing.T) {
	checker := newTestRBACChecker()
	checker.AssignRole("user1", "admin")
	checker.AssignRole("user1", "editor")

	t.Run("HasRole_Existing", func(t *testing.T) {
		assert.True(t, checker.HasRole("user1", "admin"))
	})

	t.Run("HasRole_NotExisting", func(t *testing.T) {
		assert.False(t, checker.HasRole("user1", "viewer"))
	})
}

func TestRBACChecker_HasAnyRole(t *testing.T) {
	checker := newTestRBACChecker()
	checker.AssignRole("user1", "admin")
	checker.AssignRole("user1", "editor")

	t.Run("HasAnyRole_Match", func(t *testing.T) {
		assert.True(t, checker.HasAnyRole("user1", []string{"admin", "viewer"}))
	})

	t.Run("HasAnyRole_NoMatch", func(t *testing.T) {
		assert.False(t, checker.HasAnyRole("user1", []string{"viewer", "guest"}))
	})
}

func TestRBACChecker_HasAllRoles(t *testing.T) {
	checker := newTestRBACChecker()
	checker.AssignRole("user1", "admin")
	checker.AssignRole("user1", "editor")

	t.Run("HasAllRoles_HasAll", func(t *testing.T) {
		assert.True(t, checker.HasAllRoles("user1", []string{"admin", "editor"}))
	})

	t.Run("HasAllRoles_MissingOne", func(t *testing.T) {
		assert.False(t, checker.HasAllRoles("user1", []string{"admin", "viewer"}))
	})
}

func TestRBACChecker_LoadFromMap(t *testing.T) {
	checker := newTestRBACChecker()

	config := map[string]interface{}{
		"role_permissions": map[string][]string{
			"admin": {"*:*", "project:read"},
			"editor": {"document:write"},
		},
		"user_roles": map[string][]string{
			"user1": {"admin"},
			"user2": {"editor"},
		},
	}

	err := checker.LoadFromMap(config)
	assert.NoError(t, err)

	// 验证角色权限
	adminPerms := checker.GetRolePermissions("admin")
	assert.Contains(t, adminPerms, "*:*")

	// 验证用户角色
	user1Roles := checker.GetUserRoles("user1")
	assert.Contains(t, user1Roles, "admin")
}

func TestRBACChecker_Stats(t *testing.T) {
	checker := newTestRBACChecker()

	checker.GrantPermission("admin", "*:*")
	checker.GrantPermission("editor", "project:read")
	checker.AssignRole("user1", "admin")
	checker.AssignRole("user2", "editor")

	stats := checker.Stats()
	assert.Equal(t, 2, stats["total_roles"])
	assert.Equal(t, 2, stats["total_users"])
	assert.Equal(t, 2, stats["total_permissions"])
}

func TestRBACChecker_Close(t *testing.T) {
	checker := newTestRBACChecker()
	checker.AssignRole("user1", "admin")
	checker.GrantPermission("admin", "*:*")

	err := checker.Close()
	assert.NoError(t, err)

	// 验证数据已清空
	stats := checker.Stats()
	assert.Equal(t, 0, stats["total_roles"])
	assert.Equal(t, 0, stats["total_users"])
}

// ========== PermissionMiddleware 测试 ==========

func TestPermissionConfig_Default(t *testing.T) {
	config := DefaultPermissionConfig()
	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.Equal(t, "rbac", config.Strategy)
	assert.Equal(t, 403, config.StatusCode)
	assert.Contains(t, config.SkipPaths, "/health")
}

func TestNewPermissionMiddleware(t *testing.T) {
	t.Run("WithValidConfig", func(t *testing.T) {
		config := DefaultPermissionConfig()
		config.Strategy = "noop" // 使用空检查器

		middleware, err := NewPermissionMiddleware(config, nil)
		assert.NoError(t, err)
		assert.NotNil(t, middleware)
		assert.Equal(t, "permission", middleware.Name())
		assert.Equal(t, 10, middleware.Priority())
	})

	t.Run("WithNilConfig", func(t *testing.T) {
		middleware, err := NewPermissionMiddleware(nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, middleware)
	})
}

func TestPermissionMiddleware_LoadConfig(t *testing.T) {
	config := DefaultPermissionConfig()
	config.Strategy = "noop"
	middleware, _ := NewPermissionMiddleware(config, nil)

	newConfig := map[string]interface{}{
		"enabled":      false,
		"message":      "自定义消息",
		"status_code":  403,
	}

	err := middleware.LoadConfig(newConfig)
	assert.NoError(t, err)
	assert.False(t, middleware.config.Enabled)
	assert.Equal(t, "自定义消息", middleware.config.Message)
}

func TestPermissionMiddleware_ValidateConfig(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := DefaultPermissionConfig()
		config.Strategy = "noop"
		middleware, _ := NewPermissionMiddleware(config, nil)

		err := middleware.ValidateConfig()
		assert.NoError(t, err)
	})

	t.Run("InvalidStatusCode", func(t *testing.T) {
		config := &PermissionConfig{
			Enabled:    true,
			Strategy:   "noop",
			StatusCode: 999, // 无效状态码
		}
		middleware, _ := NewPermissionMiddleware(config, nil)

		err := middleware.ValidateConfig()
		assert.Error(t, err)
	})
}

func TestPermissionMiddleware_Reload(t *testing.T) {
	config := DefaultPermissionConfig()
	config.Strategy = "noop"
	middleware, _ := NewPermissionMiddleware(config, nil)

	reloadConfig := map[string]interface{}{
		"enabled": false,
		"message": "新消息",
	}

	err := middleware.Reload(reloadConfig)
	assert.NoError(t, err)
	assert.False(t, middleware.config.Enabled)
	assert.Equal(t, "新消息", middleware.config.Message)
}

// ========== 辅助函数测试 ==========

func TestGetResourceFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/api/v1/projects", "project"},
		{"/api/v1/projects/123", "project"},
		{"/api/v1/users", "user"},
		{"/api/v1/documents/456", "document"}, // 单数形式
		{"/api/v1/books", "book"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := getResourceFromPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetActionFromMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "read"},
		{"POST", "create"},
		{"PUT", "update"},
		{"PATCH", "update"},
		{"DELETE", "delete"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := getActionFromMethod(tt.method)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchPermission(t *testing.T) {
	t.Run("ExactMatch", func(t *testing.T) {
		assert.True(t, matchPermission("project:read", "project:read"))
	})

	t.Run("WildcardAll", func(t *testing.T) {
		assert.True(t, matchPermission("*:*", "project:read"))
	})

	t.Run("WildcardResource", func(t *testing.T) {
		assert.True(t, matchPermission("project:*", "project:read"))
		assert.True(t, matchPermission("project:*", "project:write"))
		assert.False(t, matchPermission("project:*", "document:read"))
	})

	t.Run("WildcardAction", func(t *testing.T) {
		assert.True(t, matchPermission("*:read", "project:read"))
		assert.True(t, matchPermission("*:read", "document:read"))
		assert.False(t, matchPermission("*:read", "project:write"))
	})

	t.Run("NoMatch", func(t *testing.T) {
		assert.False(t, matchPermission("project:read", "document:write"))
	})
}
