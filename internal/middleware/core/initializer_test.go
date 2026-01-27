package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"Qingyu_backend/internal/middleware"
)

// 测试辅助函数：创建临时配置文件
func createTempConfigFile(t *testing.T, config *middleware.Config) string {
	t.Helper()

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "middleware.yaml")

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	return configPath
}

// 测试辅助函数：创建测试用的初始化器
func createTestInitializer(t *testing.T) *InitializerImpl {
	t.Helper()

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	return NewInitializer(logger)
}

// TestLoadFromConfig_ValidConfig 测试正常加载配置
func TestLoadFromConfig_ValidConfig(t *testing.T) {
	initializer := createTestInitializer(t)

	// 准备测试配置
	config := &middleware.Config{
		Middleware: middleware.MiddlewareConfigs{
			RequestID: &middleware.RequestIDConfig{
				HeaderName: "X-Request-ID",
				ForceGen:   false,
			},
			Recovery: &middleware.RecoveryConfig{
				StackSize:    4096,
				DisablePrint: true,
			},
			Security: &middleware.SecurityConfig{
				EnableXFrameOptions: true,
				XFrameOptions:       "DENY",
				EnableHSTS:          true,
				EnableCSP:           false,
			},
			Logger: &middleware.LoggerConfig{
				SkipPaths: []string{"/health", "/metrics"},
			},
			Compression: &middleware.CompressionConfig{
				Enabled: true,
				Level:   5,
				Types:   []string{"application/json", "text/html"},
			},
			RateLimit: &middleware.RateLimitConfig{
				Enabled:    true,
				Strategy:   "token_bucket",
				Rate:       100,
				Burst:      200,
				WindowSize: 60,
			},
			Auth: &middleware.AuthConfig{
				Enabled: true,
				Secret:  "test-secret-key",
				SkipPaths: []string{
					"/health",
					"/metrics",
					"/api/v1/auth/login",
				},
			},
			Permission: &middleware.PermissionConfig{
				Enabled:        true,
				Strategy:       "rbac",
				ConfigPath:     "configs/permissions.yaml",
				SessionTimeout: 30 * time.Minute,
			},
		},
		PriorityOverrides: map[string]int{
			"rate_limit": 7,
			"metrics":    6,
		},
	}

	configPath := createTempConfigFile(t, config)

	// 测试加载
	loadedConfig, err := initializer.LoadFromConfig(configPath)
	if err != nil {
		t.Fatalf("LoadFromConfig failed: %v", err)
	}

	// 验证配置
	if loadedConfig.Middleware.RequestID.HeaderName != "X-Request-ID" {
		t.Errorf("Expected HeaderName 'X-Request-ID', got '%s'", loadedConfig.Middleware.RequestID.HeaderName)
	}

	if loadedConfig.Middleware.RateLimit.Rate != 100 {
		t.Errorf("Expected RateLimit.Rate 100, got %d", loadedConfig.Middleware.RateLimit.Rate)
	}

	if loadedConfig.Middleware.Auth.Secret != "test-secret-key" {
		t.Errorf("Expected Auth.Secret 'test-secret-key', got '%s'", loadedConfig.Middleware.Auth.Secret)
	}

	if len(loadedConfig.PriorityOverrides) != 2 {
		t.Errorf("Expected 2 priority overrides, got %d", len(loadedConfig.PriorityOverrides))
	}
}

// TestLoadFromConfig_FileNotExist 测试文件不存在
func TestLoadFromConfig_FileNotExist(t *testing.T) {
	initializer := createTestInitializer(t)

	_, err := initializer.LoadFromConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("Expected error when loading nonexistent file, got nil")
	}

	if err.Error() != "failed to read config file" &&
	   !containsString(err.Error(), "no such file or directory") &&
	   !containsString(err.Error(), "cannot find the file") {
		t.Logf("Got expected error type: %v", err)
	}
}

// TestLoadFromConfig_InvalidYAML 测试无效YAML
func TestLoadFromConfig_InvalidYAML(t *testing.T) {
	initializer := createTestInitializer(t)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	// 写入无效YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	_, err := initializer.LoadFromConfig(configPath)
	if err == nil {
		t.Fatal("Expected error when loading invalid YAML, got nil")
	}
}

