package config

import "time"

// RedisConfig Redis配置
type RedisConfig struct {
	// 连接配置
	Host     string `mapstructure:"host" json:"host"`         // Redis主机地址
	Port     int    `mapstructure:"port" json:"port"`         // Redis端口
	Password string `mapstructure:"password" json:"password"` // 密码
	DB       int    `mapstructure:"db" json:"db"`             // 数据库索引

	// 连接池配置
	PoolSize     int `mapstructure:"pool_size" json:"pool_size"`           // 连接池大小
	MinIdleConns int `mapstructure:"min_idle_conns" json:"min_idle_conns"` // 最小空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns" json:"max_idle_conns"` // 最大空闲连接数

	// 超时配置
	DialTimeout  time.Duration `mapstructure:"dial_timeout" json:"dial_timeout"`   // 连接超时
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`   // 读超时
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"` // 写超时
	PoolTimeout  time.Duration `mapstructure:"pool_timeout" json:"pool_timeout"`   // 连接池超时

	// 重连配置
	MaxRetries      int           `mapstructure:"max_retries" json:"max_retries"`             // 最大重试次数
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff" json:"min_retry_backoff"` // 最小重试间隔
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff" json:"max_retry_backoff"` // 最大重试间隔
}

// DefaultRedisConfig 返回默认Redis配置
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,

		PoolSize:     10,
		MinIdleConns: 5,
		MaxIdleConns: 10,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,

		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	}
}

// globalRedisConfig 全局Redis配置
var globalRedisConfig *RedisConfig

// GetRedisConfig 获取Redis配置
func GetRedisConfig() *RedisConfig {
	if globalRedisConfig == nil {
		globalRedisConfig = DefaultRedisConfig()
	}
	return globalRedisConfig
}

// SetRedisConfig 设置Redis配置
func SetRedisConfig(cfg *RedisConfig) {
	globalRedisConfig = cfg
}
