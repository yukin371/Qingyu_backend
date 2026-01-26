package dto

// ===========================
// 用户 DTO（符合分层架构规范）
// ===========================

// UserDTO 用户数据传输对象
// 用于：Service 层和 API 层数据传输，ID 和时间字段使用字符串类型
type UserDTO struct {
	ID        string `json:"id" validate:"required"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// 基本信息
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,e164"`

	// 角色和权限
	Roles    []string `json:"roles" validate:"required,dive,oneof=reader author admin"`
	VIPLevel int      `json:"vipLevel" validate:"min=0,max=5"`

	// 状态和资料
	Status   string `json:"status" validate:"required,oneof=active inactive banned deleted"`
	Avatar   string `json:"avatar,omitempty"`
	Nickname string `json:"nickname,omitempty" validate:"max=50"`
	Bio      string `json:"bio,omitempty" validate:"max=500"`

	// 认证相关
	EmailVerified bool   `json:"emailVerified"`
	PhoneVerified bool   `json:"phoneVerified"`
	LastLoginAt   string `json:"lastLoginAt,omitempty"`
	LastLoginIP   string `json:"lastLoginIP,omitempty"`
}
