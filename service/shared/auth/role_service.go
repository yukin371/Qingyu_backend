package auth

import (
	"context"
	"fmt"

	authModel "Qingyu_backend/models/shared/auth"
	sharedRepo "Qingyu_backend/repository/interfaces/shared"
)

// RoleServiceImpl 角色服务实现
type RoleServiceImpl struct {
	authRepo sharedRepo.AuthRepository
}

// NewRoleService 创建角色服务
func NewRoleService(authRepo sharedRepo.AuthRepository) RoleService {
	return &RoleServiceImpl{
		authRepo: authRepo,
	}
}

// ============ 角色CRUD ============

// CreateRole 创建角色
func (s *RoleServiceImpl) CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error) {
	// 1. 验证请求
	if req.Name == "" {
		return nil, fmt.Errorf("角色名称不能为空")
	}

	// 2. 检查角色是否已存在
	_, err := s.authRepo.GetRoleByName(ctx, req.Name)
	if err == nil {
		return nil, fmt.Errorf("角色已存在: %s", req.Name)
	}

	// 3. 创建角色
	role := &authModel.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		IsSystem:    false, // 非系统角色
	}

	if err := s.authRepo.CreateRole(ctx, role); err != nil {
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	// 4. 转换为响应格式
	return convertToRoleResponse(role), nil
}

// GetRole 获取角色
func (s *RoleServiceImpl) GetRole(ctx context.Context, roleID string) (*Role, error) {
	role, err := s.authRepo.GetRole(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}

	return convertToRoleResponse(role), nil
}

// UpdateRole 更新角色
func (s *RoleServiceImpl) UpdateRole(ctx context.Context, roleID string, req *UpdateRoleRequest) error {
	// 1. 检查角色是否存在
	role, err := s.authRepo.GetRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 2. 检查是否是系统角色
	if role.IsSystem {
		return fmt.Errorf("不能修改系统角色: %s", role.Name)
	}

	// 3. 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Permissions != nil {
		updates["permissions"] = req.Permissions
	}

	if len(updates) == 0 {
		return fmt.Errorf("没有要更新的内容")
	}

	// 4. 更新角色
	if err := s.authRepo.UpdateRole(ctx, roleID, updates); err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}

	return nil
}

// DeleteRole 删除角色
func (s *RoleServiceImpl) DeleteRole(ctx context.Context, roleID string) error {
	// 删除角色（Repository会检查是否是系统角色）
	if err := s.authRepo.DeleteRole(ctx, roleID); err != nil {
		return fmt.Errorf("删除角色失败: %w", err)
	}

	return nil
}

// ListRoles 列出所有角色
func (s *RoleServiceImpl) ListRoles(ctx context.Context) ([]*Role, error) {
	roles, err := s.authRepo.ListRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取角色列表失败: %w", err)
	}

	// 转换为响应格式
	result := make([]*Role, len(roles))
	for i, role := range roles {
		result[i] = convertToRoleResponse(role)
	}

	return result, nil
}

// ============ 权限管理 ============

// AssignPermissions 分配权限
func (s *RoleServiceImpl) AssignPermissions(ctx context.Context, roleID string, permissions []string) error {
	// 1. 获取角色
	role, err := s.authRepo.GetRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 2. 检查是否是系统角色
	if role.IsSystem {
		return fmt.Errorf("不能修改系统角色权限: %s", role.Name)
	}

	// 3. 合并权限（去重）
	permMap := make(map[string]bool)
	for _, p := range role.Permissions {
		permMap[p] = true
	}
	for _, p := range permissions {
		permMap[p] = true
	}

	// 4. 转换为数组
	newPermissions := make([]string, 0, len(permMap))
	for p := range permMap {
		newPermissions = append(newPermissions, p)
	}

	// 5. 更新角色
	updates := map[string]interface{}{
		"permissions": newPermissions,
	}

	if err := s.authRepo.UpdateRole(ctx, roleID, updates); err != nil {
		return fmt.Errorf("分配权限失败: %w", err)
	}

	return nil
}

// RemovePermissions 移除权限
func (s *RoleServiceImpl) RemovePermissions(ctx context.Context, roleID string, permissions []string) error {
	// 1. 获取角色
	role, err := s.authRepo.GetRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 2. 检查是否是系统角色
	if role.IsSystem {
		return fmt.Errorf("不能修改系统角色权限: %s", role.Name)
	}

	// 3. 移除指定权限
	removeMap := make(map[string]bool)
	for _, p := range permissions {
		removeMap[p] = true
	}

	newPermissions := make([]string, 0)
	for _, p := range role.Permissions {
		if !removeMap[p] {
			newPermissions = append(newPermissions, p)
		}
	}

	// 4. 更新角色
	updates := map[string]interface{}{
		"permissions": newPermissions,
	}

	if err := s.authRepo.UpdateRole(ctx, roleID, updates); err != nil {
		return fmt.Errorf("移除权限失败: %w", err)
	}

	return nil
}

// ============ 辅助函数 ============

// convertToRoleResponse 转换为响应格式
func convertToRoleResponse(role *authModel.Role) *Role {
	return &Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
