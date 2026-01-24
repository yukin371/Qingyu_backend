package circuitbreaker

import (
	"sync"
	"time"
)

// CircuitState 熔断器状态
type CircuitState int

const (
	StateClosed   CircuitState = 0 // 正常（关闭）
	StateOpen     CircuitState = 1 // 熔断（打开）
	StateHalfOpen CircuitState = 2 // 半开（尝试恢复）
)

// String 返回状态的字符串表示
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "Closed"
	case StateOpen:
		return "Open"
	case StateHalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	mu               sync.Mutex
	state            CircuitState
	failureCount     int
	failureThreshold int
	successCount     int
	successThreshold int
	lastFailureTime  time.Time
	timeout          time.Duration

	// 统计信息
	totalRequests   int64
	totalSuccesses  int64
	totalFailures   int64
	lastStateChange time.Time
}

// NewCircuitBreaker 创建熔断器
// failureThreshold: 触发熔断的失败次数阈值
// timeout: 熔断后的超时时间，超时后进入半开状态
// successThreshold: 半开状态下恢复到关闭状态需要的成功次数
func NewCircuitBreaker(failureThreshold int, timeout time.Duration, successThreshold int) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		timeout:          timeout,
		successThreshold: successThreshold,
		lastStateChange:  time.Now(),
	}
}

// AllowRequest 判断是否允许请求
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalRequests++

	// 如果熔断超时，尝试进入半开状态
	if cb.state == StateOpen &&
		time.Since(cb.lastFailureTime) > cb.timeout {
		cb.setState(StateHalfOpen)
		return true
	}

	return cb.state != StateOpen
}

// RecordSuccess 记录成功
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalSuccesses++
	cb.failureCount = 0

	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.setState(StateClosed)
		}
	}
}

// RecordFailure 记录失败
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalFailures++
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.failureThreshold {
		cb.setState(StateOpen)
	}
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// setState 设置状态（内部方法，需要持有锁）
func (cb *CircuitBreaker) setState(state CircuitState) {
	cb.state = state
	cb.lastStateChange = time.Now()

	if state == StateHalfOpen {
		cb.successCount = 0
	} else if state == StateClosed {
		cb.failureCount = 0
		cb.successCount = 0
	}
}

// GetStats 获取统计信息
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return map[string]interface{}{
		"state":             cb.state.String(),
		"failureCount":      cb.failureCount,
		"successCount":      cb.successCount,
		"failureThreshold":  cb.failureThreshold,
		"successThreshold":  cb.successThreshold,
		"totalRequests":     cb.totalRequests,
		"totalSuccesses":    cb.totalSuccesses,
		"totalFailures":     cb.totalFailures,
		"lastFailureTime":   cb.lastFailureTime,
		"lastStateChange":   cb.lastStateChange,
	}
}

// Reset 重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
	cb.totalRequests = 0
	cb.totalSuccesses = 0
	cb.totalFailures = 0
	cb.lastStateChange = time.Now()
}

// GetFailureRate 获取失败率（0.0 - 1.0）
func (cb *CircuitBreaker) GetFailureRate() float64 {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.totalRequests == 0 {
		return 0.0
	}
	return float64(cb.totalFailures) / float64(cb.totalRequests)
}

// IsOpen 判断熔断器是否打开
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.GetState() == StateOpen
}

// IsClosed 判断熔断器是否关闭
func (cb *CircuitBreaker) IsClosed() bool {
	return cb.GetState() == StateClosed
}

// IsHalfOpen 判断熔断器是否半开
func (cb *CircuitBreaker) IsHalfOpen() bool {
	return cb.GetState() == StateHalfOpen
}
