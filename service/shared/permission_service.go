package shared

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/auth"
	authrepo "Qingyu_backend/repository/interfaces/auth"
)

var (
	// ErrPermissionNotFound 权限不存在
	ErrPermissionNotFound = errors.New("permission not found")
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")
	// ErrRoleAlreadyExists 角色已存在
	ErrRoleAlreadyExists = errors.New("role already exists")
	// ErrCannotDeleteSystemRole 不能删除系统角色
	ErrCannotDeleteSystemRole = errors.New("cannot delete system role")
)

// PermissionService 权限服务接口
type PermissionService interface {
	// ==================== 权限管理 ====================

	// GetAllPermissions 获取所有权限
	GetAllPermissions(ctx context.Context) ([]*auth.Permission, error)

	// GetPermissionByCode 根据代码获取权限
	GetPermissionByCode(ctx context.Context, code string) (*auth.Permission, error)

	// CreatePermission 创建权限
	CreatePermission(ctx context.Context, permission *auth.Permission) error

	// UpdatePermission 更新权限
	UpdatePermission(ctx context.Context, permission *auth.Permission) error

	// DeletePermission 删除权限
	DeletePermission(ctx context.Context, code string) error

	// ==================== 角色管理 ====================

	// GetAllRoles 获取所有角色
	GetAllRoles(ctx context.Context) ([]*auth.Role, error)

	// GetRoleByID 根据ID获取角色
	GetRoleByID(ctx context.Context, roleID string) (*auth.Role, error)

	// GetRoleByName 根据名称获取角色
	GetRoleByName(ctx context.Context, name string) (*auth.Role, error)

	// CreateRole 创建角色
	CreateRole(ctx context.Context, role *auth.Role) error

	// UpdateRole 更新角色
	UpdateRole(ctx context.Context, role *auth.Role) error

	// DeleteRole 删除角色
	DeleteRole(ctx context.Context, roleID string) error

	// AssignPermissionToRole 为角色分配权限
	AssignPermissionToRole(ctx context.Context, roleID, permissionCode string) error

	// RemovePermissionFromRole 移除角色权限
	RemovePermissionFromRole(ctx context.Context, roleID, permissionCode string) error

	// GetRolePermissions 获取角色的所有权限
	GetRolePermissions(ctx context.Context, roleID string) ([]*auth.Permission, error)

	// ==================== 用户角色管理 ====================

	// GetUserRoles 获取用户的所有角色
	GetUserRoles(ctx context.Context, userID string) ([]string, error)

	// AssignRoleToUser 为用户分配角色
	AssignRoleToUser(ctx context.Context, userID, roleName string) error

	// RemoveRoleFromUser 移除用户角色
	RemoveRoleFromUser(ctx context.Context, userID, roleName string) error

	// ClearUserRoles 清除用户所有角色
	ClearUserRoles(ctx context.Context, userID string) error

	// ==================== 权限检查 ====================

	// UserHasPermission 检查用户是否有指定权限
	UserHasPermission(ctx context.Context, userID, permissionCode string) (bool, error)

	// UserHasAnyPermission 检查用户是否有任意一个权限
	UserHasAnyPermission(ctx context.Context, userID string, permissionCodes []string) (bool, error)

	// UserHasAllPermissions 检查用户是否拥有所有权限
	UserHasAllPermissions(ctx context.Context, userID string, permissionCodes []string) (bool, error)

	// GetUserPermissions 获取用户的所有权限
	GetUserPermissions(ctx context.Context, userID string) ([]*auth.Permission, error)

	// CheckPermission 检查权限（资源+操作）
	CheckPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

// PermissionServiceImpl 权限服务实现
type PermissionServiceImpl struct {
	permissionRepo authrepo.PermissionRepository
}

// NewPermissionService 创建权限服务
func NewPermissionService(permissionRepo authrepo.PermissionRepository) PermissionService {
	return &PermissionServiceImpl{
		permissionRepo: permissionRepo,
	}
}

// ==================== 权限管理 ====================

func (s *PermissionServiceImpl) GetAllPermissions(ctx context.Context) ([]*auth.Permission, error) {
	return s.permissionRepo.GetAllPermissions(ctx)
}

func (s *PermissionServiceImpl) GetPermissionByCode(ctx context.Context, code string) (*auth.Permission, error) {
	permission, err := s.permissionRepo.GetPermissionByCode(ctx, code)
	if err != nil {
		return nil, ErrPermissionNotFound
	}
	return permission, nil
}

func (s *PermissionServiceImpl) CreatePermission(ctx context.Context, permission *auth.Permission) error {
	// 检查权限是否已存在
	_, err := s.permissionRepo.GetPermissionByCode(ctx, permission.Code)
	if err == nil {
		return errors.New("permission already exists")
	}

	return s.permissionRepo.CreatePermission(ctx, permission)
}

func (s *PermissionServiceImpl) UpdatePermission(ctx context.Context, permission *auth.Permission) error {
	// 检查权限是否存在
	_, err := s.permissionRepo.GetPermissionByCode(ctx, permission.Code)
	if err != nil {
		return ErrPermissionNotFound
	}

	return s.permissionRepo.UpdatePermission(ctx, permission)
}

func (s *PermissionServiceImpl) DeletePermission(ctx context.Context, code string) error {
	return s.permissionRepo.DeletePermission(ctx, code)
}

// ==================== 角色管理 ====================

func (s *PermissionServiceImpl) GetAllRoles(ctx context.Context) ([]*auth.Role, error) {
	return s.permissionRepo.GetAllRoles(ctx)
}

func (s *PermissionServiceImpl) GetRoleByID(ctx context.Context, roleID string) (*auth.Role, error) {
	role, err := s.permissionRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *PermissionServiceImpl) GetRoleByName(ctx context.Context, name string) (*auth.Role, error) {
	role, err := s.permissionRepo.GetRoleByName(ctx, name)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *PermissionServiceImpl) CreateRole(ctx context.Context, role *auth.Role) error {
	// 检查角色是否已存在
	_, err := s.permissionRepo.GetRoleByName(ctx, role.Name)
	if err == nil {
		return ErrRoleAlreadyExists
	}

	return s.permissionRepo.CreateRole(ctx, role)
}

func (s *PermissionServiceImpl) UpdateRole(ctx context.Context, role *auth.Role) error {
	// 检查角色是否存在
	existingRole, err := s.permissionRepo.GetRoleByID(ctx, role.ID)
	if err != nil {
		return ErrRoleNotFound
	}

	// 系统角色不允许修改名称
	if existingRole.IsSystem && existingRole.Name != role.Name {
		return errors.New("cannot change system role name")
	}

	return s.permissionRepo.UpdateRole(ctx, role)
}

func (s *PermissionServiceImpl) DeleteRole(ctx context.Context, roleID string) error {
	role, err := s.permissionRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return ErrRoleNotFound
	}

	if role.IsSystem {
		return ErrCannotDeleteSystemRole
	}

	return s.permissionRepo.DeleteRole(ctx, roleID)
}