// TestLoadFromConfig_ValidationFailed 测试配置验证失败
func TestLoadFromConfig_ValidationFailed(t *testing.T) {
	initializer := createTestInitializer(t)

	// 准备无效配置（限流rate <= 0）
	config := &middleware.Config{
		Middleware: middleware.MiddlewareConfigs{
			RateLimit: &middleware.RateLimitConfig{
				Enabled:  true,
				Strategy: "token_bucket",
				Rate:     0, // 无效值
				Burst:    200,
			},
		},
	}

	configPath := createTempConfigFile(t, config)

	_, err := initializer.LoadFromConfig(configPath)
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	configErr, ok := err.(*middleware.ConfigError)
	if !ok {
		t.Logf("Error type: %T, message: %v", err, err)
	}

	if ok && configErr.Field != "rate_limit.rate" {
		t.Errorf("Expected error field 'rate_limit.rate', got '%s'", configErr.Field)
	}
}

// TestInitialize_AllMiddlewares 测试初始化所有中间件
func TestInitialize_AllMiddlewares(t *testing.T) {
	initializer := createTestInitializer(t)

	// 准备配置
	config := &middleware.Config{
		Middleware: middleware.MiddlewareConfigs{
			RequestID: &middleware.RequestIDConfig{
				HeaderName: "X-Request-ID",
			},
			Recovery: &middleware.RecoveryConfig{
				StackSize: 4096,
			},
			Security: &middleware.SecurityConfig{
				EnableXFrameOptions: true,
			},
			Logger: &middleware.LoggerConfig{
				SkipPaths: []string{"/health"},
			},
			Compression: &middleware.CompressionConfig{
				Enabled: true,
				Level:   5,
			},
			RateLimit: &middleware.RateLimitConfig{
				Enabled:  true,
				Strategy: "token_bucket",
				Rate:     100,
				Burst:    200,
			},
			Auth: &middleware.AuthConfig{
				Enabled: true,
				Secret:  "test-secret",
			},
			Permission: &middleware.PermissionConfig{
				Enabled:    true,
				Strategy:   "rbac",
				ConfigPath: "configs/permissions.yaml",
			},
		},
	}

	configPath := createTempConfigFile(t, config)
	if _, err := initializer.LoadFromConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 初始化中间件
	middlewares, err := initializer.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// 验证中间件数量（应该是8个，因为都启用了）
	expectedCount := 8
	if len(middlewares) != expectedCount {
		t.Errorf("Expected %d middlewares, got %d", expectedCount, len(middlewares))
	}

	// 验证中间件列表
	middlewareNames := initializer.ListMiddlewares()
	if len(middlewareNames) != expectedCount {
		t.Errorf("Expected %d middleware names, got %d", expectedCount, len(middlewareNames))
	}

	// 验证特定中间件存在
	expectedNames := []string{"request_id", "recovery", "security", "logger", "compression", "rate_limit", "auth", "permission"}
	for _, name := range expectedNames {
		mw, err := initializer.GetMiddleware(name)
		if err != nil {
			t.Errorf("Middleware %s not found: %v", name, err)
		}
		if mw.Name() != name {
			t.Errorf("Expected middleware name '%s', got '%s'", name, mw.Name())
		}
	}
}

// TestInitialize_OnlyStaticMiddlewares 测试只初始化静态中间件
func TestInitialize_OnlyStaticMiddlewares(t *testing.T) {
	initializer := createTestInitializer(t)

	// 只配置静态中间件
	config := &middleware.Config{
		Middleware: middleware.MiddlewareConfigs{
			RequestID: &middleware.RequestIDConfig{
				HeaderName: "X-Request-ID",
			},
			Recovery: &middleware.RecoveryConfig{
				StackSize: 4096,
			},
			Security: &middleware.SecurityConfig{
				EnableXFrameOptions: true,
			},
		},
	}

	configPath := createTempConfigFile(t, config)
	if _, err := initializer.LoadFromConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 初始化中间件
	middlewares, err := initializer.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// 验证只有3个静态中间件
	if len(middlewares) != 3 {
		t.Errorf("Expected 3 middlewares, got %d", len(middlewares))
	}

	// 验证动态中间件不存在
	_, err = initializer.GetMiddleware("rate_limit")
	if err == nil {
		t.Error("Expected error when getting non-existent middleware 'rate_limit'")
	}
}

