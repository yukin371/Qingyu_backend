// Package jwt owns low-level JWT signing, parsing and verification rules.
// It has no dependency on HTTP middleware or application services.
package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"
)

const (
	DefaultAccessExpiration  = 2 * time.Hour
	DefaultRefreshExpiration = 7 * 24 * time.Hour
	DefaultIssuer            = "qingyu"
)

// Claims is the wire-level token claims contract shared by the verifier and
// its adapters. Application services map it to their own response type.
type Claims struct {
	UserID    string   `json:"user_id"`
	Username  string   `json:"username,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	TokenType string   `json:"token_type"`
	jwtgo.RegisteredClaims
}

// Manager is the low-level token port consumed by middleware and services.
type Manager interface {
	GenerateTokenPair(userID string, extraClaims map[string]interface{}) (accessToken, refreshToken string, err error)
	ValidateToken(tokenString string) (*Claims, error)
	ParseToken(tokenString string) (*Claims, error)
	RefreshToken(refreshToken string) (string, error)
}

type manager struct {
	secret            []byte
	issuer            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

// NewManager creates the canonical JWT manager.
func NewManager(secret string, accessExpiration, refreshExpiration time.Duration) (Manager, error) {
	if secret == "" {
		return nil, errors.New("secret cannot be empty")
	}
	if accessExpiration <= 0 {
		accessExpiration = DefaultAccessExpiration
	}
	if refreshExpiration <= 0 {
		refreshExpiration = DefaultRefreshExpiration
	}
	return &manager{
		secret:            []byte(secret),
		issuer:            DefaultIssuer,
		accessExpiration:  accessExpiration,
		refreshExpiration: refreshExpiration,
	}, nil
}

func (m *manager) GenerateTokenPair(userID string, extraClaims map[string]interface{}) (string, string, error) {
	now := time.Now()
	accessClaims := &Claims{
		UserID: userID, TokenType: "access",
		RegisteredClaims: jwtgo.RegisteredClaims{
			Issuer: m.issuer, Subject: userID,
			ExpiresAt: jwtgo.NewNumericDate(now.Add(m.accessExpiration)),
			IssuedAt:  jwtgo.NewNumericDate(now), NotBefore: jwtgo.NewNumericDate(now),
		},
	}
	if username, ok := extraClaims["username"].(string); ok {
		accessClaims.Username = username
	}
	if roles, ok := extraClaims["roles"].([]string); ok {
		accessClaims.Roles = roles
	}
	accessToken, err := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, accessClaims).SignedString(m.secret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshClaims := &Claims{
		UserID: userID, TokenType: "refresh",
		RegisteredClaims: jwtgo.RegisteredClaims{
			Issuer: m.issuer, Subject: userID,
			ExpiresAt: jwtgo.NewNumericDate(now.Add(m.refreshExpiration)),
			IssuedAt:  jwtgo.NewNumericDate(now), NotBefore: jwtgo.NewNumericDate(now),
		},
	}
	refreshToken, err := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, refreshClaims).SignedString(m.secret)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return accessToken, refreshToken, nil
}

func (m *manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwtgo.ParseWithClaims(tokenString, &Claims{}, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

func (m *manager) ParseToken(tokenString string) (*Claims, error) {
	token, _, err := jwtgo.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

func (m *manager) RefreshToken(refreshToken string) (string, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}
	if claims.TokenType != "refresh" {
		return "", errors.New("token is not a refresh token")
	}
	now := time.Now()
	newClaims := &Claims{
		UserID: claims.UserID, Username: claims.Username, Roles: claims.Roles, TokenType: "access",
		RegisteredClaims: jwtgo.RegisteredClaims{
			Issuer: m.issuer, Subject: claims.UserID,
			ExpiresAt: jwtgo.NewNumericDate(now.Add(m.accessExpiration)),
			IssuedAt:  jwtgo.NewNumericDate(now), NotBefore: jwtgo.NewNumericDate(now),
		},
	}
	newAccessToken, err := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, newClaims).SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign new access token: %w", err)
	}
	return newAccessToken, nil
}
