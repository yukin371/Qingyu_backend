package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestABTestBenchmark_RunABTest(t *testing.T) {
	// 使用公共测试API httpbin.org
	benchmark := NewABTestBenchmark("http://httpbin.org")

	scenario := TestScenario{
		Name:       "Test Scenario",
		Requests:   10,
		Concurrent: 2,
		Endpoints:  []string{"/get", "/uuid"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := benchmark.RunABTest(ctx, scenario, true)

	// 验证没有错误
	require.NoError(t, err)

	// 验证基本字段
	assert.Equal(t, "Test Scenario", result.Scenario)
	assert.True(t, result.WithCache)
	assert.Equal(t, 10, result.TotalRequests)

	// 验证有成功的请求
	assert.Greater(t, result.SuccessCount, 0)

	// 验证平均延迟大于0
	assert.Greater(t, result.AvgLatency, time.Duration(0))

	// 验证总请求数等于成功+失败
	assert.Equal(t, result.TotalRequests, result.SuccessCount+result.ErrorCount)

	// 验证吞吐量大于0
	assert.Greater(t, result.Throughput, 0.0)
}

func TestTestResult_CalculateStatistics(t *testing.T) {
	// 创建一个已知的延迟列表
	latencies := []time.Duration{
		100 * time.Millisecond,
		150 * time.Millisecond,
		200 * time.Millisecond,
		250 * time.Millisecond,
		300 * time.Millisecond,
	}

	result := &TestResult{TotalRequests: 5}
	result.calculateStatistics(latencies)

	// 验证平均延迟
	expectedAvg := 200 * time.Millisecond
	assert.Equal(t, expectedAvg, result.AvgLatency)

	// 验证P95延迟（索引：int(5 * 0.95) = 4，即300ms）
	assert.Equal(t, 300*time.Millisecond, result.P95Latency)

	// 验证P99延迟（索引：int(5 * 0.99) = 4，即300ms）
	assert.Equal(t, 300*time.Millisecond, result.P99Latency)
}

func TestTestResult_CalculateStatistics_Empty(t *testing.T) {
	// 测试空切片情况
	result := &TestResult{TotalRequests: 0}
	result.calculateStatistics([]time.Duration{})

	// 空切片不应该导致panic，且延迟应该为0
	assert.Equal(t, time.Duration(0), result.AvgLatency)
	assert.Equal(t, time.Duration(0), result.P95Latency)
	assert.Equal(t, time.Duration(0), result.P99Latency)
}

func TestTestResult_CalculateStatistics_SingleValue(t *testing.T) {
	// 测试单个值的情况
	latencies := []time.Duration{100 * time.Millisecond}

	result := &TestResult{TotalRequests: 1}
	result.calculateStatistics(latencies)

	// 单个值时，所有延迟指标应该相同
	assert.Equal(t, 100*time.Millisecond, result.AvgLatency)
	assert.Equal(t, 100*time.Millisecond, result.P95Latency)
	assert.Equal(t, 100*time.Millisecond, result.P99Latency)
}
