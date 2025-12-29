package auth

import (
	authModel "Qingyu_backend/models/auth"
	"context"
	"fmt"
	"testing"
	"time"
)

// MockCacheClient Mock缓存客户端
type MockCacheClient struct {
	data map[string]string
}

func NewMockCacheClient() *MockCacheClient {
	return &MockCacheClient{
		data: make(map[string]string),
	}
}

func (m *MockCacheClient) Get(ctx context.Context, key string) (string, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

func (m *MockCacheClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// 将interface{}转换为string存储
	if strValue, ok := value.(string); ok {
		m.data[key] = strValue
	} else {
		m.data[key] = fmt.Sprintf("%v", value)
	}
	return nil
}

func (m *MockCacheClient) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// ============ 测试用例 ============

// TestCheckPermission 测试权限检查
func TestCheckPermission(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建角色和用户
	role := &authModel.Role{
		ID:          "role_1",
		Name:        "editor",
		Permissions: []string{"book.read", "book.write"},
	}
	repo.roles[role.ID] = role
	repo.userRoles["user_1"] = []string{"role_1"}

	// 测试有权限
	hasRead, err := service.CheckPermission(ctx, "user_1", "book.read")
	if err != nil {
		t.Fatalf("检查权限失败: %v", err)
	}
	if !hasRead {
		t.Error("用户应该有book.read权限")
	}

	// 测试无权限
	hasDelete, _ := service.CheckPermission(ctx, "user_1", "book.delete")
	if hasDelete {
		t.Error("用户不应该有book.delete权限")
	}

	t.Logf("权限检查测试通过")
}

// TestCheckPermission_Wildcard 测试通配符权限
func TestCheckPermission_Wildcard(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建具有通配符权限的角色
	role := &authModel.Role{
		ID:          "role_1",
		Name:        "admin",
		Permissions: []string{"*"}, // 全部权限
	}
	repo.roles[role.ID] = role
	repo.userRoles["user_1"] = []string{"role_1"}

	// 测试各种权限
	permissions := []string{"book.read", "book.write", "book.delete", "user.create", "anything"}
	for _, perm := range permissions {
		has, err := service.CheckPermission(ctx, "user_1", perm)
		if err != nil {
			t.Fatalf("检查权限失败: %v", err)
		}
		if !has {
			t.Errorf("用户应该有%s权限（通配符*）", perm)
		}
	}

	t.Logf("通配符权限测试通过")
}

// TestCheckPermission_PatternMatch 测试模式匹配权限
func TestCheckPermission_PatternMatch(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建具有模式匹配权限的角色
	role := &authModel.Role{
		ID:          "role_1",
		Name:        "book_manager",
		Permissions: []string{"book.*"}, // book模块的所有权限
	}
	repo.roles[role.ID] = role
	repo.userRoles["user_1"] = []string{"role_1"}

	// 测试匹配的权限
	matchPerms := []string{"book.read", "book.write", "book.delete"}
	for _, perm := range matchPerms {
		has, _ := service.CheckPermission(ctx, "user_1", perm)
		if !has {
			t.Errorf("用户应该有%s权限（匹配book.*）", perm)
		}
	}

	// 测试不匹配的权限
	notMatchPerms := []string{"user.read", "comment.write"}
	for _, perm := range notMatchPerms {
		has, _ := service.CheckPermission(ctx, "user_1", perm)
		if has {
			t.Errorf("用户不应该有%s权限（不匹配book.*）", perm)
		}
	}

	t.Logf("模式匹配权限测试通过")
}

// TestGetUserPermissions 测试获取用户权限
func TestGetUserPermissions(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建多个角色
	role1 := &authModel.Role{
		ID:          "role_1",
		Name:        "reader",
		Permissions: []string{"book.read"},
	}
	role2 := &authModel.Role{
		ID:          "role_2",
		Name:        "author",
		Permissions: []string{"book.read", "book.write"},
	}
	repo.roles[role1.ID] = role1
	repo.roles[role2.ID] = role2
	repo.userRoles["user_1"] = []string{"role_1", "role_2"}

	// 获取用户权限
	perms, err := service.GetUserPermissions(ctx, "user_1")
	if err != nil {
		t.Fatalf("获取用户权限失败: %v", err)
	}

	// 权限应该去重（book.read只出现一次）
	if len(perms) != 2 {
		t.Errorf("权限数量错误: 期望2个，实际%d个", len(perms))
	}

	t.Logf("获取用户权限成功: %v", perms)
}

// TestGetUserPermissions_Cache 测试权限缓存
func TestGetUserPermissions_Cache(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		ID:          "role_1",
		Name:        "editor",
		Permissions: []string{"book.read", "book.write"},
	}
	repo.roles[role.ID] = role
	repo.userRoles["user_1"] = []string{"role_1"}

	// 第一次获取（从数据库）
	perms1, _ := service.GetUserPermissions(ctx, "user_1")

	// 第二次获取（从缓存）
	perms2, _ := service.GetUserPermissions(ctx, "user_1")

	// 验证结果一致
	if len(perms1) != len(perms2) {
		t.Errorf("缓存前后权限数量不一致: %d vs %d", len(perms1), len(perms2))
	}

	// 验证缓存Key存在
	cacheKey := fmt.Sprintf("user:permissions:%s", "user_1")
	if _, ok := cache.data[cacheKey]; !ok {
		t.Error("权限缓存未生效")
	}

	t.Logf("权限缓存测试通过")
}

