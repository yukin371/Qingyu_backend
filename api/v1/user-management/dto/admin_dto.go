package dto

// ===========================
// 管理员用户管理相关 DTO
// ===========================

// ListUsersRequest 获取用户列表请求（查询参数）
type ListUsersRequest struct {
	Page     int    `form:"page" validate:"omitempty,min=1"`
	PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100"`
	Username string `form:"username" validate:"omitempty"`
	Email    string `form:"email" validate:"omitempty,email"`
	Role     string `form:"role" validate:"omitempty"`
	Status   string `form:"status" validate:"omitempty"`
}

// AdminUpdateUserRequest 管理员更新用户请求
type AdminUpdateUserRequest struct {
	Nickname      *string `json:"nickname,omitempty" validate:"omitempty,max=50"`
	Bio           *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	Avatar        *string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone         *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Role          *string `json:"role,omitempty" validate:"omitempty,oneof=user author admin"`
	Status        *string `json:"status,omitempty" validate:"omitempty"`
	EmailVerified *bool   `json:"email_verified,omitempty"`
	PhoneVerified *bool   `json:"phone_verified,omitempty"`
}

// BanUserRequest 封禁用户请求
type BanUserRequest struct {
	Reason       string     `json:"reason" binding:"required"`
	Duration     int        `json:"duration" binding:"required,min=1"`
	DurationUnit string     `json:"durationUnit" binding:"required,oneof=hours days weeks months"`
	BanUntil     *string    `json:"banUntil,omitempty"`
}
