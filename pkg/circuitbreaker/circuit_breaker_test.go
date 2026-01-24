package circuitbreaker

import (
	"sync"
	"testing"
	"time"
)

func TestCircuitBreaker_AllowRequest(t *testing.T) {
	cb := NewCircuitBreaker(3, 2*time.Second, 2)

	// 初始状态应该允许请求
	if !cb.AllowRequest() {
		t.Error("Expected AllowRequest to return true initially")
	}

	// 验证状态是关闭的
	if cb.GetState() != StateClosed {
		t.Errorf("Expected initial state to be Closed, got %v", cb.GetState())
	}
}

func TestCircuitBreaker_FailureTrip(t *testing.T) {
	cb := NewCircuitBreaker(3, 2*time.Second, 2)

	// 记录失败直到触发熔断
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}

	// 应该进入打开状态
	if cb.GetState() != StateOpen {
		t.Errorf("Expected StateOpen, got %v", cb.GetState())
	}

	// 不应该允许请求
	if cb.AllowRequest() {
		t.Error("Expected AllowRequest to return false when open")
	}
}

func TestCircuitBreaker_Recovery(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond, 2)

	// 触发熔断
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}

	// 验证打开状态
	if cb.GetState() != StateOpen {
		t.Errorf("Expected StateOpen after failures, got %v", cb.GetState())
	}

	// 等待超时
	time.Sleep(150 * time.Millisecond)

	// 下一次 AllowRequest 应该触发半开状态
	if !cb.AllowRequest() {
		t.Error("Expected AllowRequest to return true after timeout")
	}

	if cb.GetState() != StateHalfOpen {
		t.Errorf("Expected StateHalfOpen after timeout, got %v", cb.GetState())
	}

	// 记录成功
	cb.RecordSuccess()
	cb.RecordSuccess()

	// 应该恢复到关闭状态
	if cb.GetState() != StateClosed {
		t.Errorf("Expected StateClosed after successes, got %v", cb.GetState())
	}

	// 应该再次允许请求
	if !cb.AllowRequest() {
		t.Error("Expected AllowRequest to return true after recovery")
	}
}

func TestCircuitBreaker_HalfOpenFailure(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond, 2)

	// 触发熔断
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}

	// 等待超时
	time.Sleep(150 * time.Millisecond)

	// 进入半开
	cb.AllowRequest()

	if cb.GetState() != StateHalfOpen {
		t.Errorf("Expected StateHalfOpen, got %v", cb.GetState())
	}

	// 记录一次成功
	cb.RecordSuccess()

	// 记录失败，应该重新打开
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()

	if cb.GetState() != StateOpen {
		t.Errorf("Expected StateOpen after failure in HalfOpen, got %v", cb.GetState())
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(3, 2*time.Second, 2)

	// 记录一些操作
	cb.AllowRequest()
	cb.RecordFailure()
	cb.RecordFailure()
	cb.AllowRequest()
	cb.RecordSuccess()

	// 重置
	cb.Reset()

	// 验证状态
	stats := cb.GetStats()
	if stats["state"] != "Closed" {
		t.Errorf("Expected state to be Closed after reset, got %v", stats["state"])
	}

	if stats["failureCount"] != 0 {
		t.Errorf("Expected failureCount to be 0 after reset, got %v", stats["failureCount"])
	}

	if stats["totalRequests"].(int64) != 0 {
		t.Errorf("Expected totalRequests to be 0 after reset, got %v", stats["totalRequests"])
	}

	if stats["totalSuccesses"].(int64) != 0 {
		t.Errorf("Expected totalSuccesses to be 0 after reset, got %v", stats["totalSuccesses"])
	}

	if stats["totalFailures"].(int64) != 0 {
		t.Errorf("Expected totalFailures to be 0 after reset, got %v", stats["totalFailures"])
	}
}