// TestHasRole 测试角色检查
func TestHasRole(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		ID:          "role_1",
		Name:        "editor",
		Permissions: []string{},
	}
	repo.roles[role.ID] = role
	repo.userRoles["user_1"] = []string{"role_1"}

	// 测试有角色
	hasEditor, err := service.HasRole(ctx, "user_1", "editor")
	if err != nil {
		t.Fatalf("检查角色失败: %v", err)
	}
	if !hasEditor {
		t.Error("用户应该有editor角色")
	}

	// 测试无角色
	hasAdmin, _ := service.HasRole(ctx, "user_1", "admin")
	if hasAdmin {
		t.Error("用户不应该有admin角色")
	}

	t.Logf("角色检查测试通过")
}

// TestGetRolePermissions 测试获取角色权限
func TestGetRolePermissions(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建角色
	role := &authModel.Role{
		ID:          "role_1",
		Name:        "editor",
		Permissions: []string{"book.read", "book.write", "book.delete"},
	}
	repo.roles[role.ID] = role

	// 获取角色权限
	perms, err := service.GetRolePermissions(ctx, role.ID)
	if err != nil {
		t.Fatalf("获取角色权限失败: %v", err)
	}

	if len(perms) != 3 {
		t.Errorf("权限数量错误: %d", len(perms))
	}

	t.Logf("获取角色权限成功: %v", perms)
}

// TestMultipleRolesPermissions 测试多角色权限合并
func TestMultipleRolesPermissions(t *testing.T) {
	repo := NewMockAuthRepository()
	cache := NewMockCacheClient()
	service := NewPermissionService(repo, cache)
	ctx := context.Background()

	// 创建多个角色
	role1 := &authModel.Role{
		ID:          "role_1",
		Name:        "reader",
		Permissions: []string{"book.read"},
	}
	role2 := &authModel.Role{
		ID:          "role_2",
		Name:        "author",
		Permissions: []string{"book.write"},
	}
	role3 := &authModel.Role{
		ID:          "role_3",
		Name:        "reviewer",
		Permissions: []string{"book.review"},
	}
	repo.roles[role1.ID] = role1
	repo.roles[role2.ID] = role2
	repo.roles[role3.ID] = role3
	repo.userRoles["user_1"] = []string{"role_1", "role_2", "role_3"}

	// 获取用户权限（应该是所有角色权限的并集）
	perms, _ := service.GetUserPermissions(ctx, "user_1")

	// 应该有3个权限
	if len(perms) != 3 {
		t.Errorf("权限数量错误: 期望3个，实际%d个", len(perms))
	}

	// 检查所有权限都存在
	permMap := make(map[string]bool)
	for _, p := range perms {
		permMap[p] = true
	}

	expectedPerms := []string{"book.read", "book.write", "book.review"}
	for _, expected := range expectedPerms {
		if !permMap[expected] {
			t.Errorf("缺少权限: %s", expected)
		}
	}

	t.Logf("多角色权限合并测试通过: %v", perms)
}
