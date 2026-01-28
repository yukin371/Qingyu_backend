package config

import "time"

// CacheConfig 缓存配置
type CacheConfig struct {
	// 全局开关
	Enabled bool `mapstructure:"enabled" json:"enabled"` // 是否启用Repository缓存

	// 延迟双删配置
	DoubleDeleteDelay time.Duration `mapstructure:"double_delete_delay" json:"double_delete_delay"` // 延迟双删延迟时间

	// 空值缓存配置
	NullCacheTTL    time.Duration `mapstructure:"null_cache_ttl" json:"null_cache_ttl"`       // 空值缓存TTL
	NullCachePrefix string        `mapstructure:"null_cache_prefix" json:"null_cache_prefix"` // 空值缓存前缀

	// 熔断器配置
	BreakerMaxRequests uint32        `mapstructure:"breaker_max_requests" json:"breaker_max_requests"` // 熔断器最大请求数
	BreakerInterval    time.Duration `mapstructure:"breaker_interval" json:"breaker_interval"`         // 熔断器统计间隔
	BreakerTimeout     time.Duration `mapstructure:"breaker_timeout" json:"breaker_timeout"`           // 熔断器超时时间
	BreakerThreshold   float64       `mapstructure:"breaker_threshold" json:"breaker_threshold"`       // 熔断器失败率阈值（0-1）
}

// DefaultCacheConfig 返回默认缓存配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled:           false, // 默认关闭，需要显式启用
		DoubleDeleteDelay: 1 * time.Second,
		NullCacheTTL:      30 * time.Second,
		NullCachePrefix:   "@@NULL@@",

		BreakerMaxRequests: 3,
		BreakerInterval:    10 * time.Second,
		BreakerTimeout:     30 * time.Second,
		BreakerThreshold:   0.6, // 60%失败率触发熔断
	}
}

// globalCacheConfig 全局缓存配置
var globalCacheConfig *CacheConfig

// GetCacheConfig 获取缓存配置
func GetCacheConfig() *CacheConfig {
	if globalCacheConfig == nil {
		globalCacheConfig = DefaultCacheConfig()
	}
	return globalCacheConfig
}

// SetCacheConfig 设置缓存配置
func SetCacheConfig(cfg *CacheConfig) {
	globalCacheConfig = cfg
}
