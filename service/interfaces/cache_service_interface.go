package interfaces

import (
	"context"
	"time"
)

// CacheService 缓存服务接口
type CacheService interface {
	// 基础操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// 批量操作
	MGet(ctx context.Context, keys ...string) ([]string, error)
	MSet(ctx context.Context, kvPairs map[string]string, expiration time.Duration) error
	MDelete(ctx context.Context, keys ...string) error

	// 高级操作
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)

	// 哈希操作
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key, field, value string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDelete(ctx context.Context, key string, fields ...string) error

	// 集合操作
	SAdd(ctx context.Context, key string, members ...string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key, member string) (bool, error)
	SRemove(ctx context.Context, key string, members ...string) error

	// 有序集合操作
	ZAdd(ctx context.Context, key string, score float64, member string) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) (map[string]float64, error)
	ZRemove(ctx context.Context, key string, members ...string) error

	// 服务管理
	Ping(ctx context.Context) error
	FlushDB(ctx context.Context) error
	Close() error
}

