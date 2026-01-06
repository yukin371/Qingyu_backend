package shared

import "time"

// ===========================
// 用户相关 DTO（统一定义，消除重复）
// ===========================

// UserProfileResponse 用户信息响应（完整信息，包含敏感字段）
// 用于：用户自己查看详细信息、管理员查看用户信息
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

// PublicUserProfileResponse 用户公开信息响应（去敏感字段）
// 用于：用户主页、评论作者信息、社交互动等公开场景
type PublicUserProfileResponse struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// UserBasicInfo 基本用户信息
// 用于：登录响应、注册响应、简单用户引用
type UserBasicInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role,omitempty"`
}
