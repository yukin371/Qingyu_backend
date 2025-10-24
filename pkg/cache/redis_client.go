package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"Qingyu_backend/config"

	"github.com/go-redis/redis/v8"
)

// Redis错误定义
var (
	ErrRedisNil         = errors.New("redis: nil returned")
	ErrKeyNotFound      = errors.New("redis: key not found")
	ErrConnectionFailed = errors.New("redis: connection failed")
	ErrTimeout          = errors.New("redis: operation timeout")
)

// RedisClient Redis客户端接口
type RedisClient interface {
	// 基础操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)

	// 批量操作
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	MSet(ctx context.Context, pairs ...interface{}) error

	// 过期时间管理
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Hash操作
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// Set操作
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SRem(ctx context.Context, key string, members ...interface{}) error

	// 原子操作
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)

	// 健康检查
	Ping(ctx context.Context) error

	// 生命周期
	Close() error

	// 获取原始客户端（用于高级操作）
	GetClient() interface{}
}

// redisClientImpl Redis客户端实现
type redisClientImpl struct {
	client *redis.Client
	config *config.RedisConfig
	mu     sync.RWMutex
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg *config.RedisConfig) (RedisClient, error) {
	if cfg == nil {
		cfg = config.DefaultRedisConfig()
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,

		// 连接池配置
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxConnAge:   0, // 0表示不限制连接年龄

		// 超时配置
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,

		// 重连配置
		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
	})

	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	return &redisClientImpl{
		client: client,
		config: cfg,
	}, nil
}

// wrapRedisError 包装Redis错误
func wrapRedisError(err error) error {
	if err == nil {
		return nil
	}
	if err == redis.Nil {
		return ErrRedisNil
	}
	if strings.Contains(err.Error(), "timeout") {
		return ErrTimeout
	}
	return err
}

// ========== 基础操作 ==========

// Get 获取键值
func (r *redisClientImpl) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", wrapRedisError(err)
	}
	return val, nil
}

// Set 设置键值
func (r *redisClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	return wrapRedisError(err)
}

// Delete 删除键
func (r *redisClientImpl) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	err := r.client.Del(ctx, keys...).Err()
	return wrapRedisError(err)
}

// Exists 检查键是否存在
func (r *redisClientImpl) Exists(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	count, err := r.client.Exists(ctx, keys...).Result()
	if err != nil {
		return 0, wrapRedisError(err)
	}
	return count, nil
}

// ========== 批量操作 ==========

// MGet 批量获取
func (r *redisClientImpl) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	if len(keys) == 0 {
		return []interface{}{}, nil
	}
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, wrapRedisError(err)
	}
	return vals, nil
}

// MSet 批量设置
func (r *redisClientImpl) MSet(ctx context.Context, pairs ...interface{}) error {
	if len(pairs) == 0 {
		return nil
	}
	if len(pairs)%2 != 0 {
		return errors.New("pairs must be even number")
	}
	err := r.client.MSet(ctx, pairs...).Err()
	return wrapRedisError(err)
}

// ========== 过期时间管理 ==========

// Expire 设置过期时间
func (r *redisClientImpl) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	return wrapRedisError(err)
}

// TTL 获取剩余过期时间
func (r *redisClientImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, wrapRedisError(err)
	}
	return ttl, nil
}

// ========== Hash操作 ==========

// HGet 获取Hash字段值
func (r *redisClientImpl) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		return "", wrapRedisError(err)
	}
	return val, nil
}

// HSet 设置Hash字段
func (r *redisClientImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	if len(values) == 0 {
		return nil
	}
	err := r.client.HSet(ctx, key, values...).Err()
	return wrapRedisError(err)
}

// HGetAll 获取Hash所有字段
func (r *redisClientImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	vals, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, wrapRedisError(err)
	}
	return vals, nil
}

// HDel 删除Hash字段
func (r *redisClientImpl) HDel(ctx context.Context, key string, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}
	err := r.client.HDel(ctx, key, fields...).Err()
	return wrapRedisError(err)
}

// ========== Set操作 ==========

// SAdd 添加Set成员
func (r *redisClientImpl) SAdd(ctx context.Context, key string, members ...interface{}) error {
	if len(members) == 0 {
		return nil
	}
	err := r.client.SAdd(ctx, key, members...).Err()
	return wrapRedisError(err)
}

// SMembers 获取Set所有成员
func (r *redisClientImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, wrapRedisError(err)
	}
	return members, nil
}

// SRem 删除Set成员
func (r *redisClientImpl) SRem(ctx context.Context, key string, members ...interface{}) error {
	if len(members) == 0 {
		return nil
	}
	err := r.client.SRem(ctx, key, members...).Err()
	return wrapRedisError(err)
}

// ========== 原子操作 ==========

// Incr 自增
func (r *redisClientImpl) Incr(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, wrapRedisError(err)
	}
	return val, nil
}

// Decr 自减
func (r *redisClientImpl) Decr(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, wrapRedisError(err)
	}
	return val, nil
}

// IncrBy 增加指定值
func (r *redisClientImpl) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := r.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, wrapRedisError(err)
	}
	return val, nil
}

// DecrBy 减少指定值
func (r *redisClientImpl) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := r.client.DecrBy(ctx, key, value).Result()
	if err != nil {
		return 0, wrapRedisError(err)
	}
	return val, nil
}

// ========== 健康检查 ==========

// Ping 健康检查
func (r *redisClientImpl) Ping(ctx context.Context) error {
	err := r.client.Ping(ctx).Err()
	return wrapRedisError(err)
}

// ========== 生命周期 ==========

// Close 关闭连接
func (r *redisClientImpl) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// GetClient 获取原始客户端
func (r *redisClientImpl) GetClient() interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client
}
