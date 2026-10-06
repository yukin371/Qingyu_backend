package auth

import (
	"time"

	securityjwt "Qingyu_backend/internal/security/jwt"
)

// Compatibility aliases. New code must import internal/security/jwt directly;
// this package remains a migration facade for middleware and legacy callers.
const (
	DefaultAccessExpiration  = securityjwt.DefaultAccessExpiration
	DefaultRefreshExpiration = securityjwt.DefaultRefreshExpiration
	DefaultIssuer            = securityjwt.DefaultIssuer
)

type Claims = securityjwt.Claims

type JWTManager = securityjwt.Manager

func NewJWTManager(secret string, accessExpiration, refreshExpiration time.Duration) (JWTManager, error) {
	return securityjwt.NewManager(secret, accessExpiration, refreshExpiration)
}
