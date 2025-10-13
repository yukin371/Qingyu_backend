package user

import (
	usersModel "Qingyu_backend/models/users"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*usersModel.Role, string]
	// 继承健康检查接口
	base.HealthRepository

	// 角色查询方法
	GetByName(ctx context.Context, name string) (*usersModel.Role, error)
	GetDefaultRole(ctx context.Context) (*usersModel.Role, error)
	ExistsByName(ctx context.Context, name string) (bool, error)

	// 角色列表
	ListAllRoles(ctx context.Context) ([]*usersModel.Role, error)
	ListDefaultRoles(ctx context.Context) ([]*usersModel.Role, error)

	// 权限管理
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	UpdateRolePermissions(ctx context.Context, roleID string, permissions []string) error
	AddPermission(ctx context.Context, roleID string, permission string) error
	RemovePermission(ctx context.Context, roleID string, permission string) error

	// 统计
	CountByName(ctx context.Context, name string) (int64, error)
}
