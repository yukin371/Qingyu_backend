package shared

import (
	authModel "Qingyu_backend/models/auth"
	"context"

	"Qingyu_backend/repository/interfaces/infrastructure"
	userRepo "Qingyu_backend/repository/interfaces/user"

	"go.mongodb.org/mongo-driver/mongo"
)

// RoleRepositoryAdapter RoleRepository 适配器
// 将 AuthRepository 适配为 RoleRepository 接口
type RoleRepositoryAdapter struct {
	authRepo *AuthRepositoryImpl
}

// NewRoleRepository 创建 RoleRepository
func NewRoleRepository(db *mongo.Database) userRepo.RoleRepository {
	return &RoleRepositoryAdapter{
		authRepo: &AuthRepositoryImpl{
			db:             db,
			roleCollection: db.Collection("roles"),
		},
	}
}

// ============ CRUD 接口实现 ============

// Create 创建角色
func (r *RoleRepositoryAdapter) Create(ctx context.Context, entity *authModel.Role) error {
	return r.authRepo.CreateRole(ctx, entity)
}

// GetByID 根据ID获取角色
func (r *RoleRepositoryAdapter) GetByID(ctx context.Context, id string) (*authModel.Role, error) {
	return r.authRepo.GetRole(ctx, id)
}

// Update 更新角色
func (r *RoleRepositoryAdapter) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.authRepo.UpdateRole(ctx, id, updates)
}

// Delete 删除角色
func (r *RoleRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return r.authRepo.DeleteRole(ctx, id)
}

// List 列出所有角色
func (r *RoleRepositoryAdapter) List(ctx context.Context, filter infrastructure.Filter) ([]*authModel.Role, error) {
	// AuthRepository 的 ListRoles 不支持过滤，返回所有角色
	return r.authRepo.ListRoles(ctx)
}

// Count 统计角色数量
func (r *RoleRepositoryAdapter) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	roles, err := r.authRepo.ListRoles(ctx)
	if err != nil {
		return 0, err
	}
	return int64(len(roles)), nil
}

// Exists 检查角色是否存在
func (r *RoleRepositoryAdapter) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.authRepo.GetRole(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ============ 角色查询方法 ============

// GetByName 根据名称获取角色
func (r *RoleRepositoryAdapter) GetByName(ctx context.Context, name string) (*authModel.Role, error) {
	return r.authRepo.GetRoleByName(ctx, name)
}

// GetDefaultRole 获取默认角色
func (r *RoleRepositoryAdapter) GetDefaultRole(ctx context.Context) (*authModel.Role, error) {
	// 获取所有角色并找到默认角色
	roles, err := r.authRepo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.IsDefault {
			return role, nil
		}
	}

	// 如果没有默认角色，返回 reader 角色
	return r.authRepo.GetRoleByName(ctx, authModel.RoleReader)
}

// ExistsByName 检查角色名称是否存在
func (r *RoleRepositoryAdapter) ExistsByName(ctx context.Context, name string) (bool, error) {
	_, err := r.authRepo.GetRoleByName(ctx, name)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ============ 角色列表 ============

// ListAllRoles 列出所有角色
func (r *RoleRepositoryAdapter) ListAllRoles(ctx context.Context) ([]*authModel.Role, error) {
	return r.authRepo.ListRoles(ctx)
}

// ListDefaultRoles 列出所有默认角色
func (r *RoleRepositoryAdapter) ListDefaultRoles(ctx context.Context) ([]*authModel.Role, error) {
	roles, err := r.authRepo.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	var defaultRoles []*authModel.Role
	for _, role := range roles {
		if role.IsDefault {
			defaultRoles = append(defaultRoles, role)
		}
	}

	return defaultRoles, nil
}

// ============ 权限管理 ============

// GetRolePermissions 获取角色权限
func (r *RoleRepositoryAdapter) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	return r.authRepo.GetRolePermissions(ctx, roleID)
}

// UpdateRolePermissions 更新角色权限
func (r *RoleRepositoryAdapter) UpdateRolePermissions(ctx context.Context, roleID string, permissions []string) error {
	updates := map[string]interface{}{
		"permissions": permissions,
	}
	return r.authRepo.UpdateRole(ctx, roleID, updates)
}

// AddPermission 添加权限
func (r *RoleRepositoryAdapter) AddPermission(ctx context.Context, roleID string, permission string) error {
	// 获取当前权限
	permissions, err := r.authRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return err
	}

	// 检查权限是否已存在
	for _, p := range permissions {
		if p == permission {
			return nil // 权限已存在，直接返回
		}
	}

	// 添加新权限
	permissions = append(permissions, permission)
	return r.UpdateRolePermissions(ctx, roleID, permissions)
}

// RemovePermission 移除权限
func (r *RoleRepositoryAdapter) RemovePermission(ctx context.Context, roleID string, permission string) error {
	// 获取当前权限
	permissions, err := r.authRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return err
	}

	// 移除权限
	var newPermissions []string
	for _, p := range permissions {
		if p != permission {
			newPermissions = append(newPermissions, p)
		}
	}

	return r.UpdateRolePermissions(ctx, roleID, newPermissions)
}

// ============ 统计 ============

// CountByName 按名称统计角色数量
func (r *RoleRepositoryAdapter) CountByName(ctx context.Context, name string) (int64, error) {
	exists, err := r.ExistsByName(ctx, name)
	if err != nil {
		return 0, err
	}
	if exists {
		return 1, nil
	}
	return 0, nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *RoleRepositoryAdapter) Health(ctx context.Context) error {
	return r.authRepo.Health(ctx)
}