func (s *PermissionServiceImpl) AssignPermissionToRole(ctx context.Context, roleID, permissionCode string) error {
	// 检查角色是否存在
	_, err := s.permissionRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return ErrRoleNotFound
	}

	// 检查权限是否存在
	_, err = s.permissionRepo.GetPermissionByCode(ctx, permissionCode)
	if err != nil {
		return ErrPermissionNotFound
	}

	return s.permissionRepo.AssignPermissionToRole(ctx, roleID, permissionCode)
}

func (s *PermissionServiceImpl) RemovePermissionFromRole(ctx context.Context, roleID, permissionCode string) error {
	return s.permissionRepo.RemovePermissionFromRole(ctx, roleID, permissionCode)
}

func (s *PermissionServiceImpl) GetRolePermissions(ctx context.Context, roleID string) ([]*auth.Permission, error) {
	return s.permissionRepo.GetRolePermissions(ctx, roleID)
}

// ==================== 用户角色管理 ====================

func (s *PermissionServiceImpl) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.permissionRepo.GetUserRoles(ctx, oid)
}

func (s *PermissionServiceImpl) AssignRoleToUser(ctx context.Context, userID, roleName string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// 检查角色是否存在
	_, err = s.permissionRepo.GetRoleByName(ctx, roleName)
	if err != nil {
		return ErrRoleNotFound
	}

	return s.permissionRepo.AssignRoleToUser(ctx, oid, roleName)
}

func (s *PermissionServiceImpl) RemoveRoleFromUser(ctx context.Context, userID, roleName string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.permissionRepo.RemoveRoleFromUser(ctx, oid, roleName)
}

func (s *PermissionServiceImpl) ClearUserRoles(ctx context.Context, userID string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.permissionRepo.ClearUserRoles(ctx, oid)
}

// ==================== 权限检查 ====================

func (s *PermissionServiceImpl) UserHasPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid user ID")
	}

	return s.permissionRepo.UserHasPermission(ctx, oid, permissionCode)
}

func (s *PermissionServiceImpl) UserHasAnyPermission(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid user ID")
	}

	return s.permissionRepo.UserHasAnyPermission(ctx, oid, permissionCodes)
}

func (s *PermissionServiceImpl) UserHasAllPermissions(ctx context.Context, userID string, permissionCodes []string) (bool, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New("invalid user ID")
	}

	return s.permissionRepo.UserHasAllPermissions(ctx, oid, permissionCodes)
}

func (s *PermissionServiceImpl) GetUserPermissions(ctx context.Context, userID string) ([]*auth.Permission, error) {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return s.permissionRepo.GetUserPermissions(ctx, oid)
}

func (s *PermissionServiceImpl) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	// 构建权限代码
	permissionCode := resource + "." + action
	return s.UserHasPermission(ctx, userID, permissionCode)
}
