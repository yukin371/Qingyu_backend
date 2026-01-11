package dto

import "Qingyu_backend/api/v1/shared"

// ===========================
// 用户相关 DTO
// ===========================

// UserBasicInfo 基本用户信息（引用共享定义）
type UserBasicInfo = shared.UserBasicInfo

// UserProfileResponse 用户信息响应（引用共享定义）
type UserProfileResponse = shared.UserProfileResponse

// PublicUserProfileResponse 用户公开信息响应（引用共享定义）
type PublicUserProfileResponse = shared.PublicUserProfileResponse

// ===========================
// 请求 DTO
// ===========================

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

// GetUserBooksResponse 用户作品列表响应
type GetUserBooksResponse struct {
	Books interface{} `json:"books"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
