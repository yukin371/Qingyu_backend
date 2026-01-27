package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// TokenBucketLimiter 令牌桶限流器
//
// 使用golang.org/x/time/rate包实现
type TokenBucketLimiter struct {
	// 限流器映射，key是限流键
	limiters map[string]*limiterState
	// 读写锁
	mu sync.RWMutex
	// 配置
	config *RateLimitConfig
	// 停止通道
	stopCh chan struct{}
}

// limiterState 限流器状态
type limiterState struct {
	limiter      *rate.Limiter
	lastSeen     time.Time
	totalRequests int64
	allowedRequests int64
	rejectedRequests int64
	mu           sync.RWMutex
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(config *RateLimitConfig) (*TokenBucketLimiter, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	limiter := &TokenBucketLimiter{
		limiters: make(map[string]*limiterState),
		config:   config,
		stopCh:   make(chan struct{}),
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter, nil
}

// Allow 检查是否允许请求
func (l *TokenBucketLimiter) Allow(key string) bool {
	limiterState := l.getLimiter(key)

	limiterState.mu.Lock()
	limiterState.totalRequests++
	limiterState.mu.Unlock()

	allowed := limiterState.limiter.Allow()

	limiterState.mu.Lock()
	limiterState.lastSeen = time.Now()
	if allowed {
		limiterState.allowedRequests++
	} else {
		limiterState.rejectedRequests++
	}
	limiterState.mu.Unlock()

	return allowed
}

// Wait 等待直到可以处理请求
func (l *TokenBucketLimiter) Wait(ctx context.Context, key string) error {
	limiterState := l.getLimiter(key)

	limiterState.mu.Lock()
	limiterState.totalRequests++
	limiterState.mu.Unlock()

	err := limiterState.limiter.Wait(ctx)

	limiterState.mu.Lock()
	limiterState.lastSeen = time.Now()
	if err == nil {
		limiterState.allowedRequests++
	} else {
		limiterState.rejectedRequests++
	}
	limiterState.mu.Unlock()

	return err
}

// Reset 重置指定key的限流状态
func (l *TokenBucketLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.limiters, key)
}

// GetStats 获取指定key的统计信息
func (l *TokenBucketLimiter) GetStats(key string) *LimiterStats {
	l.mu.RLock()
	limiterState, exists := l.limiters[key]
	l.mu.RUnlock()

	if !exists {
		return nil
	}

	limiterState.mu.RLock()
	defer limiterState.mu.RUnlock()

	return &LimiterStats{
		TotalRequests:    limiterState.totalRequests,
		RejectedRequests: limiterState.rejectedRequests,
		AllowedRequests:  limiterState.allowedRequests,
		LastRequestTime:  limiterState.lastSeen,
	}
}

// getLimiter 获取或创建限流器
func (l *TokenBucketLimiter) getLimiter(key string) *limiterState {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 双重检查
	if limiterState, exists := l.limiters[key]; exists {
		return limiterState
	}

	// 创建新的限流器
	// 每秒速率和突发容量
	r := rate.Limit(l.config.Rate)
	burst := l.config.Burst

	state := &limiterState{
		limiter:  rate.NewLimiter(r, burst),
		lastSeen: time.Now(),
	}

	l.limiters[key] = state

	return state
}

// cleanup 定期清理过期的限流器
func (l *TokenBucketLimiter) cleanup() {
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
func (l *TokenBucketLimiter) doCleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	expirationTime := time.Duration(l.config.CleanupInterval) * time.Second

	for key, limiterState := range l.limiters {
		limiterState.mu.RLock()
		lastSeen := limiterState.lastSeen
		limiterState.mu.RUnlock()

		if now.Sub(lastSeen) > expirationTime {
			delete(l.limiters, key)
		}
	}
}

// Stop 停止限流器
func (l *TokenBucketLimiter) Stop() {
	close(l.stopCh)
}

// Count 返回当前限流器数量
func (l *TokenBucketLimiter) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.limiters)
}

// GetTotalStats 获取所有限流器的总统计信息
func (l *TokenBucketLimiter) GetTotalStats() *LimiterStats {
	l.mu.RLock()
	defer l.mu.RUnlock()

	totalStats := &LimiterStats{}

	for _, limiterState := range l.limiters {
		limiterState.mu.RLock()
		totalStats.TotalRequests += limiterState.totalRequests
		totalStats.RejectedRequests += limiterState.rejectedRequests
		totalStats.AllowedRequests += limiterState.allowedRequests
		if limiterState.lastSeen.After(totalStats.LastRequestTime) {
			totalStats.LastRequestTime = limiterState.lastSeen
		}
		limiterState.mu.RUnlock()
	}

	return totalStats
}

// 确保实现了Limiter接口
var _ Limiter = (*TokenBucketLimiter)(nil)
