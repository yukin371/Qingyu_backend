package system

import (
	"time"

	usersModel "Qingyu_backend/models/users"
)

// DTO (Data Transfer Object) - API层的请求和响应结构

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" validate:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// UserProfileResponse 用户信息响应
type UserProfileResponse struct {
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone,omitempty"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	Avatar        string    `json:"avatar,omitempty"`
	Nickname      string    `json:"nickname,omitempty"`
	Bio           string    `json:"bio,omitempty"`
	EmailVerified bool      `json:"email_verified"`
	PhoneVerified bool      `json:"phone_verified"`
	LastLoginAt   time.Time `json:"last_login_at,omitempty"`
	LastLoginIP   string    `json:"last_login_ip,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname *string `json:"nickname,omitempty" validate:"omitempty,max=50"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	Avatar   *string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,e164"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" validate:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6" validate:"required,min=6"`
}

// ListUsersRequest 获取用户列表请求（查询参数）
type ListUsersRequest struct {
	Page     int                   `form:"page" validate:"omitempty,min=1"`
	PageSize int                   `form:"page_size" validate:"omitempty,min=1,max=100"`
	Username string                `form:"username" validate:"omitempty"`
	Email    string                `form:"email" validate:"omitempty,email"`
	Role     string                `form:"role" validate:"omitempty"`
	Status   usersModel.UserStatus `form:"status" validate:"omitempty"`
}

// AdminUpdateUserRequest 管理员更新用户请求
type AdminUpdateUserRequest struct {
	Nickname      *string                `json:"nickname,omitempty" validate:"omitempty,max=50"`
	Bio           *string                `json:"bio,omitempty" validate:"omitempty,max=500"`
	Avatar        *string                `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone         *string                `json:"phone,omitempty" validate:"omitempty,e164"`
	Role          *string                `json:"role,omitempty" validate:"omitempty,oneof=user author admin"`
	Status        *usersModel.UserStatus `json:"status,omitempty" validate:"omitempty"`
	EmailVerified *bool                  `json:"email_verified,omitempty"`
	PhoneVerified *bool                  `json:"phone_verified,omitempty"`
}
