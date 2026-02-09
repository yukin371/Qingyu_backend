// Package auth 提供认证服务
package auth

// 为了向后兼容，重新导出所有公共符号到shared路径
// 这样现有代码可以继续使用Qingyu_backend/service/shared/auth
//
// 当所有代码迁移完成后，可以删除这个兼容层
//
// 迁移状态：第一阶段 - 创建兼容层
// 创建日期：2026-02-09
// 计划移除版本：v1.2

// 导出服务类型
type AuthService = AuthServiceImpl
type JWTService = JWTServiceImpl
type OAuthService = OAuthServiceImpl
type SessionService = SessionServiceImpl
type PermissionService = PermissionServiceImpl
type RoleService = RoleServiceImpl

// 导出适配器类型
type RedisAdapter = RedisAdapterImpl
type MemoryBlacklist = MemoryBlacklistImpl

// 导出其他类型
type PasswordValidator = PasswordValidatorImpl

// 导出请求和响应类型
type (
	RegisterRequest    = RegisterRequest
	RegisterResponse   = RegisterResponse
	LoginRequest       = LoginRequest
	LoginResponse      = LoginResponse
	OAuthLoginRequest  = OAuthLoginRequest
	TokenClaims        = TokenClaims
	UserInfo           = UserInfo
	Session            = Session
	Role               = Role
	Permission         = Permission
	CreateRoleRequest  = CreateRoleRequest
	UpdateRoleRequest  = UpdateRoleRequest
)

// 导出配置类型
type (
	JWTConfig       = JWTConfig
	OAuthConfig     = OAuthConfig
	SessionConfig   = SessionConfig
	PasswordRules   = PasswordRules
)

// 导出函数
var (
	NewAuthService       = NewAuthService
	NewJWTService        = NewJWTService
	NewOAuthService      = NewOAuthService
	NewSessionService    = NewSessionService
	NewPermissionService = NewPermissionService
	NewRoleService       = NewRoleService
	NewPasswordValidator = NewPasswordValidator
)

// 导出常量
const (
	DefaultTokenDuration     = DefaultTokenDuration
	DefaultRefreshDuration   = DefaultRefreshDuration
	MaxPasswordLength        = MaxPasswordLength
	MinPasswordLength        = MinPasswordLength
)

// 导出OAuth提供商
type (
	OAuthProvider = OAuthProvider
)

// 导出权限常量
const (
	PermissionReadPost   = PermissionReadPost
	PermissionWritePost  = PermissionWritePost
	PermissionDeletePost = PermissionDeletePost
	PermissionManageUser = PermissionManageUser
)

// 导出角色常量
const (
	RoleAdmin  = RoleAdmin
	RoleEditor = RoleEditor
	RoleReader = RoleReader
)
