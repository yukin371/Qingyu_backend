package middleware

import "time"

// Config 中间件总配置
//
// 对应 configs/middleware.yaml 的顶层结构。
type Config struct {
	Middleware MiddlewareConfigs `yaml:"middleware"`
	// PriorityOverrides 优先级覆盖配置（可选）
	PriorityOverrides map[string]int `yaml:"priority_overrides,omitempty"`
}

// MiddlewareConfigs 中间件配置集合
type MiddlewareConfigs struct {
	// ========== 静态配置中间件 ==========
	// RequestIDConfig 请求ID配置
	RequestID *RequestIDConfig `yaml:"request_id,omitempty"`
	// RecoveryConfig 异常恢复配置
	Recovery *RecoveryConfig `yaml:"recovery,omitempty"`
	// SecurityConfig 安全头配置
	Security *SecurityConfig `yaml:"security,omitempty"`
	// LoggerConfig 日志配置
	Logger *LoggerConfig `yaml:"logger,omitempty"`
	// CompressionConfig 压缩配置
	Compression *CompressionConfig `yaml:"compression,omitempty"`

	// ========== 动态配置中间件 ==========
	// RateLimitConfig 限流配置
	RateLimit *RateLimitConfig `yaml:"rate_limit,omitempty"`
	// AuthConfig 认证配置
	Auth *AuthConfig `yaml:"auth,omitempty"`
	// PermissionConfig 权限配置
	Permission *PermissionConfig `yaml:"permission,omitempty"`
}

// ========== 静态配置中间件 ==========

// RequestIDConfig 请求ID配置
type RequestIDConfig struct {
	// HeaderName 请求头名称，默认 "X-Request-ID"
	HeaderName string `yaml:"header_name,omitempty"`
	// ForceGen 是否强制生成，默认 false
	ForceGen bool `yaml:"force_gen,omitempty"`
}

// RecoveryConfig 异常恢复配置
type RecoveryConfig struct {
	// StackSize 堆栈大小，默认 4096
	StackSize int `yaml:"stack_size,omitempty"`
	// DisablePrint 是否禁用打印，默认 true
	DisablePrint bool `yaml:"disable_print,omitempty"`
}

// SecurityConfig 安全头配置
type SecurityConfig struct {
	// EnableXFrameOptions 是否启用X-Frame-Options
	EnableXFrameOptions bool `yaml:"enable_x_frame_options,omitempty"`
	// XFrameOptions X-Frame-Options值
	XFrameOptions string `yaml:"x_frame_options,omitempty"`
	// EnableHSTS 是否启用HSTS
	EnableHSTS bool `yaml:"enable_hsts,omitempty"`
	// EnableCSP 是否启用CSP
	EnableCSP bool `yaml:"enable_csp,omitempty"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	// SkipPaths 跳过记录的路径列表
	SkipPaths []string `yaml:"skip_paths,omitempty"`
}

// CompressionConfig 压缩配置
type CompressionConfig struct {
	// Enabled 是否启用压缩
	Enabled bool `yaml:"enabled,omitempty"`
	// Level 压缩级别（1-9）
	Level int `yaml:"level,omitempty"`
	// Types 压缩的内容类型列表
	Types []string `yaml:"types,omitempty"`
}

// ========== 动态配置中间件 ==========

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// Enabled 是否启用限流
	Enabled bool `yaml:"enabled,omitempty"`
	// Strategy 限流策略：token_bucket, sliding_window, redis
	Strategy string `yaml:"strategy,omitempty"`
	// Rate 速率（请求/秒）
	Rate int `yaml:"rate,omitempty"`
	// Burst 突发容量
	Burst int `yaml:"burst,omitempty"`
	// WindowSize 时间窗口大小（秒）
	WindowSize int `yaml:"window_size,omitempty"`
	// Redis Redis配置（分布式限流）
	Redis *RedisConfig `yaml:"redis,omitempty"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	// Addr Redis地址
	Addr string `yaml:"addr,omitempty"`
	// Password Redis密码
	Password string `yaml:"password,omitempty"`
	// DB Redis数据库编号
	DB int `yaml:"db,omitempty"`
	// Prefix 键前缀
	Prefix string `yaml:"prefix,omitempty"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	// Enabled 是否启用认证
	Enabled bool `yaml:"enabled,omitempty"`
	// Secret JWT密钥
	Secret string `yaml:"secret,omitempty"`
	// SkipPaths 跳过认证的路径列表
	SkipPaths []string `yaml:"skip_paths,omitempty"`
}

// PermissionConfig 权限配置
type PermissionConfig struct {
	// Enabled 是否启用权限检查
	Enabled bool `yaml:"enabled,omitempty"`
	// Strategy 权限策略：rbac, casbin
	Strategy string `yaml:"strategy,omitempty"`
	// ConfigPath 权限规则配置文件路径
	ConfigPath string `yaml:"config_path,omitempty"`
	// SkipPaths 跳过权限检查的路径列表
	SkipPaths []string `yaml:"skip_paths,omitempty"`
	// SessionTimeout 会话超时时间
	SessionTimeout time.Duration `yaml:"session_timeout,omitempty"`
}

// ========== 验证方法 ==========

// Validate 验证配置有效性
func (c *Config) Validate() error {
	// 验证限流配置
	if c.Middleware.RateLimit != nil && c.Middleware.RateLimit.Enabled {
		if c.Middleware.RateLimit.Rate <= 0 {
			return &ConfigError{
				Field:   "rate_limit.rate",
				Message: "rate must be positive",
			}
		}
		if c.Middleware.RateLimit.Burst <= 0 {
			return &ConfigError{
				Field:   "rate_limit.burst",
				Message: "burst must be positive",
			}
		}
		if c.Middleware.RateLimit.Strategy == "" {
			return &ConfigError{
				Field:   "rate_limit.strategy",
				Message: "strategy cannot be empty",
			}
		}
	}

	// 验证认证配置
	if c.Middleware.Auth != nil && c.Middleware.Auth.Enabled {
		if c.Middleware.Auth.Secret == "" {
			return &ConfigError{
				Field:   "auth.secret",
				Message: "secret cannot be empty when auth is enabled",
			}
		}
	}

	// 验证权限配置
	if c.Middleware.Permission != nil && c.Middleware.Permission.Enabled {
		if c.Middleware.Permission.Strategy == "" {
			return &ConfigError{
				Field:   "permission.strategy",
				Message: "strategy cannot be empty",
			}
		}
		if c.Middleware.Permission.ConfigPath == "" {
			return &ConfigError{
				Field:   "permission.config_path",
				Message: "config_path cannot be empty",
			}
		}
	}

	return nil
}

// ConfigError 配置错误
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Field + " - " + e.Message
}
