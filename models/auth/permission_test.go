package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPermission_HasPermission_ExactMatch 测试精确匹配
func TestPermission_HasPermission_ExactMatch(t *testing.T) {
	tests := []struct {
		name          string
		permission    Permission
		requiredCode  string
		expected      bool
	}{
		{
			name: "精确匹配-允许",
			permission: Permission{
				Code:   "user.read",
				Effect: "allow",
			},
			requiredCode: "user.read",
			expected:     true,
		},
		{
			name: "精确匹配-拒绝",
			permission: Permission{
				Code:   "user.read",
				Effect: "deny",
			},
			requiredCode: "user.read",
			expected:     false,
		},
		{
			name: "不匹配",
			permission: Permission{
				Code:   "user.read",
				Effect: "allow",
			},
			requiredCode: "user.write",
			expected:     false,
		},
		{
			name: "通配符匹配所有",
			permission: Permission{
				Code:   "*",
				Effect: "allow",
			},
			requiredCode: "any.permission",
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.permission.HasPermission(tt.requiredCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermission_HasPermission_WildcardMatch 测试通配符匹配
func TestPermission_HasPermission_WildcardMatch(t *testing.T) {
	tests := []struct {
		name          string
		permission    Permission
		requiredCode  string
		expected      bool
	}{
		{
			name: "user.* 匹配 user.read",
			permission: Permission{
				Code:   "user.*",
				Effect: "allow",
			},
			requiredCode: "user.read",
			expected:     true,
		},
		{
			name: "user.* 匹配 user.write",
			permission: Permission{
				Code:   "user.*",
				Effect: "allow",
			},
			requiredCode: "user.write",
			expected:     true,
		},
		{
			name: "user.* 不匹配 book.read",
			permission: Permission{
				Code:   "user.*",
				Effect: "allow",
			},
			requiredCode: "book.read",
			expected:     false,
		},
		{
			name: "content.* 不匹配 content",
			permission: Permission{
				Code:   "content.*",
				Effect: "allow",
			},
			requiredCode: "content",
			expected:     false,
		},
		{
			name: "admin.* 匹配 admin.manage",
			permission: Permission{
				Code:   "admin.*",
				Effect: "allow",
			},
			requiredCode: "admin.manage",
			expected:     true,
		},
		{
			name: "通配符拒绝",
			permission: Permission{
				Code:   "user.*",
				Effect: "deny",
			},
			requiredCode: "user.read",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.permission.HasPermission(tt.requiredCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermission_HasPermission_EdgeCases 测试边界情况
func TestPermission_HasPermission_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		permission    Permission
		requiredCode  string
		expected      bool
	}{
		{
			name: "空代码",
			permission: Permission{
				Code:   "",
				Effect: "allow",
			},
			requiredCode: "user.read",
			expected:     false,
		},
		{
			name: "空需求",
			permission: Permission{
				Code:   "user.read",
				Effect: "allow",
			},
			requiredCode: "",
			expected:     false,
		},
		{
			name: "多层通配符-匹配成功",
			permission: Permission{
				Code:   "book.chapter.*",
				Effect: "allow",
			},
			requiredCode: "book.chapter.read",
			expected:     true, // book.chapter.* 会匹配任何以 book.chapter. 开头的权限
		},
		{
			name: "特殊字符",
			permission: Permission{
				Code:   "user.read",
				Effect: "allow",
			},
			requiredCode: "user.read.extra",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.permission.HasPermission(tt.requiredCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermission_GetHigherPriority 测试优先级判断
func TestPermission_GetHigherPriority(t *testing.T) {
	tests := []struct {
		name     string
		p1       Permission
		p2       Permission
		expected string // 返回优先级高的权限的Code
	}{
		{
			name: "deny优先于allow",
			p1: Permission{
				Code:     "user.read",
				Effect:   "deny",
				Priority: 1,
			},
			p2: Permission{
				Code:     "user.read",
				Effect:   "allow",
				Priority: 10,
			},
			expected: "user.read", // deny优先
		},
		{
			name: "allow对比deny-deny优先",
			p1: Permission{
				Code:     "user.read",
				Effect:   "allow",
				Priority: 10,
			},
			p2: Permission{
				Code:     "user.read",
				Effect:   "deny",
				Priority: 1,
			},
			expected: "user.read", // deny优先
		},
		{
			name: "allow对比allow-高priority优先",
			p1: Permission{
				Code:     "user.read",
				Effect:   "allow",
				Priority: 5,
			},
			p2: Permission{
				Code:     "user.read",
				Effect:   "allow",
				Priority: 10,
			},
			expected: "user.read", // p2优先
		},
		{
			name: "deny对比deny-高priority优先",
			p1: Permission{
				Code:     "user.read",
				Effect:   "deny",
				Priority: 5,
			},
			p2: Permission{
				Code:     "user.read",
				Effect:   "deny",
				Priority: 10,
			},
			expected: "user.read", // p2优先
		},
		{
			name: "相同priority-返回p1",
			p1: Permission{
				Code:     "user.read",
				Effect:   "allow",
				Priority: 5,
			},
			p2: Permission{
				Code:     "user.write",
				Effect:   "allow",
				Priority: 5,
			},
			expected: "user.read", // p1优先（相同优先级返回第一个）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.p1.GetHigherPriority(&tt.p2)
			assert.Equal(t, tt.expected, result.Code)
		})
	}
}

// TestResourcePermission_IsResourceAllowed 测试资源级权限检查
func TestResourcePermission_IsResourceAllowed(t *testing.T) {
	tests := []struct {
		name              string
		resourcePermission ResourcePermission
		resourceID        string
		expected          bool
	}{
		{
			name: "允许所有资源",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  []string{}, // 空表示所有资源
				Effect:       "allow",
			},
			resourceID: "any-resource-id",
			expected:   true,
		},
		{
			name: "允许特定资源-匹配",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  []string{"resource1", "resource2"},
				Effect:       "allow",
			},
			resourceID: "resource1",
			expected:   true,
		},
		{
			name: "允许特定资源-不匹配",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  []string{"resource1", "resource2"},
				Effect:       "allow",
			},
			resourceID: "resource3",
			expected:   false,
		},
		{
			name: "拒绝特定资源",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  []string{"resource1"},
				Effect:       "deny",
			},
			resourceID: "resource1",
			expected:   false,
		},
		{
			name: "通配符资源ID",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  []string{"*"},
				Effect:       "allow",
			},
			resourceID: "any-resource-id",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resourcePermission.IsResourceAllowed(tt.resourceID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestResourcePermission_IsResourceAllowed_EdgeCases 测试资源级权限边界情况
func TestResourcePermission_IsResourceAllowed_EdgeCases(t *testing.T) {
	tests := []struct {
		name              string
		resourcePermission ResourcePermission
		resourceID        string
		expected          bool
	}{
		{
			name: "空资源ID",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  []string{},
				Effect:       "allow",
			},
			resourceID: "",
			expected:   true, // 空ResourceIDs表示所有资源
		},
		{
			name: "nil资源ID列表",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  nil,
				Effect:       "allow",
			},
			resourceID: "any-resource",
			expected:   true, // nil也表示所有资源
		},
		{
			name: "资源列表包含空字符串",
			resourcePermission: ResourcePermission{
				PermissionID: "perm1",
				ResourceIDs:  []string{""},
				Effect:       "allow",
			},
			resourceID: "",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resourcePermission.IsResourceAllowed(tt.resourceID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPermission_Validation 测试权限验证
func TestPermission_Validation(t *testing.T) {
	tests := []struct {
		name       string
		permission Permission
		wantErr    bool
		errMsg     string
	}{
		{
			name: "有效权限-allow",
			permission: Permission{
				Code:        "user.read",
				Name:        "读取用户",
				Description: "允许读取用户信息",
				Effect:      "allow",
				Priority:    1,
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "有效权限-deny",
			permission: Permission{
				Code:        "user.delete",
				Name:        "删除用户",
				Description: "禁止删除用户",
				Effect:      "deny",
				Priority:    10,
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "无效effect",
			permission: Permission{
				Code:        "user.read",
				Name:        "读取用户",
				Description: "允许读取用户信息",
				Effect:      "invalid",
				Priority:    1,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
			errMsg:  "effect must be 'allow' or 'deny'",
		},
		{
			name: "空code",
			permission: Permission{
				Code:        "",
				Name:        "空权限",
				Description: "权限代码不能为空",
				Effect:      "allow",
				Priority:    1,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
			errMsg:  "code cannot be empty",
		},
		{
			name: "负数priority",
			permission: Permission{
				Code:        "user.read",
				Name:        "读取用户",
				Description: "允许读取用户信息",
				Effect:      "allow",
				Priority:    -1,
				CreatedAt:   time.Now(),
			},
			wantErr: true,
			errMsg:  "priority cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.permission.Validate()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, strings.ToLower(err.Error()), strings.ToLower(tt.errMsg))
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestPermission_Constants 测试权限常量
func TestPermission_Constants(t *testing.T) {
	// 验证预定义权限常量格式
	assert.Equal(t, "user.read", PermUserRead)
	assert.Equal(t, "user.write", PermUserWrite)
	assert.Equal(t, "user.delete", PermUserDelete)
	assert.Equal(t, "book.read", PermBookRead)
	assert.Equal(t, "book.write", PermBookWrite)
	assert.Equal(t, "admin.access", PermAdminAccess)

	// 验证Effect常量
	assert.Equal(t, "allow", EffectAllow)
	assert.Equal(t, "deny", EffectDeny)

	// 验证错误常量存在且不为nil
	assert.NotNil(t, ErrPermissionCodeEmpty)
	assert.NotNil(t, ErrPermissionEffectInvalid)
	assert.NotNil(t, ErrPermissionPriorityNeg)

	// 验证错误消息
	assert.Contains(t, ErrPermissionCodeEmpty.Error(), "code cannot be empty")
	assert.Contains(t, ErrPermissionEffectInvalid.Error(), "effect must be")
	assert.Contains(t, ErrPermissionPriorityNeg.Error(), "priority cannot be negative")
}

// TestRole_Constants 测试角色常量
func TestRole_Constants(t *testing.T) {
	// 验证预定义角色常量
	assert.Equal(t, "reader", RoleReader)
	assert.Equal(t, "author", RoleAuthor)
	assert.Equal(t, "admin", RoleAdmin)
}
