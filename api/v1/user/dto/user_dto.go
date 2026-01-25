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
	Gender   *string `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Birthday *string `json:"birthday,omitempty"` // RFC3339格式：2006-01-02T15:04:05Z07:00，在handler中验证
	Location *string `json:"location,omitempty" validate:"omitempty,max=100"`
	Website  *string `json:"website,omitempty" validate:"omitempty,url"`
}


// GetUserBooksResponse 用户作品列表响应
type GetUserBooksResponse struct {
	Books interface{} `json:"books"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// ===========================
// 头像上传相关 DTO
// ===========================

// UploadAvatarResponse 上传头像响应
type UploadAvatarResponse struct {
	AvatarURL string `json:"avatar_url"`
	Message   string `json:"message"`
}