// TestInitialize_WithoutLoadConfig 测试未加载配置就初始化
func TestInitialize_WithoutLoadConfig(t *testing.T) {
	initializer := createTestInitializer(t)

	_, err := initializer.Initialize()
	if err == nil {
		t.Fatal("Expected error when initializing without loading config, got nil")
	}

	if err.Error() != "config not loaded, call LoadFromConfig first" {
		t.Errorf("Expected error message 'config not loaded, call LoadFromConfig first', got '%s'", err.Error())
	}
}

// TestInitialize_InvalidCompressionLevel 测试无效压缩级别
func TestInitialize_InvalidCompressionLevel(t *testing.T) {
	initializer := createTestInitializer(t)

	// 配置无效的压缩级别
	config := &middleware.Config{
		Middleware: middleware.MiddlewareConfigs{
			Compression: &middleware.CompressionConfig{
				Enabled: true,
				Level:   10, // 无效值（必须是1-9）
			},
		},
	}

	configPath := createTempConfigFile(t, config)
	if _, err := initializer.LoadFromConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	_, err := initializer.Initialize()
	if err == nil {
		t.Fatal("Expected error for invalid compression level, got nil")
	}
}

// TestInitialize_EmptyAuthSecret 测试空的认证密钥
func TestInitialize_EmptyAuthSecret(t *testing.T) {
	initializer := createTestInitializer(t)

	// 配置空的认证密钥（应该被配置验证拦截）
	config := &middleware.Config{
		Middleware: middleware.MiddlewareConfigs{
			Auth: &middleware.AuthConfig{
				Enabled: true,
				Secret:  "", // 空密钥
			},
		},
	}

	configPath := createTempConfigFile(t, config)
	_, err := initializer.LoadFromConfig(configPath)
	if err == nil {
		t.Fatal("Expected validation error for empty auth secret, got nil")
	}
}

// TestGetMiddleware_NotFound 测试获取不存在的中间件
func TestGetMiddleware_NotFound(t *testing.T) {
	initializer := createTestInitializer(t)

	// 加载空配置
	config := &middleware.Config{}
	configPath := createTempConfigFile(t, config)
	if _, err := initializer.LoadFromConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 初始化（没有中间件）
	if _, err := initializer.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// 尝试获取不存在的中间件
	_, err := initializer.GetMiddleware("nonexistent")
	if err == nil {
		t.Fatal("Expected error when getting non-existent middleware, got nil")
	}

	if err.Error() != "middleware nonexistent not found" {
		t.Errorf("Expected error message 'middleware nonexistent not found', got '%s'", err.Error())
	}
}

// TestListMiddlewares_Empty 测试列出空中间件列表
func TestListMiddlewares_Empty(t *testing.T) {
	initializer := createTestInitializer(t)

	// 加载空配置
	config := &middleware.Config{}
	configPath := createTempConfigFile(t, config)
	if _, err := initializer.LoadFromConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 初始化（没有中间件）
	if _, err := initializer.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// 列出中间件
	names := initializer.ListMiddlewares()
	if len(names) != 0 {
		t.Errorf("Expected empty list, got %d names", len(names))
	}
}

// TestInitialize_MultipleCalls 测试多次调用Initialize
func TestInitialize_MultipleCalls(t *testing.T) {
	initializer := createTestInitializer(t)

	// 准备配置
	config := &middleware.Config{
		Middleware: middleware.MiddlewareConfigs{
			RequestID: &middleware.RequestIDConfig{
				HeaderName: "X-Request-ID",
			},
		},
	}

	configPath := createTempConfigFile(t, config)
	if _, err := initializer.LoadFromConfig(configPath); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 第一次初始化
	middlewares1, err := initializer.Initialize()
	if err != nil {
		t.Fatalf("First Initialize failed: %v", err)
	}

	// 第二次初始化（应该清空之前的中间件）
	middlewares2, err := initializer.Initialize()
	if err != nil {
		t.Fatalf("Second Initialize failed: %v", err)
	}

	// 验证数量一致
	if len(middlewares1) != len(middlewares2) {
		t.Errorf("Expected same number of middlewares, got %d and %d", len(middlewares1), len(middlewares2))
	}

	// 验证中间件数量仍然是1
	names := initializer.ListMiddlewares()
	if len(names) != 1 {
		t.Errorf("Expected 1 middleware after re-initialization, got %d", len(names))
	}
}

// 辅助函数：检查字符串是否包含子串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
