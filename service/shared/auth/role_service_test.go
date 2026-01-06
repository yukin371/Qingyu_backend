package auth

import (
	authModel "Qingyu_backend/models/auth"
	"context"
	"fmt"
	"testing"
	"time"
)

// MockAuthRepository Mock认证Repository
type MockAuthRepository struct {
	roles     map[string]*authModel.Role
	userRoles map[string][]string // userID -> []roleID
	nextID    int
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		roles:     make(map[string]*authModel.Role),
		userRoles: make(map[string][]string),
		nextID:    1,
	}
}

func (m *MockAuthRepository) CreateRole(ctx context.Context, role *authModel.Role) error {
	role.ID = fmt.Sprintf("role_%d", m.nextID)
	m.nextID++
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()
	m.roles[role.ID] = role
	return nil
}

func (m *MockAuthRepository) GetRole(ctx context.Context, roleID string) (*authModel.Role, error) {
	role, ok := m.roles[roleID]
	if !ok {
		return nil, fmt.Errorf("角色不存在: %s", roleID)
	}
	return role, nil
}

func (m *MockAuthRepository) GetRoleByName(ctx context.Context, name string) (*authModel.Role, error) {
	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, fmt.Errorf("角色不存在: %s", name)
}

func (m *MockAuthRepository) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error {
	role, ok := m.roles[roleID]
	if !ok {
		return fmt.Errorf("角色不存在: %s", roleID)
	}

	if name, ok := updates["name"].(string); ok {
		role.Name = name
	}
	if desc, ok := updates["description"].(string); ok {
		role.Description = desc
	}
	if perms, ok := updates["permissions"].([]string); ok {
		role.Permissions = perms
	}
	role.UpdatedAt = time.Now()

	return nil
}

func (m *MockAuthRepository) DeleteRole(ctx context.Context, roleID string) error {
	role, ok := m.roles[roleID]
	if !ok {
		return fmt.Errorf("角色不存在: %s", roleID)
	}

	if role.IsSystem {
		return fmt.Errorf("不能删除系统角色: %s", role.Name)
	}

	delete(m.roles, roleID)
	return nil
}

func (m *MockAuthRepository) ListRoles(ctx context.Context) ([]*authModel.Role, error) {
	roles := make([]*authModel.Role, 0, len(m.roles))
	for _, role := range m.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

func (m *MockAuthRepository) AssignUserRole(ctx context.Context, userID, roleID string) error {
	if _, ok := m.roles[roleID]; !ok {
		return fmt.Errorf("角色不存在: %s", roleID)
	}

	// 检查是否已分配
	for _, rid := range m.userRoles[userID] {
		if rid == roleID {
			return nil // 已存在
		}
	}

	m.userRoles[userID] = append(m.userRoles[userID], roleID)
	return nil
}

func (m *MockAuthRepository) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	roleIDs := m.userRoles[userID]
	newRoleIDs := make([]string, 0)
	for _, rid := range roleIDs {
		if rid != roleID {
			newRoleIDs = append(newRoleIDs, rid)
		}
	}
	m.userRoles[userID] = newRoleIDs
	return nil
}

func (m *MockAuthRepository) GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error) {
	roleIDs := m.userRoles[userID]
	roles := make([]*authModel.Role, 0, len(roleIDs))
	for _, rid := range roleIDs {
		if role, ok := m.roles[rid]; ok {
			roles = append(roles, role)
		}
	}
	return roles, nil
}

func (m *MockAuthRepository) HasUserRole(ctx context.Context, userID, roleID string) (bool, error) {
	for _, rid := range m.userRoles[userID] {
		if rid == roleID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockAuthRepository) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	role, ok := m.roles[roleID]
	if !ok {
		return nil, fmt.Errorf("角色不存在: %s", roleID)
	}
	return role.Permissions, nil
}

