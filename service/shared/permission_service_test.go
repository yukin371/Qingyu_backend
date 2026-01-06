package shared

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/auth"
)

// MockPermissionRepository 模拟权限仓储
type MockPermissionRepository struct {
	permissions map[string]*auth.Permission
	roles       map[string]*auth.Role
	userRoles   map[string][]string // userID -> roleNames
}

func NewMockPermissionRepository() *MockPermissionRepository {
	return &MockPermissionRepository{
		permissions: make(map[string]*auth.Permission),
		roles:       make(map[string]*auth.Role),
		userRoles:   make(map[string][]string),
	}
}

// ==================== 权限管理 ====================

func (m *MockPermissionRepository) GetAllPermissions(ctx context.Context) ([]*auth.Permission, error) {
	result := make([]*auth.Permission, 0, len(m.permissions))
	for _, p := range m.permissions {
		result = append(result, p)
	}
	return result, nil
}

func (m *MockPermissionRepository) GetPermissionByCode(ctx context.Context, code string) (*auth.Permission, error) {
	if p, ok := m.permissions[code]; ok {
		return p, nil
	}
	return nil, ErrPermissionNotFound
}

func (m *MockPermissionRepository) CreatePermission(ctx context.Context, permission *auth.Permission) error {
	m.permissions[permission.Code] = permission
	return nil
}

func (m *MockPermissionRepository) UpdatePermission(ctx context.Context, permission *auth.Permission) error {
	if _, ok := m.permissions[permission.Code]; ok {
		m.permissions[permission.Code] = permission
		return nil
	}
	return nil
}

func (m *MockPermissionRepository) DeletePermission(ctx context.Context, code string) error {
	delete(m.permissions, code)
	return nil
}

// ==================== 角色管理 ====================

func (m *MockPermissionRepository) GetAllRoles(ctx context.Context) ([]*auth.Role, error) {
	result := make([]*auth.Role, 0, len(m.roles))
	for _, r := range m.roles {
		result = append(result, r)
	}
	return result, nil
}

func (m *MockPermissionRepository) GetRoleByID(ctx context.Context, roleID string) (*auth.Role, error) {
	if r, ok := m.roles[roleID]; ok {
		return r, nil
	}
	return nil, ErrRoleNotFound
}

func (m *MockPermissionRepository) GetRoleByName(ctx context.Context, name string) (*auth.Role, error) {
	for _, r := range m.roles {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, ErrRoleNotFound
}

func (m *MockPermissionRepository) CreateRole(ctx context.Context, role *auth.Role) error {
	role.ID = primitive.NewObjectID().Hex()
	m.roles[role.ID] = role
	return nil
}

func (m *MockPermissionRepository) UpdateRole(ctx context.Context, role *auth.Role) error {
	if _, ok := m.roles[role.ID]; ok {
		m.roles[role.ID] = role
		return nil
	}
	return nil
}

func (m *MockPermissionRepository) DeleteRole(ctx context.Context, roleID string) error {
	delete(m.roles, roleID)
	return nil
}

func (m *MockPermissionRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionCode string) error {
	if r, ok := m.roles[roleID]; ok {
		r.Permissions = append(r.Permissions, permissionCode)
		return nil
	}
	return nil
}

func (m *MockPermissionRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionCode string) error {
	if r, ok := m.roles[roleID]; ok {
		for i, p := range r.Permissions {
			if p == permissionCode {
				r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
				break
			}
		}
		return nil
	}
	return nil
}

func (m *MockPermissionRepository) GetRolePermissions(ctx context.Context, roleID string) ([]*auth.Permission, error) {
	r, err := m.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return []*auth.Permission{}, nil
	}

	result := make([]*auth.Permission, 0)
	for _, code := range r.Permissions {
		if p := m.permissions[code]; p != nil {
			result = append(result, p)
		}
	}
	return result, nil
}

// ==================== 用户角色管理 ====================

func (m *MockPermissionRepository) GetUserRoles(ctx context.Context, userID primitive.ObjectID) ([]string, error) {
	if roles, ok := m.userRoles[userID.Hex()]; ok {
		return roles, nil
	}
	return []string{}, nil
}

func (m *MockPermissionRepository) AssignRoleToUser(ctx context.Context, userID primitive.ObjectID, roleName string) error {
	userIDStr := userID.Hex()
	m.userRoles[userIDStr] = append(m.userRoles[userIDStr], roleName)
	return nil
}

