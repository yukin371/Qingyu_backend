package dto

// ===========================
// 密码相关 DTO
// ===========================

// SendPasswordResetRequest 发送密码重置验证码请求
type SendPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" validate:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8" validate:"required,min=8"`
}

// SendPasswordResetResponse 发送密码重置验证码响应
type SendPasswordResetResponse struct {
	ExpiresIn int    `json:"expires_in"`
	Message   string `json:"message"`
}

// ResetPasswordResponse 重置密码响应
type ResetPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
