package auth

import (
	"context"
	"errors"
	"sync"
	"time"
)

// 定义错误
var (
	ErrRedisNil = errors.New("key does not exist")
)

// InMemoryTokenBlacklist 内存Token黑名单(Redis不可用时使用)
//
// 这是一个降级方案，当Redis不可用时，使用内存存储被吊销的Token。
// 注意：
//   - 仅适用于单实例部署
//   - 服务器重启后黑名单会丢失
//   - 不支持分布式部署
type InMemoryTokenBlacklist struct {
	mu         sync.RWMutex
	blacklist  map[string]time.Time // token -> 过期时间
	cleanTimer *time.Ticker
	closeChan  chan struct{}
}

// NewInMemoryTokenBlacklist 创建内存Token黑名单
func NewInMemoryTokenBlacklist() *InMemoryTokenBlacklist {
	bl := &InMemoryTokenBlacklist{
		blacklist: make(map[string]time.Time),
		closeChan: make(chan struct{}),
	}

	// 启动定期清理过期token的goroutine
	bl.cleanTimer = time.NewTicker(5 * time.Minute)
	go bl.cleanupExpiredTokens()

	return bl
}

// Set 添加token到黑名单
func (bl *InMemoryTokenBlacklist) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	bl.blacklist[key] = time.Now().Add(expiration)
	return nil
}

// Get 检查token是否在黑名单中
func (bl *InMemoryTokenBlacklist) Get(ctx context.Context, key string) (string, error) {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	expireTime, exists := bl.blacklist[key]
	if !exists {
		return "", ErrRedisNil
	}

	// 检查是否已过期
	if time.Now().After(expireTime) {
		// 过期了，删除并返回不存在
		delete(bl.blacklist, key)
		return "", ErrRedisNil
	}

	return "revoked", nil
}

// Exists 检查key是否存在
func (bl *InMemoryTokenBlacklist) Exists(ctx context.Context, keys ...string) (int64, error) {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	count := int64(0)
	now := time.Now()

	for _, key := range keys {
		if expireTime, exists := bl.blacklist[key]; exists {
			// 检查是否过期
			if now.Before(expireTime) {
				count++
			}
		}
	}

	return count, nil
}

// Del 从黑名单删除key（支持多个key）
func (bl *InMemoryTokenBlacklist) Del(ctx context.Context, keys ...string) error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	for _, key := range keys {
		delete(bl.blacklist, key)
	}

	return nil
}

// Delete 从黑名单删除单个key（CacheClient接口兼容）
func (bl *InMemoryTokenBlacklist) Delete(ctx context.Context, key string) error {
	return bl.Del(ctx, key)
}

// cleanupExpiredTokens 定期清理过期的token
func (bl *InMemoryTokenBlacklist) cleanupExpiredTokens() {
	for {
		select {
		case <-bl.cleanTimer.C:
			bl.mu.Lock()
			now := time.Now()
			for key, expireTime := range bl.blacklist {
				if now.After(expireTime) {
					delete(bl.blacklist, key)
				}
			}
			bl.mu.Unlock()

		case <-bl.closeChan:
			return
		}
	}
}

// Close 关闭清理器
func (bl *InMemoryTokenBlacklist) Close() error {
	bl.cleanTimer.Stop()
	close(bl.closeChan)
	return nil
}

// Ping 健康检查（内存模式始终返回正常）
func (bl *InMemoryTokenBlacklist) Ping(ctx context.Context) error {
	return nil
}

// Stats 获取黑名单统计信息
func (bl *InMemoryTokenBlacklist) Stats() map[string]interface{} {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	total := len(bl.blacklist)
	expired := 0
	now := time.Now()

	for _, expireTime := range bl.blacklist {
		if now.After(expireTime) {
			expired++
		}
	}

	return map[string]interface{}{
		"total":       total,
		"active":      total - expired,
		"expired":     expired,
		"storage":     "memory",
		"cleanup_interval": "5m",
	}
}
