package user

import (
	"Qingyu_backend/api/v1/shared"
)

// DTO (Data Transfer Object) - 用户API层的请求和响应结构

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
	Role     string `json:"role"`
	Status   string `json:"status"`
	Token    string `json:"token"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string              `json:"token"`
	User  shared.UserBasicInfo `json:"user"`
}

// ===========================
// 用户相关 DTO（引用共享定义，保持向后兼容）
// ===========================

// UserProfileResponse 用户信息响应（引用共享定义）
type UserProfileResponse = shared.UserProfileResponse

// PublicUserProfileResponse 用户公开信息响应（引用共享定义）
type PublicUserProfileResponse = shared.PublicUserProfileResponse

// UserBasicInfo 基本用户信息（引用共享定义）
type UserBasicInfo = shared.UserBasicInfo

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

// UserBooksResponse 用户作品列表响应
type UserBooksResponse struct {
	Books interface{} `json:"books"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
