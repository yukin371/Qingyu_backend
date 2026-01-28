package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRateLimitConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, 100.0, config.RequestsPerSec)
	assert.Equal(t, 200, config.Burst)
	assert.Contains(t, config.SkipPaths, "/health")
	assert.Contains(t, config.SkipPaths, "/metrics")
}

func TestRateLimitConfig_LoadedFromTestYaml(t *testing.T) {
	// Test loading from actual test config file
	cfg, err := LoadConfig("../../config")

	assert.NoError(t, err)
	assert.NotNil(t, cfg.RateLimit)
	assert.False(t, cfg.RateLimit.Enabled, "Test environment should have rate limit disabled")
	assert.Equal(t, 10000.0, cfg.RateLimit.RequestsPerSec, "Test environment should have higher rate limit")
	assert.Equal(t, 20000, cfg.RateLimit.Burst, "Test environment should have higher burst")
	assert.Contains(t, cfg.RateLimit.SkipPaths, "/api/v1/reader", "Test environment should skip reader API paths")
}
