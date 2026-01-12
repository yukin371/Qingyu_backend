package users

import (
	"testing"

	"Qingyu_backend/models/auth"
)

// TestGetEffectiveRoles 测试角色继承逻辑
func TestGetEffectiveRoles(t *testing.T) {
	tests := []struct {
		name          string
		roles         []string
		expectedRoles []string
	}{
		{
			name:          "普通读者",
			roles:         []string{auth.RoleReader},
			expectedRoles: []string{auth.RoleReader},
		},
		{
			name:          "作者（应包含reader）",
			roles:         []string{auth.RoleAuthor},
			expectedRoles: []string{auth.RoleAuthor, auth.RoleReader},
		},
		{
			name:          "管理员（应包含author和reader）",
			roles:         []string{auth.RoleAdmin},
			expectedRoles: []string{auth.RoleAdmin, auth.RoleAuthor, auth.RoleReader},
		},
		{
			name:          "作者+读者（显式指定）",
			roles:         []string{auth.RoleReader, auth.RoleAuthor},
			expectedRoles: []string{auth.RoleReader, auth.RoleAuthor},
		},
		{
			name:          "全部角色",
			roles:         []string{auth.RoleReader, auth.RoleAuthor, auth.RoleAdmin},
			expectedRoles: []string{auth.RoleReader, auth.RoleAuthor, auth.RoleAdmin},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Roles: tt.roles}
			effectiveRoles := user.GetEffectiveRoles()

			// 验证角色数量
			if len(effectiveRoles) != len(tt.expectedRoles) {
				t.Errorf("GetEffectiveRoles() returned %d roles, expected %d",
					len(effectiveRoles), len(tt.expectedRoles))
			}

			// 验证所有预期角色都存在
			expectedSet := make(map[string]bool)
			for _, r := range tt.expectedRoles {
				expectedSet[r] = true
			}

			for _, role := range effectiveRoles {
				if !expectedSet[role] {
					t.Errorf("GetEffectiveRoles() returned unexpected role: %s", role)
				}
			}
		})
	}
}

// TestHasRole 测试角色检查方法
func TestHasRole(t *testing.T) {
	tests := []struct {
		name     string
		roles    []string
		role     string
		expected bool
	}{
		{
			name:     "普通读者检查reader",
			roles:    []string{auth.RoleReader},
			role:     auth.RoleReader,
			expected: true,
		},
		{
			name:     "普通读者检查author",
			roles:    []string{auth.RoleReader},
			role:     auth.RoleAuthor,
			expected: false,
		},
		{
			name:     "作者检查author（含继承）",
			roles:    []string{auth.RoleAuthor},
			role:     auth.RoleAuthor,
			expected: true,
		},
		{
			name:     "作者检查reader（含继承）",
			roles:    []string{auth.RoleAuthor},
			role:     auth.RoleReader,
			expected: true,
		},
		{
			name:     "管理员检查admin（含继承）",
			roles:    []string{auth.RoleAdmin},
			role:     auth.RoleAdmin,
			expected: true,
		},
		{
			name:     "管理员检查reader（含继承）",
			roles:    []string{auth.RoleAdmin},
			role:     auth.RoleReader,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Roles: tt.roles}
			result := user.HasRole(tt.role)

			if result != tt.expected {
				t.Errorf("HasRole(%s) = %v, expected %v", tt.role, result, tt.expected)
			}
		})
	}
}

// TestHasAnyRole 测试HasAnyRole方法
func TestHasAnyRole(t *testing.T) {
	tests := []struct {
		name       string
		roles      []string
		checkRoles []string
		expected   bool
	}{
		{
			name:       "普通读者检查reader或author",
			roles:      []string{auth.RoleReader},
			checkRoles: []string{auth.RoleReader, auth.RoleAuthor},
			expected:   true,
		},
		{
			name:       "普通读者检查author或admin",
			roles:      []string{auth.RoleReader},
			checkRoles: []string{auth.RoleAuthor, auth.RoleAdmin},
			expected:   false,
		},
		{
			name:       "作者检查author或admin",
			roles:      []string{auth.RoleAuthor},
			checkRoles: []string{auth.RoleAuthor, auth.RoleAdmin},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Roles: tt.roles}
			result := user.HasAnyRole(tt.checkRoles...)

			if result != tt.expected {
				t.Errorf("HasAnyRole(%v) = %v, expected %v", tt.checkRoles, result, tt.expected)
			}
		})
	}
}

