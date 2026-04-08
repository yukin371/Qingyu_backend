package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJWTConfigRejectsWeakSecrets(t *testing.T) {
	t.Run("missing secret", func(t *testing.T) {
		err := validateJWTConfig(&JWTConfig{Secret: "", ExpirationHours: 24})
		assert.ErrorContains(t, err, "JWT secret is required")
	})

	t.Run("known default secret", func(t *testing.T) {
		err := validateJWTConfig(&JWTConfig{Secret: "qingyu_secret_key", ExpirationHours: 24})
		assert.ErrorContains(t, err, "insecure default")
	})

	t.Run("placeholder secret", func(t *testing.T) {
		err := validateJWTConfig(&JWTConfig{Secret: "${JWT_SECRET}", ExpirationHours: 24})
		assert.ErrorContains(t, err, "insecure default")
	})

	t.Run("short secret", func(t *testing.T) {
		err := validateJWTConfig(&JWTConfig{Secret: "short-secret", ExpirationHours: 24})
		assert.ErrorContains(t, err, "at least 16 characters")
	})

	t.Run("strong secret", func(t *testing.T) {
		err := validateJWTConfig(&JWTConfig{Secret: "fixture-strong-jwt-placeholder-0001", ExpirationHours: 24})
		assert.NoError(t, err)
	})
}
