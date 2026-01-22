package auth

import "time"

// Session 会话模型（存储在Redis中）
type Session struct {
	ID        string                 `json:"id"`                   // 会话ID
	UserID    string                 `json:"$1$2"`              // 用户ID
	Token     string                 `json:"token,omitempty"`      // JWT Token（可选）
	Data      map[string]interface{} `json:"data,omitempty"`       // 会话数据
	IP        string                 `json:"ip,omitempty"`         // IP地址
	UserAgent string                 `json:"user_agent,omitempty"` // 用户代理
	CreatedAt time.Time              `json:"$1$2"`           // 创建时间
	ExpiresAt time.Time              `json:"$1$2"`           // 过期时间
	UpdatedAt time.Time              `json:"updated_at,omitempty"` // 最后更新时间
}

// TokenBlacklist Token黑名单（存储在Redis中）
type TokenBlacklist struct {
	Token     string    `json:"token"`      // 被吊销的Token
	UserID    string    `json:"$1$2"`    // 用户ID
	Reason    string    `json:"reason"`     // 吊销原因
	RevokedAt time.Time `json:"$1$2"` // 吊销时间
	ExpiresAt time.Time `json:"$1$2"` // 过期时间（与Token原过期时间一致）
}
