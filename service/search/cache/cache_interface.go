package cache

import (
	"context"
	"time"
)

// Cache 搜索缓存接口
type Cache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) ([]byte, error)

	// Set 设置缓存
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// DeletePattern 批量删除匹配模式的缓存
	DeletePattern(ctx context.Context, pattern string) error

	// Exists 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Clear 清空所有缓存
	Clear(ctx context.Context) error

	// Ping 健康检查
	Ping(ctx context.Context) error

	// Close 关闭连接
	Close() error
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// 默认过期时间
	DefaultTTL time.Duration
	// 热点数据过期时间（更短）
	HotTTL time.Duration
	// 缓存键前缀
	KeyPrefix string
	// 最大缓存条数
	MaxEntries int
	// 是否启用缓存
	Enabled bool
}

// CacheStats 缓存统计信息
type CacheStats struct {
	// 命中次数
	Hits int64
	// 未命中次数
	Misses int64
	// 缓存条数
	Keys int64
	// 内存使用
	MemoryUsage int64
	// 命中率
	HitRate float64
}

// CacheStatsProvider 缓存统计接口
type CacheStatsProvider interface {
	// Stats 获取缓存统计信息
	Stats(ctx context.Context) (*CacheStats, error)

	// Reset 重置统计信息
	Reset(ctx context.Context) error
}
