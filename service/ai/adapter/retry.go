package adapter

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int           // 最大重试次数
	InitialDelay    time.Duration // 初始延迟
	MaxDelay        time.Duration // 最大延迟
	BackoffFactor   float64       // 退避因子
	Jitter          bool          // 是否添加随机抖动
	RetryableErrors []string      // 可重试的错误类型
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:    3,
		InitialDelay:  time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
		RetryableErrors: []string{
			ErrorTypeRateLimit,
			ErrorTypeTimeout,
			ErrorTypeNetworkError,
			ErrorTypeServiceUnavailable,
		},
	}
}

// RetryableFunc 可重试的函数类型
type RetryableFunc func(ctx context.Context) error

// Retryer 重试器
type Retryer struct {
	config *RetryConfig
}

// NewRetryer 创建重试器
func NewRetryer(config *RetryConfig) *Retryer {
	if config == nil {
		config = DefaultRetryConfig()
	}
	return &Retryer{config: config}
}

// Execute 执行可重试的函数
func (r *Retryer) Execute(ctx context.Context, fn RetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		// 执行函数
		err := fn(ctx)
		if err == nil {
			return nil // 成功，无需重试
		}

		lastErr = err

		// 检查是否为可重试的错误
		if !r.isRetryableError(err) {
			return err // 不可重试的错误，直接返回
		}

		// 最后一次尝试失败
		if attempt == r.config.MaxRetries {
			break
		}

		// 计算延迟时间
		delay := r.calculateDelay(attempt)

		// 等待重试
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// 继续重试
		}
	}

	return fmt.Errorf("重试 %d 次后仍然失败: %w", r.config.MaxRetries, lastErr)
}

// isRetryableError 检查错误是否可重试
func (r *Retryer) isRetryableError(err error) bool {
	if adapterErr, ok := err.(*AdapterError); ok {
		for _, retryableType := range r.config.RetryableErrors {
			if adapterErr.Type == retryableType {
				return true
			}
		}
	}
	return false
}

// calculateDelay 计算延迟时间
func (r *Retryer) calculateDelay(attempt int) time.Duration {
	// 指数退避
	delay := float64(r.config.InitialDelay) * math.Pow(r.config.BackoffFactor, float64(attempt))

	// 限制最大延迟
	if delay > float64(r.config.MaxDelay) {
		delay = float64(r.config.MaxDelay)
	}

	// 添加随机抖动
	if r.config.Jitter {
		jitter := rand.Float64() * 0.1 * delay // 10% 的随机抖动
		delay += jitter
	}

	return time.Duration(delay)
}

// ExecuteWithResult 执行可重试的函数并返回结果
func ExecuteWithResult(ctx context.Context, retryer *Retryer, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	var result interface{}
	var lastErr error

	err := retryer.Execute(ctx, func(ctx context.Context) error {
		var err error
		result, err = fn(ctx)
		lastErr = err
		return err
	})

	if err != nil {
		return nil, lastErr
	}

	return result, nil
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	maxFailures     int
	resetTimeout    time.Duration
	failureCount    int
	lastFailureTime time.Time
	state           CircuitState
}

// CircuitState 熔断器状态
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        CircuitClosed,
	}
}

// Execute 执行函数（带熔断器保护）
func (cb *CircuitBreaker) Execute(ctx context.Context, fn RetryableFunc) error {
	// 检查熔断器状态
	if cb.state == CircuitOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = CircuitHalfOpen
		} else {
			return &AdapterError{
				Code:    ErrorTypeServiceUnavailable,
				Message: "熔断器开启，服务暂时不可用",
				Type:    ErrorTypeServiceUnavailable,
			}
		}
	}

	// 执行函数
	err := fn(ctx)

	if err != nil {
		cb.onFailure()
		return err
	}

	cb.onSuccess()
	return nil
}

// onFailure 处理失败
func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.maxFailures {
		cb.state = CircuitOpen
	}
}

// onSuccess 处理成功
func (cb *CircuitBreaker) onSuccess() {
	cb.failureCount = 0
	cb.state = CircuitClosed
}

// GetState 获取熔断器状态
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// RateLimiter 限流器
type RateLimiter struct {
	tokens     chan struct{}
	refillRate time.Duration
	capacity   int
	ticker     *time.Ticker
	done       chan struct{}
}

// NewRateLimiter 创建限流器
func NewRateLimiter(capacity int, refillRate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(chan struct{}, capacity),
		refillRate: refillRate,
		capacity:   capacity,
		done:       make(chan struct{}),
	}

	// 初始化令牌桶
	for i := 0; i < capacity; i++ {
		rl.tokens <- struct{}{}
	}

	// 启动令牌补充协程
	rl.ticker = time.NewTicker(refillRate)
	go rl.refill()

	return rl
}

// refill 补充令牌
func (rl *RateLimiter) refill() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
				// 成功添加令牌
			default:
				// 令牌桶已满
			}
		case <-rl.done:
			return
		}
	}
}

// Acquire 获取令牌
func (rl *RateLimiter) Acquire(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close 关闭限流器
func (rl *RateLimiter) Close() {
	close(rl.done)
	rl.ticker.Stop()
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	retryer        *Retryer
	circuitBreaker *CircuitBreaker
	rateLimiter    *RateLimiter
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(retryConfig *RetryConfig, maxFailures int, resetTimeout time.Duration, rateLimit int, refillRate time.Duration) *ErrorHandler {
	return &ErrorHandler{
		retryer:        NewRetryer(retryConfig),
		circuitBreaker: NewCircuitBreaker(maxFailures, resetTimeout),
		rateLimiter:    NewRateLimiter(rateLimit, refillRate),
	}
}

// Execute 执行函数（带完整错误处理）
func (eh *ErrorHandler) Execute(ctx context.Context, fn RetryableFunc) error {
	// 限流
	if err := eh.rateLimiter.Acquire(ctx); err != nil {
		return &AdapterError{
			Code:    ErrorTypeRateLimit,
			Message: "请求频率过高，请稍后重试",
			Type:    ErrorTypeRateLimit,
		}
	}

	// 熔断器保护
	return eh.circuitBreaker.Execute(ctx, func(ctx context.Context) error {
		// 重试机制
		return eh.retryer.Execute(ctx, fn)
	})
}

// ExecuteWithResult 执行函数并返回结果（带完整错误处理）
func (eh *ErrorHandler) ExecuteWithResult(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	var result interface{}

	err := eh.Execute(ctx, func(ctx context.Context) error {
		var err error
		result, err = fn(ctx)
		return err
	})

	return result, err
}

// Close 关闭错误处理器
func (eh *ErrorHandler) Close() {
	if eh.rateLimiter != nil {
		eh.rateLimiter.Close()
	}
}

// GetCircuitBreakerState 获取熔断器状态
func (eh *ErrorHandler) GetCircuitBreakerState() CircuitState {
	return eh.circuitBreaker.GetState()
}
