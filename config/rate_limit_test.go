package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultRateLimitConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	// DefaultRateLimitConfig only sets Enabled=true
	// Other defaults are populated by Viper's setDefaults()
	assert.True(t, config.Enabled)
	// The following fields will be populated by Viper when LoadConfig is called
	// RequestsPerSec and Burst will be set to 100 and 200 respectively by setDefaults()
}

func TestRateLimitConfig_LoadedFromTestYaml(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.test.yaml")
	yamlContent := `
rate_limit:
  enabled: false
  requests_per_sec: 10000
  burst: 20000
  skip_paths:
    - /api/v1/reader
`
	err := os.WriteFile(configFile, []byte(yamlContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configFile)
	require.NoError(t, err)
	assert.NotNil(t, cfg.RateLimit)
	assert.False(t, cfg.RateLimit.Enabled, "Test environment should have rate limit disabled")
	assert.Equal(t, 10000.0, cfg.RateLimit.RequestsPerSec, "Test environment should have higher rate limit")
	assert.Equal(t, 20000, cfg.RateLimit.Burst, "Test environment should have higher burst")
	assert.Contains(t, cfg.RateLimit.SkipPaths, "/api/v1/reader", "Test environment should skip reader API paths")
}

func TestRateLimitConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *RateLimitConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &RateLimitConfig{
				Enabled:        true,
				RequestsPerSec: 100,
				Burst:          200,
				SkipPaths:      []string{"/health", "/metrics"},
			},
			expectError: false,
		},
		{
			name: "disabled config",
			config: &RateLimitConfig{
				Enabled:        false,
				RequestsPerSec: 0,
				Burst:          0,
			},
			expectError: false,
		},
		{
			name: "invalid - negative requests per sec",
			config: &RateLimitConfig{
				Enabled:        true,
				RequestsPerSec: -10,
				Burst:          200,
			},
			expectError: true,
			errorMsg:    "rate_limit.requests_per_sec must be positive",
		},
		{
			name: "invalid - zero requests per sec",
			config: &RateLimitConfig{
				Enabled:        true,
				RequestsPerSec: 0,
				Burst:          200,
			},
			expectError: true,
			errorMsg:    "rate_limit.requests_per_sec must be positive",
		},
		{
			name: "invalid - negative burst",
			config: &RateLimitConfig{
				Enabled:        true,
				RequestsPerSec: 100,
				Burst:          -10,
			},
			expectError: true,
			errorMsg:    "rate_limit.burst cannot be negative",
		},
		{
			name: "invalid - burst too large",
			config: &RateLimitConfig{
				Enabled:        true,
				RequestsPerSec: 100,
				Burst:          10001,
			},
			expectError: true,
			errorMsg:    "rate_limit.burst too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
