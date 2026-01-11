package dto

// ===========================
// 认证相关 DTO
// ===========================

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
	Token string `json:"token"`
	User  UserBasicInfo `json:"user"`
}
