package auth

import (
	authModel "Qingyu_backend/models/auth"
	"context"
)

// RoleRepository 角色管理Repository
type RoleRepository interface {
	// 角色管理
	CreateRole(ctx context.Context, role *authModel.Role) error
	GetRole(ctx context.Context, roleID string) (*authModel.Role, error)
	GetRoleByName(ctx context.Context, name string) (*authModel.Role, error)
	UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error
	DeleteRole(ctx context.Context, roleID string) error
	ListRoles(ctx context.Context) ([]*authModel.Role, error)

	// 用户角色关联
	AssignUserRole(ctx context.Context, userID, roleID string) error
	RemoveUserRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error)
	HasUserRole(ctx context.Context, userID, roleID string) (bool, error)

	// 权限查询
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
