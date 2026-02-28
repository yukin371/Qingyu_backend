package testutil

import (
	"testing"

	"Qingyu_backend/config"
	"Qingyu_backend/pkg/logger"
)

// EnableStrictLogAssertions enables strict logging and fails the test on warn/error.
func EnableStrictLogAssertions(t *testing.T) {
	enableStrictLogAssertionsWithPolicy(t, true)
}

// EnableStrictLogAssertionsIgnoreWarn enables strict logging while tolerating warn logs.
// E2E suites may emit expected warning-level logs on successful paths.
func EnableStrictLogAssertionsIgnoreWarn(t *testing.T) {
	enableStrictLogAssertionsWithPolicy(t, false)
}

func enableStrictLogAssertionsWithPolicy(t *testing.T, failOnWarn bool) {
	t.Helper()

	EnableStrictLogging()

	t.Cleanup(func() {
		stats := logger.GetStrictStats()
		if (failOnWarn && stats.WarnCount > 0) || stats.ErrorCount > 0 || stats.PanicCount > 0 || stats.FatalCount > 0 {
			t.Fatalf(
				"strict log violations: warn=%d error=%d panic=%d fatal=%d",
				stats.WarnCount, stats.ErrorCount, stats.PanicCount, stats.FatalCount,
			)
		}
	})
}

// EnableStrictLogging enables strict logging without binding to a test instance.
func EnableStrictLogging() {
	cfg := config.GlobalConfig
	if cfg == nil {
		cfg = &config.Config{}
		config.GlobalConfig = cfg
	}

	applyStrictLogConfig(cfg)

	_ = logger.Init(&logger.Config{
		Level:       cfg.Log.Level,
		Format:      cfg.Log.Format,
		Output:      cfg.Log.Output,
		Filename:    cfg.Log.Filename,
		Development: cfg.Log.Development,
		StrictMode:  true,
	})

	logger.ResetStrictStats()
}

// CheckStrictLogViolations adjusts exit code when strict log violations exist.
func CheckStrictLogViolations(code int) int {
	stats := logger.GetStrictStats()
	if stats.WarnCount > 0 || stats.ErrorCount > 0 || stats.PanicCount > 0 || stats.FatalCount > 0 {
		return 1
	}
	return code
}

func applyStrictLogConfig(cfg *config.Config) {
	if cfg.Log == nil {
		cfg.Log = &config.LogConfig{}
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "debug"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}
	if cfg.Log.Output == "" {
		cfg.Log.Output = "stdout"
	}
	if cfg.Log.Filename == "" {
		cfg.Log.Filename = "logs/app.log"
	}
	cfg.Log.Mode = "strict"

	if cfg.Log.Request == nil {
		cfg.Log.Request = &config.LogRequestConfig{}
	}
	cfg.Log.Request.EnableBody = true
	if cfg.Log.Request.MaxBodySize == 0 {
		cfg.Log.Request.MaxBodySize = 4096
	}
	if cfg.Log.Request.SkipPaths == nil {
		cfg.Log.Request.SkipPaths = []string{"/health", "/metrics", "/swagger"}
	}
	if cfg.Log.RedactKeys == nil {
		cfg.Log.RedactKeys = []string{"authorization", "password", "token", "cookie"}
	}
}
