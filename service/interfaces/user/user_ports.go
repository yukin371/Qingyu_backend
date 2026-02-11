package user

import (
	"context"
	"Qingyu_backend/service/interfaces/base"
)

// UserManagementPort 用户管理端口
// 负责用户的 CRUD 操作
type UserManagementPort interface {
	base.BaseService

	// CreateUser 创建新用户
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)

	// GetUser 获取用户信息
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error)

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error)

	// ListUsers 列出用户
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
}

// UserAuthPort 用户认证端口
// 负责用户注册、登录、登出和令牌验证
type UserAuthPort interface {
	base.BaseService

	// RegisterUser 用户注册
	RegisterUser(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error)

	// LoginUser 用户登录
	LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error)

	// LogoutUser 用户登出
	LogoutUser(ctx context.Context, req *LogoutUserRequest) (*LogoutUserResponse, error)

	// ValidateToken 验证令牌
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
}

// PasswordManagementPort 密码管理端口
// 负责密码更新、重置和验证
type PasswordManagementPort interface {
	base.BaseService

	// UpdatePassword 更新密码
	UpdatePassword(ctx context.Context, req *UpdatePasswordRequest) (*UpdatePasswordResponse, error)

	// ResetPassword 重置密码（保持兼容性）
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) (*ResetPasswordResponse, error)

	// RequestPasswordReset 请求密码重置
	RequestPasswordReset(ctx context.Context, req *RequestPasswordResetRequest) (*RequestPasswordResetResponse, error)

	// ConfirmPasswordReset 确认密码重置
	ConfirmPasswordReset(ctx context.Context, req *ConfirmPasswordResetRequest) (*ConfirmPasswordResetResponse, error)

	// VerifyPassword 验证密码
	VerifyPassword(ctx context.Context, userID string, password string) error
}

// EmailManagementPort 邮箱管理端口
// 负责邮箱验证和管理
type EmailManagementPort interface {
	base.BaseService

	// SendEmailVerification 发送邮箱验证
	SendEmailVerification(ctx context.Context, req *SendEmailVerificationRequest) (*SendEmailVerificationResponse, error)

	// VerifyEmail 验证邮箱
	VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*VerifyEmailResponse, error)

	// UnbindEmail 解绑邮箱
	UnbindEmail(ctx context.Context, userID string) error

	// EmailExists 检查邮箱是否已存在
	EmailExists(ctx context.Context, email string) (bool, error)
}

// UserPermissionPort 用户权限端口
// 负责角色和权限管理
type UserPermissionPort interface {
	base.BaseService

	// AssignRole 分配角色
	AssignRole(ctx context.Context, req *AssignRoleRequest) (*AssignRoleResponse, error)

	// RemoveRole 移除角色
	RemoveRole(ctx context.Context, req *RemoveRoleRequest) (*RemoveRoleResponse, error)

	// GetUserRoles 获取用户角色
	GetUserRoles(ctx context.Context, req *GetUserRolesRequest) (*GetUserRolesResponse, error)

	// GetUserPermissions 获取用户权限
	GetUserPermissions(ctx context.Context, req *GetUserPermissionsRequest) (*GetUserPermissionsResponse, error)

	// DowngradeRole 角色降级
	DowngradeRole(ctx context.Context, req *DowngradeRoleRequest) (*DowngradeRoleResponse, error)
}

// UserStatusPort 用户状态端口
// 负责用户状态管理
type UserStatusPort interface {
	base.BaseService

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, req *UpdateLastLoginRequest) (*UpdateLastLoginResponse, error)

	// DeleteDevice 删除设备
	DeleteDevice(ctx context.Context, userID string, deviceID string) error

	// UnbindPhone 解绑手机
	UnbindPhone(ctx context.Context, userID string) error
}
