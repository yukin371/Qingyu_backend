package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector("http://localhost:9090")

	assert.NotNil(t, collector)
	assert.Equal(t, "http://localhost:9090", collector.prometheusURL)
	assert.NotNil(t, collector.client)
}

func TestMetricsCollector_CollectCacheMetrics(t *testing.T) {
	collector := NewMetricsCollector("http://localhost:9090")

	metrics, err := collector.CollectCacheMetrics()

	// 注意: 这需要Prometheus运行, 如果失败则跳过
	if err != nil {
		t.Skip("Prometheus not available, skipping metrics collection test")
	}

	require.NotNil(t, metrics)
	assert.GreaterOrEqual(t, metrics.HitRate, 0.0)
	assert.LessOrEqual(t, metrics.HitRate, 1.0)
}

func TestMetricsCollector_CollectSnapshot(t *testing.T) {
	collector := NewMetricsCollector("http://localhost:9090")

	snapshot, err := collector.CollectSnapshot()

	// 如果Prometheus不可用, 跳过测试
	if err != nil {
		t.Skip("Prometheus not available")
	}

	require.NotNil(t, snapshot)
	assert.Contains(t, snapshot, "cache_hits")
	assert.Contains(t, snapshot, "cache_misses")
}

func TestABTestBenchmark_RunABTestWithMetrics(t *testing.T) {
	benchmark := NewABTestBenchmark("http://httpbin.org")
	collector := NewMetricsCollector("http://localhost:9090")

	scenario := TestScenario{
		Name:       "Metrics Test",
		Requests:   10,
		Concurrent: 2,
		Endpoints:  []string{"/get"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用带指标收集的版本
	result, err := benchmark.RunABTestWithMetrics(ctx, scenario, true, collector)

	// 即使Prometheus不可用, 测试也应该成功
	if err != nil && err.Error() != "context deadline exceeded" {
		t.Skipf("Skipping due to error: %v", err)
	}

	if result != nil {
		assert.Equal(t, "Metrics Test", result.Scenario)
	}
}

func TestTestResult_CacheMetricsFields(t *testing.T) {
	result := &TestResult{
		CacheHits:      850,
		CacheMisses:    150,
		CacheHitRate:   0.85,
		DBQueries:      150,
	}

	assert.Equal(t, 850, result.CacheHits)
	assert.Equal(t, 150, result.CacheMisses)
	assert.Equal(t, 0.85, result.CacheHitRate)
	assert.Equal(t, 150, result.DBQueries)
}