func TestCircuitBreaker_GetStats(t *testing.T) {
	cb := NewCircuitBreaker(3, 2*time.Second, 2)

	// 执行一些操作（先成功后失败，避免 RecordSuccess 重置失败计数）
	cb.AllowRequest()
	cb.RecordSuccess()
	cb.AllowRequest()
	cb.RecordFailure()
	cb.RecordFailure()

	stats := cb.GetStats()

	// 验证统计信息
	if stats["state"] != "Closed" {
		t.Errorf("Expected state Closed, got %v", stats["state"])
	}

	if stats["failureCount"] != 2 {
		t.Errorf("Expected failureCount 2, got %v", stats["failureCount"])
	}

	if stats["totalRequests"] != int64(2) {
		t.Errorf("Expected totalRequests 2, got %v", stats["totalRequests"])
	}

	if stats["totalFailures"] != int64(2) {
		t.Errorf("Expected totalFailures 2, got %v", stats["totalFailures"])
	}

	if stats["totalSuccesses"] != int64(1) {
		t.Errorf("Expected totalSuccesses 1, got %v", stats["totalSuccesses"])
	}
}

func TestCircuitBreaker_GetFailureRate(t *testing.T) {
	cb := NewCircuitBreaker(5, 2*time.Second, 2)

	// 初始失败率应该是 0
	if rate := cb.GetFailureRate(); rate != 0.0 {
		t.Errorf("Expected initial failure rate 0.0, got %v", rate)
	}

	// 记录一些操作
	cb.AllowRequest()
	cb.RecordFailure()
	cb.AllowRequest()
	cb.RecordFailure()
	cb.AllowRequest()
	cb.RecordSuccess()

	// 失败率应该是 2/3
	expectedRate := 2.0 / 3.0
	if rate := cb.GetFailureRate(); rate != expectedRate {
		t.Errorf("Expected failure rate %v, got %v", expectedRate, rate)
	}
}

func TestCircuitBreaker_IsOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, 2*time.Second, 1)

	if cb.IsOpen() {
		t.Error("Expected circuit breaker to not be open initially")
	}

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	if !cb.IsOpen() {
		t.Error("Expected circuit breaker to be open after failures")
	}
}

func TestCircuitBreaker_IsClosed(t *testing.T) {
	cb := NewCircuitBreaker(2, 2*time.Second, 1)

	if !cb.IsClosed() {
		t.Error("Expected circuit breaker to be closed initially")
	}

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	if cb.IsClosed() {
		t.Error("Expected circuit breaker to not be closed after failures")
	}
}

func TestCircuitBreaker_IsHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond, 1)

	if cb.IsHalfOpen() {
		t.Error("Expected circuit breaker to not be half open initially")
	}

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	// 等待超时
	time.Sleep(150 * time.Millisecond)

	// 进入半开
	cb.AllowRequest()

	if !cb.IsHalfOpen() {
		t.Error("Expected circuit breaker to be half open after timeout")
	}
}

func TestCircuitBreaker_StateString(t *testing.T) {
	tests := []struct {
		state    CircuitState
		expected string
	}{
		{StateClosed, "Closed"},
		{StateOpen, "Open"},
		{StateHalfOpen, "HalfOpen"},
		{CircuitState(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("State.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCircuitBreaker_ConcurrentAccess(t *testing.T) {
	cb := NewCircuitBreaker(100, 1*time.Second, 10)
	var wg sync.WaitGroup

	// 并发访问
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cb.AllowRequest()
			cb.RecordFailure()
			cb.GetStats()
		}()
	}

	wg.Wait()

	// 验证状态一致性
	stats := cb.GetStats()
	if stats["totalRequests"] != int64(100) {
		t.Errorf("Expected 100 total requests, got %v", stats["totalRequests"])
	}
}

func TestCircuitBreaker_SuccessInClosedState(t *testing.T) {
	cb := NewCircuitBreaker(3, 2*time.Second, 2)

	// 在关闭状态下记录成功
	cb.RecordSuccess()
	cb.RecordSuccess()

	// 应该保持关闭状态
	if cb.GetState() != StateClosed {
		t.Errorf("Expected state to remain Closed, got %v", cb.GetState())
	}

	// 失败计数应该被重置
	if stats := cb.GetStats(); stats["failureCount"] != 0 {
		t.Errorf("Expected failureCount to be 0, got %v", stats["failureCount"])
	}
}

func TestCircuitBreaker_ZeroThreshold(t *testing.T) {
	cb := NewCircuitBreaker(0, 1*time.Second, 1)

	// 零阈值应该立即熔断
	cb.RecordFailure()

	if cb.GetState() != StateOpen {
		t.Errorf("Expected StateOpen with zero threshold, got %v", cb.GetState())
	}
}
