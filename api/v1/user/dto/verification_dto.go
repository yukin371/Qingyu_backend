package dto

// ===========================
// 验证相关 DTO
// ===========================

// SendEmailCodeRequest 发送邮箱验证码请求
type SendEmailCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// SendPhoneCodeRequest 发送手机验证码请求
type SendPhoneCodeRequest struct {
	Phone string `json:"phone" binding:"required" example:"13800138000"`
}

// VerifyEmailRequest 验证邮箱请求
type VerifyEmailRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Code      string `json:"code" binding:"required,len=6"`
	Timestamp int64  `json:"timestamp"`
}

// UnbindEmailRequest 解绑邮箱请求
type UnbindEmailRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// UnbindPhoneRequest 解绑手机请求
type UnbindPhoneRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	ExpiresIn int    `json:"expires_in"` // 有效期（秒）
	Message   string `json:"message"`   // 额外信息
}

// VerifyEmailResponse 验证邮箱响应
type VerifyEmailResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

// DeleteDeviceRequest 删除设备请求
type DeleteDeviceRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}
