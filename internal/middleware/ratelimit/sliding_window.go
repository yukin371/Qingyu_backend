package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// SlidingWindowLimiter 滑动窗口限流器
//
// 实现基于时间窗口的限流算法，比令牌桶更适合短时间突发控制
type SlidingWindowLimiter struct {
	// 窗口映射，key是限流键
	windows map[string]*slidingWindow
	// 读写锁
	mu sync.RWMutex
	// 配置
	config *RateLimitConfig
	// 停止通道
	stopCh chan struct{}
}

// slidingWindow 滑动窗口
type slidingWindow struct {
	// 请求时间戳队列（按时间递增）
	requests []time.Time
	// 窗口大小
	windowSize time.Duration
	// 最大请求数
	maxRequests int
	// 统计信息
	stats *LimiterStats
	// 互斥锁
	mu sync.Mutex
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(config *RateLimitConfig) (*SlidingWindowLimiter, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	limiter := &SlidingWindowLimiter{
		windows: make(map[string]*slidingWindow),
		config:  config,
		stopCh:  make(chan struct{}),
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter, nil
}

// Allow 检查是否允许请求
func (l *SlidingWindowLimiter) Allow(key string) bool {
	window := l.getWindow(key)

	window.mu.Lock()
	defer window.mu.Unlock()

	now := time.Now()
	windowSize := time.Duration(l.config.WindowSize) * time.Second

	// 移除窗口外的请求
	window.cleanupOldRequests(now, windowSize)

	// 检查是否超过限制
	window.stats.TotalRequests++
	if len(window.requests) >= l.config.Rate {
		window.stats.RejectedRequests++
		return false
	}

	// 添加当前请求
	window.requests = append(window.requests, now)
	window.stats.AllowedRequests++
	window.stats.LastRequestTime = now

	return true
}

// Wait 等待直到可以处理请求
func (l *SlidingWindowLimiter) Wait(ctx context.Context, key string) error {
	// 滑动窗口算法不支持等待
	// 如果需要等待，建议使用令牌桶算法
	return fmt.Errorf("sliding window limiter does not support Wait method, use token_bucket instead")
}

// Reset 重置指定key的限流状态
func (l *SlidingWindowLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.windows, key)
}

// GetStats 获取指定key的统计信息
func (l *SlidingWindowLimiter) GetStats(key string) *LimiterStats {
	l.mu.RLock()
	window, exists := l.windows[key]
	l.mu.RUnlock()

	if !exists {
		return nil
	}

	window.mu.Lock()
	defer window.mu.Unlock()

	// 返回统计信息的副本
	return &LimiterStats{
		TotalRequests:    window.stats.TotalRequests,
		RejectedRequests: window.stats.RejectedRequests,
		AllowedRequests:  window.stats.AllowedRequests,
		LastRequestTime:  window.stats.LastRequestTime,
	}
}

// getWindow 获取或创建窗口
func (l *SlidingWindowLimiter) getWindow(key string) *slidingWindow {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 双重检查
	if window, exists := l.windows[key]; exists {
		return window
	}

	// 创建新的窗口
	windowSize := time.Duration(l.config.WindowSize) * time.Second
	window := &slidingWindow{
		requests:   make([]time.Time, 0, l.config.Burst),
		windowSize: windowSize,
		maxRequests: l.config.Rate,
		stats: &LimiterStats{
			LastRequestTime: time.Now(),
		},
	}

	l.windows[key] = window

	return window
}

// cleanupOldRequests 清理窗口外的旧请求
func (w *slidingWindow) cleanupOldRequests(now time.Time, windowSize time.Duration) {
	cutoff := now.Add(-windowSize)

	// 使用二分查找找到第一个在窗口内的请求
	left, right := 0, len(w.requests)
	for left < right {
		mid := (left + right) / 2
		if w.requests[mid].Before(cutoff) {
			left = mid + 1
		} else {
			right = mid
		}
	}

	// 保留窗口内的请求
	if left > 0 {
		w.requests = w.requests[left:]
	}
}

// cleanup 定期清理过期的窗口
func (l *SlidingWindowLimiter) cleanup() {
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
func (l *SlidingWindowLimiter) doCleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	expirationTime := time.Duration(l.config.CleanupInterval) * time.Second

	for key, window := range l.windows {
		window.mu.Lock()

		// 检查窗口是否过期
		if now.Sub(window.stats.LastRequestTime) > expirationTime && len(window.requests) == 0 {
			delete(l.windows, key)
		}

		window.mu.Unlock()
	}
}

// Stop 停止限流器
func (l *SlidingWindowLimiter) Stop() {
	close(l.stopCh)
}

// Count 返回当前窗口数量
func (l *SlidingWindowLimiter) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.windows)
}

// GetTotalStats 获取所有窗口的总统计信息
func (l *SlidingWindowLimiter) GetTotalStats() *LimiterStats {
	l.mu.RLock()
	defer l.mu.RUnlock()

	totalStats := &LimiterStats{}

	for _, window := range l.windows {
		window.mu.Lock()
		totalStats.TotalRequests += window.stats.TotalRequests
		totalStats.RejectedRequests += window.stats.RejectedRequests
		totalStats.AllowedRequests += window.stats.AllowedRequests
		if window.stats.LastRequestTime.After(totalStats.LastRequestTime) {
			totalStats.LastRequestTime = window.stats.LastRequestTime
		}
		window.mu.Unlock()
	}

	return totalStats
}

// 确保实现了Limiter接口
var _ Limiter = (*SlidingWindowLimiter)(nil)
