package auth

import (
	"context"
	"time"
)

// AuthService 认证服务接口（对外暴露）
type AuthService interface {
	// 用户认证
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token string) (string, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)

	// 权限管理
	CheckPermission(ctx context.Context, userID, permission string) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	HasRole(ctx context.Context, userID, role string) (bool, error)
	GetUserRoles(ctx context.Context, userID string) ([]string, error)

	// 角色管理
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error)
	UpdateRole(ctx context.Context, roleID string, req *UpdateRoleRequest) error
	DeleteRole(ctx context.Context, roleID string) error
	AssignRole(ctx context.Context, userID, roleID string) error
	RemoveRole(ctx context.Context, userID, roleID string) error

	// 会话管理
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	DestroySession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string) error

	// 健康检查
	Health(ctx context.Context) error
}

// JWTService JWT令牌服务接口
type JWTService interface {
	GenerateToken(ctx context.Context, userID string, roles []string) (string, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	RefreshToken(ctx context.Context, token string) (string, error)
	RevokeToken(ctx context.Context, token string) error
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}

// RoleService 角色服务接口
type RoleService interface {
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error)
	GetRole(ctx context.Context, roleID string) (*Role, error)
	UpdateRole(ctx context.Context, roleID string, req *UpdateRoleRequest) error
	DeleteRole(ctx context.Context, roleID string) error
	ListRoles(ctx context.Context) ([]*Role, error)
	AssignPermissions(ctx context.Context, roleID string, permissions []string) error
	RemovePermissions(ctx context.Context, roleID string, permissions []string) error
}

// PermissionService 权限服务接口
type PermissionService interface {
	CheckPermission(ctx context.Context, userID, permission string) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	HasRole(ctx context.Context, userID, role string) (bool, error)
}

// SessionService 会话服务接口
type SessionService interface {
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error
	DestroySession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string) error
}

// ============ 请求/响应结构 ============

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // 可选，默认为 "reader"
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	User  *UserInfo `json:"user"`
	Token string    `json:"token"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  *UserInfo `json:"user"`
	Token string    `json:"token"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// ============ 数据结构 ============

// TokenClaims JWT令牌声明
type TokenClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Exp    int64    `json:"exp"`
}

// Role 角色
type Role struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Permissions []string  `json:"permissions" bson:"permissions"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Session 会话
type Session struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
}

// UserInfo 用户信息（简化版，用于响应）
type UserInfo struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}
