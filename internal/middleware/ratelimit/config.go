package ratelimit

import (
	"fmt"
	"time"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// Enabled 是否启用限流
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Strategy 限流策略
	// 可选值: token_bucket, sliding_window, redis
	Strategy string `yaml:"strategy" json:"strategy"`

	// Rate 每秒请求数
	Rate int `yaml:"rate" json:"rate"`

	// Burst 突发容量（桶大小）
	Burst int `yaml:"burst" json:"burst"`

	// WindowSize 时间窗口大小（秒），用于滑动窗口
	WindowSize int `yaml:"window_size" json:"window_size"`

	// KeyFunc 限流键类型
	// 可选值: ip, user, path, ip_path, user_path
	KeyFunc string `yaml:"key_func" json:"key_func"`

	// SkipPaths 跳过限流的路径列表
	SkipPaths []string `yaml:"skip_paths" json:"skip_paths"`

	// SkipSuccessful 是否跳过成功的请求（只对失败的请求限流）
	SkipSuccessful bool `yaml:"skip_successful" json:"skip_successful"`

	// SkipFailedRequest 是否跳过失败的请求
	SkipFailedRequest bool `yaml:"skip_failed_request" json:"skip_failed_request"`

	// Message 限流错误消息
	Message string `yaml:"message" json:"message"`

	// StatusCode 限流时返回的HTTP状态码
	StatusCode int `yaml:"status_code" json:"status_code"`

	// Redis Redis配置（用于分布式限流）
	Redis *RedisConfig `yaml:"redis" json:"redis"`

	// CleanupInterval 清理过期限流器的间隔（秒）
	CleanupInterval int `yaml:"cleanup_interval" json:"cleanup_interval"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	// Addr Redis地址
	Addr string `yaml:"addr" json:"addr"`

	// Password Redis密码
	Password string `yaml:"password" json:"password"`

	// DB Redis数据库编号
	DB int `yaml:"db" json:"db"`

	// Prefix 键前缀
	Prefix string `yaml:"prefix" json:"prefix"`

	// PoolSize 连接池大小
	PoolSize int `yaml:"pool_size" json:"pool_size"`

	// MinIdleConns 最小空闲连接数
	MinIdleConns int `yaml:"min_idle_conns" json:"min_idle_conns"`

	// MaxRetries 最大重试次数
	MaxRetries int `yaml:"max_retries" json:"max_retries"`

	// DialTimeout 连接超时
	DialTimeout time.Duration `yaml:"dial_timeout" json:"dial_timeout"`

	// ReadTimeout 读取超时
	ReadTimeout time.Duration `yaml:"read_timeout" json:"read_timeout"`

	// WriteTimeout 写入超时
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
}

// DefaultRateLimitConfig 返回默认限流配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		Enabled:          false,
		Strategy:         "token_bucket",
		Rate:             100,
		Burst:            200,
		WindowSize:       60,
		KeyFunc:          "ip",
		SkipPaths:        []string{"/health", "/metrics"},
		SkipSuccessful:   false,
		SkipFailedRequest: false,
		Message:          "请求过于频繁，请稍后再试",
		StatusCode:       429,
		Redis:            DefaultRedisConfig(),
		CleanupInterval:  300, // 5分钟
	}
}

// DefaultRedisConfig 返回默认Redis配置
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		Prefix:       "ratelimit:",
		PoolSize:     10,
		MinIdleConns: 2,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// Validate 验证配置有效性
func (c *RateLimitConfig) Validate() error {
	if c.Rate <= 0 {
		return fmt.Errorf("rate must be positive, got %d", c.Rate)
	}

	if c.Burst <= 0 {
		return fmt.Errorf("burst must be positive, got %d", c.Burst)
	}

	if c.Burst < c.Rate {
		return fmt.Errorf("burst (%d) should be >= rate (%d)", c.Burst, c.Rate)
	}

	validStrategies := map[string]bool{
		"token_bucket":  true,
		"sliding_window": true,
		"redis":         true,
	}

	if !validStrategies[c.Strategy] {
		return fmt.Errorf("invalid strategy: %s, must be one of token_bucket, sliding_window, redis", c.Strategy)
	}

	validKeyFuncs := map[string]bool{
		"ip":         true,
		"user":       true,
		"path":       true,
		"ip_path":    true,
		"user_path":  true,
	}

	if !validKeyFuncs[c.KeyFunc] {
		return fmt.Errorf("invalid key_func: %s, must be one of ip, user, path, ip_path, user_path", c.KeyFunc)
	}

	if c.Strategy == "redis" {
		if c.Redis.Addr == "" {
			return fmt.Errorf("redis addr is required when strategy is redis")
		}
	}

	if c.StatusCode <= 0 || c.StatusCode >= 600 {
		return fmt.Errorf("invalid status_code: %d", c.StatusCode)
	}

	return nil
}

// ShouldSkipPath 检查路径是否应该跳过限流
func (c *RateLimitConfig) ShouldSkipPath(path string) bool {
	for _, skipPath := range c.SkipPaths {
		if skipPath == path {
			return true
		}
	}
	return false
}

// GetKeyFunc 获取KeyFunc
func (c *RateLimitConfig) GetKeyFunc() KeyFunc {
	return GetKeyFunc(KeyFuncType(c.KeyFunc))
}
