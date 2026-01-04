package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/auth"
)

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
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
	GetUserRoles(ctx context.Context, userID primitive.ObjectID) ([]string, error)

	// AssignRoleToUser 为用户分配角色
	AssignRoleToUser(ctx context.Context, userID primitive.ObjectID, roleName string) error

	// RemoveRoleFromUser 移除用户角色
	RemoveRoleFromUser(ctx context.Context, userID primitive.ObjectID, roleName string) error

	// ClearUserRoles 清除用户所有角色
	ClearUserRoles(ctx context.Context, userID primitive.ObjectID) error

	// ==================== 权限检查 ====================

	// UserHasPermission 检查用户是否有指定权限
	UserHasPermission(ctx context.Context, userID primitive.ObjectID, permissionCode string) (bool, error)

	// UserHasAnyPermission 检查用户是否有任意一个权限
	UserHasAnyPermission(ctx context.Context, userID primitive.ObjectID, permissionCodes []string) (bool, error)

	// UserHasAllPermissions 检查用户是否拥有所有权限
	UserHasAllPermissions(ctx context.Context, userID primitive.ObjectID, permissionCodes []string) (bool, error)

	// GetUserPermissions 获取用户的所有权限
	GetUserPermissions(ctx context.Context, userID primitive.ObjectID) ([]*auth.Permission, error)
}
