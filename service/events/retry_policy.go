package events

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

// RetryPolicy 重试策略接口
type RetryPolicy interface {
	// ShouldRetry 判断是否应该重试
	ShouldRetry(err error, attempt int) bool

	// GetDelay 获取重试延迟时间
	GetDelay(attempt int) time.Duration

	// MaxRetries 最大重试次数
	MaxRetries() int
}

// ExponentialBackoffPolicy 指数退避策略
type ExponentialBackoffPolicy struct {
	maxRetries    int           // 最大重试次数
	InitialDelay  time.Duration // 初始延迟
	maxDelay      time.Duration // 最大延迟
	Multiplier    float64       // 延迟倍数
	Randomization float64       // 随机化因子（0-1）
}

// NewExponentialBackoffPolicy 创建指数退避策略
func NewExponentialBackoffPolicy(maxRetries int, initialDelay, maxDelay time.Duration) *ExponentialBackoffPolicy {
	return &ExponentialBackoffPolicy{
		maxRetries:    maxRetries,
		InitialDelay:  initialDelay,
		maxDelay:      maxDelay,
		Multiplier:    2.0, // 默认每次翻倍
		Randomization: 0.1, // 默认10%随机化
	}
}

// ShouldRetry 判断是否应该重试
func (p *ExponentialBackoffPolicy) ShouldRetry(err error, attempt int) bool {
	if attempt >= p.maxRetries {
		return false
	}

	// 不重试已取消的错误
	if errors.Is(err, context.Canceled) {
		return false
	}

	// 不重试超时错误（可选，根据业务需求）
	if errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return true
}

// GetDelay 获取重试延迟时间（指数退避 + 随机化）
func (p *ExponentialBackoffPolicy) GetDelay(attempt int) time.Duration {
	// 计算指数退避延迟
	delay := float64(p.InitialDelay) * math.Pow(p.Multiplier, float64(attempt))

	// 应用最大延迟限制
	if delay > float64(p.maxDelay) {
		delay = float64(p.maxDelay)
	}

	// 应用随机化（避免惊群效应）
	if p.Randomization > 0 {
		// 在 (1-randomization, 1+randomization) 范围内随机
		variation := delay * p.Randomization * (2*rand.Float64() - 1)
		delay += variation
	}

	return time.Duration(delay)
}

// MaxRetries 最大重试次数
func (p *ExponentialBackoffPolicy) MaxRetries() int {
	return p.maxRetries
}

// FixedDelayPolicy 固定延迟策略
type FixedDelayPolicy struct {
	maxRetries int
	delay      time.Duration
}

// NewFixedDelayPolicy 创建固定延迟策略
func NewFixedDelayPolicy(maxRetries int, delay time.Duration) *FixedDelayPolicy {
	return &FixedDelayPolicy{
		maxRetries: maxRetries,
		delay:      delay,
	}
}

// ShouldRetry 判断是否应该重试
func (p *FixedDelayPolicy) ShouldRetry(err error, attempt int) bool {
	return attempt < p.maxRetries
}

// GetDelay 获取重试延迟时间
func (p *FixedDelayPolicy) GetDelay(attempt int) time.Duration {
	return p.delay
}

// MaxRetries 最大重试次数
func (p *FixedDelayPolicy) MaxRetries() int {
	return p.maxRetries
}

// LinearBackoffPolicy 线性退避策略
type LinearBackoffPolicy struct {
	maxRetries   int
	InitialDelay time.Duration
	increment    time.Duration
	maxDelay     time.Duration
}

// NewLinearBackoffPolicy 创建线性退避策略
func NewLinearBackoffPolicy(maxRetries int, initialDelay, increment, maxDelay time.Duration) *LinearBackoffPolicy {
	return &LinearBackoffPolicy{
		maxRetries:   maxRetries,
		InitialDelay: initialDelay,
		increment:    increment,
		maxDelay:     maxDelay,
	}
}

// ShouldRetry 判断是否应该重试
func (p *LinearBackoffPolicy) ShouldRetry(err error, attempt int) bool {
	return attempt < p.maxRetries
}

// GetDelay 获取重试延迟时间
func (p *LinearBackoffPolicy) GetDelay(attempt int) time.Duration {
	delay := p.InitialDelay + time.Duration(attempt)*p.increment
	if delay > p.maxDelay {
		delay = p.maxDelay
	}
	return delay
}

// MaxRetries 最大重试次数
func (p *LinearBackoffPolicy) MaxRetries() int {
	return p.maxRetries
}
