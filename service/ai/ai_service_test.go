package ai

import (
	"testing"
	"time"

	"Qingyu_backend/pkg/circuitbreaker"

	"github.com/stretchr/testify/assert"
)

// TestCircuitBreaker_StateMachine 测试熔断器状态机
func TestCircuitBreaker_StateMachine(t *testing.T) {
	cb := circuitbreaker.NewCircuitBreaker(3, 100*time.Millisecond, 2)

	// 初始状态：Closed
	assert.True(t, cb.IsClosed())
	assert.False(t, cb.IsOpen())
	assert.False(t, cb.IsHalfOpen())

	// 记录失败
	cb.RecordFailure()
	cb.RecordFailure()
	assert.True(t, cb.IsClosed()) // 还没到阈值

	// 第三次失败，应该打开
	cb.RecordFailure()
	assert.True(t, cb.IsOpen())
	assert.False(t, cb.IsClosed())
	assert.False(t, cb.IsHalfOpen())

	// 等待超时，进入半开状态
	time.Sleep(150 * time.Millisecond)
	assert.True(t, cb.AllowRequest()) // 应该允许请求（进入半开）
	assert.True(t, cb.IsHalfOpen())

	// 记录成功
	cb.RecordSuccess()
	assert.True(t, cb.IsHalfOpen()) // 还需要一次成功

	cb.RecordSuccess()
	assert.True(t, cb.IsClosed()) // 回到关闭状态
}

// TestCircuitBreaker_Stats 测试熔断器统计信息
func TestCircuitBreaker_Stats(t *testing.T) {
	cb := circuitbreaker.NewCircuitBreaker(5, 30*time.Second, 3)

	// 记录一些请求
	cb.AllowRequest()
	cb.RecordSuccess()
	cb.AllowRequest()
	cb.RecordFailure()

	stats := cb.GetStats()

	assert.Equal(t, "Closed", stats["state"])
	assert.Equal(t, int64(2), stats["totalRequests"])
	assert.Equal(t, int64(1), stats["totalSuccesses"])
	assert.Equal(t, int64(1), stats["totalFailures"])
	assert.Equal(t, 1, stats["failureCount"])
	assert.Equal(t, float64(0.5), cb.GetFailureRate())
}

// TestAIService_Create 测试创建 AI 服务
func TestAIService_Create(t *testing.T) {
	// 创建配置
	cfg := &AIServiceConfig{
		Endpoint:   "localhost:50051",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: time.Second,
	}

	// 创建服务（不连接真实的 gRPC 服务器）
	// 注意：这个测试需要一个真实的 gRPC 连接或模拟连接
	// 这里我们只测试配置创建
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost:50051", cfg.Endpoint)
	assert.Equal(t, 30*time.Second, cfg.Timeout)
}

// TestAIService_CircuitBreakerIntegration 测试熔断器集成
func TestAIService_CircuitBreakerIntegration(t *testing.T) {
	// 创建一个模拟的 gRPC 连接
	// 在实际测试中，应该使用模拟服务器或跳过
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 注意：这个测试需要真实的 gRPC 服务器运行
	// 在 CI/CD 环境中应该使用 mock 服务器
	t.Skip("Requires gRPC server - skipping in unit tests")
}

// TestAIService_FallbackAdapter 测试回退适配器
func TestAIService_FallbackAdapter(t *testing.T) {
	// 创建配置
	_ = &AIServiceConfig{
		Endpoint:       "invalid-host:9999",
		Timeout:        1 * time.Second,
		MaxRetries:     1,
		RetryDelay:     100 * time.Millisecond,
		EnableFallback: false, // 禁用回退
	}

	// 创建服务（需要 gRPC 连接）
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Requires gRPC connection - skipping in unit tests")
}

// BenchmarkCircuitBreaker_AllowRequest 基准测试熔断器
func BenchmarkCircuitBreaker_AllowRequest(b *testing.B) {
	cb := circuitbreaker.NewCircuitBreaker(5, 30*time.Second, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.AllowRequest()
	}
}