func (m *MockPermissionRepository) RemoveRoleFromUser(ctx context.Context, userID primitive.ObjectID, roleName string) error {
	userIDStr := userID.Hex()
	if roles, ok := m.userRoles[userIDStr]; ok {
		for i, r := range roles {
			if r == roleName {
				m.userRoles[userIDStr] = append(roles[:i], roles[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (m *MockPermissionRepository) ClearUserRoles(ctx context.Context, userID primitive.ObjectID) error {
	delete(m.userRoles, userID.Hex())
	return nil
}

// ==================== 权限检查 ====================

func (m *MockPermissionRepository) UserHasPermission(ctx context.Context, userID primitive.ObjectID, permissionCode string) (bool, error) {
	roles, _ := m.GetUserRoles(ctx, userID)
	if len(roles) == 0 {
		return false, nil
	}

	for _, roleName := range roles {
		for _, r := range m.roles {
			if r.Name == roleName {
				for _, p := range r.Permissions {
					if p == permissionCode {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

func (m *MockPermissionRepository) UserHasAnyPermission(ctx context.Context, userID primitive.ObjectID, permissionCodes []string) (bool, error) {
	for _, code := range permissionCodes {
		has, _ := m.UserHasPermission(ctx, userID, code)
		if has {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockPermissionRepository) UserHasAllPermissions(ctx context.Context, userID primitive.ObjectID, permissionCodes []string) (bool, error) {
	for _, code := range permissionCodes {
		has, _ := m.UserHasPermission(ctx, userID, code)
		if !has {
			return false, nil
		}
	}
	return true, nil
}

func (m *MockPermissionRepository) GetUserPermissions(ctx context.Context, userID primitive.ObjectID) ([]*auth.Permission, error) {
	roles, _ := m.GetUserRoles(ctx, userID)
	if len(roles) == 0 {
		return []*auth.Permission{}, nil
	}

	permissionCodes := make(map[string]bool)
	for _, roleName := range roles {
		for _, r := range m.roles {
			if r.Name == roleName {
				for _, p := range r.Permissions {
					permissionCodes[p] = true
				}
			}
		}
	}

	result := make([]*auth.Permission, 0)
	for code := range permissionCodes {
		if p := m.permissions[code]; p != nil {
			result = append(result, p)
		}
	}
	return result, nil
}

// ==================== 测试用例 ====================

// TestPermissionService_GetAllPermissions 测试获取所有权限
func TestPermissionService_GetAllPermissions(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()

	// 创建测试权限
	perm1 := &auth.Permission{
		Code:        "user.read",
		Name:        "读取用户",
		Description: "查看用户信息",
		Resource:    "user",
		Action:      "read",
	}
	perm2 := &auth.Permission{
		Code:        "user.write",
		Name:        "写入用户",
		Description: "修改用户信息",
		Resource:    "user",
		Action:      "write",
	}

	repo.CreatePermission(ctx, perm1)
	repo.CreatePermission(ctx, perm2)

	// 测试获取所有权限
	permissions, err := service.GetAllPermissions(ctx)
	if err != nil {
		t.Fatalf("GetAllPermissions failed: %v", err)
	}

	if len(permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions))
	}
}

// TestPermissionService_CreatePermission 测试创建权限
func TestPermissionService_CreatePermission(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()

	perm := &auth.Permission{
		Code:        "book.read",
		Name:        "阅读书籍",
		Description: "允许阅读书籍内容",
		Resource:    "book",
		Action:      "read",
	}

	err := service.CreatePermission(ctx, perm)
	if err != nil {
		t.Fatalf("CreatePermission failed: %v", err)
	}

	// 验证权限已创建
	retrieved, _ := service.GetPermissionByCode(ctx, "book.read")
	if retrieved == nil {
		t.Error("Permission was not created")
	}
}

// TestPermissionService_CreateRole 测试创建角色
func TestPermissionService_CreateRole(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()

	role := &auth.Role{
		Name:        "editor",
		Description: "编辑角色",
		Permissions: []string{"user.read", "user.write"},
		IsSystem:    false,
	}

	err := service.CreateRole(ctx, role)
	if err != nil {
		t.Fatalf("CreateRole failed: %v", err)
	}

	// 验证角色已创建
	if role.ID == "" {
		t.Error("Role ID was not generated")
	}
}

// TestPermissionService_AssignPermissionToRole 测试为角色分配权限
func TestPermissionService_AssignPermissionToRole(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()

	// 创建权限
	perm := &auth.Permission{
		Code:     "user.delete",
		Name:     "删除用户",
		Resource: "user",
		Action:   "delete",
	}
	repo.CreatePermission(ctx, perm)

	// 创建角色
	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{},
		IsSystem:    false,
	}
	service.CreateRole(ctx, role)

	// 分配权限
	err := service.AssignPermissionToRole(ctx, role.ID, "user.delete")
	if err != nil {
		t.Fatalf("AssignPermissionToRole failed: %v", err)
	}

	// 验证权限已分配
	updatedRole, err := service.GetRoleByID(ctx, role.ID)
	if err != nil {
		t.Fatalf("GetRoleByID failed: %v", err)
	}
	if updatedRole == nil {
		t.Fatal("Updated role is nil")
	}
	if len(updatedRole.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(updatedRole.Permissions))
	}
}

// TestPermissionService_AssignRoleToUser 测试为用户分配角色
func TestPermissionService_AssignRoleToUser(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	// 创建角色
	role := &auth.Role{
		Name:        "author",
		Description: "作者",
		Permissions: []string{"book.write"},
		IsSystem:    false,
	}
	service.CreateRole(ctx, role)

	// 为用户分配角色
	err := service.AssignRoleToUser(ctx, userID, "author")
	if err != nil {
		t.Fatalf("AssignRoleToUser failed: %v", err)
	}

	// 验证用户角色
	roles, _ := service.GetUserRoles(ctx, userID)
	if len(roles) != 1 || roles[0] != "author" {
		t.Errorf("Expected user to have 'author' role, got %v", roles)
	}
}

// TestPermissionService_UserHasPermission 测试用户权限检查
func TestPermissionService_UserHasPermission(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// 创建权限和角色
	perm := &auth.Permission{
		Code:     "book.publish",
		Name:     "发布书籍",
		Resource: "book",
		Action:   "publish",
	}
	repo.CreatePermission(ctx, perm)

	role := &auth.Role{
		Name:        "author",
		Description: "作者",
		Permissions: []string{"book.publish"},
		IsSystem:    false,
	}
	service.CreateRole(ctx, role)

	repo.AssignRoleToUser(ctx, userID, "author")

	// 测试权限检查
	has, err := service.UserHasPermission(ctx, userID.Hex(), "book.publish")
	if err != nil {
		t.Fatalf("UserHasPermission failed: %v", err)
	}
	if !has {
		t.Error("Expected user to have 'book.publish' permission")
	}

	// 测试不存在的权限
	has, _ = service.UserHasPermission(ctx, userID.Hex(), "book.delete")
	if has {
		t.Error("Expected user to NOT have 'book.delete' permission")
	}
}

// TestPermissionService_UserHasAllPermissions 测试用户拥有所有权限检查
func TestPermissionService_UserHasAllPermissions(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// 创建权限和角色
	perms := []*auth.Permission{
		{Code: "user.read", Name: "读用户", Resource: "user", Action: "read"},
		{Code: "user.write", Name: "写用户", Resource: "user", Action: "write"},
	}
	for _, p := range perms {
		repo.CreatePermission(ctx, p)
	}

	role := &auth.Role{
		Name:        "admin",
		Description: "管理员",
		Permissions: []string{"user.read", "user.write"},
		IsSystem:    false,
	}
	service.CreateRole(ctx, role)
	repo.AssignRoleToUser(ctx, userID, "admin")

	// 测试拥有所有权限
	has, err := service.UserHasAllPermissions(ctx, userID.Hex(), []string{"user.read", "user.write"})
	if err != nil {
		t.Fatalf("UserHasAllPermissions failed: %v", err)
	}
	if !has {
		t.Error("Expected user to have all permissions")
	}

	// 测试部分权限
	has, _ = service.UserHasAllPermissions(ctx, userID.Hex(), []string{"user.read", "user.delete"})
	if has {
		t.Error("Expected user to NOT have all permissions")
	}
}

// TestPermissionService_GetUserPermissions 测试获取用户权限
func TestPermissionService_GetUserPermissions(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	// 创建权限和角色
	perms := []*auth.Permission{
		{Code: "book.read", Name: "读书", Resource: "book", Action: "read"},
		{Code: "book.write", Name: "写书", Resource: "book", Action: "write"},
	}
	for _, p := range perms {
		repo.CreatePermission(ctx, p)
	}

	role := &auth.Role{
		Name:        "author",
		Description: "作者",
		Permissions: []string{"book.read", "book.write"},
		IsSystem:    false,
	}
	service.CreateRole(ctx, role)
	repo.AssignRoleToUser(ctx, userID, "author")

	// 获取用户权限
	userPerms, err := service.GetUserPermissions(ctx, userID.Hex())
	if err != nil {
		t.Fatalf("GetUserPermissions failed: %v", err)
	}

	if len(userPerms) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(userPerms))
	}
}

// TestPermissionService_DeleteRole 测试删除系统角色
func TestPermissionService_DeleteRole(t *testing.T) {
	repo := NewMockPermissionRepository()
	service := NewPermissionService(repo)

	ctx := context.Background()

	// 创建系统角色
	role := &auth.Role{
		Name:        "super_admin",
		Description: "超级管理员",
		IsSystem:    true,
	}
	service.CreateRole(ctx, role)

	// 尝试删除系统角色
	err := service.DeleteRole(ctx, role.ID)
	if err == nil {
		t.Error("Expected error when deleting system role")
	}

	// 创建非系统角色
	normalRole := &auth.Role{
		Name:        "custom_role",
		Description: "自定义角色",
		IsSystem:    false,
	}
	service.CreateRole(ctx, normalRole)

	// 删除非系统角色
	err = service.DeleteRole(ctx, normalRole.ID)
	if err != nil {
		t.Errorf("Should be able to delete non-system role: %v", err)
	}
}
