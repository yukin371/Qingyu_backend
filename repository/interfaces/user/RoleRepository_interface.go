package user

import (
	usersModel "Qingyu_backend/models/users"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	base.CRUDRepository[*usersModel.Role, interface{}]

	// 角色特定方法
	GetByName(ctx context.Context, name string) (interface{}, error)
	GetDefaultRole(ctx context.Context) (interface{}, error)
	GetUserRoles(ctx context.Context, userID string) ([]interface{}, error)
	AssignRole(ctx context.Context, userID, roleID string) error
	RemoveRole(ctx context.Context, userID, roleID string) error
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
}
