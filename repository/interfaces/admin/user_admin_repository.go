package admin

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/users"
	userrepo "Qingyu_backend/repository/interfaces/user"
)

// UserAdminRepository 用户管理仓储接口
// 嵌入共享基础接口，保留admin特有的管理方法
type UserAdminRepository interface {
	// === 嵌入共享基础接口 ===
	// 提供: Create, GetByID(string), Update, Delete, GetByEmail, GetByUsername,
	//       UpdateStatus, UpdatePassword, SetEmailVerified, BatchUpdateStatus, BatchDelete,
	//       CountByStatus, CountByRole, List (所有方法使用string ID)
	userrepo.BaseUserRepository

	// === admin特有方法 ===

	// ListWithPagination 获取用户列表（带分页和筛选，admin特有）
	ListWithPagination(ctx context.Context, filter *UserFilter, page, pageSize int) ([]*users.User, int64, error)

	// BatchCreate 批量创建用户（admin特有）
	BatchCreate(ctx context.Context, usersList []*users.User) error

	// HardDelete 硬删除用户（admin特有）
	HardDelete(ctx context.Context, userID primitive.ObjectID) error

	// GetActivities 获取用户活动记录（admin特有）
	GetActivities(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*users.UserActivity, int64, error)

	// GetStatistics 获取用户统计信息（admin特有）
	GetStatistics(ctx context.Context, userID primitive.ObjectID) (*users.UserStatistics, error)

	// ResetPassword 管理员重置用户密码（admin特有）
	ResetPassword(ctx context.Context, userID primitive.ObjectID, newPassword string) error

	// UpdateRoles 更新用户角色（admin特有）
	UpdateRoles(ctx context.Context, userID primitive.ObjectID, role string) error

	// SearchUsers 搜索用户（admin特有）
	SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]*users.User, int64, error)

	// GetUsersByRole 根据角色获取用户列表（admin特有）
	GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]*users.User, int64, error)

	// CountByStatusMap 按状态统计用户数量（admin特有，返回map）
	// 注意：与 BaseUserRepository.CountByStatus(UserStatus) int64 不同
	CountByStatusMap(ctx context.Context) (map[string]int64, error)

	// GetRecentUsers 获取最近注册的用户（admin特有）
	GetRecentUsers(ctx context.Context, limit int) ([]*users.User, error)

	// GetActiveUsers 获取活跃用户（admin特有）
	GetActiveUsers(ctx context.Context, days int, limit int) ([]*users.User, error)
}

// 注意：以下方法已通过 BaseUserRepository 提供（使用string ID）：
// - Create(ctx, *users.User) error
// - GetByID(ctx, string) (*users.User, error)
// - GetByEmail(ctx, string) (*users.User, error)
// - GetByUsername(ctx, string) (*users.User, error)
// - Update(ctx, string, map[string]interface{}) error
// - UpdateStatus(ctx, string, UserStatus) error
// - UpdatePassword(ctx, string, string) error
// - SetEmailVerified(ctx, string, bool) error
// - Delete(ctx, string) error
// - BatchUpdateStatus(ctx, []string, UserStatus) error
// - BatchDelete(ctx, []string) error
// - CountByStatus(ctx, UserStatus) (int64, error)
// - CountByRole(ctx, string) (int64, error)
// - List(ctx, filter) ([]*users.User, error)

// UserFilter 用户筛选条件
type UserFilter struct {
	Keyword    string           // 用户名/邮箱模糊搜索
	Status     users.UserStatus // 状态筛选
	Role       string           // 角色筛选
	DateFrom   *time.Time       // 注册时间起
	DateTo     *time.Time       // 注册时间止
	LastActive *time.Time       // 最后活跃时间
	IsVIP      *bool            // 是否VIP
}