func (m *MockAuthRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	roles, _ := m.GetUserRoles(ctx, userID)
	permMap := make(map[string]bool)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			permMap[perm] = true
		}
	}

	permissions := make([]string, 0, len(permMap))
	for perm := range permMap {
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

func (m *MockAuthRepository) Health(ctx context.Context) error {
	return nil
}

// ============ 测试用例 ============

// TestCreateRole 测试创建角色
func TestCreateRole(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	req := &CreateRoleRequest{
		Name:        "editor",
		Description: "编辑者角色",
		Permissions: []string{"book.read", "book.write"},
	}

	role, err := service.CreateRole(ctx, req)
	if err != nil {
		t.Fatalf("创建角色失败: %v", err)
	}

	if role.Name != "editor" {
		t.Errorf("角色名称错误: %s", role.Name)
	}
	if len(role.Permissions) != 2 {
		t.Errorf("权限数量错误: %d", len(role.Permissions))
	}

	t.Logf("创建角色成功: %+v", role)
}

// TestCreateRole_Duplicate 测试创建重复角色
func TestCreateRole_Duplicate(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	req := &CreateRoleRequest{
		Name:        "editor",
		Description: "编辑者角色",
		Permissions: []string{"book.read"},
	}

	// 第一次创建
	_, err := service.CreateRole(ctx, req)
	if err != nil {
		t.Fatalf("创建角色失败: %v", err)
	}

	// 第二次创建（应该失败）
	_, err = service.CreateRole(ctx, req)
	if err == nil {
		t.Fatal("创建重复角色应该失败，但成功了")
	}

	t.Logf("正确拒绝了重复角色: %v", err)
}

// TestUpdateRole 测试更新角色
func TestUpdateRole(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	// 创建角色
	createReq := &CreateRoleRequest{
		Name:        "editor",
		Description: "编辑者",
		Permissions: []string{"book.read"},
	}
	role, _ := service.CreateRole(ctx, createReq)

	// 更新角色
	updateReq := &UpdateRoleRequest{
		Description: "高级编辑者",
		Permissions: []string{"book.read", "book.write", "book.delete"},
	}
	err := service.UpdateRole(ctx, role.ID, updateReq)
	if err != nil {
		t.Fatalf("更新角色失败: %v", err)
	}

	// 验证更新
	updated, _ := service.GetRole(ctx, role.ID)
	if updated.Description != "高级编辑者" {
		t.Errorf("描述未更新: %s", updated.Description)
	}
	if len(updated.Permissions) != 3 {
		t.Errorf("权限数量错误: %d", len(updated.Permissions))
	}

	t.Logf("更新角色成功: %+v", updated)
}

// TestDeleteRole 测试删除角色
func TestDeleteRole(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	// 创建角色
	req := &CreateRoleRequest{
		Name:        "temp_role",
		Description: "临时角色",
		Permissions: []string{},
	}
	role, _ := service.CreateRole(ctx, req)

	// 删除角色
	err := service.DeleteRole(ctx, role.ID)
	if err != nil {
		t.Fatalf("删除角色失败: %v", err)
	}

	// 验证已删除
	_, err = service.GetRole(ctx, role.ID)
	if err == nil {
		t.Fatal("角色应该已删除，但仍然存在")
	}

	t.Logf("删除角色成功")
}

// TestDeleteRole_System 测试删除系统角色
func TestDeleteRole_System(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	// 创建系统角色
	systemRole := &authModel.Role{
		ID:          "role_system",
		Name:        "admin",
		Description: "系统管理员",
		Permissions: []string{"*"},
		IsSystem:    true,
	}
	repo.roles[systemRole.ID] = systemRole

	// 尝试删除（应该失败）
	err := service.DeleteRole(ctx, systemRole.ID)
	if err == nil {
		t.Fatal("删除系统角色应该失败，但成功了")
	}

	t.Logf("正确拒绝了删除系统角色: %v", err)
}

// TestListRoles 测试列出角色
func TestListRoles(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	// 创建多个角色
	roles := []string{"reader", "author", "editor"}
	for _, name := range roles {
		req := &CreateRoleRequest{
			Name:        name,
			Description: fmt.Sprintf("%s角色", name),
			Permissions: []string{"book.read"},
		}
		_, _ = service.CreateRole(ctx, req)
	}

	// 列出角色
	list, err := service.ListRoles(ctx)
	if err != nil {
		t.Fatalf("列出角色失败: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("角色数量错误: 期望3个，实际%d个", len(list))
	}

	t.Logf("列出角色成功: %d个角色", len(list))
}

// TestAssignPermissions 测试分配权限
func TestAssignPermissions(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	// 创建角色
	req := &CreateRoleRequest{
		Name:        "editor",
		Description: "编辑者",
		Permissions: []string{"book.read"},
	}
	role, _ := service.CreateRole(ctx, req)

	// 分配权限
	newPerms := []string{"book.write", "book.delete"}
	err := service.AssignPermissions(ctx, role.ID, newPerms)
	if err != nil {
		t.Fatalf("分配权限失败: %v", err)
	}

	// 验证权限
	updated, _ := service.GetRole(ctx, role.ID)
	if len(updated.Permissions) < 3 {
		t.Errorf("权限数量错误: %d", len(updated.Permissions))
	}

	t.Logf("分配权限成功: %v", updated.Permissions)
}

// TestRemovePermissions 测试移除权限
func TestRemovePermissions(t *testing.T) {
	repo := NewMockAuthRepository()
	service := NewRoleService(repo)
	ctx := context.Background()

	// 创建角色
	req := &CreateRoleRequest{
		Name:        "editor",
		Description: "编辑者",
		Permissions: []string{"book.read", "book.write", "book.delete"},
	}
	role, _ := service.CreateRole(ctx, req)

	// 移除权限
	removePerms := []string{"book.delete"}
	err := service.RemovePermissions(ctx, role.ID, removePerms)
	if err != nil {
		t.Fatalf("移除权限失败: %v", err)
	}

	// 验证权限
	updated, _ := service.GetRole(ctx, role.ID)
	if len(updated.Permissions) != 2 {
		t.Errorf("权限数量错误: %d", len(updated.Permissions))
	}

	// 检查book.delete是否已移除
	for _, perm := range updated.Permissions {
		if perm == "book.delete" {
			t.Error("book.delete权限未被移除")
		}
	}

	t.Logf("移除权限成功: %v", updated.Permissions)
}