// TestHasAllRoles 测试HasAllRoles方法
func TestHasAllRoles(t *testing.T) {
	tests := []struct {
		name       string
		roles      []string
		checkRoles []string
		expected   bool
	}{
		{
			name:       "普通读者检查reader",
			roles:      []string{auth.RoleReader},
			checkRoles: []string{auth.RoleReader},
			expected:   true,
		},
		{
			name:       "普通读者检查reader和author",
			roles:      []string{auth.RoleReader},
			checkRoles: []string{auth.RoleReader, auth.RoleAuthor},
			expected:   false,
		},
		{
			name:       "作者检查reader和author（含继承）",
			roles:      []string{auth.RoleAuthor},
			checkRoles: []string{auth.RoleReader, auth.RoleAuthor},
			expected:   true,
		},
		{
			name:       "管理员检查全部角色（含继承）",
			roles:      []string{auth.RoleAdmin},
			checkRoles: []string{auth.RoleReader, auth.RoleAuthor, auth.RoleAdmin},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Roles: tt.roles}
			result := user.HasAllRoles(tt.checkRoles...)

			if result != tt.expected {
				t.Errorf("HasAllRoles(%v) = %v, expected %v", tt.checkRoles, result, tt.expected)
			}
		})
	}
}

// TestIsAdmin 测试IsAdmin方法
func TestIsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		roles    []string
		expected bool
	}{
		{
			name:     "普通读者",
			roles:    []string{auth.RoleReader},
			expected: false,
		},
		{
			name:     "作者",
			roles:    []string{auth.RoleAuthor},
			expected: false,
		},
		{
			name:     "管理员",
			roles:    []string{auth.RoleAdmin},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Roles: tt.roles}
			result := user.IsAdmin()

			if result != tt.expected {
				t.Errorf("IsAdmin() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestIsAuthor 测试IsAuthor方法
func TestIsAuthor(t *testing.T) {
	tests := []struct {
		name     string
		roles    []string
		expected bool
	}{
		{
			name:     "普通读者",
			roles:    []string{auth.RoleReader},
			expected: false,
		},
		{
			name:     "作者",
			roles:    []string{auth.RoleAuthor},
			expected: true,
		},
		{
			name:     "管理员（应含author）",
			roles:    []string{auth.RoleAdmin},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Roles: tt.roles}
			result := user.IsAuthor()

			if result != tt.expected {
				t.Errorf("IsAuthor() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestAddRole 测试AddRole方法
func TestAddRole(t *testing.T) {
	user := &User{Roles: []string{auth.RoleReader}}

	// 添加新角色
	user.AddRole(auth.RoleAuthor)
	if len(user.Roles) != 2 {
		t.Errorf("AddRole() failed: expected 2 roles, got %d", len(user.Roles))
	}

	// 添加已存在的角色（不应重复）
	user.AddRole(auth.RoleAuthor)
	if len(user.Roles) != 2 {
		t.Errorf("AddRole() added duplicate role: expected 2 roles, got %d", len(user.Roles))
	}
}

// TestRemoveRole 测试RemoveRole方法
func TestRemoveRole(t *testing.T) {
	user := &User{Roles: []string{auth.RoleReader, auth.RoleAuthor, auth.RoleAdmin}}

	// 移除中间的角色
	user.RemoveRole(auth.RoleAuthor)
	if len(user.Roles) != 2 {
		t.Errorf("RemoveRole() failed: expected 2 roles, got %d", len(user.Roles))
	}

	// 验证被移除的角色不存在
	for _, role := range user.Roles {
		if role == auth.RoleAuthor {
			t.Error("RemoveRole() failed: removed role still exists")
		}
	}
}

// TestGetVIPLevel 测试VIP等级标准化
func TestGetVIPLevel(t *testing.T) {
	tests := []struct {
		name     string
		vipLevel int
		expected int
	}{
		{
			name:     "正常等级3",
			vipLevel: 3,
			expected: 3,
		},
		{
			name:     "负数等级",
			vipLevel: -1,
			expected: 0,
		},
		{
			name:     "过大等级",
			vipLevel: 10,
			expected: 5,
		},
		{
			name:     "零等级",
			vipLevel: 0,
			expected: 0,
		},
		{
			name:     "最大等级",
			vipLevel: 5,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{VIPLevel: tt.vipLevel}
			result := user.GetVIPLevel()

			if result != tt.expected {
				t.Errorf("GetVIPLevel() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestIsVIP 测试IsVIP方法
func TestIsVIP(t *testing.T) {
	tests := []struct {
		name     string
		vipLevel int
		expected bool
	}{
		{
			name:     "非VIP",
			vipLevel: 0,
			expected: false,
		},
		{
			name:     "VIP等级1",
			vipLevel: 1,
			expected: true,
		},
		{
			name:     "VIP等级5",
			vipLevel: 5,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{VIPLevel: tt.vipLevel}
			result := user.IsVIP()

			if result != tt.expected {
				t.Errorf("IsVIP() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestHasVIPLevel 测试HasVIPLevel方法
func TestHasVIPLevel(t *testing.T) {
	tests := []struct {
		name          string
		vipLevel      int
		requiredLevel int
		expected      bool
	}{
		{
			name:          "等级3满足等级2要求",
			vipLevel:      3,
			requiredLevel: 2,
			expected:      true,
		},
		{
			name:          "等级2不满足等级3要求",
			vipLevel:      2,
			requiredLevel: 3,
			expected:      false,
		},
		{
			name:          "等级0不满足等级1要求",
			vipLevel:      0,
			requiredLevel: 1,
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{VIPLevel: tt.vipLevel}
			result := user.HasVIPLevel(tt.requiredLevel)

			if result != tt.expected {
				t.Errorf("HasVIPLevel(%d) = %v, expected %v", tt.requiredLevel, result, tt.expected)
			}
		})
	}
}

// TestSetVIPLevel 测试SetVIPLevel方法
func TestSetVIPLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "设置正常等级",
			input:    3,
			expected: 3,
		},
		{
			name:     "设置负数等级",
			input:    -1,
			expected: 0,
		},
		{
			name:     "设置过大等级",
			input:    10,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{}
			user.SetVIPLevel(tt.input)

			if user.VIPLevel != tt.expected {
				t.Errorf("SetVIPLevel(%d) = %v, expected %v", tt.input, user.VIPLevel, tt.expected)
			}
		})
	}
}

// TestMultiRoleScenario 测试多角色实际场景
func TestMultiRoleScenario(t *testing.T) {
	// 场景：一个用户既是读者又是作者
	user := &User{
		Roles:    []string{auth.RoleReader, auth.RoleAuthor},
		VIPLevel: 2,
	}

	// 验证角色检查
	if !user.IsReader() {
		t.Error("User should be a reader")
	}
	if !user.IsAuthor() {
		t.Error("User should be an author")
	}
	if user.IsAdmin() {
		t.Error("User should not be an admin")
	}

	// 验证VIP
	if !user.IsVIP() {
		t.Error("User should be VIP")
	}
	if !user.HasVIPLevel(2) {
		t.Error("User should have VIP level 2")
	}
	if user.HasVIPLevel(3) {
		t.Error("User should not have VIP level 3")
	}

	// 添加管理员角色
	user.AddRole(auth.RoleAdmin)

	// 验证现在拥有所有角色
	if !user.HasAllRoles(auth.RoleReader, auth.RoleAuthor, auth.RoleAdmin) {
		t.Error("User should have all roles")
	}

	// 移除作者角色
	user.RemoveRole(auth.RoleAuthor)

	// 验证仍然拥有author权限（通过admin继承）
	if !user.IsAuthor() {
		t.Error("Admin should still have author permissions through inheritance")
	}
}
