package core

import (
	"testing"

	"Qingyu_backend/config"
	"Qingyu_backend/internal/middleware/ratelimit"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBuildLoggerMiddlewareConfig(t *testing.T) {
	logCfg := &config.LogConfig{
		Mode:       "strict",
		RedactKeys: []string{"authorization", "password"},
		Request: &config.LogRequestConfig{
			SkipPaths:      []string{"/health", "/metrics"},
			BodyAllowPaths: []string{"/api/v1"},
			EnableBody:     false,
			MaxBodySize:    4096,
		},
	}

	loggerCfg := buildLoggerMiddlewareConfig(logCfg)
	if assert.NotNil(t, loggerCfg) {
		assert.Equal(t, true, loggerCfg["enable_request_body"])
		assert.Equal(t, 4096, loggerCfg["max_body_size"])
		assert.Equal(t, "strict", loggerCfg["mode"])

		skipPaths, ok := loggerCfg["skip_paths"].([]interface{})
		if assert.True(t, ok) {
			assert.Equal(t, []interface{}{"/health", "/metrics"}, skipPaths)
		}

		bodyAllowPaths, ok := loggerCfg["body_allow_paths"].([]interface{})
		if assert.True(t, ok) {
			assert.Equal(t, []interface{}{"/api/v1"}, bodyAllowPaths)
		}

		redactKeys, ok := loggerCfg["redact_keys"].([]interface{})
		if assert.True(t, ok) {
			assert.Equal(t, []interface{}{"authorization", "password"}, redactKeys)
		}
	}
}

func TestBuildLoggerMiddlewareConfigNilRequest(t *testing.T) {
	logCfg := &config.LogConfig{Mode: "normal"}
	assert.Nil(t, buildLoggerMiddlewareConfig(logCfg))
	assert.Nil(t, buildLoggerMiddlewareConfig(nil))
}

func TestBuildRateLimitMiddlewareConfig(t *testing.T) {
	cfg := &config.RateLimitConfig{
		Enabled:        true,
		RequestsPerSec: 100,
		Burst:          200,
		SkipPaths:      []string{"/health", "/metrics"},
	}

	rateLimitCfg := buildRateLimitMiddlewareConfig(cfg)
	if assert.NotNil(t, rateLimitCfg) {
		assert.Equal(t, "token_bucket", rateLimitCfg.Strategy)
		assert.Equal(t, 100, rateLimitCfg.Rate)
		assert.Equal(t, 200, rateLimitCfg.Burst)
		assert.Equal(t, "ip", rateLimitCfg.KeyFunc)
		assert.Equal(t, 429, rateLimitCfg.StatusCode)
		assert.Equal(t, "请求过于频繁，请稍后再试", rateLimitCfg.Message)
		assert.Equal(t, []string{"/health", "/metrics"}, rateLimitCfg.SkipPaths)

		_, err := ratelimit.NewRateLimitMiddleware(rateLimitCfg, zap.NewNop())
		assert.NoError(t, err)
	}
}

func TestBuildRateLimitMiddlewareConfigClampValues(t *testing.T) {
	cfg := &config.RateLimitConfig{
		Enabled:        true,
		RequestsPerSec: 0.5, // int 转换后为0，需要兜底
		Burst:          0,   // 小于rate，需要兜底
	}

	rateLimitCfg := buildRateLimitMiddlewareConfig(cfg)
	if assert.NotNil(t, rateLimitCfg) {
		assert.Equal(t, 1, rateLimitCfg.Rate)
		assert.Equal(t, 1, rateLimitCfg.Burst)
	}
}
