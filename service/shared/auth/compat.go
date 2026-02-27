package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	mwauth "Qingyu_backend/internal/middleware/auth"
	repoauth "Qingyu_backend/repository/mongodb/auth"
	newauth "Qingyu_backend/service/auth"
)

// Service-layer aliases (migrated to service/auth).
type (
	AuthService           = newauth.AuthService
	OAuthServiceInterface = newauth.OAuthServiceInterface
	JWTService            = newauth.JWTService
	RoleService           = newauth.RoleService
	PermissionService     = newauth.PermissionService
	SessionService        = newauth.SessionService

	RegisterRequest   = newauth.RegisterRequest
	RegisterResponse  = newauth.RegisterResponse
	LoginRequest      = newauth.LoginRequest
	LoginResponse     = newauth.LoginResponse
	OAuthLoginRequest = newauth.OAuthLoginRequest

	CreateRoleRequest = newauth.CreateRoleRequest
	UpdateRoleRequest = newauth.UpdateRoleRequest
	TokenClaims       = newauth.TokenClaims
	Role              = newauth.Role
	Session           = newauth.Session
	UserInfo          = newauth.UserInfo

	OAuthService            = newauth.OAuthService
	AuthServiceImpl         = newauth.AuthServiceImpl
	JWTServiceImpl          = newauth.JWTServiceImpl
	PermissionServiceImpl   = newauth.PermissionServiceImpl
	RoleServiceImpl         = newauth.RoleServiceImpl
	SessionServiceImpl      = newauth.SessionServiceImpl
	RedisAdapter            = newauth.RedisAdapter
	InMemoryTokenBlacklist  = newauth.InMemoryTokenBlacklist
	PasswordValidator       = newauth.PasswordValidator
)

var (
	NewAuthService            = newauth.NewAuthService
	NewJWTService             = newauth.NewJWTService
	NewOAuthService           = newauth.NewOAuthService
	NewSessionService         = newauth.NewSessionService
	NewPermissionService      = newauth.NewPermissionService
	NewRoleService            = newauth.NewRoleService
	NewPasswordValidator      = newauth.NewPasswordValidator
	NewRedisAdapter           = newauth.NewRedisAdapter
	NewInMemoryTokenBlacklist = newauth.NewInMemoryTokenBlacklist
)

// Repository compatibility aliases.
type AuthRepository = repoauth.AuthRepository

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return repoauth.NewAuthRepository(db)
}

// Middleware compatibility aliases (migrated to internal/middleware/auth).
type (
	PermissionConfig     = mwauth.PermissionConfig
	PermissionMiddleware = mwauth.PermissionMiddleware
)

func NewPermissionMiddleware(config *PermissionConfig, logger *zap.Logger) (*PermissionMiddleware, error) {
	return mwauth.NewPermissionMiddleware(config, logger)
}

func JWTAuth() gin.HandlerFunc { return mwauth.JWTAuth() }

func GenerateToken(userID, username string, roles []string) (string, error) {
	return mwauth.GenerateToken(userID, username, roles)
}

func RequireRole(role string) gin.HandlerFunc { return mwauth.RequireRole(role) }

func RequireAnyRole(roles ...string) gin.HandlerFunc { return mwauth.RequireAnyRole(roles...) }

func RequireAllRoles(roles ...string) gin.HandlerFunc { return mwauth.RequireAllRoles(roles...) }
