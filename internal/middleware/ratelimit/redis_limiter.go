package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisLimiter Redis分布式限流器
//
// 使用Redis实现分布式限流，支持多实例部署
type RedisLimiter struct {
	// Redis客户端
	client *redis.Client
	// 配置
	config *RateLimitConfig
	// 本地缓存（减少Redis访问）
	cache map[string]*LimiterStats
	// 缓存锁
	cacheMu sync.RWMutex
	// 停止通道
	stopCh chan struct{}
}

// NewRedisLimiter 创建Redis限流器
func NewRedisLimiter(config *RateLimitConfig) (*RedisLimiter, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if config.Redis == nil || config.Redis.Addr == "" {
		return nil, fmt.Errorf("redis config is required")
	}

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
		MaxRetries:   config.Redis.MaxRetries,
		DialTimeout:  config.Redis.DialTimeout,
		ReadTimeout:  config.Redis.ReadTimeout,
		WriteTimeout: config.Redis.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	limiter := &RedisLimiter{
		client: client,
		config: config,
		cache:  make(map[string]*LimiterStats),
		stopCh:  make(chan struct{}),
	}

	// 启动缓存清理协程
	go limiter.cleanup()

	return limiter, nil
}

// Allow 检查是否允许请求
func (l *RedisLimiter) Allow(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	redisKey := l.getRedisKey(key)

	// 使用Lua脚本确保原子性
	script := `
		local key = KEYS[1]
		local rate = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])

		-- 清理过期的请求记录
		redis.call('ZREMRANGEBYSCORE', key, '-inf', now - window)

		-- 获取当前窗口内的请求数
		local current = redis.call('ZCARD', key)

		if current < rate then
			-- 添加当前请求
			redis.call('ZADD', key, now, now)
			redis.call('EXPIRE', key, window + 1)
			return 1
		else
			return 0
		end
	`

	windowSize := time.Duration(l.config.WindowSize) * time.Second
	result, err := l.client.Eval(ctx, script, []string{redisKey},
		l.config.Rate,
		int(windowSize.Seconds()),
		time.Now().UnixMilli(),
	).Result()

	if err != nil {
		// Redis错误时，降级为允许请求
		return true
	}

	allowed := result.(int64) == 1

	// 更新本地缓存
	l.updateCache(key, allowed)

	return allowed
}

// Wait 等待直到可以处理请求
func (l *RedisLimiter) Wait(ctx context.Context, key string) error {
	// Redis限流器不支持等待，建议使用Allow方法
	return fmt.Errorf("redis limiter does not support Wait method")
}

// Reset 重置指定key的限流状态
func (l *RedisLimiter) Reset(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	redisKey := l.getRedisKey(key)
	l.client.Del(ctx, redisKey)

	// 清理本地缓存
	l.cacheMu.Lock()
	delete(l.cache, key)
	l.cacheMu.Unlock()
}

// GetStats 获取指定key的统计信息
func (l *RedisLimiter) GetStats(key string) *LimiterStats {
	l.cacheMu.RLock()
	defer l.cacheMu.RUnlock()

	if stats, exists := l.cache[key]; exists {
		return stats
	}

	return nil
}

// getRedisKey 获取Redis键
func (l *RedisLimiter) getRedisKey(key string) string {
	prefix := l.config.Redis.Prefix
	if prefix == "" {
		prefix = "ratelimit:"
	}
	return prefix + key
}

// updateCache 更新本地缓存
func (l *RedisLimiter) updateCache(key string, allowed bool) {
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	stats, exists := l.cache[key]
	if !exists {
		stats = &LimiterStats{}
		l.cache[key] = stats
	}

	stats.TotalRequests++
	stats.LastRequestTime = time.Now()

	if allowed {
		stats.AllowedRequests++
	} else {
		stats.RejectedRequests++
	}
}

// cleanup 定期清理本地缓存
func (l *RedisLimiter) cleanup() {
	ticker := time.NewTicker(time.Duration(l.config.CleanupInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.doCleanup()
		case <-l.stopCh:
			return
		}
	}
}

// doCleanup 执行清理
func (l *RedisLimiter) doCleanup() {
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	now := time.Now()
	expirationTime := time.Duration(l.config.CleanupInterval) * time.Second

	for key, stats := range l.cache {
		if now.Sub(stats.LastRequestTime) > expirationTime {
			delete(l.cache, key)
		}
	}
}

// Stop 停止限流器
func (l *RedisLimiter) Stop() {
	close(l.stopCh)
	l.client.Close()
}

// Count 返回当前缓存数量
func (l *RedisLimiter) Count() int {
	l.cacheMu.RLock()
	defer l.cacheMu.RUnlock()

	return len(l.cache)
}

// GetTotalStats 获取所有限流器的总统计信息
func (l *RedisLimiter) GetTotalStats() *LimiterStats {
	l.cacheMu.RLock()
	defer l.cacheMu.RUnlock()

	totalStats := &LimiterStats{}

	for _, stats := range l.cache {
		totalStats.TotalRequests += stats.TotalRequests
		totalStats.RejectedRequests += stats.RejectedRequests
		totalStats.AllowedRequests += stats.AllowedRequests
		if stats.LastRequestTime.After(totalStats.LastRequestTime) {
			totalStats.LastRequestTime = stats.LastRequestTime
		}
	}

	return totalStats
}

// GetCurrentUsage 获取当前使用率
func (l *RedisLimiter) GetCurrentUsage(key string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	redisKey := l.getRedisKey(key)
	windowSize := time.Duration(l.config.WindowSize) * time.Second

	// 获取当前窗口内的请求数
	count, err := l.client.ZCount(ctx, redisKey,
		strconv.FormatInt(time.Now().Add(-windowSize).UnixMilli(), 10),
		strconv.FormatInt(time.Now().UnixMilli(), 10),
	).Result()

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// 确保实现了Limiter接口
var _ Limiter = (*RedisLimiter)(nil)
