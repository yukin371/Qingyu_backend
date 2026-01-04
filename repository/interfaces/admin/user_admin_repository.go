package admin

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/users"
)

// UserAdminRepository 用户管理仓储接口
type UserAdminRepository interface {
	// List 获取用户列表（分页、筛选）
	List(ctx context.Context, filter *UserFilter, page, pageSize int) ([]*users.User, int64, error)

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, userID primitive.ObjectID) (*users.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*users.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, userID primitive.ObjectID, user *users.User) error

	// UpdateStatus 更新用户状态
	UpdateStatus(ctx context.Context, userID primitive.ObjectID, status users.UserStatus) error

	// Delete 删除用户（软删除）
	Delete(ctx context.Context, userID primitive.ObjectID) error

	// HardDelete 硬删除用户
	HardDelete(ctx context.Context, userID primitive.ObjectID) error

	// GetActivities 获取用户活动记录
	GetActivities(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*users.UserActivity, int64, error)

	// UpdateRoles 更新用户角色
	UpdateRoles(ctx context.Context, userID primitive.ObjectID, role string) error

	// GetStatistics 获取用户统计信息
	GetStatistics(ctx context.Context, userID primitive.ObjectID) (*users.UserStatistics, error)

	// ResetPassword 重置用户密码
	ResetPassword(ctx context.Context, userID primitive.ObjectID, newPassword string) error

	// BatchUpdateStatus 批量更新用户状态
	BatchUpdateStatus(ctx context.Context, userIDs []primitive.ObjectID, status users.UserStatus) error

	// BatchDelete 批量删除用户
	BatchDelete(ctx context.Context, userIDs []primitive.ObjectID) error

	// GetUsersByRole 根据角色获取用户列表
	GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]*users.User, int64, error)

	// SearchUsers 搜索用户
	SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]*users.User, int64, error)

	// CountByStatus 按状态统计用户数量
	CountByStatus(ctx context.Context) (map[string]int64, error)

	// GetRecentUsers 获取最近注册的用户
	GetRecentUsers(ctx context.Context, limit int) ([]*users.User, error)

	// GetActiveUsers 获取活跃用户
	GetActiveUsers(ctx context.Context, days int, limit int) ([]*users.User, error)
}

// UserFilter 用户筛选条件
type UserFilter struct {
	Keyword    string               // 用户名/邮箱模糊搜索
	Status     users.UserStatus     // 状态筛选
	Role       string               // 角色筛选
	DateFrom   *time.Time           // 注册时间起
	DateTo     *time.Time           // 注册时间止
	LastActive *time.Time           // 最后活跃时间
	IsVIP      *bool                // 是否VIP
}
