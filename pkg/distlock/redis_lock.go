package distlock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrLockAcquisitionFailed 获取锁失败
	ErrLockAcquisitionFailed = errors.New("failed to acquire lock")
	// ErrLockNotHeld 锁未持有
	ErrLockNotHeld = errors.New("lock not held")
	// ErrLockTimeout 获取锁超时
	ErrLockTimeout = errors.New("acquire lock timeout")
)

// Lock 分布式锁接口
type Lock interface {
	// Acquire 尝试获取锁
	// ctx: 上下文
	// ttl: 锁的过期时间
	// 返回 lockID 用于释放锁
	Acquire(ctx context.Context, ttl time.Duration) (string, error)

	// Release 释放锁
	// ctx: 上下文
	// lockID: Acquire 返回的锁标识
	Release(ctx context.Context, lockID string) error

	// Extend 延长锁的过期时间
	// ctx: 上下文
	// lockID: Acquire 返回的锁标识
	// ttl: 新的过期时间
	Extend(ctx context.Context, lockID string, ttl time.Duration) error
}

// RedisLockService Redis分布式锁实现
type RedisLockService struct {
	client *redis.Client
	prefix string
}

// NewRedisLockService 创建Redis分布式锁服务
// 接受 *redis.Client 直接传递（推荐方式）
func NewRedisLockService(client *redis.Client, prefix string) *RedisLockService {
	if prefix == "" {
		prefix = "distlock"
	}
	return &RedisLockService{
		client: client,
		prefix: prefix,
	}
}

// newRedisLockServiceWithClient 直接使用 redis.Client 创建分布式锁服务（内部使用）
func newRedisLockServiceWithClient(client *redis.Client, prefix string) *RedisLockService {
	if prefix == "" {
		prefix = "distlock"
	}
	return &RedisLockService{
		client: client,
		prefix: prefix,
	}
}

// NewRedisLockServiceFromAddr 通过地址创建Redis分布式锁服务
func NewRedisLockServiceFromAddr(addr, password string, db int, prefix string) (*RedisLockService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return newRedisLockServiceWithClient(client, prefix), nil
}

// Acquire 尝试获取锁
// 使用 SET key value NX PX ttl 实现原子性的分布式锁
func (s *RedisLockService) Acquire(ctx context.Context, lockKey string, ttl time.Duration) (string, error) {
	key := s.getKey(lockKey)

	// 生成唯一的锁标识（用于安全释放锁）
	lockID := generateLockID()

	// 使用 SETNX 风格的原子操作
	// SET key value NX PX ttl
	result, err := s.client.SetNX(ctx, key, lockID, ttl).Result()
	if err != nil {
		return "", fmt.Errorf("redis SetNX error: %w", err)
	}

	if !result {
		// 锁已被占用
		return "", ErrLockAcquisitionFailed
	}

	return lockID, nil
}

// AcquireWithRetry 获取锁，带重试
func (s *RedisLockService) AcquireWithRetry(ctx context.Context, lockKey string, ttl time.Duration, maxRetries int, retryInterval time.Duration) (string, error) {
	for i := 0; i < maxRetries; i++ {
		lockID, err := s.Acquire(ctx, lockKey, ttl)
		if err == nil {
			return lockID, nil
		}

		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		// 等待后重试
		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}

	return "", ErrLockTimeout
}

// Release 释放锁
// 使用 Lua 脚本保证原子性：只有锁的值等于 lockID 时才删除
func (s *RedisLockService) Release(ctx context.Context, lockKey string, lockID string) error {
	key := s.getKey(lockKey)

	// Lua 脚本：只有锁的值等于 lockID 时才删除
	// 这样可以防止误删其他请求创建的锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := s.client.Eval(ctx, script, []string{key}, lockID).Result()
	if err != nil {
		return fmt.Errorf("redis Eval error: %w", err)
	}

	if result.(int64) == 0 {
		return ErrLockNotHeld
	}

	return nil
}

// Extend 延长锁的过期时间
// 使用 Lua 脚本保证原子性：只有锁的值等于 lockID 时才延长
func (s *RedisLockService) Extend(ctx context.Context, lockKey string, lockID string, ttl time.Duration) error {
	key := s.getKey(lockKey)

	// Lua 脚本：只有锁的值等于 lockID 时才延长过期时间
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := s.client.Eval(ctx, script, []string{key}, lockID, ttl.Milliseconds()).Result()
	if err != nil {
		return fmt.Errorf("redis Eval error: %w", err)
	}

	if result.(int64) == 0 {
		return ErrLockNotHeld
	}

	return nil
}

// getKey 获取带前缀的键
func (s *RedisLockService) getKey(lockKey string) string {
	return fmt.Sprintf("%s:%s", s.prefix, lockKey)
}

// generateLockID 生成唯一的锁标识
func generateLockID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// Close 关闭Redis客户端
func (s *RedisLockService) Close() error {
	return s.client.Close()
}
