package lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"Qingyu_backend/pkg/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// DistributedLock 分布式锁接口
type DistributedLock interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
	Extend(ctx context.Context, key string, ttl time.Duration) (bool, error)
	IsLocked(ctx context.Context, key string) (bool, error)
	WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error
	WithLockRetry(ctx context.Context, key string, ttl time.Duration, maxRetries int, retryInterval time.Duration, fn func() error) error
}

// RedisDistributedLock Redis分布式锁实现
// 使用SET NX + Lua脚本实现可靠的分布式锁
type RedisDistributedLock struct {
	client *redis.Client
	// 锁的唯一标识，用于确保只有锁的持有者才能释放锁
	value string
}

// NewRedisDistributedLock 创建Redis分布式锁
func NewRedisDistributedLock(client *redis.Client) *RedisDistributedLock {
	return &RedisDistributedLock{
		client: client,
		value:  generateLockValue(),
	}
}

// generateLockValue 生成锁的唯一值
func generateLockValue() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Lock 获取锁
// key: 锁的键名
// ttl: 锁的过期时间，防止死锁
// 返回: 是否成功获取锁，错误信息
func (l *RedisDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := l.getLockKey(key)

	// 使用SET NX EX命令原子性地获取锁
	// NX: 只在键不存在时设置
	// EX: 设置过期时间（秒）
	result, err := l.client.SetNX(ctx, lockKey, l.value, ttl).Result()
	if err != nil {
		logger.Error("failed to acquire lock",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}

	if result {
		logger.Debug("lock acquired",
			zap.String("key", key),
			zap.Duration("ttl", ttl),
		)
	}

	return result, nil
}

// Unlock 释放锁
// 只有锁的持有者才能释放锁
func (l *RedisDistributedLock) Unlock(ctx context.Context, key string) error {
	lockKey := l.getLockKey(key)

	// 使用Lua脚本确保只有锁的持有者才能释放锁
	// 如果锁的值不匹配，说明锁已被其他客户端持有
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{lockKey}, l.value).Result()
	if err != nil {
		logger.Error("failed to release lock",
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	if result.(int64) == 0 {
		logger.Warn("lock was not held by this client",
			zap.String("key", key),
		)
		return errors.New("lock was not held by this client")
	}

	logger.Debug("lock released", zap.String("key", key))
	return nil
}

// Extend 延长锁的过期时间
// 只有锁的持有者才能延长锁
func (l *RedisDistributedLock) Extend(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := l.getLockKey(key)

	// 使用Lua脚本确保只有锁的持有者才能延长锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{lockKey}, l.value, int(ttl.Seconds())).Result()
	if err != nil {
		logger.Error("failed to extend lock",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}

	if result.(int64) == 0 {
		return false, errors.New("lock was not held by this client")
	}

	logger.Debug("lock extended",
		zap.String("key", key),
		zap.Duration("ttl", ttl),
	)
	return true, nil
}

// IsLocked 检查锁是否被持有
func (l *RedisDistributedLock) IsLocked(ctx context.Context, key string) (bool, error) {
	lockKey := l.getLockKey(key)

	exists, err := l.client.Exists(ctx, lockKey).Result()
	if err != nil {
		logger.Error("failed to check lock status",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, err
	}

	return exists > 0, nil
}

// getLockKey 获取锁的键名
func (l *RedisDistributedLock) getLockKey(key string) string {
	return "lock:" + key
}

// TryLockWithRetry 尝试获取锁，支持重试
func (l *RedisDistributedLock) TryLockWithRetry(
	ctx context.Context,
	key string,
	ttl time.Duration,
	maxRetries int,
	retryInterval time.Duration,
) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		acquired, err := l.Lock(ctx, key, ttl)
		if err != nil {
			return false, err
		}

		if acquired {
			return true, nil
		}

		// 最后一次重试不需要等待
		if i < maxRetries-1 {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-time.After(retryInterval):
				continue
			}
		}
	}

	return false, errors.New("failed to acquire lock after retries")
}

// WithLock 使用锁执行函数
// 自动获取和释放锁
func (l *RedisDistributedLock) WithLock(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fn func() error,
) error {
	// 获取锁
	acquired, err := l.Lock(ctx, key, ttl)
	if err != nil {
		return err
	}

	if !acquired {
		return errors.New("failed to acquire lock")
	}

	// 确保释放锁
	defer func() {
		if err := l.Unlock(context.Background(), key); err != nil {
			logger.Error("failed to release lock in defer",
				zap.String("key", key),
				zap.Error(err),
			)
		}
	}()

	// 执行函数
	return fn()
}

// WithLockRetry 使用锁执行函数，支持重试
func (l *RedisDistributedLock) WithLockRetry(
	ctx context.Context,
	key string,
	ttl time.Duration,
	maxRetries int,
	retryInterval time.Duration,
	fn func() error,
) error {
	for i := 0; i < maxRetries; i++ {
		err := l.WithLock(ctx, key, ttl, fn)
		if err == nil {
			return nil
		}

		// 如果是获取锁失败，继续重试
		if err.Error() == "failed to acquire lock" {
			if i < maxRetries-1 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(retryInterval):
					logger.Debug("retrying to acquire lock",
						zap.String("key", key),
						zap.Int("attempt", i+2),
					)
					continue
				}
			}
		}

		return err
	}

	return errors.New("failed to acquire lock after retries")
}

// Redlock 红锁算法实现
// 用于在多个Redis实例上获取锁，提高可靠性
type Redlock struct {
	clients []*redis.Client
	quorum  int
	value   string
}

// NewRedlock 创建红锁
func NewRedlock(clients []*redis.Client) *Redlock {
	quorum := len(clients)/2 + 1
	return &Redlock{
		clients: clients,
		quorum:  quorum,
		value:   generateLockValue(),
	}
}

// Lock 在大多数Redis实例上获取锁
func (r *Redlock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := "lock:" + key
	successful := 0
	start := time.Now()

	// 在所有实例上尝试获取锁
	for _, client := range r.clients {
		acquired, err := client.SetNX(ctx, lockKey, r.value, ttl).Result()
		if err != nil {
			logger.Warn("failed to acquire lock on redis instance",
				zap.Error(err),
			)
			continue
		}

		if acquired {
			successful++
		}
	}

	// 计算获取锁所用时间
	elapsed := time.Since(start)

	// 检查是否在大多数实例上成功获取锁
	if successful >= r.quorum {
		// 确保获取锁的时间小于锁的TTL
		if elapsed < ttl {
			logger.Debug("redlock acquired",
				zap.String("key", key),
				zap.Int("successful", successful),
				zap.Int("quorum", r.quorum),
			)
			return true, nil
		}
	}

	// 获取锁失败，释放已获取的锁
	r.Unlock(ctx, key)

	return false, errors.New("failed to acquire lock on majority of instances")
}

// Unlock 释放所有实例上的锁
func (r *Redlock) Unlock(ctx context.Context, key string) error {
	lockKey := "lock:" + key
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	for _, client := range r.clients {
		client.Eval(ctx, script, []string{lockKey}, r.value)
	}

	logger.Debug("redlock released", zap.String("key", key))
	return nil
}

// Extend 暂不支持红锁的延期
func (r *Redlock) Extend(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return false, errors.New("extend not supported for redlock")
}

// IsLocked 检查锁状态（只在第一个实例上检查）
func (r *Redlock) IsLocked(ctx context.Context, key string) (bool, error) {
	if len(r.clients) == 0 {
		return false, errors.New("no redis clients")
	}

	lockKey := "lock:" + key
	exists, err := r.clients[0].Exists(ctx, lockKey).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
