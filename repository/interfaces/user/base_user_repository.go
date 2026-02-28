package user

import (
	"context"

	"Qingyu_backend/models/users"
	base "Qingyu_backend/repository/interfaces/infrastructure"
)

// BaseUserRepository 用户仓储共享基础接口
// 定义admin和user模块共用的用户操作方法
// 统一ID类型为string，ID转换封装在Repository实现层
type BaseUserRepository interface {
	// 继承通用CRUD - 所有Repository都需要
	base.CRUDRepository[*users.User, string]

	// === 用户身份查询 ===
	GetByEmail(ctx context.Context, email string) (*users.User, error)
	GetByUsername(ctx context.Context, username string) (*users.User, error)

	// === 状态管理 ===
	UpdateStatus(ctx context.Context, id string, status users.UserStatus) error
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error

	// === 验证状态 ===
	SetEmailVerified(ctx context.Context, id string, verified bool) error

	// === 批量操作 ===
	BatchUpdateStatus(ctx context.Context, ids []string, status users.UserStatus) error
	BatchDelete(ctx context.Context, ids []string) error

	// === 统计查询 ===
	CountByStatus(ctx context.Context, status users.UserStatus) (int64, error)
	CountByRole(ctx context.Context, role string) (int64, error)
}
