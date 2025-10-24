package user

import (
	"Qingyu_backend/service/interfaces/base"
	"context"
	"time"

	usersModel "Qingyu_backend/models/users"
)

// UserService 用户服务接口
// 定义用户相关的所有业务逻辑方法
type UserService interface {
	base.BaseService

	// 用户管理
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error)
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)

	// 用户认证
	RegisterUser(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error)
	LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error)
	LogoutUser(ctx context.Context, req *LogoutUserRequest) (*LogoutUserResponse, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)

	// 用户状态管理
	UpdateLastLogin(ctx context.Context, req *UpdateLastLoginRequest) (*UpdateLastLoginResponse, error)
	UpdatePassword(ctx context.Context, req *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) (*ResetPasswordResponse, error)

	// 用户权限管理
	AssignRole(ctx context.Context, req *AssignRoleRequest) (*AssignRoleResponse, error)
	RemoveRole(ctx context.Context, req *RemoveRoleRequest) (*RemoveRoleResponse, error)
	GetUserRoles(ctx context.Context, req *GetUserRolesRequest) (*GetUserRolesResponse, error)
	GetUserPermissions(ctx context.Context, req *GetUserPermissionsRequest) (*GetUserPermissionsResponse, error)
}

// 请求和响应结构体定义

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role,omitempty"`
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
	User *usersModel.User `json:"user"`
}

// GetUserRequest 获取用户请求
type GetUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// GetUserResponse 获取用户响应
type GetUserResponse struct {
	User *usersModel.User `json:"user"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID      string                 `json:"id" validate:"required"`
	Updates map[string]interface{} `json:"updates" validate:"required"`
}

// UpdateUserResponse 更新用户响应
type UpdateUserResponse struct {
	usersModel.User
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// DeleteUserResponse 删除用户响应
type DeleteUserResponse struct {
	Deleted   bool      `json:"deleted"`
	DeletedAt time.Time `json:"deleted_at"`
}

// ListUsersRequest 列出用户请求
type ListUsersRequest struct {
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Status   string    `json:"status,omitempty"`
	FromDate time.Time `json:"from_date,omitempty"`
	ToDate   time.Time `json:"to_date,omitempty"`
	Page     int       `json:"page,omitempty"`
	PageSize int       `json:"page_size,omitempty"`
}

// ListUsersResponse 列出用户响应
type ListUsersResponse struct {
	Users      []*usersModel.User `json:"users"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// RegisterUserRequest 注册用户请求
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterUserResponse 注册用户响应
type RegisterUserResponse struct {
	User  *usersModel.User `json:"user"`
	Token string           `json:"token,omitempty"`
}

// LoginUserRequest 登录用户请求
type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginUserResponse 登录用户响应
type LoginUserResponse struct {
	User  *usersModel.User `json:"user"`
	Token string           `json:"token"`
}

// LogoutUserRequest 登出用户请求
type LogoutUserRequest struct {
	Token string `json:"token" validate:"required"`
}

// LogoutUserResponse 登出用户响应
type LogoutUserResponse struct {
	Success bool `json:"success"`
}

// ValidateTokenRequest 验证令牌请求
type ValidateTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// ValidateTokenResponse 验证令牌响应
type ValidateTokenResponse struct {
	User  *usersModel.User `json:"user,omitempty"`
	Valid bool             `json:"valid"`
}

// UpdateLastLoginRequest 更新最后登录时间请求
type UpdateLastLoginRequest struct {
	ID string `json:"id" validate:"required"`
}

// UpdateLastLoginResponse 更新最后登录时间响应
type UpdateLastLoginResponse struct {
	Updated bool `json:"updated"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	ID          string `json:"id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// UpdatePasswordResponse 更新密码响应
type UpdatePasswordResponse struct {
	Updated bool `json:"updated"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordResponse 重置密码响应
type ResetPasswordResponse struct {
	Success bool `json:"success"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	UserID string `json:"user_id" validate:"required"`
	RoleID string `json:"role_id" validate:"required"`
}

// AssignRoleResponse 分配角色响应
type AssignRoleResponse struct {
	Assigned bool `json:"assigned"`
}

// RemoveRoleRequest 移除角色请求
type RemoveRoleRequest struct {
	UserID string `json:"user_id" validate:"required"`
	RoleID string `json:"role_id" validate:"required"`
}

// RemoveRoleResponse 移除角色响应
type RemoveRoleResponse struct {
	Removed bool `json:"removed"`
}

// GetUserRolesRequest 获取用户角色请求
type GetUserRolesRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// GetUserRolesResponse 获取用户角色响应
type GetUserRolesResponse struct {
	Roles []string `json:"roles"`
}

// GetUserPermissionsRequest 获取用户权限请求
type GetUserPermissionsRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// GetUserPermissionsResponse 获取用户权限响应
type GetUserPermissionsResponse struct {
	Permissions []string `json:"permissions"`
}
